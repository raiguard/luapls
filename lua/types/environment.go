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
	Nodes  []ast.Node
}

func NewEnvironment(file *parser.File) Environment {
	return Environment{
		file:   file,
		Types:  map[ast.Node]Type{},
		Errors: []parser.ParserError{},
		Nodes:  []ast.Node{},
	}
}

func (c *Environment) ResolveTypes() {
	clear(c.Types)
	clear(c.Errors)
	c.Nodes = []ast.Node{}

	c.resolveBlockTypes(&c.file.Block)
}

func (e *Environment) resolveBlockTypes(block *ast.Block) {
	e.pushNode(block)
	defer e.popNode()
	for _, stmt := range block.Stmts {
		e.resolveStmtType(stmt)
	}
}

func (e *Environment) resolveStmtType(stmt ast.Statement) {
	e.pushNode(stmt)
	defer e.popNode()
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
		e.resolveFunctionCallType(stmt)
	case *ast.FunctionStatement:
		// typ := &Function{
		// 	Params: []NameAndType{},
		// 	Return: nil,
		// }
		// comments := stmt.GetComments()
		// for _, param := range stmt.Params {
		// 	template := fmt.Sprintf("@param %s", param.Literal)
		// 	defStart := strings.Index(comments, template)
		// 	if defStart == -1 {
		// 		e.addType(param, &Any{})
		// 		continue
		// 	}
		// 	// This is awful
		// 	def := strings.TrimSpace(strings.Split(comments[defStart+len(template):], "\n")[0])
		// 	var paramTyp Type = &Unknown{}
		// 	switch def {
		// 	case "any":
		// 		paramTyp = &Any{}
		// 	case "boolean":
		// 		paramTyp = &Boolean{}
		// 	case "number":
		// 		paramTyp = &Number{}
		// 	case "string":
		// 		paramTyp = &String{}
		// 	}
		// 	e.addType(param, paramTyp)
		// 	typ.Params = append(typ.Params, NameAndType{
		// 		Name: param.Literal,
		// 		Type: paramTyp,
		// 	})
		// 	// TODO: Custom types
		// }
		// {
		// 	template := "@return"
		// 	defStart := strings.Index(comments, template)
		// 	if defStart != -1 {
		// 		// This is awful
		// 		def := strings.TrimSpace(strings.Split(comments[defStart+len(template):], "\n")[0])
		// 		switch def {
		// 		case "any":
		// 			typ.Return = &Any{}
		// 		case "boolean":
		// 			typ.Return = &Boolean{}
		// 		case "number":
		// 			typ.Return = &Number{}
		// 		case "string":
		// 			typ.Return = &String{}
		// 		}
		// 	}
		// }
		// e.addType(stmt, typ)
		// // TODO: Parse index expression
		// e.addType(stmt.Left, typ)

		// e.resolveExprType(stmt.Left)
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
	case *ast.ReturnStatement:
		for _, expr := range stmt.Exps {
			e.resolveExprType(expr)
		}
	default:
		e.addError(stmt, "Unimplemented")
	}
}

func (e *Environment) resolveExprType(expr ast.Expression) Type {
	e.pushNode(expr)
	defer e.popNode()
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
		def := e.FindDefinition(ast.NodePath{Node: expr})
		if def != nil {
			typ := e.Types[def]
			if typ != nil {
				return e.addType(expr, typ)
			}
		} else {
			e.Errors = append(e.Errors, parser.ParserError{Message: fmt.Sprintf("Unknown variable '%s'", expr.Token.Literal), Range: ast.Range(expr)})
		}

	case *ast.FunctionExpression:
		typ := &Function{Params: []NameAndType{}}
		for _, param := range expr.Params {
			// TODO: Function parameter types - requires parsing doc comments
			typ.Params = append(typ.Params, NameAndType{Name: param.Token.Literal, Type: &Unknown{}})
		}
		return e.addType(expr, typ)
	case *ast.FunctionCall:
		return e.resolveFunctionCallType(expr)
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

		// e.resolveExprType(expr.Inner)

		var key string
		switch inner := expr.Inner.(type) {
		case *ast.Identifier:
			key = inner.Token.Literal
		case *ast.StringLiteral:
			key = inner.Unit.Token.Literal
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

		for i := len(e.Nodes) - 1; i >= 0; i-- {
			switch node := e.Nodes[i].(type) {
			case *ast.AssignmentStatement:
				for j, leftVar := range node.Vars {
					if leftVar == expr {
						tbl.Fields = append(tbl.Fields, NameAndType{
							Name: key,
							Def:  node.Exps[j], // TODO: Out of bounds checking
							Type: e.Types[node.Exps[j]],
						})
						e.addType(expr, e.Types[node.Exps[j]])
						e.addType(expr.Inner, e.Types[node.Exps[j]])
						return nil
					}
				}
			case *ast.FunctionStatement:
				tbl.Fields = append(tbl.Fields, NameAndType{
					Name: key,
					Def:  node,
					Type: e.Types[node],
				})
				e.addType(expr, e.Types[node])
				e.addType(expr.Inner, e.Types[node])
				return nil

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
				Name: ident.Token.Literal,
				Type: valueTyp,
			}
			typ.Fields = append(typ.Fields, field)
			e.addType(fieldNode, valueTyp)
			e.addType(ident, valueTyp)
		}
		e.addType(expr, &typ)
		return &typ
	default:
		e.addError(expr, "Unimplemented")
	}

	return nil
}

func (e *Environment) resolveFunctionCallType(fc *ast.FunctionCall) Type {
	typ := e.resolveExprType(fc.Name)
	if typ == nil {
		return nil
	}
	e.addType(fc, typ)
	function, ok := typ.(*Function)
	if !ok {
		e.addError(fc, "'%s' is not a function", fc.Name)
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

func (e *Environment) FindDefinition(path ast.NodePath) *ast.Identifier {
	if path.Node == nil {
		return nil
	}
	identFor, ok := path.Node.(*ast.Identifier)
	if !ok {
		return nil
	}

	pos := identFor.Pos()
	var def *ast.Identifier
	ast.Walk(&e.file.Block, func(node ast.Node) bool {
		isAfter := node.Pos() > pos && pos < node.End()
		if isAfter {
			return false
		}
		isInside := node.Pos() <= pos && pos < node.End()
		switch node := node.(type) {
		case *ast.ForInStatement:
			if isInside {
				for _, ident := range node.Names {
					if ident.Token.Literal == identFor.Token.Literal {
						def = ident
					}
				}
			}
		case *ast.ForStatement:
			if isInside {
				if node.Name != nil && node.Name.Token.Literal == identFor.Token.Literal {
					def = node.Name
				}
			}
		case *ast.FunctionExpression:
			if isInside {
				for _, ident := range node.Params {
					if ident.Token.Literal == identFor.Token.Literal {
						def = ident
					}
				}
			}
		case *ast.FunctionStatement:
			if isInside {
				for _, ident := range node.Params {
					if ident.Token.Literal == identFor.Token.Literal {
						def = ident
					}
				}
			}
			if ident, ok := node.Left.(*ast.Identifier); ok {
				if ident.Token.Literal == identFor.Token.Literal {
					def = ident
				}
			}
		case *ast.LocalStatement:
			for _, ident := range node.Names {
				if ident.Token.Literal == identFor.Token.Literal {
					def = ident
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

func (e *Environment) pushNode(node ast.Node) {
	e.Nodes = append(e.Nodes, node)
}

func (e *Environment) popNode() {
	if len(e.Nodes) == 0 {
		return
	}
	e.Nodes = e.Nodes[:len(e.Nodes)-1]
}
