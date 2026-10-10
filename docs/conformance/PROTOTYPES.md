# Prototype implementation

JavaScript objects have an internal prototype link, used when a property is not
found on the object itself. A constructor's public `.prototype` property is the
object installed as that internal link on newly constructed instances. These
are different things: `Object.getPrototypeOf(instance)` reads the internal link;
`Constructor.prototype` reads an ordinary property on the constructor.

Class inheritance also uses prototypes. They matter even when application code
never spells `.prototype`, and libraries often expose their behavior through
reflection. Replacing a function's `.prototype` affects future instances;
instances already constructed retain the old object. Mutating that shared object
is visible to all instances still linked to it.

## Current strategy

- Closed, non-reflective classes retain concrete Go fields and direct method
  calls. There is no new prototype-enabled branch on those generated accesses.
- A conservative source analysis chooses dynamic class dispatch for reflection,
  imports/exports, computed accesses, unknown calls, method extraction, and property writes that could replace methods.
  It is currently source-wide, not a complete per-class escape analysis.
- Typed numeric functions and callbacks retain their existing native ABI even
  in reflective source. Their boxed function objects are constructible when the
  original JavaScript function is constructible; arrows are not.
- Constructor prototypes are shared, lazily materialized objects. They are not
  copied into every instance. Native class fields remain backed by Go storage;
  prototype-aware reads use their existing field lenses.
- Class prototype methods are shared functions that receive the actual receiver.
  Prototype overrides and lexical `super` lookups do not use a stale bound method.
- `sync.Once` protects prototype initialization. Lazy function metadata is
  published with an atomic pointer, avoiding races between concurrent readers.
  This adds no atomic operation to direct numeric calls or closed-class fields.
  User mutations remain subject to the project's existing shared-memory rules.

The implementation covers ordinary constructor functions, shared prototype
identity, replacement versus mutation, constructor object/primitive returns,
bound construction and instanceof, Object/Function prototype roots, inherited
getters/setters, membership, inherited enumeration, `Object.create`,
`getPrototypeOf`, `setPrototypeOf`, cycle checks, legacy `__proto__`, and reflected
class method overrides. Object-literal prototype syntax is distinguished from a
computed data property named `__proto__`.

Function names, lengths and prototype descriptors are represented as metadata;
ordinary metadata is allocated only when needed. Strict comparisons with dynamic
values no longer coerce a string merely because the other operand is numeric.

## Validation and remaining work

Node differential tests exercise constructor/prototype aliases, old and new
instances, bound functions, inheritance, descriptors, enumeration, class overrides
and super. A generated-code test checks that a closed class still uses a direct
Go call. A race-enabled worker test checks concurrent reflection of one shared
constructor prototype. Existing class/static/primitive-field, Date and JSON
regressions were also exercised.

The focused Test262 Date.now audit improved from 0/12 to 8/12 passing variants;
the remaining four stop in unsupported property-verification helpers. This is
not a full compatibility claim.

Remaining prototype work includes complete reflected class metadata and primitive/array/collection prototype
surfaces, class-instance descriptor operations, Proxy and Symbol.hasInstance
behavior, full function name inference and source-preserving toString, and exact
strict-versus-sloppy write semantics. The conservative analysis can be refined
with per-class escape/effect information once those semantics are stable. Dynamic
Function/eval compilation is still outside the current runtime.

The restarted full-suite run uses a separate output directory and records a hash
of the native sources, including uncommitted changes. Old results are retained;
results from the two implementations are not combined.
