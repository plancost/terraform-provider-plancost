/*
Copyright (c) 2026 Plancost.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package provider

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerateMarkdownDiff(t *testing.T) {
	// Define prior resources
	priorResources := []CostResourceModel{
		{
			Name: "resource_modified",
			CostComponents: []CostComponentModel{
				{
					Name:            "Compute",
					MonthlyQuantity: "730",
					Unit:            "hours",
					MonthlyCost:     10.00,
				},
			},
		},
		{
			Name: "resource_removed",
			CostComponents: []CostComponentModel{
				{
					Name:            "Storage",
					MonthlyQuantity: "100",
					Unit:            "GB",
					MonthlyCost:     5.00,
				},
			},
		},
	}

	// Define new resources
	newResources := []CostResourceModel{
		{
			Name: "resource_modified",
			CostComponents: []CostComponentModel{
				{
					Name:            "Compute",
					MonthlyQuantity: "730",
					Unit:            "hours",
					MonthlyCost:     20.00, // Cost increased
				},
			},
		},
		{
			Name: "resource_added",
			CostComponents: []CostComponentModel{
				{
					Name:            "Database",
					MonthlyQuantity: "1",
					Unit:            "instance",
					MonthlyCost:     50.00,
				},
			},
		},
	}

	// Call the function
	markdown := GenerateMarkdownOutput(priorResources, newResources)

	// Verify output
	expectedStrings := []string{
		"### PlanCost Report",
		"💰 Monthly cost will increase by $55.00 (367%).",
		"| Name | Monthly Qty | Unit | Monthly Cost |",
		"|:--- |:--- |:--- |:--- |",
		// Added resource
		"| + **resource_added** | | | |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;+ Database | 1 | instance | $50.00 |",
		// Modified resource
		"| ~ **resource_modified** | | | |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;~ Compute | 730 | hours | $10.00 -> $20.00 (+$10.00) |",
		// Removed resource
		"| - ~~resource_removed~~ | | | |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;- ~~Storage~~ | | | ~~$5.00~~ |",
		// Total
		"| **Total** | | | **$15.00 -> $70.00 (+$55.00, 367%)** |",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(markdown, s) {
			t.Errorf("Markdown output missing expected string: %q\nGot:\n%s", s, markdown)
		}
		fmt.Println(s)
	}
}

func TestGenerateMarkdownDiff_Complex(t *testing.T) {
	// Define prior resources
	priorResources := []CostResourceModel{
		// Modified Resource (Prior State)
		{
			Name: "azurerm_kubernetes_cluster.aks_mod",
			CostComponents: []CostComponentModel{
				{
					Name:            "Uptime SLA",
					MonthlyQuantity: "730",
					Unit:            "hours",
					MonthlyCost:     73.00,
				},
			},
			SubResources: []CostResourceModel{
				{
					Name: "default_node_pool",
					CostComponents: []CostComponentModel{
						{
							Name:            "Instance usage",
							MonthlyQuantity: "730",
							Unit:            "hours",
							MonthlyCost:     100.00,
						},
					},
					SubResources: []CostResourceModel{
						{
							Name: "os_disk",
							CostComponents: []CostComponentModel{
								{
									Name:            "Storage",
									MonthlyQuantity: "128",
									Unit:            "GB",
									MonthlyCost:     10.00,
								},
							},
						},
					},
				},
			},
		},
		// Removed Resource
		{
			Name: "azurerm_kubernetes_cluster.aks_old",
			CostComponents: []CostComponentModel{
				{
					Name:            "Uptime SLA",
					MonthlyQuantity: "730",
					Unit:            "hours",
					MonthlyCost:     73.00,
				},
			},
			SubResources: []CostResourceModel{
				{
					Name: "default_node_pool",
					CostComponents: []CostComponentModel{
						{
							Name:            "Instance usage",
							MonthlyQuantity: "730",
							Unit:            "hours",
							MonthlyCost:     50.00,
						},
					},
				},
			},
		},
	}

	// Define new resources
	newResources := []CostResourceModel{
		// Modified Resource (New State)
		{
			Name: "azurerm_kubernetes_cluster.aks_mod",
			CostComponents: []CostComponentModel{
				{
					Name:            "Uptime SLA",
					MonthlyQuantity: "730",
					Unit:            "hours",
					MonthlyCost:     73.00, // Unchanged
				},
			},
			SubResources: []CostResourceModel{
				{
					Name: "default_node_pool",
					CostComponents: []CostComponentModel{
						{
							Name:            "Instance usage",
							MonthlyQuantity: "1460", // Doubled
							Unit:            "hours",
							MonthlyCost:     200.00, // Doubled cost
						},
						{
							Name:            "Load Balancer", // Added component
							MonthlyQuantity: "1",
							Unit:            "months",
							MonthlyCost:     5.00,
						},
					},
					SubResources: []CostResourceModel{
						{
							Name: "os_disk",
							CostComponents: []CostComponentModel{
								{
									Name:            "Storage",
									MonthlyQuantity: "128",
									Unit:            "GB",
									MonthlyCost:     10.00, // Unchanged
								},
							},
						},
					},
				},
			},
		},
		// Added Resource
		{
			Name: "azurerm_kubernetes_cluster.aks_new",
			CostComponents: []CostComponentModel{
				{
					Name:            "Uptime SLA",
					MonthlyQuantity: "730",
					Unit:            "hours",
					MonthlyCost:     73.00,
				},
			},
			SubResources: []CostResourceModel{
				{
					Name: "default_node_pool",
					CostComponents: []CostComponentModel{
						{
							Name:            "Instance usage",
							MonthlyQuantity: "730",
							Unit:            "hours",
							MonthlyCost:     150.00,
						},
					},
				},
			},
		},
	}

	// Call the function
	markdown := GenerateMarkdownOutput(priorResources, newResources)

	// Verify output
	expectedStrings := []string{
		"### PlanCost Report",
		"💰 Monthly cost will increase by $205.00 (67%).",
		"| Name | Monthly Qty | Unit | Monthly Cost |",
		"|:--- |:--- |:--- |:--- |",

		// Added Resource: azurerm_kubernetes_cluster.aks_new
		"| + **azurerm_kubernetes_cluster.aks_new** | | | |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;+ Uptime SLA | 730 | hours | $73.00 |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;+ default_node_pool | | | |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+ Instance usage | 730 | hours | $150.00 |",

		// Modified Resource: azurerm_kubernetes_cluster.aks_mod
		"| ~ **azurerm_kubernetes_cluster.aks_mod** | | | |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;  Uptime SLA | 730 | hours | $73.00 |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;~ default_node_pool | | | |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;~ Instance usage | 730 -> 1460 | hours | $100.00 -> $200.00 (+$100.00) |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+ Load Balancer | 1 | months | $5.00 |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  os_disk | | | |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  Storage | 128 | GB | $10.00 |",

		// Removed Resource: azurerm_kubernetes_cluster.aks_old
		"| - ~~azurerm_kubernetes_cluster.aks_old~~ | | | |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;- ~~Uptime SLA~~ | | | ~~$73.00~~ |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;- ~~default_node_pool~~ | | | |",
		"| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- ~~Instance usage~~ | | | ~~$50.00~~ |",

		// Total
		// Prior: 183 (mod) + 123 (old) = 306
		// New: 288 (mod) + 223 (new) = 511
		// Diff: +205
		"| **Total** | | | **$306.00 -> $511.00 (+$205.00, 67%)** |",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(markdown, s) {
			t.Errorf("Markdown output missing expected string: %q\nGot:\n%s", s, markdown)
		}
		fmt.Println(s)
	}
}
