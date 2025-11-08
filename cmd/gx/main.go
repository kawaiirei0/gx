package main

import (
	"github.com/yourusername/gx/cmd/gx/cmd"
)

// Version information (set via ldflags during build)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Set version information in cmd package
	cmd.SetVersionInfo(Version, Commit, BuildDate)
	cmd.Execute()
}
