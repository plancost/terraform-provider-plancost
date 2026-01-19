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

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstraintsAllowVersionOrAbove(t *testing.T) {
	tests := []struct {
		name            string
		constraints     string
		requiredVersion string
		want            bool
	}{
		{
			name:            "bad/empty constraints",
			constraints:     "",
			requiredVersion: "5.0.0",
			want:            true,
		},
		{
			name:            "multi constraints",
			constraints:     "> 3.0.0, < 6.0.0",
			requiredVersion: "5.0.0",
			want:            true,
		},
		{
			name:            "simple match",
			constraints:     "5.0.0",
			requiredVersion: "5.0.0",
			want:            true,
		},
		{
			name:            "simple match via =",
			constraints:     "= 5.0.0",
			requiredVersion: "5.0.0",
			want:            true,
		},
		{
			name:            "simple mismatch",
			constraints:     "5.0.0",
			requiredVersion: "5.0.1",
			want:            false,
		},
		{
			name:            "constraints require greater",
			constraints:     "> 5.0.0",
			requiredVersion: "5.0.0",
			want:            true,
		},
		{
			name:            "constraints require greater than previous version",
			constraints:     "> 1.0.0",
			requiredVersion: "5.0.0",
			want:            true,
		},
		{
			name:            "constraints require >= version",
			constraints:     ">= 5.0.0",
			requiredVersion: "5.0.0",
			want:            true,
		},
		{
			name:            "patch only match on same major",
			constraints:     "~> 5.0.0",
			requiredVersion: "5.0.2",
			want:            true,
		},
		{
			name:            "patch only mismatch on lower major",
			constraints:     "~> 1.0.0",
			requiredVersion: "5.0.2",
			want:            false,
		},
		{
			name:            "patch only match on higher major",
			constraints:     "~> 6.0.0",
			requiredVersion: "5.0.2",
			want:            true,
		},
		{
			name:            "simple mismatch on lower constraint",
			constraints:     "< 5.0.0",
			requiredVersion: "5.0.0",
			want:            false,
		}, {
			name:            "simple match on <= constraint",
			constraints:     "<= 5.0.0",
			requiredVersion: "5.0.0",
			want:            true,
		},
		{
			name:            "simple match on lower constraint",
			constraints:     "< 5.0.1",
			requiredVersion: "5.0.0",
			want:            true,
		},
		{
			name:            "example",
			constraints:     "~> 5.36",
			requiredVersion: "3.38",
			want:            true,
		},
		{
			name:            "example",
			constraints:     "~> 3.0",
			requiredVersion: "3.38",
			want:            true,
		},
		{
			name:            "example",
			constraints:     "~> 3.0.1",
			requiredVersion: "3.38",
			want:            false,
		},
		{
			name:            "example",
			constraints:     "~> 2.0",
			requiredVersion: "3.38",
			want:            false,
		},
		{
			name:            "example",
			constraints:     "~> 2",
			requiredVersion: "3.38",
			want:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := version.NewConstraint(tt.constraints)
			v, err := version.NewVersion(tt.requiredVersion)
			require.NoError(t, err)
			assert.Equal(t, tt.want, ConstraintsAllowVersionOrAbove(c, v), "constraint %s does not allow %s or greater", tt.constraints, tt.requiredVersion)
		})
	}
}
