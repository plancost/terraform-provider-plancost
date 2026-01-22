/*
Copyright (c) 2026 Plancost.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package provider

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/plancost/terraform-provider-plancost/internal/optimization"
	tfschema "github.com/plancost/terraform-provider-plancost/internal/schema"
)

func GenerateConsoleOutput(displayName string, resources []*tfschema.Resource, recommendations []optimization.OptimizationRecommendation, paidTier bool) string {
	var sb strings.Builder

	if displayName == "" {
		displayName = "main"
	}

	// Header
	sb.WriteString(fmt.Sprintf("Project: %s\n\n", displayName))
	sb.WriteString(fmt.Sprintf(" %-58s %12s  %-14s %12s\n", "Name", "Monthly Qty", "Unit", "Monthly Cost"))
	sb.WriteString("\n")

	var baselineCost float64
	var usageCost float64

	estimatedResourceCount := 0
	freeResourceCount := 0
	unsupportedResourceCount := 0
	unsupportedResourceTypes := make(map[string]bool)

	for _, res := range resources {
		if res.IsSkipped {
			if res.NoPrice {
				freeResourceCount++
			} else {
				unsupportedResourceCount++
				unsupportedResourceTypes[res.ResourceType] = true
			}
			continue
		}
		estimatedResourceCount++
		resName := res.Name

		resBaseline, resUsage := calculateResourceCosts(res)
		baselineCost += resBaseline
		usageCost += resUsage

		// Print Resource Name
		sb.WriteString(fmt.Sprintf(" %s\n", resName))

		// Print cost components and sub-resources with proper tree structure
		printResourceTree(&sb, res.CostComponents, res.SubResources, " ")
		sb.WriteString("\n")
	}

	totalCost := baselineCost + usageCost

	sb.WriteString(fmt.Sprintf(" OVERALL TOTAL%87s\n", formatAmount(totalCost)))
	sb.WriteString("\n")
	sb.WriteString("*Usage costs can be estimated by providing usage data in the plancost_estimate resource.\n")
	sb.WriteString("\n")
	sb.WriteString("──────────────────────────────────\n")

	totalDetected := estimatedResourceCount + freeResourceCount + unsupportedResourceCount
	sb.WriteString(fmt.Sprintf("%d cloud resources were detected:\n", totalDetected))
	sb.WriteString(fmt.Sprintf("∙ %d were estimated\n", estimatedResourceCount))
	if freeResourceCount > 0 {
		sb.WriteString(fmt.Sprintf("∙ %d were free\n", freeResourceCount))
	}

	if unsupportedResourceCount > 0 {
		types := make([]string, 0, len(unsupportedResourceTypes))
		for t := range unsupportedResourceTypes {
			types = append(types, t)
		}
		sort.Strings(types)
		sb.WriteString(fmt.Sprintf("∙ %d are not supported yet: [%s]\n", unsupportedResourceCount, strings.Join(types, ", ")))
	}
	sb.WriteString("\n")

	// Summary Table
	sb.WriteString("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┓\n")
	sb.WriteString("┃ Project                                            ┃ Baseline cost ┃ Usage cost* ┃ Total cost ┃\n")
	sb.WriteString("┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━┫\n")
	sb.WriteString(fmt.Sprintf("┃ %-50s ┃ %13s ┃ %11s ┃ %10s ┃\n",
		truncateString(displayName, 50),
		formatAmount(baselineCost),
		formatAmount(usageCost),
		formatAmount(totalCost)))
	sb.WriteString("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┛\n")

	// Optimization Opportunities
	// Filter recommendations, only keep the Reservation or Teaser
	var filteredRecs []optimization.OptimizationRecommendation
	for _, rec := range recommendations {
		if rec.Type == "Reservation" || rec.Type == "Teaser" {
			filteredRecs = append(filteredRecs, rec)
		}
	}
	recommendations = filteredRecs

	if len(recommendations) > 0 {
		sb.WriteString("\n")
		sb.WriteString("💡 Optimization Recommendations\n")
		sb.WriteString("\n")

		if paidTier {
			grouped, keys := optimization.GroupOptimizations(recommendations)
			for _, resAddr := range keys {
				opts := grouped[resAddr]
				sb.WriteString(fmt.Sprintf(" %s\n", resAddr))
				for i, opt := range opts {
					isLast := i == len(opts)-1
					prefix := " ├─"
					if isLast {
						prefix = " └─"
					}
					// Format: "{Term} {Type}: Save ${RoundedAmount}/mo ({Pct}%)"
					// Term might be "1 yr" -> "1-Year"
					term := strings.Replace(opt.Term, " yr", "-Year", 1)
					term = strings.Replace(term, " Yr", "-Year", 1) // just in case

					// Round amount
					amount := math.Round(opt.SavingsAmount)

					// Pct is 0.0 - 1.0, display as 58%
					pct := math.Round(opt.SavingsPercentage * 100)

					sb.WriteString(fmt.Sprintf("%s %s Reservation: Save $%s/mo (%.0f%%)\n", prefix, term, addCommas(fmt.Sprintf("%.0f", amount)), pct))
				}
				sb.WriteString("\n")
			}
		} else {
			if len(recommendations) > 0 {
				sb.WriteString(" ")
				sb.WriteString(recommendations[0].Description)
			}
		}
	}

	return sb.String()
}

func truncateString(str string, num int) string {
	if utf8.RuneCountInString(str) > num {
		runes := []rune(str)
		if len(runes) > num {
			return string(runes[0:num-3]) + "..."
		}
	}
	return str
}

func calculateResourceCosts(res *tfschema.Resource) (baseline float64, usage float64) {
	for _, cc := range res.CostComponents {
		if cc.MonthlyCost != nil {
			val := cc.MonthlyCost.InexactFloat64()
			if cc.UsageBased {
				usage += val
			} else {
				baseline += val
			}
		}
	}
	for _, sub := range res.SubResources {
		b, u := calculateResourceCosts(sub)
		baseline += b
		usage += u
	}
	return
}

func printResourceTree(sb *strings.Builder, costComponents []*tfschema.CostComponent, subResources []*tfschema.Resource, baseIndent string) {
	// Collect all items (cost components and sub-resources)
	type item struct {
		name          string
		isSubResource bool
		cc            *tfschema.CostComponent
		subRes        *tfschema.Resource
	}

	var items []item
	for _, cc := range costComponents {
		items = append(items, item{name: cc.Name, isSubResource: false, cc: cc})
	}
	for _, subRes := range subResources {
		items = append(items, item{name: subRes.Name, isSubResource: true, subRes: subRes})
	}

	// Sort items by name
	sort.Slice(items, func(i, j int) bool {
		return items[i].name < items[j].name
	})

	// Print each item
	for i, item := range items {
		isLast := i == len(items)-1
		prefix := baseIndent + "├─"
		if isLast {
			prefix = baseIndent + "└─"
		}

		if item.isSubResource {
			// Print sub-resource name
			fmt.Fprintf(sb, "%s %s\n", prefix, item.name)

			// Determine indent for children
			childIndent := baseIndent + "│  "
			if isLast {
				childIndent = baseIndent + "   "
			}

			// Recursively print sub-resource contents
			printResourceTree(sb, item.subRes.CostComponents, item.subRes.SubResources, childIndent)
		} else {
			// Calculate available width for name based on prefix length
			nameWidth := 58 - utf8.RuneCountInString(prefix)
			if nameWidth < 10 {
				nameWidth = 10
			}

			if item.cc.UsageBased && (item.cc.MonthlyQuantity == nil || item.cc.MonthlyQuantity.IsZero()) {
				unitPrice := item.cc.UnitMultiplierPrice()
				var priceStr string
				if unitPrice.IsZero() {
					priceStr = "0.00"
				} else if unitPrice.Equal(unitPrice.Round(2)) {
					priceStr = unitPrice.StringFixed(2)
				} else {
					priceStr = unitPrice.String()
				}

				truncatedName := truncateString(item.name, nameWidth+1)
				space := " "
				if utf8.RuneCountInString(truncatedName) >= nameWidth {
					space = ""
				}

				fmt.Fprintf(sb, "%s %-*s%s Monthly cost depends on usage: $%s per %s\n",
					prefix,
					nameWidth,
					truncatedName,
					space,
					priceStr,
					item.cc.Unit)
			} else {
				// Print cost component
				cost := 0.0
				if item.cc.MonthlyCost != nil {
					cost = item.cc.MonthlyCost.InexactFloat64()
				}

				monthlyQuantity := "0"
				if item.cc.MonthlyQuantity != nil {
					monthlyQuantity = formatQuantity(item.cc.MonthlyQuantity.InexactFloat64())
				}

				usageMarker := ""
				if item.cc.UsageBased {
					usageMarker = "  *"
				}

				truncatedName := truncateString(item.name, nameWidth+1)
				space := " "
				if utf8.RuneCountInString(truncatedName) >= nameWidth {
					space = ""
				}

				// Total layout: prefix + name (variable) + quantity (12) + unit (14) + cost (12)
				fmt.Fprintf(sb, "%s %-*s%s%12s  %-14s %12s%s\n",
					prefix,
					nameWidth,
					truncatedName,
					space,
					monthlyQuantity,
					item.cc.Unit,
					formatAmount(cost),
					usageMarker)
			}
		}
	}
}

func formatAmount(amount float64) string {
	s := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(s, ".")
	return "$" + addCommas(parts[0]) + "." + parts[1]
}

func formatQuantity(q float64) string {
	s := fmt.Sprintf("%.4f", q)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")

	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = "." + parts[1]
	}
	return addCommas(intPart) + decPart
}

func addCommas(s string) string {
	var result strings.Builder
	n := len(s)
	remainder := n % 3
	if remainder > 0 {
		result.WriteString(s[:remainder])
		if n > 3 {
			result.WriteString(",")
		}
	}
	for i := remainder; i < n; i += 3 {
		if i > remainder {
			result.WriteString(",")
		}
		result.WriteString(s[i : i+3])
	}
	return result.String()
}
