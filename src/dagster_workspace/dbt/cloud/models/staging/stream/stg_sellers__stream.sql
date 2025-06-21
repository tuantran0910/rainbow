{{
    config(
        materialized='view',
    )
}}

SELECT *
FROM {{ ref('raw_sellers__datastream') }}
