package azurerm_test

import (
	"testing"

	"github.com/plancost/terraform-provider-plancost/internal/testcase"
)

func TestStorageTable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	opts := testcase.DefaultGoldenFileOptions()
	opts.CaptureLogs = true
	testcase.GoldenFileResourceTestsWithOpts(t, "storage_table_test", opts)
}
