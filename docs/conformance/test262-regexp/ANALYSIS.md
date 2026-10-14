# Test262 RegExp checkpoint analysis

This is an interpretation of the saved 1,782-variant checkpoint. It does not
claim results for the unrun variants or for the full 53,616-test suite.

## Why the run took so long

The runner emitted and built a standalone Go program for each variant. Even a
small `/a/` case contained about 1.14 MB and 20,119 lines of shared runtime;
Unicode property cases grew to roughly 1.27 MB and 21,811 lines. The common
runtime was therefore recompiled thousands of times. Go builds averaged 9.37
seconds (median 8.25, 90th percentile 12.17), with several Unicode property
cases reaching the 90-second timeout. Node reference execution averaged 0.12
seconds and generated-program execution averaged 0.48 seconds, so Go compilation
was the dominant cost.

The old runner also killed the `go build` wrapper on timeout without killing its
compiler children. Those orphan compilers continued to consume CPU and memory.
Four saved build failures show the Go compiler itself was killed, all for large
generated Unicode property cases. Inventory parsing adds about ten seconds once
per invocation, but it is not the per-variant bottleneck.

## Changes made

- Go timeouts now terminate the whole process group, including compiler and
  linker children.
- The runner emits each case into a unique namespace, combines up to eight case
  functions with one shared runtime, builds the package once, then starts a
  fresh process for each case to preserve test isolation.
- A failed batch build is bisected until the bad case is isolated, so one
  unsupported case does not mark its neighbors as failed.
- Default preparation concurrency is two workers; Go package builds are
  sequential. The Go and Node executables can be selected explicitly.
- The native source digest now includes `tsc/internal/regexpengine`, and resume
  also checks the runner source digest. This prevents resuming across changes
  that the previous digest did not cover.
- Generated matcher locals that were unused after an early return are marked
  as read. A targeted rerun of the affected Unicode-set test passed both modes.
- Large character-range sets now use deduplicated UTF-16 bitmaps in emitted Go
  instead of expanding thousands of intervals into a nested boolean expression.
  Sets with at most 64 ranges retain inline or balanced comparisons.
- Numeric updates use native `.read()` only for bindings that are actually
  backed by typed cells. Dynamic bindings keep the JavaScript increment path,
  avoiding invalid native-cell code and preserving `++` coercion semantics.
- Method calls on native primitive values now box the receiver at the call
  boundary, so string methods receive their `this` value. `RegExp.prototype`
  `Symbol.matchAll` now handles custom species construction and non-global
  receivers separately from `String.prototype.matchAll`'s global check.
- Parse-negative tests now recognize a regular-expression parser diagnostic as
  the expected `SyntaxError` only when Node also confirms the parse failure.
  Unsupported Go code-generation diagnostics remain failures.

## Timing check

A targeted `RegExp.prototype.exec` slice built 14 valid generated variants in
two shared-runtime packages. The saved individual-build timings for those same
14 variants sum to 110.36 seconds. The two batch builds took 13.48 seconds in
total, about an 8.2x reduction in Go build wall time. The 16-variant slice had
two emitter failures, which were not included in either build total. This is a
small check, not a guarantee for large Unicode-property batches; the batch
fallback is intended to split those if compilation fails or exceeds its limit.

## Failure patterns

The checkpoint's raw counters are 1,022 pass, 448 runtime-fail, 240 emit-error,
38 build-error, 8 build-timeout, and 26 host-adapter limitations. The raw
`emit-error` counter overstates implementation gaps: 196 rows are parse-negative
RegExp tests whose diagnostics are parser errors and whose Node reference
result is an expected `SyntaxError`. Under the corrected classification, the
checkpoint reads 1,218 pass and 44 genuine emitter failures; the saved JSONL
retains the original labels so its raw provenance is visible.

The remaining failures fall into a few larger groups:

- **Regex semantics:** 104 runtime failures are in generated `v`-mode set
  operation cases. Their assertion messages show class union, intersection,
  subtraction, and `\q{...}` string members as the main cluster.
- **RegExp object/runtime behavior:** many failures are receiver checks and
  missing or incomplete built-ins (`RegExp` accessors, `Symbol.split`, legacy
  state, and `RegExp.escape`). These are not all matcher-graph failures.
- **Code generation/build:** the historical checkpoint had 28 unused matcher
  locals and six typed-cell lowering failures. A focused rerun of all 17
  affected test files (34 strict/sloppy variants) now passes all 34. Four other
  builds were killed under memory pressure, and eight timed out on large
  Unicode data in the historical checkpoint.
- **Parser/codegen boundary:** after correcting expected parse-negative tests,
  44 emitter failures remain. The largest groups are unsupported Unicode set
  code generation, Unicode ranges outside the current UTF-16 matcher path, and
  class inheritance limits.
- **Host APIs:** 26 tests need Test262 `$262` realm, agent, GC, or detachment
  behavior and are not matcher results.

## Next order

1. Reduce generated code size for Unicode properties of strings. The property
   matcher for the RGI Emoji 14.0 Test262 case still exceeds a 120-second Go
   build when compiled alone. Its matcher adds over 118,000 lines after the
   shared runtime marker; the 3.6 KB JavaScript test file is not the main source
   of that expansion. A runtime package split will not remove this matcher cost.
2. Work through `v`-mode set matching in layers: character set operations,
   string-set operations, then Unicode case-folding interactions.
3. Address remaining RegExp constructor/prototype/accessor semantics separately from the
   graph matcher, then retry the corresponding Test262 groups.
4. Re-run the saved selection after those grouped fixes.

The focused rerun after the current changes passes all 34 variants from the
17 files that previously had generated-Go build errors. It also passes both
modes for `Assigned`,
`ID_Continue`, `Script_Extensions=Unknown`, `Script=Unknown`, and
`General_Category=Letter`; their shared builds took 25–36 seconds. The two
`Symbol.matchAll` species cases and the 65,536-code-unit non-whitespace test
also pass. The RGI Emoji 14.0 matcher still times out even when isolated, so it
remains a distinct code-size problem.

`go test ./internal/goemit`, `go test ./internal/regexpengine`, the Test262
runner unit tests, and `git diff --check` pass with the current source. The
historical JSONL remains unchanged so its checkpoint counts stay reproducible;
the focused rerun was written under `/tmp/test262-build-error-rerun-all`.

The 1,782-row raw checkpoint was produced before the runner began hashing the
regex-engine source and runner code. Keep it as a historical baseline; the new
runner intentionally refuses to resume from that checkpoint.
