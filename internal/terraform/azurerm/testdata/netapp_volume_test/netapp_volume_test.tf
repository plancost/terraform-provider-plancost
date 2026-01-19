provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US 2"
}

resource "azurerm_netapp_account" "example" {
  name                = "example-netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_pool" "example" {
  name                = "example-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "netapp"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_volume" "example" {
  name                = "example-netappvolume"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "my-unique-file-path"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.example.id
  storage_quota_in_gb = 100
}

resource "azurerm_netapp_volume" "replicated_10min" {
  name                = "replicated-netappvolume-10min"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "my-unique-replicated-file-path-10min"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.example.id
  storage_quota_in_gb = 1000

  data_protection_replication {
    endpoint_type             = "dst"
    replication_frequency     = "10minutes"
    remote_volume_location    = "southcentralus"
    remote_volume_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.NetApp/netAppAccounts/example-account/capacityPools/example-pool/volumes/example-volume"
  }
}
resource "azurerm_netapp_volume" "replicated_hourly" {
  name                = "replicated-netappvolume-hourly"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "my-unique-replicated-file-path-hourly"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.example.id
  storage_quota_in_gb = 1000

  data_protection_replication {
    endpoint_type             = "dst"
    replication_frequency     = "hourly"
    remote_volume_location    = "southcentralus"
    remote_volume_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.NetApp/netAppAccounts/example-account/capacityPools/example-pool/volumes/example-volume-2"
  }
}

resource "azurerm_netapp_volume" "replicated_daily" {
  name                = "replicated-netappvolume-daily"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "my-unique-replicated-file-path-daily"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.example.id
  storage_quota_in_gb = 1000

  data_protection_replication {
    endpoint_type             = "dst"
    replication_frequency     = "daily"
    remote_volume_location    = "southcentralus"
    remote_volume_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.NetApp/netAppAccounts/example-account/capacityPools/example-pool/volumes/example-volume-3"
  }
}

resource "azurerm_netapp_volume" "with_usage" {
  name                = "replicated-netappvolume-with-usage"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "my-unique-replicated-file-path-usage"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.example.id
  storage_quota_in_gb = 1000

  data_protection_replication {
    endpoint_type             = "dst"
    replication_frequency     = "daily"
    remote_volume_location    = "southcentralus"
    remote_volume_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.NetApp/netAppAccounts/example-account/capacityPools/example-pool/volumes/example-volume-usage"
  }
}
