{{
    config(
        materialized='ephemeral'
    )
}}

{{ raw_model_sql(source('rainbow_datastream', 'order_items'), 'stream') }} 