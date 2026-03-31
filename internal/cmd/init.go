package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/zpershuai/dwell/internal/pkg/config"
)

func initCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize dwell configuration",
		Description: `Initialize dwell by migrating from repos.lock to dwell.yaml.

This command will:
  1. Read the existing repos.lock file from current directory
  2. Create ~/.config/dwell/ directory if needed
  3. Copy configuration to ~/.config/dwell/
  4. Preserve all module settings including links

You can also use --output to specify a custom config directory.`,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Overwrite existing dwell.yaml",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "Output directory for config (default: ~/.config/dwell/)",
			},
		},
		Action: func(c *cli.Context) error {
			var outputDir string
			if c.String("output") != "" {
				outputDir = c.String("output")
			} else {
				xdgConfig := os.Getenv("XDG_CONFIG_HOME")
				if xdgConfig == "" {
					home, err := os.UserHomeDir()
					if err != nil {
						return fmt.Errorf("failed to get home directory: %w", err)
					}
					xdgConfig = filepath.Join(home, ".config")
				}
				outputDir = filepath.Join(xdgConfig, "dwell")
			}

			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}

			yamlPath := filepath.Join(outputDir, "dwell.yaml")

			if _, err := os.Stat(yamlPath); err == nil && !c.Bool("force") {
				return fmt.Errorf("dwell.yaml already exists at %s. Use --force to overwrite", yamlPath)
			}

			loader := config.NewLoader(".")
			cfg, err := loader.Load()
			if err != nil {
				return fmt.Errorf("failed to load repos.lock from current directory: %w", err)
			}

			cfg.Version = "1.0"

			if err := loader.SaveYAML(cfg, yamlPath); err != nil {
				return fmt.Errorf("failed to save dwell.yaml: %w", err)
			}

			color.Green("✓ Created dwell.yaml at %s with %d git modules", outputDir, len(cfg.Git))

			fmt.Println("\nGenerated configuration:")
			fmt.Println("---")
			for _, gitCfg := range cfg.Git {
				fmt.Printf("  - %s: %s\n", gitCfg.Name, gitCfg.URL)
				for _, link := range gitCfg.Links {
					fmt.Printf("      link: %s -> %s\n", link.From, link.To)
				}
			}

			fmt.Printf("\nYou can now use 'dwell sync' from anywhere to manage your modules.\n")
			fmt.Printf("Config location: %s\n", outputDir)

			return nil
		},
	}
}
