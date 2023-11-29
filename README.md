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

- LSP settings
  - Also support standalone configuration file (`luapls.json`)?
- Environments system
  - Analagous to Lua environments
  - Every environment has pre-defined "root" files and library files that are pre-loaded (pulled from settings)
  - Will facilitate cross-file diagnostics and go-to, and global variables
- Go-to-references / highlight references
- Unicode support in the Lexer (Lua itself doesn't support unicode identifiers, but we should support them for better diagnostics)
- Proper autocompletion and diagnostics of table members
- Automated testing of LSP faculties
  - Utilize kak-lsp somehow?
- Increase parser test coverage to ensure it doesn't allow invalid syntax
- MOAR DIAGNOSTICS
  - Most important: Unknown field, undefined global
- Support other Lua versions (5.3, 5.4, etc)
- Type annotations
- Third-party addons / plugins
  - MVP to enable Factorio require semantics and event handler typing (after the type system is made)
