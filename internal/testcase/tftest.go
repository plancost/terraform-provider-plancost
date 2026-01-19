package testcase

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// NewGoldenFileCheck creates a plan check from a GoldenFile
func NewGoldenFileCheck(goldenFile GoldenFile) plancheck.PlanCheck {
	return &CostCheck{
		ExpectedResources: goldenFile.Resources,
	}
}

// GoldenFile represents the structure of the golden file JSON
type GoldenFile struct {
	Resources []ResourceCost `json:"resources"`
}

// GoldenFileOptions contains options for golden file tests
type GoldenFileOptions struct {
	CaptureLogs      bool
	IgnoreCLI        bool
	Currency         string
	GoldenFileSuffix string
}

// DefaultGoldenFileOptions returns default options for golden file tests
func DefaultGoldenFileOptions() *GoldenFileOptions {
	return &GoldenFileOptions{
		CaptureLogs:      false,
		IgnoreCLI:        false,
		Currency:         "USD",
		GoldenFileSuffix: "",
	}
}

// GoldenFileResourceTests runs golden file tests for a resource
func GoldenFileResourceTests(t *testing.T, testName string) {
	testDataDir := "./testdata"
	testDataDir, err := filepath.Abs(testDataDir)
	if err != nil {
		t.Fatalf("failed to get absolute path of testdata directory: %v", err)
	}
	t.Logf("Running golden file test: %s", path.Join(testDataDir, testName))
	tfFilePath := filepath.Join(testDataDir, testName, testName+".tf")
	jsonFilePath := filepath.Join(testDataDir, testName, testName+".json")
	usageFilePath := filepath.Join(testDataDir, testName, testName+".usage.yml")

	tfContent, err := os.ReadFile(tfFilePath)
	if err != nil {
		t.Fatalf("failed to read tf file: %v", err)
	}

	jsonContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	var goldenFile GoldenFile
	if err := json.Unmarshal(jsonContent, &goldenFile); err != nil {
		t.Fatalf("failed to unmarshal json file: %v", err)
	}

	usageFilePathAttr := ""
	if _, err := os.Stat(usageFilePath); err == nil {
		if absUsagePath, err := filepath.Abs(usageFilePath); err == nil {
			usageFilePathAttr = fmt.Sprintf("  usage_file = \"%s\"\n", absUsagePath)
		}
	}

	config := fmt.Sprintf(`
%s

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
%s
}
`, string(tfContent), usageFilePathAttr)

	Test(t, TestCase{
		SkipInit: true,
		Steps: []TestStep{
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						NewGoldenFileCheck(goldenFile),
					},
				},
			},
		},
	})
}

// GoldenFileResourceTestsWithOpts runs golden file tests with custom options
func GoldenFileResourceTestsWithOpts(t *testing.T, testName string, opts *GoldenFileOptions) {
	GoldenFileResourceTests(t, testName)
}

// GoldenFileHCLResourceTestsWithOpts runs golden file HCL tests with custom options
func GoldenFileHCLResourceTestsWithOpts(t *testing.T, testName string, opts *GoldenFileOptions) {
	GoldenFileResourceTests(t, testName)
}
