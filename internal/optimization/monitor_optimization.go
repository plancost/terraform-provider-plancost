package optimization

import (
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type monitorOptimization struct{}

func (r monitorOptimization) Description() string {
	return "Suggests optimizations for Azure Monitor and Log Analytics."
}

func (r monitorOptimization) Apply(input *schema.CoreResource) *string {
	if input == nil {
		return nil
	}
	resource, ok := (*input).(*azure.LogAnalyticsWorkspace)
	if !ok {
		return nil
	}

	// Monitor - consider using a retention policy to reduce storage costs
	// Free tier is 30 days. If retention is higher, suggest checking if it's needed.
	if resource.RetentionInDays > 30 {
		return stringPtr("Consider reducing Log Analytics retention period or using Archive tier for older data to reduce costs")
	}

	return nil
}

var _ Rule = &monitorOptimization{}
