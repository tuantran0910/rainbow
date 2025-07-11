{{
    config(
        materialized='view',
    )
}}

{{
    lambda_view(
        batch_relation=ref(this.name ~ '__batch'),
        stream_stg_relation=ref(this.name ~ '__stream'),
        unique_key='id',
        updated_date_field='updated_date_tz_hcm',
        cutoff_days=1,
        get_scd_type2_latest=true,
        batch_exclude_columns=['valid_from', 'valid_to', 'is_current']
    )
}}
