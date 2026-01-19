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

package azure

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

const (
	vendorName = "azure"
)

var (
	priceFilterConsumption = &schema.PriceFilter{
		PurchaseOption: strPtr("Consumption"),
	}
	priceFilterDevTestConsumption = &schema.PriceFilter{
		PurchaseOption: strPtr("DevTestConsumption"),
	}
)

func strPtr(s string) *string {
	return &s
}

func decimalPtr(d decimal.Decimal) *decimal.Decimal {
	return &d
}

func int32Ptr(i int32) *int32 {
	return &i
}

func intPtrToDecimalPtr(i *int64) *decimal.Decimal {
	if i == nil {
		return nil
	}
	return decimalPtr(decimal.NewFromInt(*i))
}

func floatPtrToDecimalPtr(f *float64) *decimal.Decimal {
	if f == nil {
		return nil
	}
	return decimalPtr(decimal.NewFromFloat(*f))
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func containsInt64(arr []int64, val int64) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func regexPtr(regex string) *string {
	return strPtr(fmt.Sprintf("/%s/i", regex))
}

func convertRegion(region string) string {
	if strings.Contains(strings.ToLower(region), "usgov") {
		return "US Gov"
	} else if strings.Contains(strings.ToLower(region), "china") {
		return "Сhina"
	} else {
		return "Global"
	}
}

func dnsZoneRegion(region string) string {
	switch {
	case strings.HasPrefix(strings.ToLower(region), "usgov"):
		return "US Gov Zone 1"
	case strings.HasPrefix(strings.ToLower(region), "germany"):
		return "DE Zone 1"
	case strings.HasPrefix(strings.ToLower(region), "china"):
		return "Zone 1 (China)"
	default:
		return "Zone 1"
	}
}
