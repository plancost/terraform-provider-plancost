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

func getActiveDirectoryDomainServiceReplicaSetRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_active_directory_domain_service_replica_set",
		CoreRFunc: NewActiveDirectoryDomainServiceReplicaSet,
		ReferenceAttributes: []string{
			"domain_service_id",
		},
	}
}
func NewActiveDirectoryDomainServiceReplicaSet(d *schema.ResourceData) schema.CoreResource {
	r := &azure.ActiveDirectoryDomainServiceReplicaSet{
		Address: d.Address,
		Region:  d.Region,
	}
	if len(d.References("domain_service_id")) > 0 {
		r.DomainServiceIDSKU = d.References("domain_service_id")[0].Get("sku").String()
	}
	return r
}
