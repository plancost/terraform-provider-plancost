package azure

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

// StorageTable struct represents Azure Table Storage.
//
// Resource information: https://azure.microsoft.com/en-gb/pricing/details/storage/tables/
// Pricing information: https://azure.microsoft.com/en-gb/pricing/details/storage/tables/#pricing
type StorageTable struct {
	Address                string
	Region                 string
	AccountKind            string
	AccountReplicationType string
	AccountEncrypted       bool

	MonthlyStorageGB                    *float64 `infracost_usage:"storage_gb"`
	MonthlyTransactions                 *int64   `infracost_usage:"transactions"`
	MonthlyWriteOperations              *int64   `infracost_usage:"write_operations"`
	MonthlyBatchWriteOperations         *int64   `infracost_usage:"batch_write_operations"`
	MonthlyReadOperations               *int64   `infracost_usage:"read_operations"`
	MonthlyScanOperations               *int64   `infracost_usage:"scan_operations"`
	MonthlyListOperations               *int64   `infracost_usage:"list_operations"`
	MonthlyDeleteOperations             *int64   `infracost_usage:"delete_operations"`
	MonthlyGeoReplicationDataTransferGB *float64 `infracost_usage:"monthly_geo_replication_data_transfer_gb"`
}

// CoreType returns the name of this resource type
func (r *StorageTable) CoreType() string {
	return "StorageTable"
}

// UsageSchema defines a list which represents the usage schema of StorageTable.
func (r *StorageTable) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "storage_gb", DefaultValue: 0.0, ValueType: schema.Float64},
		{Key: "transactions", DefaultValue: 0, ValueType: schema.Int64},
		{Key: "write_operations", DefaultValue: 0, ValueType: schema.Int64},
		{Key: "batch_write_operations", DefaultValue: 0, ValueType: schema.Int64},
		{Key: "read_operations", DefaultValue: 0, ValueType: schema.Int64},
		{Key: "scan_operations", DefaultValue: 0, ValueType: schema.Int64},
		{Key: "list_operations", DefaultValue: 0, ValueType: schema.Int64},
		{Key: "delete_operations", DefaultValue: 0, ValueType: schema.Int64},
		{Key: "monthly_geo_replication_data_transfer_gb", DefaultValue: 0.0, ValueType: schema.Float64},
	}
}

// PopulateUsage parses the u schema.UsageData into the StorageTable.
// It uses the `infracost_usage` struct tags to populate data into the StorageTable.
func (r *StorageTable) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid StorageTable struct.
// This method is called after the resource is initialised by an IaC provider.
func (r *StorageTable) BuildResource() *schema.Resource {
	if !r.isAccountKindSupported() {
		logging.Logger.Warn().Msgf("Skipping resource %s. Storage Tables don't support %s accounts", r.Address, r.AccountKind)
		return nil
	}

	if !r.isReplicationTypeSupported() {
		logging.Logger.Warn().Msgf("Skipping resource %s. Storage Tables don't support %s redundancy", r.Address, r.AccountReplicationType)
		return nil
	}

	costComponents := []*schema.CostComponent{
		r.dataStorageCostComponent(),
	}
	costComponents = append(costComponents, r.operationsCostComponents()...)
	costComponents = append(costComponents, r.geoReplicationDataTransferCostComponents()...)

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}
}

func (r *StorageTable) isAccountKindSupported() bool {
	return r.isStorageV1() || r.isStorageV2()
}

func (r *StorageTable) isReplicationTypeSupported() bool {
	var validReplicationTypes []string

	switch {
	case r.isStorageV1():
		validReplicationTypes = []string{"LRS", "GRS", "RA-GRS"}
	case r.isStorageV2():
		validReplicationTypes = []string{"LRS", "ZRS", "GRS", "RA-GRS", "GZRS", "RA-GZRS"}
	}

	if validReplicationTypes != nil {
		return contains(validReplicationTypes, strings.ToUpper(r.AccountReplicationType))
	}

	return true
}

func (r *StorageTable) isStorageV1() bool {
	return strings.EqualFold(r.AccountKind, "storage")
}

func (r *StorageTable) isStorageV2() bool {
	return strings.EqualFold(r.AccountKind, "storagev2")
}

func (r *StorageTable) productName() string {
	return "Tables"
}

