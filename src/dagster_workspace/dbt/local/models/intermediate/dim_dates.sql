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
        SELECT dr.start_date + number AS date_value
        FROM (
            SELECT
                toDate('{{ start_date }}') AS start_date,
                addYears(current_date(), {{ years_forward }}) AS end_date
        ) AS dr
        CROSS JOIN numbers(
            dateDiff('day', toDate('{{ start_date }}'), addYears(current_date(), {{ years_forward }}))
        )
    ),

    date_attributes AS (
        SELECT
            toYYYYMMDD(date_value) AS date_key,
            date_value AS date,
            toYear(date_value) AS year,
            toMonth(date_value) AS month,
            toDayOfMonth(date_value) AS day_of_month,
            toDayOfWeek(date_value) AS day_of_week,
            toDayOfYear(date_value) AS day_of_year,
            toQuarter(date_value) AS quarter,
            concat('Q', toString(toQuarter(date_value))) AS quarter_name,
            formatDateTime(date_value, '%M') AS month_name,
            concat(toString(toYear(date_value)), '-', toString(toMonth(date_value))) AS year_month,
            toISOWeek(date_value) AS week_of_year,
            toStartOfWeek(date_value) AS week_start_date,
            toStartOfWeek(date_value) + 6 AS week_end_date,
            CASE
                WHEN toMonth(date_value) >= 4 THEN toYear(date_value)
                ELSE toYear(date_value) - 1
            END AS fiscal_year,
            CASE
                WHEN toMonth(date_value) BETWEEN 4 AND 6 THEN 1
                WHEN toMonth(date_value) BETWEEN 7 AND 9 THEN 2
                WHEN toMonth(date_value) BETWEEN 10 AND 12 THEN 3
                ELSE 4
            END AS fiscal_quarter,
            CASE
                WHEN toMonth(date_value) >= 4 THEN toMonth(date_value) - 3
                ELSE toMonth(date_value) + 9
            END AS fiscal_month,
            CASE
                WHEN toMonth(date_value) = 1 AND toDayOfMonth(date_value) = 1 THEN TRUE
                WHEN toMonth(date_value) = 4 AND toDayOfMonth(date_value) = 30 THEN TRUE
                WHEN toMonth(date_value) = 5 AND toDayOfMonth(date_value) = 1 THEN TRUE
                WHEN toMonth(date_value) = 9 AND toDayOfMonth(date_value) = 2 THEN TRUE
                ELSE FALSE
            END AS is_holiday,
            coalesce(toDayOfWeek(date_value) IN (6, 7), FALSE) AS is_weekend,
            coalesce(toDayOfMonth(date_value) = 1, FALSE) AS is_first_day_of_month,
            coalesce(toDayOfMonth(date_value) = toLastDayOfMonth(date_value), FALSE) AS is_last_day_of_month,
            CASE
                WHEN toMonth(date_value) IN (1, 2, 3) THEN 'Spring'
                WHEN toMonth(date_value) IN (4, 5, 6) THEN 'Summer'
                WHEN toMonth(date_value) IN (7, 8, 9) THEN 'Autumn'
                ELSE 'Winter'
            END AS season

        FROM date_spine
    )

SELECT *
FROM date_attributes
WHERE date_value <= current_date()
