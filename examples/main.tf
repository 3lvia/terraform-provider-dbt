terraform {
  required_providers {
    dbt = {
      version = "0.2"
      source  = "hashicorp.com/edu/dbt"
    }
  }
}

provider "dbt" {
    service_token = var.service_token
    account_id = var.account_id
}

resource "dbt_group" "myGroup" {
  name = "myNewGroupCreatedByTerraform123"
  assign_by_default = true
  sso_mapping_groups = ["systemaccess-edna-developer"]
  group_permissions {
        permission_set = "readonly"
        project_id = 69915
        all_projects = false
        }
}
