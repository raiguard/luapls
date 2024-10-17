# Type checking

## Step 1: Generate ASTs

* User defines "environments" that have a specified set of root files.
* Iterate these root files in order.
* Files are represented as a connected graph.
  * Need to store parent backreferences.
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

**First pass**

First pass assembles all types that exist, following `@class`, `@type`, and `@alias` annotations as well as table literal declarations (for dynamic table types).
It does _not_ resolve the fields of these types; it only gathers which types exist.
This is to allow recursive types.

**Second pass**

Second pass is performed depth-first, and resolves fields of each type as it encounters them.
"Unknown type" diagnostics are generated here, but "unknown field" etc. diagnostics are not.
At the end of this pass, we should have a complete type graph.

**Third pass**

Third pass checks all type declarations and accesses to ensure they are well-formed.
This is where "unknown variable" and "unknown field" diagnostics are generated.

**Future work**

How will type inference fit in?

# Dynamic table types

If a table is not `@class`ed, then it is a "dynamic table":
* Assigning to the table creates a field
* Can only read fields from the table that have been previously assigned.
This only works for fields that are defined with literal strings.

Assigning a global variable does not throw a warning, though I would like to add a way to be strict about which globals you introduce.
However, attempting to _read_ an unknown global (using as parameter, indexing, doing operations) will show a warning.

When you need to resolve the type of a variable, you need to check:
* The value of the assignment
* Definition of value
*

Walk backwards up the AST:
* If index expression, add left to parents
* If infix left, get type of right then apply that type to corresponding infix metamethod.
* If identifier in expression list or namelist, then resolve definition of identifier and get that type.
