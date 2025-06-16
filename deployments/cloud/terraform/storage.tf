resource "google_storage_bucket" "cubejs" {
  name                        = "rainbow-data-production-cubejs"
  location                    = local.region
  uniform_bucket_level_access = true
  public_access_prevention    = true
}

resource "google_storage_bucket" "dbt" {
  name                        = "rainbow-data-production-dbt"
  location                    = local.region
  uniform_bucket_level_access = true
  public_access_prevention    = true
}
