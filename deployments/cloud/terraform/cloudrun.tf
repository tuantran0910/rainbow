# Cloud Run API Service
resource "google_cloud_run_v2_service" "api" {
  name         = "rainbow-api"
  location     = local.region
  ingress      = "INGRESS_TRAFFIC_ALL"
  launch_stage = "GA"

  template {
    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }

    containers {
      image = "gcr.io/cloudrun/hello:latest"

      ports {
        container_port = 5000
        name           = "http1"
      }

      env {
        name  = "GIN_MODE"
        value = "release"
      }

      # env {
      #   name  = "DB_HOST"
      #   value = google_sql_database_instance.main.private_ip_address
      # }

      # env {
      #   name  = "DB_NAME"
      #   value = google_sql_database.main.name
      # }

      # env {
      #   name  = "DB_USER"
      #   value = google_sql_user.app_user.name
      # }

      # env {
      #   name = "DB_PASSWORD"
      #   value_source {
      #     secret_key_ref {
      #       secret  = google_secret_manager_secret.db_password.secret_id
      #       version = "latest"
      #     }
      #   }
      # }

      # env {
      #   name  = "DB_PORT"
      #   value = "5432"
      # }

      # env {
      #   name = "JWT_SECRET"
      #   value_source {
      #     secret_key_ref {
      #       secret  = google_secret_manager_secret.jwt_secret.secret_id
      #       version = "latest"
      #     }
      #   }
      # }

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
    google_sql_database_instance.main,
    google_sql_user.api_user,
    data.google_secret_manager_secret_version.api_password
  ]
}

# VPC Access Connector for Cloud Run to connect to private resources
resource "google_vpc_access_connector" "main" {
  name          = "${local.project_id}-vpc-connector"
  region        = local.region
  ip_cidr_range = "10.8.0.0/28"
  network       = google_compute_network.main.name
  max_instances = 3
  min_instances = 2
  machine_type  = "e2-micro"

  depends_on = [google_project_service.required_apis]
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
