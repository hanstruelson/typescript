# Selected TypeScript → Go benchmarks

These benchmarks compare emitted, compiled programs with Node.js and Bun. The
headline workloads deliberately focus on native byte storage and inexpensive
worker creation. They are selected examples of wins, not a general runtime
ranking. All completed workloads from this suite appear below, including the
callback workload where Bun is faster.

## Results

| Selected workload | This fork → Go | Node.js | Bun | Speedup vs Node / Bun |
| --- | ---: | ---: | ---: | ---: |
| Initialize/fill 4 million byte elements | 6.07 ms | 74.24 ms | 16.36 ms | 12.2× / 2.7× |
| Eight parallel array-initialization tasks | 25.79 ms | 298.33 ms | 203.24 ms | 11.6× / 7.9× |

Measured on 2026-10-06 with Go 1.27.1, Node 24.19.0, and Bun 1.4.2.
Medians of five complete-process runs; compilation excluded, startup and setup
included. The container has a two-core CPU quota. The parallel row uses eight
goroutines here and eight real worker threads in Node and Bun.

| Workload | This fork → Go | Node.js | Bun |
| --- | ---: | ---: | ---: |
| Byte-array initialization | 6.07 ms | 74.24 ms | 16.36 ms |
| Eight tasks, sequential | 37.77 ms | 304.35 ms | 139.28 ms |
| Eight tasks, parallel | 25.79 ms | 298.33 ms | 203.24 ms |
| Number-array callback summation | 62.56 ms | 81.66 ms | 14.39 ms |
| TSValue-array callback summation | 556.62 ms | 81.66 ms | 14.39 ms |

For the worker workload, running with `go` rather than sequentially improves this
fork's own time by 1.46×. The much larger differences against
Node and Bun combine native storage, runtime execution, and startup costs;
they do not establish that goroutines alone produce those ratios.

## Workloads

- [compact-array.ts](compact-array.ts): allocate an ordinary growable `int8[]`
  with four million elements, fill it with 7, then print a result. Go stores
  one byte per element. Node and Bun erase the annotation and use their ordinary
  Array representation. This is not a comparison with buffer-backed `Int8Array`;
  JavaScript can explicitly use that different API to obtain compact storage.
- [workers.ts](workers.ts): eight independent tasks each allocate and fill eight
  million byte elements. The fork runs them with `go`, each on its own event loop.
  The sequential comparison replaces `go compute(seed)` with
  `Promise.resolve(compute(seed))`. Node and Bun also run a parallel version using
  eight actual worker threads via [node-workers.cjs](node-workers.cjs).
  The worker function body is extracted from the same TypeScript source using
  Node's type stripping, rather than reimplementing the algorithm.
- [arrays.ts](arrays.ts): sum a 4,096-element array through `forEach` 512 times.
  The runner emits separate `number[]` and `any[]` programs to show native versus
  TSValue dispatch. Their JavaScript is equivalent, so they share the same
  Node/Bun baseline. Bun wins this workload; the dynamic version is considerably
  slower than either JavaScript runtime.

## Measurement

The runner builds the compiler and all emitted programs before timing.
Each workload gets one untimed complete-process warm-up, followed by five timed
complete-process runs. Execution order rotates between samples. Timings measure
wall time from process launch to exit and include startup, parsing/type stripping
in JavaScript processes, JIT warm-up, array setup, worker creation, and shutdown.
There is no warmed, persistent worker pool and no steady-state JIT benchmark.
All printed results are checked for agreement between equivalent workloads.

Machine: AMD EPYC 9V74 virtual CPU, Linux x86-64; three logical CPUs exposed and a
cgroup quota of two CPUs (`200000 100000`). Go uses `GOMAXPROCS=2`; eight tasks
compete for the same CPU budget in every runtime. CPU sharing, allocation, and
startup noise affect these small timings. No Bun JIT-disabling flags or Node
optimization-disabling flags are used.

Raw samples, checked outputs, versions, and environment:
[results.json](results.json). The complete run transcript is
[raw-run.txt](raw-run.txt).

## Reproduce

Install Go 1.27+, Node 24+, Python 3, and Bun. If Bun is not already available,
you can install the measured version locally:

```sh
npm install --prefix benchmarks/.build/tools --no-audit --no-fund bun@1.4.2
```

From the repository root:

```sh
python3 benchmarks/run.py --samples 5 --cpus 2
```

The runner finds Bun on PATH or in that local installation. Generated sources,
executables, and local tooling stay in the ignored `.build/` directory.
Measurements are written to `benchmarks/results.json`; use `--output` to save
another run without replacing the published samples. On a machine with a
different CPU budget, set `--cpus` accordingly and report that difference.

The older [dynamic representation microbenchmarks](../examples/go-target/benchmarks/README.md)
isolate tagged containers versus Go interfaces; they are separate from these
complete-program comparisons.
