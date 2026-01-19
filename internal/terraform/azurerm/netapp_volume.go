package azurerm

import (
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getNetAppVolumeRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_netapp_volume",
		CoreRFunc: newNetAppVolume,
	}
}

func newNetAppVolume(d *schema.ResourceData) schema.CoreResource {
	return &azure.NetAppVolume{
		Address:              d.Address,
		Region:               d.Region,
		StorageQuotaInGB:     d.Get("storage_quota_in_gb").Int(),
		ReplicationFrequency: d.Get("data_protection_replication.0.replication_frequency").String(),
		RemoteVolumeLocation: d.Get("data_protection_replication.0.remote_volume_location").String(),
	}
}
