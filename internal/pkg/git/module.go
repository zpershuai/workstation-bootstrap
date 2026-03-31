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

type Module struct {
	name        string
	url         string
	path        string
	ref         string
	links       []Link
	postSync    string
	description string
}

type Link struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type Config struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Path     string `yaml:"path"`
	Ref      string `yaml:"ref,omitempty"`
	PostSync string `yaml:"post_sync,omitempty"`
	Links    []Link `yaml:"links,omitempty"`
}

func NewModule(cfg Config) (*Module, error) {
	path, err := expandPath(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand path %q: %w", cfg.Path, err)
	}
	
	return &Module{
		name:        cfg.Name,
		url:         cfg.URL,
		path:        path,
		ref:         cfg.Ref,
		postSync:    cfg.PostSync,
		links:       cfg.Links,
		description: fmt.Sprintf("Git repository: %s", cfg.Name),
	}, nil
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Type() string {
	return "git"
}

func (m *Module) Description() string {
	return m.description
}

func (m *Module) Status(ctx context.Context) (*modules.State, error) {
	state := &modules.State{
		Name: m.name,
		Type: m.Type(),
		Ref:  m.ref,
	}

	if _, err := os.Stat(m.path); os.IsNotExist(err) {
		state.Status = modules.StatusMissing
		state.Message = "Repository not cloned"
		return state, nil
	}

	gitDir := filepath.Join(m.path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		state.Status = modules.StatusError
		state.Message = "Path exists but is not a git repository"
		return state, nil
	}

	statusCmd := exec.CommandContext(ctx, "git", "-C", m.path, "status", "--porcelain")
	output, _ := statusCmd.Output()
	if len(output) > 0 {
		state.Status = modules.StatusDirty
		state.Message = "Local changes detected"
		return state, nil
	}

	if m.ref != "" {
		fetchCmd := exec.CommandContext(ctx, "git", "-C", m.path, "fetch", "--all", "--tags")
		if err := fetchCmd.Run(); err != nil {
			state.Status = modules.StatusUnknown
			state.Message = fmt.Sprintf("Failed to fetch: %v", err)
			return state, nil
		}

		behindCmd := exec.CommandContext(ctx, "git", "-C", m.path, "rev-list", "HEAD..@{upstream}", "--count")
		behind, _ := behindCmd.Output()
		behindCount := strings.TrimSpace(string(behind))
		if behindCount != "" && behindCount != "0" {
			state.Status = modules.StatusBehind
			state.Message = fmt.Sprintf("%s commits behind remote", behindCount)
			return state, nil
		}

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

func (m *Module) Sync(ctx context.Context) error {
	if err := m.ensureParentDir(); err != nil {
		return err
	}

	if err := m.cloneOrUpdate(ctx); err != nil {
		return err
	}

	if err := m.fetchAndCheckout(ctx); err != nil {
		return err
	}

	if err := m.createSymlinks(); err != nil {
		return err
	}

	if err := m.runPostSync(ctx); err != nil {
		return err
	}

	return nil
}

func (m *Module) ensureParentDir() error {
	parent := filepath.Dir(m.path)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}
	return nil
}

func (m *Module) cloneOrUpdate(ctx context.Context) error {
	if _, err := os.Stat(filepath.Join(m.path, ".git")); os.IsNotExist(err) {
		cmd := exec.CommandContext(ctx, "git", "clone", m.url, m.path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}
	return nil
}

func (m *Module) fetchAndCheckout(ctx context.Context) error {
	fetchCmd := exec.CommandContext(ctx, "git", "-C", m.path, "fetch", "--all", "--tags")
	fetchCmd.Stderr = os.Stderr
	if err := fetchCmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch updates: %w", err)
	}

	if m.ref != "" {
		checkoutCmd := exec.CommandContext(ctx, "git", "-C", m.path, "checkout", m.ref)
		checkoutCmd.Stderr = os.Stderr
		if err := checkoutCmd.Run(); err != nil {
			return fmt.Errorf("failed to checkout %s: %w", m.ref, err)
		}

		showRefCmd := exec.CommandContext(ctx, "git", "-C", m.path, "show-ref", "--quiet", "refs/heads/"+m.ref)
		if err := showRefCmd.Run(); err == nil {
			pullCmd := exec.CommandContext(ctx, "git", "-C", m.path, "pull", "--ff-only")
			pullCmd.Stderr = os.Stderr
			if err := pullCmd.Run(); err != nil {
				return fmt.Errorf("failed to pull updates: %w", err)
			}
		}
	} else {
		pullCmd := exec.CommandContext(ctx, "git", "-C", m.path, "pull", "--ff-only")
		pullCmd.Stderr = os.Stderr
		if err := pullCmd.Run(); err != nil {
			return fmt.Errorf("failed to pull updates: %w", err)
		}
	}

	return nil
}

func (m *Module) createSymlinks() error {
	for _, link := range m.links {
		from, err := expandPath(link.From)
		if err != nil {
			return fmt.Errorf("failed to expand 'from' path %q: %w", link.From, err)
		}
		to, err := expandPath(link.To)
		if err != nil {
			return fmt.Errorf("failed to expand 'to' path %q: %w", link.To, err)
		}
		
		if err := os.Remove(to); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove existing file at %s: %w", to, err)
		}
		
		toParent := filepath.Dir(to)
		if err := os.MkdirAll(toParent, 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", to, err)
		}
		
		if err := os.Symlink(from, to); err != nil {
			return fmt.Errorf("failed to create symlink %s -> %s: %w", to, from, err)
		}
	}
	return nil
}

func (m *Module) runPostSync(ctx context.Context) error {
	if m.postSync == "" {
		return nil
	}
	
	postSyncPath, err := expandPath(m.postSync)
	if err != nil {
		return fmt.Errorf("failed to expand post-sync path %q: %w", m.postSync, err)
	}
	
	if _, err := os.Stat(postSyncPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to stat post-sync script %s: %w", postSyncPath, err)
	}
	
	cmd := exec.CommandContext(ctx, "bash", postSyncPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("post-sync script failed: %w", err)
	}
	
	return nil
}

func (m *Module) Check(ctx context.Context) []modules.CheckResult {
	results := []modules.CheckResult{}

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

	for _, link := range m.links {
		to, err := expandPath(link.To)
		if err != nil {
			results = append(results, modules.CheckResult{
				Name:    fmt.Sprintf("symlink-%s", filepath.Base(link.To)),
				Status:  modules.StatusError,
				Message: fmt.Sprintf("Failed to expand path: %v", err),
			})
			continue
		}
		
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
			from, err := expandPath(link.From)
			if err != nil {
				results = append(results, modules.CheckResult{
					Name:    fmt.Sprintf("symlink-%s", filepath.Base(to)),
					Status:  modules.StatusError,
					Message: fmt.Sprintf("Failed to expand source path: %v", err),
				})
				continue
			}
			
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
