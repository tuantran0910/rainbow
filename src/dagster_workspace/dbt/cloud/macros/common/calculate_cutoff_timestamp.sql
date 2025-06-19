{% macro calculate_cutoff_timestamp(cutoff_days=1) %}
    {% set cutoff_timestamp = modules.datetime.datetime.now() - modules.datetime.timedelta(days=cutoff_days) %}
    {% set cutoff_date_str = cutoff_timestamp.strftime('%Y-%m-%d') %}
    {{ return("'" + cutoff_date_str + "'") }}
{% endmacro %}
