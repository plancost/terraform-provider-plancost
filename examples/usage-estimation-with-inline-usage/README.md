# Usage Estimation with Inline Usage Example

This example demonstrates how to define resource usage directly within the Terraform code.

Instead of using an external file, the usage parameters (like `storage_gb`, `monthly_executions`) are defined inside the `usage` block of the `plancost_estimate` resource. This keeps the usage logic close to the infrastructure definition.
