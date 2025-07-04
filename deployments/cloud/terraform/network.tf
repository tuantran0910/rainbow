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

# resource "google_compute_subnetwork" "datastream" {
#   name          = "${local.project_id}-datastream-subnet"
#   ip_cidr_range = "10.3.0.0/28"
#   region        = local.cloud_run_region
#   network       = data.google_compute_network.main.id
# }

# data "http" "google_ip_ranges" {
#   url = "https://www.gstatic.com/ipranges/goog.json"
# }

# locals {
#   google_cloud_ip_ranges = [
#     for r in jsondecode(data.http.google_ip_ranges.response_body).prefixes : r.ipv4Prefix if try(r.ipv4Prefix, null) != null
#   ]
# }

# resource "google_compute_firewall" "allow_datastream_proxy" {
#   name    = "${local.project_id}-allow-datastream-proxy"
#   network = data.google_compute_network.main.id
#   direction = "INGRESS"

#   allow {
#     protocol = "tcp"
#     ports    = ["5432"]
#   }

#   source_ranges = local.google_cloud_ip_ranges
#   target_tags   = ["datastream-proxy"]
#   description   = "Allow ingress from Datastream public IPs to the proxy VM."
# }
