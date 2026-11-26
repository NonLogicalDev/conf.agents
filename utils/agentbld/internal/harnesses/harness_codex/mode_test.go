package harness_codex

import (
	"encoding/json"
	"testing"
)

func TestModeRejectsInvalidPluginNames(t *testing.T) {
	for _, name := range []string{"", "..", "../outside", "/tmp/escape", "Uppercase", "bad name", "bad.name"} {
		t.Run(name, func(t *testing.T) {
			value, _ := json.Marshal(name)
			mode := `codex:{"agents":false,"plugin":"full","plugin_name":` + string(value) + "}"
			if _, err := parseSyncMode(&mode); err == nil {
				t.Fatalf("accepted plugin name %q", name)
			}
		})
	}
	for _, mode := range []string{
		`codex:{"plugin_name":"valid"}`,
		`codex:{"agents":false,"plugin_name":"valid"}`,
		`codex:{"plugin":"full","plugin_name":null}`,
		`codex:{"plugin":"full","plugin_name":true}`,
	} {
		if _, err := parseSyncMode(&mode); err == nil {
			t.Fatalf("accepted %s", mode)
		}
	}
}
