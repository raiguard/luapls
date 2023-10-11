package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/raiguard/luapls/lsp"
	"github.com/raiguard/luapls/lua/ast"
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
		lexFile(args[2])
	case "lsp":
		lsp.Run()
	case "parse":
		parseFile(args[2], len(args) > 3 && args[3] == "json")
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

func parseFile(filename string, printJson bool) {
	before := time.Now()
	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	l := lexer.New(string(src))
	p := parser.New(l)
	file := p.ParseFile()
	if !printJson {
		fmt.Printf("Time taken: %s\n", time.Since(before))
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
	}
	if printJson {
		bytes, err := json.MarshalIndent(file, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bytes))
	} else {
		fmt.Println(file.String())
		ast.Walk(&file, func(n ast.Node) bool {
			fmt.Printf("%T: {%d, %d}\n", n, n.Pos(), n.End())
			return true
		})
	}
}
