package ast

type Environment map[string]VariableDeclaration

func (e *Environment) Add(name string, definition Expression) {
	(*e)[name] = NewVariable(name, definition)
}
