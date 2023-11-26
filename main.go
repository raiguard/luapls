package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/raiguard/luapls/lsp"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/token"
	"github.com/raiguard/luapls/repl"
)

func main() {
	args := os.Args
	task := "lsp"
	if len(args) > 1 {
		task = args[1]
	}

	switch task {
	case "lex":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Did not provide a filename")
			os.Exit(1)
		}
		lexFile(args[2])
	case "lsp":
		lsp.Run()
	case "parse":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Did not provide a filename")
			os.Exit(1)
		}
		parseFile(args[2])
	case "make-test":
		if len(args) < 4 {
			fmt.Fprintln(os.Stderr, "Not enough arguments, requires a label and input string")
			os.Exit(1)
		}
		makeTest(args[2], args[3], args[4])
	case "repl":
		repl.Run()
	}
}

func lexFile(filename string) {
	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	l := lexer.New(string(src))
	for {
		tok := l.NextToken()
		if tok.Type == token.EOF || tok.Type == token.INVALID {
			break
		}
		fmt.Println(tok.String())
	}
	fmt.Println(l.GetLineBreaks())
}

func parseFile(filename string) {
	before := time.Now()
	src, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	p := parser.New(string(src))
	file := p.ParseFile()
	bytes, err := json.MarshalIndent(struct {
		Duration string
		File     parser.File
	}{
		Duration: time.Since(before).String(),
		File:     file,
	}, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(bytes))
}

type testSpec struct {
	Label string
	Input string
	AST   json.RawMessage
}

func makeTest(suite string, label string, input string) {
	block := parseBlockToJSON(input)
	newSpec := testSpec{label, input, block}
	path := filepath.Join("lua/parser/test_specs", suite+".json")
	specs := readOrMakeSpecs(path)
	specs = append(specs, newSpec)
	output, _ := json.MarshalIndent(specs, "", "  ")
	os.WriteFile(path, output, 0644)
}

func readOrMakeSpecs(path string) []testSpec {
	specs := []testSpec{}
	file, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(file, &specs)
	}
	return specs
}

func parseBlockToJSON(input string) []byte {
	p := parser.New(input)
	block := p.ParseBlock()
	bytes, err := json.MarshalIndent(&block, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return bytes
}
