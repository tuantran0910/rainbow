# Router
resource "google_compute_router" "router" {
  name    = "${local.project_id}-router"
  region  = local.region
  network = data.google_compute_network.main.id

  depends_on = [google_project_service.required_apis]
}

resource "google_compute_router_nat" "nat" {
  name                               = "${local.project_id}-nat"
  router                             = google_compute_router.router.name
  region                             = local.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }

  depends_on = [google_compute_router.router]
}
