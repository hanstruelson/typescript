# Bun Node built-in port inventory

**Bun snapshot:** `e40a58ad` (`/workspace/bun`; clean worktree when inspected).
**Scope:** the Node module source and Bun runtime bindings at that snapshot. This document inventories source; it does not claim complete API support or behavioral parity.

## Port checkpoint

| Field | Value |
|---|---|
| Modules with registered runtime adapters | 57 |
| Bun source lines mapped in tracked checkpoints | 109117 |
| Active module | `node:net` |
| Active source file | `src/js/node/net.ts` |
| Next source line | 448 |

`node:sqlite` is paused at `NodeSqlite.cpp:1`: its remaining implementation enters JavaScriptCore and Bun's bundled SQLite C engine. The Go adapter retains the public facade and an in-memory SQL subset. `node:net.ts` is mapped through line 447, including option flags, BlockList/SocketAddress constructors, connection entry points, TLS key-log setup, socket-write error completion, and TLS bridge substitutions; source review pauses at line 448's adopted TLS raw-handle close path. `node_net_binding.rs` is mapped through line 184 onto the Go socket adapter. BlockList matching and rule descriptions are implemented, but its structured-clone hooks still cross the JavaScriptCore boundary.

Each source index row starts at `ported through line 0`. Update the row after each file is translated. The module-level counter advances only when a module is fully translated; keep the per-file line checkpoint even when work pauses mid-module.

Current stack checkpoint: `node:tls` uses Go’s standard `crypto/tls` without OpenSSL or cgo. `node:_http_server` is mapped through line 4306, including response completion, header accessors, writes, and queued pipelined operations. `node:http` and `node:https` are registered through their full wrappers; HTTP parser, agent, native binding, and platform details remain prototype approximations. `internal/dgram.ts` and `node:dgram.ts` are mapped through their full source ranges onto Go `net.UDPConn` and `syscall` adapters; cluster sharing, the native send queue, custom async lookup, and source-specific multicast membership remain approximations. `node:dns.ts` and its Promise facade are mapped through their full source ranges onto Go `net.Resolver`; uncommon record types remain limited by the standard library. `node:fs.ts` is mapped through all 1,509 lines, and `node:fs/promises.ts` through all 1,674 lines. The `node:readline` and `node:readline/promises` wrappers are mapped through their final lines. Their internal helpers cover the interface, history, key decoding, cursor callbacks, and promise cursor builder. `node:tty` is mapped through all 204 wrapper lines; stream controls and cursor/color APIs are exposed, while raw terminal mode and real window-size probing are approximated. `node:repl` and its internal Acorn, await, completion, mode, and utility dependencies are mapped through their final source lines and executed through the local VM; runtime behavior is unverified. `node:wasi` is mapped through all 1,775 lines and exposes the complete 47-entry preview1 import table; core memory, clock, random, descriptor, and path operations are adapted to Go, while socket support and several uncommon semantics are approximated. `node:v8` uses a local object-graph serialization format rather than V8 wire data, and `node:inspector` exposes Session, runtime console events, the inspector WebSocket endpoint, profiler/coverage commands, and its Promise facade; engine-backed debugger state, samples, and coverage remain approximations. `node:http2` maps its complete 6,794-line JavaScript wrapper and private TLS-upgrade helper to Go session, stream, request, response, server, settings, and validation surfaces; the Bun native frame parser is not translated, so protocol operation remains approximate. `node:quic` maps its public wrapper and internal symbol, diagnostics, header, binding, state, statistics, and class surfaces. Its connect/listen/session APIs are placeholders because Go’s standard library has no QUIC transport and no external QUIC dependency is enabled. `node:sqlite` now exposes DatabaseSync, StatementSync, Session, SQLTagStore, backup, constants, statement properties, and symbol disposal with a small in-memory SQL subset. `node:test` now maps the full 2,874-line wrapper and its mock-timer/assert utility dependencies; nested tests, hooks, contexts, mocks, assertions, and event-shaped run output are exposed, while file subprocess execution and snapshot persistence are stubs. `node:stream/iter` now maps its wrapper, ring buffer, push, duplex, pull, consumer, broadcast, share, and classic-stream adapters; backpressure and async scheduling remain approximate. `node:zlib/iter` exposes all 16 async/sync transform factories through the local transform adapter; gzip/deflate use Go standard codecs, while Brotli and Zstd remain unsupported. Bun’s native SQLite binding and bundled C engine are not translated; file-backed databases, full SQL semantics, changesets, serialization, and backup are placeholders or absent. The public module shim is mapped through 38/38 lines; the SQLite native binding is the active source at `NodeSqlite.cpp:1`, with its JavaScriptCore and bundled SQLite C engine boundary recorded in the module row. `node:fs/promises.ts` maps `FileHandle.readLines`, `pull`, `pullSync`, and `writer`; `internal/streams/iter/pull.ts` is mapped through all 977 lines and supplies the pull pipeline adapter. The fs dependencies `internal/fs/glob`, `cp-sync`, `cp`, `watch`, `watchfile`, `streams`, and `internal/streams/fast-utf8-stream` route through Go filesystem helpers. Glob uses the local simpler matcher in place of vendored minimatch; cp uses the shared Go copy walker instead of Node’s separate sync/async tree algorithms; watch helpers use the existing watcher state; filesystem streams use the local stream state machine. The UTF-8 stream maps buffering, flush, reopen, and lifecycle APIs, while Bun’s custom fs overrides, precise retry behavior, and zero-copy writer path are approximated. Through values remain per-file source checkpoints.

## Line-by-line dependency-stack strategy

1. Start at line 1 of a selected module entry point and translate its code in source order. Map a continuous source range to one or more named Go helpers rather than copying syntax mechanically.
2. If the next source range needs a Node module that is not available, stop before that range. Record the blocked file, next source line, and missing specifier in the active stack.
3. Push that dependency as the new active module and repeat from its first implementation line. When its APIs are available, mark its source files and pop back to the exact saved line in the caller.
4. Repeat until the original entry point reaches its end. Record each module’s dependencies and per-file `Through` line so no file silently jumps ahead.
5. Keep one Go module (`go.mod`) and the existing `goemit` package. Give each Node builtin its own file, named `internal/goemit/node_<builtin>_runtime.go`, containing that builtin's generated-runtime fragment. Go files in this package can share the private runtime value types and helpers without Go imports or package cycles. A build-tagged registry joins the fragments into the generated single-file runtime and dispatches module names to their constructors.

The completed `node:querystring`, `node:punycode`, `node:console`, `node:os`, `node:timers`, `node:timers/promises`, `node:util/types`, `node:module`, `node:perf_hooks`, `node:events`, `node:domain`, `node:diagnostics_channel`, `node:trace_events`, `node:assert`, `node:assert/strict`, `node:async_hooks`, `node:crypto`, `node:stream`, `node:stream/promises`, `node:stream/consumers`, and `node:stream/web` wrappers are registered in the generated runtime. The `node:module` adapter follows Bun's 76-entry builtin table and exposes Module construction, CommonJS wrapper helpers, lookup entry points, and compile-cache APIs; host resolver, source-map, and compile-cache behavior use local approximations. `node:async_hooks` implements `AsyncLocalStorage` and `AsyncResource` against persistent event-loop context frames; promise reactions and timer task completions capture those frames, while hook callbacks and async IDs remain stubs. The assertion facade provides the standard assertion method family, strict aliases, error matching, and promise assertions through local runtime helpers; deep object comparison and deprecated `CallTracker` reporting are approximations. The timers implementation shares the transpiler event loop; its promise interval uses the local async-iterator machinery. Timer handle refresh/ref APIs are not modeled yet. `node:util/types` maps predicates to generated runtime value tags; brands for weak collections, shared buffers, crypto keys, and native externals are unavailable and report false. `node:perf_hooks` uses Go's monotonic clock and exposes performance entries, observers, timerify, and a basic recordable histogram; Node-only engine observation and high resolution histogram quantization are approximated. Its `monitorEventLoopDelay` wrapper is translated to a local histogram adapter; the native C++ sampling substrate remains unported. `node:events` implements EventEmitter listeners, emission, static helpers, one-shot Promises, and async iteration; EventTarget bridges are adapted. `node:domain` builds on `node:events`; it tracks enter/run/exit and the common listener wrappers, but does not update a process global. `node:diagnostics_channel` provides cached channels, subscriptions, stores, and tracing wrappers; weak-reference cleanup and async-context propagation use strong references and direct calls. `node:trace_events` tracks category enablement and writes local trace files through an internal agent adapter; worker-specific native hooks and exact trace timing remain approximate. `node:os` uses Go's standard library and procfs in place of Bun's Rust host binding; OS constants and some non-Linux host details remain approximate. Its `src/runtime/node/node_os.rs` source rows remain at line 0 because the binding was replaced rather than translated line by line. Newly ported modules are enabled by default and can be omitted from a compiler build with `-tags no_node_modules`. The existing `fs`, `buffer`, and `path` runtime support predates this registry and is not controlled by that tag yet.

`node:process` source review is complete through `BunProcess.cpp:5172`; `ProcessBindingConstants.cpp/.h`, `NodeProcessModule.h`, and `node_process.rs` are also mapped through their final lines. No unported Node builtin dependency was encountered. The binding's TLS/OpenSSL constants are numeric compatibility values only; the Go runtime adds no OpenSSL or cgo dependency. The `node:path` wrappers, native bridge, Rust host wrapper, and portable Rust algorithm crate are mapped through their final lines; they route to the existing Go path runtime with standard-library host functions and local Windows path logic. `path.matchesGlob` uses the local glob matcher, so advanced Bun glob features remain approximate. The `node:buffer` module wrapper, Rust helper, JSBuffer bindings, encoding parser, constructor, and prototype sources have all been mapped. The Go runtime now exposes the constructor/statics, string and view operations, search and compare, typed accessors, numeric/BigInt readers and writers, encoding-specific methods, inspection, and module exports. Object URL lookup remains a placeholder; JSC memory mapping, resizable-buffer side effects, custom inspect context properties, and some exact Node error ordering/messages are approximated. The `node:string_decoder` native module handoff, declaration, and 631-line implementation are now translated into stateful UTF-8, UTF-16LE, and base64/base64url chunk decoding with the `write`, `end`, `text`, `lastChar`, `lastNeed`, `lastTotal`, and `encoding` API. Its encoders reuse the Buffer runtime; invalid receiver/error wording and a few malformed chunk boundary details are approximate. The `node:url` wrapper, internal IDNA helpers, NodeURL native bridge, and punycode/host-validation helpers are mapped through their final lines; the Go adapter uses `net/url`, `filepath`, and the existing punycode runtime, so ICU UTS #46 validation and detailed `URLPattern` matching remain approximate. Its URL, URLSearchParams, legacy `Url`, file URL, IDNA, and HTTP-options entry points are exposed. The crypto dependency stack has mapped `readable.ts` through line 1734, `writable.ts` through 1118, `duplex.ts` through 166, `duplexify.ts` through 395, `webstreams_adapters.ts` through 864, the TypeScript logic of `native-readable.ts` through 306, `iter/from.ts` through 576, `iter/classic.ts` through 898, `transform.ts` through 173, `lazy_transform.ts` through 53, `passthrough.ts` through 20, `pipeline.ts` through 448, and `compose.ts` through 221. These Go adapters cover the classic/Web stream conversion surfaces with approximate scheduling, destroy, and strategy behavior; Bun's `BunStreamConsumers.cpp` native-handle transfer remains a no-op.

The active stack has completed `internal/abort_listener.ts` (31 lines), `internal/async_hooks.ts` (28), `internal/errors.ts` (23), `internal/streams/end-of-stream.ts` (384), `internal/streams/add-abort-signal.ts` (58), `internal/streams/destroy.ts` (343), `internal/streams/state.ts` (52), `internal/streams/from.ts` (200), `internal/streams/iter/types.ts` (87), `internal/streams/iter/utils.ts` (270), `internal/streams/iter/from.ts` (576), and `internal/streams/iter/classic.ts` (898). Readable resumes after its classic iterator import on line 1286. Writable, Duplex, Duplexify, native-readable TypeScript, and Web stream adapters are mapped through their source ends; the C++ native-readable transfer is represented by a placeholder. Scheduling, Web stream strategies, native transfer, and stream error ordering remain approximate.

`node:crypto` now maps the 562-line JavaScript wrapper surface into a Go standard-library adapter. The C++ and Rust binding sources were replaced at the boundary, not translated line by line; Argon2, scrypt, classic Diffie-Hellman, SPKAC, FIPS/secure-heap controls, and subtle WebCrypto are explicit stubs or errors. The adapter imports no OpenSSL and uses no cgo. The `node:stream`, `stream/promises`, `stream/consumers`, and `stream/web` facades are registered. Readable operators and `duplexPair` are exposed with simplified concurrency and scheduling; `stream/web` forwards global constructors, and names missing from the runtime remain undefined. Experimental `stream/iter` and `zlib/iter` are registered with their wrappers and shared iterator helpers. Gzip/deflate transforms use Go standard codecs; Brotli/Zstd are exposed with unsupported-codec errors. Broadcast/Share and native async scheduling remain prototype adaptations.

## Port ledger

The stream dependency pass has mapped `internal/abort_listener.ts` through line 31, `internal/async_hooks.ts` through 28, `internal/errors.ts` through 23, `internal/streams/end-of-stream.ts` through 384, `internal/streams/add-abort-signal.ts` through 58, `internal/streams/destroy.ts` through 343, `internal/streams/state.ts` through 52, and `internal/streams/from.ts` through 200. `readable.ts` is at line 29 after its imports were resolved. The stream adapters expose EOS and abort handling, iterator-to-readable construction, shared state flags, and destroy lifecycle methods. The performance and error helpers use local runtime state, so Node native stream and host behavior are approximate.

