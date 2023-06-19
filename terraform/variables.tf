variable "project_id" {
  type        = string
  description = "Project ID"
}

variable "region" {
  type        = string
  description = "Region"
  default     = "us-central1"
}

variable "zone" {
  type        = string
  description = "Zone"
  default     = "us-central1-a"
}

variable "backend_bucket" {
  type        = string
  description = "Backend bucket name"
  default     = "resizenator-state"
}

variable "bucket_name" {
  type        = string
  description = "Bucket name"
  default     = "resizenator-images"
}

variable "bucket_language" {
  type        = string
  description = "Bucket language"
  default     = "US"
}

variable "function_name" {
  type        = string
  description = "Function name"
  default     = "resizenator"
}
