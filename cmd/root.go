package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	
	"github.com/acme-corp/rs-cli/internal/config"
	"github.com/acme-corp/rs-cli/internal/container"
	"github.com/acme-corp/rs-cli/internal/logging"
	"github.com/acme-corp/rs-cli/internal/ui"
)

var (
	// di is the global dependency injection container populated in PersistentPreRunE
	di *container.Container
)

var (
	verbose       bool
	quiet         bool
	noColor       bool
	nonInteractive bool
	awsProfile    string
	awsRegion     string
	outputJSON    bool
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "rs",
	Short: "RS Platform CLI",
	Long: `RS Platform CLI (rs) is the official enterprise developer platform tool.
It handles developer onboarding, authentication, and package manager configuration.

Use 'rs bootstrap' to configure your machine for the first time.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return bootstrapCmd.RunE(cmd, args)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Suppress color if requested
		if noColor {
			ui.DisableColor()
		}

		// Initialize logger early
		logLevel := "error" // Default: hide background logs so they don't clash with UI
		if verbose {
			logLevel = "debug"
		} else if quiet {
			logLevel = "error"
		}
		logger := logging.Setup(logLevel, outputJSON)

		// Load config
		loader := config.NewLoader()
		
		// If profile flag is passed, we should set it in the environment
		if awsProfile != "" {
			os.Setenv("RS_AWS_PROFILE_NAME", awsProfile)
		}
		
		cfg, err := loader.Load()
		if err != nil {
			// Don't fail immediately, some commands (like doctor or bootstrap) might want to fix it.
			logger.Warn("Failed to load config completely", "error", err)
			cfg = config.DefaultConfig // Fallback
		}

		// Create DI Container
		c, err := container.Build(cfg, logger, nonInteractive)
		if err != nil {
			return fmt.Errorf("failed to initialize dependencies: %w", err)
		}

		di = c
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.ExecuteContext(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	RootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress all output except errors")
	RootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colour output")
	RootCmd.PersistentFlags().BoolVar(&nonInteractive, "non-interactive", false, "disable interactive prompts")
	RootCmd.PersistentFlags().StringVar(&awsProfile, "profile", "", "override AWS profile")
	RootCmd.PersistentFlags().StringVar(&awsRegion, "region", "", "override AWS region")
	RootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "output in JSON format")
}
