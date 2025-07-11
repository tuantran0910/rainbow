{% macro lambda_view(
    batch_relation,
    stream_stg_relation,
    unique_key,
    updated_date_field,
    cutoff_days=1,
    scd_type2_field='is_current',
    get_scd_type2_latest=false,
    batch_exclude_columns=[]
) %}
    {# Calculate cutoff timestamp for lambda view #}
    {% set cutoff_timestamp = calculate_cutoff_timestamp(cutoff_days=cutoff_days) %}

    WITH
        -- Batch layer data (T-1): Data older than cutoff
        batch_data AS (
            SELECT * {% if batch_exclude_columns %} EXCEPT({{ batch_exclude_columns | join(', ') }}) {% endif %}
            FROM {{ batch_relation }}
            WHERE {{ updated_date_field }} < DATE({{ cutoff_timestamp }})
            {% if get_scd_type2_latest %}
                AND {{ scd_type2_field }} = TRUE
            {% endif %}
        ),

        -- Speed layer data (T): Recent streaming data
        stream_data AS (
            SELECT *
            FROM {{ stream_stg_relation }}
            WHERE {{ updated_date_field }} >= DATE({{ cutoff_timestamp }})
        ),

        -- Unified view: Stream data takes precedence over batch for overlapping periods
        unified AS (
            SELECT *
            FROM stream_data

            UNION ALL BY NAME

            SELECT batch.*
            FROM batch_data AS batch
            LEFT JOIN stream_data AS stream
                ON batch.{{ unique_key }} = stream.{{ unique_key }}
            WHERE stream.{{ unique_key }} IS NULL
        )

    SELECT *
    FROM unified

{% endmacro %}
