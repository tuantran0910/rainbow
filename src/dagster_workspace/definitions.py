import dagster as dg

from documentation import definitions as documentation_definitions
from ingestion import definitions as ingestion_definitions
from mocking import definitions as mocking_definitions
from transformation import definitions as transformation_definitions


main_defs = dg.Definitions.merge(
    ingestion_definitions.defs,
    transformation_definitions.defs,
    mocking_definitions.defs,
    documentation_definitions.defs,
)
