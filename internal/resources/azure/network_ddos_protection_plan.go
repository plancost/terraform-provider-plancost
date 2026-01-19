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
	"github.com/shopspring/decimal"

	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

// NetworkDdosProtectionPlan struct represents Azure DDoS Protection Plan.
// DDoS Protection Plan is a resource that provides DDoS protection for virtual networks and IPs.
//
// Resource information: https://azure.microsoft.com/en-us/products/ddos-protection/
// Pricing information: https://azure.microsoft.com/en-us/pricing/details/ddos-protection/#pricing
type NetworkDdosProtectionPlan struct {
	Address       string
	Region        string
	OverageAmount *int64 `infracost_usage:"overage_amount"`
}

// CoreType returns the name of this resource type.
func (r *NetworkDdosProtectionPlan) CoreType() string {
	return "NetworkDdosProtectionPlan"
}

// UsageSchema defines a list which represents the usage schema of
// NetworkDdosProtectionPlan. There is only one usage item, `overage_amount`,
// which represents the number of resources that fall outside the base ddos
// coverage.
func (r *NetworkDdosProtectionPlan) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "overage_amount", DefaultValue: 0, ValueType: schema.Int64},
	}
}

// PopulateUsage parses the u schema.UsageData into the
// NetworkDdosProtectionPlan. It uses the `infracost_usage` struct tags to
// populate data into the NetworkDdosProtectionPlan.
func (r *NetworkDdosProtectionPlan) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid NetworkDdosProtectionPlan
// struct. This method is called after the resource is initialised by an IaC
// provider.
//
// BuildResource returns two cost components:
//
//  1. DDoS Protection Plan: The cost of the DDoS Protection Plan.
//  2. Overage charges: The cost of the overage charges for the DDoS Protection Plan.
//     This is the number of resources that fall outside the base coverage offered by
//     the protection plan (100). This amount is defined in the usage file as it is
//     difficult to infer the number of resources that fall outside the base coverage
//     from the IaC.
func (r *NetworkDdosProtectionPlan) BuildResource() *schema.Resource {
	var overageAmount *decimal.Decimal
	overageUnit := "resource"
	if r.OverageAmount != nil {
		overageAmount = decimalPtr(decimal.NewFromInt(*r.OverageAmount))
		if overageAmount.GreaterThan(decimal.NewFromInt(1)) {
			overageUnit = "resources"
		}
	}

	return &schema.Resource{
		Name:        r.Address,
		UsageSchema: r.UsageSchema(),
		CostComponents: []*schema.CostComponent{
			{
				Name:           "DDoS Protection Plan",
				Unit:           "hours",
				UnitMultiplier: decimal.NewFromInt(1),
				HourlyQuantity: decimalPtr(decimal.NewFromInt(1)),
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr(vendorName),
					Region:        strPtr(r.Region),
					Service:       strPtr("Azure DDOS Protection"),
					ProductFamily: strPtr("Networking"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "skuName", Value: strPtr("Network Protection")},
						{Key: "meterName", Value: strPtr("Network Protection Plan")},
					},
				},
			},
			{
				Name:           "Overage charges",
				Unit:           overageUnit,
				UnitMultiplier: schema.HourToMonthUnitMultiplier,
				HourlyQuantity: overageAmount,
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr(vendorName),
					Region:        strPtr(r.Region),
					Service:       strPtr("Azure DDOS Protection"),
					ProductFamily: strPtr("Networking"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "skuName", Value: strPtr("Network Protection")},
						{Key: "meterName", Value: strPtr("Network Protection Resource")},
					},
				},
				UsageBased: true,
			},
		},
	}
}
