package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestSpec struct {
	Label  string
	Input  string `json:"input"`
	AST    json.RawMessage
	Errors json.RawMessage `json:",omitempty"`
}

func TestParser(t *testing.T) {
	dir, err := os.ReadDir("test_specs")
	require.NoError(t, err)
	for _, entry := range dir {
		bytes, err := os.ReadFile(filepath.Join("test_specs", entry.Name()))
		var specs []TestSpec
		err = json.Unmarshal(bytes, &specs)
		if !assert.NoError(t, err) {
			continue
		}
		for _, spec := range specs {
			specLabel := strings.TrimSuffix(entry.Name(), ".json") + "/" + spec.Label
			t.Run(specLabel, func(t *testing.T) { testSpec(t, &spec) })
		}
	}
}

func testSpec(t *testing.T, spec *TestSpec) {
	p := New(spec.Input)
	file := p.ParseFile()
	ast, err := json.Marshal(&file.Block)
	if !assert.NoError(t, err) {
		return
	}
	assert.JSONEq(t, string(spec.AST), string(ast))

	if len(spec.Errors) == 0 {
		assert.Empty(t, p.Errors())
		return
	}

	errors, err := json.Marshal(&file.Errors)
	if !assert.NoError(t, err) {
		return
	}
	assert.JSONEq(t, string(spec.Errors), string(errors))
}
