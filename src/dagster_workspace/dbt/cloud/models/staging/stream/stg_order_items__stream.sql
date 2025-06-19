{{
    config(
        materialized='view',
    )
}}

SELECT *
FROM {{ ref('raw_order_items__stream') }}
