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

	"github.com/tidwall/gjson"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getVirtualMachineScaleSetRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_virtual_machine_scale_set",
		CoreRFunc: NewVirtualMachineScaleSet,
	}
}
func NewVirtualMachineScaleSet(d *schema.ResourceData) schema.CoreResource {
	r := &azure.VirtualMachineScaleSet{
		Address:     d.Address,
		Region:      d.Region,
		SKUName:     d.Get("sku.0.name").String(),
		SKUCapacity: d.Get("sku.0.capacity").Int(),
		LicenseType: d.Get("license_type").String(),
		IsDevTest:   d.ProjectMetadata["isProduction"] == "false",
	}

	if !d.IsEmpty("os_profile_windows_config") {
		r.IsWindows = true
	}
	if d.Get("storage_profile_os_disk.0.os_type").Type != gjson.Null {
		if strings.ToLower(d.Get("storage_profile_os_disk.0.os_type").String()) == "windows" {
			r.IsWindows = true
		}
	}
	if d.Get("storage_profile_image_reference.0.offer").Type != gjson.Null {
		if strings.ToLower(d.Get("storage_profile_image_reference.0.offer").String()) == "windowsserver" {
			r.IsWindows = true
		}
	}

	if len(d.Get("storage_profile_os_disk").Array()) > 0 {
		storageData := d.Get("storage_profile_os_disk").Array()[0]
		r.StorageProfileOSDiskData = &azure.ManagedDiskData{
			DiskType:   storageData.Get("managed_disk_type").String(),
			DiskSizeGB: storageData.Get("disk_size_gb").Int(),
		}
	}

	if len(d.Get("storage_profile_data_disk").Array()) > 0 {
		for _, s := range d.Get("storage_profile_data_disk").Array() {
			if s.Get("managed_disk_type").Type == gjson.Null {
				continue
			}
			r.StorageProfileOSDisksData = append(r.StorageProfileOSDisksData, &azure.ManagedDiskData{
				DiskType:   s.Get("managed_disk_type").String(),
				DiskSizeGB: s.Get("disk_size_gb").Int(),
			})
		}
	}

	return r
}
