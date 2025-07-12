{{
    config(
        materialized='ephemeral'
    )
}}

SELECT
    * EXCEPT (datastream_metadata),
    DATE(created_at, 'Asia/Ho_Chi_Minh') AS created_date_tz_hcm,
    DATE(updated_at, 'Asia/Ho_Chi_Minh') AS updated_date_tz_hcm,
    STRUCT(
        'SPEED' AS layer,
        CURRENT_TIMESTAMP() AS dbt_loaded_at
    ) AS _dbt_metadata
FROM {{ source('rainbow__datastream', 'public_book_authors') }}
