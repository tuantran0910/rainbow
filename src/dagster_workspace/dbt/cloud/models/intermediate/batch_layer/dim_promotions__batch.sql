{% set surrogate_key_field_name = generate_surrogate_field_name(prefix='promotion') %}

{{
    config(
        materialized='incremental',
        unique_key=surrogate_key_field_name,
        incremental_strategy='merge',
        partition_by={
            "field": "created_date_tz_hcm",
            "data_type": "date",
            "granularity": "day"
        },
        cluster_by=[
            "updated_date_tz_hcm",
            "id"
        ]
    )
}}

{{
    dim_scd_type_1(
        stg_relation=ref('raw_promotions'),
        unique_keys=['id'],
        updated_at_field='updated_at',
        lookback_in_days=1,
        surrogate_key_field_name=surrogate_key_field_name
    )
}}
