---
page_title: "plancost_estimate"
subcategory: ""
description: |-
  Cost estimation resource for Terraform modules
---

# plancost_estimate (Resource)

The `plancost_estimate` resource estimates the cost of cloud resources within a Terraform module. It integrates cost estimation, policy enforcement, and optimization recommendations directly into your Terraform workflow.

> **Note:** Currently, `plancost` only supports the **Azure** provider (`azurerm`). Support for **AWS** and **GCP** is planned.

## Example Usage

```terraform
terraform {
  required_providers {
    plancost = {
      source = "plancost/plancost"
    }
  }
}

provider "plancost" {}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  project_name      = "my-terraform-project"
}
```

## Schema

### Required

- `working_directory` (String) Absolute path to the Terraform module directory (e.g., `abspath(path.module)`).

### Optional

- `project_name` (String) The name of the project to create in PlanCost. If not specified, the directory name of `working_directory` will be used, e.g., `my-terraform-project`.

- `var_file` (String) Absolute path to the variables file (e.g., `abspath("${path.module}/variables.tfvars")`). The provider automatically loads variables from the following sources:
  1. The file specified in `var_file`.
  2. The file path specified in the `PLANCOST_VAR_FILE` environment variable.
  3. `terraform.tfvars` in the `working_directory`.
  4. Environment variables starting with `TF_VAR_`.


- `usage_file` (String) Absolute path to the usage file (e.g., `abspath("${path.module}/usage.yml")`). 

- `usage` (Dynamic) Usage data for resources. 

  More details can be found in the [Usage Guide](../guides/usage.md).

  Example:
  ```hcl
  resource "plancost_estimate" "this" {
    working_directory = abspath(path.module)

    usage = {
      "azurerm_storage_account.example" = {
        storage_gb = 100
      }
      "azurerm_linux_function_app.example" = {
        monthly_execution_count = 1000000
        execution_duration_ms   = 200
      }
    }
  }
  ```


- `recommendations_enabled` (Boolean) Enable optimization recommendations. Note: This is a paid feature.

- `export_markdown_file` (String) Absolute path to the output markdown file (e.g., `abspath("${path.module}/estimate.md")`). If specified, the cost estimate report will be written to this file.

- `export_usage_file` (String) Absolute path to the output usage file (e.g., `abspath("${path.module}/usage.yml")`). If specified, the provider will generate a usage file containing the usage schema for all resources in the module. This is useful for discovering available usage parameters and creating a baseline for customization.

