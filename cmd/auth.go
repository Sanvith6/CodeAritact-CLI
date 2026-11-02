package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/acme-corp/rs-cli/internal/ui"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with AWS SSO",
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		spinner := di.UIManager.StartSpinner("Authenticating...")
		
		err := di.AuthService.Login(cmd.Context(), force)
		if err != nil {
			spinner.StopError("Authentication failed")
			return err
		}

		spinner.StopSuccess("Successfully authenticated")
		
		// Fetch CodeArtifact token and configure uv
		spinner = di.UIManager.StartSpinner("Configuring package manager...")
		endpoint, token, err := di.RepoService.GetAuthInfo(cmd.Context())
		if err != nil {
			spinner.StopError("Failed to get repository token")
			return err
		}
		
		err = di.PackageProvider.ConfigureAuth(cmd.Context(), "internal", endpoint, token)
		if err != nil {
			spinner.StopError("Failed to configure uv")
			return err
		}
		spinner.StopSuccess("Package manager configured")
		
		// Fetch identity for display
		identity, err := di.AuthService.GetIdentity(cmd.Context())
		if err == nil {
			ui.PrintInfo(fmt.Sprintf("Logged in as %s in %s", identity.RoleName, identity.AccountID))
		}

		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of AWS SSO",
	RunE: func(cmd *cobra.Command, args []string) error {
		spinner := di.UIManager.StartSpinner("Logging out...")
		
		err := di.AuthService.Logout(cmd.Context())
		if err != nil {
			spinner.StopError("Logout failed")
			return err
		}

		spinner.StopSuccess("Successfully logged out")
		return nil
	},
}

func init() {
	loginCmd.Flags().BoolP("force", "f", false, "Force a new login flow even if currently authenticated")
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(logoutCmd)
}
