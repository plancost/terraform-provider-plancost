package azurerm

import (
	"github.com/tidwall/gjson"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getDashboardGrafanaRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_dashboard_grafana",
		CoreRFunc: newAzureRMDashboardGrafana,
	}
}

func newAzureRMDashboardGrafana(d *schema.ResourceData) schema.CoreResource {
	sku := "Standard"
	if d.Get("sku").Type != gjson.Null {
		sku = d.Get("sku").String()
	}

	return &azure.DashboardGrafana{
		Address:               d.Address,
		Region:                d.Region,
		SKU:                   sku,
		ZoneRedundancyEnabled: d.Get("zone_redundancy_enabled").Bool(),
	}
}
