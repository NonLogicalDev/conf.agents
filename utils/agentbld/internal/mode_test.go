package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func modeFixture(t *testing.T) *contractFixture {
	t.Helper()
	f := newContractFixture(t)
	f.put(f.layer, "skills/Agent_Guidance/SKILL.md", "---\nname: Agent_Guidance\ndescription: Guidance\n---\n{{%_agents_%}}\n")
	return f
}

func modeRun(t *testing.T, f *contractFixture, want int, extra ...string) string {
	t.Helper()
	args := []string{"render", "--sync", "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", f.codex, "--home", f.home}
	args = append(args, extra...)
	out, err := exec.Command(contractBinary, args...).CombinedOutput()
	code := 0
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			code = e.ExitCode()
		} else {
			t.Fatal(err)
		}
	}
	if code != want {
		t.Fatalf("%v: exit %d, want %d\n%s", args, code, want, out)
	}
	return string(out)
}

func modePlugin(f *contractFixture) string {
	return filepath.Join(f.home, ".agents/plugins/skillset")
}

func TestModeDefaultCodex(t *testing.T) {
	for _, args := range [][]string{nil, {"--mode", "codex"}, {"--mode", "codex:{}"}} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			f := modeFixture(t)
			modeRun(t, f, 0, args...)
			for _, p := range []string{"AGENTS.md", "skills/Alpha", "skills/Agent_Guidance", "resources/guide", "rules/default.rules"} {
				f.symlink(filepath.Join(f.codex, p))
			}
			f.absent(modePlugin(f))
		})
	}
}

