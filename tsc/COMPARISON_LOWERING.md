# Comparison and operator lowering matrix

Status: semantic inventory of implemented lowering, the shared plan rules, and the remaining measured optimization boundaries.

## Scope

The binary comparison operators are partitioned into three semantic operations:

| Operation | Operators | Shared result rule |
| --- | --- | --- |
| Loose equality | `==`, `!=` | Compute equality once; `!=` negates it. |
| Strict equality | `===`, `!==` | Compute equality once; `!==` negates it. |
| Relational | `<`, `<=`, `>`, `>=` | Compute an ordered result or unordered; map it to the operator. Unordered comparisons are false. |

JavaScript has no strict relational operators: there is no strict `<`, `>`, `<=`, or `>=`. Strictness applies to `===` and `!==`. The remaining operators in this document are adjacent expression families with different rules, not extra forms of comparison:

| Family | Operators | Core behavior |
| --- | --- | --- |
| Logical negation | `!` | Convert the operand by truthiness and return the opposite boolean. |
| Short-circuit logical | `&&`, `||`, `??` | Conditionally evaluate the right side and return one of the original operand values. |
| Logical assignment | `&&=`, `||=`, `??=` | Apply the corresponding short-circuit rule, then assign only when the rule selects the right value. |
| Bitwise unary | `~` | Convert according to the Number, typed integer, or BigInt family, then complement. |
| Bitwise binary | `&`, `|`, `^`, `<<`, `>>`, `>>>` | Use Number 32-bit rules, typed integer rules, or BigInt rules, according to the operand family. |

The property-membership operators `in` and `instanceof` use property/prototype semantics and belong in a separate operator plan; they are not primitive comparison pairs covered by this matrix.

Several built-ins reuse equality-like relations with different Number edge cases:

| Consumer | Equality relation | Difference from `===` |
| --- | --- | --- |
| `switch` case matching, `Array#indexOf`, `Array#lastIndexOf` | Strict equality | `NaN` does not match itself; `+0` and `-0` match. |
| `Object.is` | SameValue | `NaN` matches itself; `+0` and `-0` differ. |
| `Array#includes`, `Map`, `Set` keys | SameValueZero | `NaN` matches itself; `+0` and `-0` match. |

For nonnumeric values, these relations all use the same equality: primitive values compare by their strict value rule, Symbols and references by identity. Reuse the nonnumeric strict payload logic, but carry an equality flavor for Number edge cases. Array hole behavior and collection key normalization remain separate from the equality relation.

String content comparison is necessarily linear in the worst case for distinct strings with equal contents or a long common prefix. The current helper makes alias comparisons constant time by checking pointer identity, rejects different lengths before scanning, then compares UTF-16 units. A cached hash can reject some unequal large strings sooner, but adds per-string memory and still needs a content check for hash collisions; only add it if benchmarks justify that cost.

The emitted plan must evaluate each source operand once and preserve left-to-right evaluation and any observable conversion side effects. For pure native numeric operands, the relational operators can share one three-way comparator. Dynamic coercion must retain the runtime's specified operand order.

## Complete operand domain

The checker exposes these statically typed primitive families to the Go emitter:

| Family | Source type spellings | Native representation |
| --- | --- | --- |
| Numeric | `number`, `float32`, `float64`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `number` and `float64` both use Go `float64`; other spellings use their corresponding Go type. `int` and `uint` have platform width `W`. |
| String | `string` | `*tsString` |
| Boolean | `boolean` | `bool` |

There are 15 typed base spellings. For each base, the nullable mask is one of `0` (nonnullable), `1` (nullable with null), `2` (nullable with undefined), or `3` (nullable with both). A symbolic typed operand is therefore `T(k,m)`, where `k` is one of the 15 bases and `m` is one of those four masks. This gives 60 typed shapes before representation equivalences are folded together. `number` and `float64` share a comparison representation, but their spelling should remain available to conversion policy where needed.

The emitter also has two statically recognizable nullish constants, `Cnull` and `Cundefined`, and one dynamic representation, `D = tsValue`. `D` can contain any of the 35 runtime kind tags:

| Dynamic family | Runtime kinds |
| --- | --- |
| Nullish | undefined, null |
| Numeric | number, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64 |
| Other primitives | boolean, string, bigint, symbol |
| References | object, array, function, class, instance, regexp, error, promise, task, namespace, import, module, iterator, collection, typed array, array buffer, data view |

Objects, class instances, functions, and other nonprimitive TypeScript types generally reach comparison lowering as `D`; the Go emitter does not assign each such TypeScript type a distinct native comparison representation. BigInt is likewise a dynamic kind, not one of the native typed numeric families above.

The complete ordered matrix is the product of these symbolic domains, not a hand-maintained list of rows:

```text
S x S, where
S = {T(k,m) | k in typed bases, m in masks}
    union {Cnull, Cundefined, D}
```

Whenever an operand is `D`, its runtime kind is selected from the 35-kind table above; `D × D` therefore expands to the complete runtime-kind cross product. This formulation accounts for every left/right combination while allowing identical lowerings to share one rule.

## Comparison plan representation

The systematic emitter should first build a symbolic plan, then render it. A plan needs only these fields:

