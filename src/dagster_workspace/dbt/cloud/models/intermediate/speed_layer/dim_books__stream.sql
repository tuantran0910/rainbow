{% set surrogate_key_field_name = generate_surrogate_field_name(prefix='book') %}

{{
    config(
        materialized='view',
    )
}}

{{
    dim_scd_type_1(
        stg_relation=ref('raw_books__datastream'),
        unique_keys=['id'],
        updated_at_field='updated_at',
        surrogate_key_field_name=surrogate_key_field_name
    )
}}
