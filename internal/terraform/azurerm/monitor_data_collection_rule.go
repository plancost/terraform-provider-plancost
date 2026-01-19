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

package azurerm

import (
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getMonitorDataCollectionRuleRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_monitor_data_collection_rule",
		CoreRFunc: newMonitorDataCollectionRule,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}

func newMonitorDataCollectionRule(d *schema.ResourceData) schema.CoreResource {
	region := d.Region
	return &azure.MonitorDataCollectionRule{
		Address: d.Address,
		Region:  region,
	}
}
