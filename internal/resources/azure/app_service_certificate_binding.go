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
	"github.com/plancost/terraform-provider-plancost/internal/resources"
	"github.com/plancost/terraform-provider-plancost/internal/schema"

	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type AppServiceCertificateBinding struct {
	Address  string
	Region   string
	SSLState string
}

func (r *AppServiceCertificateBinding) CoreType() string {
	return "AppServiceCertificateBinding"
}

func (r *AppServiceCertificateBinding) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{}
}

func (r *AppServiceCertificateBinding) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *AppServiceCertificateBinding) BuildResource() *schema.Resource {
	region := r.Region

	var sslType string

	sslState := strings.ToUpper(r.SSLState)

	if strings.HasPrefix(sslState, "IP") {
		sslType = "IP"
	} else {

		return &schema.Resource{
			Name:        r.Address,
			NoPrice:     true,
			IsSkipped:   true,
			UsageSchema: r.UsageSchema(),
		}
	}

	var instanceCount int64 = 1

	costComponents := []*schema.CostComponent{
		{
			Name:            "IP SSL certificate",
			Unit:            "months",
			UnitMultiplier:  decimal.NewFromInt(1),
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(instanceCount)),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("azure"),
				Region:        strPtr(region),
				Service:       strPtr("Azure App Service"),
				ProductFamily: strPtr("Compute"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "skuName", Value: strPtr(fmt.Sprintf("%s SSL", sslType))},
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
