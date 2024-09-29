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
	for _, stmt := range block.Pairs {
		e.resolveStmtType(stmt.Node)
	}
}

func (e *Environment) resolveStmtType(stmt ast.Statement) {
	e.pushNode(stmt)
	defer e.popNode()
	switch stmt := stmt.(type) {
	case *ast.AssignmentStatement:
		for i := 0; i < len(stmt.Vars.Pairs); i++ {
			leftVar := stmt.Vars.Pairs[i]
			if i >= len(stmt.Exps.Pairs) {
				e.addType(leftVar.Node, &Unknown{})
				continue
			}
			exp := stmt.Exps.Pairs[i]
			typ := e.resolveExprType(exp.Node)
			leftTyp := e.resolveExprType(leftVar.Node)
			if leftTyp != nil && typ != nil && leftTyp != typ {
				e.Errors = append(e.Errors, parser.ParserError{Message: fmt.Sprintf("Cannot assign '%s' to '%s'", typ, leftTyp), Range: ast.Range(leftVar.Node)})
				e.addType(leftVar.Node, leftTyp)
				continue
			}
			if typ != nil {
				e.addType(leftVar.Node, typ)
			} else {
				e.addType(leftVar.Node, &Unknown{})
			}
		}
	case *ast.ForStatement:
		typ := e.resolveExprType(stmt.Start.Node)
		if typ == nil {
			e.addType(stmt.Name, &Unknown{})
			return
		}
		if _, ok := typ.(*Number); !ok {
			e.addError(&stmt.Start, "Range expressions must be of type 'number'")
			return
		}
		e.addType(stmt.Name, typ)
		finishTyp := e.resolveExprType(stmt.Finish.Node)
		if typ != finishTyp {
			e.addError(&stmt.Finish, "Range end must be of type '%s'", typ)
		}
		if stmt.Step != nil {
			stepTyp := e.resolveExprType(stmt.Step.Node)
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
		for i := 0; i < len(stmt.Names.Pairs); i++ {
			ident := stmt.Names.Pairs[i]
			if stmt.Exps == nil {
				continue
			}
			if i >= len(stmt.Exps.Pairs) {
				e.addType(ident.Node, &Unknown{})
				continue
			}
			exp := stmt.Exps.Pairs[i]
			typ := e.resolveExprType(exp.Node)
			leftTyp := e.resolveExprType(ident.Node)
			if leftTyp != nil && typ != nil && leftTyp != typ {
				e.Errors = append(e.Errors, parser.ParserError{Message: fmt.Sprintf("Cannot assign '%s' to '%s'", typ, leftTyp), Range: ast.Range(ident.Node)})
				e.addType(ident.Node, leftTyp)
				continue
			}
			if typ != nil {
				e.addType(ident.Node, typ)
			} else {
				e.addType(ident.Node, &Unknown{})
			}
		}
	case *ast.ReturnStatement:
		for _, expr := range stmt.Exps.Pairs {
			e.resolveExprType(expr.Node)
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
		for _, param := range expr.Params.Pairs {
			// TODO: Function parameter types - requires parsing doc comments
			typ.Params = append(typ.Params, NameAndType{Name: param.Node.Token.Literal, Type: &Unknown{}})
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
				for j, leftVar := range node.Vars.Pairs {
					if leftVar.Node == expr {
						tbl.Fields = append(tbl.Fields, NameAndType{
							Name: key,
							Def:  node.Exps.Pairs[j].Node, // TODO: Out of bounds checking
							Type: e.Types[node.Exps.Pairs[j].Node],
						})
						e.addType(expr, e.Types[node.Exps.Pairs[j].Node])
						e.addType(expr.Inner, e.Types[node.Exps.Pairs[j].Node])
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
	for i := 0; i < len(fc.Args.Pairs); i++ {
		arg := fc.Args.Pairs[i]
		argTyp := e.resolveExprType(arg.Node)
		if i >= len(function.Params) {
			e.addError(arg.Node, "Unused parameter")
			break // TODO:
		}
		if argTyp == nil {
			continue // TODO:
		}
		if argTyp != function.Params[i].Type {
			e.addError(arg.Node, "Cannot use '%s' as '%s' in argument.", argTyp, function.Params[i].Type)
		}
		e.addType(arg.Node, argTyp)
	}
	if len(fc.Args.Pairs) < len(function.Params) {
		e.addError(fc, "Too few function parameters, expected %v, got %v", len(function.Params), len(fc.Args.Pairs))
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
				for _, ident := range node.Names.Pairs {
					if ident.Node.Token.Literal == identFor.Token.Literal {
						def = ident.Node
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
				for _, ident := range node.Params.Pairs {
					if ident.Node.Token.Literal == identFor.Token.Literal {
						def = ident.Node
					}
				}
			}
		case *ast.FunctionStatement:
			if isInside {
				for _, ident := range node.Params.Pairs {
					if ident.Node.Token.Literal == identFor.Token.Literal {
						def = ident.Node
					}
				}
			}
			if ident, ok := node.Name.(*ast.Identifier); ok {
				if ident.Token.Literal == identFor.Token.Literal {
					def = ident
				}
			}
		case *ast.LocalStatement:
			for _, ident := range node.Names.Pairs {
				if ident.Node.Token.Literal == identFor.Token.Literal {
					def = ident.Node
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
