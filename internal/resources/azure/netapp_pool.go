package azure

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type NetAppPool struct {
	Address           string
	Region            string
	SizeInTB          int64
	ServiceLevel      string
	EncryptionType    string
	CoolAccessEnabled bool

	CustomThroughputMiBps float64

	CoolAccessDataPercentage      *float64 `infracost_usage:"capacity_in_cool_access_percentage"`
	CoolAccessReadWritePercentage *float64 `infracost_usage:"cool_access_data_read_write_percentage"`
}

func (r *NetAppPool) CoreType() string {
	return "NetAppPool"
}

func (r *NetAppPool) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "capacity_in_cool_access_percentage", DefaultValue: 0, ValueType: schema.Float64},
		{Key: "cool_access_data_read_write_percentage", DefaultValue: 0, ValueType: schema.Float64},
	}
}

func (r *NetAppPool) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *NetAppPool) BuildResource() *schema.Resource {
	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: r.poolCostComponents(),
	}
}

func (r *NetAppPool) poolCostComponents() []*schema.CostComponent {
	totalCapacity := decimal.NewFromInt(r.SizeInTB).Mul(decimal.NewFromInt(1024))
	quantity := totalCapacity.Mul(schema.HourToMonthUnitMultiplier)

	if strings.EqualFold(r.ServiceLevel, "Flexible") {
		costComponents := []*schema.CostComponent{
			{
				Name:            "Storage (flexible)",
				Unit:            "GiB",
				UnitMultiplier:  schema.HourToMonthUnitMultiplier,
				MonthlyQuantity: &quantity,
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("azure"),
					Region:        strPtr(r.Region),
					Service:       strPtr("Azure NetApp Files"),
					ProductFamily: strPtr("Storage"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "skuName", Value: strPtr("Flexible Service Level")},
						{Key: "meterName", Value: strPtr("Flexible Service Level Capacity")},
					},
				},
				PriceFilter: &schema.PriceFilter{
					PurchaseOption: strPtr("Consumption"),
				},
			},
		}

		// 128 MiB/sec per pool is available free of charge
		freeThroughput := 128.0
		if r.CustomThroughputMiBps > freeThroughput {
			excess := decimal.NewFromFloat(r.CustomThroughputMiBps).Sub(decimal.NewFromFloat(freeThroughput))
			excessQuantity := excess.Mul(schema.HourToMonthUnitMultiplier)
			costComponents = append(costComponents, &schema.CostComponent{
				Name:            "Throughput (overage)",
				Unit:            "MiB/s",
				UnitMultiplier:  schema.HourToMonthUnitMultiplier,
				MonthlyQuantity: &excessQuantity,
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("azure"),
					Region:        strPtr(r.Region),
					Service:       strPtr("Azure NetApp Files"),
					ProductFamily: strPtr("Storage"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "skuName", Value: strPtr("Flexible Service Level")},
						{Key: "meterName", Value: strPtr("Flexible Service Level Throughput MiBps")},
					},
				},
				PriceFilter: &schema.PriceFilter{
					PurchaseOption: strPtr("Consumption"),
				},
			})
		}
		return costComponents
	}

	sku := r.ServiceLevel

	// Handle Double Encryption
	if strings.EqualFold(r.EncryptionType, "Double") {
		if r.CoolAccessEnabled {
			logging.Logger.Warn().Msgf("Azure NetApp Files capacity pools with Double Encryption do not support Cool Access. Calculating cost for Double Encryption only for resource %s", r.Address)
		}

		meterNameRegex := fmt.Sprintf("/%s Double Encrypted Capacity/i", sku)

		return []*schema.CostComponent{
			{
				Name:            fmt.Sprintf("Storage (%s, double encryption)", strings.ToLower(r.ServiceLevel)),
				Unit:            "GiB",
				UnitMultiplier:  schema.HourToMonthUnitMultiplier,
				MonthlyQuantity: &quantity,
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("azure"),
					Region:        strPtr(r.Region),
					Service:       strPtr("Azure NetApp Files"),
					ProductFamily: strPtr("Storage"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "skuName", Value: strPtr(fmt.Sprintf("%s Double Encrypted", sku))},
						{Key: "meterName", ValueRegex: strPtr(meterNameRegex)},
					},
				},
				PriceFilter: &schema.PriceFilter{
					PurchaseOption: strPtr("Consumption"),
				},
			},
		}
	}

	// Handle Cool Access
	if r.CoolAccessEnabled {
		coolFactor := decimal.Zero
		if r.CoolAccessDataPercentage != nil {
			coolFactor = decimal.NewFromFloat(*r.CoolAccessDataPercentage).Div(decimal.NewFromInt(100))
		}

		hotCapacity := totalCapacity.Mul(decimal.NewFromInt(1).Sub(coolFactor))
		coolCapacity := totalCapacity.Mul(coolFactor)

		transferFactor := decimal.Zero
		if r.CoolAccessReadWritePercentage != nil {
			transferFactor = decimal.NewFromFloat(*r.CoolAccessReadWritePercentage).Div(decimal.NewFromInt(100))
		}
		transferQuantity := coolCapacity.Mul(transferFactor)

		costComponents := []*schema.CostComponent{}

		// 1. Hot Tier Capacity (Standard/Premium/Ultra)
		//    Uses the regular capacity meter (e.g. "Standard Capacity")
		hotQuantity := hotCapacity.Mul(schema.HourToMonthUnitMultiplier)
		costComponents = append(costComponents, &schema.CostComponent{
			Name:            fmt.Sprintf("Storage (%s)", strings.ToLower(r.ServiceLevel)),
			Unit:            "GiB",
			UnitMultiplier:  schema.HourToMonthUnitMultiplier,
			MonthlyQuantity: &hotQuantity,
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(r.Region),
				Service:       strPtr("Azure NetApp Files"),
				ProductFamily: strPtr("Storage"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "skuName", Value: strPtr(sku)},
					{Key: "meterName", ValueRegex: strPtr(fmt.Sprintf("/%s Capacity/i", sku))},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Consumption"),
			},
		})

		// 2. Cool Tier Capacity
		//    Meter: "{Sku} Storage with Cool Access Capacity"
		//    Note: The pricing JSON shows "Standard Storage with Cool Access Capacity" and SKU "Standard Storage with Cool Access"
		//    We will dynamically construct this for Premium/Ultra as well.
		coolSku := fmt.Sprintf("%s Storage with Cool Access", sku)
		coolQuantity := coolCapacity.Mul(schema.HourToMonthUnitMultiplier)
		costComponents = append(costComponents, &schema.CostComponent{
			Name:            "Storage (cool access)",
			Unit:            "GiB",
			UnitMultiplier:  schema.HourToMonthUnitMultiplier,
			MonthlyQuantity: &coolQuantity,
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(r.Region),
				Service:       strPtr("Azure NetApp Files"),
				ProductFamily: strPtr("Storage"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "skuName", Value: strPtr(coolSku)},
					{Key: "meterName", ValueRegex: strPtr(fmt.Sprintf("/%s Capacity/i", coolSku))},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Consumption"),
			},
			UsageBased: true,
		})

		// 3. Initial data transfer to cool access
		//    Meter: "{Sku} Storage with Cool Access Data Transfer"
		if coolCapacity.GreaterThan(decimal.Zero) {
			costComponents = append(costComponents, &schema.CostComponent{
				Name:            "Initial data transfer to cool access",
				Unit:            "GiB",
				UnitMultiplier:  decimal.NewFromInt(1),
				MonthlyQuantity: &coolCapacity,
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("azure"),
					Region:        strPtr(r.Region),
					Service:       strPtr("Azure NetApp Files"),
					ProductFamily: strPtr("Storage"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "skuName", Value: strPtr(coolSku)},
						{Key: "meterName", ValueRegex: strPtr(fmt.Sprintf("/%s Data Transfer/i", coolSku))},
					},
				},
				PriceFilter: &schema.PriceFilter{
					PurchaseOption: strPtr("Consumption"),
				},
				UsageBased: true,
			})
		}

		// 4. Recurring Data Transfer (Read/Write)
		//    Meter: "{Sku} Storage with Cool Access Data Transfer"
		if transferQuantity.GreaterThan(decimal.Zero) {
			costComponents = append(costComponents, &schema.CostComponent{
				Name:            "Recurring data transfer (cool access)",
				Unit:            "GiB",
				UnitMultiplier:  decimal.NewFromInt(1),
				MonthlyQuantity: &transferQuantity,
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("azure"),
					Region:        strPtr(r.Region),
					Service:       strPtr("Azure NetApp Files"),
					ProductFamily: strPtr("Storage"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "skuName", Value: strPtr(coolSku)},
						{Key: "meterName", ValueRegex: strPtr(fmt.Sprintf("/%s Data Transfer/i", coolSku))},
					},
				},
				PriceFilter: &schema.PriceFilter{
					PurchaseOption: strPtr("Consumption"),
				},
				UsageBased: true,
			})
		}

		return costComponents
	}

	// Standard Capacity (No Cool Access)
	return []*schema.CostComponent{
		{
			Name:            fmt.Sprintf("Storage (%s)", strings.ToLower(r.ServiceLevel)),
			Unit:            "GiB",
			UnitMultiplier:  schema.HourToMonthUnitMultiplier,
			MonthlyQuantity: &quantity,
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(r.Region),
				Service:       strPtr("Azure NetApp Files"),
				ProductFamily: strPtr("Storage"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "skuName", Value: strPtr(sku)},
					{Key: "meterName", ValueRegex: strPtr(fmt.Sprintf("/%s Capacity/i", sku))},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Consumption"),
			},
		},
	}
}
