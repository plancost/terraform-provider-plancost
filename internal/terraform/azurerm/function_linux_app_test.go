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

func TestAzureRMLinuxAppFunctionGoldenFile(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	opts := testcase.DefaultGoldenFileOptions()
	// ignore the CLI as the Terraform provider cannot traverse the each.value references correctly
	// meaning that the HCL provider is more accurate here.
	opts.IgnoreCLI = true
	testcase.GoldenFileResourceTestsWithOpts(t, "function_linux_app_test", opts)
}
