package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// openCodePath resolves the destination without following its symlinks.
func openCodePath(path string) (string, error) {
	if path == "" {
		path = os.Getenv("OPENCODE_CONFIG_DIR")
		if path == "" {
			base := os.Getenv("XDG_CONFIG_HOME")
			if base == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return "", err
				}
				base = filepath.Join(home, ".config")
			}
			path = filepath.Join(base, "opencode")
		}
	}
	return absolute(path)
}

// discoverOpenCodePlugins selects flat native modules, with later layers winning.
func discoverOpenCodePlugins(layers []string) (map[string]string, error) {
	plugins := map[string]string{}
	for _, layer := range layers {
		root := filepath.Join(layer, "opencode")
		dir := filepath.Join(root, "plugins")
		if isLink(root) || isLink(dir) {
			return nil, fmt.Errorf("symlinked OpenCode plugin roots are not supported: %s", dir)
		}
		if _, err := os.Lstat(dir); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if !isDir(dir) {
			return nil, fmt.Errorf("OpenCode plugin root is not a directory: %s", dir)
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			path := filepath.Join(dir, entry.Name())
			ext := filepath.Ext(path)
			if entry.Type()&os.ModeSymlink != 0 || !entry.Type().IsRegular() || (ext != ".js" && ext != ".ts") {
				return nil, fmt.Errorf("OpenCode plugins must be regular .js or .ts files: %s", path)
			}
			plugins[entry.Name()] = path
		}
	}
	return plugins, nil
}

// buildOpenCodePlan adds native plugins to the shared, destination-independent plan.
func buildOpenCodePlan(sections []string, includePlugins bool) (*renderPlan, error) {
	plan, err := buildPlan(sections, true)
	if err != nil {
		return nil, err
	}
	if includePlugins {
		plugins, err := discoverOpenCodePlugins(plan.layers)
		if err != nil {
			return nil, err
		}
		plan.dirs["plugins"] = true
		for name, source := range plugins {
			plan.files[filepath.Join("plugins", name)] = source
		}
	}
	plan.agents = string(resolveSkillReferences([]byte(plan.agents), "local", ""))
	if err := bundleResources(plan, plan.agents); err != nil {
		return nil, err
	}
	return renderSkillReferences(plan, "opencode", "", false)
}

// openCodeMode selects instructions/resources, skills and native plugins.
type openCodeMode struct{ agents, skills, plugins bool }

type openCodeSnapshot struct {
	Version       int               `json:"version"`
	Agents        bool              `json:"agents"`
	Skills        bool              `json:"skills"`
	Plugins       bool              `json:"plugins"`
	SectionsDirs  []string          `json:"sections_dirs"`
	RenderDir     string            `json:"render_dir"`
	LayerCount    int               `json:"layer_count"`
	FragmentCount int               `json:"fragment_count"`
	FileModes     map[string]uint32 `json:"file_modes"`
}

func parseOpenCodeMode(raw *string, _ syncFlags) (any, error) {
	mode := openCodeMode{agents: true, skills: true, plugins: true}
	if raw == nil || *raw == "opencode" {
		return mode, nil
	}
	name, data, ok := strings.Cut(*raw, ":")
	if !ok || name != "opencode" {
		return nil, fmt.Errorf("OpenCode mode must be opencode or opencode:<JSON object>")
	}
	decoder := json.NewDecoder(strings.NewReader(data))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return nil, fmt.Errorf("OpenCode mode JSON must be an object")
	}
	seen := map[string]bool{}
	for decoder.More() {
		token, err = decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("invalid OpenCode mode JSON: %w", err)
		}
		key, ok := token.(string)
		if !ok {
			return nil, fmt.Errorf("invalid OpenCode mode field")
		}
		if seen[key] {
			return nil, fmt.Errorf("duplicate OpenCode mode field %q", key)
		}
		seen[key] = true
		if key != "agents" && key != "skills" && key != "plugins" {
			return nil, fmt.Errorf("unknown OpenCode mode field %q", key)
		}
		token, err = decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("invalid OpenCode mode JSON: %w", err)
		}
		value, ok := token.(bool)
		if !ok {
			return nil, fmt.Errorf("OpenCode %s must be a boolean", key)
		}
		switch key {
		case "agents":
			mode.agents = value
		case "skills":
			mode.skills = value
		case "plugins":
			mode.plugins = value
		}
	}
	if _, err = decoder.Token(); err != nil {
		return nil, fmt.Errorf("invalid OpenCode mode JSON: %w", err)
	}
	if _, err = decoder.Token(); err != io.EOF {
		return nil, fmt.Errorf("trailing or invalid content after OpenCode mode JSON")
	}
	if !mode.agents && !mode.skills && !mode.plugins {
		return nil, fmt.Errorf("OpenCode mode must select agents, skills, or plugins")
	}
	return mode, nil
}

