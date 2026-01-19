/*
Copyright (c) 2026 Plancost.
Licensed under the Business Source License 1.1.
See LICENSE file in the project root for full license information.
*/

package provider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/plancost/terraform-provider-plancost/internal/optimization"
	"github.com/plancost/terraform-provider-plancost/internal/prices"
	tfschema "github.com/plancost/terraform-provider-plancost/internal/schema"
	"github.com/shopspring/decimal"
)

// TaggingPolicies enforces tagging policies on the provided resources. Only available in the paid tier.
func TaggingPolicies(paidTier bool, TaggingPolicy []TaggingPolicyModel, allParsedResources []*tfschema.Resource) diag.Diagnostics {
	var resp diag.Diagnostics
	if !paidTier {
		if len(TaggingPolicy) > 0 {
			resp.AddWarning("Tagging Policies Disabled", "Tagging policies are a paid feature. Please upgrade to the paid tier at https://plancost.io to enable this feature.")
		}
		return resp
	}
	// Tagging Policy Logic
	for _, policy := range TaggingPolicy {
		key := policy.Key.ValueString()
		allowedValues := []string{}
		for _, v := range policy.AllowedValues {
			allowedValues = append(allowedValues, v.ValueString())
		}
		pattern := policy.Pattern.ValueString()
		resourceTypes := []string{}
		for _, v := range policy.ResourceTypes {
			resourceTypes = append(resourceTypes, v.ValueString())
		}
		action := policy.Action.ValueString()

		for _, res := range allParsedResources {
			// Check if resource supports tags (if Tags is nil, assume it doesn't or no tags are set)
			// If resourceTypes is NOT specified, we only check resources that have Tags != nil (meaning they support tags)
			if res.Tags == nil {
				continue
			}

			match := len(resourceTypes) == 0
			for _, rt := range resourceTypes {
				if res.ResourceType == rt {
					match = true
					break
				}
			}
			if !match {
				continue
			}

			tags := res.Tags
			var val string
			var ok bool

			if tags != nil {
				val, ok = (*tags)[key]
			}

			if !ok {
				// Violation: Tag key missing
				msg := fmt.Sprintf("Resource %s (%s) missing required tag '%s'", res.Name, res.ResourceType, key)
				if action == "block" {
					resp.AddError("Tagging Policy Violation", msg)
				} else {
					resp.AddWarning("Tagging Policy Violation", msg)
				}
				continue
			}

			// Check allowed values
			if len(allowedValues) > 0 {
				valid := false
				for _, av := range allowedValues {
					if strings.EqualFold(av, val) {
						valid = true
						break
					}
				}
				if !valid {
					msg := fmt.Sprintf("Resource %s (%s) tag '%s' has invalid value '%s'. Allowed: %v", res.Name, res.ResourceType, key, val, allowedValues)
					if action == "block" {
						resp.AddError("Tagging Policy Violation", msg)
					} else {
						resp.AddWarning("Tagging Policy Violation", msg)
					}
				}
			}

			// Check pattern
			if pattern != "" {
				matched, err := regexp.MatchString(pattern, val)
				if err != nil {
					resp.AddError("Invalid Regex Pattern", fmt.Sprintf("Pattern '%s' is invalid: %s", pattern, err))
					continue
				}
				if !matched {
					msg := fmt.Sprintf("Resource %s (%s) tag '%s' value '%s' does not match pattern '%s'", res.Name, res.ResourceType, key, val, pattern)
					if action == "block" {
						resp.AddError("Tagging Policy Violation", msg)
					} else {
						resp.AddWarning("Tagging Policy Violation", msg)
					}
				}
			}
		}
	}

	return resp
}

// Guardrails enforces cost guardrails based on the provided configurations. Free tier only allows 1 guardrail and does not enforce "block" actions.
func Guardrails(paidTier bool, Guardrail []GuardrailModel, totalCost, previousCost float64) diag.Diagnostics {
	var resp diag.Diagnostics

	diffAmount := totalCost - previousCost
	diffPercent := 0.0
	if previousCost > 0 {
		diffPercent = (diffAmount / previousCost) * 100
	}

	if !paidTier && len(Guardrail) > 1 {
		// only allow 1 guardrail for free tier
		resp.AddWarning("Guardrails Limited", "Guardrails enforcement is a paid feature. Please upgrade to the paid tier at https://plancost.io to enable this feature. Only the first guardrail will be evaluated.")
	}

	for _, guardrail := range Guardrail {
		condition := guardrail.Condition.ValueString()
		threshold, _ := guardrail.Threshold.ValueBigFloat().Float64()
		action := guardrail.Action.ValueString()

		triggered := false
		var msg string

		switch condition {
		case "monthly_cost_increase_amount":
			if diffAmount > threshold {
				triggered = true
				msg = fmt.Sprintf("Monthly cost increase amount $%.2f exceeds threshold $%.2f.", diffAmount, threshold)
			}
		case "monthly_cost_increase_percentage":
			if previousCost > 0 && diffPercent > threshold {
				triggered = true
				msg = fmt.Sprintf("Monthly cost increase percentage %.2f%% exceeds threshold %.2f%%.", diffPercent, threshold)
			}
		case "monthly_cost_budget":
			if totalCost > threshold {
				triggered = true
				msg = fmt.Sprintf("Monthly cost $%.2f exceeds budget $%.2f.", totalCost, threshold)
			}
		}

		if triggered {
			if action == "block" {
				if !paidTier {
					msg += " Guardrails enforcement is a paid feature. Please upgrade to the paid tier at https://plancost.io to enable this feature."
					resp.AddWarning("Guardrail Violation", msg)
				} else {
					resp.AddError("Guardrail Violation", msg)
				}
			} else {
				resp.AddWarning("Guardrail Violation", msg)
			}
		}

		if !paidTier {
			// only evaluate the first guardrail for free tier
			break
		}
	}

	return resp
}

