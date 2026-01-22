/*
Copyright (c) 2026 Plancost.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/imdario/mergo"
	"github.com/plancost/terraform-provider-plancost/internal/apiclient"
	"github.com/plancost/terraform-provider-plancost/internal/dynamic"
	"github.com/plancost/terraform-provider-plancost/internal/hclparser/hcl"
	"github.com/plancost/terraform-provider-plancost/internal/prices"
	"github.com/plancost/terraform-provider-plancost/internal/provider/myvalidator"
	tfschema "github.com/plancost/terraform-provider-plancost/internal/schema"
	"github.com/plancost/terraform-provider-plancost/internal/usage"
	"github.com/shopspring/decimal"
)

var _ resource.Resource = &EstimateResource{}
var _ resource.ResourceWithModifyPlan = &EstimateResource{}

func NewEstimateResource() resource.Resource {
	return &EstimateResource{}
}

// EstimateResource defines the resource implementation.
type EstimateResource struct {
	priceFetcher *prices.PriceFetcher
	client       *apiclient.APIClient
}

// EstimateResourceModel describes the resource data model.
type EstimateResourceModel struct {
	WorkingDirectory types.String `tfsdk:"working_directory"`
	ProjectName      types.String `tfsdk:"project_name"`

	UsageFile types.String  `tfsdk:"usage_file"`
	Usage     types.Dynamic `tfsdk:"usage"`
	VarFile   types.String  `tfsdk:"var_file"`

	Resources     types.Dynamic        `tfsdk:"resources"`
	MonthlyCost   types.Number         `tfsdk:"monthly_cost"`
	View          types.String         `tfsdk:"view"`
	Id            types.String         `tfsdk:"id"`
	Guardrail     []GuardrailModel     `tfsdk:"guardrail"`
	Discount      []DiscountModel      `tfsdk:"discount"`
	TaggingPolicy []TaggingPolicyModel `tfsdk:"tagging_policy"`

	RecommendationsEnabled types.Bool `tfsdk:"recommendations_enabled"`
	Recommendations        types.List `tfsdk:"recommendations"`

	ExportMarkdownFile types.String `tfsdk:"export_markdown_file"`
	ExportUsageFile    types.String `tfsdk:"export_usage_file"`
}

type TaggingPolicyModel struct {
	Key           types.String   `tfsdk:"key"`
	AllowedValues []types.String `tfsdk:"allowed_values"`
	Pattern       types.String   `tfsdk:"pattern"`
	ResourceTypes []types.String `tfsdk:"resource_types"`
	Action        types.String   `tfsdk:"action"`
}

type DiscountModel struct {
	Percentage   types.Number `tfsdk:"percentage"`
	ResourceType types.String `tfsdk:"resource_type"`
}

type GuardrailModel struct {
	Condition types.String `tfsdk:"condition"`
	Threshold types.Number `tfsdk:"threshold"`
	Action    types.String `tfsdk:"action"`
}

type CostResourceModel struct {
	Name           string               `json:"name"`
	CostComponents []CostComponentModel `json:"cost_components"`
	SubResources   []CostResourceModel  `json:"sub_resources"`
}

type CostComponentModel struct {
	Name            string  `json:"name"`
	MonthlyQuantity string  `json:"monthly_quantity"`
	Unit            string  `json:"unit"`
	MonthlyCost     float64 `json:"monthly_cost"`
}

func (r *EstimateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_estimate"
}

func (r *EstimateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The `plancost_estimate` resource estimates the cost of cloud resources within a Terraform module. It integrates cost estimation, policy enforcement, and optimization recommendations directly into your Terraform workflow.\n\n> **Note:** Currently, `plancost` only supports the **Azure** provider (`azurerm`). Support for **AWS** and **GCP** is planned.",

		Attributes: map[string]schema.Attribute{
			"working_directory": schema.StringAttribute{
				MarkdownDescription: "Absolute path to the Terraform module directory (e.g., `abspath(path.module)`)",
				Required:            true,
				WriteOnly:           true,
			},

			"project_name": schema.StringAttribute{
				MarkdownDescription: "The name of the project to create in PlanCost. If not specified, the directory name of `working_directory` will be used, e.g., `my-terraform-project`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"usage_file": schema.StringAttribute{
				MarkdownDescription: "Absolute path to the usage file (e.g., `abspath(\"${path.module}/usage.yml\")`)",
				Optional:            true,
				WriteOnly:           true,
			},

			"usage": schema.DynamicAttribute{
				MarkdownDescription: "Usage data for resources. More details can be found in the [Usage Guide](../guides/usage.md).",
				Optional:            true,
			},

			"var_file": schema.StringAttribute{
				MarkdownDescription: "Absolute path to the variables file (e.g., `abspath(\"${path.module}/variables.tfvars\")`). The provider automatically loads variables from the following sources:\n" +
					"  1. The file specified in `var_file`.\n" +
					"  2. The file path specified in the `PLANCOST_VAR_FILE` environment variable.\n" +
					"  3. `terraform.tfvars` in the `working_directory`.\n" +
					"  4. Environment variables starting with `TF_VAR_`.",
				Optional:  true,
				WriteOnly: true,
			},

			"resources": schema.DynamicAttribute{
				Computed: true,
			},

			"monthly_cost": schema.NumberAttribute{
				MarkdownDescription: "The estimated monthly cost (numeric value)",
				Computed:            true,
			},

			"view": schema.StringAttribute{
				MarkdownDescription: "The pretty printed output of the estimate",
				Computed:            true,
			},

			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"recommendations_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable optimization recommendations. Note: This is a paid feature.",
				Optional:            true,
			},

			"recommendations": schema.ListAttribute{
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"resource_address":   types.StringType,
						"description":        types.StringType,
						"type":               types.StringType,
						"term":               types.StringType,
						"savings_amount":     types.NumberType,
						"savings_percentage": types.NumberType,
					},
				},
				MarkdownDescription: "List of optimization recommendations.",
				Computed:            true,
			},

			"export_markdown_file": schema.StringAttribute{
				MarkdownDescription: "Absolute path to the output markdown file (e.g., `abspath(\"${path.module}/estimate.md\")`). If specified, the cost estimate report will be written to this file.",
				Optional:            true,
				WriteOnly:           true,
			},

			"export_usage_file": schema.StringAttribute{
				MarkdownDescription: "Absolute path to the output usage file (e.g., `abspath(\"${path.module}/usage.yml\")`). If specified, the provider will generate a usage file containing the usage schema for all resources in the module. This is useful for discovering available usage parameters and creating a baseline for customization.",
				Optional:            true,
				WriteOnly:           true,
			},
		},

		Blocks: map[string]schema.Block{
			"guardrail": schema.ListNestedBlock{
				MarkdownDescription: "List of guardrail policies to enforce cost limits. Note: This is a paid feature. Free tier users are limited to 1 guardrail and cannot use 'block' actions.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"condition": schema.StringAttribute{
							MarkdownDescription: "The condition to trigger the guardrail. Valid values: 'monthly_cost_increase_amount', 'monthly_cost_increase_percentage', 'monthly_cost_budget'.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"monthly_cost_increase_amount",
									"monthly_cost_increase_percentage",
									"monthly_cost_budget",
								),
							},
						},

						"threshold": schema.NumberAttribute{
							MarkdownDescription: "The numeric value for the condition (amount or percentage).",
							Required:            true,
						},

						"action": schema.StringAttribute{
							MarkdownDescription: "The action to take when the threshold is breached. Valid values: 'warning', 'block'.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("warning", "block"),
							},
						},
					},
				},
			},

			"discount": schema.ListNestedBlock{
				MarkdownDescription: "List of discounts to apply.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"percentage": schema.NumberAttribute{
							MarkdownDescription: "The discount percentage (0.0 to 1.0).",
							Required:            true,
							Validators: []validator.Number{
								myvalidator.NumberBetween(0.0, 1.0),
							},
						},

						"resource_type": schema.StringAttribute{
							MarkdownDescription: "The resource type to apply the discount to (e.g., 'azurerm_virtual_machine'). If not specified, applies to all resources.",
							Optional:            true,
						},
					},
				},
			},

			"tagging_policy": schema.ListNestedBlock{
				MarkdownDescription: "List of tagging policies to enforce. Note: This is a paid feature.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "The tag key that must exist.",
							Required:            true,
						},

						"allowed_values": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of allowed values for the tag. If specified, the tag value must be one of these.",
							Optional:            true,
						},

						"pattern": schema.StringAttribute{
							MarkdownDescription: "Regex pattern that the tag value must match.",
							Optional:            true,
							Validators: []validator.String{
								myvalidator.ValidRegex(),
							},
						},

						"resource_types": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of resource types to apply this policy to (e.g.  ['azurerm_resource_group']). If empty, applies to all resources that support tags.",
							Optional:            true,
						},

						"action": schema.StringAttribute{
							MarkdownDescription: "The action to take when the policy is violated. Valid values: 'warning', 'block'.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("warning", "block"),
							},
						},
					},
				},
			},
		},
	}
}

func (r *EstimateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*PlanCostProviderData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *PlanCostProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	if data.PriceFetcher == nil {
		resp.Diagnostics.AddError(
			"Price Fetcher Not Configured",
			"The price fetcher was not properly configured. Please report this issue to the provider developers.",
		)
		return
	}

	r.priceFetcher = data.PriceFetcher
	r.client = data.Client
}

func (r *EstimateResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var config *EstimateResourceModel
	var state *EstimateResourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	// destroy doesn't need to modify plan
	if config == nil {
		return
	}

	// Preserve ID from state if exists
	config.Id = types.StringUnknown()
	if state != nil {
		config.Id = state.Id
	}

	// Load usage data if usage file path is provided
	usageMap, err := expandUsageMap(config.UsageFile, config.Usage)
	if err != nil {
		resp.Diagnostics.AddError(
			"Usage Data Initialization Error",
			fmt.Sprintf("Failed to initialize usage data: %s", err.Error()),
		)
		return
	}

	// API Key provided
	paidTier := false
	quota, err := getQuota(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Project Quota Retrieval Error",
			fmt.Sprintf("Failed to retrieve project quota: %s", err.Error()),
		)
		return
	}
	if quota != nil {
		paidTier = quota.IsPaid
	}

	if isNewResource := state == nil; isNewResource {
		if quota != nil && quota.Remaining <= 0 {
			resp.Diagnostics.AddError(
				"Project Quota Exceeded",
				"You have reached the maximum number of projects allowed in your account. Please visit https://plancost.io to upgrade your plan or delete existing projects.",
			)
			return
		}
	}

	// Parse the module
	workingDir := config.WorkingDirectory.ValueString()
	options := expandVariableOptions(config.VarFile.ValueString(), workingDir)
	allParsedResources, coreResources, err := ParseModule(workingDir, usageMap, options...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Module Calculation Error",
			fmt.Sprintf("Failed to calculate module %s: %s", "", err.Error()),
		)
	}
	allCostResources := make([]*tfschema.Resource, 0)
	for _, res := range allParsedResources {
		if len(res.CostComponents) == 0 && len(res.SubResources) == 0 {
			continue
		}
		allCostResources = append(allCostResources, res)
	}

	// Populate prices for all cost components
	if err = r.priceFetcher.PopulatePrices(allCostResources); err != nil {
		resp.Diagnostics.AddError(
			"Pricing Data Population Error",
			fmt.Sprintf("Failed to populate pricing data: %s", err.Error()),
		)
		return
	}

	// Apply discounts
	if len(config.Discount) > 0 {
		for _, costResource := range allCostResources {
			ApplyDiscount(costResource, config.Discount)
		}
	}

	totalCost := 0.0
	for _, costResource := range allCostResources {
		// Calculate costs based on populated prices
		costResource.CalculateCosts()

		// Add resource costs to total
		if costResource.MonthlyCost != nil {
			totalCost += costResource.MonthlyCost.InexactFloat64()
		}
	}

	// Guardrail Logic
	previousCost := 0.0
	if state != nil && !state.MonthlyCost.IsNull() {
		previousCost, _ = state.MonthlyCost.ValueBigFloat().Float64()
	}
	resp.Diagnostics.Append(Guardrails(paidTier, config.Guardrail, totalCost, previousCost)...)

	// Tagging Policy Logic
	resp.Diagnostics.Append(TaggingPolicies(paidTier, config.TaggingPolicy, allParsedResources)...)

	// Optimization Recommendations
	recommendations := Optimization(paidTier, config.RecommendationsEnabled.ValueBool(), coreResources, allCostResources, r.priceFetcher)

	// Convert structured recommendations to object list for schema compatibility
	config.Recommendations = ConvertRecommendationsToAttrValue(recommendations)

	flattenedResources := flattenResources(allCostResources)
	if err != nil {
		resp.Diagnostics.AddError(
			"Resource Flattening Error",
			fmt.Sprintf("Failed to flatten resources: %s", err.Error()),
		)
		return
	}
	if v, err := dynamic.ToDynamic(flattenedResources); err != nil {
		resp.Diagnostics.AddError(
			"Resource Flattening Error",
			fmt.Sprintf("Failed to convert flattened resources to dynamic: %s", err.Error()),
		)
		return
	} else {
		config.Resources = v
	}

	// Write markdown file if export_markdown_file is set
	if !config.ExportMarkdownFile.IsNull() && config.ExportMarkdownFile.ValueString() != "" {
		priorResources := make([]CostResourceModel, 0)
		if state != nil {
			if err := dynamic.Unmarshal(state.Resources, &priorResources); err != nil {
				resp.Diagnostics.AddError("Failed to unmarshal prior resources", err.Error())
				return
			}
		}

		markdownContent := GenerateMarkdownOutput(priorResources, flattenedResources)
		err = os.WriteFile(config.ExportMarkdownFile.ValueString(), []byte(markdownContent), 0644)
		if err != nil {
			resp.Diagnostics.AddError("Failed to write markdown file", err.Error())
			return
		}
	}

	// Write usage file if export_usage_file is set
	if !config.ExportUsageFile.IsNull() && config.ExportUsageFile.ValueString() != "" {
		usageContent, err := GenerateUsageYAML(allParsedResources)
		if err != nil {
			resp.Diagnostics.AddError("Failed to generate usage file", err.Error())
			return
		}
		err = os.WriteFile(config.ExportUsageFile.ValueString(), usageContent, 0644)
		if err != nil {
			resp.Diagnostics.AddError("Failed to write usage file", err.Error())
			return
		}
	}

	// Set the modified plan
	config.WorkingDirectory = types.StringNull()
	config.UsageFile = types.StringNull()
	config.VarFile = types.StringNull()
	config.ExportMarkdownFile = types.StringNull()
	config.ExportUsageFile = types.StringNull()
	config.MonthlyCost = types.NumberValue(decimal.NewFromFloat(totalCost).Round(2).BigFloat())
	config.View = types.StringValue(GenerateConsoleOutput(config.ProjectName.ValueString(), allParsedResources, recommendations, paidTier))
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &config)...)
}

func (r *EstimateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config, plan EstimateResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(req.Config.Get(ctx, &config)...); resp.Diagnostics.HasError() {
		return
	}

	if r.client != nil {
		friendlyName := config.ProjectName.ValueString()
		if friendlyName == "" {
			workingDirectory := config.WorkingDirectory.ValueString()
			friendlyName = filepath.Base(workingDirectory)
		}
		if friendlyName == "" || friendlyName == "." || friendlyName == "/" {
			friendlyName = "PlanCost Estimate Project"
		}
		project, err := r.client.CreateProject(friendlyName)
		if err != nil {
			resp.Diagnostics.AddError("Failed to create project", err.Error())
			return
		}
		plan.Id = types.StringValue(project.ID)
	} else {
		plan.Id = types.StringValue("plancost-local")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EstimateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EstimateResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EstimateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EstimateResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	if r.client != nil {
		_, err := r.client.GetProject(data.Id.ValueString())
		if err != nil {
			if strings.Contains(err.Error(), "status 404") {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("Failed to get project", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EstimateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EstimateResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	if r.client != nil {
		err := r.client.DeleteProject(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Failed to delete project", err.Error())
			return
		}
	}
}

func getQuota(client *apiclient.APIClient) (*apiclient.ProjectQuota, error) {
	if client == nil {
		return nil, nil
	}
	quota, err := client.GetProjectQuota()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve project quota: %w", err)
	}
	return quota, nil
}

func expandUsageMap(usageFilePath types.String, usageDyn types.Dynamic) (tfschema.UsageMap, error) {
	combinedMap := make(map[string]interface{})

	if !usageFilePath.IsNull() && usageFilePath.ValueString() != "" {
		usageFile, err := usage.LoadUsageFile(usageFilePath.ValueString())
		if err != nil {
			return tfschema.UsageMap{}, fmt.Errorf("could not load usage file from path %s: %w", usageFilePath.ValueString(), err)
		}
		combinedMap = usageFile.ToMap()
	}

	if !usageDyn.IsNull() {
		var rawMap map[string]interface{}
		if err := dynamic.Unmarshal(usageDyn, &rawMap); err != nil {
			return tfschema.UsageMap{}, fmt.Errorf("failed to unmarshal usage: %w", err)
		}

		if err := mergo.Merge(&combinedMap, rawMap, mergo.WithOverride); err != nil {
			return tfschema.UsageMap{}, fmt.Errorf("failed to merge usage: %w", err)
		}
	}

	if len(combinedMap) == 0 {
		return tfschema.UsageMap{}, nil
	}

	return tfschema.NewUsageMapFromInterface(combinedMap), nil
}

func expandVariableOptions(varFile string, workingDirectory string) []hcl.Option {
	options := make([]hcl.Option, 0)
	tfVarsPaths := make([]string, 0)

	// 1. VarFile
	if varFile != "" {
		tfVarsPaths = append(tfVarsPaths, varFile)
	}

	// 2. variable file path provided by environment variable: PLANCOST_VAR_FILE
	if envVarPath := os.Getenv("PLANCOST_VAR_FILE"); envVarPath != "" {
		tfVarsPaths = append(tfVarsPaths, envVarPath)
	}

	// 3. terraform.tfvars in the WorkingDirectory
	tfvarsPath := filepath.Join(workingDirectory, "terraform.tfvars")
	if _, err := os.Stat(tfvarsPath); err == nil {
		tfVarsPaths = append(tfVarsPaths, tfvarsPath)
	}

	if len(tfVarsPaths) > 0 {
		options = append(options, hcl.OptionWithTFVarsPaths(tfVarsPaths, false))
	}

	// 4. TF_VAR_XXXX environment variables
	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "TF_VAR_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimPrefix(parts[0], "TF_VAR_")
				envVars[key] = parts[1]
			}
		}
	}
	if len(envVars) > 0 {
		options = append(options, hcl.OptionWithTFEnvVars(envVars))
	}
	return options
}

func flattenResources(resources []*tfschema.Resource) []CostResourceModel {
	// Sort resources by name to ensure deterministic order
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})

	costResources := make([]CostResourceModel, 0)
	for _, r := range resources {
		costComps := flattenCostComponents(r.CostComponents)
		subRes := flattenSubResources(r.SubResources)
		if len(costComps) == 0 && len(subRes) == 0 {
			continue
		}
		costResources = append(costResources, CostResourceModel{
			Name:           r.Name,
			CostComponents: costComps,
			SubResources:   subRes,
		})
	}
	return costResources
}

// Helper to convert cost components to a native list
func flattenCostComponents(components []*tfschema.CostComponent) []CostComponentModel {
	// Sort components by name to ensure deterministic order
	sort.Slice(components, func(i, j int) bool {
		return components[i].Name < components[j].Name
	})

	result := make([]CostComponentModel, 0)
	for _, c := range components {
		unit := c.Unit
		qty := "0"
		if q := c.UnitMultiplierMonthlyQuantity(); q != nil {
			qty = q.Round(4).String()
		} else if c.UsageBased {
			qty = "0"
			unit = "usage_based"
		}
		monthlyCost := 0.0
		if c.MonthlyCost != nil {
			monthlyCost = c.MonthlyCost.Round(2).InexactFloat64()
		} else {
			unit = "usage_based"
		}
		if monthlyCost == 0.0 && qty == "0" && unit != "usage_based" {
			continue
		}
		result = append(result, CostComponentModel{
			Name:            c.Name,
			MonthlyQuantity: qty,
			Unit:            unit,
			MonthlyCost:     monthlyCost,
		})
	}
	return result
}

// Helper to convert sub-resources to a native list
func flattenSubResources(subResources []*tfschema.Resource) []CostResourceModel {
	// Sort sub-resources by name to ensure deterministic order
	sort.Slice(subResources, func(i, j int) bool {
		return subResources[i].Name < subResources[j].Name
	})

	result := make([]CostResourceModel, 0)
	for _, r := range subResources {
		costComps := flattenCostComponents(r.CostComponents)
		subRes := flattenSubResources(r.SubResources)
		if len(costComps) == 0 && len(subRes) == 0 {
			continue
		}
		result = append(result, CostResourceModel{
			Name:           r.Name,
			CostComponents: costComps,
			SubResources:   subRes,
		})
	}
	return result
}
