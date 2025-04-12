from collections.abc import Mapping
from collections.abc import Sequence
from typing import Any
from typing import Optional

from dagster import AssetKey
from dagster import AssetSpec
from dagster_dbt import DagsterDbtTranslator
from dagster_dlt import DagsterDltTranslator
from dagster_dlt.translator import DltResourceTranslatorData

from constants import DAGSTER_ASSETS_OWNER


class CustomDagsterDltTranslator(DagsterDltTranslator):
    def get_asset_spec(self, data: DltResourceTranslatorData) -> AssetSpec:
        """
        Overrides asset spec to override asset key to be the dlt resource name and
        override deps to hide the default upstream assets.
        """
        default_spec = super().get_asset_spec(data)
        return default_spec.replace_attributes(
            key=AssetKey(f"{data.resource.name}"),
            deps=[],
        )


class CustomDagsterDbtTranslator(DagsterDbtTranslator):
    @staticmethod
    def defaul_asset_key_fn(dbt_resource_props: Mapping[str, Any]) -> AssetKey:
        """
        Default asset key function to return the dbt resource name.

        Args:
            dbt_resource_props (Mapping[str, Any]): The dbt resource properties.

        Returns:
            AssetKey: The asset key for the dbt resource.
        """
        dbt_meta = dbt_resource_props.get("config", {}).get("meta", {}) or dbt_resource_props.get(
            "meta", {}
        )
        dagster_metadata = dbt_meta.get("dagster", {})
        asset_key_config = dagster_metadata.get("asset_key", [])
        if asset_key_config:
            return AssetKey(asset_key_config)

        if dbt_resource_props.get("version"):
            components = [dbt_resource_props["alias"]]
        else:
            components = [dbt_resource_props["name"]]

        return AssetKey(components)

    def get_asset_key(self, dbt_resource_props: Mapping[str, Any]) -> AssetKey:
        """
        Overrides the get_asset_key method to return a custom asset key.
        This is used to set the asset key in Dagster.
        """
        resource_type = dbt_resource_props["resource_type"]
        name = dbt_resource_props["name"]
        if resource_type == "source":
            return AssetKey(f"dlt_raw__{name}")
        else:
            return self.defaul_asset_key_fn(dbt_resource_props)

    def get_owners(self, dbt_resource_props: Mapping[str, Any]) -> Optional[Sequence[str]]:
        """
        Overrides the get_owners method to return a list of emails.
        This is used to set the owner of the asset in Dagster.
        """
        return [DAGSTER_ASSETS_OWNER]
