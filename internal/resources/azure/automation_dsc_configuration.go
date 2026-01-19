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

package azure

import (
	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"

	"fmt"

	"github.com/shopspring/decimal"
)

type AutomationDSCConfiguration struct {
	Address                 string
	Region                  string
	NonAzureConfigNodeCount *int64 `infracost_usage:"non_azure_config_node_count"`
}

func (r *AutomationDSCConfiguration) CoreType() string {
	return "AutomationDSCConfiguration"
}

func (r *AutomationDSCConfiguration) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{{Key: "non_azure_config_node_count", ValueType: schema.Int64, DefaultValue: 0}}
}

func (r *AutomationDSCConfiguration) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *AutomationDSCConfiguration) BuildResource() *schema.Resource {
	return &schema.Resource{
		Name:           r.Address,
		CostComponents: automationDSCNodesCostComponent(&r.Region, r.NonAzureConfigNodeCount),
		UsageSchema:    r.UsageSchema(),
	}
}

func automationDSCNodesCostComponent(location *string, nonAzureConfigNodeCount *int64) []*schema.CostComponent {
	var nonAzureConfigNodeCountDec *decimal.Decimal

	if nonAzureConfigNodeCount != nil {
		nonAzureConfigNodeCountDec = decimalPtr(decimal.NewFromInt(*nonAzureConfigNodeCount))
	}

	costComponents := make([]*schema.CostComponent, 0)
	costComponents = append(costComponents, nonautomationDSCNodesCostComponent(*location, "5", "Non-Azure Node", "Non-Azure", nonAzureConfigNodeCountDec))

	return costComponents
}

func nonautomationDSCNodesCostComponent(location, startUsage, meterName, skuName string, monthlyQuantity *decimal.Decimal) *schema.CostComponent {
	return &schema.CostComponent{

		Name:            "Non-azure config nodes",
		Unit:            "nodes",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: monthlyQuantity,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(location),
			Service:       strPtr("Automation"),
			ProductFamily: strPtr("Management and Governance"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "meterName", ValueRegex: strPtr(fmt.Sprintf("/^%s$/i", meterName))},
				{Key: "skuName", ValueRegex: strPtr(fmt.Sprintf("/^%s$/i", skuName))},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption:   strPtr("Consumption"),
			StartUsageAmount: strPtr(startUsage),
		},
		UsageBased: true,
	}
}
