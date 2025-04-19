import random
from datetime import datetime
from datetime import timedelta
from datetime import timezone
from enum import Enum
from typing import Any
from typing import Optional

import dagster as dg

from assets.helpers import AuthTokenManager
from assets.helpers import make_http_request
from constants import ADMIN_EMAIL
from constants import ADMIN_PASSWORD
from constants import API_BASE_URL
from constants import BOOK_CLUB_SALE_MAX_DAYS
from constants import BOOK_CLUB_SALE_MIN_DAYS
from constants import DAGSTER_METADATA
from constants import DAGSTER_MOCKING_ASSET_GROUP
from constants import DAGSTER_TAGS
from constants import FLASH_SALE_HOURS
from constants import FLASH_SALE_MAX_USES
from constants import FLASH_SALE_MIN_USES
from constants import MAX_FIXED_DISCOUNT_VND
from constants import MAX_PERCENTAGE_DISCOUNT
from constants import MAX_PROMOTION_USES
from constants import MAX_PROMOTIONS_PER_REQUEST
from constants import MIN_FIXED_DISCOUNT_VND
from constants import MIN_PERCENTAGE_DISCOUNT
from constants import MIN_PROMOTION_USES
from constants import REGULAR_SALE_MAX_DAYS
from constants import REGULAR_SALE_MIN_DAYS
from constants import SEASONAL_SALE_MAX_DAYS
from constants import SEASONAL_SALE_MIN_DAYS
from constants import WEEKEND_SALE_DAYS

logger = dg.get_dagster_logger()


class DiscountType(str, Enum):
    PERCENTAGE = "PERCENTAGE"
    FIXED = "FIXED"


class PromotionType(str, Enum):
    FLASH_SALE = "FLASH_SALE"
    WEEKEND_DEAL = "WEEKEND_DEAL"
    SEASONAL = "SEASONAL"
    BOOK_CLUB = "BOOK_CLUB"
    REGULAR = "REGULAR"


class PromotionGenerationOpConfig(dg.Config):
    num_promotions: int = MAX_PROMOTIONS_PER_REQUEST


@dg.op(
    name="initialize_auth_token_manager",
    description="Initialize the AuthTokenManager for use in other operations.",
    out=dg.Out(description="Initialized AuthTokenManager"),
)
def initialize_auth_token_manager() -> AuthTokenManager:
    """
    Initialize the AuthTokenManager for use in other operations.

    Returns:
        AuthTokenManager: An initialized AuthTokenManager instance.
    """
    auth_api_url = f"{API_BASE_URL}/auth/login"
    auth_token_manager = AuthTokenManager(
        email=ADMIN_EMAIL, password=ADMIN_PASSWORD, auth_api_url=auth_api_url
    )
    logger.info("Initialized AuthTokenManager")
    return auth_token_manager


def generate_promotion_name() -> tuple[str, PromotionType]:
    """
    Generate a random promotion name appropriate for a book marketplace.

    Returns:
        tuple[str, PromotionType]: A random promotion name and its type.
    """
    # Determine current month for seasonal promotions
    current_month = datetime.now().month
    seasonal_prefix = ""

    # Spring: January-March (1-3)
    if current_month in [1, 2, 3]:
        seasonal_prefix = "Spring"
    # Summer: April-June (4-6)
    elif current_month in [4, 5, 6]:
        seasonal_prefix = "Summer"
    # Fall: July-September (7-9)
    elif current_month in [7, 8, 9]:
        seasonal_prefix = "Fall"
    # Winter: October-December (10-12)
    elif current_month in [10, 11, 12]:
        seasonal_prefix = "Winter"

    # Determine type of promotion
    promo_type = random.choices(
        [
            PromotionType.FLASH_SALE,
            PromotionType.WEEKEND_DEAL,
            PromotionType.SEASONAL,
            PromotionType.BOOK_CLUB,
            PromotionType.REGULAR,
        ],
        weights=[0.2, 0.25, 0.2, 0.15, 0.2],
        k=1,
    )[0]

    # Generate name based on promotion type
    if promo_type == PromotionType.FLASH_SALE:
        names = ["Flash Sale", "Lightning Deal", "Quick Deal", "Hourly Special", "Rush Sale"]
        return random.choice(names), promo_type

    elif promo_type == PromotionType.WEEKEND_DEAL:
        prefixes = ["Weekend", "Friday", "Saturday", "Sunday", "Holiday"]
        suffixes = ["Special", "Deal", "Offer", "Sale"]
        return f"{random.choice(prefixes)} {random.choice(suffixes)}", promo_type

    elif promo_type == PromotionType.SEASONAL:
        suffixes = ["Collection", "Reading", "Special", "Event", "Celebration", "Picks"]
        return f"{seasonal_prefix} {random.choice(suffixes)}", promo_type

    elif promo_type == PromotionType.BOOK_CLUB:
        prefixes = ["Book Lovers", "Readers", "Book Club", "Literature", "Bookworm", "Knowledge"]
        suffixes = ["Membership", "Rewards", "Program", "Club", "Benefits", "Advantage"]
        return f"{random.choice(prefixes)} {random.choice(suffixes)}", promo_type
    else:
        prefixes = [
            "Bestseller",
            "Reading",
            "Author",
            "Paperback",
            "Novel",
            "Bookish",
            "Library",
            "Literature",
            "Readers",
        ]
        suffixes = ["Discount", "Deal", "Offer", "Promo", "Savings", "Event", "Promotion"]
        return f"{random.choice(prefixes)} {random.choice(suffixes)}", promo_type


