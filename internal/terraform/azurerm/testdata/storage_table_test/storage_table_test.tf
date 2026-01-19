provider "azurerm" {
  skip_provider_registration = true
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "eastus"
}

locals {
  storagev1_options = [
    {
      account_kind              = "Storage"
      account_tier              = "Standard"
      account_replication_types = ["LRS", "GRS", "RAGRS"],
    },
  ]

  storagev2_options = [
    {
      account_kind              = "StorageV2"
      account_tier              = "Standard"
      account_replication_types = ["LRS", "GRS", "RAGRS", "ZRS", "GZRS", "RAGZRS"],
    },
  ]

  storagev1_permutations = distinct(flatten([
    for storagev1_option in local.storagev1_options : [
      for account_replication_type in storagev1_option.account_replication_types : {
        account_kind             = storagev1_option.account_kind
        account_tier             = storagev1_option.account_tier
        account_replication_type = account_replication_type
      }
    ]
  ]))

  storagev2_permutations = distinct(flatten([
    for storagev2_option in local.storagev2_options : [
      for account_replication_type in storagev2_option.account_replication_types : {
        account_kind             = storagev2_option.account_kind
        account_tier             = storagev2_option.account_tier
        account_replication_type = account_replication_type
      }
    ]
  ]))
}

resource "azurerm_storage_account" "storagev1" {
  for_each = { for entry in local.storagev1_permutations : "${entry.account_kind}.${entry.account_tier}.${entry.account_replication_type}" => entry }

  name                     = substr(lower("ic${each.value.account_kind}${each.value.account_tier}${each.value.account_replication_type}"), 0, 24)
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_kind             = each.value.account_kind
  account_tier             = each.value.account_tier
  account_replication_type = each.value.account_replication_type
}

resource "azurerm_storage_table" "storagev1-table" {
  for_each = { for entry in local.storagev1_permutations : "${entry.account_kind}.${entry.account_tier}.${entry.account_replication_type}" => entry }

  name                 = substr(lower("${each.value.account_kind}${each.value.account_tier}${each.value.account_replication_type}"), 0, 24)
  storage_account_name = azurerm_storage_account.storagev1["${each.value.account_kind}.${each.value.account_tier}.${each.value.account_replication_type}"].name
}

resource "azurerm_storage_account" "storagev2" {
  for_each = { for entry in local.storagev2_permutations : "${entry.account_kind}.${entry.account_tier}.${entry.account_replication_type}" => entry }

  name                     = substr(lower("ic${each.value.account_kind}${each.value.account_tier}${each.value.account_replication_type}"), 0, 24)
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_kind             = each.value.account_kind
  account_tier             = each.value.account_tier
  account_replication_type = each.value.account_replication_type
}

resource "azurerm_storage_account" "storagev2_encrypted" {
  for_each = { for entry in local.storagev2_permutations : "${entry.account_kind}.${entry.account_tier}.${entry.account_replication_type}" => entry }

  name                      = substr(lower("enc${each.value.account_kind}${each.value.account_tier}${each.value.account_replication_type}"), 0, 24)
  resource_group_name       = azurerm_resource_group.example.name
  location                  = azurerm_resource_group.example.location
  account_kind              = each.value.account_kind
  account_tier              = each.value.account_tier
  account_replication_type  = each.value.account_replication_type
  table_encryption_key_type = "Account"
}

resource "azurerm_storage_table" "storagev2-table-encrypted" {
  for_each = { for entry in local.storagev2_permutations : "${entry.account_kind}.${entry.account_tier}.${entry.account_replication_type}" => entry }

  name                 = substr(lower("enc${each.value.account_kind}${each.value.account_tier}${each.value.account_replication_type}"), 0, 24)
  storage_account_name = azurerm_storage_account.storagev2_encrypted["${each.value.account_kind}.${each.value.account_tier}.${each.value.account_replication_type}"].name
}

resource "azurerm_storage_table" "storagev2-table" {
  for_each = { for entry in local.storagev2_permutations : "${entry.account_kind}.${entry.account_tier}.${entry.account_replication_type}" => entry }

  name                 = substr(lower("ic${each.value.account_kind}${each.value.account_tier}${each.value.account_replication_type}"), 0, 24)
  storage_account_name = azurerm_storage_account.storagev2["${each.value.account_kind}.${each.value.account_tier}.${each.value.account_replication_type}"].name
}
