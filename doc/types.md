# Type checking

## Step 1: Generate ASTs

* User defines "environments" that have a specified set of root files.
* Iterate these root files in order.
* Gather a list of all the files it `require`s in order.
* Parse those files.
  * Detect if a file has already been parsed and don't do it again.
* Files are stored in memory with their AST and line breaks, then pointers to dependent files.
  * Also keep a list of dependee files?
* Lua itself doesn't allow require loops but I guess we have to check them in the future. Think about this.

## Step 2: Check all the types

End goals:
* Show all diagnostics for files in your environments at all times.
  * How to handle orphaned files? We need to check them somehow.
  * Have a way to mark a file as "ours"?
* Gather types from library files but only show diagnostics on those if they are actually open.
* If a file is not open, keep its diagnostics, line breaks, and exports in memory but throw away the AST itself.
  * AST can be regenerated as needed, it's fairly cheap to do so.
* Type annotations are a single global namespace so they need to be resolved first.

* First pass, top-down: Assemble info for all `@class`, `@type`, and `@alias` annotations.
  * These are in a shared global namespace and are entirely static, so they are the launchpad for the rest of the checking.
  * This can be performed for each file in parallel.
  * Only support fully defined `---@class`es for now.
* Second pass, bottom-up: Resolve file exports and global variables.
  * Fully typecheck the root scope and check for any global exports in non-executed child functions.

# Dynamic table types

If a table is not `@class`ed, then it is a "dynamic table":
* Assigning to the table creates a field
* Can only read fields from the table that have been previously assigned.
This only works for fields that are defined with literal strings.

Assigning a global variable does not throw a warning, though I would like to add a way to be strict about which globals you introduce.
However, attempting to _read_ an unknown global (using as parameter, indexing, doing operations) will show a warning.
