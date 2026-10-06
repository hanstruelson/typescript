# TypeScript → Go

This is a fork of TypeScript with an emitter to Go. The goal is to run all
TypeScript and JavaScript code as native programs, except runtime code generation
through `eval` or `new Function`.

This is a work in progress, not a claim of complete compatibility today. The
[Go target documentation](examples/go-target/README.md) and
[implementation plan](examples/go-target/IMPLEMENTATION_PLAN.md) describe current
coverage and remaining work. JSX is outside the current scope.

## Build and run

Install Go 1.27 or later. From this repository:

```sh
go run ./tsc/cmd/tsc --project examples/go-target
go build -o app examples/go-target/out/input.go
./app
```

For your own project, set `"target": "go"` in `tsconfig.json`, or use
`--target go`. The emitter produces Go source; the Go compiler already turns it
into machine code. Compilation is ahead of time, with no JavaScript JIT required.

## Native types and fast paths

Alongside ordinary TypeScript types, this fork supports:

- `float32`; ordinary `number` uses Go's 64-bit floating-point storage.
- `int8`, `int16`, `int32`, and `int64`.
- `uint8`, `uint16`, `uint32`, and `uint64`.
- Native-width `int` and `uint`, plus `float64` as an alias for `number`.

When the compiler knows the types, it selects native Go values and operations.
Booleans use `bool`, strings use our UTF-16 representation, and numeric types use
their corresponding Go types. Matching types need no boxing or conversion.
Guaranteed lossless widening is automatic. Other conversions use checked helpers
that report invalid input or overflow rather than silently corrupting values.

Ordinary growable arrays such as `int32[]` and `number[]` use native Go slices and
support `push`, `pop`, and other array methods. Known synchronous callbacks also
use native element and result types. These are separate from JavaScript's
fixed-length, buffer-backed typed arrays.

Nullable primitives have their own path. Nullable arrays keep native values and
one byte of presence information per element, distinguishing `null` from
`undefined` without padding every element. Eight nullable `int64` values require
72 bytes of element storage, excluding slice headers and spare capacity.

Dynamic `any` values work through our 24-byte TSValue: an eight-byte type tag, an
eight-byte inline value payload, and a separate GC-visible pointer field.
`any[]` and mixed unions use TSValue elements. Dynamic operations check the tag
and perform the required conversion; they do not store application values in
Go's `any` interface.

Implicit coercion defaults to enabled. `--coerceAny false` requires explicit
conversion when a conversion is not guaranteed lossless; safe widening still
works. Explicit methods include `.int64()`, `.int32()`, `.uint64()`, `.float32()`,
`.number()`, and `.string()`. `.toString()` remains supported. Numeric
conversions check bounds in either mode.

## The `go` keyword

`go worker(args)` starts a Go goroutine and returns a Promise for the result:

```ts
function square(value: number): number {
    return value * value;
}

async function main() {
    const results = await Promise.all([go square(6), go square(7)]);
    console.log(results[0], results[1]); // 36 49
}
main();
```

Each user goroutine has its own event loop, so callbacks within that worker run
on one goroutine while separate workers can use separate CPU cores. Workers use
the same task and Promise plumbing as the rest of the runtime. Module imports
initialize once, with concurrent initialization synchronized.

Values are passed normally: primitives by value, reference values by reference.
Objects and arrays are not deep-copied. Code sharing mutable references across
workers must avoid races.

## Benchmarks

We compare compiled output with Node.js and Bun on selected workloads intended
to show the benefits of native types and parallel workers. These are targeted
examples, not a claim that every TypeScript program runs faster. The benchmark
sources, runner, raw samples, versions, and methodology are in
[benchmarks/](benchmarks/README.md).

| Selected workload | This fork → Go | Node.js | Bun | Speedup vs Node / Bun |
| --- | ---: | ---: | ---: | ---: |
| Initialize/fill 4 million byte elements | 6.07 ms | 74.24 ms | 16.36 ms | 12.2× / 2.7× |
| Eight parallel array-initialization tasks | 25.79 ms | 298.33 ms | 203.24 ms | 11.6× / 7.9× |

Measured on 2026-10-06 with Go 1.27.1, Node 24.19.0, and Bun 1.4.2.
Medians of five complete-process runs; compilation excluded, startup and setup
included. The container has a two-core CPU quota. The parallel row uses eight
goroutines here and eight real worker threads in Node and Bun.

These byte-array workloads benefit from native storage as well as execution and
worker-startup costs; the ratios do not isolate the goroutine keyword alone.
Node and Bun erase the `int8` annotation and execute the same ordinary Array
source. This is not a comparison with JavaScript `Int8Array`.


## To Jarred Sumner, creator of Bun

Jarred, we hope you might consider adding WebKit/JavaScriptCore bindings to this
project and exploring its use in Bun. The aim is true, type-aware TypeScript
execution with native types and ahead-of-time compilation, beyond stripping
type annotations from JavaScript.

Would you be willing to fork the Go compiler and integrate this emitter with it
to provide a direct TypeScript-to-machine-code toolchain? Emitted Go already
compiles to native machine code today; tighter compiler integration could make
that path more seamless. We hope you can help make it happen.
