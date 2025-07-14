{% macro dim_scd_type_2(stg_relation, unique_key, surrogate_key_field_name, updated_at_field=none, except_columns_to_compare=[], lookback_in_days=1) %}

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
            -- First, deduplicate staging data to get the latest record for each unique_key
            stg_data AS (
                SELECT *
                FROM {{ stg_relation }}
                WHERE DATE({{ updated_at_field }}) IN {{ model_params.incremental_dates_quoted_tz_hcm }}
                QUALIFY ROW_NUMBER() OVER (
                    PARTITION BY {{ unique_key }}
                    ORDER BY {{ updated_at_field }} DESC
                ) = 1
            ),

            -- Identify records that are new or have changes
            records_to_process AS (
                SELECT
                    stg_data.*,
                    existing.{{ surrogate_key_field_name }} AS existing_surrogate_key,
                    CASE
                        WHEN existing.{{ unique_key }} IS NULL THEN 'NEW'
                        {% if compare_columns %}
                        WHEN
                        {% for column in compare_columns %}
                            COALESCE(existing.{{ column }}, '') != COALESCE(stg_data.{{ column }}, ''){% if not loop.last %} OR {% endif %}
                        {% endfor %} THEN 'CHANGED'
                        {% endif %}
                        ELSE 'UNCHANGED'
                    END AS record_status
                FROM stg_data
                LEFT JOIN {{ this }} AS existing
                    ON existing.{{ unique_key }} = stg_data.{{ unique_key }}
                    AND existing.is_current = TRUE
            ),

            -- Create new records for NEW and CHANGED records
            new_records AS (
                SELECT
                    {{ dbt_utils.generate_surrogate_key(['rtp.' ~ unique_key, 'rtp.' ~ updated_at_field]) }} AS {{ surrogate_key_field_name }},
                    rtp.* EXCEPT (existing_surrogate_key, record_status),
                    CURRENT_TIMESTAMP() AS valid_from,
                    TIMESTAMP('9999-12-31 23:59:59.999999') AS valid_to,
                    TRUE AS is_current
                FROM records_to_process AS rtp
                WHERE record_status IN ('NEW'{% if compare_columns %}, 'CHANGED'{% endif %})
            ),

            -- Update existing records to set is_current = FALSE for CHANGED records
            existing_records_to_expire AS (
                SELECT
                    existing.* EXCEPT (valid_to, is_current),
                    CURRENT_TIMESTAMP() AS valid_to,
                    FALSE AS is_current
                FROM {{ this }} AS existing
                INNER JOIN records_to_process AS rtp
                    ON existing.{{ unique_key }} = rtp.{{ unique_key }}
                    AND existing.is_current = TRUE
                    {% if compare_columns %}
                    AND rtp.record_status = 'CHANGED'
                    {% else %}
                    AND FALSE  -- No columns to compare, so no records will be changed
                    {% endif %}
            ),

            -- Keep all existing records that are not being processed (touched) at all
            unchanged_existing_records AS (
                SELECT existing.*
                FROM {{ this }} AS existing
                LEFT JOIN records_to_process AS rtp
                    ON existing.{{ unique_key }} = rtp.{{ unique_key }}
                WHERE rtp.{{ unique_key }} IS NULL
            )

        -- Combine all records
        SELECT *
        FROM unchanged_existing_records

        UNION ALL

        SELECT *
        FROM existing_records_to_expire

        UNION ALL

        SELECT *
        FROM new_records

    {% else %}
        -- Initial load (or full refresh)
        WITH stg_data AS (
            SELECT *
            FROM {{ stg_relation }}
            QUALIFY ROW_NUMBER() OVER (
                PARTITION BY {{ unique_key }}
                ORDER BY {{ updated_at_field }} DESC
            ) = 1
        )

        SELECT
            {{ dbt_utils.generate_surrogate_key([unique_key, updated_at_field]) }} AS {{ surrogate_key_field_name }},
            *,
            CURRENT_TIMESTAMP() AS valid_from,
            TIMESTAMP('9999-12-31 23:59:59.999999') AS valid_to,
            TRUE AS is_current
        FROM stg_data
    {% endif %}

{% endmacro %}
