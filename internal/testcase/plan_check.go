package testcase

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// ResourceCost represents a resource with its cost components
type ResourceCost struct {
	Name           string          `json:"name"`
	CostComponents []CostComponent `json:"costComponents"`
	SubResources   []ResourceCost  `json:"subResources"`
}

// CostComponent represents a single cost component
type CostComponent struct {
	Name            string `json:"name"`
	MonthlyQuantity string `json:"monthlyQuantity"`
	Unit            string `json:"unit"`
	MonthlyCost     string `json:"monthlyCost"`
}

// CostCheck checks the plan against expected costs
type CostCheck struct {
	ExpectedResources []ResourceCost
}

var _ plancheck.PlanCheck = &CostCheck{}

// NewResourceCostCheck creates a plan check from a list of ResourceCost
func NewResourceCostCheck(resources []ResourceCost) plancheck.PlanCheck {
	return &CostCheck{
		ExpectedResources: resources,
	}
}

func (c *CostCheck) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var resourceChange *tfjson.ResourceChange
	for _, rc := range req.Plan.ResourceChanges {
		if rc.Address == "plancost_estimate.this" {
			resourceChange = rc
			break
		}
	}

	if resourceChange == nil {
		resp.Error = fmt.Errorf("resource plancost_estimate.this not found in plan")
		return
	}

	after, ok := resourceChange.Change.After.(map[string]interface{})
	if !ok {
		resp.Error = fmt.Errorf("failed to cast resource change After to map[string]interface{}")
		return
	}

	resourcesRaw, ok := after["resources"]
	if !ok {
		resp.Error = fmt.Errorf("resources attribute not found in plancost_estimate.this")
		return
	}

	resourcesList, ok := resourcesRaw.([]interface{})
	if !ok {
		resp.Error = fmt.Errorf("resources attribute is not a list")
		return
	}

	// Convert actual resources to a map for easier lookup
	actualResources := make(map[string]map[string]interface{})
	for _, r := range resourcesList {
		rMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		name, ok := rMap["name"].(string)
		if ok {
			actualResources[name] = rMap
		}
	}

	var errors []string

	for _, expectedRes := range c.ExpectedResources {
		actualRes, ok := actualResources[expectedRes.Name]
		if !ok {
			errors = append(errors, fmt.Sprintf("Resource '%s' missing in plan", expectedRes.Name))
			continue
		}

		// Check Cost Components
		if errs := checkCostComponents(expectedRes.Name, expectedRes.CostComponents, actualRes["cost_components"]); len(errs) > 0 {
			errors = append(errors, errs...)
		}

		// Check Sub Resources
		if errs := checkSubResources(expectedRes.Name, expectedRes.SubResources, actualRes["sub_resources"]); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	if len(errors) > 0 {
		resp.Error = fmt.Errorf("golden file check failed:\n%s", strings.Join(errors, "\n"))
	}
}

func checkCostComponents(resourceName string, expected []CostComponent, actualRaw interface{}) []string {
	var errors []string
	actualList, ok := actualRaw.([]interface{})
	if !ok {
		if len(expected) == 0 {
			return nil
		}
		return []string{fmt.Sprintf("Resource '%s': cost_components is not a list", resourceName)}
	}

	actualMap := make(map[string]map[string]interface{})
	for _, item := range actualList {
		m, ok := item.(map[string]interface{})
		if ok {
			name, _ := m["name"].(string)
			actualMap[name] = m
		}
	}

	for _, exp := range expected {
		act, ok := actualMap[exp.Name]
		if !ok {
			errors = append(errors, fmt.Sprintf("Resource '%s': Cost component '%s' missing", resourceName, exp.Name))
			continue
		}

		// Check fields
		// Unit
		if actUnit, ok := act["unit"].(string); !ok || actUnit != exp.Unit {
			errors = append(errors, fmt.Sprintf("Resource '%s' Component '%s': Unit mismatch. Expected '%s', got '%v'", resourceName, exp.Name, exp.Unit, act["unit"]))
		}
		// Monthly Quantity
		expQty := strings.ReplaceAll(exp.MonthlyQuantity, ",", "")
		if actQty, ok := act["monthly_quantity"].(string); !ok || actQty != expQty {
			errors = append(errors, fmt.Sprintf("Resource '%s' Component '%s': MonthlyQuantity mismatch. Expected '%s', got '%v'", resourceName, exp.Name, expQty, act["monthly_quantity"]))
		}
		// Monthly Cost
		expCostStr := strings.ReplaceAll(exp.MonthlyCost, ",", "")
		expCost, _ := strconv.ParseFloat(expCostStr, 64)

		actCostVal := act["monthly_cost"]
		var actCost float64
		switch v := actCostVal.(type) {
		case float64:
			actCost = v
		case string:
			actCost, _ = strconv.ParseFloat(v, 64)
		case int:
			actCost = float64(v)
		}

		if diff := actCost - expCost; diff < -0.0001 || diff > 0.0001 {
			errors = append(errors, fmt.Sprintf("Resource '%s' Component '%s': MonthlyCost mismatch. Expected %f, got %f", resourceName, exp.Name, expCost, actCost))
		}
	}
	return errors
}

func checkSubResources(parentName string, expected []ResourceCost, actualRaw interface{}) []string {
	var errors []string
	actualList, ok := actualRaw.([]interface{})
	if !ok {
		if len(expected) == 0 {
			return nil
		}
		return []string{fmt.Sprintf("Resource '%s': sub_resources is not a list", parentName)}
	}

	actualMap := make(map[string]map[string]interface{})
	for _, item := range actualList {
		m, ok := item.(map[string]interface{})
		if ok {
			name, _ := m["name"].(string)
			actualMap[name] = m
		}
	}

	for _, exp := range expected {
		act, ok := actualMap[exp.Name]
		if !ok {
			errors = append(errors, fmt.Sprintf("Resource '%s': Sub-resource '%s' missing", parentName, exp.Name))
			continue
		}

		fullName := parentName + "." + exp.Name
		if errs := checkCostComponents(fullName, exp.CostComponents, act["cost_components"]); len(errs) > 0 {
			errors = append(errors, errs...)
		}
		if errs := checkSubResources(fullName, exp.SubResources, act["sub_resources"]); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	return errors
}
