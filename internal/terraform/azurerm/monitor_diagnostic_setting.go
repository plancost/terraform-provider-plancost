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

func getMonitorDiagnosticSettingRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_monitor_diagnostic_setting",
		CoreRFunc: newMonitorDiagnosticSetting,
		ReferenceAttributes: []string{
			"target_resource_id",
		},
		GetRegion: func(defaultRegion string, d *schema.ResourceData) string {
			return lookupRegion(d, []string{"target_resource_id"})
		},
	}
}

func newMonitorDiagnosticSetting(d *schema.ResourceData) schema.CoreResource {
	return &azure.MonitorDiagnosticSetting{
		Address: d.Address,
		Region:  d.Region,

		EventHubTarget:        !d.IsEmpty("eventhub_authorization_rule_id"),
		PartnerSolutionTarget: !d.IsEmpty("partner_solution_id"),
		StorageAccountTarget:  !d.IsEmpty("storage_account_id"),
	}
}
