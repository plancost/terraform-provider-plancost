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

func getNetworkConnectionMonitorRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_network_connection_monitor",
		CoreRFunc: newNetworkConnectionMonitor,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}

func newNetworkConnectionMonitor(d *schema.ResourceData) schema.CoreResource {
	tests := 0

	for _, testGroup := range d.Get("test_group").Array() {
		if !testGroup.Get("enabled").Exists() || testGroup.Get("enabled").Bool() {
			destinationCount := len(testGroup.Get("destination_endpoints").Array())
			sourceCount := len(testGroup.Get("source_endpoints").Array())
			testConfigCount := len(testGroup.Get("test_configuration_names").Array())
			tests += sourceCount * destinationCount * testConfigCount
		}
	}

	region := d.Region
	return &azure.NetworkConnectionMonitor{
		Address: d.Address,
		Region:  region,
		Tests:   intPtr(int64(tests)),
	}
}