func decodeOpenCodeSnapshot(raw json.RawMessage) (any, error) {
	var c openCodeSnapshot
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, err
	}
	if c.Version != 1 || (!c.Agents && !c.Skills && !c.Plugins) || !filepath.IsAbs(c.RenderDir) || filepath.Clean(c.RenderDir) != c.RenderDir || len(c.SectionsDirs) == 0 || c.LayerCount < 0 || c.LayerCount > 1000000 || c.FragmentCount < 0 || c.FragmentCount > 1000000 || c.FileModes == nil {
		return nil, fmt.Errorf("invalid OpenCode generation configuration")
	}
	for _, section := range c.SectionsDirs {
		if !filepath.IsAbs(section) || filepath.Clean(section) != section {
			return nil, fmt.Errorf("invalid OpenCode section path: %s", section)
		}
	}
	for rel, mode := range c.FileModes {
		clean := filepath.Clean(filepath.FromSlash(rel))
		if rel == "" || clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || filepath.ToSlash(clean) != rel || mode > 0777 {
			return nil, fmt.Errorf("invalid OpenCode generation path: %q", rel)
		}
	}
	return c, nil
}

func openCodeSnapshotFor(plan *renderPlan, mode openCodeMode, sections []string, root string) (openCodeSnapshot, error) {
	c := openCodeSnapshot{Version: 1, Agents: mode.agents, Skills: mode.skills, Plugins: mode.plugins, SectionsDirs: append([]string{}, sections...), RenderDir: root, LayerCount: len(plan.layers), FragmentCount: plan.fragments, FileModes: map[string]uint32{}}
	for path, source := range plan.files {
		info, err := os.Stat(source)
		if err != nil {
			return c, err
		}
		c.FileModes[filepath.ToSlash(path)] = uint32(info.Mode().Perm())
	}
	return c, nil
}

func openCodeHandler() harnessHandler {
	return harnessHandler{
		parse:  parseOpenCodeMode,
		decode: decodeOpenCodeSnapshot,
		stub: func(value any, _ syncFlags, sections []string, root string) (any, error) {
			mode, ok := value.(openCodeMode)
			if !ok {
				return nil, fmt.Errorf("invalid OpenCode mode")
			}
			return openCodeSnapshot{Version: 1, Agents: mode.agents, Skills: mode.skills, Plugins: mode.plugins, SectionsDirs: append([]string{}, sections...), RenderDir: root, FileModes: map[string]uint32{}}, nil
		},
		build: func(value any, _ syncFlags, sections []string, root string) (*renderPlan, any, error) {
			mode, ok := value.(openCodeMode)
			if !ok {
				return nil, nil, fmt.Errorf("invalid OpenCode mode")
			}
			plan, err := buildOpenCodePlan(sections, mode.plugins)
			if err != nil {
				return nil, nil, err
			}
			plan.instructionsName = "AGENTS.md"
			config, err := openCodeSnapshotFor(plan, mode, sections, root)
			return plan, config, err
		},
		restore: func(dir string, saved snapshot) (*renderPlan, error) {
			config, ok := saved["opencode"].(openCodeSnapshot)
			if !ok {
				return nil, fmt.Errorf("generation config has no OpenCode mode")
			}
			return restorePlan(dir, config.FileModes, config.LayerCount, config.FragmentCount, "AGENTS.md")
		},
		prepare: prepareOpenCode,
	}
}

// prepareOpenCode checks every destination before any harness applies changes.
func prepareOpenCode(gen, root string, saved snapshot, plan *renderPlan, destinations runtimeDestinations) (func() error, error) {
	config, ok := saved["opencode"].(openCodeSnapshot)
	if !ok {
		return nil, fmt.Errorf("generation config has no OpenCode mode")
	}
	destination := destinations.OpenCode
	payload := filepath.Join(gen, "opencode")
	target := filepath.Join(root, ".targets", "opencode", targetHash(destination), filepath.Base(gen))
	var instructions agentsPlan
	var err error
	if config.Agents {
		instructions, err = planHarnessInstructions(destination, target, root, "opencode", "AGENTS.md")
		if err != nil {
			return nil, err
		}
	}
	plans := make([]linkPlan, 0, 3)
	components := []struct {
		name    string
		enabled bool
	}{{"skills", config.Skills}, {"resources", config.Agents}, {"plugins", config.Plugins}}
	for _, component := range components {
		if !component.enabled {
			continue
		}
		p, err := planHarnessComponent(destination, payload, target, root, "opencode", component.name, true)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	if err := preflightHarnessTarget(root, gen, "opencode", destination, plan); err != nil {
		return nil, err
	}
	return func() error {
		actual, _, err := materializeHarnessTarget(root, gen, "opencode", destination, plan)
		if err != nil {
			return err
		}
		if actual != target {
			return fmt.Errorf("unexpected OpenCode target: %s", actual)
		}
		if config.Agents {
			if err := applyHarnessInstructions(instructions, destination); err != nil {
				return err
			}
		}
		for _, component := range plans {
			if err := applyLinks(component); err != nil {
				return err
			}
		}
		fmt.Printf("synced OpenCode: generation=%s destination=%s\n", filepath.Base(gen), destination)
		return nil
	}, nil
}
