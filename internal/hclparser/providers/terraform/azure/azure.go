// Copyright 2021 Infracost Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azure

import (
	"github.com/plancost/terraform-provider-plancost/internal/hclparser/provider_schemas"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

var DefaultProviderRegion = "eastus"

func GetDefaultRefIDFunc(d *schema.ResourceData) []string {
	return []string{d.Get("id").String()}
}

func DefaultCloudResourceIDFunc(d *schema.ResourceData) []string {
	return []string{}
}

func GetSpecialContext(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{}
}

func ParseTags(externalTags map[string]string, r *schema.ResourceData) (map[string]string, []string) {
	_, supportsTags := provider_schemas.AzureTagsSupport[r.Type]
	rTags := r.Get("tags").Map()
	missing := schema.ExtractMissingVarsCausingMissingAttributeKeys(r, "tags")
	if !supportsTags && len(rTags) == 0 {
		return nil, missing
	}
	tags := make(map[string]string)
	for k, v := range rTags {
		tags[k] = v.String()
	}
	for k, v := range externalTags {
		tags[k] = v
	}
	return tags, missing
}
