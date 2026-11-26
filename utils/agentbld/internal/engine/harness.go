package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ModeParser func(*string) (any, error)

// Registry binds shared generation operations to the configured harness adapters.
type Registry map[string]Handler

// Handler describes a harness's mode, snapshot, render, and activation operations.
type Handler struct {
	Default  bool
	Summary  string
	ModeHelp string
	Examples string

	Parse   ModeParser
	Decode  func(json.RawMessage) (any, error)
	Stub    func(any, []string, string) (any, error)
	Build   func(Registry, any, []string, string) (*RenderPlan, any, error)
	Restore func(string, Snapshot) (*RenderPlan, error)
	Prepare func(Registry, string, string, Snapshot, *RenderPlan, Destination) (func() error, error)

	ResolveDestination func(any) (Destination, error)
	LoadProfile        func(map[string]json.RawMessage) error
	ConfigReport       func(any, Destination) (fields map[string]any, destinations map[string]string)
	Reset              func(Destination, string) error

	LinkComponents []string
}

// Destination contains the resolved paths for one harness.
type Destination struct {
	Config, Marketplace string
	UsesConfig          bool
	Profile             map[string]string
}

func (registry Registry) DefaultName() string {
	for _, name := range SortedKeys(registry) {
		if registry[name].Default {
			return name
		}
	}
	return ""
}

func (registry Registry) restoreHarnessPlan(gen, harness string, s Snapshot) (*RenderPlan, error) {
	h, ok := registry[harness]
	if !ok || h.Restore == nil {
		return nil, fmt.Errorf("unsupported harness in generation: %s", harness)
	}
	return h.Restore(filepath.Join(gen, harness), s)
}

func (registry Registry) RestoreHarnessPlans(gen string, s Snapshot) (map[string]*RenderPlan, error) {
	if err := validateGenerationLayout(gen, s); err != nil {
		return nil, err
	}
	plans := map[string]*RenderPlan{}
	for _, name := range SortedKeys(s) {
		p, err := registry.restoreHarnessPlan(gen, name, s)
		if err != nil {
			return nil, err
		}
		plans[name] = p
	}
	return plans, nil
}

func validateGenerationLayout(gen string, s Snapshot) error {
	entries, err := os.ReadDir(gen)
	if err != nil {
		return err
	}
	payloads := 0
	for _, entry := range entries {
		if entry.Name() == SnapshotFile && entry.Type().IsRegular() {
			continue
		}
		if _, ok := s[entry.Name()]; !ok || !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("unexpected generation entry: %s", filepath.Join(gen, entry.Name()))
		}
		payloads++
	}
	if payloads != len(s) {
		return fmt.Errorf("generation payloads do not match config: %s", gen)
	}
	return nil
}

// ApplyProfilePath fills a path option when the command line did not set it.
func ApplyProfilePath(fields map[string]json.RawMessage, key string, value *string) error {
	if data, ok := fields[key]; ok && *value == "" {
		if err := json.Unmarshal(data, value); err != nil {
			return fmt.Errorf("invalid %s: %w", key, err)
		}
	}
	return nil
}
