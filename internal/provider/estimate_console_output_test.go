/*
Copyright (c) 2026 Plancost.
Licensed under the Business Source License 1.1.
See LICENSE file in the project root for full license information.
*/

package provider

import (
	"testing"

	"github.com/plancost/terraform-provider-plancost/internal/optimization"
	tfschema "github.com/plancost/terraform-provider-plancost/internal/schema"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGenerateConsoleOutput(t *testing.T) {
	// Helper to create decimal pointer
	d := func(f float64) *decimal.Decimal {
		v := decimal.NewFromFloat(f)
		return &v
	}

	// Helper to create cost component with price
	ccWithPrice := func(name, unit string, price float64) *tfschema.CostComponent {
		cc := &tfschema.CostComponent{
			Name:           name,
			Unit:           unit,
			UsageBased:     true,
			UnitMultiplier: decimal.NewFromInt(1),
		}
		cc.SetPrice(decimal.NewFromFloat(price))
		return cc
	}

	tests := []struct {
		name            string
		displayName     string
		resources       []*tfschema.Resource
		recommendations []optimization.OptimizationRecommendation
		paidTier        bool
		expected        string
	}{
		{
			name:        "Resource without usage",
			displayName: "main",
			resources: []*tfschema.Resource{
				{
					Name: "azurerm_application_gateway.waf_without_usage",
					CostComponents: []*tfschema.CostComponent{
						{
							Name:            "Gateway usage (WAF, medium)",
							Unit:            "hours",
							MonthlyQuantity: d(1460.0),
							MonthlyCost:     d(183.96),
							UnitMultiplier:  decimal.NewFromInt(1),
						},
						{
							Name:           "Data processing (0-10TB)",
							Unit:           "GB",
							UsageBased:     true,
							UnitMultiplier: decimal.NewFromInt(1),
						},
					},
				},
			},
			expected: `Project: main

 Name                                                        Monthly Qty  Unit           Monthly Cost

 azurerm_application_gateway.waf_without_usage
 ├─ Data processing (0-10TB)                                 Monthly cost depends on usage: $0.00 per GB
 └─ Gateway usage (WAF, medium)                                    1,460  hours               $183.96

 OVERALL TOTAL                                                                                $183.96

*Usage costs can be estimated by providing usage data in the plancost_estimate resource.

──────────────────────────────────
1 cloud resources were detected:
∙ 1 were estimated

┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
┃ main                                               ┃       $183.96 ┃       $0.00 ┃    $183.96 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛
`,
		},
		{
			name:        "Simple resource",
			displayName: "main",
			resources: []*tfschema.Resource{
				{
					Name: "azurerm_resource_group.example",
					CostComponents: []*tfschema.CostComponent{
						{
							Name:            "Cost Component 1",
							Unit:            "GB",
							MonthlyQuantity: d(10.0),
							MonthlyCost:     d(5.0),
							UnitMultiplier:  decimal.NewFromInt(1),
						},
					},
				},
			},
			expected: `Project: main

 Name                                                        Monthly Qty  Unit           Monthly Cost

 azurerm_resource_group.example
 └─ Cost Component 1                                                  10  GB                    $5.00

 OVERALL TOTAL                                                                                  $5.00

*Usage costs can be estimated by providing usage data in the plancost_estimate resource.

──────────────────────────────────
1 cloud resources were detected:
∙ 1 were estimated

┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
┃ main                                               ┃         $5.00 ┃       $0.00 ┃      $5.00 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛
`,
		},
		{
			name:        "Usage based resource",
			displayName: "main",
			resources: []*tfschema.Resource{
				{
					Name: "azurerm_storage_account.example",
					CostComponents: []*tfschema.CostComponent{
						{
							Name:            "Storage",
							Unit:            "GB",
							MonthlyQuantity: d(100.0),
							MonthlyCost:     d(20.0),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
					},
				},
			},
			expected: `Project: main

 Name                                                        Monthly Qty  Unit           Monthly Cost

 azurerm_storage_account.example
 └─ Storage                                                          100  GB                   $20.00  *

 OVERALL TOTAL                                                                                 $20.00

*Usage costs can be estimated by providing usage data in the plancost_estimate resource.

──────────────────────────────────
1 cloud resources were detected:
∙ 1 were estimated

┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
┃ main                                               ┃         $0.00 ┃      $20.00 ┃     $20.00 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛
`,
		},
		{
			name:        "Nested resources",
			displayName: "main",
			resources: []*tfschema.Resource{
				{
					Name: "azurerm_kubernetes_cluster.example",
					SubResources: []*tfschema.Resource{
						{
							Name: "default_node_pool",
							CostComponents: []*tfschema.CostComponent{
								{
									Name:            "VM usage",
									Unit:            "hours",
									MonthlyQuantity: d(730.0),
									MonthlyCost:     d(100.0),
									UnitMultiplier:  decimal.NewFromInt(1),
								},
							},
						},
					},
				},
			},
			expected: `Project: main

 Name                                                        Monthly Qty  Unit           Monthly Cost

 azurerm_kubernetes_cluster.example
 └─ default_node_pool
    └─ VM usage                                                      730  hours               $100.00

 OVERALL TOTAL                                                                                $100.00

*Usage costs can be estimated by providing usage data in the plancost_estimate resource.

──────────────────────────────────
1 cloud resources were detected:
∙ 1 were estimated

┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
┃ main                                               ┃       $100.00 ┃       $0.00 ┃    $100.00 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛
`,
		},
		{
			name:        "Complex nested resources",
			displayName: "main",
			resources: []*tfschema.Resource{
				{
					Name: "azurerm_kubernetes_cluster.usage_ephemeral",
					CostComponents: []*tfschema.CostComponent{
						{
							Name:            "Uptime SLA",
							Unit:            "hours",
							MonthlyQuantity: d(730.0),
							MonthlyCost:     d(73.0),
							UnitMultiplier:  decimal.NewFromInt(1),
						},
					},
					SubResources: []*tfschema.Resource{
						{
							Name: "default_node_pool",
							CostComponents: []*tfschema.CostComponent{
								{
									Name:            "Instance usage (Linux, pay as you go, Standard_D2_v2)",
									Unit:            "hours",
									MonthlyQuantity: d(900.0),
									MonthlyCost:     d(131.40),
									UnitMultiplier:  decimal.NewFromInt(1),
								},
							},
						},
						{
							Name: "Load Balancer",
							CostComponents: []*tfschema.CostComponent{
								{
									Name:            "Data processed",
									Unit:            "GB",
									MonthlyQuantity: d(100.0),
									MonthlyCost:     d(0.50),
									UsageBased:      true,
									UnitMultiplier:  decimal.NewFromInt(1),
								},
							},
						},
						{
							Name: "DNS",
							CostComponents: []*tfschema.CostComponent{
								{
									Name:            "Hosted zone",
									Unit:            "months",
									MonthlyQuantity: d(1.0),
									MonthlyCost:     d(0.50),
									UnitMultiplier:  decimal.NewFromInt(1),
								},
							},
						},
					},
				},
			},
			expected: `Project: main

 Name                                                        Monthly Qty  Unit           Monthly Cost

 azurerm_kubernetes_cluster.usage_ephemeral
 ├─ DNS
 │  └─ Hosted zone                                                     1  months                $0.50
 ├─ Load Balancer
 │  └─ Data processed                                                100  GB                    $0.50  *
 ├─ Uptime SLA                                                       730  hours                $73.00
 └─ default_node_pool
    └─ Instance usage (Linux, pay as you go, Standard_D2_v2)         900  hours               $131.40

 OVERALL TOTAL                                                                                $205.40

*Usage costs can be estimated by providing usage data in the plancost_estimate resource.

──────────────────────────────────
1 cloud resources were detected:
∙ 1 were estimated

┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
┃ main                                               ┃       $204.90 ┃       $0.50 ┃    $205.40 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛
`,
		},
		{
			name:        "Very long resource name",
			displayName: "main",
			resources: []*tfschema.Resource{
				{
					Name: "azurerm_resource_group.this_is_a_very_long_resource_name_that_should_be_truncated_in_the_output",
					CostComponents: []*tfschema.CostComponent{
						{
							Name:            "Cost Component with a very long name that should also be truncated",
							Unit:            "GB",
							MonthlyQuantity: d(10.0),
							MonthlyCost:     d(5.0),
							UnitMultiplier:  decimal.NewFromInt(1),
						},
					},
				},
			},
			expected: `Project: main

 Name                                                        Monthly Qty  Unit           Monthly Cost

 azurerm_resource_group.this_is_a_very_long_resource_name_that_should_be_truncated_in_the_output
 └─ Cost Component with a very long name that should also...          10  GB                    $5.00

 OVERALL TOTAL                                                                                  $5.00

*Usage costs can be estimated by providing usage data in the plancost_estimate resource.

──────────────────────────────────
1 cloud resources were detected:
∙ 1 were estimated

┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
┃ main                                               ┃         $5.00 ┃       $0.00 ┃      $5.00 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛
`,
		},
		{
			name:        "Unsupported resources",
			displayName: "main",
			resources: []*tfschema.Resource{
				{
					Name: "azurerm_resource_group.example",
					CostComponents: []*tfschema.CostComponent{
						{
							Name:            "Cost Component 1",
							Unit:            "GB",
							MonthlyQuantity: d(10.0),
							MonthlyCost:     d(5.0),
							UnitMultiplier:  decimal.NewFromInt(1),
						},
					},
				},
				{
					Name:         "azurerm_unsupported_resource.example1",
					ResourceType: "azurerm_unsupported_resource",
					IsSkipped:    true,
					NoPrice:      false,
				},
				{
					Name:         "azurerm_another_unsupported.example2",
					ResourceType: "azurerm_another_unsupported",
					IsSkipped:    true,
					NoPrice:      false,
				},
			},
			expected: `Project: main

 Name                                                        Monthly Qty  Unit           Monthly Cost

 azurerm_resource_group.example
 └─ Cost Component 1                                                  10  GB                    $5.00

 OVERALL TOTAL                                                                                  $5.00

*Usage costs can be estimated by providing usage data in the plancost_estimate resource.

──────────────────────────────────
3 cloud resources were detected:
∙ 1 were estimated
∙ 2 are not supported yet: [azurerm_another_unsupported, azurerm_unsupported_resource]

┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
┃ main                                               ┃         $5.00 ┃       $0.00 ┃      $5.00 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛
`,
		},
		{
			name:        "Complex usage scenario",
			displayName: "main",
			resources: []*tfschema.Resource{
				{
					Name: "module.main.azurerm_postgresql_flexible_server.db",
					CostComponents: []*tfschema.CostComponent{
						{
							Name:            "Additional backup storage",
							Unit:            "GB",
							MonthlyQuantity: d(5000),
							MonthlyCost:     d(540.00),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
						{
							Name:            "Compute (B_Standard_B1ms)",
							Unit:            "hours",
							MonthlyQuantity: d(730),
							MonthlyCost:     d(14.02),
							UnitMultiplier:  decimal.NewFromInt(1),
						},
						{
							Name:            "Storage",
							Unit:            "GB",
							MonthlyQuantity: d(32),
							MonthlyCost:     d(4.16),
							UnitMultiplier:  decimal.NewFromInt(1),
						},
					},
				},
				{
					Name: "module.main.azurerm_storage_account.main",
					CostComponents: []*tfschema.CostComponent{
						{
							Name:            "All other operations",
							Unit:            "10k operations",
							MonthlyQuantity: d(100),
							MonthlyCost:     d(0.40),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
						{
							Name:            "Blob index",
							Unit:            "10k tags",
							MonthlyQuantity: d(10),
							MonthlyCost:     d(0.37),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
						{
							Name:            "Capacity (first 50TB)",
							Unit:            "GB",
							MonthlyQuantity: d(51200),
							MonthlyCost:     d(942.08),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
						{
							Name:            "Capacity (next 450TB)",
							Unit:            "GB",
							MonthlyQuantity: d(512000),
							MonthlyCost:     d(9043.97),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
						{
							Name:            "Capacity (over 500TB)",
							Unit:            "GB",
							MonthlyQuantity: d(436800),
							MonthlyCost:     d(7394.15),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
						{
							Name:            "List and create container operations",
							Unit:            "10k operations",
							MonthlyQuantity: d(100),
							MonthlyCost:     d(5.00),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
						{
							Name:            "Read operations",
							Unit:            "10k operations",
							MonthlyQuantity: d(10),
							MonthlyCost:     d(0.04),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
						{
							Name:            "Write operations",
							Unit:            "10k operations",
							MonthlyQuantity: d(100),
							MonthlyCost:     d(5.00),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
					},
				},
				{
					Name: "module.main.azurerm_application_insights.app_insights",
					CostComponents: []*tfschema.CostComponent{
						{
							Name:            "Data ingested",
							Unit:            "GB",
							MonthlyQuantity: d(1000),
							MonthlyCost:     d(2300.00),
							UsageBased:      true,
							UnitMultiplier:  decimal.NewFromInt(1),
						},
					},
				},
				{
					Name: "module.main.azurerm_linux_function_app.pricing_function",
					CostComponents: []*tfschema.CostComponent{
						ccWithPrice("Execution time", "GB-seconds", 0.000016),
						ccWithPrice("Executions", "1M requests", 0.20),
					},
				},
			},
			expected: `Project: main

 Name                                                        Monthly Qty  Unit           Monthly Cost

 module.main.azurerm_postgresql_flexible_server.db
 ├─ Additional backup storage                                      5,000  GB                  $540.00  *
 ├─ Compute (B_Standard_B1ms)                                        730  hours                $14.02
 └─ Storage                                                           32  GB                    $4.16

 module.main.azurerm_storage_account.main
 ├─ All other operations                                             100  10k operations        $0.40  *
 ├─ Blob index                                                        10  10k tags              $0.37  *
 ├─ Capacity (first 50TB)                                         51,200  GB                  $942.08  *
 ├─ Capacity (next 450TB)                                        512,000  GB                $9,043.97  *
 ├─ Capacity (over 500TB)                                        436,800  GB                $7,394.15  *
 ├─ List and create container operations                             100  10k operations        $5.00  *
 ├─ Read operations                                                   10  10k operations        $0.04  *
 └─ Write operations                                                 100  10k operations        $5.00  *

 module.main.azurerm_application_insights.app_insights
 └─ Data ingested                                                  1,000  GB                $2,300.00  *

 module.main.azurerm_linux_function_app.pricing_function
 ├─ Execution time                                           Monthly cost depends on usage: $0.000016 per GB-seconds
 └─ Executions                                               Monthly cost depends on usage: $0.20 per 1M requests

 OVERALL TOTAL                                                                             $20,249.19

*Usage costs can be estimated by providing usage data in the plancost_estimate resource.

──────────────────────────────────
4 cloud resources were detected:
∙ 4 were estimated

┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
┃ main                                               ┃        $18.18 ┃  $20,231.01 ┃ $20,249.19 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛
`,
		},
		{
			name:        "Optimization Recommendations (Paid)",
			displayName: "main",
			resources: []*tfschema.Resource{
				{
					Name: "azurerm_linux_virtual_machine.vm",
					CostComponents: []*tfschema.CostComponent{
						ccWithPrice("Compute", "hours", 0.5),
					},
				},
			},
			recommendations: []optimization.OptimizationRecommendation{
				{
					ResourceAddress:   "azurerm_linux_virtual_machine.vm",
					Description:       "Save $369/mo (58%) on azurerm_linux_virtual_machine.vm with a 1-year Reservation",
					Type:              "Reservation",
					Term:              "1 yr",
					SavingsAmount:     369.48,
					SavingsPercentage: 0.58,
				},
			},
			paidTier: true,
			expected: `Project: main

 Name                                                        Monthly Qty  Unit           Monthly Cost

 azurerm_linux_virtual_machine.vm
 └─ Compute                                                  Monthly cost depends on usage: $0.50 per hours

 OVERALL TOTAL                                                                                  $0.00

*Usage costs can be estimated by providing usage data in the plancost_estimate resource.

──────────────────────────────────
1 cloud resources were detected:
∙ 1 were estimated

┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
┃ main                                               ┃         $0.00 ┃       $0.00 ┃      $0.00 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛

💡 Optimization Recommendations

 azurerm_linux_virtual_machine.vm
 └─ 1-Year Reservation: Save $369/mo (58%)

`},
		{
			name:        "Optimization Recommendations (Free)",
			displayName: "main",
			resources: []*tfschema.Resource{
				{
					Name: "azurerm_linux_virtual_machine.vm",
					CostComponents: []*tfschema.CostComponent{
						ccWithPrice("Compute", "hours", 0.5),
					},
				},
			},
			recommendations: []optimization.OptimizationRecommendation{
				{
					Description: "We found reservations that could save you ~$369/year (0%).\n [ Upgrade to Pro to see details ]\n",
					Type:        "Teaser",
				},
			},
			paidTier: false,
			expected: `Project: main

 Name                                                        Monthly Qty  Unit           Monthly Cost

 azurerm_linux_virtual_machine.vm
 └─ Compute                                                  Monthly cost depends on usage: $0.50 per hours

 OVERALL TOTAL                                                                                  $0.00

*Usage costs can be estimated by providing usage data in the plancost_estimate resource.

──────────────────────────────────
1 cloud resources were detected:
∙ 1 were estimated

┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓
┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫
┃ main                                               ┃         $0.00 ┃       $0.00 ┃      $0.00 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛

💡 Optimization Recommendations

 We found reservations that could save you ~$369/year (0%).
 [ Upgrade to Pro to see details ]
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := GenerateConsoleOutput(tt.displayName, tt.resources, tt.recommendations, tt.paidTier)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0.0, "$0.00"},
		{1.0, "$1.00"},
		{1234.56, "$1,234.56"},
		{1234567.89, "$1,234,567.89"},
		{0.004, "$0.00"}, // Rounds down
		{0.006, "$0.01"}, // Rounds up
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, formatAmount(tt.input))
	}
}

func TestFormatQuantity(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0.0, "0"},
		{1.0, "1"},
		{1000.0, "1,000"},
		{1234.5678, "1,234.5678"},
		{0.12345, "0.1235"}, // Rounds to 4 decimal places
		{10.5000, "10.5"},   // Trims trailing zeros
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, formatQuantity(tt.input))
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		num      int
		expected string
	}{
		{"short", 10, "short"},
		{"exactlength", 11, "exactlength"},
		{"verylongstring", 10, "verylon..."},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, truncateString(tt.input, tt.num))
	}
}
