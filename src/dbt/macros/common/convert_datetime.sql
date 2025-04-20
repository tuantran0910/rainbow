{% macro convert_datetime(datetime_str, format_str) %}
    {% set datetime_obj = modules.datetime.datetime.strptime(datetime_str, format_str) %}
    {{ return(datetime_obj) }}
{% endmacro %}
