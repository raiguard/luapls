package main

import (
	"fmt"
	"os"

	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/token"
)

func main() {
	input, _ := os.ReadFile(os.Args[1])
	l := lexer.New(string(input))
	for {
		tok := l.NextToken()
		if tok.Type == token.EOF {
			break
		}
		fmt.Printf("%v\n", tok)
	}
}
