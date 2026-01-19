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
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"github.com/zclconf/go-cty/cty"
)

type VertexVariable struct {
	logger        zerolog.Logger
	moduleConfigs *ModuleConfigs
	block         *Block
}

func (v *VertexVariable) ID() string {
	return v.block.FullName()
}

func (v *VertexVariable) ModuleAddress() string {
	return v.block.ModuleAddress()
}

func (v *VertexVariable) References() []VertexReference {
	return v.block.VerticesReferenced()
}

func (v *VertexVariable) Visit(mutex *sync.Mutex) error {

	moduleInstances := v.moduleConfigs.Get(v.block.ModuleAddress())
	if len(moduleInstances) == 0 {
		return fmt.Errorf("no module instances found for module address %q", v.block.ModuleAddress())
	}

	for _, moduleInstance := range moduleInstances {
		e := moduleInstance.evaluator
		blockInstance := e.module.Blocks.FindLocalName(v.block.LocalName())

		if blockInstance == nil {
			return fmt.Errorf("could not find block %q in module %q", v.ID(), moduleInstance.name)
		}

		// Re-evaluate the matching module input variables for this variable block
		// to ensure we have the most up-to-date value
		inputVars := e.inputVars

		if moduleInstance.moduleCall != nil {
			attrName := v.block.TypeLabel()
			attr := moduleInstance.moduleCall.Definition.GetAttribute(attrName)
			if attr != nil {
				inputVars = map[string]cty.Value{
					attrName: attr.Value(),
				}
			}
		}

		val, err := e.evaluateVariable(blockInstance, inputVars)
		if err != nil {
			return fmt.Errorf("could not evaluate variable %s: %w", v.ID(), err)
		}

		v.logger.Debug().Msgf("adding variable %s to the evaluation context", v.ID())
		key := fmt.Sprintf("var.%s", blockInstance.Label())
		e.ctx.SetByDot(val, key)

		e.AddFilteredBlocks(blockInstance)
	}

	return nil
}
