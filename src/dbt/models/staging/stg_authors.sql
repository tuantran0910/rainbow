{{
    config(
        materialized='ephemeral'
    )
}}

SELECT *
EXCEPT(_dlt_load_id, _dlt_id),
    toDate(created_at, 'Asia/Ho_Chi_Minh') AS created_at_tz_hcm,
    toDate(updated_at, 'Asia/Ho_Chi_Minh') AS updated_at_tz_hcm
FROM {{ source('rainbow', 'authors') }}
