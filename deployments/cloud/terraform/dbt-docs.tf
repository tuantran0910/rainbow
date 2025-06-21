# Global IP address for dbt docs
resource "google_compute_global_address" "dbt_docs_ip" {
  name         = "dbt-docs-ip"
  address_type = "EXTERNAL"
}

# SSL certificate for dbt-docs.tuantrann.work
resource "google_compute_managed_ssl_certificate" "dbt_docs_ssl" {
  name = "dbt-docs-ssl-cert"

  managed {
    domains = ["dbt-docs.tuantrann.work"]
  }
}

# URL map for dbt docs
resource "google_compute_url_map" "dbt_docs" {
  name        = "dbt-docs-url-map"
  description = "URL map for dbt docs static site"

  default_service = google_compute_backend_bucket.dbt_docs.id

  host_rule {
    hosts        = ["dbt-docs.tuantrann.work"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = google_compute_backend_bucket.dbt_docs.id
  }
}

# HTTPS proxy
resource "google_compute_target_https_proxy" "dbt_docs" {
  name             = "dbt-docs-https-proxy"
  url_map          = google_compute_url_map.dbt_docs.id
  ssl_certificates = [google_compute_managed_ssl_certificate.dbt_docs_ssl.id]
  ssl_policy       = google_compute_ssl_policy.rainbow_ssl_policy.id
}

# HTTP to HTTPS redirect
resource "google_compute_url_map" "dbt_docs_http_redirect" {
  name = "dbt-docs-http-redirect"

  default_url_redirect {
    https_redirect         = true
    redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
    strip_query            = false
  }
}

resource "google_compute_target_http_proxy" "dbt_docs_http" {
  name    = "dbt-docs-http-proxy"
  url_map = google_compute_url_map.dbt_docs_http_redirect.id
}

# Global forwarding rule for HTTPS
resource "google_compute_global_forwarding_rule" "dbt_docs_https" {
  name                  = "dbt-docs-https-forwarding-rule"
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL"
  port_range            = "443"
  target                = google_compute_target_https_proxy.dbt_docs.id
  ip_address            = google_compute_global_address.dbt_docs_ip.id
}

# Global forwarding rule for HTTP (redirect to HTTPS)
resource "google_compute_global_forwarding_rule" "dbt_docs_http" {
  name                  = "dbt-docs-http-forwarding-rule"
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL"
  port_range            = "80"
  target                = google_compute_target_http_proxy.dbt_docs_http.id
  ip_address            = google_compute_global_address.dbt_docs_ip.id
}
