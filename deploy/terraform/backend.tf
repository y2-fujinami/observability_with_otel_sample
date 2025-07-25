terraform {
  backend "gcs" {
    bucket = "tfstate-f8d77d8a-b043-1085-af22-6f6faa5217b9"
  }
}

# tfstateを保存するためのGCSバケット
resource "google_storage_bucket" "tfstate" {
  name          = "tfstate-f8d77d8a-b043-1085-af22-6f6faa5217b9"
  location      = var.default_region
  storage_class = "STANDARD"
  versioning {
    enabled = true
  }
}