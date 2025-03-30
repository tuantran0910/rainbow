import os
from pathlib import Path

DAGSTER_ASSETS_CONFIG_DIR = Path(os.getenv("DAGSTER_ASSETS_CONFIG_DIR", "/opt/dagster/app/configs"))
DAGSTER_DBT_TARGET_PROFILE = os.getenv("DAGSTER_DBT_TARGET_PROFILE", "production")
DAGSTER_ASSETS_OWNER = ["tntuan0910@gmail.com"]
