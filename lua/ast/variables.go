package ast

type VariableDeclaration struct {
	Name       string
	Definition Expression
	Type       Type
}

func NewVariable(name string, definition Expression) VariableDeclaration {
	return VariableDeclaration{Name: name, Definition: definition}
}
