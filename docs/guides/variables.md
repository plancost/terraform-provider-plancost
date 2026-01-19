---
page_title: "Variables"
description: |-
  Learn how to pass Terraform variables to plancost for accurate cost estimation.
---

# Variables

When your Terraform configuration uses input variables, `plancost` needs to know their values to accurately evaluate your code and estimate costs. If a resource attribute (like `sku` or `instance_type`) depends on a variable, and that variable is not defined, `plancost` may not be able to determine the cost.

The `plancost` provider supports multiple ways to pass variables, following a specific precedence order.

## Precedence Order

`plancost` loads variables from the following sources, in order of precedence (highest to lowest):

1.  **`var_file` argument**: Explicitly defined in the `plancost_estimate` resource.
2.  **`PLANCOST_VAR_FILE` environment variable**: Path to a variables file set in the environment.
3.  **`terraform.tfvars`**: Automatically loaded if present in the `working_directory`.
4.  **`TF_VAR_` environment variables**: Standard Terraform environment variables.

## 1. Using `var_file`

You can specify the path to a `.tfvars` file directly in the `plancost_estimate` resource using the `var_file` argument. This is the most explicit method.

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  var_file          = abspath("${path.module}/prod.tfvars")
}
```

## 2. Using `PLANCOST_VAR_FILE`

If you want to set the variables file path dynamically without modifying your Terraform code (e.g., in a CI/CD pipeline), you can use the `PLANCOST_VAR_FILE` environment variable.

```bash
export PLANCOST_VAR_FILE="/path/to/your/prod.tfvars"
terraform plan
```

## 3. Using `terraform.tfvars`

If you have a `terraform.tfvars` file in the same directory as your module (specified by `working_directory`), `plancost` will automatically load it. This mimics standard Terraform behavior.

```terraform
# terraform.tfvars
instance_type = "Standard_D2s_v3"
environment   = "production"
```

## 4. Using `TF_VAR_` Environment Variables

`plancost` also respects standard Terraform environment variables starting with `TF_VAR_`. This is useful for passing individual variable values.

```bash
export TF_VAR_instance_type="Standard_D2s_v3"
export TF_VAR_environment="production"
terraform plan
```

## Troubleshooting

If you see unexpected costs or "unknown" values in your estimate, check the following:

*   Ensure your variables are defined in one of the supported sources.
*   Verify the path to your `.tfvars` file is correct (using `abspath` is recommended).
*   Check the precedence order to ensure a lower-priority source isn't being overridden unexpectedly.
