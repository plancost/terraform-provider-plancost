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
	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func GetAzureRMCosmosdbMongoCollectionRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_cosmosdb_mongo_collection",
		RFunc: NewAzureRMCosmosdbMongoCollection,
		ReferenceAttributes: []string{
			"account_name",
			"database_name",
			"resource_group_name",
		},
	}
}

func NewAzureRMCosmosdbMongoCollection(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	if len(d.References("database_name")) > 0 {
		mongoDB := d.References("database_name")[0]
		if len(mongoDB.References("account_name")) > 0 {
			account := mongoDB.References("account_name")[0]
			return &schema.Resource{
				Name:           d.Address,
				CostComponents: cosmosDBCostComponents(d, u, account),
			}
		}
		logging.Logger.Warn().Msgf("Skipping resource %s as its 'database_name.account_name' property could not be found.", d.Address)
		return nil
	}
	logging.Logger.Warn().Msgf("Skipping resource %s as its 'database_name' property could not be found.", d.Address)
	return nil
}
