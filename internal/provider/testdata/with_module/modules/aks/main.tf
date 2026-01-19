terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

variable "location" {
  type    = string
}

variable "vm_size" {
  type    = string
}

module "avm-res-resources-resourcegroup" {
  source  = "Azure/avm-res-resources-resourcegroup/azurerm"
  version = "0.2.1"
  name     = "example-resources-rg"
  location = var.location
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks"
  location            = var.location
  resource_group_name = module.avm-res-resources-resourcegroup.name
  dns_prefix          = "exampleaks"

  default_node_pool {
    name       = "default"
    node_count = 3
    vm_size    = var.vm_size
  }

  identity {
    type = "SystemAssigned"
  }
}
