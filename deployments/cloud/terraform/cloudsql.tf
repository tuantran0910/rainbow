# Private IP configuration for CloudSQL
resource "google_compute_global_address" "private_ip_address" {
  name          = "${local.project_id}-private-ip"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.main.id
  depends_on    = [google_project_service.required_apis]
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.main.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]

  lifecycle {
    prevent_destroy = false
  }
}

# CloudSQL instance
resource "google_sql_database_instance" "main" {
  name                = "${local.project_id}-db"
  database_version    = "POSTGRES_16"
  region              = local.region
  deletion_protection = false

  settings {
    edition = "ENTERPRISE"
    tier    = "db-f1-micro"

    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }

    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.main.id
    }
  }

  depends_on = [
    google_service_networking_connection.private_vpc_connection,
    google_project_service.required_apis
  ]

  lifecycle {
    prevent_destroy = false
  }
}
