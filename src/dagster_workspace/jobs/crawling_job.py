import dagster as dg

from constants import DAGSTER_CRAWLING_ASSET_GROUP
from constants import DAGSTER_CRAWLING_JOB_NAME
from constants import DAGSTER_METADATA
from constants import DAGSTER_TAGS

crawling_job = dg.define_asset_job(
    name=DAGSTER_CRAWLING_JOB_NAME,
    description="This job is responsible for crawling tiki resources",
    selection=dg.AssetSelection.groups(DAGSTER_CRAWLING_ASSET_GROUP),
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
)
