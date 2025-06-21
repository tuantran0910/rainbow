# Cloud Run API Service
resource "google_cloud_run_v2_service" "api" {
  name                = "rainbow-api"
  location            = local.cloud_run_region
  ingress             = "INGRESS_TRAFFIC_ALL"
  launch_stage        = "GA"
  deletion_protection = false

  template {
    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }

    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"

      ports {
        container_port = 5000
        name           = "http1"
      }

      env {
        name  = "GIN_MODE"
        value = "release"
      }

      env {
        name  = "ENV"
        value = "production"
      }

      env {
        name  = "USE_CLOUD_SQL"
        value = "true"
      }

      env {
        name  = "DB_HOST"
        value = data.google_sql_database_instance.main.private_ip_address
      }

      env {
        name  = "DB_NAME"
        value = data.google_sql_database.rainbow.name
      }

      env {
        name  = "DB_USER"
        value = "rainbow"
      }

      env {
        name  = "DB_PASSWORD"
        value = data.google_secret_manager_secret_version.api_password.secret_data
      }

      env {
        name  = "DB_PORT"
        value = "5432"
      }

      env {
        name  = "MIGRATIONS_DIR"
        value = "/migrations"
      }

      env {
        name  = "LOG_LEVEL"
        value = "info"
      }

      env {
        name  = "ENABLE_CONSOLE"
        value = "false"
      }

      env {
        name  = "JWT_SECRET"
        value = data.google_secret_manager_secret_version.jwt_secret.secret_data
      }

      env {
        name  = "SERVER_PORT"
        value = "5000"
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
        startup_cpu_boost = true
      }
    }

    timeout                          = "300s"
    max_instance_request_concurrency = 100

    vpc_access {
      connector = google_vpc_access_connector.main.id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    service_account = google_service_account.api.email
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }

  depends_on = [
    google_project_service.required_apis,
    data.google_sql_database_instance.main,
    data.google_secret_manager_secret_version.api_password,
    data.google_secret_manager_secret_version.jwt_secret
  ]
}

# Service Account for Cloud Run
resource "google_service_account" "api" {
  account_id   = "rainbow-api"
  display_name = "Rainbow API Service Account"
  description  = "Service account for Rainbow API Cloud Run service"
}

# IAM bindings for the service account
resource "google_project_iam_member" "api_sql_client" {
  project = local.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.api.email}"
}

resource "google_project_iam_member" "api_secret_accessor" {
  project = local.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.api.email}"
}

# Output the Cloud Run service name
output "api_service_name" {
  description = "The name of the Cloud Run service"
  value       = google_cloud_run_v2_service.api.name
}

# VPC Access Connector for Cloud Run to connect to private resources
resource "google_vpc_access_connector" "main" {
  name          = "cloud-run-vpc-connector"
  region        = local.cloud_run_region
  ip_cidr_range = "10.9.0.0/28"
  network       = data.google_compute_network.main.name
  max_instances = 3
  min_instances = 2
  machine_type  = "f1-micro"

  depends_on = [google_project_service.required_apis]
}

# Reference the existing JWT secret from Google Secret Manager
data "google_secret_manager_secret_version" "jwt_secret" {
  secret = "jwt-secret"
}
