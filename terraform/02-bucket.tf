resource "google_storage_bucket" "bucket" {
  name                        = var.bucket_name
  location                    = var.bucket_language
  uniform_bucket_level_access = true
  project                     = var.project_id
}

resource "google_storage_bucket_object" "function-source" {
  name   = "function.zip"
  bucket = google_storage_bucket.bucket.name
  source = "./../function.zip"
}
