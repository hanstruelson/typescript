# Dynamic value representation benchmark

Measured on 2026-10-05, Go 1.27.1, Linux amd64, Intel Xeon Platinum 8573C
virtual CPU. One Go execution thread (`GOMAXPROCS=1`, `-cpu=1`). Each result is
the median of five 200 ms runs. These are representation microbenchmarks, not
measurements of complete transpiled programs or npm packages.

The backend now uses the Inline24 layout for dynamic storage. The recommendations
below record the original experiment; they do not describe the current ABI.
Complete generated-program performance still needs measurement as legacy
interface adapters are removed.

## What was compared

- **Any (16 bytes):** the Go interface compatibility baseline. Numbers may need
  separate storage when returned through a non-inlined dynamic call or written to
  an escaping dynamic cell. Object pointers do not require that numeric box.
- **Inline24 (24 bytes):** our own inline tag, float64 payload, and a dedicated
  `unsafe.Pointer` field. The pointer is visible to the collector; it is never
  encoded into integer bits. Correct tag/pointer pairing remains our responsibility.
- **Header16 (16 bytes):** our own float64 payload and pointer to a tagged header.
  Numeric values need no header; null/undefined/boolean use shared sentinels.
  Other references have embedded headers, costing object memory and indirection.
- **Safe32 (32 bytes):** our own inline tag and float64 payload, with Go's `any`
  used only for references. Primitive numbers are unboxed. Reference assertions
  remain checked by Go, with no unsafe casts. This is a custom dynamic container,
  even though it delegates reference storage to a Go interface.
- **Native:** ordinary float64 values, where relevant.

The benchmark string/object types are deliberately small. They isolate value
representation and do not implement all runtime semantics. Candidate semantics
checks preserve separate null/undefined, float64 bit patterns, and references
across forced garbage collections. Object-read factories prevent compile-time
proof of the dynamic type/tag. Inputs are constructed outside read/dispatch timing.

## Median time per operation

| Workload | Native | Any | Inline24 | Header16 | Safe32 |
| --- | ---: | ---: | ---: | ---: | ---: |
| TypeCheck | — | 2.04 ns | 1.47 ns | 1.66 ns | 2.38 ns |
| NumericRead | 3.11 ns | 3.11 ns | 5.76 ns | 6.05 ns | 5.14 ns |
| NumericUpdate | 3.33 ns | 24.18 ns | 6.80 ns | 3.77 ns | 3.94 ns |
| MixedDispatch | — | 3.77 ns | 3.26 ns | 3.29 ns | 3.09 ns |
| Call | 3.03 ns | 16.99 ns | 2.30 ns | 2.30 ns | 2.59 ns |
| ObjectRead | — | 3.35 ns | 3.40 ns | 3.97 ns | 3.13 ns |
| BuildNumbers | 2.41 µs | 23.78 µs | 7.02 µs | 4.79 µs | 9.70 µs |
| BuildObjects | — | 5.31 µs | 5.32 µs | 9.11 µs | 5.90 µs |

BuildNumbers creates a slice containing 1,024 numbers. BuildObjects creates
256 new objects plus their dynamic-value slice. The other rows measure one
operation. Call uses a non-inlined two-argument addition to expose the calling
convention and returning a dynamic number, without an application function table.
NumericUpdate writes a heap-escaping cell. It models mutable captured/state-machine
storage, not every possible local variable.

## Allocations

| Workload | Any | Inline24 | Header16 | Safe32 |
| --- | ---: | ---: | ---: | ---: |
| NumericUpdate | 8 B / 1 allocs | 0 B / 0 allocs | 0 B / 0 allocs | 0 B / 0 allocs |
| Call | 8 B / 1 allocs | 0 B / 0 allocs | 0 B / 0 allocs | 0 B / 0 allocs |
| BuildNumbers | 26,624 B / 1,025 allocs | 27,264 B / 1 allocs | 18,432 B / 1 allocs | 32,768 B / 1 allocs |
| BuildObjects | 6,912 B / 257 allocs | 8,576 B / 257 allocs | 11,008 B / 257 allocs | 11,520 B / 257 allocs |

## Recommendation

Keep native float64/bool/UTF-16 references for statically typed code. For a custom
dynamic representation, start with **Safe32**, not a packed pointer-in-integer
scheme or the header design. It removes the numeric boxing allocations in these
workloads and retains checked, GC-managed references. The exact tag-check cost is
not a compelling reason to switch: it is already small, and Safe32's type check
was slower than Any's here.

The measured tradeoff is real: Safe32 doubles the value-container size, uses more memory for reference-heavy arrays. See the object-creation timings
above: the custom designs do not have a universal throughput advantage.
Inline24 is a potentially useful later optimization, but puts pointer casts and
representation invariants into our runtime. It has not been established as the
best production choice by these isolated tests.

Before changing the production ABI, compare complete numeric-heavy and
reference-heavy transpiled workloads with equivalent helpers and optimization
settings. Carry the custom value consistently through cells, calls, arrays,
objects, promises, and modules; repeatedly wrapping it in Go interfaces would
undermine the allocation benefit. Use the Any implementation as the compatibility
baseline during migration. No production representation was switched for this
benchmark experiment.

These measurements establish that avoiding numeric boxing can matter much more
than the type check. They do not establish one universal winner on all CPUs,
all compilers, or all npm applications. Absolute nanosecond timings include loop
cost and VM noise; raw samples and allocation counts are supplied for review.

## Reproduce

From the repository root:

```sh
GOMAXPROCS=1 go test ./tsc/internal/goemit -run '^TestValue' \
  -bench '^BenchmarkValue' -benchmem -benchtime=200ms -count=5 -cpu=1
```

Authoritative combined samples: [raw-combined.txt](raw-combined.txt).
Preliminary samples: [initial candidates](raw-values.txt),
[checked reference candidate and revised object-read measurements](raw-safe-values.txt).
Machine-readable medians: [summary.json](summary.json).
Benchmark source: `tsc/internal/goemit/value_bench_test.go`.
