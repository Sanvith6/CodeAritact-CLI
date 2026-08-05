package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss/table"
)

// PrintTable displays tabular data using lipgloss.
func PrintTable(headers []string, rows [][]string) {
	t := table.New().
		BorderStyle(StyleBrand).
		Headers(headers...).
		Rows(rows...)

	fmt.Fprintln(os.Stdout, t.Render())
}
