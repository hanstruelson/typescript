# RegExp engine port inventory

This document tracks the planned replacement of the Goja RegExp boundary with a
standalone Go implementation based on the semantic portions of V8's RegExp
implementation. It is intentionally a source-level inventory: every upstream
RegExp file is listed, including architecture-specific V8 JIT glue that will
be explicitly reviewed and omitted from the portable semantic implementation.

Upstream snapshot: V8 commit `a74948bb81a597a6affabff4a867cc661f5c2e69`
(`chromium.googlesource.com/v8/v8`, shallow sparse checkout on 2026-10-11).
The V8 BSD license is retained at
`tsc/internal/regexpengine/LICENSE.v8` for the Go adaptation of the
experimental engine. No V8 C++ file is copied verbatim into the repository.

The snapshot contains 77 C++ headers/translation units and approximately
50,000 lines in total. Representative dependencies include
`src/execution/isolate.h`, `src/heap/factory.h`, `src/zone/*`,
`src/objects/js-regexp-inl.h`, V8's macro assembler, and ICU Unicode headers.
Those dependencies are why a mechanical C++→Go rewrite is not an independent
engine: it would also require porting V8's heap and code-generation ABI.

## Port boundary

The standalone engine will contain the ECMAScript-facing parser, AST, Unicode
classification, bytecode compiler, interpreter, capture bookkeeping, and
RegExp result semantics. These are the parts that determine JavaScript
observable behavior. Porting is tracked file by file: each upstream file is
read in full, its executable and declarative lines are mapped to Go or given an
explicit disposition, and uncovered semantics are recorded before moving on.
The goal is not textual C++ translation; it is to account for every source line
and preserve every relevant behavior while replacing V8-specific machinery.

V8's `Zone`, `Isolate`, heap handles, Torque-generated objects, ICU wrappers,
macro assembler, instruction buffers, and architecture entry points are not
portable RegExp semantics. They are replaced by Go-owned slices, indexes,
explicit program objects, and Unicode tables. The architecture-specific files
remain in this checklist so they can be accounted for, but are not mechanically translated;
they implement V8's native-code ABI rather than a reusable regex algorithm.

## Source checklist

Legend: `[ ]` not yet fully read; `[~]` fully read with implementation or
semantic gaps still open; `[x]` fully read and every material behavior is
either mapped to Go or explicitly dispositioned. A checked file is not by
itself a claim of ECMAScript conformance; differential validation is tracked
separately.

### Completed line audit: `src/regexp/regexp-flags.h`

Reviewed against the complete 82-line file at the snapshot commit above and
against `tsc/internal/regexpengine/engine.go`. The table has one entry for each
source line, including blank lines. The file is audited, but the `[~]` status
below remains because the audit exposed concrete gaps.

