{{
    config(
        materialized='view',
    )
}}

{{
    lambda_view(
        batch_relation=ref(this.name ~ '__batch'),
        stream_stg_relation=ref('stg_promotions__stream'),
        unique_key='id',
        updated_date_field='updated_at_tz_hcm',
        cutoff_days=1
    )
}}
