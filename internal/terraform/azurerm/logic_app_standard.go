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

func getLogicAppStandardRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_logic_app_standard",
		CoreRFunc: newLogicAppStandard,
		ReferenceAttributes: []string{
			"resource_group_name",
			"app_service_plan_id",
		},
	}
}

func newLogicAppStandard(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	var sku *string
	appServicePlans := d.References("app_service_plan_id")
	if len(appServicePlans) > 0 {
		sku = strPtr(appServicePlans[0].Get("sku.0.size").String())
	}

	return &azure.LogicAppStandard{
		Address: d.Address,
		Region:  region,
		SKU:     sku,
	}
}
