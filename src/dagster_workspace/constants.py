import os
from pathlib import Path

DAGSTER_ASSETS_CONFIG_DIR = Path(os.getenv("DAGSTER_ASSETS_CONFIG_DIR", "/opt/dagster/app/configs"))
