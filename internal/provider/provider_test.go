/*
Copyright (c) 2026 Plancost.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/plancost/terraform-provider-plancost/internal/testcase"
)

func TestAccProvider_MissingAPIKey_InCI(t *testing.T) {
	// Unset the environment variable to ensure the test fails if not provided in config
	t.Setenv("PLANCOST_API_KEY", "")
	// Simulate CI environment
	t.Setenv("CI", "true")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testcase.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "plancost" {
  api_key = ""
}
resource "plancost_estimate" "this" {
  working_directory = "/tmp"
}
				`,
				ExpectError: regexp.MustCompile("Missing API Key in CI/CD Environment"),
			},
		},
	})
}

func TestAccProvider_MissingAPIKey_Local(t *testing.T) {
	// Unset the environment variable
	t.Setenv("PLANCOST_API_KEY", "")
	// Ensure CI is not set
	t.Setenv("CI", "")
	t.Setenv("GITHUB_ACTIONS", "")
	t.Setenv("GITLAB_CI", "")
	t.Setenv("TF_BUILD", "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testcase.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "plancost" {
  api_key = ""
}
resource "plancost_estimate" "this" {
  working_directory = "/tmp"
}
				`,
			},
		},
	})
}
