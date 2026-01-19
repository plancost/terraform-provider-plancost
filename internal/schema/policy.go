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

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type Policy struct {
	ID                 string           `json:"id"`
	Title              string           `json:"title"`
	Description        string           `json:"description"`
	ResourceType       string           `json:"resource_type"`
	ResourceAttributes json.RawMessage  `json:"resource_attributes"`
	Address            string           `json:"address"`
	Suggested          string           `json:"suggested"`
	NoCost             bool             `json:"no_cost"`
	Cost               *decimal.Decimal `json:"cost"`
}

type Policies []Policy

func (r Policies) Len() int {
	return len(r)
}

func (r Policies) Less(i, j int) bool {
	iSug := r[i]
	jSug := r[j]

	if iSug.Cost == nil && jSug.Cost == nil {
		return iSug.Address < jSug.Address
	}

	if iSug.Cost == nil {
		return false
	}

	if jSug.Cost == nil {
		return true
	}

	if iSug.Cost.Equal(*jSug.Cost) {
		return iSug.Address < jSug.Address
	}

	return iSug.Cost.GreaterThan(*jSug.Cost)
}

func (r Policies) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
