from typing import Any

from core.builder import DbtDagBuilder
from core.factory.base import BaseDagFactory


class DbtDagFactory(BaseDagFactory):
    """
    Factory class for dbt Airflow DAGs.

    Args:
        airflow_type (str, optional): The type of Airflow DAG to build. Defaults to "dbt".
    """

    def __init__(self, airflow_type: str = "dbt", **kwargs: Any):
        super().__init__(airflow_type=airflow_type, **kwargs)

    def _get_builder(self) -> type[DbtDagBuilder]:
        """
        Get the builder class for dbt Airflow DAGs.

        Returns:
            type[DbtDagBuilder]: The builder class for dbt Airflow DAGs.
        """
        return DbtDagBuilder
