{% macro dim_scd_type_1(
    stg_relation,
    unique_keys,
    surrogate_key_field_name,
    updated_at_field,
    time_zone='UTC',
    lookback_in_days=1
) %}
    {% set model_params = get_common_params(lookback_in_days) %}

    WITH
        stg_data AS (
            SELECT *
            FROM {{ stg_relation }}
            {% if is_incremental() %}
                WHERE DATE({{ updated_at_field }}, '{{ time_zone }}') IN {{ model_params.incremental_dates_quoted_tz_hcm }}
            {% endif %}
            QUALIFY ROW_NUMBER() OVER (
                PARTITION BY {{ unique_keys | join(', ') }}
                ORDER BY {{ updated_at_field }} DESC
            ) = 1
        )

    SELECT
        {{ dbt_utils.generate_surrogate_key(unique_keys) }} AS {{ surrogate_key_field_name }},
        *
    FROM stg_data

{% endmacro %}
