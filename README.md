# luapls

luapls ("Lua please") is a language server for the Lua programming language.

This project is under very heavy development and is not ready for use.

## Build

**Dependencies:**
- [go](https://go.dev)

```
go build
```

## Features

- Parsing of every Lua file in the directory where it is opened.
- Extremely basic, locals-only, file-local missing variable diagnostics, go-to-definition, and autocompletion.

## TODO (always in flux, not necessarily in order)

- Bring parser into alignment with the Lua specification (don't successfully parse illegal things)
- Add parent back-references to AST nodes
- Improve parser error handling
- Instrumentation for resolving types in the parser
- Basic types (literals, simple expressions)
- Table types (structs/classes)
- Union types
- Intersection types
- Inheritance
- Metatable handling
- As much type inference as possible
- Useful LSP integrations
- Third-party addons/plugins to modify AST
