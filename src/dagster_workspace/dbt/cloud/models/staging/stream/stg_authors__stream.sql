{{
    config(
        materialized='view',
    )
}}

SELECT *
FROM {{ ref('raw_authors__datastream') }}