`Through` tracks translated Bun source lines, not verified behavioral parity. The current runtime fragments are prototype implementations; the console fragment currently supplies common logging, count, and timing methods, not the full Node `Console` surface.

| Module | Bun source | Bun lines | Through | Status | Dependencies |
|---|---|---:|---:|---|---|
| `node:querystring` | `src/js/node/querystring.ts` | 523 | 523 | source translated | none |
| `node:punycode` | `src/js/node/punycode.ts` | 414 | 414 | source translated | none |
| `node:console` | `src/js/node/console.ts` | 1 | 1 | source translated; API partial | none |
| `node:os` | `src/js/node/os.ts` | 179 | 179 | wrapper translated; native bridge replaced | Go standard library |
| `node:timers` | `src/js/node/timers.ts` | 41 | 41 | wrapper translated; handle methods partial | `node:timers/promises` (lazy property) |
| `node:timers/promises` | `src/js/node/timers.promises.ts` | 253 | 253 | wrapper translated; async iterator adapted | `node:timers` (global timer primitives) |
| `node:util/types` | `src/jsc/modules/NodeUtilTypesModule.cpp`, `.h` | 1314 | 1314 | predicate table translated; host-only brands partial | generated runtime value tags |
| `node:perf_hooks` | `src/js/node/perf_hooks.ts`; `src/js/internal/perf_hooks/monitorEventLoopDelay.ts` | 452 | 452 | wrappers translated; native event-loop sampling approximated | Go monotonic clock, generated runtime |
| `node:events` | `src/js/node/events.ts` | 1057 | 1057 | EventEmitter and static APIs adapted; host integrations partial | `node:util/types`, local AbortSignal |
| `node:child_process` | `src/js/node/child_process.ts` | 2071 | 0 | broad spawn/exec/fork and synchronous API adapter uses `os/exec`; IPC, stream timing, and source-line sweep remain approximate | Go standard `os/exec`, local streams and events |
| `node:domain` | `src/js/node/domain.ts` | 96 | 96 | wrapper translated | `node:events` |
| `node:cluster` | `src/js/node/cluster.ts`; `src/js/internal/cluster/{primary,child,Worker,isPrimary,RoundRobinHandle,SharedHandle}.ts`; `src/js/internal/linkedlist.ts` | 1183 | 1183 | primary and worker APIs mapped; native process IPC and descriptor transfer are approximated | `node:events`, Go `net`, `os/exec`, local cluster handles |
| `node:net` | `src/js/node/net.ts`; `src/runtime/node/node_net_binding.rs`; `src/runtime/node/net/BlockList.rs` | 5342 | net.ts:447/4634; node_net_binding.rs:184/184; BlockList.rs:0/524 | Go `net` adapter exposes Socket, Server, BlockList, SocketAddress, IP predicates, connection controls, process flags, and TLS key-log file output; write failures reach stream destruction, while adopted TLS handle ownership, cluster listeners, and BlockList structured cloning remain approximate | Go standard `net`, `crypto/tls`, `os`, `syscall`; local Duplex, events, DNS, validators |
| `node:tls` | `src/js/node/tls.ts` | 1757 | 1757 | exported client/server/context and socket surfaces mapped through Go `crypto/tls`; native certificate caches, OpenSSL cipher strings, and several compatibility details are approximate | Go standard `crypto/tls`, `crypto/x509`, `net`; local `node:net` sockets |
| `node:_http_agent` | `src/js/node/_http_agent.ts` | 626 | 626 | agent pools, request queues, socket reuse, and proxy option handling adapted; CONNECT tunnel behavior and host networking details remain approximate | `node:events`, `internal/http`, validators, `node:tls`, `node:net`, async context |
| `node:http` | `src/js/node/http.ts` | 138 | 138 | public request/server facade, header controls, shared agent accessors, and proxy environment setup mapped | local `_http_*` adapters, `internal/http`, validators; lazy `node:https` edge |
| `node:https` | `src/js/node/https.ts` | 636 | 636 | HTTPS request/server facade, TLS agent session keys, proxy CONNECT handshake, and global agent mapped over Go TLS/net adapters | `node:http`, `node:tls`, `node:net`, `internal/http`, URL and validation helpers |
| `node:fs` | `src/js/node/fs.ts`; `src/js/node/fs.promises.ts` | 3183 | fs.ts:1509/1509; fs.promises.ts:1674/1674 | callback, sync, Promise, directory, watcher, copy, stream, glob, FileHandle, readLines, and disposable APIs map onto shared Go filesystem helpers; FileHandle pulls and writer APIs route through `internal/streams/iter/pull`; handle locking and write backpressure are approximated | Go filesystem operation, stream, watcher, directory, glob, copy, and handle adapters |
| `node:readline` | `src/js/node/readline.js` | 546 | 546 | legacy and promise-question interfaces, compatibility accessors, terminal callbacks, keypress events, and async line iteration adapted | `internal/readline/*`, history, events, validators, string decoder |
| `node:readline/promises` | `src/js/node/readline.promises.js`; `src/js/internal/readline/promises.js` | 198 | 60/60; 138/138 | Promise question interface and queued cursor-control `Readline` API mapped; terminal writes use local stream/event-loop hooks | `node:events`, internal readline interface and cursor utilities |
| `node:tty` | `src/js/node/tty.ts` | 204 | 204 | ReadStream/WriteStream wrappers, ref methods, cursor controls, color depth, and TTY predicates exposed; raw mode and terminal sizing are simplified | `internal/fs/streams`, `internal/tty`, validators, `node:readline` |
| `node:repl` | `src/js/node/repl.js` | 1543 | 1543 | Bun’s REPL entry point and internal dependency stack embedded and routed through the local VM; runtime behavior remains unverified | `node:vm`, `node:readline`, `node:console`, `node:async_hooks`, `node:path`, `node:fs`, internal REPL helpers |
| `node:wasi` | `src/js/node/wasi.ts` | 1775 | 1775 | `WASI` constructor and all 47 preview1 imports exposed; common memory, clock, random, descriptor, and path calls mapped, with sockets and several semantics approximated | Go `os`, `filepath`, `time`, `crypto/rand`; `node:tty`, `node:fs`, `node:path` |
| `node:v8` | `src/js/node/v8.ts` | 450 | 450 | heap counters, GCProfiler shape, snapshot entry points, and all exports exposed; serializer preserves common JS graphs in a local format that is not V8 wire compatible | Go runtime/JSON, `node:stream`, `node:os`, `node:worker_threads`, `internal/serialization_buffers` |
| `node:inspector` | `src/js/node/inspector.ts` | 760 | 760 | Session and CDP command surface, console notifications, inspector WebSocket listener, lifecycle, and profiler/coverage entry points exposed; engine-backed debugging and profile data are approximated | Go `net`, WebSocket framing, `node:events`, `node:worker_threads`, `internal/trace_events` |
| `node:inspector/promises` | `src/js/node/inspector.promises.ts` | 28 | 28 | Promise-based Session facade exposed | `node:inspector` |
| `node:http2` | `src/js/node/http2.ts`; `src/js/node/_http2_upgrade.ts` | 7199 | 6794/6794; 405/405 | constants, settings, headers, sessions, streams, requests, responses, servers, and TLS upgrade surfaces exposed; Bun's native HTTP/2 frame engine is not translated and protocol behavior is an API-shaped prototype | Go streams, `crypto/tls`, sockets, `node:events`, private upgrade helper |
| `node:quic` | `src/js/node/quic.ts`; `src/js/internal/quic/{binding,diagnostics,http2util,quic,state,stats,symbols}.ts` | 7828 | 55/55; 2/2; 70/70; 197/197; 5565/5565; 886/886; 983/983; 70/70 | public classes, state/stat objects, callback names, diagnostics channels, settings/header helpers, and internal binding constants mapped; wire protocol and native transport are unavailable | Go `crypto/tls`, local streams, diagnostics channels; transport placeholder |
| `node:sqlite` | `src/jsc/modules/NodeSqliteModule.h`; `src/jsc/bindings/sqlite/NodeSqlite.{h,cpp}`, `JSSQLStatement.{h,cpp}`, SQLite engine | 293,545 | NodeSqliteModule.h:38/38; native binding:0 | DatabaseSync/StatementSync/Session/SQLTagStore and module exports are exposed over an in-memory SQL subset; the native SQLite engine is unavailable in the Go standard library | local SQL facade; no third-party Go driver |
| `node:test` | `src/js/node/test.ts`; `src/js/internal/test_runner/mock_timers.ts`; `src/js/internal/assert/utils.ts` | 3,690 | 2874/2874; 781/781; 35/35 | test/suite/hook/context/mock/assert and run stream surfaces mapped; file subprocess runner and snapshot persistence stubbed | local `node:stream`, `node:assert`, `node:async_hooks`, timer scheduler |
| `node:dgram` | `src/js/node/dgram.ts` | 1546 | 1546 | broad UDP Socket API adapter mapped onto Go `net.UDPConn` and `syscall`; broadcast, TTL, loopback, interface, and common membership options are wired, while cluster sharing, native queue behavior, custom async lookup, and source-specific membership remain approximations | `internal/dgram` (146/146), validators, `net.BlockList`, events, DNS, diagnostics channel |
| `node:dns` | `src/js/node/dns.ts`; `src/js/node/dns.promises.ts` | 1073 | 1071/1071; 2/2 | callback and Promise resolver APIs, per-instance resolver state, server configuration, callback promisify hooks, and result ordering mapped over Go’s standard resolver; unsupported record types return Node-shaped errors | Go standard `net.Resolver`, local validators and IP helpers |
| `node:_http_incoming` | `src/js/node/_http_incoming.ts` | 737 | 737 | IncomingMessage constructor, lazy headers/trailers, read and destroy lifecycle, and header folding mapped onto local streams | `internal/streams/readable`, EOS, `internal/http`; local stream and native handle adapters |
| `node:_http_common` | `src/js/node/_http_common.ts` | 361 | 361 | validation, parser pool, callback lifecycle, leniency flags, and parser interface exposed; native llhttp behavior and ConnectionsList are approximated | `internal/validators`, `internal/freelist`, `node:_http_incoming`; Go byte/string parser adapter |
| `node:_http_outgoing` | `src/js/node/_http_outgoing.ts` | 1335 | 1335 | OutgoingMessage construction, header APIs, serialization, write buffering, trailers, finish/flush, and socket accessors mapped; deprecation warnings and exact error behavior are approximate | `internal/stream`, validators, HTTP shared symbols, `_http_common`, `internal/http/FakeSocket`, events |
| `node:diagnostics_channel` | `src/js/node/diagnostics_channel.ts` | 419 | 419 | channel and tracing APIs adapted; weak cleanup partial | generated runtime values |
| `node:trace_events` | `src/js/node/trace_events.ts`; `src/js/internal/trace_events.ts` | 961 | 87/87; agent API mapped | category buffers, trace writing, and common instrumentation adapted; worker/native hooks partial | Go file and crypto/zlib instrumentation, generated runtime |
| `node:assert` | `src/js/node/assert.ts` | 981 | 981 | public assertion APIs adapted; deep comparison and CallTracker partial | `node:util/types`, generated errors and promises |
| `node:assert/strict` | `src/js/node/assert.strict.ts` | 2 | 2 | strict facade translated | `node:assert` |
| `node:async_hooks` | `src/js/node/async_hooks.ts` | 634 | 634 | AsyncLocalStorage and AsyncResource adapted; native hook callbacks partial | event-loop async-context frames |
| `node:constants` | `src/jsc/modules/NodeConstantsModule.h` | 981 | 981 | source translated; host POSIX and crypto tables emitted; platform-conditional symbols vary by host | Go standard `syscall` constants |
| `node:process` | `src/jsc/modules/NodeProcessModule.h`; `src/runtime/node/node_process.rs`; `src/jsc/bindings/BunProcess.cpp`; `src/jsc/bindings/ProcessBindingConstants.cpp`; `.h` | 7092 | 24/24; 595/595; BunProcess.cpp:5172/5172; binding constants:1301/1301 | source sweep complete; process APIs adapted, some platform and engine host APIs remain approximations or stubs | Go standard `os`, `runtime`, `syscall`, `os/user`; local events runtime |
| `node:buffer` | `src/jsc/modules/NodeBufferModule.cpp`; `.h`; `src/runtime/node/buffer.rs`; `src/jsc/bindings/JSBuffer.cpp`; `.h`; `JSBufferEncodingType.cpp`; `.h`; `BufferEncodingType.h`; `src/js/builtins/JSBufferConstructor.ts`; `JSBufferPrototype.ts` | 5011 | NodeBufferModule.cpp:282/282; .h:210/210; buffer.rs:141/141; JSBuffer.cpp:4016/4016; JSBuffer.h:100/100; encoding sources:196/196; TS sources:66/66 | source sweep complete; API mapped with noted runtime approximations | Go typed arrays, standard `bytes`, `encoding/*`, `math/big`, and local buffer methods |
| `node:module` | `src/jsc/modules/NodeModuleModule.cpp/.h`; `src/jsc/NodeModuleModule.rs`; `src/jsc/NodeModuleModule.bind.ts` | 1549 | 1257/1257; 45/45; 239/239; 8/8 | source sweep complete; Node's builtin table and CommonJS Module facade adapted; resolver and host callbacks approximated | Go standard `os`, `path/filepath`, and generated addon registry |
| `node:path` | `src/js/node/path.ts`; `src/jsc/bindings/Path.cpp`; `src/runtime/node/path.rs`; `src/node_path/lib.rs` | 3919 | 81/81; 168/168; 956/956; 2714/2714 | source sweep complete; all exported POSIX/Windows operations mapped, host cwd and glob behavior adapted | Go standard `os`, `path`, `filepath`, `strings`, `unicode/utf16`; local path validation and glob matcher |
| `node:path/posix` | `src/js/node/path.posix.ts` | 2 | 2 | source translated; routes to `node:path`.posix | `node:path` |
| `node:path/win32` | `src/js/node/path.win32.ts` | 2 | 2 | source translated; routes to `node:path`.win32 | `node:path` |
| `node:crypto` | `src/js/node/crypto.ts` | 562 | 562 | API surface mapped onto Go standard crypto; unsupported native facilities remain explicit stubs | Go standard `crypto/*`; no OpenSSL/cgo |
| `node:stream` | `src/js/node/stream.ts`; `src/js/internal/stream.ts` | 120 | 120 | classic stream facade, constructors, predicates, pipeline, operators, and pairing adapted; concurrency and scheduling are approximate | mapped internal stream runtime |
| `node:stream/promises` | `src/js/node/stream.promises.ts`; `src/js/internal/stream.promises.ts` | 47 | 47 | promise `pipeline` and `finished` entry points exposed | `node:stream`, EOS, pipeline |
| `node:stream/consumers` | `src/js/node/stream.consumers.ts` | 54 | 54 | six consume methods exposed; Blob construction and decoding use local runtime values | async iterator, Buffer, Blob |
| `node:stream/web` | `src/js/node/stream.web.ts` | 20 | 20 | Web Stream constructor names forwarded from the global object | generated Web Stream values |
| `node:zlib/iter` | `src/js/node/zlib.iter.ts`; `src/js/internal/streams/iter/transform.ts` | 811 | 51/51; 760/760 | 16 async/sync gzip, deflate, Brotli, and Zstd transform factories exposed; Go standard library handles gzip/deflate, while Brotli/Zstd report unsupported codec errors and native streaming behavior is approximated | Go standard `compress/gzip`, `compress/zlib`, `compress/flate`; `node:buffer`, `node:util/types` |
| `node:stream/iter` | `src/js/node/stream.iter.ts` and `src/js/internal/streams/iter/{types,utils,from,pull,consumers,classic,ringbuffer,push,duplex,broadcast,share}.ts` | 6,021 | entry 165/165; types 87/87; utils 270/270; from 576/576; pull 977/977; consumers 495/495; classic 898/898; ringbuffer 147/147; push 747/747; duplex 142/142; broadcast 785/785; share 732/732 | exports, push/duplex/consumer/classic adapters, and broadcast/share cursor buffers mapped; backpressure and scheduling remain approximate | local iterator, timer, stream, validator, and Buffer runtimes |

