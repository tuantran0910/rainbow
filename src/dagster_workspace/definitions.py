import dagster as dg

from assets.factory import build_dlt_pipelines


defs = dg.Definitions.merge(
    build_dlt_pipelines(),
)
