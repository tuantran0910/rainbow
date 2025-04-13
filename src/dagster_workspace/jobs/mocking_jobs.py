import dagster as dg

from assets.promotions import promotion_creations
from assets.transactions import order_transactions
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

orders_mocking_job = dg.define_asset_job(
    name="orders_mocking_job",
    description="This job is responsible for mocking orders",
    selection=dg.AssetSelection.assets(order_transactions),
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
)

promotions_mocking_job = dg.define_asset_job(
    name="promotions_mocking_job",
    description="This job is responsible for mocking promotions",
    selection=dg.AssetSelection.assets(promotion_creations),
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
)
