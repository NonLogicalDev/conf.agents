package internal

import (
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSkillsPluginSkipsGuidanceBeforeValidation(t *testing.T) {
	cases := []struct{ name, body string }{
		{"invalid metadata", "not frontmatter\n"},
		{"missing guidance resource", "---\nname: Agent_Guidance\ndescription: Guidance\n---\n{{%_agents_%}}\n{{%_resources_%}}/missing.md\n"},
	}
	for _, test := range cases {
		for _, agents := range []bool{false, true} {
			name := test.name + " plugin only"
			if agents {
				name = test.name + " with agents"
			}
			t.Run(name, func(t *testing.T) {
				f := newContractFixture(t)
				f.put(f.layer, "skills/Agent_Guidance/SKILL.md", test.body)
				f.run("sync", 2, "--mode", `codex:{"plugin":"full"}`)
				f.run("sync", 2)
				config := `codex:{"plugin":"skills","agents":false}`
				if agents {
					config = `codex:{"plugin":"skills"}`
				}
				args := []string{"render", "--sync", "--config-dir", f.config, "--render-dir", f.render, "--codex-dir", f.codex, "--home", f.home, "--mode", config}
				if out, err := exec.Command(contractBinary, args...).CombinedOutput(); err != nil {
					t.Fatalf("skills sync failed: %v\n%s", err, out)
				}
				f.absent(filepath.Join(f.home, ".agents/plugins", marketplacePlugin, "skills/Agent_Guidance"))
				if got := f.read(filepath.Join(f.home, ".agents/plugins", marketplacePlugin, "skills/Alpha/SKILL.md")); got == "" {
					t.Fatal("other skill missing")
				}
				if agents {
					f.absent(filepath.Join(f.latest(), "skills/Agent_Guidance"))
				} else {
					_ = f.latest()
					f.absent(f.codex)
				}
			})
		}
	}
}

func TestSkillsPluginStillValidatesOtherSkills(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "skills/Alpha/SKILL.md", "invalid metadata\n")
	f.run("sync", 2, "--mode", `codex:{"plugin":"skills","agents":false}`)
	f.absent(filepath.Join(f.home, ".agents/plugins", marketplacePlugin))
}

func TestCheckSkillsGenerationWithoutChangingPlugin(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "skills/Agent_Guidance/SKILL.md", "---\nname: Agent_Guidance\ndescription: Guidance\n---\n{{%_agents_%}}\n")
	mode := `codex:{"plugin":"skills"}`
	f.run("sync", 0, "--mode", mode)
	plugin := filepath.Join(f.home, ".agents/plugins", marketplacePlugin)
	before, err := marketplaceInventory(plugin)
	if err != nil {
		t.Fatal(err)
	}
	catalog := filepath.Join(f.home, ".agents/plugins/marketplace.json")
	beforeCatalog := f.read(catalog)
	beforeMarker := f.read(filepath.Join(plugin, marketplaceMarker))
	f.run("check", 0, "--mode", mode)
	f.run("check", 1)
	after, err := marketplaceInventory(plugin)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, after) || f.read(catalog) != beforeCatalog || f.read(filepath.Join(plugin, marketplaceMarker)) != beforeMarker {
		t.Fatal("check changed plugin or catalog")
	}
	f.run("check", 1, "--mode", `codex:{"agents":false,"plugin":"skills"}`)
}

func TestPluginOnlySkipsUnselectedRules(t *testing.T) {
	for _, plugin := range []string{"full", "skills"} {
		t.Run(plugin, func(t *testing.T) {
			f := newContractFixture(t)
			f.put(f.layer, "rules/invalid.txt", "not a rule\n")
			f.run("sync", 2, "--skill-marketplace")
			f.absent(filepath.Join(f.home, ".agents/plugins", marketplacePlugin))
			f.run("sync", 0, "--mode", `codex:{"agents":false,"plugin":"`+plugin+`"}`)
			root := filepath.Join(f.home, ".agents/plugins", marketplacePlugin)
			f.absent(filepath.Join(root, "rules"))
			if f.read(filepath.Join(root, "skills/Alpha/SKILL.md")) == "" {
				t.Fatal("skill missing")
			}
			f.absent(f.codex)
			_ = f.latest()
		})
	}
}

func TestCheckFullPluginGenerationWithoutChangingPlugin(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "skills/Agent_Guidance/SKILL.md", "---\nname: Agent_Guidance\ndescription: Guidance\n---\n{{%_agents_%}}\n")
	mode := `codex:{"plugin":"full","agents":true}`
	f.run("sync", 0, "--mode", mode)
	plugin := filepath.Join(f.home, ".agents/plugins", marketplacePlugin)
	before, err := marketplaceInventory(plugin)
	if err != nil {
		t.Fatal(err)
	}
	catalog := filepath.Join(f.home, ".agents/plugins/marketplace.json")
	beforeCatalog := f.read(catalog)
	beforeMarker := f.read(filepath.Join(plugin, marketplaceMarker))
	f.run("check", 0, "--mode", mode)
	f.run("check", 1)
	after, err := marketplaceInventory(plugin)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, after) || f.read(catalog) != beforeCatalog || f.read(filepath.Join(plugin, marketplaceMarker)) != beforeMarker {
		t.Fatal("check changed plugin or catalog")
	}
	f.run("check", 1, "--mode", `codex:{"agents":false,"plugin":"full"}`)
}

func TestStandaloneSkillsAndPluginSkipUnselectedRules(t *testing.T) {
	f := newContractFixture(t)
	f.put(f.layer, "rules/invalid.txt", "not a rule\n")
	mode := `codex:{"agents":false,"skills":true,"plugin":"skills"}`
	f.run("sync", 0, "--mode", mode)
	f.symlink(filepath.Join(f.codex, "skills/Alpha"))
	f.absent(filepath.Join(f.codex, "AGENTS.md"))
	f.absent(filepath.Join(f.codex, "rules"))
	f.run("check", 0, "--mode", mode)
}
