package harness_codex

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

const codexManagedMarker = ".managed-skillset.json"

// loadMarketplace validates entries before a build while preserving unrelated metadata and policies.
func loadMarketplace(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]any{"name": "personal", "interface": map[string]any{"displayName": "Personal"}, "plugins": []any{}}, nil
	}
	if err != nil {
		return nil, err
	}
	var catalog map[string]any
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("read marketplace %s: %w", path, err)
	}
	if name, _ := catalog["name"].(string); strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("marketplace must have a nonempty name: %s", path)
	}
	plugins, ok := catalog["plugins"].([]any)
	if !ok {
		return nil, fmt.Errorf("marketplace plugins must be a list: %s", path)
	}
	if value, present := catalog["interface"]; present {
		iface, _ := value.(map[string]any)
		if name, _ := iface["displayName"].(string); strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("marketplace interface must have a displayName: %s", path)
		}
	}
	seen := map[string]bool{}
	for _, value := range plugins {
		entry, _ := value.(map[string]any)
		name, _ := entry["name"].(string)
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("marketplace plugin must have a nonempty name: %s", path)
		}
		if seen[name] {
			return nil, fmt.Errorf("duplicate marketplace plugin: %s", name)
		}
		seen[name] = true
		source, ok := entry["source"].(map[string]any)
		if text, isString := entry["source"].(string); isString {
			source, ok = map[string]any{"source": "local", "path": text}, true
		}
		if !ok {
			return nil, fmt.Errorf("invalid marketplace source: %s", name)
		}
		kind, _ := source["source"].(string)
		if kind != "local" && kind != "url" && kind != "git-subdir" && kind != "npm" {
			return nil, fmt.Errorf("invalid marketplace source: %s", name)
		}
		if kind == "local" || kind == "git-subdir" {
			relative, valid := source["path"].(string)
			if !valid || !strings.HasPrefix(relative, "./") || marketHasDotDot(relative) {
				return nil, fmt.Errorf("marketplace source path must stay within its root: %s", name)
			}
		}
		required := ""
		if kind == "url" || kind == "git-subdir" {
			required = "url"
		}
		if kind == "npm" {
			required = "package"
		}
		if value, _ := source[required].(string); required != "" && strings.TrimSpace(value) == "" {
			return nil, fmt.Errorf("marketplace source requires %s: %s", required, name)
		}
		if value, exists := entry["policy"]; exists {
			policy, ok := value.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid marketplace installation policy: %s", name)
			}
			installation, _ := policy["installation"].(string)
			if installation != "AVAILABLE" && installation != "INSTALLED_BY_DEFAULT" && installation != "NOT_AVAILABLE" {
				return nil, fmt.Errorf("invalid marketplace installation policy: %s", name)
			}
			if auth, _ := policy["authentication"].(string); strings.TrimSpace(auth) == "" {
				return nil, fmt.Errorf("invalid marketplace authentication policy: %s", name)
			}
		}
		if value, exists := entry["category"]; exists {
			if category, _ := value.(string); strings.TrimSpace(category) == "" {
				return nil, fmt.Errorf("invalid marketplace category: %s", name)
			}
		}
	}
	return catalog, nil
}

func marketHasDotDot(path string) bool {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == ".." {
			return true
		}
	}
	return false
}

