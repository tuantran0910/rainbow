import dagster as dg

from assets import build_dlt_pipelines
from assets import defs
from dbt_project import build_dbt_project


main_defs = dg.Definitions.merge(
    build_dlt_pipelines(),
    build_dbt_project(),
    defs,
)
