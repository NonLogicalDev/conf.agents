package harness_codex

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

type Snapshot struct {
	Version       int               `json:"version"`
	Agents        bool              `json:"agents"`
	Skills        bool              `json:"skills"`
	Plugin        string            `json:"plugin"`
	PluginName    string            `json:"plugin_name"`
	RenderDir     string            `json:"render_dir"`
	SectionsDirs  []string          `json:"sections_dirs"`
	LayerCount    int               `json:"layer_count"`
	FragmentCount int               `json:"fragment_count"`
	FileModes     map[string]uint32 `json:"file_modes"`
}

func makeSnapshot(p *engine.RenderPlan, mode syncMode, sections []string, renderDir string) (engine.Snapshot, error) {
	s := engine.Snapshot{}
	c := &Snapshot{}
	c.Version = 1
	c.Agents = mode.agents
	c.Skills = mode.skills
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
	c.LayerCount = len(p.Layers)
	c.FragmentCount = p.Fragments
	c.FileModes = map[string]uint32{}
	for path, source := range p.Files {
		info, e := os.Stat(source)
		if e != nil {
			return s, e
		}
		c.FileModes[filepath.ToSlash(path)] = uint32(info.Mode().Perm())
	}
	s["codex"] = *c
	return s, nil
}

func decodeCodexSnapshot(data json.RawMessage) (any, error) {
	var c Snapshot
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
	if !c.Agents && !c.Skills && c.Plugin == "none" {
		return nil, fmt.Errorf("generation config selects no output: %s", path)
	}
	for _, section := range c.SectionsDirs {
		if !filepath.IsAbs(section) || filepath.Clean(section) != section {
			return nil, fmt.Errorf("invalid source path in generation config: %s", path)
		}
	}
	for rel, mode := range c.FileModes {
		clean := filepath.Clean(filepath.FromSlash(rel))
		if rel == "" || clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || filepath.ToSlash(clean) != rel || mode > 0777 {
			return nil, fmt.Errorf("invalid file in generation config: %q", rel)
		}
	}
	return c, nil
}
