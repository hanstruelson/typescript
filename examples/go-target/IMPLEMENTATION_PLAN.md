# Go backend implementation plan

This is a roadmap, not a claim that these features already work. The README
records the current subset. Each phase should land with focused behavioral
comparisons against JavaScript and structural checks of the generated Go.

## Performance and representations

The current production representation is Go's `any`. On 64-bit Go, an interface
occupies two words (16 bytes): type information and data. This already handles
references with garbage-collector visibility. It is not a nine-byte tagged union. The [representation benchmarks](benchmarks/README.md)
compare it with custom containers. The recommended starting point for a custom
implementation is an inline tag/native numeric payload plus checked Go-managed
reference storage (32 bytes on amd64). Migration is paused for this experiment;
the production ABI has not changed.
Nullable primitive values have a native payload and a null/undefined/value tag;
Go alignment means their actual size depends on the payload. Typed primitive
bindings have additional initialization and const tracking.

The optimization order is:

1. Keep statically proven primitive expressions in native Go temporaries,
   parameters, returns, and fields. Extend the native calling convention to
   eligible class methods and imported functions. Preserve a dynamic adapter for
   indirect calls. Resolve method implementations and field slots at compile time
   wherever possible.
2. Separate storage type, flow-narrowed expression type, and boundary policy in a
   small backend representation plan. Derive facts through the checker resolver;
   do not infer safety from source spelling or erase null/undefined distinctions.
3. Avoid allocating cells for uncaptured locals. Captured mutable bindings and
   suspended state-machine locals need persistent storage; ordinary locals do
   not. Preserve TDZ, per-iteration bindings, and closure lifetime.
4. Specialize arrays whose element representation is stable. Use native slices
   for proven dense typed arrays, with adapters for dynamic views. Mutation
   through an alias must not silently corrupt an optimized representation.
5. Benchmark dynamic numeric loops, object property access, function calls,
   heterogeneous arrays, and async continuations. Measure allocations, throughput,
   memory, and output size. Use `any` as the compatibility baseline while evaluating the custom value;
   its numeric allocation benefit must survive complete runtime calls and storage
   paths. A custom tag alone does not guarantee faster code.

Never hide Go pointers in integer payload bits without a complete GC-safe
ownership scheme. NaN boxing and handle-table designs would be separate research,
not prerequisites for this backend.

## One typed-boundary policy

Today typed parameters, assignments, initializers, and returns have runtime
checks or conversions. Comparisons still follow the implemented JavaScript
operator helpers; mixed-any comparison checks described below are planned.

Build one conversion/check planner used by every lowering site: declarations,
assignments, parameter defaults, function calls, return values, fields, bracket
setters, destructuring leaves, and mixed dynamic/typed operators. Evaluate source
operands once, left to right, then check or convert saved values. Primitive-to-
primitive operations with proven representations bypass dynamic checks.

For a mixed `any`/typed primitive comparison, apply the requested policy to the
`any` operand: default checks its runtime type; `coerceAny` converts it. This
includes `==`, `!=`, `===`, `!==`, and ordering comparisons. Nullable destinations
retain their allowed tags. Comparing against null/undefined sentinels needs its
own tag rule rather than an arbitrary primitive conversion. Dynamic arithmetic
at a typed operand boundary uses the same planner. Short-circuit operators retain
lazy evaluation and do not eagerly convert operands that are never evaluated.

This deliberately changes ordinary JavaScript equality: normally `10 === "10"`
is false, not a conversion or an exception. Fully dynamic operators must still
implement ECMAScript's distinct strict equality, abstract equality, relational
comparison, and ToPrimitive rules. Imported JavaScript npm code needs an explicit
JavaScript-compatible compilation policy, while typed application code can use
these checked/coercing boundaries. Applying the custom comparison policy to all
npm source would break common type-discrimination code. Carry the selected policy
per source/module and in build-cache keys rather than switching a process-global
runtime flag.

Tests must cover both operand orders, null/undefined, booleans, NaN, negative zero,
objects with conversion hooks, nullable destinations, coercion failures, evaluation
order, and side effects. A thrown conversion error must be a catchable TypeError
with the option hint. Check native output separately from behavioral correctness.

## Class work

