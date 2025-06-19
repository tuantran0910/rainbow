{% macro raw_model_sql(source, layer) %}

    SELECT *
    EXCEPT(_dlt_load_id, _dlt_id),
        DATE(created_at, 'Asia/Ho_Chi_Minh') AS created_date_tz_hcm,
        DATE(updated_at, 'Asia/Ho_Chi_Minh') AS updated_date_tz_hcm,
        {% if layer == 'batch' %}
        STRUCT(
            'batch' AS layer,
            CURRENT_TIMESTAMP() AS dbt_loaded_at
        ) AS _dbt_metadata
        {% else %}
        STRUCT(
            'stream' AS layer,
            CURRENT_TIMESTAMP() AS dbt_loaded_at
        ) AS _dbt_metadata
        {% endif %}
    FROM {{ source }}

{% endmacro %}
