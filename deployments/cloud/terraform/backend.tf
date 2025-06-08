terraform {
  backend "gcs" {
    bucket = "rainbow-data-production-terraform-state"
    prefix = "terraform/state"
  }
}
