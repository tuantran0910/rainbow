{{
    config(
        materialized='ephemeral'
    )
}}

{{ raw_model_sql(source('rainbow', 'users')) }}
