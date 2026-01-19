/*
Copyright (c) 2026 Plancost.
Licensed under the Business Source License 1.1.
See LICENSE file in the project root for full license information.
*/

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/plancost/terraform-provider-plancost/internal/terraform/azurerm"
)

func main() {
	// Map to store resource name -> pricing ("Free" or "Paid")
	resourceMap := make(map[string]string)

	// 1. Add explicitly free resources from the registry
	for _, name := range azurerm.FreeResources {
		resourceMap[name] = "Free"
	}

	// 2. Process resources in the registry
	for _, item := range azurerm.ResourceRegistry {
		// If already marked as free, skip (but ensure it's in the map)
		if val, ok := resourceMap[item.Name]; ok && val == "Free" {
			continue
		}

		pricing := "Paid"
		if item.NoPrice {
			pricing = "Free"
		}
		resourceMap[item.Name] = pricing
	}

	// Sort resources by name
	var names []string
	for name := range resourceMap {
		names = append(names, name)
	}
	sort.Strings(names)

	var supportedResources []string
	for _, name := range names {
		supportedResources = append(supportedResources, fmt.Sprintf("| %s | %s |", name, resourceMap[name]))
	}

	totalCount := len(supportedResources)

	content := fmt.Sprintf(`---
page_title: "Supported Resources"
description: "A list of all resources supported by the plancost provider."
---

# Supported Resources

> **Note:** This page is auto-generated.

The following %d resources are supported by the `+"`plancost`"+` provider.

| Resource Name | Pricing |
| :--- | :--- |
%s
`, totalCount, strings.Join(supportedResources, "\n"))

	err := os.WriteFile("docs/guides/supported-resources.md", []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated docs/guides/supported-resources.md")
}
