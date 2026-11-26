package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

type profileCommands struct {
	List struct{} `cmd:"" help:"List saved sync profile names."`
	Show struct {
		Name string `arg:"" help:"Name of the profile to show."`
	} `cmd:"" help:"Print a saved sync profile as JSON."`
}

var profileNamePattern = regexp.MustCompile("^[A-Za-z0-9][A-Za-z0-9_-]*$")

func profileFile(name string) (string, error) {
	if !profileNamePattern.MatchString(name) {
		return "", fmt.Errorf("invalid profile name %q: use letters, digits, hyphens, or underscores, starting with a letter or digit", name)
	}
	dir, err := profileDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".json"), nil
}

func profileDirectory() (string, error) {
	config := os.Getenv("XDG_CONFIG_HOME")
	if !filepath.IsAbs(config) {
		home, err := engine.DefaultHome()
		if err != nil {
			return "", err
		}
		config = filepath.Join(home, ".config")
	}
	return filepath.Join(config, "agentbld"), nil
}

func listProfiles() error {
	dir, err := profileDirectory()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name, jsonFile := strings.CutSuffix(entry.Name(), ".json")
		if jsonFile && profileNamePattern.MatchString(name) && engine.IsRegular(filepath.Join(dir, entry.Name())) {
			fmt.Println(name)
		}
	}
	return nil
}

func showProfile(name string) error {
	path, err := profileFile(name)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read profile: %w", err)
	}
	fmt.Println(strings.TrimSuffix(string(data), "\n"))
	return nil
}

func loadSyncProfile(a *cli, registry engine.Registry, name string) error {
	path, err := profileFile(name)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read profile: %w", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return fmt.Errorf("invalid profile %s: %w", path, err)
	}
	if fields == nil {
		return fmt.Errorf("invalid profile %s: expected a JSON object", path)
	}
	if data, ok := fields["modes"]; ok && len(a.Mode) == 0 {
		if err := json.Unmarshal(data, &a.Mode); err != nil {
			return fmt.Errorf("invalid profile %s: invalid modes: %w", path, err)
		}
	}
	if err := engine.ApplyProfilePath(fields, "render_dir", &a.RenderDir); err != nil {
		return fmt.Errorf("invalid profile %s: %w", path, err)
	}
	for _, harness := range engine.SortedKeys(registry) {
		if load := registry[harness].LoadProfile; load != nil {
			if err := load(fields); err != nil {
				return fmt.Errorf("invalid profile %s: %w", path, err)
			}
		}
	}
	return nil
}

func resolvedSyncProfile(root string, selected engine.Snapshot, destinations map[string]engine.Destination) map[string]any {
	// A profile can still record absolute destinations if HOME is unavailable.
	home, _ := engine.DefaultHome()
	portablePath := func(path string) string {
		if home == "" || !filepath.IsAbs(path) {
			return path
		}
		rel, err := filepath.Rel(home, path)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return path
		}
		if rel == "." {
			return "~/"
		}
		return "~/" + filepath.ToSlash(rel)
	}
	profile := map[string]any{"render_dir": portablePath(root), "modes": engine.SortedKeys(selected)}
	for _, destination := range destinations {
		for key, path := range destination.Profile {
			profile[key] = portablePath(path)
		}
	}
	return profile
}

func saveSyncProfile(path string, profile map[string]any) error {
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	// Publish a complete JSON file even when replacing an existing profile.
	file, err := os.CreateTemp(dir, ".profile-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	_, err = file.Write(append(data, '\n'))
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}
