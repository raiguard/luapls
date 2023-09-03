package main

import (
	"fmt"
	"os"

	"github.com/raiguard/luapls/lsp"
	"github.com/raiguard/luapls/lua/lexer"
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
