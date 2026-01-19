package terraform

import (
	"github.com/plancost/terraform-provider-plancost/internal/schema"
	"github.com/plancost/terraform-provider-plancost/internal/terraform/azurerm"
)

type RegistryItemMap map[string]*schema.RegistryItem

var (
	ResourceRegistryMap = buildResourceRegistryMap()
)

var otherFreeResources = []string{
	"plancost_estimate",
	"random_password",
	"random_string",
	"random_uuid",
	"random_id",
	"random_pet",
	"random_shuffle",
	"random_integer",
	"time_rotating",
	"time_static",
	"time_sleep",
	"time_offset",
	"null_resource",
	"local_file",
	"template_file",
	"archive_file",
	"tls_private_key",
	"tls_self_signed_cert",
	"tls_cert_request",
	"tls_locally_signed_cert",
}

func buildResourceRegistryMap() *RegistryItemMap {
	resourceRegistryMap := make(RegistryItemMap)

	// Merge all resource registries

	for _, registryItem := range azurerm.ResourceRegistry {
		if registryItem.CloudResourceIDFunc == nil {
			registryItem.CloudResourceIDFunc = azurerm.DefaultCloudResourceIDFunc
		}
		resourceRegistryMap[registryItem.Name] = registryItem
		resourceRegistryMap[registryItem.Name].DefaultRefIDFunc = azurerm.GetDefaultRefIDFunc
	}
	for _, registryItem := range createFreeResources(azurerm.FreeResources, azurerm.GetDefaultRefIDFunc, azurerm.DefaultCloudResourceIDFunc) {
		resourceRegistryMap[registryItem.Name] = registryItem
	}
	for _, registryItem := range createFreeResources(otherFreeResources, defaultRefIDFunc, defaultCloudResourceIDFunc) {
		resourceRegistryMap[registryItem.Name] = registryItem
	}

	return &resourceRegistryMap
}

// GetRegion returns the region lookup function for the given resource data type if it exists.
func (r *RegistryItemMap) GetRegion(resourceDataType string) schema.RegionLookupFunc {
	item, ok := (*r)[resourceDataType]
	if ok {
		return item.GetRegion
	}

	return nil
}

func (r *RegistryItemMap) GetReferenceAttributes(resourceDataType string) []string {
	var refAttrs []string
	item, ok := (*r)[resourceDataType]
	if ok {
		refAttrs = item.ReferenceAttributes
	}
	return refAttrs
}

func (r *RegistryItemMap) GetCustomRefIDFunc(resourceDataType string) schema.ReferenceIDFunc {
	item, ok := (*r)[resourceDataType]
	if ok {
		return item.CustomRefIDFunc
	}
	return nil
}

func (r *RegistryItemMap) GetDefaultRefIDFunc(resourceDataType string) schema.ReferenceIDFunc {
	item, ok := (*r)[resourceDataType]
	if ok {
		return item.DefaultRefIDFunc
	}
	return func(d *schema.ResourceData) []string {
		return []string{d.Get("id").String()}
	}
}

func GetUsageOnlyResources() []string {
	r := []string{}
	r = append(r, azurerm.UsageOnlyResources...)
	return r
}

func createFreeResources(l []string, defaultRefsFunc schema.ReferenceIDFunc, resourceIdFunc schema.CloudResourceIDFunc) []*schema.RegistryItem {
	freeResources := make([]*schema.RegistryItem, 0)
	for _, resourceName := range l {
		freeResources = append(freeResources, &schema.RegistryItem{
			Name:                resourceName,
			NoPrice:             true,
			Notes:               []string{"Free resource."},
			DefaultRefIDFunc:    defaultRefsFunc,
			CloudResourceIDFunc: resourceIdFunc,
		})
	}
	return freeResources
}

func defaultRefIDFunc(d *schema.ResourceData) []string {
	return []string{d.Get("id").String()}
}

func defaultCloudResourceIDFunc(d *schema.ResourceData) []string {
	return []string{d.Get("id").String()}
}
