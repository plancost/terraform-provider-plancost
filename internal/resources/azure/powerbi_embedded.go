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

	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
	"github.com/shopspring/decimal"
)

// PowerBIEmbedded struct represents a Power BI Embedded resource.
//
// Resource information: https://learn.microsoft.com/en-us/power-bi/developer/embedded/
// Pricing information: https://azure.microsoft.com/en-gb/pricing/details/power-bi-embedded/
type PowerBIEmbedded struct {
	Address string
	SKU     string
	Region  string
}

func (r *PowerBIEmbedded) CoreType() string {
	return "PowerBIEmbedded"
}

func (r *PowerBIEmbedded) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{}
}

// PopulateUsage parses the u schema.UsageData into the PowerBIEmbedded.
// It uses the `infracost_usage` struct tags to populate data into the PowerBIEmbedded.
func (r *PowerBIEmbedded) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid PowerBIEmbedded struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *PowerBIEmbedded) BuildResource() *schema.Resource {
	return &schema.Resource{
		Name:           r.Address,
		CostComponents: []*schema.CostComponent{r.instanceUsageCostComponent()},
	}
}

func (r *PowerBIEmbedded) instanceUsageCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:           fmt.Sprintf("Node usage (%s)", r.SKU),
		Unit:           "hours",
		UnitMultiplier: decimal.NewFromInt(1),
		HourlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Power BI Embedded"),
			ProductFamily: strPtr("Analytics"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "skuName", Value: strPtr(r.SKU)},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	}
}
