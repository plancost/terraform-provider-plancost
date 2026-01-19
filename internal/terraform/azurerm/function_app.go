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

func getFunctionAppRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name: "azurerm_function_app",
		ReferenceAttributes: []string{
			"app_service_plan_id",
		},
		CoreRFunc: func(d *schema.ResourceData) schema.CoreResource {
			return newFunctionApp(d)
		},
	}
}

func newFunctionApp(d *schema.ResourceData) *azure.FunctionApp {
	appServicePlan := d.References("app_service_plan_id")
	servicePlan := d.References("service_plan_id")
	region := d.Region

	if len(appServicePlan) == 0 && len(servicePlan) == 0 {
		return &azure.FunctionApp{
			Address: d.Address,
			Region:  region,
			Tier:    "standard",
		}
	}

	if len(appServicePlan) > 0 {
		data := appServicePlan[0]
		return newFunctionAppFromAppServicePlanRef(d, data)
	}

	data := servicePlan[0]
	return newFunctionAppFromAppServicePlanRef(d, data)
}

func newFunctionAppFromAppServicePlanRef(d *schema.ResourceData, data *schema.ResourceData) *azure.FunctionApp {
	region := d.Region

	tier := "standard"
	// support for the legacy azurerm_app_service_plan resource. This is only applicable for the legacy azurerm_function_app resource.
	if data.Get("sku").Exists() {
		skuTier := strings.ToLower(data.Get("sku.0.tier").String())
		skuSize := strings.ToLower(data.Get("sku.0.size").String())
		kind := strings.ToLower(data.Get("kind").String())

		if strings.ToLower(skuSize) != "y1" && (strings.ToLower(kind) == "elastic" || strings.ToLower(skuTier) == "elasticpremium") {
			tier = "premium"
		}

		return &azure.FunctionApp{
			Address: d.Address,
			Region:  region,
			SKUName: skuSize,
			Tier:    tier,
			OSType:  kind,
		}
	}

	skuName := data.Get("sku_name").String()
	if strings.HasPrefix(strings.ToLower(skuName), "ep") {
		tier = "premium"
	}

	return &azure.FunctionApp{
		Address: d.Address,
		Region:  region,
		SKUName: strings.ToLower(skuName),
		Tier:    tier,
		OSType:  strings.ToLower(data.Get("os_type").String()),
	}
}
