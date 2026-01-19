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

	"github.com/stretchr/testify/assert"
)

func TestEnvFileMatcher_EnvName(t *testing.T) {
	type fields struct {
		envNames []string
	}
	type args struct {
		file string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "directly matches env name",
			fields: fields{
				envNames: []string{"dev", "prod"},
			},
			args: args{
				file: "dev.tfvars",
			},
			want: "dev",
		},
		{
			name: "directly matches env which collides",
			fields: fields{
				envNames: []string{"dev", "dev-legacy", "prod"},
			},
			args: args{
				file: "dev-legacy.tfvars",
			},
			want: "dev-legacy",
		},
		{
			name: "returns filename when no match",
			fields: fields{
				envNames: []string{"dev", "prod"},
			},
			args: args{
				file: "foo.tfvars",
			},
			want: "foo",
		},
		{
			name: "returns prefix",
			fields: fields{
				envNames: []string{"dev", "prod"},
			},
			args: args{
				file: "prod-defaults.tfvars",
			},
			want: "prod",
		},
		{
			name: "returns longest prefix match",
			fields: fields{
				envNames: []string{"dev", "prod", "prod-legacy"},
			},
			args: args{
				file: "prod-legacy-defaults.tfvars",
			},
			want: "prod-legacy",
		},
		{
			name: "returns suffix",
			fields: fields{
				envNames: []string{"dev", "prod"},
			},
			args: args{
				file: "defaults-prod.tfvars",
			},
			want: "prod",
		},
		{
			name: "returns longest suffix match",
			fields: fields{
				envNames: []string{"dev", "prod", "legacy-prod"},
			},
			args: args{
				file: "defaults-legacy-prod.tfvars",
			},
			want: "legacy-prod",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := CreateEnvFileMatcher(tt.fields.envNames, nil)
			assert.Equalf(t, tt.want, e.EnvName(tt.args.file), "EnvName(%v)", tt.args.file)
		})
	}
}
