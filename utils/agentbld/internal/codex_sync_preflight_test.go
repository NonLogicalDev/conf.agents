package internal

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMarketplaceSyncPreflightsCodexConflicts(t *testing.T) {
	for _, component := range []string{"instructions", "resources", "rules"} {
		t.Run(component, func(t *testing.T) {
			f := newContractFixture(t)
			f.run("sync", 0, "--skill-marketplace")
			plugin := filepath.Join(f.home, ".agents/plugins", marketplacePlugin)
			catalog := filepath.Join(f.home, ".agents/plugins/marketplace.json")
			marker := filepath.Join(plugin, marketplaceMarker)
			inventory, err := marketplaceInventory(plugin)
			if err != nil {
				t.Fatal(err)
			}
			beforeCatalog := f.read(catalog)
			beforeMarker := f.read(marker)
			agents := filepath.Join(f.codex, "AGENTS.md")
			if component == "instructions" {
				if err := os.Remove(agents); err != nil {
					t.Fatal(err)
				}
				f.put(f.codex, "AGENTS.md", "user instructions\n")
			} else {
				name := "guide"
				if component == "rules" {
					name = "default.rules"
				}
				path := filepath.Join(f.codex, component, name)
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				f.put(f.codex, filepath.Join(component, name), "user content\n")
			}
			// Changing the source would replace the package if the Codex conflict were
			// discovered only after marketplace publication.
			f.put(f.layer, "skills/Alpha/SKILL.md", "---\nname: Alpha\ndescription: Changed\n---\nnew content\n")
			f.run("sync", 2, "--skill-marketplace")
			after, err := marketplaceInventory(plugin)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(inventory, after) || f.read(marker) != beforeMarker || f.read(catalog) != beforeCatalog {
				t.Fatal("Codex conflict changed the installed plugin or catalog")
			}
			switch component {
			case "instructions":
				if f.read(agents) != "user instructions\n" {
					t.Fatal("instructions were changed")
				}
			case "resources":
				if f.read(filepath.Join(f.codex, "resources/guide")) != "user content\n" {
					t.Fatal("resource conflict changed")
				}
			case "rules":
				if f.read(filepath.Join(f.codex, "rules/default.rules")) != "user content\n" {
					t.Fatal("rule conflict changed")
				}
			}
		})
	}
}
