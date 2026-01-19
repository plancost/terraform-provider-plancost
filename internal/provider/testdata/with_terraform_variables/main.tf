provider "azurerm" {
  features {}
}

variable "resource_group_name" {
  type    = string
}

variable "location" {
  type    = string
}

variable "public_ip_allocation_method" {
  type    = string
}

variable "public_ip_sku" {
  type    = string
}

resource "azurerm_resource_group" "example" {
  name     = var.resource_group_name
  location = var.location
}

resource "azurerm_public_ip" "example" {
  name                = "example-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = var.public_ip_allocation_method
  sku                 = var.public_ip_sku
}