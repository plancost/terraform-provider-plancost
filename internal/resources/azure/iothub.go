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

const (
	iotHubFreeSku = "F1"
)

// IoTHub struct represents an IoT Hub
//
// Resource information: https://azure.microsoft.com/en-us/services/iot-hub/
// Pricing information: https://azure.microsoft.com/en-us/pricing/details/iot-hub/
type IoTHub struct {
	Address  string
	Region   string
	Sku      string
	Capacity int64
}

func (r *IoTHub) CoreType() string {
	return "IoTHub"
}

func (r *IoTHub) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{}
}

func (r *IoTHub) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *IoTHub) BuildResource() *schema.Resource {
	if r.Sku == iotHubFreeSku {
		return &schema.Resource{
			Name:      r.Address,
			IsSkipped: true,
			NoPrice:   true,
		}
	}

	return &schema.Resource{
		Name:           r.Address,
		CostComponents: r.costComponents(),
	}
}

func (r *IoTHub) costComponents() []*schema.CostComponent {
	return r.iotHubCostComponent()
}

func (r *IoTHub) iotHubCostComponent() []*schema.CostComponent {
	costComponents := []*schema.CostComponent{
		{
			Name:            "IoT Hub",
			Unit:            "units",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(r.Capacity)),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(r.Region),
				Service:       strPtr("IoT Hub"),
				ProductFamily: strPtr("Internet of Things"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "skuName", Value: strPtr(r.Sku)},
					{Key: "meterName", ValueRegex: regexPtr("Unit$")},
				},
			},
		},
	}

	return costComponents
}

// IoTHubDPS struct represents an IoT Hub DPS
//
// Resource information: https://azure.microsoft.com/en-us/services/iot-hub/
// Pricing information: https://azure.microsoft.com/en-us/pricing/details/iot-hub/
type IoTHubDPS struct {
	Address string
	Region  string
	Sku     string

	MonthlyOperations *int64 `infracost_usage:"monthly_operations"`
}

func (r *IoTHubDPS) CoreType() string {
	return "IoTHubDPS"
}

func (r *IoTHubDPS) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "monthly_operations", DefaultValue: 0, ValueType: schema.Int64},
	}
}

func (r *IoTHubDPS) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *IoTHubDPS) BuildResource() *schema.Resource {
	var monthlyOperations *decimal.Decimal

	if r.MonthlyOperations != nil {
		value := decimal.NewFromInt((*r.MonthlyOperations))
		monthlyOperations = decimalPtr(value.Div(decimal.NewFromInt(1000)))
	}

	return &schema.Resource{
		Name:        r.Address,
		UsageSchema: r.UsageSchema(),
		CostComponents: []*schema.CostComponent{
			{
				Name:            "Device provisioning",
				Unit:            "1k operations",
				UnitMultiplier:  decimal.NewFromInt(1),
				MonthlyQuantity: monthlyOperations,
				ProductFilter: &schema.ProductFilter{
					VendorName:    strPtr("azure"),
					Region:        strPtr(r.Region),
					Service:       strPtr("IoT Hub"),
					ProductFamily: strPtr("Internet of Things"),
					AttributeFilters: []*schema.AttributeFilter{
						{Key: "skuName", Value: strPtr(r.Sku)},
						{Key: "meterName", ValueRegex: regexPtr("Operations$")},
					},
				},
				UsageBased: true,
			},
		},
	}
}
