package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
)

type cli struct {
	SectionsDirs         []string        `name:"sections-dir" type:"path" sep:"none" help:"Sections directory; repeat in precedence order."`
	ConfigDirs           []string        `name:"config-dir" type:"path" sep:"none" hidden:"" help:"Legacy repository root."`
	RenderDir            string          `type:"path" help:"Rendered generations directory (default: $XDG_STATE_HOME/agentbld or $HOME/.local/state/agentbld)."`
	OpenCodeDir          string          `name:"opencode-dir" type:"path" help:"OpenCode destination directory (default: $OPENCODE_CONFIG_DIR or $HOME/.config/opencode)."`
	ClaudeDir            string          `type:"path" help:"Claude Code destination directory (default: $HOME/.claude)."`
	ClaudeMarketplaceDir string          `type:"path" help:"Claude Code marketplace root (default: $XDG_DATA_HOME/agentbld/claude-marketplace or $HOME/.local/share/agentbld/claude-marketplace)."`
	CodexDir             string          `type:"path" help:"Codex destination directory (default: $CODEX_HOME or $HOME/.codex)."`
	CodexMarketplaceHome *string         `type:"path" help:"Directory for marketplace.json and generated plugins (default: $HOME/.agents/plugins)."`
	Home                 *string         `type:"path" hidden:"" help:"Legacy home directory; maps to <home>/.agents/plugins."`
	Mode                 []string        `name:"mode" sep:"none" help:"Harness mode; repeat once per harness. For sync, bare names select saved harnesses. Run 'agentbld modes' for options."`
	Render               renderFlags     `cmd:"" help:"Create or reuse an immutable generation."`
	Check                struct{}        `cmd:"" help:"Check whether the latest generation matches the source."`
	Sync                 generationFlags `cmd:"" help:"Activate an existing generation by number."`
	Prune                struct{}        `cmd:"" help:"Remove old, inactive generations."`
	Reset                struct{}        `cmd:"" help:"Back up and reset Codex links."`
	Modes                struct{}        `cmd:"" help:"Explain the available modes and their options."`
	operation            string
}
type syncFlags struct {
	SkillMarketplace bool `help:"Populate the personal skill marketplace instead of the skills directory."`
	InstructionsOnly bool `help:"Sync instructions and resources and remove managed skill links."`
	RulesOnly        bool `help:"Sync only managed Codex rule links."`
}
type renderFlags struct {
	Sync   bool      `help:"Activate the rendered generation."`
	Legacy syncFlags `embed:""`
}
type generationFlags struct {
	Number int `arg:"" help:"Generation ordinal to activate."`
}

func Main() int {
	var args cli
	parser := kong.Must(&args, kong.Name("agentbld"), kong.Description("Render read-only agent-config generations and sync them to selected harnesses."))
	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	args.operation = strings.Fields(ctx.Command())[0]
	if args.operation == "modes" {
		printModes()
		return 0
	}
	code, err := run(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return code
}

// defaultHome resolves the user's home alias so default paths stay consistent
// on systems where /home is a symlink. Explicit paths retain their spelling.
func defaultHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(home)
	if err != nil {
		if os.IsNotExist(err) {
			return home, nil
		}
		return "", err
	}
	return resolved, nil
}

