package buildinfo

import "runtime"

var (
	// Version is the current version of the CLI. Injected via ldflags.
	Version = "dev"
	
	// Commit is the git commit hash. Injected via ldflags.
	Commit = "unknown"
	
	// Date is the build timestamp. Injected via ldflags.
	Date = "unknown"
	
	// Organization is the company name.
	Organization = "Acme Corp"
	
	// BuildNumber is the CI build number. Injected via ldflags.
	BuildNumber = "local"

	
	// GoVersion is the version of Go used to build the CLI.
	GoVersion = runtime.Version()
)