## What counts as a module source

- `src/resolve_builtins/HardcodedModule.rs` is the runtime resolver list. It contains 72 `node:` specifiers, including private `_http*`, `_tls*`, and `_stream*` routes.
- 65 of those specifiers have a JavaScript entry file in `src/js/node`. Seven are native or synthetic module registrations. `_http2_upgrade.ts` is an additional private helper loaded by `node:http2`, not a separate `node:` resolver entry.
- Bun’s module registry also contains literal `require()` dependencies under `src/js/internal`. The edge index below records their transitive graph without duplicating a shared dependency list for every public module.
- Native bridges referenced from reachable JavaScript and Node integration sources are listed separately. The 45 files under `src/runtime/node` and the `src/node_path` Rust crate are indexed as runtime implementations.
- `src/node-fallbacks` is excluded: its README describes browser-target polyfills, not Bun’s runtime `node:` modules. Bun’s module README also says builtin sources are preprocessed into an internal registry, with literal `require()` calls rewritten to module IDs; these are not ordinary Go package imports.
- Generic JavaScriptCore, operating-system, TLS, networking, and other engine dependencies are not recursively listed as module source. They are platform/runtime substrate. Direct `$cpp` and `$rust` handoff targets are listed so the Go port has the boundary visible.

## Runtime `node:` specifier to source route

| Runtime specifier | Bun source route |
|---|---|
| `node:assert` | `src/js/node/assert.ts` |
| `node:assert/strict` | `src/js/node/assert.strict.ts` |
| `node:async_hooks` | `src/js/node/async_hooks.ts` |
| `node:buffer` | `src/jsc/modules/NodeBufferModule.cpp` and `.h` |
| `node:child_process` | `src/js/node/child_process.ts` |
| `node:console` | `src/js/node/console.ts` |
| `node:constants` | `src/jsc/modules/NodeConstantsModule.h` |
| `node:crypto` | `src/js/node/crypto.ts` |
| `node:dns` | `src/js/node/dns.ts` |
| `node:dns/promises` | `src/js/node/dns.promises.ts` |
| `node:domain` | `src/js/node/domain.ts` |
| `node:events` | `src/js/node/events.ts` |
| `node:fs` | `src/js/node/fs.ts` |
| `node:fs/promises` | `src/js/node/fs.promises.ts` |
| `node:http` | `src/js/node/http.ts` |
| `node:https` | `src/js/node/https.ts` |
| `node:module` | `src/jsc/modules/NodeModuleModule.cpp` and `.h`; `src/jsc/NodeModuleModule.rs`; `src/jsc/NodeModuleModule.bind.ts` |
| `node:net` | `src/js/node/net.ts` |
| `node:os` | `src/js/node/os.ts` |
| `node:path` | `src/js/node/path.ts` |
| `node:path/posix` | `src/js/node/path.posix.ts` |
| `node:path/win32` | `src/js/node/path.win32.ts` |
| `node:perf_hooks` | `src/js/node/perf_hooks.ts` |
| `node:process` | `src/jsc/modules/NodeProcessModule.h` and `src/runtime/node/node_process.rs` |
| `node:querystring` | `src/js/node/querystring.ts` |
| `node:readline` | `src/js/node/readline.js` |
| `node:readline/promises` | `src/js/node/readline.promises.js` |
| `node:sqlite` | `src/jsc/modules/NodeSqliteModule.h`, `src/jsc/bindings/sqlite/NodeSqlite.{h,cpp}`, `JSSQLStatement.{h,cpp}`, and local SQLite engine sources |
| `node:stream` | `src/js/node/stream.ts` |
| `node:stream/consumers` | `src/js/node/stream.consumers.ts` |
| `node:stream/iter` | `src/js/node/stream.iter.ts` |
| `node:zlib/iter` | `src/js/node/zlib.iter.ts` |
| `node:stream/promises` | `src/js/node/stream.promises.ts` |
| `node:stream/web` | `src/js/node/stream.web.ts` |
| `node:string_decoder` | `src/jsc/modules/NodeStringDecoderModule.h` (native module) |
| `node:test` | `src/js/node/test.ts` |
| `node:timers` | `src/js/node/timers.ts` |
| `node:timers/promises` | `src/js/node/timers.promises.ts` |
| `node:tls` | `src/js/node/tls.ts` |
| `node:tty` | `src/js/node/tty.ts` |
| `node:url` | `src/js/node/url.ts` |
| `node:util` | `src/js/node/util.ts` |
| `node:util/types` | `src/jsc/modules/NodeUtilTypesModule.cpp` and `.h` |
| `node:vm` | `src/js/node/vm.ts` |
| `node:wasi` | `src/js/node/wasi.ts` |
| `node:zlib` | `src/js/node/zlib.ts` |
| `node:worker_threads` | `src/js/node/worker_threads.ts` |
| `node:punycode` | `src/js/node/punycode.ts` |
| `node:v8` | `src/js/node/v8.ts` |
| `node:trace_events` | `src/js/node/trace_events.ts` |
| `node:repl` | `src/js/node/repl.js` |
| `node:inspector` | `src/js/node/inspector.ts` |
| `node:inspector/promises` | `src/js/node/inspector.promises.ts` |
| `node:http2` | `src/js/node/http2.ts` |
| `node:quic` | `src/js/node/quic.ts` |
| `node:diagnostics_channel` | `src/js/node/diagnostics_channel.ts` |
| `node:dgram` | `src/js/node/dgram.ts` |
| `node:cluster` | `src/js/node/cluster.ts` |
| `node:_stream_duplex` | `src/js/node/_stream_duplex.ts` |
| `node:_stream_passthrough` | `src/js/node/_stream_passthrough.ts` |
| `node:_stream_readable` | `src/js/node/_stream_readable.ts` |
| `node:_stream_transform` | `src/js/node/_stream_transform.ts` |
| `node:_stream_wrap` | `src/js/node/_stream_wrap.ts` |
| `node:_stream_writable` | `src/js/node/_stream_writable.ts` |
| `node:_tls_common` | `src/js/node/_tls_common.ts` |
| `node:_tls_wrap` | `src/js/node/_tls_wrap.ts` |
| `node:_http_agent` | `src/js/node/_http_agent.ts` |
| `node:_http_client` | `src/js/node/_http_client.ts` |
| `node:_http_common` | `src/js/node/_http_common.ts` |
| `node:_http_incoming` | `src/js/node/_http_incoming.ts` |
| `node:_http_outgoing` | `src/js/node/_http_outgoing.ts` |
| `node:_http_server` | `src/js/node/_http_server.ts` |

Bare names alias their corresponding `node:` entry except `sqlite`, `test`, and `quic`, which resolve only with the `node:` prefix. The private HTTP, TLS, and stream entries are also wired as resolver aliases. `sys` and `node:sys` alias `node:util`. This list is from Bun’s resolver, not the browser fallback table.

## Local JavaScript dependency graph

Each row lists the literal `require()` specifiers in one reachable source file. `internal/...` resolves under `src/js/internal`; `node:...` resolves through the route table above; entries with no local wrapper are synthetic/native imports. Following these edges from any `src/js/node/...` route yields that module’s local transitive source closure. The graph includes all 176 reachable JavaScript sources (66 `src/js/node` files and 110 `src/js/internal` files).

