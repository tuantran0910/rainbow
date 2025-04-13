import random
from typing import Any

import dagster as dg
from faker import Faker

from assets.helpers import AuthTokenManager
from assets.helpers import make_http_request
from constants import API_BASE_URL
from constants import DAGSTER_METADATA
from constants import DAGSTER_MOCKING_ASSET_GROUP
from constants import DAGSTER_TAGS
from constants import DEFAULT_USER_PASSWORD
from constants import MAX_BOOKS_PER_REQUEST
from constants import MAX_ITEMS_PER_ORDER
from constants import MAX_QUANTITY_PER_BOOK
from resources.psql_resource import PostgresResource

logger = dg.get_dagster_logger()


@dg.op(
    name="get_users",
    description="Get all users from the PostgreSQL database.",
    out=dg.Out(description="A list of user emails."),
)
def get_users(rainbow_psql_resource: PostgresResource) -> list[str]:
    """
    Get all users from the PostgreSQL database.

    Args:
        rainbow_psql_resource (PostgresResource): A resource for interacting with the PostgreSQL database.

    Returns:
        list[str]: A list of user emails.

    Raises:
        Exception: If database query fails.
    """
    try:
        query = """
            SELECT email
            FROM users
            WHERE deleted_at IS NULL
        """
        users = rainbow_psql_resource.fetchall(query=query)
        if not users:
            raise dg.DagsterError("No users found in the database")

        return [user[0] for user in users]
    except Exception as e:
        raise dg.DagsterError(f"Failed to fetch users: {e}")


@dg.op(
    name="get_available_books",
    description="Fetch available books from the Rainbow API.",
    out=dg.Out(description="A list of book objects with their details."),
)
def get_available_books() -> list[dict[str, Any]]:
    """
    Fetch available books from the Rainbow API.

    Returns:
        list[dict[str, Any]]: A list of book objects with their details.

    Raises:
        Exception: If API request fails.
    """
    books_url = f"{API_BASE_URL}/api/books"

    try:
        # Get total number of books
        response_data = make_http_request(f"{books_url}/total", method="GET")
        total_books = response_data["data"]["total"]

        # Get random books at random offset, maximum MAX_BOOKS_PER_REQUEST books
        if total_books > 0:
            max_offset = max(0, total_books - 1)
            offset = random.randint(0, max_offset)
            response_data = make_http_request(
                books_url,
                method="GET",
                params={"offset": offset, "limit": MAX_BOOKS_PER_REQUEST},
            )
            return response_data["data"]["books"]
        return []
    except Exception as e:
        raise dg.DagsterError(f"Failed to fetch books: {e}")


@dg.op(
    name="get_available_payments",
    description="Fetch all available payment methods from the Rainbow API.",
    out=dg.Out(description="A list of payment objects."),
)
def get_available_payments() -> list[dict[str, Any]]:
    """
    Fetch all available payment methods from the Rainbow API.

    Returns:
        list[dict[str, Any]]: A list of payment objects.

    Raises:
        Exception: If API request fails.
    """
    payments_url = f"{API_BASE_URL}/api/payments"
    try:
        response_data = make_http_request(payments_url, method="GET")
        return response_data["data"]["payments"]
    except Exception as e:
        raise dg.DagsterError(f"Failed to fetch payments: {e}")


@dg.op(
    name="get_available_promotions",
    description="Fetch available promotions from the Rainbow API.",
    out=dg.Out(description="A list of promotion objects."),
)
def get_available_promotions() -> list[dict[str, Any]]:
    """
    Fetch available promotions from the Rainbow API.

    Returns:
        list[dict[str, Any]]: A list of promotion objects.

    Raises:
        Exception: If API request fails.
    """
    promotions_url = f"{API_BASE_URL}/api/promotions"
    try:
        response_data = make_http_request(promotions_url, method="GET")
        return response_data["data"]["promotions"]
    except Exception as e:
        raise dg.DagsterError(f"Failed to fetch promotions: {e}")


