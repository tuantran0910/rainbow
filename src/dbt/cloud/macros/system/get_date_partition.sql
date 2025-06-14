{% macro get_date_partition(start_date_dt, end_date_dt, timestamp_format) %}
    {% set num_days = (end_date_dt - start_date_dt).days %}
    {% set partitions = [] %}
    {% for i in range(num_days + 1) %}
        {% set current_date = (start_date_dt + modules.datetime.timedelta(days=i)) %}
        {% do partitions.append(current_date.strftime('%Y-%m-%d')) %}
    {% endfor %}

    {{ return(partitions) }}
{% endmacro %}
