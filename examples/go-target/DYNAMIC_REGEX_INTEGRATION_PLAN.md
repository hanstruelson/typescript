# Dynamic RegExp integration plan

## Scope

This plan moves the temporary Go-target JIT regex path from a focused fixture
to real imported applications such as Express. It is separate from the general
Go `jit.Eval` lifetime and linker plan: passing this plan does not make
arbitrary evaluated Go code safe to unload.

## Verified baseline

Commit `4b0ed16f1` on `work` and `main` contains a runtime JIT regex bridge for
Linux/amd64. A dynamic pattern reaches the RegExp runtime at execution time,
the runtime generates a matcher expression and calls `jit.Eval`, and the
resulting matcher is retained by the RegExp object. The source string is not
parsed into a matcher during TypeScript transpilation.

The `dynamic-regex-stress` fixture spans four TypeScript files. It exercises
dynamic constructors in an imported factory and directly in the caller,
capture extraction, a reused RegExp with 16 calls to `.compile`, global
`lastIndex` behavior, and mutable state shared through an imported module. The
transpiled Go executable was run successfully three times. Each run produced
32/32 factory matches, 32/32 direct matches, 16/16 recompilation matches,
3/4 global checks, and 83/83 shared-state matches. The fixture's call graph
requests 82 runtime compilations per process; that count is derived from the
fixture rather than reported by JIT instrumentation.

This proves a narrow runtime execution path with imported TypeScript modules.
It does not prove broad ECMAScript RegExp compatibility, Express compatibility,
general imported Go package linking, GC-managed unloading, or support beyond
Linux/amd64 with the custom Go toolchain and `gojit` build tag.

## Ordered integration work

1. **Make the target application reproducible.** Check in or identify the
   actual TypeScript/JavaScript entry points, exact dependency versions and
   lockfile, TypeScript configuration, and one command that runs the
   transpiler. The tracked `tsc/.express-experiment/main.go` is a generated
   10 MB file; it contains the older RegExp runtime shape and has no
   `TsDynamicRegExpCompiler` bridge. It is not a valid JIT integration test or
   a reproducible source input. Regenerate the Go output from application
   sources after the bridge is present.
2. **Trace dynamic regex creation in the real dependency graph.** Record each
   path, the source and flags supplied, whether creation happens during route
   registration or request handling, and which APIs consume the result
   (`exec`, `test`, `match`, `replace`, `search`, `split`, or others). Separate
   dynamic constructors from regex literals, which continue through the
   ahead-of-time path. Do not infer the application's requirements from a
   package name alone.
3. **Build a differential regex corpus from those call sites.** For every
   required pattern shape, compare the source JavaScript runtime and the Go
   target on matches, captures, named captures, flags, `lastIndex`, empty
   matches, and invalid-pattern errors. Include route-like separators,
   optional/repeated segments, escaped characters, and only the flag and
   syntax features the actual dependency uses. Keep unsupported cases as
   explicit failures rather than silently changing their meaning.
4. **Exercise the real import path.** Put dynamic pattern construction in
   imported helper modules and use the resulting RegExp from the importing
   application. Verify that imported modules initialize once, their shared
   state remains shared, and the same generated executable runs the JIT
   matcher without a compiler subprocess, source file, or Go installation at
   runtime. Retain the current four-file fixture as the small regression
   baseline.
5. **Add test-only JIT accounting.** Count successful and failed runtime
   compilation requests and executable bytes. Assert that a fixed set of
   routes compiles during setup and repeated requests reuse those RegExp
   objects without compiling again. This distinguishes route setup from
   request-time compilation and gives later memory work a measurable baseline.
6. **Run the smallest real application end to end.** Transpile and build it
   with the custom Linux/amd64 toolchain and `-tags gojit`, start the generated
   executable, send representative requests, and verify status, route
   selection, captures, and response output. Preserve the exact inputs and
   commands so this remains repeatable after compiler changes.
7. **Set an operational bound while unloading is absent.** Each successful
   dynamic compile retains executable memory until process exit; `.compile`
   can replace the handle without retiring the old mapping. Do not compile
   unbounded user-controlled patterns in a request loop. For the temporary
   integration, constrain compilation to bounded application setup or enforce
   a clear compilation limit and report it. Measure route-count startup cost
   and memory before expanding the workload.
8. **Use failures to sequence runtime work.** Fix only the demonstrated
   semantic or module integration gaps, add each to the differential corpus,
   and keep the broader owner/GC, symbol manifest, relocation, and unload
   obligations in `jit-in-memory-linker-plan.md`. This regex bridge is a
   restricted client of that work, not evidence that the general lifetime
   contract is complete.

## Exit criteria for the Express checkpoint

- The real application sources and dependency graph are reproducible from a
  clean checkout.
- All dynamic regex forms it uses pass the differential corpus on Linux/amd64.
- The transpiled application imports multiple modules and runs in one
  self-contained executable with runtime-only pattern compilation.
- Route setup has a known compile count and memory bound; repeated request
  handling does not create new matchers for already-compiled routes.
- Unsupported platforms, syntax, and exhausted compile limits fail clearly.
- The no-unload limitation remains documented until the separate GC owner and
  retirement protocol is implemented and verified.
