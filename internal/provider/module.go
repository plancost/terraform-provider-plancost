/*
Copyright (c) 2026 Plancost.
Licensed under the Business Source License 1.1.
See LICENSE file in the project root for full license information.
*/

package provider

import (
	"fmt"

	"github.com/plancost/terraform-provider-plancost/internal/hclparser/hcl"
	"github.com/plancost/terraform-provider-plancost/internal/hclparser/terraform"
	tfschema "github.com/plancost/terraform-provider-plancost/internal/schema"
)

func ParseModule(moduleSourceDir string, usageDataMap tfschema.UsageMap, variableOptions ...hcl.Option) ([]*tfschema.Resource, []tfschema.CoreResource, error) {
	provider, err := terraform.NewHCLProvider(moduleSourceDir, &terraform.HCLProviderConfig{}, variableOptions...)
	if err != nil {
		return nil, nil, err
	}

	projects, err := provider.LoadResources(usageDataMap)
	if err != nil {
		return nil, nil, err
	}

	if len(projects) == 0 {
		return nil, nil, fmt.Errorf("LoadResources returns empty projects")
	}

	res := make([]*tfschema.Resource, 0)
	coreResources := make([]tfschema.CoreResource, 0)
	for _, rd := range projects[0].PartialResources {
		if rd.Resource != nil {
			rd.Resource.ResourceType = rd.Type
			rd.Resource.Tags = rd.Tags
			res = append(res, rd.Resource)
			continue
		}
		if rd.CoreResource != nil {
			coreResources = append(coreResources, rd.CoreResource)
			if rd.UsageData != nil {
				rd.CoreResource.PopulateUsage(rd.UsageData)
			}
			costResource := rd.CoreResource.BuildResource()
			if costResource == nil {
				continue
			}
			costResource.ResourceType = rd.Type
			costResource.Tags = rd.Tags
			res = append(res, costResource)
			continue
		}
	}
	return res, coreResources, nil
}
