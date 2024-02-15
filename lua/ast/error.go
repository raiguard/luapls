package ast

import "github.com/raiguard/luapls/lua/token"

type Error struct {
	Message string
	Range   token.Range
}
