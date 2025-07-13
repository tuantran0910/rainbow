import dagster as dg

from documentation.docs_jobs import dbt_docs_job


dbt_docs_daily_schedule = dg.ScheduleDefinition(
    name="dbt_docs_daily_schedule",
    job=dbt_docs_job,
    cron_schedule="40 0,4,8,12,16,20 * * *",
    description="Generate dbt documentation daily",
    tags={"team": "data_engineering", "schedule_type": "daily"},
)
