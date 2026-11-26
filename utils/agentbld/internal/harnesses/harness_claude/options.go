package harness_claude

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

const Name = "claude"

type Options struct {
	ClaudeDir            string `type:"path" help:"Claude Code destination directory (default: $HOME/.claude)."`
	ClaudeMarketplaceDir string `type:"path" help:"Claude Code marketplace root (default: $XDG_DATA_HOME/agentbld/claude-marketplace or $HOME/.local/share/agentbld/claude-marketplace)."`
}

func (o *Options) resolveDestination(value any) (engine.Destination, error) {
	config, err := Destination(o.ClaudeDir)
	if err != nil {
		return engine.Destination{}, err
	}
	saved := value.(Snapshot)
	dest := engine.Destination{Config: config, UsesConfig: saved.Agents || saved.Skills,
		Profile: map[string]string{"claude_dir": config}}
	if saved.Plugin != "none" {
		dest.Marketplace, err = claudeMarketplacePath(o.ClaudeMarketplaceDir)
		if err != nil {
			return engine.Destination{}, err
		}
		dest.Profile["claude_marketplace_dir"] = dest.Marketplace
	}
	return dest, nil
}

func (o *Options) loadProfile(fields map[string]json.RawMessage) error {
	if err := engine.ApplyProfilePath(fields, "claude_dir", &o.ClaudeDir); err != nil {
		return err
	}
	return engine.ApplyProfilePath(fields, "claude_marketplace_dir", &o.ClaudeMarketplaceDir)
}

func claudeMarketplacePath(path string) (string, error) {
	if path == "" {
		data := os.Getenv("XDG_DATA_HOME")
		if data == "" {
			home, err := engine.DefaultHome()
			if err != nil {
				return "", err
			}
			data = filepath.Join(home, ".local", "share")
		}
		path = filepath.Join(data, "agentbld", "claude-marketplace")
	}
	return engine.Absolute(path)
}
