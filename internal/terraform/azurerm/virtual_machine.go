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
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getVirtualMachineRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_virtual_machine",
		CoreRFunc: NewVirtualMachine,
	}
}
func NewVirtualMachine(d *schema.ResourceData) schema.CoreResource {
	r := &azure.VirtualMachine{
		Address:                    d.Address,
		Region:                     d.Region,
		StorageImageReferenceOffer: d.Get("storage_image_reference.0.offer").String(),
		StorageOSDiskOSType:        d.Get("storage_os_disk.0.os_type").String(),
		LicenseType:                d.Get("license_type").String(),
		VMSize:                     d.Get("vm_size").String(),
		StoragesDiskData:           make([]*azure.ManagedDiskData, 0),
	}

	if len(d.Get("storage_os_disk").Array()) > 0 {
		storageData := d.Get("storage_os_disk").Array()[0]
		r.StorageOSDiskData = &azure.ManagedDiskData{
			DiskType:   storageData.Get("managed_disk_type").String(),
			DiskSizeGB: storageData.Get("disk_size_gb").Int(),
		}
	}

	if len(d.Get("storage_data_disk").Array()) > 0 {
		for _, s := range d.Get("storage_data_disk").Array() {
			r.StoragesDiskData = append(r.StoragesDiskData, &azure.ManagedDiskData{
				DiskType:   s.Get("managed_disk_type").String(),
				DiskSizeGB: s.Get("disk_size_gb").Int(),
			})
		}
	}

	return r
}
