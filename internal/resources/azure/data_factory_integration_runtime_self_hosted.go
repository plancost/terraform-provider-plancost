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
)

// DataFactoryIntegrationRuntimeSelfHosted struct represents Data Factory's
// Self-hosted runtime.
//
// Resource information: https://azure.microsoft.com/en-us/services/data-factory/
// Pricing information: https://azure.microsoft.com/en-us/pricing/details/data-factory/data-pipeline/
type DataFactoryIntegrationRuntimeSelfHosted struct {
	Address string
	Region  string

	// "usage" args
	MonthlyOrchestrationRuns *int64 `infracost_usage:"monthly_orchestration_runs"`
}

func (r *DataFactoryIntegrationRuntimeSelfHosted) CoreType() string {
	return "DataFactoryIntegrationRuntimeSelfHosted"
}

func (r *DataFactoryIntegrationRuntimeSelfHosted) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "monthly_orchestration_runs", DefaultValue: 0, ValueType: schema.Int64},
	}
}

// PopulateUsage parses the u schema.UsageData into the DataFactoryIntegrationRuntimeSelfHosted.
// It uses the `infracost_usage` struct tags to populate data into the DataFactoryIntegrationRuntimeSelfHosted.
func (r *DataFactoryIntegrationRuntimeSelfHosted) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid DataFactoryIntegrationRuntimeSelfHosted struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *DataFactoryIntegrationRuntimeSelfHosted) BuildResource() *schema.Resource {
	runtimeFilter := "Self Hosted"

	costComponents := []*schema.CostComponent{
		dataFactoryOrchestrationCostComponent(r.Region, runtimeFilter, r.MonthlyOrchestrationRuns),
		dataFactoryDataMovementCostComponent(r.Region, runtimeFilter),
		dataFactoryPipelineCostComponent(r.Region, runtimeFilter),
		dataFactoryExternalPipelineCostComponent(r.Region, runtimeFilter),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}
}
