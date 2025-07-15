import os
import subprocess

from cube import config  # type: ignore[attr-defined]
from cube import file_repository  # type: ignore[attr-defined]

GIT_REPO_PATH = "/tmp/git/rainbow"


@config("repository_factory")
def repository_factory(ctx: dict) -> list[dict]:
    return file_repository("model")


@config("schema_version")
def schema_version(ctx: dict) -> str:
    env = os.environ.get("ENVIRONMENT", "development")
    if env == "production":
        try:
            process = subprocess.run(
                ["git", "rev-parse", "HEAD"],
                cwd=GIT_REPO_PATH,
                capture_output=True,
                text=True,
                check=True,
                timeout=5,
            )
            return process.stdout.strip()
        except Exception:
            return "unknown"
    else:
        import uuid

        return str(uuid.uuid4())[:8]
