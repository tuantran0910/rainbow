import json
import logging
import time
from typing import Any
from typing import Optional

from pydantic import BaseModel

# from constants import (
#     TIKI_CATEGORIES,
#     TIKI_BASE_PRODUCT_LISTINGS,
#     TIKI_BASE_SPECIFIC_PRODUCT,
#     TIKI_REQUEST_DELAY,
#     TIKI_HEADERS,
#     TIKI_PRODUCT_LISTINGS_PAGE_PARAMS,
#     API_URL,
# )
# from assets.helpers import make_http_request

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

import os

API_URL = os.getenv("API_URL", "http://127.0.0.1:5000")

import requests


class AuthTokenManager:
    """
    Manages authentication tokens for API requests.
    """

    def __init__(
        self,
        api_url: str,
        admin_email: str = "admin@example.com",
        admin_password: str = "Admin123!",
    ):
        self.api_url = api_url
        self.admin_email = admin_email
        self.admin_password = admin_password
        self.token = None
        self.token_expiry = 0.0  # Unix timestamp when token expires

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
        logging.info("Getting new authentication token")
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
                logging.error(f"Unexpected response format: {token_data}")
                return None

            # Assume token is valid for 24 hours
            self.token_expiry = current_time + (24 * 60 * 60)
            logging.info("Successfully obtained authentication token")
            return self.token
        except requests.exceptions.RequestException as e:
            logging.error(f"Failed to get authentication token: {str(e)}")
            if hasattr(e, "response") and e.response:
                logging.error(f"Response status: {e.response.status_code}, Body: {e.response.text}")
            return None


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
    Makes an HTTP request to the provided URL with the given method, params, and headers.

    Args:
        url (str): The URL to make the request to.
        method (str, optional): The HTTP method to use (GET, POST, PUT, PATCH, DELETE). Default: "GET".
        params (dict[str, Any], optional): The query parameters to send with the request.
        data (dict[str, Any], optional): The JSON data to send in the request body for POST, PUT, PATCH.
        headers (dict[str, Any], optional): The headers to send with the request.
        timeout (int): The timeout for the request (default: 30 seconds).
        max_retries (int): Maximum number of retry attempts (default: 3).
        retry_delay (float): Initial delay between retries in seconds (default: 1.0).
        retry_backoff (float): Multiplier for increasing retry delay with each attempt (default: 2.0).
        use_auth (bool): Whether to use authentication for this request (default: False).
        auth_token_manager (AuthTokenManager, optional): Token manager for authentication.

    Returns:
        Optional[dict[str, Any]]: The response from the request, or None if an error occurs.
    """
    method = method.upper()
    attempts = 0
    current_delay = retry_delay

    # Create a copy of the headers to avoid modifying the original
    request_headers = headers.copy() if headers else {}

    # Add authentication if required - check base URL more flexibly
    if use_auth and auth_token_manager:
        token = auth_token_manager.get_token()
        if token:
            request_headers["Authorization"] = f"Bearer {token}"
        else:
            logger.error("Authentication required but couldn't get token")
            return None

    # Ensure we have Content-Type for POST/PUT/PATCH
    if method in ["POST", "PUT", "PATCH"] and "Content-Type" not in request_headers:
        request_headers["Content-Type"] = "application/json"

    while attempts < max_retries:
        try:
            if method == "GET":
                response = requests.get(
                    url, params=params, headers=request_headers, timeout=timeout
                )
            elif method == "POST":
                response = requests.post(
                    url, params=params, json=data, headers=request_headers, timeout=timeout
                )
            elif method == "PUT":
                response = requests.put(
                    url, params=params, json=data, headers=request_headers, timeout=timeout
                )
            elif method == "PATCH":
                response = requests.patch(
                    url, params=params, json=data, headers=request_headers, timeout=timeout
                )
            elif method == "DELETE":
                response = requests.delete(
                    url, params=params, json=data, headers=request_headers, timeout=timeout
                )
            else:
                raise ValueError(f"Unsupported method: {method}")

            response.raise_for_status()
            return response.json()

        except requests.exceptions.Timeout:
            attempts += 1
            if attempts >= max_retries:
                logger.error(f"Request timed out after {max_retries} attempts: {url}")
                break
            logger.warning(
                f"Request timeout (attempt {attempts}/{max_retries}), retrying in {current_delay}s: {url}"
            )

        except requests.exceptions.HTTPError as e:
            status_code = e.response.status_code
            # Handle authentication error - token might be expired
            if status_code == 401 and use_auth and auth_token_manager and attempts < max_retries:
                logger.warning("Authentication token might be expired, refreshing and retrying...")
                # Reset token to force new login
                auth_token_manager.token = None
                token = auth_token_manager.get_token()
                if token:
                    request_headers["Authorization"] = f"Bearer {token}"
                    attempts += 1
                else:
                    logger.error("Failed to refresh authentication token")
                    break
            # Only retry on certain status codes (server errors)
            elif 500 <= status_code < 600 and attempts < max_retries:
                attempts += 1
                logger.warning(
                    f"HTTP Error {status_code} (attempt {attempts}/{max_retries}), retrying in {current_delay}s: {url}"
                )
            else:
                logger.error(f"HTTP Error: {status_code} for URL: {url}")
                break

        except requests.exceptions.RequestException as e:
            attempts += 1
            if attempts >= max_retries:
                logger.error(f"Request exception after {max_retries} attempts: {url} - {e}")
                break
            logger.warning(
                f"Request exception (attempt {attempts}/{max_retries}), retrying in {current_delay}s: {url} - {e}"
            )

        except Exception as e:
            logger.error(f"An unexpected error occurred during request to {url}: {e}")
            break

        # Wait before retrying with exponential backoff
        if attempts < max_retries:
            time.sleep(current_delay)
            current_delay *= retry_backoff

    return None


logger = logging.getLogger(__name__)
logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)


class TikiSeller(BaseModel):
    secondary_id: str
    name: str
    link: str
    logo: str


class TikiBook(BaseModel):
    secondary_id: str
    name: str
    price: float
    original_price: float
    rating_average: float
    review_count: int
    page_count: int
    author_ids: list[int] = []


class TikiAuthor(BaseModel):
    secondary_id: str
    name: str
    slug: str


class TikiCategory(BaseModel):
    secondary_id: str
    name: str


class TikiCrawler:
    def __init__(self, admin_email: str = "admin@example.com", admin_password: str = "Admin123!"):
        self.categories = TIKI_CATEGORIES
        self.request_params_base = TIKI_PRODUCT_LISTINGS_PAGE_PARAMS.copy()
        self.request_headers = TIKI_HEADERS
        self.request_listings_url = TIKI_BASE_PRODUCT_LISTINGS
        self.request_product_url_base = TIKI_BASE_SPECIFIC_PRODUCT
        self.request_delay = TIKI_REQUEST_DELAY

        # API Configuration
        self.api_base_url = API_URL
        self.api_url = f"{API_URL}/api"  # For resource endpoints
        self.auth_token_manager = AuthTokenManager(self.api_base_url, admin_email, admin_password)

        logger.info(f"Initialized TikiCrawler with API URL: {self.api_url}")
        logger.info(f"Auth base URL: {self.api_base_url}")

    def _get_total_pages(self, category_id: int, url_key: str) -> int:
        """
        Gets the total number of pages for a given category.

        Args:
            category_id (int): The ID of the category to get the total pages for.
            url_key (str): The URL key of the category to get the total pages for.

        Returns:
            int: The total number of pages for the given category.

        Raises:
            JSONDecodeError: If the JSON response is invalid.
            Exception: If an error occurs while getting the total pages.
        """
        logger.info(f"Getting total pages for category {category_id} ({url_key})...")

        params = self.request_params_base.copy()
        params["category"] = category_id
        params["urlKey"] = url_key
        # Always check the first page for total pages
        params["page"] = 1

        data = make_http_request(
            url=self.request_listings_url, headers=self.request_headers, params=params
        )
        if not data:
            logger.error("Failed to get initial page response.")
            return 0

        try:
            last_page = data.get("paging", {}).get("last_page", 0)
            if last_page > 0:
                logger.info(f"Total pages found: {last_page}")
                return last_page
            else:
                logger.warning("Could not find 'last_page' in response paging data.")
                return 0
        except json.JSONDecodeError:
            logger.exception("Failed to parse JSON response for total pages.")
            raise
        except Exception as e:
            logger.exception(f"Error processing response for total pages: {e}")
            raise

    def _extract_product_ids(self, page: int, category_id: int, url_key: str) -> list[int]:
        """
        Extracts product IDs from a specific listings page.

        Args:
            page (int): The page number to extract product IDs from.
            category_id (int): The ID of the category to extract product IDs from.
            url_key (str): The URL key of the category to extract product IDs from.

        Returns:
            list[int]: A list of product IDs.

        Raises:
            JSONDecodeError: If the JSON response is invalid.
            Exception: If an error occurs while extracting product IDs.
        """
        logger.info(f"Extracting product IDs from page {page} for category {category_id}...")

        params = self.request_params_base.copy()
        params["category"] = category_id
        params["urlKey"] = url_key
        params["page"] = page

        data = make_http_request(
            url=self.request_listings_url, headers=self.request_headers, params=params
        )
        if not data:
            return []

        try:
            items = data.get("data", [])
            ids = [item.get("id") for item in items if item.get("id") is not None]
            logger.info(f"Found {len(ids)} product IDs on page {page}.")
            return ids
        except json.JSONDecodeError:
            logger.error(f"Failed to parse JSON response for product IDs on page {page}.")
        except Exception as e:
            logger.error(f"Error processing response for product IDs on page {page}: {e}")

        return []

    def _extract_seller_info(self, seller_data: dict) -> dict:
        """Extract and format seller information from product data."""
        seller_id = str(seller_data.get("id"))
        seller_name = seller_data.get("name")
        seller_link = seller_data.get("link")
        seller_logo = seller_data.get("logo")

        # Fix logo URL if needed
        if seller_logo and not seller_logo.startswith("https://"):
            seller_logo = f"https://vcdn.tikicdn.com/cache/w100/ts/seller/{seller_logo}"

        return {
            "seller_id": seller_id,
            "seller_name": seller_name,
            "seller_link": seller_link,
            "seller_logo": seller_logo,
        }

    def _extract_page_count(self, specifications: list) -> Optional[int]:
        """Extract page count from product specifications."""
        if not specifications:
            return None

        try:
            spec_data = specifications[0]
            attributes = spec_data.get("attributes", [])

            if not attributes:
                return None

            pages_attr = next(
                (attr for attr in attributes if attr.get("code") == "number_of_page"),
                None,
            )

            if pages_attr and pages_attr.get("value"):
                return int(pages_attr.get("value"))
        except Exception as e:
            logger.error(f"Error extracting page count: {e}")

        return None

    def _extract_authors_data(self, authors: list) -> list:
        """Extract and format authors data."""
        authors_data = []
        for author in authors:
            authors_data.append(
                {
                    "id": str(author.get("id")),
                    "name": author.get("name"),
                    "slug": author.get("slug"),
                }
            )
        return authors_data

    def _extract_category_info(self, data: dict) -> dict:
        """Extract category information from product data."""
        category = data.get("categories", {})
        if category:
            return {
                "category_id": str(category.get("id")),
                "category_name": category.get("name"),
            }

        # Fallback to breadcrumbs if categories is empty
        breadcrumbs = data.get("breadcrumbs", [])
        if breadcrumbs and len(breadcrumbs) > 1:
            return {
                "category_id": str(breadcrumbs[-2].get("id")),
                "category_name": breadcrumbs[-2].get("name"),
            }

        return {"category_id": None, "category_name": None}

    def _extract_product_details(self, product_id: int) -> Optional[dict]:
        """
        Extracts detailed information for a single product ID.

        Args:
            product_id (int): The ID of the product to extract details for.

        Returns:
            Optional[dict]: Dictionary containing the product details.
        """
        logger.info(f"Extracting details for product ID: {product_id}")
        url = f"{self.request_product_url_base}{product_id}"
        data = make_http_request(url=url, headers=self.request_headers)

        if not data:
            return None

        try:
            # Extract data using helper methods
            seller_info = self._extract_seller_info(data.get("current_seller", {}))
            page_count = self._extract_page_count(data.get("specifications", []))
            authors_data = self._extract_authors_data(data.get("authors", []))
            category_info = self._extract_category_info(data)

            # Combine all information into a single details dictionary
            details = {
                "id": str(data.get("id")),
                "name": data.get("name"),
                "description": data.get("short_description"),
                "price": data.get("price"),
                "original_price": data.get("original_price"),
                "rating_average": data.get("rating_average"),
                "review_count": data.get("review_count"),
                "page_count": page_count,
                "authors": authors_data,
                **seller_info,
                **category_info,
            }
            return details

        except json.JSONDecodeError:
            logger.error(f"Failed to parse JSON response for product details (ID: {product_id}).")
        except Exception as e:
            logger.error(f"Error processing product details for ID {product_id}: {e}")

        return None

    def _extract_products_details(self, product_ids: list[int]) -> list[dict]:
        """
        Extracts details for a list of product IDs.

        Args:
            product_ids (list[int]): A list of product IDs to extract details for.

        Returns:
            list[dict]: A list of dictionaries containing the product details.
        """
        products_details = []
        for product_id in product_ids:
            details = self._extract_product_details(product_id)
            if details:
                products_details.append(details)
            time.sleep(self.request_delay)

        logger.info(
            f"Successfully extracted details for {len(products_details)} out of {len(product_ids)} products."
        )
        return products_details

    def _upsert_resource(
        self, resource: str, data: TikiSeller | TikiBook | TikiAuthor | TikiCategory
    ) -> None:
        """
        Upserts a resource via API.

        Args:
            resource (str): The resource to upsert.
            data (TikiSeller | TikiBook | TikiAuthor | TikiCategory): The data to upsert.
        """
        existing = make_http_request(
            url=f"{self.api_url}/{resource}/{data.secondary_id}",
            method="GET",
            params={"secondary": "true"},
            use_auth=True,
            auth_token_manager=self.auth_token_manager,
        )

        # Handle the case when existing is None
        existing_resource = {} if existing is None else existing.get("data", {})

        # Create data payload, handling special case for books with author_ids
        payload = data.model_dump()

        if existing_resource:
            logger.info(
                f"{resource[:-1].capitalize()} {data.secondary_id} already exists. Performing update..."
            )
            make_http_request(
                url=f"{self.api_url}/{resource}/{existing_resource.get('id')}",
                method="PUT",
                data=payload,
                use_auth=True,
                auth_token_manager=self.auth_token_manager,
            )
        else:
            logger.info(
                f"{resource[:-1].capitalize()} {data.secondary_id} does not exist. Creating new {resource[:-1]}..."
            )
            make_http_request(
                url=f"{self.api_url}/{resource}",
                method="POST",
                data=payload,
                use_auth=True,
                auth_token_manager=self.auth_token_manager,
            )

    def _upsert_specific_data(
        self, data: TikiSeller | TikiBook | TikiAuthor | TikiCategory
    ) -> None:
        """
        Upserts specific data into the database.

        Args:
            data (TikiSeller | TikiBook | TikiAuthor | TikiCategory): The data to upsert.
        """
        resource_map = {
            TikiSeller: "sellers",
            TikiBook: "books",
            TikiAuthor: "authors",
            TikiCategory: "categories",
        }

        for model_cls, resource in resource_map.items():
            if isinstance(data, model_cls):
                self._upsert_resource(resource=resource, data=data)

    def _upsert_data(self, data: list[dict]) -> None:
        """
        Upserts data into the database via API.

        Args:
            data (list[dict]): List of dictionaries containing the product details.
        """
        if not data:
            logger.info("No data provided to be upserted.")
            return

        # Use sets to avoid duplicates within this batch
        seller_ids = set()
        product_ids = set()
        category_ids = set()
        author_ids = set()

        for record in data:
            # Seller
            if record["seller_id"] and record["seller_id"] not in seller_ids:
                seller_to_upsert = TikiSeller(
                    secondary_id=record["seller_id"],
                    name=record["seller_name"],
                    link=record["seller_link"],
                    logo=record["seller_logo"],
                )
                self._upsert_specific_data(data=seller_to_upsert)
                seller_ids.add(record["seller_id"])

            # Product & Inventory
            if record["id"] and record["id"] not in product_ids:
                product_to_upsert = TikiBook(
                    secondary_id=record["id"],
                    name=record["name"],
                    price=record["price"],
                    original_price=record["original_price"],
                    rating_average=record["rating_average"],
                    review_count=record["review_count"],
                    page_count=record["page_count"],
                    author_ids=[author["id"] for author in record["authors"]],
                )
                self._upsert_specific_data(data=product_to_upsert)
                product_ids.add(record["id"])

            # Authors - Handle all authors instead of just the first one
            if record["authors"]:
                # Process all authors
                for author in record["authors"]:
                    if author["id"] not in author_ids:
                        author_to_upsert = TikiAuthor(
                            secondary_id=author["id"],
                            name=author["name"],
                            slug=author["slug"],
                        )
                        self._upsert_specific_data(data=author_to_upsert)
                        author_ids.add(author["id"])

            # Category
            if record["category_id"] and record["category_id"] not in category_ids:
                category_to_upsert = TikiCategory(
                    secondary_id=record["category_id"],
                    name=record["category_name"],
                )
                self._upsert_specific_data(data=category_to_upsert)
                category_ids.add(record["category_id"])

    def _process_category(self, category_id: int, url_key: str) -> None:
        """
        Handles crawling and inserting data for a single category.

        Args:
            category_id (int): The ID of the category to process.
            url_key (str): The URL key of the category to process.
        """
        total_pages = self._get_total_pages(category_id=category_id, url_key=url_key)
        if total_pages == 0:
            logger.warning(
                f"No pages found or error occurred for category {category_id}. Skipping."
            )
            return

        logger.info(
            f"Starting data extraction for {total_pages} pages in category {category_id}..."
        )

        for page in range(1, total_pages + 1):
            logger.info(f"Processing page {page}/{total_pages} for category {category_id}...")
            product_ids = self._extract_product_ids(
                page=page, category_id=category_id, url_key=url_key
            )
            if not product_ids:
                logger.warning(f"No product IDs found on page {page}. Moving to next.")
                time.sleep(self.request_delay)
                continue

            products_details = self._extract_products_details(product_ids)
            print(products_details[-3])
            if products_details:
                self._upsert_data(data=[products_details[-3]])
            else:
                logger.info(f"No valid product details extracted for page {page}.")

            logger.info(f"Waiting {self.request_delay}s before next page...")
            time.sleep(self.request_delay)

            break

    def run(self) -> None:
        """
        Main execution method for the crawler pipeline.
        """
        logger.info("Tiki Crawler starting run...")
        start_time = time.time()

        for category_id, url_key in self.categories.items():
            logger.info(f"Processing Category ID: {category_id}, URL Key: {url_key}")
            self._process_category(category_id=category_id, url_key=url_key)
            logger.info(f"Finished Category ID: {category_id}, URL Key: {url_key}")

        end_time = time.time()
        logger.info(f"Tiki Crawler finished run. Total time: {end_time - start_time:.2f} seconds.")


if __name__ == "__main__":
    try:
        # Set up logging to debug level for more detailed information
        logging.getLogger().setLevel(logging.DEBUG)
        logging.info("Starting Tiki Crawler...")

        # Create a crawler with admin credentials - you can override these from command line if needed
        import sys

        admin_email = sys.argv[1] if len(sys.argv) > 1 else "admin@example.com"
        admin_password = sys.argv[2] if len(sys.argv) > 2 else "Admin123!"

        logging.info(f"Using credentials: {admin_email}")

        # Test authentication explicitly first
        crawler = TikiCrawler(admin_email=admin_email, admin_password=admin_password)

        # Try to get a token to verify authentication works
        token = crawler.auth_token_manager.get_token()
        if token:
            logging.info("✅ Authentication successful!")
            crawler.run()
        else:
            logging.error("❌ Authentication failed. Please check your credentials and API URL.")
            logging.info(f"API URL: {API_URL}")

    except Exception as e:
        logging.exception(f"Error in crawler main execution: {e}")
