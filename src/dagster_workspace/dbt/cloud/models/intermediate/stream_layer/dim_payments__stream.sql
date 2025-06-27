{% set surrogate_key_field_name = generate_surrogate_field_name(prefix='payment') %}

{{
    config(
        materialized='view'
    )
}}

{{ dim_scd_type_1(
    stg_relation=ref('stg_payments__stream'),
    unique_key='id',
    updated_at_field='updated_at_tz_hcm',
    surrogate_key_field_name=surrogate_key_field_name,
    is_stream_mode=true
) }}
