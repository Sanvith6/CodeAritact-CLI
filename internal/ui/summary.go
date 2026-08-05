package ui

import (
	"fmt"
	"os"
)

// PrintSummary prints a final summary box at the end of a long workflow.
func PrintSummary(title string, items map[string]string) {
	PrintTitle(title)
	
	for key, val := range items {
		// Just a simple padding for alignment in this skeleton
		fmt.Fprintf(os.Stdout, "  %-20s : %s\n", StyleStepTitle.Render(key), val)
	}
	fmt.Println()
}
