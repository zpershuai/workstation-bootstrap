package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/zpershuai/dwell/internal/pkg/git"
)

// Loader handles configuration loading from multiple sources
type Loader struct {
	rootDir string
}

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

func NewLoader(rootDir string) *Loader {
	return &Loader{rootDir: rootDir}
}

// Load attempts to load configuration from dwell.yaml or falls back to repos.lock
func (l *Loader) Load() (*Config, error) {
	// Try dwell.yaml first
	yamlPath := filepath.Join(l.rootDir, "dwell.yaml")
	if _, err := os.Stat(yamlPath); err == nil {
		return l.loadYAML(yamlPath)
	}

	// Fall back to repos.lock
	lockPath := filepath.Join(l.rootDir, "repos", "repos.lock")
	if _, err := os.Stat(lockPath); err == nil {
		return l.loadReposLock(lockPath)
	}

	return nil, fmt.Errorf("no configuration found (tried dwell.yaml and repos/repos.lock)")
}

// Config is the unified configuration
type Config struct {
	Version string        `yaml:"version"`
	Git     []git.Config  `yaml:"git,omitempty"`
}

func (l *Loader) loadYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dwell.yaml: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse dwell.yaml: %w", err)
	}

	if cfg.Version == "" {
		cfg.Version = "1.0"
	}

	return &cfg, nil
}

// ReposLockEntry represents a single entry in repos.lock
type ReposLockEntry struct {
	Name string
	URL  string
	Dest string
	Ref  string
}

func (l *Loader) loadReposLock(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repos.lock: %w", err)
	}
	defer file.Close()

	cfg := &Config{
		Version: "1.0",
		Git:     []git.Config{},
	}

	// Module-specific link configurations
	linkConfigs := map[string][]git.Link{
		"nvim": {{From: "~/.dotfiles.d/repos/nvim", To: "~/.config/nvim"}},
		"tmux": {{From: "~/.dotfiles.d/repos/tmux", To: "~/.tmux"}},
		"claudecode_dotfiles": {{From: "~/.dotfiles.d/repos/claudecode_dotfiles", To: "~/.claude"}},
		"tpm": {{From: "~/.dotfiles.d/repos/tpm", To: "~/.tmux/plugins/tpm"}},
		"zsh-syntax-highlighting": {{From: "~/.dotfiles.d/repos/zsh-syntax-highlighting", To: "~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting"}},
		"zsh-navigation-tools": {{From: "~/.dotfiles.d/repos/zsh-navigation-tools", To: "~/.oh-my-zsh/custom/plugins/zsh-navigation-tools"}},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse: name url dest [ref]
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue // Skip invalid lines
		}

		entry := ReposLockEntry{
			Name: fields[0],
			URL:  fields[1],
			Dest: fields[2],
		}
		if len(fields) >= 4 {
			entry.Ref = fields[3]
		}

		gitCfg := git.Config{
			Name:  entry.Name,
			URL:   entry.URL,
			Path:  entry.Dest,
			Ref:   entry.Ref,
			Links: linkConfigs[entry.Name],
		}

		// Add post-sync for tmux
		if entry.Name == "tmux" {
			gitCfg.PostSync = "~/.tmux/install.sh"
		}

		cfg.Git = append(cfg.Git, gitCfg)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read repos.lock: %w", err)
	}

	return cfg, nil
}

// SaveYAML saves the configuration to dwell.yaml
func (l *Loader) SaveYAML(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write dwell.yaml: %w", err)
	}

	return nil
}
