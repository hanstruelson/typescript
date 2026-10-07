# Filesystem implementation checklist

**Status: incomplete.** The native core and common API paths are implemented;
this is not yet a complete replacement for Node's `fs` module.

[fs-api.json](fs-api.json) records the complete installed Node 24 export inventory
and the pinned Bun source files. Its per-export status is an implementation
inventory, not a claim that every option or overload is conformant.

## Implemented families

| Family | Go implementation | Coverage |
| --- | --- | --- |
| Imports | `node_modules.go`, compiler bundler | Literal `require`, named/default/namespace imports, dynamic import; `node:` and plain names |
| File content | `fs_operations.go`, `fs_encoding.go` | read/write/append; encoding options; native byte buffers; descriptor borrowing; flush |
| Descriptors | `fs_handles.go` | open/close, read/write overloads, vectored I/O, explicit and current positions |
| FileHandle | `fs_handles.go` | Promise I/O and metadata methods, close waits for submitted operations, repeated close, close event, stream factories |
| Metadata | `fs_operations.go` | stat/lstat/fstat, statfs, ordinary numeric fields and file-kind predicates |
| Directories | `fs_directory.go` | readdir with Dirents and recursive mode; opendir read/close; recursive traversal; async iteration and early close |
| Mutation | `fs_operations.go` | mkdir, mkdtemp, rename, unlink, truncate, chmod/chown variants, timestamps, rm/rmdir with recursive retry options |
| Links | `fs_operations.go` | hard links, symbolic links, readlink, realpath and `.native` |
| Copying | `fs_copy.go` | copyFile flags/clone ioctl, recursive cp, filters on caller loop, async filters, common collision errors |
| Globbing | `fs_glob.go` | sync/callback arrays, promise async iterator, `*`, `?`, brackets, `**`, basic braces, exclusions, Dirents |
| Watching | `fs_watch.go` | inotify watch, recursive additions, stat polling, unwatchFile, listener removal, ref/unref, promise watch iterator |
| File streams | `fs_streams.go` | ordered writes, read ranges, native chunks, flow/pause/resume, high-water-mark backpressure, pipe, async iteration, basic destruction, UTF-8/UTF-16 chunk decoding |
| Buffer subset | `fs_encoding.go` | from/alloc/concat/isBuffer/byteLength, views, string conversion, write/copy/fill/equals/compare |
| Errors | `fs_runtime.go` | Error name/message/code/errno/syscall/path/dest; argument and bounds errors |

Filesystem objects and module exports remain native `tsValue` values. Byte data
uses native `uint8` storage rather than Go `any`. The descriptor table is shared
and synchronized. Each async operation snapshots options on the caller loop and
uses its task queue; callbacks and filters run on that loop. Module objects
initialize once and are shared across user goroutines.

## Remaining work before marking complete

- [ ] Full argument/option validation and every overload, including observable synchronous callback errors.
- [ ] BigInt runtime support and BigInt Stats/StatFs (currently rejected explicitly).
- [ ] Correct birthtime/statx behavior, complete Date identity/prototype behavior, and timestamp edge cases.
- [ ] Public Stats, Dirent, Dir, ReadStream, WriteStream constructors/prototypes, aliases, and inheritance.
- [ ] `openAsBlob`, Blob, Web ReadableStream, FileHandle `readableWebStream` and `readLines`.
- [ ] AbortController/AbortSignal integration for file operations, watchers, and streams.
- [ ] FileHandle reference counting across stream lifetime, locking, disposal symbols, finalizers, and cross-loop use.
- [ ] Directory disposal symbols and complete async iterator abrupt-completion semantics.
- [ ] Complete Readable/Writable interface (`read`, `unpipe`, `cork`, `uncork`, listener behavior, custom `fs` providers, queue/error edge cases).
- [ ] Stream descriptor ownership/close options, event ordering, destruction during I/O, and encoding edge cases beyond current fixtures.
- [ ] Promise watcher queue limits/overflow options, cancellation, recursive rename races, and all watcher lifecycle events.
- [ ] Full glob grammar and Node-specific matching/order/path behavior; current implementation scans the tree and buffers matches.
- [ ] cp symlink-alias/cycle handling, atime preservation, complete clone/collision behavior, and platform-specific errors.
- [ ] Full Node Buffer API and `buffer` exports; Buffer constructor/prototype identity.
- [ ] Real file URL construction and validation through a native URL implementation.
- [ ] Import-equals, native re-exports, general CommonJS require resolution and caching.
- [ ] macOS/Windows providers and omission of Linux runtime dependencies when unused.
- [ ] Differential coverage of the full inventory and applicable Bun/Node compatibility tests.

## Verification

`tsc/internal/compiler/gofs_test.go` checks emitted code, not just runtime helpers:

- sync/callback/promise I/O, binary data, metadata, directories, links, copying, errors;
- early directory iterator close, glob iteration, filesystem calls on user worker loops;
- descriptor offset/vector behavior and callback buffer identity;
- closing a FileHandle with outstanding I/O;
- native watchers, stat polling and unreferenced watcher shutdown;
- streams, pipe backpressure, ranged reads, and multibyte text split across one-byte chunks;
- a shared fixture compared directly with Node when Node is installed.

The generated-program checks use Go's race detector. Runtime checks cover the
value ABI and existing worker/task behavior. Run from the repository root:

```sh
go test ./tsc/internal/goemit -run '^TestRuntime(InterfaceBoundaries|Workers)$'
go test ./tsc/internal/compiler -run '^TestGoNodeFS'
```

Passing these checks establishes the listed paths, not full module conformance.
