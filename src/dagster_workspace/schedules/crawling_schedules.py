import dagster as dg

from constants import DAGSTER_CRAWLING_SCHEDULE_NAME
from constants import DAGSTER_METADATA
from constants import DAGSTER_TAGS
from jobs.crawling_job import crawling_job

crawling_schedule = dg.ScheduleDefinition(
    name=DAGSTER_CRAWLING_SCHEDULE_NAME,
    description="This schedule is responsible for crawling tiki resources",
    cron_schedule="0 * * * *",
    job=crawling_job,
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
)