```text
ComparePlan {
    operation: LooseEqual | StrictEqual | SameValue | SameValueZero | Relational
    left, right: OperandShape
    payload: Native(target) | ExactNumeric | String | Boolean | BigInt | Reference | Dynamic
    nullish: NoNulls | EqualityTags | RelationalNullishConversion | DynamicNullish
    relationalResult: Less | Equal | Greater | Unordered | Throw
}

OperandShape {
    family: Numeric | String | Boolean | BigInt | Symbol | Reference | Dynamic | Null | Undefined
    numericKind: optional source numeric spelling
    width: optional integer bit width (including platform W)
    nullableMask: 0..3
    runtimeTagSet: set of tsKind values, only for Dynamic
}
```

The planner should normalize aliases and then apply a small set of rules. It should not contain one bespoke helper name for each numeric pair and operator spelling.

## Numeric pair matrix

Every typed numeric pair `T(a,m) × T(b,n)` uses the same payload decision procedure. Null handling is layered on top and described separately below.

1. If both payloads have the same native type, emit the native Go operator directly.
2. Otherwise search for a common native target that represents **every value** of both source types exactly. Prefer the target with the fewest required conversions; emit conversions and the native operator.
3. If no such target exists, use the exact numeric comparison implementation. Never cast one operand to the other merely because Go permits that cast.

For integer source `S` and target `T`, the cast is safe only when the complete range of `S` is contained in `T`. For integer-to-float conversion, every source value must fit the target mantissa precision. Float32 has 24 bits of precision; float64 has 53. Float32 widens exactly to float64. A float type cannot be converted to an integer type for a general comparison because it can contain fractions, NaN, or infinity.

The safe automatic widening relation for a statically known payload is:

| Source | Lossless targets |
| --- | --- |
| Signed integer of width `b` | Signed integers of width at least `b`; float32 if `b <= 24`; float64/`number` if `b <= 53`. Never an unsigned target for the full source range. |
| Unsigned integer of width `b` | Unsigned integers of width at least `b`; signed integers strictly wider than `b`; float32 if `b <= 24`; float64/`number` if `b <= 53`. |
| float32 | float32 and float64/`number`. |
| float64/`number` | float64/`number` only. |

`int` and `uint` use the platform width `W` in those rules. “Width” means the source type's full range, not the current operand's value. For example, a particular int64 may happen to fit in float64 exactly, but the type pair cannot assume that. Same-type comparison needs no conversion. A value-specific check may choose a narrower representation only when it proves the actual value is exact; the current planner does not need that extra branch to be correct.

When both operands need conversion, the common-target selector considers every native numeric type, including types that are neither input. It ranks candidates by number of casts, then target width, then prefers an integer target over a floating target on a tie. That yields int64 for int32/uint32 and float64 for int32/float32. If no target represents both complete ranges exactly, the plan uses the exact numeric comparator. This is a comparison-local representation choice; whether the checker permits implicit assignments remains a separate policy.

Examples of the common-target rule:

| Pair | Best plan | Emitter result |
| --- | --- | --- |
| `float32` / `float64` | Convert float32 to float64; direct Go comparison. | Selected as a one-cast common target. |
| `int32` / `float32` | Convert both to float64; direct Go comparison. | The selector finds float64 as a third common target. |
| `int32` / `uint32` | Convert both to int64; direct Go comparison. | The selector finds int64 as a third common target. |
| `int32` / `int64` | Convert int32 to int64; direct Go comparison. | Selected as a one-cast common target. |
| `int64` / `float64` | Exact numeric comparison; do not cast either operand. | Uses the exact numeric runtime comparator. |
| `int64` / `uint64` | Exact signed/unsigned comparison; no common native integer type covers both ranges. | Uses the exact numeric runtime comparator. |

The grouped rules below cover the full integer/float submatrix, including platform-sized `int` and `uint`:

| Source pair group | Safe target rule |
| --- | --- |
| Signed integer / signed integer | A signed target whose width is at least the larger source width. |
| Unsigned integer / unsigned integer | An unsigned target whose width is at least the larger source width. |
| Signed integer / unsigned integer | A signed target wide enough for the signed source and wider than the unsigned source; otherwise use exact signed/unsigned comparison. |
| Float32 / float32 | Float32. |
| Float32 / float64 (including `number`) | Float64. |
| Float32 / integer | Float32 if every integer value fits its 24-bit precision; otherwise float64 if every integer value fits its 53-bit precision; otherwise exact comparison. |
| Float64 / integer | Float64 if every integer value fits its 53-bit precision; otherwise exact comparison. |

For integer-only pairs the target search includes native integer types beyond the two input types. For mixed integer/float pairs it includes both float types. This is what finds third targets without a hand-written pair table.

The rule is symmetric in the operand types: exchanging the types does not change the selected common target or exact fallback. Rendering still preserves operand order for any non-native conversions.

### Exact int64/float64 comparison

An IEEE float64 can represent finite integral values beyond the int64 range. The runtime must check range bounds before converting the float to int64. The current implementation does this; it does not need to construct a BigInt for this pair.

For a signed integer `i` and float `f`, the important boundaries are:

```text
minInt64 = -2^63             // representable and inclusive
maxInt64Exclusive = 2^63     // representable; int64 cannot contain it

if f < minInt64:             i is greater than f
if f >= maxInt64Exclusive:   i is less than f
otherwise:
    j = int64(f)             // safe after the bounds check
    compare i and j
    if equal, compare f with float64(j) to account for a fractional remainder
```

