{{
    config(
        materialized='view',
    )
}}

{{
    lambda_fact_sales(
        batch_orders_relation=ref('stg_orders'),
        batch_order_items_relation=ref('stg_order_items'),
        stream_orders_relation=ref('stg_orders__stream'),
        stream_order_items_relation=ref('stg_order_items__stream'),
        cutoff_days=1
    )
}}
