import os
import re
import time
import unicodedata
from pathlib import Path
from typing import Any
from typing import Optional

import dagster as dg
import dlt
import requests
import yaml
from dlt.extract.source import DltSource
from dlt.sources.sql_database import sql_database

from constants import DAGSTER_ASSETS_CONFIG_DIR
from exceptions import AuthenticationError

logger = dg.get_dagster_logger(__name__)


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


class AuthTokenManager:
    """
    Manages authentication tokens for API requests.

    Args:
        api_url (str): The base URL of the API.
        admin_email (str): The email of the admin user.
        admin_password (str): The password of the admin user.
    """

    def __init__(
        self,
        api_url: str,
        admin_email: str,
        admin_password: str,
    ):
        self.api_url = api_url
        self.admin_email = admin_email
        self.admin_password = admin_password
        self.token = None
        self.token_expiry = 0.0

    def get_token(self) -> Optional[str]:
        """
        Gets a valid authentication token, logging in if necessary.

        Returns:
            Optional[str]: The authentication token or None if login fails
        """
        current_time = time.time()

        # Check if token exists and is not expired (with 5 min buffer)
        if self.token and current_time < (self.token_expiry - 300):
            return self.token

        # Need to login and get a new token
        logger.info("Getting new authentication token")
        login_url = f"{self.api_url}/auth/login"
        login_data = {"email": self.admin_email, "password": self.admin_password}
        headers = {"Content-Type": "application/json"}

        try:
            response = requests.post(login_url, json=login_data, headers=headers, timeout=30)
            response.raise_for_status()
            token_data = response.json()

            if token_data and "data" in token_data and "token" in token_data["data"]:
                self.token = token_data["data"]["token"]
            elif token_data and "token" in token_data:
                self.token = token_data["token"]
            else:
                logger.error(f"Unexpected response format: {token_data}")
                return None

            self.token_expiry = current_time + (24 * 60 * 60)
            logger.info("Successfully obtained authentication token")
            return self.token
        except requests.exceptions.RequestException as e:
            logger.error(f"Failed to get authentication token: {str(e)}")
            if hasattr(e, "response") and e.response:
                logger.error(f"Response status: {e.response.status_code}, Body: {e.response.text}")
            raise e


