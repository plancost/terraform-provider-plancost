package optimization

// OptimizationRecommendation represents a structured cost optimization recommendation.
type OptimizationRecommendation struct {
	ResourceAddress   string
	Description       string  // Description of the optimization opportunity
	Type              string  // e.g., "Reservation", "Advisor"
	Term              string  // e.g., "1 Year", "3 Year"
	SavingsAmount     float64 // The estimated monthly savings amount in USD
	SavingsPercentage float64
}
