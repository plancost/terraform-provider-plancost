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
	"testing"

	"github.com/rs/zerolog"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
)

func TestBlock_LocalName(t *testing.T) {
	tests := []struct {
		name  string
		block *Block
		want  string
	}{
		{
			name: "resource Block with empty labels will return empty local name",
			block: &Block{
				HCLBlock: &hcl.Block{
					Type:   "resource",
					Labels: nil,
				},
				logger: newDiscardLogger(),
			},
			want: "",
		},
		{
			name: "resource Block with valid labels will return reference without resource type",
			block: &Block{
				HCLBlock: &hcl.Block{
					Type:   "resource",
					Labels: []string{"my-resource", "my-name"},
				},
				logger: newDiscardLogger(),
			},
			want: "my-resource.my-name",
		},
		{
			name: "data Block with valid labels will return reference with Block type",
			block: &Block{
				HCLBlock: &hcl.Block{
					Type:   "data",
					Labels: []string{"my-block", "my-name"},
				},
				logger: newDiscardLogger(),
			},
			want: "data.my-block.my-name",
		},
		{
			name: "dynamic block inside resource blocks will return reference with parent blocks",
			block: &Block{
				HCLBlock: &hcl.Block{
					Type:   "content",
					Labels: []string{},
				},
				parent: &Block{
					HCLBlock: &hcl.Block{
						Type:   "dynamic",
						Labels: []string{"my-dynamic-block"},
					},
					logger: newDiscardLogger(),
					parent: &Block{
						HCLBlock: &hcl.Block{
							Type:   "resource",
							Labels: []string{"my-resource", "my-name"},
						},
						logger: newDiscardLogger(),
					},
				},
				logger: newDiscardLogger(),
			},
			want: "my-resource.my-name.dynamic.my-dynamic-block.content",
		},
		{
			name: "dynamic block inside data blocks will return reference with parent blocks",
			block: &Block{
				HCLBlock: &hcl.Block{
					Type:   "content",
					Labels: []string{},
				},
				parent: &Block{
					HCLBlock: &hcl.Block{
						Type:   "dynamic",
						Labels: []string{"my-dynamic-block"},
					},
					logger: newDiscardLogger(),
					parent: &Block{
						HCLBlock: &hcl.Block{
							Type:   "data",
							Labels: []string{"my-block", "my-name"},
						},
						logger: newDiscardLogger(),
					},
				},
				logger: newDiscardLogger(),
			},
			want: "data.my-block.my-name.dynamic.my-dynamic-block.content",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.block.LocalName(), "LocalName()")
		})
	}
}

func newDiscardLogger() zerolog.Logger {
	return zerolog.Nop()
}
