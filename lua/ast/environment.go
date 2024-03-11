package ast

type Environment map[string]*VariableDeclaration

func (e *Environment) Add(name string, definition Expression) {
	decl := NewVariable(name, definition)
	(*e)[name] = &decl
}
