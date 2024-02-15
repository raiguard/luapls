// Reference: https://mukulrathi.com/create-your-own-programming-language/intro-to-type-checking/

package typechecker

import "github.com/raiguard/luapls/lua/ast"

type TypeChecker struct {
	File   *ast.Chunk
	Errors []ast.Error
}

func New(file *ast.Chunk) TypeChecker {
	return TypeChecker{file, []ast.Error{}}
}

func (tc *TypeChecker) Check() {
}
