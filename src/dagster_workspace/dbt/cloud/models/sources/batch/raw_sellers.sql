{{
    config(
        materialized='ephemeral'
    )
}}

SELECT
    * EXCEPT (_dlt_id, _dlt_load_id),
    DATE(created_at, 'Asia/Ho_Chi_Minh') AS created_date_tz_hcm,
    DATE(updated_at, 'Asia/Ho_Chi_Minh') AS updated_date_tz_hcm,
    STRUCT(
        'BATCH' AS layer,
        CURRENT_TIMESTAMP() AS dbt_loaded_at
    ) AS _dbt_metadata
FROM {{ source('rainbow', 'sellers') }}
