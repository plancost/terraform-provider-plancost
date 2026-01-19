/*
Copyright (c) 2026 Plancost.
Licensed under the Business Source License 1.1.
See LICENSE file in the project root for full license information.
*/

package provider

import (
	"strings"
	"testing"

	"github.com/plancost/terraform-provider-plancost/internal/schema"
	"github.com/plancost/terraform-provider-plancost/internal/usage"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestGenerateUsageYAML(t *testing.T) {
	// Create resources
	resources := []*schema.Resource{
		{
			Name:         "azurerm_app_configuration.test",
			ResourceType: "azurerm_app_configuration",
			UsageSchema: []*schema.UsageItem{
				{
					Key:          "monthly_additional_requests",
					DefaultValue: 0, // Code default is 0, but YAML default is 30000
					ValueType:    schema.Int64,
				},
			},
		},
		{
			Name:         "azurerm_linux_virtual_machine.test",
			ResourceType: "azurerm_linux_virtual_machine",
			UsageSchema: []*schema.UsageItem{
				{
					Key:          "monthly_hrs",
					DefaultValue: 730, // Code default is 730, but YAML default is 450
					ValueType:    schema.Int64,
				},
				{
					Key:       "os_disk",
					ValueType: schema.SubResourceUsage,
					DefaultValue: &usage.ResourceUsage{
						Items: []*schema.UsageItem{
							{
								Key:          "monthly_disk_operations",
								DefaultValue: 0, // Code default is 0, but YAML default is 2000000
								ValueType:    schema.Int64,
							},
						},
					},
				},
			},
		},
		{
			Name:         "fake_resource.test",
			ResourceType: "fake_resource",
			UsageSchema: []*schema.UsageItem{
				{
					Key:          "fake_key",
					DefaultValue: 123,
					ValueType:    schema.Int64,
				},
			},
		},
	}

	// Generate YAML
	output, err := GenerateUsageYAML(resources)
	assert.NoError(t, err)
	assert.NotNil(t, output)

	outputStr := string(output)

	// 1. Verify Version
	assert.Contains(t, outputStr, "version: 0.1")

	// 2. Verify Resource Names and Sorting
	// GenerateUsageYAML sorts by Name.
	// Expected order:
	// 1. azurerm_app_configuration.test
	// 2. azurerm_linux_virtual_machine.test
	// 3. fake_resource.test

	idx1 := strings.Index(outputStr, "azurerm_app_configuration.test:")
	idx2 := strings.Index(outputStr, "azurerm_linux_virtual_machine.test:")
	idx3 := strings.Index(outputStr, "fake_resource.test:")

	assert.NotEqual(t, -1, idx1, "azurerm_app_configuration.test not found")
	assert.NotEqual(t, -1, idx2, "azurerm_linux_virtual_machine.test not found")
	assert.NotEqual(t, -1, idx3, "fake_resource.test not found")

	assert.True(t, idx1 < idx2, "azurerm_app_configuration should come before azurerm_linux_virtual_machine")
	assert.True(t, idx2 < idx3, "azurerm_linux_virtual_machine should come before fake_resource")

	// 3. Verify Defaults and Comments for azurerm_app_configuration
	// Default in YAML is 30000
	assert.Contains(t, outputStr, "monthly_additional_requests: 30000")
	assert.Contains(t, outputStr, "# Monthly number of requests which are above the included 200,000 per day per replica.")

	// 4. Verify Defaults and Comments for azurerm_linux_virtual_machine
	// monthly_hrs default in YAML is 450
	assert.Contains(t, outputStr, "monthly_hrs: 450")
	assert.Contains(t, outputStr, "# Monthly number of hours the instance ran for.")

	// Nested os_disk -> monthly_disk_operations: 2000000
	assert.Contains(t, outputStr, "monthly_disk_operations: 2000000")
	assert.Contains(t, outputStr, "# Number of disk operations (writes, reads, deletes) using a unit size of 256KiB.")

	// 5. Verify Fallback for fake_resource
	assert.Contains(t, outputStr, "fake_key: 123")

	// Verify valid YAML structure
	var node yaml.Node
	err = yaml.Unmarshal(output, &node)
	assert.NoError(t, err, "Output should be valid YAML")
}
