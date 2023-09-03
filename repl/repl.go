package repl

import (
	"fmt"
	"os"

	"github.com/raiguard/luapls/lua/lexer"
	"github.com/raiguard/luapls/lua/parser"

	"github.com/chzyer/readline"
)

func Run() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		// l := lexer.New(line)

		// for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		// 	fmt.Printf("%+v\n", tok)
		// }

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseBlock()
		for _, err := range p.Errors() {
			fmt.Fprintln(os.Stderr, err)
		}
		if len(p.Errors()) > 0 {
			continue
		}
		fmt.Println(program.String())
	}
}
