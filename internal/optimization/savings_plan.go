package optimization

import (
	"fmt"
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/prices"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
	"github.com/shopspring/decimal"
)

type candidate struct {
	ResourceName         string
	OriginalComponent    *schema.CostComponent
	ReservationComponent *schema.CostComponent
	TermLength           string // "1 yr" or "3 yr"
}

func GetSavingPlans(resources []*schema.Resource, priceFetcher *prices.PriceFetcher) []OptimizationRecommendation {
	recommendations := make([]OptimizationRecommendation, 0)

	// Create dummy resources for reservation lookups
	dummyResources := make([]*schema.Resource, 0)
	candidates := make([]candidate, 0)

	for _, res := range resources {
		for _, cc := range res.CostComponents {
			// Check if candidate
			if cc.Price().IsZero() || !isReservationCandidate(cc) {
				continue
			}

			candidateResources := candidatesForComponent(res, cc)
			for term, candidateResource := range candidateResources {
				dummyResources = append(dummyResources, candidateResource)
				var reservationComponent *schema.CostComponent
				if candidateResource != nil && len(candidateResource.CostComponents) != 0 {
					reservationComponent = candidateResource.CostComponents[0]
				}
				candidates = append(candidates, candidate{
					ResourceName:         res.Name,
					OriginalComponent:    cc,
					ReservationComponent: reservationComponent,
					TermLength:           term,
				})
			}
		}

		for _, subRes := range res.FlattenedSubResources() {
			for _, cc := range subRes.CostComponents {
				// Check if candidate
				if cc.Price().IsZero() || !isReservationCandidate(cc) {
					continue
				}

				candidateResources := candidatesForComponent(res, cc)
				for term, candidateResource := range candidateResources {
					dummyResources = append(dummyResources, candidateResource)
					var reservationComponent *schema.CostComponent
					if candidateResource != nil && len(candidateResource.CostComponents) != 0 {
						reservationComponent = candidateResource.CostComponents[0]
					}
					candidates = append(candidates, candidate{
						ResourceName:         res.Name,
						OriginalComponent:    cc,
						ReservationComponent: reservationComponent,
						TermLength:           term,
					})
				}
			}
		}
	}

	if len(dummyResources) == 0 {
		return recommendations
	}

	// Fetch prices
	err := priceFetcher.PopulatePrices(dummyResources)
	if err != nil {
		return recommendations
	}

	// Calculate savings
	for _, cand := range candidates {
		// Check if price found
		if cand.ReservationComponent.Price().IsZero() {
			continue
		}

		// Calculate savings
		// Consumption Cost (Monthly)
		// Hourly Price * 730
		hourlyPrice := cand.OriginalComponent.Price()
		monthlyConsumption := hourlyPrice.Mul(decimal.NewFromInt(730))

		// Reservation Cost (Monthly)
		// 1 Year: Price / 12
		// 3 Year: Price / 36
		resPrice := cand.ReservationComponent.Price()
		var monthlyResCost decimal.Decimal
		switch cand.TermLength {
		case "1 yr":
			monthlyResCost = resPrice.Div(decimal.NewFromInt(12))
		case "3 yr":
			monthlyResCost = resPrice.Div(decimal.NewFromInt(36))
		case "5 yr":
			monthlyResCost = resPrice.Div(decimal.NewFromInt(60))
		}

		savingsAmount := monthlyConsumption.Sub(monthlyResCost)
		if savingsAmount.LessThanOrEqual(decimal.Zero) {
			continue
		}

		savingsPct := savingsAmount.Div(monthlyConsumption).Round(2)

		// Format: "Save $369/mo (58%) on azurerm_linux_virtual_machine.vm with a 1-year Reservation"
		// Term formatting: "1 yr" -> "1-year"
		formattedTerm := strings.Replace(cand.TermLength, " yr", "-year", 1)

		description := fmt.Sprintf("Save $%.0f/mo (%.0f%%) on %s with a %s Reservation",
			savingsAmount.InexactFloat64(),
			savingsPct.InexactFloat64()*100,
			cand.ResourceName,
			formattedTerm,
		)

		recommendations = append(recommendations, OptimizationRecommendation{
			ResourceAddress:   cand.ResourceName,
			Description:       description,
			Type:              "Reservation",
			Term:              cand.TermLength,
			SavingsAmount:     savingsAmount.InexactFloat64(),
			SavingsPercentage: savingsPct.InexactFloat64(),
		})
	}

	return recommendations
}

func candidatesForComponent(res *schema.Resource, cc *schema.CostComponent) map[string]*schema.Resource {
	out := make(map[string]*schema.Resource)
	terms := []string{"1 yr", "3 yr", "5 yr"}
	for _, term := range terms {
		// Create dummy component
		dummyCC := &schema.CostComponent{
			Name:          cc.Name + " (" + term + " Reservation)",
			Unit:          cc.Unit,
			ProductFilter: cc.ProductFilter, // Same product
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Reservation"),
				TermLength:     strPtr(term),
			},
			IgnoreIfMissingPrice: true,
		}
		// Create a dummy resource for this candidate
		dummyRes := &schema.Resource{
			Name:           res.Name + " (Reservation Check)",
			ResourceType:   res.ResourceType,
			CostComponents: []*schema.CostComponent{dummyCC},
		}
		out[term] = dummyRes
	}
	return out
}

func isReservationCandidate(cc *schema.CostComponent) bool {
	if cc.PriceFilter != nil && cc.PriceFilter.PurchaseOption != nil && *cc.PriceFilter.PurchaseOption == "Reservation" {
		return false
	}

	// serviceFamily check
	if cc.ProductFilter == nil || cc.ProductFilter.ProductFamily == nil {
		return false
	}

	var unit string
	if cc.PriceFilter != nil && cc.PriceFilter.Unit != nil {
		unit = *cc.PriceFilter.Unit
	} else {
		// Map from cc.Unit
		if strings.EqualFold(cc.Unit, "hours") {
			unit = "1 Hour"
		} else if strings.EqualFold(cc.Unit, "months") {
			unit = "1/Month"
		}
	}

	serviceFamily := *cc.ProductFilter.ProductFamily
	switch unit {
	case "1 Hour":
		// Compute, Databases, Analytics
		return serviceFamily == "Compute" || serviceFamily == "Databases" || serviceFamily == "Analytics"
	case "1/Month":
		// Storage
		return serviceFamily == "Storage"
	default:
		return false
	}
}

func strPtr(s string) *string {
	return &s
}