| V8 line | Source responsibility | Go mapping / disposition |
|---:|---|---|
| 1 | V8 copyright comment | Legal notice; retain if copying or substantially deriving code. |
| 2 | BSD license notice | Legal notice; retain if copying or substantially deriving code. |
| 3 | Points to the repository `LICENSE` file | Legal notice; retain if copying or substantially deriving code. |
| 4 | Blank | Formatting only. |
| 5 | Opens the include guard | C++ header-build mechanism; Go has no equivalent. |
| 6 | Defines the include guard | C++ header-build mechanism; Go has no equivalent. |
| 7 | Blank | Formatting only. |
| 8 | Includes `std::optional` | Used by `TryFlagFromChar` to represent an unknown flag; Go uses a bool/enum-style switch and error return instead. |
| 9 | Includes `std::ostream` | Used only for the stream-formatting declaration at line 78; diagnostic formatting, not matching behavior. |
| 10 | Blank | Formatting only. |
| 11 | Includes V8's `base::Flags` | Supplies V8's packed bit operations; Go uses a bool-field `Flags` struct. |
| 12 | Blank | Formatting only. |
| 13 | Opens the C++ namespace | Replaced by Go package qualification. |
| 14 | Blank | Formatting only. |
| 15 | TODO: decouple more code from `JSRegExp::Flags` | V8 integration debt; the standalone Go type avoids a `JSRegExp` dependency. |
| 16 | TODO: consider removing `JSRegExp::Flags` | Same integration note; no independent matching algorithm. |
| 17 | Blank | Formatting only. |
| 18 | Says flag order is important | Ordering constraint for V8's macro-generated enum/printing logic. |
| 19 | Says entries are sorted by flag character | Maintenance invariant; Go's parser uses a switch and does not rely on sort order. |
| 20 | Warns that bit positions are shuffled but must stay contiguous | V8 representation invariant; not externally observable in the Go bool representation. |
| 21 | Starts the single-source flag-list macro | The following entries are the source of V8 flag metadata. |
| 22 | Defines `d` / has-indices at bit 7 | Go has `Flags.HasIndices` and parses `d`; exposure of match `indices` results is outside this header and remains a separate audit item. |
| 23 | Defines `g` / global at bit 0 | Go has `Flags.Global`; `FindAll` uses it. |
| 24 | Defines `i` / ignore-case at bit 1 | Go has `Flags.IgnoreCase`; matching consumes it. |
| 25 | Defines V8 `l` / linear at bit 6 | Go `ParseFlags` accepts it and `CompileUTF16` selects a prioritized Thompson-NFA path. Its accepted syntax intentionally follows only a subset of V8 `CanBeHandled`; remaining gaps are recorded below. |
| 26 | Defines `m` / multiline at bit 2 | Go has `Flags.Multiline`; anchor behavior consumes it. |
| 27 | Defines `s` / dotAll at bit 5 | Go has `Flags.DotAll`; dot behavior consumes it. |
| 28 | Defines `u` / Unicode at bit 4 | Go has `Flags.Unicode`; parser and matcher use Unicode mode. |
| 29 | Defines `v` / Unicode sets at bit 8 | Go has `Flags.UnicodeSets`; `ParseFlags` also enables Unicode mode while preserving this distinction. |
| 30 | Defines `y` / sticky at bit 3 | Go has `Flags.Sticky`; `FindUTF16` restricts matching to the supplied start. |
| 31 | Blank | Formatting only. |
| 32 | Macro that maps each list entry to a bit-valued enum member | V8 packed representation generation; Go's explicit struct fields replace it. |
| 33 | Declares the `Flag` enum from the macro | V8 representation; no Go enum is needed for current bool fields. |
| 34 | Undefines the temporary macro | C++ preprocessor hygiene only. |
| 35 | Blank | Formatting only. |
| 36 | Macro that counts flag-list entries | Compile-time bookkeeping only. |
| 37 | Computes `kFlagCount` | Compile-time bookkeeping only; no Go runtime equivalent. |
| 38 | Undefines the temporary macro | C++ preprocessor hygiene only. |
| 39 | Blank | Formatting only. |
| 40 | Comment introducing alphabetical-order assertion | Documents a maintenance invariant. |
| 41 | Macro that expands to pairwise character comparisons | Compile-time assertion machinery only. |
| 42 | Asserts flag characters are alphabetically ordered | V8 maintenance invariant; Go's switch has no generated-order dependency. |
| 43 | Undefines the temporary macro | C++ preprocessor hygiene only. |
| 44 | Blank | Formatting only. |
| 45 | Comment introducing contiguous-bit assertion | Documents a representation invariant. |
| 46 | Macro that expands each entry to a bit mask | Compile-time assertion machinery only. |
| 47 | Begins the contiguous-bit `static_assert` | V8 representation invariant. |
| 48 | Completes the assertion and diagnostic text | V8 representation invariant; unnecessary for Go bool fields. |
| 49 | Undefines the temporary macro | C++ preprocessor hygiene only. |
| 50 | Blank | Formatting only. |
| 51 | Aliases `Flags` to `base::Flags<Flag>` | Replaced by the Go `Flags` struct. |
| 52 | Defines bitwise operators for V8 `Flags` | Replaced by direct boolean reads and assignments. |
| 53 | Blank | Formatting only. |
| 54 | Macro generating `Is<Name>` predicates | Accessor-generation machinery, not standalone semantics. |
| 55 | Defines each predicate as a bit test | Replaced by direct reads such as `flags.Global`. |
| 56 | Expands predicates for every flag | V8 convenience API; Go currently reads fields directly. |
| 57 | Undefines the temporary macro | C++ preprocessor hygiene only. |
| 58 | Blank | Formatting only. |
| 59 | Declares `IsEitherUnicode` | Semantically means either `u` or `v` mode. |
| 60 | Implements it as `IsUnicode(f) || IsUnicodeSets(f)` | Go preserves `UnicodeSets` and sets `Unicode` for `v`; consumers must treat both as Unicode mode. |
| 61 | Closes the helper | No extra behavior. |
| 62 | Blank | Formatting only. |
| 63 | Says the next helper may rewind an index | This is UTF-16 match-start behavior, not merely representation detail. |
| 64 | Links the behavior to `OptionallyStepBackToLeadSurrogate` | Points to the implementation/call-site reasoning to follow. |
| 65 | Declares the rewind predicate | Implemented by `ShouldOptionallyStepBackToLeadSurrogate` and used by `Program.FindUTF16`. |
| 66 | Requires Unicode (`u` or `v`) and global or sticky (`g` or `y`) | The Go matcher rewinds a start index that points at a low surrogate preceded by a lead surrogate; UTF-16 regression tests cover `uy`, `ug`, `vy`, and `vg`. |
| 67 | Closes the helper | No extra behavior. |
| 68 | Blank | Formatting only. |
| 69 | Disables clang-format for the generated ternary | Formatting directive only. |
| 70 | Macro maps a matching character to a flag enum value | Implements character lookup over the canonical flag list. |
| 71 | Completes each ternary branch | Lookup implementation detail. |
| 72 | Declares `TryFlagFromChar` | Implemented by Go's `TryFlagFromChar`, used by `ParseFlags` and canonical flag formatting. |
| 73 | Expands the character lookup and returns empty optional if unmatched | The character lookup is implemented and reused by Go's scoped `i`, `m`, `s` group-modifier parser. |
| 74 | Closes `TryFlagFromChar` | No extra behavior. |
| 75 | Undefines the temporary macro | C++ preprocessor hygiene only. |
| 76 | Re-enables clang-format | Formatting directive only. |
| 77 | Blank | Formatting only. |
| 78 | Declares stream formatting for flags | Defined in V8 `regexp.cc:1907`, which emits characters in canonical list order. Go's `Program.String` omits flags; diagnostic formatting, not match behavior, absent a Go API requirement. |
| 79 | Blank | Formatting only. |
| 80 | Closes the C++ namespace | Namespace syntax only. |
| 81 | Blank | Formatting only. |
| 82 | Closes the include guard | C++ header-build mechanism; Go has no equivalent. |

The flag helpers are now wired into the Go matcher and compiler. The upstream
call sites were checked: the surrogate rewind is covered by UTF-16 tests, `l`
selects the standalone NFA for its supported syntax subset, and scoped `i`,
`m`, and `s` modifiers are parsed and attached to the relevant AST operations.

### Parser, AST, flags, and Unicode semantics

- [~] `regexp-parser.cc`
- [~] `regexp-parser.h`
- [~] `regexp-ast.cc`
- [~] `regexp-ast.h`
- [x] `regexp-ast-printer.cc`
- [x] `regexp-ast-printer.h`
- [x] `regexp-flags.h`
- [x] `regexp-error.cc`
- [x] `regexp-error.h`
- [~] `regexp-utils.cc`
- [~] `regexp-utils.h`
- [~] `special-case.h`
- [~] `gen-regexp-special-case.cc`

