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

type ActiveDirectoryDomainServiceReplicaSet struct {
	Address            string
	Region             string
	DomainServiceIDSKU string
}

func (r *ActiveDirectoryDomainServiceReplicaSet) CoreType() string {
	return "ActiveDirectoryDomainServiceReplicaSet"
}

func (r *ActiveDirectoryDomainServiceReplicaSet) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{}
}

func (r *ActiveDirectoryDomainServiceReplicaSet) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *ActiveDirectoryDomainServiceReplicaSet) BuildResource() *schema.Resource {
	region := r.Region

	costComponents := activeDirectoryDomainServiceCostComponents("Active directory domain service replica set", region, r.DomainServiceIDSKU)

	return &schema.Resource{
		Name:           r.Address,
		CostComponents: costComponents,
		UsageSchema:    r.UsageSchema(),
	}
}
