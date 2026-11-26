package harness_codex

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const defaultPluginName = "skillset"

var pluginNamePattern = regexp.MustCompile("^[a-z0-9][a-z0-9_-]*$")

type syncMode struct {
	agents          bool
	skills          bool
	plugin          bool
	includeGuidance bool
	pluginName      string
}

// parseSyncMode validates the complete request before any destination changes.
func parseSyncMode(raw *string) (syncMode, error) {
	mode := syncMode{agents: true, skills: true, includeGuidance: true, pluginName: defaultPluginName}
	if raw == nil || *raw == Name {
		return mode, nil
	}
	harness, data, ok := strings.Cut(*raw, ":")
	if !ok || harness != "codex" {
		return mode, fmt.Errorf("mode must be codex or codex:<JSON object>")
	}
	decoder := json.NewDecoder(strings.NewReader(data))
	token, err := decoder.Token()
	if err != nil {
		return mode, fmt.Errorf("invalid mode JSON: %w", err)
	}
	if token != json.Delim('{') {
		return mode, fmt.Errorf("mode JSON must be an object")
	}
	seen := map[string]bool{}
	for decoder.More() {
		token, err = decoder.Token()
		if err != nil {
			return mode, fmt.Errorf("invalid mode JSON: %w", err)
		}
		key, ok := token.(string)
		if !ok {
			return mode, fmt.Errorf("invalid mode JSON key")
		}
		if seen[key] {
			return mode, fmt.Errorf("duplicate mode field %q", key)
		}
		seen[key] = true
		if key != "agents" && key != "skills" && key != "plugin" && key != "plugin_name" {
			return mode, fmt.Errorf("unknown mode field %q", key)
		}
		token, err = decoder.Token()
		if err != nil {
			return mode, fmt.Errorf("invalid mode JSON: %w", err)
		}
		if key == "plugin_name" {
			value, ok := token.(string)
			if !ok || !pluginNamePattern.MatchString(value) {
				return mode, fmt.Errorf("plugin_name must begin with a lowercase letter or digit and contain only lowercase letters, digits, hyphens, or underscores")
			}
			mode.pluginName = value
		} else if key == "plugin" {
			value, ok := token.(string)
			if !ok {
				return mode, fmt.Errorf("plugin must be a string")
			}
			switch value {
			case "full":
				mode.plugin = true
			case "skills":
				mode.plugin = true
				mode.includeGuidance = false
			default:
				return mode, fmt.Errorf("unknown plugin %q", value)
			}
		} else {
			value, ok := token.(bool)
			if !ok {
				return mode, fmt.Errorf("%s must be a boolean", key)
			}
			if key == "agents" {
				mode.agents = value
			} else {
				mode.skills = value
			}
		}
	}
	if _, err = decoder.Token(); err != nil {
		return mode, fmt.Errorf("invalid mode JSON: %w", err)
	}
	if _, err = decoder.Token(); err != io.EOF {
		if err != nil {
			return mode, fmt.Errorf("invalid mode JSON: %w", err)
		}
		return mode, fmt.Errorf("trailing content after mode JSON")
	}
	if seen["plugin_name"] && !mode.plugin {
		return mode, fmt.Errorf("plugin_name requires plugin output")
	}
	if mode.plugin && !seen["skills"] {
		mode.skills = false
	}
	if mode.plugin && mode.includeGuidance && !seen["agents"] {
		mode.agents = false
	}
	if !mode.agents && !mode.skills && !mode.plugin {
		return mode, fmt.Errorf("mode must select agents, skills, or a plugin")
	}
	return mode, nil
}
