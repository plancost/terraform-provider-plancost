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

package resources

import (
	"github.com/plancost/terraform-provider-plancost/internal/resources/types"
)

// ResourceRegistry manages the collection of cost estimable resources
type ResourceRegistry struct {
	resources map[string]types.Resource
}

// NewResourceRegistry creates a new resource registry
func NewResourceRegistry() *ResourceRegistry {
	return &ResourceRegistry{
		resources: make(map[string]types.Resource),
	}
}

// Register adds a resource to the registry
func (r *ResourceRegistry) Register(resource types.Resource) {
	r.resources[resource.CoreType()] = resource
}

// Get retrieves a resource by core type
func (r *ResourceRegistry) Get(coreType string) (types.Resource, bool) {
	res, ok := r.resources[coreType]
	return res, ok
}

var defaultRegistry *ResourceRegistry

// DefaultRegistry returns the default resource registry with all resources registered
func DefaultRegistry() *ResourceRegistry {
	return defaultRegistry
}

func init() {
	// Register all resources when the package is initialized
	defaultRegistry = NewResourceRegistry()

	// TODO: Register Azure resources
	// Note: Direct import of azure package creates circular dependency
	// Resources should be registered from outside this package

	// Register GCP resources when available
	// for _, r := range gcp.Resources() {
	// 	defaultRegistry.Register(r)
	// }

	// Register AWS resources when available
	// for _, r := range aws.Resources() {
	// 	defaultRegistry.Register(r)
	// }
}
