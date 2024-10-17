package types

import "github.com/raiguard/luapls/lua/token"

// Named represents a named type, constructed with `@class`.
type Named struct {
	Name  string
	Range token.Range
}

func (n *Named) isType() {}

func (n *Named) String() string {
	return n.Name
}
