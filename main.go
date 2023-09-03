package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raiguard/luapls/lsp"
	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/parser"
	"github.com/raiguard/luapls/lua/token"
	"github.com/raiguard/luapls/repl"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("No argument provided, expected 'lsp' or 'repl'")
		os.Exit(1)
	}

	switch args[1] {
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
}

func parseFile(filename string, printJson bool) {
	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	l := lexer.New(string(src))
	p := parser.New(l)
	block := p.ParseBlock()
	for _, err := range p.Errors() {
		fmt.Println(err)
	}
	if len(p.Errors()) > 0 {
		return
	}
	if printJson {
		bytes, err := json.MarshalIndent(block, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bytes))
	} else {
		fmt.Println(block.String())
	}
}
