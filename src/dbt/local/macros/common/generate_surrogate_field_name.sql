{% macro generate_surrogate_field_name(prefix='') %}

    {% set surrogate_key_field = 'surrogate_key' %}
    {% if prefix != '' %}
        {% set surrogate_key_field = prefix ~ '_key' %}
    {% endif %}

    {{ return(surrogate_key_field) }}

{% endmacro %}
