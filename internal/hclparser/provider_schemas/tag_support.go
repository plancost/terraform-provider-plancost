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

package provider_schemas

import (
	_ "embed"

	json "github.com/json-iterator/go"
)

////go:embed aws.tags.json
//var awsTagsJSON []byte
//
////go:embed aws.tags_all.json
//var awsTagsAllJSON []byte
//
////go:embed aws.tag_block.json
//var awsTagBlockJSON []byte

//go:embed azurerm.tags.json
var azurermTagsJSON []byte

////go:embed google.labels.json
//var googleLabelsJSON []byte
//
////go:embed google.user_labels.json
//var googleUserLabelsJSON []byte
//
////go:embed google.settings_user_labels.json
//var googleSettingsUserLabelsJSON []byte

var AWSTagsSupport map[string]bool
var AWSTagsAllSupport map[string]bool
var AWSTagBlockSupport map[string]bool
var AzureTagsSupport map[string]bool
var GoogleLabelsSupport map[string]bool
var GoogleUserLabelsSupport map[string]bool
var GoogleSettingsUserLabelsSupport map[string]bool

func init() {
	//err := json.Unmarshal(awsTagsJSON, &AWSTagsSupport)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = json.Unmarshal(awsTagsAllJSON, &AWSTagsAllSupport)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = json.Unmarshal(awsTagBlockJSON, &AWSTagBlockSupport)
	//if err != nil {
	//	panic(err)
	//}

	err := json.Unmarshal(azurermTagsJSON, &AzureTagsSupport)
	if err != nil {
		panic(err)
	}

	//err = json.Unmarshal(googleLabelsJSON, &GoogleLabelsSupport)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = json.Unmarshal(googleUserLabelsJSON, &GoogleUserLabelsSupport)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = json.Unmarshal(googleSettingsUserLabelsJSON, &GoogleSettingsUserLabelsSupport)
	//if err != nil {
	//	panic(err)
	//}
}
