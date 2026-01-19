package azurerm

import (
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getNetAppPoolRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_netapp_pool",
		CoreRFunc: newNetAppPool,
	}
}

func newNetAppPool(d *schema.ResourceData) schema.CoreResource {
	return &azure.NetAppPool{
		Address:               d.Address,
		Region:                d.Region,
		SizeInTB:              d.Get("size_in_tb").Int(),
		ServiceLevel:          d.Get("service_level").String(),
		EncryptionType:        d.Get("encryption_type").String(),
		CoolAccessEnabled:     d.Get("cool_access_enabled").Bool(),
		CustomThroughputMiBps: d.Get("custom_throughput_mibps").Float(),
	}
}
