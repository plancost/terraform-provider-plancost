# Usage Estimation with Inline Overrides Example

This example demonstrates how to combine a usage file with inline overrides in the `plancost_estimate` resource.

It loads base usage definitions from `usage.yaml` and then overrides specific parameters (e.g., `storage_gb` for `azurerm_storage_account.example`) directly within the `usage` block in `main.tf`. This allows for flexible usage configuration.
