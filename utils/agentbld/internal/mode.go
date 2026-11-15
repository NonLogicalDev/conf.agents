package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type syncMode struct {
	agents          bool
	skills          bool
	plugin          bool
	includeGuidance bool
	pluginName      string
}

var pluginNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

// parseSyncMode validates the complete request before any destination changes.
func parseSyncMode(raw *string, flags syncFlags) (syncMode, error) {
	mode := syncMode{agents: true, skills: true, includeGuidance: true, pluginName: marketplacePlugin}
	if raw == nil {
		mode.plugin = flags.SkillMarketplace
		if mode.plugin {
			mode.skills = false
		}
		return mode, nil
	}
	if flags.SkillMarketplace || flags.InstructionsOnly || flags.RulesOnly {
		return mode, fmt.Errorf("--mode cannot be combined with legacy sync flags")
	}
	if *raw == "codex" {
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

type modeParser func(*string, syncFlags) (any, error)

type resolvedModes = map[string]any

// resolveModes preserves each selected harness and rejects duplicate names.
func resolveModes(values []string, flags syncFlags, parsers map[string]modeParser) (resolvedModes, error) {
	if len(values) == 0 {
		values = []string{"codex"}
	}
	resolved := resolvedModes{}
	for _, value := range values {
		harness, _, _ := strings.Cut(value, ":")
		parser, ok := parsers[harness]
		if !ok {
			return nil, fmt.Errorf("unsupported harness %q", harness)
		}
		if _, exists := resolved[harness]; exists {
			return nil, fmt.Errorf("duplicate mode for harness %q", harness)
		}
		mode, err := parser(&value, flags)
		if err != nil {
			return nil, err
		}
		resolved[harness] = mode
	}
	return resolved, nil
}

func parseModes(values []string, flags syncFlags) (resolvedModes, error) {
	handlers := harnessRegistry()
	parsers := map[string]modeParser{}
	for name, handler := range handlers {
		parsers[name] = handler.parse
	}
	if len(values) == 0 {
		parser, ok := parsers["codex"]
		if !ok {
			return nil, fmt.Errorf("default codex harness is unavailable")
		}
		mode, err := parser(nil, flags)
		if err != nil {
			return nil, err
		}
		return resolvedModes{"codex": mode}, nil
	}
	return resolveModes(values, flags, parsers)
}

// printModes does not resolve paths or inspect the filesystem.
func printModes() {
	fmt.Println(`config: {"command":"modes"}`)
	fmt.Println(`Mode syntax:
  --mode=<harness>[:<JSON object>]

Available harnesses:
  codex      Codex instructions, skills, and plugin.
  opencode   OpenCode instructions, skills, and native plugins.
  claude     Claude Code instructions, skills, and plugin.

The JSON object is optional and each harness has its own options. Repeat --mode
once per harness to include several in one generation. The default is codex.

Codex JSON options:
  agents       Boolean. Sync instructions, resources, and rules. Default: true,
               except plugin "full" defaults to false.
  skills       Boolean. Sync standalone Codex skills. Default: true without a
               plugin; false when a plugin is selected.
  plugin       String: "full" includes Agent_Guidance; "skills" omits it.
               Omit this option to skip plugin output. Skills keeps agents on
               by default; full turns agents off by default.
  plugin_name  String. Name the generated plugin (default: "skillset"). Use
               lowercase letters, digits, hyphens, and underscores; the first
               character must be a letter or digit. Requires a plugin.

Explicit agents and skills values override their defaults. At least one of
agents, skills, or plugin must be selected. Render records the selected mode;
check compares that mode with the latest generation, and sync activates it.
For sync, repeat --mode with bare harness names to activate only those saved
harnesses. Without a filter, sync activates all harnesses in the generation.

OpenCode JSON options:
  agents       Boolean. Sync instructions and resources. Default: true.
  skills       Boolean. Sync standalone skills. Default: true.
  plugins      Boolean. Sync native plugins. Default: true.
  At least one option must be enabled.

Claude JSON options:
  agents       Boolean. Sync instructions and resources. Default: true,
               except plugin "full" defaults to false.
  skills       Boolean. Sync standalone skills. Default: true without a
               plugin; false when a plugin is selected.
  plugin       String: "full" includes Agent_Guidance; "skills" omits it.
               Omit to skip plugin output.
  plugin_name  String. Generated plugin name (default: "skillset").
  At least one of agents, skills, or plugin must be selected.

Examples (run from the repository root):
  agentbld render --sync --sections-dir ./sections --mode=codex
  agentbld render --sync --sections-dir ./sections --mode=codex --mode=opencode --mode=claude
  agentbld sync 01 --mode=opencode
  agentbld render --sections-dir ./sections --mode='opencode:{"agents":true,"skills":false,"plugins":true}'
  agentbld render --sections-dir ./sections --mode='claude:{"agents":true,"skills":false}'
  agentbld render --sync --sections-dir ./sections --mode='claude:{"plugin":"skills"}'
  agentbld render --sync --sections-dir ./sections --mode='codex:{"plugin":"full"}'
  agentbld render --sync --sections-dir ./sections --mode='codex:{"plugin":"skills"}'
  agentbld render --sync --sections-dir ./sections --mode='codex:{"agents":true,"skills":true,"plugin":"full"}'
  agentbld render --sync --sections-dir ./sections --mode='codex:{"plugin":"skills","plugin_name":"my-skills"}'`)
}