def make_http_request(
    url: str,
    method: str = "GET",
    params: Optional[dict[str, Any]] = None,
    data: Optional[dict[str, Any]] = None,
    headers: Optional[dict[str, Any]] = None,
    timeout: int = 30,
    max_retries: int = 3,
    retry_delay: float = 1.0,
    retry_backoff: float = 2.0,
    use_auth: bool = False,
    auth_token_manager: Optional[AuthTokenManager] = None,
) -> Optional[dict[str, Any]]:
    """
    Makes an HTTP request with retry logic and optional authentication.

    Args:
        url (str): The URL to make the request to.
        method (str, optional): The HTTP method to use. Default: "GET".
        params (dict[str, Any], optional): Query parameters.
        data (dict[str, Any], optional): JSON data for request body.
        headers (dict[str, Any], optional): Request headers.
        timeout (int): Request timeout in seconds (default: 30).
        max_retries (int): Maximum retry attempts (default: 3).
        retry_delay (float): Initial delay between retries (default: 1.0).
        retry_backoff (float): Multiplier for increasing retry delay (default: 2.0).
        use_auth (bool): Whether to use authentication (default: False).
        auth_token_manager (AuthTokenManager, optional): Token manager.

    Returns:
        Optional[dict[str, Any]]: JSON response or None if request failed.
    """
    method = method.upper()
    request_headers = headers.copy() if headers else {}

    # Set content type for request methods with body
    if method in ["POST", "PUT", "PATCH"] and "Content-Type" not in request_headers:
        request_headers["Content-Type"] = "application/json"

    attempts = 0
    current_delay = retry_delay

    while attempts < max_retries:
        # Add auth token if required (refresh on each attempt to handle expiration)
        if use_auth and auth_token_manager:
            token = auth_token_manager.get_token()
            if token:
                request_headers["Authorization"] = f"Bearer {token}"
            else:
                logger.error("Authentication required but couldn't get token")
                raise AuthenticationError("Authentication failed")

        try:
            # Create session and prepare request
            session = requests.Session()
            request = requests.Request(
                method=method,
                url=url,
                params=params,
                json=data if method != "GET" else None,
                headers=request_headers,
            )
            prepped = request.prepare()

            # Execute request
            response = session.send(prepped, timeout=timeout)
            response_json: dict[str, Any] = response.json()

            is_warning_response = False
            if 400 <= response.status_code < 600:
                response_message: Optional[str] = response_json.get("message")
                if response_json.get("error"):
                    logger.error(
                        f"Error response: {response.json()['error']} - URL: {url} - Status Code: {response.status_code}"
                    )
                elif response_message is not None and "not found" in response_message.lower():
                    is_warning_response = True
                    logger.warning(
                        f"Resource not found: {response_message} - URL: {url} - Status Code: {response.status_code}"
                    )

            if not is_warning_response:
                response.raise_for_status()
            return response.json()

        except (requests.exceptions.Timeout, requests.exceptions.ConnectionError) as e:
            # Network-related errors
            attempts += 1
            if attempts >= max_retries:
                logger.error(f"Request failed after {max_retries} attempts: {url} - {e}")
                raise dg.DagsterError(f"Request failed: {e}")

        except requests.exceptions.HTTPError as e:
            status_code = e.response.status_code
            # Handle auth token refresh for 401 errors
            if status_code == 401 and use_auth and auth_token_manager:
                # Force token refresh
                auth_token_manager.token = None
                attempts += 1
            # Only retry server errors
            elif 500 <= status_code < 600:
                attempts += 1
            else:
                logger.error(f"HTTP Error {status_code} for URL: {url}")
                raise dg.DagsterError(f"HTTP Error {status_code}: {e}")

        except Exception as e:
            logger.error(f"Unexpected error during request to {url}: {e}")
            raise dg.DagsterError(f"Unexpected error occurred: {e}")

        # Apply backoff delay before retry
        if attempts < max_retries:
            logger.warning(
                f"Retrying request (attempt {attempts}/{max_retries}) in {current_delay}s"
            )
            time.sleep(current_delay)
            current_delay *= retry_backoff

    return None


def sanitize_text(text: Optional[str]) -> Optional[str]:
    """
    Sanitize text by handling special characters, normalizing Unicode, and cleaning whitespace.
    This is a general-purpose function that can handle various text issues.

    Args:
        text (Optional[str]): The text to sanitize

    Returns:
        Optional[str]: The sanitized text
    """
    if not text:
        return text

    try:
        text = unicodedata.normalize("NFKC", text)
        replacements = {
            "\xa0": " ",  # Non-breaking space
            "\u200b": "",  # Zero-width space
            "\u200c": "",  # Zero-width non-joiner
            "\u200d": "",  # Zero-width joiner
            "\u2028": " ",  # Line separator
            "\u2029": " ",  # Paragraph separator
            "\u202f": " ",  # Narrow no-break space
            "\u205f": " ",  # Medium mathematical space
            "\u3000": " ",  # Ideographic space
            "\ufeff": "",  # Byte order mark
            "\u200e": "",  # Left-to-right mark
            "\u200f": "",  # Right-to-left mark
            "\u202a": "",  # Left-to-right embedding
            "\u202b": "",  # Right-to-left embedding
            "\u202c": "",  # Pop directional formatting
            "\u202d": "",  # Left-to-right override
            "\u202e": "",  # Right-to-left override
        }

        for char, replacement in replacements.items():
            text = text.replace(char, replacement)

        text = "".join(ch for ch in text if unicodedata.category(ch)[0] != "C")
        text = re.sub(r"\s+", " ", text)
        return text.strip()
    except Exception as e:
        logger.warning(f"Error sanitizing text: {e}. Returning original text.")
        return text
