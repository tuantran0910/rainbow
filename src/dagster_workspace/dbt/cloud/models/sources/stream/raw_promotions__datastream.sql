{{
    config(
        materialized='ephemeral'
    )
}}

{{
    raw_model_sql(
        source=source('rainbow_datastream', 'rainbow_promotions'),
        layer='stream'
    )
}}
