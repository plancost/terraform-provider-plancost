// Copyright 2021 Infracost Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azurerm_test

import (
	"testing"

	"github.com/plancost/terraform-provider-plancost/internal/testcase"
)

func TestMSSQLDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	opts := testcase.DefaultGoldenFileOptions()
	opts.CaptureLogs = true
	opts.IgnoreCLI = true

	t.Run("base price", func(t *testing.T) {
		testcase.GoldenFileResourceTestsWithOpts(t, "mssql_database_test", opts)
	})

	t.Run("dev/test price", func(t *testing.T) {
		opts.GoldenFileSuffix = "dev_test_price"
		testcase.GoldenFileResourceTestsWithOpts(t, "mssql_database_test", opts)
	})
}

func TestMSSQLDatabaseWithBlankLocation(t *testing.T) {
	t.Skip("Error: Resource Group \"test\" was not found")
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	opts := testcase.DefaultGoldenFileOptions()
	opts.CaptureLogs = true

	testcase.GoldenFileHCLResourceTestsWithOpts(t, "mssql_database_test_with_blank_location", opts)
}
