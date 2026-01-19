---
page_title: "Tagging Policy"
description: |-
  Enforce required tags and tag value rules using plancost tagging_policy.
---

# Tagging Policy

The `tagging_policy` feature lets you enforce tag requirements (presence + allowed values / regex patterns) as part of a `plancost_estimate` run.

You configure tagging policies inside the `plancost_estimate` resource. During `terraform plan`/`apply`, the provider parses the module at `working_directory` and evaluates each policy against resources it can read tags for.

> **Note:** Tagging policy enforcement is a **paid feature**. If the account is not on the paid tier, the provider emits **Tagging Policies Disabled** and skips enforcement. You can upgrade your plan at https://plancost.io.

## Overview

This guide documents how to configure `tagging_policy` and what results to expect from `terraform plan` / `terraform apply`.

## Prerequisites

- Terraform installed.
- AzureRM provider available.
- The `plancost` provider installed/configured.

If you run in CI, set the API key via `PLANCOST_API_KEY` (or `api_key` in the provider block).

## Step 1: Create a module to evaluate

Create a folder (for example `tagging-policy-demo/`) with a `main.tf` like this:

```terraform
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0"
    }
    plancost = {
      source = "plancost/plancost"
    }
  }
}

provider "azurerm" {
  features {}
  subscription_id = "your-subscription-id"
}

provider "plancost" {
  # Prefer env vars in real usage:
  # export PLANCOST_API_KEY=...
  # export PLANCOST_API_ENDPOINT=https://api.plancost.io
  # api_key = "..."
}

resource "azurerm_resource_group" "example" {
  name     = "exampleRG1"
  location = "eastus"

  tags = {
    Environment = "Dev"
    Owner       = "platform-team"
  }
}

resource "azurerm_public_ip" "example" {
  name                = "example-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"

  # No tags defined (used to demonstrate policy violations).
}
```

Result: one tagged resource and one untagged resource.

## Step 2: Add a baseline `plancost_estimate`

Add this resource to `main.tf`:

```terraform
resource "plancost_estimate" "this" {
  # Terraform module path to parse and evaluate.
  working_directory = abspath(path.module)
}
```

Run:

```bash
terraform init
terraform plan
```

Result:

- You should see the `plancost_estimate` resource planned.
- No tagging policy is enforced yet.

## Step 3: Enforce a required tag key

Configure a policy that requires an `Environment` tag key.

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  tagging_policy {
    key    = "Environment"
    action = "warning"
  }
}
```

Run:

```bash
terraform plan
```

Result:

- When `action = "warning"`, the diagnostic is emitted as a warning and the plan can continue.

Example plan output:

```text
╷
│ Warning: Tagging Policy Violation
│
│   with plancost_estimate.this,
│   on main.tf line 0, in resource "plancost_estimate" "this":
│
│ Resource azurerm_public_ip.example (azurerm_public_ip) missing required tag 'Environment'
╵
```

## Step 4: Restrict allowed values

Restrict the allowed values for the tag.

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  tagging_policy {
    key            = "Environment"
    allowed_values = ["Dev", "Prod"]
    action         = "block"
  }
}
```

Run:

```bash
terraform plan
```

Result:

- Missing tag keys become a policy violation.
- Values not in the allow-list become a policy violation.
- With `action = "block"`, violations appear as an **error** and the plan fails.

Example plan output (blocked):

```text
╷
│ Error: Tagging Policy Violation
│
│   with plancost_estimate.this,
│   on main.tf line 0, in resource "plancost_estimate" "this":
│
│ Resource azurerm_public_ip.example (azurerm_public_ip) missing required tag 'Environment'
╵
```

## Step 5: Enforce a value format using regex

Use `pattern` to validate formatting (regex must be valid):

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  tagging_policy {
    key = "Owner"
    # Example: basic email format (local@domain).
    pattern = "^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$"
    action  = "warning"
  }
}
```

Result:

- Values that do not match the regex produce a **Tagging Policy Violation**.

Example plan output (regex mismatch):

```text
╷
│ Warning: Tagging Policy Violation
│
│   with plancost_estimate.this,
│   on main.tf line 0, in resource "plancost_estimate" "this":
│
│ Resource azurerm_resource_group.example (azurerm_resource_group) tag 'Owner' value 'platform-team' does not match pattern '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$'
╵
```

## Step 6: Scope the policy to specific resource types

Scope enforcement to a specific set of Terraform resource types by setting `resource_types`.

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  tagging_policy {
    key            = "Environment"
    allowed_values = ["Dev", "Prod"]
    resource_types = ["azurerm_resource_group"]
    action         = "block"
  }
}
```

Result:

- Only resources whose Terraform type matches entries in `resource_types` are evaluated.
- This is useful for incremental rollout and for avoiding enforcement on resources you do not want to evaluate.

## Configuration reference

The `tagging_policy` block supports:

- `key` (required): Tag key that must exist.
- `allowed_values` (optional): Allowed values (case-insensitive compare).
- `pattern` (optional): Regex that tag value must match.
- `resource_types` (optional): List of Terraform resource types the policy applies to. If omitted/empty, policy applies to resources that support tags.
- `action` (required): `warning` or `block`.

## Tips

- Start with `action = "warning"`, then switch to `block`.
- Keep policies few and explicit (for example: `Environment`, `Owner`, `CostCenter`).
- Use `resource_types` when violations are expected for unrelated resources.
