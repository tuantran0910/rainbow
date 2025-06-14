# Persistent Disk for Dagster workspace persistence volume
resource "google_compute_disk" "dagster_workspace_disk" {
  name = "${local.project_id}-dagster-workspace-disk"
  type = "pd-ssd"
  zone = local.zone
  size = 5

  labels = {
    component = "dagster"
    env       = "production"
    owner     = "tuan.tran"
  }

  depends_on = [
    google_project_service.required_apis
  ]
}
