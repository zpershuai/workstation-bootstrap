package modules

import (
	"context"
	"fmt"
	"time"
)

type Status string

const (
	StatusOK       Status = "ok"
	StatusPending  Status = "pending"
	StatusSyncing  Status = "syncing"
	StatusError    Status = "error"
	StatusDirty    Status = "dirty"
	StatusBehind   Status = "behind"
	StatusAhead    Status = "ahead"
	StatusMissing  Status = "missing"
	StatusUnknown  Status = "unknown"
)

type State struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    Status    `json:"status"`
	LastSync  time.Time `json:"last_sync,omitempty"`
	Version   string    `json:"version,omitempty"`
	Ref       string    `json:"ref,omitempty"`
	Message   string    `json:"message,omitempty"`
}

type Module interface {
	Name() string
	Type() string
	Description() string
	Status(ctx context.Context) (*State, error)
	Sync(ctx context.Context) error
	Check(ctx context.Context) []CheckResult
}

type CheckResult struct {
	Name    string
	Status  Status
	Message string
}

type Registry struct {
	modules map[string]Module
}

func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

func (r *Registry) Register(m Module) error {
	if m == nil {
		return fmt.Errorf("cannot register nil module")
	}
	name := m.Name()
	if _, exists := r.modules[name]; exists {
		return fmt.Errorf("module %q already registered", name)
	}
	r.modules[name] = m
	return nil
}

func (r *Registry) Get(name string) (Module, bool) {
	m, ok := r.modules[name]
	return m, ok
}

func (r *Registry) List() []Module {
	list := make([]Module, 0, len(r.modules))
	for _, m := range r.modules {
		list = append(list, m)
	}
	return list
}

type Config struct {
	Version string          `yaml:"version"`
	Modules []ModuleConfig  `yaml:"modules"`
}

type ModuleConfig struct {
	Name     string                 `yaml:"name"`
	Type     string                 `yaml:"type"`
	Config   map[string]interface{} `yaml:",inline"`
}
