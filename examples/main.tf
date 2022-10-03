terraform {
  required_providers {
    dbt = {
      version = "0.2"
      source  = "hashicorp.com/3lvia/dbt"
    }
  }
}

provider "dbt" {
    service_token = var.service_token
    account_id = var.account_id
}

# resource "dbt_user_group" "myUserGroup" {
#   name = "myNewUserGroupCreatedByTerraform"
#   assign_by_default = true
#   sso_mapping_groups = ["systemaccess-edna-developer"]
#   group_permissions {
#     permission_set = "readonly"
#     all_projects = true
#     }
# }

# resource "dbt_license_map" "myLicensemap" {
#   license_type = "developer"
#   sso_license_mapping_groups = ["systemaccess-kunde-developer"]
# }

# resource "dbt_license_map" "myLicensemap2" {
#   license_type = "developer"
#   sso_license_mapping_groups = ["systemaccess-edna-developer"]
# }

# resource "dbt_license_map" "myLicensemap3" {
#   license_type = "developer"
#   sso_license_mapping_groups = ["systemaccess-edna-developer"]
# }
