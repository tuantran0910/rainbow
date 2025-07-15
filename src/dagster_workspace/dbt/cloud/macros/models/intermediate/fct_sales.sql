{% macro fact_sales(
    unique_key,
    orders_relation,
    order_items_relation,
    dim_dates_relation,
    dim_users_relation,
    dim_books_relation,
    dim_sellers_relation,
    dim_promotions_relation,
    dim_payments_relation,
    dim_categories_relation
) %}
    {% set model_params = get_common_params(lookback_in_days=1) %}

    WITH
        stg_orders AS (
            SELECT *
            FROM {{ orders_relation }}
            {% if is_incremental() %}
                WHERE created_date_tz_hcm IN {{ model_params.incremental_dates_quoted_tz_hcm }}
            {% endif %}
        ),

        stg_order_items AS (
            SELECT *
            FROM {{ order_items_relation }}
            {% if is_incremental() %}
                WHERE created_date_tz_hcm IN {{ model_params.incremental_dates_quoted_tz_hcm }}
            {% endif %}
        ),

        dim_dates AS (
            SELECT *
            FROM {{ dim_dates_relation }}
        ),

        dim_users AS (
            SELECT *
            FROM {{ dim_users_relation }}
        ),

        dim_books AS (
            SELECT *
            FROM {{ dim_books_relation }}
        ),

        dim_sellers AS (
            SELECT *
            FROM {{ dim_sellers_relation }}
        ),

        dim_promotions AS (
            SELECT *
            FROM {{ dim_promotions_relation }}
        ),

        dim_payments AS (
            SELECT *
            FROM {{ dim_payments_relation }}
        ),

        dim_categories AS (
            SELECT *
            FROM {{ dim_categories_relation }}
        )

    SELECT
        {{ dbt_utils.generate_surrogate_key([
            'orders.id',
            'order_items.id'
        ]) }} AS {{ unique_key }},
        dates.date_key,
        users.user_key,
        books.book_key,
        sellers.seller_key,
        promotions.promotion_key,
        payments.payment_key,
        categories.category_key,
        orders.id AS order_id,
        order_items.quantity AS sales_quantity,
        order_items.unit_price AS regular_unit_price,
        CASE
            WHEN promotions.id IS NOT NULL AND promotions.discount_type = 'PERCENTAGE'
                THEN order_items.unit_price * promotions.discount_value / 100
            WHEN promotions.id IS NOT NULL AND promotions.discount_type = 'FIXED'
                THEN promotions.discount_value * 1000
            ELSE 0
        END AS discount_unit_price,
        CASE
            WHEN promotions.id IS NOT NULL AND promotions.discount_type = 'PERCENTAGE'
                THEN order_items.unit_price - order_items.unit_price * promotions.discount_value / 100
            WHEN promotions.id IS NOT NULL AND promotions.discount_type = 'FIXED'
                THEN order_items.unit_price - promotions.discount_value * 1000
            ELSE order_items.unit_price
        END AS net_unit_price,
        CASE
            WHEN promotions.id IS NOT NULL AND promotions.discount_type = 'PERCENTAGE'
                THEN order_items.quantity * (order_items.unit_price * promotions.discount_value / 100)
            WHEN promotions.id IS NOT NULL AND promotions.discount_type = 'FIXED'
                THEN order_items.quantity * (promotions.discount_value * 1000)
            ELSE order_items.unit_price
        END AS extended_discount_amount,
        order_items.quantity * order_items.unit_price AS extended_sales_amount,
        orders.created_date_tz_hcm
    FROM stg_orders AS orders
    JOIN stg_order_items AS order_items
        ON orders.id = order_items.order_id
    LEFT JOIN dim_dates AS dates
        ON orders.created_date_tz_hcm = dates.date
    LEFT JOIN dim_users AS users
        ON orders.user_id = users.id
    LEFT JOIN dim_books AS books
        ON order_items.book_id = books.id
    LEFT JOIN dim_sellers AS sellers
        ON books.seller_id = sellers.id
    LEFT JOIN dim_promotions AS promotions
        ON orders.promotion_id = promotions.id
    LEFT JOIN dim_payments AS payments
        ON orders.payment_id = payments.id
    LEFT JOIN dim_categories AS categories
        ON books.category_id = categories.id

{% endmacro %}
