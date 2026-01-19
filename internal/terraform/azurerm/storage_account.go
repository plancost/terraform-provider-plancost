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

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getStorageAccountRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_storage_account",
		CoreRFunc: newAzureRMStorageAccount,
		ReferenceAttributes: []string{
			"azurerm_storage_management_policy.storage_account_id",
		},
		CustomRefIDFunc: func(d *schema.ResourceData) []string {
			return []string{d.Get("name").String()}
		},
	}
}

func newAzureRMStorageAccount(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	accountKind := "StorageV2"
	if !d.IsEmpty("account_kind") {
		accountKind = d.Get("account_kind").String()
	}

	accountReplicationType := d.Get("account_replication_type").String()
	switch strings.ToLower(accountReplicationType) {
	case "ragrs":
		accountReplicationType = "RA-GRS"
	case "ragzrs":
		accountReplicationType = "RA-GZRS"
	}

	accountTier := d.Get("account_tier").String()

	accessTier := "Hot"
	if !d.IsEmpty("access_tier") {
		accessTier = d.Get("access_tier").String()
	}

	nfsv3 := false
	if !d.IsEmpty("nfsv3_enabled") {
		nfsv3 = d.Get("nfsv3_enabled").Bool()
	}

	return &azure.StorageAccount{
		Address:                d.Address,
		Region:                 region,
		AccessTier:             accessTier,
		AccountKind:            accountKind,
		AccountReplicationType: accountReplicationType,
		AccountTier:            accountTier,
		NFSv3:                  nfsv3,
	}
}