`float64(math.MaxInt64)` rounds to `2^63`, so the upper test must use the exclusive `2^63` boundary; converting `math.MaxInt64` to float and treating that as an inclusive maximum would be wrong. A CPU can compare the float against these representable constants directly. The unsigned path analogously checks `f < 0` and `f >= 2^64` before converting to uint64. NaN is unordered; infinities are handled by the range checks.

The exact comparison handles values such as `int64(9007199254740993)` against float64 `9007199254740992` without rounding the integer down first. The separate JavaScript BigInt path uses Go's `math/big` support; in particular, BigInt versus float comparison is exact there as well. The current BigInt/float path constructs a `big.Rat`, so it may allocate; an exact bit-level comparison could be considered if profiling shows this case matters. BigInt is not a substitute for the int64/float64 algorithm.

## Nullability layer

The current typed nullable representation is `tsOptional[T]{value T, tag uint8}`: tag 0 means a payload, 1 means null, and 2 means undefined. It is not boxed as `tsValue`. **Zero is not null**: `value == 0` with `tag == 0` is an ordinary numeric zero; the tag determines whether the payload is present. When the tag is nonzero, the payload is ignored; constructors may leave zero there, while a reused typed cell can retain its previous payload. A typed binding also tracks initialization separately for temporal-dead-zone behavior. Dynamic `tsValue` instead carries a runtime kind tag.

For any typed pair, first apply the operator's null rule, then use the payload plan from the preceding sections:

| Operation | Both payloads present | One or both nullish |
| --- | --- | --- |
| Loose equality | Apply the payload loose-equality rule. | Any nullish left equals any nullish right, including null versus undefined. A nullish value does not equal a present value. |
| Strict equality | Apply the payload strict-equality rule. | Equal only when both tags are the same nullish tag. |
| Relational | Apply the payload relational rule. | Null converts to numeric zero; undefined makes the relational result unordered (false for all four relational operators). |

For nullable values of different payload types, these tag rules still apply. For example, strict comparison of two nullable values with different payload families can be true when both are null, even though their present payloads would never be strictly equal. When both tags indicate present values, numeric pairs use the numeric plan; unrelated strict payload families return false.

This tag-first rule is why a nullable float should not compare its payload before checking its tag: a nullish value's payload is not its value. A bitmask does not remove the need for the tag. The hot present-value path can be ordered first in emitted branches, but that instruction ordering should be chosen from generated-code benchmarks.

The representation is chosen to preserve typed payloads without paying for the universal value on every operation. `tsValue` has a kind, numeric payload, and pointer field (24 bytes on the current 64-bit Go ABI); typed nullable values keep the native `T` plus a one-byte tag. Dense nullable arrays go further and store native payloads and tags in parallel slices, avoiding per-element struct padding. Boxing every nullable primitive into `tsValue` would simplify dispatch but add size, tag switches, and boundary conversions to typed reads and writes. It would also make direct operations on known payloads harder to emit. The tradeoff is that each typed nullable operation must correctly apply its tag rule before using the payload.

## Typed nonnumeric and mixed-family matrix

The statically native, nonnumeric value families are `string` and `boolean`. All object-like source types, symbols, and BigInts are dynamic `tsValue`s. For present typed payloads, the semantic pair rules are:

| Left / right | Loose equality | Strict equality | Relational |
| --- | --- | --- | --- |
| Numeric / numeric | Numeric plan; exact cross-family equality if no safe common target. | Same numeric-value rule across numeric representations. | Common native target or exact ordered numeric comparison. |
| String / string | Compare UTF-16 contents. | Compare UTF-16 contents. | Lexicographic UTF-16 code-unit order. |
| Boolean / boolean | Compare the boolean values. | Compare the boolean values. | Convert `false`/`true` to Number `0`/`1`, then order. |
| Number / string | Apply `ToNumber` to the string; compare numeric values. | False. | Apply `ToNumber` to the string, then numeric ordering. |
| Number / boolean | Convert boolean to Number; compare numeric values. | False. | Convert boolean to Number, then numeric ordering. |
| String / boolean | Convert boolean to Number and string with `ToNumber`; compare. | False. | Convert both to Number, then numeric ordering. |
| Typed family / dynamic | Apply the dynamic loose-equality rules below after preserving typed null tags. | Apply the dynamic strict-equality rules below after preserving typed null tags. | Apply the dynamic relational rules below, including throwing conversions. |
| Dynamic / dynamic | Factored runtime loose-equality dispatch. | Factored runtime strict-equality dispatch. | Factored runtime relational dispatch. |

The table describes present payloads. If either side is nullable, first materialize its semantic value from its tag (`null`, `undefined`, or the present payload); the nullish comparison rules still apply. In particular, a typed nullable string compared with a number does not compare its unused zero-value payload.

For typed `string`/`boolean` against each dynamic value, the same semantic rules apply after unboxing the known operand. A direct typed fast path is valid only when it preserves those conversions and errors. For example, typed string versus dynamic string can compare contents directly after checking the runtime kind; typed string versus object must still allow object conversion and its side effects.

The dynamic `tsValue` domain has 35 runtime tags. Its cross product is covered by the factored rules below, rather than 1,225 bespoke rows. The runtime tags first collapse into semantic classes: nullish (`undefined`, `null`), Number (12 native numeric tags), Boolean, String, BigInt, Symbol, and references (17 object-like tags). `==`/`!=`, `===`/`!==`, and the four relational operators share the three operation classes above, but equality and ordering use different conversion rules.

