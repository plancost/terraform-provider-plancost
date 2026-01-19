provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = "example-resources"
  location                 = "East US"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "plancost_estimate" "this" {
  working_directory     = abspath(path.module)
  usage_file = abspath("${path.module}/usage.yml")
}
