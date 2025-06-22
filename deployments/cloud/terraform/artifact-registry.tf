# Artifact Registry Repository for Docker images
data "google_artifact_registry_repository" "docker_repo" {
  repository_id = "rainbow-docker"
  location      = local.region
}

# IAM binding for the GitHub Actions service account to push to Artifact Registry
resource "google_project_iam_member" "github_actions_artifact_registry_writer" {
  project = local.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:github-actions-deployer@${local.project_id}.iam.gserviceaccount.com"
}
