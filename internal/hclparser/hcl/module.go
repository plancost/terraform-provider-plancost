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
	"strconv"

	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

// ModuleCall represents a call to a defined Module by a parent Module.
type ModuleCall struct {
	// Name the name of the module as specified a the point of definition.
	Name string
	// Path is the path to the local directory containing the HCL for the Module.
	Path string
	// Definition is the actual Block where the ModuleCall happens in a hcl.File
	Definition *Block
	// Module contains the parsed root module that represents this ModuleCall.
	Module *Module
}

// Module encapsulates all the Blocks that are part of a Module in a Terraform project.
type Module struct {
	Name   string
	Source string

	Blocks Blocks
	// RawBlocks are the Blocks that were built when the module was loaded from the filesystem.
	// These are safe to pass to the child module calls as they are yet to be expanded.
	RawBlocks  Blocks
	RootPath   string
	ModulePath string

	Modules  []*Module
	Parent   *Module
	Warnings []*schema.ProjectDiag

	HasChanges         bool
	TerraformVarsPaths []string

	// ModuleSuffix is a unique name that can be optionally appended to the Module's
	// project name. This is only applicable to root modules.
	ModuleSuffix string

	// SourceURL is the discovered remote url for the module. This will only be
	// filled if the module is a remote module.
	SourceURL string

	// ProviderReferences is a map of provider names (relative to the module) to
	// the provider block that defines that provider. We keep track of this so we
	// can re-evaluate the provider blocks when we need to.
	ProviderReferences  map[string]*Block
	TerraformVersion    string
	ProviderConstraints *ProviderConstraints
}

// Index returns the count index of the Module using the name.
// Index returns nil if the Module has no count.
func (m *Module) Index() *int64 {
	matches := countRegex.FindStringSubmatch(m.Name)

	if len(matches) > 0 {
		i, _ := strconv.ParseInt(matches[1], 10, 64)

		return &i
	}

	return nil
}

// Key returns the foreach key of the Module using the name.
// Key returns nil if the Module has no each key.
func (m *Module) Key() *string {
	matches := foreachRegex.FindStringSubmatch(m.Name)

	if len(matches) > 0 {
		return &matches[1]
	}

	return nil
}