| Feature | Implementation | Critical verification |
| --- | --- | --- |
| Static fields and methods | Implemented: a separate generated static struct, one shared object per evaluated class, native primitive fields, inherited reads through parent storage, and child-owned writes. Methods specialize on the receiving class while retaining lexical `super`. | Shared state across instances; parent mutation; child shadowing; four-level super; bracket access; async methods; initialization once. |
| Static blocks | Execute generated block bodies interleaved with field initializers in source order. Each block gets its own locals and can reuse structured or machine lowering. | Multiple blocks, scoped bindings, loops, exceptions, inheritance, and initialization side effects. |
| Accessors and descriptors | Add property descriptors with optional getter/setter, writable/enumerable/configurable attributes. Known accessors become direct Go calls; dynamic access uses descriptors. `super` chooses the lexical owner's descriptor but passes the current receiver. | Getter/setter overrides, getter-only writes, side effects, descriptors, receiver identity, and multi-level super. |
| ECMAScript private fields/methods | Resolve private names by declaration identity, emit distinct slots for every declaring class, and track brands separately from ordinary keys. Private accesses bypass hashing. Keep static private brands attached to their declaring singleton. | Parent and child `#x` names, brand errors on wrong receivers, private methods, and static private inheritance. |
| Parameter properties | Create field slots from constructor parameter declarations and initialize them at TypeScript's defined point after base construction. Reuse parameter boundary checks. | Defaults, field initializer order, derived constructors, nullable and dynamic parameters. |
| Computed members | Evaluate computed keys once during class definition, in source order; store keys in class metadata. Use direct slots when a key is statically fixed and dynamic descriptors otherwise. | Computed side effects, symbol keys, duplicate keys, and initializer ordering. |
| Cross-file inheritance | Build a bundle-wide class graph using resolved declaration identities before preparing layouts. Generate all specializations in a coordinated pass; keep module execution order separate from layout planning. | Imported/re-exported bases, duplicate names, three-file inheritance, and module cycles. |
| Class expressions/dynamic bases | Evaluate a base constructor once. Use static specialization for known bases and a dynamic class/prototype representation for runtime-selected bases. Preserve named class-expression scope. | Anonymous and named expressions, factories, conditional bases, and lexical name isolation. |
| Rest/destructured method parameters | Share parameter lowering with ordinary functions instead of retaining a second restricted path. | Defaults referencing earlier parameters, TDZ, rest, nested patterns, and super constructors. |
| Method receiver behavior | Give dynamic calls an explicit receiver argument. Method lookup returns the method identity; a member call supplies its receiver, an extracted function does not automatically bind. Arrow closures capture lexical `this`. | `const f=obj.method`, call/apply/bind, inherited methods, and strict-mode undefined receivers. |

A singleton is per evaluated class value. Re-evaluating a class declaration in a
factory or loop creates a new class and new static storage, as JavaScript does;
instances of the same evaluated class share its singleton.

## Objects, arrays, and syntax

| Limitation | Implementation | Critical verification |
| --- | --- | --- |
| Object methods/accessors | Reuse explicit receiver calls and descriptors. Lower method bodies through existing function machinery. Preserve concise-method lexical super via a home-object reference. | Extracted methods, getter/setter receiver, prototype super, enumeration. |
| Prototypes/dynamic properties | Add prototype links and descriptor lookup/set algorithms. Known class slots stay direct. Prototype mutation and dynamic overrides require guards or fallback paths, not stale direct dispatch. | Object.create/setPrototypeOf, inherited setters, property deletion, hasOwn vs `in`, for-in. |
| Symbols/property keys | Use a property-key type distinguishing UTF-16 string keys and unique symbol identities. Intern string keys for repeated dynamic access. Keep direct typed fields independent of this lookup. | Distinct same-description symbols, well-known symbols, key ordering, lone surrogates. |
| Sparse arrays | Track length separately from storage, presence with a bitmap for dense regions, and a sparse map for remote indices. Holes differ from explicit undefined. Preserve insertion order for non-index keys. | Holes, delete, large indices, length truncation, keys, map/filter/reduce, and iterator behavior. |
| Array/collection methods | Add ECMAScript Array, Map, Set, WeakMap, and WeakSet APIs over native storage. Use SameValueZero for collection keys and preserve order. Typed-array views need buffer ownership, bounds, and aliasing rules. | NaN/-0 keys, mutation during iteration, callbacks, holes, buffer views, weak-key lifetime. |
| Destructuring assignments | Generalize pattern lowering to lvalues. Snapshot RHS, computed keys, and targets in specification order; handle defaults, rest, and iterator closing on abrupt completion. | Getters, aliasing, nested assignments, skipped elements, throwing iterators. |
| Optional chaining | Lower through branch blocks and saved receivers. Skip arguments/computed keys when nullish and preserve member-call receivers. Optional delete has separate semantics. | Nested chains, calls, computed-key side effects, and grouping boundaries. |
| Missing operators | Add ToPrimitive and separate operator algorithms, then support loose equality, exponentiation, bitwise shifts, logical assignments, `in`, delete, and void. Use native arithmetic only when the representation plan proves it safe. | Signed/unsigned 32-bit conversion, infinities, NaN, string/object conversion order, lazy assignments. |
| Generators | Extend the existing resumable machine with yield results and next/throw/return inputs. Preserve protected-region unwinding and iterator closing. Implement yield-star delegation explicitly. | Yield inside loops/try/finally, sent values, delegation, abrupt close. |
| Async iterators/generators | Add the async iterator protocol and serialize each generator's pending requests through the event loop. Lower for-await-of with awaited next/return operations. | Sync-iterator adaptation, rejected next, early exit, finally, queued requests. |
| Named function expressions | Create an internal self-binding scoped to the function, distinct from its surrounding variable. Extend function metadata with name/length and the explicit receiver ABI. | Recursion after outer reassignment, scope isolation, call/apply/bind. |
| using/await using | Reuse TypeScript resource-management lowering or translate its disposal stack into the Go runtime. Preserve reverse disposal and suppressed-error semantics. | Normal/abrupt exits, multiple errors, awaited disposal, nested scopes. |
| Namespaces/enums/import-equals | Reuse established TypeScript erasure/lowering concepts. Emit persistent namespace objects, correct numeric enum reverse mappings, and resolved import aliases. | Declaration merging, enum initialization order, const enums, ambient erasure. |
| BigInt | Add an immutable BigInt wrapper over math/big, with operators, conversions, and separate dynamic tags. Never silently convert BigInt arithmetic to float64. | Precision, signed division/remainder, mixing with Number, equality, serialization errors. |
| Decorators | Reuse TypeScript's decorator lowering and implement the helper contract, including initializer lists and metadata. | Evaluation/application order, field/method replacement, static/instance initializers. |
| JSX | Apply the configured JSX transform before Go emission and implement imports/calls for the selected runtime. Keep preserve-mode JSX unsupported until an output contract exists. | Classic vs automatic factories, spread props, fragments, evaluation order. |
| Source maps | Track generated spans against original AST positions. Emit mapping data for diagnostics/debug tooling; optionally use Go line directives where useful without corrupting generated-file tooling. | Multi-file bundles, lowered awaits, synthetic helpers, source locations. |

