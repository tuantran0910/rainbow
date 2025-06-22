import dagster as dg

from assets.docs_generation import dbt_docs_generation_asset
from constants import DAGSTER_METADATA
from constants import DAGSTER_TAGS


dbt_docs_job = dg.define_asset_job(
    name="dbt_docs_generation_job",
    description="This job is responsible for generating and uploading dbt documentation to GCS",
    selection=dg.AssetSelection.assets(dbt_docs_generation_asset),
    tags=DAGSTER_TAGS,
    metadata=DAGSTER_METADATA,
)
