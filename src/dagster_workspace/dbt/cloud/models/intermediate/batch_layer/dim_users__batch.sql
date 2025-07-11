{% set surrogate_key_field_name = generate_surrogate_field_name(prefix='user') %}

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
    dim_scd_type_2(
        stg_relation=ref('raw_users'),
        except_columns_to_compare=['updated_at', 'created_at', 'deleted_at', 'created_at_tz_hcm', 'updated_at_tz_hcm'],
        unique_key='id',
        updated_at_field='updated_at',
        lookback_in_days=1,
        surrogate_key_field_name=surrogate_key_field_name
    )
}}
