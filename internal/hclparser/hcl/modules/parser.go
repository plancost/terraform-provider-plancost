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

package modules

import (
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type SharedHCLParser struct {
	parser *hclparse.Parser
	mu     *sync.Mutex
}

func NewSharedHCLParser() *SharedHCLParser {
	return &SharedHCLParser{
		parser: hclparse.NewParser(),
		mu:     &sync.Mutex{},
	}
}

func (p *SharedHCLParser) ParseHCLFile(filename string) (*hcl.File, hcl.Diagnostics) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.parser.ParseHCLFile(filename)
}

func (p *SharedHCLParser) ParseJSONFile(filename string) (*hcl.File, hcl.Diagnostics) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.parser.ParseJSONFile(filename)
}
