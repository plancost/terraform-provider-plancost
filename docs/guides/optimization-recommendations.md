---
page_title: "Optimization Recommendations"
description: |-
  Generate cost optimization recommendations during plancost estimation.
---

# Optimization Recommendations

Optimization recommendations are generated as part of a `plancost_estimate` evaluation when `recommendations_enabled = true`.

The provider parses the Terraform module at `working_directory`, estimates costs, and analyzes supported resources for optimization recommendations. Recommendations are returned as a list of strings in the `recommendations` attribute.

> **Note:** Optimization recommendations are a **paid feature**. When the account is not on the paid tier, the provider returns a single recommendation indicating the feature is paid. You can upgrade your plan at https://plancost.io.

## Overview

This guide documents how to enable optimization recommendations and what results to expect from `terraform plan` / `terraform apply`.

## Prerequisites

- Terraform installed.
- AzureRM provider available.
- The `plancost` provider installed/configured.

If you run in CI, set the API key via `PLANCOST_API_KEY` (or `api_key` in the provider block).

## Step 1: Create a module to evaluate

Create a folder (for example `optimization-demo/`) with a `main.tf` like this:

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
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "example" {
  name                  = "example-vm"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  network_interface_ids = [azurerm_network_interface.example.id]
  vm_size               = "Standard_DS1_v2"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
    os_type           = "Windows"
  }

  os_profile {
    computer_name  = "hostname"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {}
}
```

## Step 2: Enable recommendations in `plancost_estimate`

Add a `plancost_estimate` resource:

```terraform
resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  recommendations_enabled = true
}
```

Run:

```bash
terraform init
terraform plan
```

Result:

- The plan includes a `plancost_estimate` resource.
- The `recommendations` attribute is populated with provider-generated strings.
- The `view` attribute includes the optimization recommendations section.

  Example:
  ```text
  💡 Optimization Recommendations

   azurerm_virtual_machine.example
   └─ 1-Year Reservation: Save $369/mo (58%)
  ```

Example plan output (recommendations attribute excerpt):

```text
recommendations = [
  {
    resource_address   = "azurerm_virtual_machine.example"
    description        = "Consider using the latest generation version 5 of DS series for better performance and cost efficiency"
    type               = "Advisory"
  },
  {
    resource_address   = "azurerm_virtual_machine.example"
    description        = "Consider using Azure Hybrid Benefit for Windows VMs if you have eligible on-premises licenses"
    type               = "Advisory"
  },
  {
    resource_address   = "azurerm_virtual_machine.example"
    description        = "Save $369/mo (58%) on azurerm_virtual_machine.example with a 1-year Reservation"
    type               = "Reservation"
    term               = "1 yr"
    savings_amount     = 369.48
    savings_percentage = 0.58
  }
]
```

## Configuration reference

- `recommendations_enabled` (Boolean): Enables recommendation generation.
- `recommendations` (List of Object, read-only): List of optimization recommendations.
  - `resource_address` (String): The address of the resource.
  - `description` (String): Description of the recommendation.
  - `type` (String): Type of recommendation (e.g., "Reservation", "Advisory").
  - `term` (String): Term length for reservations (e.g., "1 yr", "3 yr").
  - `savings_amount` (Number): Estimated monthly savings.
  - `savings_percentage` (Number): Estimated savings percentage.

## Notes

- Recommendations are best-effort and depend on supported resource types and available metadata.
- Use `recommendations` as input for review/approval workflows (e.g., PR checks) rather than automatic changes.
