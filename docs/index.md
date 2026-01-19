---
page_title: "plancost Provider"
description: |-
  The plancost provider allows you to estimate, track, and optimize cloud costs directly within your Terraform workflow.
---

# plancost Provider

**Website:** [https://plancost.io](https://plancost.io)

The `plancost` provider empowers engineering teams to estimate, track, and optimize cloud costs directly within their Terraform workflow. By treating cost as a first-class citizen in your infrastructure code, you gain immediate visibility, enforce budget guardrails, and ensure compliance before a single resource is deployed.

> **Note:** Currently, `plancost` supports the **Azure** provider (`azurerm`). Support for **AWS** and **GCP** is coming soon.

## Why plancost?

- **💰 Instant Cost Estimates:** See the price tag of your infrastructure changes in every plan.
- **🛡️ Budget Guardrails:** Prevent bill shock by setting cost thresholds that fail the build if exceeded.
- **🏷️ Tagging Compliance:** Enforce tagging policies to ensure every resource is properly allocated.
- **🚀 Optimization Recommendations:** Get actionable advice to right-size resources and save money.
- **🤖 CI/CD Integration:** Seamlessly integrate cost checks into your GitHub Actions, GitLab CI, or Jenkins pipelines.

## Get Started

Simply copy-paste this into your Terraform configuration:

```terraform
terraform {
  required_providers {
    plancost = {
      source = "plancost/plancost"
    }
  }
}

# Optional: Configure the plancost provider
provider "plancost" {
  # Optional: Set the API key for the pricing service.
  # This is required when running in CI/CD environments.
  # Can also be set via the PLANCOST_API_KEY environment variable.
  # Get your API key at https://plancost.io
  # api_key = "your-api-key"
}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
}
```

## Example Output

When you run `terraform plan`, `plancost` generates a detailed cost breakdown directly in your terminal, highlighting the cost impact of your changes:

```diff
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

## Documentation & Guides

Ready to get started? Check out our comprehensive guides:

- **[Quickstart](guides/quickstart.md)**: Get your first cost estimate in minutes.
- **[Native Diff View](guides/diff-view.md)**: See cost changes alongside infrastructure changes.
- **[CI/CD Integration](guides/cicd.md)**: Automate cost checks in your pipelines.
- **[Pull Request Comments](guides/post-comment.md)**: Post cost summaries directly to your PRs.
- **[Usage-Based Estimation](guides/usage.md)**: Refine estimates with expected usage data.
- **[Variables](guides/variables.md)**: Learn how to pass Terraform variables.
- **[Tagging Policies](guides/tagging-policy.md)**: Enforce mandatory tags and values.
- **[Cost Guardrails](guides/guardrails.md)**: Set limits on total monthly costs.
- **[OPA Integration](guides/opa.md)**: Enforce advanced cost policies with Open Policy Agent.
- **[Optimization](guides/optimization-recommendations.md)**: Discover recommendations to save.
- **[Security & Privacy](guides/security-and-privacy.md)**: Understand how we protect your data.
- **[Troubleshooting](guides/troubleshooting.md)**: Solutions to common issues.

## Examples

Explore our [examples directory](https://github.com/plancost/terraform-provider-plancost/tree/main/examples) for complete use cases:

- **[Quickstart](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/quickstart)**: Basic usage with Azure Public IP.
- **[Cost Guardrails](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/cost-guardrails)**: Enforce budget limits and cost increase policies.
- **[Tagging Policies](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/tagging_policy)**: Validate required tags and allowed values.
- **[Usage Estimation](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/usage-estimation-with-usage-file)**: Improve accuracy with usage files (also see [Inline Usage](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/usage-estimation-with-inline-usage) and [Wildcards](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/usage-estimation-with-wildcards)).
- **[CI/CD Integration](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/ci-demo)**: GitHub Actions workflow example.
- **[Pull Request Comments](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/pull-request-comment)**: Export estimates for PR comments.
- **[Policy as Code (OPA)](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/policy_as_code)**: Advanced policies using Rego.
- **[Cost Optimization](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/cost-optimaization)**: Get recommendations for cost savings.
- **[Simple Discounting](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/simple_discounting)**: Apply negotiated discounts to estimates.
- **[Variables](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/variables)**: Working with Terraform variables and `.tfvars` files.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_key` (String, Sensitive) The API key for the pricing service. Can also be set via the `PLANCOST_API_KEY` environment variable. This is required when running in CI/CD environments. You can obtain an API key by signing up at [plancost.io](https://plancost.io).

- `api_endpoint` (String) The API endpoint for the pricing service. Defaults to `https://api.plancost.io`. Can also be set via the `PLANCOST_API_ENDPOINT` environment variable.
