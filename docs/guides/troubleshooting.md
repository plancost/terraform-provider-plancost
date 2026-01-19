---
page_title: "Troubleshooting"
description: "Common issues and troubleshooting tips for the plancost provider."
---

# Troubleshooting

## "No changes" after a failed Terraform run

You might encounter a situation where a `terraform plan` or `terraform apply` fails due to errors in other resources, but a subsequent `terraform plan` shows "No changes" for your `plancost_estimate` resource.

### Cause

The `plancost_estimate` resource typically has no dependencies on other resources in your Terraform workspace. This means it is evaluated independently during the plan/apply phase.

Even if other resources fail to plan or apply, the `plancost_estimate` resource may successfully complete its estimation and update its state (or be marked as planned).

### Solution

This is **expected behavior**. Since the `plancost_estimate` resource was successfully processed and there have been no changes to your configuration since then, Terraform correctly reports that there are no pending changes for this resource.

If you need to force the estimate to be regenerated (for example, if you suspect the external pricing data has changed or you want to see the output again), you can taint the resource:

```bash
terraform taint plancost_estimate.this
terraform plan
```

Replace `plancost_estimate.this` with the actual address of your resource.

## Resource recreation when adding an API key

When adding an API key to a configuration that previously ran without one, `terraform plan` indicates that the `plancost_estimate` resource will be created.

### Cause

This is **by design**.

*   **No API Key (Local Mode):** The resource uses a local placeholder ID and is not tracked by the plancost.io platform.
*   **With API Key (Cloud Mode):** The resource requires a valid project ID from the plancost.io platform for tracking and dashboard integration.

Switching modes changes the underlying identity of the resource, causing Terraform to recreate it to establish the correct association with the plancost.io platform.

