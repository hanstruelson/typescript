# Go target

Use Go 1.27 or later. From the repository root:

```sh
go run ./tsc/cmd/tsc --project examples/go-target
go run examples/go-target/out/input.go
```

The numeric example prints `10`. The regular CLI also accepts `--target go
input.ts`. The target uses the existing TypeScript checker and output handling,
including `outDir`, `noEmit`, and `noEmitOnError`. Declaration output remains
available. Projects with local imports or exports emit one Go executable containing a
persistent state machine for each module. Its filename follows the first entry
module, and `outDir` still controls the destination.

## Timers

```sh
go run ./tsc/cmd/tsc --project examples/go-target/timer
go run examples/go-target/timer/out/main.go
```

Output:

```text
sync 0
promise 0
timer 1
```

`setTimeout(callback, milliseconds, ...arguments)` registers a native timer and
returns an opaque handle. Only the worker sleeps; the callback runs on the event
loop and can safely mutate captured variables. Missing, negative, nonfinite, or
out-of-range delays become zero; finite delays are truncated to milliseconds.
Timer callbacks run in deadline order, with Promise reactions drained between
callbacks. `clearTimeout(handle)` suppresses the callback and wakes a canceled
worker, so canceled long timers do not keep the process alive.

## Async file reading

```sh
go run ./tsc/cmd/tsc --project examples/go-target/async
go run examples/go-target/async/out/main.go
```

Output: `sync`, `10 contents`, then `0`, `1`, and `2`, on separate lines.
The example's declaration file describes the text-file adapter implemented by
this target; it is not a replacement for full Node.js type definitions. Named
`readFile` imports from `node:fs/promises` or `fs/promises` require an explicit
`utf8` encoding. Ambient intrinsics `readFile(path)` and `delay(milliseconds)`
return Promises for text and timer completion respectively.

## Lowering and runtime ownership

Async functions run their synchronous prefix immediately and return a Promise.
Numbered resume blocks preserve branches, loops, and awaits. The compiler's
reference resolver identifies bindings by declaration rather than spelling.
Binding cells preserve shadowing, initialization timing, undefined `var` values,
and shared mutable captures. Blocks allocate fresh lexical cells; `for (let ...)`
creates fresh iteration bindings; closures capture the appropriate cells.
Protected-region and pending-completion stacks preserve return, rejection,
`break`, and `continue` through `try`/`catch`/`finally`, including awaited cleanup.

Each event loop runs its callbacks and async function steps on one goroutine.
Native tasks publish completion through that loop's task channel. User workers
have separate event loops; their reference values and captured bindings remain
shared. Users are responsible for avoiding concurrent mutation of those values.

After synchronous execution, the loop drains Promise reactions and completed
timers. If tasks remain, it waits on the completion channel. It exits when ready
work and registered tasks are empty. Pending Promises without active operations
do not keep the process alive.

The worker wrapper converts recoverable panics and early worker termination into
rejected results and always publishes completion. Promise failures reach await
catch handlers or rejection callbacks. Uncaught callback failures are collected
while remaining tasks drain; they and unhandled rejections produce a clear
stderr message and exit status 1. This does not protect against Go's unrecoverable
process failures, such as out-of-memory termination. Native adapters must not
launch unguarded goroutines or mutate application-owned state.

## Supported subset

The backend supports primitive values; identifier declarations; ordinary and
async functions and arrows; arrays; assignments; numeric arithmetic and
comparisons; strict equality; short-circuit and conditional expressions; `if`;
`for`, `while`, `do`, and array/string `for-of`; labeled loop exits; return and
throw; and `try`/`catch`/`finally`. Promises support an executor, `resolve`,
`reject`, array `all`, `then`, and `catch`. Every output includes the runtime.

Objects, destructuring, spread, classes, regular expressions, optional chaining,
named function expressions, logical assignments, exponentiation, bitwise
operators, loose equality, property membership, and numeric/string enums are
also supported. Enums include numeric reverse mappings and declaration merging.
Array elisions and gaps created by index or length assignments remain absent
properties. Iteration yields `undefined` for those slots, while spread creates
present `undefined` elements; property membership and key enumeration distinguish
these cases.

For `number`, bitwise operations truncate and wrap operands to 32 bits; shift
counts are masked to five bits, and `>>>` returns an unsigned 32-bit result as a
`number`. Known `number` operands use native Go numeric helpers, including
compound assignments and exponentiation. Dynamic operands perform JavaScript
primitive conversion first. Explicit native integer types retain their native
width and Go-style shift behavior.


