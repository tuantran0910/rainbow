# Private IP configuration for CloudSQL
resource "google_compute_global_address" "private_ip_address" {
  name          = "${local.project_id}-private-ip"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.main.id
  depends_on    = [google_project_service.required_apis]
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = data.google_compute_network.main.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]

  lifecycle {
    prevent_destroy = false
  }
}

data "google_sql_database_instance" "main" {
  name = "${local.project_id}-db"
}

data "google_sql_database" "rainbow" {
  name     = "rainbow"
  instance = data.google_sql_database_instance.main.name
}

data "google_secret_manager_secret_version" "api_password" {
  secret = "api-password"
}
