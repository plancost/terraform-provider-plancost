provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_netapp_account" "example" {
  name                = "example-netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_pool" "example" {
  name                = "example-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
}

resource "azurerm_netapp_pool" "premium" {
  name                = "premium-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Premium"
  size_in_tb          = 4
}

resource "azurerm_netapp_pool" "ultra" {
  name                = "ultra-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Ultra"
  size_in_tb          = 4
}

resource "azurerm_netapp_pool" "double_encryption_standard" {
  name                = "double-encryption-standard-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
  encryption_type     = "Double"
}

resource "azurerm_netapp_pool" "double_encryption_premium" {
  name                = "double-encryption-premium-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Premium"
  size_in_tb          = 4
  encryption_type     = "Double"
}

resource "azurerm_netapp_pool" "double_encryption_ultra" {
  name                = "double-encryption-ultra-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Ultra"
  size_in_tb          = 4
  encryption_type     = "Double"
}

resource "azurerm_netapp_pool" "cool_access" {
  name                = "cool-access-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
  cool_access_enabled = true
}

resource "azurerm_netapp_pool" "flexible" {
  name                = "flexible-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Flexible"
  size_in_tb          = 4
  qos_type            = "Manual"
}

resource "azurerm_netapp_pool" "flexible_overage" {
  name                    = "flexible-overage-netapppool"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  account_name            = azurerm_netapp_account.example.name
  service_level           = "Flexible"
  size_in_tb              = 4
  qos_type                = "Manual"
  custom_throughput_mibps = 200
}

resource "azurerm_netapp_pool" "flexible_under_limit" {
  name                    = "flexible-under-limit-netapppool"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  account_name            = azurerm_netapp_account.example.name
  service_level           = "Flexible"
  size_in_tb              = 4
  qos_type                = "Manual"
  custom_throughput_mibps = 128
}
