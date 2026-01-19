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
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getCognitiveAccountRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_cognitive_account",
		CoreRFunc: newCognitiveAccount,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}

func newCognitiveAccount(d *schema.ResourceData) schema.CoreResource {
	region := d.Region
	kind := d.Get("kind").String()

	if strings.EqualFold(kind, "speechservices") {
		return &azure.CognitiveAccountSpeech{
			Address: d.Address,
			Region:  region,
			Sku:     d.Get("sku_name").String(),
		}
	}

	if strings.EqualFold(kind, "luis") {
		return &azure.CognitiveAccountLUIS{
			Address: d.Address,
			Region:  region,
			Sku:     d.Get("sku_name").String(),
		}
	}

	if strings.EqualFold(kind, "textanalytics") {
		return &azure.CognitiveAccountLanguage{
			Address: d.Address,
			Region:  region,
			Sku:     d.Get("sku_name").String(),
		}
	}

	if strings.EqualFold(kind, "openai") {
		// OpenAI costs are counted as part of a Cognitive Deployment so
		// this resource is counted as free
		return schema.BlankCoreResource{
			Name: d.Address,
			Type: d.Type,
		}
	}

	logging.Logger.Warn().Msgf("Skipping resource %s. Kind %q is not supported", d.Address, kind)

	return nil
}
