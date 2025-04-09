import os
from pathlib import Path

DAGSTER_ASSETS_CONFIG_DIR = Path(os.getenv("DAGSTER_ASSETS_CONFIG_DIR", "/opt/dagster/app/configs"))
DAGSTER_DBT_TARGET_PROFILE = os.getenv("DAGSTER_DBT_TARGET_PROFILE", "production")
DAGSTER_ASSETS_OWNER = ["tntuan0910@gmail.com"]

TIKI_REQUEST_DELAY = 0.5
TIKI_HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0",
    "Accept": "application/json, text/plain, */*",
    "Accept-Language": "vi-VN,vi;q=0.8,en-US;q=0.5,en;q=0.3",
    "Accept-Encoding": "gzip, deflate, br",
    "Connection": "keep-alive",
}
TIKI_PRODUCT_LISTINGS_PAGE_PARAMS = {
    "limit": "10",
}
TIKI_CATEGORIES = {
    8322: "nha-sach-tiki",
}
TIKI_BASE_PRODUCT_LISTINGS = "https://tiki.vn/api/personalish/v1/blocks/listings"
TIKI_BASE_SPECIFIC_PRODUCT = "https://tiki.vn/api/v2/products/"

API_URL = os.getenv("API_URL", "http://api:5000/api")

ADMIN_EMAIL = os.getenv("ADMIN_EMAIL", "admin@example.com")
ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD", "Admin123!")
