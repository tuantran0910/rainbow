import os


def get_environment() -> str:
    """
    Get the current environment of Airflow (e.g. development, production)

    Returns:
        str: The current environment of Airflow.
    """
    return os.getenv("ENV", "development")
