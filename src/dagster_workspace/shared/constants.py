import os
from pathlib import Path
from typing import Any

# Tiki Crawler
TIKI_REQUEST_DELAY = 0.5
TIKI_MAX_PAGES_PER_CATEGORY = int(os.getenv("TIKI_MAX_PAGES_PER_CATEGORY", "0"))  # 0 means no limit
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

INVENTORY_MIN_STOCK = int(os.getenv("INVENTORY_MIN_STOCK", "10"))
INVENTORY_MAX_STOCK = int(os.getenv("INVENTORY_MAX_STOCK", "50"))

# Enhanced Mocking Configuration
MAX_BOOKS_PER_REQUEST = int(os.getenv("MAX_BOOKS_PER_REQUEST", "10"))
MAX_ITEMS_PER_ORDER = int(os.getenv("MAX_ITEMS_PER_ORDER", "5"))
MAX_QUANTITY_PER_BOOK = int(os.getenv("MAX_QUANTITY_PER_BOOK", "3"))

MAX_USERS_PER_REQUEST = int(os.getenv("MAX_USERS_PER_REQUEST", "5"))
DEFAULT_USER_PASSWORD = os.getenv("DEFAULT_USER_PASSWORD", "default")

MAX_PROMOTIONS_PER_REQUEST = int(os.getenv("MAX_PROMOTIONS_PER_REQUEST", "2"))
MIN_PERCENTAGE_DISCOUNT = float(os.getenv("MIN_PERCENTAGE_DISCOUNT", "5.0"))
MAX_PERCENTAGE_DISCOUNT = float(os.getenv("MAX_PERCENTAGE_DISCOUNT", "50.0"))
MIN_FIXED_DISCOUNT_VND = float(os.getenv("MIN_FIXED_DISCOUNT_VND", "15000.0"))
MAX_FIXED_DISCOUNT_VND = float(os.getenv("MAX_FIXED_DISCOUNT_VND", "30000.0"))
MIN_PROMOTION_USES = int(os.getenv("MIN_PROMOTION_USES", "10"))
MAX_PROMOTION_USES = int(os.getenv("MAX_PROMOTION_USES", "100"))
FLASH_SALE_MIN_USES = int(os.getenv("FLASH_SALE_MIN_USES", "10"))
FLASH_SALE_MAX_USES = int(os.getenv("FLASH_SALE_MAX_USES", "30"))
FLASH_SALE_HOURS = int(os.getenv("FLASH_SALE_HOURS", "3"))
WEEKEND_SALE_DAYS = int(os.getenv("WEEKEND_SALE_DAYS", "3"))
SEASONAL_SALE_MIN_DAYS = int(os.getenv("SEASONAL_SALE_MIN_DAYS", "7"))
SEASONAL_SALE_MAX_DAYS = int(os.getenv("SEASONAL_SALE_MAX_DAYS", "14"))
REGULAR_SALE_MIN_DAYS = int(os.getenv("REGULAR_SALE_MIN_DAYS", "1"))
REGULAR_SALE_MAX_DAYS = int(os.getenv("REGULAR_SALE_MAX_DAYS", "7"))
BOOK_CLUB_SALE_MIN_DAYS = int(os.getenv("BOOK_CLUB_SALE_MIN_DAYS", "14"))
BOOK_CLUB_SALE_MAX_DAYS = int(os.getenv("BOOK_CLUB_SALE_MAX_DAYS", "30"))

# Enhanced HTTP & Authentication Configuration
HTTP_REQUEST_TIMEOUT = int(os.getenv("HTTP_REQUEST_TIMEOUT", "30"))
HTTP_MAX_RETRIES = int(os.getenv("HTTP_MAX_RETRIES", "3"))
HTTP_RETRY_DELAY = float(os.getenv("HTTP_RETRY_DELAY", "1.0"))
HTTP_RETRY_BACKOFF = float(os.getenv("HTTP_RETRY_BACKOFF", "2.0"))

# Rate Limiting & Throttling
THROTTLE_DELAY_BETWEEN_REQUESTS = float(os.getenv("THROTTLE_DELAY_BETWEEN_REQUESTS", "0.1"))

# Data Quality & Validation
ENABLE_DATA_VALIDATION = os.getenv("ENABLE_DATA_VALIDATION", "true").lower() == "true"
MAX_DUPLICATE_EMAIL_ATTEMPTS = int(os.getenv("MAX_DUPLICATE_EMAIL_ATTEMPTS", "10"))

# Dagster Resources
ENVIRONMENT = os.getenv("ENVIRONMENT", "local")
DAGSTER_ASSETS_CONFIG_DIR = Path(
    os.getenv("DAGSTER_ASSETS_CONFIG_DIR", str(Path(__file__).absolute().parent / "configs"))
)
DAGSTER_ASSETS_OWNER = "tntuan0910@gmail.com"
DAGSTER_METADATA = {
    "owner": DAGSTER_ASSETS_OWNER,
    "team": "data_engineering",
}
DAGSTER_TAGS = {"team": "data_engineering"}

DAGSTER_CRAWLING_ASSET_GROUP = "crawling"
DAGSTER_MOCKING_ASSET_GROUP = "mocking"
DAGSTER_DBT_ASSET_GROUP = "dbt"

DB_HOST = os.getenv("DB_HOST", "localhost")
DB_PORT = int(os.getenv("DB_PORT", "5432"))

RAINBOW_DB_NAME = os.getenv("API_DB_NAME", "rainbow")
RAINBOW_DB_USER = os.getenv("API_DB_USER", "rainbow")
RAINBOW_DB_PASSWORD = os.getenv("API_DB_PASSWORD")

# Dbt
DBT_PROJECT_DIR = Path(__file__).absolute().parent.parent / "dbt" / ENVIRONMENT
DBT_PROFILES_DIR = DBT_PROJECT_DIR
DBT_TARGET_PROFILE = os.getenv("DBT_TARGET_PROFILE", "production")
DBT_TARGET_PATH = DBT_PROJECT_DIR / "target"

# Dbt Docs Generation
DBT_GCS_PROJECT = os.getenv("DBT_GCS_PROJECT", "rainbow-data-production")
DBT_DOCS_BUCKET_NAME = os.getenv("DBT_DOCS_BUCKET_NAME", "rainbow-data-production-dbt")
DBT_DOCS_BUCKET_FOLDER = "docs"
DBT_DOCS_BASE_URL = os.getenv("DBT_DOCS_BASE_URL", "https://dbt-docs.tuantrann.work")
