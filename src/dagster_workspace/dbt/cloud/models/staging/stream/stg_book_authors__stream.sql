{{
    config(
        materialized='view',
    )
}}

SELECT *
FROM {{ ref('raw_book_authors__datastream') }}
