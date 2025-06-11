resource "google_compute_global_address" "dagster_ingress_ip" {
  name         = "dagster-ingress-ip"
  address_type = "EXTERNAL"
}

output "dagster_ingress_ip" {
  value       = google_compute_global_address.dagster_ingress_ip.address
  description = "External IP address for Dagster ingress"
}

resource "google_dns_managed_zone" "tuantrann_work" {
  name        = "tuantrann-work-zone"
  dns_name    = "tuantrann.work."
  description = "DNS zone for tuantrann.work domain"
}

resource "google_dns_record_set" "dagster_dns" {
  name         = "dagster.tuantrann.work."
  managed_zone = google_dns_managed_zone.tuantrann_work.name
  type         = "A"
  ttl          = 300
  rrdatas      = [google_compute_global_address.dagster_ingress_ip.address]
}

output "google_dns_nameservers" {
  value       = google_dns_managed_zone.tuantrann_work.name_servers
  description = "Nameservers to configure in GoDaddy"
}
