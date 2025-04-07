import json
import logging
import time
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
TIKI_PRODUCT_LISTINGS_PAGE_PARAMS = {
    "limit": "10",
}
TIKI_CATEGORIES = {
    8322: "nha-sach-tiki",
}
TIKI_BASE_PRODUCT_LISTINGS = "https://tiki.vn/api/personalish/v1/blocks/listings"
TIKI_BASE_SPECIFIC_PRODUCT = "https://tiki.vn/api/v2/products/"

import os

API_URL = os.getenv("API_URL", "http://api:5000/api")

import requests
from typing import Any


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


logger = logging.getLogger(__name__)


class TikiSeller(BaseModel):
    secondary_id: int
    name: str
    link: str
    logo: str


class TikiProduct(BaseModel):
    secondary_id: int
    name: str
    price: float
    original_price: float
    rating_average: float
    review_count: int
    page_count: int


class TikiInventory(BaseModel):
    secondary_id: int
    quantity: int


class TikiAuthor(BaseModel):
    secondary_id: int
    name: str
    slug: str


class TikiCategory(BaseModel):
    secondary_id: int
    name: str


class TikiCrawler:
    def __init__(self):
        self.categories = TIKI_CATEGORIES
        self.request_params_base = TIKI_PRODUCT_LISTINGS_PAGE_PARAMS.copy()
        self.request_headers = TIKI_HEADERS
        self.request_listings_url = TIKI_BASE_PRODUCT_LISTINGS
        self.request_product_url_base = TIKI_BASE_SPECIFIC_PRODUCT
        self.request_delay = TIKI_REQUEST_DELAY

        self.api_url = API_URL

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
        seller_id = seller_data.get("id")
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
                    "id": author.get("id"),
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
                "category_id": category.get("id"),
                "category_name": category.get("name"),
            }

        # Fallback to breadcrumbs if categories is empty
        breadcrumbs = data.get("breadcrumbs", [])
        if breadcrumbs and len(breadcrumbs) > 1:
            return {
                "category_id": breadcrumbs[-2].get("id"),
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
                "id": data.get("id"),
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
        self, resource: str, data: TikiSeller | TikiProduct | TikiAuthor | TikiCategory
    ) -> None:
        """
        Upserts a resource via API.

        Args:
            resource (str): The resource to upsert.
            data (TikiSeller | TikiProduct | TikiAuthor | TikiCategory): The data to upsert.
        """
        existing = make_http_request(
            url=f"{self.api_url}/{resource}/{data.id}",
            method="GET",
            params={"secondary": True},
        )

        # Handle the case when existing is None
        existing_resource = {} if existing is None else existing.get("data", {})

        if existing_resource:
            logger.info(
                f"{resource[:-1].capitalize()} {data.id} already exists. Performing update..."
            )
            make_http_request(
                url=f"{self.api_url}/{resource}/{existing_resource.get('id')}",
                method="PUT",
                # json=data.model_dump(),
            )
        else:
            logger.info(
                f"{resource[:-1].capitalize()} {data.id} does not exist. Creating new {resource[:-1]}..."
            )
            make_http_request(
                url=f"{self.api_url}/{resource}",
                method="POST",
                # json=data.model_dump(),
            )

    def _upsert_specific_data(
        self, data: TikiSeller | TikiProduct | TikiAuthor | TikiCategory
    ) -> None:
        """
        Upserts specific data into the database.

        Args:
            data (TikiSeller | TikiProduct | TikiAuthor | TikiCategory): The data to upsert.
        """
        resource_map = {
            TikiSeller: "sellers",
            TikiProduct: "products",
            TikiAuthor: "authors",
            TikiCategory: "categories",
        }

        for model_cls, resource in resource_map.items():
            if isinstance(data, model_cls):
                self._upsert_resource(resource=resource, data=data)

    def _upsert_data(self, data: list[dict]) -> None:
        """

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
            # --- Seller ---
            if record["seller_id"] and record["seller_id"] not in seller_ids:
                seller_to_upsert = TikiSeller(
                    secondary_id=record["seller_id"],
                    name=record["seller_name"],
                    link=record["seller_link"],
                    logo=record["seller_logo"],
                )
                self._upsert_specific_data(data=seller_to_upsert)
                seller_ids.add(record["seller_id"])

            # --- Product & Inventory ---
            if record["id"] and record["id"] not in product_ids:
                product_to_upsert = TikiProduct(
                    secondary_id=record["id"],
                    name=record["name"],
                    price=record["price"],
                    original_price=record["original_price"],
                    rating_average=record["rating_average"],
                    review_count=record["review_count"],
                    page_count=record["page_count"],
                )
                self._upsert_specific_data(data=product_to_upsert)
                inventory_to_upsert = TikiInventory(
                    secondary_id=record["id"],  # Matches product ID
                    quantity=10000,
                )
                self._upsert_specific_data(data=inventory_to_upsert)
                product_ids.add(record["id"])

            # --- Author ---
            if record["authors"] and record["authors"] not in author_ids:
                author_to_upsert = TikiAuthor(
                    secondary_id=record["authors"][0]["id"],
                    name=record["authors"][0]["name"],
                    slug=record["authors"][0]["slug"],
                )
                self._upsert_specific_data(data=author_to_upsert)
                author_ids.add(record["authors"][0]["id"])

            # --- Category ---
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
            # if products_details:
            #     self._upsert_data(data=products_details)
            # else:
            #     logger.info(f"No valid product details extracted for page {page}.")

            # logger.info(f"Waiting {self.request_delay}s before next page...")
            # time.sleep(self.request_delay)

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
    crawler = TikiCrawler()
    crawler.run()
