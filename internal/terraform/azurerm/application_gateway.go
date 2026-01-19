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

func getApplicationGatewayRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_application_gateway",
		CoreRFunc: NewApplicationGateway,
	}
}

func NewApplicationGateway(d *schema.ResourceData) schema.CoreResource {
	var autoscalingMinCapacity *int64
	if d.Get("autoscale_configuration.0.min_capacity").Exists() {
		autoscalingMinCapacity = intPtr(d.Get("autoscale_configuration.0.min_capacity").Int())
	}

	r := &azure.ApplicationGateway{
		Address:                d.Address,
		SKUName:                d.Get("sku.0.name").String(),
		SKUCapacity:            d.Get("sku.0.capacity").Int(),
		AutoscalingMinCapacity: autoscalingMinCapacity,
		Region:                 d.Region,
	}

	return r
}
