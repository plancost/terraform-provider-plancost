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

func getVirtualNetworkPeeringRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_virtual_network_peering",
		CoreRFunc: newVirtualNetworkPeering,
		ReferenceAttributes: []string{
			"virtual_network_name",
			"remote_virtual_network_id",
			"resource_group_name",
		},
		GetRegion: func(defaultRegion string, d *schema.ResourceData) string {
			return lookupRegion(d, []string{"virtual_network_name"})
		},
	}
}

func newVirtualNetworkPeering(d *schema.ResourceData) schema.CoreResource {
	sourceRegion := d.Region
	destinationRegion := lookupRegion(d, []string{"remote_virtual_network_id"})

	sourceZone := virtualNetworkPeeringConvertRegion(sourceRegion)
	destinationZone := virtualNetworkPeeringConvertRegion(destinationRegion)

	r := &azure.VirtualNetworkPeering{
		Address:           d.Address,
		DestinationRegion: destinationRegion,
		SourceRegion:      sourceRegion,
		DestinationZone:   destinationZone,
		SourceZone:        sourceZone,
	}
	return r
}

func virtualNetworkPeeringConvertRegion(region string) string {
	zone := regionToVNETZone(region)

	if strings.HasPrefix(strings.ToLower(region), "china") {
		zone = "CN Zone 1"
	}

	return zone
}
