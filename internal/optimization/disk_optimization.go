package optimization

import (
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type diskOptimization struct{}

func (r diskOptimization) Description() string {
	return "Suggests optimizations for Azure Managed Disks."
}

func (r diskOptimization) Apply(input *schema.CoreResource) *string {
	if input == nil {
		return nil
	}
	resource, ok := (*input).(*azure.ManagedDisk)
	if !ok {
		return nil
	}

	// Disk Storage - consider updating Premium SSD disk to v2
	if strings.EqualFold(resource.DiskType, "Premium_LRS") {
		return stringPtr("Consider updating Premium SSD disk to Premium SSD v2 for better performance and cost")
	}

	return nil
}

var _ Rule = &diskOptimization{}