| Bun source file | Literal `require()` dependencies |
|---|---|
| `src/js/internal/abort_listener.ts` | `internal/validators`, `internal/shared` |
| `src/js/internal/assert/assertion_error.ts` | `internal/primordials`, `internal/util/inspect`, `internal/util/colors`, `internal/validators`, `internal/assert/myers_diff` |
| `src/js/internal/assert/calltracker.ts` | `internal/primordials`, `internal/assert/assertion_error`, `internal/validators` |
| `src/js/internal/assert/myers_diff.ts` | `internal/util/colors` |
| `src/js/internal/assert/utils.ts` | `internal/assert/assertion_error` |
| `src/js/internal/async_context_frame.ts` | — |
| `src/js/internal/async_hooks.ts` | — |
| `src/js/internal/async_hooks_tick.ts` | — |
| `src/js/internal/cluster/RoundRobinHandle.ts` | `internal/linkedlist`, `internal/shared`, `node:net` |
| `src/js/internal/cluster/SharedHandle.ts` | `node:fs` |
| `src/js/internal/cluster/Worker.ts` | `node:events` |
| `src/js/internal/cluster/child.ts` | `node:events`, `internal/cluster/Worker`, `node:path`, `internal/shared` |
| `src/js/internal/cluster/isPrimary.ts` | — |
| `src/js/internal/cluster/primary.ts` | `node:events`, `internal/cluster/Worker`, `internal/shared`, `node:child_process`, `internal/cluster/RoundRobinHandle`, `internal/cluster/SharedHandle`, `node:path` |
| `src/js/internal/dgram.ts` | — |
| `src/js/internal/errors.ts` | `internal/primordials` |
| `src/js/internal/fixed_queue.ts` | — |
| `src/js/internal/freelist.ts` | — |
| `src/js/internal/fs/binding.ts` | — |
| `src/js/internal/fs/cp-sync.ts` | `internal/validators`, `node:fs`, `node:path` |
| `src/js/internal/fs/cp.ts` | `internal/fs/cp-sync`, `node:fs/promises`, `node:path` |
| `src/js/internal/fs/glob.ts` | `internal/validators`, `node:path`, `internal/shared`, `node:fs`, `node:fs/promises` |
| `src/js/internal/fs/streams.ts` | `node:stream`, `node:fs`, `internal/validators`, `node:path`, `internal/streams/destroy` |
| `src/js/internal/fs/watch.ts` | `node:events`, `node:path`, `internal/fs/binding`, `internal/shared` |
| `src/js/internal/fs/watchfile.ts` | `node:events`, `internal/async_context_frame`, `internal/validators`, `internal/fs/binding` |
| `src/js/internal/http.ts` | `internal/net/isIP` |
| `src/js/internal/http/FakeSocket.ts` | `internal/http`, `internal/shared`, `internal/stream` |
| `src/js/internal/http1_server_fallback.ts` | `internal/http`, `internal/primordials`, `internal/async_context_frame`, `node:_http_common`, `node:http`, `internal/shared` |
| `src/js/internal/linkedlist.ts` | — |
| `src/js/internal/net/isIP.ts` | — |
| `src/js/internal/net/server_abort_signal.ts` | `internal/validators`, `internal/abort_listener` |
| `src/js/internal/net/symbols.ts` | — |
| `src/js/internal/perf_hooks/monitorEventLoopDelay.ts` | `internal/validators` |
| `src/js/internal/primordials.js` | — |
| `src/js/internal/promisify.ts` | `internal/validators`, `internal/shared` |
| `src/js/internal/quic/binding.ts` | — |
| `src/js/internal/quic/diagnostics.ts` | `node:diagnostics_channel` |
| `src/js/internal/quic/http2util.ts` | — |
| `src/js/internal/quic/quic.ts` | `internal/primordials`, `node:util`, `node:buffer`, `node:util/types`, `node:net`, `node:crypto`, `internal/quic/binding`, `node:fs/promises`, `internal/streams/iter/types`, `internal/streams/iter/utils`, `internal/streams/iter/from`, `internal/validators`, `internal/quic/http2util`, `internal/quic/symbols`, `internal/quic/stats`, `internal/quic/state`, `internal/shared`, `internal/quic/diagnostics` |
| `src/js/internal/quic/state.ts` | `internal/primordials`, `node:util/types`, `node:util`, `node:os`, `internal/quic/symbols`, `internal/quic/binding` |
| `src/js/internal/quic/stats.ts` | `internal/primordials`, `node:util/types`, `node:util`, `internal/quic/symbols`, `internal/quic/binding` |
| `src/js/internal/quic/symbols.ts` | — |
| `src/js/internal/readline/callbacks.js` | `internal/repl/node-primordials`, `internal/repl/node-errors`, `internal/validators`, `internal/readline/utils` |
| `src/js/internal/readline/emitKeypressEvents.js` | `internal/repl/node-primordials`, `internal/readline/utils`, `internal/readline/interface`, `node:timers`, `node:string_decoder` |
| `src/js/internal/readline/interface.js` | `internal/repl/node-primordials`, `internal/repl/node-errors`, `internal/validators`, `internal/shared`, `internal/repl/node-inspect`, `node:events`, `internal/abort_listener`, `internal/readline/utils`, `internal/readline/callbacks`, `node:string_decoder`, `internal/readline/emitKeypressEvents`, `internal/repl/history` |
| `src/js/internal/readline/promises.js` | `internal/repl/node-primordials`, `internal/readline/utils`, `internal/validators`, `internal/repl/node-errors` |
| `src/js/internal/readline/utils.js` | `internal/repl/node-primordials` |
| `src/js/internal/repl.js` | `internal/repl/node-primordials`, `node:repl`, `internal/repl/utils` |
| `src/js/internal/repl/acorn-walk.js` | `node:vm` |
| `src/js/internal/repl/acorn.js` | `node:vm` |
| `src/js/internal/repl/await.js` | `internal/repl/acorn`, `internal/repl/acorn-walk`, `node:repl` |
| `src/js/internal/repl/completion.js` | `internal/repl/node-primordials`, `internal/repl/utils`, `internal/repl/node-shims`, `node:path`, `node:fs`, `internal/repl/acorn`, `internal/repl/acorn-walk` |
| `src/js/internal/repl/history.js` | `internal/repl/node-primordials`, `internal/validators`, `node:fs`, `node:os`, `node:path`, `internal/repl/node-shims`, `node:timers`, `internal/readline/utils` |
| `src/js/internal/repl/mode.js` | — |
| `src/js/internal/repl/node-errors.js` | — |
| `src/js/internal/repl/node-inspect.js` | `internal/util/inspect` |
| `src/js/internal/repl/node-primordials.js` | `internal/primordials` |
| `src/js/internal/repl/node-shims.js` | `node:util`, `node:module`, `node:path`, `internal/repl/node-primordials`, `internal/util/colors`, `node:vm` |
| `src/js/internal/repl/utils.js` | `internal/repl/node-primordials`, `internal/repl/acorn`, `internal/repl/node-shims`, `internal/repl/node-errors`, `internal/readline/callbacks`, `internal/readline/interface`, `internal/readline/utils`, `internal/repl/node-inspect`, `node:vm`, `internal/repl/mode`, `node:path` |
| `src/js/internal/serialization_buffers.ts` | `node:buffer`, `node:util/types` |
| `src/js/internal/shared.ts` | `node:util`, `internal/primordials`, `internal/async_context_frame` |
| `src/js/internal/stream.promises.ts` | `internal/streams/utils`, `internal/streams/pipeline`, `internal/streams/end-of-stream` |
| `src/js/internal/stream.ts` | `internal/streams/operators`, `internal/streams/compose`, `internal/streams/state`, `internal/streams/pipeline`, `internal/streams/destroy`, `internal/streams/end-of-stream`, `internal/streams/utils`, `node:util/types`, `internal/streams/legacy`, `internal/streams/readable`, `internal/streams/writable`, `internal/streams/duplex`, `internal/streams/transform`, `internal/streams/passthrough`, `internal/streams/duplexpair`, `internal/streams/add-abort-signal`, `internal/stream.promises` |
| `src/js/internal/streams/add-abort-signal.ts` | `internal/streams/utils`, `internal/streams/end-of-stream`, `internal/abort_listener` |
| `src/js/internal/streams/compose.ts` | `internal/streams/pipeline`, `internal/streams/duplex`, `internal/streams/destroy`, `internal/streams/utils`, `internal/streams/end-of-stream` |
| `src/js/internal/streams/destroy.ts` | `internal/errors`, `internal/streams/utils` |
| `src/js/internal/streams/duplex.ts` | `internal/streams/legacy`, `internal/streams/readable`, `internal/streams/writable`, `internal/streams/add-abort-signal`, `internal/streams/destroy`, `internal/streams/utils`, `internal/webstreams_adapters`, `internal/streams/duplexify` |
| `src/js/internal/streams/duplexify.ts` | `internal/streams/utils`, `internal/streams/end-of-stream`, `internal/streams/destroy`, `internal/streams/duplex`, `internal/streams/readable`, `internal/streams/writable`, `internal/streams/from` |
| `src/js/internal/streams/duplexpair.ts` | `internal/streams/duplex` |
| `src/js/internal/streams/end-of-stream.ts` | `internal/shared`, `internal/validators`, `internal/streams/utils`, `node:async_hooks`, `internal/async_hooks`, `internal/abort_listener` |
| `src/js/internal/streams/fast-utf8-stream.ts` | `internal/validators`, `node:events`, `node:path`, `node:fs` |
| `src/js/internal/streams/from.ts` | `internal/errors` |
| `src/js/internal/streams/iter/broadcast.ts` | `internal/validators`, `internal/streams/iter/types`, `internal/streams/iter/from`, `internal/streams/iter/pull`, `internal/streams/iter/utils`, `internal/streams/iter/ringbuffer` |
| `src/js/internal/streams/iter/classic.ts` | `internal/errors`, `internal/validators`, `internal/streams/end-of-stream`, `internal/streams/add-abort-signal`, `internal/streams/iter/types`, `internal/streams/iter/utils`, `internal/streams/destroy`, `internal/streams/readable`, `internal/streams/writable`, `internal/streams/iter/from`, `node:util/types` |
| `src/js/internal/streams/iter/consumers.ts` | `internal/validators`, `internal/streams/iter/from`, `internal/streams/iter/utils`, `internal/streams/iter/types`, `node:util/types` |
| `src/js/internal/streams/iter/duplex.ts` | `internal/streams/iter/push`, `internal/validators` |
| `src/js/internal/streams/iter/from.ts` | `node:util/types`, `internal/streams/iter/types`, `internal/streams/iter/utils` |
| `src/js/internal/streams/iter/pull.ts` | `internal/validators`, `node:util/types`, `internal/streams/iter/from`, `internal/streams/iter/utils`, `internal/streams/iter/types` |
| `src/js/internal/streams/iter/push.ts` | `node:util/types`, `internal/validators`, `internal/streams/iter/types`, `internal/streams/iter/utils`, `internal/streams/iter/pull`, `internal/streams/iter/ringbuffer` |
| `src/js/internal/streams/iter/ringbuffer.ts` | — |
| `src/js/internal/streams/iter/share.ts` | `internal/streams/iter/types`, `internal/streams/iter/from`, `internal/streams/iter/pull`, `internal/streams/iter/utils`, `internal/streams/iter/ringbuffer`, `internal/validators` |
| `src/js/internal/streams/iter/transform.ts` | `node:buffer`, `node:util/types`, `internal/streams/iter/types`, `internal/validators` |
| `src/js/internal/streams/iter/types.ts` | — |
| `src/js/internal/streams/iter/utils.ts` | `node:util/types`, `internal/validators`, `node:util` |
| `src/js/internal/streams/lazy_transform.ts` | `internal/streams/transform` |
| `src/js/internal/streams/legacy.ts` | `node:events`, `node:util/types` |
| `src/js/internal/streams/native-readable.ts` | `internal/streams/readable`, `internal/streams/destroy` |
| `src/js/internal/streams/operators.ts` | `internal/validators`, `internal/shared`, `internal/streams/end-of-stream` |
| `src/js/internal/streams/passthrough.ts` | `internal/streams/transform` |
| `src/js/internal/streams/pipeline.ts` | `internal/streams/end-of-stream`, `internal/shared`, `internal/streams/destroy`, `internal/streams/duplex`, `internal/errors`, `internal/validators`, `internal/streams/utils`, `internal/streams/readable`, `internal/abort_listener`, `internal/streams/passthrough` |
| `src/js/internal/streams/readable.ts` | `node:events`, `internal/streams/legacy`, `internal/streams/add-abort-signal`, `internal/streams/end-of-stream`, `internal/streams/destroy`, `internal/streams/state`, `internal/streams/utils`, `internal/errors`, `internal/validators`, `node:string_decoder`, `internal/streams/from`, `internal/primordials`, `internal/shared`, `internal/streams/iter/classic`, `internal/streams/iter/types`, `internal/streams/compose`, `internal/webstreams_adapters` |
| `src/js/internal/streams/state.ts` | `internal/validators` |
| `src/js/internal/streams/transform.ts` | `internal/streams/duplex`, `internal/streams/state` |
| `src/js/internal/streams/utils.ts` | — |
| `src/js/internal/streams/writable.ts` | `node:events`, `internal/streams/legacy`, `internal/streams/destroy`, `internal/streams/end-of-stream`, `internal/streams/add-abort-signal`, `internal/streams/state`, `internal/streams/utils`, `internal/webstreams_adapters` |
| `src/js/internal/test_runner/mock_timers.ts` | `internal/validators`, `internal/abort_listener`, `node:timers`, `node:timers/promises`, `node:events` |
| `src/js/internal/timers.ts` | `internal/validators` |
| `src/js/internal/tls.ts` | `node:util/types`, `internal/validators`, `node:crypto`, `internal/shared` |
| `src/js/internal/trace_events.ts` | `internal/fs/binding`, `node:worker_threads`, `node:fs`, `node:path`, `node:zlib`, `node:crypto` |
| `src/js/internal/tty.ts` | `node:os` |
| `src/js/internal/url.ts` | — |
| `src/js/internal/util/colors.ts` | `internal/tty` |
| `src/js/internal/util/comparisons.ts` | — |
| `src/js/internal/util/deprecate.ts` | `internal/validators` |
| `src/js/internal/util/inspect.js` | `internal/primordials`, `node:util/types`, `internal/validators`, `node:module`, `node:url`, `node:buffer` |
| `src/js/internal/util/mime.ts` | — |
| `src/js/internal/validators.ts` | `internal/shared`, `node:path` |
| `src/js/internal/webstreams_adapters.ts` | `internal/primordials`, `internal/streams/writable`, `internal/streams/readable`, `internal/streams/duplex`, `internal/streams/destroy`, `internal/streams/utils`, `internal/shared`, `internal/validators`, `node:util/types`, `internal/streams/end-of-stream`, `internal/streams/native-readable` |
| `src/js/internal/worker/messaging.ts` | `internal/validators`, `internal/primordials` |
| `src/js/internal/worker/stdio.ts` | `internal/streams/readable`, `internal/streams/writable` |
| `src/js/node/_http2_upgrade.ts` | `node:stream`, `internal/tls` |
| `src/js/node/_http_agent.ts` | `node:events`, `internal/http`, `internal/shared`, `internal/validators`, `internal/net/isIP`, `internal/net/symbols`, `internal/async_context_frame`, `node:tls`, `node:net`, `node:async_hooks` |
| `src/js/node/_http_client.ts` | `internal/shared`, `node:_http_common`, `node:_http_outgoing`, `node:_http_agent`, `internal/url`, `internal/http`, `internal/validators`, `internal/timers`, `internal/streams/add-abort-signal`, `internal/streams/end-of-stream`, `node:diagnostics_channel`, `internal/async_context_frame`, `node:tls`, `node:net`, `internal/trace_events` |
| `src/js/node/_http_common.ts` | `internal/validators`, `internal/freelist`, `node:_http_incoming` |
| `src/js/node/_http_incoming.ts` | `internal/streams/readable`, `internal/streams/end-of-stream`, `internal/http` |
| `src/js/node/_http_outgoing.ts` | `internal/stream`, `internal/validators`, `internal/util/deprecate`, `internal/streams/state`, `internal/http`, `node:_http_common`, `internal/http/FakeSocket`, `node:events` |
| `src/js/node/_http_server.ts` | `node:events`, `node:stream`, `internal/net/isIP`, `node:_http_common`, `internal/validators`, `internal/shared`, `internal/cluster/isPrimary`, `internal/net/server_abort_signal`, `internal/http`, `internal/http/FakeSocket`, `node:_http_incoming`, `node:_http_outgoing`, `internal/http1_server_fallback`, `internal/tls`, `node:tls`, `node:net`, `node:cluster`, `internal/trace_events`, `internal/timers`, `internal/util/inspect` |
| `src/js/node/_stream_duplex.ts` | `internal/streams/duplex` |
| `src/js/node/_stream_passthrough.ts` | `internal/streams/passthrough` |
| `src/js/node/_stream_readable.ts` | `internal/streams/readable` |
| `src/js/node/_stream_transform.ts` | `internal/streams/transform` |
| `src/js/node/_stream_wrap.ts` | `node:stream` |
| `src/js/node/_stream_writable.ts` | `internal/streams/writable` |
| `src/js/node/_tls_common.ts` | `node:tls` |
| `src/js/node/_tls_wrap.ts` | `node:tls` |
| `src/js/node/assert.strict.ts` | `node:assert` |
| `src/js/node/assert.ts` | `node:util/types`, `internal/assert/utils`, `internal/validators`, `internal/util/comparisons`, `internal/util/inspect`, `internal/assert/assertion_error`, `internal/util/deprecate`, `internal/assert/calltracker` |
| `src/js/node/async_hooks.ts` | `internal/validators`, `bun:jsc`, `node:path`, `internal/async_hooks_tick`, `internal/async_hooks` |
| `src/js/node/child_process.ts` | `node:events`, `internal/shared`, `internal/validators`, `node:os`, `node:util`, `internal/streams/writable`, `internal/fs/streams`, `internal/streams/readable`, `internal/streams/native-readable`, `node:net`, `internal/util/inspect` |
| `src/js/node/cluster.ts` | `internal/cluster/isPrimary`, `internal/cluster/child`, `internal/cluster/primary` |
| `src/js/node/console.ts` | — |
| `src/js/node/crypto.ts` | `node:string_decoder`, `internal/streams/lazy_transform`, `internal/promisify`, `internal/streams/writable`, `internal/validators`, `internal/util/deprecate`, `node:util/types` |
| `src/js/node/dgram.ts` | `internal/dgram`, `internal/shared`, `internal/validators`, `internal/net/isIP`, `node:events`, `internal/util/deprecate`, `node:util`, `node:cluster`, `node:dns`, `node:diagnostics_channel` |
| `src/js/node/diagnostics_channel.ts` | `internal/validators` |
| `src/js/node/dns.promises.ts` | `node:dns` |
| `src/js/node/dns.ts` | `internal/net/isIP`, `internal/shared`, `internal/validators` |
| `src/js/node/domain.ts` | `node:events` |
| `src/js/node/events.ts` | `internal/validators`, `internal/abort_listener`, `internal/shared`, `node:util/types`, `internal/util/inspect`, `internal/fixed_queue`, `node:async_hooks` |
| `src/js/node/fs.promises.ts` | `node:util/types`, `node:events`, `internal/fs/binding`, `internal/validators`, `internal/fs/watch`, `internal/fs/cp-sync`, `internal/fs/cp`, `node:fs`, `internal/fs/glob`, `node:path`, `node:readline`, `internal/fs/streams`, `internal/streams/iter/utils`, `internal/streams/iter/pull`, `internal/shared` |
| `src/js/node/fs.ts` | `node:util/types`, `internal/validators`, `internal/fs/binding`, `internal/fs/glob`, `internal/shared`, `node:fs/promises`, `node:path`, `internal/fs/cp-sync`, `internal/fs/watch`, `internal/fs/watchfile`, `internal/fs/streams`, `internal/streams/fast-utf8-stream` |
| `src/js/node/http.ts` | `internal/validators`, `node:_http_agent`, `node:_http_client`, `node:_http_common`, `node:_http_incoming`, `node:_http_outgoing`, `node:_http_server`, `internal/http`, `node:util`, `node:https` |
| `src/js/node/http2.ts` | `node:util/types`, `internal/shared`, `internal/http`, `internal/timers`, `node:tls`, `node:net`, `node:fs`, `node:fs/promises`, `node:stream`, `node:diagnostics_channel`, `node:events`, `internal/primordials`, `internal/promisify`, `node:_http2_upgrade`, `internal/async_context_frame`, `internal/validators`, `node:async_hooks`, `internal/abort_listener`, `node:util`, `internal/http1_server_fallback` |
| `src/js/node/https.ts` | `node:http`, `internal/net/isIP`, `internal/url`, `internal/shared`, `internal/validators`, `internal/http`, `node:_http_common`, `node:util/types`, `node:tls`, `node:net`, `node:_http_agent` |
| `src/js/node/inspector.promises.ts` | `node:inspector` |
| `src/js/node/inspector.ts` | `internal/shared`, `internal/validators`, `internal/primordials`, `node:events`, `node:url`, `node:path`, `node:worker_threads`, `internal/trace_events` |
| `src/js/node/net.ts` | `internal/streams/duplex`, `node:util/types`, `internal/streams/state`, `node:events`, `internal/shared`, `internal/timers`, `internal/validators`, `internal/net/isIP`, `internal/net/symbols`, `internal/net/server_abort_signal`, `node:fs`, `internal/async_hooks`, `internal/tls`, `internal/trace_events`, `internal/abort_listener`, `node:dns`, `internal/cluster/isPrimary`, `node:cluster` |
| `src/js/node/os.ts` | — |
| `src/js/node/path.posix.ts` | `node:path` |
| `src/js/node/path.ts` | `internal/validators` |
| `src/js/node/path.win32.ts` | `node:path` |
| `src/js/node/perf_hooks.ts` | `internal/shared`, `internal/validators`, `internal/perf_hooks/monitorEventLoopDelay` |
| `src/js/node/punycode.ts` | — |
| `src/js/node/querystring.ts` | — |
| `src/js/node/quic.ts` | `internal/quic/quic` |
| `src/js/node/readline.js` | `internal/repl/node-primordials`, `internal/readline/callbacks`, `internal/readline/emitKeypressEvents`, `internal/repl/node-errors`, `internal/repl/node-inspect`, `internal/shared`, `internal/validators`, `internal/readline/interface`, `internal/abort_listener`, `node:readline/promises`, `internal/readline/utils`, `internal/util/inspect`, `node:util` |
| `src/js/node/readline.promises.js` | `internal/repl/node-primordials`, `internal/readline/promises`, `internal/readline/interface`, `internal/repl/node-errors`, `internal/validators`, `internal/shared`, `internal/abort_listener` |
| `src/js/node/repl.js` | `internal/repl/node-primordials`, `internal/repl/mode`, `internal/repl/node-shims`, `internal/repl/acorn`, `internal/repl/acorn-walk`, `internal/repl/node-inspect`, `node:vm`, `node:path`, `node:fs`, `node:readline`, `internal/readline/utils`, `node:console`, `node:async_hooks`, `internal/repl/node-errors`, `internal/validators`, `internal/repl/utils`, `internal/readline/interface`, `node:url`, `internal/repl/await`, `internal/repl/completion`, `internal/repl` |
| `src/js/node/stream.consumers.ts` | — |
| `src/js/node/stream.iter.ts` | `internal/streams/iter/types`, `internal/streams/iter/push`, `internal/streams/iter/duplex`, `internal/streams/iter/from`, `internal/streams/iter/pull`, `internal/streams/iter/consumers`, `internal/streams/iter/classic`, `internal/streams/iter/broadcast`, `internal/streams/iter/share` |
| `src/js/node/stream.promises.ts` | `node:stream` |
| `src/js/node/stream.ts` | `internal/stream`, `internal/streams/end-of-stream` |
| `src/js/node/stream.web.ts` | — |
| `src/js/node/test.ts` | `internal/shared`, `internal/validators`, `node:stream`, `node:path`, `internal/test_runner/mock_timers`, `node:assert`, `internal/assert/utils`, `node:async_hooks`, `internal/abort_listener` |
| `src/js/node/timers.promises.ts` | `internal/validators`, `internal/shared` |
| `src/js/node/timers.ts` | `internal/shared`, `node:timers/promises` |
| `src/js/node/tls.ts` | `node:util/types`, `node:net`, `internal/streams/duplex`, `node:events`, `internal/shared`, `internal/url`, `internal/tls`, `internal/validators`, `internal/net/symbols`, `node:crypto` |
| `src/js/node/trace_events.ts` | `internal/trace_events` |
| `src/js/node/tty.ts` | `internal/validators`, `internal/fs/streams`, `node:readline`, `internal/tty` |
| `src/js/node/url.ts` | `internal/url`, `internal/validators`, `internal/shared`, `node:path` |
| `src/js/node/util.ts` | `node:util/types`, `internal/promisify`, `internal/validators`, `internal/shared`, `internal/util/mime`, `internal/util/deprecate`, `internal/primordials`, `internal/util/comparisons`, `internal/util/inspect`, `internal/streams/utils`, `internal/util/colors`, `node:os` |
| `src/js/node/v8.ts` | `internal/shared`, `internal/validators`, `node:util/types`, `bun:jsc`, `node:stream`, `node:os`, `internal/serialization_buffers`, `node:worker_threads`, `node:fs` |
| `src/js/node/vm.ts` | `internal/validators`, `bun:jsc`, `node:util`, `internal/primordials` |
| `src/js/node/wasi.ts` | `node:tty`, `node:fs`, `node:path` |
| `src/js/node/worker_threads.ts` | `node:events`, `internal/primordials`, `internal/shared`, `internal/validators`, `node:path`, `internal/worker/stdio`, `node:console`, `node:fs`, `internal/worker/messaging`, `internal/trace_events`, `internal/streams/readable` |
| `src/js/node/zlib.iter.ts` | `internal/streams/iter/transform` |
| `src/js/node/zlib.ts` | `node:buffer`, `node:stream`, `internal/validators` |

