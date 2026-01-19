resource "azurerm_resource_group" "example" {
  name     = "my-example-resources"
  location = "West US"
}

resource "azurerm_public_ip" "example" {
  name                = "my-example-ip-1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_public_ip" "example2" {
  name                = "my-example-ip-2"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"

  tags = {
    environment = "Production"
  }
}


resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  export_markdown_file = abspath("${path.module}/estimate.md")
}
