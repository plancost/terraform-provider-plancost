package azurerm

import (
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getCDNFrontDoorProfileRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_cdn_frontdoor_profile",
		CoreRFunc: newCDNFrontDoorProfile,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}

func newCDNFrontDoorProfile(d *schema.ResourceData) schema.CoreResource {
	// The sku_name property is required for this resource
	skuName := d.Get("sku_name").String()

	return &azure.CDNFrontDoorProfile{
		Address: d.Address,
		Region:  d.Region,
		SkuName: skuName,
	}
}
