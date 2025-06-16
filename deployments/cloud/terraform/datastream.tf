locals {
  datastream_databases = {
    # "postgres" : {
    #   "tables" : ["*"],
    #   "write_disposition" : "merge",
    # }
  }

  datastream_databases_mapping = {
    for key, value in local.datastream_databases : key => {
      "tables"            = value.tables
      "write_disposition" = value.write_disposition
    }
  }
}

# Google Secret Manager secret for datastream password
resource "google_secret_manager_secret" "datastream_password" {
  secret_id = "datastream-password"

  replication {
    auto {}
  }

  labels = {
    "owner" = "tuan-tran"
  }
}

data "google_secret_manager_secret_version" "datastream_password" {
  secret = google_secret_manager_secret.datastream_password.secret_id
}

# Datastream Private Connectivity
resource "google_datastream_private_connection" "datastream_private_connection" {
  location     = "us-central1"
  display_name = "Private Connection for Datastream"

  vpc_peering_config {
    vpc    = google_compute_network.main.id
    subnet = google_compute_subnetwork.main.id
  }

  private_connection_id = "datastream-private-connection"
}

# Connection Profile
resource "google_datastream_connection_profile" "cloudsql_source" {
  for_each = local.datastream_databases_mapping

  display_name          = "CloudSQL Source for ${each.key}"
  location              = local.region
  connection_profile_id = "${each.key}-source-connection-profile"

  postgresql_profile {
    hostname = google_sql_database_instance.main.private_ip_address
    username = local.datastream_username
    password = data.google_secret_manager_secret_version.datastream_password.secret_data
    database = each.key
  }

  private_connectivity {
    private_connection = google_datastream_private_connection.datastream_private_connection.id
  }

  depends_on = [
    google_datastream_private_connection.datastream_private_connection,
    google_sql_database_instance.main,
    google_sql_database.databases,
    google_sql_user.datastream_user
  ]

  labels = {
    "owner" = "tuan-tran"
  }
}

resource "google_datastream_connection_profile" "bq_destination" {
  location              = local.region
  display_name          = "BigQuery Destination"
  connection_profile_id = "bigquery-destination"

  bigquery_profile {}

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
      replication_slot              = "${each.key}_replication_slot"
      include_objects {
        postgresql_schemas {
          schema = "public"
        }
      }
    }
  }

  destination_config {
    destination_connection_profile = google_datastream_connection_profile.bq_destination.id
    bigquery_destination_config {
      data_freshness = "1d"
      single_target_dataset {
        dataset_id = google_bigquery_dataset.datastream_dataset[each.key].dataset_id
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

  depends_on = [
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
}

# Create databases on CloudSQL instance
resource "google_sql_database" "databases" {
  for_each = local.datastream_databases_mapping

  name     = each.key
  instance = google_sql_database_instance.main.name
}

# Create datastream user with replication permissions
resource "google_sql_user" "datastream_user" {
  name     = local.datastream_username
  instance = google_sql_database_instance.main.name
  password = data.google_secret_manager_secret_version.datastream_password.secret_data
}
