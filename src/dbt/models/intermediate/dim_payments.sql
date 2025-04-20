{{
    config(
        materialized='incremental',
        unique_key='surrogate_key',
        incremental_strategy='delete+insert'
    )
}}

{{ dim_scd_type_2(
    stg_relation=ref('stg_payments'),
    except_columns_to_compare=['updated_at', 'created_at', 'deleted_at', 'created_at_tz_hcm', 'updated_at_tz_hcm'],
    unique_key='id',
    updated_at_field='updated_at_tz_hcm',
    lookback_in_days=1,
) }}
