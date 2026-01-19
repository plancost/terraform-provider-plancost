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

type ProductFilter struct {
	VendorName       *string            `json:"vendorName,omitempty"`
	Service          *string            `json:"service,omitempty"`
	ProductFamily    *string            `json:"productFamily,omitempty"`
	Region           *string            `json:"region,omitempty"`
	Sku              *string            `json:"sku,omitempty"`
	AttributeFilters []*AttributeFilter `json:"attributeFilters,omitempty"`
}

type PriceFilter struct {
	PurchaseOption     *string `json:"purchaseOption,omitempty"`
	Unit               *string `json:"unit,omitempty"`
	Description        *string `json:"description,omitempty"`
	DescriptionRegex   *string `json:"description_regex,omitempty"`
	StartUsageAmount   *string `json:"startUsageAmount,omitempty"`
	EndUsageAmount     *string `json:"endUsageAmount,omitempty"`
	TermLength         *string `json:"termLength,omitempty"`
	TermPurchaseOption *string `json:"termPurchaseOption,omitempty"`
	TermOfferingClass  *string `json:"termOfferingClass,omitempty"`
}

type AttributeFilter struct {
	Key        string  `json:"key"`
	Value      *string `json:"value,omitempty"`
	ValueRegex *string `json:"value_regex,omitempty"`
}
