import random
import time
import unicodedata
from datetime import datetime
from enum import Enum
from typing import Any

import dagster as dg
from faker import Faker

from shared.constants import API_BASE_URL
from shared.constants import DAGSTER_METADATA
from shared.constants import DAGSTER_MOCKING_ASSET_GROUP
from shared.constants import DAGSTER_TAGS
from shared.constants import DEFAULT_USER_PASSWORD
from shared.constants import ENABLE_DATA_VALIDATION
from shared.constants import HTTP_MAX_RETRIES
from shared.constants import HTTP_REQUEST_TIMEOUT
from shared.constants import HTTP_RETRY_BACKOFF
from shared.constants import HTTP_RETRY_DELAY
from shared.constants import MAX_DUPLICATE_EMAIL_ATTEMPTS
from shared.constants import MAX_USERS_PER_REQUEST
from shared.constants import THROTTLE_DELAY_BETWEEN_REQUESTS
from shared.helpers import make_http_request

logger = dg.get_dagster_logger()


class UserRole(str, Enum):
    USER = "USER"
    ADMIN = "ADMIN"


class UserGenerationOpConfig(dg.Config):
    num_users: int = MAX_USERS_PER_REQUEST
    enable_validation: bool = ENABLE_DATA_VALIDATION
    max_duplicate_attempts: int = MAX_DUPLICATE_EMAIL_ATTEMPTS


def generate_phone_number() -> str:
    """
    Generate a random Vietnamese phone number with proper validation.

    Returns:
        str: A valid Vietnamese phone number.
    """
    phone_start = [
        "086",
        "096",
        "097",
        "098",
        "032",
        "033",
        "034",
        "035",
        "036",
        "037",
        "038",
        "039",
        "090",
        "093",
        "091",
        "094",
        "083",
        "084",
        "085",
    ]
    start = random.choice(phone_start)
    end = "".join([str(random.randint(0, 9)) for _ in range(7)])
    return f"{start}{end}"


def remove_accents(text: str) -> str:
    """
    Remove accents from a Vietnamese string.

    Args:
        text (str): The string to remove accents from.

    Returns:
        str: The string without accents.
    """
    return "".join(c for c in unicodedata.normalize("NFD", text) if unicodedata.category(c) != "Mn")


def validate_user_data(user: dict[str, Any]) -> bool:
    """
    Validate user data before submission.

    Args:
        user (dict[str, Any]): User data to validate

    Returns:
        bool: True if valid, False otherwise
    """
    if not ENABLE_DATA_VALIDATION:
        return True

    # Email validation
    email = user.get("email", "")
    if not email or "@" not in email or "." not in email:
        logger.warning(f"Invalid email format: {email}")
        return False

    # Phone number validation (Vietnamese format)
    phone = user.get("phone_number", "")
    if not phone or len(phone) != 10 or not phone.isdigit():
        logger.warning(f"Invalid phone number format: {phone}")
        return False

    # Name validation
    first_name = user.get("first_name", "")
    last_name = user.get("last_name", "")
    if not first_name or not last_name or len(first_name) < 2 or len(last_name) < 2:
        logger.warning(f"Invalid name: {first_name} {last_name}")
        return False

    return True


@dg.op(
    name="generate_fake_users",
    description="Generate a specified number of fake users with enhanced validation.",
    out=dg.Out(description="List of generated user data"),
)
def generate_fake_users(config: UserGenerationOpConfig) -> list[dict[str, Any]]:
    """
    Generate a specified number of fake users with Faker and enhanced validation.

    Args:
        config (UserGenerationOpConfig): A configuration for user generation.

    Returns:
        list[dict[str, Any]]: List of generated user data.
    """
    faker = Faker(locale="vi_VN")
    users: list[dict[str, str]] = []
    generated_emails = set()
    duplicate_attempts = 0

    while len(users) < config.num_users and duplicate_attempts < config.max_duplicate_attempts:
        # Randomly choose a gender to generate a fake name
        gender = random.choice(["male", "female"])
        if gender == "male":
            first_name = faker.first_name_male()
            last_name = faker.last_name_male() + " " + faker.middle_name()
        else:
            first_name = faker.first_name_female()
            last_name = faker.last_name_female() + " " + faker.middle_name()

        # Email is generated from the first name and last name
        first_name_without_accents = remove_accents(first_name).lower()
        last_name_without_accents = remove_accents(last_name).lower().replace(" ", "")
        email_template = "{first_name}.{last_name}{random_number}@gmail.com"
        email = email_template.format(
            first_name=first_name_without_accents,
            last_name=last_name_without_accents,
            random_number=random.randint(1, 9999),  # Increased range to reduce duplicates
        )

        # Check for duplicate emails in this batch
        if email in generated_emails:
            duplicate_attempts += 1
            continue

        phone_number = generate_phone_number()

        # Construct user data
        user = {
            "email": email,
            "password": DEFAULT_USER_PASSWORD,
            "first_name": first_name,
            "last_name": last_name,
            "phone_number": phone_number,
            "role": UserRole.USER,
        }

        # Validate user data
        if validate_user_data(user):
            users.append(user)
            generated_emails.add(email)
        else:
            duplicate_attempts += 1

    if duplicate_attempts >= config.max_duplicate_attempts:
        logger.warning(
            f"Reached max duplicate attempts ({config.max_duplicate_attempts}). Generated {len(users)} users."
        )

    logger.info(f"Generated {len(users)} valid users")
    return users


