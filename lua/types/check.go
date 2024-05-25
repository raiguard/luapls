package types

import (
	"cmp"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/parser"
)

func min[T cmp.Ordered](a, b T) T {
	if cmp.Compare(a, b) <= 0 {
		return a
	}
	return b
}

type Checker struct {
	file  *parser.File
	Types map[ast.Node]Type
}

func NewChecker(file *parser.File) Checker {
	return Checker{
		file:  file,
		Types: map[ast.Node]Type{},
	}
}

func (t *Checker) Run() {
	ast.Walk(&t.file.Block, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.LocalStatement:
			for i := 0; i < len(node.Names); i++ {
				ident := node.Names[i]
				if i >= len(node.Exps) {
					t.Types[ident] = &Unknown{}
					continue
				}
				exp := node.Exps[i]
				switch exp.(type) {
				case *ast.BooleanLiteral:
					t.Types[ident] = &Boolean{}
				// case *ast.FunctionCall:
				// 	t.types[ident] = &Unknown{}
				case *ast.FunctionExpression:
					t.Types[ident] = &Function{Params: []Type{}}
				// case *ast.Identifier:
				// 	t.types[ident] = &Unknown{}
				// case *ast.IndexExpression:
				// 	t.types[ident] = &Unknown{}
				// case *ast.InfixExpression:
				// 	t.types[ident] = &Unknown{}
				// case *ast.Invalid:
				// 	t.types[ident] = &Unknown{}
				case *ast.NumberLiteral:
					t.Types[ident] = &Number{}
				// case *ast.PrefixExpression:
				// 	t.types[ident] = &Unknown{}
				case *ast.StringLiteral:
					t.Types[ident] = &String{}
				// case *ast.TableLiteral:
				// 	t.types[ident] = &Unknown{}
				// case *ast.Vararg:
				// 	t.types[ident] = &Unknown{}
				default:
					t.Types[ident] = &Unknown{}
				}
			}
		}
		return true
	})
}
