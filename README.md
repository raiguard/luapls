# luapls

luapls is a language server for the Lua programming language.

This project is under very heavy development and is not ready for use.

## Build

**Dependencies:**
- [go](https://go.dev)

```
go build
```

## TODO (always in flux)

- Refactor parser to better support syntax errors.
  - Currently, when an unexpected token is encountered the parser will just keep
    going and emit a bunch of errors before ceasing parsing of the file.
- Validate that the AST is valid Lua syntax.
  - Don't replicate Lua compiler errors 1:1, be more helpful.
  - Still allow as many features as possible even within invalid sections.
- Rewrite parser tests to not rely on `String()` output.
- Remove `String()` methods from AST nodes - now that we have range highlighting in-editor it's not nearly
  as helpful any more, and is full of inconsistencies.
- Implement basic autocompletion of identifiers.
- Implement basic go-to definition for identifiers.
