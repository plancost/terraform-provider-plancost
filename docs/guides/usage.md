---
page_title: "Usage"
description: |-
  Provide usage inputs to improve plancost estimates via usage or usage_file.
---

# Usage

Usage inputs improve cost estimation accuracy for usage-based pricing components (for example: storage GB-months, requests, execution counts). Usage can be provided to `plancost_estimate` in one of two ways:

- `usage` (inline dynamic map)
- `usage_file` (YAML file)

You can use `usage_file` to load a base set of usage values (e.g. defaults) and override specific values using the inline `usage` map.

Alternatively, you can [generate a usage file automatically](#generating-a-usage-file) based on your configuration.

## Overview

This guide documents how to supply usage inputs, how to generate a usage file, and what configuration is loaded by the provider.

The repository includes ready-to-use examples under `examples/usage-file/`:

- [usage-example.yml](https://github.com/plancost/terraform-provider-plancost/blob/main/examples/usage-file/usage-example.yml): a comprehensive reference of supported usage keys
- [usage-defaults.small.yml](https://github.com/plancost/terraform-provider-plancost/blob/main/examples/usage-file/usage-defaults.small.yml): a smaller curated set
- [usage-defaults.medium.yml](https://github.com/plancost/terraform-provider-plancost/blob/main/examples/usage-file/usage-defaults.medium.yml): a medium curated set
- [usage-defaults.large.yml](https://github.com/plancost/terraform-provider-plancost/blob/main/examples/usage-file/usage-defaults.large.yml): a larger curated set

## Prerequisites

- Terraform installed.
- The `plancost` provider installed/configured.

## Step 1: Create a module with usage-based resources

Example module skeleton (add resources that have usage-based cost dimensions):

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

provider "plancost" {}
```

## Step 2: Provide usage inline

Inline usage is a map keyed by Terraform resource address. Values are per-resource usage attributes.

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  usage = {
    "azurerm_storage_account.example" = {
      storage_gb = 100
    }

    "azurerm_function_app.example" = {
      monthly_executions    = 100000
      execution_duration_ms = 500
      memory_mb             = 128
      instances             = 1
    }
  }
}
```

Result:

- The provider reads `usage` as a dynamic map and passes it into estimation.

## Step 3: Provide usage via `usage_file`

Create a YAML file (for example `usage.yml`) and reference it from `plancost_estimate`.

`main.tf`:

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  usage_file        = abspath("${path.module}/usage.yml")
}
```

`usage.yml`:

```yaml
version: 0.1
resource_usage:
  azurerm_storage_account.example:
    storage_gb: 100
  azurerm_function_app.example:
    monthly_executions: 100000
    execution_duration_ms: 500
    memory_mb: 128
    instances: 1
```

Result:

- The provider reads the YAML file and loads `resource_usage` into the estimator.

## Generating a Usage File

You can automatically generate a usage file containing all resources in your module and their available usage parameters. This is useful for discovering what usage parameters are supported and creating a baseline for your `usage.yml`.

To generate a usage file, specify the `export_usage_file` attribute in your `plancost_estimate` resource:

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  export_usage_file = abspath("${path.module}/usage.yml")
}
```

When you run `terraform plan` or `terraform apply`, the provider will create the file at the specified path. The generated file will contain the usage schema for all resources found in the module, populated with default values where applicable.

You can then modify the values in the generated file as needed. Once you are satisfied with the usage values, you can switch to using this file for input by changing `export_usage_file` to `usage_file` in your configuration:

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  usage_file        = abspath("${path.module}/usage.yml")
}
```

## Advanced: Usage Precedence

The provider merges usage data from multiple sources. The precedence order (highest priority first) is:

1. **Inline HCL `usage`**: Values defined directly in the `plancost_estimate` resource.
2. **Usage File (`usage.yml`)**:
    a. **Exact resource address**: e.g., `azurerm_storage_account.example["blue"]` in `resource_usage`.
    b. **Wildcard usage keys**: e.g., `azurerm_storage_account.example[*]` in `resource_usage`.
    c. **Resource type defaults**: e.g., `azurerm_storage_account` in `resource_type_default_usage`.

### Resource Type Defaults

You can define default usage for all resources of a specific type in the YAML file using `resource_type_default_usage`.

```yaml
version: 0.1
resource_type_default_usage:
  azurerm_storage_account:
    storage_gb: 50
resource_usage:
  azurerm_storage_account.large_app:
    storage_gb: 500
```

### Overriding File Values with HCL

You can load a usage file and override specific values in your Terraform configuration.

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  usage_file        = abspath("${path.module}/usage.yml")

  # Overrides values from usage.yml
  usage = {
    "azurerm_storage_account.example" = {
      storage_gb = 1000
    }
  }
}
```

### Wildcards (`[*]`)

Wildcards are used heavily in the provider test fixtures (for example [app_configuration_test.usage.yml](https://github.com/plancost/terraform-provider-plancost/blob/main/internal/terraform/azurerm/testdata/app_configuration_test/app_configuration_test.usage.yml)):

```yaml
version: 0.1
resource_usage:
  azurerm_app_configuration.usage[*]:
    monthly_additional_requests: 30000
```

This is intended for cases where the module creates multiple instances (via `count`/`for_each`) and you want one usage definition to apply to all instances.

### Nested usage maps

Some resources are modeled as children of other resources for cost estimation. In these cases, the usage for the child resource is nested under the parent resource's usage key.

For example, `azurerm_backup_protected_vm` costs are aggregated under the `azurerm_recovery_services_vault`. Therefore, you define the VM usage inside the vault's usage block:

```yaml
version: 0.1
resource_usage:
  # Parent resource
  azurerm_recovery_services_vault.example[*]:
    # Child resource nested inside
    azurerm_backup_protected_vm.with_usage[*]:
      disk_utilization_gb: 300
      additional_backup_retention_gb: 100
```
