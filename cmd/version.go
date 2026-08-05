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

		fmt.Printf("RS Platform CLI (rs) v%s\n", buildinfo.Version)
		fmt.Printf("Commit:  %s\n", buildinfo.Commit)
		fmt.Printf("Date:    %s\n", buildinfo.Date)
		fmt.Printf("Go:      %s\n", buildinfo.GoVersion)
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
