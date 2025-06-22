import dagster as dg

from mocking.mocking_jobs import orders_mocking_job
from mocking.mocking_jobs import promotions_mocking_job
from mocking.mocking_jobs import users_mocking_job
from mocking.mocking_schedules import orders_mocking_schedule
from mocking.mocking_schedules import promotions_mocking_schedule
from mocking.mocking_schedules import users_mocking_schedule
from mocking.promotions import promotion_creations
from mocking.transactions import order_transactions
from mocking.users import user_registrations
from shared.constants import DB_HOST
from shared.constants import DB_PORT
from shared.constants import RAINBOW_DB_NAME
from shared.constants import RAINBOW_DB_PASSWORD
from shared.constants import RAINBOW_DB_USER
from shared.resources.psql_resource import PostgresResource


# Combine all mocking definitions
defs = dg.Definitions(
    assets=[
        user_registrations,
        order_transactions,
        promotion_creations,
    ],
    jobs=[
        users_mocking_job,
        orders_mocking_job,
        promotions_mocking_job,
    ],
    schedules=[
        users_mocking_schedule,
        orders_mocking_schedule,
        promotions_mocking_schedule,
    ],
    resources={
        "rainbow_psql_resource": PostgresResource(
            host=DB_HOST,
            port=DB_PORT,
            database=RAINBOW_DB_NAME,
            user=RAINBOW_DB_USER,
            password=RAINBOW_DB_PASSWORD,
        ),
    },
)
