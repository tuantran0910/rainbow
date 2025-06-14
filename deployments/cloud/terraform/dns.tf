resource "google_compute_global_address" "dagster_ingress_ip" {
  name         = "dagster-ingress-ip"
  address_type = "EXTERNAL"
}

# Managed zone for tuantrann.work
data "google_dns_managed_zone" "tuantrann_work" {
  name = "tuantrann-work-zone"
}

# DNS record for Dagster ingress
resource "google_dns_record_set" "dagster_dns" {
  name         = "dagster.tuantrann.work."
  managed_zone = data.google_dns_managed_zone.tuantrann_work.name
  type         = "A"
  ttl          = 300
  rrdatas      = [google_compute_global_address.dagster_ingress_ip.address]
}

# SSL Policy
resource "google_compute_ssl_policy" "rainbow_ssl_policy" {
  name = "rainbow-ssl-policy"
  profile = "MODERN"
  min_tls_version = "TLS_1_2"
}
