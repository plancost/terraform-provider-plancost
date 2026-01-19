package azure

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
	"github.com/plancost/terraform-provider-plancost/internal/usage"
)

// CDNFrontDoorProfile struct represents Azure's Front Door (Standard/Premium) profile.
type CDNFrontDoorProfile struct {
	Address string
	Region  string
	SkuName string

	MonthlyOutboundDataTransferGB         *float64 `infracost_usage:"monthly_outbound_data_transfer_gb"`
	MonthlyOutboundDataTransferToOriginGB *float64 `infracost_usage:"monthly_outbound_data_transfer_to_origin_gb"`
	MonthlyRequestsMillions               *float64 `infracost_usage:"monthly_requests_millions"`
}

// CoreType returns the name of this resource type
func (r *CDNFrontDoorProfile) CoreType() string {
	return "CDNFrontDoorProfile"
}

// UsageSchema defines a list which represents the usage schema of CDNFrontDoorProfile.
func (r *CDNFrontDoorProfile) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "monthly_outbound_data_transfer_gb", DefaultValue: 0, ValueType: schema.Float64},
		{Key: "monthly_outbound_data_transfer_to_origin_gb", DefaultValue: 0, ValueType: schema.Float64},
		{Key: "monthly_requests_millions", DefaultValue: 0, ValueType: schema.Float64},
	}
}

// PopulateUsage parses the usage schema for the resource
func (r *CDNFrontDoorProfile) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid CDNFrontDoorProfile struct.
func (r *CDNFrontDoorProfile) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{}

	sku := r.SkuName
	if sku == "" {
		sku = "Standard_AzureFrontDoor"
	}
	tier := "Standard"
	if strings.HasPrefix(strings.ToLower(sku), "premium") {
		tier = "Premium"
	}

	// Base Fee
	costComponents = append(costComponents, r.baseFeeCostComponent(tier))

	// Outbound Data Transfer (Edge to Client)
	if r.MonthlyOutboundDataTransferGB != nil {
		costComponents = append(costComponents, r.outboundDataTransferCostComponents(tier)...)
	}

	// Outbound Data Transfer (Edge to Origin)
	if r.MonthlyOutboundDataTransferToOriginGB != nil {
		costComponents = append(costComponents, r.outboundDataTransferToOriginCostComponent(tier))
	}

	// Requests
	if r.MonthlyRequestsMillions != nil {
		costComponents = append(costComponents, r.requestsCostComponents(tier)...)
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}
}

func (r *CDNFrontDoorProfile) baseFeeCostComponent(tier string) *schema.CostComponent {
	zone := mapRegionToFrontDoorZone(r.Region)
	return &schema.CostComponent{
		Name:            fmt.Sprintf("Base fee (%s)", tier),
		Unit:            "months",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			ProductFamily: strPtr("Networking"),
			Service:       strPtr("Azure Front Door Service"),
			Region:        strPtr(zone),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr("Azure Front Door")},
				{Key: "skuName", Value: strPtr(tier)},
				{Key: "meterName", ValueRegex: strPtr(fmt.Sprintf("/%s Base Fees?/i", tier))},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	}
}

func (r *CDNFrontDoorProfile) outboundDataTransferCostComponents(tier string) []*schema.CostComponent {
	zone := mapRegionToFrontDoorZone(r.Region)
	costComponents := []*schema.CostComponent{}

	// Tiers in GB
	tierLimits := []int64{10000, 40000, 100000, 350000, 500000, 4000000}
	tiers := []struct {
		name       string
		startUsage string
	}{
		{name: "Outbound data transfer (first 10TB)", startUsage: "0"},
		{name: "Outbound data transfer (next 40TB)", startUsage: "10000"},
		{name: "Outbound data transfer (next 100TB)", startUsage: "50000"},
		{name: "Outbound data transfer (next 350TB)", startUsage: "150000"},
		{name: "Outbound data transfer (next 500TB)", startUsage: "500000"},
		{name: "Outbound data transfer (next 4PB)", startUsage: "1000000"},
		{name: "Outbound data transfer (over 5PB)", startUsage: "5000000"},
	}

	buckets := usage.CalculateTierBuckets(decimal.NewFromFloat(*r.MonthlyOutboundDataTransferGB), tierLimits)

	for i, bucket := range buckets {
		if i >= len(tiers) {
			break
		}

		if bucket.IsZero() {
			continue
		}

		costComponents = append(costComponents, &schema.CostComponent{
			Name:            tiers[i].name,
			Unit:            "GB",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(bucket),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				ProductFamily: strPtr("Networking"),
				Service:       strPtr("Azure Front Door Service"),
				Region:        strPtr(zone),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "productName", Value: strPtr("Azure Front Door")},
					{Key: "skuName", Value: strPtr(tier)},
					{Key: "meterName", Value: strPtr(fmt.Sprintf("%s Data Transfer Out", tier))},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption:   strPtr("Consumption"),
				StartUsageAmount: strPtr(tiers[i].startUsage),
			},
			UsageBased: true,
		})
	}

	return costComponents
}

