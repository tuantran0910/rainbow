import random
import unicodedata
from enum import Enum
from typing import Any

import dagster as dg
from faker import Faker

from assets.helpers import make_http_request
from constants import API_BASE_URL
from constants import DAGSTER_METADATA
from constants import DAGSTER_MOCKING_ASSET_GROUP
from constants import DAGSTER_TAGS
from constants import DEFAULT_USER_PASSWORD
from constants import MAX_USERS_PER_REQUEST

logger = dg.get_dagster_logger()


class UserRole(str, Enum):
    USER = "USER"
    ADMIN = "ADMIN"


class UserGenerationOpConfig(dg.Config):
    num_users: int = MAX_USERS_PER_REQUEST


def generate_phone_number() -> str:
    """
    Generate a random phone number.

    Returns:
        str: A random phone number.
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


@dg.op(
    name="generate_fake_users",
    description="Generate a specified number of fake users with Faker.",
    out=dg.Out(description="List of generated user data"),
)
def generate_fake_users(config: UserGenerationOpConfig) -> list[dict[str, Any]]:
    """
    Generate a specified number of fake users with Faker.

    Args:
        config (UserGenerationOpConfig): A configuration for the number of users to generate.

    Returns:
        list[dict[str, Any]]: List of generated user data.
    """
    faker = Faker(locale="vi_VN")
    users = []

    for _ in range(config.num_users):
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
            random_number=random.randint(1, 999),
        )

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
        users.append(user)

    return users


@dg.op(
    name="register_users",
    description="Register a list of users via the Rainbow API.",
    ins={
        "user_data": dg.In(description="List of user data to register"),
    },
    out=dg.Out(description="List of registration results"),
)
def register_users(user_data: list[dict[str, Any]]) -> list[dict[str, Any]]:
    """
    Register a list of users via the Rainbow API.

    Args:
        user_data (list[dict[str, Any]]): List of user data to register.

    Returns:
        list[dict[str, Any]]: List of registration results.
    """
    register_url = f"{API_BASE_URL}/auth/register"
    registration_results = []

    for user in user_data:
        user_email = user["email"]

        try:
            logger.info(f"Registering user with email: {user_email}")
            registration_result = make_http_request(register_url, method="POST", data=user)
            if registration_result:
                registration_results.append(
                    {
                        "email": user_email,
                        "status_code": registration_result.get("status_code", 201),
                        "message": registration_result.get(
                            "message", "User registered successfully"
                        ),
                    }
                )
                logger.info(f"Successfully registered user: {user_email}")
            else:
                logger.error(
                    f"Failed to register user {user_email}: {registration_result.get('error')}"
                )
                registration_results.append(
                    {
                        "email": user_email,
                        "status_code": registration_result.get("status_code", 400),
                        "message": registration_result.get("message", "Failed to register user"),
                    }
                )
        except Exception as e:
            logger.error(f"Failed to register user {user_email}: {e}")
            registration_results.append(
                {
                    "email": user_email,
                    "status_code": 500,
                    "message": str(e),
                }
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
def analyze_registration_results(registration_results: list[dict[str, Any]]) -> dict[str, int]:
    """
    Analyze registration results and count successes and failures.

    Args:
        registration_results (list[dict[str, Any]]): List of registration results.

    Returns:
        dict[str, int]: Analysis of registration results.
    """
    successful = [r for r in registration_results if r["status_code"] == 201]
    failed = [r for r in registration_results if r["status_code"] != 201]
    return {
        "users_generated": len(registration_results),
        "successful_registrations": len(successful),
        "failed_registrations": len(failed),
    }


@dg.graph_asset(
    name="user_registrations",
    description="A graph asset that creates mock users and registers them via the API.",
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
    group_name=DAGSTER_MOCKING_ASSET_GROUP,
    kinds={"python"},
)
def user_registrations() -> dict[str, int]:
    """
    A graph asset that creates mock users and registers them via the API.

    This asset:
    1. Generates fake user data
    2. Registers new users via the API. If the user already exists, it will return a failed status.
    3. Returns the number of users registered and the number of users that failed to register.

    Args:
        config (UserRegistrationsConfig): A configuration for the number of users to register.

    Returns:
        dict[str, int]: A dictionary containing the number of users registered and the number of users that failed to register.
    """
    fake_users = generate_fake_users()
    registration_results = register_users(fake_users)
    analysis = analyze_registration_results(registration_results)
    return analysis