### Strict equality: complete kind-pair rules

After nullability tags are resolved, strict equality has a small complete partition:

| Left / right semantic class | Result |
| --- | --- |
| Number / Number | Exact numeric-value equality across the numeric tags; `NaN` is unequal and `+0` equals `-0`. |
| BigInt / BigInt | Integer-value equality. |
| String / String | UTF-16 code-unit content equality. |
| Boolean / Boolean | Boolean-value equality. |
| Symbol / Symbol | Symbol identity, not description text. |
| Null / Null; Undefined / Undefined | True. |
| Any two reference values | Reference identity, never structural equality. |
| Any two different classes, including Null / Undefined | False. |

The runtime uses identity for ordinary objects, arrays, functions, classes, instances, errors, regexps, promises, iterators, collections, typed arrays, ArrayBuffers, and DataViews. It also uses pointer identity for tasks and modules, module identity for namespace wrappers, and `(module, export name)` identity for import references. A namespace wrapper is allocated on conversion, so its Go wrapper address is not its identity key.

### Loose equality: complete coercion partition

Use the following ordered rules and restart after any conversion that produces a new primitive. “Number” below includes every numeric runtime tag and uses exact cross-tag numeric comparison.

| Pair | Rule |
| --- | --- |
| Same primitive class | Use that class's strict-value equality, except Number uses exact cross-tag numeric equality. |
| Nullish / nullish | `null == undefined` and same-kind nullish pairs are true. |
| Nullish / any non-nullish value | False; do not coerce either operand. |
| Boolean / any other value | Convert the boolean to Number `0` or `1`, then restart. This includes BigInt and references. |
| Number / String | Convert the string with `ToNumber`; compare Number values. Invalid numeric text becomes `NaN`, hence false. |
| BigInt / String | Parse with `StringToBigInt`; compare integer values. Invalid text is false, not an exception. |
| Number / BigInt | Exact mixed numeric equality. Do not round the BigInt through float64; nonintegral, nonfinite, and out-of-range Number values cannot equal an integer BigInt. |
| Reference / primitive (except nullish) | Convert the reference with `ToPrimitive` and the default hint, then restart. A Date-like object uses its default string hint; custom `Symbol.toPrimitive`, `valueOf`, or `toString` can have side effects or throw. |
| Reference / reference | Compare identity; do not call either object's conversion hooks. |
| Symbol / different primitive class | False; do not convert the Symbol. |

These rules cover symmetric operand order, but do not justify arbitrary reordering: conversions can execute user code. In loose equality, a Symbol can compare equal to the identical Symbol, and an object can compare equal to a Symbol if its primitive conversion returns that exact Symbol. A Symbol compared loosely with a Number or String remains false without a numeric/string conversion of the Symbol.

The transpiler's exact integer kinds extend JavaScript Number semantics: when a dynamic `int*`/`uint*` tag is compared loosely or relationally with a numeric string, the runtime compares its exact decimal or radix-prefixed value instead of rounding the integer through float64. Float-backed `number`/`float32` values still use JavaScript `ToNumber` behavior. With `CoerceAny` disabled, a dynamic string/integer or string/Number loose or relational pair throws a checked-conversion `TypeError`; strict equality remains a type-sensitive comparison. Keep this project policy explicit in addition to the ECMAScript rules above.

### Relational comparison: complete conversion and result partition

For `<`, `<=`, `>`, and `>=`, evaluate the operand expressions once in source order, then apply `ToPrimitive` with the number hint in left-to-right order. The comparison result is a partial order: it is less, equal, greater, unordered, or an abrupt completion (throw). All four operators return false for unordered results; they do not turn `NaN` into a true result by negating an ordinary boolean.

After primitive conversion:

1. If both primitives are strings, compare UTF-16 code units lexicographically. No numeric parsing happens in this branch.
2. If exactly one primitive is BigInt and the other is a String, parse the string with `StringToBigInt`. Invalid text produces an unordered result; otherwise compare the two BigInt values exactly.
3. Otherwise apply `ToNumeric`: Boolean becomes Number `0`/`1`; String becomes Number (invalid text becomes `NaN`); null becomes Number `0`; undefined becomes `NaN`; Number remains Number; BigInt remains BigInt; Symbol throws `TypeError`.
4. Compare Number/Number with IEEE ordering (`NaN` is unordered), BigInt/BigInt with integer ordering, and BigInt/Number with exact mixed ordering. Never round a BigInt through float64.
5. Map a valid ordered result to `<`, `<=`, `>`, or `>=`; map unordered to false for all four operators.

This gives a useful compact pair table:

| Numeric primitives after conversion | Ordered comparison |
| --- | --- |
| Number / Number | IEEE comparison; unordered if either side is `NaN`. |
| BigInt / BigInt | Exact integer comparison. |
| BigInt / String or String / BigInt | Parse the string with `StringToBigInt`; invalid text is unordered. |
| BigInt / Number or Number / BigInt | Exact mixed comparison; unordered if the Number is `NaN`. |
| Symbol / anything requiring `ToNumeric` | Throw `TypeError`. |

