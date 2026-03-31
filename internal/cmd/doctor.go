package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/zpershuai/dwell/internal/pkg/config"
	"github.com/zpershuai/dwell/internal/pkg/modules"
)

func doctorCommand() *cli.Command {
	return &cli.Command{
		Name:  "doctor",
		Usage: "Run health checks on your environment",
		Description: `Perform comprehensive health checks on all modules and system requirements.

Checks include:
  - Required binaries are installed
  - Git repositories are accessible
  - Symlinks are correctly configured
  - Files and directories exist`,
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			
			appCfg, err := config.LoadConfig(config.GetRootDir())
			if err != nil {
				return err
			}

			color.Blue("Running health checks...\n")

			mods := appCfg.Registry.List()
			totalChecks := 0
			passedChecks := 0
			failedChecks := 0

			for _, mod := range mods {
				fmt.Printf("\n[%s] %s\n", mod.Type(), mod.Name())
				
				results := mod.Check(ctx)
				for _, result := range results {
					totalChecks++
					
					switch result.Status {
					case modules.StatusOK:
						color.Green("  ✓ %s", result.Name)
						passedChecks++
					case modules.StatusMissing:
						color.Yellow("  ⚠ %s: %s", result.Name, result.Message)
						failedChecks++
					case modules.StatusError:
						color.Red("  ✗ %s: %s", result.Name, result.Message)
						failedChecks++
					default:
						fmt.Printf("  ? %s: %s\n", result.Name, result.Message)
					}
				}
			}

			fmt.Println()
			fmt.Println("----------------------------------------")
			
			if failedChecks == 0 {
				color.Green("✓ All checks passed (%d/%d)", passedChecks, totalChecks)
				return nil
			} else {
				color.Yellow("⚠ Checks completed with issues (%d passed, %d failed)", 
					passedChecks, failedChecks)
				os.Exit(1)
			}
			
			return nil
		},
	}
}
