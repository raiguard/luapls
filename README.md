# luapls will not be finished. Check out [emmylua_ls](https://github.com/CppCXY/emmylua-analyzer-rust) instead!

# luapls

luapls ("Lua please") is a language server for the Lua programming language.

This project is under very heavy development and is not ready for use.

## Build

**Dependencies:**
- [go](https://go.dev)

```
go build
```

## TODO (always in flux, not necessarily in order)

- Rewrite parser for better error tolerance, semantic info (whitespace, comments), and unicode support
- Expand lexer and parser tests to cover as many cases as I can think of
- Implement doc comment parsing in the AST
- Basic types (literals, simple expressions)
- Table types (structs/classes)
- Union types
- Intersection types
- Interface types
- Inheritance
- Metatable handling
- As much type inference as possible
- Useful LSP integrations
- Third-party addons/plugins to modify AST