func TestModeCodexAndLocalPlugin(t *testing.T) {
	for _, tc := range []struct {
		name, mode           string
		guidance, standalone bool
	}{
		{"default-guidance", `codex:{"plugin":"full","agents":true}`, true, false},
		{"explicit-guidance", `codex:{"plugin":"full","agents":true,"skills":true}`, true, true},
		{"exclude-guidance", `codex:{"plugin":"skills"}`, false, false},
		{"skills-and-standalone", `codex:{"skills":true,"plugin":"skills"}`, false, true},
		{"full-and-standalone-reverse-order", `codex:{"skills":true,"agents":true,"plugin":"full"}`, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := modeFixture(t)
			modeRun(t, f, 0)
			f.put(f.codex, "skills/Local/keep", "local")
			modeRun(t, f, 0, "--mode", tc.mode)
			for _, p := range []string{"AGENTS.md", "resources/guide", "rules/default.rules"} {
				f.symlink(filepath.Join(f.codex, p))
			}
			if tc.standalone {
				f.symlink(filepath.Join(f.codex, "skills/Alpha"))
				f.symlink(filepath.Join(f.codex, "skills/Agent_Guidance"))
			} else {
				f.absent(filepath.Join(f.codex, "skills/Alpha"))
				f.absent(filepath.Join(f.codex, "skills/Agent_Guidance"))
			}
			if f.read(filepath.Join(f.codex, "skills/Local/keep")) != "local" {
				t.Fatal("unmanaged skill changed")
			}
			plugin := modePlugin(f)
			if !strings.Contains(f.read(filepath.Join(plugin, "skills/Alpha/SKILL.md")), "Example skill") {
				t.Fatal("plugin skill missing")
			}
			guidance := filepath.Join(plugin, "skills/Agent_Guidance/SKILL.md")
			if tc.guidance {
				if !strings.Contains(f.read(guidance), "Global agent instructions") {
					t.Fatal("plugin guidance missing")
				}
			} else {
				f.absent(guidance)
			}
			if f.json(filepath.Join(plugin, ".codex-plugin/plugin.json"))["version"] != "0.1.0" {
				t.Fatal("plugin manifest missing")
			}
			if _, err := os.Stat(filepath.Join(f.home, ".agents/plugins/marketplace.json")); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestModePluginOnlyLeavesGenerationAndCodexUntouched(t *testing.T) {
	for _, tc := range []struct {
		name, mode string
		guidance   bool
	}{
		{"default-guidance", `codex:{"plugin":"full"}`, true},
		{"explicit-agents-first", `codex:{"agents":false,"plugin":"full"}`, true},
		{"explicit-agents-last", `codex:{"plugin":"full","agents":false}`, true},
		{"explicit-guidance", `codex:{"skills":false,"agents":false,"plugin":"full"}`, true},
		{"exclude-guidance", `codex:{"plugin":"skills","agents":false}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := modeFixture(t)
			f.put(f.codex, "AGENTS.md", "private instructions")
			f.put(f.codex, "skills/Alpha/SKILL.md", "local skill")
			f.put(f.codex, "rules/default.rules", "local rule")
			modeRun(t, f, 0, "--mode", tc.mode)
			_ = f.latest()
			if f.read(filepath.Join(f.codex, "AGENTS.md")) != "private instructions" || f.read(filepath.Join(f.codex, "skills/Alpha/SKILL.md")) != "local skill" || f.read(filepath.Join(f.codex, "rules/default.rules")) != "local rule" {
				t.Fatal("Codex destination changed")
			}
			f.absent(filepath.Join(f.codex, "resources"))
			plugin := modePlugin(f)
			if !strings.Contains(f.read(filepath.Join(plugin, "skills/Alpha/SKILL.md")), "Example skill") {
				t.Fatal("plugin missing skill")
			}
			guidance := filepath.Join(plugin, "skills/Agent_Guidance/SKILL.md")
			if tc.guidance {
				if _, err := os.Stat(guidance); err != nil {
					t.Fatal(err)
				}
			} else {
				f.absent(guidance)
			}
		})
	}
}

func TestModeStandaloneSkillsLeavesAgentsUntouched(t *testing.T) {
	f := modeFixture(t)
	f.put(f.codex, "AGENTS.md", "private instructions")
	f.put(f.codex, "resources/guide/local", "private resource")
	f.put(f.codex, "rules/default.rules", "private rule")
	modeRun(t, f, 0, "--mode", `codex:{"agents":false,"skills":true}`)
	if f.read(filepath.Join(f.codex, "AGENTS.md")) != "private instructions" || f.read(filepath.Join(f.codex, "resources/guide/local")) != "private resource" || f.read(filepath.Join(f.codex, "rules/default.rules")) != "private rule" {
		t.Fatal("agent components changed")
	}
	f.symlink(filepath.Join(f.codex, "skills/Alpha"))
	f.symlink(filepath.Join(f.codex, "skills/Agent_Guidance"))
	f.absent(modePlugin(f))
}

func TestModeSkillsAndPluginWithoutAgents(t *testing.T) {
	f := modeFixture(t)
	f.put(f.codex, "AGENTS.md", "private instructions")
	f.put(f.codex, "resources/guide/local", "private resource")
	f.put(f.codex, "rules/default.rules", "private rule")
	modeRun(t, f, 0, "--mode", `codex:{"agents":false,"skills":true,"plugin":"skills"}`)
	if f.read(filepath.Join(f.codex, "AGENTS.md")) != "private instructions" || f.read(filepath.Join(f.codex, "resources/guide/local")) != "private resource" || f.read(filepath.Join(f.codex, "rules/default.rules")) != "private rule" {
		t.Fatal("agent components changed")
	}
	f.symlink(filepath.Join(f.codex, "skills/Alpha"))
	f.symlink(filepath.Join(f.codex, "skills/Agent_Guidance"))
	if _, err := os.Stat(filepath.Join(modePlugin(f), "skills/Alpha/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	f.absent(filepath.Join(modePlugin(f), "skills/Agent_Guidance"))
}

func TestModeAgentsOnlyPreservesSkills(t *testing.T) {
	f := modeFixture(t)
	modeRun(t, f, 0)
	old := f.resolved(filepath.Join(f.codex, "skills/Alpha"))
	f.put(f.codex, "skills/Local/keep", "local skill")
	f.put(f.layer, "AGENTS.md.d/10-base.md", "Updated instructions\n")
	modeRun(t, f, 0, "--mode", `codex:{"agents":true,"skills":false}`)
	if got := f.resolved(filepath.Join(f.codex, "skills/Alpha")); got != old {
		t.Fatalf("skill link changed from %s to %s", old, got)
	}
	if f.read(filepath.Join(f.codex, "skills/Local/keep")) != "local skill" {
		t.Fatal("unmanaged skill changed")
	}
	if !strings.Contains(f.read(filepath.Join(f.codex, "AGENTS.md")), "Updated instructions") {
		t.Fatal("instructions not updated")
	}
	f.symlink(filepath.Join(f.codex, "resources/guide"))
	f.symlink(filepath.Join(f.codex, "rules/default.rules"))
	f.absent(modePlugin(f))
}

func TestModeRejectsInvalidCombinationsBeforeSideEffects(t *testing.T) {
	for _, args := range [][]string{
		{"--mode", `other:{}`},
		{"--mode", `codex:{"publish":true}`},
		{"--mode", `codex:{"plugin":"publish"}`},
		{"--mode", `codex:{"extra":true}`},
		{"--mode", `codex:{"plugin":{}}`},
		{"--mode", `codex:{"plugin":"full","plugin":"skills"}`},
		{"--mode", `codex:{"agents":true,"agents":false}`},
		{"--mode", `codex:{"skills":true,"skills":false}`},
		{"--mode", `codex:{"plugin":null}`},
		{"--mode", `codex:{"plugin":true}`},
		{"--mode", `codex:{"plugin":"unknown"}`},
		{"--mode", `codex:{"plugin":"core"}`},
		{"--mode", `codex:{"agents":null}`},
		{"--mode", `codex:{"skills":"false"}`},
		{"--mode", `codex:{"agents":false,"skills":false}`},
		{"--mode", `codex:{"output":"plugin"}`},
		{"--mode", `codex:[]`},
		{"--mode", `codex:null`},
		{"--mode", `codex:{} trailing`},
		{"--mode", `codex:{}{}`},
		{"--mode", "codex", "--instructions-only"},
		{"--mode", `codex:{"plugin":"full"}`, "--skill-marketplace"},
		{"--mode", `codex:{"plugin":"skills","agents":false}`, "--rules-only"},
	} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			f := modeFixture(t)
			modeRun(t, f, 2, args...)
			f.absent(f.render)
			f.absent(f.codex)
			f.absent(f.home)
		})
	}
}

func TestModeLegacyMarketplaceFlag(t *testing.T) {
	f := modeFixture(t)
	modeRun(t, f, 0, "--skill-marketplace")
	f.symlink(filepath.Join(f.codex, "AGENTS.md"))
	f.symlink(filepath.Join(f.codex, "rules/default.rules"))
	f.absent(filepath.Join(f.codex, "skills/Alpha"))
	if _, err := os.Stat(filepath.Join(modePlugin(f), "skills/Alpha/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(modePlugin(f), "skills/Agent_Guidance/SKILL.md")); err != nil {
		t.Fatal(err)
	}
}

func TestModeGlobalOnRenderAndCheck(t *testing.T) {
	f := modeFixture(t)
	sections := filepath.Join(f.config, "sections")
	for _, args := range [][]string{
		{"--mode=codex", "render", "--sections-dir", sections, "--render-dir", f.render, "--codex-dir", f.codex},
		{"check", "--mode=codex:{}", "--sections-dir", sections, "--render-dir", f.render, "--codex-dir", f.codex},
		{"render", "--mode=codex:{}", "--sections-dir", sections, "--render-dir", f.render, "--codex-dir", f.codex},
		{"--mode=codex", "check", "--sections-dir", sections, "--render-dir", f.render, "--codex-dir", f.codex},
	} {
		out, err := exec.Command(contractBinary, args...).CombinedOutput()
		if err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
	}
}

func TestModePluginRejectedOutsideSync(t *testing.T) {
	for _, op := range []string{"prune", "reset"} {
		t.Run(op, func(t *testing.T) {
			f := modeFixture(t)
			mode := `--mode=codex:{"plugin":"full"}`
			if op == "check" {
				mode = `--mode=codex:{"agents":false,"plugin":"full"}`
			}
			args := []string{mode, op, "--sections-dir", filepath.Join(f.config, "sections"), "--render-dir", f.render, "--codex-dir", f.codex, "--home", f.home}
			out, err := exec.Command(contractBinary, args...).CombinedOutput()
			exit, ok := err.(*exec.ExitError)
			if !ok || exit.ExitCode() != 2 {
				t.Fatalf("%v: want exit 2, got %v\n%s", args, err, out)
			}
			f.absent(f.render)
			f.absent(f.codex)
			f.absent(f.home)
		})
	}
}

func TestModeFullPluginWithStandaloneSkillsLeavesAgentsUntouched(t *testing.T) {
	for _, mode := range []string{`codex:{"plugin":"full","skills":true}`, `codex:{"skills":true,"plugin":"full"}`} {
		t.Run(mode, func(t *testing.T) {
			f := modeFixture(t)
			f.put(f.codex, "AGENTS.md", "private instructions")
			f.put(f.codex, "rules/default.rules", "private rule")
			modeRun(t, f, 0, "--mode", mode)
			if f.read(filepath.Join(f.codex, "AGENTS.md")) != "private instructions" || f.read(filepath.Join(f.codex, "rules/default.rules")) != "private rule" {
				t.Fatal("agent components changed")
			}
			f.symlink(filepath.Join(f.codex, "skills/Alpha"))
			f.symlink(filepath.Join(f.codex, "skills/Agent_Guidance"))
			if _, err := os.Stat(filepath.Join(modePlugin(f), "skills/Agent_Guidance/SKILL.md")); err != nil {
				t.Fatal(err)
			}
		})
	}
}
