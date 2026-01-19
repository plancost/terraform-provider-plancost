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

// DataFactory struct represents Azure Data Factory resource.
//
// Resource information: https://azure.microsoft.com/en-us/services/data-factory/
// Pricing information: https://azure.microsoft.com/en-us/pricing/details/data-factory/data-pipeline/
type DataFactory struct {
	Address string
	Region  string

	// "usage" args
	MonthlyReadWriteOperationEntities  *int64 `infracost_usage:"monthly_read_write_operation_entities"`
	MonthlyMonitoringOperationEntities *int64 `infracost_usage:"monthly_monitoring_operation_entities"`
}

func (r *DataFactory) CoreType() string {
	return "DataFactory"
}

func (r *DataFactory) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "monthly_read_write_operation_entities", DefaultValue: 0, ValueType: schema.Int64},
		{Key: "monthly_monitoring_operation_entities", DefaultValue: 0, ValueType: schema.Int64},
	}
}

// PopulateUsage parses the u schema.UsageData into the DataFactory.
// It uses the `infracost_usage` struct tags to populate data into the DataFactory.
func (r *DataFactory) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid DataFactory struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *DataFactory) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.readWriteOperationsCostComponent(),
		r.monitoringOperationsCostComponent(),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}
}

// readWriteOperationsCostComponent returns a cost component for
// Data Factory's Read/Write operations.
//
// The pricing is identical for all integration runtimes.
func (r *DataFactory) readWriteOperationsCostComponent() *schema.CostComponent {
	var quantity *decimal.Decimal
	divider := decimal.NewFromInt(50000)

	if r.MonthlyReadWriteOperationEntities != nil {
		quantity = decimalPtr(decimal.NewFromInt(*r.MonthlyReadWriteOperationEntities).Div(divider))
	}

	return &schema.CostComponent{
		Name:            "Read/Write operations",
		Unit:            "50k entities",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Azure Data Factory v2"),
			ProductFamily: strPtr("Analytics"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "meterName", ValueRegex: regexPtr("^Cloud Read Write Operations$")},
				{Key: "skuName", ValueRegex: regexPtr("^Cloud$")},
				{Key: "productName", ValueRegex: regexPtr("^Azure Data Factory v2$")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
		UsageBased: true,
	}
}

// monitoringOperationsCostComponent returns a cost component for
// Data Factory's Monitoring operations.
//
// The pricing is identical for all integration runtimes.
func (r *DataFactory) monitoringOperationsCostComponent() *schema.CostComponent {
	var quantity *decimal.Decimal
	divider := decimal.NewFromInt(50000)

	if r.MonthlyMonitoringOperationEntities != nil {
		quantity = decimalPtr(decimal.NewFromInt(*r.MonthlyMonitoringOperationEntities).Div(divider))
	}

	return &schema.CostComponent{
		Name:            "Monitoring operations",
		Unit:            "50k entities",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Azure Data Factory v2"),
			ProductFamily: strPtr("Analytics"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "meterName", ValueRegex: regexPtr("^Cloud Monitoring Operations$")},
				{Key: "skuName", ValueRegex: regexPtr("^Cloud$")},
				{Key: "productName", ValueRegex: regexPtr("^Azure Data Factory v2$")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
		UsageBased: true,
	}
}
