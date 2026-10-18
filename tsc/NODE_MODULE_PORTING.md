# Bun Node built-in port inventory

**Bun reference:** `/workspace/bun` local source tree (`main` at `e40a58ad875f91bb95fb05fbd87bd17c6a93566b`).
**Scope:** line-by-line mapping of Node module source and Bun runtime bindings into the Go adapters, with source coverage and behavior gaps recorded in this file.

## Port checkpoint

| Field | Value |
|---|---|
| Modules with registered runtime adapters | 65 |
| Bun source lines mapped in tracked checkpoints | 190682 |
| Active module | `node:dgram` |
| Active source file | `src/js/node/dgram.ts` |
| Active source reviewed through | `dgram.ts:1546/1546`; the UDP surface is mapped; continue closing cluster sharing, native queue, async lookup, and source-specific membership gaps in the Go adapter |
| Go backend reviewed through | `node_dgram_runtime.go:180`; `node_dns_runtime.go:129`; `node_worker_threads_runtime.go:171`; `node_v8_runtime.go:268`; `runtime.go:501`; `node_net_runtime.go:681`; `internal_serialization_buffers_runtime.go` |
| Next source line | `src/js/node/dgram.ts:1547` (EOF); `node:test` is parked per the request to skip this dynamic test-runner module. |

`node:sqlite` binding sources `NodeSqlite.cpp:4135/4135`, `NodeSqlite.h:902/902`, and `JSSQLStatement.cpp/.h:3166/3166` are mapped. The mapped sections cover SQLite library selection/version reporting, text/blob/row conversion, database and statement lifecycle/configuration, scalar and aggregate callbacks, sessions, changeset validation, authorizers, serialization, SQL tag-store creation and LRU behavior, backup validation, path/options validation, named and positional binding, and StatementSync run/get/all/iterate/columns behavior. The Go adapter validates database and backup paths and constructor options, matches method lengths, provides own database/statement/tag-store getters and stable live limit accessors, checks closed-session and tag-store states, uses an LRU cache, and reports finalized/reentrant state errors. The runtime now uses a local cgo-enabled fork of MIT go-sqlite3 v1.14.24 with its bundled SQLite C engine, so SQL parsing/execution, files, parameters, live cursors, limits, transactions, UDFs/window aggregates, authorizers, sessions/changesets, serialization, backup, and tag storage use a real engine. The cgo path requires a C compiler; SQLite’s upstream C implementation is replaced wholesale by that engine and is recorded as an explicit skipped range. `JSSQLStatement.cpp/.h` are Bun’s separate `bun:sqlite` statement API; this pass mapped their binding, result conversion, database registry, and lifecycle sections and recorded the shared SQLite initialization/termination boundary. The first 3,500 lines of `sqlite3_local.h` cover the connection/result-code/flags/VFS/file-control API, library and allocator configuration, database options, rowid and change counters, busy handling, formatting, and authorizer declarations. Lines 3,501–14,350 of `sqlite3_local.h` and all 269,669 lines of `sqlite3.c` are explicitly replaced by the local bundled C engine rather than translated into Go. SQLite storage classes are preserved for Node conversion; the driver adds column-origin metadata, expanded SQL, direct-only scalar/window functions, defensive/DQS configuration, extension loading, sessions, and changeset callbacks. JSC ownership, exact Bun SQLite error wording, and Bun-only `setCustomSQLite` remain host gaps. The `node:net.ts` pass reaches line 4634, and `node_net_binding.rs` reaches 184/184. `BlockList.rs` (524 lines) and its `SocketAddress.rs` dependency (1,084 lines) are now mapped into the Go adapter. SocketAddress option validation, getters, `toJSON`, exact-instance recognition, numeric/string family handling, IPv4/IPv6 validation, and common URL-style address parsing are represented. Its `parse` uses Go `net/url` and a local IPv4 canonicalizer in place of Bun’s WHATWG parser; native sockaddr DTO helpers and named-interface IPv6 scopes use Go address handling or remain absent. BlockList rule mutation/matching is lock-protected; the v8 serializer tags BlockList references and restores the same object through a nonce registry, but the registry has no destruction cleanup. Worker MessagePort delivery clones through the v8 value graph serializer, and transfer lists now move MessagePort endpoints and ArrayBuffer backing stores, preserve shared views, and detach transferred senders after a successful clone.

The `net.ts` source pass reaches line 4634, with its final export surface mapped. The Go adapter routes worker listen requests through the local cluster module and can hand off accepted round-robin sockets when both sides share this runtime. Native descriptor transfer across real process boundaries, exact kernel backlog control, and the private TLS handshake hook remain gaps. The `node:child_process.ts` pass is complete through line 2071, including validators, URL/path helpers, `SystemError`, and the default export. Its Go adapter covers spawn/exec entry points, event lifecycle, sync results, promisify hooks, stdio normalization, and error metadata; `fork()` throws `NotImplementedError` because it launches a JavaScript module in another runtime. Real IPC transport, descriptor-backed extra stdio, detached process groups, uid/gid application, exact Windows argv quoting, maxBuffer termination, lazy native streams, and some error details remain approximate. The assertion import chain now has its UTF-16-aware Myers algorithm, comma-tolerant line matching, diff list output, color rendering, skipped-line collapse, AssertionError message/metadata formatting, and CallTracker call recording/report/reset/verify behavior mapped across its JS and Rust helper files. Deep equality still uses a local structural comparator instead of Bun's full native comparison semantics. The `internal/async_hooks_tick.ts` dependency is mapped through all 20 lines, including snapshotted `TickObject` init-hook delivery. `internal/fixed_queue.ts` is mapped through all 130 lines, including segmented circular buffers and the `FixedQueue` export. `internal/fs/binding.ts` is mapped through all 123 lines onto the shared fs binding adapter. The public Rust binding dispatch and test-only host helpers are mapped through line 527. The `fs.rm` parser validates own options before snapshotting, `stringToFlags` shares the filesystem flag parser, and the memfd test helper follows Bun's Linux flag and fallback path. `node_fs.rs` imports, error/path/time helpers, the 52-bit Blob size bound, IOVec dispatch, and async operation aliases map through line 637. The Windows-only `UVFSRequest` operation setup, buffer slicing, request callbacks/cleanup, completion settlement, `FsArgument`/`FsReturn` conversions, and generic `AsyncFSTask` lifecycle map through line 1459 onto the shared Go filesystem worker and event loop. The Go adapter rechecks the three Bun AbortSignal-bearing file operations (`readFile`, `writeFile`, `appendFile`) at completion; Windows libuv ref accounting, context tags, and native callback/task ownership remain host-boundary gaps. Bun's `NewAsyncCpTask`/`CpSingleTask` async traversal and promise settlement map through line 2254 onto `fs_copy.go`; Go now validates cp options and paths before dispatch, handles source/destination identity and directory nesting, merges directories, follows `dereference`, resolves or preserves symlink targets, applies source modes/timestamps, and mirrors `force`/`errorOnExist` behavior. Go copies tree entries sequentially instead of Bun's per-file fanout, and native clonefile/libuv/Windows copy handling remains approximate. Go routes descriptor I/O through `os.File` and `unix.Readv`/`Writev`; it caps readv at 1,024 buffers and batches writev, while the exact platform IOV_MAX and descriptor ownership remain differences. `AsyncReaddirRecursiveTask`, its result batching/type conversions, and promise/callback completion map through line 2736 onto `tsFSReadDirectory`; Go now follows recursive directory symlinks, skips vanished/not-enterable descendants, honors Buffer encoding precedence, and reflects Bun's async recursive encoding path. Bun fans directory scans out across worker tasks and joins lock-free result batches; Go scans sequentially with `os.File.ReadDir` to preserve filesystem order, but uses path based traversal, so fd-root safety, race behavior, and cancellation-on-first-error remain approximate. The argument structs map through line 4554, including null/default behavior for optional access/copy flags, u32 mode parsing, optional `exists` paths, and pre-scheduling `read`/`write` validation. Go handles DataView byte windows, read length coercion and zero-length behavior, buffer offset/length bounds, and write encoding validation/caps. Return conversions and handlers for `readlink`, realpath, rename, rmdir/rm, statfs/stat, symlink, truncate, unlink, watch, and time updates map through line 8109. Sync `cp` traversal, symlink handling, single-file copy behavior, and generated function dispatch/Dirent conversion map through line 9450. Recursive delete/error tables, `zig_delete_tree` and its bounded-stack fallback, and the function enum map through line 10116. Go now preserves the distinct syscall tags for `unlink`, `rm`, and `rmdir`, translates recursive-rmdir EACCES to EPERM, and applies Bun's narrow nonrecursive-rm errno table. The Go filesystem uses `os.RemoveAll` instead of Bun's fd-relative delete walker; descriptor-root race resistance and exact retry races remain approximate, along with platform-specific clone/copy paths, standalone graph lookup, early `readFile` OOM limits, large-write preallocation, and native Windows libuv ownership. The 821-line `dir_iterator.rs` platform implementations and wrapped iterator initialization map onto Go `os.File.ReadDir`; raw getdents records, Windows NT directory queries, borrowed filename buffers, and Windows UTF-16 names remain standard-library boundary differences. `Stat.rs` maps both Stats variants, Date-backed timestamp access, high inode conversion, and its test-only constructor through line 212; Go emits the same fields and lazy Date values, with statx supplying Linux birthtime. `StatFS.rs` and `node_fs_constant.rs` map through their final lines. The 465-line `NodeFSStatFSBinding.cpp` object layout, constructors, number/BigInt fields, and prototype tags are mapped onto the Go StatFs adapter; the Go result now includes `frsize` in Bun’s field order and falls back to `bsize` on hosts without a separate fragment size. `NodeFSStatBinding.cpp` (967 lines) and `NodeDirent.cpp` (375 lines) now map through their final lines. Go covers Stats/BigIntStats fields, lazy Date getters, predicates, BigInt timestamp fields, Dirent type predicates and constructors, and native `toStringTag` values; Stats and StatFs constructor calls are represented. `Dirent` type discovery for `os.ReadDir` uses Go metadata fallback when the host does not report a type, so unknown entries can be classified more precisely than Bun’s raw `DT_UNKNOWN` path. Continue with `node_fs_stat_watcher.rs`.

The `node:fs` watcher source stack is mapped through the end of each dependency: `node_fs_watcher.rs` 1256/1256, `node_fs_stat_watcher.rs` 1110/1110, `path_watcher.rs` 1618/1618, `win_watcher.rs` 547/547, `fs_events.rs` 939/939, and `node.classes.ts` 384/384. The Go adapter now returns a Node EventEmitter for `fs.watch`, handles the wrapper's listener, ignore, encoding shorthand, `throwIfNoEntry`, abort, error, close, ref/unref, and `_handle` behavior, and recursively adds Linux inotify watches with dynamic directory discovery and overflow `(change, null)` events. Go's standard library has no portable filesystem notification API, so Linux uses inotify through the Unix syscall adapter. Windows libuv, macOS FSEvents, FreeBSD kqueue, Bun's process-wide watcher deduplication and native cross-event batching remain host/backend gaps. `watchFile` still stores one watcher per path per loop rather than Bun's per-realm/per-ModuleGraph ownership. `ProcessBindingFs.cpp` (216 lines) and `ProcessBindingFs.h` (36 lines) also map through their final lines: Bun's raw process binding table deliberately exposes throwing stubs for filesystem calls plus four typed-array metadata views; this Go runtime currently routes `process.binding("fs")` to the functional filesystem adapter instead. `src/runtime/node/types.rs` is mapped through all 1,779 lines: argument ownership/flavor, string and typed-array parsing, encoding aliases and conversion, file URL and byte-path validation, uint32 mode and open-flag parsing, vector I/O buffers, Dirent conversion, path-or-fd, and path-or-Blob helpers are traced to the Go filesystem, Buffer, URL, and Blob adapters. Go GC replaces Bun's explicit JSC root/pin lifetime; native Windows path-prefix normalization and deferred async path-too-long errors remain differences. `src/runtime/node/util/validators.rs` is mapped through all 486 lines: integer/int32/uint32/number/string/boolean/function/object/array checks, array element validation, and string enum validation route through the local validators, filesystem argument checks, and `node:util.parseArgs` implementation. Parsing bounds and supported object flags are represented; formatted error details and JSC-specific type naming remain approximate. The `node:fs` dependency stack is complete. `node:crypto` is mapped through its 562-line wrapper, 403-line C++ binding, and 1,750-line Rust binding. The Go adapter uses standard crypto packages for hashes, HMAC, ciphers, key parsing/generation, signatures, ECDH, PBKDF2, HKDF, random bytes, and scrypt; scrypt uses a bounded RFC 7914 implementation, and Argon2d/i/id use a local BSD-licensed Argon2 fork with worker callbacks. Argon2 currently caps memory at 256 MiB and output at 64 KiB to avoid oversized allocations. SPKAC now handles base64/PEM input and common RSA, ECDSA, and Ed25519 signatures; uncommon signature algorithms, classic Diffie-Hellman, OpenSSL cipher/hash enumeration, custom engines, secure heap, and subtle WebCrypto remain unavailable or partial. OpenSSL error details and exact curve/cipher tables differ. Bun's `node_cluster_binding.rs` is mapped through all 810 lines: IPC setup, internal-message dispatch and acknowledgments, callback tracking, channel ref control, raw TCP/UDP/Unix socket binding, descriptor validation and closure are represented by the cluster and network adapters. The Go shared-handle adapter now validates socket descriptors, honors address family and IPv6/reuse flags where Go's socket controls expose them, reports the actual ephemeral port, and closes owned descriptors. Real child-process IPC/descriptor transfer, Windows socket export, and exact libuv error translation remain gaps. The 4,135-line `src/jsc/bindings/sqlite/NodeSqlite.cpp` is mapped through its final line.

