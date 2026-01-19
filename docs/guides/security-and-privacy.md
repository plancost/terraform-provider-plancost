---
page_title: "Security and Privacy"
description: "How the plancost provider handles your data and ensures privacy."
---

# Security and Privacy

Security is of paramount importance to us. This page outlines how the `plancost` provider handles your data and ensures privacy.

## What data is sent to the Pricing API?

**No cloud credentials, secrets, or sensitive data are sent to the Pricing API.**

The `plancost` provider parses your Terraform configuration to extract only the parameters needed to determine the cost of a resource. For example, to estimate the cost of an Azure Virtual Machine, we send:

- **Region** (e.g., `eastus`)
- **Instance Type** (e.g., `Standard_D2s_v3`)
- **Operating System** (e.g., `Linux`)

We do **not** send:
- Your Terraform state file.
- Your full Terraform plan.
- Any variable values that are not directly related to cost attributes.
- Cloud credentials (access keys, service principal secrets, etc.).
- Resource names or identifiers (unless they are required for pricing, which is rare).

### Example Request

Here is a simplified example of the data sent to the API for a pricing lookup:

```json
{
  "vendor": "azure",
  "service": "Virtual Machines",
  "region": "eastus",
  "attributes": {
    "instanceType": "Standard_D2s_v3",
    "os": "Linux"
  }
}
```

## Does plancost need cloud credentials?

**No.** The `plancost` provider works by analyzing your Terraform configuration and plan. It does not need to authenticate with your cloud provider (Azure, AWS, GCP) to generate an estimate. It uses the `plancost` API to retrieve public pricing data.

## Do you sell my data?

**No.** We do not sell your data.

## Network Requirements

The `plancost` provider needs to communicate with the `plancost` API to fetch pricing data and track project usage.

- **Hostname:** `api.plancost.io`
- **Hostname:** `plancost.io`
- **Port:** `443` (HTTPS)

If your environment has strict network policies, please allow outbound traffic to these hosts.
