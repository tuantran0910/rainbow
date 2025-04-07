import logging
import os
from pathlib import Path
from typing import Any
from typing import Optional

import dlt
import requests
import yaml
from dlt.extract.source import DltSource
from dlt.sources.sql_database import sql_database

from constants import DAGSTER_ASSETS_CONFIG_DIR


logger = logging.getLogger(__name__)


def load_assets_configs(
    dagster_product: str, config_dir: str | Path = DAGSTER_ASSETS_CONFIG_DIR
) -> dict[str, Any] | None:
    """
    Load Dagster assets' configurations from YAML files in a directory.

    Args:
        dagster_product (str): Name of the Dagster product (e.g., "dlt").
        config_dir (str | Path): Path to the base directory containing product-specific YAML files.

    Returns:
        dict[str, Any] | None: A dictionary containing the loaded configurations, or None if no valid configurations are found.

    Raises:
        FileNotFoundError: If the specified directory does not exist.
    """
    config_dir = Path(config_dir) / dagster_product
    if not config_dir.exists():
        raise FileNotFoundError(f"Assets configurations directory '{config_dir}' does not exist.")

    loaded_configs = {}
    # Look for files with both .yaml and .yml extensions
    for yaml_file in config_dir.glob("*.y*ml"):
        try:
            with open(yaml_file, encoding="utf-8") as file:
                file_configs = yaml.safe_load(file)
                if file_configs:
                    loaded_configs.update(file_configs)
        except yaml.YAMLError as e:
            logger.error(f"Error parsing YAML file {yaml_file}: {e}")
            raise

    return loaded_configs if loaded_configs else None


def make_dlt_resources(dlt_resources_config: dict[str, Any]) -> tuple[dict[str, DltSource], str]:
    """
    Initializes dlt resources (including general configs, sources, and destinations).

    Args:
        dlt_resources_config (dict[str, Any]): A dictionary containing dlt resources configuration.

    Returns:
        tuple[dict[str, DltSource], str]: A tuple containing the table source's name and the DltSource object.

    Raises:
        ValueError: If no dlt sources or destinations are provided.
        KeyError: If the 'type' key is missing from a source or destination dictionary.
    """
    config: dict[str, Any] = dlt_resources_config.get("config", {})
    if config:
        set_dlt_object(dlt.config, config)

    source_secrets: dict[str, Any] = {}
    destination_secrets: dict[str, Any] = {}

    # Configure dlt source
    source: dict[str, Any] = dlt_resources_config.get("source", {})
    if not source:
        raise ValueError("No dlt source provided")

    source_type = source.get("type")
    if not source_type:
        raise KeyError(
            "The type of source is missing. Currently supported types are: 'sql_database'"
        )

    # Configure credentials and configs for the source
    source_credentials = {
        "credentials": source.get("credentials"),
    }
    source_secrets.setdefault("sources", {}).setdefault(source_type, {}).update(source_credentials)
    source_secrets["sources"][source_type].update(source.get("configs", {}))
    set_dlt_object(dlt.secrets, source_secrets)

    # Configure the source tables
    schema: dict[str, Any] = source.get("schema", {})
    if not schema:
        raise ValueError("Source schema is missing")
    schema_name = schema.get("name") if schema else None
    tables = schema.get("tables")
    if not tables:
        raise ValueError("No tables defined in schema")

    dlt_sources: dict[str, DltSource] = {}
    for table in tables:
        table_name = table["name"]
        columns = table.get("columns")
        incremental_field = table.get("incremental_field")
        initial_value = table.get("initial_value")
        chunk_size = int(table.get("chunk_size", "50000"))
        write_disposition = table.get("write_disposition", "append")
        dlt_source = sql_database(schema=schema_name, chunk_size=chunk_size).with_resources(
            table_name
        )
        if incremental_field:
            table_dlt_source = getattr(dlt_source, table_name)
            table_dlt_source.apply_hints(
                columns=columns,
                incremental=dlt.sources.incremental(
                    cursor_path=incremental_field, initial_value=initial_value
                ),
                write_disposition=write_disposition,
            )
        dlt_sources[table_name] = dlt_source

    # Configure the dlt destination
    destination: dict[str, Any] = dlt_resources_config.get("destination", {})
    if not destination:
        raise ValueError("No dlt destinations provided")

    destination_type = destination.get("type")
    if not destination_type:
        raise KeyError("Destination 'type' is missing")

    # Configure credentials and configs for the destination
    destination_credentials = {
        "credentials": destination.get("credentials"),
    }
    destination_secrets.setdefault("destination", {}).setdefault(destination_type, {}).update(
        destination_credentials
    )
    destination_secrets["destination"][destination_type].update(destination.get("configs", {}))
    set_dlt_object(dlt.secrets, destination_secrets)

    return dlt_sources, destination_type


def set_dlt_object(dlt_object: dict[str, Any], config: dict[str, Any], *, prefix: str = "") -> None:
    """
    Recursively sets attributes on a DLT object from a configuration dictionary.

    Supports resolving environment variables for values prefixed with 'env:'.

    Args:
        dlt_object (dict[str, Any]): The DLT object on which to set attributes.
        config (dict[str, Any]): Configuration dictionary containing values.
        prefix (str, optional): Prefix for nested attributes (default: "").
    """
    try:
        for key, value in config.items():
            full_key = f"{prefix}.{key}" if prefix else key

            if isinstance(value, dict):
                set_dlt_object(dlt_object, value, prefix=full_key)
                continue

            if isinstance(value, str) and value.startswith("env:"):
                env_var_name = value[4:]
                env_var = os.getenv(env_var_name)
                if env_var is None:
                    logger.warning(f"Environment variable {env_var_name} not found")
                value = env_var if env_var is not None else value

            dlt_object[full_key] = value
    except Exception as e:
        logger.error(f"Error setting attribute {full_key}: {e}")
        raise


def make_http_request(
    url: str,
    method: str = "GET",
    params: Optional[dict[str, Any]] = None,
    headers: Optional[dict[str, Any]] = None,
    timeout: int = 30,
) -> Optional[dict[str, Any]]:
    """
    Makes an HTTP request to the provided URL with the given method, params, and headers.

    Args:
        url (str): The URL to make the request to.
        method (str, optional): The HTTP method to use (default: "GET").
        params (dict[str, Any], optional): The parameters to send with the request.
        headers (dict[str, Any], optional): The headers to send with the request.
        timeout (int): The timeout for the request (default: 30 seconds).

    Returns:
        Optional[dict[str, Any]]: The response from the request, or None if an error occurs.
    """
    try:
        if method == "GET":
            response = requests.get(url, params=params, headers=headers, timeout=timeout)
        elif method == "POST":
            response = requests.post(url, json=params, headers=headers, timeout=timeout)
        else:
            raise ValueError(f"Unsupported method: {method}")

        response.raise_for_status()
        return response.json()
    except requests.exceptions.Timeout:
        logger.error(f"Request timed out: {url}")
    except requests.exceptions.HTTPError as e:
        logger.error(f"HTTP Error: {e.response.status_code} for URL: {url}")
    except requests.exceptions.RequestException as e:
        logger.error(f"Error making request to {url}: {e}")
    except Exception as e:
        logger.error(f"An unexpected error occurred during request to {url}: {e}")

    return None
