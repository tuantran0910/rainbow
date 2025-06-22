import dagster as dg

from ingestion.crawling_jobs import crawling_job
from shared.constants import DAGSTER_METADATA
from shared.constants import DAGSTER_TAGS

crawling_schedule = dg.ScheduleDefinition(
    name="tiki_crawling_schedule",
    description="This schedule is responsible for crawling tiki resources",
    cron_schedule="0 * * * *",
    job=crawling_job,
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
)