The `node:net` pass includes option flags, connection entry points, connect trace and performance entries, TLS key-log file output and client `keylog` events, socket-write error completion and string encoding, TLS upgrades over an existing socket, transport close collapse, socket AbortSignal cleanup, the peer-FIN `EPIPE` write path, default server-name verification, protocol-error conversion, certificate-chain and identity checks, TLS authorization/event state, standalone wrapping, pre-handshake EOF `ECONNRESET` reporting, established-session TLS protocol-error conversion, data listener exception routing, read backpressure, client `connect`/`ready` event ordering, accepted-socket `pauseOnConnect`, timeout setup and pending-write suppression, read-error teardown, exact/wildcard server contexts, asynchronous SNI callbacks, ALPN selection callbacks, TLS server handshake deadlines and per-socket `tlsClientError` routing, requested client-certificate verification, accepted-socket server properties and connection-error routing, TLS client session-cache `session` events, Socket option validation, kernel-backed IP_TOS/IPV6_TCLASS controls, `onread` typed-array callbacks with sliced delivery, pause/resume tail preservation and deferred EOF, Socket address and state accessors, write-finalization waiting for connect, and handle close/reset routing through Go `crypto/tls` and `crypto/x509`. Source lines 2800–3702 additionally map descriptor-backed synchronous write helpers, reset/keepalive/no-delay/timeout controls, writev and pending-write handling, connect argument validation and DNS lookup, address alternation and deduplication, connection-attempt/BlockList/trace/performance events, and the `afterConnect` success/failure path. Go and Bun both try one family-selected address at a time, advancing immediately on failure or after the per-attempt deadline; the final address has no attempt deadline. The Go path preserves this flow, while exact libuv cancellation and error ordering remain approximations. Server construction, listen overload normalization, Unix socket permissions, fd listener adoption, asynchronous listen errors, close draining, async disposal, and EventEmitter rejection handling are mapped. The `listenInCluster`, faux-listen, bind-error-check, and round-robin connection paths now route through the local cluster adapter. Same-runtime round-robin handoff can pass an accepted Socket directly. Cross-process descriptor adoption and IPC handle transfer remain incomplete; Go does not expose the requested kernel listen backlog. TCP `reusePort` and `ipv6Only` socket options are applied through Go `ListenConfig`; private native-handle callbacks remain approximate. For lines 448–596, Go interrupts the raw read pump, replays bytes already read, and wraps that same `net.Conn` with `crypto/tls`; this folds Bun's separate raw/TLS handle pair and event callbacks into one Go connection lifecycle. Lines 597–604 use the existing abort-listener helper and remove it when the socket closes; lines 605–632 install the socket-specific write error after remote EOF. Lines 807–1088 map Bun's callback table and read-stop/end/error flow onto the Go read pump and local Duplex lifecycle; a read failure reaches `error` listeners once when present and otherwise ends the readable side like EOF. The TLS listener accepts without blocking the next accept, performs handshakes as separate tasks, reports failed handshakes through `tlsClientError`, and keeps accepting connections. Accepted sockets receive `server`, `_server`, and `isServer`; throwing connection/TLS event listeners use the `uncaughtException` path. Go collapses the raw/TLS handle pair into one `net.Conn`; server-side TLS ticket `session` events, descriptor-backed network reads and complete fd adoption/close semantics, exact kernel read sizing for `onread`, precise close-reset noise handling, and exact loop-reference accounting remain unimplemented or approximate. The TLS client cache is local to a socket, so session tickets are exposed as buffers but are not interoperable with Node/OpenSSL session state. SNI and ALPN callbacks execute on the owning event loop; the Go server dispatches its config callback on each handshake rather than selecting a reduced handler table, and ALPN callback `this` is the server object rather than the accepted socket. IPv4/IPv6 family autoselection follows Bun’s sequential retry shape; exact libuv timeout/error formatting and per-attempt TLS option rebuilding, cross-process cluster descriptor sharing, kernel backlog control and some connect/listen option validation remain gaps; `reusePort`/`ipv6Only` use Go socket options. Some callback error/uncaught-exception ordering remains approximate. The `internal/tls.ts` is mapped through all 476 lines. The TLS native dependency pass now traces `SecureContext.rs` (460/460), `runtime/socket/SSLConfig.rs` (531/531), `http/ssl_config.rs` (699/699), and the 116-line `SSLConfig.bindv2.ts` schema. The Go adapter now provides private secure contexts, `_external`, `addCACert`, certificate/key/CA file reads, encrypted legacy PEM and PBES2 PKCS#8 handling, a `SecureContext.parsePkcs12` adapter backed by `x/crypto/pkcs12`, selected TLS version flags, common ECDH curves, and CRL parsing/revoked-chain checks. Shared BoringSSL context interning, DH parameters, signature algorithm selection, TLS session timeout, and exact OpenSSL cipher expressions remain unavailable or approximate. Standalone client wrapping skips hostname verification, preserves `authorized=false`, emits only `secure`, and rejects an untrusted chain when requested. `node:net` source dependency coverage is complete; continue with the remaining TLS and network behavior gaps before selecting the next module.

The `node:zlib` source pass now covers the 1,021-line wrapper and all four native files (3,144 Rust lines). Deflate/gzip use Go standard codecs, Brotli uses `github.com/andybalholm/brotli`, and Zstandard uses the dependency already declared in go.mod. Brotli dictionaries, most encoder/decoder parameters, native worker scheduling, exact zlib window/memory/strategy controls, and incremental decompression remain tracked behavior gaps. `src/js/internal/worker/stdio.ts` is mapped through line 152 onto the local MessagePort stream adapter. `src/js/internal/worker/messaging.ts` is mapped through line 301 with local MessagePort routing, listener results, and timeout handling. `src/js/internal/trace_events.ts` is mapped through all 874 lines; category enablement, output, instrumentation, inspector collection, and metadata use the shared Go adapter, with separate worker VMs and native hooks remaining gaps. Its event timestamps now use the runtime's monotonic performance clock. The `node:worker_threads` wrapper is now mapped through line 1296; its Go adapter exposes the API but worker execution still lacks a separate JavaScript VM.

Each source index row starts at `ported through line 0`. Update the row after each file is translated. The module-level counter advances only when a module is fully translated; keep the per-file line checkpoint even when work pauses mid-module.

Dynamic source execution boundary: global `eval`, `Function`/generator constructors, `node:vm` Script and SourceTextModule compilation/evaluation, `node:module` CommonJS `_compile`/`runMain`/loader registration/preloading, `node:test.run` file dispatch, `node:repl`, `Worker`, `child_process.fork`, and inspector CDP commands that evaluate/compile source (`Runtime.evaluate`, `Runtime.callFunctionOn`, `Runtime.compileScript`, `Runtime.runScript`, `Debugger.evaluateOnCallFrame`, `Debugger.setScriptSource`, and `Page.addScriptToEvaluateOnNewDocument`) now throw `NotImplementedError`. Static Go-compiled modules, JSON module parsing, and regular external process spawning retain their existing routes.

Current stack checkpoint: `node:tls` uses Go’s standard `crypto/tls` without OpenSSL or cgo; its `validateCiphers` rules now reject invalid-only lists and `@SECLEVEL`, and supported literal OpenSSL names map to Go cipher IDs, while OpenSSL selector evaluation remains approximate. `node:_http_server` is mapped through line 4306, including response completion, header accessors, writes, and queued pipelined operations. `node:http` and `node:https` are registered through their full wrappers; HTTP parser, agent, native binding, and platform details remain prototype approximations. `internal/dgram.ts` and `node:dgram.ts` are mapped through their full source ranges onto Go `net.UDPConn` and `syscall` adapters; cluster sharing, the native send queue, custom async lookup, and source-specific multicast membership remain approximations. `node:dns.ts` and its Promise facade are mapped through their full source ranges onto Go `net.Resolver`; uncommon record types remain limited by the standard library. `node:fs.ts` is mapped through all 1,509 lines, and `node:fs/promises.ts` through all 1,674 lines. The `node:readline` and `node:readline/promises` wrappers are mapped through their final lines. Their internal helpers cover the interface, history, key decoding, cursor callbacks, and promise cursor builder. `node:tty` maps all 204 wrapper lines, the 533-line native binding, its 18-line header, the 26-line mode-state header, and the TTY helper section of `wtf-bindings.cpp` through line 223. `x/term` now supplies actual terminal detection, dimensions, per-stream raw-mode snapshots, and restoration; process-exit restoration, Bun’s global stdio reset, exact TCSADRAIN behavior, and Windows shared-handler/libuv lifecycle remain host differences. `node:repl` and its internal Acorn, await, completion, mode, and utility dependencies are source-mapped; the REPL entry point throws `NotImplementedError` because loading and evaluating its JavaScript implementation requires dynamic source compilation. The VM-based Acorn and REPL input paths are explicitly unsupported. `node:wasi` is mapped through all 1,775 lines and exposes the complete 47-entry preview1 import table; core memory, clock, random, descriptor, and path operations are adapted to Go, while socket support and several uncommon semantics are approximated. `node:v8` uses a local object-graph serialization format rather than V8 wire data, and `node:inspector` exposes Session, runtime console events, the inspector WebSocket endpoint, profiler/coverage commands, and its Promise facade; source-evaluation CDP methods return explicit unsupported errors; engine-backed debugger state, samples, and coverage remain approximations. `node:http2` maps its complete 6,794-line JavaScript wrapper and private TLS-upgrade helper to Go session, stream, request, response, server, settings, and validation surfaces. The generated `h2.classes.ts` host-method registration is mapped through line 151. The parser is reviewed through EOF with its prefix mapped sequentially through line 5008; this includes frame/settings helpers, DATA batch writes and queue coalescing, session-state fields, GC-backed ownership, AbortSignal listeners, stream-ID snapshot traversal, pending-write cleanup, control-frame output, event-loop-deferred HPACK errors, Writable backpressure/drain handling, per-session cork-buffer batching, frame-complete reentrant custom-socket writes, received stream/header/data/end/reset callbacks, settings parsing, local-window resizing, counters, GOAWAY/PING, ORIGIN/ALTSVC, and end-after-headers state. Selected protocol/session paths and the complete 210-line native header materializer are represented in Go. Inbound and outbound HPACK use `golang.org/x/net/http2/hpack` v0.48.0 in place of Bun’s C-backed `HpackHandle` (`h2_frame_parser.rs:1144,7731–7795`; `c-bindings.cpp:433–528`, `lshpack.h`). The Go buffer batches writes until a posted flush and clones bytes before dispatch. It does not reproduce Bun’s shared thread-local cork slot, native ref guards/scratch-buffer reuse, zero-copy TCP vectors, partial native write recovery, fixed TLS record boundaries, or fatal errno classification. Multi-stream flush still sweeps sorted stream IDs rather than using Bun’s resumable round-robin iterator. Lines 3346–4083 map native writes, posted auto-flush, session sink callbacks, reentrant input, settings/events, and session memory checks onto Go sockets and the event loop. The native socket syscall bridge and Bun’s JSC allocation-failure latch remain host gaps; multi-stream flush still sweeps sorted IDs. The production HTTP/2 parser surface is complete. `node:quic` maps its full public wrapper and internal symbol, diagnostics, headers, binding, state, statistics, class registration, endpoint, session, stream, TLS, and FFI sources onto a `quic-go` transport with Go TLS, datagrams, stream bodies, and basic HTTP/3 framing/control handling. The TLS pass traced Bun’s BoringSSL keylog and early-data callbacks, ALPN choice, certificate-chain loading, CA/CRL trust-store setup, hostname and peer verification, TLS 1.3 ciphersuite policy, and named-group selection. The Go backend preserves endpoint/session/stream method and state shapes plus graceful close and transport/application close types; remaining gaps are keylog callback delivery, 0-RTT setup, CRL enforcement, BoringSSL-specific verification details, ciphersuite policy selection, packet-level filtering and early-session callback timing, datagram ACK/loss status, in-flight/retransmit statistics, Retry/token controls, migration/path validation, dynamic QPACK, exact native error/status mapping, and endpoint UDP controls. `node:sqlite` now exposes DatabaseSync, StatementSync, Session, SQLTagStore, backup, constants, statement properties, and symbol disposal over the bundled SQLite C engine, including persistent files and full SQL execution through the local cgo driver. `node:test` source review now reaches the end of its 2,874-line wrapper and its mock-timer/assert utility dependencies. The Go adapter covers nested tests, hooks, contexts, mocks, assertions, canonical inherited tags, expected-failure matching, hook timeout/signal races, polling `waitFor`, run-option validation, and unsupported-mode errors. Test bodies now race the configured timeout/signal across callback settlement, registered child tests, and assertion-plan waits, with the timeout armed before invocation; hooks and callback-style `done` follow their two-argument and returned-Promise rules. Exact abort-listener lifecycle and late asynchronous duplicate-callback reporting remain runtime differences. `run({ files })` and default file discovery throw `NotImplementedError`: the statically compiled Go program cannot dispatch arbitrary source files into another runtime. Snapshot methods deliberately mirror Bun’s NotImplementedError behavior. `node:stream/iter` now maps its wrapper, ring buffer, push, duplex, pull, consumer, broadcast, share, and classic-stream adapters; backpressure and async scheduling remain approximate. `node:zlib/iter` exposes all 16 async/sync transform factories through the local transform adapter; gzip/deflate use Go standard codecs, Brotli uses `andybalholm/brotli`, and Zstd uses the existing klauspost codec. Brotli dictionaries and most native parameters remain unsupported. The Bun C++ wrapper is mapped into the Go adapter. The SQLite C engine is replaced wholesale by the bundled engine in the local MIT driver fork; the engine source skip is recorded in the source index, and `cgo` plus a C compiler are required. The public module shim is mapped through 38/38 lines; the `NodeSqlite.cpp` and `NodeSqlite.h` passes are complete, and `JSSQLStatement.cpp/.h` have been mapped as Bun `bun:sqlite` infrastructure with only its shared engine setup relevant here. The first 3,500 lines of `sqlite3_local.h` are mapped as engine API context; lines 3,501–14,350 and all of `sqlite3.c` are replaced by the local CGo SQLite engine. The SQLite binding source skip ranges are explicit in the source index. `node:fs/promises.ts` maps `FileHandle.readLines`, `pull`, `pullSync`, and `writer`; `internal/streams/iter/pull.ts` is mapped through all 977 lines and supplies the pull pipeline adapter. The fs dependencies `internal/fs/glob`, `cp-sync`, `cp`, `watch`, `watchfile`, `streams`, and `internal/streams/fast-utf8-stream` route through Go filesystem helpers. Glob uses the local simpler matcher in place of vendored minimatch; cp uses the shared Go copy walker instead of Node’s separate sync/async tree algorithms; watch helpers use the existing watcher state; filesystem streams use the local stream state machine. The UTF-8 stream maps buffering, flush, reopen, and lifecycle APIs, while Bun’s custom fs overrides, precise retry behavior, and zero-copy writer path are approximated. Through values remain per-file source checkpoints.

