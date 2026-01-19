package testcase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/plancost/terraform-provider-plancost/internal/provider"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"
)

var ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"plancost": providerserver.NewProtocol6WithError(&provider.PlanCostProvider{}),
}

// TestCase is a simplified test case for running terraform plan
type TestCase struct {
	Steps                    []TestStep
	ProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
	SkipInit                 bool // Skip terraform init (can be set via TF_ACC_SKIP_INIT env var)
}

// TestStep is a single step in the test case
type TestStep struct {
	ConfigDirectory  string
	Config           string
	ConfigPlanChecks resource.ConfigPlanChecks
	ExpectError      *regexp.Regexp
	Check            func(t *testing.T, workDir string)
}

// Test runs the test case
func Test(t *testing.T, c TestCase) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping acceptance test; TF_ACC not set")
	}
	reattachInfo := tfexec.ReattachInfo{}

	c.ProtoV6ProviderFactories = ProviderFactories

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	for name, factory := range c.ProtoV6ProviderFactories {
		//reattachCh := make(chan *plugin.ReattachConfig)
		//closeCh := make(chan struct{})

		provider, err := factory()
		if err != nil {
			t.Fatalf("unable to create provider %q from factory: %v", name, err)
		}

		providerAddress := fmt.Sprintf("registry.terraform.io/hashicorp/%s", name)
		opts := &plugin.ServeOpts{
			GRPCProviderV6Func: func() tfprotov6.ProviderServer {
				return provider
			},
			Logger: hclog.New(&hclog.LoggerOptions{
				Name:   "plugintest",
				Level:  hclog.Trace,
				Output: io.Discard,
			}),
			NoLogOutputOverride: true,
			UseTFLogSink:        t,
			ProviderAddr:        providerAddress,
		}

		config, closeCh, err := plugin.DebugServe(ctx, opts)
		if err != nil {
			t.Fatalf("unable to start debug serve for provider %q: %v", name, err)
		}

		rc := tfexec.ReattachConfig{
			Protocol:        config.Protocol,
			ProtocolVersion: config.ProtocolVersion,
			Pid:             config.Pid,
			Test:            config.Test,
			Addr: tfexec.ReattachConfigAddr{
				Network: config.Addr.Network,
				String:  config.Addr.String,
			},
		}

		// Add to reattach config map
		reattachInfo[name] = rc
		// Add default hashicorp namespace mapping
		reattachInfo["registry.terraform.io/hashicorp/"+name] = rc

		// Track that server stopped
		go func() {
			<-closeCh
		}()

	}

	var workDir string
	if c.Steps[0].ConfigDirectory != "" {
		workDir = c.Steps[0].ConfigDirectory
	} else {
		workDir = t.TempDir()
	}

	tf, err := tfexec.NewTerraform(workDir, "terraform")
	require.NoError(t, err)

	// Check if init should be skipped from env var
	skipInit := c.SkipInit || os.Getenv("TF_ACC_SKIP_INIT") != ""

	for i, step := range c.Steps {
		if step.ConfigDirectory != "" {
			workDir = step.ConfigDirectory
		} else {
			err := os.WriteFile(filepath.Join(workDir, "main.tf"), []byte(step.Config), 0644)
			require.NoError(t, err)
		}

		// Run terraform init unless skipped
		if !skipInit {
			err = tf.Init(context.TODO(), tfexec.Reattach(reattachInfo))
			require.NoError(t, err, "terraform init failed")
		} else {
			t.Log("Skipping terraform init (TF_ACC_SKIP_INIT is set)")
		}

		// Run terraform plan and capture the plan file
		planFile := filepath.Join(workDir, fmt.Sprintf("plan-%d.tfplan", i))
		hasChanges, err := tf.Plan(context.TODO(), tfexec.Reattach(reattachInfo), tfexec.Out(planFile))

		if step.ExpectError != nil {
			require.Error(t, err)
			require.Regexp(t, step.ExpectError, err.Error())
			continue
		}
		require.NoError(t, err, "terraform plan failed")

		// Show the plan in JSON format for checking
		planJSON, err := tf.ShowPlanFile(context.TODO(), planFile, tfexec.Reattach(reattachInfo))

		for _, resourceChange := range planJSON.ResourceChanges {
			if resourceChange.Type != "plancost_estimate" {
				continue
			}
			var afterValue interface{}
			if resourceChange.Change != nil && resourceChange.Change.After != nil {
				afterValue = resourceChange.Change.After
			}
			afterValueJson, _ := json.MarshalIndent(afterValue, "", "  ")
			t.Logf("plancost_estimate resource %s after value: %s", resourceChange.Address, string(afterValueJson))
		}

		require.NoError(t, err, "terraform show plan failed")

		// Run the config plan checks if provided
		if step.ConfigPlanChecks.PreApply != nil {
			for _, check := range step.ConfigPlanChecks.PreApply {
				req := plancheck.CheckPlanRequest{
					Plan: planJSON,
				}
				resp := &plancheck.CheckPlanResponse{}
				check.CheckPlan(context.TODO(), req, resp)
				if resp.Error != nil {
					t.Errorf("plan check failed: %v", resp.Error)
				}
			}
		}

		if step.Check != nil {
			step.Check(t, workDir)
		}

		t.Logf("Plan has changes: %v", hasChanges)
	}
}
