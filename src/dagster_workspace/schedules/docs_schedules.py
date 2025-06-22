import dagster as dg

from jobs.docs_jobs import dbt_docs_job


dbt_docs_daily_schedule = dg.ScheduleDefinition(
    name="dbt_docs_daily_schedule",
    job=dbt_docs_job,
    cron_schedule="0 2 * * *",
    description="Generate dbt documentation daily",
    tags={"team": "data_engineering", "schedule_type": "daily"},
)


@dg.asset_sensor(
    name="dbt_docs_after_models_sensor",
    asset_selection=dg.AssetSelection.groups("dbt_models"),
    job=dbt_docs_job,
    description="Generate dbt docs when dbt models are updated",
)
def dbt_docs_sensor(context: dg.SensorEvaluationContext, asset_event: dg.EventLogEntry):
    """
    Sensor that triggers dbt docs generation when dbt models are materialized.
    This ensures docs are always up-to-date with the latest model runs.
    """
    return dg.RunRequest(
        run_key=f"dbt_docs_{context.cursor}",
        tags={
            "triggered_by": "dbt_models_update",
            "model_run_id": asset_event.dagster_event.event_specific_data.materialization.asset_key,
        },
    )
