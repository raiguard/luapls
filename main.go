package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/raiguard/luapls/lsp"
	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/token"
	"github.com/raiguard/luapls/repl"
	"github.com/tliron/kutil/util"
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
		level := int64(0)
		if len(args) > 2 {
			level, _ = strconv.ParseInt(args[2], 0, 8)
		}
		lsp.Run(int(level))
	case "parse":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Did not provide a filename")
			os.Exit(1)
		}
		parseFile(args[2])
	case "make-test":
		if len(args) < 4 {
			fmt.Fprintln(os.Stderr, "Not enough arguments: luapls make-test <suite> <label> <input string>")
			os.Exit(1)
		}
		makeTest(args[2], args[3], args[4])
	case "repl":
		repl.Run()
	case "check":
		checkFile(args[2])
	default:
		fmt.Fprintf(os.Stderr, "%s: unrecognized subcommand\n", task)
	}

	util.Exit(0)
}

func lexFile(filename string) {
	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	l := lexer.New(string(src))
	for {
		tok := l.Next()
		fmt.Println(tok.String())
		if tok.Type == token.EOF || tok.Type == token.INVALID {
			break
		}
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
	// units, _ := parser.Run(string(src))
	// bytes, err := json.MarshalIndent(units, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(bytes))
	p := parser.New(string(src))
	file := p.ParseFile()
	duration := time.Since(before)
	bytes, err := json.MarshalIndent(struct {
		Duration string
		File     ast.File
	}{
		Duration: duration.String(),
		File:     file,
	}, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(bytes))
}

type testSpec struct {
	Label  string
	Input  string
	AST    json.RawMessage
	Errors json.RawMessage `json:",omitempty"`
}

func makeTest(suite string, label string, input string) {
	p := parser.New(input)
	file := p.ParseFile()
	ast, _ := json.Marshal(&file.Block)
	errors, _ := json.Marshal(p.Errors())
	newSpec := testSpec{label, input, ast, errors}
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

func checkFile(path string) {
	fmt.Println("Unimplemented")
	// src, err := os.ReadFile(path)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }
	// p := parser.New(string(src))
	// file := p.ParseFile()
	// if len(file.Errors) > 0 {
	// 	fmt.Println("Parsing errors:")
	// 	for _, err := range file.Errors {
	// 		fmt.Println(&err)
	// 	}
	// 	return // TODO: Support partial type checking
	// }

	// env := types.NewLegacyEnvironment(&file)
	// env.ResolveTypes()
	// if len(env.Errors) > 0 {
	// 	fmt.Println("ERRORS:")
	// 	for _, err := range env.Errors {
	// 		fmt.Println(err.String())
	// 	}
	// 	os.Exit(1)
	// }
}
