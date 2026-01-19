package azure

import "strings"

var netAppRegionIDs = map[string]string{
	"uaecentral":         "AE Central",
	"uaenorth":           "AE North",
	"eastasia":           "AP East",
	"southeastasia":      "AP Southeast",
	"australiaeast":      "AU East",
	"australiasoutheast": "AU Southeast",
	"australiacentral":   "AU Central",
	"australiacentral2":  "AU Central 2",
	"brazilsouth":        "BR South",
	"brazilsoutheast":    "BR Southeast",
	"canadacentral":      "CA Central",
	"canadaeast":         "CA East",
	"switzerlandnorth":   "CH North",
	"switzerlandwest":    "CH West",
	"germanynorth":       "DE North",
	"germanywestcentral": "DE West Central",
	"spaincentral":       "ES Central",
	"northeurope":        "EU North",
	"westeurope":         "EU West",
	"francecentral":      "FR Central",
	"francesouth":        "FR South",
	"israelcentral":      "IL Central",
	"centralindia":       "IN Central",
	"southindia":         "IN South",
	"italynorth":         "IT North",
	"japaneast":          "JA East",
	"japanwest":          "JA West",
	"koreacentral":       "KR Central",
	"koreasouth":         "KR South",
	"malaysiawest":       "MY West",
	"norwayeast":         "NO East",
	"norwaywest":         "NO West",
	"newzealandnorth":    "NZ North",
	"qatarcentral":       "QA Central",
	"swedencentral":      "SE Central",
	"swedensouth":        "SE South",
	"taiwannorth":        "TW North",
	"uksouth":            "UK South",
	"ukwest":             "UK West",
	"centralus":          "US Central",
	"eastus":             "US East",
	"eastus2":            "US East 2",
	"usgovarizona":       "US Gov AZ",
	"usgovtexas":         "US Gov TX",
	"usgovvirginia":      "US Gov VA",
	"westus":             "US West",
	"westus2":            "US West 2",
	"westus3":            "US West 3",
	"northcentralus":     "US North Central",
	"southcentralus":     "US South Central",
}

// netAppRegionIDToShortName maps the Azure Region ID to the "Short Name" used in the metering log.
func netAppRegionIDToShortName(regionID string) string {
	// Normalize the region ID to handle display names (e.g. "West Europe" -> "westeurope")
	normalized := strings.ToLower(strings.ReplaceAll(regionID, " ", ""))
	if name, ok := netAppRegionIDs[normalized]; ok {
		return name
	}
	return ""
}
