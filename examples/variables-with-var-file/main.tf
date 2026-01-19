variable "location" {
    type = string
}

variable "public_ip_sku" {
    type = string
}

resource "azurerm_resource_group" "test" {
  name     = "plancost-quickstart"
  location = var.location
}

resource "azurerm_public_ip" "test" {
  name                = "plancost-quickstart"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = var.public_ip_sku
}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  var_file = abspath("${path.module}/tf.tfvars")
}