`node:vm` is mapped through the 674-line wrapper and native source chain. Context creation and static SyntheticModule callbacks remain exposed. `Script`, `runIn*`, `compileFunction`, SourceTextModule construction/evaluation, and dynamic source evaluation throw `NotImplementedError` because the Go program cannot compile or execute JavaScript source strings; the earlier literal-only evaluator has been removed from all reachable execution paths.

## Line-by-line dependency-stack strategy

1. Start at line 1 of a selected module entry point and translate its code in source order. Map a continuous source range to one or more named Go helpers rather than copying syntax mechanically.
2. If the next source range needs a Node module that is not available, stop before that range. Record the blocked file, next source line, and missing specifier in the active stack.
3. Push that dependency as the new active module and repeat from its first implementation line. When its APIs are available, mark its source files and pop back to the exact saved line in the caller.
4. Repeat until the original entry point reaches its end. Record each module’s dependencies and per-file `Through` line so no file silently jumps ahead.
5. Keep one Go module (`go.mod`) and the existing `goemit` package. Give each Node builtin its own file, named `internal/goemit/node_<builtin>_runtime.go`, containing that builtin's generated-runtime fragment. Go files in this package can share the private runtime value types and helpers without Go imports or package cycles. A build-tagged registry joins the fragments into the generated single-file runtime and dispatches module names to their constructors.

The completed `node:querystring`, `node:punycode`, `node:console`, `node:os`, `node:timers`, `node:timers/promises`, `node:util/types`, `node:module`, `node:perf_hooks`, `node:events`, `node:domain`, `node:diagnostics_channel`, `node:trace_events`, `node:assert`, `node:assert/strict`, `node:async_hooks`, `node:crypto`, `node:stream`, `node:stream/promises`, `node:stream/consumers`, and `node:stream/web` wrappers are registered in the generated runtime. The `node:module` adapter follows Bun's 76-entry builtin table and exposes Module construction, CommonJS wrapper helpers, lookup entry points, and compile-cache APIs; host resolver, source-map, and compile-cache behavior use local approximations. `node:async_hooks` implements `AsyncLocalStorage` and `AsyncResource` against persistent event-loop context frames; promise reactions and timer task completions capture those frames; its engine context-tracking switch is represented; `createHook` delivers `TickObject` init events for `process.nextTick`, while other hook phases and non-tick resource events remain unsupported. The assertion facade provides the standard assertion method family, strict aliases, error matching, and promise assertions through local runtime helpers; `partialDeepStrictEqual` handles recursive subset matching for arrays, collections, typed arrays, errors, and cycles; ordinary deep comparison and deprecated `CallTracker` reporting remain approximations. The timers implementation shares the transpiler event loop; its promise interval uses the local async-iterator machinery. Timeout and interval cancellation stop pending work through the loop cancellation channel; timer handle refresh/ref APIs remain partial. `node:util/types` maps predicates to generated runtime tags, including async functions and Node crypto KeyObjects; WeakMap, WeakSet, SharedArrayBuffer, WebCrypto CryptoKey, and N-API External values are not represented. `node:perf_hooks` uses Go's monotonic clock and exposes performance entries, observers, timerify, recordable histogram methods, and periodic event-loop delay sampling; a quantized bucket map approximates HDR histogram precision, and timer sampling approximates Bun's native hook. `node:events` implements EventEmitter listeners, emission, static helpers, one-shot Promises, and async iteration; EventTarget bridges are adapted. `node:domain` builds on `node:events`; it tracks enter/run/exit and the common listener wrappers, but does not update a process global. `node:diagnostics_channel` uses weakly cached channels with generation-safe cleanup, active/inactive prototypes, live subscriber mutation, Map-ordered store composition, process-nextTick error reporting, async-context wrapping, and trace wrappers that normalize thenables; Go finalization timing and exact JavaScript engine behavior remain different. `node:trace_events` tracks category enablement and writes local trace files through an internal agent adapter; worker-specific native hooks and exact trace timing remain approximate. `node:os` uses Go's standard library and procfs in place of Bun's Rust host binding; OS constants and some non-Linux host details remain approximate. This pass mapped `node_os.bind.ts` through 94/94 and `node_os.rs` through 1614/1614, correcting Linux iowait/IRQ CPU accounting, physical free-memory and system-uptime semantics, IPv6-only `scopeid`, interface filtering/MAC width, and priority validation. Host-specific CPU/network details, non-Linux physical memory/uptime/release data, and platform constants remain approximate. Newly ported modules are enabled by default and can be omitted from a compiler build with `-tags no_node_modules`. The existing `fs`, `buffer`, and `path` runtime support predates this registry and is not controlled by that tag yet.

`node:process` source review is complete through `BunProcess.cpp:5172`; `ProcessBindingConstants.cpp/.h`, `NodeProcessModule.h`, and `node_process.rs` are also mapped through their final lines. No unported Node builtin dependency was encountered. The binding's TLS/OpenSSL constants are numeric compatibility values only; the Go runtime adds no OpenSSL or cgo dependency. The `node:path` wrappers, native bridge, Rust host wrapper, and portable Rust algorithm crate are mapped through their final lines; they route to the existing Go path runtime with standard-library host functions and local Windows path logic. `path.matchesGlob` uses the local glob matcher, so advanced Bun glob features remain approximate. The `node:buffer` module wrapper, Rust helper, JSBuffer bindings, encoding parser, constructor, and prototype sources have all been mapped. The Go runtime now exposes the constructor/statics, string and view operations, search and compare, typed accessors, numeric/BigInt readers and writers, encoding-specific methods, inspection, and module exports. Object URL lookup remains a placeholder; JSC memory mapping, resizable-buffer side effects, custom inspect context properties, and some exact Node error ordering/messages are approximated. The `node:string_decoder` native module handoff, declaration, and 631-line implementation are now translated into stateful UTF-8, UTF-16LE, and base64/base64url chunk decoding with the `write`, `end`, `text`, `lastChar`, `lastNeed`, `lastTotal`, and `encoding` API. Its encoders reuse the Buffer runtime; invalid receiver/error wording and a few malformed chunk boundary details are approximate. The `node:url` wrapper, internal IDNA helpers, NodeURL native bridge, and punycode/host-validation helpers are mapped through their final lines; the Go adapter uses `net/url`, `filepath`, and the existing punycode runtime, so ICU UTS #46 validation and detailed `URLPattern` matching remain approximate. Its URL, URLSearchParams, legacy `Url`, file URL, IDNA, and HTTP-options entry points are exposed. The crypto dependency stack has mapped `readable.ts` through line 1734, `writable.ts` through 1118, `duplex.ts` through 166, `duplexify.ts` through 395, `webstreams_adapters.ts` through 864, the TypeScript logic of `native-readable.ts` through 306, `iter/from.ts` through 576, `iter/classic.ts` through 898, `transform.ts` through 173, `lazy_transform.ts` through 53, `passthrough.ts` through 20, `pipeline.ts` through 448, and `compose.ts` through 221. These Go adapters cover the classic/Web stream conversion surfaces with approximate scheduling, destroy, and strategy behavior; Bun's `BunStreamConsumers.cpp` native-handle transfer remains a no-op.

The active stack has completed `internal/abort_listener.ts` (31 lines), `internal/async_hooks.ts` (28), `internal/errors.ts` (23), `internal/streams/end-of-stream.ts` (384), `internal/streams/add-abort-signal.ts` (58), `internal/streams/destroy.ts` (343), `internal/streams/state.ts` (52), `internal/streams/from.ts` (200), `internal/streams/iter/types.ts` (87), `internal/streams/iter/utils.ts` (270), `internal/streams/iter/from.ts` (576), and `internal/streams/iter/classic.ts` (898). Readable resumes after its classic iterator import on line 1286. Writable, Duplex, Duplexify, native-readable TypeScript, and Web stream adapters are mapped through their source ends; the C++ native-readable transfer is represented by a placeholder. Scheduling, Web stream strategies, native transfer, and stream error ordering remain approximate.

`node:crypto` maps the 562-line JavaScript wrapper and its C++/Rust binding sources into Go adapters. Most C++ and Rust implementation internals were replaced at the binding boundary; the Diffie-Hellman wrappers and shared methods were traced into the Go adapter. Go standard crypto covers the common operations, scrypt uses a bounded RFC 7914 implementation, and Argon2d/i/id use a local BSD-licensed Argon2 fork. The Argon2 adapter caps memory at 256 MiB and output at 64 KiB. Classic Diffie-Hellman now supports explicit prime/generator parameters, the `modp14` group, generated DH KeyObjects, and shared-secret calculation; other named groups, DH key encoding/import, safe-prime generation, and constant-time exponentiation remain gaps. FIPS/secure-heap controls and subtle WebCrypto remain explicit stubs or errors; SPKAC supports common RSA, ECDSA, and Ed25519 signature forms. The adapter imports no OpenSSL and uses no cgo. The `node:stream`, `stream/promises`, `stream/consumers`, and `stream/web` facades are registered. Readable operators and `duplexPair` are exposed with simplified concurrency and scheduling; `stream/web` forwards global constructors, and names missing from the runtime remain undefined. Experimental `stream/iter` and `zlib/iter` are registered with their wrappers and shared iterator helpers. Gzip/deflate transforms use Go standard codecs; Brotli uses `andybalholm/brotli`, and Zstd uses the existing `klauspost/compress` codec; stream scheduling and some native parameters remain approximate. Broadcast/Share and native async scheduling remain prototype adaptations.

## Port ledger

The stream dependency pass has mapped `internal/abort_listener.ts` through line 31, `internal/async_hooks.ts` through 28, `internal/errors.ts` through 23, `internal/streams/end-of-stream.ts` through 384, `internal/streams/add-abort-signal.ts` through 58, `internal/streams/destroy.ts` through 343, `internal/streams/state.ts` through 52, and `internal/streams/from.ts` through 200. `readable.ts` is at line 29 after its imports were resolved. The stream adapters expose EOS and abort handling, iterator-to-readable construction, shared state flags, and destroy lifecycle methods. The performance and error helpers use local runtime state, so Node native stream and host behavior are approximate.

`Through` tracks translated Bun source lines, not verified behavioral parity. The current runtime fragments are prototype implementations; the console fragment currently supplies common logging, count, and timing methods, not the full Node `Console` surface.

