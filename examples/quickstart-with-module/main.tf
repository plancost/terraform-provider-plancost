module "avm-res-storage-storageaccount" {
  source  = "Azure/avm-res-storage-storageaccount/azurerm"
  version = "0.6.7"

  name                = "uniquestorageaccount"
  location            = "eastus"
  resource_group_name = "test"
}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
}

output "monthly_cost" {
  value = plancost_estimate.this.monthly_cost
}

output "estimate_view" {
  value = plancost_estimate.this.view
}
