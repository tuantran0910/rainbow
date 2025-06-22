import dagster as dg

from ingestion.crawler import TikiCrawler
from ingestion.crawling_jobs import crawling_job
from ingestion.crawling_schedules import crawling_schedule
from ingestion.factory import build_dlt_pipelines
from shared.constants import ADMIN_EMAIL
from shared.constants import ADMIN_PASSWORD
from shared.constants import DAGSTER_CRAWLING_ASSET_GROUP
from shared.constants import DAGSTER_METADATA
from shared.constants import DAGSTER_TAGS

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


# Build DLT pipeline definitions
dlt_definitions = build_dlt_pipelines()

# Combine all ingestion definitions
defs = dg.Definitions(
    assets=[
        tiki_resources_asset,
        *(dlt_definitions.assets or []),
    ],
    jobs=[
        crawling_job,
        *(dlt_definitions.jobs or []),
    ],
    schedules=[
        crawling_schedule,
        *(dlt_definitions.schedules or []),
    ],
    resources={
        **(dlt_definitions.resources or {}),
    },
)
