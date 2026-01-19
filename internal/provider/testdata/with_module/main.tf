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


module "aks" {
  source = "./modules/aks"
  location = "westus"
  vm_size  = "Standard_D2s_v3"
}

module "avm-res-resources-resourcegroup" {
  source  = "Azure/avm-res-resources-resourcegroup/azurerm"
  version = "0.2.1"
  name     = "example-resources-rg"
  location = "westus"
}

output "o1" {
  value = module.aks
}

output "o2" {
  value = module.avm-res-resources-resourcegroup
}