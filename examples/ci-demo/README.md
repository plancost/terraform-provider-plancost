# PlanCost CI/CD Integration Demo

This example demonstrates how to integrate the **PlanCost** Terraform provider into a GitHub Actions workflow. It calculates the estimated cost of your Terraform plan and posts the results as a comment on your Pull Request.

## Files

- **`main.tf`**: A sample Terraform configuration creating Azure resources and a `plancost_estimate` resource.
- **`provider.tf`**: Configuration for the `azurerm` and `plancost` providers.
- **`github-action-plancost-post-comment.yml`**: A sample GitHub Action workflow file.

## Setup

To use this in your own repository:

1.  **Configure Secrets**:
    Add the following secrets to your GitHub repository (or Environment):
    - `ARM_CLIENT_ID`
    - `ARM_CLIENT_SECRET`
    - `ARM_SUBSCRIPTION_ID`
    - `ARM_TENANT_ID`
    - `PLANCOST_API_KEY` (Get this from the https://plancost.io/ dashboard after signing up)

2.  **Add the Workflow**:
    Copy `github-action-plancost-post-comment.yml` to `.github/workflows/plancost-estimate.yml` in your repository.

3.  **Configure Environment (Optional)**:
    If you are using GitHub Environments (e.g., `azure` or `production`), ensure the `environment:` field in the workflow matches your setup.

## How it works

1.  The workflow triggers on a Pull Request.
2.  It runs `terraform init` and `terraform plan`.
3.  The `plancost` provider generates an `estimate.md` file during the plan phase.
4.  The workflow uses the `thollander/actions-comment-pull-request` action to post the content of `estimate.md` to the PR.

## Example Output

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
