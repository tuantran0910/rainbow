resource "google_storage_bucket" "cubejs" {
  name                        = "rainbow-data-production-cubejs"
  location                    = local.region
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"
}

# resource "google_storage_bucket" "dbt" {
#   name                        = "rainbow-data-production-dbt"
#   location                    = local.region
#   uniform_bucket_level_access = true
#   # Allow public access for static website hosting
#   public_access_prevention = "inherited"

#   # Enable static website hosting
#   website {
#     main_page_suffix = "docs/static_index.html"
#     not_found_page   = "docs/static_index.html"
#   }

#   lifecycle {
#     prevent_destroy = true
#   }
# }

data "google_storage_bucket" "dbt" {
  name = "${local.project_id}-dbt"
}

resource "google_storage_bucket_iam_member" "dbt_docs_public_read" {
  bucket = data.google_storage_bucket.dbt.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}
