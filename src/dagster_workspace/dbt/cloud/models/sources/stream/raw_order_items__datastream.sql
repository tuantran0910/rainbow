{{
    config(
        materialized='ephemeral'
    )
}}

{{
    raw_model_sql(
        source=source('rainbow_datastream', 'rainbow_order_items'),
        layer='stream'
    )
}}
