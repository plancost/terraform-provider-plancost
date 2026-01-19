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

func getAutomationJobScheduleRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_automation_job_schedule",
		CoreRFunc: NewAutomationJobSchedule,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}
func NewAutomationJobSchedule(d *schema.ResourceData) schema.CoreResource {
	r := &azure.AutomationJobSchedule{Address: d.Address, Region: d.Region}
	return r
}
