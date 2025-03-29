import logging

import dagster as dg
from dagster import AssetExecutionContext
from dagster_dlt import DagsterDltResource
from dagster_dlt import dlt_assets
from dlt import pipeline

from assets.helpers import load_assets_configs
from assets.helpers import make_dlt_resources
from assets.translators import CustomDagsterDltTranslator


logger = logging.getLogger(__name__)


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
        logger.error("No asset configurations found")
        raise ValueError("No asset configurations found")

    base_asset_name = assets_configs.get("name")
    asset_group = assets_configs.get("group_name", "dlt")
    pipeline_progress_mode = assets_configs.get("progress", "log")

    assets = []
    dlt_sources, destination_type = make_dlt_resources(
        dlt_resources_config=assets_configs.get("resources", {}),
    )
    for dlt_source in dlt_sources:
        asset_name = f"{base_asset_name}__{dlt_source.name}"
        logger.info(f"Creating asset: {asset_name}")

        @dlt_assets(
            dlt_source=dlt_source,
            dlt_pipeline=pipeline(
                pipeline_name=asset_name,
                destination=destination_type,
                progress=pipeline_progress_mode,
            ),
            name=asset_name,
            group_name=asset_group,
            dagster_dlt_translator=CustomDagsterDltTranslator(),
        )
        def dagster_dlt_asset(context: AssetExecutionContext, dlt: DagsterDltResource):
            """
            Transfer data from source to destination using dlt.
            """
            yield from dlt.run(context=context)

        assets.append(dagster_dlt_asset)

    return dg.Definitions(
        assets=assets,
    )
