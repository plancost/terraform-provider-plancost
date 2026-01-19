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

func getWindowsVirtualMachineScaleSetRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_windows_virtual_machine_scale_set",
		RFunc: NewWindowsVirtualMachineScaleSet,
	}
}
func NewWindowsVirtualMachineScaleSet(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	r := &azure.WindowsVirtualMachineScaleSet{
		Address:                               d.Address,
		Region:                                d.Region,
		SKU:                                   d.Get("sku").String(),
		LicenseType:                           d.Get("license_type").String(),
		IsDevTest:                             d.ProjectMetadata["isProduction"] == "false",
		AdditionalCapabilitiesUltraSSDEnabled: d.Get("additional_capabilities.0.ultra_ssd_enabled").Bool(),
	}
	if len(d.Get("os_disk").Array()) > 0 {
		diskData := d.Get("os_disk").Array()[0]
		r.OSDiskData = &azure.ManagedDiskData{
			DiskType:   diskData.Get("storage_account_type").String(),
			DiskSizeGB: diskData.Get("disk_size_gb").Int(),
		}
	}

	r.PopulateUsage(u)

	if u == nil || u.IsEmpty("instances") {
		r.Instances = intPtr(d.Get("instances").Int())
	}

	return r.BuildResource()
}
