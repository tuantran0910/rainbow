import textwrap
from pathlib import Path
from typing import Any

import jinja2

from core.builder.base import BaseDagBuilder
from core.builder.exceptions import DagBuilderException
from core.builder.exceptions import DagBuilderTemplateException
from core.models import DbtParams


class DbtDagBuilder(BaseDagBuilder):
    """
    Builds Airflow DAGs for dbt.

    Args:
        airflow_type (str): Identifier for the DAG type (e.g., 'dbt', 'spark').
    """

    def __init__(self, airflow_type: str = "dbt", **kwargs: Any):
        super().__init__(airflow_type=airflow_type, **kwargs)

        self.dbt_params: dict[str, Any] = {}
        self._initialize_template_dbt_params()

    def _validate_dbt_root(self, dbt_project_path: Path) -> Path:
        """
        Validate dbt project root directory.

        Args:
            dbt_project_path (Path): Path to the dbt project root directory.

        Returns:
            Path: Validated dbt project root directory.

        Raises:
            DagBuilderException: If the directory does not exist.
        """
        if not dbt_project_path.exists():
            raise DagBuilderException(
                f"Dbt project directory not found: {dbt_project_path}. Verify the path exists."
            )

        return dbt_project_path

    def _initialize_template_dbt_params(self) -> None:
        """
        Initialize the Dbt Astronomer Cosmos DAG parameters from the configuration.

        Raises:
            DagBuilderException: If dbt configurations are not provided.
        """
        dbt_configs = self.dag_configs.get("dbt", {})
        if not dbt_configs:
            raise DagBuilderException("Dbt configurations are required.")

        self.dbt_params = DbtParams(**dbt_configs).model_dump()

    def _get_dbt_docs_url(self) -> str:
        """
        Return Dbt static documentation link of the project.

        Returns:
            str: The link to the Dbt static documentation.
        """
        # return f"https://dbt-docs.your-company.com/projects/{self.dbt_params.get('profile_name')}"
        return "https://docs.getdbt.com/"

    def _generate_dbt_doc(self) -> str:
        """
        Generate dbt-specific documentation section in Markdown format.

        Returns:
            str: The formatted DBT documentation section.
        """
        models_selected = [
            f"`{model}`" for model in self.dbt_params.get("select", [])
        ] or "All models"
        models_excluded = [f"`{model}`" for model in self.dbt_params.get("exclude", [])] or "None"

        dbt_doc = f"""
        #### Dbt Configuration
        - **Project Directory**: `{self.dbt_params.get("project_dir")}`
        - **Profile Name**: `{self.dbt_params.get("profile_name")}`
        - **Profile Target**: `{self.dbt_params.get("profile_target")}`
        - **Selected Models**: {models_selected}
        - **Excluded Models**: {models_excluded}
        - **DBT Documentation**: [View Docs]({self._get_dbt_docs_url()})
        """

        return textwrap.dedent(dbt_doc).strip()

    def _generate_doc_md(self):
        """
        Generate Markdown documentation for the DAG based on its configuration plus
        additional dbt-specific documentation.
        """
        return super()._generate_doc_md() + "\n\n" + self._generate_dbt_doc()

    def build_dag(self) -> Path:
        """
        Render the DAG template and write to the output path.

        Returns:
            Path: Path to the generated DAG file.

        Raises:
            DagBuilderTemplateException: If template rendering fails.
            DagBuilderException: If file operations fail.
        """
        try:
            print(self.dbt_params)
            print(self.dag_params)
            rendered_dag = self.template.render(
                **self.dag_params, **self.dbt_params, doc_md=self._generate_doc_md()
            )

            output_path = self._get_output_path()
            self.log.info(f"Generating Airflow DAG file at {output_path}")
            output_path.parent.mkdir(parents=True, exist_ok=True)
            output_path.write_text(rendered_dag)
            return output_path
        except jinja2.TemplateError as e:
            raise DagBuilderTemplateException(f"Template rendering failed: {str(e)}") from e
        except OSError as e:
            raise DagBuilderException(f"File operation failed: {str(e)}") from e
