{{
    config(
        materialized='ephemeral'
    )
}}

SELECT
    * EXCEPT (deleted_at, datastream_metadata)
    REPLACE (
        CAST(secondary_id AS INT64) AS secondary_id
    ),
    DATE(created_at, 'Asia/Ho_Chi_Minh') AS created_date_tz_hcm,
    DATE(updated_at, 'Asia/Ho_Chi_Minh') AS updated_date_tz_hcm,
    STRUCT(
        'SPEED' AS layer,
        CURRENT_TIMESTAMP() AS dbt_loaded_at
    ) AS _dbt_metadata
FROM {{ source('rainbow__datastream', 'public_authors') }}
