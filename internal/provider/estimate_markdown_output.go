/*
Copyright (c) 2026 Plancost.
Licensed under the Business Source License 1.1.
See LICENSE file in the project root for full license information.
*/

package provider

import (
	"fmt"
	"sort"
	"strings"
)

func GenerateMarkdownOutput(priorResources, newResources []CostResourceModel) string {
	priorMap := make(map[string]CostResourceModel)
	totalPriorCost := 0.0
	for _, r := range priorResources {
		priorMap[r.Name] = r
		totalPriorCost += calculateResourceCost(r)
	}

	newMap := make(map[string]CostResourceModel)
	totalNewCost := 0.0
	for _, r := range newResources {
		newMap[r.Name] = r
		totalNewCost += calculateResourceCost(r)
	}

	// Collect all resource names
	allNames := make(map[string]bool)
	for k := range priorMap {
		allNames[k] = true
	}
	for k := range newMap {
		allNames[k] = true
	}

	sortedNames := make([]string, 0, len(allNames))
	for k := range allNames {
		sortedNames = append(sortedNames, k)
	}
	sort.Strings(sortedNames)

	var sb strings.Builder
	sb.WriteString("### PlanCost Report\n\n")

	diff := totalNewCost - totalPriorCost
	pct := 0.0
	if totalPriorCost > 0 {
		pct = (diff / totalPriorCost) * 100
	} else if totalNewCost > 0 {
		pct = 100
	}

	if diff > 0 {
		sb.WriteString(fmt.Sprintf("💰 Monthly cost will increase by $%.2f (%.0f%%).\n\n", diff, pct))
	} else if diff < 0 {
		sb.WriteString(fmt.Sprintf("💰 Monthly cost will decrease by $%.2f (%.0f%%).\n\n", -diff, -pct))
	} else {
		sb.WriteString("💰 Monthly cost will remain unchanged.\n\n")
	}

	sb.WriteString("| Name | Monthly Qty | Unit | Monthly Cost |\n")
	sb.WriteString("|:--- |:--- |:--- |:--- |\n")

	for _, name := range sortedNames {
		prior, hasPrior := priorMap[name]
		newRes, hasNew := newMap[name]

		if !hasPrior && hasNew {
			// Added resource
			sb.WriteString(fmt.Sprintf("| + **%s** | | | |\n", name))
			printMarkdownDiffTree(&sb, nil, nil, newRes.CostComponents, newRes.SubResources, 0)
		} else if hasPrior && !hasNew {
			// Deleted resource
			sb.WriteString(fmt.Sprintf("| - ~~%s~~ | | | |\n", name))
			printMarkdownDiffTree(&sb, prior.CostComponents, prior.SubResources, nil, nil, 0)
		} else {
			// Modified or Unchanged
			priorCost := calculateResourceCost(prior)
			newCost := calculateResourceCost(newRes)
			icon := "  "
			if priorCost != newCost {
				icon = "~ "
			}
			sb.WriteString(fmt.Sprintf("| %s**%s** | | | |\n", icon, name))
			printMarkdownDiffTree(&sb, prior.CostComponents, prior.SubResources, newRes.CostComponents, newRes.SubResources, 0)
		}
	}

	diffStr := ""
	if diff > 0 {
		diffStr = fmt.Sprintf(" (+$%.2f, %.0f%%)", diff, pct)
	} else if diff < 0 {
		diffStr = fmt.Sprintf(" (-$%.2f, %.0f%%)", -diff, -pct)
	}

	sb.WriteString(fmt.Sprintf("| **Total** | | | **$%.2f -> $%.2f%s** |\n", totalPriorCost, totalNewCost, diffStr))

	return sb.String()
}

