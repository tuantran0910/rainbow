{% macro dim_scd_type_1(stg_relation, unique_key, updated_at_field, surrogate_key_field_name, lookback_in_days=1, is_stream_mode=false) %}

    {% if not is_stream_mode %}
        {% set model_params = get_common_params(lookback_in_days) %}
    {% endif %}

    WITH
        stg_data AS (
            SELECT *
            FROM {{ stg_relation }}
            {% if is_incremental() and not is_stream_mode %}
                WHERE {{ updated_at_field }} IN {{ model_params.incremental_dates_quoted_tz_hcm }}
            {% endif %}
        )

    SELECT
        {{ dbt_utils.generate_surrogate_key(['stg_data.' ~ unique_key]) }} AS {{ surrogate_key_field_name }},
        stg_data.*,
        {% if is_stream_mode %}
        TRUE AS is_current
        {% else %}
        TRUE AS is_current
        {% endif %}
    FROM stg_data

{% endmacro %}
