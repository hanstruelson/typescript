# JavaScript conformance checklist

Test262 revision: `2e0a56762801e275a9fdf96dc49d90ba0cddcf63`.
Inventory: 53,616 tests; 2,321 explicitly excluded for dynamic code.
Execution: 1,782 variants from 1,630 selected tests; selection: all.

**Full-suite run: results below are checkpoints until every selected variant completes.** Untested cases are not passes. The default coercion policy is unchanged; failures may expose differences from ECMAScript.

Module and host-adapter limitations are tracked separately from compiler/runtime failures. Node reference failures are recorded separately and do not prevent native execution. Negative tests pass only on syntactic rejection, never on an arbitrary build failure.

**Checkpoint interpretation:** the raw JSONL predates regex-parser error classification. Of its 240 `emit-error` rows, 196 are parse-negative RegExp cases whose Node result is the expected `SyntaxError`; the current runner classifies those as passes. The effective counts are 1,218 passes and 44 genuine emitter failures. See [ANALYSIS.md](ANALYSIS.md). This historical checkpoint cannot be resumed with the updated source digest.

## Results

- build-error: 38
- build-timeout: 8
- emit-error: 240
- pass: 1,022
- runner-host-unimplemented: 26
- runtime-fail: 448

## Category checklist

| Category | Suite tests | Dynamic exclusions | Executed variants | Pass | Native failures | Adapter/reference limits |
|---|---:|---:|---:|---:|---:|---:|
| annexB/built-ins | 241 | 34 | 68 | 14 | 42 | 12 |
| annexB/language | 845 | 309 | 0 | 0 | 0 | 0 |
| built-ins/AbstractModuleSource | 8 | 0 | 0 | 0 | 0 | 0 |
| built-ins/AggregateError | 25 | 2 | 0 | 0 | 0 | 0 |
| built-ins/Array | 3083 | 26 | 0 | 0 | 0 | 0 |
| built-ins/ArrayBuffer | 221 | 1 | 0 | 0 | 0 | 0 |
| built-ins/ArrayIteratorPrototype | 27 | 0 | 0 | 0 | 0 | 0 |
| built-ins/AsyncDisposableStack | 104 | 1 | 0 | 0 | 0 | 0 |
| built-ins/AsyncFromSyncIteratorPrototype | 38 | 0 | 0 | 0 | 0 | 0 |
| built-ins/AsyncFunction | 18 | 4 | 0 | 0 | 0 | 0 |
| built-ins/AsyncGeneratorFunction | 23 | 14 | 0 | 0 | 0 | 0 |
| built-ins/AsyncGeneratorPrototype | 48 | 0 | 0 | 0 | 0 | 0 |
| built-ins/AsyncIteratorPrototype | 13 | 0 | 0 | 0 | 0 | 0 |
| built-ins/Atomics | 389 | 0 | 0 | 0 | 0 | 0 |
| built-ins/BigInt | 77 | 0 | 0 | 0 | 0 | 0 |
| built-ins/Boolean | 51 | 3 | 0 | 0 | 0 | 0 |
| built-ins/DataView | 561 | 2 | 0 | 0 | 0 | 0 |
| built-ins/Date | 594 | 3 | 0 | 0 | 0 | 0 |
| built-ins/DisposableStack | 93 | 1 | 0 | 0 | 0 | 0 |
| built-ins/Error | 93 | 2 | 0 | 0 | 0 | 0 |
| built-ins/FinalizationRegistry | 47 | 1 | 0 | 0 | 0 | 0 |
| built-ins/Function | 509 | 151 | 0 | 0 | 0 | 0 |
| built-ins/GeneratorFunction | 23 | 14 | 0 | 0 | 0 | 0 |
| built-ins/GeneratorPrototype | 61 | 0 | 0 | 0 | 0 | 0 |
| built-ins/Infinity | 6 | 0 | 0 | 0 | 0 | 0 |
| built-ins/Iterator | 654 | 1 | 0 | 0 | 0 | 0 |
| built-ins/JSON | 166 | 0 | 0 | 0 | 0 | 0 |
| built-ins/Map | 204 | 2 | 0 | 0 | 0 | 0 |
| built-ins/MapIteratorPrototype | 11 | 0 | 0 | 0 | 0 | 0 |
| built-ins/Math | 327 | 0 | 0 | 0 | 0 | 0 |
| built-ins/NaN | 6 | 0 | 0 | 0 | 0 | 0 |
| built-ins/NativeErrors | 94 | 6 | 0 | 0 | 0 | 0 |
| built-ins/Number | 340 | 1 | 0 | 0 | 0 | 0 |
| built-ins/Object | 3411 | 57 | 0 | 0 | 0 | 0 |
| built-ins/Promise | 732 | 1 | 0 | 0 | 0 | 0 |
| built-ins/Proxy | 311 | 6 | 0 | 0 | 0 | 0 |
| built-ins/Reflect | 153 | 1 | 0 | 0 | 0 | 0 |
| built-ins/RegExp | 1879 | 287 | 1706 | 1008 | 684 | 14 |
| built-ins/RegExpStringIteratorPrototype | 17 | 0 | 0 | 0 | 0 | 0 |
| built-ins/Set | 383 | 1 | 0 | 0 | 0 | 0 |
| built-ins/SetIteratorPrototype | 11 | 0 | 0 | 0 | 0 | 0 |
| built-ins/ShadowRealm | 64 | 1 | 0 | 0 | 0 | 0 |
| built-ins/SharedArrayBuffer | 104 | 1 | 0 | 0 | 0 | 0 |
| built-ins/String | 1223 | 33 | 0 | 0 | 0 | 0 |
| built-ins/StringIteratorPrototype | 7 | 0 | 0 | 0 | 0 | 0 |
| built-ins/SuppressedError | 22 | 2 | 0 | 0 | 0 | 0 |
| built-ins/Symbol | 98 | 0 | 0 | 0 | 0 | 0 |
| built-ins/Temporal | 4605 | 0 | 0 | 0 | 0 | 0 |
| built-ins/ThrowTypeError | 14 | 1 | 0 | 0 | 0 | 0 |
| built-ins/TypedArray | 1453 | 0 | 0 | 0 | 0 | 0 |
| built-ins/TypedArrayConstructors | 738 | 12 | 0 | 0 | 0 | 0 |
| built-ins/Uint8Array | 70 | 0 | 0 | 0 | 0 | 0 |
| built-ins/WeakMap | 141 | 2 | 0 | 0 | 0 | 0 |
| built-ins/WeakRef | 29 | 1 | 0 | 0 | 0 | 0 |
| built-ins/WeakSet | 85 | 1 | 0 | 0 | 0 | 0 |
| built-ins/decodeURI | 55 | 0 | 0 | 0 | 0 | 0 |
| built-ins/decodeURIComponent | 56 | 0 | 0 | 0 | 0 | 0 |
| built-ins/encodeURI | 31 | 0 | 0 | 0 | 0 | 0 |
| built-ins/encodeURIComponent | 31 | 0 | 0 | 0 | 0 | 0 |
| built-ins/eval | 10 | 10 | 0 | 0 | 0 | 0 |
| built-ins/global | 29 | 8 | 0 | 0 | 0 | 0 |
| built-ins/isFinite | 15 | 0 | 0 | 0 | 0 | 0 |
| built-ins/isNaN | 15 | 0 | 0 | 0 | 0 | 0 |
| built-ins/parseFloat | 54 | 0 | 0 | 0 | 0 | 0 |
| built-ins/parseInt | 55 | 0 | 0 | 0 | 0 | 0 |
| built-ins/undefined | 8 | 1 | 0 | 0 | 0 | 0 |
| harness/assert-false.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-notsamevalue-nan.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-notsamevalue-notsame.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-notsamevalue-objects.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-notsamevalue-tostring.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-notsamevalue-zeros.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-obj.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-samevalue-nan.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-samevalue-objects.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-samevalue-same.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-samevalue-tostring.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-samevalue-zeros.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-custom-typeerror.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-custom.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-incorrect-ctor.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-native.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-no-arg.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-no-error.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-null-fn.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-primitive.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-same-realm.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-throws-single-arg.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-tostring.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assert-true.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/assertRelativeDateMs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-asyncTest-func-throws-sync.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-asyncTest-rejects-non-callable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-asyncTest-return-not-thenable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-asyncTest-returns-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-asyncTest-then-rejects.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-asyncTest-then-resolves.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-asyncTest-without-async-flag.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-custom-typeerror.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-custom.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-func-never-settles.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-func-throws-sync.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-incorrect-ctor.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-invalid-func.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-native.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-no-arg.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-no-error.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-primitive.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-resolved-error.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-same-realm.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/asyncHelpers-throwsAsync-single-arg.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/byteConversionValues.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-arguments.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-arraylike.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-different-elements.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-different-length.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-empty.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-falsy-arguments.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-message.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-same-elements-different-order.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-same-elements-same-order.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-samevalue.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-sparse.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/compare-array-symbol.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/dateConstants.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/decimalToHexString.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/deepEqual-array.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/deepEqual-circular.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/deepEqual-deep.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/deepEqual-mapset.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/deepEqual-object.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/deepEqual-primitives-bigint.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/deepEqual-primitives.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/detachArrayBuffer-host-detachArrayBuffer.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/detachArrayBuffer.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/fnGlobalObject.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/isConstructor.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/nans.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/nativeFunctionMatcher.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/promiseHelper.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifyconfigurable-configurable-object.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifyconfigurable-configurable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifyconfigurable-not-configurable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifyenumerable-enumerable-symbol.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifyenumerable-enumerable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifyenumerable-not-enumerable-symbol.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifyenumerable-not-enumerable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifynotconfigurable-configurable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifynotconfigurable-not-configurable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifynotenumerable-enumerable-symbol.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifynotenumerable-enumerable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifynotenumerable-not-enumerable-symbol.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifynotenumerable-not-enumerable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifynotwritable-not-writable-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifynotwritable-writable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifywritable-array-length.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifywritable-not-writable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/propertyhelper-verifywritable-writable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/proxytrapshelper-default.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/proxytrapshelper-overrides.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/sta.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/tcoHelper.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/testTypedArray-conversions-call-error.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/testTypedArray-conversions.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/testTypedArray.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-arguments.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-configurable-object.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-desc-is-not-object.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-noproperty.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-restore-accessor-symbol.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-restore-accessor.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-restore-symbol.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-restore.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-same-value.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-string-prop.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-symbol-prop.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-undefined-desc.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-value-error.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/verifyProperty-value.js | 1 | 0 | 0 | 0 | 0 | 0 |
| harness/wellKnownIntrinsicObjects.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/Array | 2 | 0 | 0 | 0 | 0 | 0 |
| intl402/BigInt | 11 | 0 | 0 | 0 | 0 | 0 |
| intl402/Collator | 65 | 1 | 0 | 0 | 0 | 0 |
| intl402/Date | 12 | 0 | 0 | 0 | 0 | 0 |
| intl402/DateTimeFormat | 245 | 1 | 0 | 0 | 0 | 0 |
| intl402/DisplayNames | 57 | 3 | 0 | 0 | 0 | 0 |
| intl402/DurationFormat | 110 | 0 | 0 | 0 | 0 | 0 |
| intl402/FallbackSymbol | 2 | 0 | 0 | 0 | 0 | 0 |
| intl402/Intl | 66 | 0 | 0 | 0 | 0 | 0 |
| intl402/ListFormat | 81 | 1 | 0 | 0 | 0 | 0 |
| intl402/Locale | 190 | 1 | 0 | 0 | 0 | 0 |
| intl402/Number | 7 | 0 | 0 | 0 | 0 | 0 |
| intl402/NumberFormat | 251 | 2 | 0 | 0 | 0 | 0 |
| intl402/PluralRules | 53 | 1 | 0 | 0 | 0 | 0 |
| intl402/RelativeTimeFormat | 80 | 1 | 0 | 0 | 0 | 0 |
| intl402/Segmenter | 79 | 4 | 0 | 0 | 0 | 0 |
| intl402/String | 19 | 0 | 0 | 0 | 0 | 0 |
| intl402/Temporal | 2029 | 0 | 0 | 0 | 0 | 0 |
| intl402/TypedArray | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/constructors-string-and-single-element-array.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/constructors-taint-Object-prototype-2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/constructors-taint-Object-prototype.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/default-locale-is-canonicalized.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/default-locale-is-supported.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/fallback-locales-are-supported.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/language-tags-canonicalized.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/language-tags-invalid.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/language-tags-valid.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/language-tags-with-underscore.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-consistent-with-resolvedOptions.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-default-locale-and-zxx-locale.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-duplicate-elements-removed.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-empty-and-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-locales-arg-coered-to-object.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-locales-arg-empty-array.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-returned-array-elements-are-not-frozen.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-taint-Array-2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-taint-Array.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-test-option-localeMatcher.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-throws-if-element-not-string-or-object.js | 1 | 0 | 0 | 0 | 0 | 0 |
| intl402/supportedLocalesOf-unicode-extensions-ignored.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.5-1-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.5-1gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.5-7-b-1-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.5-7-b-2-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.5-7-b-3-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.5-7-b-4-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-10-c-ii-1-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-10-c-ii-1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-10-c-ii-2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-11-b-1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-12-1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-12-2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-13-a-1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-13-a-2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-13-a-3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-13-c-1-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-13-c-2-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-13-c-3-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-14-c-1-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-14-c-4-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-2gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-5-1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-6-1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-6-2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-6-3-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-6-3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-6-4-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-6-4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.6-7-1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.1.6_A1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A3_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A3_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A3_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A3_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A5_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A5_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A5_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A5_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/S10.6_A7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/arguments-caller.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/async-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/async-gen-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/async-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/async-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/async-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/async-gen-named-func-expr-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/async-gen-named-func-expr-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/async-gen-named-func-expr-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/async-gen-named-func-expr-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/async-gen-named-func-expr-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-gen-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-gen-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/func-decl-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/func-decl-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/func-decl-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/func-decl-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/func-decl-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/func-expr-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/func-expr-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/func-expr-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/func-expr-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/func-expr-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-func-decl-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-func-decl-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-func-decl-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-func-decl-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-func-decl-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-func-expr-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-func-expr-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-func-expr-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-func-expr-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-func-expr-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/mapped | 43 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/meth-args-trailing-comma-multiple.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/meth-args-trailing-comma-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/meth-args-trailing-comma-single-args.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/meth-args-trailing-comma-spread-operator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/meth-args-trailing-comma-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/non-strict-arguments-object-is-immutable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/arguments-object/unmapped | 5 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9.2_A1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9.2_A1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9.2_A1_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9.2_A1_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9.2_A1_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9.2_A1_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9.2_A1_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T10.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T11.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T12.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T8.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A10_T9.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T10.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T11.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T8.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A11_T9.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.2_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.3_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.4_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.5_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.5_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.5_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.5_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.5_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.6_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.6_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.7_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.8_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A5.9_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T10.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T11.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T12.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T13.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T8.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.1_T9.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.2_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.2_T10.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.2_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.2_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.2_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.2_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.2_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.2_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.2_T8.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.2_T9.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.3_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.3_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.3_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.3_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.3_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.3_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.3_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.4_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A6.4_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A7_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A7_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A7_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A7_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A7_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A7_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A7_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A7_T8.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A7_T9.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A8_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A8_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A8_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A8_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A8_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A9_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A9_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A9_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A9_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A9_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A9_T8.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/S7.9_A9_T9.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/asi/do-while-same-line.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/block-scope/leave | 15 | 1 | 0 | 0 | 0 | 0 |
| language/block-scope/return-from | 2 | 0 | 0 | 0 | 0 | 0 |
| language/block-scope/shadowing | 15 | 0 | 0 | 0 | 0 | 0 |
| language/block-scope/syntax | 113 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A2_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A2_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A4_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A4_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A4_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A4_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A4_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A4_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A4_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A5.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/comments/S7.4_A6.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/comments/hashbang | 29 | 2 | 0 | 0 | 0 | 0 |
| language/comments/mongolian-vowel-separator-multi.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/mongolian-vowel-separator-single-eval.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/comments/mongolian-vowel-separator-single.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/multi-line-asi-carriage-return.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/multi-line-asi-line-feed.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/multi-line-asi-line-separator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/multi-line-asi-paragraph-separator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/multi-line-html-close-extra.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/comments/single-line-html-close-without-lt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/computed-property-names/basics | 3 | 0 | 0 | 0 | 0 | 0 |
| language/computed-property-names/class | 29 | 0 | 0 | 0 | 0 | 0 |
| language/computed-property-names/object | 12 | 0 | 0 | 0 | 0 | 0 |
| language/computed-property-names/to-name-side-effects | 4 | 0 | 0 | 0 | 0 | 0 |
| language/destructuring/binding | 19 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-1-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-10-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-11-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-12-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-13-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-14-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-28-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-29-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-2gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-3-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-30-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-31-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-32-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-4-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-5-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-5gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-6-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-7-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-8-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-8gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/10.1.1-9-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-1-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-10-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-11-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-12-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-13-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-14-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-15-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-16-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-17-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-2-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-3-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-4-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-4gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-5-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-5gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-6-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-7-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-8-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/14.1-9-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-decl-final-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-decl-inside-func-decl-parse.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-decl-inside-func-decl-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-decl-no-semi-parse.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-decl-no-semi-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-decl-not-first-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-decl-parse.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-decl-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-expr-final-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-expr-inside-func-decl-parse.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-expr-inside-func-decl-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-expr-no-semi-parse.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-expr-no-semi-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-expr-not-first-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-expr-parse.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/func-expr-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/get-accsr-inside-func-expr-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/get-accsr-not-first-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/get-accsr-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/set-accsr-inside-func-expr-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/set-accsr-not-first-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/directive-prologue/set-accsr-runtime.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/eval-code/direct | 286 | 286 | 0 | 0 | 0 | 0 |
| language/eval-code/indirect | 61 | 1 | 0 | 0 | 0 | 0 |
| language/export/escaped-as-export-specifier.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/export/escaped-default.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/export/escaped-from.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/addition | 48 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/array | 52 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/arrow-function | 343 | 10 | 0 | 0 | 0 | 0 |
| language/expressions/assignment | 485 | 2 | 0 | 0 | 0 | 0 |
| language/expressions/assignmenttargettype | 324 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/async-arrow-function | 60 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/async-function | 93 | 5 | 0 | 0 | 0 | 0 |
| language/expressions/async-generator | 623 | 6 | 0 | 0 | 0 | 0 |
| language/expressions/await | 22 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/bitwise-and | 30 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/bitwise-not | 16 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/bitwise-or | 30 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/bitwise-xor | 30 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/call | 92 | 16 | 0 | 0 | 0 | 0 |
| language/expressions/class | 4059 | 99 | 0 | 0 | 0 | 0 |
| language/expressions/coalesce | 24 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/comma | 6 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/compound-assignment | 454 | 20 | 0 | 0 | 0 | 0 |
| language/expressions/concatenation | 5 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/conditional | 22 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/delete | 69 | 2 | 0 | 0 | 0 | 0 |
| language/expressions/division | 45 | 2 | 0 | 0 | 0 | 0 |
| language/expressions/does-not-equals | 38 | 2 | 0 | 0 | 0 | 0 |
| language/expressions/dynamic-import | 1005 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/equals | 47 | 2 | 0 | 0 | 0 | 0 |
| language/expressions/exponentiation | 44 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/function | 264 | 13 | 0 | 0 | 0 | 0 |
| language/expressions/generators | 290 | 9 | 0 | 0 | 0 | 0 |
| language/expressions/greater-than | 49 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/greater-than-or-equal | 43 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/grouping | 9 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/import.meta | 22 | 3 | 0 | 0 | 0 | 0 |
| language/expressions/in | 36 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/instanceof | 43 | 12 | 0 | 0 | 0 | 0 |
| language/expressions/left-shift | 45 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/less-than | 45 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/less-than-or-equal | 47 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/logical-and | 18 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/logical-assignment | 78 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/logical-not | 19 | 3 | 0 | 0 | 0 | 0 |
| language/expressions/logical-or | 18 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/member-expression | 1 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/modulus | 40 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/multiplication | 40 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/new | 59 | 2 | 0 | 0 | 0 | 0 |
| language/expressions/new.target | 14 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/object | 1170 | 25 | 0 | 0 | 0 | 0 |
| language/expressions/optional-chaining | 38 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/postfix-decrement | 37 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/postfix-increment | 38 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/prefix-decrement | 34 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/prefix-increment | 33 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/property-accessors | 21 | 2 | 0 | 0 | 0 | 0 |
| language/expressions/relational | 1 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/right-shift | 37 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/strict-does-not-equals | 30 | 2 | 0 | 0 | 0 | 0 |
| language/expressions/strict-equals | 30 | 2 | 0 | 0 | 0 | 0 |
| language/expressions/subtraction | 38 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/super | 94 | 5 | 0 | 0 | 0 | 0 |
| language/expressions/tagged-template | 27 | 6 | 0 | 0 | 0 | 0 |
| language/expressions/tco-pos.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/template-literal | 57 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/this | 6 | 4 | 0 | 0 | 0 | 0 |
| language/expressions/typeof | 16 | 4 | 0 | 0 | 0 | 0 |
| language/expressions/unary-minus | 14 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/unary-plus | 17 | 2 | 0 | 0 | 0 | 0 |
| language/expressions/unsigned-right-shift | 45 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/void | 9 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/yield | 63 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-1-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-10-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-100-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-100gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-101-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-101gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-102-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-102gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-103.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-104.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-105.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-106.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-10gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-11-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-11gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-12-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-12gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-13-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-13gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-14-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-14gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-15-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-15gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-16-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-16gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-17-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-17gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-18gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-19-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-19gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-2-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-20-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-20gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-21-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-21gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-22-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-22gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-23-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-23gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-24-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-24gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-25-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-25gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-26-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-26gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-27-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-27gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-28-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-28gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-29-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-29gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-3-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-30-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-30gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-31-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-31gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-32-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-32gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-33-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-33gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-34-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-34gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-35-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-35gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-36-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-36gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-37-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-37gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-38-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-38gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-39-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-39gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-4-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-40-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-40gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-41-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-41gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-42-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-42gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-43-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-43gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-44-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-44gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-45-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-45gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-46-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-46gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-47-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-47gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-48-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-48gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-49-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-49gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-5-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-50-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-50gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-51-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-51gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-52-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-52gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-53-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-53gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-54-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-54gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-55-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-55gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-56-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-56gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-57-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-57gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-58-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-58gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-59-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-59gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-60-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-60gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-61-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-61gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-62-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-62gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-63-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-63gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-64-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-64gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-65-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-65gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-66-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-66gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-67-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-67gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-68-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-68gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-69-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-69gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-7-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-70-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-70gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-71-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-71gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-72-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-72gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-73-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-73gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-74-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-74gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-75-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-75gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-76-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-76gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-77-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-77gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-78-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-78gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-79-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-79gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-7gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-8-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-80-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-80gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-81-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-81gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-82-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-82gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-83-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-83gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-84-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-84gs.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-85-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-85gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-86-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-86gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-87-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-87gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-88-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-88gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-89-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-89gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-8gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-9-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-90-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-90gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-91-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-91gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-92-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-92gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-93-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-93gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-94-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-94gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-95-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-95gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-96-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-96gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-97-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-97gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-98-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-98gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-99-s.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-99gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/10.4.3-1-9gs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.1.6_A1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.2.1_A1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.2.1_A2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.2.1_A3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.2.1_A4_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.2.1_A4_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.2.1_A5.1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.2.1_A5.1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.2.1_A5.2_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.4.3_A1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.4A1.1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/S10.4_A1.1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/block-decl-onlystrict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/eval-param-env-with-computed-key.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/eval-param-env-with-prop-initializer.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/function-code/switch-case-decl-onlystrict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/function-code/switch-dflt-decl-onlystrict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/_implements.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/abstract.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/boolean.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/byte.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/char.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/debugger.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/double.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/enum.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/export.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/extends.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/final.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/float.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/goto.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/implement.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/implements-strict-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/implements-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/implements-titlecase.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/implements-uppercase.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/implements.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/implements0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/implementss.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/import.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/int.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/interface-strict-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/interface-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/interface.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/let-strict-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/let-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/long.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/native.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/package-strict-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/package-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/package.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/private-strict-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/private-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/private.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/protected-strict-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/protected-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/protected.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/public-strict-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/public-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/public.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/short.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/static-strict-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/static-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/static.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/super.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/synchronized.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/throws.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/transient.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/volatile.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/yield-strict-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/future-reserved-words/yield-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/S10.1.7_A1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/S10.4.1_A1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/S10.4.1_A1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/block-decl-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/decl-func-dup.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/decl-func.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/decl-lex-configurable-global.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/decl-lex-deletion.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/decl-lex-restricted-global.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/decl-lex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/decl-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/export.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/import.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/invalid-private-names-call-expression-bad-reference.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/invalid-private-names-call-expression-this.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/invalid-private-names-member-expression-bad-reference.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/invalid-private-names-member-expression-this.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/new.target-arrow.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/new.target.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/return.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-func-dups.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-func-err-non-configurable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-func-err-non-extensible.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-func.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-lex-deletion.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-lex-lex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-lex-restricted-global.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-lex-var-declared-via-eval.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-lex-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-lex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-var-collision.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-var-err.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/script-decl-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/super-call-arrow.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/super-call.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/super-prop-arrow.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/super-prop.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/switch-case-decl-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/switch-dflt-decl-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/unscopables-ignored.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/yield-non-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/global-code/yield-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S10.2.2_A1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S10.2.2_A1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S10.2.2_A1_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S10.2.2_A1_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S10.2.2_A1_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S10.2.2_A1_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S10.2.2_A1_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S10.2.2_A1_T8.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S10.2.2_A1_T9.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S11.1.2_A1_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/S11.1.2_A1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/assign-to-global-undefined.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/static-init-invalid-await.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifier-resolution/unscopables.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/other_id_continue-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/other_id_continue.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/other_id_start-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/other_id_start.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-digits-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-digits-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-digits.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-10.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-10.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-10.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-10.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-11.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-11.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-11.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-11.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-12.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-12.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-12.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-12.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-13.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-13.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-13.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-13.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-14.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-14.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-14.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-14.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-15.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-15.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-15.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-15.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-15.1.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-15.1.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-15.1.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-15.1.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-16.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-16.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-16.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-16.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-17.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-17.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-17.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-17.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-5.2.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-5.2.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-5.2.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-5.2.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-6.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-6.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-6.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-6.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-6.1.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-6.1.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-6.1.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-6.1.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-7.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-7.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-7.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-7.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-8.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-8.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-8.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-8.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-9.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-9.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-9.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-unicode-9.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/part-zwj-zwnj-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-dollar-sign.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-escape-seq.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-underscore.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-10.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-10.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-10.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-10.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-11.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-11.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-11.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-11.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-12.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-12.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-12.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-12.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-13.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-13.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-13.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-13.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-14.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-14.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-14.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-14.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-15.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-15.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-15.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-15.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-15.1.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-15.1.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-15.1.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-15.1.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-16.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-16.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-16.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-16.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-17.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-17.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-17.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-17.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-5.2.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-5.2.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-5.2.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-5.2.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-6.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-6.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-6.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-6.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-6.1.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-6.1.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-6.1.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-6.1.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-7.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-7.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-7.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-7.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-8.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-8.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-8.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-8.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-9.0.0-class-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-9.0.0-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-9.0.0-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-9.0.0.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-unicode-ltr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-zwj-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/start-zwnj-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/unicode-escape-nls-err.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-break-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-break-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-break.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-case-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-case-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-case.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-catch-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-catch-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-catch.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-class-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-class-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-class.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-const-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-const-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-continue-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-continue-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-continue.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-debugger-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-debugger-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-debugger.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-default-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-default-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-default.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-delete-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-delete-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-delete.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-do-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-do-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-do.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-dollar-sign-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-dollar-sign-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-dollar-sign.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-else-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-else-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-else.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-enum-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-enum-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-enum.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-export-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-export-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-export.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-extends-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-extends-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-extends.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-false-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-false-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-false.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-finally-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-finally-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-finally.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-for-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-for-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-for.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-function-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-function-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-function.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-if-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-if-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-if.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-import-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-import-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-import.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-in-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-in-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-in.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-instanceof-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-instanceof-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-instanceof.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-new-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-new-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-new.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-null-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-null-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-return-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-return-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-return.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-super-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-super-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-super.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-switch-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-switch-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-switch.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-this-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-this-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-this.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-throw-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-throw-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-throw.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-true-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-true-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-true.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-try-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-try-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-try.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-typeof-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-typeof-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-typeof.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-underscore-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-underscore-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-underscore.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-var-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-var-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-void-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-void-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-void.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-while-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-while-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-while.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-with-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-with-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-with.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/val-yield-strict.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-cjk-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-cjk.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-eng-alpha-lower-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-eng-alpha-lower-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-eng-alpha-lower.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-eng-alpha-upper-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-eng-alpha-upper-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-eng-alpha-upper.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-rus-alpha-lower-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-rus-alpha-lower-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-rus-alpha-lower.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-rus-alpha-upper-via-escape-hex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-rus-alpha-upper-via-escape-hex4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vals-rus-alpha-upper.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vertical-tilde-continue-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vertical-tilde-continue.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vertical-tilde-start-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/identifiers/vertical-tilde-start.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/import/dup-bound-names.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/import/escaped-as-import-specifier.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/import/escaped-as-namespace-import.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/import/escaped-from.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/import/import-attributes | 17 | 0 | 0 | 0 | 0 | 0 |
| language/import/import-bytes | 5 | 0 | 0 | 0 | 0 | 0 |
| language/import/import-defer | 103 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-break.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-case.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-catch.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-continue.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-default.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-delete.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-do.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-else.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-finally.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-for.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-function.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-if.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-in.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-instanceof.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-new.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-return.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-switch.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-this.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-throw.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-try.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-typeof.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-void.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-while.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/keywords/ident-ref-with.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/7.3-15.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/7.3-5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/7.3-6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A2.1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A2.2_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A3.2_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A5.4.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A6_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A6_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A6_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A6_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A7_T1.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A7_T2.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A7_T3.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A7_T4.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A7_T5.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A7_T6.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A7_T7.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/line-terminators/S7.3_A7_T8.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/line-terminators/between-tokens-cr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/between-tokens-lf.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/between-tokens-ls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/between-tokens-ps.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/comment-multi-cr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/comment-multi-lf.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/comment-multi-ls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/comment-multi-ps.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/comment-single-cr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/comment-single-lf.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/comment-single-ls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/comment-single-ps.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/invalid-comment-single-cr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/invalid-comment-single-lf.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/invalid-comment-single-ls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/invalid-comment-single-ps.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/invalid-regexp-cr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/invalid-regexp-lf.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/invalid-regexp-ls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/invalid-regexp-ps.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/invalid-string-cr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/line-terminators/invalid-string-lf.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/literals/bigint | 59 | 0 | 0 | 0 | 0 | 0 |
| language/literals/boolean | 4 | 0 | 0 | 0 | 0 | 0 |
| language/literals/null | 3 | 1 | 0 | 0 | 0 | 0 |
| language/literals/numeric | 157 | 1 | 0 | 0 | 0 | 0 |
| language/literals/regexp | 240 | 24 | 0 | 0 | 0 | 0 |
| language/literals/string | 73 | 3 | 0 | 0 | 0 | 0 |
| language/module-code/ambiguous-export-bindings | 10 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/comment-multi-line-html-close.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/comment-single-line-html-close.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/comment-single-line-html-open.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-export-as-star-as.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-export-decl.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-export-dflt-id.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-export-dflt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-export-id-as.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-export-id.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-export-star-as-dflt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-lables.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-lex.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-top-function-async-generator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-top-function-async.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-top-function-generator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-dup-top-function.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-export-global.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-export-ill-formed-string.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-export-unresolvable.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-import-arguments.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-import-as-arguments.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-import-as-eval.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-import-eval.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-lex-and-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-new-target.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-strict-mode.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-super.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-undef-break.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/early-undef-continue.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-cls-semi.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-cls-anon-semi.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-cls-anon.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-cls-name-meth.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-cls-named-semi.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-cls-named.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-expr-cls-anon.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-expr-cls-name-meth.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-expr-cls-named.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-expr-err-eval.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-expr-err-get-value.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-expr-fn-anon.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-expr-fn-named.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-expr-gen-anon.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-expr-gen-named.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-expr-in.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-fun-anon-semi.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-fun-named-semi.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-gen-anon-semi.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-dflt-gen-named-semi.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-fun-semi.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-export-gen-semi.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-gtbndng-indirect-trlng-comma.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-gtbndng-indirect-update-as.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-gtbndng-indirect-update-dflt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-gtbndng-indirect-update.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-gtbndng-local-bndng-cls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-gtbndng-local-bndng-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-gtbndng-local-bndng-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-gtbndng-local-bndng-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-rqstd-abrupt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-rqstd-once.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-rqstd-order.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-self-abrupt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-self-once.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/eval-this.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-default-asyncfunction-declaration-binding-exists.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-default-asyncfunction-declaration-binding.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-default-asyncgenerator-declaration-binding-exists.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-default-asyncgenerator-declaration-binding.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-default-function-declaration-binding-exists.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-default-function-declaration-binding.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-default-generator-declaration-binding-exists.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-default-generator-declaration-binding.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-binding-index.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-binding-string.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-from-as-unpaired-surrogate.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-from-binding-string.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-from-star-string.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-from-star-unpaired-surrogate.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-from-star.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-from-string-binding.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-from-string-string.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-from-string.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-from-unpaired-surrogate.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-import-string-binding.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-import-unpaired-surrogate.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-string-binding.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-expname-unpaired-surrogate.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/export-star-as-dflt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/import-attributes | 13 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-bndng-cls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-bndng-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-bndng-fun.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-bndng-gen.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-bndng-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-bndng-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-err-circular-as.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-err-circular.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-err-dflt-thru-star-as.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-err-dflt-thru-star.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-err-not-found-as.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-err-not-found.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-iee-cycle.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-star-cycle.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-iee-trlng-comma.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-cls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-export-cls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-export-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-export-fun.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-export-gen.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-export-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-export-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-for-dup.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-for.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-fun.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-gen.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-var-dup.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-local-bndng-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-cls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-dflt-cls.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-dflt-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-dflt-fun-anon.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-dflt-fun-named.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-dflt-gen-anon.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-dflt-gen-named.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-dflt-named.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-dflt-star.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-fun.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-gen.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-trlng-comma.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-bndng-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-err-dflt-thru-star-as.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-err-dflt-thru-star-dflt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-err-not-found-as.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-err-not-found-dflt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-err-not-found.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-id-name.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-iee-cycle.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-named-star-cycle.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-once.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-resolve-empty-export.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-resolve-empty-import.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-resolve-err-syntax-1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-resolve-err-syntax-2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-resolve-order-depth.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-resolve-order-src.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-same-global.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-as-props-dflt-skip.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-binding.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-equality.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-err-not-found.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-id-name.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-iee-cycle.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-iee-multi-cycle-same-name.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-iee-single-cycle-same-name.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-props-circular.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-props-dflt-keep-indirect.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-props-dflt-keep-local.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-props-dflt-skip.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-props-nrml.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-star-star-cycle.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/instn-uniq-env-rec.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/invalid-private-names-call-expression-bad-reference.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/invalid-private-names-call-expression-this.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/invalid-private-names-member-expression-bad-reference.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/invalid-private-names-member-expression-this.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/namespace | 38 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-arrow-function.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-block-stmt-list.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-block-stmt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-class-decl-meth-static.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-class-decl-meth.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-class-decl-method-gen-static.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-class-decl-method-gen.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-class-expr-meth-gen-static.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-class-expr-meth-gen.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-class-expr-meth-static.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-class-expr-meth.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-do-while.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-in-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-in-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-in-lhs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-in-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-lhs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-of-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-of-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-of-lhs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-of-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-for-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-function-decl.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-function-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-generator-decl.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-generator-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-if-else.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-if-if.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-labeled.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-object-gen-method.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-object-getter.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-object-method.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-object-setter.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-switch-case-dflt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-switch-case.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-switch-dftl.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-try-catch-finally.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-try-catch.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-try-finally.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-try-try.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-export-while.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-arrow-function.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-block-stmt-list.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-block-stmt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-class-decl-meth-static.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-class-decl-meth.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-class-decl-method-gen-static.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-class-decl-method-gen.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-class-expr-meth-gen-static.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-class-expr-meth-gen.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-class-expr-meth-static.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-class-expr-meth.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-do-while.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-in-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-in-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-in-lhs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-in-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-lhs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-of-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-of-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-of-lhs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-of-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-for-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-function-decl.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-function-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-generator-decl.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-generator-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-if-else.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-if-if.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-labeled.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-object-gen-method.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-object-getter.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-object-method.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-object-setter.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-switch-case-dflt.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-switch-case.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-switch-dftl.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-try-catch-finally.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-try-catch.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-try-finally.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-try-try.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-decl-pos-import-while.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-export-dflt-const.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-export-dflt-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-export-dflt-let.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-export-dflt-var.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-hoist-lex-fun.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-hoist-lex-gen.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-invoke-anon-fun-decl.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-invoke-anon-gen-decl.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-return.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-semi-dflt-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-semi-export-star.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-semi-name-space-export.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-semi-named-export-from.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-semi-named-export.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-syntax-1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-syntax-2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-err-yield.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/parse-export-empty.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/private-identifiers-not-empty.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/privatename-not-valid-earlyerr-module-1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/privatename-not-valid-earlyerr-module-2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/privatename-not-valid-earlyerr-module-3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/privatename-not-valid-earlyerr-module-4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/privatename-not-valid-earlyerr-module-5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/privatename-not-valid-earlyerr-module-6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/privatename-not-valid-earlyerr-module-7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/privatename-not-valid-earlyerr-module-8.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/privatename-valid-no-earlyerr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/source-phase-import | 3 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/top-level-await | 251 | 0 | 0 | 0 | 0 | 0 |
| language/module-code/verify-dfs.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A2_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A2_T10.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A2_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A2_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A2_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A2_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A2_T6.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A2_T7.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A2_T8.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/punctuators/S7.7_A2_T9.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/await-module.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/await-script.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-global-property-accessor.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-global-property-memberexpr-str.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-global-property-memberexpr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-global-property-prop-name.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-keyword-accessor.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-keyword-memberexpr-str.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-keyword-memberexpr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-keyword-prop-name.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-reserved-word-literal-accessor.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-reserved-word-literal-memberexpr-str.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-reserved-word-literal-memberexpr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-name-reserved-word-literal-prop-name.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-reference-false-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-reference-false.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-reference-null-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-reference-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-reference-true-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/ident-reference-true.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/label-ident-false-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/label-ident-false.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/label-ident-null-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/label-ident-null.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/label-ident-true-escaped.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/label-ident-true.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/reserved-words/unreserved-words.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/array-pattern.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/arrow-function.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/expected-argument-count.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/no-alias-arguments.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/object-pattern.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/position-invalid.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/rest-index.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/rest-parameters-apply.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/rest-parameters-call.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/rest-parameters-produce-an-array.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/rest-parameters/with-new-target.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/source-text/6.1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-array-literal-with-item.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-array-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-arrow-function-assignment-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-arrow-function-functionbody.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-block-with-labels.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-block.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-expr-arrow-function-boolean-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-let-declaration.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-regexp-literal-flags.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-regexp-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-with-statment-array-literal-with-item.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-with-statment-array-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-with-statment-arrow-function-assignment-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-with-statment-arrow-function-functionbody.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-with-statment-block-with-labels.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-with-statment-block.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-with-statment-expr-arrow-function-boolean-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-with-statment-let-declaration.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-with-statment-regexp-literal-flags.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/block-with-statment-regexp-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/class-array-literal-with-item.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/class-array-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/class-arrow-function-assignment-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/class-arrow-function-functionbody.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/class-block-with-labels.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/class-block.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/class-expr-arrow-function-boolean-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/class-let-declaration.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/class-regexp-literal-flags.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/class-regexp-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-array-literal-with-item.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-array-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-arrow-function-assignment-expr.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-arrow-function-functionbody.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-block-with-labels.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-block.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-expr-arrow-function-boolean-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-let-declaration.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-regexp-literal-flags.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-regexp-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-with-statment-array-literal-with-item.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-with-statment-array-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-with-statment-arrow-function-assignment-expr.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-with-statment-arrow-function-functionbody.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-with-statment-block-with-labels.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-with-statment-block.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-with-statment-expr-arrow-function-boolean-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-with-statment-let-declaration.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-with-statment-regexp-literal-flags.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-block-with-statment-regexp-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-class-array-literal-with-item.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-class-array-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-class-arrow-function-assignment-expr.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-class-arrow-function-functionbody.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-class-block-with-labels.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-class-block.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-class-expr-arrow-function-boolean-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-class-let-declaration.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-class-regexp-literal-flags.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-class-regexp-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-fn-array-literal-with-item.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-fn-array-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-fn-arrow-function-assignment-expr.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-fn-arrow-function-functionbody.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-fn-block-with-labels.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-fn-block.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-fn-expr-arrow-function-boolean-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-fn-let-declaration.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-fn-regexp-literal-flags.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/eval-fn-regexp-literal.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/statementList/fn-array-literal-with-item.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/fn-array-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/fn-arrow-function-assignment-expr.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/fn-arrow-function-functionbody.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/fn-block-with-labels.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/fn-block.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/fn-expr-arrow-function-boolean-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/fn-let-declaration.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/fn-regexp-literal-flags.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statementList/fn-regexp-literal.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/statements/async-function | 74 | 3 | 0 | 0 | 0 | 0 |
| language/statements/async-generator | 301 | 1 | 0 | 0 | 0 | 0 |
| language/statements/await-using | 98 | 1 | 0 | 0 | 0 | 0 |
| language/statements/block | 21 | 0 | 0 | 0 | 0 | 0 |
| language/statements/break | 20 | 1 | 0 | 0 | 0 | 0 |
| language/statements/class | 4367 | 109 | 0 | 0 | 0 | 0 |
| language/statements/const | 136 | 1 | 0 | 0 | 0 | 0 |
| language/statements/continue | 24 | 1 | 0 | 0 | 0 | 0 |
| language/statements/debugger | 2 | 0 | 0 | 0 | 0 | 0 |
| language/statements/do-while | 36 | 6 | 0 | 0 | 0 | 0 |
| language/statements/empty | 2 | 1 | 0 | 0 | 0 | 0 |
| language/statements/expression | 3 | 2 | 0 | 0 | 0 | 0 |
| language/statements/for | 385 | 9 | 0 | 0 | 0 | 0 |
| language/statements/for-await-of | 1234 | 0 | 0 | 0 | 0 | 0 |
| language/statements/for-in | 119 | 13 | 0 | 0 | 0 | 0 |
| language/statements/for-of | 751 | 7 | 0 | 0 | 0 | 0 |
| language/statements/function | 451 | 54 | 0 | 0 | 0 | 0 |
| language/statements/generators | 266 | 7 | 0 | 0 | 0 | 0 |
| language/statements/if | 69 | 9 | 0 | 0 | 0 | 0 |
| language/statements/labeled | 24 | 2 | 0 | 0 | 0 | 0 |
| language/statements/let | 145 | 1 | 0 | 0 | 0 | 0 |
| language/statements/return | 16 | 0 | 0 | 0 | 0 | 0 |
| language/statements/switch | 111 | 23 | 0 | 0 | 0 | 0 |
| language/statements/throw | 14 | 1 | 0 | 0 | 0 | 0 |
| language/statements/try | 201 | 16 | 0 | 0 | 0 | 0 |
| language/statements/using | 80 | 2 | 0 | 0 | 0 | 0 |
| language/statements/variable | 178 | 14 | 0 | 0 | 0 | 0 |
| language/statements/while | 38 | 7 | 0 | 0 | 0 | 0 |
| language/statements/with | 181 | 20 | 0 | 0 | 0 | 0 |
| language/types/boolean | 5 | 0 | 0 | 0 | 0 | 0 |
| language/types/list | 3 | 0 | 0 | 0 | 0 | 0 |
| language/types/null | 4 | 0 | 0 | 0 | 0 | 0 |
| language/types/number | 21 | 0 | 0 | 0 | 0 | 0 |
| language/types/object | 19 | 0 | 0 | 0 | 0 | 0 |
| language/types/reference | 29 | 3 | 0 | 0 | 0 | 0 |
| language/types/string | 24 | 4 | 0 | 0 | 0 | 0 |
| language/types/undefined | 8 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A2.1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A2.2_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A2.3_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A2.4_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A2.5_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A3.1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A3.2_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A3.3_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A3.4_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A3.5_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A4.1_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A4.2_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A4.3_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A4.4_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A4.5_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A5_T1.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A5_T2.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A5_T3.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A5_T4.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/S7.2_A5_T5.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-carriage-return.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-em-quad.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-em-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-en-quad.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-en-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-figure-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-form-feed.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-four-per-em-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-hair-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-ideographic-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-line-feed.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-line-separator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-medium-mathematical-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-nbsp.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-nnbsp.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-ogham-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-paragraph-separator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-punctuation-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-six-per-em-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-tab.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-thin-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-three-per-em-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-vertical-tab.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/after-regular-expression-literal-zwnbsp.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/between-form-feed.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/between-horizontal-tab.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/between-nbsp.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/between-space.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/between-vertical-tab.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/comment-multi-form-feed.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/comment-multi-horizontal-tab.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/comment-multi-nbsp.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/comment-multi-space.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/comment-multi-vertical-tab.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/comment-single-form-feed.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/comment-single-horizontal-tab.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/comment-single-nbsp.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/comment-single-space.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/comment-single-vertical-tab.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/mongolian-vowel-separator-eval.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/mongolian-vowel-separator.js | 1 | 0 | 0 | 0 | 0 | 0 |
| language/white-space/string-form-feed.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/string-horizontal-tab.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/string-nbsp.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/string-space.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/white-space/string-vertical-tab.js | 1 | 1 | 0 | 0 | 0 | 0 |
| staging/Temporal | 2 | 0 | 0 | 0 | 0 | 0 |
| staging/Uint8Array | 1 | 0 | 0 | 0 | 0 | 0 |
| staging/built-ins | 8 | 0 | 8 | 0 | 8 | 0 |
| staging/decorators | 3 | 2 | 0 | 0 | 0 | 0 |
| staging/explicit-resource-management | 53 | 0 | 0 | 0 | 0 | 0 |
| staging/set-is-subset-of-empty-index.js | 1 | 0 | 0 | 0 | 0 | 0 |
| staging/set-is-subset-on-set-like.js | 1 | 0 | 0 | 0 | 0 | 0 |
| staging/set-is-subset-table-receiver-cleared.js | 1 | 0 | 0 | 0 | 0 | 0 |
| staging/set-is-subset-table-transition.js | 1 | 0 | 0 | 0 | 0 | 0 |
| staging/set-methods | 3 | 0 | 0 | 0 | 0 | 0 |
| staging/sm | 1406 | 269 | 0 | 0 | 0 | 0 |
| staging/source-phase-imports | 2 | 0 | 0 | 0 | 0 | 0 |
| staging/top-level-await | 1 | 0 | 0 | 0 | 0 | 0 |

