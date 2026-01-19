package optimization

import (
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type vmOptimization struct{}

func (r vmOptimization) Description() string {
	return "Suggests optimizations for Azure Virtual Machines."
}

func (r vmOptimization) Apply(input *schema.CoreResource) *string {
	if input == nil {
		return nil
	}
	resource, ok := (*input).(*azure.VirtualMachine)
	if !ok {
		return nil
	}

	// Virtual Machines - consider using Azure Hybrid Benefit for Windows VMs
	// Check if OS is Windows and LicenseType is not set (or None)
	isWindows := strings.EqualFold(resource.StorageOSDiskOSType, "Windows")
	hasLicense := resource.LicenseType != "" && !strings.EqualFold(resource.LicenseType, "None")

	if isWindows && !hasLicense {
		return stringPtr("Consider using Azure Hybrid Benefit for Windows VMs if you have eligible on-premises licenses")
	}

	return nil
}

var _ Rule = &vmOptimization{}
