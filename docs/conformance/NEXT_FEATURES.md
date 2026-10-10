# Next compatibility milestones

This is an implementation plan, not a claim of completed support. Scope:
Error constructors, Math, Symbol/property keys, delete, computed class members,
ordinary-function receivers and arguments, and synchronous/asynchronous generators.
Runtime code generation through eval/Function and JSX remain outside scope.

## Shared rules

- Reuse the TypeScript parser and AST. These constructs are already parsed;
  implementation is primarily in Go emission, analysis, and runtime support.
- Known numeric values, class fields, and direct calls retain native Go storage
  and operations. Dynamic values use the existing 24-byte TSValue, not Go any.
- Preserve JavaScript evaluation order, coercion, strictness, property descriptors,
  inheritance, and observable identity. Type annotations alone do not prove that
  reflection or deletion is impossible.
- Analyze aliases, inheritance, escapes, and reflective mutations. When a receiver
  cannot be proven closed, select prototype-aware property operations. Refine the
  current source-wide analysis toward class- and field-specific effects.
- Reuse common semantics across syntax and reflective APIs; avoid separate rules
  for Object, Reflect, assignment, deletion, and class fields.

## 1. Errors and Math

- Expose stable Error, TypeError, RangeError, ReferenceError, SyntaxError,
  URIError, EvalError, and AggregateError constructors, including call/new,
  prototypes, instanceof, subclassing, name/message, cause and descriptors.
- Unify native runtime failures with these public error objects. Stack reporting
  needs a documented native implementation; do not promise engine-identical stacks.
- Implement the complete standard Math function/constant inventory, including
  coercion, signed zero, NaN, infinities, overflow, imul/clz32/fround, hypot,
  random, and standardized additions supported by the pinned suite.
- Direct known Math calls can lower to Go math or specialized helpers only when
  name binding and mutation analysis prove the intrinsic unchanged. Aliasing,
  extraction, rebinding and mutation use the same stable public Math object.

## 2. Property keys and Symbol

- Introduce an internal property-key representation distinguishing UTF-16 string
  names from identity-bearing symbols. Never stringify symbols into map keys.
- Cover unique identity, description, Symbol.for/keyFor, well-known symbols,
  typeof, explicit string conversion, and forbidden implicit conversions.
- Carry keys through objects, descriptors, prototypes, classes, membership,
  deletion, assignment and reflection. Retain existing fast string-key access.
- Implement Object.getOwnPropertySymbols and Reflect.ownKeys ordering; exclude
  symbol keys from Object.keys, for-in and JSON.stringify as specified.
- Connect well-known-symbol protocols to iteration, instanceof, coercion and
  RegExp operations. Sharing an identity does not require locks on every lookup;
  registry insertion does require synchronization across user goroutines.

## 3. Delete and selective field presence

- A deleted property is absent, not present with value undefined. Reads can reveal
  an inherited property; in, hasOwn, descriptors, keys and JSON must agree.
- For a proven literal-key delete on a known class, mark only affected fields as
  potentially absent. Keep each field's native value plus presence information;
  do not box every field or confuse absence with null/undefined.
- Dynamic keys or unknown reflective mutation conservatively mark reachable
  fields. A dynamic key affects its possible receiver, not every class globally.
- A deleted native pointer/reference field must release its stored reference.
  Reassignment restores presence, subject to extensibility and descriptors.
- Preserve own-property configurability: deletion returns false in sloppy mode
  or throws TypeError in strict mode when forbidden. Missing properties succeed.
- Evaluate receivers/keys once in source order; implement optional-chain delete,
  primitive receivers, and non-property delete expressions. Respect early errors
  for strict identifier deletion and private-field deletion.
- General JavaScript arrays must preserve sparse deletion semantics and length.
  Dense native array lowering is eligible only when the program cannot observe
  holes; deletion can require a sparse-capable representation for that receiver.
- Reflect.deleteProperty shares the primitive delete operation but returns its
  boolean result rather than applying the syntax's strict-mode throwing rule.

## 4. Computed class members

- Support fields, methods and accessors such as [name] and [Symbol.iterator],
  including static members, duplicate keys and descriptor rules.
- Evaluate and convert computed names once during class definition, in specified
  order. Instance field initializers run for each constructed instance; computed
  names do not. Preserve static initialization order, TDZ and derived/super rules.
- Proven constant string keys can use native fields. Truly dynamic/symbol keys
  use keyed storage for those members without forcing unrelated fields to box.
- Integrate deletion/presence, inheritance and reflection before optimizing.

## 5. Ordinary-function this and arguments

- Pass a receiver only where needed, preserving existing typed direct-call ABIs
  for functions that do not observe this or arguments. Distinguish strict/sloppy
  receivers, arrow lexical capture, call/apply/bind and constructor calls.
- Implement arguments as an arguments object, not an ordinary Array: length,
  indexed properties, iteration, descriptors and strict callee restrictions.
- Sloppy functions with simple parameter lists need mapped parameter aliasing;
  strict functions and non-simple parameter lists use unmapped arguments.
  Deletion/redefinition can break an index's mapping. Cover duplicate parameters,
  extra/missing arguments, defaults/rest/destructuring and arrow capture.
- Allocate the observable object only when required. Direct arguments.length
  and arguments[index] operations can use call-frame data when equivalence is
  proven; observing arguments must not silently discard extra supplied values.

## 6. Generators and iterator protocols

- Extend the existing state-machine approach with caller-driven suspension.
  Async tasks provide useful plumbing but do not already implement generators.
- Implement lazy invocation, yield/resume values, next/throw/return, completion,
  reentrancy errors, nested try/finally and yields during finally execution.
