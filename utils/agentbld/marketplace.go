package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const marketplacePlugin = "skillset-uoleg"
const marketplaceMarker = ".managed-skillset.json"

func marketNonempty(value any) bool {
	s, ok := value.(string)
	return ok && strings.TrimSpace(s) != ""
}

func marketObject(value any) (map[string]any, bool) {
	m, ok := value.(map[string]any)
	return m, ok
}

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
	if !marketNonempty(catalog["name"]) {
		return nil, fmt.Errorf("marketplace must have a nonempty name: %s", path)
	}
	plugins, ok := catalog["plugins"].([]any)
	if !ok {
		return nil, fmt.Errorf("marketplace plugins must be a list: %s", path)
	}
	if value, present := catalog["interface"]; present {
		iface, ok := marketObject(value)
		if !ok || !marketNonempty(iface["displayName"]) {
			return nil, fmt.Errorf("marketplace interface must have a displayName: %s", path)
		}
	}
	seen := map[string]bool{}
	for _, value := range plugins {
		entry, ok := marketObject(value)
		if !ok || !marketNonempty(entry["name"]) {
			return nil, fmt.Errorf("marketplace plugin must have a nonempty name: %s", path)
		}
		name := entry["name"].(string)
		if seen[name] {
			return nil, fmt.Errorf("duplicate marketplace plugin: %s", name)
		}
		seen[name] = true
		source, ok := marketObject(entry["source"])
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
		if required != "" && !marketNonempty(source[required]) {
			return nil, fmt.Errorf("marketplace source requires %s: %s", required, name)
		}
		if value, exists := entry["policy"]; exists {
			policy, ok := marketObject(value)
			if !ok {
				return nil, fmt.Errorf("invalid marketplace installation policy: %s", name)
			}
			installation, _ := policy["installation"].(string)
			if installation != "AVAILABLE" && installation != "INSTALLED_BY_DEFAULT" && installation != "NOT_AVAILABLE" {
				return nil, fmt.Errorf("invalid marketplace installation policy: %s", name)
			}
			if !marketNonempty(policy["authentication"]) {
				return nil, fmt.Errorf("invalid marketplace authentication policy: %s", name)
			}
		}
		if category, exists := entry["category"]; exists && !marketNonempty(category) {
			return nil, fmt.Errorf("invalid marketplace category: %s", name)
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

// marketplaceInventory records package contents without following symlinks; the marker cannot hash itself.
func marketplaceInventory(root string) (map[string]string, error) {
	files := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlinks are not supported: %s", path)
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if filepath.ToSlash(rel) == marketplaceMarker {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hash := sha256.Sum256(data)
		files[filepath.ToSlash(rel)] = hex.EncodeToString(hash[:])
		return nil
	})
	return files, err
}

func marketplaceWriteJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
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
		if _, err := readSkillName(path); err != nil {
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

func marketplaceNoSymlinkAncestors(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	for {
		info, err := os.Lstat(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing a symlinked destination: %s", path)
		}
		parent := filepath.Dir(path)
		if parent == path {
			return nil
		}
		path = parent
	}
}

func populateMarketplace(plan *renderPlan, home string) (map[string]any, error) {
	if len(plan.skills) == 0 {
		return nil, fmt.Errorf("refusing to publish an empty skill collection")
	}
	marketplace := filepath.Join(home, ".agents/plugins/marketplace.json")
	destination := filepath.Join(home, "plugins", marketplacePlugin)
	for _, path := range []string{marketplace, destination} {
		if err := marketplaceNoSymlinkAncestors(path); err != nil {
			return nil, err
		}
	}
	catalog, err := loadMarketplace(marketplace)
	if err != nil {
		return nil, err
	}
	plugins := catalog["plugins"].([]any)
	entry := map[string]any{"name": marketplacePlugin, "source": map[string]any{"source": "local", "path": "./plugins/" + marketplacePlugin}, "policy": map[string]any{"installation": "AVAILABLE", "authentication": "ON_INSTALL"}, "category": "Productivity"}
	matched := false
	for _, value := range plugins {
		plugin := value.(map[string]any)
		if plugin["name"] != marketplacePlugin {
			continue
		}
		matched = true
		sourcePath, sourceIsString := plugin["source"].(string)
		if !reflect.DeepEqual(plugin["source"], entry["source"]) && !(sourceIsString && sourcePath == "./plugins/"+marketplacePlugin) {
			return nil, fmt.Errorf("existing %s entry points elsewhere", marketplacePlugin)
		}
	}
	version := "0.1.0"
	if _, err := os.Lstat(destination); err == nil {
		marker := filepath.Join(destination, marketplaceMarker)
		info, err := os.Lstat(marker)
		if err != nil || !info.Mode().IsRegular() {
			return nil, fmt.Errorf("refusing to replace an unmanaged plugin: %s", destination)
		}
		data, err := os.ReadFile(marker)
		if err != nil {
			return nil, err
		}
		// A matching marker establishes that replacing the generated package will not discard local edits.
		var existing struct {
			Files map[string]string `json:"files"`
		}
		if err := json.Unmarshal(data, &existing); err != nil {
			return nil, fmt.Errorf("invalid plugin marker: %w", err)
		}
		inventory, err := marketplaceInventory(destination)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(existing.Files, inventory) {
			return nil, fmt.Errorf("generated plugin has local edits: %s", destination)
		}
		data, err = os.ReadFile(filepath.Join(destination, ".codex-plugin/plugin.json"))
		if err != nil {
			return nil, err
		}
		var manifest struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, err
		}
		match := marketVersion.FindStringSubmatch(manifest.Version)
		if match == nil {
			return nil, fmt.Errorf("cannot increment marketplace plugin version: %s", manifest.Version)
		}
		patch, err := strconv.ParseUint(match[3], 10, 64)
		if err != nil || patch == ^uint64(0) {
			return nil, fmt.Errorf("cannot increment marketplace plugin version: %s", manifest.Version)
		}
		version = fmt.Sprintf("%s.%s.%d", match[1], match[2], patch+1)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
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
	stage := filepath.Join(temp, marketplacePlugin)
	if err := os.Mkdir(stage, 0755); err != nil {
		return nil, err
	}
	for relative := range plan.files {
		target := filepath.Join(stage, relative)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return nil, err
		}
		if err := copyPlanFile(plan, relative, target); err != nil {
			return nil, err
		}
	}
	manifest := map[string]any{"name": marketplacePlugin, "version": version, "description": "Oleg's complete managed skill collection", "author": map[string]any{"name": "uoleg"}, "skills": "./skills/", "interface": map[string]any{"displayName": "Skillset(uoleg)", "shortDescription": "Oleg's managed skills", "longDescription": "The effective personal and private skills from agent-config, with shared reference files.", "developerName": "uoleg", "category": "Productivity", "capabilities": []any{}, "defaultPrompt": "Use my managed skill collection."}}
	if err := marketplaceWriteJSON(filepath.Join(stage, ".codex-plugin/plugin.json"), manifest); err != nil {
		return nil, err
	}
	if err := validateMarketplaceSkills(stage, plan.skills); err != nil {
		return nil, err
	}
	inventory, err := marketplaceInventory(stage)
	if err != nil {
		return nil, err
	}
	if err := marketplaceWriteJSON(filepath.Join(stage, marketplaceMarker), map[string]any{"files": inventory}); err != nil {
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
	return map[string]any{"plugin": destination, "marketplace": marketplace, "marketplace_name": catalog["name"], "skills": len(plan.skills), "version": version}, nil
}
