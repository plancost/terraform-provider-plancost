---
page_title: "CI/CD"
description: |-
  Running plancost estimates in CI/CD.
---

# CI/CD

This guide describes how to run `plancost_estimate` in CI/CD.

## Key requirement: API key

In CI/CD you must provide an API key.

> **Note:** On the free tier, you can run `plancost` in CI/CD for up to two projects. You can upgrade your plan at https://plancost.io.

## Provider configuration patterns

### Use env vars (recommended)

Prefer storing the key in your CI system secret store and injecting it as an environment variable.

```shell
export PLANCOST_API_KEY="..."
```

```terraform
provider "plancost" {}
```

### Use explicit provider arguments

```terraform
provider "plancost" {
  api_key = var.plancost_api_key
}
```

## Terraform execution notes

- Ensure the CI runner has network access to the pricing endpoint (`https://api.plancost.io`) and project management endpoint (`https://plancost.io`).

## How CI/CD project usage is tracked

The `plancost` service tracks CI/CD usage so it can apply plan limits (for example the free-tier “two projects in CI/CD” limit).

Simply put, the limit is based on how many `plancost_estimate` resources exist across all your projects.

What you need to know:

- A “project” roughly corresponds to one Terraform state running estimates (i.e., the same `plancost_estimate` resource).
- Multiple people (or CI jobs) using the same shared Terraform state still count as one project.
- When `plancost_estimate` resource is destroyed, the project is freed up immediately.
- Please use one `plancost_estimate` resource per Terraform state. Multiple `plancost_estimate` resources in the same Terraform workspace/state count as multiple projects.

If you hit a project limit unexpectedly:

- Check that you’re not creating multiple `plancost_estimate` resources in the same workspace.
- Manage your projects at https://plancost.io/dashboard.
- Upgrade your plan at https://plancost.io if you need more CI/CD projects.