@dg.op(
    name="register_users",
    description="Register a list of users via the Rainbow API with enhanced error handling.",
    ins={
        "user_data": dg.In(description="List of user data to register"),
    },
    out=dg.Out(description="List of registration results"),
)
def register_users(user_data: list[dict[str, Any]]) -> list[dict[str, Any]]:
    """
    Register a list of users via the Rainbow API with enhanced error handling and rate limiting.

    Args:
        user_data (list[dict[str, Any]]): List of user data to register.

    Returns:
        list[dict[str, Any]]: List of registration results with detailed metrics.
    """
    register_url = f"{API_BASE_URL}/auth/register"
    registration_results = []
    success_count = 0
    error_count = 0
    duplicate_count = 0

    for i, user in enumerate(user_data):
        user_email = user["email"]

        try:
            # Add throttling between requests to be respectful to the API
            if i > 0 and THROTTLE_DELAY_BETWEEN_REQUESTS > 0:
                time.sleep(THROTTLE_DELAY_BETWEEN_REQUESTS)

            logger.info(f"Registering user {i + 1}/{len(user_data)}: {user_email}")

            registration_result = make_http_request(
                url=register_url,
                method="POST",
                data=user,
                timeout=HTTP_REQUEST_TIMEOUT,
                max_retries=HTTP_MAX_RETRIES,
                retry_delay=HTTP_RETRY_DELAY,
                retry_backoff=HTTP_RETRY_BACKOFF,
            )

            if registration_result:
                status_code = registration_result.get("status_code", 201)
                message = registration_result.get("message", "User registered successfully")

                if status_code == 201:
                    success_count += 1
                    logger.info(f"Successfully registered user: {user_email}")
                elif "already exists" in message.lower() or "email" in message.lower():
                    duplicate_count += 1
                    logger.warning(f"User already exists: {user_email}")
                else:
                    error_count += 1
                    logger.error(f"Failed to register user {user_email}: {message}")

                registration_results.append(
                    {
                        "email": user_email,
                        "status_code": status_code,
                        "message": message,
                        "timestamp": datetime.now().isoformat(),
                        "success": status_code == 201,
                    }
                )
            else:
                error_count += 1
                logger.error(f"No response received for user registration: {user_email}")
                registration_results.append(
                    {
                        "email": user_email,
                        "status_code": 500,
                        "message": "No response received",
                        "timestamp": datetime.now().isoformat(),
                        "success": False,
                    }
                )

        except Exception as e:
            error_count += 1
            error_message = str(e)
            logger.error(f"Failed to register user {user_email}: {error_message}")
            registration_results.append(
                {
                    "email": user_email,
                    "status_code": 500,
                    "message": error_message,
                    "timestamp": datetime.now().isoformat(),
                    "success": False,
                }
            )

    # Log summary statistics
    total_attempts = len(user_data)
    success_rate = (success_count / total_attempts * 100) if total_attempts > 0 else 0

    logger.info(
        f"Registration Summary: {success_count}/{total_attempts} successful ({success_rate:.1f}%), "
        f"{duplicate_count} duplicates, {error_count} errors"
    )

    return registration_results


@dg.op(
    name="analyze_registration_results",
    description="Analyze registration results and count successes and failures.",
    ins={
        "registration_results": dg.In(description="List of registration results"),
    },
    out=dg.Out(description="Analysis of registration results"),
)
def analyze_registration_results(registration_results: list[dict[str, Any]]) -> dict[str, Any]:
    """
    Analyze registration results and count successes and failures.

    Args:
        registration_results (list[dict[str, Any]]): List of registration results.

    Returns:
        dict[str, Any]: Analysis of registration results.
    """
    successful = [r for r in registration_results if r.get("success", False)]
    failed = [
        r
        for r in registration_results
        if not r.get("success", False) and "already exists" not in r.get("message", "").lower()
    ]
    duplicates = [
        r for r in registration_results if "already exists" in r.get("message", "").lower()
    ]

    total_attempts = len(registration_results)
    success_rate = (len(successful) / total_attempts * 100) if total_attempts > 0 else 0

    return {
        "users_generated": total_attempts,
        "successful_registrations": len(successful),
        "failed_registrations": len(failed),
        "duplicate_registrations": len(duplicates),
        "success_rate": success_rate,
    }


@dg.graph_asset(
    name="user_registrations",
    description="A graph asset that creates mock users and registers them via the API.",
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
    group_name=DAGSTER_MOCKING_ASSET_GROUP,
    kinds={"python"},
)
def user_registrations() -> dict[str, Any]:
    """
    A graph asset that creates mock users and registers them via the API.

    This asset:
    1. Generates fake user data with validation
    2. Registers new users via the API with rate limiting and retry logic
    3. Analyzes results and returns basic statistics
    4. Returns registration statistics

    Returns:
        dict[str, Any]: A dictionary containing registration statistics.
    """
    fake_users = generate_fake_users()
    registration_results = register_users(fake_users)
    analysis = analyze_registration_results(registration_results)
    return analysis