### Program representation, compiler, and interpreter

- [x] `regexp-bytecodes.cc`
- [x] `regexp-bytecodes.h`
- [x] `regexp-bytecodes-inl.h`
- [~] `regexp-bytecode-generator.cc`
- [~] `regexp-bytecode-generator.h`
- [~] `regexp-bytecode-generator-inl.h`
- [x] `regexp-bytecode-iterator.cc`
- [x] `regexp-bytecode-iterator.h`
- [x] `regexp-bytecode-iterator-inl.h`
- [x] `regexp-bytecode-peephole.cc`
- [x] `regexp-bytecode-peephole.h`
- [~] `regexp-interpreter.cc`
- [~] `regexp-interpreter.h`
- [~] `regexp-stack.cc`
- [~] `regexp-stack.h`
- [x] `regexp-result-vector.cc`
- [x] `regexp-result-vector.h`
- [~] `regexp.cc`
- [~] `regexp.h`

### Graph compiler and diagnostics

- [~] `regexp-compiler-tonode.cc`
- [~] `regexp-compiler.cc`
- [~] `regexp-compiler.h`
- [x] `regexp-code-generator.cc`
- [x] `regexp-code-generator.h`
- [x] `regexp-nodes.h`
- [x] `regexp-node-printer.cc`
- [x] `regexp-node-printer.h`
- [x] `regexp-printer.cc`
- [x] `regexp-printer.h`
- [x] `regexp-dotprinter.cc`
- [x] `regexp-dotprinter.h`
- [x] `regexp-graph-printer.cc`
- [x] `regexp-graph-printer.h`
- [~] `regexp-macro-assembler.cc`
- [~] `regexp-macro-assembler.h`
- [x] `regexp-macro-assembler-arch.h`
- [x] `regexp-macro-assembler-tracer.cc`
- [x] `regexp-macro-assembler-tracer.h`

### Experimental engine

- [~] `experimental/experimental.cc`
- [~] `experimental/experimental.h`
- [x] `experimental/experimental-bytecode.cc`
- [~] `experimental/experimental-bytecode.h`
- [~] `experimental/experimental-compiler.cc`
- [~] `experimental/experimental-compiler.h`
- [~] `experimental/experimental-interpreter.cc`
- [~] `experimental/experimental-interpreter.h`

### Architecture-specific native code generators (reviewed, not translated)

These files are included so the inventory is complete. They emit machine code
using V8's macro assembler and depend on V8's calling convention, heap object
layout, write barriers, and isolate state. They are not required by the Go
bytecode interpreter.

- [ ] `arm/regexp-macro-assembler-arm.cc`
- [ ] `arm/regexp-macro-assembler-arm.h`
- [ ] `arm64/regexp-macro-assembler-arm64.cc`
- [ ] `arm64/regexp-macro-assembler-arm64.h`
- [ ] `ia32/regexp-macro-assembler-ia32.cc`
- [ ] `ia32/regexp-macro-assembler-ia32.h`
- [ ] `loong64/regexp-macro-assembler-loong64.cc`
- [ ] `loong64/regexp-macro-assembler-loong64.h`
- [ ] `mips64/regexp-macro-assembler-mips64.cc`
- [ ] `mips64/regexp-macro-assembler-mips64.h`
- [ ] `ppc/regexp-macro-assembler-ppc.cc`
- [ ] `ppc/regexp-macro-assembler-ppc.h`
- [ ] `riscv/regexp-macro-assembler-riscv.cc`
- [ ] `riscv/regexp-macro-assembler-riscv.h`
- [ ] `s390/regexp-macro-assembler-s390.cc`
- [ ] `s390/regexp-macro-assembler-s390.h`
- [ ] `x64/regexp-macro-assembler-x64.cc`
- [ ] `x64/regexp-macro-assembler-x64.h`

## Validation gates

1. Differential tests compare parser errors, flags, captures, indices, and
   replacement behavior against a current ECMAScript host.
2. The generated runtime keeps the existing compatibility adapter until the
   standalone interpreter passes the differential suite; this avoids silently
   changing JavaScript behavior while the port is incomplete.
3. Once parity is reached, the RegExp path is switched to the standalone Go
   engine and the Goja dependency is audited separately for non-RegExp uses.

## Current implementation

`tsc/internal/regexpengine` has UTF-16 source/input entry points, flag parsing,
the existing backtracking interpreter, and a separate `l`-selected prioritized
NFA. The NFA covers the V8 experimental engine's basic supported flags and
captureless lookbehind (including nested lookbehind) with one-pass endpoint
tables; it rejects backreferences, lookahead, captured lookaround, Unicode and
ignore-case modes, class-set expressions, and unsupported quantifier bounds.
Quantifier threads track whether they consumed input since entering a repeated
body, corresponding to V8's `BEGIN_LOOP`/`END_LOOP` behavior; nullable-loop
capture cases are tested in both the normal backtracking path and against
Node's V8 `l` engine. The backtracking VM backtracks through empty alternatives
before leaving a nullable loop and preserves the last meaningful capture. It
is not a byte-for-byte port of V8's bytecode compiler/interpreter: V8's capture-clock/
filter machinery, optional captured-lookaround engine, interrupt handling,
per-isolate resource policy, and full bytecode format are not reproduced. The
standalone package offers bounded search APIs, but its older `Find`/`FindUTF16`
conveniences remain unbounded and the Goja runtime has not adopted the bounded
API.

