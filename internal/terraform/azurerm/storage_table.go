package azurerm

import (
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getStorageTableRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_storage_table",
		CoreRFunc: newStorageTable,
		ReferenceAttributes: []string{
			"storage_account_name",
		},
		GetRegion: func(defaultRegion string, d *schema.ResourceData) string {
			return lookupRegion(d, []string{"storage_account_name"})
		},
	}
}

func newStorageTable(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	accountReplicationType := "LRS"
	accountKind := "StorageV2"
	accountEncrypted := false

	if len(d.References("storage_account_name")) > 0 {
		storageAccount := d.References("storage_account_name")[0]

		accountTier := storageAccount.Get("account_tier").String()
		if strings.EqualFold(accountTier, "premium") {
			logging.Logger.Warn().Msgf("Skipping resource %s. Storage Tables don't support %s tier", d.Address, accountTier)
			return nil
		}

		accountReplicationType = storageAccount.Get("account_replication_type").String()
		accountKind = storageAccount.Get("account_kind").String()

		if strings.EqualFold(storageAccount.Get("table_encryption_key_type").String(), "Account") {
			accountEncrypted = true
		}
	}

	switch strings.ToLower(accountReplicationType) {
	case "ragrs":
		accountReplicationType = "RA-GRS"
	case "ragzrs":
		accountReplicationType = "RA-GZRS"
	}

	return &azure.StorageTable{
		Address:                d.Address,
		Region:                 region,
		AccountKind:            accountKind,
		AccountReplicationType: accountReplicationType,
		AccountEncrypted:       accountEncrypted,
	}
}
