package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/zpershuai/dwell/internal/pkg/config"
	"github.com/zpershuai/dwell/internal/pkg/git"
	"github.com/zpershuai/dwell/internal/pkg/modules"
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
		Compiled: time.Now(),
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

func getRootDir() string {
	execPath, err := os.Executable()
	if err != nil {
		execPath = os.Args[0]
	}
	execPath, _ = filepath.Abs(execPath)
	
	// If running from bin/, go up one level
	if filepath.Base(filepath.Dir(execPath)) == "bin" {
		return filepath.Dir(filepath.Dir(execPath))
	}
	
	// Otherwise assume we're in the repo root
	return filepath.Dir(execPath)
}

func loadConfig() (*config.Config, *modules.Registry, error) {
	rootDir := getRootDir()
	loader := config.NewLoader(rootDir)
	
	cfg, err := loader.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	registry := modules.NewRegistry()
	
	// Register git modules
	for _, gitCfg := range cfg.Git {
		module := git.NewModule(gitCfg)
		if err := registry.Register(module); err != nil {
			return nil, nil, fmt.Errorf("failed to register module %s: %w", gitCfg.Name, err)
		}
	}

	return cfg, registry, nil
}
