{{
    config(
        materialized='ephemeral'
    )
}}

{{
    raw_model_sql(
        source=source('rainbow_datastream', 'book_authors'),
        layer='stream'
    )
}}
