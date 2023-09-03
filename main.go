package main

import (
	"fmt"
	"os"

	"github.com/raiguard/luapls/lsp"
	"github.com/raiguard/luapls/repl"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("No argument provided, expected 'lsp' or 'repl'")
		os.Exit(1)
	}

	switch args[1] {
	case "lsp":
		lsp.Run()
	case "repl":
		repl.Run()
	}
}
