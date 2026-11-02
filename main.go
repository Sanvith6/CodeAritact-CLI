package main

import (
	"github.com/acme-corp/rs-cli/cmd"
)

//go:generate go run scripts/generate_config.go

func main() {
	cmd.Execute()
}
