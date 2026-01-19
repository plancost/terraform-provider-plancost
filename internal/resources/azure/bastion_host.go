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
	"math"

	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"

	"fmt"

	"github.com/shopspring/decimal"

	"github.com/plancost/terraform-provider-plancost/internal/usage"
)

type BastionHost struct {
	Address               string
	Region                string
	MonthlyOutboundDataGB *float64 `infracost_usage:"monthly_outbound_data_gb"`
}

func (r *BastionHost) CoreType() string {
	return "BastionHost"
}

func (r *BastionHost) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{{Key: "monthly_outbound_data_gb", ValueType: schema.Float64, DefaultValue: 0}}
}

func (r *BastionHost) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *BastionHost) BuildResource() *schema.Resource {
	productType := "Basic"

	costComponents := []*schema.CostComponent{
		{
			Name:           "Bastion host",
			Unit:           "hours",
			UnitMultiplier: decimal.NewFromInt(1),
			HourlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(r.Region),
				Service:       strPtr("Azure Bastion"),
				ProductFamily: strPtr("Networking"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "skuName", Value: strPtr(productType)},
					{Key: "meterName", Value: strPtr(fmt.Sprintf("%s Gateway", productType))},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Consumption"),
			},
		},
	}

	costComponents = append(costComponents, r.outboundDataTransferComponents(productType)...)

	return &schema.Resource{
		Name:           r.Address,
		CostComponents: costComponents,
		UsageSchema:    r.UsageSchema(),
	}
}

func (r *BastionHost) outboundDataTransferComponents(productType string) []*schema.CostComponent {
	costComponents := make([]*schema.CostComponent, 0)
	name := "Outbound data transfer "

	type dataTier struct {
		name       string
		startUsage string
	}

	data := []dataTier{
		{name: fmt.Sprintf("%s%s", name, "(first 10TB)"), startUsage: "5"},
		{name: fmt.Sprintf("%s%s", name, "(next 40TB)"), startUsage: "10240"},
		{name: fmt.Sprintf("%s%s", name, "(next 100TB)"), startUsage: "51200"},
		{name: fmt.Sprintf("%s%s", name, "(next 350TB)"), startUsage: "153600"},
		{name: fmt.Sprintf("%s%s", name, "(over 500TB)"), startUsage: "512000"},
	}

	if r.MonthlyOutboundDataGB != nil {
		// Note-plancost: fix the tier limit
		// the first tier limit is changed from 10000 to 9995 because first tier covers up to 10TB (10240GB) but first 5GB is free
		// the second tier limit is changed from 50000 to 40000 because it's next 40TB after first 10TB
		// the third tier limit is changed from 150000 to 100000 because it's next 100TB
		// the fourth tier limit is changed from 500000 to 350000 because it's next 350TB
		tierLimits := []int64{9995, 40000, 100000, 350000}
		// Subtract the free 5GB from the total monthly outbound data
		monthlyOutboundDataGb := decimal.NewFromFloat(math.Max(0, *r.MonthlyOutboundDataGB-5))
		tiers := usage.CalculateTierBuckets(monthlyOutboundDataGb, tierLimits)
		for i, d := range data {
			if i < len(tiers) && tiers[i].GreaterThan(decimal.Zero) {
				costComponents = append(costComponents, r.outboundDataTransferSingleComponent(
					d.name,
					productType,
					d.startUsage,
					decimalPtr(tiers[i])))
			}
		}
	} else {
		costComponents = append(costComponents, r.outboundDataTransferSingleComponent(
			data[0].name,
			productType,
			data[0].startUsage,
			nil))
	}
	return costComponents
}

func (r *BastionHost) outboundDataTransferSingleComponent(name, productType, startUsage string, quantity *decimal.Decimal) *schema.CostComponent {
	return &schema.CostComponent{
		Name:            name,
		Unit:            "GB",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Azure Bastion"),
			ProductFamily: strPtr("Networking"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "skuName", Value: strPtr(productType)},
				{Key: "meterName", Value: strPtr(fmt.Sprintf("%s Data Transfer Out", productType))},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption:   strPtr("Consumption"),
			StartUsageAmount: strPtr(startUsage),
		},
		UsageBased: true,
	}
}
