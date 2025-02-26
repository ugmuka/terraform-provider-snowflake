---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "snowflake_resource_monitor_grant Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# snowflake_resource_monitor_grant (Resource)



## Example Usage

```terraform
resource snowflake_monitor_grant grant {
  monitor_name      = "monitor"
  privilege         = "MODIFY"
  roles             = ["role1"]
  with_grant_option = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **monitor_name** (String) Identifier for the resource monitor; must be unique for your account.

### Optional

- **enable_multiple_grants** (Boolean) When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.
- **id** (String) The ID of this resource.
- **privilege** (String) The privilege to grant on the resource monitor.
- **roles** (Set of String) Grants privilege to these roles.
- **with_grant_option** (Boolean) When this is set to true, allows the recipient role to grant the privileges to other roles.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_resource_monitor_grant.example name
```