- `discount` (Block List) List of discounts to apply. (see [below for nested schema](#nestedblock--discount))

- `guardrail` (Block List) List of guardrail policies to enforce cost limits. Note: This is a paid feature. Free tier users are limited to 1 guardrail and cannot use 'block' actions. (see [below for nested schema](#nestedblock--guardrail))

- `tagging_policy` (Block List) List of tagging policies to enforce. Note: This is a paid feature. (see [below for nested schema](#nestedblock--tagging_policy))



<a id="nestedblock--discount"></a>
### Nested Schema for `discount`

Required:

- `percentage` (Number) The discount percentage (0.0 to 1.0).

Optional:

- `resource_type` (String) The resource type to apply the discount to (e.g., 'azurerm_virtual_machine'). If not specified, applies to all resources.

Example:
```hcl
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  discount {
    percentage    = 0.2
    resource_type = "azurerm_virtual_machine"
  }
}
```



<a id="nestedblock--guardrail"></a>
### Nested Schema for `guardrail`

Required:

- `action` (String) The action to take when the threshold is breached. Valid values: 'warning', 'block'.
- `condition` (String) The condition to trigger the guardrail. Valid values: 'monthly_cost_increase_amount', 'monthly_cost_increase_percentage', 'monthly_cost_budget'.
- `threshold` (Number) The numeric value for the condition (amount or percentage).

Example:
```hcl
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  # Block if the monthly cost budget exceeds $1000
  guardrail {
    condition = "monthly_cost_budget"
    threshold = 1000
    action    = "block"
  }
}
```

**Example Plan Output (Blocked):**
```text
╷
│ Error: Guardrail Violation
│ 
│   with plancost_estimate.this,
│   on main.tf line 5, in resource "plancost_estimate" "this":
│    5:   guardrail {
│ 
│ Monthly cost $1250.00 exceeds budget $1000.00.
╵
```

<a id="nestedblock--tagging_policy"></a>
### Nested Schema for `tagging_policy`

Required:

- `action` (String) The action to take when the policy is violated. Valid values: 'warning', 'block'.

- `key` (String) The tag key that must exist.

Optional:

- `allowed_values` (List of String) List of allowed values for the tag. If specified, the tag value must be one of these values (case-insensitive). Only one of `allowed_values` or `pattern` should be specified.

- `pattern` (String) Regex pattern that the tag value must match. Only one of `allowed_values` or `pattern` should be specified.

- `resource_types` (List of String) List of resource types to apply this policy to (e.g. ['azurerm_resource_group']). If empty, applies to all resources that support tags.

Example:
```hcl
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  # Ensure all resources have an 'Environment' tag with value 'Dev' or 'Prod'
  tagging_policy {
    key            = "Environment"
    allowed_values = ["Dev", "Prod"]
    action         = "block"
  }
}
```

**Example Plan Output (Blocked):**
```text
╷
│ Error: Tagging Policy Violation
│ 
│   with plancost_estimate.this,
│   on main.tf line 5, in resource "plancost_estimate" "this":
│    5:   tagging_policy {
│ 
│ Resource azurerm_resource_group.example (azurerm_resource_group) missing required tag 'Environment'
╵
```

### Read-Only

- `id` (String) Resource identifier.

- `monthly_cost` (Number) The estimated monthly cost (numeric value).

- `recommendations` (List of Object) List of optimization recommendations.

  Structure:
  - `resource_address` (String): The address of the resource.
  - `description` (String): Description of the recommendation.
  - `type` (String): Type of recommendation (e.g., "Reservation", "Advisory").
  - `term` (String): Term length for reservations (e.g., "1 yr", "3 yr").
  - `savings_amount` (Number): Estimated monthly savings.
  - `savings_percentage` (Number): Estimated savings percentage.

  Example:
  ```text
  recommendations = [
    {
      resource_address   = "azurerm_linux_virtual_machine.example"
      description        = "Save $30/mo (58%) on azurerm_linux_virtual_machine.example with a 1-year Reservation"
      type               = "Reservation"
      term               = "1 yr"
      savings_amount     = 30.5
      savings_percentage = 0.58
    },
    {
      resource_address   = "azurerm_linux_virtual_machine.example"
      description        = "Consider using the latest generation version 5 of D series for better performance and cost efficiency"
      type               = "Advisory"
    }
  ]
  ```

- `resources` (Dynamic) Detailed cost breakdown per resource. This is a list of objects containing resource details, cost components, and sub-resources.

  Structure:
  - `name` (String): The name of the resource.
  - `cost_components` (List): List of cost components.
    - `name` (String): Name of the cost component.
    - `monthly_quantity` (String): Monthly quantity.
    - `unit` (String): Unit of measurement.
    - `monthly_cost` (Number): Estimated monthly cost.
  - `sub_resources` (List): List of sub-resources (recursive structure).

  Example:
  ```json
  [
    {
      "name": "azurerm_public_ip.test",
      "cost_components": [
        {
          "name": "IP address (static, regional)",
          "monthly_quantity": "730",
          "unit": "hours",
          "monthly_cost": 3.65
        }
      ],
      "sub_resources": []
    }
  ]
  ```

- `view` (String) The pretty printed output of the estimate.

  Example:
  ```text
  Project: main
  
   Name                                                         Monthly Qty  Unit    Monthly Cost
  
   azurerm_public_ip.test
   └─ IP address (static, regional)                             730          hours   $3.65
  
   OVERALL TOTAL                                                                     $3.65

  💡 Optimization Recommendations

   azurerm_linux_virtual_machine.vm
   └─ 1-Year Reservation: Save $369/mo (58%)
  ```

## Import

Import is not supported for this resource.
