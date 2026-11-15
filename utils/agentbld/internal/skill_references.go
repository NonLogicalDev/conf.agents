package internal

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Portable skill links are resolved while rendering each harness payload.
var portableSkillReference = regexp.MustCompile(`\[\$@:([A-Za-z0-9_-]+)\]\(skill://@/([A-Za-z0-9_-]+)\)`)

func resolveSkillReferences(data []byte, namespace, labelPrefix string) []byte {
	prefix := strings.ReplaceAll(labelPrefix, "$", "$$")
	return portableSkillReference.ReplaceAll(data, []byte("["+prefix+"${1}](skill://"+namespace+"/${2})"))
}

func resolvePlanReferences(p *renderPlan, instructionsNamespace, instructionsPrefix, filesNamespace, filesPrefix string) (*renderPlan, error) {
	result := *p
	result.agents = string(resolveSkillReferences([]byte(p.agents), instructionsNamespace, instructionsPrefix))
	result.contents = make(map[string][]byte, len(p.contents))
	for name, data := range p.contents {
		result.contents[name] = data
	}
	for name, source := range p.files {
		if !strings.EqualFold(filepath.Ext(name), ".md") {
			continue
		}
		data, ok := p.contents[name]
		if !ok {
			var err error
			data, err = os.ReadFile(source)
			if err != nil {
				return nil, err
			}
		}
		result.contents[name] = resolveSkillReferences(data, filesNamespace, filesPrefix)
	}
	return &result, nil
}

func resolvedHarnessTarget(p *renderPlan, destination, harness string, snap snapshot) (*renderPlan, error) {
	resolved, err := targetPlan(p, destination)
	if err != nil {
		return nil, err
	}
	namespace, instructionsPrefix, filesPrefix := "local", "", ""
	switch harness {
	case "codex":
		instructionsPrefix, filesPrefix = "$", "$"
		if c, ok := snap["codex"].(snapshotCodex); ok && c.Plugin != "none" {
			namespace = c.PluginName
			instructionsPrefix = "$" + namespace + ":"
		}
	case "claude":
		instructionsPrefix, filesPrefix = "/", "/"
	}
	return resolvePlanReferences(resolved, namespace, instructionsPrefix, "local", filesPrefix)
}

// renderSkillReferences stores the two Codex skill variants independently so
// plugin and standalone references are fully resolved before activation.
func renderSkillReferences(p *renderPlan, harness, pluginName string, pluginGuidance bool) (*renderPlan, error) {
	localPrefix := ""
	if harness == "codex" {
		localPrefix = "$"
	}
	if harness == "claude" {
		localPrefix = "/"
	}
	local, err := resolvePlanReferences(p, "local", localPrefix, "local", localPrefix)
	if err != nil {
		return nil, err
	}
	if pluginName == "" {
		return local, nil
	}
	pluginPrefix := "$" + pluginName + ":"
	if harness == "claude" {
		pluginPrefix = "/" + pluginName + ":"
	}
	plugin, err := resolvePlanReferences(p, pluginName, pluginPrefix, pluginName, pluginPrefix)
	if err != nil {
		return nil, err
	}
	local.files = make(map[string]string, len(p.files)*2)
	local.dirs = make(map[string]bool, len(p.dirs)*2)
	for key, value := range p.files {
		local.files[key] = value
	}
	for key, value := range p.dirs {
		local.dirs[key] = value
	}
	local.dirs["plugin-skills"] = true
	for rel := range p.dirs {
		if rel == "skills" || !within(rel, "skills") {
			continue
		}
		part, _ := filepath.Rel("skills", rel)
		if !pluginGuidance && (part == "Agent_Guidance" || strings.HasPrefix(part, "Agent_Guidance"+string(filepath.Separator))) {
			continue
		}
		local.dirs[filepath.Join("plugin-skills", part)] = true
	}
	for rel, source := range p.files {
		if !within(rel, "skills") || rel == "skills" {
			continue
		}
		part, _ := filepath.Rel("skills", rel)
		if !pluginGuidance && (part == "Agent_Guidance" || strings.HasPrefix(part, "Agent_Guidance"+string(filepath.Separator))) {
			continue
		}
		target := filepath.Join("plugin-skills", part)
		local.files[target] = source
		if data, ok := plugin.contents[rel]; ok {
			local.contents[target] = data
		}
	}
	return local, nil
}

// projectPluginSkills replaces standalone skills with the rendered plugin copy.
// A missing plugin-skills directory denotes an older portable generation.
func projectPluginSkills(p *renderPlan, namespace string) (*renderPlan, error) {
	if !p.dirs["plugin-skills"] {
		return resolvePlanReferences(p, namespace, "$"+namespace+":", namespace, "$"+namespace+":")
	}
	result := *p
	result.dirs = map[string]bool{"skills": true}
	result.files = map[string]string{}
	result.contents = map[string][]byte{}
	result.fileModes = map[string]os.FileMode{}
	for dir := range p.dirs {
		if dir == "skills" || within(dir, "skills") || dir == "plugin-skills" {
			continue
		}
		if within(dir, "plugin-skills") {
			part, _ := filepath.Rel("plugin-skills", dir)
			result.dirs[filepath.Join("skills", part)] = true
		} else {
			result.dirs[dir] = true
		}
	}
	for rel, src := range p.files {
		if within(rel, "skills") {
			continue
		}
		target := rel
		if within(rel, "plugin-skills") {
			part, _ := filepath.Rel("plugin-skills", rel)
			target = filepath.Join("skills", part)
		}
		result.files[target] = src
		if data, ok := p.contents[rel]; ok {
			result.contents[target] = data
		}
		if mode, ok := p.fileModes[rel]; ok {
			result.fileModes[target] = mode
		}
	}
	result.skills = map[string]string{}
	for name, source := range p.skills {
		if result.dirs[filepath.Join("skills", name)] {
			result.skills[name] = source
		}
	}
	return &result, nil
}

// targetWithoutPluginSkills prevents the plugin-only copy leaking into a harness target.
func targetWithoutPluginSkills(p *renderPlan) *renderPlan {
	if !p.dirs["plugin-skills"] {
		return p
	}
	result := *p
	result.dirs, result.files, result.contents, result.fileModes = map[string]bool{}, map[string]string{}, map[string][]byte{}, map[string]os.FileMode{}
	for key, value := range p.dirs {
		if !within(key, "plugin-skills") {
			result.dirs[key] = value
		}
	}
	for key, value := range p.files {
		if within(key, "plugin-skills") {
			continue
		}
		result.files[key] = value
		if data, ok := p.contents[key]; ok {
			result.contents[key] = data
		}
		if mode, ok := p.fileModes[key]; ok {
			result.fileModes[key] = mode
		}
	}
	return &result
}
