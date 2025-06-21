resource "google_storage_bucket" "cubejs" {
  name                        = "rainbow-data-production-cubejs"
  location                    = local.region
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_storage_bucket" "dbt" {
  name                        = "rainbow-data-production-dbt"
  location                    = local.region
  uniform_bucket_level_access = true
  # Allow public access for static website hosting
  public_access_prevention = "inherited"

  # Enable static website hosting
  website {
    main_page_suffix = "docs/static_index.html"
    not_found_page   = "docs/static_index.html"
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_storage_bucket_iam_member" "dbt_docs_public_read" {
  bucket = google_storage_bucket.dbt.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

# Backend bucket for load balancer
resource "google_compute_backend_bucket" "dbt_docs" {
  name        = "dbt-docs-backend"
  bucket_name = google_storage_bucket.dbt.name
  enable_cdn  = true

  cdn_policy {
    cache_mode       = "CACHE_ALL_STATIC"
    default_ttl      = 3600
    max_ttl          = 86400
    client_ttl       = 3600
    negative_caching = true
  }
}
