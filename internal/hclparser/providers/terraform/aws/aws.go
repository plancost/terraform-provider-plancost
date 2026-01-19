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

package aws

import (
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/logging"
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

var DefaultProviderRegion = "us-east-1"
var arnAttributeMap = map[string]string{
	"aws_cloudwatch_dashboard":     "dashboard_arn",
	"aws_db_snapshot":              "db_snapshot_arn",
	"aws_db_cluster_snapshot":      "db_cluster_snapshot_arn",
	"aws_ecs_service":              "id",
	"aws_neptune_cluster_snapshot": "db_cluster_snapshot_arn",
	"aws_docdb_cluster_snapshot":   "db_cluster_snapshot_arn",
	"aws_dms_certificate":          "certificate_arn",
	"aws_dms_endpoint":             "endpoint_arn",
	"aws_dms_replication_instance": "replication_instance_arn",
	"aws_dms_replication_task":     "replication_task_arn",
}

func GetDefaultRefIDFunc(d *schema.ResourceData) []string {
	defaultRefs := []string{d.Get("id").String(), d.Get("name").String()}

	arnAttr, ok := arnAttributeMap[d.Type]
	if !ok {
		arnAttr = "arn"
	}

	if d.Get(arnAttr).Exists() {
		defaultRefs = append(defaultRefs, d.Get(arnAttr).String())
	}

	return defaultRefs
}

func DefaultCloudResourceIDFunc(d *schema.ResourceData) []string {
	var ids []string

	id := d.Get("id").String()
	if id != "" && id != "none" && !strings.HasPrefix(id, "hcl-") {
		ids = append(ids, id)
	}

	arnAttr, ok := arnAttributeMap[d.Type]
	if !ok {
		arnAttr = "arn"
	}

	arn := d.Get(arnAttr).String()
	if strings.HasPrefix(arn, "arn:aws:") && !strings.HasPrefix(arn, "arn:aws:hcl") {
		ids = append(ids, arn)
	}

	return ids
}

func GetSpecialContext(d *schema.ResourceData) map[string]interface{} {

	specialContexts := make(map[string]interface{})

	if strings.HasPrefix(d.Get("region").String(), "cn-") {
		specialContexts["isAWSChina"] = true
	}

	return specialContexts
}

func GetResourceRegion(d *schema.ResourceData) string {
	v := d.RawValues

	// If a region key exists in the values use that
	if v.Get("region").Exists() && v.Get("region").String() != "" {
		return v.Get("region").String()
	}

	// Otherwise try and parse the ARN from the values
	arnAttr, ok := arnAttributeMap[d.Type]
	if !ok {
		arnAttr = "arn"
	}

	if !v.Get(arnAttr).Exists() {
		return ""
	}

	arn := v.Get(arnAttr).String()
	p := strings.Split(arn, ":")
	if len(p) < 4 {
		logging.Logger.Debug().Msgf("Unexpected ARN format for %s", arn)
		return ""
	}

	return p[3]
}
