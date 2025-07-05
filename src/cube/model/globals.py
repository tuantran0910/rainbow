import json
import os
from typing import Any
from typing import Optional

from cube import TemplateContext
from cube_dbt import Dbt
from google.cloud import storage

# Configuration
USE_REMOTE_MANIFEST = os.getenv("USE_REMOTE_MANIFEST", "false").lower() == "true"
DBT_GCS_PROJECT = os.getenv("DBT_GCS_PROJECT")
DBT_GCS_BUCKET = os.getenv("DBT_GCS_BUCKET")
MANIFEST_PATH = os.getenv("MANIFEST_PATH", "manifest.json")

template = TemplateContext()
_manifest: Optional[dict[str, Any]] = None


def load_manifest_from_gcs(
    project: Optional[str], bucket_name: Optional[str], path: str
) -> dict[str, Any]:
    """
    Loads a manifest from a bucket.

    Args:
        project: The project ID.
        bucket: The bucket name.
        path: The path to the manifest.

    Returns:
        The manifest.
    """
    if not project:
        raise ValueError("Project must be provided")
    if not bucket_name:
        raise ValueError("Bucket name must be provided")

    client = storage.Client(project=project)
    bucket = client.bucket(bucket_name=bucket_name)
    blob = bucket.blob(blob_name=path)
    manifest = json.loads(blob.download_as_bytes())
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

    The manifest is loaded only once and cached in memory.

    Args:
        name: The name of the model to load.

    Returns:
        The Dbt object.
    """
    global _manifest

    if _manifest is None:
        if USE_REMOTE_MANIFEST:
            _manifest = load_manifest_from_gcs(
                project=DBT_GCS_PROJECT,
                bucket_name=DBT_GCS_BUCKET,
                path=MANIFEST_PATH,
            )
        else:
            _manifest = load_manifest_from_local(MANIFEST_PATH)

    dbt = Dbt(manifest=_manifest)
    return dbt.model(name)
