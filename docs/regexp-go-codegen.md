# RegExp graph-to-Go implementation

## Boundary

Statically known RegExp patterns follow this pipeline:

```text
pattern and flags
  -> V8-shaped UTF-16 parser and AST
  -> V8-shaped AST-to-graph lowering and preprocessing
  -> custom traversal that emits a matcher-specific Go function
  -> Go compiler
```

The generated function contains ordinary Go control flow for that pattern.
It has no opcode array, instruction dispatcher, regex interpreter, Goja, or
third-party regex engine. The matcher-local backtracking stack stores runtime
alternatives; Go source labels and branches encode the pattern's control flow.

V8 behavior above the graph boundary is the compatibility target. The graph
traversal and Go source generation below that boundary are custom. V8's
post-graph analyses, assembler, bytecode generator, and architecture-specific
machine-code backends are outside the port.

## Current implementation

- `internal/regexpengine/regexp_parser.go` and `regexp_parser_api_v8.go` parse
  ECMAScript patterns using UTF-16 code units and produce the V8-shaped AST.
- `internal/regexpengine/regexp_ast.go` and `regexp_ast_behavior_v8.go` model
  V8 AST nodes, visitor behavior, range operations, and class-set behavior.
- `internal/regexpengine/regexp_compiler.go` contains graph nodes, V8-shaped
  AST-to-graph lowering, graph preprocessing, and the state needed to build
  that graph. It has no assembler or runtime execution backend.
- `internal/regexpengine/regexp_go_codegen.go` emits a separate Go matcher
  function for each statically known pattern. The function consumes UTF-16
  code units and returns UTF-16 match and capture offsets.
- `internal/goemit/regex_codegen.go` connects static regex literals and
  statically known constructors to code generation. Runtime-computed pattern
  sources are diagnosed during compilation.
- `internal/goemit/regex_runtime.go` supplies the JavaScript RegExp object
  surface and String methods around those generated matcher functions.

The graph-to-Go backend supports forward text, character ranges, ordered
choices and repeats, captures, anchors, basic word boundaries, and
case-sensitive backreferences. It reports unsupported graph nodes as compile
errors. In particular, lookaround execution, ignore-case matching, and some
Unicode forms still need Go-source generation support. The parser and graph
lowerer accept more syntax than the current emitter can execute.

This is a working vertical slice, not complete JavaScript RegExp
compatibility. The parser, AST, and graph-lowering files still need a
systematic comparison with their V8 sources. The comparison progress is kept
in `docs/conformance/regexp-v8-port.md`; a source line advances only after its
logic has been ported thought for thought.

## Next work

1. Finish the source-by-source audit of parser, AST, flags/errors, and
   AST-to-graph construction. Keep V8 file names and source ordering where
   practical, and record the last fully ported source line for each file.
2. Extend the custom Go emitter for every graph node the frontend can create,
   while keeping generated matchers as ordinary compiled Go functions.
3. Complete the JavaScript RegExp object and String-method behavior around the
   generated functions.
4. Compare against ECMAScript and V8 behavior after each compatibility area is
   implemented.
