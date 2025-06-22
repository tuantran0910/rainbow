import dagster as dg

from shared.constants import DAGSTER_CRAWLING_ASSET_GROUP
from shared.constants import DAGSTER_METADATA
from shared.constants import DAGSTER_TAGS

crawling_job = dg.define_asset_job(
    name="tiki_crawling_job",
    description="This job is responsible for crawling tiki resources",
    selection=dg.AssetSelection.groups(DAGSTER_CRAWLING_ASSET_GROUP),
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
)
