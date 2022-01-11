terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.81.0"
    }
  }
  backend "gcs" {
    bucket = "provided_at_runtime_by_alis"
    prefix = "provided_at_runtime_by_alis"
  }
}

provider "google-beta" {
  project = var.ALIS_OS_PROJECT
}

provider "google" {
  project = var.ALIS_OS_PROJECT
}