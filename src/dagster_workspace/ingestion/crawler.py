import json
import random
import time
from typing import Optional

import dagster as dg
from pydantic import BaseModel

from shared.constants import API_BASE_URL
from shared.constants import INVENTORY_MAX_STOCK
from shared.constants import INVENTORY_MIN_STOCK
from shared.constants import TIKI_BASE_PRODUCT_LISTINGS
from shared.constants import TIKI_BASE_SPECIFIC_PRODUCT
from shared.constants import TIKI_CATEGORIES
from shared.constants import TIKI_HEADERS
from shared.constants import TIKI_MAX_PAGES_PER_CATEGORY
from shared.constants import TIKI_PRODUCT_LISTINGS_PAGE_PARAMS
from shared.constants import TIKI_REQUEST_DELAY
from shared.helpers import AuthTokenManager
from shared.helpers import make_http_request
from shared.helpers import sanitize_text


logger = dg.get_dagster_logger(__name__)


class TikiSeller(BaseModel):
    secondary_id: str
    name: str
    link: str
    logo: str


class TikiBook(BaseModel):
    secondary_id: str
    category_id: str
    seller_id: str
    name: str
    description: str
    price: float
    original_price: float
    rating_average: float
    review_count: Optional[int] = None
    page_count: Optional[int] = None
    author_ids: list[str] = []
    stock: Optional[int] = None


class TikiAuthor(BaseModel):
    secondary_id: str
    name: str
    slug: str


class TikiCategory(BaseModel):
    secondary_id: str
    name: str


