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

	"github.com/fatih/camelcase"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getMSSQLElasticPoolRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_mssql_elasticpool",
		CoreRFunc: newMSSQLElasticPool,
		ReferenceAttributes: []string{
			"server_name",
			"resource_group_name",
		},
		GetRegion: func(defaultRegion string, d *schema.ResourceData) string {
			return lookupRegion(d, []string{"server_name", "resource_group_name"})
		},
	}
}

func newMSSQLElasticPool(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	sku := d.Get("sku.0.name").String()
	capacity := d.Get("sku.0.capacity").Int()
	tier := strings.Join(camelcase.Split(d.Get("sku.0.tier").String()), " ")
	family := fmt.Sprintf("Compute %s", d.Get("sku.0.family").String())

	var maxSizeGB float64
	if !d.IsEmpty("max_size_gb") {
		maxSizeGB = d.Get("max_size_gb").Float()
	}
	if !d.IsEmpty("max_size_bytes") {
		maxSizeGB = d.Get("max_size_bytes").Float() / 1024.0 / 1024.0 / 1024.0
	}

	licenseType := d.GetStringOrDefault("license_type", "LicenseIncluded")

	r := &azure.MSSQLElasticPool{
		Address:       d.Address,
		Region:        region,
		SKU:           sku,
		Tier:          tier,
		Family:        family,
		LicenseType:   licenseType,
		MaxSizeGB:     &maxSizeGB,
		ZoneRedundant: d.Get("zone_redundant").Bool(),
	}

	s := strings.ToLower(r.SKU)
	if s == "basicpool" || s == "standardpool" || s == "premiumpool" {
		r.DTUCapacity = &capacity
	} else {
		r.Cores = &capacity
	}

	return r
}
