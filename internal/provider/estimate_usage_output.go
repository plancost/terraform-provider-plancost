/*
Copyright (c) 2026 Plancost.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package provider

import (
	"fmt"
	"sort"

	"github.com/plancost/terraform-provider-plancost/internal/schema"
	"github.com/plancost/terraform-provider-plancost/internal/usage"
	"gopkg.in/yaml.v3"
)

func GenerateUsageYAML(resources []*schema.Resource) ([]byte, error) {
	defaults, err := GetUsageDefaults()
	if err != nil {
		return nil, fmt.Errorf("failed to load usage defaults: %w", err)
	}

	rootNode := &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "version"},
					{Kind: yaml.ScalarNode, Value: "0.1"},
					{Kind: yaml.ScalarNode, Value: "resource_usage"},
					{Kind: yaml.MappingNode, Content: []*yaml.Node{}},
				},
			},
		},
	}

	resourceUsageNode := rootNode.Content[0].Content[3]

	// Sort resources by name for deterministic output
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})

	for _, res := range resources {
		if len(res.UsageSchema) == 0 {
			continue
		}

		resDefaults := defaults[res.ResourceType]
		resNode := buildResourceUsageNode(res.UsageSchema, resDefaults)

		resourceUsageNode.Content = append(resourceUsageNode.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: res.Name},
			resNode,
		)
	}

	return yaml.Marshal(rootNode)
}

func buildResourceUsageNode(usageSchema []*schema.UsageItem, defaults map[string]UsageDefault) *yaml.Node {
	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{},
	}

	for _, item := range usageSchema {
		key := item.Key
		var defaultVal UsageDefault
		if defaults != nil {
			if v, ok := defaults[key]; ok {
				defaultVal = v
			}
		}

		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key}
		var valNode *yaml.Node

		if item.ValueType == schema.SubResourceUsage {
			var subDefaults map[string]UsageDefault
			if defaultVal.Items != nil {
				subDefaults = defaultVal.Items
			}

			if item.DefaultValue != nil {
				if subRes, ok := item.DefaultValue.(*usage.ResourceUsage); ok {
					valNode = buildResourceUsageNode(subRes.Items, subDefaults)
				}
			}
			if valNode == nil {
				valNode = &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{}}
			}
		} else {
			val := item.DefaultValue
			if defaultVal.Value != nil {
				val = defaultVal.Value
			}
			valStr := fmt.Sprintf("%v", val)
			valNode = &yaml.Node{Kind: yaml.ScalarNode, Value: valStr}
			if defaultVal.Comment != "" {
				valNode.LineComment = "# " + defaultVal.Comment
			}
		}

		node.Content = append(node.Content, keyNode, valNode)
	}

	return node
}