The parser source (`regexp-parser.cc` and `.h`) has now been read in full. This
pass implemented scoped `i`, `m`, and `s` modifiers, including flag removal,
nested restoration, and mode-sensitive literals, classes, backreferences,
anchors, boundaries, and dot. It also fixed legacy character classes with a
trailing hyphen so the closing bracket is not accidentally added to the class.
Named capture names and named backreferences also decode Unicode escapes,
including braced code points and escaped surrogate pairs outside `u` mode.
The parser pre-scans for named captures so legacy `\k` remains an identity
escape when the pattern has no named groups, even when the apparent backref is
unterminated; forward named references still resolve. These behaviors have
regression tests; modifier and escaped-name cases were checked against Node.
The capture-group limit is also capped at V8's 32,767 explicit groups, with a
regression test that checks the corresponding “Too many captures” diagnostic.
The normal matcher covers all three modifiers; the `l`
subset accepts scoped `m`/`s` and explicitly rejects scoped `i`, which its
case-sensitive NFA cannot preserve. Parser files remain `[~]` because the full
V8 grammar and early-error behavior are not yet implemented.

The duplicate-capture path analysis from `regexp-parser.cc` is now represented
in the Go parser: named captures with the same name are accepted only when
their enclosing disjunction branches prove they cannot both participate.
Named backreferences retain all alternative capture indices and consume the
one whose slot was set, or match empty if none participated. The compiled
program exposes `GroupCaptureIndices` for all such slots while keeping
`GroupIndices` as a stable first-index view. Tests cover each participating
branch, an unmatched duplicate name, and sequential/co-occurring duplicates
that must still be rejected. Unicode property names now use exact,
case-sensitive aliases generated from UCD 15.1 `PropertyAliases.txt` and
`PropertyValueAliases.txt`; loose case and separator normalization is rejected.
`Script_Extensions` now checks Unicode 15.1 override ranges before consulting
the default Script property. UCD 15.1 range tables now back the V8 binary
property allowlist, general categories, Script values, and Script_Extensions.
Emoji 15.1 data backs the emoji character properties and the `Basic_Emoji`,
keycap, RGI flag, modifier, tag, ZWJ, and aggregate `RGI_Emoji` string
properties in `/v` mode. Case-folding still uses Go Unicode 17 tables. Exact
V8 syntax/error parity and unresolved grammar cases remain open; this does not
complete the parser's semantic port.

`regexp-parser.h` is the public parser boundary rather than another parser
algorithm: its heap-string entry point and no-GC syntax-verification template
are represented by the Go `Compile`/`CompileUTF16` API, which owns its UTF-16
input and returns a compiled program or syntax error. The V8 `Zone`, `Isolate`,
heap handles, and `DisallowGarbageCollection` lifetime constraints are not
portable parser semantics and are intentionally replaced by Go ownership.

`regexp-ast.h` and `regexp-ast.cc` have also been read in full. The Go AST's
`node` and `classExpr` forms cover the source's executable tree types; visitor
dispatch maps to opcode/type switches. Capture-register intervals, anchored
subtree summaries, and min/max-length summaries are compiler optimizations and
are not needed for the direct interpreter. A 2,048-level parser nesting guard
for groups and nested `v` classes provides a Go-side recursion safety bound,
not an exact port of V8's native-stack check. Empty
classes, negated classes, class-set operations, and longest-first string-class
matching are represented directly in the parser/matcher rather than by V8's
in-place AST normalization. The AST checklist remains `[~]`: complete parser
compatibility and V8's range/case-fold tables are still open.

The main engine also parses numeric/named backreferences and lookarounds,
including reverse execution for lookbehind; parses nested `v` classes and
code-point union/intersection/subtraction; and has partial Unicode-property
support based on Go's Unicode tables, which can differ from V8's ICU version.
Unicode properties whose value is a set of strings (such as RGI emoji
properties), complete property aliases, and several grammar/early-error cases
remain unsupported. The generated runtime still routes RegExp operations
through Goja, so the standalone package is not yet a replacement. No checklist
item is marked complete merely because a partial implementation exists.

### Experimental interpreter source review: `experimental-interpreter.cc` and `.h`

Both files have now been read in full (the implementation is 1,607 lines).
The Go NFA preserves V8's priority-sensitive accept rule: it records the best
accepted thread but keeps higher-priority consuming threads alive, using
ordered thread slices and per-position `(pc, consumed-since-loop)` deduplication.
V8's capture registers, fork copies, and loop-consumption state map to Go
capture slices, ordered split expansion, and the thread's consumed bit. The Go
interpreter also clears captures at repeat boundaries, reproducing the common
last-iteration behavior without V8's optional clock arrays and filter-tree
bytecode.

The source's lookbehind automata, thread scheduling, repeated-match advancement,
and UTF-16 code-unit assertions were compared directly. This review found and
fixed an observable anchor mismatch: multiline `^`/`$` now recognize the
boundary between CR and LF, and non-multiline `$` no longer incorrectly
matches before a trailing line terminator. Regression tests cover both. Go's
Unicode global zero-width advancement is more complete than this V8
experimental implementation, which explicitly has Unicode disabled.

V8's one-byte/two-byte storage specialization, heap handles, GC-safe raw
pointers, interrupt/retry protocol, sandbox checks, thread recycling
allocators, output-register buffers, memory-accounting flags, and advanced
capture/lookaround clocks are not reproduced. The Go implementation uses
owned UTF-16 slices, explicit execution limits, and a simpler captureless
lookbehind table. The compiler/interpreter pair still lacks V8's captured
lookaround and lookahead-table mode and remains `[~]`.

### Experimental compiler source review: `experimental-compiler.cc` and `.h`

Both files have now been read in full (the implementation is 1,226 lines).
V8's supported-flag, backreference, class-string, lookaround, and quantifier
replication checks map to `linearUnsupportedReason`. V8's patch-list labels map
to Go fragment patch lists; ordered disjunctions map to prioritized NFA splits;
character classes map to normalized range membership; capture-register clears
map to `linearClearCapture`; and empty-loop protection maps to the
`linearBeginLoop`/`linearEndLoop` consumed-input bit. The NFA injects a new
start thread at each position rather than emitting V8's non-greedy search
preamble, with sticky mode restricting injection to the caller's start.

