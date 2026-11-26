package harness_claude

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

// claudeMode selects Claude instructions, standalone skills, and plugin output.
const defaultPluginName = "skillset"

var pluginNamePattern = regexp.MustCompile("^[a-z0-9][a-z0-9_-]*$")

type claudeMode struct {
	agents, skills, plugin, includeGuidance bool
	pluginName                              string
}

func parseClaudeMode(raw *string) (any, error) {
	mode := claudeMode{agents: true, skills: true, includeGuidance: true, pluginName: defaultPluginName}
	if raw == nil || *raw == "claude" {
		return mode, nil
	}
	harness, data, ok := strings.Cut(*raw, ":")
	if !ok || harness != "claude" {
		return nil, fmt.Errorf("mode must be claude or claude:<JSON object>")
	}
	dec := json.NewDecoder(strings.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("invalid Claude mode JSON: %w", err)
	}
	if tok != json.Delim('{') {
		return nil, fmt.Errorf("Claude mode JSON must be an object")
	}
	seen := map[string]bool{}
	for dec.More() {
		tok, err = dec.Token()
		if err != nil {
			return nil, fmt.Errorf("invalid Claude mode JSON: %w", err)
		}
		key, ok := tok.(string)
		if !ok {
			return nil, fmt.Errorf("invalid Claude mode field")
		}
		if seen[key] {
			return nil, fmt.Errorf("duplicate Claude mode field %q", key)
		}
		seen[key] = true
		if key != "agents" && key != "skills" && key != "plugin" && key != "plugin_name" {
			return nil, fmt.Errorf("unknown Claude mode field %q", key)
		}
		tok, err = dec.Token()
		if err != nil {
			return nil, fmt.Errorf("invalid Claude mode JSON: %w", err)
		}
		switch key {
		case "plugin_name":
			value, ok := tok.(string)
			if !ok || !pluginNamePattern.MatchString(value) {
				return mode, fmt.Errorf("invalid Claude plugin_name")
			}
			mode.pluginName = value
		case "plugin":
			value, ok := tok.(string)
			if !ok || (value != "full" && value != "skills") {
				return mode, fmt.Errorf("Claude plugin must be full or skills")
			}
			mode.plugin = true
			mode.includeGuidance = value == "full"
		default:
			value, ok := tok.(bool)
			if !ok {
				return mode, fmt.Errorf("Claude %s must be a boolean", key)
			}
			if key == "agents" {
				mode.agents = value
			} else {
				mode.skills = value
			}
		}
	}
	if _, err = dec.Token(); err != nil {
		return nil, fmt.Errorf("invalid Claude mode JSON: %w", err)
	}
	if _, err = dec.Token(); err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("invalid Claude mode JSON: %w", err)
		}
		return nil, fmt.Errorf("trailing content after Claude mode JSON")
	}
	if seen["plugin_name"] && !mode.plugin {
		return nil, fmt.Errorf("Claude plugin_name requires plugin output")
	}
	if mode.plugin && !seen["skills"] {
		mode.skills = false
	}
	if mode.plugin && mode.includeGuidance && !seen["agents"] {
		mode.agents = false
	}
	if !mode.agents && !mode.skills && !mode.plugin {
		return nil, fmt.Errorf("Claude mode must select agents, skills, or a plugin")
	}
	return mode, nil
}

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

func decodeClaudeSnapshot(data json.RawMessage) (any, error) {
	var c Snapshot
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if !pluginNamePattern.MatchString(c.PluginName) || (c.Plugin != "none" && c.Plugin != "full" && c.Plugin != "skills") {
		return nil, fmt.Errorf("invalid Claude plugin config")
	}
	if c.Version != 1 || (!c.Agents && !c.Skills && c.Plugin == "none") || !filepath.IsAbs(c.RenderDir) || filepath.Clean(c.RenderDir) != c.RenderDir || len(c.SectionsDirs) == 0 || c.FileModes == nil || c.LayerCount < 0 || c.LayerCount > 1000000 || c.FragmentCount < 0 || c.FragmentCount > 1000000 {
		return nil, fmt.Errorf("invalid Claude generation config")
	}
	for _, p := range c.SectionsDirs {
		if !filepath.IsAbs(p) || filepath.Clean(p) != p {
			return nil, fmt.Errorf("invalid Claude source path: %s", p)
		}
	}
	for rel, mode := range c.FileModes {
		clean := filepath.Clean(filepath.FromSlash(rel))
		if rel == "" || clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || filepath.ToSlash(clean) != rel || mode > 0777 {
			return nil, fmt.Errorf("invalid Claude generation file: %q", rel)
		}
	}
	return c, nil
}

// claudePath follows Claude Code's configuration-directory override.
func Destination(path string) (string, error) {
	if path == "" {
		path = os.Getenv("CLAUDE_CONFIG_DIR")
	}
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, ".claude")
	}
	return engine.Absolute(path)
}
