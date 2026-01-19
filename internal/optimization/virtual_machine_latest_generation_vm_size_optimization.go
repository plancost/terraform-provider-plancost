package optimization

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type latestGenerationVMSizeOptimization struct{}

func (l latestGenerationVMSizeOptimization) Description() string {
	return "Suggests using the latest generation of VM sizes for better performance and cost efficiency."
}

func (l latestGenerationVMSizeOptimization) Apply(input *schema.CoreResource) *string {
	if input == nil {
		return nil
	}
	resource, ok := (*input).(*azure.VirtualMachine)
	if !ok {
		return nil
	}

	vmSize := resource.VMSize

	// Specific upgrades
	if strings.HasPrefix(vmSize, "Standard_G") && !strings.HasPrefix(vmSize, "Standard_GS") {
		return stringPtr("Consider upgrading G series machines to Ev5")
	}
	if strings.HasPrefix(vmSize, "Standard_GS") {
		return stringPtr("Consider upgrading GS series machines to Esv5 or Mv3")
	}
	if strings.HasPrefix(vmSize, "Standard_ND") {
		return stringPtr("Consider upgrading ND series machines to NCas T4 v3 or NDamsr A100 v4")
	}

	// Generic latest generation checks
	family := getVMFamily(vmSize)
	if family != nil {
		version := getVMVersion(vmSize)
		if version < family.Version {
			return stringPtr(fmt.Sprintf("Consider using the latest generation version %d of %s series for better performance and cost efficiency", family.Version, family.Family))
		}
	}

	return nil
}

func getVMVersion(vmSize string) int {
	re := regexp.MustCompile(`_v(\d+)$`)
	matches := re.FindStringSubmatch(vmSize)
	if len(matches) > 1 {
		v, _ := strconv.Atoi(matches[1])
		return v
	}
	return 1
}

func getVMFamily(vmSize string) *VMInfo {
	if vmSize == "" {
		return nil
	}
	for _, info := range latestGenerations {
		if strings.HasPrefix(vmSize, info.NamePrefix) {
			return &info
		}
	}
	return nil
}

type VMInfo struct {
	NamePrefix string
	Version    int
	Family     string
}

var latestGenerations = []VMInfo{
	{
		NamePrefix: "Standard_DCa",
		Version:    6,
		Family:     "DCa",
	},
	{
		NamePrefix: "Standard_DC",
		Version:    3,
		Family:     "DC",
	},
	{
		NamePrefix: "Standard_Da",
		Version:    7,
		Family:     "Da",
	},
	{
		NamePrefix: "Standard_DS",
		Version:    5,
		Family:     "DS",
	},
	{
		NamePrefix: "Standard_D",
		Version:    5,
		Family:     "D",
	},
	{
		NamePrefix: "Standard_Ea",
		Version:    7,
		Family:     "Ea",
	},
	{
		NamePrefix: "Standard_E",
		Version:    5,
		Family:     "E",
	},
	{
		NamePrefix: "Standard_Fa",
		Version:    7,
		Family:     "Fa",
	},
	{
		NamePrefix: "Standard_F",
		Version:    2,
		Family:     "F",
	},
	{
		NamePrefix: "Standard_A",
		Version:    2,
		Family:     "A",
	},
	{
		NamePrefix: "Standard_L",
		Version:    3,
		Family:     "Ls",
	},
}

func stringPtr(s string) *string {
	return &s
}

var _ Rule = &latestGenerationVMSizeOptimization{}
