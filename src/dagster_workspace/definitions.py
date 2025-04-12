import dagster as dg

import assets
from assets import build_dlt_pipelines
from dbt_project import build_dbt_project


main_defs = dg.Definitions.merge(
    build_dlt_pipelines(),
    build_dbt_project(),
    assets.defs,
)
