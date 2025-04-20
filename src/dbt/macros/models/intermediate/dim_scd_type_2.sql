{% macro dim_scd_type_2(stg_relation, unique_key, updated_at_field, except_columns_to_compare=[], lookback_in_days=1) %}

{# Get the columns from the staging model relation #}
{% set stg_columns = adapter.get_columns_in_relation(stg_relation) %}

{# Create a list of columns to compare, excluding the specified ones #}
{% set compare_columns = [] %}
{% for column in stg_columns %}
{% if column.name not in except_columns_to_compare and column.name != unique_key %}
{% do compare_columns.append(column.name) %}
{% endif %}
{% endfor %}

{% set model_params = get_common_params(lookback_in_days) %}

{% if is_incremental() %}
        WITH
            stg_data AS (
                SELECT *
                FROM {{ stg_relation }}
                WHERE {{ updated_at_field }} IN {{ model_params.incremental_dates_quoted_tz_hcm }}
            ),

            new_records AS (
                SELECT
                    {{ dbt_utils.generate_surrogate_key(['current.' ~ unique_key, 'current.' ~ updated_at_field]) }} AS surrogate_key,
                    current.* EXCEPT (valid_from, valid_to, is_current),
                    CURRENT_TIMESTAMP() AS valid_from,
                    toDateTime64('9999-12-31 23:59:59.999999', 6) AS valid_to,
                    TRUE AS is_current
                FROM stg_data AS current
                LEFT JOIN {{ this }} AS existing
                    ON
                        existing.{{ unique_key }} = current.{{ unique_key }}
                        AND existing.is_current = TRUE
                WHERE
                    existing.{{ unique_key }} IS NULL
                    {% for column in compare_columns %}
                    OR COALESCE(existing.{{ column }}, '') != COALESCE(current.{{ column }}, '')
                    {% endfor %}
            ),

            existing_records AS (
                SELECT
                    existing.* EXCEPT (valid_to, is_current),
                    CURRENT_TIMESTAMP() AS valid_to,
                    FALSE AS is_current
                FROM {{ this }} AS existing
                JOIN new_records AS incoming
                    ON
                        existing.{{ unique_key }} = incoming.{{ unique_key }}
                        AND existing.is_current = TRUE
            )

        SELECT *
        FROM existing_records

        UNION ALL

        SELECT *
        FROM new_records

    {% else %}
        -- Initial load (or full refresh)
        WITH stg_data AS (
            SELECT *
            FROM {{ stg_relation }}
        )

        SELECT
            {{ dbt_utils.generate_surrogate_key(['stg_data.' ~ unique_key, 'stg_data.' ~ updated_at_field]) }} AS surrogate_key,
            stg_data.* EXCEPT (valid_from, valid_to, is_current),
            CURRENT_TIMESTAMP() AS valid_from,
            toDateTime64('9999-12-31 23:59:59.999999', 6) AS valid_to,
            TRUE AS is_current
        FROM stg_data
    {% endif %}

{% endmacro %}
