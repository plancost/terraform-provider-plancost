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

func TestMariaDBServer(t *testing.T) {
	t.Skip("`azurerm_mariadb_server` resource is removed from the latest provider")
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	opts := testcase.DefaultGoldenFileOptions()
	// Ignore the CLI because the resource has been removed from the provider in favour of azurerm_mysql_flexible_server
	opts.IgnoreCLI = true

	testcase.GoldenFileResourceTestsWithOpts(t, "mariadb_server_test", opts)
}
