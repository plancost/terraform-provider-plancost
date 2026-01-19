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

package prices

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/plancost/terraform-provider-plancost/internal/apiclient"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func Test_notFound_Add(t *testing.T) {
	type args struct {
		results []apiclient.PriceQueryResult
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test aggregates resource/cost component with correct keys",
			args: args{results: []apiclient.PriceQueryResult{
				{
					PriceQueryKey: apiclient.PriceQueryKey{
						Resource:      &schema.Resource{ResourceType: "aws_instance"},
						CostComponent: &schema.CostComponent{Name: "Compute (on-demand, foo)"},
					},
				},
				{
					PriceQueryKey: apiclient.PriceQueryKey{
						Resource:      &schema.Resource{ResourceType: "aws_instance"},
						CostComponent: &schema.CostComponent{Name: "Data Storage"},
					},
				},
				{
					PriceQueryKey: apiclient.PriceQueryKey{
						Resource:      &schema.Resource{ResourceType: "aws_instance"},
						CostComponent: &schema.CostComponent{Name: "Compute (on-demand, bar)"},
					},
				},
			}},
			want: []string{"aws_instance.compute", "aws_instance.compute", "aws_instance.data_storage"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PriceFetcher{
				resources:  make(map[string]*notFoundData),
				components: make(map[string]int),
				mux:        &sync.RWMutex{},
			}
			for _, res := range tt.args.results {
				p.addNotFoundResult(res)

			}

			actual := p.MissingPricesComponents()
			assert.EqualValues(t, tt.want, actual)
		})
	}
}
