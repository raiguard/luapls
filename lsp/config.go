package lsp

import (
	"encoding/json"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Config struct {
	Roots *[]string `json:"roots"`
}

func (s *Server) didChangeConfiguration(ctx *glsp.Context, params *protocol.DidChangeConfigurationParams) error {
	return s.updateConfig(params.Settings)
}

func (s *Server) updateConfig(settings any) error {
	// Because GLSP gives it to us as `any`, we have to re-marshal it to JSON then unmarshal it again.
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	s.config = config
	return nil
}