## Failure checklist

### annexB/built-ins

- [ ] `annexB/built-ins/RegExp/legacy-accessors/index/this-subclass-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/index/this-subclass-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/leftContext/this-cross-realm-constructor.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/leftContext/this-cross-realm-constructor.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastParen/this-subclass-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastParen/this-subclass-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/input/this-subclass-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/input/this-subclass-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/prototype/flags/order-after-compile.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `annexB/built-ins/RegExp/prototype/flags/order-after-compile.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/input/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/input/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastMatch/this-subclass-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastMatch/this-subclass-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/index/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/index/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastMatch/this-not-regexp-constructor.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastMatch/this-not-regexp-constructor.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/index/this-not-regexp-constructor.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/index/this-not-regexp-constructor.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastParen/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastParen/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/leftContext/this-not-regexp-constructor.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/leftContext/this-not-regexp-constructor.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/rightContext/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/input/this-cross-realm-constructor.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/input/this-cross-realm-constructor.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/rightContext/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/prototype/Symbol.split/toint32-limit-recompiles-source.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected a number

```

- [ ] `annexB/built-ins/RegExp/prototype/Symbol.split/toint32-limit-recompiles-source.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected a number

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/index/this-cross-realm-constructor.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/index/this-cross-realm-constructor.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/rightContext/this-not-regexp-constructor.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/rightContext/this-not-regexp-constructor.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastParen/this-cross-realm-constructor.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastParen/this-cross-realm-constructor.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/rightContext/this-subclass-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/rightContext/this-subclass-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastParen/this-not-regexp-constructor.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastParen/this-not-regexp-constructor.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/named-groups/non-unicode-malformed.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `annexB/built-ins/RegExp/named-groups/non-unicode-malformed.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/input/this-not-regexp-constructor.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/input/this-not-regexp-constructor.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastMatch/this-cross-realm-constructor.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastMatch/this-cross-realm-constructor.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/leftContext/this-subclass-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/leftContext/this-subclass-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/rightContext/this-cross-realm-constructor.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/rightContext/this-cross-realm-constructor.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastMatch/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/lastMatch/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/leftContext/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/leftContext/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

### built-ins/RegExp

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_N.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_N.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-07.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid character in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-07.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid character in character class"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-union-character-class-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[dd]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-union-character-class-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[dd]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/multiline/this-val-non-obj.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/multiline/this-val-non-obj.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-26.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-26.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-13-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 13: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-13-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 13: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/exec/name.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Default_Ignorable_Code_Point-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 38: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Default_Ignorable_Code_Point-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 38: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/exec/name.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-02.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-02.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-difference-character-property-escape.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-difference-character-property-escape.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/prototype/global/S15.10.7.2_A10.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/global/S15.10.7.2_A10.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-union-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-union-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/builtin-success-y-set-lastindex-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/builtin-success-y-set-lastindex-err.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/named-groups/unicode-property-names-valid.js` (sloppy): **emit-error**

