import dagster as dg

from mocking.mocking_jobs import orders_mocking_job
from mocking.mocking_jobs import promotions_mocking_job
from mocking.mocking_jobs import users_mocking_job
from shared.constants import DAGSTER_METADATA
from shared.constants import DAGSTER_TAGS

users_mocking_schedule = dg.ScheduleDefinition(
    name="users_mocking_schedule",
    description="This schedule is responsible for mocking users",
    cron_schedule="0 * * * *",
    job=users_mocking_job,
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
)

orders_mocking_schedule = dg.ScheduleDefinition(
    name="orders_mocking_schedule",
    description="This schedule is responsible for mocking orders",
    cron_schedule="*/5 * * * *",
    job=orders_mocking_job,
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
)

promotions_mocking_schedule = dg.ScheduleDefinition(
    name="promotions_mocking_schedule",
    description="This schedule is responsible for mocking promotions",
    cron_schedule="0 0,12 * * *",
    job=promotions_mocking_job,
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
)
