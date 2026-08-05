package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

// PrintBanner prints the CLI banner.
func PrintBanner(version string) {
	banner := `
 █▀█ █▀
 █▀▄ ▄█  Platform CLI
`
	
	styledBanner := lipgloss.NewStyle().
		Foreground(ColorBrand).
		Bold(true).
		Render(banner)

	versionStyle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		MarginLeft(2).
		Render(fmt.Sprintf("v%s", version))

	fmt.Fprintln(os.Stdout, styledBanner)
	fmt.Fprintln(os.Stdout, versionStyle)
	fmt.Fprintln(os.Stdout, "")
}
