{% set surrogate_key_field_name = generate_surrogate_field_name(prefix='category') %}

{{
    config(
        materialized='incremental',
        unique_key=surrogate_key_field_name,
        incremental_strategy='delete+insert'
    )
}}

{{ dim_scd_type_1(
    stg_relation=ref('stg_categories'),
    unique_key='id',
    updated_at_field='updated_at_tz_hcm',
    lookback_in_days=1,
    surrogate_key_field_name=surrogate_key_field_name
) }}
