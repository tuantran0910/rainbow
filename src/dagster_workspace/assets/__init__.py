import dagster as dg

from assets.crawler import TikiCrawler
from assets.factory import build_dlt_pipelines
from assets.users import user_registrations
from constants import ADMIN_EMAIL
from constants import ADMIN_PASSWORD
from constants import DAGSTER_CRAWLING_ASSET_GROUP
from constants import DAGSTER_METADATA
from constants import DAGSTER_TAGS
from jobs.crawling_jobs import crawling_job
from jobs.mocking_jobs import users_mocking_job
from schedules.crawling_schedules import crawling_schedule
from schedules.mocking_schedules import users_mocking_schedule

__all__ = ["build_dlt_pipelines"]


logger = dg.get_dagster_logger(__name__)


@dg.asset(
    name="tiki_resources",
    description="Crawl tiki resources",
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
    group_name=DAGSTER_CRAWLING_ASSET_GROUP,
    kinds={"python"},
)
def tiki_resources_asset():
    """
    Dagster asset that handles the crawling of Tiki resources.

    This asset authenticates with Tiki using provided admin credentials,
    then executes the crawler to fetch and process resources.

    Returns:
        None

    Raises:
        Exception: If any step in the crawling process fails
    """
    logger.info("Starting Tiki Crawler...")
    try:
        crawler = TikiCrawler(admin_email=ADMIN_EMAIL, admin_password=ADMIN_PASSWORD)
        crawler.run()
        logger.info("Tiki Crawler completed successfully.")

    except Exception as e:
        logger.exception(f"Crawler execution failed: {e}")
        raise dg.DagsterError(f"Tiki crawler failed: {e}")


defs = dg.Definitions(
    assets=[tiki_resources_asset, user_registrations],
    jobs=[crawling_job, users_mocking_job],
    schedules=[crawling_schedule, users_mocking_schedule],
)
