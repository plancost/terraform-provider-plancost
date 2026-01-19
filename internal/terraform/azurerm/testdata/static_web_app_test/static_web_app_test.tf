provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US 2"
}

resource "azurerm_static_web_app" "standard" {
  name                = "example-static-web-app-standard"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_tier            = "Standard"
  sku_size            = "Standard"
}

resource "azurerm_static_web_app" "free" {
  name                = "example-static-web-app-free"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_tier            = "Free"
  sku_size            = "Free"
}

resource "azurerm_static_web_app" "usage_standard" {
  name                = "example-static-web-app-usage"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_tier            = "Standard"
  sku_size            = "Standard"
}
