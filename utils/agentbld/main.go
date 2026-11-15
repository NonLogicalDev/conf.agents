package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
)

type cli struct {
	Operation        string `arg:"" enum:"render,check,sync,prune,reset" help:"Operation to perform."`
	ConfigDir        string `required:"" type:"path" help:"Agent configuration source directory."`
	RenderDir        string `required:"" type:"path" help:"Rendered generations directory."`
	CodexDir         string `type:"path" help:"Codex destination directory."`
	Home             string `type:"path" help:"Marketplace destination home."`
	SkillMarketplace bool   `help:"Populate the personal skill marketplace instead of the skills directory during sync."`
	InstructionsOnly bool   `help:"Sync instructions and resources and remove managed skill links."`
	RulesOnly        bool   `help:"Sync only managed Codex rule links."`
}

func main() {
	var args cli
	parser := kong.Must(&args, kong.Name("agentbld"), kong.Description("Render read-only agent-config generations and link them into Codex."))
	_, err := parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	code, err := run(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	os.Exit(code)
}

func absolute(p string) (string, error) {
	if len(p) > 0 && p[0] == '~' && (len(p) == 1 || p[1] == '/') {
		h, e := os.UserHomeDir()
		if e != nil {
			return "", e
		}
		p = filepath.Join(h, p[1:])
	}
	return filepath.Abs(p)
}
func run(a cli) (int, error) {
	var err error
	if a.ConfigDir, err = absolute(a.ConfigDir); err != nil {
		return 2, err
	}
	if a.RenderDir, err = absolute(a.RenderDir); err != nil {
		return 2, err
	}
	if a.CodexDir != "" {
		if a.CodexDir, err = absolute(a.CodexDir); err != nil {
			return 2, err
		}
	}
	if a.Home == "" {
		a.Home, err = os.UserHomeDir()
	} else {
		a.Home, err = absolute(a.Home)
	}
	if err != nil {
		return 2, err
	}
	if a.SkillMarketplace && a.Operation != "sync" {
		return 2, fmt.Errorf("--skill-marketplace requires sync")
	}
	if a.InstructionsOnly && a.Operation != "sync" {
		return 2, fmt.Errorf("--instructions-only requires sync")
	}
	if a.InstructionsOnly && a.SkillMarketplace {
		return 2, fmt.Errorf("--skill-marketplace and --instructions-only are mutually exclusive")
	}
	if a.RulesOnly && (a.Operation != "sync" || a.SkillMarketplace || a.InstructionsOnly) {
		return 2, fmt.Errorf("--rules-only requires sync without another sync mode")
	}
	if a.Operation == "reset" || a.Operation == "prune" {
		if a.CodexDir == "" {
			return 2, fmt.Errorf("%s requires --codex-dir", a.Operation)
		}
		if a.Operation == "reset" {
			return 0, reset(a.CodexDir, filepath.Join(filepath.Dir(a.ConfigDir), "agent-config-rendered"))
		}
		return 0, prune(a.RenderDir, a.CodexDir)
	}
	plan, err := buildPlan(a.ConfigDir, a.CodexDir)
	if err != nil {
		return 2, err
	}
	if a.Operation == "check" {
		gens, e := discoverGenerations(a.RenderDir)
		if e != nil {
			return 2, e
		}
		latest := ""
		if len(gens) > 0 {
			latest = gens[len(gens)-1].path
		}
		ok := latest != "" && renderedMatches(latest, plan) && aliasMatches(a.RenderDir, latest)
		status := "out-of-date"
		if ok {
			status = "up-to-date"
		}
		summary(plan, latest, status)
		if !ok {
			return 1, nil
		}
		return 0, nil
	}
	gen, err := performRender(a.RenderDir, plan)
	if err != nil {
		return 2, err
	}
	if a.Operation == "sync" {
		if a.CodexDir == "" {
			return 2, fmt.Errorf("sync requires --codex-dir")
		}
		if a.RulesOnly {
			return 0, syncRules(a.CodexDir, gen, a.RenderDir)
		}
		if a.SkillMarketplace {
			result, e := populateMarketplace(plan, a.Home)
			if e != nil {
				return 2, e
			}
			out, e := json.MarshalIndent(result, "", "  ")
			if e != nil {
				return 2, e
			}
			fmt.Println(string(out))
		}
		return 0, syncAll(a.CodexDir, gen, a.RenderDir, plan, !a.SkillMarketplace && !a.InstructionsOnly, a.InstructionsOnly)
	}
	return 0, nil
}
