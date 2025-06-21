{{
    config(
        materialized='ephemeral'
    )
}}

{{
    raw_model_sql(
        source=source('rainbow_datastream', 'order_items'),
        layer='stream'
    )
}}