Relational BigInt/string uses the same `StringToBigInt` parse as loose equality, but produces an ordering result instead of testing equality. Thus `1n < "2"` is true, while `1n < "1.5"` and `1n < "not a number"` are false because parsing is invalid and the comparison is unordered.

Null and undefined are converted only in the non-string branch: null compares as Number zero; undefined produces `NaN` and therefore all four relational operators are false. A reference must first produce a primitive; conversion hooks can throw or mutate state. When either primitive is a Symbol, conversion to numeric throws rather than yielding an unordered result.

`tsBinary` currently performs left-then-right `ToPrimitive` for the relational family, which is the required observable order. It then handles two strings lexically, numeric values with `tsNumericCompare`, and BigInt through `tsBigIntCompare`. The BigInt relational branch parses strings with `StringToBigInt`, treats invalid text as unordered, and throws `TypeError` for Symbol. Keep “unordered” distinct from “conversion failed with an exception” in the plan/runtime contract.

## Short-circuit logical operators

Logical operators are not eager boolean comparisons. Unary `!` is a truthiness conversion; binary logical operators use a control-flow plan with one left-side test and a conditional right-side block:

| Operator | Plan |
| --- | --- |
| `!value` | Evaluate `value` once, test truthiness, return the opposite boolean. |

For statically known values, truthiness can usually be emitted directly from the native representation. For dynamic values it uses kind-specific truthiness. In either case, preserve observable conversion behavior and evaluate the operand once.

The binary logical plan is:

| Operator | Evaluate right when | Result when right is skipped |
| --- | --- | --- |
| `left && right` | `left` is truthy | The original `left` value. |
| `left || right` | `left` is falsy | The original `left` value. |
| `left ?? right` | `left` is null or undefined | The original `left` value. |

When the right side is evaluated, its original value is the result. Do not lower `&&` or `||` to eager Go boolean operators unless both operands are statically boolean and short-circuit evaluation is preserved. Do not use truthiness for `??`.

Truthiness can be selected by operand shape: booleans use the value; integers test nonzero; floats test nonzero and not NaN; strings test nonempty; nullable values check the tag and then the payload; dynamic values use kind-specific truthiness; objects and symbols are truthy; BigInt is false only at zero. For `??`, a typed nullable value only needs its tag, while a known nonnullable value never selects the right side. The emitter now uses direct tests for known native booleans, numbers, strings, and optional payloads; only dynamic values use `tsTruthy`/`tsNullish`.

Across all 35 dynamic kinds, only these values are falsy: undefined, null, false, numeric `+0`/`-0`, numeric `NaN`, the empty string, and `0n`. Every Symbol and every reference is truthy, including empty arrays and objects. `??` is deliberately narrower: it selects the right side only for null and undefined. This is a complete truthiness partition; do not route it through equality or numeric comparison.

The compound forms `&&=`, `||=`, and `??=` share those tests but must evaluate an lvalue reference once and only write when the right side is selected. Their implementation belongs to the same logical plan with an assignment action, not a new type-pair matrix.

## Bitwise operators

Bitwise operators have three distinct operand families. They must not be merged into the relational numeric widening rules.

| Operand family | `&`, `|`, `^` | `<<`, `>>`, `>>>` | `~` |
| --- | --- | --- | --- |
| JavaScript Number (`number`, `float32`, `float64`, dynamic Number/Float32) | Convert both operands to uint32 bits; return signed int32 result as a Number. | Convert to 32 bits; mask shift count with 31. `>>` is signed, `>>>` is unsigned and returns 0..2^32−1. | Convert to int32, complement, return signed int32 as a Number. |
| Typed integer extensions (`int`, fixed-width signed/unsigned integers) | Operate on the selected integer width and signedness. | Apply that integer family's shift rules; signed `>>` is arithmetic and `>>>` shifts the same width as unsigned. | Complement at the selected width. |
| BigInt | Arbitrary-precision two's-complement operations. | `<<`/`>>` operate on BigInts; negative shift counts reverse direction. `>>>` throws. | Arbitrary-precision complement. |

For Number operands, `ToUint32` maps zero, NaN, infinities, null, and undefined to zero; otherwise it truncates toward zero and reduces modulo 2^32. The two input conversions and the 32-bit operation are required semantics; they are separate from the exact-preserving numeric comparison plan. For BigInt, both operands must be BigInts; mixing BigInt and Number throws. The runtime uses `tsNumberBitwise`/`tsUint32Bitwise` for Number and `tsBigIntBinary` for BigInt. Static same-type integer `&`, `|`, and `^` lower directly to Go operators; shifts use a shared width-aware helper so negative shifts keep the checked error. Statically mixed integer types follow the runtime's existing Number32 rule, and native float/boolean/string pairs use the Number32 path.

For dynamic Number bitwise operations, each operand first undergoes `ToPrimitive` and `ToNumeric`. Strings, booleans, and nullish values can therefore become Number operands; Symbols throw. A BigInt paired with any value that converts to Number throws for mixing BigInt and Number. The `&&`/`||`/`??` result-preserving rules do not apply to `&`/`|`/`^`.

Nullable bitwise operands use the operator's conversion rules after inspecting the optional tag. For Number semantics, null and undefined both ultimately produce zero bits. For BigInt semantics, a nullish operand mixed with a BigInt becomes a Number and therefore throws as a mixed-family operation. A typed fast path may handle the present/present case directly, but it must retain the specified nullish and error cases.

