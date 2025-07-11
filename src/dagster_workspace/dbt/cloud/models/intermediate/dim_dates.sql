{{
    config(
        materialized='table',
        unique_key='date_key'
    )
}}

{% set start_date = var('dim_date_start_date', '2020-01-01') %}
{% set years_forward = var('dim_date_years_forward', 10) %}

WITH
    date_spine AS (
        SELECT
            DATE_ADD(DATE('{{ start_date }}'), INTERVAL day_num DAY) AS date_value
        FROM UNNEST(
            GENERATE_ARRAY(0, DATE_DIFF(DATE_ADD(CURRENT_DATE(), INTERVAL {{ years_forward }} YEAR), DATE('{{ start_date }}'), DAY))
        ) AS day_num
    ),

    date_attributes AS (
        SELECT
            FORMAT_DATE('%Y%m%d', date_value) AS date_key,
            date_value AS date,
            EXTRACT(YEAR FROM date_value) AS year,
            EXTRACT(MONTH FROM date_value) AS month,
            EXTRACT(DAY FROM date_value) AS day_of_month,
            EXTRACT(DAYOFWEEK FROM date_value) AS day_of_week, -- 1=Sunday
            EXTRACT(DAYOFYEAR FROM date_value) AS day_of_year,
            EXTRACT(QUARTER FROM date_value) AS quarter,
            CONCAT('Q', CAST(EXTRACT(QUARTER FROM date_value) AS STRING)) AS quarter_name,
            FORMAT_DATE('%B', date_value) AS month_name,
            FORMAT_DATE('%Y-%m', date_value) AS year_month,
            EXTRACT(ISOWEEK FROM date_value) AS week_of_year,
            DATE_TRUNC(date_value, WEEK(MONDAY)) AS week_start_date,
            DATE_ADD(DATE_TRUNC(date_value, WEEK(MONDAY)), INTERVAL 6 DAY) AS week_end_date,
            CASE
                WHEN EXTRACT(MONTH FROM date_value) >= 4 THEN EXTRACT(YEAR FROM date_value)
                ELSE EXTRACT(YEAR FROM date_value) - 1
            END AS fiscal_year,
            CASE
                WHEN EXTRACT(MONTH FROM date_value) BETWEEN 4 AND 6 THEN 1
                WHEN EXTRACT(MONTH FROM date_value) BETWEEN 7 AND 9 THEN 2
                WHEN EXTRACT(MONTH FROM date_value) BETWEEN 10 AND 12 THEN 3
                ELSE 4
            END AS fiscal_quarter,
            CASE
                WHEN EXTRACT(MONTH FROM date_value) >= 4 THEN EXTRACT(MONTH FROM date_value) - 3
                ELSE EXTRACT(MONTH FROM date_value) + 9
            END AS fiscal_month,
            CASE
                WHEN EXTRACT(MONTH FROM date_value) = 1 AND EXTRACT(DAY FROM date_value) = 1 THEN TRUE
                WHEN EXTRACT(MONTH FROM date_value) = 4 AND EXTRACT(DAY FROM date_value) = 30 THEN TRUE
                WHEN EXTRACT(MONTH FROM date_value) = 5 AND EXTRACT(DAY FROM date_value) = 1 THEN TRUE
                WHEN EXTRACT(MONTH FROM date_value) = 9 AND EXTRACT(DAY FROM date_value) = 2 THEN TRUE
                ELSE FALSE
            END AS is_holiday,
            IF(EXTRACT(DAYOFWEEK FROM date_value) IN (1, 7), TRUE, FALSE) AS is_weekend,
            IF(EXTRACT(DAY FROM date_value) = 1, TRUE, FALSE) AS is_first_day_of_month,
            IF(EXTRACT(DAY FROM date_value) = EXTRACT(DAY FROM LAST_DAY(date_value)), TRUE, FALSE) AS is_last_day_of_month,
            CASE
                WHEN EXTRACT(MONTH FROM date_value) IN (1, 2, 3) THEN 'Spring'
                WHEN EXTRACT(MONTH FROM date_value) IN (4, 5, 6) THEN 'Summer'
                WHEN EXTRACT(MONTH FROM date_value) IN (7, 8, 9) THEN 'Autumn'
                ELSE 'Winter'
            END AS season
        FROM date_spine
    ),

    final AS (
        SELECT *
        FROM date_attributes
        WHERE date <= CURRENT_DATE()
    )

SELECT *
FROM final
