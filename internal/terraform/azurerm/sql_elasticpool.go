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
	"fmt"
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getSQLElasticPoolRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_sql_elasticpool",
		CoreRFunc: newSQLElasticPool,
		ReferenceAttributes: []string{
			"server_name",
			"resource_group_name",
		},
		GetRegion: func(defaultRegion string, d *schema.ResourceData) string {
			return lookupRegion(d, []string{"server_name", "resource_group_name"})
		},
	}
}

func newSQLElasticPool(d *schema.ResourceData) schema.CoreResource {
	tier := d.Get("edition").String()
	sku := fmt.Sprintf("%sPool", strings.ToTitle(tier))
	dtu := d.Get("dtu").Int()

	region := d.Region
	r := &azure.MSSQLElasticPool{
		Address:       d.Address,
		Region:        region,
		SKU:           sku,
		Family:        "",
		Tier:          tier,
		DTUCapacity:   &dtu,
		LicenseType:   "LicenseIncluded",
		ZoneRedundant: d.Get("zone_redundant").Bool(),
	}

	if !d.IsEmpty("pool_size") {
		maxSizeGB := d.Get("pool_size").Float() / 1024.0
		r.MaxSizeGB = &maxSizeGB
	}

	return r
}