The current plan preserves the runtime's boundary rule for integer extensions: equal native integer kinds use that width; differing numeric kinds use Number32. `>>>` is emitted through a helper because Go has no `>>>` token, and float32 never receives a Go bitwise operator. Nullable and dynamic operands retain the runtime conversion path so nullish values, BigInt mixing, Symbols, and user-defined `ToPrimitive` behavior remain intact.

The efficient implementation does not need a type-pair helper for every bitwise operation. Normalize to those three families, select a common typed integer width only under a documented safe rule, and share small emit patterns for native bitwise, Number32, and BigInt. Keep `>>>` as a distinct operation because it has different signedness for Number and is invalid for BigInt.

## Membership and instance tests

`in` and `instanceof` look comparison-like but are not equality or relational pairs. They need independent plans, not rows added to the primitive matrix.

| Operator | Left-side work | Right-side requirement and work | Observable behavior to preserve |
| --- | --- | --- | --- |
| `key in object` | Evaluate and retain the key value; do not convert it yet. | First require an object; then convert the retained key with `ToPropertyKey` (Symbol stays Symbol, other keys become strings) and run `HasProperty`. | A primitive right operand throws before key conversion. Proxy `has` traps and prototype-chain lookup can execute code or throw. |
| `value instanceof constructor` | Preserve the original value. | First consult `@@hasInstance` when present; otherwise require a callable constructor and use its ordinary instance rule. | Custom `@@hasInstance`, proxy/property access, bound functions, mutable prototypes, and invalid constructors are observable. A primitive left value normally yields false after the constructor checks. |

The current runtime routes `in` through `tsHas` and `instanceof` through `tsInstanceOf`/`tsOrdinaryInstanceOf`. `tsHas` now rejects a primitive right operand before converting the key, preserving the required TypeError and preventing an invalid left-key conversion hook from running. The implementation remains general for proxies, prototype lookup, custom hooks, and error behavior; a faster static property path would need to prove a closed prototype/property shape.

## Lowering principles

These rules reduce the complete operator space to a few plan builders:

1. **Separate semantic families first.** Equality, relational, short-circuit logical, Number bitwise, typed integer bitwise, and BigInt bitwise have different conversions and errors.
2. **Factor operators by result relation.** `!=` and `!==` negate their corresponding equality result. `<`, `<=`, `>`, and `>=` share one ordered/unordered comparison result. They do not need separate numeric algorithms.
3. **Use static facts before runtime dispatch.** Same native types and safe common native targets emit direct Go operations. A target is safe only if it preserves the full source ranges.
4. **Keep nullability orthogonal.** Apply one equality/relational/logical null rule, then reuse the selected present-payload plan. Never infer null from a zero payload.
5. **Keep exact fallbacks exact.** If no safe native target exists, use a shared exact comparator. Do not force a lossy cast or arbitrary-precision allocation when a bounded exact algorithm exists.
6. **Dispatch once for dynamic values.** Known/dynamic pairs can specialize the known side and switch on the dynamic kind. Dynamic/dynamic values use factored runtime dispatch. Avoid duplicating the entire kind cross product in emitted source.
7. **Preserve effects and error behavior.** Evaluate operands once, preserve short-circuiting and conversion order, and distinguish false/unordered from a thrown conversion.
8. **Measure code shape, not helper count.** Direct native operators are the fast path, but larger inline branches can cost code size and instruction-cache performance. Benchmark representative generated programs before expanding specialization.

The emitter's cost ladder should choose the first valid plan for the operand shapes:

1. A statically certain result, such as strict equality across two different known primitive families.
2. Same present-value family: direct numeric, boolean, string-content, or identity operation.
3. Same family with nullability: inline tag checks, then the same present-value operation.
4. Different but known primitive families: apply only the required conversion, then reuse numeric/string comparison logic. Keep exact mixed numeric fallbacks when a native cast would lose information.
5. One dynamic side: check its runtime tag once and use the known-side operation; send only exceptional kinds to the conversion helper.
6. Two dynamic sides: use the complete abstract comparison algorithm, retaining nullish, BigInt, Symbol, reference identity, conversion side effects, and throws.

The steps are semantic tiers, not an assumption that one giant inlined branch is always faster. In particular, a dynamic-side kind switch can be faster than repeated helpers, but can also increase generated code size; measure it.

The symbolic plan can include an operation family and then reuse a small number of selectors:

```text
ComparePlan: equality kind | relational kind | payload plan | null plan
LogicalPlan: truthiness-or-nullish test | conditional evaluation | result representation | optional writeback
BitwisePlan: JSNumber32 | NativeInteger(width, signedness) | BigInt | conversion/error rule
```

Each builder consumes operand descriptors and produces a plan. The source emitter then renders shared patterns; coverage comes from the descriptors' Cartesian products and rule tests, not a handwritten matrix of generated function names.

## Dynamic lowering audit

The preceding strict, loose, and relational partitions define every dynamic tag pair after tags are grouped by semantic class. Implementation-specific observations:

