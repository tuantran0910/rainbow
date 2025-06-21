{{
    config(
        materialized='view',
    )
}}

SELECT *
FROM {{ ref('raw_inventories__datastream') }}
