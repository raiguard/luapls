package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Config struct {
	Environments map[string]EnvConfig `json:"environments"`
}

func (s *Server) getConfiguration(ctx *glsp.Context) {
	var res []Config
	ctx.Call(protocol.ServerWorkspaceConfiguration, protocol.ConfigurationParams{
		Items: []protocol.ConfigurationItem{{Section: ptr("luapls")}},
	}, &res)
	if res == nil || len(res) == 0 {
		return
	}
	s.config = res[0]

	s.log.Debugf("Configuration loaded: %s", toJSON(s.config))

	if s.config.Environments == nil || len(s.config.Environments) == 0 {
		s.config.Environments = map[string]EnvConfig{}
		s.log.Notice("No environments were specified, using root path")
		s.config.Environments["<fallback>"] = EnvConfig{
			Libraries: []string{},
			Roots:     []string{s.rootPath},
		}
	}

	for name, config := range s.config.Environments {
		s.envs[name] = newEnv(name, config)
		s.log.Debugf("Created environment %s", name)
	}
}