class TikiCrawler:
    def __init__(self, admin_email: str, admin_password: str):
        """
        Initializes the TikiCrawler with the provided admin credentials.

        Args:
            admin_email (str): Admin email for authentication.
            admin_password (str): Admin password for authentication.
        """
        self.categories = TIKI_CATEGORIES
        self.request_params_base = TIKI_PRODUCT_LISTINGS_PAGE_PARAMS.copy()
        self.request_headers = TIKI_HEADERS
        self.request_listings_url = TIKI_BASE_PRODUCT_LISTINGS
        self.request_product_url_base = TIKI_BASE_SPECIFIC_PRODUCT
        self.request_delay = TIKI_REQUEST_DELAY

        # API Configuration
        self.api_base_url = API_BASE_URL
        self.api_url = f"{API_BASE_URL}/api"
        self.auth_api_url = f"{API_BASE_URL}/auth/login"

        # For resource endpoints
        self.auth_token_manager = AuthTokenManager(
            email=admin_email, password=admin_password, auth_api_url=self.auth_api_url
        )

        logger.info(f"Initialized TikiCrawler with API URL: {self.api_url}")
        logger.info(f"Auth base URL: {self.auth_api_url}")

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
            else:
                logger.warning("Could not find 'last_page' in response paging data.")
            return last_page
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
                    "slug": author.get("slug", ""),
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

            # Clean text fields that may contain special characters
            name = sanitize_text(data.get("name"))
            description = sanitize_text(data.get("short_description"))

            # Combine all information into a single details dictionary
            details = {
                "id": str(data.get("id")),
                "name": name,
                "description": description,
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
    ) -> Optional[str]:
        """
        Upserts a resource via API.

        Args:
            resource (str): The resource to upsert.
            data (TikiSeller | TikiBook | TikiAuthor | TikiCategory): The data to upsert.

        Returns:
            Optional[str]: The UUID of the created/updated resource, or None if the operation failed.
        """
        existing = make_http_request(
            url=f"{self.api_url}/{resource}/{data.secondary_id}",
            method="GET",
            params={"secondary": "true"},
        )
        existing_resource = {} if existing is None else existing.get("data", {})

        # Create data payload, handling special case for books with author_ids
        payload = data.model_dump()

        if existing_resource:
            logger.info(
                f"{resource[:-1].capitalize()} {data.secondary_id} already exists. Performing update..."
            )
            response = make_http_request(
                url=f"{self.api_url}/{resource}/{existing_resource.get('id')}",
                method="PATCH",
                data=payload,
                use_auth=True,
                auth_token_manager=self.auth_token_manager,
            )
            if response and "data" in response:
                return response["data"].get("id")
        else:
            logger.info(
                f"{resource[:-1].capitalize()} {data.secondary_id} does not exist. Creating new {resource[:-1]}..."
            )
            response = make_http_request(
                url=f"{self.api_url}/{resource}",
                method="POST",
                data=payload,
                use_auth=True,
                auth_token_manager=self.auth_token_manager,
            )
            if response and "data" in response:
                return response["data"].get("id")

        logger.error(f"Failed to upsert {resource} with secondary_id {data.secondary_id}")
        return None

    def _upsert_specific_data(
        self, data: TikiSeller | TikiBook | TikiAuthor | TikiCategory
    ) -> Optional[str]:
        """
        Upserts specific data into the database.

        Args:
            data (TikiSeller | TikiBook | TikiAuthor | TikiCategory): The data to upsert.

        Returns:
            Optional[str]: The UUID of the created/updated resource, or None if the operation failed.
        """
        resource_map = {
            TikiSeller: "sellers",
            TikiBook: "books",
            TikiAuthor: "authors",
            TikiCategory: "categories",
        }

        for model_cls, resource in resource_map.items():
            if isinstance(data, model_cls):
                return self._upsert_resource(resource=resource, data=data)

        return None

    def _upsert_data(self, data: list[dict]) -> None:  # noqa: C901
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
        book_ids = set()
        category_ids = set()
        author_ids = set()

        # Store UUIDs for each resource type
        seller_uuids = {}
        category_uuids = {}
        author_uuids = {}

        for record in data:
            # Seller
            if record["seller_id"] and record["seller_id"] not in seller_ids:
                # Simple check: skip if seller name is empty
                if not record["seller_name"] or record["seller_name"].strip() == "":
                    logger.warning(f"Skipping seller {record['seller_id']} - empty name")
                    continue

                seller_to_upsert = TikiSeller(
                    secondary_id=record["seller_id"],
                    name=sanitize_text(record["seller_name"]),
                    link=record["seller_link"],
                    logo=record["seller_logo"],
                )
                seller_uuid = self._upsert_specific_data(data=seller_to_upsert)
                if seller_uuid:
                    seller_uuids[record["seller_id"]] = seller_uuid
                    seller_ids.add(record["seller_id"])

            # Authors
            if record["authors"]:
                # Process all authors
                for author in record["authors"]:
                    if author["id"] not in author_ids:
                        if not author["name"] or author["name"].strip() == "":
                            logger.warning(f"Skipping author {author['id']} - empty name")
                            continue

                        author_to_upsert = TikiAuthor(
                            secondary_id=author["id"],
                            name=sanitize_text(author["name"]),
                            slug=author["slug"],
                        )
                        author_uuid = self._upsert_specific_data(data=author_to_upsert)
                        if author_uuid:
                            author_uuids[author["id"]] = author_uuid
                            author_ids.add(author["id"])

            # Category
            if record["category_id"] and record["category_id"] not in category_ids:
                if not record["category_name"] or record["category_name"].strip() == "":
                    logger.warning(f"Skipping category {record['category_id']} - empty name")
                    continue

                category_to_upsert = TikiCategory(
                    secondary_id=record["category_id"],
                    name=sanitize_text(record["category_name"]),
                )
                category_uuid = self._upsert_specific_data(data=category_to_upsert)
                if category_uuid:
                    category_uuids[record["category_id"]] = category_uuid
                    category_ids.add(record["category_id"])

            # Product & Inventory
            if record["id"] and record["id"] not in book_ids:
                # First, check if the book already exists by secondary ID
                existing_book = make_http_request(
                    url=f"{self.api_url}/books/{record['id']}",
                    method="GET",
                    params={"secondary": "true"},
                )
                existing_book_data = existing_book.get("data", {})

                # Generate a random stock value in case the book is new or needs to be restocked
                new_stock = random.randint(INVENTORY_MIN_STOCK, INVENTORY_MAX_STOCK)

                # Get the UUIDs for the related entities
                category_uuid = category_uuids.get(record["category_id"])
                seller_uuid = seller_uuids.get(record["seller_id"])
                author_uuids_list = [
                    author_uuids.get(author["id"])
                    for author in record["authors"]
                    if author["id"] in author_uuids
                ]

                # Only create the book if we have all required UUIDs
                if category_uuid and seller_uuid:
                    if not record["name"] or record["name"].strip() == "":
                        logger.warning(f"Skipping book {record['id']} - empty name")
                        continue

                    book_to_upsert = TikiBook(
                        secondary_id=record["id"],
                        category_id=category_uuid,
                        seller_id=seller_uuid,
                        name=sanitize_text(record["name"]),
                        description=sanitize_text(record["description"]),
                        price=record["price"],
                        original_price=record["original_price"],
                        rating_average=record["rating_average"],
                        review_count=record["review_count"],
                        page_count=record["page_count"],
                        author_ids=author_uuids_list,
                        stock=new_stock if not existing_book_data else None,
                    )
                    book_uuid = self._upsert_specific_data(data=book_to_upsert)
                    if book_uuid:
                        book_ids.add(record["id"])

                        # Check inventory status if the book exists
                        inventory_needs_update = False

                        if existing_book_data:
                            existing_inventory = existing_book_data.get("inventory", {})
                            if existing_inventory:
                                current_stock = existing_inventory.get("stock", 0)

                                # Check if stock is low and needs replenishment
                                if current_stock <= 0:
                                    inventory_needs_update = True
                                    logger.info(
                                        f"Book {record['id']} has low inventory ({current_stock}). Restocking..."
                                    )

                        # Update inventory if needed when it's low
                        if inventory_needs_update:
                            logger.info(
                                f"Updating inventory for book {record['id']} with new stock: {new_stock}"
                            )
                            make_http_request(
                                url=f"{self.api_url}/books/{book_uuid}",
                                method="PATCH",
                                data={
                                    "stock": new_stock,
                                },
                                use_auth=True,
                                auth_token_manager=self.auth_token_manager,
                            )
                else:
                    logger.warning(
                        f"Missing required UUIDs for book {record['id']}. Skipping book creation."
                    )
                    if not category_uuid:
                        logger.warning(f"Missing category UUID for book {record['id']}")
                    if not seller_uuid:
                        logger.warning(f"Missing seller UUID for book {record['id']}")

    def _process_crawling(self, category_id: int, url_key: str) -> None:
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
            if products_details:
                self._upsert_data(data=products_details)
            else:
                logger.info(f"No valid product details extracted for page {page}.")

            logger.info(f"Waiting {self.request_delay}s before next page...")
            time.sleep(self.request_delay)

            # Check if we've reached the maximum pages limit (if configured)
            if TIKI_MAX_PAGES_PER_CATEGORY > 0 and page >= TIKI_MAX_PAGES_PER_CATEGORY:
                logger.info(
                    f"Reached maximum pages limit ({TIKI_MAX_PAGES_PER_CATEGORY}). Stopping crawl for this category."
                )
                break

    def run(self) -> None:
        """
        Main execution method for the crawler pipeline.
        """
        logger.info("Tiki Crawler starting run...")
        start_time = time.time()

        for category_id, url_key in self.categories.items():
            logger.info(f"Processing Category ID: {category_id}, URL Key: {url_key}")
            self._process_crawling(category_id=category_id, url_key=url_key)
            logger.info(f"Finished Category ID: {category_id}, URL Key: {url_key}")

        end_time = time.time()
        logger.info(f"Tiki Crawler finished run. Total time: {end_time - start_time:.2f} seconds.")


if __name__ == "__main__":
    try:
        logger.info("Starting Tiki Crawler...")

        # Create a crawler with admin credentials - you can override these from command line if needed
        import sys

        admin_email = sys.argv[1] if len(sys.argv) > 1 else "admin@example.com"
        admin_password = sys.argv[2] if len(sys.argv) > 2 else "Admin123!"

        logger.info(f"Using credentials: {admin_email}")

        # Test authentication explicitly first
        crawler = TikiCrawler(admin_email=admin_email, admin_password=admin_password)

        # Try to get a token to verify authentication works
        token = crawler.auth_token_manager.get_token()
        if token:
            logger.info("✅ Authentication successful!")
            crawler.run()
        else:
            logger.error("❌ Authentication failed. Please check your credentials and API URL.")
            logger.info(f"API URL: {API_BASE_URL}")

    except Exception as e:
        logger.exception(f"Error in crawler main execution: {e}")
