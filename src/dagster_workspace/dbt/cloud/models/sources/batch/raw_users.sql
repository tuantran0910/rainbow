{{
    config(
        materialized='ephemeral'
    )
}}

{{
    raw_model_sql(
        source=source('rainbow', 'users'),
        layer='batch'
    )
}}
