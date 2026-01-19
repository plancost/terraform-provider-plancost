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

func TestAppConfiguration(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	testcase.GoldenFileResourceTestsWithOpts(t, "app_configuration_test", &testcase.GoldenFileOptions{
		Currency:    "USD",
		CaptureLogs: false,
		// ignore the CLI as this throws errors for the test case with an empty sku we
		// want to test for this case as well as this is a valid case for vscode cli
		// users who are yet to run the terraform plan/apply step.
		IgnoreCLI: true,
	})
}
