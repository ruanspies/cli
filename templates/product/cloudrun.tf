// The Cloud Run service agent should have view access to artifact registries in the ProductDeployment resource
resource "google_artifact_registry_repository_iam_member" "neurons" {
  provider = google-beta
  project = data.google_project.product.project_id
  location = "europe-west1"
  repository = "neurons"
  role = "roles/artifactregistry.reader"
  member = "serviceAccount:service-${google_project.product_deployment.number}@serverless-robot-prod.iam.gserviceaccount.com"
  depends_on = [
    google_project_service.run_googleapis_com]
}