Union types retain TypeScript's control-flow narrowing. Intersection objects use
one flat property store and preserve object identity. Generic functions, classes,
and aliases are checked by TypeScript. Eligible named synchronous generic
functions now receive cached concrete Go implementations for checked call-site
types. Primitive arguments, locals, and results stay native; TypeScript `any`
uses `tsValue`. Small bodies emit direct Go expressions. More complex supported
bodies reuse the lowering machinery with concrete substitutions, including
indexed-access and conditional types. Captures are passed as binding pointers.
This is our own specialization implementation, rather than Go type parameters.

Generic classes, first-class function calls, reassigned functions, async generic
functions, and bodies with unsupported specialization control flow retain the
existing tagged ABI. Structural objects still use the existing flat property
store; specialization does not yet create native structs for their fields.
The specializer limits expansion to 64 implementations per function and 16 nested
specializations, then falls back to the ordinary function ABI. Primitive templates
share implementations by Go representation; dependent bodies use checked type
identity to avoid conflating different shapes. `keyof`, indexed access, mapped
and conditional types, and `satisfies` remain supported.

Generators, async iterators, decorators, namespaces, and several advanced class
features still need backend work. This is not a complete JavaScript or Node.js
runtime: additional built-ins and host APIs remain incomplete. JSX and dynamic
code execution are excluded from the project scope. Timer handles are opaque,
not browser numeric IDs.

## Demand-driven modules

```sh
go run ./tsc/cmd/tsc --project examples/go-target/modules
go run examples/go-target/modules/out/main.go
```

This target implements demand-driven module execution: a named import requests
only the selected exports. Each module has an export table with readiness flags,
one saved machine, and per-request missing-name sets. Export publication suspends
the module. When a request's set becomes empty, its importer resumes and the
module remainder stays suspended unless another outstanding request needs it.
Cached exports return without restarting the module. Imports hold live binding
references, so later assignments remain visible.

A side-effect import requests full completion and drives the remaining states.
Namespace imports also request completion. Relative module paths use TypeScript's
normal resolver, including `.js` specifiers resolving to `.ts` sources. Named and
default imports, named exports, default exports, and explicit named re-exports
are supported. Use explicit type-only imports for purely type dependencies.
Export-star currently requires an explicit export list.

The circular example publishes `B` before importing `A`. The helper can therefore
read cached `B`, publish `A`, and stop. Main receives `A` before the helper's
remainder executes; its later side-effect import resumes that remainder. The
intended output is:

```text
helper begins
main 30 20
helper remainder
main complete
```

A circular dependency with no ready export to satisfy its requests is reported
as an unresolved module dependency, rather than restarting modules or hanging.
These demand-driven semantics intentionally follow this backend's module policy;
they differ from standard JavaScript's eager module evaluation.


## Classes

```sh
go run ./tsc/cmd/tsc --project examples/go-target/classes
go run examples/go-target/classes/out/main.go
```

Ordinary instance fields with `number` types become concrete `float64` Go struct
fields. String and boolean fields use `*tsString` and `bool`; nullable fields use
a native payload plus a tag, and `any` fields use dynamic storage. Inherited fields share one slot per ordinary property name. Each concrete
class receives specialized copies of inherited method bodies; `super.method()`
is resolved using the body's original declaring class, even across several
inheritance levels. Calls through `this` select the descendant's override.
Constructors retain base-to-derived field initialization order. Uninitialized
fields retain JavaScript's `undefined` behavior through initialization flags.
Simple synchronous methods emit structured Go bodies; asynchronous methods and
complex control flow reuse the existing state-machine machinery.

Known class dot access uses concrete fields/accessors and generated callback views,
without property-name hashing. Parent-typed references can hold different concrete
descendants. Bracket reads and writes use an instance-owned property table with
generated accessors into those same fields; dynamically added properties use
separate entries. No reflection is used. Declared primitive fields validate dynamic
values through either access path. Class ancestry metadata supports `instanceof`.

Static fields and methods use a separate generated struct and one shared static
object per evaluated class. Static blocks and fields execute in source order.
Inherited static fields read parent storage until a child writes its own value;
static methods retain lexical `super` and the current class receiver.

The initial class subset requires named classes and statically resolved base
classes in the same source file. Constructor parameter properties use the same
native slots as ordinary fields.
Synthetic property declarations precede regular fields; parameter values are
assigned after base field initialization or after `super()` in derived classes.
Validated native and nullable parameters write their slots directly. Constructors
and methods support rest and destructured parameters, and constructor calls
support spread arguments. Defaults preserve parameter temporal dead zones and
exact native integer literals, including declared numeric rest arguments. Accessors, ECMAScript private fields, and computed
declarations are diagnosed. Ordinary array storage still uses dynamic values.
Typed primitive bindings and
suitable ordinary functions use native payloads and signatures.


## Types and strings

