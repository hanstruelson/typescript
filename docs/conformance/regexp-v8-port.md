# V8 RegExp compiler-to-Go port progress

| V8 source file / scope | Go target | Last completed V8 source line |
| --- | --- | ---: |
| `src/regexp/regexp-parser.cc` | `internal/regexpengine/regexp_parser.go`, `internal/regexpengine/regexp_parser_support.go`, `internal/regexpengine/regexp_unicode_properties.go`, `internal/regexpengine/regexp_unicode_identifiers.go`, `internal/regexpengine/regexp_unicode_aliases.go`, `internal/regexpengine/regexp_unicode_categories_scripts.go`, `internal/regexpengine/regexp_unicode_binary_properties.go`, `internal/regexpengine/regexp_unicode_script_extensions.go`, `internal/regexpengine/regexp_unicode_emoji_properties.go` | 3641 |
| `src/regexp/regexp-parser.h` | `internal/regexpengine/regexp_parser_api_v8.go` | 39 |
| `src/regexp/regexp-flags.h` | `internal/regexpengine/regexp_flags.go` | 82 |
| `src/regexp/regexp-utils.h` and `src/regexp/regexp-utils.cc` | `internal/goemit/regex_runtime.go`, `internal/goemit/strings.go` | 225 |
| `src/regexp/regexp-result-vector.h` and `src/regexp/regexp-result-vector.cc` | `internal/goemit/regex_runtime.go`, `internal/goemit/runtime.go`, `internal/regexpengine/regexp_go_codegen.go` | 84 |
| `src/regexp/regexp-error.h` | `internal/regexpengine/regexp_errors.go` | 70 |
| `src/regexp/regexp-error.cc` | `internal/regexpengine/regexp_errors.go` | 24 |
| `src/regexp/regexp-ast.h` | `internal/regexpengine/regexp_ast.go`, `internal/regexpengine/regexp_ast_interval.go`, `internal/regexpengine/regexp_ast_behavior_v8.go`, `internal/regexpengine/regexp_compiler.go` | 808 |
| `src/regexp/regexp-ast.cc` | `internal/regexpengine/regexp_ast.go`, `internal/regexpengine/regexp_ast_behavior_v8.go` | 280 |
| `src/regexp/special-case.h` | `internal/regexpengine/regexp_casefold.go`, `regexp_casefold_data.go`, `regexp_compiler.go`, `regexp_parser.go` | 122 |
| `src/regexp/gen-regexp-special-case.cc` | `tools/genregexp-casefold.js`, `internal/regexpengine/regexp_casefold_data.go` | 222 |
| `src/regexp/regexp-nodes.h` — graph node data and edges | `internal/regexpengine/regexp_compiler.go`, `internal/regexpengine/regexp_go_codegen.go` | 1073 |
| `src/regexp/regexp-compiler.h` — compiler state, graph-analysis declarations, and constants | `internal/regexpengine/regexp_compiler_constants.go`, `internal/regexpengine/regexp_quick_check.go`, `internal/regexpengine/regexp_boyer_moore.go`, `internal/regexpengine/regexp_graph_analysis.go`, `internal/regexpengine/regexp_compiler.go`, `internal/regexpengine/regexp_go_codegen.go` | 789 |
| `src/regexp/regexp-compiler-tonode.cc` | `internal/regexpengine/regexp_compiler.go`, `internal/regexpengine/regexp_ast_behavior_v8.go`, `internal/regexpengine/regexp_casefold.go`, `internal/regexpengine/regexp_casefold_data.go`, `internal/regexpengine/regexp_compiler_constants.go` | 2421 |
| `src/regexp/regexp-compiler.cc` — graph metadata, quick-checks, first-character filters, and atomic-loop analysis (assembler emission is custom Go code generation) | `internal/regexpengine/regexp_compiler_constants.go`, `internal/regexpengine/regexp_ast.go`, `internal/regexpengine/regexp_ast_behavior_v8.go`, `internal/regexpengine/regexp_compiler.go`, `internal/regexpengine/regexp_compiler_boyer_moore.go`, `internal/regexpengine/regexp_compiler_first_char.go`, `internal/regexpengine/regexp_compiler_atomic_loop.go`, `internal/regexpengine/regexp_quick_check.go`, `internal/regexpengine/regexp_go_codegen.go` | 3571 |
| `src/regexp/regexp-compiler.cc` — Boyer-Moore lookahead range selection, frequency scoring, skip-table data, and generated-Go search filtering (`EmitSkipInstructions` machine-code output is replaced) | `internal/regexpengine/regexp_boyer_moore.go`, `internal/regexpengine/regexp_compiler_boyer_moore.go`, `internal/regexpengine/regexp_go_codegen_boyer_moore.go`, `internal/regexpengine/regexp_go_codegen.go` | 4172 |
| `src/regexp/regexp-compiler.cc` — text matching, reverse reads, case folding, word-boundary checks (`TextNode::Emit`, `TextEmitPass`, `EmitWordCheck`) | `internal/regexpengine/regexp_go_codegen_text.go`, `internal/regexpengine/regexp_go_codegen.go` | 2906 |
| `src/regexp/regexp-compiler.cc` — choice/loop branch and backtracking semantics; V8 assembler layout is replaced by generated Go labels and a backtrack switch | `internal/regexpengine/regexp_go_codegen.go` | 5880 |
| `src/regexp/regexp-compiler.cc` — positive and negative submatch actions, negative-submatch success, and empty-match checks | `internal/regexpengine/regexp_go_codegen_submatch.go`, `internal/regexpengine/regexp_go_codegen.go` | 5985 |
| `src/regexp/regexp-compiler.cc` — backreference matching in forward/reverse and ignore-case modes | `internal/regexpengine/regexp_go_codegen_text.go`, `internal/regexpengine/regexp_go_codegen.go` | 6049 |
| `src/regexp/regexp-compiler.cc` — graph analysis, Boyer-Moore propagation, surrogate preprocessing, and compiler stack checks; diagnostics-only tail omitted | `internal/regexpengine/regexp_graph_analysis.go`, `internal/regexpengine/regexp_compiler_boyer_moore.go`, `internal/regexpengine/regexp_compiler.go`, `internal/regexpengine/regexp_go_codegen.go` | 6530 |
| `src/regexp/regexp-node-printer.h` — diagnostics-only graph printer interface | — | 44 |
| `src/regexp/regexp-node-printer.cc` — diagnostics-only graph printer implementation; not `regexp.cc` | — | 238 |
| `src/regexp/regexp-code-generator.h` — bytecode-to-assembler interface (not used when emitting directly from the graph) | — | 86 |
| `src/regexp/regexp-code-generator.cc` — bytecode decoding, jump-target binding, and bytecode-to-assembler instruction dispatch (replaced by direct graph-to-Go source emission) | `internal/regexpengine/regexp_go_codegen.go`, `internal/regexpengine/regexp_go_codegen_text.go`, `internal/regexpengine/regexp_go_codegen_submatch.go`, `internal/regexpengine/regexp_go_codegen_boyer_moore.go` | 552 |
| `src/regexp/regexp-bytecode-generator.h` — bytecode writer and assembler declarations (no Go bytecode layer) | — | 248 |
| `src/regexp/regexp-bytecode-generator-inl.h` — byte emission, resizing, alignment, and debug bookkeeping (no Go bytecode layer) | — | 99 |
| `src/regexp/regexp-bytecode-generator.cc` — operand serialization, jump relocation, macro-assembler adapter, and peephole handoff (no Go bytecode layer) | — | 607 |
| `src/regexp/regexp-bytecode-iterator.h` — bytecode-array iterator API (unused without an encoded program) | — | 54 |
| `src/regexp/regexp-bytecode-iterator-inl.h` — instruction stepping and dispatch iteration (unused without an encoded program) | — | 54 |
| `src/regexp/regexp-bytecode-iterator.cc` — iterator pointer repair for V8 garbage collection (no Go heap handle equivalent needed) | — | 63 |
| `src/regexp/regexp-bytecodes.h` — opcode and operand schema for the bytecode representation (not part of direct graph-to-Go generation) | — | 404 |
| `src/regexp/regexp-bytecodes-inl.h` — compile-time operand layout, packing, and dispatch for bytecode (not part of direct graph-to-Go generation) | — | 387 |
| `src/regexp/regexp-bytecodes.cc` — bytecode disassembly and debug output | — | 102 |
| `src/regexp/regexp-stack.h` | `internal/regexpengine/regexp_go_codegen.go` | 214 |
| `src/regexp/regexp-stack.cc` | `internal/regexpengine/regexp_go_codegen.go` | 167 |
| `src/regexp/regexp.cc` — flag validation, syntax verification, and RegExp compile errors (experimental engine excluded) | `internal/regexpengine/regexp_flags.go`, `internal/regexpengine/regexp_errors.go`, `internal/goemit/regex_codegen.go` | 194 |
| `src/regexp/regexp.cc` — `EscapeRegExpSource` | `internal/goemit/regex_runtime.go` | 383 |
| `src/regexp/regexp.cc` — compile cache, atom-selection heuristic, and literal initialization (AOT matcher deduplication and direct-Go atom/graph selection) | `internal/regexpengine/regexp.go`, `internal/goemit/regex_codegen.go` | 493 |
| `src/regexp/regexp.cc` — ensure-compiled dispatch (Go matchers are emitted before runtime) | `internal/regexpengine/regexp_go_codegen.go`, `internal/goemit/regex_codegen.go` | 521 |
| `src/regexp/regexp.cc` — experimental one-shot dispatch (intentionally excluded) | — | 531 |
| `src/regexp/regexp.cc` — generic RegExp execution and single-match result construction (experimental dispatch excluded) | `internal/regexpengine/regexp_go_codegen.go`, `internal/goemit/regex_runtime.go` | 586 |
| `src/regexp/regexp.cc` — raw atom execution entry and capture-name map ordering | `internal/regexpengine/regexp.go`, `internal/regexpengine/regexp_go_codegen_atom.go`, `internal/goemit/regex_codegen.go` | 652 |
| `src/regexp/regexp.cc` — last-match capture registers, subject, and input state | `internal/goemit/regex_runtime.go`, `internal/goemit/runtime.go` | 682 |
| `src/regexp/regexp.cc` — atom compile and execution | `internal/regexpengine/regexp_go_codegen_atom.go`, `internal/regexpengine/regexp_go_codegen.go`, `internal/goemit/regex_runtime.go` | 824 |
| `src/regexp/regexp.cc` — compile-to-Go orchestration and global execution (V8 tier-up, bytecode, isolate, and heap-cache machinery excluded) | `internal/regexpengine/regexp.go`, `internal/regexpengine/regexp_go_codegen.go`, `internal/goemit/regex_codegen.go`, `internal/goemit/regex_runtime.go` | 1685 |
| `src/regexp/regexp.cc` — V8 fixed-array result caches (optimization only; no matching semantics) | — | 1902 |
| `src/regexp/regexp.cc` — canonical flag formatting | `internal/regexpengine/regexp_flags.go` | 1912 |
| `src/regexp/regexp.cc` — optional tracing and execution diagnostics (not part of matching behavior) | — | 2038 |
| `src/regexp/regexp-macro-assembler-arch.h` — architecture-specific native assembler selection (replaced by portable Go source emission) | — | 45 |
| `src/regexp/regexp-macro-assembler.h` — character checks, position/register operations, backtracking, scans, and result contract (the graph emitter lowers these semantics directly; V8 ABI, architecture/SIMD, tracing, and experimental-fallback hooks are not used) | `internal/regexpengine/regexp_go_codegen.go`, `internal/regexpengine/regexp_go_codegen_text.go`, `internal/regexpengine/regexp_go_codegen_submatch.go`, `internal/regexpengine/regexp_go_codegen_boyer_moore.go`, `internal/goemit/regex_runtime.go` | 531 |
| `src/regexp/regexp.h` — parse/compile metadata, execution contract, match-info updates, named-capture order, and global matching (V8 heap/cache/diagnostic APIs replaced or omitted) | `internal/regexpengine/regexp.go`, `internal/regexpengine/regexp_go_codegen.go`, `internal/goemit/regex_codegen.go`, `internal/goemit/regex_runtime.go` | 338 |
| `src/regexp/regexp-macro-assembler.cc` — case-insensitive backreference comparison, range-array semantics, binary class dispatch, surrogate-boundary checks, Unicode unanchored advance | `internal/regexpengine/regexp_go_codegen_text.go`, `internal/regexpengine/regexp_go_codegen.go` | 321 |
| `src/regexp/regexp-macro-assembler.cc` — current-character loads and offset bounds checks | `internal/regexpengine/regexp_go_codegen_text.go`, `internal/regexpengine/regexp_go_codegen.go` | 354 |
| `src/regexp/regexp-macro-assembler.cc` — scalar SkipUntil candidate and stop-position behavior; generated Go uses a conservative first-unit scan in place of V8’s specialized scan encodings | `internal/regexpengine/regexp_go_codegen.go`, `internal/regexpengine/regexp_go_codegen_boyer_moore.go` | 523 |
| `src/regexp/regexp-macro-assembler.cc` — backtracking stack growth/limit and ASCII word-character map (V8 isolate, native-code entry, GC, and tracing wrappers replaced by generated Go/runtime behavior) | `internal/regexpengine/regexp_go_codegen.go`, `internal/goemit/regex_runtime.go` | 781 |
