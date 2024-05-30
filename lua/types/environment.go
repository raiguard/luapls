package types

import (
	"fmt"
	"strings"

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
	clear(c.Types)
	clear(c.Errors)

	c.resolveBlockTypes(&c.file.Block)
}

func (e *Environment) resolveBlockTypes(block *ast.Block) {
	for _, stmt := range block.Stmts {
		e.resolveStmtType(stmt)
	}
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
	case *ast.ForStatement:
		typ := e.resolveExprType(stmt.Start)
		if typ == nil {
			e.addType(stmt.Name, &Unknown{})
			return
		}
		if _, ok := typ.(*Number); !ok {
			e.addError(stmt.Start, "Range expressions must be of type 'number'")
			return
		}
		e.addType(stmt.Name, typ)
		finishTyp := e.resolveExprType(stmt.Finish)
		if typ != finishTyp {
			e.addError(stmt.Finish, "Range end must be of type '%s'", typ)
		}
		if stmt.Step != nil {
			stepTyp := e.resolveExprType(stmt.Step)
			if typ != stepTyp {
				e.addError(stmt.Step, "Range step must be of type '%s'", typ)
			}
		}
	case *ast.FunctionCall:
		e.resolveFunctionCall(stmt)
	case *ast.FunctionStatement:
		typ := &Function{
			Params: []NameAndType{},
			Return: nil,
		}
		comments := stmt.GetComments()
		for _, param := range stmt.Params {
			template := fmt.Sprintf("@param %s", param.Literal)
			defStart := strings.Index(comments, template)
			if defStart == -1 {
				e.addType(param, &Any{})
				continue
			}
			// This is awful
			def := strings.TrimSpace(strings.Split(comments[defStart+len(template):], "\n")[0])
			var paramTyp Type = &Unknown{}
			switch def {
			case "any":
				paramTyp = &Any{}
			case "boolean":
				paramTyp = &Boolean{}
			case "number":
				paramTyp = &Number{}
			case "string":
				paramTyp = &String{}
			}
			e.addType(param, paramTyp)
			typ.Params = append(typ.Params, NameAndType{
				Name: param.Literal,
				Type: paramTyp,
			})
			// TODO: Custom types
		}
		{
			template := "@return"
			defStart := strings.Index(comments, template)
			if defStart != -1 {
				// This is awful
				def := strings.TrimSpace(strings.Split(comments[defStart+len(template):], "\n")[0])
				switch def {
				case "any":
					typ.Return = &Any{}
				case "boolean":
					typ.Return = &Boolean{}
				case "number":
					typ.Return = &Number{}
				case "string":
					typ.Return = &String{}
				}
			}
		}
		e.addType(stmt, typ)
		// TODO: Parse index expression
		e.addType(stmt.Left, typ)
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
		typ := &Function{Params: []NameAndType{}}
		for _, param := range expr.Params {
			// TODO: Function parameter types - requires parsing doc comments
			typ.Params = append(typ.Params, NameAndType{Name: param.Literal, Type: &Unknown{}})
		}
		return e.addType(expr, typ)
	case *ast.FunctionCall:
		return e.resolveFunctionCall(expr)
	case *ast.IndexExpression:
		leftTyp := e.resolveExprType(expr.Left)
		if leftTyp == nil {
			return e.addType(expr, &Unknown{})
		}
		tbl, ok := leftTyp.(*Table)
		if !ok {
			e.addError(expr.Inner, "Attempting to index a non-table")
			return nil
		}

		var key string
		switch inner := expr.Inner.(type) {
		case *ast.Identifier:
			key = inner.Literal
		case *ast.StringLiteral:
			key = inner.Literal
		default:
			e.addError(inner, "Unimplemented")
			return nil
		}

		for i := 0; i < len(tbl.Fields); i++ {
			if tbl.Fields[i].Name == key {
				typ := tbl.Fields[i].Type
				e.addType(expr, typ)
				e.addType(expr.Inner, typ)
				return typ
			}
		}
		e.addError(expr.Inner, "Unknown field '%s'", key)
	case *ast.TableLiteral:
		typ := Table{
			Fields: []NameAndType{},
		}
		for _, fieldNode := range expr.Fields {
			ident, ok := fieldNode.Key.(*ast.Identifier)
			if !ok {
				// TODO:
				continue
			}
			valueTyp := e.resolveExprType(fieldNode.Value)
			field := NameAndType{
				Name: ident.Literal,
				Type: valueTyp,
			}
			typ.Fields = append(typ.Fields, field)
			e.addType(fieldNode, valueTyp)
			e.addType(ident, valueTyp)
		}
		e.addType(expr, &typ)
		return &typ
	}

	return nil
}

func (e *Environment) resolveFunctionCall(fc *ast.FunctionCall) Type {
	typ := e.resolveExprType(fc.Left)
	if typ == nil {
		return nil
	}
	e.addType(fc, typ)
	function, ok := typ.(*Function)
	if !ok {
		e.addError(fc, "'%s' is not a function", fc.Left)
		return typ
	}
	for i := 0; i < len(fc.Args); i++ {
		arg := fc.Args[i]
		argTyp := e.resolveExprType(arg)
		if i >= len(function.Params) {
			e.addError(arg, "Unused parameter")
			break // TODO:
		}
		if argTyp == nil {
			continue // TODO:
		}
		if argTyp != function.Params[i].Type {
			e.addError(arg, "Cannot use '%s' as '%s' in argument.", argTyp, function.Params[i].Type)
		}
		e.addType(arg, argTyp)
	}
	if len(fc.Args) < len(function.Params) {
		e.addError(fc, "Too few function parameters, expected %v, got %v", len(function.Params), len(fc.Args))
	}
	return function.Return
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

func (e *Environment) addError(node ast.Node, messageFmt string, messageArgs ...any) {
	e.Errors = append(e.Errors, parser.ParserError{
		Message: fmt.Sprintf(messageFmt, messageArgs...),
		Range:   ast.Range(node),
	})
}
