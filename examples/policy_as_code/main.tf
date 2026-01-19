resource "azurerm_resource_group" "test" {
  name     = "plancost-quickstart"
  location = "eastus"
}

resource "azurerm_public_ip" "test" {
  name                = "plancost-quickstart"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Basic"
}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
}