```text
xBacktrack }\ngoto rxBacktrack\nrxBacktrack:\nif len(stack) \u003c rxStackFrameSlots { return 0, 0, false }\nrxFrame = len(stack)-rxStackFrameSlots; choiceTarget = int(stack[rxFrame]); pos = int(stack[rxFrame+1]); for rxIndex := range regs { regs[rxIndex] = int(stack[rxFrame+2+rxIndex]) }; stack = stack[:rxFrame]\nswitch choiceTarget {\ncase 1:\ngoto rxNode2\ncase 2:\ngoto rxNode9\ncase 3:\ngoto rxNode14\ncase 4:\ngoto rxNode13\ncase 5:\ngoto rxNode12\ncase 6:\ngoto rxNode21\ncase 7:\ngoto rxNode7\ndefault: return 0, 0, false\n}\n}\n"}
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: format generated Go matcher: 73:379: expected ';', found 'if'\nconst tsRegexMatcher887CaptureSlots = 8\nfunc tsRegexMatcher887(input []uint16, start int) (matchStart, matchEnd int, captures []int, ok bool) {\noffsets := make([]int32, tsRegexMatcher887CaptureSlots)\nmatchStart, matchEnd, ok = tsRegexMatcher887Into(input, start, offsets)\nif !ok { return 0, 0, nil, false }; captures = make([]int, len(offsets)); for index, offset := range offsets { captures[index] = int(offset) }; return matchStart, matchEnd, captures, true\n}\nfunc tsRegexMatcher887Into(input []uint16, start int, captures []int32) (matchStart, matchEnd int, ok bool) {\nconst captureSlots = 8\nif start \u003c 0 { start = 0 }; if start \u003e len(input) { return 0, 0, false }; if len(captures) \u003c captureSlots { return 0, 0, false }\nregs := make([]int, 10)\nfor i := range regs { regs[i] = -1 }\npos := start\nattemptStart := start\nconst rxStackFrameSlots = 12\nconst rxMaxStackSlots = 64 * 1024 * 1024 /
```

