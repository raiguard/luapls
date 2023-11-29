package lsp

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func uriToPath(uri protocol.URI) (string, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to parse file URI: %s", err))
	}
	return u.Path, nil
}

func pathToURI(path string) (protocol.URI, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to convert filepath to URI: %s", err))
	}
	return "file://" + abs, nil
}
