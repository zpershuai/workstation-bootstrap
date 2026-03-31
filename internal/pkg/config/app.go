package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zpershuai/dwell/internal/pkg/git"
	"github.com/zpershuai/dwell/internal/pkg/modules"
)

type AppConfig struct {
	Cfg      *Config
	Registry *modules.Registry
}

func LoadConfig() (*AppConfig, error) {
	rootDir, err := findConfigDir()
	if err != nil {
		return nil, err
	}

	loader := NewLoader(rootDir)
	cfg, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	registry := modules.NewRegistry()

	for _, gitCfg := range cfg.Git {
		module, err := git.NewModule(gitCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create module %s: %w", gitCfg.Name, err)
		}
		if err := registry.Register(module); err != nil {
			return nil, fmt.Errorf("failed to register module %s: %w", gitCfg.Name, err)
		}
	}

	return &AppConfig{
		Cfg:      cfg,
		Registry: registry,
	}, nil
}

func findConfigDir() (string, error) {
	if envDir := os.Getenv("DWELL_CONFIG"); envDir != "" {
		if hasConfig(envDir) {
			return envDir, nil
		}
	}

	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		xdgConfig = filepath.Join(home, ".config")
	}

	dwellConfig := filepath.Join(xdgConfig, "dwell")
	if hasConfig(dwellConfig) {
		return dwellConfig, nil
	}

	cwd, err := os.Getwd()
	if err == nil && hasConfig(cwd) {
		return cwd, nil
	}

	return "", fmt.Errorf("no configuration found in ~/.config/dwell/ or current directory (set DWELL_CONFIG to override)")
}

func hasConfig(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, "dwell.yaml")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(dir, "repos", "repos.lock")); err == nil {
		return true
	}
	return false
}
