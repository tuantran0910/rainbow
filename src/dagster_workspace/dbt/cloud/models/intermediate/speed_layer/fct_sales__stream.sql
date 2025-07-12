{% set surrogate_key_field_name = 'sale_key' %}

{{
    config(
        materialized='view',
    )
}}

{{
    fact_sales(
        unique_key=surrogate_key_field_name,
        orders_relation=ref('raw_orders__datastream'),
        order_items_relation=ref('raw_order_items__datastream'),
        dim_dates_relation=ref('dim_dates'),
        dim_users_relation=ref('dim_users'),
        dim_books_relation=ref('dim_books'),
        dim_sellers_relation=ref('dim_sellers'),
        dim_promotions_relation=ref('dim_promotions'),
        dim_payments_relation=ref('dim_payments'),
        dim_categories_relation=ref('dim_categories'),
        is_stream=true
    )
}}
