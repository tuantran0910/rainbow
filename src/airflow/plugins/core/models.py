from pathlib import Path
from typing import Any
from typing import Optional

import pendulum
from pydantic import BaseModel
from pydantic import field_serializer

from core.constants import DBT_PROJECT_DIR
from core.utils import get_environment


class DefaultArgs(BaseModel):
    """
    Defines the default arguments for Airflow DAGs.

    These arguments are applied to all tasks within a DAG unless overridden at the task level.

    Attributes:
        owner (str): The owner of the DAG, typically used for tracking.
        depends_on_past (bool): Whether each task instance depends on the success of the previous run.
        email_on_failure (bool): Whether to send an email when a task fails.
        email_on_retry (bool): Whether to send an email when a task retries.
        retries (int): The number of retry attempts in case of failure.
        retry_delay (int): The delay (in seconds) between retry attempts.
        pool (str, Optional): The Airflow pool to use for task instances.
    """

    owner: str = "rainbow"
    depends_on_past: bool = False
    email_on_failure: bool = False
    email_on_retry: bool = False
    retries: int = 3
    retry_delay: int = 300
    pool: Optional[str] = None


class DagParams(BaseModel):
    """
    Represents the configuration parameters for defining an Airflow DAG.

    This model ensures proper validation and consistency of DAG properties when dynamically creating DAGs.

    Attributes:
        dag_id (str): The unique identifier for the DAG.
        description (str, Optional): A brief description of the DAG's purpose.
        schedule_interval (str, Optional): The schedule for DAG execution (e.g., cron expression or preset like "@daily").
        start_date (str, Optional): The date and time from which the DAG starts running.
        timezone (str): The timezone in which the DAG operates (default: "UTC").
        catchup (bool): Whether past DAG runs should be scheduled if they were missed.
        max_active_runs (int): The maximum number of concurrently running instances of the DAG.
        default_args (Dict[str, Any], Optional): Default arguments applied to all tasks within the DAG.
        tags (List[str], Optional): Tags for categorizing and filtering DAGs in the Airflow UI.
    """

    dag_id: str
    description: Optional[str] = None
    schedule_interval: Optional[str] = None
    start_date: str = pendulum.now().subtract(months=1).to_date_string()
    timezone: str = "UTC"
    catchup: bool = False
    max_active_runs: int = 16
    default_args: Optional[dict[str, Any]] = None
    tags: Optional[list[str]] = None

    def model_dump(self, embed_timezone: bool = False, **kwargs: Any) -> dict[str, Any]:
        """
        Serialize the model to a dictionary for use in Airflow DAGs.

        Returns:
            Dict[str, Any]: The serialized model as a dictionary.
        """
        if embed_timezone:
            return {
                **super().model_dump(exclude={"start_date", "timezone"}, **kwargs),
                "start_date": pendulum.parse(self.start_date).in_timezone(self.timezone),
            }

        return super().model_dump(**kwargs)


class DbtParams(BaseModel):
    """
    Represents the configuration parameters for dbt DAGs.

    This model ensures proper validation and consistency of dbt properties when dynamically creating DAGs.

    Attributes:
        project_dir (str): The path to the dbt project root directory.
        profile_name (str): The name of the dbt profile to use.
        profile_target (str): The target profile within the dbt profile to use.
        select (List[str]): A list of dbt models to include in the DAG.
        exclude (List[str]): A list of dbt models to exclude from the DAG.
    """

    project_dir: str | Path = DBT_PROJECT_DIR
    profile_name: str = "default"
    profile_target: str = get_environment()
    select: list[str] = []
    exclude: list[str] = []

    @field_serializer("project_dir")
    def serialize_project_dir(self, value: str | Path) -> Path:
        """
        Serialize the project directory to a Path object.

        Args:
            value (str | Path): The project directory string.

        Returns:
            Path: The project directory as a Path object.
        """
        return Path(value)
