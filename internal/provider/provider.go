/*
Copyright (c) 2026 Plancost.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/plancost/terraform-provider-plancost/internal/apiclient"
	"github.com/plancost/terraform-provider-plancost/internal/prices"
)

// Ensure PlanCostProvider satisfies various provider interfaces.
var _ provider.Provider = &PlanCostProvider{}
var _ provider.ProviderWithValidateConfig = &PlanCostProvider{}
var _ provider.ProviderWithFunctions = &PlanCostProvider{}
var _ provider.ProviderWithEphemeralResources = &PlanCostProvider{}

// PlanCostProvider defines the provider implementation.
type PlanCostProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// PlanCostProviderModel describes the provider data model.
type PlanCostProviderModel struct {
	ApiEndpoint types.String `tfsdk:"api_endpoint"`
	ApiKey      types.String `tfsdk:"api_key"`
}

type PlanCostProviderData struct {
	PriceFetcher *prices.PriceFetcher
	Client       *apiclient.APIClient
}

func (p *PlanCostProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "plancost"
	resp.Version = p.version
}

func (p *PlanCostProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_endpoint": schema.StringAttribute{
				MarkdownDescription: "The API endpoint for the pricing service. Defaults to https://api.plancost.io. Can also be set via the PLANCOST_API_ENDPOINT environment variable.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for the pricing service. Can also be set via the PLANCOST_API_KEY environment variable. This is required when running in CI/CD environments.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *PlanCostProvider) ValidateConfig(ctx context.Context, request provider.ValidateConfigRequest, response *provider.ValidateConfigResponse) {
	var config PlanCostProviderModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &config)...); response.Diagnostics.HasError() {
		return
	}

	apiKey := ""
	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	} else if v := os.Getenv("PLANCOST_API_KEY"); v != "" {
		apiKey = v
	}

	// Check if running in CI/CD environment
	isCI := os.Getenv("CI") != "" || // Generic CI
		os.Getenv("GITHUB_ACTIONS") != "" || // GitHub Actions
		os.Getenv("GITLAB_CI") != "" || // GitLab CI
		os.Getenv("TF_BUILD") != "" || // Azure DevOps
		os.Getenv("TFC_RUN_ID") != "" || // Terraform Cloud / HCP Terraform
		os.Getenv("JENKINS_URL") != "" || // Jenkins
		os.Getenv("CIRCLECI") != "" || // CircleCI
		os.Getenv("TRAVIS") != "" || // Travis CI
		os.Getenv("BITBUCKET_COMMIT") != "" || // Bitbucket Pipelines
		os.Getenv("TF_IN_AUTOMATION") != "" // Terraform in Automation

	if isCI && apiKey == "" {
		response.Diagnostics.AddError(
			"Missing API Key in CI/CD Environment",
			"The provider cannot be used in CI/CD pipelines without an API key. Please sign up at https://plancost.io to obtain one.\n\n"+
				"Configuration options:\n"+
				"1. Set the 'PLANCOST_API_KEY' environment variable (Recommended).\n"+
				"2. Set the 'api_key' argument in the provider configuration.",
		)
		return
	}
}

func (p *PlanCostProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var config PlanCostProviderModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &config)...); response.Diagnostics.HasError() {
		return
	}

	endpoint := "https://api.plancost.io"
	if !config.ApiEndpoint.IsNull() {
		endpoint = config.ApiEndpoint.ValueString()
	} else if v := os.Getenv("PLANCOST_API_ENDPOINT"); v != "" {
		endpoint = v
	}

	apiKey := ""
	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	} else if v := os.Getenv("PLANCOST_API_KEY"); v != "" {
		apiKey = v
	}

	// Example client configuration for data sources and resources
	var client *apiclient.APIClient
	if apiKey != "" {
		client = apiclient.NewAPIClient("https://plancost.io/api", apiKey)
	}

	providerData := &PlanCostProviderData{
		PriceFetcher: prices.NewPriceFetcher(endpoint, ""),
		Client:       client,
	}
	response.DataSourceData = providerData
	response.ResourceData = providerData
}

func (p *PlanCostProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEstimateResource,
	}
}

func (p *PlanCostProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *PlanCostProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *PlanCostProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PlanCostProvider{
			version: version,
		}
	}
}
