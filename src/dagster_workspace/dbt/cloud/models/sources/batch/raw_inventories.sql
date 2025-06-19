{{
    config(
        materialized='ephemeral'
    )
}}

{{
    raw_model_sql(
        source=source('rainbow', 'inventories'),
        layer='batch'
    )
}}
