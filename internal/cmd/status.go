package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/zpershuai/dwell/internal/pkg/config"
	"github.com/zpershuai/dwell/internal/pkg/modules"
)

func getCurrentShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknown"
	}
	return filepath.Base(shell)
}

func isZshOnlyModule(moduleName string) bool {
	zshOnlyModules := []string{"zsh-syntax-highlighting", "zsh-navigation-tools"}
	for _, zshModule := range zshOnlyModules {
		if strings.EqualFold(moduleName, zshModule) {
			return true
		}
	}
	return false
}

func filterModulesByShell(mods []modules.Module, currentShell string) []modules.Module {
	currentShell = strings.ToLower(currentShell)

	if currentShell != "zsh" {
		filtered := make([]modules.Module, 0, len(mods))
		for _, mod := range mods {
			if !isZshOnlyModule(mod.Name()) {
				filtered = append(filtered, mod)
			}
		}
		return filtered
	}

	return mods
}

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

			appCfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			mods := appCfg.Registry.List()
			if len(mods) == 0 {
				color.Yellow("No modules configured")
				return nil
			}

			currentShell := getCurrentShell()
			filteredMods := filterModulesByShell(mods, currentShell)

			fmt.Printf("Current Shell: %s\n\n", color.CyanString(currentShell))

			// Setup table output
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "MODULE\tTYPE\tSTATUS\tREF\tMESSAGE")
			fmt.Fprintln(w, "------\t----\t------\t---\t-------")

			for _, mod := range filteredMods {
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

			fmt.Printf("\nTotal: %d modules", len(filteredMods))
			if currentShell != "zsh" {
				fmt.Printf(" (%d zsh-only modules hidden)", len(mods)-len(filteredMods))
			}
			fmt.Println()

			return nil
		},
	}
}
