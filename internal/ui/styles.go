package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Brand colors
	ColorBrand     = lipgloss.Color("#0070f3")
	ColorSuccess   = lipgloss.Color("#10b981")
	ColorWarning   = lipgloss.Color("#f59e0b")
	ColorError     = lipgloss.Color("#ef4444")
	ColorMuted     = lipgloss.Color("#6b7280")
	ColorHighlight = lipgloss.Color("#8b5cf6")

	// Styles
	StyleTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(ColorBrand).
			Bold(true).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)

	StyleStepTitle = lipgloss.NewStyle().
			Foreground(ColorBrand).
			Bold(true)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	StyleMuted = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleHighlight = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Bold(true)

	// Check/Cross icons
	IconCheck = StyleSuccess.Render("✓")
	IconCross = StyleError.Render("✗")
	IconWarn  = StyleWarning.Render("!")
	IconArrow = StyleBrand.Render("➜")

	StyleBrand = lipgloss.NewStyle().Foreground(ColorBrand)
)

// NoColor disables all color output if set to true.
func DisableColor() {
	// Not implemented with this lipgloss version for simplicity
}

// PrintTitle prints a styled title.
func PrintTitle(title string) {
	fmt.Fprintln(os.Stdout, StyleTitle.Render(title))
}

// PrintSuccess prints a success message.
func PrintSuccess(msg string) {
	fmt.Fprintf(os.Stdout, "%s %s\n", IconCheck, msg)
}

// PrintError prints an error message.
func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", IconCross, StyleError.Render(msg))
}

// PrintWarning prints a warning message.
func PrintWarning(msg string) {
	fmt.Fprintf(os.Stdout, "%s %s\n", IconWarn, StyleWarning.Render(msg))
}

// PrintInfo prints an informational message.
func PrintInfo(msg string) {
	fmt.Fprintf(os.Stdout, "%s %s\n", IconArrow, msg)
}

// PrintMuted prints a muted/subdued message.
func PrintMuted(msg string) {
	fmt.Fprintln(os.Stdout, StyleMuted.Render(msg))
}