- [ ] `built-ins/RegExp/named-groups/unicode-property-names-valid.js` (strict): **emit-error**

```text
xBacktrack }\ngoto rxBacktrack\nrxBacktrack:\nif len(stack) \u003c rxStackFrameSlots { return 0, 0, false }\nrxFrame = len(stack)-rxStackFrameSlots; choiceTarget = int(stack[rxFrame]); pos = int(stack[rxFrame+1]); for rxIndex := range regs { regs[rxIndex] = int(stack[rxFrame+2+rxIndex]) }; stack = stack[:rxFrame]\nswitch choiceTarget {\ncase 1:\ngoto rxNode2\ncase 2:\ngoto rxNode9\ncase 3:\ngoto rxNode14\ncase 4:\ngoto rxNode13\ncase 5:\ngoto rxNode12\ncase 6:\ngoto rxNode21\ncase 7:\ngoto rxNode7\ndefault: return 0, 0, false\n}\n}\n"}
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: format generated Go matcher: 73:379: expected ';', found 'if'\nconst tsRegexMatcher888CaptureSlots = 8\nfunc tsRegexMatcher888(input []uint16, start int) (matchStart, matchEnd int, captures []int, ok bool) {\noffsets := make([]int32, tsRegexMatcher888CaptureSlots)\nmatchStart, matchEnd, ok = tsRegexMatcher888Into(input, start, offsets)\nif !ok { return 0, 0, nil, false }; captures = make([]int, len(offsets)); for index, offset := range offsets { captures[index] = int(offset) }; return matchStart, matchEnd, captures, true\n}\nfunc tsRegexMatcher888Into(input []uint16, start int, captures []int32) (matchStart, matchEnd int, ok bool) {\nconst captureSlots = 8\nif start \u003c 0 { start = 0 }; if start \u003e len(input) { return 0, 0, false }; if len(captures) \u003c captureSlots { return 0, 0, false }\nregs := make([]int, 10)\nfor i := range regs { regs[i] = -1 }\npos := start\nattemptStart := start\nconst rxStackFrameSlots = 12\nconst rxMaxStackSlots = 64 * 1024 * 1024 /
```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A2_T9.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A2_T9.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/test/y-fail-lastindex-no-write.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/test/y-fail-lastindex-no-write.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/lookBehind/nested-lookaround.js` (sloppy): **emit-error**

```text
panic: unknown V8 graph node in FirstCharacterSetBuilder [recovered, repanicked]

