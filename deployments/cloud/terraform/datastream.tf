locals {
  datastream_databases = {
    "rainbow" : {
      "tables" : [
        "authors",
        "book_authors",
        "books",
        "categories",
        "order_items",
        "orders",
        "payments",
        "promotions",
        "sellers",
        "users"
      ],
      "write_disposition" : "merge",
    }
  }

  datastream_databases_mapping = {
    for key, value in local.datastream_databases : key => {
      "tables"            = value.tables
      "write_disposition" = value.write_disposition
    }
  }
}

# Google Secret Manager secret for datastream password
data "google_secret_manager_secret" "datastream_password" {
  secret_id = "datastream-password"
  project   = local.project_id
}

data "google_secret_manager_secret_version" "datastream_password" {
  secret  = data.google_secret_manager_secret.datastream_password.secret_id
  project = local.project_id
}

# Connection Profile
resource "google_datastream_connection_profile" "cloudsql_source" {
  for_each = local.datastream_databases_mapping

  display_name          = "CloudSQL Source for ${each.key}"
  location              = local.region
  connection_profile_id = "${each.key}-postgres-source-connection-profile"

  postgresql_profile {
    hostname = data.google_sql_database_instance.main.public_ip_address
    username = local.datastream_username
    password = data.google_secret_manager_secret_version.datastream_password.secret_data
    database = each.key
  }

  depends_on = [google_project_service.required_apis]

  labels = {
    "owner" = "tuan-tran"
  }
}

resource "google_datastream_connection_profile" "bq_destination" {
  location              = local.region
  display_name          = "BigQuery Destination"
  connection_profile_id = "bigquery-destination-connection-profile"

  bigquery_profile {}

  depends_on = [google_project_service.required_apis]

  labels = {
    "owner" = "tuan-tran"
  }
}

# Datastream Stream
resource "google_datastream_stream" "stream" {
  for_each = local.datastream_databases_mapping

  stream_id     = "${each.key}-stream"
  location      = local.region
  display_name  = "CloudSQL Streaming from ${each.key} to BigQuery"
  desired_state = "RUNNING"

  labels = {
    "owner" = "tuan-tran"
  }

  backfill_all {}

  source_config {
    source_connection_profile = google_datastream_connection_profile.cloudsql_source[each.key].id
    postgresql_source_config {
      max_concurrent_backfill_tasks = 2
      publication                   = "${each.key}_publication"
      replication_slot              = "${each.key}_replication"
      include_objects {
        postgresql_schemas {
          schema = "public"
          dynamic "postgresql_tables" {
            for_each = toset(each.value.tables)
            content {
              table = postgresql_tables.value
            }
          }
        }
      }
    }
  }

  destination_config {
    destination_connection_profile = google_datastream_connection_profile.bq_destination.id
    bigquery_destination_config {
      single_target_dataset {
        dataset_id = "${local.project_id}:${google_bigquery_dataset.datastream_dataset[each.key].dataset_id}"
      }

      dynamic "merge" {
        for_each = each.value.write_disposition == "merge" ? [1] : []
        content {}
      }

      dynamic "append_only" {
        for_each = each.value.write_disposition == "append" ? [1] : []
        content {}
      }
    }
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [
    google_project_service.required_apis,
    google_datastream_connection_profile.cloudsql_source,
    google_datastream_connection_profile.bq_destination,
    google_bigquery_dataset.datastream_dataset
  ]
}

resource "google_bigquery_dataset" "datastream_dataset" {
  for_each = local.datastream_databases_mapping

  project     = local.project_id
  location    = local.region
  dataset_id  = "${each.key}__datastream"
  description = "Dataset for real-time ingestion from CloudSQL ${each.key} to BigQuery"

  labels = {
    "owner" = "tuan-tran"
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.required_apis]
}
