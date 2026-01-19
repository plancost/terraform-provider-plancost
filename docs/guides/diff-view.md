---
page_title: "Native Diff View"
description: |-
  Learn how plancost integrates with Terraform's native plan output to show cost changes alongside infrastructure changes.
---

# Native Diff View

One of the most powerful features of `plancost` is its seamless integration with Terraform's native CLI output. Because `plancost` operates as a standard Terraform provider, cost estimates are treated as resource attributes.

When you run `terraform plan`, `plancost` calculates the cost impact of your changes and presents them as an in-place update to the `plancost_estimate` resource. This allows you to see the financial impact of your infrastructure changes directly alongside the technical changes.

## How it Works

1.  **Resource Changes**: You modify your Terraform configuration (e.g., changing a VM size, adding a storage account).
2.  **Plan Calculation**: Terraform calculates the changes required for your infrastructure.
3.  **Cost Estimation**: `plancost` analyzes the plan and updates the `plancost_estimate` resource with the new cost data.
4.  **Unified Diff**: The `terraform plan` output shows both the infrastructure changes and the resulting cost changes in a single, unified view.

## Example

Consider a scenario where you are adding Geo-replication to an Azure Container Registry.

The `terraform plan` output will show the changes to the `azurerm_container_registry` resource, and immediately above or below it, the `plancost_estimate` resource will show the corresponding cost increase.

```diff
Terraform will perform the following actions:

  # module.containerregistry.azurerm_container_registry.this will be updated in-place
  ~ resource "azurerm_container_registry" "this" {
        id                            = "/subscriptions/000000/resourceGroups/rg-a190/providers/Microsoft.ContainerRegistry/registries/acra190"
        name                          = "acra190"
        tags                          = {}
        # (19 unchanged attributes hidden)

      + georeplications {
          + location                  = "australiasoutheast"
          + regional_endpoint_enabled = true
          + zone_redundancy_enabled   = true
        }
      + georeplications {
          + location                  = "westus"
          + regional_endpoint_enabled = true
          + zone_redundancy_enabled   = true
        }
    }

  # plancost_estimate.this will be updated in-place
  ~ resource "plancost_estimate" "this" {
        id                   = "plancost-local"
      ~ monthly_cost         = 104 -> 203.99
      ~ resources            = [
          ~ {
              ~ cost_components = [
                    {
                        monthly_cost     = 54
                        monthly_quantity = "540000"
                        name             = "Build vCPU"
                        unit             = "seconds"
                    },
                  + {
                      + monthly_cost     = 100
                      + monthly_quantity = "60"
                      + name             = "Geo replication (2 locations)"
                      + unit             = "days"
                    },
                    {
                        monthly_cost     = 50
                        monthly_quantity = "30"
                        name             = "Registry usage (Premium)"
                        unit             = "days"
                    },
                ]
                name            = "module.containerregistry.azurerm_container_registry.this"
                # (1 unchanged attribute hidden)
            },
        ]
      ~ view                 = <<-EOT
            Project: main
            
             Name                                                        Monthly Qty  Unit           Monthly Cost
            
             module.containerregistry.azurerm_container_registry.this
             ├─ Build vCPU                                                   540,000  seconds              $54.00  *
          +  ├─ Geo replication (2 locations)                                     60  days                $100.00
             └─ Registry usage (Premium)                                          30  days                 $50.00
            
          -  OVERALL TOTAL                                                                                $104.00
          +  OVERALL TOTAL                                                                                $204.00
            
            *Usage costs can be estimated by providing usage data in the plancost_estimate resource.
            
            ──────────────────────────────────
            7 cloud resources were detected:
            ∙ 1 were estimated
            ∙ 5 were free
            ∙ 1 are not supported yet: [modtm_telemetry]
            
            ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
            ┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
            ┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
          - ┃ main                                               ┃        $50.00 ┃      $54.00 ┃    $104.00 ┃
          + ┃ main                                               ┃       $149.99 ┃      $54.00 ┃    $204.00 ┃
            ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛
        EOT
        # (6 unchanged attributes hidden)
    }

Plan: 0 to add, 2 to change, 0 to destroy.
```

In this example:
- The **`azurerm_container_registry`** resource shows the addition of two `georeplications` blocks.
- The **`plancost_estimate`** resource immediately reflects this with a new cost component "Geo replication (2 locations)" costing $100.00.
- The **OVERALL TOTAL** and the summary table show the total monthly cost increasing from $104.00 to $204.00.

This immediate feedback loop ensures that engineers are aware of the cost implications of their changes *before* they apply them.
