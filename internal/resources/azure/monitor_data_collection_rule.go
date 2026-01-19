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

// MonitorDataCollectionRule struct represents an Azure Monitor Data Collection Rule.
//
// Resource information: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/monitor_data_collection_rule
// Pricing information: https://azure.microsoft.com/en-in/pricing/details/monitor/
type MonitorDataCollectionRule struct {
	Address string
	Region  string

	MonthlyCustomMetricsSamplesGB *int64 `infracost_usage:"monthly_custom_metrics_samples"`
}

// CoreType returns the name of this resource type
func (r *MonitorDataCollectionRule) CoreType() string {
	return "MonitorDataCollectionRule"
}

// UsageSchema defines a list which represents the usage schema of MonitorDataCollectionRule.
func (r *MonitorDataCollectionRule) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "monthly_custom_metrics_samples", ValueType: schema.Int64, DefaultValue: 0},
	}
}

// PopulateUsage parses the u schema.UsageData into the MonitorDataCollectionRule.
// It uses the `infracost_usage` struct tags to populate data into the MonitorDataCollectionRule.
func (r *MonitorDataCollectionRule) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid MonitorDataCollectionRule struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *MonitorDataCollectionRule) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.metricsIngestionCostComponent(r.MonthlyCustomMetricsSamplesGB),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}
}

func (r *MonitorDataCollectionRule) metricsIngestionCostComponent(quantity *int64) *schema.CostComponent {
	var q *decimal.Decimal
	if quantity != nil {
		q = decimalPtr(decimal.NewFromInt(*quantity).Div(decimal.NewFromInt(10000000)))
	}

	return &schema.CostComponent{
		Name:            "Metrics ingestion",
		Unit:            "10M samples",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Azure Monitor"),
			ProductFamily: strPtr("Management and Governance"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "skuName", Value: strPtr("Metrics ingestion")},
				{Key: "meterName", Value: strPtr("Metrics ingestion Metric samples")},
			},
		},
		UsageBased: true,
	}
}
