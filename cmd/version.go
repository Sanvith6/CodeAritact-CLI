package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/acme-corp/rs-cli/internal/buildinfo"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of RS Platform CLI",
	Run: func(cmd *cobra.Command, args []string) {
		if outputJSON {
			info := map[string]string{
				"version":    buildinfo.Version,
				"commit":     buildinfo.Commit,
				"date":       buildinfo.Date,
				"go_version": buildinfo.GoVersion,
				"os":         runtime.GOOS,
				"arch":       runtime.GOARCH,
			}
			data, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Println("RS CLI")
		fmt.Println()
		fmt.Println("Version:")
		fmt.Println(buildinfo.Version)
		fmt.Println()
		fmt.Println("Organization:")
		fmt.Println(buildinfo.Organization)
		fmt.Println()
		fmt.Println("Build Number:")
		fmt.Println(buildinfo.BuildNumber)
		fmt.Println()
		fmt.Println("Build Date:")
		fmt.Println(buildinfo.Date)
		fmt.Println()
		fmt.Println("Commit:")
		fmt.Println(buildinfo.Commit)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
