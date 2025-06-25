from cube import config  # type: ignore[attr-defined]
from cube import file_repository  # type: ignore[attr-defined]


@config("repository_factory")
def repository_factory(ctx: dict) -> list[dict]:
    return file_repository("model")