@dg.op(
    name="get_active_promotions",
    description="Get all currently active promotions from the API.",
    ins={"auth_token_manager": dg.In(description="Initialized AuthTokenManager")},
    out=dg.Out(description="List of active promotions"),
)
def get_active_promotions(auth_token_manager: AuthTokenManager) -> list[dict[str, Any]]:
    """
    Get all currently active promotions from the API.

    Args:
        auth_token_manager (AuthTokenManager): The AuthTokenManager instance used for authentication.

    Returns:
        List[Dict[str, Any]]: List of active promotions.
    """
    promotions_url = f"{API_BASE_URL}/api/promotions"

    try:
        logger.info("Fetching active promotions")
        response_data = make_http_request(
            url=promotions_url, method="GET", use_auth=True, auth_token_manager=auth_token_manager
        )

        if response_data:
            promotions = response_data.get("data", {}).get("promotions", [])
            logger.info(f"Successfully fetched {len(promotions)} active promotions")
            return promotions
        else:
            logger.error(
                f"Failed to fetch active promotions: {response_data.get('message', 'Unknown error')}"
            )
            return []
    except Exception as e:
        raise dg.DagsterError(f"Failed to fetch active promotions: {e}")


def get_promotion_type_from_name(name: str) -> Optional[PromotionType]:
    """
    Determine the promotion type based on its name.

    Args:
        name (str): The promotion name.

    Returns:
        Optional[PromotionType]: The determined promotion type or None.
    """
    name_lower = name.lower()

    # Flash Sale indicators
    if any(term in name_lower for term in ["flash", "lightning", "quick", "rush", "hourly"]):
        return PromotionType.FLASH_SALE

    # Weekend Deal indicators
    if any(term in name_lower for term in ["weekend", "friday", "saturday", "sunday", "holiday"]):
        return PromotionType.WEEKEND_DEAL

    # Seasonal indicators
    if any(term in name_lower for term in ["spring", "summer", "fall", "winter"]):
        return PromotionType.SEASONAL

    # Book Club indicators
    if any(
        term in name_lower
        for term in ["book club", "membership", "program", "readers club", "book lovers"]
    ):
        return PromotionType.BOOK_CLUB

    return PromotionType.REGULAR


@dg.op(
    name="determine_available_promotion_types",
    description="Determine which promotion types are available for creation based on active promotions.",
    ins={"active_promotions": dg.In(description="List of active promotions")},
    out=dg.Out(description="Set of available promotion types"),
)
def determine_available_promotion_types(
    active_promotions: list[dict[str, Any]],
) -> set[PromotionType]:
    """
    Determine which promotion types are available for creation based on active promotions.

    Args:
        active_promotions (List[Dict[str, Any]]): List of active promotions.

    Returns:
        Set[PromotionType]: Set of promotion types that can be created.
    """
    # Get all promotion types
    all_promotion_types = {type_value for type_value in PromotionType}

    # Determine existing promotion types
    existing_types = set()
    for promotion in active_promotions:
        promo_type = get_promotion_type_from_name(promotion.get("name", ""))
        if promo_type:
            existing_types.add(promo_type)

    # Available types are those not currently active
    available_types = all_promotion_types - existing_types

    logger.info(f"Active promotion types: {existing_types}")
    logger.info(f"Available promotion types: {available_types}")

    return available_types


