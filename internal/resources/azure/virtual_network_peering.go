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

// VirtualNetworkPeering struct represents a VNET peering.
//

// Resource information: https://azure.microsoft.com/en-us/services/virtual-network/#overview
// Pricing information: https://azure.microsoft.com/en-us/pricing/details/virtual-network/
type VirtualNetworkPeering struct {
	Address           string
	SourceRegion      string
	DestinationRegion string
	SourceZone        string
	DestinationZone   string

	MonthlyDataTransferGB *float64 `infracost_usage:"monthly_data_transfer_gb"`
}

func (r *VirtualNetworkPeering) CoreType() string {
	return "VirtualNetworkPeering"
}

func (r *VirtualNetworkPeering) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "monthly_data_transfer_gb", DefaultValue: 0, ValueType: schema.Float64},
	}
}

func (r *VirtualNetworkPeering) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *VirtualNetworkPeering) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.ingressDataProcessedCostComponent(),
		r.egressDataProcessedCostComponent(),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}
}

func (r *VirtualNetworkPeering) egressDataProcessedCostComponent() *schema.CostComponent {
	if r.DestinationRegion == r.SourceRegion {
		return &schema.CostComponent{
			Name:            "Outbound data transfer",
			Unit:            "GB",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: floatPtrToDecimalPtr(r.MonthlyDataTransferGB),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr("Global"),
				Service:       strPtr("Virtual Network"),
				ProductFamily: strPtr("Networking"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "meterName", Value: strPtr("Intra-Region Egress")},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Consumption"),
			},
			UsageBased: true,
		}
	}

	return &schema.CostComponent{
		Name:            "Outbound data transfer",
		Unit:            "GB",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: floatPtrToDecimalPtr(r.MonthlyDataTransferGB),
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("azure"),
			Region:     strPtr(r.SourceZone),
			Service:    strPtr("VPN Gateway"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "serviceFamily", ValueRegex: regexPtr("Networking")},
				{Key: "productName", ValueRegex: regexPtr("VPN Gateway Bandwidth")},
				{Key: "meterName", ValueRegex: regexPtr("Inter-Virtual Network Data Transfer Out")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
		UsageBased: true,
	}
}

func (r *VirtualNetworkPeering) ingressDataProcessedCostComponent() *schema.CostComponent {
	if r.DestinationRegion == r.SourceRegion {
		return &schema.CostComponent{
			Name:            "Inbound data transfer",
			Unit:            "GB",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: floatPtrToDecimalPtr(r.MonthlyDataTransferGB),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr("Global"),
				Service:       strPtr("Virtual Network"),
				ProductFamily: strPtr("Networking"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "meterName", Value: strPtr("Intra-Region Ingress")},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Consumption"),
			},
			UsageBased: true,
		}
	}

	return &schema.CostComponent{
		Name:            "Inbound data transfer",
		Unit:            "GB",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: floatPtrToDecimalPtr(r.MonthlyDataTransferGB),
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("azure"),
			Region:     strPtr(r.DestinationZone),
			Service:    strPtr("VPN Gateway"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "serviceFamily", ValueRegex: regexPtr("Networking")},
				{Key: "productName", ValueRegex: regexPtr("VPN Gateway Bandwidth")},
				{Key: "meterName", ValueRegex: regexPtr("Inter-Virtual Network Data Transfer Out")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
		UsageBased: true,
	}
}