## JavaScript source line checkpoints

These are the reachable files to translate from Bun’s JavaScript module layer. “Lines” is the current source-file line count at the snapshot; “Through” is the last source line translated (0 means untouched).

| Source file | Lines | Through |
|---|---:|---:|
| `src/js/internal/abort_listener.ts` | 31 | 31 |
| `src/js/internal/assert/assertion_error.ts` | 448 | 0 |
| `src/js/internal/assert/calltracker.ts` | 136 | 0 |
| `src/js/internal/assert/myers_diff.ts` | 42 | 0 |
| `src/js/internal/assert/utils.ts` | 35 | 35 |
| `src/js/internal/async_context_frame.ts` | 61 | 0 |
| `src/js/internal/async_hooks.ts` | 28 | 28 |
| `src/js/internal/async_hooks_tick.ts` | 20 | 0 |
| `src/js/internal/cluster/RoundRobinHandle.ts` | 222 | 222 |
| `src/js/internal/cluster/SharedHandle.ts` | 71 | 71 |
| `src/js/internal/cluster/Worker.ts` | 60 | 60 |
| `src/js/internal/cluster/child.ts` | 376 | 376 |
| `src/js/internal/cluster/isPrimary.ts` | 4 | 4 |
| `src/js/internal/cluster/primary.ts` | 380 | 380 |
| `src/js/internal/dgram.ts` | 146 | 146 |
| `src/js/internal/errors.ts` | 23 | 23 |
| `src/js/internal/fixed_queue.ts` | 130 | 0 |
| `src/js/internal/freelist.ts` | 27 | 27 |
| `src/js/internal/fs/binding.ts` | 123 | 0 |
| `src/js/internal/fs/cp-sync.ts` | 555 | 555 |
| `src/js/internal/fs/cp.ts` | 404 | 404 |
| `src/js/internal/fs/glob.ts` | 2982 | 2982 |
| `src/js/internal/fs/streams.ts` | 846 | 846 |
| `src/js/internal/fs/watch.ts` | 272 | 272 |
| `src/js/internal/fs/watchfile.ts` | 112 | 112 |
| `src/js/internal/http.ts` | 568 | 568 |
| `src/js/internal/http/FakeSocket.ts` | 141 | 141 |
| `src/js/internal/http1_server_fallback.ts` | 732 | 732 |
| `src/js/internal/linkedlist.ts` | 47 | 47 |
| `src/js/internal/net/isIP.ts` | 37 | 37 |
| `src/js/internal/net/server_abort_signal.ts` | 26 | 0 |
| `src/js/internal/net/symbols.ts` | 15 | 0 |
| `src/js/internal/perf_hooks/monitorEventLoopDelay.ts` | 71 | 71 |
| `src/js/internal/primordials.js` | 176 | 176 |
| `src/js/internal/promisify.ts` | 81 | 1 |
| `src/js/internal/quic/binding.ts` | 2 | 2 |
| `src/js/internal/quic/diagnostics.ts` | 70 | 70 |
| `src/js/internal/quic/http2util.ts` | 197 | 197 |
| `src/js/internal/quic/quic.ts` | 5565 | 5565 |
| `src/js/internal/quic/state.ts` | 886 | 886 |
| `src/js/internal/quic/stats.ts` | 983 | 983 |
| `src/js/internal/quic/symbols.ts` | 70 | 70 |
| `src/js/internal/readline/callbacks.js` | 122 | 122 |
| `src/js/internal/readline/emitKeypressEvents.js` | 94 | 94 |
| `src/js/internal/readline/interface.js` | 1555 | 1555 |
| `src/js/internal/readline/promises.js` | 138 | 138 |
| `src/js/internal/readline/utils.js` | 595 | 595 |
| `src/js/internal/repl.js` | 62 | 62 |
| `src/js/internal/repl/acorn-walk.js` | 33 | 33 |
| `src/js/internal/repl/acorn.js` | 40 | 40 |
| `src/js/internal/repl/await.js` | 228 | 228 |
| `src/js/internal/repl/completion.js` | 839 | 839 |
| `src/js/internal/repl/history.js` | 455 | 455 |
| `src/js/internal/repl/mode.js` | 6 | 6 |
| `src/js/internal/repl/node-errors.js` | 99 | 99 |
| `src/js/internal/repl/node-inspect.js` | 34 | 34 |
| `src/js/internal/repl/node-primordials.js` | 136 | 136 |
| `src/js/internal/repl/node-shims.js` | 388 | 388 |
| `src/js/internal/repl/utils.js` | 836 | 836 |
| `src/js/internal/serialization_buffers.ts` | 74 | 74 |
| `src/js/internal/async_context_frame.ts` | 61 | 61 |
| `src/js/internal/async_hooks.ts` | 28 | 28 |
| `src/js/internal/shared.ts` | 456 | 456 |
| `src/js/internal/stream.promises.ts` | 45 | 45 |
| `src/js/internal/stream.ts` | 112 | 112 |
| `src/js/internal/streams/add-abort-signal.ts` | 58 | 58 |
| `src/js/internal/streams/end-of-stream.ts` | 384 | 384 |
| `src/js/internal/streams/compose.ts` | 221 | 221 |
| `src/js/internal/streams/destroy.ts` | 343 | 343 |
| `src/js/internal/streams/duplex.ts` | 166 | 166 |
| `src/js/internal/streams/duplexify.ts` | 395 | 395 |
| `src/js/internal/streams/duplexpair.ts` | 78 | 78 |
| `src/js/internal/streams/end-of-stream.ts` | 384 | 384 |
| `src/js/internal/streams/fast-utf8-stream.ts` | 862 | 862 |
| `src/js/internal/streams/from.ts` | 200 | 200 |
| `src/js/internal/streams/iter/broadcast.ts` | 785 | 785 |
| `src/js/internal/streams/iter/classic.ts` | 898 | 898 |
| `src/js/internal/streams/iter/consumers.ts` | 495 | 495 |
| `src/js/internal/streams/iter/duplex.ts` | 142 | 142 |
| `src/js/internal/streams/iter/from.ts` | 576 | 576 |
| `src/js/internal/streams/iter/pull.ts` | 977 | 977 |
| `src/js/internal/streams/iter/push.ts` | 747 | 747 |
| `src/js/internal/streams/iter/ringbuffer.ts` | 147 | 147 |
| `src/js/internal/streams/iter/share.ts` | 732 | 732 |
| `src/js/internal/streams/iter/transform.ts` | 760 | 760 |
| `src/js/internal/streams/iter/types.ts` | 87 | 87 |
| `src/js/internal/streams/iter/utils.ts` | 270 | 270 |
| `src/js/internal/streams/lazy_transform.ts` | 53 | 53 |
| `src/js/internal/streams/legacy.ts` | 146 | 146 |
| `src/js/internal/shared.ts` | 456 | 62 |
| `src/js/internal/streams/native-readable.ts` | 306 | 306 |
| `src/js/internal/streams/operators.ts` | 374 | 374 |
| `src/js/internal/streams/passthrough.ts` | 20 | 20 |
| `src/js/internal/streams/pipeline.ts` | 448 | 448 |
| `src/js/internal/streams/readable.ts` | 1734 | 1734 |
| `src/js/internal/streams/state.ts` | 52 | 52 |
| `src/js/internal/streams/transform.ts` | 173 | 173 |
| `src/js/internal/streams/utils.ts` | 305 | 305 |
| `src/js/internal/validators.ts` | 180 | 140 |
| `src/js/internal/streams/writable.ts` | 1118 | 1118 |
| `src/js/internal/test_runner/mock_timers.ts` | 781 | 781 |
| `src/js/internal/timers.ts` | 32 | 32 |
| `src/js/internal/tls.ts` | 476 | 0 |
| `src/js/internal/trace_events.ts` | 874 | 0 |
| `src/js/internal/tty.ts` | 205 | 205 |
| `src/js/internal/url.ts` | 33 | 33 |
| `src/js/internal/util/colors.ts` | 38 | 38 |
| `src/js/internal/util/comparisons.ts` | 9 | 9 |
| `src/js/internal/util/deprecate.ts` | 60 | 0 |
| `src/js/internal/util/inspect.js` | 2828 | 2828 |
| `src/js/internal/util/mime.ts` | 6 | 0 |
| `src/js/internal/validators.ts` | 180 | 140 |
| `src/js/internal/webstreams_adapters.ts` | 864 | 864 |
| `src/js/internal/worker/messaging.ts` | 301 | 0 |
| `src/js/internal/worker/stdio.ts` | 152 | 0 |
| `src/js/node/_http2_upgrade.ts` | 405 | 405 |
| `src/js/node/_http_agent.ts` | 626 | 626 |
| `src/js/node/_http_client.ts` | 1187 | 1187 |
| `src/js/node/_http_common.ts` | 361 | 361 |
| `src/js/node/_http_incoming.ts` | 737 | 737 |
| `src/js/node/_http_outgoing.ts` | 1335 | 1335 |
| `src/js/node/_http_server.ts` | 4306 | 4306 |
| `src/js/node/_stream_duplex.ts` | 3 | 0 |
| `src/js/node/_stream_passthrough.ts` | 3 | 0 |
| `src/js/node/_stream_readable.ts` | 3 | 0 |
| `src/js/node/_stream_transform.ts` | 3 | 0 |
| `src/js/node/_stream_wrap.ts` | 5 | 0 |
| `src/js/node/_stream_writable.ts` | 3 | 0 |
| `src/js/node/_tls_common.ts` | 43 | 0 |
| `src/js/node/_tls_wrap.ts` | 14 | 0 |
| `src/js/node/assert.strict.ts` | 2 | 2 |
| `src/js/node/assert.ts` | 981 | 981 |
| `src/js/node/async_hooks.ts` | 634 | 634 |
| `src/js/node/child_process.ts` | 2071 | 0 |
| `src/js/node/cluster.ts` | 23 | 23 |
| `src/js/node/console.ts` | 1 | 1 |
| `src/js/node/crypto.ts` | 562 | 562 |
| `src/js/internal/streams/lazy_transform.ts` | 53 | 53 |
| `src/js/node/dgram.ts` | 1546 | 1546 |
| `src/js/node/diagnostics_channel.ts` | 419 | 419 |
| `src/js/node/dns.promises.ts` | 2 | 2 |
| `src/js/node/dns.ts` | 1071 | 1071 |
| `src/js/node/domain.ts` | 96 | 96 |
| `src/js/node/events.ts` | 1057 | 1057 |
| `src/js/node/fs.promises.ts` | 1674 | 1674 |
| `src/js/node/fs.ts` | 1509 | 1509 |
| `src/js/node/http.ts` | 138 | 138 |
| `src/js/node/http2.ts` | 6794 | 6794 |
| `src/js/node/https.ts` | 636 | 636 |
| `src/js/node/inspector.promises.ts` | 28 | 28 |
| `src/js/node/inspector.ts` | 760 | 760 |
| `src/js/node/net.ts` | 4634 | 447 |
| `src/js/node/os.ts` | 179 | 179 |
| `src/js/node/path.posix.ts` | 2 | 2 |
| `src/js/node/path.ts` | 81 | 81 |
| `src/js/node/path.win32.ts` | 2 | 2 |
| `src/js/node/perf_hooks.ts` | 381 | 381 |
| `src/js/node/punycode.ts` | 414 | 414 |
| `src/js/node/querystring.ts` | 523 | 523 |
| `src/js/node/quic.ts` | 55 | 55 |
| `src/js/node/readline.js` | 546 | 546 |
| `src/js/node/readline.promises.js` | 60 | 60 |
| `src/js/node/repl.js` | 1543 | 1543 |
| `src/js/node/stream.consumers.ts` | 54 | 54 |
| `src/js/node/stream.iter.ts` | 165 | 165 |
| `src/js/node/stream.promises.ts` | 2 | 2 |
| `src/js/node/stream.ts` | 8 | 8 |
| `src/js/node/stream.web.ts` | 20 | 20 |
| `src/js/node/test.ts` | 2874 | 2874 |
| `src/js/node/timers.promises.ts` | 253 | 253 |
| `src/js/node/timers.ts` | 41 | 41 |
| `src/js/node/tls.ts` | 1757 | 1757 |
| `src/js/node/trace_events.ts` | 87 | 87 |
| `src/js/node/tty.ts` | 204 | 204 |
| `src/js/node/url.ts` | 1382 | 1382 |
| `src/js/node/util.ts` | 786 | 786 |
| `src/js/node/v8.ts` | 450 | 450 |
| `src/js/node/vm.ts` | 674 | 674 |
| `src/js/node/wasi.ts` | 1775 | 1775 |
| `src/js/node/worker_threads.ts` | 1296 | 0 |
| `src/js/node/zlib.iter.ts` | 51 | 51 |
| `src/js/node/zlib.ts` | 1021 | 0 |