Dynamic values use a 24-byte tagged container on 64-bit Go: an eight-byte
type tag, an eight-byte numeric payload, and an eight-byte GC-visible
`unsafe.Pointer`. The numeric payload remains `float64` to preserve JavaScript
Number semantics. Runtime dispatch and generated class views use tags and concrete
pointers. Go interfaces are confined to host APIs (Goja, panic recovery, formatting)
and selection of the static generic primitive type; JavaScript values are never
boxed into interfaces for internal dispatch.

Non-nullable primitive bindings use native payloads (`float64`, `bool`, or a
pointer to our UTF-16 string). Nullable primitive bindings use a payload plus a
small tag distinguishing value, null, and undefined. Dynamic values entering a
typed parameter, assignment, return, or mixed typed/dynamic comparison are
converted at runtime by default. `--coerceAny false` (or `"coerceAny": false`
in tsconfig) instead checks tags and raises a catchable TypeError for an
incompatible value. Null
and undefined are accepted only when permitted by the destination type, or
converted when coercion is enabled. Hoisted variable reads and lexical temporal
dead zones retain separate initialization tracking.

Suitable synchronous top-level functions have native Go parameter and return
signatures. Other functions use the dynamic calling convention and typed cells;
this includes async functions and functions requiring complex control flow.

Strings store immutable UTF-16 code units. Length, indexing, slicing, comparison,
search, padding, trimming, splitting, iteration, and the ordinary String methods
operate on this representation. Unicode casing and normalization use x/text on
valid scalar segments while preserving isolated surrogates. Locale collation
uses x/text; this does not provide full Intl/ICU equivalence for every locale or
collation option.

RegExp uses the Goja ECMAScript implementation and its regexp2 dependency rather
than Go's regexp package. Regex literals, constructors, exec/test, captures,
replacement callbacks, and String regex methods are supported. The engine is
initialized lazily and application source is never passed to eval or Function.
Generated files therefore need these dependencies: run them from this repository's
Go workspace, or create an output module with the dependencies in `tsc/go.mod`.
Regex conformance follows the dependency's supported ECMAScript version.

Object literals, object/array spread, binding destructuring (including defaults
and rest), rest parameters, template literals, tagged templates, for-in,
for-of binding patterns, and switch statements are supported. Object methods,
accessors, destructuring assignments, full prototype behavior, sparse-array
semantics, and async iteration remain incomplete. This is an experimental
backend and does not yet cover all TypeScript or ECMAScript features.

The `types` example demonstrates typed boundaries, UTF-16 indexing, regex
captures, and destructuring:

```sh
go run ./tsc/cmd/tsc --project examples/go-target/types
go run examples/go-target/types/out/main.go
```

The dynamic `"42"` argument converts by default. Add `--coerceAny false`
to report a type mismatch instead.

The [implementation plan](IMPLEMENTATION_PLAN.md) describes the remaining features,
performance work, dynamic representation, and proposed comparison policy.

The [dynamic-value benchmarks](benchmarks/README.md) compare Go interfaces with
custom tagged containers; the production dynamic representation is still `any`.

## Native numeric types and user workers

The parser and checker recognize `float32`, `float64`, `int`, `int8`, `int16`,
`int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, and `uint64`. `number`
continues to mean `float64`. Statically typed bindings, parameters, results, and
fields use the corresponding Go types. Dynamic values retain exact integer bits
in the 24-byte tagged representation, including all `int64` and `uint64` values.

Conversions always check representability, regardless of `coerceAny`. Invalid
constants produce an emit diagnostic; runtime overflow, fractional integer
conversion, NaN, or infinity converted to an integer raise `RangeError`.
`float32` rounds to native precision, rejecting finite values that overflow.
`coerceAny` controls conversion from nonnumeric dynamic values and defaults to
true. Same-type native integer arithmetic uses Go's integer arithmetic semantics.
Mixed native arithmetic promotes to a common native type and checks operand
conversions; explicit `number` operands select `float64`. Inferred arithmetic
results retain their native type.

```ts
declare function delay(ms: number): Promise<void>;
let count: int64 = 9223372036854775807;
async function worker(value: int64): Promise<int64> {
    await delay(1);
    return value;
}
async function main() {
    const result = await go worker(count);
    console.log(result);
}
main();
```

`go worker(args)` evaluates the function and arguments once, starts a task with
its own event loop, and returns a Promise for the worker's result or error.
Nested calls, timers, async functions, and further `go` calls use the executing
worker's loop. Scalars pass by value; object, array, function, and class references
pass by pointer. No arguments, captures, or results are deep copied.

Local literal dynamic imports such as `await import("./helper.js")` reuse the
bundle's module registry. Initialization runs once on the registry's owner loop,
protected by an internal once lock; concurrent workers wait through task
completion. Export lookup uses an internal read lock. No channels or user mutex
primitives are introduced.
