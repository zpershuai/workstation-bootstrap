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

func LoadConfig(rootDir string) (*AppConfig, error) {
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

func GetRootDir() string {
	cwd, err := os.Getwd()
	if err == nil {
		if hasConfig(cwd) {
			return cwd
		}
	}

	if envDir := os.Getenv("DWELL_ROOT"); envDir != "" {
		if hasConfig(envDir) {
			return envDir
		}
	}

	execPath, err := os.Executable()
	if err != nil {
		execPath = os.Args[0]
	}
	execPath, _ = filepath.Abs(execPath)

	if filepath.Base(filepath.Dir(execPath)) == "bin" {
		execDir := filepath.Dir(filepath.Dir(execPath))
		if hasConfig(execDir) {
			return execDir
		}
	}

	execDir := filepath.Dir(execPath)
	if hasConfig(execDir) {
		return execDir
	}

	return cwd
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
