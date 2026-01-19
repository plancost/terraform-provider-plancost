package azurerm

import (
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getStaticWebAppRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_static_web_app",
		CoreRFunc: newStaticWebApp,
	}
}

func newStaticWebApp(d *schema.ResourceData) schema.CoreResource {
	return &azure.StaticWebApp{
		Address: d.Address,
		Region:  d.Region,
		SkuTier: d.Get("sku_tier").String(),
		SkuSize: d.Get("sku_size").String(),
	}
}
