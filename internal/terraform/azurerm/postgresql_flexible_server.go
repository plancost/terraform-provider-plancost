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
	"regexp"
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getPostgreSQLFlexibleServerRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_postgresql_flexible_server",
		CoreRFunc: newPostgreSQLFlexibleServer,
	}
}

func newPostgreSQLFlexibleServer(d *schema.ResourceData) schema.CoreResource {
	region := d.Region
	sku := d.Get("sku_name").String()
	storage := d.Get("storage_mb").Int()

	var tier, size, version string

	s := strings.Split(sku, "_")
	if len(s) < 3 || len(s) > 4 {
		logging.Logger.Warn().Msgf("Unrecognised PostgreSQL Flexible Server SKU format for resource %s: %s", d.Address, sku)
		return nil
	}

	if len(s) > 2 {
		tier = strings.ToLower(s[0])
		size = s[2]
	}

	if len(s) > 3 {
		version = s[3]
	}

	supportedTiers := []string{"b", "gp", "mo"}
	if !contains(supportedTiers, tier) {
		logging.Logger.Warn().Msgf("Unrecognised PostgreSQL Flexible Server tier prefix for resource %s: %s", d.Address, sku)
		return nil
	}

	if tier != "b" {
		coreRegex := regexp.MustCompile(`(\d+)`)
		match := coreRegex.FindStringSubmatch(size)
		if len(match) < 1 {
			logging.Logger.Warn().Msgf("Unrecognised PostgreSQL Flexible Server size for resource %s: %s", d.Address, sku)
			return nil
		}
	}

	// Check if high availability is enabled
	highAvailability := false
	if ha := d.Get("high_availability"); ha.Exists() {
		if mode := ha.Get("mode"); mode.Exists() && mode.String() != "" {
			highAvailability = true
		}
	}

	r := &azure.PostgreSQLFlexibleServer{
		Address:          d.Address,
		Region:           region,
		SKU:              sku,
		Tier:             tier,
		InstanceType:     size,
		InstanceVersion:  version,
		Storage:          storage,
		HighAvailability: highAvailability,
	}
	return r
}
