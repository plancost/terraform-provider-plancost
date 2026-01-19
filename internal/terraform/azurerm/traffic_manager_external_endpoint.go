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

func getTrafficManagerExternalEndpointRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_traffic_manager_external_endpoint",
		CoreRFunc: newTrafficManagerExternalEndpoint,
		ReferenceAttributes: []string{
			"profile_id",
		},
		GetRegion: func(defaultRegion string, d *schema.ResourceData) string {
			if len(d.References("profile_id")) > 0 {
				profile := d.References("profile_id")[0]
				return lookupRegion(profile, []string{"resource_group_name"})
			}

			return ""
		},
	}
}

func newTrafficManagerExternalEndpoint(d *schema.ResourceData) schema.CoreResource {
	region := d.Region
	healthCheckInterval := int64(30)
	profileEnabled := false

	if len(d.References("profile_id")) > 0 {
		profile := d.References("profile_id")[0]
		healthCheckInterval = profile.GetInt64OrDefault("monitor_config.0.interval_in_seconds", 30)
		profileEnabled = trafficManagerProfileEnabled(profile)
	}

	return &azure.TrafficManagerEndpoint{
		Address:             d.Address,
		Region:              region,
		ProfileEnabled:      profileEnabled,
		External:            true,
		HealthCheckInterval: healthCheckInterval,
	}
}
