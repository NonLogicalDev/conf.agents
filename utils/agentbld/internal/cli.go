package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

type cli struct {
	Harness      harnessOptions  `embed:""`
	SectionsDirs []string        `name:"sections-dir" type:"path" sep:"none" help:"Sections directory; repeat in precedence order."`
	RenderDir    string          `type:"path" help:"Rendered generations directory (default: $XDG_STATE_HOME/agentbld or $HOME/.local/state/agentbld)."`
	Mode         []string        `name:"mode" sep:"none" help:"Harness mode; repeat once per harness. For sync, bare names select saved harnesses. Run 'agentbld modes' for options."`
	Render       renderFlags     `cmd:"" help:"Create or reuse an immutable generation."`
	Check        struct{}        `cmd:"" help:"Check whether the latest generation matches the source."`
	Sync         generationFlags `cmd:"" help:"Activate an existing generation by number."`
	Profile      profileCommands `cmd:"" help:"Inspect saved sync profiles."`
	Prune        struct{}        `cmd:"" help:"Remove old, inactive generations."`
	Reset        struct{}        `cmd:"" help:"Back up configuration and reset managed links."`
	Modes        struct{}        `cmd:"" help:"Explain the available modes and their options."`
	operation    string
}
type renderFlags struct {
	Sync bool `help:"Activate the rendered generation."`
}
type generationFlags struct {
	Number      int     `arg:"" help:"Generation ordinal to activate."`
	Profile     *string `placeholder:"NAME" help:"Load sync settings from the named profile."`
	SaveProfile *string `placeholder:"NAME" help:"Save resolved sync settings to a named JSON profile after activation."`
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
	if args.operation == "profile" {
		args.operation += " " + ctx.Selected().Name
	}
	if args.operation == "modes" {
		printModes(configuredHarnesses(&args))
		return 0
	}
	code, err := run(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	return code
}

func renderRoot(path string) (string, error) {
	if path == "" {
		state := os.Getenv("XDG_STATE_HOME")
		if state == "" {
			home, e := engine.DefaultHome()
			if e != nil {
				return "", e
			}
			state = filepath.Join(home, ".local", "state")
		}
		path = filepath.Join(state, "agentbld")
	}
	return engine.Absolute(path)
}
func printConfig(command string, registry engine.Registry, s engine.Snapshot, number int, resolved map[string]engine.Destination) error {
	fields := map[string]any{"command": command, "generation": number, "modes": s}
	destinations := map[string]string{}
	for _, name := range engine.SortedKeys(s) {
		handler := registry[name]
		if handler.ConfigReport != nil {
			values, paths := handler.ConfigReport(s[name], resolved[name])
			for key, value := range values {
				fields[key] = value
			}
			for key, path := range paths {
				destinations[key] = path
			}
		}
	}
	fields["destinations"] = destinations
	data, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	fmt.Printf("config: %s\n", data)
	return nil
}

func run(a cli) (int, error) {
	switch a.operation {
	case "profile list":
		return 0, listProfiles()
	case "profile show":
		return 0, showProfile(a.Profile.Show.Name)
	}
	registry := configuredHarnesses(&a)
	if a.operation == "sync" && a.Sync.Profile != nil {
		if err := loadSyncProfile(&a, registry, *a.Sync.Profile); err != nil {
			return 2, err
		}
	}
	root, err := renderRoot(a.RenderDir)
	if err != nil {
		return 2, err
	}
	if a.operation == "sync" {
		return runNumberedSync(a, registry, root)
	}
	if (a.operation == "prune" || a.operation == "reset") && len(a.Mode) > 0 {
		return 2, fmt.Errorf("--mode applies to render and check")
	}
	if len(a.SectionsDirs) == 0 && (a.operation == "render" || a.operation == "check") {
		return 2, fmt.Errorf("at least one --sections-dir is required")
	}
	var sections []string
	for _, section := range a.SectionsDirs {
		p, e := engine.Absolute(section)
		if e != nil {
			return 2, e
		}
		sections = append(sections, p)
	}
	modes, e := parseModes(a.Mode, registry)
	if e != nil {
		return 2, e
	}
	stub := engine.Snapshot{}
	for _, name := range engine.SortedKeys(modes) {
		handler, ok := registry[name]
		if !ok {
			return 2, fmt.Errorf("unsupported harness %q", name)
		}
		value, err := handler.Stub(modes[name], sections, root)
		if err != nil {
			return 2, err
		}
		stub[name] = value
	}
	destinations, e := resolveSyncDestinations(registry, stub)
	if e != nil {
		return 2, e
	}
	if e = printConfig(a.operation, registry, stub, 0, destinations); e != nil {
		return 2, e
	}
	if a.operation == "reset" {
		for _, p := range sections {
			if len(a.SectionsDirs) > 0 && (engine.IsLink(p) || !engine.IsDir(p)) {
				return 2, fmt.Errorf("missing or invalid sections directory: %s", p)
			}
		}
		for _, name := range engine.SortedKeys(stub) {
			if reset := registry[name].Reset; reset != nil {
				if err := reset(destinations[name], root); err != nil {
					return 2, err
				}
			}
		}
		return 0, nil
	}
	if a.operation == "prune" {
		// A missing root has nothing to prune; do not create it just for a lock.
		if _, err := os.Lstat(root); os.IsNotExist(err) {
			return 0, registry.Prune(root)
		} else if err != nil {
			return 2, err
		}
		unlock, err := engine.LockRenderRoot(root, false)
		if err != nil {
			return 2, err
		}
		defer unlock()
		return 0, registry.Prune(root)
	}
	plans := map[string]*engine.RenderPlan{}
	snap := engine.Snapshot{}
	for _, name := range engine.SortedKeys(modes) {
		handler, ok := registry[name]
		if !ok {
			return 2, fmt.Errorf("unsupported harness %q", name)
		}
		next, value, err := handler.Build(registry, modes[name], sections, root)
		if err != nil {
			return 2, err
		}
		plans[name] = next
		snap[name] = value
	}
	if len(plans) == 0 {
		return 2, fmt.Errorf("no supported harness selected")
	}
	metadata, e := engine.SnapshotBytes(snap)
	if e != nil {
		return 2, e
	}
	if a.operation == "check" {
		gens, e := registry.DiscoverGenerations(root)
		if e != nil {
			return 2, e
		}
		latest := ""
		if len(gens) > 0 {
			latest = gens[len(gens)-1].Path
		}
		ok := latest != "" && engine.RenderedPlansMatch(latest, plans, metadata) && engine.AliasMatches(root, latest)
		status := "out-of-date"
		if ok {
			status = "up-to-date"
		}
		engine.SummaryPlans(plans, latest, status)
		if !ok {
			return 1, nil
		}
		return 0, nil
	}
	unlock, e := engine.LockRenderRoot(root, true)
	if e != nil {
		return 2, e
	}
	defer unlock()
	gen, e := registry.PerformRenderPlans(root, plans, metadata)
	if e != nil {
		return 2, e
	}
	if a.Render.Sync {
		destinations, err := resolveSyncDestinations(registry, snap)
		if err != nil {
			return 2, err
		}
		return 0, activateHarnesses(registry, gen, root, snap, snap, destinations)
	}
	return 0, nil
}
func runNumberedSync(a cli, registry engine.Registry, root string) (int, error) {
	if a.Sync.Number <= 0 {
		return 2, fmt.Errorf("sync requires a positive generation number")
	}
	if len(a.SectionsDirs) > 0 {
		return 2, fmt.Errorf("sync uses the generation's saved sources; do not supply --sections-dir")
	}
	var profilePath string
	if a.Sync.SaveProfile != nil {
		path, err := profileFile(*a.Sync.SaveProfile)
		if err != nil {
			return 2, err
		}
		profilePath = path
	}
	unlock, err := engine.LockRenderRoot(root, false)
	if err != nil {
		return 2, err
	}
	defer unlock()
	gens, e := registry.DiscoverGenerations(root)
	if e != nil {
		return 2, e
	}
	gen := ""
	for _, g := range gens {
		if g.Ordinal == a.Sync.Number {
			gen = g.Path
			break
		}
	}
	if gen == "" {
		return 2, fmt.Errorf("generation number %d not found in %s", a.Sync.Number, root)
	}
	snap, e := registry.ReadSnapshot(gen)
	if e != nil {
		return 2, e
	}
	selected := engine.Snapshot{}
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
	destinations, err := resolveSyncDestinations(registry, selected)
	if err != nil {
		return 2, err
	}
	if e = printConfig("sync", registry, selected, a.Sync.Number, destinations); e != nil {
		return 2, e
	}
	if err := activateHarnesses(registry, gen, root, snap, selected, destinations); err != nil {
		return 2, err
	}
	if profilePath != "" {
		if err := saveSyncProfile(profilePath, resolvedSyncProfile(root, selected, destinations)); err != nil {
			return 2, fmt.Errorf("sync succeeded but saving profile %s failed: %w", profilePath, err)
		}
		fmt.Printf("saved profile: %s\n", profilePath)
	}
	return 0, nil
}
func resolveSyncDestinations(registry engine.Registry, selected engine.Snapshot) (map[string]engine.Destination, error) {
	destinations := map[string]engine.Destination{}
	for _, name := range engine.SortedKeys(selected) {
		if resolve := registry[name].ResolveDestination; resolve != nil {
			dest, err := resolve(selected[name])
			if err != nil {
				return nil, err
			}
			destinations[name] = dest
		}
	}
	return destinations, nil
}

func activateHarnesses(registry engine.Registry, gen, root string, s, selected engine.Snapshot, destinations map[string]engine.Destination) error {
	paths := map[string]string{}
	for name, dest := range destinations {
		if dest.UsesConfig {
			paths[name] = dest.Config
		}
	}
	names := engine.SortedKeys(paths)
	for i, name := range names {
		for _, other := range names[i+1:] {
			a, b := paths[name], paths[other]
			if engine.Within(a, b) || engine.Within(b, a) {
				return fmt.Errorf("harness destinations overlap: %s (%s) and %s (%s)", name, a, other, b)
			}
		}
	}
	plans, err := registry.RestoreHarnessPlans(gen, s)
	if err != nil {
		return err
	}
	var apply []func() error
	for _, name := range engine.SortedKeys(selected) {
		handler, ok := registry[name]
		if !ok || handler.Prepare == nil {
			return fmt.Errorf("unsupported harness in generation: %s", name)
		}
		commit, err := handler.Prepare(registry, gen, root, s, plans[name], destinations[name])
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
