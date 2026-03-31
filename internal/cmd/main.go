package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	app := &cli.App{
		Name:     "dwell",
		Usage:    "Manage your development environment",
		Version:  fmt.Sprintf("%s (built %s)", Version, BuildTime),
		Commands: []*cli.Command{
			syncCommand(),
			statusCommand(),
			doctorCommand(),
			initCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
