{% set surrogate_key_field_name = generate_surrogate_field_name(prefix='book_author') %}

{{
    config(
        materialized='view',
    )
}}

{{
    dim_scd_type_1(
        stg_relation=ref('raw_book_authors__datastream'),
        unique_keys=['book_id', 'author_id'],
        updated_at_field='updated_at',
        surrogate_key_field_name=surrogate_key_field_name
    )
}}
