package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const snapshotFile = "config.json"

// snapshot preserves every selected harness under its top-level name.
type snapshot map[string]any
type snapshotCodex struct {
	Version       int      `json:"version"`
	Agents        bool     `json:"agents"`
	Skills        bool     `json:"skills"`
	Plugin        string   `json:"plugin"`
	PluginName    string   `json:"plugin_name"`
	RenderDir     string   `json:"render_dir"`
	SectionsDirs  []string `json:"sections_dirs"`
	LayerCount    int      `json:"layer_count"`
	FragmentCount int      `json:"fragment_count"`
	Legacy        struct {
		SkillMarketplace bool `json:"skill_marketplace"`
		InstructionsOnly bool `json:"instructions_only"`
		RulesOnly        bool `json:"rules_only"`
	} `json:"legacy"`
	FileModes map[string]uint32 `json:"file_modes"`
}

func snapshotBytes(s snapshot) ([]byte, error) {
	b, e := json.Marshal(s)
	if e != nil {
		return nil, e
	}
	return append(b, '\n'), nil
}
func makeSnapshot(p *renderPlan, mode syncMode, flags syncFlags, sections []string, renderDir string) (snapshot, error) {
	s := snapshot{}
	c := &snapshotCodex{}
	c.Version = 1
	c.Agents = mode.agents && !flags.RulesOnly
	c.Skills = mode.skills && !flags.InstructionsOnly && !flags.RulesOnly
	c.Plugin = "none"
	if mode.plugin {
		c.Plugin = "full"
		if !mode.includeGuidance {
			c.Plugin = "skills"
		}
	}
	c.PluginName = mode.pluginName
	c.RenderDir = renderDir
	c.SectionsDirs = append([]string{}, sections...)
	c.LayerCount = len(p.layers)
	c.FragmentCount = p.fragments
	c.Legacy.SkillMarketplace = flags.SkillMarketplace
	c.Legacy.InstructionsOnly = flags.InstructionsOnly
	c.Legacy.RulesOnly = flags.RulesOnly
	c.FileModes = map[string]uint32{}
	for path, source := range p.files {
		info, e := os.Stat(source)
		if e != nil {
			return s, e
		}
		c.FileModes[filepath.ToSlash(path)] = uint32(info.Mode().Perm())
	}
	s["codex"] = *c
	return s, nil
}

// readSnapshot accepts only the canonical format written by this version. This
// rejects duplicate fields, unknown fields, omitted fields, and trailing data.
func readSnapshot(gen string) (snapshot, error) {
	var s snapshot
	path := filepath.Join(gen, snapshotFile)
	data, e := os.ReadFile(path)
	if os.IsNotExist(e) {
		return nil, fmt.Errorf("generation has no config.json; rerender it before syncing")
	}
	if e != nil {
		return s, e
	}
	s, e = decodeSnapshot(data, harnessRegistry())
	if e != nil {
		return nil, fmt.Errorf("invalid generation config %s: %w", path, e)
	}
	return s, nil
}

// decodeSnapshot preserves every registered harness and rejects unknown fields,
// duplicate keys, missing fields, and noncanonical JSON.
func decodeSnapshot(data []byte, registry map[string]harnessHandler) (snapshot, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil || len(raw) == 0 {
		return nil, fmt.Errorf("invalid generation config")
	}
	s := snapshot{}
	for _, name := range sortedKeys(raw) {
		handler, ok := registry[name]
		if !ok {
			return nil, fmt.Errorf("unsupported harness in generation config: %s", name)
		}
		value, err := handler.decode(raw[name])
		if err != nil {
			return nil, fmt.Errorf("invalid %s generation config: %w", name, err)
		}
		s[name] = value
	}
	canonical, err := snapshotBytes(s)
	if err != nil || !bytes.Equal(data, canonical) {
		return nil, fmt.Errorf("invalid or unsupported generation config")
	}
	return s, nil
}

