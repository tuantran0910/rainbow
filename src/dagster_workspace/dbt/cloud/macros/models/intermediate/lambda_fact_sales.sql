{% macro lambda_fact_sales(batch_fact_relation, stream_orders_relation, stream_order_items_relation, cutoff_days=1) %}

    {# Calculate cutoff timestamp for lambda view #}
    {% set cutoff_timestamp = calculate_cutoff_timestamp(cutoff_days=cutoff_days) %}

    WITH
        -- Batch layer facts (T-1): Pre-processed fact table data
        batch_facts AS (
            SELECT *
            FROM {{ batch_fact_relation }}
            WHERE created_date_tz_hcm < DATE({{ cutoff_timestamp }})
        ),

        -- Stream layer: Real-time processing of orders and order_items (T)
        stream_orders AS (
            SELECT *
            FROM {{ stream_orders_relation }}
            WHERE created_date_tz_hcm >= DATE({{ cutoff_timestamp }})
        ),

        stream_order_items AS (
            SELECT *
            FROM {{ stream_order_items_relation }}
            WHERE created_date_tz_hcm >= DATE({{ cutoff_timestamp }})
        ),

        -- Current dimensions for stream processing
        dim_dates AS (
            SELECT *
            FROM {{ ref('dim_dates') }}
        ),

        dim_users AS (
            SELECT *
            FROM {{ ref('dim_users') }}
            WHERE is_current = TRUE
        ),

        dim_books AS (
            SELECT *
            FROM {{ ref('dim_books') }}
            WHERE is_current = TRUE
        ),

        bri_book_authors AS (
            SELECT *
            FROM {{ ref('bri_book_authors') }}
        ),

        dim_sellers AS (
            SELECT *
            FROM {{ ref('dim_sellers') }}
        ),

        dim_promotions AS (
            SELECT *
            FROM {{ ref('dim_promotions') }}
        ),

        dim_payments AS (
            SELECT *
            FROM {{ ref('dim_payments') }}
        ),

        dim_categories AS (
            SELECT *
            FROM {{ ref('dim_categories') }}
        ),

        -- Stream facts: Process current data using same logic as batch
        stream_facts AS (
            SELECT
                {{ dbt_utils.generate_surrogate_key([
                    'dates.date_key',
                    'users.user_key',
                    'books.book_key',
                    'sellers.seller_key',
                    'promotions.promotion_key',
                    'payments.payment_key',
                    'categories.category_key'
                ]) }} AS sale_key,
                dates.date_key,
                users.user_key,
                books.book_key,
                sellers.seller_key,
                promotions.promotion_key,
                payments.payment_key,
                categories.category_key,
                orders.id AS order_id,
                order_items.quantity AS sales_quantity,
                order_items.price AS regular_unit_price,
                CASE
                    WHEN promotions.promotion_id IS NOT NULL AND promotions.discount_type = 'PERCENTAGE'
                        THEN order_items.price * promotions.discount_value / 100
                    WHEN promotions.promotion_id IS NOT NULL AND promotions.discount_type = 'FIXED'
                        THEN promotions.price * 1000
                    ELSE 0
                END AS discount_unit_price,
                CASE
                    WHEN promotions.promotion_id IS NOT NULL AND promotions.discount_type = 'PERCENTAGE'
                        THEN order_items.price - order_items.price * promotions.discount_value / 100
                    WHEN promotions.promotion_id IS NOT NULL AND promotions.discount_type = 'FIXED'
                        THEN order_items.price - promotions.price * 1000
                    ELSE order_items.price
                END AS net_unit_price,
                CASE
                    WHEN promotions.promotion_id IS NOT NULL AND promotions.discount_type = 'PERCENTAGE'
                        THEN order_items.quantity * (order_items.price * promotions.discount_value / 100)
                    WHEN promotions.promotion_id IS NOT NULL AND promotions.discount_type = 'FIXED'
                        THEN order_items.quantity * (promotions.price * 1000)
                    ELSE order_items.price
                END AS extended_discount_amount,
                order_items.quantity * order_items.price AS extended_sales_amount,
                orders.created_date_tz_hcm
            FROM stream_orders AS orders
            LEFT JOIN stream_order_items AS order_items
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
            LEFT JOIN bri_book_authors AS book_authors
                ON books.id = book_authors.book_id
        ),

        -- Unified facts: Stream data takes precedence, batch fills gaps
        unified_facts AS (
            SELECT *
            FROM stream_facts

            UNION ALL

            SELECT batch.*
            FROM batch_facts AS batch
            LEFT JOIN stream_facts AS stream
                ON batch.sale_key = stream.sale_key
            WHERE stream.sale_key IS NULL
        )

    SELECT *
    FROM unified_facts

{% endmacro %}
