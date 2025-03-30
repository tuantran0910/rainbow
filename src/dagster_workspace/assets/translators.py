from dagster import AssetKey
from dagster import AssetSpec
from dagster_dlt import DagsterDltTranslator
from dagster_dlt.translator import DltResourceTranslatorData


class CustomDagsterDltTranslator(DagsterDltTranslator):
    def get_asset_spec(self, data: DltResourceTranslatorData) -> AssetSpec:
        """
        Overrides asset spec to:
            - Override asset key to be the dlt resource name.
            - Override deps to hide the default upstream assets.
        """
        default_spec = super().get_asset_spec(data)
        return default_spec.replace_attributes(
            key=AssetKey(f"{data.resource.name}"),
            deps=[],
        )
