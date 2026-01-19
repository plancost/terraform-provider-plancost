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
	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"

	"fmt"

	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type DatabricksWorkspace struct {
	Address                        string
	Region                         string
	SKU                            string
	MonthlyAllPurposeComputeDBUHrs *int64 `infracost_usage:"monthly_all_purpose_compute_dbu_hrs"`
	MonthlyJobsComputeDBUHrs       *int64 `infracost_usage:"monthly_jobs_compute_dbu_hrs"`
	MonthlyJobsLightComputeDBUHrs  *int64 `infracost_usage:"monthly_jobs_light_compute_dbu_hrs"`
}

func (r *DatabricksWorkspace) CoreType() string {
	return "DatabricksWorkspace"
}

func (r *DatabricksWorkspace) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "monthly_all_purpose_compute_dbu_hrs", ValueType: schema.Int64, DefaultValue: 0},
		{Key: "monthly_jobs_compute_dbu_hrs", ValueType: schema.Int64, DefaultValue: 0},
		{Key: "monthly_jobs_light_compute_dbu_hrs", ValueType: schema.Int64, DefaultValue: 0},
	}
}

func (r *DatabricksWorkspace) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *DatabricksWorkspace) BuildResource() *schema.Resource {
	var costComponents []*schema.CostComponent

	sku := cases.Title(language.English).String(r.SKU)

	if sku == "Trial" {
		return &schema.Resource{
			Name:        r.Address,
			NoPrice:     true,
			IsSkipped:   true,
			UsageSchema: r.UsageSchema(),
		}
	}

	var allPurpose, jobs, jobsLight *decimal.Decimal

	if r.MonthlyAllPurposeComputeDBUHrs != nil {
		allPurpose = decimalPtr(decimal.NewFromInt(*r.MonthlyAllPurposeComputeDBUHrs))
	}
	costComponents = append(costComponents, r.databricksCostComponent(
		"All-purpose compute DBUs",
		fmt.Sprintf("%s All-purpose Compute", sku),
		allPurpose,
	))

	if r.MonthlyJobsComputeDBUHrs != nil {
		jobs = decimalPtr(decimal.NewFromInt(*r.MonthlyJobsComputeDBUHrs))
	}
	costComponents = append(costComponents, r.databricksCostComponent(
		"Jobs compute DBUs",
		fmt.Sprintf("%s Jobs Compute", sku),
		jobs,
	))

	if r.MonthlyJobsLightComputeDBUHrs != nil {
		jobsLight = decimalPtr(decimal.NewFromInt(*r.MonthlyJobsLightComputeDBUHrs))
	}
	costComponents = append(costComponents, r.databricksCostComponent(
		"Jobs light compute DBUs",
		fmt.Sprintf("%s Jobs Light Compute", sku),
		jobsLight,
	))

	return &schema.Resource{
		Name:           r.Address,
		CostComponents: costComponents,
		UsageSchema:    r.UsageSchema(),
	}
}

func (r *DatabricksWorkspace) databricksCostComponent(name, skuName string, quantity *decimal.Decimal) *schema.CostComponent {
	return &schema.CostComponent{
		Name:            name,
		Unit:            "DBU-hours",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: quantity,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Azure Databricks"),
			ProductFamily: strPtr("Analytics"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "skuName", Value: strPtr(skuName)},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
		UsageBased: true,
	}
}
