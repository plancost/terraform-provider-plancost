resource "azurerm_resource_group" "example" {
  name     = "exampleRG1"
  location = "eastus"

  tags = {
    Environment = "Dev"
    Owner       = "platform-team"
  }
}

resource "azurerm_public_ip" "example" {
  name                = "example-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"

  # No tags defined (used to demonstrate policy violations).
}

resource "azurerm_public_ip" "example2" {
  name                = "example-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"

  tags = {
    Environment = "Staging"
  }
}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  # Validate that Owner tag is set.
  tagging_policy {
    key    = "Owner"
    action = "warning"
  }

  # Validate that Environment tag is either "Dev" or "Prod".
  tagging_policy {
    key            = "Environment"
    allowed_values = ["Dev", "Prod"]
    action         = "block"
  }

  # Validate that Owner tag is a valid email address.
  tagging_policy {
    key            = "Owner"
    resource_types = ["azurerm_resource_group"]
    pattern        = "^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$"
    action         = "block"
  }
}