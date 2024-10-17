package util

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Ptr returns a pointer to the given value.
func Ptr[T any](value T) *T {
	return &value
}

// FileExists returns whether the given file exists on the filesystem.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// URIToPath returns a path from the given URI.
func URIToPath(uri protocol.URI) (string, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to parse file URI: %s", err))
	}
	return u.Path, nil
}

// PathToURI returns an absolute URI from a file path.
func PathToURI(path string) (protocol.URI, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to convert filepath to URI: %s", err))
	}
	return "file://" + abs, nil
}