| Module | Bun source | Bun lines | Through | Status | Dependencies |
|---|---|---:|---:|---|---|
| `node:querystring` | `src/js/node/querystring.ts` | 523 | 523 | source translated | none |
| `node:punycode` | `src/js/node/punycode.ts` | 414 | 414 | source translated | none |
| `node:console` | `src/js/node/console.ts` | 1 | 1 | default export and `globalThis.console` share one object; standard logging, count, timer, assertion, and group methods are present; Node's `Console` constructor and custom writable streams remain unavailable | global object, runtime writers |
| `node:os` | `src/js/node/os.ts` | 179 | 179 | wrapper translated; native bridge replaced | Go standard library |
| `node:timers` | `src/js/node/timers.ts` | 41 | 41 | wrapper translated; timeout and interval cancellation stop pending work through the loop cancellation channel; `refresh`/handle APIs remain partial | `node:timers/promises` (lazy property) |
| `node:timers/promises` | `src/js/node/timers.promises.ts` | 253 | 253 | wrapper translated; async iterator adapted | `node:timers` (global timer primitives) |
| `node:util/types` | `src/jsc/modules/NodeUtilTypesModule.cpp`, `.h` | 1314 | 1314 | native predicate and comparison helpers mapped; weak collection/shared buffer/crypto native brands unavailable | generated runtime tags, `node:assert`, `node:crypto` |
| `node:perf_hooks` | `src/js/node/perf_hooks.ts`; `src/js/internal/perf_hooks/monitorEventLoopDelay.ts`; `JSNodePerformanceHooksHistogram.cpp/.h`; `JSNodePerformanceHooksHistogramConstructor.cpp/.h`; `JSNodePerformanceHooksHistogramPrototype.cpp/.h` | 1627 | 381/381; 71/71; 274/274; 165/165; 112/112; 49/49; 529/529; 46/46 | histogram record/add/reset/percentile/statistics APIs and event-loop delay sampling adapted; Go uses a bounded quantized bucket map in place of HDR's exact bucket layout, and timer sampling approximates Bun's native hook | Go monotonic clock, event-loop timers, generated runtime |
| `node:events` | `src/js/node/events.ts`; `JSEventTarget.cpp/.h`; `JSEventTargetNode.cpp/.h`; `JSEventTargetCustom.cpp/.h`; `EventTarget.cpp/.h`; `JSEvent.cpp/.h`; `Event.cpp/.h` | 3393 | events.ts:1057/1057; native:2336/2336 | EventEmitter and static helpers adapted; EventTarget listener identity/options, capture ordering, Event attributes, cancellation, dispatch, AbortSignal integration, and listener queries mapped. Single-target dispatch has no DOM ancestor path or WebCore exception delivery; JavaScriptCore GC rooting, event scopes, and wrapper conversion are replaced with Go values | `node:util/types`, local AbortSignal, Go EventTarget runtime |
| `node:child_process` | `src/js/node/child_process.ts` | 2071 | 2071/2071 | Go `os/exec` adapter maps spawn, exec, execFile, sync calls, fork option validation (fork execution throws `NotImplementedError` because it launches a JavaScript module in another runtime), the `ChildProcess` class, kill/dispose, timeout/AbortSignal handling, custom promisify symbols and `.child`, shell argv normalization, stdio/environment helpers, validators, URL/path conversion, `SystemError`, and exports. Real IPC transport, descriptor-backed extra stdio, detached process groups, uid/gid application, exact Windows argv quoting, maxBuffer termination, lazy native streams, and some error details remain approximate | Go standard `os/exec`, local streams and events, `node:os`, `node:util`, `internal/promisify`, `internal/fs/streams`, `node:url`, `node:net` |
| `node:domain` | `src/js/node/domain.ts` | 96 | 96 | wrapper translated | `node:events` |
| `node:cluster` | `src/js/node/cluster.ts`; `src/js/internal/cluster/{primary,child,Worker,isPrimary,RoundRobinHandle,SharedHandle}.ts`; `src/js/internal/linkedlist.ts`; `src/runtime/node/node_cluster_binding.rs` | 1993 | 1183/1183; 810/810 | primary and worker APIs mapped; Go `net` backs shared TCP/UDP/Unix listeners and socket metadata; real child-process IPC/descriptor transfer, Windows socket export, and exact libuv errors remain approximate | `node:events`, Go `net`, `os/exec`, local cluster handles |
| `node:net` | `src/js/node/net.ts`; `src/runtime/node/node_net_binding.rs`; `src/runtime/node/net/BlockList.rs`; `src/runtime/socket/SocketAddress.rs` | 6426 | net.ts:4634/4634; node_net_binding.rs:184/184; BlockList.rs:524/524; SocketAddress.rs:1084/1084 | Go `net` adapter exposes Socket, Server, BlockList, SocketAddress, IP predicates, connection controls, process flags, connect trace events, TLS key-log file and client event output, existing-socket and standalone TLS wrapping, socket AbortSignal cleanup, peer-FIN `EPIPE` writes, certificate-chain/identity verification, pre-handshake EOF reset errors, established-session TLS error conversion, data listener exception reporting, read backpressure, client `connect`/`ready` ordering, accepted-socket pausing, timeout suppression for queued writes, read-error teardown, exact/wildcard TLS server contexts, event-loop-dispatched SNI/ALPN callbacks, handshake deadlines, per-socket TLS failure events, client session-cache events, SocketAddress validation/getters/JSON/common parse behavior, option checks, kernel TOS controls, custom `onread` buffer callbacks and tail/EOF handling, Socket state accessors, write-finalization waiting, optional client-certificate validation, accepted-socket server state, BlockList rules and matching, and BlockList reference serialization; server TLS session-ticket events, descriptor-backed network reads and complete fd adoption/close semantics, exact kernel read sizing for `onread`, precise close-reset noise handling, exact loop-reference accounting, Bun's separate adopted TLS handles, custom socket classes, cross-process cluster descriptor adoption/IPC and kernel backlog control, exact libuv IPv4/IPv6 retry cancellation/error formatting and per-attempt TLS option rebuilding, exact WHATWG URL parsing, named-interface IPv6 scope ids, low-level SocketAddress DTO helpers, clone-registry cleanup, and worker structured cloning remain approximate or absent | Go standard `net`, `net/url`, `crypto/tls`, `crypto/x509`, `os`, `syscall`; local Duplex, events, DNS, validators |
| `node:tls` | `src/js/node/tls.ts`; `src/jsc/bindings/NodeTLS.cpp`; `src/runtime/api/bun/SecureContext.rs`; `src/runtime/socket/SSLConfig.rs`; `src/http/ssl_config.rs`; `src/runtime/socket/SSLConfig.bindv2.ts` | 3925 | tls.ts:1757/1757; NodeTLS.cpp:362/362; SecureContext.rs:460/460; SSLConfig.rs:531/531; http/ssl_config.rs:699/699; SSLConfig.bindv2.ts:116/116 | exported client/server/context/socket surfaces, certificate enumeration and overrides (including the read-only `rootCertificates` property), private secure-context CA mutation, key/CA file loading, encrypted legacy PEM and PBES2 PKCS#8 keys, PKCS#12 option conversion/CA chains, common curves, and selected secureOptions mapped onto Go `crypto/tls`; Bun-style invalid cipher-list and `@SECLEVEL` validation plus supported literal cipher aliases are represented; hostname checks use Go `x509.VerifyHostname` with the Bun IDNA helper and CN fallback; ALPN arrays and views follow Bun’s encoded buffer conversion and 255-byte validation. Pre-handshake writes queue while connecting, and `TLSSocket._final` waits for the private handshake-complete event before delegating to the net finalizer. Ticket-key methods and `parseCertString` mirror Bun’s explicit unsupported errors; renegotiation errors are scheduled through `nextTick`. `setKeyCert` reports unsupported because Go TLS cannot replace a per-connection certificate; finished-message and shared-signature introspection remain unavailable. Decoder limits from `x/crypto/pkcs12`, BoringSSL interning, DH params, sigalgs, session timeout, OpenSSL selector evaluation, PSK cipher negotiation, and exact cipher naming remain gaps | Go `crypto/tls`, `crypto/x509`, `crypto/pem`, `x/crypto/pkcs12`, `net`; local `node:net` sockets |
| `node:_tls_common` | `src/js/node/_tls_common.ts` | 43 | 43/43 | re-exports SecureContext APIs, translates legacy peer certificate fields, and emits `DEP0192` once | `node:tls`, `node:process`, JSON |
| `node:_tls_wrap` | `src/js/node/_tls_wrap.ts` | 14 | 14/14 | re-exports TLSSocket/server/connect APIs and emits `DEP0192` once | `node:tls`, `node:process` |
| `node:_http_agent` | `src/js/node/_http_agent.ts` | 626 | 626 | agent pools, request queues, socket reuse, and proxy option handling adapted; CONNECT tunnel behavior and host networking details remain approximate | `node:events`, `internal/http`, validators, `node:tls`, `node:net`, async context |
| `node:http` | `src/js/node/http.ts`; `src/jsc/bindings/NodeHTTP.cpp`; `src/jsc/bindings/NodeHTTP.h`; `NodeHTTPRequestHeaderValues.h`; `src/jsc/bindings/node/JSNodeHTTPServerSocket.{h,cpp}`; `JSNodeHTTPServerSocketPrototype.{h,cpp}`; `src/runtime/node/node_http_binding.rs`; `src/runtime/server/NodeHTTPResponse.rs` | 6,186 | http.ts:138/138; NodeHTTP.cpp:1023/1023; NodeHTTP.h:8/8; NodeHTTPRequestHeaderValues.h:48/48; socket:216/216+1141/1141; prototype:46/46+631/631; node_http_binding.rs:110/110; NodeHTTPResponse.rs:2825/2825 | public request/server facade and native HTTP dispatch/socket/response lifecycle mapped; fallback response length checks and the parser binding/ConnectionsList API are represented; incremental llhttp grammar, raw socket framing, TLS/uWS callbacks, GC pinning, and engine ownership remain approximate or host-bound | Go standard `net/http`, `net`, `time`; local `_http_*`, stream, and event adapters |
| `node:https` | `src/js/node/https.ts` | 636 | 636 | HTTPS request/server facade, TLS agent session keys, proxy CONNECT handshake, and global agent mapped over Go TLS/net adapters | `node:http`, `node:tls`, `node:net`, `internal/http`, URL and validation helpers |
| `node:fs` | `src/js/node/fs.ts`; `src/js/node/fs.promises.ts` | 3183 | fs.ts:1509/1509; fs.promises.ts:1674/1674 | callback, sync, Promise, directory, watcher, copy, stream, glob, FileHandle, readLines, and disposable APIs map onto shared Go filesystem helpers; timestamp conversion accepts JavaScript numeric strings and preserves Bun’s negative-number behavior; FileHandle pulls and writer APIs route through `internal/streams/iter/pull`; handle locking and write backpressure are approximated | Go filesystem operation, stream, watcher, directory, glob, copy, and handle adapters |
| `node:readline` | `src/js/node/readline.js` | 546 | 546 | legacy and promise-question interfaces, compatibility accessors, terminal callbacks, keypress events, and async line iteration adapted | `internal/readline/*`, history, events, validators, string decoder |
| `node:readline/promises` | `src/js/node/readline.promises.js`; `src/js/internal/readline/promises.js` | 198 | 60/60; 138/138 | Promise question interface and queued cursor-control `Readline` API mapped; terminal writes use local stream/event-loop hooks | `node:events`, internal readline interface and cursor utilities |
| `node:tty` | `src/js/node/tty.ts`; `ProcessBindingTTYWrap.cpp/.h`; `BunTTYState.h`; `wtf-bindings.cpp:1-223` | 1004 | tty.ts:204/204; ProcessBindingTTYWrap.cpp:533/533; ProcessBindingTTYWrap.h:18/18; BunTTYState.h:26/26; wtf-bindings.cpp:223/320 | ReadStream/WriteStream wrappers, ref methods, cursor controls, color depth, terminal detection and sizing, and per-stream raw-mode restore are mapped onto the local stream adapter and `x/term`; process-exit terminal reset, exact `TCSADRAIN`, and Bun’s Windows uv/control-handler lifecycle remain host gaps | `internal/fs/streams`, `internal/tty`, validators, `node:readline`, `golang.org/x/term` |
| `node:repl` | `src/js/node/repl.js` | 1543 | 1543 | the entry point throws `NotImplementedError` because loading its implementation and evaluating REPL input require dynamic JavaScript compilation; the internal Acorn parser and REPL evaluation paths are unsupported | `node:vm`, `node:readline`, `node:console`, `node:async_hooks`, `node:path`, `node:fs`, internal REPL helpers |
| `node:wasi` | `src/js/node/wasi.ts` | 1775 | 1775 | `WASI` constructor and all 47 preview1 imports exposed; common memory, clock, random, descriptor, and path calls mapped, with sockets and several semantics approximated | Go `os`, `filepath`, `time`, `crypto/rand`; `node:tty`, `node:fs`, `node:path` |
| `node:v8` | `src/js/node/v8.ts`; `NodeV8.cpp` | 558 | v8.ts:450/450; NodeV8.cpp:108/108 | heap counters and GCProfiler are mapped; promise hooks report init with parent promises, before/after reactions, and settlement, with removable `createHook` and single-phase hooks. `Serializer`/`Deserializer` support value framing, fixed-width/raw byte methods, transfer IDs, shared ArrayBuffer backing and offsets for views, RegExp, Error names/messages and own `cause`/`errors`; the local graph format accepts its prior version but is not V8 wire compatible. GCProfiler reports Go runtime GC cycles rather than JSC heap events and can restart after stop | Go runtime/JSON, `node:stream`, `node:os`, `node:worker_threads`, `internal/serialization_buffers` |
| `node:vm` | `src/js/node/vm.ts`; `NodeVM.cpp/.h`; `NodeVMScript.cpp/.h`; `NodeVMModule.cpp/.h`; `NodeVMSourceTextModule.cpp/.h`; `NodeVMSyntheticModule.cpp/.h`; VM fetcher, option, termination, and code-generation headers | 6002 | vm.ts:674/674; NodeVM.cpp:2056/2056; NodeVM.h:184/184; NodeVMScript.h:119/119; NodeVMScript.cpp:679/679; NodeVMModule.h:130/130; NodeVMModule.cpp:818/818; NodeVMSourceTextModule.h:70/70; NodeVMSourceTextModule.cpp:566/566; NodeVMSyntheticModule.h:60/60; NodeVMSyntheticModule.cpp:254/254; NodeVMScriptFetcher.h:104/104; NodeVMOptionNames.h:47/47; NodeVMRunTermination.h:62/62; NodeVMRunTermination.cpp:109/109; CodeGenerationFromStrings.h:70/70 | context operations and SyntheticModule callbacks remain available; `Script`, `runIn*`, `compileFunction`, `SourceTextModule` compilation, and dynamic source evaluation throw `NotImplementedError` because the static Go program cannot compile or execute JavaScript source strings. SourceTextModule parsing/linking, bytecode caches, real microtask queues, TLA, and SIGINT/time-bound execution are outside the runtime | Go runtime values, local Promise/event loop, `node:util`, `node:stream`, validators |
| `node:inspector` | `src/js/node/inspector.ts`; `JSInspectorProfiler.cpp` | 919 | 760/760; cpp:159/159 | Session and CDP command surface, console notifications, inspector WebSocket listener, lifecycle, and profiler/coverage entry points exposed. Source-evaluation/compilation CDP methods return explicit unsupported errors for the static Go runtime. `JSInspectorProfiler.cpp:1-58` relies on BunCPUProfiler/JSC VM APIs; `60-83` toggles JSC control-flow instrumentation; `85-159` walks JSC executable/source heaps and reads engine counters. Go keeps protocol-shaped CPU profile responses and returns empty precise coverage because its compiler has no matching engine hooks | Go `net`, WebSocket framing, `node:events`, `node:worker_threads`, `internal/trace_events` |
| `node:inspector/promises` | `src/js/node/inspector.promises.ts` | 28 | 28 | Promise-based Session facade exposed | `node:inspector` |
| `node:http2` | `src/js/node/http2.ts`; `src/js/node/_http2_upgrade.ts`; `src/runtime/api/bun/h2*`; HPACK FFI | 19918 | wrappers 7199/7199; wire 350/350; settings 120/120; flow control 153/153; HPACK wrapper 350/350; stream state 145/145; `mod.rs` 12/12; connection.rs 2462/2462 (production through 2117; remaining test fixtures reviewed as reference-only); h2.classes.ts 151/151; parser reviewed 7873/7873 | constants, settings, headers, sessions, streams, requests, responses, servers, and TLS upgrade surfaces exposed; Go maps protocol framing/control handlers and bridges the connection preface, explicit-only initial SETTINGS, retained frame input, session callbacks, flow-controlled DATA output, padding, request priority frames, and parser-backed settings limits to sockets. Parser lines 1–420 map frame/error constants, 31-bit reserved-bit serialization, frame headers, and priority wire encoding to the Go HTTP/2 runtime. Lines 421–655 map settings defaults/projection, pseudo-header and single-value tables, malformed-header checks, and outbound Latin-1-or-UTF-8 value encoding. Lines 662–839 map callback dispatch through the Go engine sink and event loop; Bun's internal parser-constructor callback validation is not a public Node API. Lines 841–897 use Go GC and socket writes in place of Bun's parser pool, thread-local corking, native refcounts, and vectored-write buffers. Lines 899–991 map DATA batching, empty DATA frames, and complete outbound frame-block writes; Go coalesces copied bytes into one socket write rather than borrowing TCP payloads for writev. Lines 993–1190 map parser/session state into Go session, connection, and stream state; JSC roots, intrusive refcounts, native keepalives, and raw stream borrows are replaced by Go values and GC. Lines 1191–1356 map parser release/ref guards to GC ownership and resumable stream iteration to ID snapshots with re-lookup; the stream state and flow-control fields map to the Go stream/session state. Lines 1358–1380 map per-stream AbortSignal listeners through the local abort-listener helper. Lines 1382–1687 map the pending DATA queue, frame padding, and flow-controlled flush path to the Go outbound queue and batched DATA encoder. Lines 1689–1829 map queue coalescing and callback replacement; the Go queue retains per-frame length while its backing slice advances on partial sends. Lines 1831–1968 map stream initialization, send/receive state checks, AbortSignal attachment, and pending-write cleanup; cleanup drains callbacks before the close event. Lines 1974–2315 reuse the Go HPACK encoder and explicit-settings serializer and map stream abort/reset, GOAWAY, ALTSVC, and PING output to session/connection helpers. Lines 2316–2477 map connection preface/settings, WINDOW_UPDATE, and event/callback dispatch to the Go socket sink and event loop. The Go header materializer applies Bun's Latin-1, raw-header, sensitive-header, duplicate, cookie, set-cookie, status, and trailer rules. The generated host-method registration is mapped, including stream read backpressure and trailer completion behavior. `golang.org/x/net/http2/hpack` v0.48.0 now supplies a stateful encoder/decoder with indexed fields, dynamic tables, Huffman strings, and compression-error validation; sensitive fields use a never-indexed literal path so they cannot be indexed by the library encoder; Bun’s native `HpackHandle` route (`h2_frame_parser.rs:1144,7731–7795`, `h2/hpack.rs:1–112`, `src/http/lshpack.rs:1–238`, C ABI `c-bindings.cpp:433–528`) is mapped as the source dependency this Go codec replaces. Lines 2479–2490 defer outbound compression failures through the event loop before reporting a session error. Lines 2492–2897 map socket output through cloned writes, a transport queue for `Writable.write` backpressure, drain replay, and flow-controlled DATA flushes. Lines 2899–2917 map stable payload ownership by copying bytes before socket dispatch. Lines 2919–3155 map cork/uncork and batch-buffer output to a per-session byte buffer flushed once by a posted event-loop callback. Lines 3157–3244 map the native vectored batch to one contiguous Go socket write; Go does not expose Bun’s zero-copy TCP vectors or partial native write recovery. Lines 3246–3344 map auto-cork, frame-complete custom-socket writes, and reentrant output to the posted flush buffer; each write receives a stable byte copy and reentrant output schedules after the in-flight batch. Go does not reproduce Bun’s shared thread-local cork slot, native ref guards/scratch-buffer reuse, fixed TLS record boundary, fatal errno classification, or explicit tx_tracker. Multi-stream flush still sweeps sorted stream IDs rather than using Bun’s resumable round-robin iterator. Bun’s native socket writer bridge and JSC callback-allocation failure/left-exception state remain host-boundary details. Lines 3346–4083 map native writes, posted auto-flush, parser sink bookkeeping, settings projection, received stream creation, pending-input/reentrant input handling, frame counters, error callbacks, PING/GOAWAY/WINDOW_UPDATE/ALTSVC/ORIGIN callbacks, and session memory checks to Go socket/event-loop/session state. Lines 4085–4280 map stream/header/data/end/reset callbacks; CANCEL alone emits `aborted`, while other resets report a deferred stream error. Lines 4283–4324 map direct/padded writes; lines 4330–4663 map explicit/custom settings and local-window growth; lines 4665–5008 map counters, GOAWAY/PING, ORIGIN/ALTSVC, and end-after-headers; lines 5010–5148 map stream state and RST_STREAM; lines 5154–5917 map memory/data/trailer/read-window behavior; lines 5920–6430 map writable data, IDs, PUSH_PROMISE, contexts, and stream iteration; lines 6433–7359 map abort/error fan-out, request options, header encoding/fragmentation, and stream state transitions; lines 7362–7873 map read/native transport entry points, constructor settings, context roots, and teardown. Go now defers stream reset/size errors, rejects invalid request priority options, orders request pseudo-headers, enforces outbound session-memory limits, and removes closed sessions/streams from singleton ownership maps. Remaining host gaps are Bun’s native socket syscall/refcount bridge, JSC allocation-failure exception latch, C HPACK implementation, exact writev/TLS behavior, and reentrant callback semantics. Raw header-pair duplicates are grouped by name, and Go does not retain the wrapper’s pending-request queue for connecting/concurrency-limited sessions. Rust test fixtures remain reference-only. | Go `net`, `crypto/tls`, `golang.org/x/net/http2/hpack`, sockets, `node:events`, private upgrade helper |
| `node:quic` | `src/js/node/quic.ts`; `src/js/internal/quic/{binding,diagnostics,http2util,quic,state,stats,symbols}.ts` | 7828 | 55/55; 2/2; 70/70; 197/197; 5565/5565; 886/886; 983/983; 70/70 | UDP listen/dial, endpoint reuse and idle close, TLS identities/ALPN/groups/CA validation, stream bodies/writer/backpressure/idle timeout, datagrams, H3 control/header/data frames, ORIGIN, GOAWAY, PRIORITY_UPDATE, graceful drain, block-list policy, native binding constants/state offsets, and lifecycle diagnostics use `quic-go`; basic byte/packet/loss/RTT/congestion stats come from `ConnectionStats`. Endpoint graceful close stops new accepts and waits for existing sessions; graceful close preserves transport/application close type. TLS 1.3 groups, CA roots, hostname verification, server certificate chains, and TLS callbacks/state shapes are mapped. Gaps: keylog callback delivery, 0-RTT setup, CRL enforcement, BoringSSL-specific verification, ciphersuite policy selection, packet-level block filtering and early-session callback timing, datagram ACK/loss status, in-flight/retransmit stats, Retry/token controls, migration/path validation, dynamic QPACK, exact native error/status mapping, and endpoint UDP controls | `github.com/quic-go/quic-go`, `github.com/quic-go/qpack`, Go `crypto/tls`, local streams, diagnostics channels |
| `node:sqlite` | `src/jsc/modules/NodeSqliteModule.h`; `src/jsc/bindings/sqlite/NodeSqlite.{h,cpp}`, `JSSQLStatement.{h,cpp}`, SQLite engine | 293,545 | NodeSqliteModule.h:38/38; NodeSqlite.cpp:4135/4135; NodeSqlite.h:902/902; JSSQLStatement.cpp/.h:3166/3166 (Bun `bun:sqlite` binding, shared engine setup only); sqlite3.c:0/269669 (replaced by bundled SQLite C engine in local MIT go-sqlite3 v1.14.24 fork); sqlite3_local.h:3500/14350; sqlite3_error_codes.h and lazy_sqlite3.h replaced by driver constants/initialization | DatabaseSync/StatementSync/Session/SQLTagStore and module exports use the bundled SQLite engine for file-backed/full SQL execution, bindings, transactions, cursors, limits, scalar and window functions, authorizers, sessions/changesets, serialization, backup, and tag-store LRU; cgo and a C compiler are required; Bun-only custom-engine selection and native error/ownership details remain gaps | local MIT `github.com/mattn/go-sqlite3` fork with bundled C engine; no system SQLite needed |
| `node:test` | `src/js/node/test.ts`; `src/js/internal/test_runner/mock_timers.ts`; `src/js/internal/assert/utils.ts` | 3,690 | 2874/2874; 781/781; 35/35 | Further porting is parked per the request to skip this dynamic test-runner module. Existing static test/suite/hook/context/mock/assert surfaces are mapped, including option defaults and validation, inherited canonical tags, expected-failure matching, assertion plans, hook timeout/signal races, polling `waitFor`, nested-test draining, test-body timeout/signal races covering child tests and assertion-plan waits, timeout arming before invocation, callback-style `done` settlement including returned-Promise conflict, nested `runOnly()` filtering for `only` tests, and snapshot serializer/path-resolver configuration validation and storage. `mock.module()` and file-run discovery/execution require a dynamic module loader; Bun's snapshot assertion functions are reference stubs. Exact abort-listener lifecycle and late asynchronous duplicate-callback reporting remain runtime differences | local `node:stream`, `node:assert`, `node:async_hooks`, timer scheduler |
| `node:dgram` | `src/js/node/dgram.ts` | 1546 | 1546 | broad UDP Socket API adapter mapped onto Go `net.UDPConn` and `syscall`; broadcast, TTL, loopback, interface, and common membership options are wired, while cluster sharing, native queue behavior, custom async lookup, and source-specific membership remain approximations | `internal/dgram` (146/146), validators, `net.BlockList`, events, DNS, diagnostics channel |
| `node:dns` | `src/js/node/dns.ts`; `src/js/node/dns.promises.ts` | 1073 | 1071/1071; 2/2 | callback and Promise resolver APIs, per-instance resolver state, server configuration, callback promisify hooks, and result ordering mapped over Go’s standard resolver; unsupported record types return Node-shaped errors | Go standard `net.Resolver`, local validators and IP helpers |
| `node:_http_incoming` | `src/js/node/_http_incoming.ts` | 737 | 737 | IncomingMessage constructor, lazy headers/trailers, read and destroy lifecycle, and header folding mapped onto local streams | `internal/streams/readable`, EOS, `internal/http`; local stream and native handle adapters |
| `node:_http_common` | `src/js/node/_http_common.ts`; `ProcessBindingHTTPParser.cpp`; `node/http/JSHTTPParser*.{cpp,h}`; `node/http/JSConnectionsList*.{cpp,h}`; `node/http/NodeHTTPParser.{cpp,h}`; `node/http/llhttp/{llhttp.h,api.c,http.c}` | 4,209 | _http_common.ts:361/361; process binding:117/117; HTTPParser and ConnectionsList bindings/impl:2101/2101; llhttp.h:925/925; api.c:517/517; http.c:170/170; generated llhttp.c:0/10269 (entire file reviewed) | parser methods/constants, separate 35/47 method arrays, private ConnectionsList sets, parser connection deadlines, header-pair flushes/limits, callback lifecycle, duplicate Content-Length rejection, transfer-encoding ordering and conflict checks, supported HTTP versions and HTTP/RTSP/ICE method constraints, incremental body delivery, same-buffer message dispatch, strict/lenient post-close handling, chunk size overflow checks, chunk-extension token/quoted-value validation, trailer validation, and optional CR/LF chunk framing flags are represented; the generated llhttp machine's exact state/error behavior remains untranslated, including HTTP/2 preface handling, URL grammar, chunk-boundary callback segmentation, close-delimited EOF edge cases, callback offsets, and exact error positions | `internal/validators`, `internal/freelist`, `node:_http_incoming`; Go byte/string parser adapter and standard HTTP framing |
| `node:_http_outgoing` | `src/js/node/_http_outgoing.ts` | 1335 | 1335 | OutgoingMessage construction, header APIs, serialization, write buffering, trailers, finish/flush, and socket accessors mapped; deprecation warnings and exact error behavior are approximate | `internal/stream`, validators, HTTP shared symbols, `_http_common`, `internal/http/FakeSocket`, events |
| `node:diagnostics_channel` | `src/js/node/diagnostics_channel.ts` | 425 | 425 | channel and tracing APIs follow Bun's active/inactive prototypes, live subscriber array, Map-ordered stores, async-context wrapping, next-tick error reporting, and weak cache cleanup; Go runtime cleanup timing and JavaScript engine finalization details remain different | generated runtime values, Go `weak` references and `runtime.AddCleanup`, local async hooks |
| `node:trace_events` | `src/js/node/trace_events.ts`; `src/js/internal/trace_events.ts` | 961 | 87/87; 874/874 | category buffers, monotonic trace timestamps, trace writing, instrumentation, and inspector event collection adapted; separate worker VMs and native hooks remain partial | Go file and crypto/zlib instrumentation, generated runtime |
| `node:worker_threads` | `src/js/node/worker_threads.ts`; `src/js/internal/worker/stdio.ts`; `src/js/internal/worker/messaging.ts`; `src/js/internal/trace_events.ts`; `src/jsc/bindings/webcore/Worker.cpp` | 3093 | 1296/1296; 152/152; 301/301; 874/874; Worker.cpp:470/470 | Worker option validation is mapped, then construction throws an explicit unsupported error because the static Go runtime cannot launch another JavaScript runtime or execute a worker source file. Message-port queue operations, environment-data cloning, and local channels use the Go `node:v8` value serializer; transfer lists move MessagePort endpoints and ArrayBuffer backing stores, preserve shared views, and detach transferred senders after a successful clone. Marked-untransferable buffers clone and marked-uncloneable values (including view backing stores) reject. Worker bootstrap/stdio, FileHandle transfer, host entry hooks, and native CPU/heap profiling remain host gaps | local event/stream adapters, Go worker messaging and V8 serializer |
| `node:assert` | `src/js/node/assert.ts` | 981 | 981 | public assertion APIs adapted; `partialDeepStrictEqual` uses recursive subset comparison; ordinary deep comparison remains a local structural comparator; AssertionError formatting and CallTracker behavior use mapped helpers | `node:util/types`, generated errors and promises, assertion error/diff helpers |
| `node:assert/strict` | `src/js/node/assert.strict.ts` | 2 | 2 | strict facade translated | `node:assert` |
| `node:async_hooks` | `src/js/node/async_hooks.ts`; `src/jsc/bindings/NodeAsyncHooks.cpp/.h` | 664 | 634/634; cpp:21/21; h:9/9 | AsyncLocalStorage and AsyncResource adapted; ALS options validate as objects, `run` replaces its prior binding and restores it across nested scopes and disable mutations, `exit` restores via `run(undefined, ...)`, and `withScope` enters then restores its saved store; `bind` preserves receiver/argument behavior and function lengths. AsyncResource validates safe `triggerAsyncId` values and empty types after hooks are enabled. `createHook` validates callbacks, warns once for unsupported phases, and dispatches TickObject `init` from `process.nextTick`; execution/trigger IDs remain zero and executionAsyncResource returns `process.stdin`, with Bun's warning. Other lifecycle phases stay unsupported to match Bun's wrapper | event-loop async-context frames, `internal/async_hooks_tick` |
| `node:constants` | `src/jsc/modules/NodeConstantsModule.h` | 981 | 981 | source translated; host POSIX and crypto tables emitted; platform-conditional symbols vary by host | Go standard `syscall` constants |
| `node:process` | `src/jsc/modules/NodeProcessModule.h`; `src/runtime/node/node_process.rs`; `src/runtime/node/memory_pressure.rs`; `src/runtime/node/uv_signal_handle_windows.rs`; `src/jsc/bindings/BunProcess.cpp`; `src/jsc/bindings/ProcessBindingConstants.cpp`; `.h` | 7740 | 24/24; 595/595; memory_pressure.rs:582/582; uv_signal_handle_windows.rs:66/66; BunProcess.cpp:5172/5172; binding constants:1301/1301 | process APIs adapted; Linux memory-pressure PSI is tied to listener lifecycle and emits `critical`; signal listeners start and stop Go signal subscriptions and receive Bun’s signal-name/number arguments. macOS kqueue, Windows low-memory notifications, exact Windows `uv_signal_t` lifetime, watch-mode sticky signals, and some engine/host APIs remain unavailable or approximate | Go standard `os`, `runtime`, `syscall`, `os/user`; `golang.org/x/sys/unix` PSI polling; local events runtime |
| `node:buffer` | `src/jsc/modules/NodeBufferModule.cpp`; `.h`; `src/runtime/node/buffer.rs`; `src/jsc/bindings/JSBuffer.cpp`; `.h`; `JSBufferEncodingType.cpp`; `.h`; `BufferEncodingType.h`; `src/js/builtins/JSBufferConstructor.ts`; `JSBufferPrototype.ts` | 5011 | NodeBufferModule.cpp:282/282; .h:210/210; buffer.rs:141/141; JSBuffer.cpp:4016/4016; JSBuffer.h:100/100; encoding sources:196/196; TS sources:66/66 | source sweep complete; API mapped with noted runtime approximations | Go typed arrays, standard `bytes`, `encoding/*`, `math/big`, and local buffer methods |
| `node:module` | `src/jsc/modules/NodeModuleModule.cpp/.h`; `src/jsc/NodeModuleModule.rs`; `src/jsc/NodeModuleModule.bind.ts` | 1549 | 1257/1257; 45/45; 239/239; 8/8 | source sweep complete; JSON CommonJS dependencies load through the Go filesystem; JavaScript CommonJS source loading through `_compile`, `runMain`, loader registration, and preloading throws `NotImplementedError` because it requires dynamic source execution. Mutable `_resolveFilename` and `_extensions` hooks and the public cache remain exposed, and SourceMap VLQ lookup is available; package `exports`, native `.node` loading, source-map registry, persistent compile cache, and exact host errors remain gaps | Go `path/filepath`, `os`, generated builtin registry; Bun JSC/N-API addon loader not ported |
| `node:path` | `src/js/node/path.ts`; `src/jsc/bindings/Path.cpp`; `src/runtime/node/path.rs`; `src/node_path/lib.rs` | 3919 | 81/81; 168/168; 956/956; 2714/2714 | source sweep complete; all exported POSIX/Windows operations mapped, host cwd and glob behavior adapted; Windows drive-specific CWD lookup, Bun cached-CWD behavior, and exact native path error mapping remain approximate | Go standard `os`, `path`, `filepath`, `strings`, `unicode/utf16`; local path validation and glob matcher |
| `node:path/posix` | `src/js/node/path.posix.ts` | 2 | 2 | source translated; routes to `node:path`.posix | `node:path` |
| `node:path/win32` | `src/js/node/path.win32.ts` | 2 | 2 | source translated; routes to `node:path`.win32 | `node:path` |
| `node:crypto` | `src/js/node/crypto.ts`; `src/jsc/bindings/node/crypto/node_crypto_binding.cpp`; `src/runtime/node/node_crypto_binding.rs`; DiffieHellman C++ helpers; `src/jsc/bindings/ncrypto.cpp` | 5,367 | 562/562; 403/403; 1,750/1,750; 832/832; 3,754/3,754 | all tracked crypto sources mapped; Go imports common PEM/DER/JWK RSA, EC, Ed25519, and X25519 keys, supports legacy encrypted PEM and PBES2 PKCS#8 AES/3DES, all RFC 2409/3526 named MODP groups, plus hashes, ciphers, signatures, ECDH, PBKDF2, HKDF, random operations, scrypt, and Argon2d/i/id. Gaps: RSA-PSS per-key/options details and uncommon signatures, DSA/Ed448/ML-DSA, PKCS#8 PBES2 PRFs/ciphers beyond implemented set, RSA/EC JWK optional metadata, DH private-key CRT parameters and PKCS#3 encodings beyond PKCS#8/SPKI, X509 SAN ordering/exact formatting, full OpenSSL algorithm enumeration, FIPS/secure heap, and subtle WebCrypto | Go `crypto/*`, `crypto/rand`, `crypto/x509`, x/crypto PBKDF2 and local BSD-licensed Argon2 fork; no OpenSSL/cgo |
| `node:_stream_duplex` | `src/js/node/_stream_duplex.ts` | 3 | 3/3 | exports the mapped internal Duplex constructor | `internal/streams/duplex` |
| `node:_stream_passthrough` | `src/js/node/_stream_passthrough.ts` | 3 | 3/3 | exports the mapped internal PassThrough constructor | `internal/streams/passthrough` |
| `node:_stream_readable` | `src/js/node/_stream_readable.ts` | 3 | 3/3 | exports the mapped internal Readable constructor | `internal/streams/readable` |
| `node:_stream_transform` | `src/js/node/_stream_transform.ts` | 3 | 3/3 | exports the mapped internal Transform constructor | `internal/streams/transform` |
| `node:_stream_wrap` | `src/js/node/_stream_wrap.ts` | 5 | 5/5 | exports `node:stream` and emits the `DEP0125` deprecation warning once | `node:stream`, `node:process` |
| `node:_stream_writable` | `src/js/node/_stream_writable.ts` | 3 | 3/3 | exports the mapped internal Writable constructor | `internal/streams/writable` |
| `node:stream` | `src/js/node/stream.ts`; `src/js/internal/stream.ts` | 120 | 120 | classic stream facade, constructors, predicates, pipeline, operators, and pairing adapted; concurrency and scheduling are approximate | mapped internal stream runtime |
| `node:stream/promises` | `src/js/node/stream.promises.ts`; `src/js/internal/stream.promises.ts` | 47 | 47 | promise `pipeline` and `finished` entry points exposed | `node:stream`, EOS, pipeline |
| `node:stream/consumers` | `src/js/node/stream.consumers.ts`; `src/jsc/bindings/webcore/streams/BunStreamConsumers.cpp` | 1824 | 54/54; BunStreamConsumers.cpp:1770/1770 | six consumers use the local async-iterator pump and Blob implementation for string, Blob, ArrayBuffer, view, and Buffer chunks; read errors reject and JSON parsing errors reject. Bun's direct-pull sink, native buffered fast path, stream lock/disturbance error ordering, and native-readable pointer transfer remain partial | async iterator, Buffer, Blob, Web Streams adapters |
| `node:stream/web` | `src/js/node/stream.web.ts` | 20 | 20 | Web Stream constructor names forwarded from the global object | generated Web Stream values |
| `node:zlib/iter` | `src/js/node/zlib.iter.ts`; `src/js/internal/streams/iter/transform.ts` | 811 | 51/51; 760/760 | 16 async/sync gzip, deflate, Brotli, and Zstd transform factories exposed; Go standard library handles gzip/deflate, `andybalholm/brotli` handles Brotli streams, and `klauspost/compress` handles Zstd; Brotli dictionaries and most native Brotli params are unsupported, and streaming/backpressure is approximate | Go standard `compress/gzip`, `compress/zlib`, `compress/flate`; `github.com/andybalholm/brotli`; `github.com/klauspost/compress/zstd`; `node:buffer`, `node:util/types` |
| `node:zlib` | `src/js/node/zlib.ts`; `src/runtime/node/node_zlib_binding.rs`; `src/runtime/node/zlib/NativeZlib.rs`; `src/runtime/node/zlib/NativeBrotli.rs`; `src/runtime/node/zlib/NativeZstd.rs` | 4165 | 1021/1021; 1066/1066; 654/654; 737/737; 687/687 | All wrapper and native source ranges mapped. Go standard deflate/gzip, `andybalholm/brotli`, and the existing `klauspost/compress` Zstd codec back the public API. Brotli quality/window are supported; dictionaries, Brotli mode/block/decoder parameters, exact native error codes, native worker scheduling, zlib window/memory/strategy controls, and incremental decompression remain gaps | Go standard `compress/gzip`, `compress/zlib`, `compress/flate`; `github.com/andybalholm/brotli`; `github.com/klauspost/compress/zstd`; `node:buffer`, `node:stream`, `internal/validators` |
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
| `src/js/internal/assert/assertion_error.ts` | 448 | 448 |
| `src/js/internal/assert/calltracker.ts` | 136 | 136 |
| `src/js/internal/assert/myers_diff.ts` | 42 | 42 |
| `src/js/internal/assert/utils.ts` | 35 | 35 |
| `src/js/internal/async_context_frame.ts` | 61 | 61 |
| `src/js/internal/async_hooks.ts` | 28 | 28 |
| `src/js/internal/async_hooks_tick.ts` | 20 | 20 |
| `src/js/internal/cluster/RoundRobinHandle.ts` | 222 | 222 |
| `src/js/internal/cluster/SharedHandle.ts` | 71 | 71 |
| `src/js/internal/cluster/Worker.ts` | 60 | 60 |
| `src/js/internal/cluster/child.ts` | 376 | 376 |
| `src/js/internal/cluster/isPrimary.ts` | 4 | 4 |
| `src/js/internal/cluster/primary.ts` | 380 | 380 |
| `src/js/internal/dgram.ts` | 146 | 146 |
| `src/js/internal/errors.ts` | 23 | 23 |
| `src/js/internal/fixed_queue.ts` | 130 | 130 |
| `src/js/internal/freelist.ts` | 27 | 27 |
| `src/js/internal/fs/binding.ts` | 123 | 123 |
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
| `src/js/internal/net/server_abort_signal.ts` | 26 | 26 |
| `src/js/internal/net/symbols.ts` | 15 | 15 |
| `src/js/internal/perf_hooks/monitorEventLoopDelay.ts` | 71 | 71 |
| `src/js/internal/primordials.js` | 176 | 176 |
| `src/js/internal/promisify.ts` | 81 | 81 |
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
| `src/js/internal/shared.ts` | 456 | 456 |
| `src/js/internal/streams/native-readable.ts` | 306 | 306 |
| `src/js/internal/streams/operators.ts` | 374 | 374 |
| `src/js/internal/streams/passthrough.ts` | 20 | 20 |
| `src/js/internal/streams/pipeline.ts` | 448 | 448 |
| `src/js/internal/streams/readable.ts` | 1734 | 1734 |
| `src/js/internal/streams/state.ts` | 52 | 52 |
| `src/js/internal/streams/transform.ts` | 173 | 173 |
| `src/js/internal/streams/utils.ts` | 305 | 305 |
| `src/js/internal/validators.ts` | 180 | 180 |
| `src/js/internal/streams/writable.ts` | 1118 | 1118 |
| `src/js/internal/test_runner/mock_timers.ts` | 781 | 781 |
| `src/js/internal/timers.ts` | 32 | 32 |
| `src/js/internal/tls.ts` | 476 | 476 |
| `src/js/internal/trace_events.ts` | 874 | 874 |
| `src/js/internal/tty.ts` | 205 | 205 |
| `src/js/internal/url.ts` | 33 | 33 |
| `src/js/internal/util/colors.ts` | 38 | 38 |
| `src/js/internal/util/comparisons.ts` | 9 | 9 |
| `src/js/internal/util/deprecate.ts` | 60 | 60 |
| `src/js/internal/util/inspect.js` | 2830 | 2830 |
| `src/js/internal/util/mime.ts` | 6 | 6 |
| `src/js/internal/validators.ts` | 180 | 180 |
| `src/js/internal/webstreams_adapters.ts` | 864 | 864 |
| `src/js/internal/worker/messaging.ts` | 301 | 301 |
| `src/js/internal/worker/stdio.ts` | 152 | 152 |
| `src/js/node/_http2_upgrade.ts` | 405 | 405 |
| `src/js/node/_http_agent.ts` | 626 | 626 |
| `src/js/node/_http_client.ts` | 1187 | 1187 |
| `src/js/node/_http_common.ts` | 361 | 361 |
| `src/js/node/_http_incoming.ts` | 737 | 737 |
| `src/js/node/_http_outgoing.ts` | 1335 | 1335 |
| `src/js/node/_http_server.ts` | 4306 | 4306 |
| `src/js/node/_stream_duplex.ts` | 3 | 3 |
| `src/js/node/_stream_passthrough.ts` | 3 | 3 |
| `src/js/node/_stream_readable.ts` | 3 | 3 |
| `src/js/node/_stream_transform.ts` | 3 | 3 |
| `src/js/node/_stream_wrap.ts` | 5 | 5 |
| `src/js/node/_stream_writable.ts` | 3 | 3 |
| `src/js/node/_tls_common.ts` | 43 | 43 |
| `src/js/node/_tls_wrap.ts` | 14 | 14 |
| `src/js/node/assert.strict.ts` | 2 | 2 |
| `src/js/node/assert.ts` | 981 | 981 |
| `src/js/node/async_hooks.ts` | 634 | 634 |
| `src/js/node/child_process.ts` | 2071 | 2071 |
| `src/js/node/cluster.ts` | 23 | 23 |
| `src/js/node/console.ts` | 1 | 1 |
| `src/js/node/crypto.ts` | 562 | 562 |
| `src/js/internal/streams/lazy_transform.ts` | 53 | 53 |
| `src/js/node/dgram.ts` | 1546 | 1546 |
| `src/js/node/diagnostics_channel.ts` | 425 | 425 |
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
| `src/js/node/net.ts` | 4634 | 4634 |
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
| `src/js/node/worker_threads.ts` | 1296 | 1296 |
| `src/js/node/zlib.iter.ts` | 51 | 51 |
| `src/js/node/zlib.ts` | 1021 | 1021 |

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
| `src/js/internal/cluster/SharedHandle.ts` | Rust | `node_cluster_binding.rs` | `src/runtime/node/node_cluster_binding.rs` |
| `src/js/internal/cluster/child.ts` | Rust | `node_cluster_binding.rs` | `src/runtime/node/node_cluster_binding.rs` |
| `src/js/internal/cluster/primary.ts` | Rust | `node_cluster_binding.rs` | `src/runtime/node/node_cluster_binding.rs` |
| `src/js/internal/cluster/RoundRobinHandle.ts` | Rust | `node_cluster_binding.rs` | `src/runtime/node/node_cluster_binding.rs` |
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
| `src/js/node/events.ts` | C++ | `EventTarget.cpp/.h` | `src/jsc/bindings/webcore/EventTarget.cpp`, `EventTarget.h` |
| `src/js/node/events.ts` | C++ | `JSEventTarget.h` | `src/jsc/bindings/webcore/JSEventTarget.h` |
| `src/js/node/events.ts` | C++ | `JSEventTargetCustom.cpp/.h` | `src/jsc/bindings/webcore/JSEventTargetCustom.cpp`, `JSEventTargetCustom.h` |
| `src/js/node/events.ts` | C++ | `JSEventTargetNode.h` | `src/jsc/bindings/webcore/JSEventTargetNode.h` |
| `src/js/node/events.ts` | C++ | `JSEvent.cpp/.h` | `src/jsc/bindings/webcore/JSEvent.cpp`, `JSEvent.h` |
| `src/js/node/events.ts` | C++ | `Event.cpp/.h` | `src/jsc/bindings/webcore/Event.cpp`, `Event.h` |
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
| `node:sqlite` | `src/jsc/bindings/sqlite/NodeSqlite.cpp` | 4135 | 4135 |
| `node:sqlite` | `src/jsc/bindings/sqlite/NodeSqlite.h` | 902 | 902 |
| `node:sqlite` | `src/jsc/bindings/sqlite/JSSQLStatement.cpp` | 3077 | 3077 |
| `node:sqlite` | `src/jsc/bindings/sqlite/JSSQLStatement.h` | 89 | 89 |
| `node:sqlite` | `src/jsc/bindings/sqlite/sqlite3.c` | 269669 | 0 (replaced by local bundled SQLite C engine; source lines 1–269669 skipped) |
| `node:sqlite` | `src/jsc/bindings/sqlite/sqlite3_local.h` | 14350 | 3500 (lines 3501–14350 replaced by driver declarations/engine) |
| `node:sqlite` | `src/jsc/bindings/sqlite/sqlite3_error_codes.h` | 764 | 0 (replaced by driver error codes and Go runtime constants) |
| `node:sqlite` | `src/jsc/bindings/sqlite/lazy_sqlite3.h` | 521 | 0 (replaced by bundled driver initialization; lines 1–521 skipped) |
| `node:string_decoder` | `src/jsc/modules/NodeStringDecoderModule.h`; `src/jsc/bindings/JSStringDecoder.h`; `.cpp` | 774 | NodeStringDecoderModule.h:17/17; JSStringDecoder.h:126/126; JSStringDecoder.cpp:631/631 | source sweep complete; stateful decoder API adapted | Go UTF-8/UTF-16 decoders, Buffer encoding helpers |
| `node:stream` | `src/js/node/stream.ts`; `src/js/internal/stream.ts`; `src/js/internal/stream.promises.ts`; `src/js/internal/streams/operators.ts`; `src/js/internal/streams/duplexpair.ts` | 617 | 8/8; 112/112; 45/45; 374/374; 78/78 | API surface mapped; stream operator concurrency, abort cleanup, and pair scheduling are approximate | shared stream runtime, iterators, event loop |
| `node:stream/consumers` | `src/js/node/stream.consumers.ts`; `BunStreamConsumers.cpp` | 1824 | 54/54; 1770/1770 | arrayBuffer, bytes, text, JSON, Buffer, and Blob consumers use a rejection-safe iterator pump and local Blob conversions; direct-pull sources, buffered native fast paths, lock/disturbed ordering, and native-readable transfer remain partial | async iterator, Buffer, global Blob, `BunStreamConsumers.cpp` |
| `node:stream/promises` | `src/js/node/stream.promises.ts` | 2 | 2/2 | default exports promise `pipeline` and `finished` | `node:stream` |
| `node:stream/web` | `src/js/node/stream.web.ts` | 20 | 20/20 | constructor exports forward global Web Stream values; unavailable globals stay undefined | global Web Stream values |
| `node:crypto` | `src/js/node/crypto.ts` | 562 | 562/562 | exported wrapper API mapped; RSA/EC/Ed25519/X25519 JWK import/export, PEM/DER key parsing, legacy encrypted PEM and common PBES2 PKCS#8, Argon2d/i/id, common SPKAC signatures, scrypt, and RFC MODP Diffie-Hellman groups use Go implementations; DSA/Ed448/ML-DSA, uncommon SPKAC algorithms, DH private-key CRT parameters and non-PKCS#8/SPKI encodings, FIPS/secure heap, and subtle WebCrypto remain partial or unavailable | Go `crypto/*`, x/crypto PBKDF2, local BSD-licensed Argon2 fork, random source; C++/Rust binding behavior replaced rather than line-translated |
| `internal/streams/lazy_transform` | `src/js/internal/streams/lazy_transform.ts` | 53 | 53/53 | lazy Transform initialization, accessors, and constructor surface adapted | `internal/streams/transform` |
| `internal/streams/transform` | `src/js/internal/streams/transform.ts` | 173 | 173/173 | constructor, flush, prefinish, callback backpressure, and read/write hooks adapted | `internal/streams/duplex`, `internal/streams/state` |
| `internal/streams/duplex` | `src/js/internal/streams/duplex.ts` | 166 | 166/166 | constructor, combined prototype, duplex-side flags, and Web statics adapted | `internal/streams/legacy`, Readable, Writable, Web adapters, Duplexify |
| `internal/streams/duplexify` | `src/js/internal/streams/duplexify.ts` | 395 | 395/395 | Node/Web stream, iterable, promise, and readable/writable pair conversion adapted; async-generator handling is approximate | Duplex, Readable, Writable, `internal/streams/from`, EOS, destroy |
| `internal/streams/writable` | `src/js/internal/streams/writable.ts` | 1118 | 1118/1118 | constructor, queue, writev batching, lifecycle, properties, and Web statics adapted; callback ordering and WritableState internals are approximate | legacy, destroy, EOS, abort signal, state, utils, Web adapters |
| `internal/streams/native-readable` | `src/js/internal/streams/native-readable.ts` | 306 | 306/306 | Readable wrapper maps pointer methods and pull results; Bun C++ transfer entry is a placeholder | Readable, destroy, `BunStreamConsumers.cpp` |
| `internal/webstreams_adapters` | `src/js/internal/webstreams_adapters.ts` | 864 | 864/864 | all six classic/Web conversion entry points exposed; strategy and native transfer behavior are approximate | Writable, Readable, Duplex, destroy, utils, shared, validators, util/types, EOS, native-readable |
| `internal/streams/legacy` | `src/js/internal/streams/legacy.ts` | 146 | 146/146 | source mapped; EventEmitter bridge, pipe, and Buffer helpers adapted | `node:events`, `node:util/types`, Buffer |
| `internal/abort_listener` | `src/js/internal/abort_listener.ts` | 31 | 31/31 | addAbortListener and Disposable removal adapted | AbortSignal event methods, shared stop-propagation option |
| `internal/async_hooks` (private) | `src/js/internal/async_hooks.ts` | 28 | 28/28 | enabled hook counter and owner symbol adapted; `TickObject` init delivery is wired through `node:async_hooks`; other hook phases remain unsupported | `node:async_hooks` adapter state |
| `internal/async_hooks_tick` | `src/js/internal/async_hooks_tick.ts` | 20 | 20/20 | stable hook array and incrementing async ID helper mapped; `process.nextTick` delivers snapshotted init hooks before enqueueing | `node:async_hooks`, `node:process` next-tick queue |
| `internal/fixed_queue` | `src/js/internal/fixed_queue.ts` | 130 | 130/130 | segmented 2,048-slot circular buffers, null sentinel behavior, buffer links, queue push/shift, and exported `FixedQueue` class mapped | `node:events` |
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
| `node:module` | `NodeModuleModule.cpp/.h`, `NodeModuleModule.rs`, `NodeModuleModule.bind.ts` | 1549 | 1549/1549 | source sweep complete; JSON CommonJS dependencies load through the Go filesystem; JavaScript CommonJS source loading through `_compile`, `runMain`, loader registration, and preloading throws `NotImplementedError` because it requires dynamic source execution. Mutable `_resolveFilename` and `_extensions` hooks and the public cache remain exposed, and SourceMap VLQ lookup is available; package `exports`, native `.node` loading, source-map registry, persistent compile cache, and exact host errors remain gaps | Go `path/filepath`, `os`, generated builtin registry; Bun JSC/N-API addon loader not ported |
| `node:url` | `src/js/node/url.ts`; `src/js/internal/url.ts`; NodeURL C++ bridge and helper headers | 1969 | 1382/1382; 33/33; 308/308; 10/10; 14/14; 222/222 | source sweep complete; WHATWG and legacy URL surfaces adapted; ICU UTS #46 and detailed URLPattern matching approximated | Go `net/url`, `path/filepath`, local punycode and Buffer runtimes |
| `internal/util/inspect` | `src/js/internal/util/inspect.js` | 2830 | 2830/2830 | source sweep complete; object formatting, Map/Set iterator previews via private runtime tags, printf formats, error stacks, ANSI stripping, custom inspection, and helper APIs adapted; exact grapheme widths and some property details remain approximate | primordials, `node:util/types`, validators, `node:module`, `node:url`, `node:buffer`, `UtilInspect.cpp`, `stringWidth.cpp` |
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
| `src/jsc/bindings/JSInspectorProfiler.cpp` | 159 | 159 |
| `src/jsc/bindings/JSNodePerformanceHooksHistogram.cpp` | 274 | 274 |
| `src/jsc/bindings/JSNodePerformanceHooksHistogram.h` | 165 | 165 |
| `src/jsc/bindings/JSNodePerformanceHooksHistogramConstructor.cpp` | 112 | 112 |
| `src/jsc/bindings/JSNodePerformanceHooksHistogramConstructor.h` | 49 | 49 |
| `src/jsc/bindings/JSNodePerformanceHooksHistogramPrototype.cpp` | 529 | 529 |
| `src/jsc/bindings/JSNodePerformanceHooksHistogramPrototype.h` | 46 | 46 |
| `src/jsc/bindings/ModuleGraph.cpp` | 855 | 0 |
| `src/jsc/bindings/NodeAsyncHooks.cpp` | 21 | 21 |
| `src/jsc/bindings/NodeAsyncHooks.h` | 9 | 9 |
| `src/jsc/bindings/NodeDirent.cpp` | 375 | 375 |
| `src/jsc/bindings/NodeFSStatBinding.cpp` | 967 | 967 |
| `src/jsc/bindings/NodeFSStatFSBinding.cpp` | 465 | 465 |
| `src/jsc/bindings/ProcessBindingFs.cpp` | 216 | 216 |
| `src/jsc/bindings/ProcessBindingFs.h` | 36 | 36 |
| `src/jsc/bindings/NodeHTTP.cpp` | 1023 | 1023 |
| `src/jsc/bindings/NodeHTTP.h` | 8 | 8 |
| `src/jsc/bindings/NodeHTTPRequestHeaderValues.h` | 48 | 48 |
| `src/jsc/bindings/ProcessBindingHTTPParser.cpp` | 84 | 84 |
| `src/jsc/bindings/ProcessBindingHTTPParser.h` | 33 | 33 |
| `src/jsc/bindings/node/JSNodeHTTPServerSocket.cpp` | 1141 | 1141 |
| `src/jsc/bindings/node/JSNodeHTTPServerSocket.h` | 216 | 216 |
| `src/jsc/bindings/node/JSNodeHTTPServerSocketPrototype.cpp` | 631 | 631 |
| `src/jsc/bindings/node/JSNodeHTTPServerSocketPrototype.h` | 46 | 46 |
| `src/jsc/bindings/node/http/JSHTTPParser.cpp` | 52 | 52 |
| `src/jsc/bindings/node/http/JSHTTPParser.h` | 66 | 66 |
| `src/jsc/bindings/node/http/JSHTTPParserConstructor.cpp` | 61 | 61 |
| `src/jsc/bindings/node/http/JSHTTPParserConstructor.h` | 46 | 46 |
| `src/jsc/bindings/node/http/JSHTTPParserPrototype.cpp` | 381 | 381 |
| `src/jsc/bindings/node/http/JSHTTPParserPrototype.h` | 44 | 44 |
| `src/jsc/bindings/node/http/JSConnectionsList.cpp` | 192 | 192 |
| `src/jsc/bindings/node/http/JSConnectionsList.h` | 67 | 67 |
| `src/jsc/bindings/node/http/JSConnectionsListConstructor.cpp` | 35 | 35 |
| `src/jsc/bindings/node/http/JSConnectionsListConstructor.h` | 50 | 50 |
| `src/jsc/bindings/node/http/JSConnectionsListPrototype.cpp` | 123 | 123 |
| `src/jsc/bindings/node/http/JSConnectionsListPrototype.h` | 44 | 44 |
| `src/jsc/bindings/node/http/NodeHTTPParser.cpp` | 725 | 725 |
| `src/jsc/bindings/node/http/NodeHTTPParser.h` | 215 | 215 |
| `src/jsc/bindings/node/http/llhttp/README.md` | 18 | 18 |
| `src/jsc/bindings/node/http/llhttp/llhttp.h` | 925 | 925 |
| `src/jsc/bindings/node/http/llhttp/api.c` | 517 | 517 |
| `src/jsc/bindings/node/http/llhttp/http.c` | 170 | 170 |
| `src/jsc/bindings/node/http/llhttp/llhttp.c` | 10269 | 0 |
| `src/jsc/bindings/NodeTLS.cpp` | 362 | 362 |
| `src/runtime/api/bun/SecureContext.rs` | 460 | 460 |
| `src/runtime/socket/SSLConfig.rs` | 531 | 531 |
| `src/http/ssl_config.rs` | 699 | 699 |
| `src/runtime/socket/SSLConfig.bindv2.ts` | 116 | 116 |
| `src/jsc/bindings/NodeURL.cpp` | 308 | 308 |
| `src/jsc/bindings/NodeURL.h` | 10 | 10 |
| `src/jsc/bindings/NodeURLHelpers.h` | 14 | 14 |
| `src/jsc/bindings/ASCIIHostPunycodeCheck.h` | 222 | 222 |
| `src/jsc/bindings/NodeV8.cpp` | 108 | 108 |
| `src/jsc/bindings/NodeVM.cpp` | 2056 | 2056 |
| `src/jsc/bindings/NodeVM.h` | 184 | 184 |
| `src/jsc/bindings/CodeGenerationFromStrings.h` | 70 | 70 |
| `src/jsc/bindings/NodeVMScript.h` | 119 | 119 |
| `src/jsc/bindings/NodeVMScript.cpp` | 679 | 679 |
| `src/jsc/bindings/NodeVMModule.h` | 130 | 130 |
| `src/jsc/bindings/NodeVMModule.cpp` | 818 | 818 |
| `src/jsc/bindings/NodeVMSourceTextModule.h` | 70 | 70 |
| `src/jsc/bindings/NodeVMSourceTextModule.cpp` | 566 | 566 |
| `src/jsc/bindings/NodeVMSyntheticModule.h` | 60 | 60 |
| `src/jsc/bindings/NodeVMSyntheticModule.cpp` | 254 | 254 |
| `src/jsc/bindings/NodeVMScriptFetcher.h` | 104 | 104 |
| `src/jsc/bindings/NodeVMOptionNames.h` | 47 | 47 |
| `src/jsc/bindings/vm/NodeVMRunTermination.h` | 62 | 62 |
| `src/jsc/bindings/vm/NodeVMRunTermination.cpp` | 109 | 109 |
| `src/jsc/bindings/NodeValidator.cpp` | 655 | 655 |
| `src/jsc/bindings/NodeValidator.h` | 48 | 48 |
| `src/runtime/node/node_util_binding.rs` | 246 | 246 |
| `src/runtime/node/util/parse_args.rs` | 1049 | 1049 |
| `src/runtime/node/util/parse_args_utils.rs` | 149 | 149 |
| `src/jsc/bindings/Path.cpp` | 168 | 168 |
| `src/jsc/bindings/BunTTYState.h` | 26 | 26 |
| `src/jsc/bindings/ProcessBindingTTYWrap.cpp` | 533 | 533 |
| `src/jsc/bindings/ProcessBindingTTYWrap.h` | 18 | 18 |
| `src/jsc/bindings/wtf-bindings.cpp` | 320 | 223 |
| `src/jsc/bindings/UtilInspect.cpp` | 121 | 121 |
| `src/jsc/bindings/node/crypto/node_crypto_binding.cpp` | 403 | 403 |
| `src/jsc/bindings/ncrypto.cpp` | 3754 | 3754 |
| `src/jsc/bindings/node/crypto/DiffieHellmanFunctions.h` | 286 | 286 |
| `src/jsc/bindings/node/crypto/JSDiffieHellmanConstructor.cpp` | 195 | 195 |
| `src/jsc/bindings/node/crypto/JSDiffieHellmanGroupConstructor.cpp` | 57 | 57 |
| `src/jsc/bindings/node/crypto/JSDiffieHellmanPrototype.cpp` | 84 | 84 |
| `src/jsc/bindings/node/crypto/JSDiffieHellmanGroupPrototype.cpp` | 70 | 70 |
| `src/jsc/bindings/node/crypto/JSDiffieHellman.cpp` | 71 | 71 |
| `src/jsc/bindings/node/crypto/JSDiffieHellmanGroup.cpp` | 69 | 69 |
| `src/jsc/bindings/stringWidth.cpp` | 1178 | 1178 |
| `src/jsc/bindings/webcore/JSEventTarget.cpp` | 370 | 370 |
| `src/jsc/bindings/webcore/JSEventTargetNode.cpp` | 49 | 49 |
| `src/jsc/bindings/webcore/JSEventTarget.h` | 102 | 102 |
| `src/jsc/bindings/webcore/JSEventTargetCustom.h` | 86 | 86 |
| `src/jsc/bindings/webcore/JSEventTargetCustom.cpp` | 88 | 88 |
| `src/jsc/bindings/webcore/JSEventTargetNode.h` | 7 | 7 |
| `src/jsc/bindings/webcore/EventTarget.cpp` | 353 | 353 |
| `src/jsc/bindings/webcore/EventTarget.h` | 236 | 236 |
| `src/jsc/bindings/webcore/JSEvent.cpp` | 613 | 613 |
| `src/jsc/bindings/webcore/JSEvent.h` | 93 | 93 |
| `src/jsc/bindings/webcore/Event.cpp` | 156 | 156 |
| `src/jsc/bindings/webcore/Event.h` | 183 | 183 |
| `src/jsc/bindings/webcore/JSMIMEParams.cpp` | 722 | 722 |
| `src/jsc/bindings/webcore/JSMIMEParams.h` | 89 | 89 |
| `src/jsc/bindings/webcore/JSMIMEType.cpp` | 614 | 614 |
| `src/jsc/bindings/webcore/JSMIMEType.h` | 102 | 102 |
| `src/jsc/bindings/webcore/Worker.cpp` | 470 | 470 |
| `src/jsc/bindings/webcore/streams/BunStreamConsumers.cpp` | 1770 | 1770 |
| `src/jsc/modules/NodeModuleModule.cpp` | 1257 | 1257 |
| `src/jsc/modules/NodeUtilTypesModule.cpp` | 1289 | 1289 |
| `src/runtime/api/bun/h2_frame_parser.rs` | 7873 | 7873 |
| `src/runtime/api/h2.classes.ts` | 151 | 151 |
| `src/jsc/bindings/H2HeadersMaterializer.cpp` | 210 | 210 |
| `src/runtime/api/bun/h2/mod.rs` | 12 | 12 |
| `src/runtime/api/bun/h2/wire.rs` | 350 | 350 |
| `src/runtime/api/bun/h2/settings.rs` | 120 | 120 |
| `src/runtime/api/bun/h2/flow_control.rs` | 153 | 153 |
| `src/runtime/api/bun/h2/hpack.rs` | 112 | 112 |
| `src/runtime/api/bun/h2/stream.rs` | 145 | 145 |
| `src/runtime/api/bun/h2/connection.rs` | 2462 | 2462 |
| `src/http/lshpack.rs` | 238 | 238 |
| `src/jsc/bindings/c-bindings.cpp` | 1254 | 1254 |
| `src/runtime/node/node_assert_binding.rs` | 147 | 147 |
| `src/runtime/node/node_crypto_binding.rs` | 1750 | 1750 |
| `src/runtime/node/node_fs_binding.rs` | 527 | 527 |
| `src/runtime/node/node_net_binding.rs` | 184 | 184 |
| `src/runtime/node/node_os.rs` | 1614 | 1614 |
| `src/runtime/node/path.rs` | 956 | 956 |
| `src/runtime/node/node_quic_binding.rs` | 284 | 284 |
| `src/runtime/node/node_zlib_binding.rs` | 1066 | 1066 |

