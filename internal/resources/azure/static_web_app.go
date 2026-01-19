package azure

import (
	"github.com/shopspring/decimal"

	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type StaticWebApp struct {
	Address string
	Region  string
	SkuTier string
	SkuSize string

	MonthlyDataTransferGB *float64 `infracost_usage:"monthly_data_transfer_gb"`
}

func (r *StaticWebApp) CoreType() string {
	return "StaticWebApp"
}

func (r *StaticWebApp) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "monthly_data_transfer_gb", DefaultValue: 0, ValueType: schema.Float64},
	}
}

func (r *StaticWebApp) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *StaticWebApp) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{}

	if r.SkuTier == "Standard" || r.SkuSize == "Standard" {
		costComponents = append(costComponents, r.standardPlanCostComponent())
		costComponents = append(costComponents, r.dataTransferCostComponent())
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}
}

func (r *StaticWebApp) standardPlanCostComponent() *schema.CostComponent {
	return &schema.CostComponent{
		Name:            "Standard plan",
		Unit:            "app",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Azure App Service"),
			ProductFamily: strPtr("Compute"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr("Static Web Apps")},
				{Key: "skuName", Value: strPtr("Standard")},
				{Key: "meterName", Value: strPtr("Standard App")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	}
}

func (r *StaticWebApp) dataTransferCostComponent() *schema.CostComponent {
	var quantity *decimal.Decimal
	if r.MonthlyDataTransferGB != nil {
		q := decimal.NewFromFloat(*r.MonthlyDataTransferGB)
		// 100GB is free
		q = q.Sub(decimal.NewFromInt(100))
		if q.LessThan(decimal.Zero) {
			q = decimal.Zero
		}
		quantity = &q
	}

	return &schema.CostComponent{
		Name:            "Data transfer (over 100GB)",
		Unit:            "GB",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Azure App Service"),
			ProductFamily: strPtr("Compute"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr("Static Web Apps")},
				{Key: "skuName", Value: strPtr("Standard")},
				{Key: "meterName", Value: strPtr("Standard Bandwidth Usage")},
			},
		},
		UsageBased: true,
		PriceFilter: &schema.PriceFilter{
			PurchaseOption:   strPtr("Consumption"),
			StartUsageAmount: strPtr("100"),
		},
	}
}
