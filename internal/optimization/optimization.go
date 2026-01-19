package optimization

import tfschema "github.com/plancost/terraform-provider-plancost/internal/schema"

type Rule interface {
	Description() string
	Apply(resource *tfschema.CoreResource) *string
}

func GetRules() []Rule {
	return []Rule{
		&latestGenerationVMSizeOptimization{},
		&appServicePlanOptimization{},
		&storageAccountOptimization{},
		&diskOptimization{},
		&vmOptimization{},
		&databaseOptimization{},
		&monitorOptimization{},
	}
}
