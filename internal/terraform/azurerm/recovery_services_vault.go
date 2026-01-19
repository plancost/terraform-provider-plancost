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
	"sort"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getRecoveryServicesVaultRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_recovery_services_vault",
		CoreRFunc: newRecoveryServicesVault,
		ReferenceAttributes: []string{
			"resource_group_name",
			"azurerm_backup_protected_vm.recovery_vault_name",
		},
		CustomRefIDFunc: func(d *schema.ResourceData) []string {
			name := d.Get("name").String()
			if name != "" {
				return []string{name}
			}

			return nil
		},
	}
}

func newRecoveryServicesVault(d *schema.ResourceData) schema.CoreResource {
	region := d.Region
	vms := d.References("azurerm_backup_protected_vm.recovery_vault_name")

	var protectedVMs []*azure.BackupProtectedVM
	for _, vm := range vms {
		protectedVm := newBackupProtectedVm(vm)
		if protectedVm != nil {
			protectedVMs = append(protectedVMs, protectedVm)
		}
	}

	sort.Slice(protectedVMs, func(i, j int) bool {
		return protectedVMs[i].Address < protectedVMs[j].Address
	})

	return &azure.RecoveryServicesVault{
		Address:      d.Address,
		Region:       region,
		ProtectedVMs: protectedVMs,
	}
}
