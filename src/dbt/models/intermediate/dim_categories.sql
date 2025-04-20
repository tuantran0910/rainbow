{% set model_params = get_common_params() %}

{{
    config(
        materialized='incremental',
        unique_key='surrogate_key',
        incremental_strategy='delete+insert'
    )
}}

WITH
    stg_categories AS (
        SELECT *
        FROM {{ ref('stg_categories') }}
        {% if is_incremental() %}
            WHERE updated_at_tz_hcm IN {{ model_params.incremental_dates_quoted_tz_hcm }}
        {% endif %}
    ),

    new_records AS (
        SELECT
            {{ dbt_utils.generate_surrogate_key(['current.id', 'current.updated_at']) }} AS surrogate_key,
            current.* EXCEPT (valid_from, valid_to, is_current),
            CURRENT_TIMESTAMP() AS valid_from,
            toDateTime64('9999-12-31 23:59:59.999999', 6) AS valid_to,
            TRUE AS is_current
        FROM stg_categories AS current
        LEFT JOIN {{ this }} AS existing
            ON
                existing.id = current.id
                AND existing.is_current = TRUE
        WHERE
            existing.id IS NULL
            OR existing.name != current.name
            OR existing.slug != current.slug
            OR COALESCE(existing.secondary_id, '') != COALESCE(current.secondary_id, '')
    ),

    existing_records AS (
        SELECT
            existing.* EXCEPT (valid_to, is_current),
            CURRENT_TIMESTAMP() AS valid_to,
            FALSE AS is_current
        FROM {{ this }} AS existing
        JOIN new_records AS incoming
            ON
                existing.id = incoming.id
                AND existing.is_current = TRUE
    )

{% if is_incremental() %}
    SELECT *
    FROM existing_records

    UNION ALL
{% endif %}

SELECT *
FROM new_records
