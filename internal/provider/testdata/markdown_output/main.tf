provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_group" "example" {
  name     = "exampleRG-markdown"
  location = "eastus"
}

resource "azurerm_public_ip" "example" {
  name                = "example-public-ip-markdown"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "plancost_estimate" "this" {
  working_directory          = abspath(path.module)
  export_markdown_file = abspath("${path.module}/estimate.md")
}
