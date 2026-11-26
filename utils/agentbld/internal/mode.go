package internal

import (
	"fmt"
	"strings"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

type resolvedModes = map[string]any

// resolveModes preserves each selected harness and rejects duplicate names.
func resolveModes(values []string, parsers map[string]engine.ModeParser) (resolvedModes, error) {
	resolved := resolvedModes{}
	for _, value := range values {
		harness, _, _ := strings.Cut(value, ":")
		parser, ok := parsers[harness]
		if !ok {
			return nil, fmt.Errorf("unsupported harness %q", harness)
		}
		if _, exists := resolved[harness]; exists {
			return nil, fmt.Errorf("duplicate mode for harness %q", harness)
		}
		mode, err := parser(&value)
		if err != nil {
			return nil, err
		}
		resolved[harness] = mode
	}
	return resolved, nil
}

func parseModes(values []string, registry engine.Registry) (resolvedModes, error) {
	if len(values) == 0 {
		name := registry.DefaultName()
		if name == "" {
			return nil, fmt.Errorf("no default harness is configured")
		}
		mode, err := registry[name].Parse(nil)
		if err != nil {
			return nil, err
		}
		return resolvedModes{name: mode}, nil
	}
	parsers := map[string]engine.ModeParser{}
	for name, handler := range registry {
		parsers[name] = handler.Parse
	}
	return resolveModes(values, parsers)
}

func printModes(registry engine.Registry) {
	fmt.Println(`config: {"command":"modes"}`)
	fmt.Println(`Mode syntax:
  --mode=<harness>[:<JSON object>]

Available harnesses:`)
	names := engine.SortedKeys(registry)
	for _, name := range names {
		fmt.Printf("  %-10s %s\n", name, registry[name].Summary)
	}
	fmt.Printf("\nThe JSON object is optional. Repeat --mode once per harness. The default is %s.\n", registry.DefaultName())
	fmt.Println(`For sync, repeat --mode with bare harness names to activate only those saved
harnesses. Without a filter, sync activates all harnesses in the generation.`)
	for _, name := range names {
		fmt.Println("\n" + registry[name].ModeHelp)
	}
	fmt.Println("\nExamples (run from the repository root):")
	for _, name := range names {
		fmt.Println(registry[name].Examples)
	}
	fmt.Printf("  agentbld render --sync --sections-dir ./sections --mode=%s\n", strings.Join(names, " --mode="))
}
