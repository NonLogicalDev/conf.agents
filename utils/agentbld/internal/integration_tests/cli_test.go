package integration_tests

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLICommandHelp(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		want, absent []string
	}{
		{"root", []string{"--help"}, []string{"<command>", "render", "check", "sync", "prune", "reset", "profile"}, []string{"<operation>"}},
		{"sync", []string{"sync", "--help"}, []string{"agentbld sync", "--profile", "--save-profile"}, []string{"--home"}},
		{"profile", []string{"profile", "--help"}, []string{"agentbld profile", "list", "show"}, nil},
		{"render", []string{"render", "--help"}, []string{"agentbld render", "immutable generation", "--sync"}, []string{"--instructions-only", "--rules-only", "--skill-marketplace"}},
		{"check", []string{"check", "--help"}, []string{"agentbld check", "latest generation"}, []string{"--skill-marketplace"}},
		{"prune", []string{"prune", "--help"}, []string{"agentbld prune", "inactive generations"}, []string{"--rules-only"}},
		{"reset", []string{"reset", "--help"}, []string{"agentbld reset", "Back up"}, []string{"--rules-only"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := exec.Command(contractBinary, tt.args...).CombinedOutput()
			if err != nil {
				t.Fatalf("help failed: %v\n%s", err, out)
			}
			text := strings.Join(strings.Fields(string(out)), " ")
			for _, value := range tt.want {
				if !strings.Contains(text, value) {
					t.Errorf("help missing %q:\n%s", value, text)
				}
			}
			for _, value := range tt.absent {
				if strings.Contains(text, value) {
					t.Errorf("help unexpectedly contains %q:\n%s", value, text)
				}
			}
		})
	}
}

func TestCLICommandRestrictions(t *testing.T) {
	f := newContractFixture(t)
	tests := []struct {
		name string
		args []string
	}{
		{"unknown command", []string{"unknown"}},
		{"missing command", nil},
		{"removed config-dir", []string{"render", "--config-dir", "."}},
		{"removed home", []string{"render", "--home", "."}},
		{"removed skill-marketplace", []string{"render", "--skill-marketplace"}},
		{"removed instructions-only", []string{"render", "--instructions-only"}},
		{"removed rules-only", []string{"render", "--rules-only"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{}, tt.args...)
			args = append(args, "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", f.codex)
			out, err := exec.Command(contractBinary, args...).CombinedOutput()
			exit, ok := err.(*exec.ExitError)
			if !ok || exit.ExitCode() != 2 {
				t.Fatalf("expected exit 2, got %v:\n%s", err, out)
			}
		})
	}
}

// Existing callers may place shared flags on either side of the command.
func TestCLISharedFlagsBeforeAndAfterCommand(t *testing.T) {
	f := newContractFixture(t)
	common := []string{"--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", f.codex}
	before := append(append([]string{}, common...), "render")
	if out, err := exec.Command(contractBinary, before...).CombinedOutput(); err != nil {
		t.Fatalf("flags before command: %v\n%s", err, out)
	}
	after := append([]string{"check"}, common...)
	if out, err := exec.Command(contractBinary, after...).CombinedOutput(); err != nil {
		t.Fatalf("flags after command: %v\n%s", err, out)
	}
}
