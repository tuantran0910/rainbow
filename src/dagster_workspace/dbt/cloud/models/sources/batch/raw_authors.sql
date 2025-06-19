{{
    config(
        materialized='ephemeral'
    )
}}

{{
    raw_model_sql(
        source=source('rainbow', 'authors'),
        layer='batch'
    )
}}
