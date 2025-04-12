import random
from datetime import datetime
from typing import Any
from typing import Optional

import dagster as dg

from assets.helpers import AuthTokenManager
from assets.helpers import make_http_request
from constants import API_BASE_URL
from constants import DEFAULT_USER_PASSWORD
from constants import MAX_BOOKS_PER_REQUEST
from constants import MAX_ITEMS_PER_ORDER
from resources.psql_resource import PostgresResource

logger = dg.get_dagster_logger()


@dg.op
def get_users(psql_resource: PostgresResource) -> list[str]:
    """
    Get all users from the PostgreSQL database.

    Args:
        psql_resource (PostgresResource): A resource for interacting with the PostgreSQL database.

    Returns:
        list[str]: A list of user emails.

    Raises:
        Exception: If database query fails.
    """
    try:
        users = psql_resource.fetchall("SELECT email FROM users")
        return [user[0] for user in users]
    except Exception as e:
        raise dg.DagsterError(f"Failed to fetch users: {e}")


@dg.op
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


@dg.op
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


@dg.op
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


@dg.op
def create_mock_order(
    user: str,
    order_items: list[dict[str, Any]],
    payment_id: str,
    promotion_id: Optional[str] = None,
) -> dict[str, Any]:
    """
    Create a mock order via the Rainbow API.

    Args:
        user (str): User email to create the order for.
        order_items (list[dict[str, Any]]): List of items to include in the order.
        payment_id (str): ID of the payment method to use.
        promotion_id (Optional[str]): Optional promotion ID to apply.

    Returns:
        dict[str, Any]: The created order object.

    Raises:
        Exception: If API request fails.
    """
    logger.info(f"Creating mock order for user {user}...")

    # Create order payload
    order_payload = {
        "payment_id": payment_id,
        "shipping_address": f"Mock Address {datetime.now().strftime('%Y%m%d%H%M%S')}",
        "order_items": order_items,
    }

    # Optionally add promotion
    if promotion_id:
        order_payload["promotion_id"] = promotion_id

    # Make the API request
    orders_url = f"{API_BASE_URL}/api/orders"
    try:
        auth_token_manager = AuthTokenManager(
            email=user, password=DEFAULT_USER_PASSWORD, api_url=API_BASE_URL
        )
        response_data = make_http_request(
            orders_url,
            method="POST",
            json=order_payload,
            use_auth=True,
            auth_token_manager=auth_token_manager,
        )
        order_id = response_data["data"]["order"]["id"]
        logger.info(f"Order created successfully: {order_id}")
        return response_data["data"]["order"]
    except Exception as e:
        raise dg.DagsterError(f"Failed to mock order transaction: {e}")


@dg.graph_asset
def order_transactions(psql_resource: PostgresResource) -> dg.MaterializeResult:
    """
    A graph asset that creates mock orders for all users in the database.

    This asset:
    1. Retrieves all users from the database
    2. Fetches available books, payment methods, and promotions
    3. Creates a random order for each user

    Args:
        psql_resource (PostgresResource): A resource for interacting with the PostgreSQL database.

    Returns:
        dg.MaterializeResult: A materialize result containing the number of orders created.
    """
    users = get_users(psql_resource)
    books = get_available_books()
    payments = get_available_payments()
    promotions = get_available_promotions()

    created_orders = []

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

    for user in users:
        logger.info(f"Processing order for user {user}")

        try:
            # Randomly select items
            num_items = random.randint(1, min(MAX_ITEMS_PER_ORDER, len(books)))
            selected_books = random.sample(books, num_items)

            # Create order items
            order_items = [
                {
                    "book_id": book["id"],
                    "quantity": random.randint(1, 3),
                }
                for book in selected_books
                if "id" in book  # Validate book has required fields
            ]

            # Randomly select payment
            payment = random.choice(payments)
            payment_id = payment.get("id")

            if not payment_id:
                logger.warning(f"Payment method missing ID, skipping order for {user}")
                continue

            # Randomly select promotion (if available)
            promotion_id = None
            if promotions:
                promotion = random.choice(promotions)
                promotion_id = promotion.get("id")

            order = create_mock_order(
                user=user,
                order_items=order_items,
                payment_id=payment_id,
                promotion_id=promotion_id,
            )

            created_orders.append(order)

        except Exception as e:
            raise dg.DagsterError(f"Failed to create order for user {user}: {e}")

    return dg.MaterializeResult(
        metadata={
            "number_of_orders": len(created_orders),
        }
    )
