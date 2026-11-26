package harness_codex

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

const Name = "codex"

type Options struct {
	CodexDir             string  `type:"path" help:"Codex destination directory (default: $CODEX_HOME or $HOME/.codex)."`
	CodexMarketplaceHome *string `type:"path" help:"Directory for marketplace.json and generated plugins (default: $HOME/.agents/plugins)."`
}

func (o *Options) resolveDestination(any) (engine.Destination, error) {
	config, err := codexPath(o.CodexDir)
	if err != nil {
		return engine.Destination{}, err
	}
	marketplace, err := marketplacePath(o)
	if err != nil {
		return engine.Destination{}, err
	}
	return engine.Destination{Config: config, Marketplace: marketplace, UsesConfig: true,
		Profile: map[string]string{"codex_dir": config, "codex_marketplace_home": marketplace}}, nil
}

func (o *Options) loadProfile(fields map[string]json.RawMessage) error {
	if err := engine.ApplyProfilePath(fields, "codex_dir", &o.CodexDir); err != nil {
		return err
	}
	if o.CodexMarketplaceHome != nil {
		return nil
	}
	var market string
	if err := engine.ApplyProfilePath(fields, "codex_marketplace_home", &market); err != nil {
		return err
	}
	if market != "" {
		o.CodexMarketplaceHome = &market
	}
	return nil
}

func configReport(saved any, target engine.Destination) (map[string]any, map[string]string) {
	c := saved.(Snapshot)
	return map[string]any{
		"sections_dirs": c.SectionsDirs, "render_dir": c.RenderDir,
		"codex_dir": target.Config, "codex_marketplace_home": target.Marketplace,
		"agents": c.Agents, "skills": c.Skills,
		"plugin": c.Plugin, "plugin_name": c.PluginName,
	}, nil
}

func marketplacePath(a *Options) (string, error) {
	path := ""
	if a.CodexMarketplaceHome != nil {
		path = *a.CodexMarketplaceHome
	}
	if path == "" {
		home, e := engine.DefaultHome()
		if e != nil {
			return "", e
		}
		path = filepath.Join(home, ".agents", "plugins")
	}
	return engine.Absolute(path)
}
func codexPath(path string) (string, error) {
	if path == "" {
		path = os.Getenv("CODEX_HOME")
		if path == "" {
			home, e := engine.DefaultHome()
			if e != nil {
				return "", e
			}
			path = filepath.Join(home, ".codex")
		}
	}
	return engine.Absolute(path)
}
