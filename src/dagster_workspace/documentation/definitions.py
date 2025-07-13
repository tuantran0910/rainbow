import dagster as dg

from documentation.docs_generation import dbt_docs_generation_asset
from documentation.docs_jobs import dbt_docs_job
from documentation.docs_schedules import dbt_docs_daily_schedule


# Combine all documentation definitions
defs = dg.Definitions(
    assets=[
        dbt_docs_generation_asset,
    ],
    jobs=[
        dbt_docs_job,
    ],
    schedules=[
        dbt_docs_daily_schedule,
    ],
)
