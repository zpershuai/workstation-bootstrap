package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/zpershuai/dwell/internal/pkg/config"
)

func statusCommand() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Check the status of all modules",
		Description: `Display the current status of all configured modules.

The status column shows:
  ok      - Module is up to date
  dirty   - Module has local changes
  behind  - Module is behind remote
  ahead   - Module is ahead of remote
  missing - Module is not installed
  error   - Module has errors`,
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			
			appCfg, err := config.LoadConfig(config.GetRootDir())
			if err != nil {
				return err
			}

			mods := appCfg.Registry.List()
			if len(mods) == 0 {
				color.Yellow("No modules configured")
				return nil
			}

			// Setup table output
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "MODULE\tTYPE\tSTATUS\tREF\tMESSAGE")
			fmt.Fprintln(w, "------\t----\t------\t---\t-------")

			for _, mod := range mods {
				state, err := mod.Status(ctx)
				if err != nil {
					fmt.Fprintf(w, "%s\t%s\terror\t-\t%v\n", 
						mod.Name(), mod.Type(), err)
					continue
				}

				// Color-code the status
				statusStr := string(state.Status)
				switch state.Status {
				case "ok":
					statusStr = color.GreenString(string(state.Status))
				case "dirty", "behind", "ahead":
					statusStr = color.YellowString(string(state.Status))
				case "missing", "error":
					statusStr = color.RedString(string(state.Status))
				}

				ref := state.Ref
				if ref == "" {
					ref = "-"
				}

				message := state.Message
				if message == "" {
					message = "-"
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					mod.Name(), mod.Type(), statusStr, ref, message)
			}

			w.Flush()
			
			fmt.Printf("\nTotal: %d modules\n", len(mods))
			
			return nil
		},
	}
}
