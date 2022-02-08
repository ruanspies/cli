// The project hosting the ProductDeployment resource
resource "google_project" "product_deployment" {
  name = var.ALIS_OS_PROJECT
  project_id = var.ALIS_OS_PROJECT
  folder_id = (length(regexall("^[a-z0-4]+-[a-z]{2}-(prod)-[a-z0-9]+$", var.ALIS_OS_PROJECT)) > 0 ? data.terraform_remote_state.product.outputs.folder_prod : data.terraform_remote_state.product.outputs.folder_dev)
  labels = {
    "managed-by-alis" : true
  }
  billing_account = var.ALIS_OS_BILLING
}

// place the ProductDeployment in the product folder
data "terraform_remote_state" "product" {
  backend = "gcs"
  config = {
    bucket = var.ALIS_OS_ORG_BACKEND_BUCKET
    prefix = var.ALIS_OS_ORG_BACKEND_PRODUCT_PREFIX
  }
}

// The project hosting the Product resource
data "google_project" "product" {
  project_id = var.ALIS_OS_PRODUCT_PROJECT
}