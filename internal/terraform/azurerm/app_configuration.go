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

func getAppConfigurationRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_app_configuration",
		CoreRFunc: newAppConfiguration,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}

func newAppConfiguration(d *schema.ResourceData) schema.CoreResource {
	region := d.Region
	sku := strings.ToLower(strings.TrimSpace(d.Get("sku").String()))
	if sku == "" {
		sku = "free"
	}
	array := d.Get("replica").Array()
	return &azure.AppConfiguration{
		Address:  d.Address,
		Region:   region,
		SKU:      sku,
		Replicas: int64(len(array)),
	}
}