func printMarkdownDiffTree(sb *strings.Builder, priorComponents []CostComponentModel, priorSubResources []CostResourceModel, newComponents []CostComponentModel, newSubResources []CostResourceModel, level int) {
	prefix := strings.Repeat("&nbsp;&nbsp;&nbsp;&nbsp;", level+1)

	// Map components by name
	priorCCMap := make(map[string]CostComponentModel)
	for _, c := range priorComponents {
		priorCCMap[c.Name] = c
	}
	newCCMap := make(map[string]CostComponentModel)
	for _, c := range newComponents {
		newCCMap[c.Name] = c
	}

	// Collect all component names
	allCCNames := make(map[string]bool)
	for k := range priorCCMap {
		allCCNames[k] = true
	}
	for k := range newCCMap {
		allCCNames[k] = true
	}
	sortedCCNames := make([]string, 0, len(allCCNames))
	for k := range allCCNames {
		sortedCCNames = append(sortedCCNames, k)
	}
	sort.Strings(sortedCCNames)

	for _, name := range sortedCCNames {
		prior, hasPrior := priorCCMap[name]
		newCC, hasNew := newCCMap[name]

		if !hasPrior && hasNew {
			// Added component
			fmt.Fprintf(sb, "| %s+ %s | %s | %s | $%.2f |\n", prefix, name, newCC.MonthlyQuantity, newCC.Unit, newCC.MonthlyCost)
		} else if hasPrior && !hasNew {
			// Deleted component
			fmt.Fprintf(sb, "| %s- ~~%s~~ | | | ~~$%.2f~~ |\n", prefix, name, prior.MonthlyCost)
		} else {
			// Modified or Unchanged
			qtyDiff := ""
			if prior.MonthlyQuantity != newCC.MonthlyQuantity {
				qtyDiff = fmt.Sprintf("%s -> %s", prior.MonthlyQuantity, newCC.MonthlyQuantity)
			} else {
				qtyDiff = newCC.MonthlyQuantity
			}

			costDiff := ""
			changed := false
			if prior.MonthlyCost != newCC.MonthlyCost {
				changed = true
				diff := newCC.MonthlyCost - prior.MonthlyCost
				sign := "+"
				if diff < 0 {
					sign = "-"
					diff = -diff
				}
				costDiff = fmt.Sprintf("$%.2f -> $%.2f (%s$%.2f)", prior.MonthlyCost, newCC.MonthlyCost, sign, diff)
			} else {
				costDiff = fmt.Sprintf("$%.2f", newCC.MonthlyCost)
			}

			if changed {
				fmt.Fprintf(sb, "| %s~ %s | %s | %s | %s |\n", prefix, name, qtyDiff, newCC.Unit, costDiff)
			} else {
				fmt.Fprintf(sb, "| %s  %s | %s | %s | %s |\n", prefix, name, qtyDiff, newCC.Unit, costDiff)
			}
		}
	}

	// Handle SubResources similarly
	priorSubMap := make(map[string]CostResourceModel)
	for _, r := range priorSubResources {
		priorSubMap[r.Name] = r
	}
	newSubMap := make(map[string]CostResourceModel)
	for _, r := range newSubResources {
		newSubMap[r.Name] = r
	}

	allSubNames := make(map[string]bool)
	for k := range priorSubMap {
		allSubNames[k] = true
	}
	for k := range newSubMap {
		allSubNames[k] = true
	}
	sortedSubNames := make([]string, 0, len(allSubNames))
	for k := range allSubNames {
		sortedSubNames = append(sortedSubNames, k)
	}
	sort.Strings(sortedSubNames)

	for _, name := range sortedSubNames {
		prior, hasPrior := priorSubMap[name]
		newSub, hasNew := newSubMap[name]

		if !hasPrior && hasNew {
			fmt.Fprintf(sb, "| %s+ %s | | | |\n", prefix, name)
			printMarkdownDiffTree(sb, nil, nil, newSub.CostComponents, newSub.SubResources, level+1)
		} else if hasPrior && !hasNew {
			fmt.Fprintf(sb, "| %s- ~~%s~~ | | | |\n", prefix, name)
			printMarkdownDiffTree(sb, prior.CostComponents, prior.SubResources, nil, nil, level+1)
		} else {
			priorCost := calculateResourceCost(prior)
			newCost := calculateResourceCost(newSub)
			icon := "  "
			if priorCost != newCost {
				icon = "~ "
			}
			fmt.Fprintf(sb, "| %s%s%s | | | |\n", prefix, icon, name)
			printMarkdownDiffTree(sb, prior.CostComponents, prior.SubResources, newSub.CostComponents, newSub.SubResources, level+1)
		}
	}
}

func calculateResourceCost(res CostResourceModel) float64 {
	cost := 0.0
	for _, c := range res.CostComponents {
		cost += c.MonthlyCost
	}
	for _, sub := range res.SubResources {
		cost += calculateResourceCost(sub)
	}
	return cost
}
