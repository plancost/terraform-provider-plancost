terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0.0"
    }
    plancost = {
      source = "plancost/plancost"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "plancost" {
}