## Direct native bridge references

These rows come from `$cpp`, `$newCppFunction`, and `$rust` source literals in reachable module files. They identify the first native handoff; dependencies called by those bindings belong to Bun’s runtime substrate and are outside this source inventory.

| Importing JS source | Bridge | Literal target | Resolved source path |
|---|---|---|---|
| `src/js/internal/assert/myers_diff.ts` | Rust | `node_assert_binding.rs` | `src/runtime/node/node_assert_binding.rs` |
| `src/js/internal/async_context_frame.ts` | C++ | `ModuleGraph.cpp` | `src/jsc/bindings/ModuleGraph.cpp` |
| `src/js/internal/fs/binding.ts` | Rust | `node_fs_binding.rs` | `src/runtime/node/node_fs_binding.rs` |
| `src/js/internal/http.ts` | C++ | `NodeHTTP.cpp` | `src/jsc/bindings/NodeHTTP.cpp` |
| `src/js/internal/perf_hooks/monitorEventLoopDelay.ts` | C++ | `JSNodePerformanceHooksHistogramPrototype.cpp` | `src/jsc/bindings/JSNodePerformanceHooksHistogramPrototype.cpp` |
| `src/js/internal/quic/binding.ts` | Rust | `node_quic_binding.rs` | `src/runtime/node/node_quic_binding.rs` |
| `src/js/internal/repl/node-inspect.js` | C++ | `stringWidth.cpp` | `src/jsc/bindings/stringWidth.cpp` |
| `src/js/internal/repl/node-shims.js` | C++ | `UtilInspect.cpp` | `src/jsc/bindings/UtilInspect.cpp` |
| `src/js/internal/shared.ts` | C++ | `BunProcess.cpp` | `src/jsc/bindings/BunProcess.cpp` |
| `src/js/internal/shared.ts` | C++ | `ModuleGraph.cpp` | `src/jsc/bindings/ModuleGraph.cpp` |
| `src/js/internal/streams/iter/transform.ts` | Rust | `node_zlib_binding.rs` | `src/runtime/node/node_zlib_binding.rs` |
| `src/js/internal/streams/native-readable.ts` | C++ | `streams/BunStreamConsumers.cpp` | `src/jsc/bindings/webcore/streams/BunStreamConsumers.cpp` |
| `src/js/internal/streams/readable.ts` | C++ | `NodeModuleModule.cpp` | `src/jsc/modules/NodeModuleModule.cpp` |
| `src/js/internal/tls.ts` | Rust | `SecureContext.rs` | `src/runtime/api/bun/SecureContext.rs` |
| `src/js/internal/url.ts` | C++ | `NodeURL.cpp` | `src/jsc/bindings/NodeURL.cpp` |
| `src/js/internal/util/comparisons.ts` | C++ | `NodeUtilTypesModule.cpp` | `src/jsc/modules/NodeUtilTypesModule.cpp` |
| `src/js/internal/util/inspect.js` | C++ | `UtilInspect.cpp` | `src/jsc/bindings/UtilInspect.cpp` |
| `src/js/internal/util/inspect.js` | C++ | `stringWidth.cpp` | `src/jsc/bindings/stringWidth.cpp` |
| `src/js/internal/util/mime.ts` | C++ | `JSMIMEParams.cpp` | `src/jsc/bindings/webcore/JSMIMEParams.cpp` |
| `src/js/internal/validators.ts` | C++ | `NodeValidator.cpp` | `src/jsc/bindings/NodeValidator.cpp` |
| `src/js/node/assert.ts` | C++ | `NodeUtilTypesModule.cpp` | `src/jsc/modules/NodeUtilTypesModule.cpp` |
| `src/js/node/async_hooks.ts` | C++ | `NodeAsyncHooks.cpp` | `src/jsc/bindings/NodeAsyncHooks.cpp` |
| `src/js/node/crypto.ts` | C++ | `node_crypto_binding.cpp` | `src/jsc/bindings/node/crypto/node_crypto_binding.cpp` |
| `src/js/node/crypto.ts` | Rust | `node_crypto_binding.rs` | `src/runtime/node/node_crypto_binding.rs` |
| `src/js/node/dgram.ts` | Rust | `node_net_binding.rs` | `src/runtime/node/node_net_binding.rs` |
| `src/js/node/events.ts` | C++ | `JSEventTarget.cpp` | `src/jsc/bindings/webcore/JSEventTarget.cpp` |
| `src/js/node/events.ts` | C++ | `JSEventTargetNode.cpp` | `src/jsc/bindings/webcore/JSEventTargetNode.cpp` |
| `src/js/node/http2.ts` | Rust | `h2_frame_parser.rs` | `src/runtime/api/bun/h2_frame_parser.rs` |
| `src/js/node/https.ts` | Rust | `node_net_binding.rs` | `src/runtime/node/node_net_binding.rs` |
| `src/js/node/inspector.ts` | C++ | `BunDebugger.cpp` | `src/jsc/bindings/BunDebugger.cpp` |
| `src/js/node/inspector.ts` | C++ | `JSInspectorProfiler.cpp` | `src/jsc/bindings/JSInspectorProfiler.cpp` |
| `src/js/node/net.ts` | Rust | `node_net_binding.rs` | `src/runtime/node/node_net_binding.rs` |
| `src/js/node/os.ts` | Rust | `node_os.rs` | `src/runtime/node/node_os.rs` |
| `src/js/node/path.ts` | C++ | `Path.cpp` | `src/jsc/bindings/Path.cpp` |
| `src/js/node/perf_hooks.ts` | C++ | `JSNodePerformanceHooksHistogram.cpp` | `src/jsc/bindings/JSNodePerformanceHooksHistogram.cpp` |
| `src/js/node/tls.ts` | C++ | `NodeTLS.cpp` | `src/jsc/bindings/NodeTLS.cpp` |
| `src/js/node/tls.ts` | Rust | `SecureContext.rs` | `src/runtime/api/bun/SecureContext.rs` |
| `src/js/node/tty.ts` | C++ | `ProcessBindingTTYWrap.cpp` | `src/jsc/bindings/ProcessBindingTTYWrap.cpp` |
| `src/js/node/util.ts` | C++ | `NodeUtilTypesModule.cpp` | `src/jsc/modules/NodeUtilTypesModule.cpp` |
| `src/js/node/v8.ts` | C++ | `NodeV8.cpp` | `src/jsc/bindings/NodeV8.cpp` |
| `src/js/node/vm.ts` | C++ | `NodeVM.cpp` | `src/jsc/bindings/NodeVM.cpp` |
| `src/js/node/worker_threads.ts` | C++ | `Worker.cpp` | `src/jsc/bindings/webcore/Worker.cpp` |
| `src/js/node/zlib.ts` | Rust | `node_zlib_binding.rs` | `src/runtime/node/node_zlib_binding.rs` |

