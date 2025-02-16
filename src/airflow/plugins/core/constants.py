import os
from pathlib import Path

# Path configurations
CONFIGS_DIR = Path(os.getenv("CONFIGS_DIR", "/opt/airflow/configs"))
DAGS_DIR = Path(os.getenv("DAGS_DIR", "/opt/airflow/dags"))
PLUGINS_DIR = Path(os.getenv("PLUGINS_DIR", "/opt/airflow/plugins"))
TEMPLATES_DIR = PLUGINS_DIR / "templates"
TEMPLATED_DAGS_DIR = DAGS_DIR / "templated"
DBT_PROJECT_DIR = Path(os.getenv("DBT_PROJECT_DIR", "/opt/airflow/dbt"))

# Other configurations
CONFIG_FILE_EXT = {".yml", ".yaml"}

# Airflow DAG configurations
DBT_TASK_GROUP_ID = "rainbow_dbt_task_group"
