package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zpershuai/dwell/internal/pkg/modules"
)

// Module manages external git repositories
type Module struct {
	name        string
	url         string
	path        string
	ref         string
	links       []Link
	postSync    string
	description string
}

// Link represents a symlink to create after syncing
type Link struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

// Config represents the configuration for a git module
type Config struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Path     string `yaml:"path"`
	Ref      string `yaml:"ref,omitempty"`
	PostSync string `yaml:"post_sync,omitempty"`
	Links    []Link `yaml:"links,omitempty"`
}

// NewModule creates a new git module from config
func NewModule(cfg Config) *Module {
	// Expand ~ to home directory
	path := expandPath(cfg.Path)
	
	return &Module{
		name:        cfg.Name,
		url:         cfg.URL,
		path:        path,
		ref:         cfg.Ref,
		postSync:    cfg.PostSync,
		links:       cfg.Links,
		description: fmt.Sprintf("Git repository: %s", cfg.Name),
	}
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

// Name returns the module name
func (m *Module) Name() string {
	return m.name
}

// Type returns the module type
func (m *Module) Type() string {
	return "git"
}

// Description returns the module description
func (m *Module) Description() string {
	return m.description
}

// Status returns the current state of the git repository
func (m *Module) Status(ctx context.Context) (*modules.State, error) {
	state := &modules.State{
		Name: m.name,
		Type: m.Type(),
		Ref:  m.ref,
	}

	// Check if directory exists
	if _, err := os.Stat(m.path); os.IsNotExist(err) {
		state.Status = modules.StatusMissing
		state.Message = "Repository not cloned"
		return state, nil
	}

	// Check if it's a git repo
	gitDir := filepath.Join(m.path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		state.Status = modules.StatusError
		state.Message = "Path exists but is not a git repository"
		return state, nil
	}

	// Check for local changes
	statusCmd := exec.CommandContext(ctx, "git", "-C", m.path, "status", "--porcelain")
	output, _ := statusCmd.Output()
	if len(output) > 0 {
		state.Status = modules.StatusDirty
		state.Message = "Local changes detected"
		return state, nil
	}

	// Check sync status with remote
	if m.ref != "" {
		// Fetch latest info
		fetchCmd := exec.CommandContext(ctx, "git", "-C", m.path, "fetch", "--all", "--tags")
		fetchCmd.Run() // Ignore errors for status check

		// Check if we're behind
		behindCmd := exec.CommandContext(ctx, "git", "-C", m.path, "rev-list", "HEAD..@{upstream}", "--count")
		behind, _ := behindCmd.Output()
		behindCount := strings.TrimSpace(string(behind))
		if behindCount != "" && behindCount != "0" {
			state.Status = modules.StatusBehind
			state.Message = fmt.Sprintf("%s commits behind remote", behindCount)
			return state, nil
		}

		// Check if we're ahead
		aheadCmd := exec.CommandContext(ctx, "git", "-C", m.path, "rev-list", "@{upstream}..HEAD", "--count")
		ahead, _ := aheadCmd.Output()
		aheadCount := strings.TrimSpace(string(ahead))
		if aheadCount != "" && aheadCount != "0" {
			state.Status = modules.StatusAhead
			state.Message = fmt.Sprintf("%s commits ahead of remote", aheadCount)
			return state, nil
		}
	}

	state.Status = modules.StatusOK
	state.Message = "Up to date"
	return state, nil
}

// Sync synchronizes the git repository
func (m *Module) Sync(ctx context.Context) error {
	// Ensure parent directory exists
	parent := filepath.Dir(m.path)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Check if already cloned
	if _, err := os.Stat(filepath.Join(m.path, ".git")); os.IsNotExist(err) {
		// Clone the repository
		cmd := exec.CommandContext(ctx, "git", "clone", m.url, m.path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}

	// Fetch all updates
	fetchCmd := exec.CommandContext(ctx, "git", "-C", m.path, "fetch", "--all", "--tags")
	fetchCmd.Stderr = os.Stderr
	if err := fetchCmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch updates: %w", err)
	}

	// Checkout specific ref if provided
	if m.ref != "" {
		checkoutCmd := exec.CommandContext(ctx, "git", "-C", m.path, "checkout", m.ref)
		checkoutCmd.Stderr = os.Stderr
		if err := checkoutCmd.Run(); err != nil {
			return fmt.Errorf("failed to checkout %s: %w", m.ref, err)
		}

		// Pull if it's a branch
		showRefCmd := exec.CommandContext(ctx, "git", "-C", m.path, "show-ref", "--quiet", "refs/heads/"+m.ref)
		if err := showRefCmd.Run(); err == nil {
			pullCmd := exec.CommandContext(ctx, "git", "-C", m.path, "pull", "--ff-only")
			pullCmd.Stderr = os.Stderr
			if err := pullCmd.Run(); err != nil {
				return fmt.Errorf("failed to pull updates: %w", err)
			}
		}
	} else {
		// Pull default branch
		pullCmd := exec.CommandContext(ctx, "git", "-C", m.path, "pull", "--ff-only")
		pullCmd.Stderr = os.Stderr
		if err := pullCmd.Run(); err != nil {
			return fmt.Errorf("failed to pull updates: %w", err)
		}
	}

	// Create symlinks
	for _, link := range m.links {
		from := expandPath(link.From)
		to := expandPath(link.To)
		
		// Remove existing file/link
		os.Remove(to)
		
		// Ensure parent directory exists
		toParent := filepath.Dir(to)
		os.MkdirAll(toParent, 0755)
		
		// Create symlink
		if err := os.Symlink(from, to); err != nil {
			return fmt.Errorf("failed to create symlink %s -> %s: %w", to, from, err)
		}
	}

	// Run post-sync script if provided
	if m.postSync != "" {
		postSyncPath := expandPath(m.postSync)
		if _, err := os.Stat(postSyncPath); err == nil {
			cmd := exec.CommandContext(ctx, "bash", postSyncPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("post-sync script failed: %w", err)
			}
		}
	}

	return nil
}

// Check performs health checks on the git module
func (m *Module) Check(ctx context.Context) []modules.CheckResult {
	results := []modules.CheckResult{}

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		results = append(results, modules.CheckResult{
			Name:    "git-binary",
			Status:  modules.StatusError,
			Message: "Git is not installed or not in PATH",
		})
		return results
	}

	results = append(results, modules.CheckResult{
		Name:   "git-binary",
		Status: modules.StatusOK,
	})

	// Check repository access
	if m.url != "" {
		cmd := exec.CommandContext(ctx, "git", "ls-remote", "--heads", m.url)
		if err := cmd.Run(); err != nil {
			results = append(results, modules.CheckResult{
				Name:    "repo-access",
				Status:  modules.StatusError,
				Message: fmt.Sprintf("Cannot access repository %s", m.url),
			})
		} else {
			results = append(results, modules.CheckResult{
				Name:   "repo-access",
				Status: modules.StatusOK,
			})
		}
	}

	// Check symlinks
	for _, link := range m.links {
		to := expandPath(link.To)
		info, err := os.Lstat(to)
		if err != nil {
			results = append(results, modules.CheckResult{
				Name:    fmt.Sprintf("symlink-%s", filepath.Base(to)),
				Status:  modules.StatusMissing,
				Message: fmt.Sprintf("Symlink missing: %s", to),
			})
		} else if info.Mode()&os.ModeSymlink == 0 {
			results = append(results, modules.CheckResult{
				Name:    fmt.Sprintf("symlink-%s", filepath.Base(to)),
				Status:  modules.StatusError,
				Message: fmt.Sprintf("Path exists but is not a symlink: %s", to),
			})
		} else {
			target, _ := os.Readlink(to)
			from := expandPath(link.From)
			if target != from {
				results = append(results, modules.CheckResult{
					Name:    fmt.Sprintf("symlink-%s", filepath.Base(to)),
					Status:  modules.StatusError,
					Message: fmt.Sprintf("Symlink points to wrong target: %s -> %s (expected %s)", to, target, from),
				})
			} else {
				results = append(results, modules.CheckResult{
					Name:   fmt.Sprintf("symlink-%s", filepath.Base(to)),
					Status: modules.StatusOK,
				})
			}
		}
	}

	return results
}
