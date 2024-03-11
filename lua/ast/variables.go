package ast

type VariableDeclaration struct {
	Name       string
	Definition Expression
	Type       Type
}

func NewVariable(name string, def Expression) VariableDeclaration {
	return VariableDeclaration{Name: name, Definition: def, Type: &Unknown{}}
}
