package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package management information",
	Long:  "RS CLI prepares your development environment, but does not manage packages directly.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("--------------------------------------------------")
		fmt.Println("RS CLI prepares your development environment.")
		fmt.Println("Package management is performed using uv.")
		fmt.Println()
		fmt.Println("Install packages using:")
		fmt.Println()
		fmt.Println("    uv add <package-name>==<version>")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println()
		fmt.Printf("    uv add company-web==1.2.3\n")
		fmt.Println()
		fmt.Println("Use the package name and version")
		fmt.Println("provided by your project or team.")
		fmt.Println("--------------------------------------------------")
		
		return nil
	},
}

func init() {
	RootCmd.AddCommand(packageCmd)
}
