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

	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"

	"github.com/shopspring/decimal"
)

type AppServiceCertificateOrder struct {
	Address     string
	ProductType string
}

func (r *AppServiceCertificateOrder) CoreType() string {
	return "AppServiceCertificateOrder"
}

func (r *AppServiceCertificateOrder) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{}
}

func (r *AppServiceCertificateOrder) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *AppServiceCertificateOrder) BuildResource() *schema.Resource {
	region := "Global"

	if strings.HasPrefix(region, "usgov") {
		region = "US Gov"
	}

	productType := "Standard"
	if r.ProductType != "" {
		productType = r.ProductType
	}
	productType = strings.ToLower(productType)

	costComponents := []*schema.CostComponent{
		{
			Name:           fmt.Sprintf("SSL certificate (%s)", productType),
			Unit:           "years",
			UnitMultiplier: decimal.NewFromInt(1),

			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1).Div(decimal.NewFromInt(12))),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(region),
				Service:       strPtr("Azure App Service"),
				ProductFamily: strPtr("Compute"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "skuName", ValueRegex: strPtr(fmt.Sprintf("/%s SSL - 1 Year/i", productType))},
				},
			},
			PriceFilter: &schema.PriceFilter{
				PurchaseOption: strPtr("Consumption"),
			},
		},
	}

	return &schema.Resource{
		Name:           r.Address,
		CostComponents: costComponents,
		UsageSchema:    r.UsageSchema(),
	}
}
