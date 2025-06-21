{{
    config(
        materialized='ephemeral'
    )
}}

{{
    raw_model_sql(
        source=source('rainbow_datastream', 'inventories'),
        layer='stream'
    )
}}
