/*
Copyright (c) 2026 Plancost.
Licensed under the Business Source License 1.1.
See LICENSE file in the project root for full license information.
*/

package provider_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/plancost/terraform-provider-plancost/internal/testcase"
)

func TestAccEstimateResource_Basic(t *testing.T) {
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "basic"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue("plancost_estimate.this", tfjsonpath.New("monthly_cost"), testcase.Float64Exact(3.65)),
						testcase.NewResourceCostCheck([]testcase.ResourceCost{
							{
								Name: "azurerm_public_ip.example",
								CostComponents: []testcase.CostComponent{
									{
										Name:            "IP address (static, regional)",
										MonthlyCost:     "3.65",
										MonthlyQuantity: "730",
										Unit:            "hours",
									},
								},
							},
						}),
					},
				},
			},
		},
	})
}

func TestAccEstimateResource_WithIndex(t *testing.T) {
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "with_index"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue("plancost_estimate.this", tfjsonpath.New("monthly_cost"), testcase.Float64Exact(13.14)),
						testcase.NewResourceCostCheck([]testcase.ResourceCost{
							{
								Name: "azurerm_public_ip.example2[0]",
								CostComponents: []testcase.CostComponent{
									{
										Name:            "IP address (static, regional)",
										MonthlyCost:     "3.65",
										MonthlyQuantity: "730",
										Unit:            "hours",
									},
								},
							},
							{
								Name: "azurerm_public_ip.example2[1]",
								CostComponents: []testcase.CostComponent{
									{
										Name:            "IP address (static, regional)",
										MonthlyCost:     "3.65",
										MonthlyQuantity: "730",
										Unit:            "hours",
									},
								},
							},
							{
								Name: "azurerm_public_ip.example[\"one\"]",
								CostComponents: []testcase.CostComponent{
									{
										Name:            "IP address (dynamic, regional)",
										MonthlyCost:     "2.92",
										MonthlyQuantity: "730",
										Unit:            "hours",
									},
								},
							},
							{
								Name: "azurerm_public_ip.example[\"two\"]",
								CostComponents: []testcase.CostComponent{
									{
										Name:            "IP address (dynamic, regional)",
										MonthlyCost:     "2.92",
										MonthlyQuantity: "730",
										Unit:            "hours",
									},
								},
							},
						}),
					},
				},
			},
		},
	})
}

func TestAccEstimateResource_WithTerraformVariables(t *testing.T) {
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "with_terraform_variables"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue("plancost_estimate.this", tfjsonpath.New("monthly_cost"), testcase.Float64Exact(3.65)),
						testcase.NewResourceCostCheck([]testcase.ResourceCost{
							{
								Name: "azurerm_public_ip.example",
								CostComponents: []testcase.CostComponent{
									{
										Name:            "IP address (static, regional)",
										MonthlyCost:     "3.65",
										MonthlyQuantity: "730",
										Unit:            "hours",
									},
								},
							},
						}),
					},
				},
			},
		},
	})
}

func TestAccEstimateResource_WithModule(t *testing.T) {
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "with_module"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue("plancost_estimate.this", tfjsonpath.New("monthly_cost"), testcase.Float64Exact(315.36)),
						testcase.NewResourceCostCheck([]testcase.ResourceCost{
							{
								Name:           "module.aks.azurerm_kubernetes_cluster.example",
								CostComponents: []testcase.CostComponent{},
								SubResources: []testcase.ResourceCost{
									{
										Name: "default_node_pool",
										SubResources: []testcase.ResourceCost{
											{
												Name: "os_disk",
												CostComponents: []testcase.CostComponent{
													{
														Name:            "Storage (P10, LRS)",
														MonthlyCost:     "59.13",
														MonthlyQuantity: "3",
														Unit:            "months",
													},
												},
											},
										},
										CostComponents: []testcase.CostComponent{
											{
												Name:            "Instance usage (Linux, pay as you go, Standard_D2s_v3)",
												MonthlyCost:     "256.23",
												MonthlyQuantity: "2190",
												Unit:            "hours",
											},
										},
									},
								},
							},
						}),
					},
				},
			},
		},
	})
}

func TestAccEstimateResource_UsageDefaults(t *testing.T) {
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "usage_file_defaults"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						testcase.NewResourceCostCheck([]testcase.ResourceCost{
							{
								Name: "azurerm_storage_account.example",
								CostComponents: []testcase.CostComponent{
									{
										Name:            "Capacity (first 50TB)",
										MonthlyQuantity: "100",
										Unit:            "GB",
										MonthlyCost:     "2.08",
									},
								},
							},
						}),
					},
				},
			},
		},
	})
}

func TestAccEstimateResource_UsageOverrideYaml(t *testing.T) {
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "usage_file_override"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						testcase.NewResourceCostCheck([]testcase.ResourceCost{
							{
								Name: "azurerm_storage_account.example",
								CostComponents: []testcase.CostComponent{
									{
										Name:            "Capacity (first 50TB)",
										MonthlyQuantity: "200",
										Unit:            "GB",
										MonthlyCost:     "4.16",
									},
								},
							},
						}),
					},
				},
			},
		},
	})
}

func TestAccEstimateResource_UsageOverrideHcl(t *testing.T) {
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "usage_hcl_override"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
						testcase.NewResourceCostCheck([]testcase.ResourceCost{
							{
								Name: "azurerm_storage_account.example",
								CostComponents: []testcase.CostComponent{
									{
										Name:            "Capacity (first 50TB)",
										MonthlyQuantity: "300",
										Unit:            "GB",
										MonthlyCost:     "6.24",
									},
								},
							},
						}),
					},
				},
			},
		},
	})
}

func TestAccEstimateResource_MarkdownOutput(t *testing.T) {
	wd, _ := os.Getwd()
	testcase.Test(t, testcase.TestCase{
		SkipInit: true,
		Steps: []testcase.TestStep{
			{
				ConfigDirectory: path.Join(wd, "testdata", "markdown_output"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("plancost_estimate.this", plancheck.ResourceActionCreate),
					},
				},
				Check: func(t *testing.T, workDir string) {
					markdownPath := path.Join(workDir, "estimate.md")
					if _, err := os.Stat(markdownPath); os.IsNotExist(err) {
						t.Errorf("Markdown file was not created at %s", markdownPath)
						return
					}
					// Optionally check content
					content, err := os.ReadFile(markdownPath)
					if err != nil {
						t.Errorf("Failed to read markdown file: %v", err)
						return
					}
					if len(content) == 0 {
						t.Errorf("Markdown file is empty")
					}
					// Check for expected content
					if !strings.Contains(string(content), "### PlanCost Report") {
						t.Errorf("Markdown file does not contain expected title")
					}
				},
			},
		},
	})
}
