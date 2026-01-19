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

func getIoTHubRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_iothub",
		CoreRFunc: newIoTHub,
	}
}

func getIoTHubDPSRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_iothub_dps",
		CoreRFunc: newIoTHubDPS,
	}
}

func newIoTHub(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	sku := d.Get("sku.0.name").String()
	capacity := d.Get("sku.0.capacity").Int()

	r := &azure.IoTHub{
		Address:  d.Address,
		Region:   region,
		Sku:      sku,
		Capacity: capacity,
	}

	return r
}

func newIoTHubDPS(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	sku := d.Get("sku.0.name").String()

	r := &azure.IoTHubDPS{
		Address: d.Address,
		Region:  region,
		Sku:     sku,
	}

	return r
}
