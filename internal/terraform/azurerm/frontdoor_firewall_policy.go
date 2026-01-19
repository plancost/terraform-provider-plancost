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
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

// getFrontdoorFirewallPolicyRegistryItem returns a registry item for the
// resource
func getFrontdoorFirewallPolicyRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_frontdoor_firewall_policy",
		CoreRFunc: newFrontdoorFirewallPolicy,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}

// newFrontdoorFirewallPolicy parses Terraform's data and uses it to build
// a new resource
func newFrontdoorFirewallPolicy(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	if strings.HasPrefix(strings.ToLower(region), "usgov") {
		region = "US Gov Zone 1"
	} else {
		region = regionToCDNZone(region)
	}

	customRules := 0
	if rules := d.Get("custom_rule"); rules.Exists() {
		customRules = len(rules.Array())
	}

	managedRulesets := 0
	if rules := d.Get("managed_rule"); rules.Exists() {
		managedRulesets = len(rules.Array())
	}

	r := &azure.FrontdoorFirewallPolicy{
		Address:         d.Address,
		Region:          region,
		CustomRules:     customRules,
		ManagedRulesets: managedRulesets,
	}
	return r
}
