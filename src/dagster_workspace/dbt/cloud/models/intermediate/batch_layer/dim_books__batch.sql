{% set surrogate_key_field_name = generate_surrogate_field_name(prefix='book') %}

{{
    config(
        materialized='incremental',
        incremental_strategy='merge',
        unique_key=surrogate_key_field_name,
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
        stg_relation=ref('raw_books'),
        except_columns_to_compare=['updated_at', 'created_at', 'deleted_at', 'created_date_tz_hcm', 'updated_date_tz_hcm', '_dbt_metadata'],
        unique_key='id',
        updated_at_field='updated_at',
        lookback_in_days=1,
        surrogate_key_field_name=surrogate_key_field_name
    )
}}
