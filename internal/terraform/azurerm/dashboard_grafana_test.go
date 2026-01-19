package azurerm_test

import (
	"testing"

	"github.com/plancost/terraform-provider-plancost/internal/testcase"
)

func TestAzureDashboardGrafanaGoldenFile(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	testcase.GoldenFileResourceTests(t, "dashboard_grafana_test")
}
