import os
from pathlib import Path
from typing import Any

# Tiki Crawler
TIKI_REQUEST_DELAY = 0.5
TIKI_HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0",
    "Accept": "application/json, text/plain, */*",
    "Accept-Language": "vi-VN,vi;q=0.8,en-US;q=0.5,en;q=0.3",
    "Accept-Encoding": "gzip, deflate, br",
    "Connection": "keep-alive",
}
TIKI_PRODUCT_LISTINGS_PAGE_PARAMS: dict[str, Any] = {
    "limit": 10,
}
TIKI_CATEGORIES = {
    8322: "nha-sach-tiki",
}
TIKI_BASE_PRODUCT_LISTINGS = "https://tiki.vn/api/personalish/v1/blocks/listings"
TIKI_BASE_SPECIFIC_PRODUCT = "https://tiki.vn/api/v2/products/"

API_BASE_URL = os.getenv("API_BASE_URL", "http://127.0.0.1:5000")

ADMIN_EMAIL = os.getenv("ADMIN_EMAIL", "admin@example.com")
ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD", "Admin123!")

# Tiki Mocker
MAX_BOOKS_PER_REQUEST = int(os.getenv("MAX_BOOKS_PER_REQUEST", "10"))
MAX_ITEMS_PER_ORDER = int(os.getenv("MAX_ITEMS_PER_ORDER", "5"))

MAX_USERS_PER_REQUEST = int(os.getenv("MAX_USERS_PER_REQUEST", "5"))
DEFAULT_USER_PASSWORD = os.getenv("DEFAULT_USER_PASSWORD", "default")

# Dagster Resources
DAGSTER_ASSETS_CONFIG_DIR = Path(os.getenv("DAGSTER_ASSETS_CONFIG_DIR", "/opt/dagster/app/configs"))
DAGSTER_DBT_TARGET_PROFILE = os.getenv("DAGSTER_DBT_TARGET_PROFILE", "production")
DAGSTER_ASSETS_OWNER = "tntuan0910@gmail.com"
DAGSTER_METADATA = {
    "owner": DAGSTER_ASSETS_OWNER,
    "team": "data_engineering",
}
DAGSTER_TAGS = {"team": "data_engineering"}

DAGSTER_CRAWLING_ASSET_GROUP = "crawling"
DAGSTER_MOCKING_ASSET_GROUP = "mocking"
