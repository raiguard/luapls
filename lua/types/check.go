package types

import (
	"cmp"
	"fmt"

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
	file   *parser.File
	Types  map[ast.Node]Type
	Errors []parser.ParserError
}

func NewChecker(file *parser.File) Checker {
	return Checker{
		file:   file,
		Types:  map[ast.Node]Type{},
		Errors: []parser.ParserError{},
	}
}

func (c *Checker) Run() {
	ast.Walk(&c.file.Block, func(node ast.Node) bool {
		c.resolveType(node)
		return true
	})
}

func (c *Checker) resolveType(node ast.Node) Type {
	switch node := node.(type) {
	// Expressions
	case *ast.BooleanLiteral:
		return c.addType(node, &Boolean{})
	case *ast.FunctionExpression:
		typ := &Function{Params: []FunctionParameter{}}
		for _, param := range node.Params {
			// TODO: Function parameter types - requires parsing doc comments
			typ.Params = append(typ.Params, FunctionParameter{Name: param.Literal, Type: &Unknown{}})
		}
		return c.addType(node, typ)
	case *ast.NilLiteral:
		return c.addType(node, &Unknown{})
	case *ast.NumberLiteral:
		return c.addType(node, &Number{})
	case *ast.StringLiteral:
		return c.addType(node, &String{})

	case *ast.Identifier:
		def := c.FindDefinition(node, true)
		if def != nil {
			typ := c.Types[def]
			if typ != nil {
				return c.addType(node, typ)
			}
		} else {
			c.Errors = append(c.Errors, parser.ParserError{Message: fmt.Sprintf("Unknown variable '%s'", node.Literal), Range: ast.Range(node)})
		}

	// Statements
	case *ast.LocalStatement:
		for i := 0; i < len(node.Names); i++ {
			ident := node.Names[i]
			if i >= len(node.Exps) {
				c.addType(ident, &Unknown{})
				continue
			}
			exp := node.Exps[i]
			typ := c.resolveType(exp)
			if typ != nil {
				c.addType(ident, typ)
			} else {
				c.addType(ident, &Unknown{})
			}
		}
	}

	return nil
}

func (c *Checker) addType(node ast.Node, typ Type) Type {
	c.Types[node] = typ
	return typ
}

func (c *Checker) FindDefinition(identFor *ast.Identifier, includeSelf bool) *ast.Identifier {
	pos := identFor.StartPos
	var def *ast.Identifier
	ast.Walk(&c.file.Block, func(node ast.Node) bool {
		isAfter := node.Pos() > pos && pos < node.End()
		if isAfter {
			return false
		}
		isBefore := node.Pos() <= pos && pos > node.End()
		isInside := node.Pos() <= pos && pos < node.End()
		switch node := node.(type) {
		case *ast.ForInStatement:
			if isInside {
				for _, ident := range node.Names {
					if ident.Literal == identFor.Literal {
						def = ident
					}
				}
			}
		case *ast.ForStatement:
			if isInside {
				if node.Name != nil && node.Name.Literal == identFor.Literal {
					def = node.Name
				}
			}
		case *ast.FunctionExpression:
			if isInside {
				for _, ident := range node.Params {
					if ident.Literal == identFor.Literal {
						def = ident
					}
				}
			}
		case *ast.FunctionStatement:
			if isInside {
				for _, ident := range node.Params {
					if ident.Literal == identFor.Literal {
						def = ident
					}
				}
			}
			if isBefore || includeSelf {
				if ident, ok := node.Left.(*ast.Identifier); ok {
					if ident.Literal == identFor.Literal {
						def = ident
					}
				}
			}
		case *ast.LocalStatement:
			if isBefore || includeSelf {
				for _, ident := range node.Names {
					if ident.Literal == identFor.Literal {
						def = ident
					}
				}
			}
		default:
			return isInside
		}

		return true
	})
	return def
}
