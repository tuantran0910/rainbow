import json
import os
from typing import Any

import gcsfs
from cube import TemplateContext
from cube_dbt import Dbt

USE_REMOTE_MANIFEST = os.getenv("USE_REMOTE_MANIFEST", "true").lower() == "true"
DBT_GCS_PROJECT = os.getenv("DBT_GCS_PROJECT")
DBT_GCS_BUCKET = os.getenv("DBT_GCS_BUCKET")
MANIFEST_PATH = os.getenv("MANIFEST_PATH", "manifest.json")

template = TemplateContext()


def load_manifest_from_gcs(gs_path: str) -> dict[str, Any]:
    """
    Loads a manifest from a GCS path.
    """
    fs = gcsfs.GCSFileSystem(project=DBT_GCS_PROJECT)
    with fs.open(gs_path, "rb") as f:
        manifest = json.load(f)
        return manifest


def load_manifest_from_local(local_path: str) -> dict[str, Any]:
    """
    Loads a manifest from a local path.
    """
    with open(local_path, "rb") as f:
        manifest = json.load(f)
        return manifest


@template.function("dbt_model")
def dbt_model(name: str):
    """
    Initialize the Dbt object with the manifest loaded from either local or remote.
    If USE_REMOTE_MANIFEST is true, the manifest is loaded from the remote GCS path.
    Otherwise, the manifest is loaded from the local path.

    Args:
        name: The name of the model to load.

    Returns:
        The Dbt object.
    """
    if USE_REMOTE_MANIFEST:
        manifest = load_manifest_from_gcs(
            f"gs://{DBT_GCS_BUCKET}/{DBT_GCS_PROJECT}/{MANIFEST_PATH}"
        )
    else:
        manifest = load_manifest_from_local(MANIFEST_PATH)

    dbt = Dbt(manifest=manifest)
    return dbt.model(name)
