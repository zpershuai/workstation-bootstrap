package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/zpershuai/dwell/internal/pkg/config"
	"github.com/zpershuai/dwell/internal/pkg/modules"
)

func syncCommand() *cli.Command {
	return &cli.Command{
		Name:  "sync",
		Usage: "Synchronize modules to desired state",
		Description: `Sync all modules or a specific module to their desired state.
		
Examples:
  dwell sync           # Sync all modules
  dwell sync nvim      # Sync only the nvim module
  dwell sync --dry-run # Preview changes without applying`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Preview changes without applying",
			},
		},
		Action: func(c *cli.Context) error {
			ctx := context.Background()

			appCfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			dryRun := c.Bool("dry-run")
			if dryRun {
				color.Yellow("DRY RUN MODE: No changes will be applied")
				fmt.Println()
			}

			var modulesToSync []modules.Module

			if c.NArg() == 0 {
				modulesToSync = appCfg.Registry.List()
				if len(modulesToSync) == 0 {
					color.Yellow("No modules configured")
					return nil
				}
				color.Blue("Syncing all %d modules...\n", len(modulesToSync))
			} else {
				moduleName := c.Args().First()
				module, found := appCfg.Registry.Get(moduleName)
				if !found {
					return fmt.Errorf("module %q not found", moduleName)
				}
				modulesToSync = []modules.Module{module}
				color.Blue("Syncing module: %s\n", moduleName)
			}

			successCount := 0
			errorCount := 0

			for _, mod := range modulesToSync {
				start := time.Now()

				if dryRun {
					color.Cyan("[DRY-RUN] Would sync: %s", mod.Name())
					continue
				}

				fmt.Printf("[%s] ", mod.Name())

				if err := mod.Sync(ctx); err != nil {
					color.Red("✗ %v", err)
					errorCount++
				} else {
					duration := time.Since(start).Round(time.Millisecond)
					color.Green("✓ synced (%s)", duration)
					successCount++
				}
			}

			fmt.Println()
			if errorCount == 0 {
				color.Green("✓ All modules synced successfully (%d/%d)", successCount, len(modulesToSync))
			} else {
				color.Yellow("⚠ Sync completed with errors (%d succeeded, %d failed)", successCount, errorCount)
			}

			return nil
		},
	}
}
