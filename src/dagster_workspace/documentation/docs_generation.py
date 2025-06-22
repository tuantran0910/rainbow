from pathlib import Path

import dagster as dg
from dagster import AssetExecutionContext
from dagster import AssetMaterialization
from dagster_dbt import DbtCliResource
from google.cloud import storage

from shared.constants import DAGSTER_METADATA
from shared.constants import DAGSTER_TAGS
from shared.constants import DBT_DOCS_BASE_URL
from shared.constants import DBT_DOCS_BUCKET_FOLDER
from shared.constants import DBT_DOCS_BUCKET_NAME
from shared.constants import DBT_PROFILES_DIR
from shared.constants import DBT_PROJECT_DIR


logger = dg.get_dagster_logger(__name__)


def upload_docs_to_gcs(docs_files: dict[str, Path]) -> dict[str, str]:
    """
    Upload dbt docs files to GCS bucket.

    Args:
        docs_files: Dictionary mapping destination filename to source Path

    Returns:
        Dictionary mapping local filename to GCS URI
    """
    logger.info(f"Initializing GCS client for bucket: {DBT_DOCS_BUCKET_NAME}")

    # Initialize GCS client
    storage_client = storage.Client()
    bucket = storage_client.bucket(DBT_DOCS_BUCKET_NAME)
    uploaded_files = {}
    for dest_filename, source_path in docs_files.items():
        try:
            gcs_path = f"{DBT_DOCS_BUCKET_FOLDER}/{dest_filename}"

            logger.info(f"Uploading {source_path} to {gcs_path}")

            # Create blob and upload
            blob = bucket.blob(gcs_path)

            # Set appropriate content type
            if dest_filename.endswith(".html"):
                blob.content_type = "text/html"
            elif dest_filename.endswith(".json"):
                blob.content_type = "application/json"

            # Upload the file
            blob.upload_from_filename(str(source_path))

            gcs_uri = f"gs://{DBT_DOCS_BUCKET_NAME}/{gcs_path}"
            uploaded_files[dest_filename] = gcs_uri

            logger.info(f"Successfully uploaded {dest_filename} to {gcs_uri}")

        except Exception as e:
            logger.error(f"Failed to upload {dest_filename}: {e}")
            raise dg.DagsterError(f"Failed to upload {dest_filename} to GCS: {e}")

    return uploaded_files


@dg.asset(
    name="dbt_docs_generation",
    description="Generate dbt documentation and upload to GCS bucket",
    metadata=DAGSTER_METADATA,
    tags=DAGSTER_TAGS,
    group_name="dbt_docs",
    kinds={"python", "dbt"},
    deps=["dbt_models"],
)
def dbt_docs_generation_asset(
    context: AssetExecutionContext, dbt: DbtCliResource
) -> AssetMaterialization:
    """
    Generate dbt documentation (manifest.json, static_index.html)
    and upload to GCS bucket for hosting.

    This asset runs `dbt docs generate` with proper flags and uploads the resulting files to the
    GCS bucket for static website hosting.

    Returns:
        AssetMaterialization with metadata about generated docs
    """

    logger.info("Starting dbt docs generation...")

    try:
        # Generate dbt docs using the dbt CLI resource with proper flags
        logger.info("Running dbt docs generate command with --static flag...")
        docs_invocation = dbt.cli(
            [
                "docs",
                "generate",
                "--project-dir",
                str(DBT_PROJECT_DIR),
                "--profiles-dir",
                str(DBT_PROFILES_DIR),
                "--static",
            ],
            context=context,
        )
        docs_invocation.wait()

        if not docs_invocation.is_successful():
            error = docs_invocation.get_error()
            raise dg.DagsterError(f"dbt docs generate failed: {error}")

        logger.info("dbt docs generate completed successfully")

        # Get the target path where dbt artifacts are generated
        target_path = docs_invocation.target_path
        docs_files: dict[str, Path] = {
            "manifest.json": target_path / "manifest.json",
            "static_index.html": target_path / "static_index.html",
        }

        # Verify all required files exist
        missing_files = []
        for filename, filepath in docs_files.items():
            if not filepath.exists():
                missing_files.append(filename)

        if missing_files:
            raise dg.DagsterError(f"Missing dbt docs files: {missing_files}")

        logger.info(f"Found all required docs files: {list(docs_files.keys())}")

        # Upload files to GCS bucket
        uploaded_files = upload_docs_to_gcs(docs_files=docs_files)

        logger.info(
            f"Successfully uploaded {len(uploaded_files)} files to GCS bucket: {DBT_DOCS_BUCKET_NAME}"
        )

        # Create metadata for the asset materialization
        metadata = {
            "bucket_name": DBT_DOCS_BUCKET_NAME,
            "uploaded_files": uploaded_files,
            "docs_url": DBT_DOCS_BASE_URL,
            "total_files": len(uploaded_files),
            "manifest_size_mb": round(docs_files["manifest.json"].stat().st_size / 1024 / 1024, 2),
            "static_index_size_mb": round(
                docs_files["static_index.html"].stat().st_size / 1024 / 1024, 2
            ),
        }

        logger.info("dbt docs generation and upload completed successfully")

        return AssetMaterialization(
            asset_key=context.asset_key,
            metadata=metadata,
            description="dbt documentation generated and uploaded to GCS",
        )

    except Exception as e:
        logger.exception(f"dbt docs generation failed: {e}")
        raise dg.DagsterError(f"dbt docs generation failed: {e}")
