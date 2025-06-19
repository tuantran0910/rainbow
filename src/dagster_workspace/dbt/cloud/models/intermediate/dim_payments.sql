{{
    config(
        materialized='view',
    )
}}

{{
    lambda_view(
        batch_relation=model.alias ~ '__batch',
        stream_stg_relation=ref('stg_payments__stream'),
        unique_key='id',
        updated_date_field='updated_at_tz_hcm',
        cutoff_days=1
    )
}}
