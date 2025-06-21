{{
    config(
        materialized='view',
    )
}}

SELECT *
FROM {{ ref('raw_users__datastream') }}
