package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/zpershuai/dwell/internal/pkg/config"
	"github.com/zpershuai/dwell/internal/pkg/git"
)

func initCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize dwell configuration",
		Description: `Initialize dwell by migrating from repos.lock to dwell.yaml.

This command will:
  1. Read the existing repos.lock file
  2. Generate a new dwell.yaml with the same configuration
  3. Preserve all module settings including links`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Overwrite existing dwell.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			rootDir := getRootDir()
			yamlPath := filepath.Join(rootDir, "dwell.yaml")

			// Check if dwell.yaml already exists
			if _, err := os.Stat(yamlPath); err == nil && !c.Bool("force") {
				return fmt.Errorf("dwell.yaml already exists. Use --force to overwrite")
			}

			// Load from repos.lock
			loader := config.NewLoader(rootDir)
			cfg, err := loader.Load()
			if err != nil {
				return fmt.Errorf("failed to load repos.lock: %w", err)
			}

			// Update to new format
			cfg.Version = "1.0"
			
			// Save as dwell.yaml
			if err := loader.SaveYAML(cfg, yamlPath); err != nil {
				return fmt.Errorf("failed to save dwell.yaml: %w", err)
			}

			color.Green("✓ Created dwell.yaml with %d git modules", len(cfg.Git))
			
			// Print sample config
			fmt.Println("\nGenerated configuration:")
			fmt.Println("---")
			for _, gitCfg := range cfg.Git {
				fmt.Printf("  - %s: %s\n", gitCfg.Name, gitCfg.URL)
				for _, link := range gitCfg.Links {
					fmt.Printf("      link: %s -> %s\n", link.From, link.To)
				}
			}
			
			fmt.Println("\nYou can now use 'dwell sync' to manage your modules.")
			
			return nil
		},
	}
}

// Helper to print git config details
func printGitConfig(cfg git.Config) {
	fmt.Printf("git:\n")
	fmt.Printf("  - name: %s\n", cfg.Name)
	fmt.Printf("    url: %s\n", cfg.URL)
	fmt.Printf("    path: %s\n", cfg.Path)
	if cfg.Ref != "" {
		fmt.Printf("    ref: %s\n", cfg.Ref)
	}
	if len(cfg.Links) > 0 {
		fmt.Printf("    links:\n")
		for _, link := range cfg.Links {
			fmt.Printf("      - from: %s\n", link.From)
			fmt.Printf("        to: %s\n", link.To)
		}
	}
}
