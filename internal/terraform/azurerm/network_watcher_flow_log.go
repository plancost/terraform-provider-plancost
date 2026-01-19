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

func getNetworkWatcherFlowLogRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_network_watcher_flow_log",
		CoreRFunc: newNetworkWatcherFlowLog,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}

func newNetworkWatcherFlowLog(d *schema.ResourceData) schema.CoreResource {
	if !d.Get("enabled").Bool() {
		return schema.BlankCoreResource{
			Name: d.Address,
			Type: d.Type,
		}
	}

	trafficAnalyticsEnabled := false
	trafficAnalyticsAcceleratedProcessing := false

	if len(d.Get("traffic_analytics").Array()) > 0 {
		trafficAnalyticsEnabled = d.Get("traffic_analytics.0.enabled").Bool()
		trafficAnalyticsAcceleratedProcessing = d.Get("traffic_analytics.0.interval_in_minutes").Int() == int64(10)
	}

	region := d.Region
	return &azure.NetworkWatcherFlowLog{
		Address:                               d.Address,
		Region:                                region,
		TrafficAnalyticsEnabled:               trafficAnalyticsEnabled,
		TrafficAnalyticsAcceleratedProcessing: trafficAnalyticsAcceleratedProcessing,
	}
}
