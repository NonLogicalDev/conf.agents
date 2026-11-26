package internal

import (
	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/harnesses/harness_claude"
	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/harnesses/harness_codex"
	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/harnesses/harness_opencode"
)

type harnessOptions struct {
	OpenCode harness_opencode.Options `embed:""`
	Claude   harness_claude.Options   `embed:""`
	Codex    harness_codex.Options    `embed:""`
}

func configuredHarnesses(a *cli) engine.Registry {
	return engine.Registry{
		harness_codex.Name:    harness_codex.Handler(&a.Harness.Codex),
		harness_claude.Name:   harness_claude.Handler(&a.Harness.Claude),
		harness_opencode.Name: harness_opencode.Handler(&a.Harness.OpenCode),
	}
}
