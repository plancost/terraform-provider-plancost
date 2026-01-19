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

type ApplicationInsights struct {
	Address               string
	Region                string
	RetentionInDays       int64
	MonthlyDataIngestedGB *float64 `infracost_usage:"monthly_data_ingested_gb"`
}

func (r *ApplicationInsights) CoreType() string {
	return "ApplicationInsights"
}

func (r *ApplicationInsights) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{{Key: "monthly_data_ingested_gb", ValueType: schema.Float64, DefaultValue: 0}}
}

func (r *ApplicationInsights) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *ApplicationInsights) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{}

	var dataIngested *decimal.Decimal
	if r.MonthlyDataIngestedGB != nil {
		dataIngested = decimalPtr(decimal.NewFromFloat(*r.MonthlyDataIngestedGB))
	}
	costComponents = append(costComponents, &schema.CostComponent{
		Name:            "Data ingested",
		Unit:            "GB",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: dataIngested,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Application Insights"),
			ProductFamily: strPtr("Management and Governance"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "meterName", ValueRegex: strPtr(fmt.Sprintf("/^%s$/i", "Enterprise Overage Data"))},
				{Key: "skuName", ValueRegex: strPtr(fmt.Sprintf("/^%s$/i", "Enterprise"))},
			},
		},
		UsageBased: true,
	})

	var dataRetentionDays *decimal.Decimal
	if r.RetentionInDays != 0 {
		dataRetentionDays = decimalPtr(decimal.NewFromInt(r.RetentionInDays))

		if dataRetentionDays.GreaterThan(decimal.NewFromInt(90)) && dataIngested != nil {
			days := dataRetentionDays.Sub(decimal.NewFromInt(90)).Div(decimal.NewFromInt(30))
			qty := decimalPtr(dataIngested.Mul(days))
			costComponents = append(costComponents, &schema.CostComponent{
				Name:            fmt.Sprintf("Data retention (%s days)", dataRetentionDays.String()),
				Unit:            "GB",
				UnitMultiplier:  decimal.NewFromInt(1),
				MonthlyQuantity: qty,
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("azure"),
					Region:        strPtr(r.Region),
					Service:       strPtr("Application Insights"),
					ProductFamily: strPtr("Management and Governance"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "meterName", ValueRegex: strPtr(fmt.Sprintf("/^%s$/i", "Data Retention"))},
						{Key: "skuName", ValueRegex: strPtr(fmt.Sprintf("/^%s$/i", "Enterprise"))},
					},
				},
			})
		}
	}

	return &schema.Resource{
		Name:           r.Address,
		CostComponents: costComponents,
		UsageSchema:    r.UsageSchema(),
	}
}