@dg.op(
    name="process_mock_orders",
    description="Process mock orders for users.",
    ins={
        "users": dg.In(description="List of user emails to create orders for."),
        "books": dg.In(description="List of available books."),
        "payments": dg.In(description="List of available payment methods."),
        "promotions": dg.In(description="List of available promotions."),
    },
    out=dg.Out(description="List of created order objects."),
)
def process_mock_orders(
    users: list[str],
    books: list[dict[str, Any]],
    payments: list[dict[str, Any]],
    promotions: list[dict[str, Any]],
) -> list[dict[str, Any]]:
    """
    Create a mock order via the Rainbow API.

    Args:
        users (list[str]): List of user emails to create orders for.
        books (list[dict[str, Any]]): List of available books.
        payments (list[dict[str, Any]]): List of available payment methods.
        promotions (list[dict[str, Any]]): List of available promotions.

    Returns:
        list[dict[str, Any]]: List of created order objects.

    Raises:
        Exception: If API request fails.
    """
    created_orders = []
    faker = Faker(locale="vi_VN")

    for user in users:
        logger.info(f"Processing mock order for user {user}...")

        try:
            # Randomly select items
            num_items = random.randint(1, min(MAX_ITEMS_PER_ORDER, len(books)))
            selected_books = random.sample(books, num_items)

            # Create order items
            order_items = [
                {
                    "book_id": book["id"],
                    "quantity": random.randint(1, MAX_QUANTITY_PER_BOOK),
                }
                for book in selected_books
            ]

            # Randomly select payment
            payment = random.choice(payments)
            payment_id = payment.get("id")

            # Randomly select promotion (if available)
            promotion_id = None
            if promotions:
                promotion = random.choice(promotions)
                promotion_id = promotion.get("id")

            # Create order payload
            order_payload = {
                "payment_id": payment_id,
                "shipping_address": faker.state(),
                "order_items": order_items,
            }

            # Optionally add promotion
            if promotion_id:
                order_payload["promotion_id"] = promotion_id

            # Make the API request
            orders_url = f"{API_BASE_URL}/api/orders"
            auth_token_manager = AuthTokenManager(
                email=user, password=DEFAULT_USER_PASSWORD, api_url=API_BASE_URL
            )
            response_data = make_http_request(
                orders_url,
                method="POST",
                data=order_payload,
                use_auth=True,
                auth_token_manager=auth_token_manager,
            )
            order_id = response_data["data"]["id"]
            logger.info(f"Order created successfully: {order_id}")
            created_orders.append(response_data["data"])
        except Exception as e:
            logger.error(f"Failed to mock order transaction for user {user}: {e}")
            continue

    return created_orders


@dg.graph_asset(
    name="order_transactions",
    description="A graph asset that creates mock orders for all users in the database.",
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
    group_name=DAGSTER_MOCKING_ASSET_GROUP,
    kinds={"python"},
)
def order_transactions() -> list[dict[str, Any]]:
    """
    A graph asset that creates mock orders for all users in the database.

    This asset:
    1. Retrieves all users from the database
    2. Fetches available books, payment methods, and promotions
    3. Creates a random order for each user

    Returns:
        list[dict[str, Any]]: A list of created order objects.
    """
    users = get_users()
    books = get_available_books()
    payments = get_available_payments()
    promotions = get_available_promotions()

    # Validate we have necessary data
    if not users:
        logger.warning("No users found to create orders for")
        return dg.MaterializeResult(metadata={"number_of_orders": 0})

    if not books:
        logger.warning("No books available for orders")
        return dg.MaterializeResult(metadata={"number_of_orders": 0})

    if not payments:
        logger.warning("No payment methods available")
        return dg.MaterializeResult(metadata={"number_of_orders": 0})

    created_orders = process_mock_orders(
        users=users, books=books, payments=payments, promotions=promotions
    )

    return created_orders
