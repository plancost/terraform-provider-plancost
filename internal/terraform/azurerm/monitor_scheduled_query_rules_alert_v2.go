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
	duration "github.com/channelmeter/iso8601duration"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getMonitorScheduledQueryRulesAlertV2RegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_monitor_scheduled_query_rules_alert_v2",
		CoreRFunc: newMonitorScheduledQueryRulesAlertV2,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}

func newMonitorScheduledQueryRulesAlertV2(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	freq := int64(1)
	ef, err := duration.FromString(d.Get("evaluation_frequency").String())
	if err != nil {
		logging.Logger.Warn().Str(
			"resource", d.Address,
		).Msgf("failed to parse ISO8061 duration string '%s' using 1 minute frequency", d.Get("evaluation_frequency").String())
	} else {
		freq = int64(ef.ToDuration().Minutes())
	}

	scopeCount := 1 // default scope is the azure subscription, so count == 1
	if !d.IsEmpty("scopes") {
		scopeCount = len(d.Get("scopes").Array())
	}

	criteriaDimensionsCount := 0
	for _, c := range d.Get("criteria").Array() {
		criteriaDimensionsCount += len(c.Get("dimension").Array())
	}

	return &azure.MonitorScheduledQueryRulesAlert{
		Address:          d.Address,
		Region:           region,
		Enabled:          d.GetBoolOrDefault("enabled", true),
		TimeSeriesCount:  int64(scopeCount * criteriaDimensionsCount),
		FrequencyMinutes: freq,
	}
}