- Dynamic comparison dispatch is in `tsStrictEqual`, `tsLooseEqual`, `tsBinary`, `tsNumericCompare`, and `tsBigIntCompare`.
- `tsPolicyCompare` keeps compiler conversion policy separate from the comparison algorithm. With `CoerceAny` off, dynamic string/numeric loose and relational pairs reject implicit conversion; strict comparisons still use strict equality. With it on, typed/dynamic pairs perform checked conversion, while dynamic/dynamic pairs retain factored runtime semantics and exact numeric payloads.
- The BigInt relational route parses a paired string with `StringToBigInt` and throws for Symbol. Loose BigInt/Symbol equality remains false; do not share the relational throw with the equality comparator.
- The strict-equality switch has no `task`, `namespace`, `import`, or `module` identity cases. Confirm which are observable comparison values, and use stable underlying identity where they are. In particular, a namespace wrapper is allocated on conversion; wrapper address cannot be its identity key.
- Every reference/primitive loose equality or relational comparison may run `Symbol.toPrimitive`, `valueOf`, or `toString`; both-reference loose equality does not run those hooks, while relational comparison does. This is why “reference/reference” has different equality and relational plans.
- Error behavior belongs in the plan. Invalid numeric text becomes unordered `NaN` for Number comparisons, but invalid BigInt text is simply false for loose BigInt/string equality. Symbol numeric conversion is a thrown `TypeError`, not unordered.

The full dynamic cross product is therefore covered by a seven-class domain and three algorithmic partitions, not a per-tag hand switch: (1) strict same-class value/identity rules, (2) ordered loose-equality rewrites with recursion after conversion, and (3) relational ToPrimitive/ToNumeric followed by a three-way Number/BigInt comparator. Runtime tags that are implementation-only should either be excluded from source-visible `tsValue` pairs or assigned an explicit identity rule.

## Current lowering coverage and remaining boundaries

The emitter and runtime cover the main fast paths with a small set of shared plans:

- Every typed numeric pair uses the safe common-target selector or the exact mixed numeric comparison. Planner tests enumerate the full supported numeric-kind product and reject lossy targets.
- Nullable numeric equality and ordering test tags before reusing the payload plan. Null becomes numeric zero for relational comparisons; undefined makes all four results false.
- Same-family nullable strings and booleans use direct content/value comparisons plus tag handling. Nullable string ordering compares present strings by UTF-16 and applies null-to-number conversion for a null/present pair.
- Boolean ordering and Boolean/numeric loose or relational comparisons use direct boolean/zero/one expressions. Known strict comparisons across incompatible primitive families are constants, except matching nullish tags on two optional operands.
- String equality checks pointer identity first, then length, then UTF-16 content. String ordering compares UTF-16 code units.
- Every statically known native primitive family, including optional forms, has shared typed-versus-dynamic dispatch. It preserves the known native operand until runtime dispatch, applies `CoerceAny` checked conversion, and keeps exact numeric tags.
- Dynamic/dynamic comparisons use factored runtime rules. Exact integer tags retain precision against other numeric tags; integer-tag/numeric-string comparisons parse recognized decimal and radix-prefixed values exactly instead of rounding through float64.
- `&&`, `||`, `??`, and their assignment forms preserve short-circuiting and return the selected original value. Known native truthiness and nullish tests avoid boxing; dynamic values retain runtime dispatch.
- Same-type native integer `&`, `|`, and `^` emit directly; shifts use one width-aware helper. Mixed numeric types and native float/boolean/string operands use Number32 conversion. Nullable, dynamic, and BigInt inputs keep runtime conversion and error behavior.
- Strict identity includes task, module, namespace, and import values. `in` rejects primitive right operands before converting the key. `instanceof` stays on the general runtime path.

The remaining boundaries are about avoiding unsafe specialization rather than missing type-pair rows:

1. Dynamic values require runtime-kind dispatch. A type switch is necessary when source types do not provide stronger facts; adding a generated branch for every runtime kind is only justified by a benchmark.
2. Nullable/dynamic bitwise pairs retain runtime dispatch because nullish conversion, BigInt mixing, Symbols, and object `ToPrimitive` hooks affect results. A present/present shortcut should include tests for every fallback and a measured performance benefit.
3. `in` and `instanceof` retain proxy, prototype, and custom-hook behavior. A static fast path needs a closed-shape proof that the current emitter does not have.
4. Integer-tag/numeric-string comparisons preserve exact typed integer values; float-backed `number`/string comparisons still use ECMAScript Number conversion. Keep these paths explicit so an optimization cannot silently round an integer.
5. The generic dynamic helper has not been benchmarked against a larger inline dispatch tree. Current emitted-shape checks establish which path is selected; performance claims beyond direct native operators need measurements.

## Planner structure

The numeric comparison planner derives plans from range and representation metadata rather than a handwritten pair switch:

1. Normalize the source type to a comparison representation: signed integer width, unsigned integer width, float precision, string, boolean, dynamic, null, or undefined. Keep the original null mask and source spelling.
2. For two native numeric payloads, enumerate candidate native targets and retain only targets that losslessly represent both complete source domains, using the widening relation above.
3. Rank safe targets by emitted cost: same native type first; then one cast; then two casts. Break ties by target width and integer-before-float. If there is no target, select the exact numeric fallback.
4. Wrap the payload plan in a nullability plan for equality or relational semantics.
5. Render one of a small set of code patterns: native operator, tag-guarded native operator, exact numeric helper, string/boolean helper, typed-versus-dynamic dispatch, or dynamic-versus-dynamic runtime operation.

