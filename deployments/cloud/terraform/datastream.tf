# locals {
#   datastream_databases = {
#     "rainbow" : {
#       "tables" : ["users"],
#       "write_disposition" : "merge",
#     }
#   }

#   datastream_databases_mapping = {
#     for key, value in local.datastream_databases : key => {
#       "tables"            = value.tables
#       "write_disposition" = value.write_disposition
#     }
#   }
# }

# # Google Secret Manager secret for datastream password
# data "google_secret_manager_secret" "datastream_password" {
#   secret_id = "datastream-password"
#   project   = local.project_id
# }

# # Datastream Private Connectivity
# resource "google_datastream_private_connection" "datastream_private_connection" {
#   location     = "us-central1"
#   display_name = "Private Connection for Datastream"

#   vpc_peering_config {
#     vpc    = data.google_compute_network.main.id
#     subnet = "10.4.0.0/29"
#   }

#   private_connection_id = "datastream-private-connection"

#   depends_on = [google_project_service.required_apis]
# }

# data "google_secret_manager_secret_version" "datastream_password" {
#   secret  = data.google_secret_manager_secret.datastream_password.secret_id
#   project = local.project_id
# }

# # Connection Profile
# resource "google_datastream_connection_profile" "cloudsql_source" {
#   for_each = local.datastream_databases_mapping

#   display_name          = "CloudSQL Source for ${each.key}"
#   location              = local.region
#   connection_profile_id = "${each.key}-source-connection-profile"

#   postgresql_profile {
#     hostname = google_compute_instance.datastream_proxy.network_interface[0].network_ip
#     username = local.datastream_username
#     password = data.google_secret_manager_secret_version.datastream_password.secret_data
#     database = each.key
#   }

#   private_connectivity {
#     private_connection = google_datastream_private_connection.datastream_private_connection.id
#   }

#   depends_on = [
#     google_project_service.required_apis,
#     google_compute_instance.datastream_proxy,
#     google_datastream_private_connection.datastream_private_connection
#   ]

#   labels = {
#     "owner" = "tuan-tran"
#   }
# }

# resource "google_datastream_connection_profile" "bq_destination" {
#   location              = local.region
#   display_name          = "BigQuery Destination"
#   connection_profile_id = "bigquery-destination-connection-profile"

#   bigquery_profile {}

#   depends_on = [google_project_service.required_apis]

#   labels = {
#     "owner" = "tuan-tran"
#   }
# }

# # Datastream Stream
# resource "google_datastream_stream" "stream" {
#   for_each = local.datastream_databases_mapping

#   stream_id     = "${each.key}-stream"
#   location      = local.region
#   display_name  = "CloudSQL Streaming from ${each.key} to BigQuery"
#   desired_state = "RUNNING"

#   labels = {
#     "owner" = "tuan-tran"
#   }

#   backfill_all {}

#   source_config {
#     source_connection_profile = google_datastream_connection_profile.cloudsql_source[each.key].id
#     postgresql_source_config {
#       max_concurrent_backfill_tasks = 2
#       publication                   = "${each.key}_publication"
#       replication_slot              = "${each.key}_replication_slot"
#       include_objects {
#         postgresql_schemas {
#           schema = "public"
#           dynamic "postgresql_tables" {
#             for_each = toset(each.value.tables)
#             content {
#               table = postgresql_tables.value
#             }
#           }
#         }
#       }
#     }
#   }

#   destination_config {
#     destination_connection_profile = google_datastream_connection_profile.bq_destination.id
#     bigquery_destination_config {
#       data_freshness = "1d"
#       single_target_dataset {
#         dataset_id = google_bigquery_dataset.datastream_dataset[each.key].dataset_id
#       }

#       dynamic "merge" {
#         for_each = each.value.write_disposition == "merge" ? [1] : []
#         content {}
#       }

#       dynamic "append_only" {
#         for_each = each.value.write_disposition == "append" ? [1] : []
#         content {}
#       }
#     }
#   }

#   depends_on = [
#     google_project_service.required_apis,
#     google_datastream_connection_profile.cloudsql_source,
#     google_datastream_connection_profile.bq_destination,
#     google_bigquery_dataset.datastream_dataset
#   ]
# }

# resource "google_bigquery_dataset" "datastream_dataset" {
#   for_each = local.datastream_databases_mapping

#   project     = local.project_id
#   location    = local.region
#   dataset_id  = "${each.key}__datastream"
#   description = "Dataset for real-time ingestion from CloudSQL ${each.key} to BigQuery"

#   depends_on = [google_project_service.required_apis]

#   labels = {
#     "owner" = "tuan-tran"
#   }
# }

# resource "google_compute_instance" "datastream_proxy" {
#   name = "datastream-proxy"
#   machine_type = "e2-micro"
#   zone = local.cloud_run_zone
#   tags = ["datastream-proxy"]

#   boot_disk {
#     initialize_params {
#       image = "debian-cloud/debian-11"
#     }
#   }

#   network_interface {
#     subnetwork = google_compute_subnetwork.datastream.id
#     access_config {
#       // This ensures an external IP is assigned
#     }
#   }

#   metadata_startup_script = <<-EOF
#     #!/bin/bash
#     apt-get update
#     curl -o /cloud_sql_proxy https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64
#     chmod +x /cloud_sql_proxy
#     ./cloud_sql_proxy -instances=${data.google_sql_database_instance.main.connection_name}=tcp:0.0.0.0:5432 &
#   EOF

#   depends_on = [google_project_service.required_apis]
# }
