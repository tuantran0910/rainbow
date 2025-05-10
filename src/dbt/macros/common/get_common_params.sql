{% macro get_common_params(lookback_in_days=1) %}

    {% set model_params = namespace() %}

    {# Define timestamp format #}
    {% set timestamp_format = '%Y-%m-%dT%H:%M:%S+00:00' %}

    {# Calculate current timestamp #}
    {% set current_timestamp = convert_datetime(var('execution_datetime', modules.datetime.datetime.now().strftime(timestamp_format)), timestamp_format) %}
    {% set current_timestamp_tz_hcm = current_timestamp + modules.datetime.timedelta(hours=7) %}

    {# Calculate lookback interval #}
    {% set lookback_timestamp_tz_hcm = current_timestamp_tz_hcm + modules.datetime.timedelta(days=-lookback_in_days) %}
    {% set model_params.incremental_dates_quoted_tz_hcm = get_date_partition(
        start_date_dt=lookback_timestamp_tz_hcm,
        end_date_dt=current_timestamp_tz_hcm,
        timestamp_format=timestamp_format
    ) %}

    {% do return(model_params) %}

{% endmacro %}
