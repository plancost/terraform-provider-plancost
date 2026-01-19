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
	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
	"github.com/plancost/terraform-provider-plancost/internal/usage"
)

// RecoveryServicesVault struct represents a storage vault that can azure users can back up
// various vms into.
//
// See the ProtectedVM struct for more information about backup services are charged.
//
// Resource information: https://learn.microsoft.com/en-us/azure/backup/backup-overview
// Pricing information: https://azure.microsoft.com/en-gb/pricing/details/backup/
type RecoveryServicesVault struct {
	Address      string
	Region       string
	ProtectedVMs []*BackupProtectedVM
}

func (r *RecoveryServicesVault) CoreType() string {
	return "RecoveryServicesVault"
}

// UsageSchema dynamically constructs a list of UsageItems based on the ProtectedVM sub resources.
func (r *RecoveryServicesVault) UsageSchema() []*schema.UsageItem {
	items := make([]*schema.UsageItem, len(r.ProtectedVMs))
	for i, pm := range r.ProtectedVMs {
		items[i] = &schema.UsageItem{
			Key:          pm.Address,
			DefaultValue: &usage.ResourceUsage{Name: pm.Address, Items: pm.UsageSchema()},
			ValueType:    schema.SubResourceUsage,
		}
	}

	return items
}

// PopulateUsage parses the u schema.UsageData into the RecoveryServicesVault's sub resources.
//
// RecoveryServicesVault does not have any actual usage associated with itself and instead relies on
// users specifying usage for child ProtectedVM resources.
func (r *RecoveryServicesVault) PopulateUsage(u *schema.UsageData) {
	if u == nil {
		return
	}

	// build a new UsageMap so that we get the wildcard support.
	data := map[string]*schema.UsageData{}
	for s, result := range u.Attributes {
		data[s] = schema.NewUsageData(s, result.Map())
	}
	um := schema.NewUsageMap(data)

	for _, pm := range r.ProtectedVMs {
		pm.PopulateUsage(um.Get(pm.Address))
	}

	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid RecoveryServicesVault struct.
//
// RecoveryServicesVault does not have any top level costs associated with it and instead returns a
// list of sub resources where the costs are encapsulated.
func (r *RecoveryServicesVault) BuildResource() *schema.Resource {
	if len(r.ProtectedVMs) == 0 {
		logging.Logger.Warn().Msgf("recovery services vault %s has been marked as free as no associated protected VMs were found", r.Address)
		return &schema.Resource{Name: r.Address, NoPrice: true}
	}

	subResources := make([]*schema.Resource, len(r.ProtectedVMs))
	for i, pvm := range r.ProtectedVMs {
		subResources[i] = pvm.BuildResource()
	}

	return &schema.Resource{
		Name:         r.Address,
		UsageSchema:  r.UsageSchema(),
		SubResources: subResources,
	}
}
