{{
    config(
        materialized='view',
    )
}}

SELECT *
FROM {{ ref('raw_promotions__datastream') }}
