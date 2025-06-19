{% macro get_column_names(relation) %}

    {% set columns = adapter.get_columns_in_relation(relation) %}
    {% set column_names = [] %}
    {% for column in columns %}
        {% do column_names.append(column.name) %}
    {% endfor %}

    {{ return(column_names) }}

{% endmacro %}
