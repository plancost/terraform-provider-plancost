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

// ReferenceIDFunc is used to let references be built using non-standard IDs (anything other than d.Get("id").string)
type ReferenceIDFunc func(d *ResourceData) []string

// CloudResourceIDFunc is used to calculate the cloud resource ids (AWS ARN, Google HREF, etc...) associated with the resource
type CloudResourceIDFunc func(d *ResourceData) []string

// RegionLookupFunc is used to look up the region of a resource, this is used to
// calculate the region of a resource if the region requires a lookup from
// reference attributes.
type RegionLookupFunc func(defaultRegion string, d *ResourceData) string

type RegistryItem struct {
	Name                string
	Notes               []string
	RFunc               ResourceFunc
	CoreRFunc           CoreResourceFunc
	ReferenceAttributes []string
	CustomRefIDFunc     ReferenceIDFunc
	DefaultRefIDFunc    ReferenceIDFunc
	CloudResourceIDFunc CloudResourceIDFunc
	NoPrice             bool
	// GetRegion is used to look up the region of a resource if it has a region that
	// cannot be calculated from the default resource/provider data. If the GetRegion
	// is nil or the return result is empty the region will be calculated from the
	// default resource/provider data.
	GetRegion RegionLookupFunc
}
