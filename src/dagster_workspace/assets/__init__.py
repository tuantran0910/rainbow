import logging

import dagster as dg

from assets.crawler import TikiCrawler
from assets.factory import build_dlt_pipelines
from constants import ADMIN_EMAIL
from constants import ADMIN_PASSWORD
from constants import DAGSTER_CRAWLING_ASSET_GROUP
from constants import DAGSTER_CRAWLING_ASSET_NAME
from constants import DAGSTER_METADATA
from constants import DAGSTER_TAGS
from jobs.crawling_job import crawling_job
from schedules.crawling_schedules import crawling_schedule

__all__ = ["build_dlt_pipelines"]


logger = logging.getLogger(__name__)


@dg.asset(
    name=DAGSTER_CRAWLING_ASSET_NAME,
    description="Crawl tiki resources",
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
    group_name=DAGSTER_CRAWLING_ASSET_GROUP,
)
def tiki_resources_asset():
    try:
        logger.info("Starting Tiki Crawler...")
        crawler = TikiCrawler(admin_email=ADMIN_EMAIL, admin_password=ADMIN_PASSWORD)

        # Try to get a token to verify authentication works
        token = crawler.auth_token_manager.get_token()
        if token:
            logger.info("✅ Authentication successful!")
            crawler.run()

    except Exception as e:
        logger.exception(f"Error in crawler main execution: {e}")


defs = dg.Definitions(
    assets=[tiki_resources_asset],
    jobs=[crawling_job],
    schedules=[crawling_schedule],
)
