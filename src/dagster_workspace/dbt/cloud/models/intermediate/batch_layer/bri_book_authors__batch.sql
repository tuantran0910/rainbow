{% set surrogate_key_field_name = generate_surrogate_field_name(prefix='book_author') %}

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
            "book_id",
            "author_id"
        ]
    )
}}

{{
    dim_scd_type_1(
        stg_relation=ref('raw_book_authors'),
        unique_keys=['book_id', 'author_id'],
        updated_at_field='updated_at',
        lookback_in_days=1,
        surrogate_key_field_name=surrogate_key_field_name
    )
}}
