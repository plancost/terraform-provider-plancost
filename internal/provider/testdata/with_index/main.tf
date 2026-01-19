provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_group" "example" {
  name     = "exampleRG1"
  location = "eastus"
}
resource "azurerm_public_ip" "example" {
  for_each = {
    "one" = "Standard"
    "two" = "Basic"
  }
  name                = "example-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
  sku                 = each.value
}

resource "azurerm_public_ip" "example2" {
  count = 2
  name                = "example-public-ip${count.index}"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}


resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
}