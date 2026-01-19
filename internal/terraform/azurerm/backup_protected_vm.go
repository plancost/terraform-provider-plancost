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

func getBackupProtectedVmRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name: "azurerm_backup_protected_vm",
		ReferenceAttributes: []string{
			"resource_group_name",
			"backup_policy_id",
			"source_vm_id",
			"recovery_vault_name",
		},
		CoreRFunc: func(d *schema.ResourceData) schema.CoreResource {
			return schema.BlankCoreResource{
				Name: d.Address,
				Type: d.Type,
			}
		},
	}
}

// newBackupProtectedVm returns a azure.BackupProtectedVM with attributes parsed from HCL.
// Note: archive tier not supported https://github.com/hashicorp/terraform-provider-azurerm/issues/21051 by the provider.
func newBackupProtectedVm(d *schema.ResourceData) *azure.BackupProtectedVM {
	region := d.Region
	vms := d.References("source_vm_id")
	if len(vms) == 0 {
		logging.Logger.Warn().Msgf("skipping resource %s as cannot find referenced source vm", d.Address)
		return nil
	}

	vm := vms[0]
	var osDiskSizeGB int64 = 128
	if vm.Get("storage_os_disk.0.disk_size_gb").Exists() {
		osDiskSizeGB = vm.Get("storage_os_disk.0.disk_size_gb").Int()
	}

	var dataDiskSizeGB int64 = 0
	for _, dd := range vm.Get("storage_data_disk").Array() {
		dataDiskSizeGB += dd.Get("disk_size_gb").Int()
	}

	diskSizeGB := osDiskSizeGB + dataDiskSizeGB
	storageType := "GRS"
	recoveryVaults := d.References("recovery_vault_name")
	if len(recoveryVaults) > 0 {
		vault := recoveryVaults[0]
		mode := strings.ToLower(vault.Get("storage_mode_type").String())
		switch mode {
		case "locallyredundant":
			storageType = "LRS"
		case "zoneredundant":
			storageType = "ZRS"
		}

		crossRegion := vault.GetBoolOrDefault("cross_region_restore_enabled", false)
		if storageType == "GRS" && crossRegion {
			storageType = "RA-GRS"
		}
	}

	return &azure.BackupProtectedVM{
		Address:     d.Address,
		Region:      region,
		StorageType: storageType,
		DiskSizeGB:  float64(diskSizeGB),
	}
}
