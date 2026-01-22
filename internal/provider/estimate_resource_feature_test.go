/*
Copyright (c) 2026 Plancost.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package provider_test

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/plancost/terraform-provider-plancost/internal/testcase"
)

func TestAccEstimateResource_Optimization(t *testing.T) {
	if os.Getenv("PLANCOST_API_KEY") == "" {
		t.Skip("PLANCOST_API_KEY environment variable is not set")
	}
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "optimization"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue("plancost_estimate.this", tfjsonpath.New("recommendations"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("Consider using the latest generation version 5 of DS series for better performance and cost efficiency"),
							knownvalue.StringExact("Consider using Azure Hybrid Benefit for Windows VMs if you have eligible on-premises licenses"),
						})),
					},
				},
			},
		},
	})
}

func TestAccEstimateResource_SavingPlans(t *testing.T) {
	if os.Getenv("PLANCOST_API_KEY") == "" {
		t.Skip("PLANCOST_API_KEY environment variable is not set")
	}
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "saving_plans"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue("plancost_estimate.this", tfjsonpath.New("recommendations"), knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"description": knownvalue.StringRegexp(regexp.MustCompile(`Save \$[0-9.]+/mo \([0-9]+%\) on .* with a [135]-year Reservation`)),
							}),
							1: knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"description": knownvalue.StringRegexp(regexp.MustCompile(`Save \$[0-9.]+/mo \([0-9]+%\) on .* with a [135]-year Reservation`)),
							}),
						})),
					},
				},
			},
		},
	})
}

func TestAccEstimateResource_TaggingPolicy(t *testing.T) {
	if os.Getenv("PLANCOST_API_KEY") == "" {
		t.Skip("PLANCOST_API_KEY environment variable is not set")
	}
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "tagging_policy"),
				ExpectError:     regexp.MustCompile("Tagging Policy Violation"),
			},
		},
	})
}

func TestAccEstimateResource_Guardrails(t *testing.T) {
	if os.Getenv("PLANCOST_API_KEY") == "" {
		t.Skip("PLANCOST_API_KEY environment variable is not set")
	}
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "guardrails"),
				ExpectError:     regexp.MustCompile("Guardrail Violation"),
			},
		},
	})
}

func TestAccEstimateResource_SimpleDiscount(t *testing.T) {
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "simple_discount"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						// 3.65 * 0.5 = 1.825 -> 1.83
						plancheck.ExpectKnownValue("plancost_estimate.this", tfjsonpath.New("monthly_cost"), testcase.Float64Exact(1.83)),
					},
				},
			},
		},
	})
}
