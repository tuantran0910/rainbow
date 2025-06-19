{{
    config(
        materialized='view',
    )
}}

SELECT *
FROM {{ ref('raw_categories__stream') }}
