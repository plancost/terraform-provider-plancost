---
page_title: "Quickstart"
description: |-
  Minimal setup for generating a cost estimate with plancost.
---

# Quickstart

This guide shows the minimal configuration to generate a cost estimate using `plancost_estimate`.

> **Note:** Currently, `plancost` estimates Azure resources via the `azurerm` provider. Support for AWS and GCP is planned.

## Prerequisites

- Terraform installed.
- Network access to `https://api.plancost.io` (or a custom `api_endpoint`).
- A Terraform module to estimate (this can be a module that only contains `resource` blocks; it does not need to be applied).

## Step 1: Configure the provider

```terraform
terraform {
  required_providers {
    plancost = {
      source = "plancost/plancost"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {}
  subscription_id = "your-subscription-id"
}
```

## Step 2: Add a small Azure resource to estimate

For a minimal, predictable example, add a Public IP.

```terraform
resource "azurerm_resource_group" "test" {
  name     = "plancost-quickstart"
  location = "eastus"
}

resource "azurerm_public_ip" "test" {
  name                = "plancost-quickstart"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}
```

## Step 3: Add a `plancost_estimate`

Add a `plancost_estimate` resource to your Terraform configuration.

Point `working_directory` at the directory containing the Terraform configuration you want to estimate. Using `abspath(path.module)` sets it to the current module directory (the directory Terraform is evaluating).

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
}
```

## Step 4: Run a plan

```shell
terraform init
terraform plan
```

Example plan output:

```text
Terraform will perform the following actions:

  # azurerm_public_ip.test will be created
  + resource "azurerm_public_ip" "test" {
      ...
    }

  # plancost_estimate.this will be created
  + resource "plancost_estimate" "this" {
      + monthly_cost = 3.65
      + resources    = {
          + "azurerm_public_ip.test" = {
              + cost_components = {
                  + "IP address (static, regional)" = {
                      + monthly_cost     = 3.65
                      + monthly_quantity = "730"
                      + unit             = "hours"
                    }
                }
              + sub_resources   = {}
            }
        }
      + view = <<-EOT
        Project: main

         Name                                                         Monthly Qty  Unit    Monthly Cost

         azurerm_public_ip.test
         └─ IP address (static, regional)                             730          hours   $3.65

         OVERALL TOTAL                                                                     $3.65
    EOT
    }

Plan: 3 to add, 0 to change, 0 to destroy.
```

## Step 5: Change the Public IP SKU and re-plan

This step assumes the initial configuration has been successfully applied at least once, so the Public IP exists in state.

Change the Public IP SKU from `Standard` to `Basic`:

```terraform
resource "azurerm_public_ip" "test" {
  name                = "plancost-quickstart"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Basic"
}
```

Run `terraform plan` again.

Example plan output (cost diff):

```text
  # azurerm_public_ip.test will be updated in-place
  ~ resource "azurerm_public_ip" "test" {
      ~ sku = "Standard" -> "Basic"
      # (other unchanged attributes hidden)
    }

  # plancost_estimate.this will be updated in-place
  ~ resource "plancost_estimate" "this" {
      ~ monthly_cost = 3.65 -> 2.63
      ~ view = <<-EOT
        Project: main

         Name                                                         Monthly Qty  Unit    Monthly Cost

         azurerm_public_ip.test
      -  └─ IP address (static, regional)                             730          hours   $3.65
      +  └─ IP address (static, regional)                             730          hours   $2.63

      -  OVERALL TOTAL                                                                     $3.65
      +  OVERALL TOTAL                                                                     $2.63
    EOT
  }
```

## Estimating a workspace that uses Terraform modules

No extra steps are required if your Terraform workspace is composed of modules.

Place a single `plancost_estimate` in your root module and point `working_directory` at the root module directory. The estimate will include resources declared in child modules referenced by that root.

Example structure:

```text
.
├── main.tf          # root module
└── modules/
    └── networking/
        └── main.tf  # child module
```

Root module (`main.tf`):

```terraform
module "networking" {
  source   = "./modules/networking"
  location = "eastus"
}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
}
```

## Outputting the estimate

To output the estimated cost and view, add output blocks for `plancost_estimate` attributes.

```terraform
output "monthly_cost" {
  value = plancost_estimate.this.monthly_cost
}

output "estimate_view" {
  value = plancost_estimate.this.view
}
```

Run the plan:

```bash
terraform init
terraform plan
```

The plan output will show the estimated cost and a pretty-printed view of the estimate.

## Inspecting the state

You can inspect the calculated estimate details directly from the Terraform state:

```shell
terraform state show plancost_estimate.this
```

Example output:

```
# plancost_estimate.this:
resource "plancost_estimate" "this" {
    id                   = "plancost-local"
    working_directory    = (write-only attribute)
    monthly_cost         = 215.42
    export_markdown_file = (write-only attribute)
    recommendations      = []
    resources            = [
        {
            cost_components = []
            name            = "azurerm_kubernetes_cluster.example"
            sub_resources   = [
                {
                    cost_components = [
                        {
                            monthly_cost     = 0
                            monthly_quantity = "0"
                            name             = "Data processed"
                            unit             = "usage_based"
                        },
                    ]
                    name            = "Load Balancer"
                    sub_resources   = []
                },
                {
                    cost_components = [
                        {
                            monthly_cost     = 197.76
                            monthly_quantity = "2190"
                            name             = "Instance usage (Linux, pay as you go, Standard_A2_v2)"
                            unit             = "hours"
                        },
                    ]
                    name            = "default_node_pool"
                    sub_resources   = [
                        {
                            cost_components = [
                                {
                                    monthly_cost     = 17.66
                                    monthly_quantity = "3"
                                    name             = "Storage (S10, LRS)"
                                    unit             = "months"
                                },
                            ]
                            name            = "os_disk"
                            sub_resources   = []
                        },
                    ]
                },
            ]
        },
    ]
    usage_file           = (write-only attribute)
    var_file             = (write-only attribute)
    view                 = <<-EOT
        Project: main
        
         Name                                                         Monthly Qty  Unit         Monthly Cost  
        
         azurerm_kubernetes_cluster.example
         ├─ Load Balancer
         │  └─ Data processed                                                    0 GB                    $0.00
         └─ default_node_pool
            ├─ Instance usage (Linux, pay as you go, Standard_A2_v2)          2190 hours               $197.76
            └─ os_disk
               └─ Storage (S10, LRS)                                             3 months               $17.66
        
         OVERALL TOTAL                                                                         $215.42
        
        *Usage costs can be estimated by providing usage data in the plancost_estimate resource.
        
        ──────────────────────────────────
        1 cloud resources were detected:
        ∙ 1 was estimated
        ∙ 0 were free
        
        ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
        ┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
        ┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
        ┃ main                                               ┃          $215 ┃           - ┃       $215 ┃
        ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛
    EOT
}
```


