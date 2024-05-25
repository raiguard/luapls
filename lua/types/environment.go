package types

import (
	"fmt"

	"github.com/raiguard/luapls/lua/ast"
	"github.com/raiguard/luapls/lua/parser"
)

type Environment struct {
	file   *parser.File
	Types  map[ast.Node]Type
	Errors []parser.ParserError
}

func NewEnvironment(file *parser.File) Environment {
	return Environment{
		file:   file,
		Types:  map[ast.Node]Type{},
		Errors: []parser.ParserError{},
	}
}

func (c *Environment) ResolveTypes() {
	ast.Walk(&c.file.Block, func(node ast.Node) bool {
		switch node := node.(type) {
		case ast.Expression:
			c.resolveExprType(node)
		case ast.Statement:
			c.resolveStmtType(node)
		}
		return true
	})
}

func (e *Environment) resolveStmtType(stmt ast.Statement) {
	switch stmt := stmt.(type) {
	case *ast.AssignmentStatement:
		for i := 0; i < len(stmt.Vars); i++ {
			leftVar := stmt.Vars[i]
			if i >= len(stmt.Exps) {
				e.addType(leftVar, &Unknown{})
				continue
			}
			exp := stmt.Exps[i]
			typ := e.resolveExprType(exp)
			leftTyp := e.resolveExprType(leftVar)
			if leftTyp != nil && typ != nil && leftTyp != typ {
				e.Errors = append(e.Errors, parser.ParserError{Message: fmt.Sprintf("Cannot assign '%s' to '%s'", typ, leftTyp), Range: ast.Range(leftVar)})
				e.addType(leftVar, leftTyp)
				continue
			}
			if typ != nil {
				e.addType(leftVar, typ)
			} else {
				e.addType(leftVar, &Unknown{})
			}
		}
	case *ast.LocalStatement:
		for i := 0; i < len(stmt.Names); i++ {
			ident := stmt.Names[i]
			if i >= len(stmt.Exps) {
				e.addType(ident, &Unknown{})
				continue
			}
			exp := stmt.Exps[i]
			typ := e.resolveExprType(exp)
			leftTyp := e.resolveExprType(ident)
			if leftTyp != nil && typ != nil && leftTyp != typ {
				e.Errors = append(e.Errors, parser.ParserError{Message: fmt.Sprintf("Cannot assign '%s' to '%s'", typ, leftTyp), Range: ast.Range(ident)})
				e.addType(ident, leftTyp)
				continue
			}
			if typ != nil {
				e.addType(ident, typ)
			} else {
				e.addType(ident, &Unknown{})
			}
		}
	}
}

func (e *Environment) resolveExprType(expr ast.Expression) Type {
	switch expr := expr.(type) {
	// Literals
	case *ast.BooleanLiteral:
		return e.addType(expr, &Boolean{})
	case *ast.NilLiteral:
		return e.addType(expr, &Unknown{})
	case *ast.NumberLiteral:
		return e.addType(expr, &Number{})
	case *ast.StringLiteral:
		return e.addType(expr, &String{})

	case *ast.Identifier:
		def := e.FindDefinition(expr, true)
		if def != nil {
			typ := e.Types[def]
			if typ != nil {
				return e.addType(expr, typ)
			}
		} else {
			e.Errors = append(e.Errors, parser.ParserError{Message: fmt.Sprintf("Unknown variable '%s'", expr.Literal), Range: ast.Range(expr)})
		}

	case *ast.FunctionExpression:
		typ := &Function{Params: []FunctionParameter{}}
		for _, param := range expr.Params {
			// TODO: Function parameter types - requires parsing doc comments
			typ.Params = append(typ.Params, FunctionParameter{Name: param.Literal, Type: &Unknown{}})
		}
		return e.addType(expr, typ)
	}

	return nil
}

func (e *Environment) addType(node ast.Node, typ Type) Type {
	e.Types[node] = typ
	return typ
}

func (e *Environment) FindDefinition(identFor *ast.Identifier, includeSelf bool) *ast.Identifier {
	pos := identFor.StartPos
	var def *ast.Identifier
	ast.Walk(&e.file.Block, func(node ast.Node) bool {
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
