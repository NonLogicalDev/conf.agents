package harness_claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

const claudeManagedMarker = ".managed-skillset.json"

type claudeMarketplacePlan struct {
	root, name string
	catalog    map[string]any
	plan       *engine.RenderPlan
}

func inspectClaudeMarketplace(plan *engine.RenderPlan, root, name string) (*claudeMarketplacePlan, error) {
	if len(plan.Skills) == 0 {
		return nil, fmt.Errorf("refusing to generate an empty Claude skill collection")
	}
	if !filepath.IsAbs(root) || filepath.Clean(root) != root {
		return nil, fmt.Errorf("invalid Claude marketplace directory: %s", root)
	}
	if err := engine.RejectSymlinkAncestors(root); err != nil {
		return nil, err
	}
	for _, path := range []string{root, filepath.Join(root, ".claude-plugin"), filepath.Join(root, "plugins")} {
		if err := engine.ValidDir(path, "Claude marketplace directory"); err != nil {
			return nil, err
		}
	}
	catalogPath := filepath.Join(root, ".claude-plugin", "marketplace.json")
	if engine.IsLink(catalogPath) {
		return nil, fmt.Errorf("Claude marketplace catalog is a symlink: %s", catalogPath)
	}
	catalog := map[string]any{"name": "agentbld", "owner": map[string]any{"name": "agentbld"}, "plugins": []any{}}
	data, err := os.ReadFile(catalogPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err == nil {
		catalog = nil
		if err := json.Unmarshal(data, &catalog); err != nil {
			return nil, fmt.Errorf("invalid Claude marketplace catalog: %w", err)
		}
		if name, _ := catalog["name"].(string); strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("Claude marketplace requires a name")
		}
		owner, _ := catalog["owner"].(map[string]any)
		if name, _ := owner["name"].(string); strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("Claude marketplace requires an owner name")
		}
	}
	entries, ok := catalog["plugins"].([]any)
	if !ok {
		return nil, fmt.Errorf("Claude marketplace plugins must be a list")
	}
	destination := filepath.Join(root, "plugins", name)
	if engine.IsLink(destination) {
		return nil, fmt.Errorf("Claude plugin destination is a symlink: %s", destination)
	}
	managed := false
	if engine.Exists(destination) {
		if !engine.IsDir(destination) {
			return nil, fmt.Errorf("Claude plugin destination is not a directory: %s", destination)
		}
		markerPath := filepath.Join(destination, claudeManagedMarker)
		if engine.IsLink(markerPath) {
			return nil, fmt.Errorf("Claude plugin marker is a symlink: %s", markerPath)
		}
		data, err := os.ReadFile(markerPath)
		if err != nil {
			return nil, fmt.Errorf("refusing unmanaged Claude plugin %s: %w", destination, err)
		}
		var saved struct {
			Name  string            `json:"name"`
			Files map[string]string `json:"files"`
		}
		if json.Unmarshal(data, &saved) != nil || saved.Name != name || saved.Files == nil {
			return nil, fmt.Errorf("invalid Claude plugin ownership marker: %s", destination)
		}
		actual, err := engine.FileHashes(destination, claudeManagedMarker)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(saved.Files, actual) {
			return nil, fmt.Errorf("refusing locally modified Claude plugin: %s", destination)
		}
		managed = true
	}
	found := false
	for _, entry := range entries {
		item, _ := entry.(map[string]any)
		entryName, _ := item["name"].(string)
		if strings.TrimSpace(entryName) == "" {
			return nil, fmt.Errorf("invalid Claude marketplace entry")
		}
		if entryName != name {
			continue
		}
		if found {
			return nil, fmt.Errorf("duplicate Claude marketplace entry: %s", name)
		}
		found = true
		if !managed || item["source"] != "./plugins/"+name {
			return nil, fmt.Errorf("refusing unmanaged Claude marketplace entry: %s", name)
		}
	}
	if managed && !found {
		return nil, fmt.Errorf("managed Claude plugin has no marketplace entry: %s", name)
	}
	return &claudeMarketplacePlan{root: root, name: name, catalog: catalog, plan: plan}, nil
}

