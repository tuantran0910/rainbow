import hashlib
import json
import os
import shutil
import tempfile
import time
from abc import ABC
from abc import abstractmethod
from pathlib import Path
from typing import Any

import yaml
from airflow.models import Variable
from airflow.utils.log.logging_mixin import LoggingMixin
from pydantic import ValidationError

from core.builder.base import BaseDagBuilder
from core.constants import CONFIG_FILE_EXT
from core.constants import CONFIGS_DIR
from core.constants import DAGS_DIR
from core.constants import TEMPLATED_DAGS_DIR
from core.constants import TEMPLATES_DIR
from core.factory.exceptions import DagFactoryConfigException
from core.factory.exceptions import DagFactoryException
from core.models import DefaultArgs


class BaseDagFactory(ABC, LoggingMixin):
    """
    Abstract base class for dynamically generating Airflow DAGs from YAML configurations.

    Subclasses must implement the `_get_builder` method to provide a concrete
    DAG builder implementation specific to their workflow type.

    The factory handles:
    - Configuration validation
    - Environment variable substitution
    - Template processing
    - Atomic deployment of generated DAGs

    Args:
        airflow_type (str): Identifier for the DAG type (e.g., 'dbt', 'spark').
        dags_dir (Path): The directory where DAGs are stored. Defaults to the folder `dags`.
        configs_dir (Path): Base directory containing type-specific configs. Defaults to the folder `configs`.
        config_file_ext (Set[str]): The set of valid config file extensions. Defaults to file with extensions `.yml` and `.yaml`.
        templates_dir (Path): The directory where DAG templates are stored. Defaults to the folder `plugins/templates`.
        templated_dags_dir (Path): The directory where templated DAGs are stored. Defaults to the folder `dags/templated`.
        default_args (Dict[str, Any]): The default arguments for the DAGs.
    """

    def __init__(
        self,
        airflow_type: str,
        dags_dir: Path = DAGS_DIR,
        configs_dir: Path = CONFIGS_DIR,
        config_file_ext: set[str] = CONFIG_FILE_EXT,
        templates_dir: Path = TEMPLATES_DIR,
        templated_dags_dir: Path = TEMPLATED_DAGS_DIR,
        default_args: dict[str, Any] = DefaultArgs().model_dump(),
        full_refresh: bool = False,
    ):
        self.airflow_type = airflow_type
        self.dags_dir = dags_dir
        self.configs_dir = configs_dir
        self.config_file_ext = config_file_ext
        self.templates_dir = templates_dir
        self.templated_dags_dir = templated_dags_dir
        self.default_args = default_args
        self.full_refresh = full_refresh

        self.hash_key = f"{self.airflow_type}_factory_configs_hash"

        # Validate directory paths
        self._validate_directory(path=self.templates_dir, info="Templates directory")
        self._validate_directory(path=self.configs_dir, info="Configs directory")

    @staticmethod
    def _validate_directory(path: Path, info: str) -> None:
        """
        Validate that a directory exists.

        Args:
            path (Path): The directory path to validate.
            info (str): The information message to display.

        Raises:
            DagFactoryException: If the directory does not exist or is not a directory.
        """
        if not path.exists():
            raise DagFactoryException(f"{info} '{path}' does not exist.")
        if not path.is_dir():
            raise DagFactoryException(f"{info} '{path}' is not a directory.")

    @staticmethod
    def _validate_config_dir(configs_dir: Path, airflow_type: str) -> Path:
        """
        Validate and return the specific config directory for the airflow type.

        Args:
            configs_dir (Path): The directory where YAML configs are stored.
            airflow_type (str): The type of Airflow DAG.

        Returns:
            Path: The specific config directory for the airflow type.

        Raises:
            DagFactoryException: If the 'airflow_type' field is missing.
            DagFactoryConfigException: If the specific config directory does not exist.
        """
        if not airflow_type:
            raise DagFactoryException("The 'airflow_type' field is required.")

        specific_config_dir = configs_dir / airflow_type
        if not specific_config_dir.exists():
            raise DagFactoryConfigException(
                f"Config directory '{specific_config_dir}' does not exist."
            )

        return specific_config_dir

    def _load_config(self, config_file_path: Path) -> dict[str, Any]:
        """
        Load and validate a YAML config file.

        Args:
            config_file_path (Path): The path to the YAML config file.

        Returns:
            Dict[str, Any]: The loaded and validated config.

        Raises:
            DagFactoryConfigException: If the config is invalid or missing required fields.
        """
        try:
            with open(config_file_path, encoding="utf-8") as fp:
                config_with_env = os.path.expandvars(fp.read())
                config: dict[str, Any] = yaml.safe_load(config_with_env)
                airflow_config = config.get("airflow")
                if airflow_config is None:
                    raise DagFactoryConfigException("Missing airflow configuration in config")

                if airflow_config.get("dag_id") is None:
                    raise DagFactoryConfigException("Airflow DAG ID is required in config")

            # Merge and validate default_args
            merged_args = self.default_args.copy()
            if airflow_config.get("default_args") is not None:
                merged_args.update(airflow_config["default_args"])

            airflow_config["dag_id"] = f"{self.airflow_type}__{airflow_config['dag_id']}"
            airflow_config["default_args"] = merged_args
            config["airflow"] = airflow_config

            return config
        except (ValidationError, yaml.YAMLError) as e:
            raise DagFactoryConfigException(f"Invalid config in '{config_file_path}'") from e
        except Exception as e:
            raise DagFactoryConfigException(
                f"Error processing config file '{config_file_path}'"
            ) from e

    def _load_configs(self) -> dict[str, dict[str, Any]]:
        """
        Load all YAML config files for the current airflow type.

        Returns:
            dict[str, dict[str, Any]]: The loaded and validated configs.
        """
        specific_config_dir = self._validate_config_dir(
            configs_dir=self.configs_dir, airflow_type=self.airflow_type
        )
        config_file_paths = [
            f
            for f in specific_config_dir.rglob("*")
            if f.is_file() and f.suffix in self.config_file_ext
        ]

        configs = dict()
        for config_file_path in config_file_paths:
            try:
                self.log.info(f"Loading config from {config_file_path}")
                config = self._load_config(config_file_path=config_file_path)
                configs[str(config_file_path)] = config
            except DagFactoryConfigException as e:
                self.log.error(f"Skipping invalid config: {str(e)}")
                continue

        if not configs:
            self.log.warning(f"No valid configs found in {specific_config_dir}")

        return configs

    def _deploy_dags(self, tmp_dir: str) -> None:
        """
        Atomically deploy generated DAGs from temporary directory to target.

        Args:
            tmp_dir (str): The temporary directory containing the generated DAGs.

        Raises:
            DagFactoryException: When DAG deployment fails.
        """
        target_dir = Path(self.templated_dags_dir)
        if not target_dir.exists():
            target_dir.mkdir(parents=True)

        try:
            for item in target_dir.iterdir():
                if item.is_file():
                    item.unlink()
                elif item.is_dir():
                    shutil.rmtree(item)

            # Copy back all DAGs to the target directory
            shutil.copytree(src=tmp_dir, dst=target_dir, dirs_exist_ok=True)
            self.log.info(
                f"Successfully deployed {len(list(target_dir.glob('*')))} DAGs to {target_dir}"
            )
        except Exception as e:
            self.log.error(f"Failed to deploy DAGs: {str(e)}")
            raise DagFactoryException("DAG deployment failed") from e

    def _calculate_configs_hash(self, configs: list[dict[str, Any]]) -> str:
        """
        Calculate a hash of the loaded configs to determine if they have changed.
        Won't re-generate DAGs if the configs haven't changed.

        Args:
            configs (List[Dict[str, Any]]): The loaded and validated configs.

        Returns:
            str: The hash of the configs.
        """
        configs_str = json.dumps(configs, sort_keys=True, default=str)
        return hashlib.sha256(configs_str.encode()).hexdigest()

    def _determine_if_dags_need_rebuilding(self, configs: list[dict[str, Any]]) -> tuple[str, bool]:
        """
        Determine if the DAGs need to be rebuilt based on the hash of the configs.

        Args:
            configs (List[Dict[str, Any]]): The loaded and validated configs.

        Returns:
            Tuple[str, bool]: The hash of the configs and whether the DAGs need rebuilding
        """
        configs_hash = self._calculate_configs_hash(configs=configs)
        last_configs_hash = Variable.get(key=self.hash_key, default_var="")
        return configs_hash, configs_hash != last_configs_hash

    def _generate_build_report(self, total_success: int, total_failed: int) -> str:
        """
        Generate a report of the DAG building process.

        Args:
            total_success (int): The total number of successful DAG builds.
            total_failed (int): The total number of failed DAG builds.

        Returns:
            str: The report of the DAG building process.
        """
        report = f"""
        DAG building report:
        - Total DAGs built: {total_success + total_failed}
        - Successful DAGs: {total_success}
        - Failed DAGs: {total_failed}
        """
        return report

    def build_dags(self) -> None:
        """
        Build Airflow DAGs from configs and deploy to target directory.
        """
        self.log.info(f"Running in {self.full_refresh} full refresh mode")
        # Load all configs for the current airflow type
        dags_configs = self._load_configs()
        configs_hash, need_rebuilding = self._determine_if_dags_need_rebuilding(
            configs=list(dags_configs.values())
        )
        if not need_rebuilding and not self.full_refresh:
            self.log.info("No changes in configs. Skipping DAG building...")
            return

        with tempfile.TemporaryDirectory() as tmp_dir:
            total_success = 0
            total_failed = 0
            for dag_configs_file_path, dag_configs in dags_configs.items():
                self.log.info(f"Building Airflow DAG: {dag_configs['airflow']['dag_id']}")
                try:
                    builder_cls = self._get_builder()
                    builder = builder_cls(
                        airflow_type=self.airflow_type,
                        templates_dir=self.templates_dir,
                        templated_dags_dir=Path(tmp_dir),
                        dag_configs=dag_configs,
                        dag_configs_file_path=dag_configs_file_path,
                    )
                    builder.build_dag()
                    total_success += 1
                except Exception as e:
                    self.log.exception(f"An unexpected error occurred while building DAG: {str(e)}")
                    total_failed += 1
                    continue

            built_report = self._generate_build_report(
                total_success=total_success, total_failed=total_failed
            )
            self.log.info(built_report)

            # Move generated DAGs to the templated DAGs directory
            self._deploy_dags(tmp_dir=tmp_dir)

        # Update the configs hash in the Airflow Variables
        if total_failed != 0:
            self.log.warning("Some DAGs failed to build. Skipping configs hash update...")
        else:
            Variable.set(
                key=self.hash_key,
                value=configs_hash,
                description=f"Configs hash generated on {time.ctime()}",
            )

    @abstractmethod
    def _get_builder(self) -> type[BaseDagBuilder]:
        """
        Abstract method to get the DAG builder implementation.
        """
        raise NotImplementedError