goroutine 155 [running]:
sync.(*WaitGroup).Go.func1.1()
	/workspace/toolchains/go/src/sync/waitgroup.go:251 +0x45
panic({0x1b48ca8?, 0x1289850?})
	/workspace/toolchains/go/src/runtime/panic.go:859 +0x125
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8CompilerFirstCharacterSetBuilder).addNode(0x3?, {0x1be96d0?, 0x38cd5cd1570?}, 0x38cd5ccba18?)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_compiler_first_char.go:53 +0x1b9
github.com/microsoft/TypeScript/tsc/internal/regexpengine.v8CompilerComputeFirstCharacterSet({0x1be96d0?, 0x38cd5cd1570?}, 0x38cd5c75780?)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_compiler_first_char.go:18 +0x5b
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8GoGraphEmitter).emitChoice(0x38cd5c9a3a0, 0x1be9870?, 0xd)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:754 +0x10f
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8GoGraphEmitter).emitNode(0x38cd5c9a3a0, 0xd, {0x1be9940, 0x38cd5c75780})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:524 +0x7fc
github.com/microsoft/TypeScript/tsc/internal/regexpengine.v8EmitGoMatcher(0x38cd5d42600, {0x38cd5c653c8, 0x11})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:243 +0x2b9
github.com/microsoft/TypeScript/tsc/internal/regexpengine.CompileUTF16ToGo({0x38cd5c627d0, 0x24, 0x24}, {0x38cd45b3c40, 0x0}, {0x38cd5c653c8, 0x11})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:43 +0x3e6
gi
```

- [ ] `built-ins/RegExp/lookBehind/nested-lookaround.js` (strict): **emit-error**

```text
panic: unknown V8 graph node in FirstCharacterSetBuilder [recovered, repanicked]

goroutine 154 [running]:
sync.(*WaitGroup).Go.func1.1()
	/workspace/toolchains/go/src/sync/waitgroup.go:251 +0x45
panic({0x1b48ca8?, 0x1289850?})
	/workspace/toolchains/go/src/runtime/panic.go:859 +0x125
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8CompilerFirstCharacterSetBuilder).addNode(0x3?, {0x1be96d0?, 0x10b14fbad030?}, 0x10b14fb853e0?)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_compiler_first_char.go:53 +0x1b9
github.com/microsoft/TypeScript/tsc/internal/regexpengine.v8CompilerComputeFirstCharacterSet({0x1be96d0?, 0x10b14fbad030?}, 0x10b14f230d80?)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_compiler_first_char.go:18 +0x5b
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8GoGraphEmitter).emitChoice(0x10b14fb1a3a0, 0x1be9870?, 0xd)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:754 +0x10f
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8GoGraphEmitter).emitNode(0x10b14fb1a3a0, 0xd, {0x1be9940, 0x10b14f230d80})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:524 +0x7fc
github.com/microsoft/TypeScript/tsc/internal/regexpengine.v8EmitGoMatcher(0x10b14fbf0600, {0x10b14f22ed20, 0x11})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:243 +0x2b9
github.com/microsoft/TypeScript/tsc/internal/regexpengine.CompileUTF16ToGo({0x10b14fb96410, 0x24, 0x24}, {0x10b14e411c4e, 0x0}, {0x10b14f22ed20, 0x11})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:
```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-difference-character-property-escape.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-difference-character-property-escape.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-result-coerce-length.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-result-coerce-length.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A2_T5.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A2_T5.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}

```

- [ ] `built-ins/RegExp/property-escapes/non-existent-binary-property-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/non-existent-binary-property-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-intersection-character-class-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}&&d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-intersection-character-class-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}&&d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/property-escapes/generated/Assigned.js` (strict): **build-error**

```text
command-line-arguments: /workspace/toolchains/go/pkg/tool/linux_amd64/compile: signal: killed

```

- [ ] `built-ins/RegExp/property-escapes/generated/Assigned.js` (sloppy): **build-timeout**

```text
(no output)
```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A6.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A6.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/get-global-err.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/get-global-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-General_Category.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 20: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-General_Category.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 20: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-Script-equals.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-Script-equals.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/coerce-flags-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/coerce-flags-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/source/this-val-non-obj.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/source/this-val-non-obj.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/toString/S15.10.6.4_A6.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/toString/S15.10.6.4_A6.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-difference-character.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}--_]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-difference-character.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}--_]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_ZWJ_Sequence-negative-CharacterClass.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 29: Negated character class may contain strings"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_ZWJ_Sequence-negative-CharacterClass.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 29: Negated character class may contain strings"}

```

- [ ] `built-ins/RegExp/lookBehind/sticky.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Actual argument [null] shouldn't be primitive. #1

```

- [ ] `built-ins/RegExp/lookBehind/sticky.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Actual argument [null] shouldn't be primitive. #1

```

- [ ] `built-ins/RegExp/prototype/unicode/length.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/unicode/length.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/generated/ID_Continue.js` (sloppy): **build-error**

```text
command-line-arguments: /workspace/toolchains/go/pkg/tool/linux_amd64/compile: signal: killed

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-intersection-character-class-escape.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-intersection-character-class-escape.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-08.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 7: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-08.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 7: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-intersection-character-property-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]&&p{ASCII_Hex_Digit}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-intersection-character-property-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]&&p{ASCII_Hex_Digit}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A9.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A9.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_Flag_Sequence-negative-u.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 27: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_Flag_Sequence-negative-u.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 27: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/hasIndices/this-val-regexp-prototype.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-In-prefix-Script.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 18: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-In-prefix-Script.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 18: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/hasIndices/this-val-regexp-prototype.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-constructor.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:24353:33: tsTemp852.read undefined (type *tsCell has no field or method read)
<case>/main.go:24355:23: cannot use tsTemp852.set(loop, tsNumberValue(tsTemp854)) (value of struct type tsValue) as float64 value in assignment

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-constructor.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:24356:33: tsTemp853.read undefined (type *tsCell has no field or method read)
<case>/main.go:24358:23: cannot use tsTemp853.set(loop, tsNumberValue(tsTemp855)) (value of struct type tsValue) as float64 value in assignment

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-character-class-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-character-class-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Prepended_Concatenation_Mark.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 32: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Prepended_Concatenation_Mark.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 32: Invalid property name"}

```

- [ ] `built-ins/RegExp/S15.10.2.13_A2_T2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/S15.10.2.13_A2_T2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-difference-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}--q{0|2|4|9️⃣}]` should match 1 (U+000049)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-difference-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}--q{0|2|4|9️⃣}]` should match 1 (U+000049)

```

- [ ] `built-ins/RegExp/lookBehind/do-not-backtrack.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/do-not-backtrack.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-difference-character-class-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}--d]` should match A (U+000065)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-difference-character-class-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}--d]` should match A (U+000065)

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-09.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-09.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/get-exec-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/get-exec-err.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/Basic_Emoji-negative-u.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 15: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/Basic_Emoji-negative-u.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 15: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/multiline/length.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/multiline/length.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-set-lastindex-match.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/sticky/cross-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/sticky/cross-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-set-lastindex-match.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-20.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-20.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A1_T11.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A1_T11.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/g-init-lastindex-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/g-init-lastindex-err.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/builtin-success-g-set-lastindex.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/builtin-success-g-set-lastindex.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/global/name.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/global/name.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/source/cross-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/source/cross-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-intersection-character-class-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d&&d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-intersection-character-class-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d&&d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-intersection-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d&&q{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-intersection-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d&&q{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-property-of-strings-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}p{Emoji_Keycap_Sequence}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-property-of-strings-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}p{Emoji_Keycap_Sequence}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/lookBehind/captures-negative.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/captures-negative.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-difference-property-of-strings-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]--p{Emoji_Keycap_Sequence}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-difference-property-of-strings-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]--p{Emoji_Keycap_Sequence}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-Script_Extensions-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 21: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-Script_Extensions-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 21: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-unclosed.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-unclosed.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/flags/coercion-dotall.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/flags/coercion-dotall.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-FC_NFKC_Closure.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-FC_NFKC_Closure.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/lookBehind/misc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/misc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/prototype/hasIndices/this-val-invalid-obj.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/hasIndices/this-val-invalid-obj.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-separator-and-value-only-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-separator-and-value-only-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-11.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 23: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-11.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 23: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-16.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-16.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-intersection-character-class.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-intersection-character-class.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-difference-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]--q{0|2|4|9️⃣}]` should match 1 (U+000049)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-difference-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]--q{0|2|4|9️⃣}]` should match 1 (U+000049)

```

- [ ] `built-ins/RegExp/prototype/sticky/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/sticky/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/g-init-lastindex-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/g-init-lastindex-err.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-Script_Extensions-equals.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 22: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/this-get-flags-throws.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-Script_Extensions-equals.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 22: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/this-get-flags-throws.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-12.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 22: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-12.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 22: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/Basic_Emoji-negative-CharacterClass.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 18: Negated character class may contain strings"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/Basic_Emoji-negative-CharacterClass.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 18: Negated character class may contain strings"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-union-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-union-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/multiline/S15.10.7.4_A10.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/multiline/S15.10.7.4_A10.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/toString/S15.10.6.4_A7.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/toString/S15.10.6.4_A7.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/match-indices/indices-array-properties.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: 0 descriptor should be configurable

```

- [ ] `built-ins/RegExp/match-indices/indices-array-properties.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: 0 descriptor should be configurable

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/species-ctor-species-get-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/species-ctor-species-get-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-15.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-15.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A6.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A6.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A8.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A8.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/alternations.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/alternations.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/prototype/dotAll/cross-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/dotAll/cross-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/ignoreCase/cross-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/ignoreCase/cross-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/multiline/name.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/multiline/name.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-difference-property-of-strings-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d--p{Emoji_Keycap_Sequence}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-difference-property-of-strings-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d--p{Emoji_Keycap_Sequence}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-property-Line_Break.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 14: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-property-Line_Break.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 14: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/flags/name.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/flags/name.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Expands_On_NFD-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 18: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Expands_On_NFD-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 18: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-property-Block-with-value-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 15: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-property-Block-with-value-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 15: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_ID_Continue.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 21: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_ID_Continue.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 21: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/get-flags-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a CustomError but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/get-flags-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a CustomError but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/y-fail-lastindex-no-write.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/source/name.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/source/name.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/y-fail-lastindex-no-write.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A2_T6.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A2_T6.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A2_T4.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A2_T4.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-difference-character-class-escape.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-difference-character-class-escape.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/builtin-failure-y-set-lastindex-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/builtin-failure-y-set-lastindex-err.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/flags/coercion-unicode.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/flags/coercion-unicode.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/hasIndices/cross-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/hasIndices/cross-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-03.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Unterminated character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-03.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Unterminated character class"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-02-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-02-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Invalid property name"}

```

- [ ] `built-ins/RegExp/S15.10.2.8_A3_T19.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0xff00,0x10ffff]"}

