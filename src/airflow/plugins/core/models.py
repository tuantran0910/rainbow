from datetime import timedelta
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

    Args:
        owner (str): The owner of the DAG, typically used for tracking and access control.
        depends_on_past (bool): Whether each task instance depends on the success of the previous run.
        email (str, optional): The email address for notifications.
        email_on_failure (bool): Whether to send an email when a task fails.
        email_on_retry (bool): Whether to send an email when a task retries.
        retries (int): The number of retry attempts in case of task failure.
        retry_delay (int): The delay (in seconds) between retry attempts.
        retry_exponential_backoff (bool): Whether to apply exponential backoff for retries.
        max_retry_delay (int, optional): The maximum delay (in seconds) for exponential backoff retries.
        sla (int | optional): The Service Level Agreement (SLA) deadline for task completion.
        execution_timeout (int | optional): The maximum runtime allowed for a task before it is forcibly marked as failed.
        queue (str | optional): The execution queue for task scheduling.
        priority_weight (int): The priority of the task when scheduling (higher values indicate higher priority).
        wait_for_downstream (bool): Whether to wait for downstream tasks to complete before marking this task as successful.
        trigger_rule (str): Defines how this task is triggered based on upstream task status (e.g., "all_success", "one_failed").
        pool (str | optional): The Airflow pool to use for task instances, which limits concurrency.
    """

    owner: str = "rainbow"
    depends_on_past: bool = False
    email: Optional[str] = None
    email_on_failure: bool = False
    email_on_retry: bool = False
    retries: int = 3
    retry_delay: int = 300
    retry_exponential_backoff: bool = False
    max_retry_delay: Optional[int] = None
    sla: Optional[int] = None
    execution_timeout: Optional[int] = None
    queue: Optional[str] = None
    priority_weight: int = 1
    wait_for_downstream: bool = False
    trigger_rule: str = "all_success"
    pool: Optional[str] = None

    @field_serializer("retry_delay", "sla", "execution_timeout")
    def serialize_timedelta(self, value: int) -> timedelta:
        """
        Serialize the integer value to a timedelta object.

        Args:
            value (int): The integer value to convert to a timedelta.

        Returns:
            timedelta: The retry delay as a timedelta object.
        """
        return timedelta(seconds=value)


class DagParams(BaseModel):
    """
    Defines the parameters for Airflow DAGs.

    These parameters are used to configure the DAG's metadata, schedule, and behavior.

    Args:
        dag_id (str): The unique identifier for the DAG.
        description (str, optional): A brief description of the DAG's purpose.
        schedule_interval (str, optional): The interval at which the DAG should run (e.g., "0 0 * * *").
        start_date (str): The start date for the DAG's schedule.
        end_date (str, optional): The end date for the DAG's schedule.
        timezone (str): The timezone to use for DAG scheduling and execution.
        catchup (bool): Whether to backfill historical DAG runs for the schedule interval.
        max_active_runs (int): The maximum number of active DAG runs allowed.
        default_args (dict[str, Any], optional): The default arguments for tasks within the DAG.
        tags (List[str], optional): A list of tags to categorize the DAG.
    """

    dag_id: str
    description: Optional[str] = None
    schedule_interval: Optional[str] = None
    start_date: str = pendulum.now().subtract(months=1).to_date_string()
    end_date: Optional[str] = None
    timezone: str = "UTC"
    catchup: bool = False
    max_active_runs: int = 16
    default_args: Optional[dict[str, Any]] = None
    tags: Optional[list[str]] = None

    def model_dump(self, embed_timezone: bool = False, **kwargs: Any) -> dict[str, Any]:
        """
        Serialize the model to a dictionary for use in Airflow DAGs.

        Args:
            embed_timezone (bool): If True, ensures `start_date` and `end_date` are timezone-aware.

        Returns:
            dict[str, Any]: The serialized model as a dictionary.
        """
        if embed_timezone:
            data = super().model_dump(exclude={"timezone"}, **kwargs)
            data["start_date"] = pendulum.parse(self.start_date).in_timezone(self.timezone)
            if self.end_date:
                data["end_date"] = pendulum.parse(self.end_date).in_timezone(self.timezone)
            return data

        return super().model_dump(**kwargs)


class DbtParams(BaseModel):
    """
    Represents the configuration parameters for dbt DAGs.

    This model ensures proper validation and consistency of dbt properties when dynamically creating DAGs.

    Args:
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
