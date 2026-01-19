# Usage Estimation with Defaults Example

This example demonstrates how to set default usage values for resource types.

The `usage.yaml` file uses `resource_type_default_usage` to define default values for all resources of a specific type (e.g., `azurerm_storage_account`). Specific resources can still have their own usage defined in `resource_usage` to override these defaults.
