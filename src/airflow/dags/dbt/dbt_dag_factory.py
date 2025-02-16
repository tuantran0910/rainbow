"""
This DAG is a factory for dbt DAGs. It is responsible for creating dbt DAGs dynamically
based on the dbt project directory, configurations for the dbt models, ...
"""  # noreorder

from typing import Any

from airflow.decorators import dag
from airflow.decorators import task
from airflow.models import Param

from core.factory import DbtDagFactory
from core.models import DagParams
from core.models import DefaultArgs


@dag(
    **DagParams(
        dag_id="dbt__dag_factory",
        schedule_interval="*/15 * * * *",
        default_args=DefaultArgs().model_dump(),
        tags=["dbt", "factory"],
    ).model_dump(embed_timezone=True),
    doc_md=__doc__,
    params={
        "full_refresh": Param(
            title="Full Refresh",
            default=False,
            type="boolean",
            description="Whether to perform a full refresh DAG generation from the Dbt Factory.",
        ),
    },
)
def dbt_factory_dag() -> None:
    """
    Factory DAG for dbt DAGs.
    """

    @task(task_id="dbt_dag_factory")
    def dbt_factory_task(params: dict[str, Any]) -> None:
        """
        Task to generate dbt DAGs dynamically.

        Args:
            params (dict[str, Any]): Task parameters.
        """
        factory = DbtDagFactory(full_refresh=params["full_refresh"])
        factory.build_dags()

    dbt_factory_task()


dbt_factory_dag()
