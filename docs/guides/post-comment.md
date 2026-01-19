---
page_title: "Post Estimate Comment"
description: |-
  How to post the cost estimate as a comment on a Pull Request.
---

# Post Estimate Comment

This guide explains how to post the generated cost estimate as a comment on a Pull Request (PR) or Merge Request (MR) in your CI/CD pipeline.

While this guide uses **GitHub Actions** as a primary example, the core concept is simple and can be easily adapted to **GitLab CI**, **Jenkins**, **Azure DevOps**, or any other CI/CD system.


## Example Comment

When the workflow runs, it will post a comment on the Pull Request similar to the following:

### PlanCost Report
💰 Monthly cost will increase by $205.00 (67%).
| Name | Monthly Qty | Unit | Monthly Cost |
|:--- |:--- |:--- |:--- |
| + **azurerm_kubernetes_cluster.aks_new** | | | |
| &nbsp;&nbsp;&nbsp;&nbsp;+ Uptime SLA | 730 | hours | $73.00 |
| &nbsp;&nbsp;&nbsp;&nbsp;+ default_node_pool | | | |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+ Instance usage | 730 | hours | $150.00 |
| ~ **azurerm_kubernetes_cluster.aks_mod** | | | |
| &nbsp;&nbsp;&nbsp;&nbsp;  Uptime SLA | 730 | hours | $73.00 |
| &nbsp;&nbsp;&nbsp;&nbsp;~ default_node_pool | | | |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;~ Instance usage | 730 -> 1460 | hours | $100.00 -> $200.00 (+$100.00) |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+ Load Balancer | 1 | months | $5.00 |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  os_disk | | | |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  Storage | 128 | GB | $10.00 |
| - ~~azurerm_kubernetes_cluster.aks_old~~ | | | |
| &nbsp;&nbsp;&nbsp;&nbsp;- ~~Uptime SLA~~ | | | ~~$73.00~~ |
| &nbsp;&nbsp;&nbsp;&nbsp;- ~~default_node_pool~~ | | | |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- ~~Instance usage~~ | | | ~~$50.00~~ |
| **Total** | | | **$306.00 -> $511.00 (+$205.00, 67%)** |


## How it Works

The integration relies on a simple file-based workflow:

1.  **Generate**: The `plancost` provider generates a Markdown file (e.g., `estimate.md`) during the `terraform plan` phase.
2.  **Read**: Your CI pipeline reads the content of this file.
3.  **Post**: Your CI pipeline uses a script or plugin to post that content as a comment to your code review system.

## Configuration

First, configure your `plancost_estimate` resource to output the estimate to a file.

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  # This path is relative to where you run terraform plan
  export_markdown_file = abspath("${path.module}/estimate.md")
}
```

## GitHub Actions Example

In GitHub Actions, you can use a community action like `thollander/actions-comment-pull-request` to handle the commenting logic.

### Prerequisites

1.  **API Key**: Ensure `PLANCOST_API_KEY` is set in your repository secrets.
2.  **Permissions**: The workflow needs write permissions to pull-requests.

### Workflow Configuration

```yaml
name: PlanCost Estimate

on:
  pull_request:

permissions:
  contents: read
  pull-requests: write # Required to post comments

jobs:
  estimate:
    runs-on: ubuntu-latest
    env:
      PLANCOST_API_KEY: ${{ secrets.PLANCOST_API_KEY }}
      # Add other provider credentials (e.g., ARM_CLIENT_ID) here

    steps:
      - uses: actions/checkout@v4

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3

      - name: Terraform Init
        run: terraform init

      - name: Terraform Plan
        # The estimate.md file is created during this step
        run: terraform plan

      - name: Comment PR
        uses: thollander/actions-comment-pull-request@v2
        with:
          filePath: estimate.md
          comment_tag: plancost_estimate # Updates existing comment instead of creating new ones
```

For a complete working example, see the [examples/ci-demo](https://github.com/plancost/terraform-provider-plancost/tree/main/examples/ci-demo) directory in the repository.

## Other CI/CD Systems

For other systems like GitLab CI, Jenkins, or Azure DevOps, the process is similar:

1.  **Run Terraform Plan**: Execute `terraform plan`.
2.  **Script the Comment**: Use `curl` or a CLI tool to post the content of `estimate.md` to your version control system's API.