The Go compiler deliberately does not reproduce V8's optional capture-group
filter bytecode, AST breadth-first filter tree, or linear-size optimization for
non-nullable pluses; it clears captures in the NFA and emits direct fragments.
It also supports only captureless lookbehind, not V8 configurations that can
compile lookahead or captured lookaround with tables. The emitted Go program is
not an 8-byte V8 instruction stream. These are subset/optimization gaps, so the
compiler pair remains `[~]` despite the mapped supported subset.

### Experimental bytecode source review: `experimental-bytecode.cc` and `.h`

Both files have now been read in full. The `.cc` file only formats bytecode and
lookaround payloads for diagnostics; the Go engine has no V8-bytecode dump
format, so that printer is deliberately omitted. The header defines V8's
fixed-width instruction encoding and its abstract NFA semantics. The Go linear
program maps ordered `FORK`/`JMP` to prioritized split/jump threads, range
consumption to rune/class instructions, register writes/clears to capture
instructions, `ACCEPT` to `linearMatch`, and loop-clock behavior to
`linearBeginLoop`/`linearEndLoop`. Source-order thread priority and capture
restoration are covered by regression tests.

V8's byte layout, union payload packing, disassembler helpers, and filter
instructions are not copied. The Go compiler evaluates supported character
classes directly and implements captureless lookbehind with endpoint tables
instead of V8's start/end/read/write lookaround-table bytecodes. It does not
implement the experimental engine's full lookaround/capture feature set, so the
header remains `[~]`; its matching behavior is only ported for the explicit Go
subset.

### Experimental engine entry-point source review: `experimental.cc` and `.h`

Both files have now been read in full. V8's eligibility walk maps to
`linearUnsupportedReason`/`compileLinear`; `Initialize` and lazy bytecode
compilation map to creating a Go `Program` at compile time; `Exec` and `ExecRaw`
map to the UTF-16 search APIs; and capture-name metadata is exposed as
`Program.GroupCaptureIndices`. V8's heap-owned regexp data, lazy bytecode
installation, flattened heap strings, raw result-register buffers, isolate
exceptions, sandbox/GC scopes, and JS-call ABI have no standalone matcher
counterpart.

The V8 one-shot retry after excessive backtracking is now available through
the opt-in `FindWithLinearFallback`/`FindUTF16WithLinearFallback` API. It
preserves existing `FindWithLimits` error semantics and retries only patterns
accepted by the linear subset; unsupported patterns keep the original limit
error. Each attempt uses the caller's resource ceiling. This is a Go API
adaptation, not an automatic JS runtime fallback: the generated RegExp adapter
still uses Goja, and V8's retry status, interrupt hooks, and call-origin
protocol are not reproduced. Those integration differences keep both source
rows at `[~]`.

### Main interpreter source review: `regexp-interpreter.cc` and `.h`

Both files have now been read in full at the pinned V8 revision. The Go engine
does not reproduce V8's packed bytecode format or its native dispatch loop; it
compiles the parsed AST to an instruction slice and executes that with a Go VM.
The following interpreter behavior is represented in that VM:

| V8 interpreter behavior | Go implementation |
|---|---|
| Push/pop backtrack targets and positions; restore captures on a failed branch | `opSplit`, `opJump`, and `executionState` in `runAtWithCaptures` |
| Read/write/clear capture registers and return match spans | `opSaveStart`, `opSaveEnd`, `opClearCaptures`, and `opMatch` |
| Character, dot, class, anchor, and word-boundary checks | `consumeAt`, `classEndsAt`, `lineStartAssertion`, `lineEndAssertion`, and boundary handling in `runAtWithCaptures` |
| Forward/reverse backreferences, including an unmatched capture matching empty | `consumeAt` and reverse subprograms used for lookbehind |
| Lookaround success/failure and positive-lookaround captures | `opAssert` plus nested `runAtWithCaptures` |
| Repeated captures not participating in the final iteration | `opClearCaptures`, added based on the compiler/interpreter interaction and covered by regression tests |
| Multiple global matches and advancing past an empty match by a code point in Unicode mode | `FindAll` and `advanceStringIndex` |

The source's bytecode dispatch table, packed operand decoding, bytecode tracing,
one-byte/two-byte template specialization, V8 stack/interrupt hooks, GC handle
repair, tier-up/retry protocol, and counters are runtime or representation
machinery rather than matching rules. The Go VM intentionally replaces them.
There are still behaviorally significant gaps: the Go backtracking VM snapshots
capture arrays at split points instead of using V8's register/backtrack-stack
protocol; the legacy `Find`/`FindUTF16` APIs are unbounded, and the new bounded
API uses Go-owned limits rather than V8's interrupt/backtrack-fallback
protocol; and Unicode case folding/property classification depends on Go's
Unicode tables rather than V8's ICU version. The generated runtime also still
uses Goja for JavaScript RegExp operations. For those reasons these source rows
remain `[~]`, not `[x]`.

The Go VM is bytecode-like, but its instructions are Go structs rather than
V8's byte-packed stream. That keeps control flow explicit and gives one
interpreter implementation for the supported feature set. It is a sensible
portable baseline; it is not automatically the fastest execution strategy.
V8's graph compiler can specialize and emit native code, while an interpreter
has dispatch and instruction-fetch overhead. This port should first establish
semantics and differential tests; optimizing or adding a JIT is a separate
measured decision.

### Graph-node declaration source review: `regexp-nodes.h`

