# Launch

Launch is a TypeScript-compatible language and ahead-of-time compiler that emits
native Go programs. Its compiler builds on TypeScript's parser and checker, with
Launch syntax for native workers, control-flow headers, operator overrides, and
slice indexing.

This is a work in progress, not a claim of complete compatibility today. The
[Go target documentation](examples/go-target/README.md) and
[implementation plan](examples/go-target/IMPLEMENTATION_PLAN.md) describe current
coverage and remaining work. JSX is outside the current scope.

## Build and run

Install Go 1.27 or later. From this repository:

```sh
go run ./tsc/cmd/launch --project examples/go-target
go build -o app examples/go-target/out/input.go
./app
```

For your own project, set `"target": "go"` in `tsconfig.json`, or use
`--target go`. The emitter produces Go source; the Go compiler already turns it
into machine code. Compilation is ahead of time, with no JavaScript JIT required.

## Native types and fast paths

Alongside ordinary TypeScript types, Launch supports:

- `float32`; ordinary `number` uses Go's 64-bit floating-point storage.
- `int8`, `int16`, `int32`, and `int64`.
- `uint8`, `uint16`, `uint32`, and `uint64`.
- Native-width `int` and `uint`, plus `float64` as an alias for `number`.

When Launch knows the types, it selects native Go values and operations.
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

## Control-flow headers

Launch allows `if`, `while`, and `for` headers to omit parentheses when the
statement body is a brace-delimited block. The following `{` ends the condition
or loop header, so expressions can use the ordinary operator rules without a
separate line-based parsing rule:

```ts
if ready && count < limit {
    count++;
}

while count < limit {
    count++;
}

for let i = 0; i < count; i++ {
    console.log(i);
}

for const item of items {
    console.log(item);
}
```

Parentheses remain valid. An unparenthesized header requires a block body;
parenthesized headers can still use a single statement body. `do ... while`
also accepts an unparenthesized condition.

## Operator overrides

Classes can override arithmetic, bitwise, comparison, equality, logical, comma,
assignment, and update operators. A binary override declares both source
operands as parameters. Dispatch checks the left operand's type first, then the
right operand's type, while preserving the original argument order. This lets
one class define overloads for tensor–tensor, tensor–number, and number–tensor
operations. Unary overrides take no parameters. Overloads for the same operator
can use ordinary TypeScript method overload declarations.

When the emitter has a concrete Go class receiver for an override, it calls the
generated implementation directly, bypassing the runtime method table and the
generic `tsGet`/`tsCallReceiver` path. Base-typed receivers use a generated
layout switch to preserve virtual dispatch to subclass overrides. Unknown
runtime classes fall back to dynamic property lookup.

Unary `++` and `--` overrides also take no parameters and return the updated
value, which is stored back into the operand. A postfix expression evaluates
to the value captured before the update; update overrides should return a new
value when they need that old-value behavior.

```ts
class Vector {
    constructor(public x: number, public y: number) {}

    operator override +(left: Vector, right: Vector): Vector;
    operator override +(left: Vector, right: number): Vector;
    operator override +(left: number, right: Vector): Vector;
    operator override +(left: Vector | number, right: Vector | number): Vector {
        if (left instanceof Vector && right instanceof Vector) {
            return new Vector(left.x + right.x, left.y + right.y);
        }
        if (left instanceof Vector) {
            return new Vector(left.x + right, left.y + right);
        }
        if (right instanceof Vector) {
            return new Vector(left + right.x, left + right.y);
        }
        throw new TypeError("Vector addition needs a Vector operand");
    }

    operator override -(left: number, right: Vector): Vector {
        return new Vector(left - right.x, left - right.y);
    }
}
```

For example, `operator override -` can define `number - Vector` separately from
`Vector - number`; the overload method receives the original operands in source
order even when dispatch uses the right operand as its receiver. Compound
assignments can define their own override, such as `operator override +=`; when
they do not, they use the matching binary override and assign its result back
to the left-hand side. `operator override =` receives both operands and is
responsible for updating its receiver; the assignment expression evaluates to
the left-hand value.

Overloaded `&&`, `||`, and `??` evaluate both operands before calling the
override, as other overloaded binary operators do. Without an override they
keep their normal short-circuit behavior. Conditional expressions and keyword
operators such as `typeof`, `delete`, `void`, and `await` are not overloadable.

## Slice indexing

`operator override slice` receives one `Slice[]` argument. A single number
becomes an index item; Python-style `start:stop:step` becomes a range item. Commas
separate dimensions, and omitted bounds or steps are `undefined`:

```ts
class Tensor {
    operator override slice(parts: Slice[]): Tensor {
        // Each item is a SliceIndex or SliceRange.
        return this.select(parts);
    }

    select(parts: Slice[]): Tensor {
        return this;
    }
}

const tensor = new Tensor();
const row = tensor[4];
const window = tensor[1:12:2, :, 0];
```

When a value's type defines `operator override slice`, every bracket access on
that value uses slice dispatch, including `tensor[0]`. Range parts retain their
typed `start`, `stop`, and `step` values; the range is half-open, and omitted
components remain `undefined` for the tensor implementation to interpret.

## The `launch` keyword

`launch worker(args)` starts a Go goroutine and returns a Promise for the result:

```ts
function square(value: number): number {
    return value * value;
}

async function main() {
    const results = await Promise.all([launch square(6), launch square(7)]);
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

| Selected workload | Launch → Go | Node.js | Bun | Speedup vs Node / Bun |
| --- | ---: | ---: | ---: | ---: |
| Initialize/fill 4 million byte elements | 6.07 ms | 74.24 ms | 16.36 ms | 12.2× / 2.7× |
| Eight parallel array-initialization tasks | 25.79 ms | 298.33 ms | 203.24 ms | 11.6× / 7.9× |

Measured on 2026-10-06 with Go 1.27.1, Node 24.19.0, and Bun 1.4.2.
Medians of five complete-process runs; compilation excluded, startup and setup
included. The container has a two-core CPU quota. The parallel row uses eight
goroutines here and eight real worker threads in Node and Bun.

These byte-array workloads benefit from native storage as well as execution and
worker-startup costs; the ratios do not isolate the `launch` keyword alone.
Node and Bun erase the `int8` annotation and execute the same ordinary Array
source. This is not a comparison with JavaScript `Int8Array`.


## To Jarred Sumner, creator of Bun

Jarred, we hope you might consider adding WebKit/JavaScriptCore bindings to this
project and exploring its use in Bun. The aim is true, type-aware Launch
execution with native types and ahead-of-time compilation, beyond stripping
type annotations from JavaScript.

Would you be willing to fork the Go compiler and integrate this emitter with it
to provide a direct Launch-to-machine-code toolchain? Emitted Go already
compiles to native machine code today; tighter compiler integration could make
that path more seamless. We hope you can help make it happen.

## Node module compatibility

Native Go implementations of `fs` and `fs/promises` are in progress, using Bun as a behavior reference. Common filesystem operations, descriptors, directory iteration, watchers, and basic file streams are implemented. Full Node compatibility is not yet complete; see the [module checklist](docs/node-compat/README.md) and [filesystem coverage](docs/node-compat/fs.md). The current filesystem provider targets Linux.
