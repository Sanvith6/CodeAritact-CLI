package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/acme-corp/rs-cli/internal/runner"
	"github.com/acme-corp/rs-cli/internal/state"
	"github.com/acme-corp/rs-cli/internal/ui"
	"github.com/acme-corp/rs-cli/internal/workflow"
	"github.com/acme-corp/rs-cli/internal/workflow/steps"
)

var bootstrapCmd = &cobra.Command{
	Use:   "auth",
	Short: "Configure your machine for development",
	RunE: func(cmd *cobra.Command, args []string) error {
		st, _ := di.State.Load()
		version := "dev"
		if st != nil {
			version = st.CLIVersion
		}
		ui.PrintBanner(version)
		ui.PrintTitle("Bootstrapping Developer Environment")
		fmt.Println("This process will ensure your machine is ready for development.")
		fmt.Println()

		// Build workflow
		workflowSteps := []workflow.Step{
			steps.NewCheckNetworkStep(di.Network, di.Config.AWS.Region),
			steps.NewValidateMachineStep(di.Validator),
			steps.NewCheckSSOProfileStep(di.AuthProvider),
			steps.NewEnsureSSOLoginStep(di.AuthService),
			steps.NewConfigureUVStep(di.RepoService, di.PackageProvider),
			steps.NewVerifyValidationPackageStep(di.RepoService),
			steps.NewSyncDependenciesStep(di.Runner),
		}

		engine := workflow.NewEngine(workflowSteps, di.UIManager, di.Logger)
		ctx := workflow.NewStepContext(cmd.Context())
		
		ctx.Ctx = context.WithValue(ctx.Ctx, "validation_package", di.Config.Package.ValidationPackage)

		if err := engine.Run(ctx); err != nil {
			fmt.Println()
			ui.PrintError(fmt.Sprintf("Bootstrap failed: %v", err))
			return err
		}

		// Update state
		err := di.State.Update(func(s *state.State) {
			s.BootstrapCompleted = true
			s.BootstrapTimestamp = time.Now()
		})
		if err != nil {
			di.Logger.Warn("Failed to save bootstrap state", "error", err)
		}

		fmt.Println()
		printSuccessScreen(cmd.Context(), di.Runner)
		
		return nil
	},
}

func printSuccessScreen(ctx context.Context, r runner.Runner) {
	fmt.Println("══════════════════════════════════════")
	fmt.Println("        RS Platform CLI")
	fmt.Println("══════════════════════════════════════")
	fmt.Println("✓ AWS SSO Login")
	fmt.Println("✓ AWS Identity Verified")
	fmt.Println("✓ CodeArtifact Connected")
	fmt.Println("✓ uv Configured")
	fmt.Println("✓ Repository Verified")
	fmt.Println("✓ Dependencies Synced")
	fmt.Println("══════════════════════════════════════")
	fmt.Println("Developer environment is ready.")
	fmt.Println()

	// Fetch installed packages
	out, err := r.RunSilent(ctx, "uv", "pip", "list")
	if err == nil && out.Stdout != "" {
		fmt.Println("Installed Packages:")
		fmt.Println(out.Stdout)
	}

	fmt.Println("If you want any other packages then add them to requirements.txt or pyproject.toml and re-run me again.")
	fmt.Println("══════════════════════════════════════")
}

func init() {
	RootCmd.AddCommand(bootstrapCmd)
}