This 1,073-line header has been read in full. It declares V8's native graph
compiler IR, not the parser AST: sequence/action nodes represent deferred
register and position changes, text/assertion/backreference nodes represent
consuming and zero-width operations, choice/loop nodes represent ordered
alternatives and repetition, and end nodes represent accept/backtrack and
negative-lookaround completion. The Go AST-to-program compiler lowers these
same observable operations directly to `opSave*`, `opClearCaptures`,
`opAssert`, `opSplit`/`opJump`, consuming instructions, and `opMatch`.
`EatsAtLeast`, `NodeInfo`, quick checks, Boyer–Moore lookahead, trace reuse,
preloading, and the fixed-loop `DrainMode`/`AtomicLoopKind`/`ParkedGrant`
analysis exist to reduce native code size and rescanning; they do not change
which match wins. The Go KMP atom path covers the simplest literal scan, while
the general Go VM has no equivalent of V8's graph/JIT search fusions. Thus the
header's matching operations are mapped to the direct VM and its graph-analysis
and native optimization machinery is explicitly omitted; the performance gap
is not treated as a compatibility success claim.

### Graph compiler interface source review: `regexp-compiler.h`

This 789-line header has been read in full. Its ECMAScript class-range
constants for whitespace, word, digits, line terminators, and surrogates are
represented by the Go class/property helpers; `NeedsUnicodeCaseEquivalents`
selects the Unicode-folding rules now implemented by the matcher. Compiler
register allocation and capture-zero wrapping map to the Go capture layout,
and `OptionallyStepBackToLeadSurrogate` maps to the UTF-16 start correction
tested on `u`/`v` global and sticky searches. Unicode code-point decoding is
represented directly by `inputCharacter`, including lone surrogates; a new
test checks that a Unicode class can match a lone lead surrogate but cannot
match one half of an actual pair.

`QuickCheckDetails`, Boyer–Moore maps/frequency collection, `Trace` deferred
code-generation actions, preloads, register-stack plumbing, and
`SpecialLoopState` only exist to build/reuse optimized native traces. The Go
compiler translates its simpler AST directly to VM instructions instead.
Compiler-size and internal-register overflow are still not mapped to V8's
exact `Regular expression too large`/overflow errors, and the full compiler
driver remains under review; this header stays `[~]` pending those limits and
end-to-end parser/compiler error mapping.

### Graph compiler implementation source review: `regexp-compiler.cc`

The complete 6,546-line implementation was read in consecutive source chunks.
Its responsibilities include character-range lowering, case-fold expansion,
quick-check and Boyer–Moore analyses, graph/trace emission, fixed-length and
atomic-loop handling, unanchored-search scans, bytecode/native compilation,
and assertion/eats-at-least propagation. The Go compiler emits typed VM
instructions directly from its AST rather than constructing V8's graph and
virtualized `Trace` objects.

Two directly useful optimizations are now represented in the Go path. A
conservative AST first-character analysis builds a 256-bit ASCII filter plus
Unicode-aware predicates and skips impossible unanchored start positions;
unknown or nullable forms fail open. Deterministic greedy repetitions of one
literal, dot, or ordinary character class lower to `opGreedyLoop`: it scans
UTF-16 code points and retains one retry frame instead of a state per
iteration, including reverse-mode lookbehind. Tests compare filtered and
unfiltered match results and exercise long greedy input under a small memory
ceiling.

The V8 `QuickCheckDetails`/trace machinery, multi-character lookahead tables,
Boyer–Moore skip tables, machine-specific SIMD scans, masked alternative
dispatch, exact class-range branch/table emission, and graph-version reuse are
not yet reproduced. The Go matcher instead uses exact normalized range
membership, a KMP literal path, the ASCII first-character filter, and the
compact single-code-point greedy loop. These remaining differences are chiefly
performance/compiler architecture, not a claim that the corresponding V8
source lines were copied. Multi-node fixed-length loops, atomic drain proofs,
compiler recursion/stack overflow mapping, and V8's exact compile-error
catalog wiring remain open, so this file stays `[~]`.

### AST-to-graph source review: `regexp-compiler-tonode.cc`

The complete 2,421-line file was read in consecutive chunks. `Tree::ToNode`
dispatches each AST form into graph nodes; character classes canonicalize and
case-fold ranges, Unicode classes split BMP/non-BMP/lone-surrogate cases,
Unicode-set operands perform union/intersection/subtraction, disjunctions
sort/factor literal runs and fold single-character alternatives, and quantifier
lowering decides unrolling, capture clearing, counters, and empty-iteration
checks. Assertions and lookarounds are converted into zero-width graph paths
with saved position/capture state.

The Go engine evaluates Unicode code points directly from UTF-16 (including
lone-surrogate handling), canonicalizes class ranges, executes Unicode-set
expressions and string alternatives, and lowers lookarounds/captures into its
VM. The graph builder's literal-alternative optimizations now factor shared
literal prefixes and lower same-mode single-character branches into normalized
classes while preserving source order; tests cover priority (`a|ab`) and outer
captures. Quantifier lowering uses mandatory unrolling plus the compact greedy
single-character loop where safe; general/capture-bearing loops retain the
backtracking path.

V8's graph-level Unicode range splitting is replaced by code-point-aware Go
matching, and its Unicode case-fold tables are still Go's Unicode version plus
the explicit supplementary special case noted above. Class-set expressions
remain direct runtime AST evaluation rather than being flattened into V8's
range/string operands, and V8's expansion-factor limiter, graph register
allocation, and several unrolling heuristics are not yet reproduced. Thus
this is a complete source review with compiler/performance gaps still open,
matching `[~]` rather than a conformance-complete claim.

### Graph diagnostics source review: node/printer/dot/graph printers

The node-printer, shared printer, DOT printer, and graph-printer `.cc`/`.h`
files have all been read in full. They visit the compiler graph to produce
colored text labels, capture/action descriptions, AST fragments, graph edges,
loop arrows, traces, and Boyer–Moore tables for diagnostics builds only. Their
`visited`/label bookkeeping is used to avoid repeating graph nodes while
printing; it is not part of matching. The Go package has no graph/JIT
diagnostics API, and its public matching results do not expose this internal
representation, so these facilities are deliberately omitted as
non-JavaScript-observable tooling.

### Error catalog source review: `regexp-error.cc` and `.h`

