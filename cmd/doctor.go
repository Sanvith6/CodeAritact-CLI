package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/acme-corp/rs-cli/internal/ui"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check your environment for potential issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintTitle("Environment Diagnostics (rs doctor)")
		fmt.Println("Checking system requirements and configuration...")
		fmt.Println()

		results, err := di.Validator.CheckAll(cmd.Context())
		if err != nil {
			return err
		}

		var tableRows [][]string
		hasErrors := false

		for _, res := range results {
			status := ui.IconCheck
			msg := res.Version
			
			if !res.Found {
				status = ui.IconCross
				msg = "Not installed"
				hasErrors = true
			} else if !res.MeetsMin {
				status = ui.IconWarn
				msg = fmt.Sprintf("%s (need >= %s)", res.Version, res.Tool.MinVersion)
				hasErrors = true
			}

			tableRows = append(tableRows, []string{
				status,
				res.Tool.Name,
				msg,
			})
		}

		ui.PrintTable([]string{"Status", "Tool", "Version"}, tableRows)

		fmt.Println()
		if hasErrors {
			ui.PrintError("Some checks failed. Please review the errors above.")
			return fmt.Errorf("doctor checks failed")
		}

		ui.PrintSuccess("All checks passed. Your environment is ready.")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(doctorCmd)
}
