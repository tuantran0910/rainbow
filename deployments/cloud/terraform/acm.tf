# GKE Hub Membership
resource "google_gke_hub_membership" "membership" {
  membership_id = "${local.project_id}-membership"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }

  depends_on = [
    google_container_cluster.primary,
    google_project_service.required_apis
  ]
}

# Enable Config Management Feature
resource "google_gke_hub_feature" "config_management" {
  name     = "configmanagement"
  location = "global"

  depends_on = [
    google_gke_hub_membership.membership,
    google_project_service.required_apis
  ]
}

resource "google_service_account" "config_connector" {
  account_id   = "cnrm-sa"
  display_name = "Config Connector Service Account"
}

resource "google_project_iam_member" "config_connector_bindings" {
  for_each = toset([
    "roles/editor"
  ])

  project = local.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.config_connector.email}"
}

resource "google_service_account_iam_member" "config_connector_workload_identity" {
  service_account_id = google_service_account.config_connector.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${local.project_id}.svc.id.goog[cnrm-system/cnrm-controller-manager]"
}