var marketSkillName = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
var marketVersion = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:\+\S+)?$`)

func validateMarketplaceSkills(stage string, skills map[string]string) error {
	allowed := map[string]bool{"name": true, "description": true, "license": true, "allowed-tools": true, "metadata": true}
	names := make([]string, 0, len(skills))
	for name := range skills {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		path := filepath.Join(stage, "skills", name, "SKILL.md")
		if _, err := engine.ReadSkillName(path); err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lines := strings.Split(string(data), "\n")
		if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
			return fmt.Errorf("skill has no YAML frontmatter: %s", path)
		}
		end := 0
		for i := 1; i < len(lines); i++ {
			if s := strings.TrimSpace(lines[i]); s == "---" || s == "..." {
				end = i
				break
			}
		}
		if end == 0 {
			return fmt.Errorf("skill has unterminated YAML frontmatter: %s", path)
		}
		var metadata map[string]any
		if err := yaml.Unmarshal([]byte(strings.Join(lines[1:end], "\n")), &metadata); err != nil {
			return fmt.Errorf("invalid skill YAML: %s: %w", path, err)
		}
		if metadata == nil {
			return fmt.Errorf("skill frontmatter must be a mapping: %s", path)
		}
		var unexpected []string
		for key := range metadata {
			if !allowed[key] {
				unexpected = append(unexpected, key)
			}
		}
		if len(unexpected) > 0 {
			sort.Strings(unexpected)
			return fmt.Errorf("unsupported skill frontmatter keys in %s: %s", path, strings.Join(unexpected, ", "))
		}
		if metadata["name"] != name || len(name) > 64 || !marketSkillName.MatchString(name) {
			return fmt.Errorf("invalid packaged skill name: %s", path)
		}
		description, ok := metadata["description"].(string)
		if !ok || strings.TrimSpace(description) == "" || len([]rune(description)) > 1024 {
			return fmt.Errorf("skill description must be a nonempty string of at most 1024 characters: %s", path)
		}
	}
	return nil
}

// managedMarketplaceVersion checks ownership and computes the next patch version.
func managedMarketplaceVersion(destination string) (string, bool, error) {
	if _, err := os.Lstat(destination); err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}
	marker := filepath.Join(destination, codexManagedMarker)
	info, err := os.Lstat(marker)
	if err != nil || !info.Mode().IsRegular() {
		return "", false, fmt.Errorf("refusing to replace an unmanaged plugin: %s", destination)
	}
	data, err := os.ReadFile(marker)
	if err != nil {
		return "", false, err
	}
	var existing struct {
		Files map[string]string `json:"files"`
	}
	if err := json.Unmarshal(data, &existing); err != nil {
		return "", false, fmt.Errorf("invalid plugin marker: %w", err)
	}
	inventory, err := engine.FileHashes(destination, codexManagedMarker)
	if err != nil {
		return "", false, err
	}
	if !reflect.DeepEqual(existing.Files, inventory) {
		return "", false, fmt.Errorf("generated plugin has local edits: %s", destination)
	}
	data, err = os.ReadFile(filepath.Join(destination, ".codex-plugin/plugin.json"))
	if err != nil {
		return "", false, err
	}
	var manifest struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", false, err
	}
	match := marketVersion.FindStringSubmatch(manifest.Version)
	if match == nil {
		return "", false, fmt.Errorf("cannot increment marketplace plugin version: %s", manifest.Version)
	}
	patch, err := strconv.ParseUint(match[3], 10, 64)
	if err != nil || patch == ^uint64(0) {
		return "", false, fmt.Errorf("cannot increment marketplace plugin version: %s", manifest.Version)
	}
	return fmt.Sprintf("%s.%s.%d", match[1], match[2], patch+1), true, nil
}

type marketplaceSetup struct {
	marketplace, destination, version string
	catalog                           map[string]any
	plugins                           []any
	entry                             map[string]any
	matched                           bool
}

// inspectMarketplace validates all existing catalog and package ownership before any writes.
func inspectMarketplace(plan *engine.RenderPlan, marketplaceDir, pluginName string) (marketplaceSetup, error) {
	if len(plan.Skills) == 0 {
		return marketplaceSetup{}, fmt.Errorf("refusing to publish an empty skill collection")
	}
	marketplace := filepath.Join(marketplaceDir, "marketplace.json")
	destination := filepath.Join(marketplaceDir, pluginName)
	// The personal marketplace resolves local sources from the directory two
	// levels above its .agents/plugins catalog.
	discoveryRoot := filepath.Dir(filepath.Dir(marketplaceDir))
	relative, err := filepath.Rel(discoveryRoot, destination)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return marketplaceSetup{}, fmt.Errorf("plugin destination is outside marketplace root: %s", destination)
	}
	sourcePath := "./" + filepath.ToSlash(relative)
	for _, path := range []string{marketplace, destination} {
		if err := engine.RejectSymlinkAncestors(path); err != nil {
			return marketplaceSetup{}, err
		}
	}
	catalog, err := loadMarketplace(marketplace)
	if err != nil {
		return marketplaceSetup{}, err
	}
	plugins := catalog["plugins"].([]any)
	entry := map[string]any{"name": pluginName, "source": map[string]any{"source": "local", "path": sourcePath}, "policy": map[string]any{"installation": "AVAILABLE", "authentication": "ON_INSTALL"}, "category": "Productivity"}
	matched := false
	for _, value := range plugins {
		plugin := value.(map[string]any)
		if plugin["name"] != pluginName {
			continue
		}
		matched = true
		source := plugin["source"]
		path, stringSource := source.(string)
		localSource, objectSource := source.(map[string]any)
		if objectSource && localSource["source"] == "local" {
			path, _ = localSource["path"].(string)
		}
		if !stringSource && !objectSource {
			return marketplaceSetup{}, fmt.Errorf("existing %s entry points elsewhere", pluginName)
		}
		if path != sourcePath || (objectSource && !reflect.DeepEqual(source, entry["source"])) {
			return marketplaceSetup{}, fmt.Errorf("existing %s entry points elsewhere", pluginName)
		}
	}
	version := "0.1.0"
	next, present, err := managedMarketplaceVersion(destination)
	if err != nil {
		return marketplaceSetup{}, err
	}
	if present {
		version = next
	}
	return marketplaceSetup{marketplace, destination, version, catalog, plugins, entry, matched}, nil
}

func populateMarketplace(plan *engine.RenderPlan, marketplaceDir, pluginName string) (map[string]any, error) {
	setup, err := inspectMarketplace(plan, marketplaceDir, pluginName)
	if err != nil {
		return nil, err
	}
	marketplace, destination, version := setup.marketplace, setup.destination, setup.version
	catalog, plugins, entry, matched := setup.catalog, setup.plugins, setup.entry, setup.matched

	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return nil, err
	}
	// Stage and validate the complete package before replacing the installed one.
	temp, err := os.MkdirTemp(filepath.Dir(destination), ".skillset-")
	if err != nil {
		return nil, err
	}
	cleanupTemp := true
	defer func() {
		if cleanupTemp {
			_ = os.RemoveAll(temp)
		}
	}()
	stage := filepath.Join(temp, pluginName)
	if err := os.Mkdir(stage, 0755); err != nil {
		return nil, err
	}
	for relative := range plan.Files {
		target := filepath.Join(stage, relative)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return nil, err
		}
		if err := engine.CopyPlanFile(plan, relative, target); err != nil {
			return nil, err
		}
	}
	if err := os.MkdirAll(filepath.Join(stage, ".codex-plugin"), 0755); err != nil {
		return nil, err
	}
	manifest := map[string]any{"name": pluginName, "version": version, "description": "Managed skill collection", "skills": "./skills/", "interface": map[string]any{"displayName": pluginName, "shortDescription": "Managed skills", "longDescription": "Skills from agent-config, with shared reference files.", "developerName": "Personal", "category": "Productivity", "capabilities": []any{}, "defaultPrompt": "Use my managed skill collection."}}
	if err := engine.WriteJSON(filepath.Join(stage, ".codex-plugin/plugin.json"), manifest); err != nil {
		return nil, err
	}
	if err := validateMarketplaceSkills(stage, plan.Skills); err != nil {
		return nil, err
	}
	inventory, err := engine.FileHashes(stage, codexManagedMarker)
	if err != nil {
		return nil, err
	}
	if err := engine.WriteJSON(filepath.Join(stage, codexManagedMarker), map[string]any{"files": inventory}); err != nil {
		return nil, err
	}
	if !matched {
		catalog["plugins"] = append(plugins, entry)
	}
	if err := os.MkdirAll(filepath.Dir(marketplace), 0755); err != nil {
		return nil, err
	}
	// Keep the catalog candidate on its destination filesystem so its rename is atomic.
	catalogTemp, err := os.CreateTemp(filepath.Dir(marketplace), ".skillset-")
	if err != nil {
		return nil, err
	}
	catalogTempPath := catalogTemp.Name()
	defer os.Remove(catalogTempPath)
	data, err := json.MarshalIndent(catalog, "", "  ")
	if err == nil {
		_, err = catalogTemp.Write(append(data, '\n'))
	}
	if closeErr := catalogTemp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return nil, err
	}
	// Retain the old package until the catalog is replaced, and restore it on failure.
	// If restoration fails, leave the temporary directory in place for recovery.
	old := filepath.Join(temp, "previous")
	hadOld := false
	if _, err := os.Lstat(destination); err == nil {
		if err := os.Rename(destination, old); err != nil {
			return nil, err
		}
		hadOld = true
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if err := os.Rename(stage, destination); err != nil {
		if hadOld {
			if restoreErr := os.Rename(old, destination); restoreErr != nil {
				cleanupTemp = false
				return nil, fmt.Errorf("replace plugin: %v; restore previous plugin from %s: %w", err, old, restoreErr)
			}
		}
		return nil, err
	}
	if err := os.Rename(catalogTempPath, marketplace); err != nil {
		if removeErr := os.RemoveAll(destination); removeErr != nil {
			cleanupTemp = false
			return nil, fmt.Errorf("replace catalog: %v; remove new plugin: %w", err, removeErr)
		}
		if hadOld {
			if restoreErr := os.Rename(old, destination); restoreErr != nil {
				cleanupTemp = false
				return nil, fmt.Errorf("replace catalog: %v; restore previous plugin from %s: %w", err, old, restoreErr)
			}
		}
		return nil, err
	}
	return map[string]any{"plugin": destination, "plugin_name": pluginName, "marketplace": marketplace, "marketplace_name": catalog["name"], "skills": len(plan.Skills), "version": version}, nil
}
