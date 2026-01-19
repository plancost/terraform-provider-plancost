package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUsageDefaults(t *testing.T) {
	defaults, err := GetUsageDefaults()
	assert.NoError(t, err)
	assert.NotNil(t, defaults)

	// Test simple value and comment
	t.Run("Simple Value", func(t *testing.T) {
		resType := "azurerm_app_configuration"
		key := "monthly_additional_requests"

		resDefaults, ok := defaults[resType]
		assert.True(t, ok, "Resource type %s not found", resType)

		val, ok := resDefaults[key]
		assert.True(t, ok, "Key %s not found in %s", key, resType)

		assert.Equal(t, "30000", val.Value)
		assert.Equal(t, "Monthly number of requests which are above the included 200,000 per day per replica.", val.Comment)
	})

	// Test nested structure
	t.Run("Nested Structure", func(t *testing.T) {
		resType := "azurerm_linux_virtual_machine"
		parentKey := "os_disk"
		childKey := "monthly_disk_operations"

		resDefaults, ok := defaults[resType]
		assert.True(t, ok, "Resource type %s not found", resType)

		parentVal, ok := resDefaults[parentKey]
		assert.True(t, ok, "Parent key %s not found in %s", parentKey, resType)
		assert.NotNil(t, parentVal.Items)

		childVal, ok := parentVal.Items[childKey]
		assert.True(t, ok, "Child key %s not found in %s -> %s", childKey, resType, parentKey)

		assert.Equal(t, "2000000", childVal.Value)
		assert.Contains(t, childVal.Comment, "Number of disk operations")
	})

	// Test another resource to ensure iteration works
	t.Run("Another Resource", func(t *testing.T) {
		resType := "azurerm_api_management"
		key := "monthly_api_calls"

		resDefaults, ok := defaults[resType]
		assert.True(t, ok, "Resource type %s not found", resType)

		val, ok := resDefaults[key]
		assert.True(t, ok, "Key %s not found in %s", key, resType)

		assert.Equal(t, "10000000", val.Value)
		assert.Contains(t, val.Comment, "Monthly number of api calls")
	})
}
