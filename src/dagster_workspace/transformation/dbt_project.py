import dagster as dg
from dagster_dbt import build_schedule_from_dbt_selection
from dagster_dbt import DagsterDbtTranslatorSettings
from dagster_dbt import dbt_assets
from dagster_dbt import DbtCliResource
from dagster_dbt import DbtProject

from shared.constants import DBT_PROFILES_DIR
from shared.constants import DBT_PROJECT_DIR
from shared.constants import DBT_TARGET_PROFILE
from transformation.translators import CustomDagsterDbtTranslator


def get_dbt_project() -> DbtProject:
    """
    Initialize and prepare the dbt project.

    Returns:
        DbtProject: The initialized dbt project.
    """
    dbt_project = DbtProject(
        project_dir=DBT_PROJECT_DIR,
        profiles_dir=DBT_PROFILES_DIR,
        target=DBT_TARGET_PROFILE,
    )

    # Prepare the project if in development mode
    # This should be safe now since we generate manifest during build
    dbt_project.prepare_if_dev()

    return dbt_project


def build_dbt_project() -> dg.Definitions:
    """
    Build the dbt project and return the Dagster definitions.

    Returns:
        dg.Definitions: Dagster definitions containing the configured assets.
    """

    # Initialize dbt project and resource inside the function
    dbt_project = get_dbt_project()
    dbt_resource = DbtCliResource(project_dir=dbt_project)

    # Create a dbt asset using the dbt project
    @dbt_assets(
        manifest=dbt_project.manifest_path,
        dagster_dbt_translator=CustomDagsterDbtTranslator(
            settings=DagsterDbtTranslatorSettings(
                enable_duplicate_source_asset_keys=True,
            ),
        ),
    )
    def dbt_models(context: dg.AssetExecutionContext, dbt: DbtCliResource):
        """
        This function is used to compile the dbt project and allow Dagster to build an asset graph.
        """
        yield from (
            dbt.cli(["run"], context=context).stream().fetch_row_counts().fetch_column_metadata()
        )

    # Create a schedule for the dbt models
    dbt_schedule = build_schedule_from_dbt_selection(
        [dbt_models],
        job_name="dbt_materialization_job",
        cron_schedule="0 2 * * *",
        dbt_select="fqn:*",
    )

    return dg.Definitions(
        assets=[dbt_models],
        resources={"dbt": dbt_resource},
        schedules=[dbt_schedule],
    )
