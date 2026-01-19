package azure

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type NetAppVolume struct {
	Address              string
	Region               string
	StorageQuotaInGB     int64
	ReplicationFrequency string
	RemoteVolumeLocation string

	MonthlyReplicatedDataGB *float64 `infracost_usage:"monthly_replicated_data_gb"`
}

func (r *NetAppVolume) CoreType() string {
	return "NetAppVolume"
}

func (r *NetAppVolume) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "monthly_replicated_data_gb", DefaultValue: 0, ValueType: schema.Float64},
	}
}

func (r *NetAppVolume) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *NetAppVolume) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{}

	if r.ReplicationFrequency != "" {
		costComponents = append(costComponents, r.replicationCostComponent())
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}
}

func (r *NetAppVolume) replicationCostComponent() *schema.CostComponent {
	frequency := r.ReplicationFrequency
	if strings.EqualFold(frequency, "10minutes") {
		frequency = "Minutes"
	} else if strings.EqualFold(frequency, "hourly") {
		frequency = "Hours"
	} else if strings.EqualFold(frequency, "daily") {
		frequency = "Days"
	}

	meterNameRegex := fmt.Sprintf("/Replication.*%s$/i", frequency)
	if r.RemoteVolumeLocation != "" {
		sourceShort := netAppRegionIDToShortName(r.Region)
		destShort := netAppRegionIDToShortName(r.RemoteVolumeLocation)

		if sourceShort != "" && destShort != "" {
			meterNameRegex = fmt.Sprintf("/(Cross Region Replication|CRR) - %s to %s %s$/i", sourceShort, destShort, frequency)
		}
	}

	var quantity *decimal.Decimal
	if r.MonthlyReplicatedDataGB != nil {
		q := decimal.NewFromFloat(*r.MonthlyReplicatedDataGB)
		quantity = &q
	} else {
		q := decimal.NewFromInt(r.StorageQuotaInGB)
		quantity = &q
	}

	return &schema.CostComponent{
		Name:            fmt.Sprintf("Replication (%s)", strings.ToLower(r.ReplicationFrequency)),
		Unit:            "GiB",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Azure NetApp Files"),
			ProductFamily: strPtr("Storage"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "meterName", ValueRegex: strPtr(meterNameRegex)},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
		UsageBased: true,
	}
}