func populateClaudeMarketplace(m *claudeMarketplacePlan) error {
	if _, err := inspectClaudeMarketplace(m.plan, m.root, m.name); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(m.root, "plugins"), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(m.root, ".claude-plugin"), 0755); err != nil {
		return err
	}
	localStage, err := os.MkdirTemp(filepath.Join(m.root, "plugins"), ".agentbld-")
	if err != nil {
		return err
	}
	cleanupStage := true
	defer func() {
		if cleanupStage {
			_ = os.RemoveAll(localStage)
		}
	}()
	plugin := filepath.Join(localStage, m.name)
	if err = os.MkdirAll(filepath.Join(plugin, ".claude-plugin"), 0755); err != nil {
		return err
	}
	if err = engine.WriteJSON(filepath.Join(plugin, ".claude-plugin", "plugin.json"), map[string]any{"name": m.name}); err != nil {
		return err
	}
	for _, dir := range engine.SortedKeys(m.plan.Dirs) {
		if dir == "skills" || strings.HasPrefix(dir, "skills"+string(filepath.Separator)) {
			if err = os.MkdirAll(filepath.Join(plugin, dir), 0755); err != nil {
				return err
			}
		}
	}
	for _, rel := range engine.SortedKeys(m.plan.Files) {
		if !strings.HasPrefix(rel, "skills"+string(filepath.Separator)) {
			continue
		}
		dst := filepath.Join(plugin, rel)
		if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		if err = engine.CopyPlanFile(m.plan, rel, dst); err != nil {
			return err
		}
	}
	inventory, err := engine.FileHashes(plugin, claudeManagedMarker)
	if err != nil {
		return err
	}
	if err = engine.WriteJSON(filepath.Join(plugin, claudeManagedMarker), map[string]any{"name": m.name, "files": inventory}); err != nil {
		return err
	}
	// Recheck ownership immediately before publication.
	current, err := inspectClaudeMarketplace(m.plan, m.root, m.name)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Join(m.root, "plugins"), 0755); err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Join(m.root, ".claude-plugin"), 0755); err != nil {
		return err
	}
	entries := current.catalog["plugins"].([]any)
	found := false
	for _, item := range entries {
		if item.(map[string]any)["name"] == m.name {
			found = true
		}
	}
	if !found {
		entries = append(entries, map[string]any{"name": m.name, "source": "./plugins/" + m.name})
		current.catalog["plugins"] = entries
	}
	// Prepare the new catalog before replacing the previous package.
	catalogTmp, err := os.CreateTemp(filepath.Join(m.root, ".claude-plugin"), ".marketplace-")
	if err != nil {
		return err
	}
	defer os.Remove(catalogTmp.Name())
	b, err := json.MarshalIndent(current.catalog, "", "  ")
	if err == nil {
		_, err = catalogTmp.Write(append(b, '\n'))
	}
	if closeErr := catalogTmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	candidate := plugin
	dest := filepath.Join(m.root, "plugins", m.name)
	backup := filepath.Join(localStage, "previous")
	hadOld := engine.Exists(dest)
	if hadOld {
		if err = os.Rename(dest, backup); err != nil {
			return err
		}
	}
	restore := func(cause error, removeNew bool) error {
		if removeNew {
			if e := os.RemoveAll(dest); e != nil {
				cleanupStage = false
				return fmt.Errorf("%w; cannot remove replacement: %v; previous plugin retained at %s", cause, e, backup)
			}
		}
		if hadOld {
			if e := os.Rename(backup, dest); e != nil {
				cleanupStage = false
				return fmt.Errorf("%w; cannot restore previous plugin: %v; previous plugin retained at %s", cause, e, backup)
			}
		}
		return cause
	}
	if err = os.Rename(candidate, dest); err != nil {
		return restore(err, false)
	}
	if err = os.Rename(catalogTmp.Name(), filepath.Join(m.root, ".claude-plugin", "marketplace.json")); err != nil {
		return restore(err, true)
	}
	fmt.Printf("Claude marketplace ready: %s (register: claude plugin marketplace add %q; install: claude plugin install %s@%s)\n", m.root, m.root, m.name, current.catalog["name"])
	return nil
}
