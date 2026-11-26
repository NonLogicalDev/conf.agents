package harness_opencode

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

// openCodePath resolves the destination without following its symlinks.
func Destination(path string) (string, error) {
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
	return engine.Absolute(path)
}

// discoverOpenCodePlugins selects flat native modules, with later layers winning.
func discoverOpenCodePlugins(layers []string) (map[string]string, error) {
	plugins := map[string]string{}
	for _, layer := range layers {
		root := filepath.Join(layer, "opencode")
		dir := filepath.Join(root, "plugins")
		if engine.IsLink(root) || engine.IsLink(dir) {
			return nil, fmt.Errorf("symlinked OpenCode plugin roots are not supported: %s", dir)
		}
		if _, err := os.Lstat(dir); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if !engine.IsDir(dir) {
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
func buildOpenCodePlan(registry engine.Registry, sections []string, includePlugins bool) (*engine.RenderPlan, error) {
	plan, err := registry.BuildPlan(sections, true, "opencode")
	if err != nil {
		return nil, err
	}
	if includePlugins {
		plugins, err := discoverOpenCodePlugins(plan.Layers)
		if err != nil {
			return nil, err
		}
		plan.Dirs["plugins"] = true
		for name, source := range plugins {
			plan.Files[filepath.Join("plugins", name)] = source
		}
	}
	plan.InstructionsName = "AGENTS.md"
	plan.Agents = string(engine.ResolveSkillReferences([]byte(plan.Agents), "local", ""))
	if err := engine.BundleResources(plan, plan.Agents); err != nil {
		return nil, err
	}
	return engine.RenderSkillReferences(plan, "", "", "", false)
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

func parseOpenCodeMode(raw *string) (any, error) {
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

func openCodeSnapshotFor(plan *engine.RenderPlan, mode openCodeMode, sections []string, root string) (openCodeSnapshot, error) {
	c := openCodeSnapshot{Version: 1, Agents: mode.agents, Skills: mode.skills, Plugins: mode.plugins, SectionsDirs: append([]string{}, sections...), RenderDir: root, LayerCount: len(plan.Layers), FragmentCount: plan.Fragments, FileModes: map[string]uint32{}}
	for path, source := range plan.Files {
		info, err := os.Stat(source)
		if err != nil {
			return c, err
		}
		c.FileModes[filepath.ToSlash(path)] = uint32(info.Mode().Perm())
	}
	return c, nil
}

func Handler(options *Options) engine.Handler {
	return engine.Handler{
		Summary: "OpenCode instructions, skills, and native plugins.",
		ModeHelp: `OpenCode JSON options:
  agents       Boolean. Sync instructions and resources. Default: true.
  skills       Boolean. Sync standalone skills. Default: true.
  plugins      Boolean. Sync native plugins. Default: true.
  At least one option must be enabled.`,
		Examples: `  agentbld sync 01 --mode=opencode
  agentbld render --sections-dir ./sections --mode='opencode:{"agents":true,"skills":false,"plugins":true}'`,
		ResolveDestination: options.resolveDestination,
		LoadProfile:        options.loadProfile,
		ConfigReport: func(_ any, dest engine.Destination) (map[string]any, map[string]string) {
			return nil, dest.Profile
		},
		LinkComponents: []string{"skills", "resources", "plugins"},

		Parse:  parseOpenCodeMode,
		Decode: decodeOpenCodeSnapshot,
		Stub: func(value any, sections []string, root string) (any, error) {
			mode, ok := value.(openCodeMode)
			if !ok {
				return nil, fmt.Errorf("invalid OpenCode mode")
			}
			return openCodeSnapshot{Version: 1, Agents: mode.agents, Skills: mode.skills, Plugins: mode.plugins, SectionsDirs: append([]string{}, sections...), RenderDir: root, FileModes: map[string]uint32{}}, nil
		},
		Build: func(registry engine.Registry, value any, sections []string, root string) (*engine.RenderPlan, any, error) {
			mode, ok := value.(openCodeMode)
			if !ok {
				return nil, nil, fmt.Errorf("invalid OpenCode mode")
			}
			plan, err := buildOpenCodePlan(registry, sections, mode.plugins)
			if err != nil {
				return nil, nil, err
			}
			plan.InstructionsName = "AGENTS.md"
			config, err := openCodeSnapshotFor(plan, mode, sections, root)
			return plan, config, err
		},
		Restore: func(dir string, saved engine.Snapshot) (*engine.RenderPlan, error) {
			config, ok := saved["opencode"].(openCodeSnapshot)
			if !ok {
				return nil, fmt.Errorf("generation config has no OpenCode mode")
			}
			return engine.RestorePlan(dir, config.FileModes, config.LayerCount, config.FragmentCount, "AGENTS.md")
		},
		Prepare: prepareOpenCode,
	}
}

// prepareOpenCode checks every destination before any harness applies changes.
func prepareOpenCode(registry engine.Registry, gen, root string, saved engine.Snapshot, plan *engine.RenderPlan, destinations engine.Destination) (func() error, error) {
	config, ok := saved["opencode"].(openCodeSnapshot)
	if !ok {
		return nil, fmt.Errorf("generation config has no OpenCode mode")
	}
	destination := destinations.Config
	payload := filepath.Join(gen, "opencode")
	target := filepath.Join(root, ".targets", "opencode", engine.TargetHash(destination), filepath.Base(gen))
	var instructions engine.AgentsPlan
	var err error
	if config.Agents {
		instructions, err = registry.PlanHarnessInstructions(destination, target, root, "opencode", "AGENTS.md")
		if err != nil {
			return nil, err
		}
	}
	plans := make([]engine.LinkPlan, 0, 3)
	components := []struct {
		name    string
		enabled bool
	}{{"skills", config.Skills}, {"resources", config.Agents}, {"plugins", config.Plugins}}
	for _, component := range components {
		if !component.enabled {
			continue
		}
		p, err := registry.PlanHarnessComponent(destination, payload, target, root, "opencode", component.name, true)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	if err := registry.PreflightHarnessTarget(root, gen, "opencode", destination, plan); err != nil {
		return nil, err
	}
	return func() error {
		actual, _, err := registry.MaterializeHarnessTarget(root, gen, "opencode", destination, plan)
		if err != nil {
			return err
		}
		if actual != target {
			return fmt.Errorf("unexpected OpenCode target: %s", actual)
		}
		if config.Agents {
			if err := engine.ApplyHarnessInstructions(instructions, destination); err != nil {
				return err
			}
		}
		for _, component := range plans {
			if err := engine.ApplyLinks(component); err != nil {
				return err
			}
		}
		fmt.Printf("synced OpenCode: generation=%s destination=%s\n", filepath.Base(gen), destination)
		return nil
	}, nil
}
