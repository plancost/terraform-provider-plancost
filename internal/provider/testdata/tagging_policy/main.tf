provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_group" "example" {
  name     = "exampleRG1"
  location = "eastus"
  tags = {
    Environment = "Dev"
  }
}

resource "azurerm_public_ip" "example" {
  name                = "example-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
  # Missing "Environment" tag
}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  tagging_policy {
    key = "Environment"
    allowed_values = ["Dev", "Prod"]
    action = "block"
  }
}