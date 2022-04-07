# Group for all developers in the product team.
resource "google_cloud_identity_group" "product-developers" {
  provider     = google-beta
  display_name = google_project.product_deployment.project_id
  description = "Group to manage 'Client Level' access to the ${google_project.product_deployment.project_id} deployment"

  parent = "customers/${var.ALIS_OS_GOOGLE_CUSTOMER_ID}"

  group_key {
    id = "${google_project.product_deployment.project_id}@identity.${var.ALIS_OS_DOMAIN}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
    "cloudidentity.googleapis.com/groups.security"         = ""
  }
  depends_on = [google_project_service.product-services]
}

# The manager of the group.
resource "google_cloud_identity_group_membership" "product-developer-manager" {
  group = google_cloud_identity_group.product-developers.id

  preferred_member_key {
    id = var.ALIS_OS_OWNER
  }

  # MEMBER role must be specified. The order of roles should not be changed.
  roles { name = "MEMBER" }
  roles { name = "MANAGER" }
}

output "groupid" {
  value       = google_cloud_identity_group.product-developers.group_key[0].id
  description = "ID of the group. For Google-managed entities, the ID is the email address the group"
}