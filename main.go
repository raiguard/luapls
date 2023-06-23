package main

import (
	"fmt"
	"luapls/lua/lexer"
	"luapls/lua/token"
	"os"
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
