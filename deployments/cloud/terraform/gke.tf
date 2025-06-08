# GKE Cluster
resource "google_container_cluster" "primary" {
  name     = "${local.project_id}-gke"
  location = local.zone

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1
  deletion_protection      = false

  network    = google_compute_network.main.name
  subnetwork = google_compute_subnetwork.main.name

  # IP allocation policy
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }

  # Workload Identity
  workload_identity_config {
    workload_pool = "${local.project_id}.svc.id.goog"
  }

  # Private cluster configuration
  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = "172.16.0.0/28"
  }

  depends_on = [
    google_project_service.required_apis,
    google_compute_subnetwork.main
  ]
}

# Primary node pool
resource "google_container_node_pool" "primary_nodes" {
  name       = "${local.project_id}-node-pool"
  location   = local.zone
  cluster    = google_container_cluster.primary.name
  node_count = 1

  node_config {
    preemptible  = false
    machine_type = "e2-small"

    service_account = google_service_account.gke_node_sa.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    labels = {
      component = "gke"
      env       = "production"
    }

    metadata = {
      disable-legacy-endpoints = "true"
    }
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  depends_on = [
    google_container_cluster.primary
  ]
}

# Service account for GKE nodes
resource "google_service_account" "gke_node_sa" {
  account_id   = "${local.project_id}-gke-sa"
  display_name = "GKE Node Service Account"
}

# IAM bindings for the GKE node service account
resource "google_project_iam_member" "gke_node_sa_bindings" {
  for_each = toset([
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles/monitoring.viewer",
    "roles/storage.objectViewer"
  ])

  project = local.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.gke_node_sa.email}"
}