```

- [ ] `built-ins/RegExp/S15.10.2.8_A3_T19.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0xff00,0x10ffff]"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Full_Composition_Exclusion-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 30: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Full_Composition_Exclusion-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 30: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_No-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 12: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_No-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 12: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-union-character.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d_]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-union-character.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d_]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-difference-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}--q{0|2|4|9️⃣}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-difference-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}--q{0|2|4|9️⃣}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Full_Composition_Exclusion.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 30: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Full_Composition_Exclusion.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 30: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/Basic_Emoji-negative-P.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 15: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/Basic_Emoji-negative-P.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 15: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/unicode/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Composition_Exclusion.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Composition_Exclusion.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-property-Line_Break-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-property-Line_Break-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/unicode/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_N-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_N-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/escape/name.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/named-groups/duplicate-names-match-indices.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/escape/name.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/named-groups/duplicate-names-match-indices.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A7.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A7.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}q{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}q{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/get-exec-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/get-exec-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/escape/escaped-solidus-character-mixed.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `built-ins/RegExp/escape/escaped-solidus-character-mixed.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-union-property-of-strings-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_p{Emoji_Keycap_Sequence}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-union-property-of-strings-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_p{Emoji_Keycap_Sequence}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-union-property-of-strings-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[dp{Emoji_Keycap_Sequence}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-union-property-of-strings-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[dp{Emoji_Keycap_Sequence}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-separator-only.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 5: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-separator-only.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 5: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/match-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/match-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-intersection-character.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_&&_]` should match _ (U+000095)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-intersection-character.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_&&_]` should match _ (U+000095)

```

- [ ] `built-ins/RegExp/property-escapes/generated/Script_Extensions_-_Unknown.js` (strict): **build-error**

```text
command-line-arguments: /workspace/toolchains/go/pkg/tool/linux_amd64/compile: signal: killed

```

- [ ] `built-ins/RegExp/property-escapes/generated/Script_Extensions_-_Unknown.js` (sloppy): **build-timeout**

```text
(no output)
```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-difference-character-class-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_--d]` should match _ (U+000095)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-difference-character-class-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_--d]` should match _ (U+000095)

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Uppercase-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Uppercase-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/isregexp-this-throws.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/isregexp-this-throws.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-intersection-property-of-strings-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}&&p{Emoji_Keycap_Sequence}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-intersection-property-of-strings-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}&&p{Emoji_Keycap_Sequence}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/success-get-index-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/success-get-index-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Grapheme_Extend.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Grapheme_Extend.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/multiline/cross-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/multiline/cross-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-04.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid character in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-04.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid character in character class"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-03.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 7: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-03.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 7: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/character-class-range-no-dash-start.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 15: Invalid character class"}

```

- [ ] `built-ins/RegExp/property-escapes/character-class-range-no-dash-start.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 15: Invalid character class"}

```

- [ ] `built-ins/RegExp/lookBehind/variable-length.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/variable-length.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_Modifier_Sequence-negative-CharacterClass.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 34: Negated character class may contain strings"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_Modifier_Sequence-negative-CharacterClass.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 34: Negated character class may contain strings"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-union-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[dq{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-union-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[dq{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/unicode/this-val-regexp-prototype.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/unicode/this-val-regexp-prototype.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-get-lastindex-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-get-lastindex-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/g-match-no-set-lastindex.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/g-match-no-set-lastindex.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/length.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/length.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-12-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 22: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-12-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 22: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-union-property-of-strings-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}p{Emoji_Keycap_Sequence}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-union-property-of-strings-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}p{Emoji_Keycap_Sequence}]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/prototype/S15.10.6_A1_T2.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/S15.10.6_A1_T2.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/regexp-modifiers/add-ignoreCase-affects-slash-upper-w.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: \W should not match ſ

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/lastindex-no-restore.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/regexp-modifiers/add-ignoreCase-affects-slash-upper-w.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: \W should not match ſ

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/lastindex-no-restore.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-result-coerce-length-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-result-coerce-length-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/property-escapes/non-existent-property-value-Script_Extensions.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Incomplete quantifier"}

```

- [ ] `built-ins/RegExp/property-escapes/non-existent-property-value-Script_Extensions.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Incomplete quantifier"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-difference-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_--q{0|2|4|9️⃣}]` should match _ (U+000095)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-difference-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_--q{0|2|4|9️⃣}]` should match _ (U+000095)

```

- [ ] `built-ins/RegExp/S15.10.2.8_A3_T17.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/S15.10.2.8_A3_T17.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/regexp-modifiers/add-ignoreCase-affects-backreferences.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: a matches a, so \1 should match A (ignores case)

```

- [ ] `built-ins/RegExp/regexp-modifiers/add-ignoreCase-affects-backreferences.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: a matches a, so \1 should match A (ignores case)

```

- [ ] `built-ins/RegExp/prototype/flags/length.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/flags/length.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/lookahead-quantifier-match-groups.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookahead-quantifier-match-groups.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/named-groups/duplicate-names-replaceall.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/named-groups/duplicate-names-replaceall.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/escape/escaped-otherpunctuators.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/escape/escaped-otherpunctuators.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_Tag_Sequence-negative-P.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 26: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_Tag_Sequence-negative-P.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 26: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A11.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A11.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_F-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_F-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/ignoreCase/this-val-invalid-obj.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/ignoreCase/this-val-invalid-obj.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A2_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A2_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-Script-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 10: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-Script-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 10: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-character.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}_]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-character.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}_]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/property-escapes/generated/Script_-_Unknown.js` (sloppy): **build-error**

```text
command-line-arguments: /workspace/toolchains/go/pkg/tool/linux_amd64/compile: signal: killed

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A2_T10.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/name.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/name.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A2_T10.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-constructor-get-species-throws.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-constructor-get-species-throws.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Alphabetic.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 20: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Alphabetic.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 20: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_Y-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_Y-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/escape/not-escaped.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `built-ins/RegExp/escape/not-escaped.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-intersection-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}&&[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-intersection-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}&&[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/global/length.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/global/length.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Math-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 14: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Math-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 14: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Hyphen-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 10: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Hyphen-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 10: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A1_T2.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/prototype/exec/S15.10.6.2_A1_T2.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_ZWJ_Sequence-negative-u.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 26: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_ZWJ_Sequence-negative-u.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 26: Invalid property name"}

```

- [ ] `built-ins/RegExp/S15.10.2.13_A2_T1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/S15.10.2.13_A2_T1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-intersection-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d&&[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-intersection-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d&&[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/lookBehind/word-boundary.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/word-boundary.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-General_Category-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 20: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-General_Category-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 20: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_T-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_T-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 11: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A2_T5.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A2_T5.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/General_Category_-_Letter.js` (sloppy): **build-timeout**

```text
(no output)
```

- [ ] `built-ins/RegExp/property-escapes/generated/General_Category_-_Letter.js` (strict): **build-timeout**

```text
(no output)
```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-intersection-character.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-intersection-character.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/property-escapes/non-existent-property-value-General_Category-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Incomplete quantifier"}

```

- [ ] `built-ins/RegExp/property-escapes/non-existent-property-value-General_Category-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Incomplete quantifier"}

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A10.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-intersection-character.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-intersection-character.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A10.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-no-braces.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 2: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-no-braces.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 2: Invalid property name"}

```

- [ ] `built-ins/RegExp/character-class-escape-non-whitespace.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:22971:29: tsTemp609.read undefined (type *tsCell has no field or method read)
<case>/main.go:22973:19: cannot use tsTemp609.set(loop, tsNumberValue(tsTemp611)) (value of struct type tsValue) as float64 value in assignment

```

- [ ] `built-ins/RegExp/character-class-escape-non-whitespace.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:22974:29: tsTemp610.read undefined (type *tsCell has no field or method read)
<case>/main.go:22976:19: cannot use tsTemp610.set(loop, tsNumberValue(tsTemp612)) (value of struct type tsValue) as float64 value in assignment

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-union-character-class-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-union-character-class-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/property-escapes/generated/Any.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20503:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/property-escapes/generated/Any.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20503:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/escape/length.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/escape/length.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-difference-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}--[0-9]]` should match 9️⃣ (U+000057 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-difference-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}--[0-9]]` should match 9️⃣ (U+000057 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/S15.10.2.13_A1_T1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/S15.10.2.13_A1_T1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/get-global-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-difference-character.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]--_]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-difference-character.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]--_]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/get-global-err.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/non-existent-binary-property.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/non-existent-binary-property.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/multiline/this-val-regexp-prototype.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/multiline/this-val-regexp-prototype.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: RegExp should be an own property

```

- [ ] `built-ins/RegExp/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: RegExp should be an own property

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-circumflex-negation-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-circumflex-negation-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-intersection-character.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-intersection-character.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji-negative-P.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 13: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji-negative-P.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 13: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-separator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-separator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/named-groups/duplicate-names-search.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/named-groups/duplicate-names-search.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/prototype/multiline/15.10.7.4-2.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/multiline/15.10.7.4-2.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/15.10.6.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: s Expected SameValue(«"[object Undefined]"», «"[object Object]"») to be true

```

- [ ] `built-ins/RegExp/prototype/15.10.6.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: s Expected SameValue(«"[object Undefined]"», «"[object Object]"») to be true

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-difference-property-of-strings-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}--p{Emoji_Keycap_Sequence}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-difference-property-of-strings-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}--p{Emoji_Keycap_Sequence}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/S15.10.5.1_A3.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Unsupported Object method

```

- [ ] `built-ins/RegExp/prototype/S15.10.5.1_A3.js` (strict): **runtime-fail**

```text
Async runtime errors: Unsupported Object method

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/this-val-invalid-obj.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/this-val-invalid-obj.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-intersection-character-property-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d&&p{ASCII_Hex_Digit}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-intersection-character-property-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[d&&p{ASCII_Hex_Digit}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-regexp-get-unicode-throws.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String.prototype.matchAll requires a global RegExp

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-regexp-get-unicode-throws.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String.prototype.matchAll requires a global RegExp

```

- [ ] `built-ins/RegExp/prototype/S15.10.6.1_A1_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/S15.10.6.1_A1_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/named-groups/unicode-match.js` (sloppy): **emit-error**

```text
]); pos = int(stack[rxFrame+1]); for rxIndex := range regs { regs[rxIndex] = int(stack[rxFrame+2+rxIndex]) }; stack = stack[:rxFrame]\nswitch choiceTarget {\ncase 1:\ngoto rxNode2\ncase 2:\ngoto rxNode8\ncase 3:\ngoto rxNode7\ncase 4:\ngoto rxNode6\ncase 5:\ngoto rxNode15\ncase 6:\ngoto rxNode9\ndefault: return 0, 0, false\n}\n}\n"}
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: format generated Go matcher: 92:379: expected ';', found 'if' (and 1 more errors)\nconst tsRegexMatcher824CaptureSlots = 10\nfunc tsRegexMatcher824(input []uint16, start int) (matchStart, matchEnd int, captures []int, ok bool) {\noffsets := make([]int32, tsRegexMatcher824CaptureSlots)\nmatchStart, matchEnd, ok = tsRegexMatcher824Into(input, start, offsets)\nif !ok { return 0, 0, nil, false }; captures = make([]int, len(offsets)); for index, offset := range offsets { captures[index] = int(offset) }; return matchStart, matchEnd, captures, true\n}\nfunc tsRegexMatcher824Into(input []uint16, start int, captures []int32) (matchStart, matchEnd int, ok bool) {\nconst captureSlots = 10\nif start \u003c 0 { start = 0 }; if start \u003e len(input) { return 0, 0, false }; if len(captures) \u003c captureSlots { return 0, 0, false }\nregs := make([]int, 12)\nfor i := range regs { regs[i] = -1 }\npos := start\nattemptStart := start\nconst rxStackFrameSlots = 14\nconst rxMaxStackSlots = 64 * 1024 * 1024 / 4\nstack := make([]int32, 0, 128)\nvar choiceTarget, rxFrame int\ngoto rxNode0\nrxNode0:\nif len(stack)+rxStackFrameSlots \u003e rxMaxStackSlots { return -1, -1, false }; if len(stack)
```

- [ ] `built-ins/RegExp/named-groups/unicode-match.js` (strict): **emit-error**

```text
]); pos = int(stack[rxFrame+1]); for rxIndex := range regs { regs[rxIndex] = int(stack[rxFrame+2+rxIndex]) }; stack = stack[:rxFrame]\nswitch choiceTarget {\ncase 1:\ngoto rxNode2\ncase 2:\ngoto rxNode8\ncase 3:\ngoto rxNode7\ncase 4:\ngoto rxNode6\ncase 5:\ngoto rxNode15\ncase 6:\ngoto rxNode9\ndefault: return 0, 0, false\n}\n}\n"}
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: format generated Go matcher: 92:379: expected ';', found 'if' (and 1 more errors)\nconst tsRegexMatcher825CaptureSlots = 10\nfunc tsRegexMatcher825(input []uint16, start int) (matchStart, matchEnd int, captures []int, ok bool) {\noffsets := make([]int32, tsRegexMatcher825CaptureSlots)\nmatchStart, matchEnd, ok = tsRegexMatcher825Into(input, start, offsets)\nif !ok { return 0, 0, nil, false }; captures = make([]int, len(offsets)); for index, offset := range offsets { captures[index] = int(offset) }; return matchStart, matchEnd, captures, true\n}\nfunc tsRegexMatcher825Into(input []uint16, start int, captures []int32) (matchStart, matchEnd int, ok bool) {\nconst captureSlots = 10\nif start \u003c 0 { start = 0 }; if start \u003e len(input) { return 0, 0, false }; if len(captures) \u003c captureSlots { return 0, 0, false }\nregs := make([]int, 12)\nfor i := range regs { regs[i] = -1 }\npos := start\nattemptStart := start\nconst rxStackFrameSlots = 14\nconst rxMaxStackSlots = 64 * 1024 * 1024 / 4\nstack := make([]int32, 0, 128)\nvar choiceTarget, rxFrame int\ngoto rxNode0\nrxNode0:\nif len(stack)+rxStackFrameSlots \u003e rxMaxStackSlots { return -1, -1, false }; if len(stack)
```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Prepended_Concatenation_Mark-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 32: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Prepended_Concatenation_Mark-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 32: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-union-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}[0-9]]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-union-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}[0-9]]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/this-tostring-flags-throws.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier ERROR"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/this-tostring-flags-throws.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier ERROR"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/builtin-infer-unicode.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/regexpcreate-this-throws.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/builtin-infer-unicode.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/regexpcreate-this-throws.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/species-ctor-ctor-get-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/species-ctor-ctor-get-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-intersection-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}&&q{0|2|4|9️⃣}]` should match 9️⃣ (U+000057 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-intersection-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}&&q{0|2|4|9️⃣}]` should match 9️⃣ (U+000057 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/last-index-exceeds-str-size.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/last-index-exceeds-str-size.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-intersection-character-class-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]&&d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-intersection-character-class-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]&&d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/S15.10.2.13_A1_T2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/S15.10.2.13_A1_T2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji-negative-CharacterClass.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 16: Negated character class may contain strings"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji-negative-CharacterClass.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 16: Negated character class may contain strings"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/rgi-emoji-14.0.js` (strict): **build-timeout**

```text
(no output)
```

- [ ] `built-ins/RegExp/unicodeSets/generated/rgi-emoji-14.0.js` (sloppy): **build-timeout**

```text
(no output)
```

- [ ] `built-ins/RegExp/prototype/unicode/name.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/unicode/name.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-intersection-property-of-strings-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}&&p{Emoji_Keycap_Sequence}]` should match 9️⃣ (U+000057 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-intersection-property-of-strings-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}&&p{Emoji_Keycap_Sequence}]` should match 9️⃣ (U+000057 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-23.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-23.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-coerce-lastindex-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-coerce-lastindex-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-union-character.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}_]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-union-character.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}_]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/prototype/multiline/S15.10.7.4_A9.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/multiline/S15.10.7.4_A9.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/S15.10.4.1_A6_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: __re.toString() must return "[object "+"RegExp"+"]" Expected SameValue(«"[object Object]"», «"[object RegExp]"») to be true

