package azure

import (
	"strings"

	"github.com/shopspring/decimal"

	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type DashboardGrafana struct {
	Address               string
	Region                string
	SKU                   string
	ZoneRedundancyEnabled bool

	ActiveUsers *int64 `infracost_usage:"active_users"`
}

func (r *DashboardGrafana) CoreType() string {
	return "DashboardGrafana"
}

func (r *DashboardGrafana) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "active_users", DefaultValue: 0, ValueType: schema.Int64},
	}
}

func (r *DashboardGrafana) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *DashboardGrafana) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{}

	// Instance cost (Standard SKU only)
	if strings.EqualFold(r.SKU, "Standard") {
		costComponents = append(costComponents, &schema.CostComponent{
			Name:           "Instance (Standard)",
			Unit:           "hours",
			UnitMultiplier: decimal.NewFromInt(1),
			HourlyQuantity: decimalPtr(decimal.NewFromInt(2)), // 2 units
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(r.Region),
				Service:       strPtr("Azure Grafana Service"),
				ProductFamily: strPtr("Developer Tools"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "meterName", Value: strPtr("Standard Node")},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Consumption"),
			},
		})

		if r.ZoneRedundancyEnabled {
			costComponents = append(costComponents, &schema.CostComponent{
				Name:           "Zone redundancy",
				Unit:           "hours",
				UnitMultiplier: decimal.NewFromInt(1),
				HourlyQuantity: decimalPtr(decimal.NewFromInt(1)), // 1 unit
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("azure"),
					Region:        strPtr(r.Region),
					Service:       strPtr("Azure Grafana Service"),
					ProductFamily: strPtr("Developer Tools"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "meterName", Value: strPtr("Standard Zone Redundancy")},
					},
				},
				PriceFilter: &schema.PriceFilter{
					PurchaseOption: strPtr("Consumption"),
				},
			})
		}
	}

	// Active Users
	var activeUsers *decimal.Decimal
	if r.ActiveUsers != nil {
		activeUsers = decimalPtr(decimal.NewFromInt(*r.ActiveUsers))
	}

	meterNameUser := "Standard User"
	if strings.EqualFold(r.SKU, "Essential") {
		meterNameUser = "Essential User"
	}

	costComponents = append(costComponents, &schema.CostComponent{
		Name:            "Active users",
		Unit:            "users",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: activeUsers,
		UsageBased:      true,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Azure Grafana Service"),
			ProductFamily: strPtr("Developer Tools"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "meterName", Value: strPtr(meterNameUser)},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	})

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}
}
