package harness_opencode

import (
	"encoding/json"

	"github.com/NonLogicalDev/conf.agents/utils/agentbld/internal/engine"
)

const Name = "opencode"

type Options struct {
	OpenCodeDir string `name:"opencode-dir" type:"path" help:"OpenCode destination directory (default: $OPENCODE_CONFIG_DIR or $HOME/.config/opencode)."`
}

func (o *Options) resolveDestination(any) (engine.Destination, error) {
	config, err := Destination(o.OpenCodeDir)
	if err != nil {
		return engine.Destination{}, err
	}
	return engine.Destination{Config: config, UsesConfig: true, Profile: map[string]string{"opencode_dir": config}}, nil
}

func (o *Options) loadProfile(fields map[string]json.RawMessage) error {
	return engine.ApplyProfilePath(fields, "opencode_dir", &o.OpenCodeDir)
}
