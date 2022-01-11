resource "google_project_service" "cloudbuild_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "cloudbuild.googleapis.com"
}
resource "google_project_service" "pubsub_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "pubsub.googleapis.com"
}
resource "google_project_service" "firestore_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "firestore.googleapis.com"
}
resource "google_project_service" "run_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "run.googleapis.com"
}
resource "google_project_service" "appengine_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "appengine.googleapis.com"
}
resource "google_project_service" "apigateway_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "apigateway.googleapis.com"
}
resource "google_project_service" "servicecontrol_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "servicecontrol.googleapis.com"
}
resource "google_project_service" "servicemanagement_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "servicemanagement.googleapis.com"
}
resource "google_project_service" "bigquery_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "bigquery.googleapis.com"
}
resource "google_project_service" "cloudscheduler_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "cloudscheduler.googleapis.com"
}

resource "google_project_service" "cloudidentity_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "cloudidentity.googleapis.com"
}

resource "google_project_service" "billing_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "cloudbilling.googleapis.com"
}

resource "google_project_service" "workflows_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "workflows.googleapis.com"
}

resource "google_project_service" "sourcerepo_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "sourcerepo.googleapis.com"
}

resource "google_project_service" "compute_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "cloudresourcemanager_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "cloudresourcemanager.googleapis.com"
}

resource "google_project_service" "serviceusage_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "serviceusage.googleapis.com"
}

resource "google_project_service" "iam_googleapis_com" {
  project = google_project.product_deployment.project_id
  service = "iam.googleapis.com"
}

