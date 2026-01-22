# Terraform Provider: plancost

[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](LICENSE)


The `plancost` provider allows you to estimate, track, and optimize cloud costs directly within your Terraform workflow. By integrating cost estimation as a native Terraform resource, you get cost visibility, policy enforcement, and security without leaving your infrastructure code.

**Website:** [https://plancost.io](https://plancost.io)

<p align="center">
  <img src="./images/terraform-plan.gif" alt="plancost demo" width="98%" style="border: 1px solid #d0d7de; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.1);">
</p>

## Overview

Unlike CLI wrappers, `plancost` works as a genuine Terraform Provider. This means cost data is stored in your state file, accessible to other resources, and verifiable by standard policy tools.

## Prerequisites

- Terraform >= 1.11.x
- `plancost` provider installed (or configured via development overrides)

## Quick Start

### 1. Setup

Install the `plancost` provider by adding it to your Terraform configuration.

```hcl
terraform {
  required_providers {
    plancost = {
      source = "plancost/plancost"
    }
  }
}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
}
```

### 2. Run Terraform Plan

Execute `terraform plan`. You will see the `plancost_estimate` resource being created (or updated) and a formatted cost breakdown in the outputs.

```sh
terraform plan
```

**Example Output:**

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

### 3. Track Cost Changes

When you modify your infrastructure (e.g., changing a SKU), `terraform plan` will show the cost difference in the `plancost_estimate` resource and the output.

**Example: Changing Public IP SKU from Standard to Basic**

```text
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

## Why plancost?

| Feature | Others | plancost |
| :--- | :--- | :--- |
| **Setup** | Binary Install + Auth + CI Scripting | **Copy-paste 6 lines of HCL. No Sign-up.** |
| **Workflow** | Separate CLI command | **Native `terraform plan` integration** |
| **Diff View** | Multi-step process / External Dashboard | **Inline in `terraform plan` output** |
| **Cost Guardrails** | External Policies (OPA/Sentinel) | **Defined in HCL (Native Resource)** |
| **Tagging Policies** | External Dashboards | **Defined in HCL (Native Resource)** |
| **Platform Support** | Custom hooks for Atlantis/TFC | **Works out-of-the-box everywhere** |
| **Policy Engines** | Requires custom input parsing | **Native HCL Resource** |
| **CI/CD** | Requires CLI Installation | **Zero-Install (Standard `terraform init`)** |
| **State** | Locked in SaaS Dashboard | **Stored in `terraform.tfstate`** |
| **Privacy** | Sends Configs, Names & Paths to SaaS | **Anonymous (SKU/Region only)** |


## Supported Providers

Currently, `plancost` focuses on **Microsoft Azure**, providing deep coverage for the top 100+ most common resources (matching and exceeding other tools).

For a complete list of supported resources and their pricing status, see the [Supported Resources](docs/guides/supported-resources.md) guide.

*   **Azure**: ✅ Full Support (500+ resources)
*   **AWS**: 🚧 Planned
*   **GCP**: 🚧 Planned


## Examples

Explore our [examples directory](examples/) for complete use cases:

- **[Quickstart](examples/quickstart)**: Basic usage with Azure Public IP.
- **[Cost Guardrails](examples/cost-guardrails)**: Enforce budget limits and cost increase policies.
- **[Tagging Policies](examples/tagging_policy)**: Validate required tags and allowed values.
- **[Usage Estimation](examples/usage-estimation-with-usage-file)**: Improve accuracy with usage files (also see [Inline Usage](examples/usage-estimation-with-inline-usage) and [Wildcards](examples/usage-estimation-with-wildcards)).
- **[CI/CD Integration](examples/ci-demo)**: GitHub Actions workflow example.
- **[Pull Request Comments](examples/pull-request-comment)**: Export estimates for PR comments.
- **[Policy as Code (OPA)](examples/policy_as_code)**: Advanced policies using Rego.
- **[Cost Optimization](examples/cost-optimaization)**: Get recommendations for cost savings.
- **[Simple Discounting](examples/simple_discounting)**: Apply negotiated discounts to estimates.
- **[Variables](examples/variables)**: Working with Terraform variables and `.tfvars` files.


## Credits

This project uses source code from [Infracost](https://github.com/infracost/infracost). We would like to thank the Infracost team for their amazing work in the cloud cost estimation space.

## Contributing 

By submitting a Pull Request, you agree that your contributions are licensed under the Mozilla Public License 2.0.

## License

This project is licensed under the **Mozilla Public License 2.0 (MPL-2.0)**. 

See the [LICENSE](LICENSE) file for specific terms and conditions.

---

This project also includes code from [Infracost](https://github.com/infracost/infracost), which is licensed under the **Apache License 2.0**. See [LICENSE_INFRACOST](LICENSE_INFRACOST) and [NOTICE](NOTICE) for details.