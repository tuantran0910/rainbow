import dagster as dg
from dagster_dlt import DagsterDltResource
from dagster_dlt import dlt_assets
from dlt import pipeline

from shared.constants import DAGSTER_METADATA
from shared.constants import DAGSTER_TAGS
from shared.helpers import load_assets_configs
from shared.helpers import make_dlt_resources
from transformation.translators import CustomDagsterDltTranslator

logger = dg.get_dagster_logger(__name__)


def build_dlt_pipelines() -> dg.Definitions:
    """
    Build Dagster definitions containing DLT pipeline assets based on configuration.

    Returns:
        dg.Definitions: Dagster definitions containing the configured assets.

    Raises:
        ValueError: If required configuration is missing.
    """
    assets_configs = load_assets_configs(dagster_product="dlt")
    if not assets_configs:
        logger.info("No asset configurations found")
        return dg.Definitions()

    base_name = assets_configs.get("name")
    asset_group = assets_configs.get("group_name", "dlt")
    pipeline_progress_mode = assets_configs.get("progress", "log")

    # Construct Dagster's assets
    assets = []
    dlt_resources_config = assets_configs.get("resources", {})
    dlt_sources, destination_type = make_dlt_resources(dlt_resources_config=dlt_resources_config)
    for table_name, dlt_source in dlt_sources.items():
        asset_name = f"{base_name}__{table_name}"

        @dlt_assets(
            dlt_source=dlt_source,
            dlt_pipeline=pipeline(
                pipeline_name=asset_name,
                destination=destination_type,
                dataset_name=dlt_resources_config.get("destination", {}).get("dataset_name"),
                progress=pipeline_progress_mode,
            ),
            name=asset_name,
            group_name=asset_group,
            dagster_dlt_translator=CustomDagsterDltTranslator(),
        )
        def dagster_dlt_asset(context: dg.AssetExecutionContext, dlt: DagsterDltResource):
            """
            Transfer data from source to destination using dlt.

            Args:
                context (dg.AssetExecutionContext): The context for the asset execution.
                dlt (DagsterDltResource): The Dagster resource for dlt pipeline.
            """
            yield from dlt.run(context=context)

        assets.append(dagster_dlt_asset)

    # Construct Dagstet's schedules
    job_config = assets_configs.get("job")
    jobs = None
    schedule_config = assets_configs.get("schedule")
    schedules = None

    if job_config:
        metadata = {**DAGSTER_METADATA, **job_config.get("metadata", {})}
        tags = {**DAGSTER_TAGS, **job_config.get("tags", {})}
        job = dg.define_asset_job(
            name=job_config.get("name", f"{base_name}__{asset_group}__job"),
            description=job_config.get("description"),
            selection=dg.AssetSelection.groups(asset_group),
            metadata=metadata,
            tags=tags,
        )
        jobs = [job]

    if schedule_config and not job_config:
        logger.warning(
            "Schedule configuration provided without a job. Schedule will not be created."
        )
    if schedule_config and job_config:
        metadata = {**DAGSTER_METADATA, **schedule_config.get("metadata", {})}
        tags = {**DAGSTER_TAGS, **schedule_config.get("tags", {})}
        schedule = dg.ScheduleDefinition(
            name=schedule_config.get("name", f"{base_name}__{asset_group}__schedule"),
            job=job,
            description=schedule_config.get("description"),
            cron_schedule=schedule_config.get("cron_schedule"),
            execution_timezone=schedule_config.get("execution_timezone"),
            metadata=metadata,
            tags=tags,
        )
        schedules = [schedule]

    return dg.Definitions(
        assets=assets,
        resources={
            "dlt": DagsterDltResource(),
        },
        jobs=jobs,
        schedules=schedules,
    )
