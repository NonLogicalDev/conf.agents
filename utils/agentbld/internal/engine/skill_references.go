package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const pluginSkillsetDir = "plugin-skillset"

// Portable skill links are resolved while rendering each harness payload.
var portableSkillReference = regexp.MustCompile(`\[\$@:([A-Za-z0-9_-]+)\]\(skill://@/([A-Za-z0-9_-]+)\)`)

func ResolveSkillReferences(data []byte, namespace, labelPrefix string) []byte {
	prefix := strings.ReplaceAll(labelPrefix, "$", "$$")
	return portableSkillReference.ReplaceAll(data, []byte("["+prefix+"${1}](skill://"+namespace+"/${2})"))
}

func resolvePlanReferences(p *RenderPlan, instructionsNamespace, instructionsPrefix, filesNamespace, filesPrefix string) (*RenderPlan, error) {
	result := *p
	result.Agents = string(ResolveSkillReferences([]byte(p.Agents), instructionsNamespace, instructionsPrefix))
	result.Contents = make(map[string][]byte, len(p.Contents))
	for name, data := range p.Contents {
		result.Contents[name] = data
	}
	for name, source := range p.Files {
		if !strings.EqualFold(filepath.Ext(name), ".md") {
			continue
		}
		data, ok := p.Contents[name]
		if !ok {
			var err error
			data, err = os.ReadFile(source)
			if err != nil {
				return nil, err
			}
		}
		result.Contents[name] = ResolveSkillReferences(data, filesNamespace, filesPrefix)
	}
	return &result, nil
}

// RenderSkillReferences stores standalone and plugin variants independently.
func RenderSkillReferences(p *RenderPlan, localPrefix, pluginName, pluginPrefix string, pluginGuidance bool) (*RenderPlan, error) {
	local, err := resolvePlanReferences(p, "local", localPrefix, "local", localPrefix)
	if err != nil {
		return nil, err
	}
	if pluginName == "" {
		return local, nil
	}
	plugin, err := resolvePlanReferences(p, pluginName, pluginPrefix, pluginName, pluginPrefix)
	if err != nil {
		return nil, err
	}
	local.Files = make(map[string]string, len(p.Files)*2)
	local.Dirs = make(map[string]bool, len(p.Dirs)*2)
	for key, value := range p.Files {
		local.Files[key] = value
	}
	for key, value := range p.Dirs {
		local.Dirs[key] = value
	}
	local.Dirs[pluginSkillsetDir] = true
	for rel := range p.Dirs {
		if rel == "skills" || !Within(rel, "skills") {
			continue
		}
		part, _ := filepath.Rel("skills", rel)
		if !pluginGuidance && (part == "Agent_Guidance" || strings.HasPrefix(part, "Agent_Guidance"+string(filepath.Separator))) {
			continue
		}
		local.Dirs[filepath.Join(pluginSkillsetDir, part)] = true
	}
	for rel, source := range p.Files {
		if !Within(rel, "skills") || rel == "skills" {
			continue
		}
		part, _ := filepath.Rel("skills", rel)
		if !pluginGuidance && (part == "Agent_Guidance" || strings.HasPrefix(part, "Agent_Guidance"+string(filepath.Separator))) {
			continue
		}
		target := filepath.Join(pluginSkillsetDir, part)
		local.Files[target] = source
		if data, ok := plugin.Contents[rel]; ok {
			local.Contents[target] = data
		}
	}
	return local, nil
}

// ProjectPluginSkills replaces standalone skills with the rendered plugin copy.
func ProjectPluginSkills(p *RenderPlan) (*RenderPlan, error) {
	skillset := pluginSkillsetDir
	if !p.Dirs[skillset] {
		return nil, fmt.Errorf("generation has no %s directory", skillset)
	}
	result := *p
	result.Dirs = map[string]bool{"skills": true}
	result.Files = map[string]string{}
	result.Contents = map[string][]byte{}
	result.fileModes = map[string]os.FileMode{}
	for dir := range p.Dirs {
		if Within(dir, "skills") || dir == skillset {
			continue
		}
		if Within(dir, skillset) {
			part, _ := filepath.Rel(skillset, dir)
			result.Dirs[filepath.Join("skills", part)] = true
		} else {
			result.Dirs[dir] = true
		}
	}
	for rel, src := range p.Files {
		if Within(rel, "skills") {
			continue
		}
		target := rel
		if Within(rel, skillset) {
			part, _ := filepath.Rel(skillset, rel)
			target = filepath.Join("skills", part)
		}
		result.Files[target] = src
		if data, ok := p.Contents[rel]; ok {
			result.Contents[target] = data
		}
		if mode, ok := p.fileModes[rel]; ok {
			result.fileModes[target] = mode
		}
	}
	result.Skills = map[string]string{}
	for name, source := range p.Skills {
		if result.Dirs[filepath.Join("skills", name)] {
			result.Skills[name] = source
		}
	}
	return &result, nil
}

// targetWithoutPluginSkills prevents the plugin-only copy leaking into a harness target.
func targetWithoutPluginSkills(p *RenderPlan) *RenderPlan {
	skillset := pluginSkillsetDir
	if !p.Dirs[skillset] {
		return p
	}
	result := *p
	result.Dirs, result.Files, result.Contents, result.fileModes = map[string]bool{}, map[string]string{}, map[string][]byte{}, map[string]os.FileMode{}
	for key, value := range p.Dirs {
		if !Within(key, skillset) {
			result.Dirs[key] = value
		}
	}
	for key, value := range p.Files {
		if Within(key, skillset) {
			continue
		}
		result.Files[key] = value
		if data, ok := p.Contents[key]; ok {
			result.Contents[key] = data
		}
		if mode, ok := p.fileModes[key]; ok {
			result.fileModes[key] = mode
		}
	}
	return &result
}
