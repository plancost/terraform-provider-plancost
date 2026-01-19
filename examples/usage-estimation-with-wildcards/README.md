# Usage Estimation with Wildcards Example

This example demonstrates how to use wildcards in usage definitions.

The `usage.yaml` file uses wildcard patterns (e.g., `azurerm_storage_account.large_app[*]`) to apply usage parameters to multiple resources that match the pattern, such as resources created with `count` or `for_each`.