// Optimization provides optimization recommendations based on the provided core resources.
func Optimization(paidTier bool, featureEnabled bool, coreResources []tfschema.CoreResource, allCostResources []*tfschema.Resource, priceFetcher *prices.PriceFetcher) []optimization.OptimizationRecommendation {
	recommendations := make([]optimization.OptimizationRecommendation, 0)
	if !featureEnabled {
		return recommendations
	}

	// Savings Plans (New Structured)
	if savingPlans := optimization.GetSavingPlans(allCostResources, priceFetcher); len(savingPlans) > 0 {
		recommendations = append(recommendations, savingPlans...)
	}

	if !paidTier && len(recommendations) > 0 {
		// Teaser logic
		// Calculate sum of maximum savings per resource
		grouped, _ := optimization.GroupOptimizations(recommendations)
		var totalSavings float64

		for _, opts := range grouped {
			maxSavings := 0.0
			for _, opt := range opts {
				if opt.SavingsAmount > maxSavings {
					maxSavings = opt.SavingsAmount
				}
			}
			totalSavings += maxSavings
		}

		// Calculate overall percentage
		totalCost := 0.0
		for _, costResource := range allCostResources {
			if costResource.MonthlyCost != nil {
				totalCost += costResource.MonthlyCost.InexactFloat64()
			}
		}

		yearlyCost := totalCost * 12
		var savingsPct float64
		if yearlyCost > 0 {
			// totalSavings is monthly now, so convert to yearly for percentage calculation or use monthly cost
			savingsPct = ((totalSavings * 12) / yearlyCost) * 100
		}

		msg := fmt.Sprintf("We found reservations that could save you ~$%s/year (%.0f%%).\n [ Upgrade to Pro to see details ]\n", addCommas(fmt.Sprintf("%.0f", totalSavings*12)), savingsPct)

		return []optimization.OptimizationRecommendation{{
			Description: msg,
			Type:        "Teaser",
		}}
	}

	if paidTier {
		// Run optimization rules
		rules := optimization.GetRules()
		for _, cr := range coreResources {
			for _, rule := range rules {
				if rec := rule.Apply(&cr); rec != nil {
					recommendations = append(recommendations, optimization.OptimizationRecommendation{
						ResourceAddress: cr.BuildResource().Name,
						Description:     *rec,
						Type:            "Advisory",
					})
				}
			}
		}
	}

	return recommendations
}

func ConvertRecommendationsToAttrValue(opportunities []optimization.OptimizationRecommendation) types.List {
	recList := make([]attr.Value, 0, len(opportunities))
	recType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"resource_address":   types.StringType,
			"description":        types.StringType,
			"type":               types.StringType,
			"term":               types.StringType,
			"savings_amount":     types.NumberType,
			"savings_percentage": types.NumberType,
		},
	}

	for _, opt := range opportunities {
		obj, _ := types.ObjectValue(
			recType.AttrTypes,
			map[string]attr.Value{
				"resource_address":   types.StringValue(opt.ResourceAddress),
				"description":        types.StringValue(opt.Description),
				"type":               types.StringValue(opt.Type),
				"term":               types.StringValue(opt.Term),
				"savings_amount":     types.NumberValue(decimal.NewFromFloat(opt.SavingsAmount).BigFloat()),
				"savings_percentage": types.NumberValue(decimal.NewFromFloat(opt.SavingsPercentage).BigFloat()),
			},
		)
		recList = append(recList, obj)
	}
	return types.ListValueMust(recType, recList)
}

// ApplyDiscount applies discounts to the resource and its sub-resources based on the provided discount configurations. Available in all tiers.
func ApplyDiscount(resource *tfschema.Resource, discounts []DiscountModel) {
	var applicableDiscount float64
	foundSpecific := false

	// First pass: look for specific resource type discount
	for _, d := range discounts {
		if !d.ResourceType.IsNull() && !d.ResourceType.IsUnknown() {
			if d.ResourceType.ValueString() == resource.ResourceType {
				perc, _ := d.Percentage.ValueBigFloat().Float64()
				applicableDiscount = perc
				foundSpecific = true
				break
			}
		}
	}

	// Second pass: if no specific found, look for global discount
	if !foundSpecific {
		for _, d := range discounts {
			if d.ResourceType.IsNull() || d.ResourceType.IsUnknown() {
				perc, _ := d.Percentage.ValueBigFloat().Float64()
				applicableDiscount = perc
				break
			}
		}
	}

	if applicableDiscount > 0 {
		for _, component := range resource.CostComponents {
			component.MonthlyDiscountPerc = applicableDiscount
		}
	}

	for _, sub := range resource.SubResources {
		ApplyDiscount(sub, discounts)
	}
}
