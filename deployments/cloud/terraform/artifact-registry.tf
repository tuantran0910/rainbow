# Artifact Registry Repository for Docker images
resource "google_artifact_registry_repository" "docker_repo" {
  repository_id = "rainbow-docker"
  location      = local.region
  format        = "DOCKER"
  description   = "Docker repository for Rainbow application images"

  depends_on = [google_project_service.required_apis]

  labels = {
    environment = "production"
    component   = "container-registry"
  }
}

# IAM binding for the GitHub Actions service account to push to Artifact Registry
resource "google_project_iam_member" "github_actions_artifact_registry_writer" {
  project = local.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:github-actions-deployer@${local.project_id}.iam.gserviceaccount.com"
}

# IAM binding for Cloud Run service account to pull from Artifact Registry
resource "google_project_iam_member" "api_artifact_registry_reader" {
  project = local.project_id
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${google_service_account.api.email}"
}

# Output the Artifact Registry repository URL
output "artifact_registry_url" {
  description = "The URL of the Artifact Registry repository"
  value       = "${local.region}-docker.pkg.dev/${local.project_id}/${google_artifact_registry_repository.docker_repo.repository_id}"
}
