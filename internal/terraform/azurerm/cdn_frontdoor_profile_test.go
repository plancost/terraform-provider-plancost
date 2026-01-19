package azurerm_test

import (
	"testing"

	"github.com/plancost/terraform-provider-plancost/internal/testcase"
)

func TestCDNFrontDoorProfileGoldenFile(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	testcase.GoldenFileResourceTestsWithOpts(t, "cdn_frontdoor_profile_test", &testcase.GoldenFileOptions{
		IgnoreCLI: true,
	})
}
