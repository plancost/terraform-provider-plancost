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
	"github.com/shopspring/decimal"
)

// MachineLearningComputeCluster struct represents a Azure Machine Learning Compute Cluster.
//
// These use the same pricing as Azure Linux Virtual Machines. We default to the minimum scale of
// the cluster, but allow the number of instances and monthly hours of each instance to be set.
//
// Resource information: https://azure.microsoft.com/en-gb/pricing/details/machine-learning/#overview
// Pricing information: https://azure.microsoft.com/en-gb/pricing/details/machine-learning/
type MachineLearningComputeCluster struct {
	Address      string
	Region       string
	InstanceType string
	MinNodeCount int64
	Instances    *int64   `infracost_usage:"instances"`
	MonthlyHours *float64 `infracost_usage:"monthly_hrs"`
}

// CoreType returns the name of this resource type
func (r *MachineLearningComputeCluster) CoreType() string {
	return "MachineLearningComputeCluster"
}

// UsageSchema defines a list which represents the usage schema of MachineLearningComputeCluster.
func (r *MachineLearningComputeCluster) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{
		{Key: "instances", ValueType: schema.Int64, DefaultValue: 0},
		{Key: "monthly_hrs", ValueType: schema.Float64, DefaultValue: 0},
	}
}

// PopulateUsage parses the u schema.UsageData into the MachineLearningComputeCluster.
// It uses the `infracost_usage` struct tags to populate data into the MachineLearningComputeCluster.
func (r *MachineLearningComputeCluster) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid MachineLearningComputeCluster struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *MachineLearningComputeCluster) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		linuxVirtualMachineCostComponent(r.Region, r.InstanceType, r.MonthlyHours),
	}

	res := &schema.Resource{
		Name:           r.Address,
		UsageSchema:    r.UsageSchema(),
		CostComponents: costComponents,
	}

	instances := r.MinNodeCount

	// If the user has set the monthly hours, but the min node count is 0,
	// we assume that the user wants to calculate the cost of 1 instance.
	if r.MonthlyHours != nil && instances == 0 {
		instances = 1
	}

	if r.Instances != nil {
		instances = *r.Instances
	}

	schema.MultiplyQuantities(res, decimal.NewFromInt(instances))

	return res
}
