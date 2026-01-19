package azurerm_test

import (
	"testing"

	tftest "github.com/plancost/terraform-provider-plancost/internal/testcase"
)

func TestAzureNetAppVolumeGoldenFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	opts := tftest.DefaultGoldenFileOptions()
	opts.CaptureLogs = true
	tftest.GoldenFileResourceTestsWithOpts(t, "netapp_volume_test", opts)
}
