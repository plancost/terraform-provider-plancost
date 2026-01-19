---
page_title: "Guardrails"
description: |-
  Enforce cost thresholds during plancost estimation using guardrail blocks.
---

# Guardrails

Guardrails enforce cost thresholds as part of a `plancost_estimate` evaluation. The provider computes a monthly cost estimate for the module at `working_directory` and evaluates each configured `guardrail` against:

- estimated total monthly cost
- (optionally) a previous baseline cost, when available

> **Note:** Guardrails enforcement is a **paid feature**. Free tier behavior is limited: only the first guardrail is evaluated, and `action = "block"` is treated as a warning. You can upgrade your plan at https://plancost.io.

## Overview

This guide documents how to configure `guardrail` blocks and what results to expect from `terraform plan` / `terraform apply`.

## Prerequisites

- Terraform installed.
- AzureRM provider available.
- The `plancost` provider installed/configured.

If you run in CI, set the API key via `PLANCOST_API_KEY` (or `api_key` in the provider block).

## Step 1: Create a module to evaluate

Create a folder (for example `guardrails-demo/`) with a `main.tf` like this:

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
}

resource "azurerm_public_ip" "example" {
  name                = "example-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}
```

## Step 2: Add a baseline `plancost_estimate`

Add this resource to `main.tf`:

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
}
```

Run:

```bash
terraform init
terraform plan
```

Result:

- The `plancost_estimate` resource is planned.
- No guardrails are enforced.

## Step 3: Add a budget guardrail

Add a monthly budget guardrail:

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  guardrail {
    condition = "monthly_cost_budget"
    threshold = 1000
    action    = "warning"
  }
}
```

Result:

- When cost exceeds the threshold, a **Guardrail Violation** diagnostic is emitted.
- With `action = "warning"`, the plan can continue.

Example plan output (warning):

```text
╷
│ Warning: Guardrail Violation
│ 
│   with plancost_estimate.this,
│   on main.tf line 0, in resource "plancost_estimate" "this":
│ 
│ Monthly cost $1250.00 exceeds budget $1000.00.
╵
```

## Step 4: Block when budget is exceeded

Change `action` to `block` to fail the plan on violation:

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  guardrail {
    condition = "monthly_cost_budget"
    threshold = 1000
    action    = "block"
  }
}
```

Example plan output (blocked):

```text
╷
│ Error: Guardrail Violation
│ 
│   with plancost_estimate.this,
│   on main.tf line 0, in resource "plancost_estimate" "this":
│ 
│ Monthly cost $1250.00 exceeds budget $1000.00.
╵
```

## Step 5: Guardrail on monthly increase

Guardrails can be configured for *increase amount* and *increase percentage*:

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  guardrail {
    # Warn when estimated cost increases by more than $200.
    condition = "monthly_cost_increase_amount"
    threshold = 200
    action    = "warning"
  }

  guardrail {
    # Warn when estimated cost increases by more than 10%.
    condition = "monthly_cost_increase_percentage"
    threshold = 10
    action    = "warning"
  }
}
```

Result:

- The provider compares the current estimated cost against a previous baseline (when available) and emits **Guardrail Violation** diagnostics when thresholds are breached.

## Configuration reference

The `guardrail` block supports:

- `condition` (required): One of:
  - `monthly_cost_increase_amount`
  - `monthly_cost_increase_percentage`
  - `monthly_cost_budget`
- `threshold` (required): Numeric threshold for the condition.
- `action` (required): `warning` or `block`.

## Tips

- Use `action = "warning"` for rollout; switch to `block` for enforcement.
- Keep guardrails scoped to a small set of actionable signals (budget and one increase threshold).