- Implement yield* delegation with iterator closing and abrupt completions;
  connect Symbol.iterator, for-of, spread and destructuring to common protocols.
- Keep statically known frame locals native; box values only at dynamic protocol
  boundaries. Synchronous generators need no goroutine/channel per iterator.
- Then implement async generators, queued requests, awaited yielded values,
  Symbol.asyncIterator and for-await-of using the owning event loop.

## Verification and completion gates

For each milestone, inventory the pinned Test262 categories and public APIs;
record supported, unsupported, failing and runner-limited behavior separately.
Add focused Node differential tests for normal and abrupt behavior, then run
the complete relevant Test262 categories and existing regression tests.

Inspect generated Go and benchmark unaffected typed calls/class-field loops:
these must retain native access without new boxing, allocations, or presence
checks on fields proven undeletable. Benchmark dynamic behavior separately.
Use race tests for shared initialization/registries, not a claim that arbitrary
user mutations are race-safe.

Keep the existing full-suite run tied to its original source hash. Build a new
driver and use a separate report directory after implementation changes; never
combine old and new results into a single compatibility score. Resource kills,
timeouts and unavailable host facilities remain distinct from semantic failures.

Declare a milestone complete only after its inventory, semantic tests and
performance checks pass. These milestones do not by themselves finish all
JavaScript compatibility: Proxy and remaining built-in/host/module behavior
still require separate work.

## Implementation status — 2026-10-09

The emitter now has native implementations of the Math inventory (including
f16round and sumPrecise), public Error constructors and native Error subclasses,
Symbol identities/property keys, property deletion, class expressions, computed class fields,
methods and accessors, ordinary function receivers, mapped/unmapped arguments,
and synchronous/asynchronous generator state machines. This is implementation
coverage, not a claim of complete ECMAScript conformance.

Focused Node differential tests cover coercion and descriptors; nullable versus
native numeric operations; selected field deletion and unaffected native loads;
computed-name evaluation, colliding field names and class-name TDZ; receiver
capture, arguments defaults/mapping/deletion; Symbol copying and reflection;
and generator next/throw/return, delegation, cleanup and async-from-sync iteration.
A generated-runtime GC test checks retention of Symbol-only property keys.

The first resumed upstream audit is pinned in `features-resumed-audit/`: 72
variants, 28 passes, 18 build failures, 11 emission failures, 7 runtime failures,
and 8 host-adapter limitations. All 18 build failures exposed the same nullable
array-index lowering issue in unmodified upstream property helpers. Follow-up
changes fix that lowering, generator function prototype roots, delegated result
identity/getter timing, arguments binding deletion and constructor validation.
Class expressions and receiver-aware Reflect get/set/defineProperty now have
passing Node differential checks as well. A new snapshot must be audited
separately; those fixes do not change the recorded
results of the earlier binary.

Remaining completion work includes private class members, complete class and
constructor metadata/newTarget behavior, RegExp well-known-symbol protocols,
Reflect exotic-object behavior and setPrototypeOf failure semantics, sparse deletion
for native arrays, full global binding behavior, and
additional generator/iterator edge cases. The full-suite run in
`prototype-full/` uses the earlier prototype binary and must not be presented
as validation of these newer changes. JSX and runtime code generation remain
excluded; host/realm and module runner gaps remain separately classified.

Ordinary typed fields proven unaffected by deletion continue to emit native Go
loads. Known nonnullable numeric array indices and unobserved arguments retain
their native paths. Dynamic reflection and nullable indices require their
corresponding operations. No new benchmark speedup is claimed for this work.

### Shared built-in roots and stack accessors

A separate follow-up snapshot (`features-followup-audit/`) completed 95 sampled
variants: 53 passed, 24 failed at runtime, 10 failed emission, and 8 required
unimplemented host facilities. There were no build failures. This sample is
not the full suite and has a different selection from the 72-variant audit.

Subsequent fixes add shared Array.prototype method identities, extracted Array
static methods, propertyIsEnumerable, authoritative function metadata after
deletion, and non-enumerable global built-in descriptors. Array construction
uses synchronized lazy initialization. Error.prototype.stack now implements
the accessor contract in the pinned suite, with lazy, implementation-defined
Go frame formatting; this is not an engine-identical Node stack format.
Focused differential checks and generated-runtime worker/interface checks pass.

The `features-stdlib-audit/` snapshot reruns the same 95-variant selection after
those changes. Its results are tied to its own emitter hash. Later prototype
mutation changes are tested separately: Reflect.setPrototypeOf returns false
for cycles and non-extensible ordinary objects, functions and native instances,
while Object.setPrototypeOf throws on failure. Arrays accept prototype mutation.
Array descriptor definition and generic array method receivers still need work.

The current upstream sample also exposes an early-error gap: strict-mode
assignment to `arguments` is not rejected during syntactic validation. The
emitter subsequently rejects the negative-test sentinel `$DONOTEVALUATE`,
which is recorded as an emission failure rather than an expected SyntaxError.
That diagnostic does not count as a conformance pass.

Final targeted checks passed for prototype mutation/Object.create and for eight
array/class regression groups, including native nullable arrays, callback
mutation/control flow, named nullable callbacks, unaffected native field loads,
and prototype-free direct class access. The regression run completed in
184.985 seconds; no benchmark performance claim is inferred from test timing.

### Completed standard-library sample — 2026-10-10

The `features-stdlib-audit/` snapshot completed all 95 planned variants: 71
passed, 10 failed emission, 6 failed at runtime, and 8 required unavailable host
facilities. The preceding sample of the same selection had 53 passes. No build
failures occurred in either of these two snapshots. This is sampled coverage,
not full-suite compatibility; the snapshot predates the final prototype mutation
and Object.create fixes, which passed separate differential tests.
