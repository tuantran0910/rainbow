import textwrap
from abc import ABC
from abc import abstractmethod
from pathlib import Path
from typing import Any

import jinja2
import pendulum
from airflow.utils.log.logging_mixin import LoggingMixin

from core.builder.exceptions import DagBuilderException
from core.builder.exceptions import DagBuilderTemplateException
from core.models import DagParams


class BaseDagBuilder(ABC, LoggingMixin):
    """
    Abstract base class for all DAG builders. It takes parameters from configs and
    generate Airflow DAGs.

    Subclasses must implement the `build_dag` method to provide a concrete
    implementation specific to their workflow type.

    Args:
        airflow_type (str): Identifier for the DAG type (e.g., 'dbt', 'spark').
        templates_dir (Path): The directory where DAG templates are stored. Defaults to the folder `plugins/templates`.
        templated_dags_dir (Path): The directory where templated DAGs are stored. Defaults to the folder `dags/templated`.
        dag_configs (Dict[str, Any]): Dictionary of DAG configuration parameters.
    """

    def __init__(
        self,
        airflow_type: str,
        templates_dir: Path,
        templated_dags_dir: Path,
        dag_configs: dict[str, Any],
        dag_configs_file_path: Path,
    ):
        self.airflow_type = airflow_type
        self.dag_configs = dag_configs
        self.templated_dags_dir = templated_dags_dir
        self.dag_configs_file_path = dag_configs_file_path

        templates_dir = self._validate_template_dir(
            templates_dir=templates_dir, airflow_type=airflow_type
        )
        self.template = None
        self._initialize_template(templates_dir=templates_dir, airflow_type=airflow_type)

        self.dag_params: dict[str, Any] = {}
        self._initialize_template_airflow_params()

    @staticmethod
    def _validate_template_dir(templates_dir: Path, airflow_type: str) -> Path:
        """
        Validate the template directory.

        Args:
            templates_dir (Path): The path to the template directory.
            airflow_type (str): Identifier for the DAG type (e.g., 'dbt', 'spark').

        Returns:
            Path: The path to the specific template directory.

        Raises:
            DagBuilderException: If the 'airflow_type' field is missing.
            DagBuilderTemplateException: If the specific template directory does not exist.
        """
        if not airflow_type:
            raise DagBuilderException("The 'airflow_type' field is required.")

        specific_templates_dir = templates_dir / airflow_type
        if not specific_templates_dir.exists():
            raise DagBuilderTemplateException(
                f"The template directory '{specific_templates_dir}' does not exist."
            )

        return specific_templates_dir

    def _initialize_template(self, templates_dir: Path, airflow_type: str) -> None:
        """
        Initialize the Jinja2 template for the Airflow DAG.

        Args:
            templates_dir (Path): The path to the template directory.
            airflow_type (str): Identifier for the DAG type (e.g., 'dbt', 'spark').

        Raises:
            DagBuilderTemplateException: If the template is not found or has syntax errors.
        """
        try:
            env = jinja2.Environment(
                loader=jinja2.FileSystemLoader(searchpath=templates_dir),
                autoescape=False,
                undefined=jinja2.StrictUndefined,
                trim_blocks=True,
                lstrip_blocks=True,
            )
            template_name = f"{airflow_type}_dag_template.py.jinja2"
            self.template = env.get_template(template_name)
        except jinja2.TemplateNotFound as e:
            available = "\n".join(env.list_templates())
            raise DagBuilderTemplateException(
                f"Template {template_name} not found. Available templates:\n{available}"
            ) from e
        except jinja2.TemplateSyntaxError as e:
            raise DagBuilderTemplateException(
                f"Syntax error in template {e.name}:{e.lineno} - {e.message}"
            ) from e

    def _initialize_template_airflow_params(self) -> None:
        """
        Initialize the Airflow DAG parameters from the configuration.
        """
        airflow_configs = self.dag_configs["airflow"]
        self.dag_params = DagParams(**airflow_configs).model_dump()

    def _generate_doc_md(self) -> str:
        """
        Generate Markdown documentation for the DAG based on its configuration.

        Returns:
            str: Markdown-formatted documentation string
        """
        dag_id = self.dag_params["dag_id"]
        description = self.dag_params.get("description", "No description provided")
        schedule = self.dag_params.get("schedule_interval", "No schedule defined")
        owner = self.dag_params.get("default_args", {}).get("owner", "Unknown owner")
        max_active_runs = self.dag_params.get("max_active_runs", 1)

        doc = f"""
        ### {dag_id}

        **Description**:
        {description}

        #### DAG Configuration
        - **Schedule**: `{schedule}`
        - **Owner**: `{owner}`
        - **Max Active Runs**: `{max_active_runs}`
        - **Catchup**: `{self.dag_params.get("catchup", False)}`

        #### Generated Information
        - **Generated At**: {pendulum.now().to_iso8601_string()}
        - **Config Source**: `{self.dag_configs_file_path}`
        """

        return textwrap.dedent(doc).strip()

    def _get_output_path(self) -> Path:
        """
        Generate output path for the rendered DAG file.

        Format: `<templated_dags_dir>/<airflow_type>/rainbow__{dag_id}.py`

        Returns:
            Path: Output path for the DAG file.
        """
        filename = f"rainbow__{self.dag_params['dag_id']}.py"

        return self.templated_dags_dir / self.airflow_type / filename

    @abstractmethod
    def build_dag(self) -> None:
        """
        Builds the Airflow DAGs using templates and configurations.
        """
        raise NotImplementedError("Subclasses must implement build method")
