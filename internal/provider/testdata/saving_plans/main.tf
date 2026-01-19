
provider "azurerm" {
  features {}
}

resource "azurerm_linux_virtual_machine" "vm" {
  name                = "basic_b1"
  resource_group_name = "fake_resource_group"
  location            = "eastus"

  size           = "Standard_DS1_v2"
  admin_username = "fakeuser"
  admin_password = "Password1234!"

  network_interface_ids = [
    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg/providers/Microsoft.Network/networkInterfaces/fakenic",
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  recommendations_enabled = true
}
