provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "eastus"
}

resource "azurerm_dashboard_grafana" "standard" {
  name                  = "grafana-std"
  resource_group_name   = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  grafana_major_version = "10"
}

resource "azurerm_dashboard_grafana" "standard_redundant" {
  name                    = "grafana-std-red"
  resource_group_name     = azurerm_resource_group.example.name
  location                = azurerm_resource_group.example.location
  zone_redundancy_enabled = true
  grafana_major_version   = "10"
}

resource "azurerm_dashboard_grafana" "essential" {
  name                  = "grafana-ess"
  resource_group_name   = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  sku                   = "Essential"
  grafana_major_version = "10"
}