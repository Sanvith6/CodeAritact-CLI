package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/acme-corp/rs-cli/internal/state"
	"github.com/acme-corp/rs-cli/internal/ui"
	"github.com/acme-corp/rs-cli/internal/workflow"
	"github.com/acme-corp/rs-cli/internal/workflow/steps"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
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
			steps.NewValidateArtifactAccessStep(di.RepoService),
		}

		engine := workflow.NewEngine(workflowSteps, di.UIManager, di.Logger)
		ctx := workflow.NewStepContext(cmd.Context())

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
		ui.PrintSuccess("Environment successfully bootstrapped!")
		fmt.Println("You can now navigate to your project and run `rs package install`")
		
		return nil
	},
}

func init() {
	RootCmd.AddCommand(bootstrapCmd)
}