```

- [ ] `built-ins/RegExp/S15.10.4.1_A6_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: __re.toString() must return "[object "+"RegExp"+"]" Expected SameValue(«"[object Object]"», «"[object RegExp]"») to be true

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-coerce-lastindex.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-coerce-lastindex.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-04.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 12: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-04.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 12: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-11.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-11.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/S15.10.2.9_A1_T5.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/S15.10.2.9_A1_T5.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/species-ctor-y.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/species-ctor-y.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/escape/not-a-constructor.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: isConstructor invoked with a non-function value

```

- [ ] `built-ins/RegExp/escape/not-a-constructor.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: isConstructor invoked with a non-function value

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_Invalid-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 17: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_Invalid-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 17: Invalid property name"}

```

- [ ] `built-ins/RegExp/escape/initial-char-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `built-ins/RegExp/escape/initial-char-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `built-ins/RegExp/prototype/toString/name.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/toString/name.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-union-character-property-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]p{ASCII_Hex_Digit}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-union-character-property-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9]p{ASCII_Hex_Digit}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/species-ctor-ctor-non-obj.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/species-ctor-ctor-non-obj.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-06-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-06-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-FC_NFKC_Closure-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-FC_NFKC_Closure-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/escape/escaped-surrogates.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `built-ins/RegExp/escape/escaped-surrogates.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `built-ins/RegExp/prototype/flags/coercion-ignoreCase.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/flags/coercion-ignoreCase.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-05.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid character in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-05.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid character in character class"}

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-Script.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 10: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-Script.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 10: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-union-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9][0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-union-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[[0-9][0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/property-escapes/character-class-range-no-dash-end.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 15: Invalid character class"}

```

- [ ] `built-ins/RegExp/property-escapes/character-class-range-no-dash-end.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 15: Invalid character class"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-constructor-is-undefined.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:24179:33: tsTemp849.read undefined (type *tsCell has no field or method read)
<case>/main.go:24181:23: cannot use tsTemp849.set(loop, tsNumberValue(tsTemp851)) (value of struct type tsValue) as float64 value in assignment

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-constructor-is-undefined.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:24182:33: tsTemp850.read undefined (type *tsCell has no field or method read)
<case>/main.go:24184:23: cannot use tsTemp850.set(loop, tsNumberValue(tsTemp852)) (value of struct type tsValue) as float64 value in assignment

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A7.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A7.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/hasIndices/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/hasIndices/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-difference-character-class-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}--d]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/property-of-strings-escape-difference-character-class-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{Emoji_Keycap_Sequence}--d]` should match #️⃣ (U+000035 U+065039 U+008419)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-difference-character-class-escape.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-difference-character-class-escape.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/coerce-flags.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/coerce-flags.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/named-groups/lookbehind.js` (sloppy): **emit-error**

```text
panic: unknown V8 graph node in FirstCharacterSetBuilder [recovered, repanicked]

goroutine 162 [running]:
sync.(*WaitGroup).Go.func1.1()
	/workspace/toolchains/go/src/sync/waitgroup.go:251 +0x45
panic({0x1b48ca8?, 0x1289850?})
	/workspace/toolchains/go/src/runtime/panic.go:859 +0x125
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8CompilerFirstCharacterSetBuilder).addNode(0x3?, {0x1be96d0?, 0x1307abb6db90?}, 0x1307abbd7888?)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_compiler_first_char.go:53 +0x1b9
github.com/microsoft/TypeScript/tsc/internal/regexpengine.v8CompilerComputeFirstCharacterSet({0x1be96d0?, 0x1307abb6db90?}, 0x1?)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_compiler_first_char.go:18 +0x5b
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8GoGraphEmitter).emitChoice(0x1307abb9a3a0, 0x1be9bb0?, 0x7)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:754 +0x10f
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8GoGraphEmitter).emitNode(0x1307abb9a3a0, 0x7, {0x1be9940, 0x1307abc6c400})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:524 +0x7fc
github.com/microsoft/TypeScript/tsc/internal/regexpengine.v8EmitGoMatcher(0x1307abc0f100, {0x1307abb55e00, 0x11})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:243 +0x2b9
github.com/microsoft/TypeScript/tsc/internal/regexpengine.CompileUTF16ToGo({0x1307abbadd70, 0x11, 0x11}, {0x1307aa1ffa98, 0x1}, {0x1307abb55e00, 0x11})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:43 +0x3e6
g
```

- [ ] `built-ins/RegExp/named-groups/lookbehind.js` (strict): **emit-error**

```text
panic: unknown V8 graph node in FirstCharacterSetBuilder [recovered, repanicked]

goroutine 46 [running]:
sync.(*WaitGroup).Go.func1.1()
	/workspace/toolchains/go/src/sync/waitgroup.go:251 +0x45
panic({0x1b48ca8?, 0x1289850?})
	/workspace/toolchains/go/src/runtime/panic.go:859 +0x125
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8CompilerFirstCharacterSetBuilder).addNode(0x3?, {0x1be96d0?, 0x2782eea1ac40?}, 0x2782ee9df508?)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_compiler_first_char.go:53 +0x1b9
github.com/microsoft/TypeScript/tsc/internal/regexpengine.v8CompilerComputeFirstCharacterSet({0x1be96d0?, 0x2782eea1ac40?}, 0x1?)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_compiler_first_char.go:18 +0x5b
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8GoGraphEmitter).emitChoice(0x2782ee99a3a0, 0x1be9bb0?, 0x7)
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:754 +0x10f
github.com/microsoft/TypeScript/tsc/internal/regexpengine.(*v8GoGraphEmitter).emitNode(0x2782ee99a3a0, 0x7, {0x1be9940, 0x2782ee9dd200})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:524 +0x7fc
github.com/microsoft/TypeScript/tsc/internal/regexpengine.v8EmitGoMatcher(0x2782edf50f00, {0x2782ee945e00, 0x11})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:243 +0x2b9
github.com/microsoft/TypeScript/tsc/internal/regexpengine.CompileUTF16ToGo({0x2782ee9addd0, 0x11, 0x11}, {0x2782ed291aa6, 0x1}, {0x2782ee945e00, 0x11})
	/workspace/TypeScript/tsc/internal/regexpengine/regexp_go_codegen.go:43 +0x3e6
gi
```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-union-character-property-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_p{ASCII_Hex_Digit}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/global/this-val-non-obj.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/unicode/cross-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/prototype/unicode/cross-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-union-character-property-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_p{ASCII_Hex_Digit}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/global/this-val-non-obj.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-intersection-character-class-escape.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_ID_Start-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 18: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_ID_Start-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 18: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-intersection-character-class-escape.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-character-property-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}p{ASCII_Hex_Digit}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-character-property-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}p{ASCII_Hex_Digit}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-intersection-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}&&[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-intersection-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}&&[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/flags/coercion-hasIndices.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/flags/coercion-hasIndices.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-empty-match-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-empty-match-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A1_T11.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A1_T11.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_ZWJ_Sequence-negative-P.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 26: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_ZWJ_Sequence-negative-P.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 26: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-14.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 13: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-14.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 13: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/global/15.10.7.2-2.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/S15.10.2.7_A4_T19.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/global/15.10.7.2-2.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/S15.10.2.7_A4_T19.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/limit-0-bail.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/limit-0-bail.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-13.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-13.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/flags/this-val-non-obj.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/flags/this-val-non-obj.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/lookBehind/greedy-loop.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/greedy-loop.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/builtin-success-g-set-lastindex-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/builtin-success-g-set-lastindex-err.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/lookBehind/start-of-line.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/start-of-line.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/negative.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/negative.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_Modifier_Sequence-negative-P.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 31: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/RGI_Emoji_Modifier_Sequence-negative-P.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 31: Invalid property name"}

```

- [ ] `built-ins/RegExp/S15.10.2.13_A2_T8.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/S15.10.2.13_A2_T8.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-property-Line_Break-with-value.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-property-Line_Break-with-value.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 25: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/rgi-emoji-13.1.js` (sloppy): **build-timeout**

```text
(no output)
```

- [ ] `built-ins/RegExp/unicodeSets/generated/rgi-emoji-13.1.js` (strict): **build-timeout**

```text
(no output)
```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_ID_Start.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 18: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_ID_Start.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 18: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/sticky/this-val-non-obj.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-match-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/sticky/this-val-non-obj.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-match-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_Invalid.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 17: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/binary-property-with-value-ASCII_-_Invalid.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 17: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A1_T2.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A1_T2.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-regexp-get-global-throws.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String.prototype.matchAll requires a global RegExp

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-regexp-get-global-throws.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String.prototype.matchAll requires a global RegExp

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/Emoji_Keycap_Sequence-negative-CharacterClass.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 28: Negated character class may contain strings"}

```

- [ ] `built-ins/RegExp/property-escapes/generated/strings/Emoji_Keycap_Sequence-negative-CharacterClass.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 28: Negated character class may contain strings"}

```

- [ ] `built-ins/RegExp/named-groups/non-unicode-references.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/named-groups/non-unicode-references.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-14-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 13: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-14-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 13: Invalid property name"}

```

- [ ] `built-ins/RegExp/S15.10.2.13_A1_T17.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/S15.10.2.13_A1_T17.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: Go regex codegen encountered a non-UTF-16 class range [0x0,0x10ffff]"}

```

- [ ] `built-ins/RegExp/escape/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: escape should be an own property

```

- [ ] `built-ins/RegExp/escape/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: escape should be an own property

```

- [ ] `built-ins/RegExp/named-groups/groups-object-undefined.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-10.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-10.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/named-groups/groups-object-undefined.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-separator-only-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 5: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-separator-only-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 5: Invalid property name"}

```

- [ ] `built-ins/RegExp/match-indices/indices-array-unicode-match.js` (sloppy): **emit-error**

```text
de27\nrxNode25:\nregs[8] = len(stack); regs[9] = pos\ngoto rxNode28\nrxNode26:\nif pos \u003e= len(input) || ((uint32(input[pos]) \u003e= 56320 \u0026\u0026 uint32(input[pos]) \u003c= 57343)) { if len(stack)+rxStackFrameSlots \u003e rxMaxStackSlots { return -1, -1, false }; if len(stack)+rxStackFrameSlots \u003e cap(stack) { rxStackCap10 := cap(stack)*2; if rxStackCap10 \u003c len(stack)+rxStackFrameSlots { rxStackCap10 = len(stack)+rxStackFrameSlots }; if rxStackCap10 \u003e rxMaxStackSlots { rxStackCap10 = rxMaxStackSlots }; rxGrownStack10 := make([]int32, len(stack), rxStackCap10); copy(rxGrownStack10, stack); stack = rxGrownStack10 }; stack = append(stack, int32(10), int32(pos), int32(regs[0]) , int32(regs[1]) , int32(regs[2]) , int32(regs[3]) , int32(regs[4]) , int32(regs[5]) , int32(regs[6]) , int32(regs[7]) , int32(regs[8]) , int32(regs[9]) ) }\ngoto rxNode29\nrxNode27:\nregs[6] = pos\ngoto rxNode31\nrxNode28:\nif pos \u003e= len(input) || ((uint32(input[pos]) \u003e= 0 \u0026\u0026 uint32(input[pos]) \u003c= 9) || (uint32(input[pos]) \u003e= 11 \u0026\u0026 uint32(input[pos]) \u003c= 12) || (uint32(input[pos]) \u003e= 14 \u0026\u0026 uint32(input[pos]) \u003c= 8231) || (uint32(input[pos]) \u003e= 8234 \u0026\u0026 uint32(input[pos]) \u003c= 65535)) { if len(stack)+rxStackFrameSlots \u003e rxMaxStackSlots { return -1, -1, false }; if len(stack)+rxStackFrameSlots \u003e cap(stack) { rxStackCap11 := cap(stack)*2; if rxStackCap11 \u003c len(stack)+rxStackFrameSlots { rxStackCap11 = len(stack)+rxStackFrameSlots }; if rxStackCap11 \u003e rxMaxStackSlots { rxStackCap11 = r
```

- [ ] `built-ins/RegExp/match-indices/indices-array-unicode-match.js` (strict): **emit-error**

```text
de27\nrxNode25:\nregs[8] = len(stack); regs[9] = pos\ngoto rxNode28\nrxNode26:\nif pos \u003e= len(input) || ((uint32(input[pos]) \u003e= 56320 \u0026\u0026 uint32(input[pos]) \u003c= 57343)) { if len(stack)+rxStackFrameSlots \u003e rxMaxStackSlots { return -1, -1, false }; if len(stack)+rxStackFrameSlots \u003e cap(stack) { rxStackCap10 := cap(stack)*2; if rxStackCap10 \u003c len(stack)+rxStackFrameSlots { rxStackCap10 = len(stack)+rxStackFrameSlots }; if rxStackCap10 \u003e rxMaxStackSlots { rxStackCap10 = rxMaxStackSlots }; rxGrownStack10 := make([]int32, len(stack), rxStackCap10); copy(rxGrownStack10, stack); stack = rxGrownStack10 }; stack = append(stack, int32(10), int32(pos), int32(regs[0]) , int32(regs[1]) , int32(regs[2]) , int32(regs[3]) , int32(regs[4]) , int32(regs[5]) , int32(regs[6]) , int32(regs[7]) , int32(regs[8]) , int32(regs[9]) ) }\ngoto rxNode29\nrxNode27:\nregs[6] = pos\ngoto rxNode31\nrxNode28:\nif pos \u003e= len(input) || ((uint32(input[pos]) \u003e= 0 \u0026\u0026 uint32(input[pos]) \u003c= 9) || (uint32(input[pos]) \u003e= 11 \u0026\u0026 uint32(input[pos]) \u003c= 12) || (uint32(input[pos]) \u003e= 14 \u0026\u0026 uint32(input[pos]) \u003c= 8231) || (uint32(input[pos]) \u003e= 8234 \u0026\u0026 uint32(input[pos]) \u003c= 65535)) { if len(stack)+rxStackFrameSlots \u003e rxMaxStackSlots { return -1, -1, false }; if len(stack)+rxStackFrameSlots \u003e cap(stack) { rxStackCap11 := cap(stack)*2; if rxStackCap11 \u003c len(stack)+rxStackFrameSlots { rxStackCap11 = len(stack)+rxStackFrameSlots }; if rxStackCap11 \u003e rxMaxStackSlots { rxStackCap11 = r
```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Expands_On_NFKC.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Expands_On_NFKC.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/named-groups/duplicate-names-split.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/named-groups/duplicate-names-split.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-24.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-24.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 1: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-result-get-length-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/str-result-get-length-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-union-character-class-escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-union-character-class-escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}d]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/cstm-exec-return-index.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.search/cstm-exec-return-index.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Default_Ignorable_Code_Point.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 38: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Other_Default_Ignorable_Code_Point.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 38: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-difference-character.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-difference-character.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-constructor-species-throws.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/Symbol.matchAll/species-constructor-species-throws.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: Expected a Test262Error but got a TypeError

