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
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

// This is a free resource but needs it's own custom registry item to specify the custom ID lookup function.
func getCosmosDBAccountRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:    "azurerm_cosmosdb_account",
		NoPrice: true,
		Notes:   []string{"Free resource."},

		CustomRefIDFunc: func(d *schema.ResourceData) []string {
			return []string{d.Get("name").String()}
		},
	}
}