func absolute(p string) (string, error) {
	if len(p) > 0 && p[0] == '~' && (len(p) == 1 || p[1] == '/') {
		h, e := defaultHome()
		if e != nil {
			return "", e
		}
		p = filepath.Join(h, p[1:])
	}
	return filepath.Abs(p)
}
func renderRoot(path string) (string, error) {
	if path == "" {
		state := os.Getenv("XDG_STATE_HOME")
		if state == "" {
			home, e := defaultHome()
			if e != nil {
				return "", e
			}
			state = filepath.Join(home, ".local", "state")
		}
		path = filepath.Join(state, "agentbld")
	}
	return absolute(path)
}
func marketplacePath(a cli) (string, error) {
	if a.CodexMarketplaceHome != nil && a.Home != nil {
		return "", fmt.Errorf("--codex-marketplace-home and --home cannot be combined")
	}
	path := ""
	if a.CodexMarketplaceHome != nil {
		path = *a.CodexMarketplaceHome
	}
	if a.Home != nil {
		home, e := absolute(*a.Home)
		if e != nil {
			return "", e
		}
		path = filepath.Join(home, ".agents", "plugins")
	}
	if path == "" {
		home, e := defaultHome()
		if e != nil {
			return "", e
		}
		path = filepath.Join(home, ".agents", "plugins")
	}
	return absolute(path)
}
func claudeMarketplacePath(path string) (string, error) {
	if path == "" {
		data := os.Getenv("XDG_DATA_HOME")
		if data == "" {
			home, err := defaultHome()
			if err != nil {
				return "", err
			}
			data = filepath.Join(home, ".local", "share")
		}
		path = filepath.Join(data, "agentbld", "claude-marketplace")
	}
	return absolute(path)
}
func codexPath(path string) (string, error) {
	if path == "" {
		path = os.Getenv("CODEX_HOME")
		if path == "" {
			home, e := defaultHome()
			if e != nil {
				return "", e
			}
			path = filepath.Join(home, ".codex")
		}
	}
	return absolute(path)
}
func printConfig(command string, s snapshot, number int, codex, marketplace string, a cli) error {
	destinations := map[string]string{}
	if _, ok := s["opencode"]; ok {
		p, err := openCodePath(a.OpenCodeDir)
		if err != nil {
			return err
		}
		destinations["opencode_dir"] = p
	}
	if _, ok := s["claude"]; ok {
		p, err := claudePath(a.ClaudeDir)
		if err != nil {
			return err
		}
		destinations["claude_dir"] = p
		if c, ok := s["claude"].(snapshotClaude); ok && c.Plugin != "none" && c.Plugin != "" {
			market, err := claudeMarketplacePath(a.ClaudeMarketplaceDir)
			if err != nil {
				return err
			}
			destinations["claude_marketplace_dir"] = market
		}
	}
	// Report the settings which will actually be applied, before any mutation.
	c, ok := s["codex"].(snapshotCodex)
	if !ok {
		b, err := json.Marshal(map[string]any{"command": command, "generation": number, "modes": s, "destinations": destinations})
		if err != nil {
			return err
		}
		fmt.Printf("config: %s\n", b)
		return nil
	}
	result := struct {
		Command      string            `json:"command"`
		Generation   int               `json:"generation,omitempty"`
		Sections     []string          `json:"sections_dirs"`
		RenderDir    string            `json:"render_dir"`
		CodexDir     string            `json:"codex_dir"`
		Marketplace  string            `json:"codex_marketplace_home"`
		Agents       bool              `json:"agents"`
		Skills       bool              `json:"skills"`
		Plugin       string            `json:"plugin"`
		PluginName   string            `json:"plugin_name"`
		Legacy       any               `json:"legacy"`
		Modes        snapshot          `json:"modes"`
		Destinations map[string]string `json:"destinations"`
	}{command, number, c.SectionsDirs, c.RenderDir, codex, marketplace, c.Agents && !c.Legacy.RulesOnly, c.Skills && !c.Legacy.InstructionsOnly && !c.Legacy.RulesOnly, c.Plugin, c.PluginName, c.Legacy, s, destinations}
	b, e := json.Marshal(result)
	if e != nil {
		return e
	}
	fmt.Printf("config: %s\n", b)
	return nil
}
func run(a cli) (int, error) {
	root, err := renderRoot(a.RenderDir)
	if err != nil {
		return 2, err
	}
	if a.operation == "sync" {
		return runNumberedSync(a, root)
	}
	if (a.operation == "prune" || a.operation == "reset") && len(a.Mode) > 0 {
		return 2, fmt.Errorf("--mode applies to render and check")
	}
	if len(a.SectionsDirs) > 0 && len(a.ConfigDirs) > 0 {
		return 2, fmt.Errorf("--sections-dir and --config-dir cannot be mixed")
	}
	if len(a.SectionsDirs) == 0 && len(a.ConfigDirs) == 0 {
		return 2, fmt.Errorf("at least one --sections-dir is required")
	}
	var sections, legacyRoots []string
	for _, config := range a.ConfigDirs {
		p, e := absolute(config)
		if e != nil {
			return 2, e
		}
		sections = append(sections, filepath.Join(p, "sections"))
		legacyRoots = append(legacyRoots, filepath.Join(filepath.Dir(p), "agent-config-rendered"))
	}
	for _, section := range a.SectionsDirs {
		p, e := absolute(section)
		if e != nil {
			return 2, e
		}
		sections = append(sections, p)
		if filepath.Base(p) == "sections" {
			legacyRoots = append(legacyRoots, filepath.Join(filepath.Dir(filepath.Dir(p)), "agent-config-rendered"))
		}
	}
	codex, e := codexPath(a.CodexDir)
	if e != nil {
		return 2, e
	}
	market, e := marketplacePath(a)
	if e != nil {
		return 2, e
	}
	flags := a.Render.Legacy
	if flags.InstructionsOnly && flags.SkillMarketplace {
		return 2, fmt.Errorf("--skill-marketplace and --instructions-only are mutually exclusive")
	}
	if flags.RulesOnly && (flags.SkillMarketplace || flags.InstructionsOnly) {
		return 2, fmt.Errorf("--rules-only requires no other sync mode")
	}
	modes, e := parseModes(a.Mode, flags)
	if e != nil {
		return 2, e
	}
	if _, ok := modes["codex"]; !ok && (flags.SkillMarketplace || flags.InstructionsOnly || flags.RulesOnly) {
		return 2, fmt.Errorf("legacy sync flags require the codex mode")
	}
	// Each resolved harness is retained for planning and snapshot creation.
	stub := snapshot{}
	registry := harnessRegistry()
	for _, name := range sortedKeys(modes) {
		handler, ok := registry[name]
		if !ok {
			return 2, fmt.Errorf("unsupported harness %q", name)
		}
		value, err := handler.stub(modes[name], flags, sections, root)
		if err != nil {
			return 2, err
		}
		stub[name] = value
	}
	if e = printConfig(a.operation, stub, 0, codex, market, a); e != nil {
		return 2, e
	}
	if a.operation == "reset" {
		for _, p := range sections {
			if len(a.SectionsDirs) > 0 && (isLink(p) || !isDir(p)) {
				return 2, fmt.Errorf("missing or invalid sections directory: %s", p)
			}
		}
		return 0, reset(codex, legacyRoots)
	}
	if a.operation == "prune" {
		// A missing root has nothing to prune; do not create it just for a lock.
		if _, err := os.Lstat(root); os.IsNotExist(err) {
			return 0, prune(root, codex)
		} else if err != nil {
			return 2, err
		}
		unlock, err := lockRenderRoot(root, false)
		if err != nil {
			return 2, err
		}
		defer unlock()
		return 0, prune(root, codex)
	}
	plans := map[string]*renderPlan{}
	snap := snapshot{}
	for _, name := range sortedKeys(modes) {
		handler, ok := registry[name]
		if !ok {
			return 2, fmt.Errorf("unsupported harness %q", name)
		}
		next, value, err := handler.build(modes[name], flags, sections, root)
		if err != nil {
			return 2, err
		}
		plans[name] = next
		snap[name] = value
	}
	if len(plans) == 0 {
		return 2, fmt.Errorf("no supported harness selected")
	}
	metadata, e := snapshotBytes(snap)
	if e != nil {
		return 2, e
	}
	if a.operation == "check" {
		gens, e := discoverGenerations(root)
		if e != nil {
			return 2, e
		}
		latest := ""
		if len(gens) > 0 {
			latest = gens[len(gens)-1].path
		}
		ok := latest != "" && renderedPlansMatch(latest, plans, metadata) && aliasMatches(root, latest)
		status := "out-of-date"
		if ok {
			status = "up-to-date"
		}
		summaryPlans(plans, latest, status)
		if !ok {
			return 1, nil
		}
		return 0, nil
	}
	unlock, e := lockRenderRoot(root, true)
	if e != nil {
		return 2, e
	}
	defer unlock()
	gen, e := performRenderPlans(root, plans, metadata)
	if e != nil {
		return 2, e
	}
	if a.Render.Sync {
		return 0, activateHarnesses(gen, root, snap, nil, a, codex, market)
	}
	return 0, nil
}
func runNumberedSync(a cli, root string) (int, error) {
	if a.Sync.Number <= 0 {
		return 2, fmt.Errorf("sync requires a positive generation number")
	}
	if len(a.SectionsDirs) > 0 || len(a.ConfigDirs) > 0 {
		return 2, fmt.Errorf("sync uses the generation's saved sources; do not supply --sections-dir or --config-dir")
	}
	unlock, err := lockRenderRoot(root, false)
	if err != nil {
		return 2, err
	}
	defer unlock()
	gens, e := discoverGenerations(root)
	if e != nil {
		return 2, e
	}
	gen := ""
	for _, g := range gens {
		if g.ordinal == a.Sync.Number {
			gen = g.path
			break
		}
	}
	if gen == "" {
		return 2, fmt.Errorf("generation number %d not found in %s", a.Sync.Number, root)
	}
	snap, e := readSnapshot(gen)
	if e != nil {
		return 2, e
	}
	selected := snapshot{}
	if len(a.Mode) == 0 {
		for name, value := range snap {
			selected[name] = value
		}
	} else {
		for _, name := range a.Mode {
			if strings.Contains(name, ":") {
				return 2, fmt.Errorf("sync mode filters must be bare harness names")
			}
			if _, ok := selected[name]; ok {
				return 2, fmt.Errorf("duplicate sync mode filter %q", name)
			}
			value, ok := snap[name]
			if !ok {
				return 2, fmt.Errorf("generation does not contain harness %q", name)
			}
			selected[name] = value
		}
	}
	codex, e := codexPath(a.CodexDir)
	if e != nil {
		return 2, e
	}
	market, e := marketplacePath(a)
	if e != nil {
		return 2, e
	}
	if e = printConfig("sync", selected, a.Sync.Number, codex, market, a); e != nil {
		return 2, e
	}
	return 0, activateHarnesses(gen, root, snap, selected, a, codex, market)
}
func activateHarnesses(gen, root string, s, selected snapshot, a cli, codex, marketplace string) error {
	if selected == nil {
		selected = s
	}
	destinations := runtimeDestinations{Codex: codex, Marketplace: marketplace}
	var err error
	if _, ok := selected["opencode"]; ok {
		destinations.OpenCode, err = openCodePath(a.OpenCodeDir)
		if err != nil {
			return err
		}
	}
	if _, ok := selected["claude"]; ok {
		destinations.Claude, err = claudePath(a.ClaudeDir)
		if err != nil {
			return err
		}
		if c, ok := selected["claude"].(snapshotClaude); ok && c.Plugin != "none" && c.Plugin != "" {
			destinations.ClaudeMarketplace, err = claudeMarketplacePath(a.ClaudeMarketplaceDir)
			if err != nil {
				return err
			}
		}
	}
	paths := map[string]string{}
	if _, ok := selected["codex"]; ok {
		paths["codex"] = destinations.Codex
	}
	if _, ok := selected["opencode"]; ok {
		paths["opencode"] = destinations.OpenCode
	}
	if c, ok := selected["claude"].(snapshotClaude); ok && (c.Agents || c.Skills) {
		paths["claude"] = destinations.Claude
	}
	names := sortedKeys(paths)
	for i, name := range names {
		for _, other := range names[i+1:] {
			a, b := paths[name], paths[other]
			if within(a, b) || within(b, a) {
				return fmt.Errorf("harness destinations overlap: %s (%s) and %s (%s)", name, a, other, b)
			}
		}
	}
	registry := harnessRegistry()
	plans, err := restoreHarnessPlans(gen, s)
	if err != nil {
		return err
	}
	var apply []func() error
	for _, name := range sortedKeys(selected) {
		handler, ok := registry[name]
		if !ok || handler.prepare == nil {
			return fmt.Errorf("unsupported harness in generation: %s", name)
		}
		commit, err := handler.prepare(gen, root, s, plans[name], destinations)
		if err != nil {
			return err
		}
		apply = append(apply, commit)
	}
	for _, commit := range apply {
		if err := commit(); err != nil {
			return err
		}
	}
	return nil
}
