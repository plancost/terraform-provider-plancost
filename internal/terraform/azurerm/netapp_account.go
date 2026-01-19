package azurerm

import (
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getNetAppAccountRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_netapp_account",
		CoreRFunc: newNetAppAccount,
	}
}

func newNetAppAccount(d *schema.ResourceData) schema.CoreResource {
	return &azure.NetAppAccount{
		Address: d.Address,
		Region:  d.Region,
	}
}
