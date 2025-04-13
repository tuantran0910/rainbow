import dagster as dg

from assets.users import user_registrations
from constants import DAGSTER_METADATA
from constants import DAGSTER_TAGS

users_mocking_job = dg.define_asset_job(
    name="users_mocking_job",
    description="This job is responsible for mocking users",
    selection=dg.AssetSelection.assets(user_registrations),
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
)