func (r *StorageTable) dataStorageCostComponent() *schema.CostComponent {
	var qty *decimal.Decimal
	if r.MonthlyStorageGB != nil {
		qty = decimalPtr(decimal.NewFromFloat(*r.MonthlyStorageGB))
	}

	replication := strings.ToUpper(r.AccountReplicationType)
	replication = strings.ReplaceAll(replication, "RAGRS", "RA-GRS")
	replication = strings.ReplaceAll(replication, "RAGZRS", "RA-GZRS")

	tier := "Standard"
	meterName := fmt.Sprintf("%s Data Stored", replication)
	if r.AccountEncrypted {
		tier = "Account Encrypted"
		if replication == "GZRS" || replication == "RA-GZRS" {
			meterName = fmt.Sprintf("Account Encrypted %s", meterName)
		}
	}

	return &schema.CostComponent{
		Name:            "Capacity",
		Unit:            "GB",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: qty,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Storage"),
			ProductFamily: strPtr("Storage"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr(r.productName())},
				{Key: "skuName", Value: strPtr(fmt.Sprintf("%s %s", tier, replication))},
				{Key: "meterName", Value: strPtr(meterName)},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption:   strPtr("Consumption"),
			StartUsageAmount: strPtr("0"),
		},
		UsageBased: true,
	}
}

func (r *StorageTable) operationsCostComponents() []*schema.CostComponent {
	if !r.AccountEncrypted {
		return []*schema.CostComponent{
			r.operationCostComponent("Transactions", "Read Operations", r.MonthlyTransactions),
		}
	}

	costComponents := []*schema.CostComponent{
		r.operationCostComponent("Batch write operations", "Batch Write Operations", r.MonthlyBatchWriteOperations),
		r.operationCostComponent("Write operations", "Write Operations", r.MonthlyWriteOperations),
		r.operationCostComponent("Read operations", "Read Operations", r.MonthlyReadOperations),
		r.operationCostComponent("Scan operations", "Scan Operations", r.MonthlyScanOperations),
		r.operationCostComponent("List operations", "List Operations", r.MonthlyListOperations),
		r.operationCostComponent("Delete operations", "Delete Operations", r.MonthlyDeleteOperations),
	}

	return costComponents
}

func (r *StorageTable) operationCostComponent(name, meterName string, quantity *int64) *schema.CostComponent {
	var qty *decimal.Decimal
	if quantity != nil {
		qty = decimalPtr(decimal.NewFromInt(*quantity).Div(decimal.NewFromInt(10000)))
	}

	replication := strings.ToUpper(r.AccountReplicationType)
	replication = strings.ReplaceAll(replication, "RAGRS", "RA-GRS")
	replication = strings.ReplaceAll(replication, "RAGZRS", "RA-GZRS")

	if strings.Contains(meterName, "Write") {
		switch replication {
		case "GRS", "RA-GRS":
			meterName = fmt.Sprintf("GRS %s", meterName)
		case "GZRS", "RA-GZRS":
			meterName = fmt.Sprintf("%s %s", replication, meterName)
		}
	}

	tier := "Standard"
	if r.AccountEncrypted {
		tier = "Account Encrypted"
		meterName = fmt.Sprintf("Account Encrypted %s", meterName)
	}

	return &schema.CostComponent{
		Name:            name,
		Unit:            "10k operations",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: qty,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Storage"),
			ProductFamily: strPtr("Storage"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr(r.productName())},
				{Key: "skuName", Value: strPtr(fmt.Sprintf("%s %s", tier, replication))},
				{Key: "meterName", Value: strPtr(meterName)},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption:   strPtr("Consumption"),
			StartUsageAmount: strPtr("0"),
		},
		UsageBased: true,
	}
}

func (r *StorageTable) geoReplicationDataTransferCostComponents() []*schema.CostComponent {
	if contains([]string{"LRS", "ZRS"}, strings.ToUpper(r.AccountReplicationType)) {
		return []*schema.CostComponent{}
	}

	var qty *decimal.Decimal
	if r.MonthlyGeoReplicationDataTransferGB != nil {
		qty = decimalPtr(decimal.NewFromFloat(*r.MonthlyGeoReplicationDataTransferGB))
	}

	return []*schema.CostComponent{
		{
			Name:            "Geo-replication data transfer",
			Unit:            "GB",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: qty,
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(r.Region),
				Service:       strPtr("Storage"),
				ProductFamily: strPtr("Storage"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "productName", Value: strPtr("Storage - Bandwidth")},
					{Key: "skuName", Value: strPtr("Geo-Replication v2")},
					{Key: "meterName", Value: strPtr("Geo-Replication v2 Data Transfer")},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption:   strPtr("Consumption"),
				StartUsageAmount: strPtr("0"),
			},
			UsageBased: true,
		},
	}
}