func (r *CDNFrontDoorProfile) outboundDataTransferToOriginCostComponent(tier string) *schema.CostComponent {
	zone := mapRegionToFrontDoorZone(r.Region)
	return &schema.CostComponent{
		Name:            "Outbound data transfer to origin",
		Unit:            "GB",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromFloat(*r.MonthlyOutboundDataTransferToOriginGB)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			ProductFamily: strPtr("Networking"),
			Service:       strPtr("Azure Front Door Service"),
			Region:        strPtr(zone),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr("Azure Front Door")},
				{Key: "skuName", Value: strPtr(tier)},
				{Key: "meterName", Value: strPtr(fmt.Sprintf("%s Data Transfer In", tier))},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
		UsageBased: true,
	}
}

func (r *CDNFrontDoorProfile) requestsCostComponents(tier string) []*schema.CostComponent {
	zone := mapRegionToFrontDoorZone(r.Region)
	costComponents := []*schema.CostComponent{}

	// Tiers in Millions
	tierLimits := []int64{250, 750, 4000}
	tiers := []struct {
		name       string
		startUsage string
	}{
		{name: "Requests (first 250M)", startUsage: "0"},
		{name: "Requests (next 750M)", startUsage: "25000"},
		{name: "Requests (next 4B)", startUsage: "100000"},
		{name: "Requests (over 5B)", startUsage: "500000"},
	}

	buckets := usage.CalculateTierBuckets(decimal.NewFromFloat(*r.MonthlyRequestsMillions), tierLimits)

	for i, bucket := range buckets {
		if i >= len(tiers) {
			break
		}

		if bucket.IsZero() {
			continue
		}

		costComponents = append(costComponents, &schema.CostComponent{
			Name:            tiers[i].name,
			Unit:            "10k requests",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(bucket.Mul(decimal.NewFromInt(100))),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				ProductFamily: strPtr("Networking"),
				Service:       strPtr("Azure Front Door Service"),
				Region:        strPtr(zone),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "productName", Value: strPtr("Azure Front Door")},
					{Key: "skuName", Value: strPtr(tier)},
					{Key: "meterName", Value: strPtr(fmt.Sprintf("%s Requests", tier))},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption:   strPtr("Consumption"),
				StartUsageAmount: strPtr(tiers[i].startUsage),
			},
			UsageBased: true,
		})
	}

	return costComponents
}

func mapRegionToFrontDoorZone(region string) string {
	r := strings.ToLower(region)
	if strings.HasPrefix(r, "usgov") || strings.HasPrefix(r, "usdod") {
		return "Zone 9"
	}

	switch r {
	case "eastus", "eastus2", "westus", "westus2", "westus3", "centralus", "northcentralus", "southcentralus", "westcentralus", "canadacentral", "canadaeast", "mexicocentral", "centraluseuap", "eastus2euap":
		return "Zone 1"
	case "eastasia", "southeastasia", "japanwest", "japaneast":
		return "Zone 2"
	case "brazilsouth", "brazilsoutheast":
		return "Zone 3"
	case "australiaeast", "australiasoutheast", "australiacentral", "australiacentral2":
		return "Zone 4"
	case "southindia", "centralindia", "westindia", "jioindiawest", "jioindiacentral":
		return "Zone 5"
	case "northeurope", "westeurope", "francecentral", "francesouth", "germanywestcentral", "germanynorth", "norwayeast", "norwaywest", "switzerlandnorth", "switzerlandwest", "uksouth", "ukwest", "swedencentral", "polandcentral", "italynorth":
		return "Zone 6"
	case "uaenorth", "uaecentral", "southafricanorth", "southafricawest", "qatarcentral", "israelcentral":
		return "Zone 7"
	case "koreacentral", "koreasouth":
		return "Zone 8"
	default:
		return "Zone 1"
	}
}