```

- [ ] `built-ins/RegExp/prototype/flags/this-val-regexp-prototype.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/flags/this-val-regexp-prototype.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-General_Category-equals-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 21: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/non-binary-property-without-value-General_Category-equals-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 21: Invalid property name"}

```

- [ ] `built-ins/RegExp/prototype/dotAll/this-val-regexp-prototype.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/dotAll/this-val-regexp-prototype.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/length.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.replace/length.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/back-references-to-captures.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/lookBehind/back-references-to-captures.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `built-ins/RegExp/prototype/global/S15.10.7.2_A9.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/global/S15.10.7.2_A9.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/property-escapes/non-existent-property-value-Script-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Incomplete quantifier"}

```

- [ ] `built-ins/RegExp/property-escapes/non-existent-property-value-Script-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 3: Incomplete quantifier"}

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/species-ctor.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/prototype/Symbol.split/species-ctor.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Cannot read properties of null or undefined

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-difference-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_--[0-9]]` should match _ (U+000095)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-difference-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[_--[0-9]]` should match _ (U+000095)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-intersection-property-of-strings-escape.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-class-escape-intersection-property-of-strings-escape.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A2_T4.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/property-escapes/character-class-range-end.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 10: Invalid character class"}

```

- [ ] `built-ins/RegExp/property-escapes/character-class-range-end.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 10: Invalid character class"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-empty.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 5: Invalid property name in character class"}

```

- [ ] `built-ins/RegExp/property-escapes/grammar-extension-empty.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 5: Invalid property name in character class"}

```

- [ ] `built-ins/RegExp/prototype/test/S15.10.6.3_A2_T4.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not a constructor

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-intersection-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}&&q{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-intersection-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}&&q{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-union-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}q{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/string-literal-union-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[q{0|2|4|9️⃣}q{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/escape/escaped-control-characters.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `built-ins/RegExp/escape/escaped-control-characters.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-union-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}[0-9]]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-intersection-string-literal.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}&&q{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-intersection-string-literal.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}&&q{0|2|4|9️⃣}]` should match 0 (U+000048)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-intersection-character-property-escape.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-intersection-character-property-escape.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:20058:2: declared and not used: pos
<case>/main.go:20059:2: declared and not used: attemptStart

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Expands_On_NFKD.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/unsupported-binary-property-Expands_On_NFKD.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 19: Invalid property name"}

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-difference-character-class.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}--[0-9]]` should match A (U+000065)

```

- [ ] `built-ins/RegExp/unicodeSets/generated/character-property-escape-difference-character-class.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: `[p{ASCII_Hex_Digit}--[0-9]]` should match A (U+000065)

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-27.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 2: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/prototype/unicodeSets/breaking-change-from-u-to-v-27.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 2: Invalid set operation in character class"}

```

- [ ] `built-ins/RegExp/escape/non-string-inputs.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Test262Error: RegExp.escape is a function Expected SameValue(«"undefined"», «"function"») to be true

```

- [ ] `built-ins/RegExp/escape/non-string-inputs.js` (strict): **runtime-fail**

```text
Async runtime errors: Test262Error: RegExp.escape is a function Expected SameValue(«"undefined"», «"function"») to be true

```

- [ ] `built-ins/RegExp/named-groups/unicode-references.js` (sloppy): **emit-error**

```text
pos\ngoto rxNode3\nrxNode2:\nif pos \u003e= len(input) { goto rxBacktrack }; if input[pos] \u003e= 0xd800 \u0026\u0026 input[pos] \u003c= 0xdbff \u0026\u0026 pos+1 \u003c len(input) \u0026\u0026 input[pos+1] \u003e= 0xdc00 \u0026\u0026 input[pos+1] \u003c= 0xdfff { pos += 2 } else { pos++ }\nfor pos \u003c len(input) \u0026\u0026 !(uint32(input[pos]) == 97) { if input[pos] \u003e= 0xd800 \u0026\u0026 input[pos] \u003c= 0xdbff \u0026\u0026 pos+1 \u003c len(input) \u0026\u0026 input[pos+1] \u003e= 0xdc00 \u0026\u0026 input[pos+1] \u003c= 0xdfff { pos += 2 } else { pos++ } }; if pos \u003e= len(input) { goto rxBacktrack }\ngoto rxNode0\nrxNode3:\nregs[2] = pos\ngoto rxNode4\nrxNode4:\nif pos \u003c 0 || pos+1 \u003e len(input) { goto rxBacktrack }\nif input[pos+0] != 97 { goto rxBacktrack }\npos += 1\ngoto rxNode5\nrxNode5:\nregs[3] = pos\ngoto rxNode6\nrxNode6:\nregs[4] = pos\ngoto rxNode7\nrxNode7:\nif pos \u003c 0 || pos+1 \u003e len(input) { goto rxBacktrack }\nif input[pos+0] != 98 { goto rxBacktrack }\npos += 1\ngoto rxNode8\nrxNode8:\nregs[5] = pos\ngoto rxNode9\nrxNode9:\n{ rxStart := regs[2]; rxEnd := regs[3]; if rxStart \u003e= 0 \u0026\u0026 rxEnd \u003e= 0 { rxLow, rxHigh := rxStart, rxEnd; if rxLow \u003e rxHigh { rxLow, rxHigh = rxHigh, rxLow }; rxLength := rxHigh-rxLow; if rxLength \u003e len(input)-pos { goto rxBacktrack }; for rxOffset := 0; rxOffset \u003c rxLength; rxOffset++ { if input[rxLow+rxOffset] != input[pos+rxOffset] { goto rxBacktrack } }; pos += rxLength } if pos \u003e 0 \u0026\u0026 pos \u003c len(input) \u0026\u0026 input[pos-1] \u003e= 0xd800 \
```

- [ ] `built-ins/RegExp/named-groups/unicode-references.js` (strict): **emit-error**

```text
pos\ngoto rxNode3\nrxNode2:\nif pos \u003e= len(input) { goto rxBacktrack }; if input[pos] \u003e= 0xd800 \u0026\u0026 input[pos] \u003c= 0xdbff \u0026\u0026 pos+1 \u003c len(input) \u0026\u0026 input[pos+1] \u003e= 0xdc00 \u0026\u0026 input[pos+1] \u003c= 0xdfff { pos += 2 } else { pos++ }\nfor pos \u003c len(input) \u0026\u0026 !(uint32(input[pos]) == 97) { if input[pos] \u003e= 0xd800 \u0026\u0026 input[pos] \u003c= 0xdbff \u0026\u0026 pos+1 \u003c len(input) \u0026\u0026 input[pos+1] \u003e= 0xdc00 \u0026\u0026 input[pos+1] \u003c= 0xdfff { pos += 2 } else { pos++ } }; if pos \u003e= len(input) { goto rxBacktrack }\ngoto rxNode0\nrxNode3:\nregs[2] = pos\ngoto rxNode4\nrxNode4:\nif pos \u003c 0 || pos+1 \u003e len(input) { goto rxBacktrack }\nif input[pos+0] != 97 { goto rxBacktrack }\npos += 1\ngoto rxNode5\nrxNode5:\nregs[3] = pos\ngoto rxNode6\nrxNode6:\nregs[4] = pos\ngoto rxNode7\nrxNode7:\nif pos \u003c 0 || pos+1 \u003e len(input) { goto rxBacktrack }\nif input[pos+0] != 98 { goto rxBacktrack }\npos += 1\ngoto rxNode8\nrxNode8:\nregs[5] = pos\ngoto rxNode9\nrxNode9:\n{ rxStart := regs[2]; rxEnd := regs[3]; if rxStart \u003e= 0 \u0026\u0026 rxEnd \u003e= 0 { rxLow, rxHigh := rxStart, rxEnd; if rxLow \u003e rxHigh { rxLow, rxHigh = rxHigh, rxLow }; rxLength := rxHigh-rxLow; if rxLength \u003e len(input)-pos { goto rxBacktrack }; for rxOffset := 0; rxOffset \u003c rxLength; rxOffset++ { if input[rxLow+rxOffset] != input[pos+rxOffset] { goto rxBacktrack } }; pos += rxLength } if pos \u003e 0 \u0026\u0026 pos \u003c len(input) \u0026\u0026 input[pos-1] \u003e= 0xd800 \
```

- [ ] `built-ins/RegExp/prototype/global/this-val-regexp-prototype.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/global/this-val-regexp-prototype.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/y-fail-lastindex-no-write.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/Symbol.match/y-fail-lastindex-no-write.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/source/prop-desc.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/prototype/source/prop-desc.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Property operation requires an object

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-05-negated.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 9: Invalid property name"}

```

- [ ] `built-ins/RegExp/property-escapes/loose-matching-05-negated.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: regular expression cannot be emitted as Go code: regular expression at 9: Invalid property name"}

```

### staging/built-ins

- [ ] `staging/built-ins/RegExp/named-groups/duplicate-named-groups-replace.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `staging/built-ins/RegExp/named-groups/duplicate-named-groups-replace.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `staging/built-ins/RegExp/named-groups/duplicate-named-groups.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `staging/built-ins/RegExp/named-groups/duplicate-named-groups.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `staging/built-ins/RegExp/named-groups/duplicate-named-groups-search.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `staging/built-ins/RegExp/named-groups/duplicate-named-groups-search.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: String method called on null or undefined

```

- [ ] `staging/built-ins/RegExp/escape/surrogate-pair.js` (strict): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```

- [ ] `staging/built-ins/RegExp/escape/surrogate-pair.js` (sloppy): **runtime-fail**

```text
Async runtime errors: TypeError: Value is not callable

```
