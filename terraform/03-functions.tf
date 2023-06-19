module "cloud_functions2" {
  source = "GoogleCloudPlatform/cloud-functions/google"

  project_id        = var.project_id
  function_name     = var.function_name
  function_location = var.region
  runtime           = "go120"
  entrypoint        = "init"
  storage_source = {
    bucket     = google_storage_bucket.bucket.name
    object     = google_storage_bucket_object.function-source.name
    generation = null
  }

  event_trigger = {
    event_type = "google.storage.object.finalize"
    resource   = google_storage_bucket.bucket.name
  }

  docker_repository = "carlinhoscamilo/resizenator:latest"
}
