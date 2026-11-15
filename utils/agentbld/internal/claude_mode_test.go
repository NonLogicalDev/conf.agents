package internal

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"
)

func TestClaudeModeParsing(t *testing.T) {
	for _, tc := range []struct {
		input          string
		agents, skills bool
	}{
		{"claude", true, true}, {`claude:{}`, true, true}, {`claude:{"agents":false}`, false, true}, {`claude:{"skills":false}`, true, false}, {`claude:{"skills":true,"agents":false}`, false, true},
	} {
		t.Run(tc.input, func(t *testing.T) {
			got, err := parseClaudeMode(&tc.input, syncFlags{})
			if err != nil {
				t.Fatal(err)
			}
			if got != (claudeMode{agents: tc.agents, skills: tc.skills, includeGuidance: true, pluginName: marketplacePlugin}) {
				t.Fatalf("got %#v", got)
			}
		})
	}
	for _, input := range []string{`claude:{"agents":false,"skills":false}`, `claude:{"agents":null}`, `claude:{"skills":"yes"}`, `claude:{"plugin":"invalid"}`, `claude:{"agents":true,"agents":false}`, `claude:[]`, `claude:{} trailing`, `other:{}`} {
		t.Run(input, func(t *testing.T) {
			if _, err := parseClaudeMode(&input, syncFlags{}); err == nil {
				t.Fatalf("accepted %s", input)
			}
		})
	}
}
func TestClaudeSnapshotRoundTripAndValidation(t *testing.T) {
	root := t.TempDir()
	want := snapshotClaude{Version: 1, Agents: true, Skills: true, RenderDir: root, SectionsDirs: []string{filepath.Join(root, "sections")}, FileModes: map[string]uint32{"skills/A/SKILL.md": 0644}}
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := decodeClaudeSnapshot(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v", got)
	}
	for _, data := range []string{`{}`, `{"version":1,"agents":false,"skills":false}`, `{"version":1,"agents":true,"skills":true,"render_dir":"/tmp","sections_dirs":["/tmp"],"file_modes":{"../escape":420}}`} {
		if _, err := decodeClaudeSnapshot(json.RawMessage(data)); err == nil {
			t.Fatalf("accepted invalid metadata %s", data)
		}
	}
}
