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

func getCognitiveDeploymentRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_cognitive_deployment",
		CoreRFunc: newCognitiveDeployment,
		ReferenceAttributes: []string{
			"cognitive_account_id",
		},
		GetRegion: func(defaultRegion string, d *schema.ResourceData) string {
			region := lookupRegion(d, []string{"cognitive_account_id"})

			cognitiveAccountRefs := d.References("cognitive_account_id")
			if region == "" && len(cognitiveAccountRefs) > 0 {
				region = lookupRegion(cognitiveAccountRefs[0], []string{"resource_group_name"})
			}

			return region
		},
	}
}

func newCognitiveDeployment(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	return &azure.CognitiveDeployment{
		Address:  d.Address,
		Region:   region,
		Model:    strings.ToLower(d.Get("model.0.name").String()),
		Version:  strings.ToLower(d.Get("model.0.version").String()),
		Tier:     strings.ToLower(d.Get("scale.0.tier").String()),
		SKU:      strings.ToLower(d.Get("sku.0.name").String()),
		Capacity: d.Get("sku.0.capacity").Int(),
	}
}
