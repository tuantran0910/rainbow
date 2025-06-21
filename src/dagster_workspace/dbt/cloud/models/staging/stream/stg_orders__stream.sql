{{
    config(
        materialized='view',
    )
}}

SELECT *
FROM {{ ref('raw_orders__datastream') }}
