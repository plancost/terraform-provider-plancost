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
	"fmt"
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"

	"github.com/shopspring/decimal"
)

type ApplicationInsightsWebTest struct {
	Address string
	Region  string
	Kind    string
	Enabled bool
}

func (r *ApplicationInsightsWebTest) CoreType() string {
	return "ApplicationInsightsWebTest"
}

func (r *ApplicationInsightsWebTest) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{}
}

func (r *ApplicationInsightsWebTest) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *ApplicationInsightsWebTest) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{}

	if r.Kind != "" {
		if strings.ToLower(r.Kind) == "multistep" && r.Enabled {
			costComponents = append(costComponents, &schema.CostComponent{
				Name:            "Multi-step web test",
				Unit:            "test",
				UnitMultiplier:  decimal.NewFromInt(1),
				MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("azure"),
					Region:        strPtr(r.Region),
					Service:       strPtr("Application Insights"),
					ProductFamily: strPtr("Management and Governance"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "meterName", ValueRegex: strPtr(fmt.Sprintf("/^%s$/i", "Multi-step Web Test"))},
						// Code change note: skuName filter is commented out to match pricing data
						//{Key: "skuName", ValueRegex: strPtr(fmt.Sprintf("/^%s$/i", "Enterprise"))},
					},
				},
			})
		}
	}

	if len(costComponents) == 0 {
		return &schema.Resource{
			Name:        r.Address,
			IsSkipped:   true,
			NoPrice:     true,
			UsageSchema: r.UsageSchema(),
		}
	}

	return &schema.Resource{
		Name:           r.Address,
		CostComponents: costComponents,
		UsageSchema:    r.UsageSchema(),
	}

}
