import dagster as dg

from assets.factory import build_dlt_pipelines
from dbt_project import build_dbt_project


defs = dg.Definitions.merge(
    build_dlt_pipelines(),
    build_dbt_project(),
)
