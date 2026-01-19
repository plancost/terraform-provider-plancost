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

// MachineLearningComputeInstance struct represents a Azure Machine Learning Compute Instance.
//
// These use the same pricing as Azure Linux Virtual Machines.
//
// Resource information: https://azure.microsoft.com/en-gb/pricing/details/machine-learning/#overview
// Pricing information: https://azure.microsoft.com/en-gb/pricing/details/machine-learning/
type MachineLearningComputeInstance struct {
	Address      string
	Region       string
	InstanceType string
	MonthlyHours *float64 `infracost_usage:"monthly_hrs"`
}

// CoreType returns the name of this resource type
func (r *MachineLearningComputeInstance) CoreType() string {
	return "MachineLearningComputeInstance"
}

// UsageSchema defines a list which represents the usage schema of MachineLearningComputeInstance.
func (r *MachineLearningComputeInstance) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "monthly_hrs", ValueType: schema.Float64, DefaultValue: 0},
	}
}

// PopulateUsage parses the u schema.UsageData into the MachineLearningComputeInstance.
// It uses the `infracost_usage` struct tags to populate data into the MachineLearningComputeInstance.
func (r *MachineLearningComputeInstance) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid MachineLearningComputeInstance struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *MachineLearningComputeInstance) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		linuxVirtualMachineCostComponent(r.Region, r.InstanceType, r.MonthlyHours),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}
}
