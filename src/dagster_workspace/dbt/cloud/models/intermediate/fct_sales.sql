{{
    config(
        materialized='view',
    )
}}

{{
    lambda_fact_sales(
        batch_fact_relation=ref(this.name ~ '__batch'),
        stream_orders_relation=ref('stg_orders__stream'),
        stream_order_items_relation=ref('stg_order_items__stream'),
        cutoff_days=1
    )
}}