These two files have been read in full and their complete ordered error enum,
message table, bounds-checked message lookup, and stack/resource-overflow
predicate are adapted in `tsc/internal/regexpengine/errors.go`. Tests verify
the enum/table cardinality, representative messages, and overflow categories.
This ports the catalog itself; the parser still returns offset-bearing Go
errors rather than assigning every failure one of these V8 error codes, so
parser-to-catalog wiring remains open with the grammar/early-error work.

### Utility source review: `regexp-utils.cc` and `.h`

Both files have been read in full. The UTF-16 `AdvanceStringIndex` operation
and matched-capture checks are implemented in `utils.go`. `Match.CaptureUTF16`
offers a zero-copy subject view for the fast path, while the package-level
`CaptureUTF16` helper returns copied UTF-16 units so lone surrogates remain
intact and callers cannot mutate the input through the result. The utility's
`lastIndex` property access,
generic `RegExpExec` dispatch (including custom `exec` validation), and
unmodified-RegExp/prototype fast-path checks are object/runtime behavior and
remain supplied by Goja. These files stay `[~]` until that adapter is replaced
and those semantics are covered by project-level tests.

### Case-fold source review: `special-case.h` and `gen-regexp-special-case.cc`

Both files have been read in full. The Go matcher now follows ECMAScript's
non-Unicode canonicalization rule that blocks non-ASCII-to-ASCII uppercase
folds, and for class ranges it walks simple-fold equivalents while filtering
by that canonical key. This replaces V8's generated `IgnoreSet` optimization
with direct equivalence checks. Unicode-mode matching walks Go simple-fold
classes and includes the supplementary sharp-s equivalence explicitly
documented by this V8 source (`U+1DF95`, `U+00DF`, `U+1E9E`). Regression tests
cover long-s, Kelvin sign, sharp-s, and that supplementary pair. The runtime
tables still come from Go Unicode 17 rather than the exact ICU version used by
the pinned V8 snapshot, and the original ICU table generator is not ported;
newer or otherwise divergent fold classes may remain. These rows therefore
stay `[~]`.

### AST printer source review: `regexp-ast-printer.cc` and `.h`

Both files have been read in full. They are compiled only with V8's regexp
diagnostics option and format internal AST nodes, graph labels, colors, and
nested trace indentation. They do not affect parsing or matching; this port has
no public AST-diagnostics API, so the visitor/color/stream integration is
explicitly omitted rather than copied into the runtime. This is a deliberate
diagnostics-only disposition, not a compatibility claim about V8's debug text.

### Bytecode-definition source review: `regexp-bytecodes.cc`, `.h`, and `.inl.h`

All three files have been read in full. V8's control/backtrack instructions map
to the Go VM's split/jump program counter and saved execution states; capture
register instructions map to save/clear operations; character, class, anchor,
and backreference instruction families map to the corresponding high-level Go
ops. Peephole `SkipUntil*` instructions and packed multi-character checks are
search optimizations, so they lower to ordinary Go character/class operations
instead of being copied. The aligned byte operand format, bit-table storage,
compile-time C++ operand traits, and bytecode disassembler are V8-internal
representation/debugging facilities; the Go VM uses typed instruction structs
and has no bytecode serialization API. Matching behavior is accounted for in
the interpreter review above; the exact V8 encoding is deliberately not ported.

### Bytecode-generator source review: `regexp-bytecode-generator.cc`, `.h`, and `.inl.h`

All three files have been read in full. The Go `compileNode` pass emits
high-level instructions directly from the parsed AST instead of implementing
V8's macro-assembler interface and byte-buffer writer. V8 label fixups map to
resolved Go program-counter edges; capture/register operations map to capture
instructions; forward/reverse and case-insensitive backreference variants map
to one operation parameterized by flags and direction. Greedy/lazy branch
ordering, empty-loop termination, global iteration, and zero-width advancement
are handled in the Go VM and covered by tests. V8's byte alignment, bounded
buffer growth, raw stream copying, labels with serialized fixups, character
load widths, `AdvanceCurrentPosition` fusion, and peephole `SkipUntil*` rewrite
are representation or speed optimizations and are replaced by Go slices and
ordinary operations. The Go compiler does not reproduce V8's graph-analysis
optimizations, native fallback protocol, or bytecode serialization, and parser
or Unicode gaps remain; these files stay `[~]`.

### Bytecode-iterator source review: `regexp-bytecode-iterator.cc`, `.h`, and `.inl.h`

All three files have been read in full. They advance over V8's packed byte
stream by operand-derived instruction sizes, optionally assert that an offset
is at an instruction boundary, dispatch visitors by opcode, and repair raw
pointers after V8 heap movement. The Go representation is a GC-managed slice
of typed instructions and is iterated directly; there is no encoded byte
offset, movable trusted byte array, or visitor-template dispatch to port. This
is representation/heap integration only and is explicitly replaced by Go's
slice iteration, with no JavaScript-visible behavior omitted.

### Native code-generator source review: `regexp-code-generator.cc` and `.h`

Both files have been read in full. The visitor translates each V8 bytecode
operand to a macro-assembler operation: position/register/backtrack behavior,
character and range checks, anchors, backreferences, and capture state map to
the Go VM's corresponding instructions and execution-state operations. The
peephole-only `SkipUntil*` instructions are search optimizations; pure literal
searches use the Go KMP atom path, while other patterns execute the ordinary
semantic operations without those fused scans. The 128-entry bit tables and
optional SIMD nibble tables encode the same ASCII membership predicate and are
not needed by the Go class matcher. Label pre-visitation, indirect-jump CFI
targets, V8 `Zone`/handles, byte-array conversion, and architecture macro
assembly emit native machine code and have no Go-interpreter counterpart.
Unsupported-bytecode reporting belongs to V8's optional native backend; the Go
compiler builds its own opcodes directly from the AST. This pair is fully
accounted for as a deliberately omitted native backend, with matching behavior
provided by the Go interpreter rather than a native-code port.

