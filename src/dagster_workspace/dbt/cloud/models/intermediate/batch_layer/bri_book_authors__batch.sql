{% set model_params = get_common_params(lookback_in_days=1) %}

{{
    config(
        materialized='incremental',
        unique_key='id',
        incremental_strategy='delete+insert'
    )
}}

WITH
    stg_book_authors AS (
        SELECT *
        FROM {{ ref('stg_book_authors') }}
        {% if is_incremental() %}
            WHERE updated_at_tz_hcm IN {{ model_params.incremental_dates_quoted_tz_hcm }}
        {% endif %}
    )

SELECT
    {{ dbt_utils.generate_surrogate_key(['book_id', 'author_id']) }} AS id,
    book_id,
    author_id
FROM stg_book_authors
