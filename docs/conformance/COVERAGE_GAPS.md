# Compatibility gaps — 2026-10-10

The completed latest standard-library sample passed 71 of 95 variants, with
10 emission failures, 6 runtime failures and 8 unavailable host cases. This is
50 selected source tests, not a random sample or a full compatibility score.
See `features-stdlib-audit/summary.json` for the exact compiler source hash.

The completed older full run passed 14,742 of 100,058 variants. Its binary
predates the newer Error, Symbol, Math, arguments, deletion, computed-class and
generator work. Do not use its remaining failures as current feature counts.
It provides a backlog inventory until a new full run is complete.

| Area | Evidence and remaining work |
| --- | --- |
| Errors | SuppressedError construction, descriptors and subclassing now have focused coverage; other error edge cases remain. |
| Generators | Delegated throw getter failures now enter body catch/finally; broader iterator edge cases still require auditing. |
| Function metadata | Anonymous/default function-name inference now has runtime coverage; broader metadata and constructor cases remain. |
| Classes | Private fields/methods and dynamic field mutation now have focused coverage. Dynamic heritage, constructor metadata and newTarget remain substantial gaps, with an upstream class reflection case still failing. |
| Early errors | Strict `arguments` writes, generator parameter defaults containing `yield`, and private identifiers outside classes now fail before emission; broader grammar early errors remain. |
| Globals | Sloppy creation/deletion of implicit identifier bindings now has focused coverage; cross-module and strict-mode binding edges remain. |
| Arrays and objects | Sparse/exotic behavior and Reflect/newTarget behavior still need broader coverage. Native typed storage remains separate. |
| ArrayBuffer and DataView | DataView reads/writes, `ArrayBuffer.isView`, detach checks, transfer-to-fixed-length and resizable `maxByteLength`/`resize` now have focused coverage. Growable view tracking and broader detached-host cases remain. |
| RegExp and strings | Coercion and legacy/exotic cases remain after focused well-known Symbol protocol coverage. |
| Temporal and Intl | Entire large categories fail in the older inventory; require dedicated native implementations. |
| Modules and host | Module graph, realm/agent, GC and detachment adapters are unavailable in this runner; Node compatibility also remains separate work. |

JSX and runtime code generation via eval/new Function are excluded from the
project goal. Unsupported host cases are not successes. No new performance
benchmark is inferred from conformance results.

Keep native typed calls, storage and unaffected class field loads direct.
Add dynamic metadata or protocol machinery only for operations that need it.
Every new compatibility implementation requires semantic tests; generated-Go
checks protect the unaffected native paths.

## Terms used in this work

Private class members are fields or methods whose names begin with `#`, such as
`#balance` or `#deposit()`. They are language-enforced capabilities: code
outside the declaring class cannot read, write or even substitute a string
property with that name. The emitter stores them under an internal key and
uses concrete class access inside the class, so ordinary public property lookup
does not expose them.

Dynamic field mutation means a field can be reached through JavaScript's
runtime property machinery, for example `obj[key] = value`,
`Object.assign(obj, source)`, or a call to an unknown function that receives the
object. A statically known `number` field can remain a native Go number only
when analysis proves those operations cannot change its value domain. Otherwise
the field uses the value representation and performs the normal conversion at
the boundary. This is deliberately a compatibility path; closed typed classes
continue using native storage.

Early errors are source-level violations that JavaScript rejects before running
the program. Examples include assigning to `arguments` in strict mode and
putting `yield` in a generator parameter default. They are different from a
runtime exception: the emitter reports them before writing Go output, so no
partially generated program can run.

## Changes after the latest sample

Focused tests pass for SuppressedError construction, descriptor attributes,
subclassing and ignored cause options; sync delegated iterator getter failures
entering body catch/finally; anonymous/default function-name inference; and
class method/field descriptor identity. Native direct-field checks still pass.
The fresh `next-gaps-audit/` reruns five specific upstream failures against a
new emitter snapshot, so its results must remain separate from the 95-variant
sample.

## Targeted upstream rerun completed

`next-gaps-audit/` completed all 10 variants: 8 passed and 2 failed at runtime.
Previously all 10 failed. The remaining class test uses propertyHelper to write
a string into an unannotated JavaScript field inferred as numeric. Native
conversion throws RangeError, exposing unsound storage inference for reflected
JavaScript fields. Fixing this requires field escape/write analysis or storage
promotion, while preserving explicitly typed native storage and proven static
loads. Do not weaken checked numeric conversions to mask this failure.

Readable runtime error reporting now has a passing test; error coercion occurs
after releasing promise locks. The latest five-test validation passed for that
reporting, SuppressedError, delegated getters, inferred function names and class
descriptor identity. The targeted audit predates the diagnostic-only change.

The current implementation also accepts private instance/static fields and
methods, stores private names under internal keys, and keeps in-class access on
concrete class paths. Early-error validation runs before Go emission for strict
`arguments` writes, updates/deletes, generator parameter defaults containing
`yield`, and private identifiers outside class bodies. Sloppy `arguments`
assignment remains valid. A dynamic-field regression confirms that an unknown
function can change an unannotated class field from a number to a string; that
field uses value storage on the reflective path.

The compatibility audit now also covers generic `Array.prototype` receivers,
array data/accessor descriptors, `globalThis` own-property behavior, custom
RegExp symbol protocols, the Float16/BigInt typed-array families, DataView
endianness and BigInt access, resizable ArrayBuffers, and sloppy implicit
globals. Typed arrays keep their native element storage; BigInt elements remain
exact values instead of passing through float64. Sparse exotic objects,
growable-view tracking, and broader host/strict binding cases remain.
