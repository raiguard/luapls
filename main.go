package main

import (
	"encoding/json"
	"fmt"
	"os"
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
		if len(args) < 2 {
			fmt.Println("Did not provide a filename")
			os.Exit(1)
		}
		lexFile(args[2])
	case "lsp":
		lsp.Run()
	case "parse":
		if len(args) < 2 {
			fmt.Println("Did not provide a filename")
			os.Exit(1)
		}
		parseFile(args[2])
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
		fmt.Println(err)
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
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(bytes))
}