## Node runtime Rust source index

The Rust files below are Bun’s Node host integration layer. A module’s concrete dependencies are shown in the JavaScript edge graph and direct bridge table; this index keeps the broader runtime layer available for selection.

| Source file | Lines | Through |
|---|---:|---:|
| `src/runtime/node/Stat.rs` | 212 | 212 |
| `src/runtime/node/StatFS.rs` | 160 | 160 |
| `src/runtime/node/assert/myers_diff.rs` | 512 | 512 |
| `src/runtime/node/buffer.rs` | 141 | 141 |
| `src/runtime/node/dir_iterator.rs` | 821 | 821 |
| `src/runtime/node/fs_events.rs` | 939 | 939 |
| `src/runtime/node/memory_pressure.rs` | 582 | 582 |
| `src/runtime/node/net/BlockList.rs` | 524 | 524 |
| `src/runtime/socket/SocketAddress.rs` | 1084 | 1084 |
| `src/runtime/node/node.classes.ts` | 384 | 384 |
| `src/runtime/node/node_assert.rs` | 365 | 365 |
| `src/runtime/node/node_assert_binding.rs` | 147 | 147 |
| `src/runtime/node/node_cluster_binding.rs` | 810 | 810 |
| `src/runtime/node/node_crypto_binding.rs` | 1750 | 1750 |
| `src/runtime/node/node_fs.rs` | 10116 | 10116 |
| `src/runtime/node/node_fs_binding.rs` | 527 | 527 |
| `src/runtime/node/node_fs_constant.rs` | 32 | 32 |
| `src/runtime/node/node_fs_stat_watcher.rs` | 1110 | 1110 |
| `src/runtime/node/node_fs_watcher.rs` | 1256 | 1256 |
| `src/runtime/node/node_http_binding.rs` | 110 | 110 |
| `src/runtime/server/NodeHTTPResponse.rs` | 2825 | 2825 |
| `src/runtime/node/node_net_binding.rs` | 184 | 184 |
| `src/runtime/node/node_os.bind.ts` | 94 | 94 |
| `src/runtime/node/node_os.rs` | 1614 | 1614 |
| `src/runtime/node/node_process.rs` | 595 | 595 |
| `src/runtime/node/node_quic_binding.rs` | 284 | 284 |
| `src/runtime/node/node_util_binding.rs` | 246 | 246 |
| `src/runtime/node/node_zlib_binding.rs` | 1066 | 1066 |
| `src/runtime/node/path.rs` | 956 | 956 |
| `src/runtime/node/path_watcher.rs` | 1618 | 1618 |
| `src/runtime/node/quic/endpoint.rs` | 2661 | 2661 |
| `src/runtime/node/quic/ffi.rs` | 103 | 103 |
| `src/runtime/node/quic/mod.rs` | 44 | 44 |
| `src/runtime/node/quic/quic.classes.ts` | 181 | 181 |
| `src/runtime/node/quic/session.rs` | 2680 | 2680 |
| `src/runtime/node/quic/stream.rs` | 1188 | 1188 |
| `src/runtime/node/quic/tls.rs` | 725 | 725 |
| `src/runtime/node/time_like.rs` | 125 | 125 |
| `src/runtime/node/types.rs` | 1779 | 1779 |
| `src/runtime/node/util/parse_args.rs` | 1049 | 1049 |
| `src/runtime/node/util/parse_args_utils.rs` | 149 | 149 |
| `src/runtime/node/util/validators.rs` | 486 | 486 |
| `src/runtime/node/uv_signal_handle_windows.rs` | 66 | 66 |
| `src/runtime/node/win_watcher.rs` | 547 | 547 |
| `src/runtime/node/zlib/NativeBrotli.rs` | 737 | 737 |
| `src/runtime/node/zlib/NativeZlib.rs` | 654 | 654 |
| `src/runtime/node/zlib/NativeZstd.rs` | 687 | 687 |

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
| Remaining focused modules | `node:sqlite`, `node:test`, `node:zlib`, `node:zlib/iter` | module-specific wrapper/native route in the tables above |

No tests were run to produce this inventory.
