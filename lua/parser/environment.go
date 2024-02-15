package parser

import (
	"slices"

	"github.com/raiguard/luapls/lua/ast"
)

// Pushes a new environment onto the environment stack.
func (p *Parser) pushEnv(env *ast.Environment) {
	p.envs = append(p.envs, env)
}

// Pops the topmost environment from the environment stack.
func (p *Parser) popEnv() {
	if len(p.envs) == 0 {
		panic("Attempted to pop a nil environment.")
	}
	p.envs = slices.Delete(p.envs, len(p.envs)-1, len(p.envs))
}

// Returns the topmost environment from the stack.
func (p *Parser) curEnv() *ast.Environment {
	if len(p.envs) == 0 {
		panic("Attempted to access a nil environment.")
	}
	env := p.envs[len(p.envs)-1]
	if env == nil {
		panic("Attempted to access a nil environment.")
	}
	return env
}
