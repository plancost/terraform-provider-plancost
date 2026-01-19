// Copyright 2021 Infracost Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azurerm

import (
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getStorageShareRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_storage_share",
		CoreRFunc: newStorageShare,
		ReferenceAttributes: []string{
			"storage_account_name",
		},
		GetRegion: func(defaultRegion string, d *schema.ResourceData) string {
			return lookupRegion(d, []string{"storage_account_name"})
		},
	}
}

func newStorageShare(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	accountReplicationType := "LRS"

	accessTier := d.Get("access_tier").String()
	if accessTier == "" {
		accessTier = "TransactionOptimized"
	}
	quota := d.Get("quota").Int()

	if len(d.References("storage_account_name")) > 0 {
		storageAccount := d.References("storage_account_name")[0]
		accountKind := storageAccount.Get("account_kind").String()
		accountReplicationType = storageAccount.Get("account_replication_type").String()

		if strings.EqualFold(accessTier, "premium") && !strings.EqualFold(accountKind, "filestorage") {
			logging.Logger.Warn().Msgf("Skipping resource %s. Premium access tier is only supported for FileStorage accounts", d.Address)
			return nil
		}
	}

	return &azure.StorageShare{
		Address:                d.Address,
		Region:                 region,
		AccountReplicationType: accountReplicationType,
		AccessTier:             accessTier,
		Quota:                  quota,
	}
}
