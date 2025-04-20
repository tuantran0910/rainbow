{{
    config(
        materialized='ephemeral'
    )
}}

{{ stg_model_sql(source('rainbow', 'inventories')) }}
