package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/acme-corp/rs-cli/internal/ui"
)

var packageCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"pkg"},
	Short:   "Manage dependencies (wraps uv)",
	Long:    `Interact with your project's dependencies via AWS CodeArtifact and uv.`,
}

var installCmd = &cobra.Command{
	Use:   "install [package]",
	Short: "Install a package",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg := args[0]
		version, _ := cmd.Flags().GetString("version")

		if !di.State.IsBootstrapped() {
			ui.PrintWarning("Your environment hasn't been bootstrapped yet.")
			if ready, _ := ui.Confirm("Would you like to run 'rs bootstrap' now?", true); ready {
				if err := bootstrapCmd.RunE(bootstrapCmd, []string{}); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("bootstrap required")
			}
		}

		err := di.PackageService.Install(cmd.Context(), pkg, version)
		if err != nil {
			ui.PrintError(fmt.Sprintf("Failed to install package: %v", err))
			return err
		}

		return nil
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync project dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !di.State.IsBootstrapped() {
			ui.PrintError("Please run 'rs bootstrap' first.")
			return fmt.Errorf("bootstrap required")
		}

		err := di.ProjectService.Verify(cmd.Context())
		if err != nil {
			ui.PrintError(fmt.Sprintf("Project verification failed: %v", err))
			return err
		}

		err = di.PackageService.Sync(cmd.Context())
		if err != nil {
			ui.PrintError(fmt.Sprintf("Failed to sync dependencies: %v", err))
			return err
		}

		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := di.ProjectService.Initialize(cmd.Context())
		if err != nil {
			ui.PrintError(fmt.Sprintf("Failed to initialize project: %v", err))
			return err
		}
		
		ui.PrintSuccess("Project initialized successfully.")
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [package]",
	Short: "Search for package versions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !di.State.IsBootstrapped() {
			ui.PrintError("Please run 'rs bootstrap' first.")
			return fmt.Errorf("bootstrap required")
		}

		pkg := args[0]
		ui.PrintInfo(fmt.Sprintf("Searching for versions of %s...", pkg))
		
		versions, err := di.RepoProvider.ListPackageVersions(cmd.Context(), pkg)
		if err != nil {
			ui.PrintError(fmt.Sprintf("Search failed: %v", err))
			return err
		}

		if len(versions) == 0 {
			ui.PrintInfo("No versions found.")
			return nil
		}

		var rows [][]string
		for _, v := range versions {
			rows = append(rows, []string{v.Version, v.Status})
		}
		
		ui.PrintTable([]string{"Version", "Status"}, rows)
		return nil
	},
}

func init() {
	installCmd.Flags().StringP("version", "v", "", "Specific version to install")
	
	packageCmd.AddCommand(installCmd)
	packageCmd.AddCommand(syncCmd)
	packageCmd.AddCommand(initCmd)
	packageCmd.AddCommand(searchCmd)
	
	RootCmd.AddCommand(packageCmd)
}
