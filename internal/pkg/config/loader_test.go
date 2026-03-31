package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zpershuai/dwell/internal/pkg/git"
)

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"~/test", filepath.Join(home, "test")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ExpandPath(tt.input)
			if result != tt.expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewLoader(t *testing.T) {
	loader := NewLoader("/tmp/test")
	if loader == nil {
		t.Fatal("NewLoader() returned nil")
	}
	if loader.rootDir != "/tmp/test" {
		t.Errorf("rootDir = %q, want %q", loader.rootDir, "/tmp/test")
	}
}

func TestLoadReposLock(t *testing.T) {
	tmpDir := t.TempDir()
	
	lockContent := `# name  url  dest  ref(optional)
nvim  git@github.com:user/nvim.git  ~/.dotfiles.d/repos/nvim  main
tmux  git@github.com:user/tmux.git  ~/.dotfiles.d/repos/tmux
`
	lockPath := filepath.Join(tmpDir, "repos.lock")
	os.WriteFile(lockPath, []byte(lockContent), 0644)
	
	loader := NewLoader(tmpDir)
	cfg, err := loader.loadReposLock(lockPath)
	
	if err != nil {
		t.Fatalf("loadReposLock() error = %v", err)
	}
	
	if len(cfg.Git) != 2 {
		t.Errorf("expected 2 git modules, got %d", len(cfg.Git))
	}
	
	if cfg.Git[0].Name != "nvim" {
		t.Errorf("first module name = %q, want nvim", cfg.Git[0].Name)
	}
	
	if cfg.Git[0].Ref != "main" {
		t.Errorf("first module ref = %q, want main", cfg.Git[0].Ref)
	}
	
	if cfg.Git[1].Ref != "" {
		t.Errorf("second module ref should be empty, got %q", cfg.Git[1].Ref)
	}
}

func TestLoadReposLockNotExist(t *testing.T) {
	loader := NewLoader("/nonexistent")
	_, err := loader.loadReposLock("/nonexistent/repos.lock")
	
	if err == nil {
		t.Error("loadReposLock() should return error for non-existent file")
	}
}

func TestLoadReposLockEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "repos.lock")
	os.WriteFile(lockPath, []byte(""), 0644)
	
	loader := NewLoader(tmpDir)
	cfg, err := loader.loadReposLock(lockPath)
	
	if err != nil {
		t.Fatalf("loadReposLock() error = %v", err)
	}
	
	if len(cfg.Git) != 0 {
		t.Errorf("expected 0 modules for empty file, got %d", len(cfg.Git))
	}
}

func TestLoadReposLockInvalidLine(t *testing.T) {
	tmpDir := t.TempDir()
	lockContent := `nvim  git@github.com:user/nvim.git
incomplete  url
`
	lockPath := filepath.Join(tmpDir, "repos.lock")
	os.WriteFile(lockPath, []byte(lockContent), 0644)
	
	loader := NewLoader(tmpDir)
	cfg, err := loader.loadReposLock(lockPath)
	
	if err != nil {
		t.Fatalf("loadReposLock() error = %v", err)
	}
	
	if len(cfg.Git) != 0 {
		t.Errorf("expected 0 valid modules, got %d", len(cfg.Git))
	}
}

func TestLoadYAML(t *testing.T) {
	tmpDir := t.TempDir()
	
	yamlContent := `version: "1.0"
git:
  - name: test-repo
    url: git@github.com:user/repo.git
    path: ~/test/repo
    ref: main
`
	yamlPath := filepath.Join(tmpDir, "dwell.yaml")
	os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	
	loader := NewLoader(tmpDir)
	cfg, err := loader.loadYAML(yamlPath)
	
	if err != nil {
		t.Fatalf("loadYAML() error = %v", err)
	}
	
	if cfg.Version != "1.0" {
		t.Errorf("version = %q, want 1.0", cfg.Version)
	}
	
	if len(cfg.Git) != 1 {
		t.Errorf("expected 1 git module, got %d", len(cfg.Git))
	}
	
	if cfg.Git[0].Name != "test-repo" {
		t.Errorf("module name = %q, want test-repo", cfg.Git[0].Name)
	}
}

func TestLoadYAMLEmptyVersion(t *testing.T) {
	tmpDir := t.TempDir()
	
	yamlContent := `git: []
`
	yamlPath := filepath.Join(tmpDir, "dwell.yaml")
	os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	
	loader := NewLoader(tmpDir)
	cfg, err := loader.loadYAML(yamlPath)
	
	if err != nil {
		t.Fatalf("loadYAML() error = %v", err)
	}
	
	if cfg.Version != "1.0" {
		t.Errorf("default version = %q, want 1.0", cfg.Version)
	}
}

func TestSaveYAML(t *testing.T) {
	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "output.yaml")
	
	loader := NewLoader(tmpDir)
	cfg := &Config{
		Version: "1.0",
		Git: []git.Config{
			{Name: "test", URL: "git@test.git", Path: "~/test"},
		},
	}
	
	err := loader.SaveYAML(cfg, yamlPath)
	if err != nil {
		t.Fatalf("SaveYAML() error = %v", err)
	}
	
	content, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("saved YAML file is empty")
	}
}
