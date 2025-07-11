{% set surrogate_key_field_name = generate_surrogate_field_name(prefix='author') %}

{{
    config(
        materialized="incremental",
        incremental_strategy="merge",
        unique_key=surrogate_key_field_name,
        partition_by={
            "field": "created_date_tz_hcm",
            "data_type": "date",
            "granularity": "day"
        },
        cluster_by=[
            "updated_date_tz_hcm",
            "name",
            "slug",
            "id"
        ]
    )
}}

{{
    dim_scd_type_1(
        stg_relation=ref('raw_authors'),
        unique_keys=['id'],
        updated_at_field='updated_at',
        lookback_in_days=1,
        surrogate_key_field_name=surrogate_key_field_name
    )
}}