@dg.op(
    name="generate_promotions",
    description="Generate a specified number of promotions.",
    ins={"available_types": dg.In(description="Set of available promotion types")},
    out=dg.Out(description="List of generated promotion data"),
)
def generate_promotions(
    config: PromotionGenerationOpConfig, available_types: set[PromotionType]
) -> list[dict[str, Any]]:
    """
    Generate a specified number of promotions from available promotion types.

    Args:
        config (PromotionGenerationOpConfig): A configuration for the number of promotions to generate.
        available_types (Set[PromotionType]): Set of promotion types that can be created.

    Returns:
        list[dict[str, Any]]: List of generated promotion data.
    """
    promotions: list[dict[str, Any]] = []

    if not available_types:
        logger.info("No promotion types available for creation")
        return promotions

    num_promotions = min(
        random.randint(1, min(config.num_promotions, MAX_PROMOTIONS_PER_REQUEST)),
        len(available_types),
    )

    promotion_types_to_create = random.sample(list(available_types), num_promotions)

    for promo_type in promotion_types_to_create:
        # Generate promotion name based on type
        if promo_type == PromotionType.FLASH_SALE:
            names = ["Flash Sale", "Lightning Deal", "Quick Deal", "Hourly Special", "Rush Sale"]
            promotion_name = random.choice(names)

        elif promo_type == PromotionType.WEEKEND_DEAL:
            prefixes = ["Weekend", "Friday", "Saturday", "Sunday", "Holiday"]
            suffixes = ["Special", "Deal", "Offer", "Sale"]
            promotion_name = f"{random.choice(prefixes)} {random.choice(suffixes)}"

        elif promo_type == PromotionType.SEASONAL:
            # Determine current season
            current_month = datetime.now().month
            seasonal_prefix = ""
            if 3 <= current_month <= 5:
                seasonal_prefix = "Spring"
            elif 6 <= current_month <= 8:
                seasonal_prefix = "Summer"
            elif 9 <= current_month <= 11:
                seasonal_prefix = "Fall"
            else:
                seasonal_prefix = "Winter"

            suffixes = ["Collection", "Reading", "Special", "Event", "Celebration", "Picks"]
            promotion_name = f"{seasonal_prefix} {random.choice(suffixes)}"

        elif promo_type == PromotionType.BOOK_CLUB:
            prefixes = [
                "Book Lovers",
                "Readers",
                "Book Club",
                "Literature",
                "Bookworm",
                "Knowledge",
            ]
            suffixes = ["Membership", "Rewards", "Program", "Club", "Benefits", "Advantage"]
            promotion_name = f"{random.choice(prefixes)} {random.choice(suffixes)}"

        else:
            prefixes = [
                "Bestseller",
                "Reading",
                "Author",
                "Paperback",
                "Novel",
                "Bookish",
                "Library",
                "Literature",
                "Readers",
            ]
            suffixes = ["Discount", "Deal", "Offer", "Promo", "Savings", "Event", "Promotion"]
            promotion_name = f"{random.choice(prefixes)} {random.choice(suffixes)}"

        # Randomly choose discount type
        discount_type = random.choice([DiscountType.PERCENTAGE, DiscountType.FIXED])

        # Set discount value based on type and promotion type
        if discount_type == DiscountType.PERCENTAGE:
            # Flash sales and weekend deals often have higher discounts
            if promo_type in [PromotionType.FLASH_SALE, PromotionType.WEEKEND_DEAL]:
                min_discount = max(20.0, MIN_PERCENTAGE_DISCOUNT)
                discount_value = int(random.uniform(min_discount, MAX_PERCENTAGE_DISCOUNT))
            else:
                discount_value = int(
                    random.uniform(MIN_PERCENTAGE_DISCOUNT, MAX_PERCENTAGE_DISCOUNT)
                )
        else:
            discount_value = int(random.uniform(MIN_FIXED_DISCOUNT_VND, MAX_FIXED_DISCOUNT_VND))

        # Set start date within the next 6 hours
        now = datetime.now()
        start_date = now + timedelta(minutes=random.randint(5, 360))

        # Determine end date based on promotion type
        if promo_type == PromotionType.FLASH_SALE:
            # Flash sales last a few hours
            end_date = start_date + timedelta(hours=FLASH_SALE_HOURS)

        elif promo_type == PromotionType.WEEKEND_DEAL:
            # Weekend deals typically last 2-3 days
            end_date = start_date + timedelta(days=WEEKEND_SALE_DAYS)

        elif promo_type == PromotionType.SEASONAL:
            # Seasonal promotions last longer
            duration_days = random.randint(SEASONAL_SALE_MIN_DAYS, SEASONAL_SALE_MAX_DAYS)
            end_date = start_date + timedelta(days=duration_days)

        elif promo_type == PromotionType.BOOK_CLUB:
            # Book club promotions are long-term
            duration_days = random.randint(BOOK_CLUB_SALE_MIN_DAYS, BOOK_CLUB_SALE_MAX_DAYS)
            end_date = start_date + timedelta(days=duration_days)

        else:
            duration_days = random.randint(REGULAR_SALE_MIN_DAYS, REGULAR_SALE_MAX_DAYS)
            end_date = start_date + timedelta(days=duration_days)

        # Set max uses based on promotion type
        if promo_type == PromotionType.FLASH_SALE:
            # Flash sales have limited uses
            max_uses = random.randint(FLASH_SALE_MIN_USES, FLASH_SALE_MAX_USES)
        elif promo_type == PromotionType.WEEKEND_DEAL:
            # Weekend deals have moderate uses
            max_uses = random.randint(MIN_PROMOTION_USES * 2, MAX_PROMOTION_USES // 2)
        elif promo_type == PromotionType.BOOK_CLUB:
            # Book club promotions have high number of uses
            max_uses = random.randint(MAX_PROMOTION_USES // 2, MAX_PROMOTION_USES)
        else:
            # Regular and seasonal promotions
            max_uses = random.randint(MIN_PROMOTION_USES, MAX_PROMOTION_USES)

        # Format dates in RFC3339 format with timezone (Z notation for UTC)
        # This matches the Go time format "2006-01-02T15:04:05Z07:00"
        start_date_utc = start_date.replace(tzinfo=timezone.utc)
        end_date_utc = end_date.replace(tzinfo=timezone.utc)

        # Construct promotion data
        promotion = {
            "name": promotion_name,
            "discount_type": discount_type,
            "discount_value": discount_value,
            "start_date": start_date_utc.strftime("%Y-%m-%dT%H:%M:%SZ"),
            "end_date": end_date_utc.strftime("%Y-%m-%dT%H:%M:%SZ"),
            "max_uses": max_uses,
        }
        promotions.append(promotion)

    return promotions


@dg.op(
    name="create_promotions",
    description="Create promotions via the Rainbow API.",
    ins={
        "promotion_data": dg.In(description="List of promotion data to create"),
        "auth_token_manager": dg.In(description="Initialized AuthTokenManager"),
    },
    out=dg.Out(description="List of creation results"),
)
def create_promotions(
    promotion_data: list[dict[str, Any]], auth_token_manager: AuthTokenManager
) -> list[dict[str, Any]]:
    """
    Create promotions via the Rainbow API.

    Args:
        promotion_data (list[dict[str, Any]]): List of promotion data to create.
        auth_token_manager (AuthTokenManager): The initialized AuthTokenManager.

    Returns:
        list[dict[str, Any]]: List of creation results.
    """
    create_url = f"{API_BASE_URL}/api/promotions"
    creation_results = []

    for promotion in promotion_data:
        promotion_name = promotion["name"]

        try:
            logger.info(f"Creating promotion: {promotion_name}")
            promotion_response = make_http_request(
                url=create_url,
                method="POST",
                data=promotion,
                use_auth=True,
                auth_token_manager=auth_token_manager,
            )
            if promotion_response:
                creation_results.append(
                    {
                        "name": promotion_name,
                        "status_code": promotion_response.get("status_code", 201),
                        "message": promotion_response.get(
                            "message", "Promotion created successfully"
                        ),
                    }
                )
                logger.info(f"Successfully created promotion: {promotion_name}")
            else:
                logger.error(
                    f"Failed to create promotion {promotion_name}: {promotion_response.get('error') if promotion_response else 'Unknown error'}"
                )
                creation_results.append(
                    {
                        "name": promotion_name,
                        "status_code": 400,
                        "message": "Failed to create promotion",
                    }
                )
        except Exception as e:
            raise dg.DagsterError(f"Failed to create promotion: {e}")

    return creation_results


@dg.graph_asset(
    name="promotion_creations",
    description="A graph asset that creates mock promotions and creates them via the API.",
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
    group_name=DAGSTER_MOCKING_ASSET_GROUP,
    kinds={"python"},
)
def promotion_creations() -> list[dict[str, Any]]:
    """
    A graph asset that creates mock promotions and creates them via the API.

    This asset:
    1. Initializes the AuthTokenManager
    2. Fetches currently active promotions
    3. Determines which promotion types are available for creation
    4. Generates promotion data for available types
    5. Creates new promotions via the API

    Returns:
        list[dict[str, Any]]: The creation results containing information about successful or failed promotions.
    """
    auth_token_manager = initialize_auth_token_manager()
    active_promotions = get_active_promotions(auth_token_manager)
    available_types = determine_available_promotion_types(active_promotions)
    promotions = generate_promotions(available_types)
    creation_results = create_promotions(promotions, auth_token_manager)
    return creation_results
