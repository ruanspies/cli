module "group" {
  source = "terraform-google-modules/group/google"
  version = "~> 0.1"
  display_name = google_project.product_deployment.project_id
  id = "${google_project.product_deployment.project_id}@alis.exchange"
  description = "Group to manage 'Client Level' access to the ${google_project.product_deployment.project_id} deployment"
  domain = "ric.co.za"
  managers = [var.ALIS_OS_OWNER]
}