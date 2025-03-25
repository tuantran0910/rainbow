import logging

from dagster import asset
from dagster import AssetExecutionContext
from dagster import Definitions

logger = logging.getLogger(__name__)


@asset
def raw_data(context: AssetExecutionContext) -> list[dict]:
    """Example asset that represents raw data ingestion."""
    logger.info("Starting raw data ingestion")
    # TODO: Implement actual data ingestion logic
    return []


@asset
def processed_data(context: AssetExecutionContext, raw_data: list[dict]) -> list[dict]:
    """Example asset that processes raw data."""
    logger.info("Starting data processing")
    # TODO: Implement actual data processing logic
    return []


@asset
def transformed_data(context: AssetExecutionContext, processed_data: list[dict]) -> list[dict]:
    """Example asset that transforms processed data."""
    logger.info("Starting data transformation")
    # TODO: Implement actual data transformation logic
    return []


# Define the Dagster definitions
defs = Definitions(
    assets=[raw_data, processed_data, transformed_data],
)
