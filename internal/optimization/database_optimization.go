package optimization

import (
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type databaseOptimization struct{}

func (r databaseOptimization) Description() string {
	return "Suggests optimizations for Azure Databases (SQL, PostgreSQL, etc.)."
}

func (r databaseOptimization) Apply(input *schema.CoreResource) *string {
	if input == nil {
		return nil
	}

	// SQL Database
	if resource, ok := (*input).(*azure.SQLDatabase); ok {
		// SQL - consider using Azure Hybrid Benefit for SQL Server
		if strings.EqualFold(resource.LicenseType, "LicenseIncluded") {
			return stringPtr("Consider using Azure Hybrid Benefit for SQL Server if you have eligible on-premises licenses")
		}

		// SQL - consider using serverless with auto-pause in non-production
		// Check if it's General Purpose Provisioned (e.g., GP_Gen5_...) and not Serverless (GP_S_...)
		// This is a heuristic.
		if strings.HasPrefix(resource.SKU, "GP_") && !strings.Contains(resource.SKU, "_S_") {
			return stringPtr("Consider using SQL Database Serverless tier with auto-pause for non-production workloads")
		}
	}

	// SQL Managed Instance
	if resource, ok := (*input).(*azure.SQLManagedInstance); ok {
		if strings.EqualFold(resource.LicenseType, "LicenseIncluded") {
			return stringPtr("Consider using Azure Hybrid Benefit for SQL Managed Instance if you have eligible on-premises licenses")
		}
		if strings.EqualFold(resource.StorageAccountType, "GRS") || strings.EqualFold(resource.StorageAccountType, "GZRS") {
			return stringPtr("Consider removing geo-redundant backups (use LRS/ZRS) in non-production projects")
		}
	}

	// MSSQL Managed Instance
	if resource, ok := (*input).(*azure.MSSQLManagedInstance); ok {
		if strings.EqualFold(resource.LicenseType, "LicenseIncluded") {
			return stringPtr("Consider using Azure Hybrid Benefit for SQL Managed Instance if you have eligible on-premises licenses")
		}
		if strings.EqualFold(resource.StorageAccountType, "GRS") || strings.EqualFold(resource.StorageAccountType, "GZRS") {
			return stringPtr("Consider removing geo-redundant backups (use LRS/ZRS) in non-production projects")
		}
	}

	return nil
}

var _ Rule = &databaseOptimization{}
