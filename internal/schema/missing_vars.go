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

package schema

func ExtractMissingVarsCausingMissingAttributeKeys(r *ResourceData, attribute string) []string {
	var missing []string
	if raw := r.Metadata["attributesWithUnknownKeys"]; raw.IsArray() {
		for _, el := range raw.Array() {
			if el.Get("attribute").String() == attribute {
				if vars := el.Get("missingVariables"); vars.IsArray() {
					for _, v := range vars.Array() {
						missing = append(missing, v.String())
					}
				}
			}
		}
	}
	return missing
}