### Native/synthetic-only module source files

These seven resolver entries have no `src/js/node` wrapper. Each source row has an independent line checkpoint. SQLite’s native implementation and bundled engine files are included because they implement the module API.

| Native/synthetic module | Bun source file | Lines | Through |
|---|---|---:|---:|
| `node:buffer` | `src/jsc/modules/NodeBufferModule.cpp` | 282 | 282 |
| `node:buffer` | `src/jsc/modules/NodeBufferModule.h` | 210 | 210 |
| `node:buffer` | `src/runtime/node/buffer.rs` | 141 | 141 |
| `node:constants` | `src/jsc/modules/NodeConstantsModule.h` | 981 | 981 |
| `node:process` | `src/jsc/bindings/ProcessBindingConstants.cpp` | 1265 | 1265 |
| `node:process` | `src/jsc/bindings/ProcessBindingConstants.h` | 36 | 36 |
| `node:module` | `src/jsc/modules/NodeModuleModule.cpp` | 1257 | 1257 |
| `node:module` | `src/jsc/modules/NodeModuleModule.h` | 45 | 45 |
| `node:module` | `src/jsc/NodeModuleModule.rs` | 239 | 239 |
| `node:module` | `src/jsc/NodeModuleModule.bind.ts` | 8 | 8 |
| `node:url` | `src/jsc/bindings/NodeURL.cpp` | 308 | 308 |
| `node:url` | `src/jsc/bindings/NodeURL.h` | 10 | 10 |
| `node:url` | `src/jsc/bindings/NodeURLHelpers.h` | 14 | 14 |
| `node:url` | `src/jsc/bindings/ASCIIHostPunycodeCheck.h` | 222 | 222 |
| `node:process` | `src/jsc/modules/NodeProcessModule.h` | 24 | 24 |
| `node:process` | `src/runtime/node/node_process.rs` | 595 | 595 |
| `node:sqlite` | `src/jsc/modules/NodeSqliteModule.h` | 38 | 38 |
| `node:sqlite` | `src/jsc/bindings/sqlite/NodeSqlite.cpp` | 4135 | 0 |
| `node:sqlite` | `src/jsc/bindings/sqlite/NodeSqlite.h` | 902 | 0 |
| `node:sqlite` | `src/jsc/bindings/sqlite/JSSQLStatement.cpp` | 3077 | 0 |
| `node:sqlite` | `src/jsc/bindings/sqlite/JSSQLStatement.h` | 89 | 0 |
| `node:sqlite` | `src/jsc/bindings/sqlite/sqlite3.c` | 269669 | 0 |
| `node:sqlite` | `src/jsc/bindings/sqlite/sqlite3_local.h` | 14350 | 0 |
| `node:sqlite` | `src/jsc/bindings/sqlite/sqlite3_error_codes.h` | 764 | 0 |
| `node:sqlite` | `src/jsc/bindings/sqlite/lazy_sqlite3.h` | 521 | 0 |
| `node:string_decoder` | `src/jsc/modules/NodeStringDecoderModule.h`; `src/jsc/bindings/JSStringDecoder.h`; `.cpp` | 774 | NodeStringDecoderModule.h:17/17; JSStringDecoder.h:126/126; JSStringDecoder.cpp:631/631 | source sweep complete; stateful decoder API adapted | Go UTF-8/UTF-16 decoders, Buffer encoding helpers |
| `node:stream` | `src/js/node/stream.ts`; `src/js/internal/stream.ts`; `src/js/internal/stream.promises.ts`; `src/js/internal/streams/operators.ts`; `src/js/internal/streams/duplexpair.ts` | 617 | 8/8; 112/112; 45/45; 374/374; 78/78 | API surface mapped; stream operator concurrency, abort cleanup, and pair scheduling are approximate | shared stream runtime, iterators, event loop |
| `node:stream/consumers` | `src/js/node/stream.consumers.ts` | 54 | 54/54 | arrayBuffer, bytes, text, JSON, Buffer, and Blob consumer methods exposed; stream errors/decoding edge cases are approximate | async iterator, Buffer, global Blob |
| `node:stream/promises` | `src/js/node/stream.promises.ts` | 2 | 2/2 | default exports promise `pipeline` and `finished` | `node:stream` |
| `node:stream/web` | `src/js/node/stream.web.ts` | 20 | 20/20 | constructor exports forward global Web Stream values; unavailable globals stay undefined | global Web Stream values |
| `node:crypto` | `src/js/node/crypto.ts` | 562 | 562/562 | exported wrapper API mapped; unsupported Argon2/scrypt, classic DH, SPKAC, FIPS/secure heap, and subtle WebCrypto remain stubs or explicit errors | Go standard `crypto/*`, random source, hash and key formats; C++/Rust binding behavior replaced rather than line-translated |
| `internal/streams/lazy_transform` | `src/js/internal/streams/lazy_transform.ts` | 53 | 53/53 | lazy Transform initialization, accessors, and constructor surface adapted | `internal/streams/transform` |
| `internal/streams/transform` | `src/js/internal/streams/transform.ts` | 173 | 173/173 | constructor, flush, prefinish, callback backpressure, and read/write hooks adapted | `internal/streams/duplex`, `internal/streams/state` |
| `internal/streams/duplex` | `src/js/internal/streams/duplex.ts` | 166 | 166/166 | constructor, combined prototype, duplex-side flags, and Web statics adapted | `internal/streams/legacy`, Readable, Writable, Web adapters, Duplexify |
| `internal/streams/duplexify` | `src/js/internal/streams/duplexify.ts` | 395 | 395/395 | Node/Web stream, iterable, promise, and readable/writable pair conversion adapted; async-generator handling is approximate | Duplex, Readable, Writable, `internal/streams/from`, EOS, destroy |
| `internal/streams/writable` | `src/js/internal/streams/writable.ts` | 1118 | 1118/1118 | constructor, queue, writev batching, lifecycle, properties, and Web statics adapted; callback ordering and WritableState internals are approximate | legacy, destroy, EOS, abort signal, state, utils, Web adapters |
| `internal/streams/native-readable` | `src/js/internal/streams/native-readable.ts` | 306 | 306/306 | Readable wrapper maps pointer methods and pull results; Bun C++ transfer entry is a placeholder | Readable, destroy, `BunStreamConsumers.cpp` |
| `internal/webstreams_adapters` | `src/js/internal/webstreams_adapters.ts` | 864 | 864/864 | all six classic/Web conversion entry points exposed; strategy and native transfer behavior are approximate | Writable, Readable, Duplex, destroy, utils, shared, validators, util/types, EOS, native-readable |
| `internal/streams/legacy` | `src/js/internal/streams/legacy.ts` | 146 | 146/146 | source mapped; EventEmitter bridge, pipe, and Buffer helpers adapted | `node:events`, `node:util/types`, Buffer |
| `internal/abort_listener` | `src/js/internal/abort_listener.ts` | 31 | 31/31 | addAbortListener and Disposable removal adapted | AbortSignal event methods, shared stop-propagation option |
| `internal/async_hooks` (private) | `src/js/internal/async_hooks.ts` | 28 | 28/28 | enabled hook counter and owner symbol adapted; public hook callbacks remain stubs | `node:async_hooks` adapter state |
| `internal/errors` | `src/js/internal/errors.ts` | 23 | 23/23 | aggregateTwoErrors mapped to local AggregateError values | internal primordials surface |
| `internal/streams/destroy` | `src/js/internal/streams/destroy.ts` | 343 | 343/343 | destroy, construct, error/close scheduling, and legacy destroyer mapped | `internal/errors`, stream flags, process nextTick |
| `internal/streams/state` | `src/js/internal/streams/state.ts` | 52 | 52/52 | default and per-stream high-water marks mapped | validators, process platform |
| `internal/streams/from` | `src/js/internal/streams/from.ts` | 200 | 200/200 | iterable-to-readable, async values, and iterator close bridge mapped | `internal/errors`, async iterator runtime, Readable constructor |
| `internal/streams/readable` | `src/js/internal/streams/readable.ts` | 1734 | 1734/1734 | iterator/compose methods, accessors, state compatibility, buffer extraction, end scheduling, Web statics, and wrap mapped with runtime approximations | classic iterator, compose, Web-stream adapters, destroy |
| `internal/streams/add-abort-signal` | `src/js/internal/streams/add-abort-signal.ts` | 58 | 58/58 | source sweep complete; AbortSignal listener and stream error adapters added | `internal/streams/utils` (complete), `internal/streams/end-of-stream` |
| `internal/streams/iter/types` | `src/js/internal/streams/iter/types.ts` | 87 | 87/87 | shared protocol symbols and private validation sentinels mapped | Symbol registry |
| `internal/streams/iter/utils` | `src/js/internal/streams/iter/utils.ts` | 270 | 270/270 | chunk conversion, cursor, signal, and transform helpers adapted | `node:util/types`, validators |
| `internal/streams/iter/from` | `src/js/internal/streams/iter/from.ts` | 576 | 576/576 | sync/async byte normalization, streamable protocols, ArrayBuffer views, and batch iterators adapted | `node:util/types`, iterator types/utils |
| `internal/streams/iter/classic` | `src/js/internal/streams/iter/classic.ts` | 898 | 898/898 | batched Readable iteration, Web conversion, Writable-to-Writer backpressure, and Writer-to-classic stream bridges adapted | iterator normalization, `node:util/types`, stream runtime |
| `internal/streams/end-of-stream` | `src/js/internal/streams/end-of-stream.ts` | 384 | 384/384 | source sweep complete; classic/web completion, cleanup, `finished()`, and abort paths adapted | `internal/shared`, validators, stream utils, async hooks, abort listener |
| `internal/async_context_frame` | `src/js/internal/async_context_frame.ts` | 61 | 61/61 | source sweep complete; current/exchange/run use event-loop frames; ModuleGraph queries are unavailable | runtime async-context frames |
| `internal/shared` | `src/js/internal/shared.ts` | 456 | 456/456 | source sweep complete; error helpers, once/guard wrappers, lazy state, and Node performance-entry observers adapted; ModuleGraph checks remain approximate | `node:util`, `internal/primordials`, `internal/async_context_frame`, timers, `node:perf_hooks` |
| `internal/tty` | `src/js/internal/tty.ts` | 205 | 205/205 | source sweep complete; terminal color depth and environment precedence adapted | `node:os`, `node:process` |
| `internal/util/colors` | `src/js/internal/util/colors.ts` | 38 | 38/38 | source sweep complete; stream color detection and ANSI values adapted | `internal/tty`, `node:process`, inspect color tables |
| `node:module` | `NodeModuleModule.cpp/.h`, `NodeModuleModule.rs`, `NodeModuleModule.bind.ts` | 1549 | 1549/1549 | source sweep complete; builtin table, Module facade, wrapper/cache/lookup APIs adapted; resolver and engine hooks are approximations | Go `path/filepath`, `os`, generated builtin registry |
| `node:url` | `src/js/node/url.ts`; `src/js/internal/url.ts`; NodeURL C++ bridge and helper headers | 1969 | 1382/1382; 33/33; 308/308; 10/10; 14/14; 222/222 | source sweep complete; WHATWG and legacy URL surfaces adapted; ICU UTS #46 and detailed URLPattern matching approximated | Go `net/url`, `path/filepath`, local punycode and Buffer runtimes |
| `internal/util/inspect` | `src/js/internal/util/inspect.js` | 2828 | 2828/2828 | source sweep complete; object formatting, printf formats, error stacks, ANSI stripping, custom inspection, and helper APIs adapted; exact grapheme widths and some property details remain approximate | primordials, `node:util/types`, validators, `node:module`, `node:url`, `node:buffer`, `UtilInspect.cpp`, `stringWidth.cpp` |
| `node:util` | `src/js/node/util.ts` | 786 | 786/786 | source sweep complete; debug, style, callback, signal, error, codec, inspect, MIME, and legacy predicate exports adapted; call-site capture, decoder state, and weak finalization are approximations | `node:util/types`, `internal/promisify`, `internal/validators`, `internal/shared`, MIME, deprecation, primordials, comparisons, inspect, streams, colors, `node:os`, `node_util_binding.rs` |
| `node:util` parseArgs binding | `src/runtime/node/util/parse_args.rs` | 1049 | 1049/1049 | source sweep complete; tokenization, option checks, defaults, and token records adapted | `parse_args_utils.rs`, validators, process argv |
| `internal/parse_args_utils` | `src/runtime/node/util/parse_args_utils.rs` | 149 | 149/149 | option metadata and token-kind classifier mapped | `parse_args.rs` |
| `internal/util/comparisons` | `src/js/internal/util/comparisons.ts` | 9 | 9/9 | source façade mapped to the local deep comparator; prototype policy adapted | generated runtime value comparison |
| `internal/promisify` | `src/js/internal/promisify.ts` | 81 | 81/81 | source mapped; custom and multi-value callback wrappers emitted | validators, shared custom promisify symbol |
| `internal/validators` | `src/js/internal/validators.ts` | 180 | 180/180 | source mapped; native calls translated through `NodeValidator.cpp` | `internal/shared`, `node:path`, `NodeValidator.cpp` |
| `internal/net/isIP` | `src/js/internal/net/isIP.ts` | 37 | 37/37 | IPv4/IPv6 validation adapted with `net.ParseIP`; IPv6 zone syntax checked locally | Go `net` |
| `internal/timers` | `src/js/internal/timers.ts` | 32 | 32/32 | timeout range validation, overflow warning, and timeout symbol adapted | validators, process warnings |
| `internal/util/mime` | `src/js/internal/util/mime.ts` | 6 | 6/6 | source mapped; MIME binding constructors adapted | JSMIMEParams and JSMIMEType bridge sources |
| `internal/util/deprecate` | `src/js/internal/util/deprecate.ts` | 60 | 60/60 | source mapped; warning-code suppression and wrapper calls adapted | validators, process warnings |
| `internal/streams/utils` | `src/js/internal/streams/utils.ts` | 305 | 305/305 | source mapped; state symbols, flags and stream predicates adapted | stream values, Web streams, shared Symbol registry |
| `node:util/types` | `src/jsc/modules/NodeUtilTypesModule.cpp` | 1289 | 1289 |
| `node:util/types` | `src/jsc/modules/NodeUtilTypesModule.h` | 25 | 25 |


