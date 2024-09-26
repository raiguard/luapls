package repl

import (
	"encoding/json"
	"fmt"

	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/parser"

	"github.com/chzyer/readline"
)

func Run() {
	rl, err := readline.New("(luapls) ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		fmt.Println("TOKENS:")
		tokens, _ := lexer.Run(line)
		for _, tok := range tokens {
			fmt.Println(tok.String())
		}

		fmt.Println("AST:")

		p := parser.New(line)
		file := p.ParseFile()
		bytes, _ := json.MarshalIndent(file, "", "  ")
		fmt.Println(string(bytes))
	}
}
