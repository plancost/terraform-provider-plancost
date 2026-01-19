package optimization

import (
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type appServicePlanOptimization struct{}

func (r appServicePlanOptimization) Description() string {
	return "Suggests upgrading App Service Plans to newer versions for better performance and cost efficiency."
}

func (r appServicePlanOptimization) Apply(input *schema.CoreResource) *string {
	if input == nil {
		return nil
	}

	var sku string

	if resource, ok := (*input).(*azure.AppServicePlan); ok {
		sku = resource.SKUSize
	} else if resource, ok := (*input).(*azure.ServicePlan); ok {
		sku = resource.SKUName
	} else {
		return nil
	}

	sku = strings.ToUpper(sku)

	// App Service - consider upgrading Isolated service plans to v2
	if sku == "I1" || sku == "I2" || sku == "I3" {
		return stringPtr("Consider upgrading Isolated service plans to v2")
	}

	// App Service - consider upgrading Premium service plans to v3
	if sku == "P1" || sku == "P2" || sku == "P3" ||
		sku == "P1V2" || sku == "P2V2" || sku == "P3V2" {
		return stringPtr("Consider upgrading Premium service plans to v3")
	}

	// App Service - consider upgrading Standard service plans to Basic or Premium v3
	if sku == "S1" || sku == "S2" || sku == "S3" {
		return stringPtr("Consider upgrading Standard service plans to Basic or Premium v3")
	}

	return nil
}

var _ Rule = &appServicePlanOptimization{}