func decodeCodexSnapshot(data json.RawMessage) (any, error) {
	var c snapshotCodex
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	path := "codex"
	if c.Version != 1 || !filepath.IsAbs(c.RenderDir) || !pluginNamePattern.MatchString(c.PluginName) || c.FileModes == nil || len(c.SectionsDirs) == 0 || c.LayerCount < 0 || c.LayerCount > 1000000 || c.FragmentCount < 0 || c.FragmentCount > 1000000 {
		return nil, fmt.Errorf("invalid or unsupported generation config: %s", path)
	}
	if c.Plugin != "none" && c.Plugin != "full" && c.Plugin != "skills" {
		return nil, fmt.Errorf("invalid plugin in generation config: %s", path)
	}
	if !c.Agents && !c.Skills && c.Plugin == "none" && !c.Legacy.RulesOnly {
		return nil, fmt.Errorf("generation config selects no output: %s", path)
	}
	for _, section := range c.SectionsDirs {
		if !filepath.IsAbs(section) || filepath.Clean(section) != section {
			return nil, fmt.Errorf("invalid source path in generation config: %s", path)
		}
	}
	if (c.Legacy.SkillMarketplace && (c.Legacy.InstructionsOnly || c.Legacy.RulesOnly || c.Plugin != "full" || !c.Agents || c.Skills)) || (c.Legacy.InstructionsOnly && (c.Legacy.RulesOnly || c.Plugin != "none")) || (c.Legacy.RulesOnly && c.Plugin != "none") {
		return nil, fmt.Errorf("conflicting legacy modes in generation config: %s", path)
	}
	for rel, mode := range c.FileModes {
		clean := filepath.Clean(filepath.FromSlash(rel))
		if rel == "" || clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || filepath.ToSlash(clean) != rel || mode > 0777 {
			return nil, fmt.Errorf("invalid file in generation config: %q", rel)
		}
	}
	return c, nil
}

// restorePlan reconstructs a harness plan solely from its immutable payload.
func restorePlan(gen string, fileModes map[string]uint32, layers, fragments int, instructionsName string) (*renderPlan, error) {
	p := &renderPlan{instructionsName: instructionsName, skills: map[string]string{}, resources: map[string]string{}, dirs: map[string]bool{}, files: map[string]string{}, contents: map[string][]byte{}, fileModes: map[string]os.FileMode{}}
	p.layers = make([]string, layers)
	p.fragments = fragments
	b, e := os.ReadFile(filepath.Join(gen, instructionsName))
	if e != nil {
		return nil, e
	}
	p.agents = string(b)
	e = walk(gen, func(path string, d os.DirEntry) error {
		rel, e := filepath.Rel(gen, path)
		if e != nil {
			return e
		}
		if rel == instructionsName || (rel == snapshotFile && generationPattern.MatchString(filepath.Base(gen))) {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("generation contains symlink: %s", path)
		}
		if d.IsDir() {
			p.dirs[rel] = true
			return nil
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("generation contains non-regular file: %s", path)
		}
		mode, ok := fileModes[filepath.ToSlash(rel)]
		if !ok {
			return fmt.Errorf("generation file is absent from config: %s", path)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if info.Mode().Perm() != os.FileMode(mode)&^0222 {
			return fmt.Errorf("generation file mode does not match config: %s", path)
		}
		p.files[rel] = path
		p.fileModes[rel] = os.FileMode(mode)
		return nil
	})
	if e != nil {
		return nil, e
	}
	if len(p.files) != len(fileModes) {
		return nil, fmt.Errorf("generation files do not match config: %s", gen)
	}
	skills := filepath.Join(gen, "skills")
	es, e := entries(skills)
	if e != nil {
		return nil, e
	}
	for _, d := range es {
		path := filepath.Join(skills, d.Name())
		if !d.IsDir() {
			return nil, fmt.Errorf("invalid generated skill: %s", path)
		}
		name, e := readSkillName(filepath.Join(path, "SKILL.md"))
		if e != nil {
			return nil, e
		}
		p.skills[name] = path
	}
	for rel, source := range p.files {
		if strings.HasPrefix(rel, "resources"+string(filepath.Separator)) {
			p.resources[strings.TrimPrefix(rel, "resources"+string(filepath.Separator))] = source
		}
	}
	return p, nil
}
