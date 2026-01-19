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

package azurerm

import (
	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func GetAzureRMKeyVaultCertificateRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_key_vault_certificate",
		RFunc: NewAzureRMKeyVaultCertificate,
		ReferenceAttributes: []string{
			"key_vault_id",
		},
		GetRegion: func(defaultRegion string, d *schema.ResourceData) string {
			return lookupRegion(d, []string{"key_vault_id"})
		},
	}
}

func NewAzureRMKeyVaultCertificate(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := d.Region

	var costComponents []*schema.CostComponent

	var skuName string
	keyVault := d.References("key_vault_id")
	if len(keyVault) > 0 {
		skuName = cases.Title(language.English).String(keyVault[0].Get("sku_name").String())
	} else {
		logging.Logger.Warn().Msgf("Skipping resource %s. Could not find its 'key_vault_id.sku_name' property.", d.Address)
		return nil
	}

	var certificateRenewals, certificateOperations *decimal.Decimal
	if u != nil && u.Get("monthly_certificate_renewal_requests").Exists() {
		certificateRenewals = decimalPtr(decimal.NewFromInt(u.Get("monthly_certificate_renewal_requests").Int()))
	}
	costComponents = append(costComponents, vaultKeysCostComponent(
		"Certificate renewals",
		region,
		"requests",
		skuName,
		"Certificate Renewal Request",
		"0",
		certificateRenewals,
		1))

	if u != nil && u.Get("monthly_certificate_other_operations").Exists() {
		certificateOperations = decimalPtr(decimal.NewFromInt(u.Get("monthly_certificate_other_operations").Int()))
	}
	costComponents = append(costComponents, vaultKeysCostComponent(
		"Certificate operations",
		region,
		"10K transactions",
		skuName,
		"Operations",
		"0",
		certificateOperations,
		10000))

	return &schema.Resource{
		Name:           d.Address,
		CostComponents: costComponents,
	}
}
