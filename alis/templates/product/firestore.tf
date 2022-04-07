// Instantiate a Firestore instance used by the resources.
resource "google_app_engine_application" "app" {
  project = google_project.product_deployment.project_id
  location_id = "europe-west"
  database_type = "CLOUD_FIRESTORE"
  depends_on = [google_project_service.product-services]
}