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

package hcl

import (
	"sync"

	"github.com/rs/zerolog"
)

type VertexModuleExit struct {
	logger        zerolog.Logger
	moduleConfigs *ModuleConfigs
	block         *Block
}

func (v *VertexModuleExit) ID() string {
	return v.block.FullName()
}

func (v *VertexModuleExit) ModuleAddress() string {
	return v.block.ModuleAddress()
}

func (v *VertexModuleExit) References() []VertexReference {
	return []VertexReference{}
}

func (v *VertexModuleExit) Visit(mutex *sync.Mutex) error {

	moduleInstances := v.moduleConfigs.Get(v.block.FullName())

	for _, moduleInstance := range moduleInstances {
		e := moduleInstance.evaluator
		e.module = *e.collectModules()
		e.module.Blocks = e.filteredBlocks

		modCall := moduleInstance.moduleCall
		if modCall == nil {
			continue
		}
		modCall.Module = &e.module
	}

	return nil
}
