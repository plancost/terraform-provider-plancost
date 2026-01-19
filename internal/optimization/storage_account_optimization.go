package optimization

import (
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type storageAccountOptimization struct{}

func (r storageAccountOptimization) Description() string {
	return "Suggests optimizations for Azure Storage Accounts."
}

func (r storageAccountOptimization) Apply(input *schema.CoreResource) *string {
	if input == nil {
		return nil
	}
	resource, ok := (*input).(*azure.StorageAccount)
	if !ok {
		return nil
	}

	// Storage Accounts - consider upgrading account kind to StorageV2
	if strings.EqualFold(resource.AccountKind, "Storage") || strings.EqualFold(resource.AccountKind, "BlobStorage") {
		return stringPtr("Consider upgrading account kind to StorageV2")
	}

	// Storage Accounts - consider using a preferred access tier
	if strings.EqualFold(resource.AccessTier, "Hot") {
		return stringPtr("Consider using Cool or Archive access tier for infrequently accessed data")
	}

	// Storage Accounts - consider using preferred replication type
	// Suggest LRS if using GRS/RA-GRS for non-critical workloads (heuristic)
	if strings.EqualFold(resource.AccountReplicationType, "GRS") ||
		strings.EqualFold(resource.AccountReplicationType, "RAGRS") ||
		strings.EqualFold(resource.AccountReplicationType, "GZRS") ||
		strings.EqualFold(resource.AccountReplicationType, "RAGZRS") {
		return stringPtr("Consider using LRS or ZRS replication for cost savings if geo-redundancy is not required")
	}

	return nil
}

var _ Rule = &storageAccountOptimization{}
