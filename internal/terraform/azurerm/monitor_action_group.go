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
	"fmt"

	"github.com/plancost/terraform-provider-plancost/internal/resources/azure"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

func getMonitorActionGroupRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:      "azurerm_monitor_action_group",
		CoreRFunc: newMonitorActionGroup,
		ReferenceAttributes: []string{
			"resource_group_name",
		},
	}
}

func newMonitorActionGroup(d *schema.ResourceData) schema.CoreResource {
	region := d.Region

	smsByCountryCode := make(map[int]int)
	for _, sms := range d.Get("sms_receiver").Array() {
		cc := int(sms.Get("country_code").Int())
		if cur, ok := smsByCountryCode[cc]; ok {
			smsByCountryCode[cc] = cur + 1
		} else {
			smsByCountryCode[cc] = 1
		}
	}

	voiceByCountryCode := make(map[int]int)
	for _, voice := range d.Get("voice_receiver").Array() {
		cc := int(voice.Get("country_code").Int())
		if cur, ok := voiceByCountryCode[cc]; ok {
			voiceByCountryCode[cc] = cur + 1
		} else {
			voiceByCountryCode[cc] = 1
		}
	}

	var secureWebhooks int
	var webhooks int
	for i := range d.Get("webhook_receiver").Array() {
		if d.IsEmpty(fmt.Sprintf("webhook_receiver.%d.aad_auth", i)) {
			webhooks += 1
		} else {
			secureWebhooks += 1
		}
	}

	return &azure.MonitorActionGroup{
		Address:                         d.Address,
		Region:                          region,
		EmailReceivers:                  len(d.Get("email_receiver").Array()),
		ITSMEventReceivers:              len(d.Get("itsm_receiver").Array()),
		PushNotificationReceivers:       len(d.Get("azure_app_push_receiver").Array()),
		SecureWebHookReceivers:          secureWebhooks,
		WebHookReceivers:                webhooks,
		SMSReceiversByCountryCode:       smsByCountryCode,
		VoiceCallReceiversByCountryCode: voiceByCountryCode,
	}
}