### Direct bridge target line checkpoints

Unique files directly named by `$cpp`, `$newCppFunction`, and `$rust` in the reachable module graph. The native-only routes above may overlap this index.

| Source file | Lines | Through |
|---|---:|---:|
| `src/jsc/bindings/BunDebugger.cpp` | 1065 | 0 |
| `src/jsc/bindings/BunProcess.cpp` | 5172 | 5172 |
| `src/jsc/bindings/JSBuffer.cpp` | 4016 | 4016 |
| `src/jsc/bindings/JSBuffer.h` | 100 | 100 |
| `src/jsc/bindings/JSBufferEncodingType.cpp` | 159 | 159 |
| `src/jsc/bindings/JSBufferEncodingType.h` | 14 | 14 |
| `src/jsc/bindings/BufferEncodingType.h` | 23 | 23 |
| `src/jsc/bindings/JSStringDecoder.cpp` | 631 | 631 |
| `src/jsc/bindings/JSStringDecoder.h` | 126 | 126 |
| `src/jsc/bindings/JSInspectorProfiler.cpp` | 159 | 0 |
| `src/jsc/bindings/JSNodePerformanceHooksHistogram.cpp` | 274 | 0 |
| `src/jsc/bindings/JSNodePerformanceHooksHistogramPrototype.cpp` | 529 | 0 |
| `src/jsc/bindings/ModuleGraph.cpp` | 855 | 0 |
| `src/jsc/bindings/NodeAsyncHooks.cpp` | 21 | 0 |
| `src/jsc/bindings/NodeHTTP.cpp` | 1023 | 0 |
| `src/jsc/bindings/NodeTLS.cpp` | 362 | 0 |
| `src/jsc/bindings/NodeURL.cpp` | 308 | 308 |
| `src/jsc/bindings/NodeURL.h` | 10 | 10 |
| `src/jsc/bindings/NodeURLHelpers.h` | 14 | 14 |
| `src/jsc/bindings/ASCIIHostPunycodeCheck.h` | 222 | 222 |
| `src/jsc/bindings/NodeV8.cpp` | 108 | 0 |
| `src/jsc/bindings/NodeVM.cpp` | 2056 | 0 |
| `src/jsc/bindings/NodeValidator.cpp` | 655 | 655 |
| `src/jsc/bindings/NodeValidator.h` | 48 | 48 |
| `src/runtime/node/node_util_binding.rs` | 246 | 246 |
| `src/runtime/node/util/parse_args.rs` | 1049 | 1049 |
| `src/runtime/node/util/parse_args_utils.rs` | 149 | 149 |
| `src/jsc/bindings/Path.cpp` | 168 | 168 |
| `src/jsc/bindings/ProcessBindingTTYWrap.cpp` | 533 | 0 |
| `src/jsc/bindings/UtilInspect.cpp` | 121 | 121 |
| `src/jsc/bindings/node/crypto/node_crypto_binding.cpp` | 403 | 0 |
| `src/jsc/bindings/stringWidth.cpp` | 1178 | 1178 |
| `src/jsc/bindings/webcore/JSEventTarget.cpp` | 370 | 0 |
| `src/jsc/bindings/webcore/JSEventTargetNode.cpp` | 49 | 0 |
| `src/jsc/bindings/webcore/JSMIMEParams.cpp` | 722 | 722 |
| `src/jsc/bindings/webcore/JSMIMEParams.h` | 89 | 89 |
| `src/jsc/bindings/webcore/JSMIMEType.cpp` | 614 | 614 |
| `src/jsc/bindings/webcore/JSMIMEType.h` | 102 | 102 |
| `src/jsc/bindings/webcore/Worker.cpp` | 470 | 0 |
| `src/jsc/bindings/webcore/streams/BunStreamConsumers.cpp` | 1770 | 0 |
| `src/jsc/modules/NodeModuleModule.cpp` | 1257 | 0 |
| `src/jsc/modules/NodeUtilTypesModule.cpp` | 1289 | 0 |
| `src/runtime/api/bun/SecureContext.rs` | 460 | 0 |
| `src/runtime/api/bun/h2_frame_parser.rs` | 7873 | 0 |
| `src/runtime/node/node_assert_binding.rs` | 147 | 0 |
| `src/runtime/node/node_crypto_binding.rs` | 1750 | 0 |
| `src/runtime/node/node_fs_binding.rs` | 527 | 0 |
| `src/runtime/node/node_net_binding.rs` | 184 | 184 |
| `src/runtime/node/node_os.rs` | 1614 | 0 |
| `src/runtime/node/path.rs` | 956 | 956 |
| `src/runtime/node/node_quic_binding.rs` | 284 | 0 |
| `src/runtime/node/node_zlib_binding.rs` | 1066 | 0 |

## Node runtime Rust source index

The Rust files below are Bun’s Node host integration layer. A module’s concrete dependencies are shown in the JavaScript edge graph and direct bridge table; this index keeps the broader runtime layer available for selection.

| Source file | Lines | Through |
|---|---:|---:|
| `src/runtime/node/Stat.rs` | 212 | 0 |
| `src/runtime/node/StatFS.rs` | 160 | 0 |
| `src/runtime/node/assert/myers_diff.rs` | 512 | 0 |
| `src/runtime/node/buffer.rs` | 141 | 141 |
| `src/runtime/node/dir_iterator.rs` | 821 | 0 |
| `src/runtime/node/fs_events.rs` | 939 | 0 |
| `src/runtime/node/memory_pressure.rs` | 582 | 0 |
| `src/runtime/node/net/BlockList.rs` | 524 | 0 |
| `src/runtime/node/node.classes.ts` | 384 | 0 |
| `src/runtime/node/node_assert.rs` | 383 | 0 |
| `src/runtime/node/node_assert_binding.rs` | 147 | 0 |
| `src/runtime/node/node_cluster_binding.rs` | 810 | 0 |
| `src/runtime/node/node_crypto_binding.rs` | 1750 | 0 |
| `src/runtime/node/node_fs.rs` | 10116 | 0 |
| `src/runtime/node/node_fs_binding.rs` | 527 | 0 |
| `src/runtime/node/node_fs_constant.rs` | 32 | 0 |
| `src/runtime/node/node_fs_stat_watcher.rs` | 1110 | 0 |
| `src/runtime/node/node_fs_watcher.rs` | 1256 | 0 |
| `src/runtime/node/node_http_binding.rs` | 110 | 0 |
| `src/runtime/node/node_net_binding.rs` | 184 | 184 |
| `src/runtime/node/node_os.bind.ts` | 94 | 0 |
| `src/runtime/node/node_os.rs` | 1614 | 0 |
| `src/runtime/node/node_process.rs` | 595 | 0 |
| `src/runtime/node/node_quic_binding.rs` | 284 | 0 |
| `src/runtime/node/node_util_binding.rs` | 246 | 0 |
| `src/runtime/node/node_zlib_binding.rs` | 1066 | 0 |
| `src/runtime/node/path.rs` | 956 | 0 |
| `src/runtime/node/path_watcher.rs` | 1618 | 0 |
| `src/runtime/node/quic/endpoint.rs` | 2661 | 0 |
| `src/runtime/node/quic/ffi.rs` | 103 | 0 |
| `src/runtime/node/quic/mod.rs` | 44 | 0 |
| `src/runtime/node/quic/quic.classes.ts` | 181 | 0 |
| `src/runtime/node/quic/session.rs` | 2680 | 0 |
| `src/runtime/node/quic/stream.rs` | 1188 | 0 |
| `src/runtime/node/quic/tls.rs` | 725 | 0 |
| `src/runtime/node/time_like.rs` | 125 | 0 |
| `src/runtime/node/types.rs` | 1779 | 0 |
| `src/runtime/node/util/parse_args.rs` | 1049 | 1049 |
| `src/runtime/node/util/parse_args_utils.rs` | 149 | 149 |
| `src/runtime/node/util/validators.rs` | 486 | 0 |
| `src/runtime/node/uv_signal_handle_windows.rs` | 66 | 0 |
| `src/runtime/node/win_watcher.rs` | 547 | 0 |
| `src/runtime/node/zlib/NativeBrotli.rs` | 737 | 0 |
| `src/runtime/node/zlib/NativeZlib.rs` | 654 | 0 |
| `src/runtime/node/zlib/NativeZstd.rs` | 687 | 0 |

## `node:path` Rust algorithm crate

The JavaScript path wrapper binds through `src/jsc/bindings/Path.cpp`; Bun’s portable path implementation is in `src/node_path` and host integration is in `src/runtime/node/path.rs`.

| Source file | Lines | Through |
|---|---:|---:|
| `src/node_path/lib.rs` | 2714 | 2714 |
| `src/node_path/Cargo.toml` | 16 | 0 |

## Source-navigation references

These files explain how the source inventory is wired into Bun. They are reference-only for translating module behavior; they are not counted as module implementation lines.

| Bun source file | Why it matters |
|---|---|
| `src/resolve_builtins/HardcodedModule.rs` | Hardcoded Node specifiers and aliases |
| `src/resolve_builtins/lib.rs` | Alias resolution and builtin wiring |
| `src/resolve_builtins/node_builtins.rs` | Node builtin name/target-compatibility lists |
| `src/codegen/internal-module-registry-scanner.ts` | Internal source scan and static module graph |
| `src/codegen/bundle-modules.ts` | Builtin bundling and registry generation |
| `src/runtime/node.rs` | Rust Node runtime module declarations |
| `src/runtime/jsc_hooks.rs` | Synthetic/native module loading hooks |
| `src/js/README.md` | Builtin source conventions and registry model |

## Suggested breadth-first module batches

These clusters keep shared dependencies together while leaving each resolver entry independently checkpointable:

| Batch | Resolver entries | Shared source center |
|---|---|---|
| Path and URL | `node:path`, `node:path/posix`, `node:path/win32`, `node:url` | `path.ts`, `internal/url.ts`, validators, `Path.cpp`, `src/node_path` |
| Events and diagnostics | `node:events`, `node:domain`, `node:diagnostics_channel`, `node:async_hooks` | shared/event/async internals and native async hooks |
| Streams | `node:stream` family, `node:zlib` family, `node:string_decoder` | `internal/stream*`, `internal/streams/*`, web streams adapters |
| Filesystem | `node:fs`, `node:fs/promises` | `internal/fs/*`, path, validators, `node_fs*` bindings |
| Network and HTTP | `node:net`, `node:dgram`, `node:dns*`, `node:http*`, `node:https`, `node:tls`, private wrappers | `internal/net/*`, `internal/http*`, streams, TLS, Rust network/HTTP bindings |
| Process and workers | `node:process`, `node:child_process`, `node:cluster`, `node:worker_threads`, `node:os` | process/worker internals, OS/process bindings |
| Developer tools | `node:util`, `node:assert*`, `node:inspector*`, `node:repl`, `node:readline*`, `node:test`, `node:v8`, `node:vm` | inspect/assert/repl internals and native util/VM/inspector bindings |
| Remaining focused modules | `node:sqlite`, `node:test`, `node:zlib/iter` | module-specific wrapper/native route in the tables above |

No tests were run to produce this inventory.
