resource "google_cloud_run_service" "default" {
  provider = google-beta
  name     = var.ALIS_OS_NEURON
  location = "europe-west1"

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale": "100"
        //        "run.googleapis.com/cpu-throttling": "false"
      }
      name = "{{.Contract}}-{{.Neuron}}-{{.VersionMajor}}-${uuid()}"
    }
    spec {
      containers {
        image = "europe-west1-docker.pkg.dev/${var.ALIS_OS_PRODUCT_PROJECT}/neurons/${var.ALIS_OS_NEURON}:${var.ALIS_OS_NEURON_VERSION_COMMIT_SHA}"
        env {
          name = "ALIS_OS_PROJECT"
          value = var.ALIS_OS_PROJECT
        }
        resources {
          limits = {
            cpu: "1000m"
            memory: "2Gi"
          }
        }
      }
      container_concurrency = 80
      timeout_seconds = 90
      service_account_name = "alis-exchange@${var.ALIS_OS_PROJECT}.iam.gserviceaccount.com"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "invoker" {
  location = google_cloud_run_service.default.location
  project = google_cloud_run_service.default.project
  service = google_cloud_run_service.default.name
  role = "roles/run.invoker"
  member = "group:${var.ALIS_OS_PROJECT}@alis.exchange"
}