## Modules and npm

Retain the requested demand-driven module policy for application modules. Expand
export-star by building export-name tables from the normal resolver, excluding
default and detecting ambiguous re-exports. Keep live binding cells and initialize
one persistent module machine per evaluated module. Add top-level await through
module promises and detect unresolved cyclic waits.

npm support additionally needs:

- JavaScript source inclusion, package exports/imports and conditional resolution,
  using the repository resolver rather than a parallel ad hoc resolver.
- CommonJS require/module/exports with a cache populated before executing a
  module, partial exports during cycles, and eager CommonJS execution. Its
  semantics cannot be replaced by demand-driven named-import execution.
- A host API layer for Buffer, process, timers, paths, filesystem, networking,
  streams, events, URLs, and supported crypto. Workers receive immutable snapshots
  or owned buffers, perform native work, and report results through the existing
  task helper; only the event loop runs JavaScript-visible callbacks.
- Promise API completion, rejection reporting policy, cancellation/resource
  cleanup, and structured catchable errors. Recover only at owned execution
  boundaries; panic recovery cannot safely fix arbitrary memory corruption.
- Explicit support declarations for native addons and packages requiring eval or
  Function. Dynamic source compilation remains excluded by project policy.

Validate representative real packages incrementally. Being dynamically typed is
only one requirement; language semantics, module behavior, and host APIs all
matter before claiming npm compatibility.

## Strings, regex, and output

Keep immutable UTF-16 storage and direct code-unit algorithms. Complete generic
String receiver conversion, object ToPrimitive hooks, constructor/wrapper behavior,
Symbol.match/replace/search/split hooks, locale option validation, and resource
limits. Compare against targeted Test262 cases, including isolated surrogates.
For full Intl behavior, select a versioned CLDR/ICU-compatible implementation;
x/text collation alone does not cover the complete Intl API.

Run regex conformance tests against the selected Goja/regexp2 versions, including
Unicode properties, advanced flags, captures, replacement protocols, lastIndex,
and invalid patterns. Implement or diagnose unsupported features explicitly.
RegExp work stays on the event loop unless an immutable operation can safely be
isolated without observable mutation.

Generate a small always-included runtime core plus feature-dependent helpers and
imports. Track helper dependencies so unused regex/collection/class machinery can
be omitted. Keep one output Go file initially; document and pin external module
dependencies. Later allow an optional generated runtime package without changing
application semantics.

## Delivery order and acceptance

1. Static objects and initialization, with targeted class and module regressions.
2. Explicit receiver ABI, descriptors, prototype lookup, and complete boundary
   planning. These are shared foundations for classes and npm code.
3. Operators, optional chains, sparse arrays, destructuring assignment, and core
   collection APIs. Add measured native optimizations alongside each feature.
4. Private/accessor/computed classes, cross-file layouts, and shared parameter
   lowering; then generators and async iteration.
5. CommonJS and Node host adapters, using representative npm packages as fixtures.
6. Remaining TypeScript lowering, Intl/regex conformance, debugging, and output
   size/startup improvements.

Use a tracked feature matrix: native implementation, dynamic fallback, explicitly
unsupported, or verified complete for the tested contract. Test each semantic
addition against JavaScript, use race checks for worker/event-loop interaction,
and benchmark only the relevant hot paths. Avoid repeatedly running unrelated
large suites. Full coverage is a measurable roadmap, not a completed promise.
