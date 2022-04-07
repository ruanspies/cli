resource "google_project_service" "product-services" {
  project  = google_project.product_deployment.project_id
  for_each = toset( [
    "cloudbuild.googleapis.com",
    "pubsub.googleapis.com",
    "firestore.googleapis.com",
    "run.googleapis.com",
    "appengine.googleapis.com",
    "apigateway.googleapis.com",
    "servicecontrol.googleapis.com",
    "servicemanagement.googleapis.com",
    "bigquery.googleapis.com",
    "cloudscheduler.googleapis.com",
    "cloudidentity.googleapis.com",
    "cloudbilling.googleapis.com",
    "workflows.googleapis.com",
    "sourcerepo.googleapis.com",
    "compute.googleapis.com",
    "bigtableadmin.googleapis.com",
    "dns.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "serviceusage.googleapis.com",
    "iam.googleapis.com",
    "iap.googleapis.com",
    "admin.googleapis.com",
    "container.googleapis.com",
    "groupssettings.googleapis.com",
    "spanner.googleapis.com",
    "artifactregistry.googleapis.com"
  ] )
  service = each.key
}