The implementation lives in `internal/core/numerictypes.go` and `internal/goemit/comparison_lowering.go`. Typed-versus-dynamic dispatch, truthiness, and bitwise lowering reuse the same operand descriptors and shared render patterns.

For nonnumeric plans, keep the renderer small and type-directed:

1. **String/string:** call content equality or UTF-16 ordering; never compare `*tsString` pointers.
2. **Boolean/boolean:** emit direct equality and the equivalent boolean formulas for ordering. For Boolean/numeric, convert the Boolean to the exact numeric value 0 or 1 and reuse the numeric planner.
3. **Strict primitive equality:** use direct Boolean equality, string content equality, BigInt value equality, Symbol identity, and tag-correct nullish equality. Strict inequality is the negation of that result.
4. **Loose equality:** first apply the ordered conversion rules above. Emit a direct plan only when both operand kinds and conversion outcomes are statically known; otherwise dispatch on the dynamic kind and preserve object hooks and conversion errors.
5. **Relational:** direct string ordering only for two known strings; otherwise normalize through the specified primitive/numeric rules. Carry unordered and throw as different outcomes.
6. **References:** use identity for strict/loose equality, but use `ToPrimitive` for mixed reference/primitive loose equality and every relational comparison. Never structurally compare two objects.
7. **Nullable values:** gate payload operations with null/undefined tags and apply the relevant operator's nullish behavior before entering the payload fast path.
8. **Typed-boundary policy:** apply `CoerceAny`/checked conversion policy outside the JavaScript comparison plan, so typed annotation behavior cannot corrupt dynamic `==`/ordering semantics.

Candidate selection must account for platform width `W`. It must not choose float32 as a common target for a full int32 domain, or float64 for a full int64 domain. It may choose float32 for integer types whose entire range fits its 24-bit precision, and float64 for integer types whose entire range fits its 53-bit precision.

For nullable numeric pairs, the emitted shape should be conceptually:

```go
// Equality: tags first, then the preselected payload plan.
if left.tag != 0 || right.tag != 0 {
    return equalityForNullishTags(left.tag, right.tag, loose)
}
return comparePresentPayloads(left.value, right.value, payloadPlan)
```

The generated form can inline this logic. The pseudocode names illustrate shared rules; they do not prescribe runtime helper calls for the fast path.

## Validation grid for later implementation

The plan generator should be checked against an independent semantic oracle for every symbolic plan class. Generate directed operand pairs from the seven semantic dynamic classes and exercise every actual tag within each class for identity/tag-specific cases. Cover each typed base/null mask and operator class; one hand-picked example is not evidence for an entire symbolic group. The useful boundary set includes:

- Integer minima/maxima for every width and signedness, including platform-dependent `int`/`uint`.
- Float32 rounding boundaries and float64 integers around ±2^53.
- Int64/float64 values around `-2^63`, `2^63`, `math.MinInt64`, `math.MaxInt64`, fractions, NaN, and infinities.
- Uint64/float64 values around `0` and `2^64`.
- Null, undefined, and a present value for every nullable mask; both equality modes and all relational operators.
- Numeric strings, invalid numeric strings, booleans, BigInts, Symbols, and objects with observable primitive-conversion hooks.
- Loose-equality cases `null == undefined`, Boolean-to-Number, Number/String, BigInt/String parsing, exact BigInt/Number, Symbol identity, object-to-primitive, and reference identity.
- Strict equality across every primitive class and every reference tag, including the task/namespace/import/module identity audit.
- Relational cases for both-string lexical order; Number, BigInt, and mixed BigInt/Number; BigInt with integer/fractional/invalid strings; null and undefined; NaN; Symbol throws; and object conversion side-effect order for all four operators.
- Both operand orders for every mixed-family rule, while checking that conversion hooks run in the spec-required order.
- Logical left operands spanning every truthiness family and nullish mask, with a right-side counter proving skipped expressions are not evaluated.
- Membership and instance tests with string/symbol keys, inherited properties, proxies, custom `@@hasInstance`, bound functions, and invalid right operands.
- For `in`, use a left key whose `Symbol.toPrimitive` records or throws alongside a primitive right operand; verify the right-operand TypeError occurs first and the key hook does not run.
- Number bitwise inputs at 0, ±1, 2^31, 2^32, NaN, infinities, fractions, numeric strings, booleans, nullish values, Symbols, and shift counts −1, 31, 32, and 33; typed integer width boundaries; BigInt mixed operands, negative shifts, and unsigned-right-shift errors.

Generated source should also be inspected for the expected shape: direct Go operators for native plans, no narrowing cast on exact numeric fallbacks, and no `tsValue` boxing on the nullable/common-target fast paths. Performance claims should be benchmarked on representative generated programs after correctness is established.

## Source map

- Primitive and null-mask normalization: `internal/goemit/types.go`, `internal/checker/emitresolver.go`.
- Safe numeric conversion and common-target selection: `internal/core/numerictypes.go`.
- Comparison plan emission: `internal/goemit/comparison_lowering.go`, called from `internal/goemit/expressions.go`.
- Exact numeric equality and ordering: `internal/goemit/type_runtime.go`, `internal/goemit/runtime.go`.
- Dynamic comparison boundaries: `internal/goemit/conversion_runtime.go`.
- BigInt comparison: `internal/goemit/bigint_runtime.go`.
- Dynamic kind inventory and value representation: `internal/goemit/value_runtime.go`.
