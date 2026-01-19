package usage

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	yamlv3 "gopkg.in/yaml.v3"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

const minUsageFileVersion = "0.1"
const maxUsageFileVersion = "0.1"

type UsageFile struct {
	Version string `yaml:"version"`
	// We represent resource type usage in using a YAML node so we have control over the comments
	RawResourceTypeUsage yamlv3.Node `yaml:"resource_type_default_usage"`
	// The raw usage is then parsed into this struct
	ResourceTypeUsages []*ResourceUsage `yaml:"-"`
	// We represent resource usage in using a YAML node so we have control over the comments
	RawResourceUsage yamlv3.Node `yaml:"resource_usage"`
	// The raw usage is then parsed into this struct
	ResourceUsages []*ResourceUsage `yaml:"-"`
}

func LoadUsageFile(path string) (*UsageFile, error) {
	blankUsage := NewBlankUsageFile()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logging.Logger.Debug().Msg("Specified usage file does not exist. Using a blank file")

		return blankUsage, nil
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return blankUsage, errors.Wrapf(err, "Error reading usage file")
	}

	usageFile, err := LoadUsageFileFromString(string(contents))
	if err != nil {
		return blankUsage, errors.Wrapf(err, "Error loading usage file")
	}

	return usageFile, nil
}

func NewBlankUsageFile() *UsageFile {
	usageFile := &UsageFile{
		Version: maxUsageFileVersion,
		RawResourceTypeUsage: yamlv3.Node{
			Kind: yamlv3.MappingNode,
		},
		RawResourceUsage: yamlv3.Node{
			Kind: yamlv3.MappingNode,
		},
	}

	return usageFile
}

func LoadUsageFileFromString(s string) (*UsageFile, error) {
	usageFile := &UsageFile{}

	err := yamlv3.Unmarshal([]byte(s), usageFile)
	if err != nil {
		return usageFile, errors.Wrap(err, "Error parsing usage YAML")
	}

	if !usageFile.checkVersion() {
		return usageFile, fmt.Errorf("invalid usage file version. Supported versions are %s ≤ x ≤ %s", minUsageFileVersion, maxUsageFileVersion)
	}

	err = usageFile.parseResourceUsages()
	if err != nil {
		return usageFile, errors.Wrap(err, "Error loading YAML file")
	}

	return usageFile, nil
}

func (u *UsageFile) ToUsageDataMap() schema.UsageMap {
	return schema.NewUsageMapFromInterface(u.ToMap())
}

func (u *UsageFile) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	for _, resourceUsage := range u.ResourceTypeUsages {
		m[resourceUsage.Name] = resourceUsage.Map()
	}

	for _, resourceUsage := range u.ResourceUsages {
		m[resourceUsage.Name] = resourceUsage.Map()
	}

	return m
}

func (u *UsageFile) checkVersion() bool {
	v := u.Version
	if !strings.HasPrefix(u.Version, "v") {
		v = "v" + u.Version
	}

	minV, _ := version.NewVersion(minUsageFileVersion)
	maxV, _ := version.NewVersion(maxUsageFileVersion)
	currV, err := version.NewVersion(v)
	if err != nil {
		return false
	}

	return currV.GreaterThanOrEqual(minV) && currV.LessThanOrEqual(maxV)
}

func (u *UsageFile) parseResourceUsages() error {
	var err error
	u.ResourceUsages, err = ResourceUsagesFromYAML(u.RawResourceUsage)
	if err != nil {
		return errors.Wrapf(err, "Error parsing usage file")
	}
	u.ResourceTypeUsages, err = ResourceUsagesFromYAML(u.RawResourceTypeUsage)
	if err != nil {
		return errors.Wrapf(err, "Error parsing usage file")
	}
	return nil
}
