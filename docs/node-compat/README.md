# Node module compatibility

Implement modules against their public behavior, using Bun as a reference and
Go packages for the implementation. Bun is not linked into generated programs.
Do not translate its repository file by file: its engine bindings, allocation,
JSC promises, and native resource management are specific to its runtime.

The reference checkout is a shallow, sparse clone at `../bun`, pinned to
`bbdc5a519e0a06d1b3133b564f91096b9ba10a31`. It includes `src`, `packages/bun-types`,
and `test/js/node/fs`. Bun's source is MIT licensed; retain attribution and its
license for any future code copied or translated directly. The implementation
here uses Bun's APIs and behavior as a reference and Go's filesystem packages.

Shared JavaScript built-ins are the current priority; module expansion is paused.
See the [built-in implementation checklist](builtins.md) and [API inventory](builtins-api.json).

## Module order

| Module | Status | Next gate |
| --- | --- | --- |
| `fs`, `fs/promises` | First implementation; incomplete | Close every item in [the filesystem checklist](fs.md) |
| `buffer` | Supporting subset for filesystem byte operations | Full Buffer API and constructor/prototype compatibility |
| `events`, `stream` | Filesystem-local plumbing only | Public modules, inheritance, complete stream contract |
| `http`, `https` | Not started | Complete `fs` and shared prerequisites first |
| npm resolution, Express | Deferred | Built-ins and CommonJS module execution |

For each module:

1. Pin reference revisions and enumerate every public export and overload.
2. Map public behavior to Go operations and identify shared runtime dependencies.
3. Implement sync, callback, and promise forms from a common native core.
4. Keep application callbacks on their owning event loop. Native workers perform I/O.
5. Record unsupported behavior explicitly; export availability is not completeness.
6. Verify module-wide fixtures, Node differential behavior, error cases, and resource lifetime.
7. Add the module to the completed list only when the outstanding checklist is empty.

The current filesystem provider targets Linux (`golang.org/x/sys/unix`, inotify,
statfs, vectored I/O, and clone ioctls). Platform adapters are required before
claiming macOS or Windows support. The embedded runtime currently includes this
provider in generated Go files, so this also limits generated-file portability.
