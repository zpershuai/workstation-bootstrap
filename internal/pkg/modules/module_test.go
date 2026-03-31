package modules

import (
	"context"
	"testing"
)

func TestStatusConstants(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusOK, "ok"},
		{StatusPending, "pending"},
		{StatusSyncing, "syncing"},
		{StatusError, "error"},
		{StatusDirty, "dirty"},
		{StatusBehind, "behind"},
		{StatusAhead, "ahead"},
		{StatusMissing, "missing"},
		{StatusUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.status)
			}
		})
	}
}

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	if r.modules == nil {
		t.Error("registry.modules is nil")
	}
}

func TestRegistryRegister(t *testing.T) {
	r := NewRegistry()
	
	mockModule := &mockModule{name: "test-module"}
	
	err := r.Register(mockModule)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}
	
	mod, ok := r.Get("test-module")
	if !ok {
		t.Error("Get() returned false for registered module")
	}
	if mod.Name() != "test-module" {
		t.Errorf("Get() returned wrong module: %s", mod.Name())
	}
}

func TestRegistryRegisterNil(t *testing.T) {
	r := NewRegistry()
	err := r.Register(nil)
	if err == nil {
		t.Error("Register(nil) should return error")
	}
}

func TestRegistryRegisterDuplicate(t *testing.T) {
	r := NewRegistry()
	
	mock1 := &mockModule{name: "duplicate"}
	mock2 := &mockModule{name: "duplicate"}
	
	r.Register(mock1)
	err := r.Register(mock2)
	
	if err == nil {
		t.Error("Register() should return error for duplicate name")
	}
}

func TestRegistryGetNotFound(t *testing.T) {
	r := NewRegistry()
	_, ok := r.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for non-existent module")
	}
}

func TestRegistryList(t *testing.T) {
	r := NewRegistry()
	
	r.Register(&mockModule{name: "module-a"})
	r.Register(&mockModule{name: "module-b"})
	
	list := r.List()
	if len(list) != 2 {
		t.Errorf("List() returned %d modules, expected 2", len(list))
	}
}

func TestRegistryListEmpty(t *testing.T) {
	r := NewRegistry()
	list := r.List()
	if len(list) != 0 {
		t.Errorf("List() returned %d modules, expected 0", len(list))
	}
}

type mockModule struct {
	name string
}

func (m *mockModule) Name() string { return m.name }
func (m *mockModule) Type() string { return "mock" }
func (m *mockModule) Description() string { return "mock module" }
func (m *mockModule) Status(ctx context.Context) (*State, error) {
	return &State{Name: m.name, Status: StatusOK}, nil
}
func (m *mockModule) Sync(ctx context.Context) error { return nil }
func (m *mockModule) Check(ctx context.Context) []CheckResult { return nil }