### Macro-assembler source review: base interface and shared implementation

`regexp-macro-assembler.h` and `.cc` have been read in full. The shared
operation contract maps to the Go matcher as follows: UTF-16 character loading
and signed position movement to `inputCharacter`/`previousInputCharacter` and
VM position updates; register and backtrack operations to captures and saved
execution states; character/range/standard-class checks to the class and
assertion matchers; backreferences to `consumeAt`; global restart and empty
match advancement to `FindAll`; and Unicode surrogate advancement to
`AdvanceStringIndexUTF16`. The 32,767-capture maximum is now enforced. The
case-insensitive compare helper exposed an important width rule: in Unicode
mode, backreference comparison must advance the capture and input by equal
UTF-16 lengths at each code point. The Go matcher now checks this in forward
and reverse comparisons, with a regression for a Unicode 18 supplementary
simple-fold pair. V8's `IsCharacterInRangeArray` also uses a binary search over
interleaved bounds; the Go parser now sorts and merges each class interval list
once, and matching uses binary search over those intervals rather than scanning
every range.

The shared `SkipUntil*` scanners, table-switch interface, Latin1/UC16 native
specialization, range-array cache/encoding, generated-code subject flattening,
interrupt/GC repair, and architecture stack growth are machine-code/runtime
details. Pure literal scanning is handled by the Go KMP path; general scans
fall back to the semantic VM. The bounded Go APIs report explicit resource
errors, but V8's configurable backtrack counter and experimental-engine
fallback are not equivalent yet, and the ordinary Find wrappers remain
unbounded. The base pair therefore stays `[~]`.

`regexp-macro-assembler-arch.h` only selects an architecture-specific native
assembler and records unaligned-load support; it is fully accounted for as
compile-time platform glue with no role in the Go interpreter. The tracer
`.cc/.h` pair wraps assembler calls to print operands and generated-code
diagnostics and is likewise fully accounted for as instrumentation-only.

### Peephole optimizer source review: `regexp-bytecode-peephole.cc` and `.h`

Both files have been read in full. The trie, operand remapping, preservation of
branch targets, and jump fixups combine known bytecode sequences into
`SkipUntil*` search opcodes. Those are interpreter throughput optimizations;
they do not accept additional syntax or alter the match result. The Go engine
does not emit V8's load/check/advance sequences or store bytecode jump offsets,
so this rewrite layer has no direct input. The semantic checks are retained in
the ordinary Go operations, while this private peephole encoding/fixup machinery
is intentionally omitted. Any performance loss from not having equivalent
Go-side search optimizations is tracked separately from conformance.

### Stack source review: `regexp-stack.cc` and `.h`

Both files have been read in full. V8 owns a thread-local regexp stack with a
1 KiB static area, dynamically grows it up to 64 MiB, preserves its relative
stack pointer across growth, resets empty scopes, and archives thread-local
state when V8 threads switch. Native assemblers access its stack/limit fields
through fixed offsets; sandbox builds have page-allocation and registration
rules. The Go VM uses heap-managed execution-state slices instead of native
assembler stack memory, so those allocation/archival details do not apply.
`FindWithLimits`, `FindUTF16WithLimits`, and `FindAllWithLimits` now expose
step, estimated active/saved VM-state memory, and nested-assertion-depth ceilings,
returning a typed error distinct from no-match. Global matching shares one
budget, and the linear engine rejects inputs whose conservative instruction
work bound exceeds the step ceiling. The estimate is not V8's exact stack
accounting; legacy `Find` APIs remain unbounded, and JS/runtime callers have
not adopted the bounded methods. Thus the safety behavior is currently an
opt-in standalone API, and these source rows remain `[~]`.

### Result-vector source review: `regexp-result-vector.cc` and `.h`

Both files have been read in full. They manage a per-isolate reusable small
`int32` match-offset vector and dynamically allocate/free larger vectors,
including sandbox-specific ownership bookkeeping. Go represents these results
as `Match` plus GC-managed `[]Span` values, so explicit pooling and sandbox
allocation are not observable parts of matching and are replaced by Go memory
ownership. The project-level JS `/d` result-object formatting is still outside
this low-level matcher and remains on the Goja adapter.

### RegExp API source review: `regexp.cc` and `.h`

Both files have been read in full. The source escaping used for a compiled
pattern is adapted by `EscapedSourceUTF16`/`RegExpStringUTF16`, and capture-free
literal patterns now use a linear-time KMP search over UTF-16 code units,
matching V8's separate atom-search boundary without copying its platform
`StringSearch` machinery. `FindAll` implements global iteration and Unicode
zero-length advancement, while `GroupIndices` and `Match.Captures` carry the
low-level named-group and offset data.

The remaining `regexp.cc/.h` responsibilities include converting compile
failures to JavaScript `SyntaxError`s, the JSRegExp object and dynamic
`lastIndex` behavior, `exec`/match-info integration, caching and global batch
execution, automatic backtrack-limit fallback to the experimental engine,
compile-time AST/graph analysis, one-byte/two-byte code caches, tier-up, and
native code generation. These depend on V8's object model or remain missing
from the standalone package; RegExp object operations still go through Goja.
The API pair therefore remains `[~]`, despite the atom/source helpers being
implemented.

The first implementation step is not a line-for-line C++ rewrite. The source
tree mixes JavaScript-visible algorithms with V8 heap/isolate integration,
diagnostics, graph/JIT code generation, and architecture-specific machine-code
backends. The port must map every upstream file to either (a) behavior to
implement in the Go semantic core, (b) V8-specific machinery to replace with a
Go-owned equivalent, or (c) non-semantic codegen that is deliberately omitted.
That mapping and semantic parity still need to be completed. The existing
inventory lists the upstream snapshot; its unchecked boxes mean the files have
not yet all been reviewed and accounted for.
