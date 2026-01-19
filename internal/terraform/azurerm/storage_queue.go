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

func getStorageQueueRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_storage_queue",
		CoreRFunc: newStorageQueue,
		ReferenceAttributes: []string{
			"storage_account_name",
		},
		GetRegion: func(defaultRegion string, d *schema.ResourceData) string {
			return lookupRegion(d, []string{"storage_account_name"})
		},
	}
}

func newStorageQueue(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	accountReplicationType := "LRS"
	accountKind := "StorageV2"

	if len(d.References("storage_account_name")) > 0 {
		storageAccount := d.References("storage_account_name")[0]

		accountTier := storageAccount.Get("account_tier").String()
		if strings.EqualFold(accountTier, "premium") {
			logging.Logger.Warn().Msgf("Skipping resource %s. Storage Queues don't support %s tier", d.Address, accountTier)
			return nil
		}

		accountReplicationType = storageAccount.Get("account_replication_type").String()
		accountKind = storageAccount.Get("account_kind").String()
	}

	switch strings.ToLower(accountReplicationType) {
	case "ragrs":
		accountReplicationType = "RA-GRS"
	case "ragzrs":
		accountReplicationType = "RA-GZRS"
	}

	return &azure.StorageQueue{
		Address:                d.Address,
		Region:                 region,
		AccountKind:            accountKind,
		AccountReplicationType: accountReplicationType,
	}
}
