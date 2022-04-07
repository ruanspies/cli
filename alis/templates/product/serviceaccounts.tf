resource "google_service_account" "alis_exchange" {
  project = google_project.product_deployment.project_id
  account_id = "alis-exchange"
  description = "Service account to manage resources and services in the ${upper(var.ALIS_OS_PRODUCT)} deployment"
  display_name = "alis_ Exchange ProductDeployment Service Account"
}

resource "google_project_iam_member" "project" {
  project = google_project.product_deployment.project_id
  role = "roles/editor"
  member = "serviceAccount:${google_service_account.alis_exchange.email}"
}

# Ensure that the group of developers are able to create a service account for local development.
resource "google_project_iam_member" "product" {
  project = google_project.product_deployment.project_id
  role = "roles/iam.serviceAccountKeyAdmin"
  member = "group:${var.ALIS_OS_PRODUCT_PROJECT}@identity.${var.ALIS_OS_DOMAIN}"
}