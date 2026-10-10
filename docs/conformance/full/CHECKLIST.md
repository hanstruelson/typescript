# JavaScript conformance checklist

Test262 revision: `2e0a56762801e275a9fdf96dc49d90ba0cddcf63`.
Inventory: 53,616 tests; 1,899 explicitly excluded for dynamic code.
Execution: 1,200 variants from 51,717 selected tests; selection: all.

**Full-suite run: results below are checkpoints until every selected variant completes.** Untested cases are not passes. The default coercion policy is unchanged; failures may expose differences from ECMAScript.

Module and host-adapter limitations are tracked separately from compiler/runtime failures. Node reference failures are recorded separately and do not prevent native execution. Negative tests pass only on syntactic rejection, never on an arbitrary build failure.

## Results

- build-error: 4
- emit-error: 919
- pass: 51
- runner-host-unimplemented: 24
- runner-module-unimplemented: 1
- runtime-fail: 201

## Category checklist

| Category | Suite tests | Dynamic exclusions | Executed variants | Pass | Native failures | Adapter/reference limits |
|---|---:|---:|---:|---:|---:|---:|
| annexB/built-ins | 241 | 8 | 202 | 0 | 190 | 12 |
| annexB/language | 845 | 309 | 1 | 0 | 1 | 0 |
| built-ins/AbstractModuleSource | 8 | 0 | 1 | 0 | 0 | 1 |
| built-ins/AggregateError | 25 | 2 | 2 | 0 | 2 | 0 |
| built-ins/Array | 3083 | 5 | 2 | 0 | 2 | 0 |
| built-ins/ArrayBuffer | 221 | 1 | 2 | 0 | 2 | 0 |
| built-ins/ArrayIteratorPrototype | 27 | 0 | 2 | 0 | 2 | 0 |
| built-ins/AsyncDisposableStack | 104 | 1 | 2 | 0 | 2 | 0 |
| built-ins/AsyncFromSyncIteratorPrototype | 38 | 0 | 2 | 0 | 2 | 0 |
| built-ins/AsyncFunction | 18 | 4 | 2 | 0 | 2 | 0 |
| built-ins/AsyncGeneratorFunction | 23 | 14 | 2 | 0 | 2 | 0 |
| built-ins/AsyncGeneratorPrototype | 48 | 0 | 2 | 0 | 2 | 0 |
| built-ins/AsyncIteratorPrototype | 13 | 0 | 2 | 0 | 2 | 0 |
| built-ins/Atomics | 389 | 0 | 2 | 0 | 2 | 0 |
| built-ins/BigInt | 77 | 0 | 2 | 0 | 2 | 0 |
| built-ins/Boolean | 51 | 3 | 2 | 0 | 2 | 0 |
| built-ins/DataView | 561 | 2 | 2 | 0 | 2 | 0 |
| built-ins/Date | 594 | 3 | 2 | 0 | 2 | 0 |
| built-ins/DisposableStack | 93 | 1 | 2 | 0 | 2 | 0 |
| built-ins/Error | 93 | 2 | 2 | 0 | 2 | 0 |
| built-ins/FinalizationRegistry | 47 | 1 | 2 | 0 | 2 | 0 |
| built-ins/Function | 509 | 151 | 2 | 0 | 2 | 0 |
| built-ins/GeneratorFunction | 23 | 14 | 2 | 0 | 2 | 0 |
| built-ins/GeneratorPrototype | 61 | 0 | 2 | 0 | 2 | 0 |
| built-ins/Infinity | 6 | 0 | 2 | 0 | 2 | 0 |
| built-ins/Iterator | 654 | 1 | 2 | 0 | 2 | 0 |
| built-ins/JSON | 166 | 0 | 2 | 0 | 2 | 0 |
| built-ins/Map | 204 | 2 | 2 | 0 | 2 | 0 |
| built-ins/MapIteratorPrototype | 11 | 0 | 2 | 0 | 2 | 0 |
| built-ins/Math | 327 | 0 | 2 | 0 | 2 | 0 |
| built-ins/NaN | 6 | 0 | 2 | 0 | 2 | 0 |
| built-ins/NativeErrors | 94 | 6 | 2 | 0 | 2 | 0 |
| built-ins/Number | 340 | 1 | 2 | 0 | 2 | 0 |
| built-ins/Object | 3411 | 7 | 2 | 0 | 2 | 0 |
| built-ins/Promise | 732 | 1 | 2 | 0 | 2 | 0 |
| built-ins/Proxy | 311 | 6 | 2 | 0 | 0 | 2 |
| built-ins/Reflect | 153 | 1 | 2 | 0 | 2 | 0 |
| built-ins/RegExp | 1879 | 9 | 2 | 0 | 2 | 0 |
| built-ins/RegExpStringIteratorPrototype | 17 | 0 | 2 | 0 | 2 | 0 |
| built-ins/Set | 383 | 1 | 2 | 0 | 2 | 0 |
| built-ins/SetIteratorPrototype | 11 | 0 | 2 | 0 | 2 | 0 |
| built-ins/ShadowRealm | 64 | 1 | 2 | 0 | 2 | 0 |
| built-ins/SharedArrayBuffer | 104 | 1 | 2 | 0 | 2 | 0 |
| built-ins/String | 1223 | 17 | 2 | 0 | 2 | 0 |
| built-ins/StringIteratorPrototype | 7 | 0 | 2 | 0 | 2 | 0 |
| built-ins/SuppressedError | 22 | 2 | 2 | 0 | 2 | 0 |
| built-ins/Symbol | 98 | 0 | 2 | 0 | 2 | 0 |
| built-ins/Temporal | 4605 | 0 | 2 | 0 | 2 | 0 |
| built-ins/ThrowTypeError | 14 | 1 | 2 | 0 | 2 | 0 |
| built-ins/TypedArray | 1453 | 0 | 2 | 0 | 2 | 0 |
| built-ins/TypedArrayConstructors | 738 | 12 | 2 | 0 | 2 | 0 |
| built-ins/Uint8Array | 70 | 0 | 2 | 0 | 2 | 0 |
| built-ins/WeakMap | 141 | 2 | 2 | 0 | 2 | 0 |
| built-ins/WeakRef | 29 | 1 | 2 | 0 | 2 | 0 |
| built-ins/WeakSet | 85 | 1 | 2 | 0 | 2 | 0 |
| built-ins/decodeURI | 55 | 0 | 2 | 0 | 2 | 0 |
| built-ins/decodeURIComponent | 56 | 0 | 2 | 0 | 2 | 0 |
| built-ins/encodeURI | 31 | 0 | 2 | 0 | 2 | 0 |
| built-ins/encodeURIComponent | 31 | 0 | 2 | 0 | 2 | 0 |
| built-ins/eval | 10 | 10 | 0 | 0 | 0 | 0 |
| built-ins/global | 29 | 8 | 2 | 0 | 2 | 0 |
| built-ins/isFinite | 15 | 0 | 2 | 0 | 2 | 0 |
| built-ins/isNaN | 15 | 0 | 2 | 0 | 2 | 0 |
| built-ins/parseFloat | 54 | 0 | 2 | 0 | 2 | 0 |
| built-ins/parseInt | 55 | 0 | 2 | 0 | 2 | 0 |
| built-ins/undefined | 8 | 1 | 2 | 0 | 2 | 0 |
| harness/assert-false.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-notsamevalue-nan.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-notsamevalue-notsame.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-notsamevalue-objects.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-notsamevalue-tostring.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-notsamevalue-zeros.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-obj.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-samevalue-nan.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-samevalue-objects.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-samevalue-same.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-samevalue-tostring.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-samevalue-zeros.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-throws-custom-typeerror.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-throws-custom.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-throws-incorrect-ctor.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-throws-native.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-throws-no-arg.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-throws-no-error.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-throws-null-fn.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-throws-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-throws-primitive.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-throws-same-realm.js | 1 | 0 | 2 | 0 | 0 | 2 |
| harness/assert-throws-single-arg.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-tostring.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assert-true.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/assertRelativeDateMs.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-asyncTest-func-throws-sync.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-asyncTest-rejects-non-callable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-asyncTest-return-not-thenable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-asyncTest-returns-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-asyncTest-then-rejects.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-asyncTest-then-resolves.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-asyncTest-without-async-flag.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-custom-typeerror.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-custom.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-func-never-settles.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-func-throws-sync.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-incorrect-ctor.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-invalid-func.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-native.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-no-arg.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-no-error.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-primitive.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-resolved-error.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/asyncHelpers-throwsAsync-same-realm.js | 1 | 0 | 2 | 0 | 0 | 2 |
| harness/asyncHelpers-throwsAsync-single-arg.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/byteConversionValues.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-arguments.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-arraylike.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-different-elements.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-different-length.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-empty.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-falsy-arguments.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-message.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-same-elements-different-order.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-same-elements-same-order.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-samevalue.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-sparse.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/compare-array-symbol.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/dateConstants.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/decimalToHexString.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/deepEqual-array.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/deepEqual-circular.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/deepEqual-deep.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/deepEqual-mapset.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/deepEqual-object.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/deepEqual-primitives-bigint.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/deepEqual-primitives.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/detachArrayBuffer-host-detachArrayBuffer.js | 1 | 0 | 2 | 0 | 0 | 2 |
| harness/detachArrayBuffer.js | 1 | 0 | 2 | 0 | 0 | 2 |
| harness/fnGlobalObject.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/isConstructor.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/nans.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/nativeFunctionMatcher.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/promiseHelper.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifyconfigurable-configurable-object.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifyconfigurable-configurable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifyconfigurable-not-configurable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifyenumerable-enumerable-symbol.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifyenumerable-enumerable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifyenumerable-not-enumerable-symbol.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifyenumerable-not-enumerable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifynotconfigurable-configurable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifynotconfigurable-not-configurable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifynotenumerable-enumerable-symbol.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifynotenumerable-enumerable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifynotenumerable-not-enumerable-symbol.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifynotenumerable-not-enumerable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifynotwritable-not-writable-strict.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifynotwritable-writable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifywritable-array-length.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifywritable-not-writable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/propertyhelper-verifywritable-writable.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/proxytrapshelper-default.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/proxytrapshelper-overrides.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/sta.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/tcoHelper.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/testTypedArray-conversions-call-error.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/testTypedArray-conversions.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/testTypedArray.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-arguments.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-configurable-object.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-desc-is-not-object.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-noproperty.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-restore-accessor-symbol.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-restore-accessor.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-restore-symbol.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-restore.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-same-value.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-string-prop.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-symbol-prop.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-undefined-desc.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-value-error.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/verifyProperty-value.js | 1 | 0 | 2 | 0 | 2 | 0 |
| harness/wellKnownIntrinsicObjects.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/Array | 2 | 0 | 2 | 0 | 2 | 0 |
| intl402/BigInt | 11 | 0 | 2 | 0 | 2 | 0 |
| intl402/Collator | 65 | 1 | 2 | 0 | 2 | 0 |
| intl402/Date | 12 | 0 | 2 | 0 | 2 | 0 |
| intl402/DateTimeFormat | 245 | 1 | 2 | 0 | 2 | 0 |
| intl402/DisplayNames | 57 | 3 | 2 | 0 | 2 | 0 |
| intl402/DurationFormat | 110 | 0 | 2 | 0 | 2 | 0 |
| intl402/FallbackSymbol | 2 | 0 | 2 | 0 | 0 | 2 |
| intl402/Intl | 66 | 0 | 2 | 0 | 2 | 0 |
| intl402/ListFormat | 81 | 1 | 2 | 0 | 2 | 0 |
| intl402/Locale | 190 | 1 | 2 | 0 | 2 | 0 |
| intl402/Number | 7 | 0 | 2 | 0 | 2 | 0 |
| intl402/NumberFormat | 251 | 1 | 2 | 0 | 2 | 0 |
| intl402/PluralRules | 53 | 1 | 2 | 0 | 2 | 0 |
| intl402/RelativeTimeFormat | 80 | 1 | 2 | 0 | 2 | 0 |
| intl402/Segmenter | 79 | 4 | 2 | 0 | 2 | 0 |
| intl402/String | 19 | 0 | 2 | 0 | 2 | 0 |
| intl402/Temporal | 2029 | 0 | 2 | 0 | 2 | 0 |
| intl402/TypedArray | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/constructors-string-and-single-element-array.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/constructors-taint-Object-prototype-2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/constructors-taint-Object-prototype.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/default-locale-is-canonicalized.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/default-locale-is-supported.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/fallback-locales-are-supported.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/language-tags-canonicalized.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/language-tags-invalid.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/language-tags-valid.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/language-tags-with-underscore.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-consistent-with-resolvedOptions.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-default-locale-and-zxx-locale.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-duplicate-elements-removed.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-empty-and-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-locales-arg-coered-to-object.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-locales-arg-empty-array.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-returned-array-elements-are-not-frozen.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-taint-Array-2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-taint-Array.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-test-option-localeMatcher.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-throws-if-element-not-string-or-object.js | 1 | 0 | 2 | 0 | 2 | 0 |
| intl402/supportedLocalesOf-unicode-extensions-ignored.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.5-1-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.5-1gs.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.5-7-b-1-s.js | 1 | 1 | 0 | 0 | 0 | 0 |
| language/arguments-object/10.5-7-b-2-s.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.5-7-b-3-s.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.5-7-b-4-s.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.6-10-c-ii-1-s.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-10-c-ii-1.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-10-c-ii-2.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-11-b-1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.6-12-1.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-12-2.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-13-a-1.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-13-a-2.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-13-a-3.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-13-c-1-s.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-13-c-2-s.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.6-13-c-3-s.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-14-c-1-s.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.6-14-c-4-s.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-2gs.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-5-1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.6-6-1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.6-6-2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.6-6-3-s.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.6-6-3.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-6-4-s.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/10.6-6-4.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/10.6-7-1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.1.6_A1_T2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.6_A1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.6_A2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.6_A3_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.6_A3_T2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.6_A3_T3.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/S10.6_A3_T4.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/S10.6_A4.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/S10.6_A5_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.6_A5_T2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.6_A5_T3.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.6_A5_T4.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.6_A6.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/S10.6_A7.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/arguments-caller.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/async-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/async-gen-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/async-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/async-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/async-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/async-gen-named-func-expr-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/async-gen-named-func-expr-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/async-gen-named-func-expr-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/async-gen-named-func-expr-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/async-gen-named-func-expr-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-gen-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-gen-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/func-decl-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/func-decl-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/func-decl-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/func-decl-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/func-decl-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/func-expr-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/func-expr-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/func-expr-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/func-expr-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/func-expr-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-func-decl-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-func-decl-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-func-decl-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-func-decl-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-func-decl-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-func-expr-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-func-expr-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-func-expr-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-func-expr-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-func-expr-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/gen-meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/mapped | 43 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/meth-args-trailing-comma-multiple.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/meth-args-trailing-comma-null.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/meth-args-trailing-comma-single-args.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/meth-args-trailing-comma-spread-operator.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/meth-args-trailing-comma-undefined.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/arguments-object/non-strict-arguments-object-is-immutable.js | 1 | 0 | 1 | 0 | 1 | 0 |
| language/arguments-object/unmapped | 5 | 0 | 1 | 0 | 1 | 0 |
| language/asi/S7.9.2_A1_T1.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9.2_A1_T2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9.2_A1_T3.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9.2_A1_T4.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9.2_A1_T5.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9.2_A1_T6.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9.2_A1_T7.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A10_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A10_T10.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A10_T11.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A10_T12.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A10_T2.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A10_T3.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A10_T4.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A10_T5.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A10_T6.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A10_T7.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A10_T8.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A10_T9.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A11_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A11_T10.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A11_T11.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A11_T2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A11_T3.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A11_T4.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A11_T5.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A11_T6.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A11_T7.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A11_T8.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A11_T9.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A3.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A4.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.1_T1.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A5.2_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.3_T1.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A5.4_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.5_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.5_T2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.5_T3.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.5_T4.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.5_T5.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.6_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.6_T2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.7_T1.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A5.8_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A5.9_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T1.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T10.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T11.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T12.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T13.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T2.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T3.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T4.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T5.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T6.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T7.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T8.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.1_T9.js | 1 | 0 | 2 | 0 | 2 | 0 |
| language/asi/S7.9_A6.2_T1.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.2_T10.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.2_T2.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.2_T3.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.2_T4.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.2_T5.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.2_T6.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.2_T7.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.2_T8.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.2_T9.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.3_T1.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.3_T2.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.3_T3.js | 1 | 0 | 2 | 2 | 0 | 0 |
| language/asi/S7.9_A6.3_T4.js | 1 | 0 | 1 | 1 | 0 | 0 |
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
| language/expressions/class | 4059 | 98 | 0 | 0 | 0 | 0 |
| language/expressions/coalesce | 24 | 0 | 0 | 0 | 0 | 0 |
| language/expressions/comma | 6 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/compound-assignment | 454 | 20 | 0 | 0 | 0 | 0 |
| language/expressions/concatenation | 5 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/conditional | 22 | 1 | 0 | 0 | 0 | 0 |
| language/expressions/delete | 69 | 1 | 0 | 0 | 0 | 0 |
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
| language/expressions/typeof | 16 | 2 | 0 | 0 | 0 | 0 |
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
| language/literals/null | 3 | 0 | 0 | 0 | 0 | 0 |
| language/literals/numeric | 157 | 1 | 0 | 0 | 0 | 0 |
| language/literals/regexp | 240 | 23 | 0 | 0 | 0 | 0 |
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
| language/statements/class | 4367 | 108 | 0 | 0 | 0 | 0 |
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
| staging/built-ins | 8 | 0 | 0 | 0 | 0 | 0 |
| staging/decorators | 3 | 2 | 0 | 0 | 0 | 0 |
| staging/explicit-resource-management | 53 | 0 | 0 | 0 | 0 | 0 |
| staging/set-is-subset-of-empty-index.js | 1 | 0 | 0 | 0 | 0 | 0 |
| staging/set-is-subset-on-set-like.js | 1 | 0 | 0 | 0 | 0 | 0 |
| staging/set-is-subset-table-receiver-cleared.js | 1 | 0 | 0 | 0 | 0 | 0 |
| staging/set-is-subset-table-transition.js | 1 | 0 | 0 | 0 | 0 | 0 |
| staging/set-methods | 3 | 0 | 0 | 0 | 0 | 0 |
| staging/sm | 1406 | 246 | 0 | 0 | 0 | 0 |
| staging/source-phase-imports | 2 | 0 | 0 | 0 | 0 | 0 |
| staging/top-level-await | 1 | 0 | 0 | 0 | 0 | 0 |

## Failure checklist

### annexB/built-ins

- [ ] `annexB/built-ins/Date/prototype/setYear/B.2.5.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/B.2.5.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/year-to-number-err.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/year-to-number-err.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/RegExp/RegExp-decimal-escape-class-range.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:18259:61: cannot convert tsTemp587 (variable of struct type tsOptional[float64]) to type float64

```

- [ ] `annexB/built-ins/RegExp/RegExp-decimal-escape-class-range.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:18256:61: cannot convert tsTemp586 (variable of struct type tsOptional[float64]) to type float64

```

- [ ] `annexB/built-ins/String/prototype/anchor/B.2.3.2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/anchor/B.2.3.2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/fontsize/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/fontsize/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/this-obj-not-regexp.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/this-obj-not-regexp.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `annexB/built-ins/String/prototype/trimLeft/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/trimLeft/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/sub/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/sub/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/fixed/this-val-tostring-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/fixed/this-val-tostring-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/substr/length-falsey.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/substr/length-falsey.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/strike/this-val-tostring-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/leftContext/this-cross-realm-constructor.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/leftContext/this-cross-realm-constructor.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/String/prototype/strike/this-val-tostring-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/substr/length-positive.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/substr/length-positive.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/bold/this-val-tostring-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/unescape/argument_types.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/argument_types.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/String/prototype/blink/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/blink/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/bold/this-val-tostring-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/Date/prototype/toGMTString/value.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/match/custom-matcher-emulates-undefined.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/String/prototype/match/custom-matcher-emulates-undefined.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/String/prototype/big/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/big/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/blink/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/blink/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/toGMTString/value.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/Date/prototype/setYear/date-value-read-before-tonumber-when-date-is-valid.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/Date/prototype/setYear/date-value-read-before-tonumber-when-date-is-valid.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/Date/prototype/setYear/year-number-relative.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/link/B.2.3.10.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/link/B.2.3.10.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/italics/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/italics/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/escape/escape-above-astral.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/escape/escape-above-astral.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/String/prototype/replace/custom-replacer-emulates-undefined.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/String/prototype/replace/custom-replacer-emulates-undefined.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/String/prototype/fontcolor/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/fontcolor/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/year-number-relative.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/Date/prototype/setYear/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/duplicate-named-capturing-groups-syntax.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/duplicate-named-capturing-groups-syntax.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}

```

- [ ] `annexB/built-ins/escape/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/escape/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/replaceAll/custom-replacer-emulates-undefined.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/String/prototype/replaceAll/custom-replacer-emulates-undefined.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/Date/prototype/setYear/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/year-number-absolute.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/Date/prototype/setYear/year-number-absolute.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/pattern-regexp-same.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/pattern-regexp-same.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/flags/order-after-compile.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/flags/order-after-compile.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/link/attr-tostring-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/link/attr-tostring-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/sub/this-val-tostring-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/unescape/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/Date/prototype/toGMTString/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/toGMTString/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/sub/this-val-tostring-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/Date/prototype/getYear/nan.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/escape/escape-below.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/escape/escape-below.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/String/prototype/fontsize/B.2.3.8.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/fontsize/B.2.3.8.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/link/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/link/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/bold/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/bold/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/getYear/nan.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/substr/this-non-obj-coerce.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/substr/this-non-obj-coerce.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/B.RegExp.prototype.compile.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/B.RegExp.prototype.compile.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/toGMTString/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/Date/prototype/toGMTString/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/escape/to-string-err-symbol.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/escape/to-string-err-symbol.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/escape/to-string-observe.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/escape/to-string-observe.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/String/prototype/italics/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/italics/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/sup/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/sup/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/flags-undefined.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/flags-undefined.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/unescape/to-primitive-observe.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/to-primitive-observe.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/input/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/input/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/this-not-date.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/this-not-date.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/bold/B.2.3.5.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/bold/B.2.3.5.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/unescape/argument_bigint.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/argument_bigint.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/four-ignore-non-hex.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/four-ignore-non-hex.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/String/prototype/small/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/small/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/trimRight/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/split/custom-splitter-emulates-undefined.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/String/prototype/split/custom-splitter-emulates-undefined.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/String/prototype/trimRight/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/fixed/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/fixed/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/small/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/small/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/escape/escape-above.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/escape/escape-above.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/this-time-nan.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/pattern-to-string-err.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/this-time-nan.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/pattern-to-string-err.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/sub/B.2.3.13.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/sub/B.2.3.13.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/sup/B.2.3.14.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/sup/B.2.3.14.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/substr/length-to-int-err.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/substr/length-to-int-err.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/italics/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/italics/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/fontcolor/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/fontcolor/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/small/B.2.3.11.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/small/B.2.3.11.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/RegExp/incomplete_hex_unicode_escape.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/incomplete_hex_unicode_escape.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/fontcolor/this-val-tostring-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/fontcolor/this-val-tostring-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/pattern-string.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/pattern-string.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/fontsize/attr-tostring-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/small/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/small/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/getYear/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/Date/prototype/getYear/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/fontsize/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/fontsize/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/escape/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/escape/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier escape"}

```

- [ ] `annexB/built-ins/String/prototype/fontsize/attr-tostring-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/big/this-val-tostring-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/flags-string-invalid.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/flags-string-invalid.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}

```

- [ ] `annexB/built-ins/String/prototype/big/this-val-tostring-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/italics/this-val-tostring-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/fontsize/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/fontsize/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/italics/this-val-tostring-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/unescape/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/to-primitive-err.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/to-primitive-err.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/pattern-regexp-distinct.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/prototype/compile/pattern-regexp-distinct.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/trimLeft/reference-trimStart.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/matchAll/custom-matcher-emulates-undefined.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/String/prototype/matchAll/custom-matcher-emulates-undefined.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `annexB/built-ins/String/prototype/trimLeft/reference-trimStart.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/Date/prototype/setYear/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/Date/prototype/setYear/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/big/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/big/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/substr/this-to-str-err.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/RegExp/RegExp-invalid-control-escape-character-class.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `annexB/built-ins/RegExp/RegExp-invalid-control-escape-character-class.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `annexB/built-ins/String/prototype/strike/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/strike/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/big/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/big/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `annexB/built-ins/String/prototype/substr/this-to-str-err.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `annexB/built-ins/String/prototype/italics/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/italics/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/unescape/four-ignore-bad-u.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/four-ignore-bad-u.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/to-string-observe.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/unescape/to-string-observe.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}
{"diagnostic":"Go target: unsupported or unresolved identifier unescape"}

```

- [ ] `annexB/built-ins/String/prototype/link/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/link/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/sup/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/String/prototype/sup/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/index/this-subclass-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `annexB/built-ins/RegExp/legacy-accessors/index/this-subclass-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

### annexB/language

- [ ] `annexB/language/global-code/if-decl-else-decl-a-global-existing-fn-update.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### built-ins/AbstractModuleSource

- [ ] `built-ins/AbstractModuleSource/throw-from-constructor.js` (module): **runner-module-unimplemented**

```text
Module graph/realm harness needs a dedicated adapter; not a pass.
```

### built-ins/AggregateError

- [ ] `built-ins/AggregateError/prototype/proto.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier AggregateError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `built-ins/AggregateError/prototype/proto.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier AggregateError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### built-ins/Array

- [ ] `built-ins/Array/prototype/copyWithin/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/Array/prototype/copyWithin/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/ArrayBuffer

- [ ] `built-ins/ArrayBuffer/prototype/resize/resize-same-size-zero-implicit.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `built-ins/ArrayBuffer/prototype/resize/resize-same-size-zero-implicit.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### built-ins/ArrayIteratorPrototype

- [ ] `built-ins/ArrayIteratorPrototype/Symbol.toStringTag/value-from-to-string.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `built-ins/ArrayIteratorPrototype/Symbol.toStringTag/value-from-to-string.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### built-ins/AsyncDisposableStack

- [ ] `built-ins/AsyncDisposableStack/prototype/defer/adds-onDisposeAsync.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier AsyncDisposableStack"}

```

- [ ] `built-ins/AsyncDisposableStack/prototype/defer/adds-onDisposeAsync.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier AsyncDisposableStack"}

```

### built-ins/AsyncFromSyncIteratorPrototype

- [ ] `built-ins/AsyncFromSyncIteratorPrototype/return/absent-value-not-passed.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `built-ins/AsyncFromSyncIteratorPrototype/return/absent-value-not-passed.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### built-ins/AsyncFunction

- [ ] `built-ins/AsyncFunction/AsyncFunctionPrototype-is-not-callable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/AsyncFunction/AsyncFunctionPrototype-is-not-callable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/AsyncGeneratorFunction

- [ ] `built-ins/AsyncGeneratorFunction/prototype/not-callable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `built-ins/AsyncGeneratorFunction/prototype/not-callable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### built-ins/AsyncGeneratorPrototype

- [ ] `built-ins/AsyncGeneratorPrototype/next/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `built-ins/AsyncGeneratorPrototype/next/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### built-ins/AsyncIteratorPrototype

- [ ] `built-ins/AsyncIteratorPrototype/Symbol.asyncDispose/return-val.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `built-ins/AsyncIteratorPrototype/Symbol.asyncDispose/return-val.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### built-ins/Atomics

- [ ] `built-ins/Atomics/or/non-shared-int-views-throws.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Atomics"}

```

- [ ] `built-ins/Atomics/or/non-shared-int-views-throws.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Atomics"}

```

### built-ins/BigInt

- [ ] `built-ins/BigInt/constructor-integer.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `built-ins/BigInt/constructor-integer.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### built-ins/Boolean

- [ ] `built-ins/Boolean/prototype/toString/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}

```

- [ ] `built-ins/Boolean/prototype/toString/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}

```

### built-ins/DataView

- [ ] `built-ins/DataView/prototype/setUint16/this-has-no-dataview-internal.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier DataView"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/DataView/prototype/setUint16/this-has-no-dataview-internal.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier DataView"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/Date

- [ ] `built-ins/Date/prototype/setUTCMilliseconds/date-value-read-before-tonumber-when-date-is-invalid.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `built-ins/Date/prototype/setUTCMilliseconds/date-value-read-before-tonumber-when-date-is-invalid.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### built-ins/DisposableStack

- [ ] `built-ins/DisposableStack/prototype/defer/throws-if-disposed.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier DisposableStack"}
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}

```

- [ ] `built-ins/DisposableStack/prototype/defer/throws-if-disposed.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier DisposableStack"}
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}

```

### built-ins/Error

- [ ] `built-ins/Error/prototype/stack/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `built-ins/Error/prototype/stack/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### built-ins/FinalizationRegistry

- [ ] `built-ins/FinalizationRegistry/newtarget-prototype-is-not-object.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}

```

- [ ] `built-ins/FinalizationRegistry/newtarget-prototype-is-not-object.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}

```

### built-ins/Function

- [ ] `built-ins/Function/S15.3.3_A1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}

```

- [ ] `built-ins/Function/S15.3.3_A1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}

```

### built-ins/GeneratorFunction

- [ ] `built-ins/GeneratorFunction/prototype/Symbol.toStringTag.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `built-ins/GeneratorFunction/prototype/Symbol.toStringTag.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### built-ins/GeneratorPrototype

- [ ] `built-ins/GeneratorPrototype/next/consecutive-yields.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `built-ins/GeneratorPrototype/next/consecutive-yields.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### built-ins/Infinity

- [ ] `built-ins/Infinity/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

- [ ] `built-ins/Infinity/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

### built-ins/Iterator

- [ ] `built-ins/Iterator/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Iterator"}

```

- [ ] `built-ins/Iterator/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Iterator"}

```

### built-ins/JSON

- [ ] `built-ins/JSON/parse/duplicate-proto.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `built-ins/JSON/parse/duplicate-proto.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### built-ins/Map

- [ ] `built-ins/Map/prototype/has/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `built-ins/Map/prototype/has/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### built-ins/MapIteratorPrototype

- [ ] `built-ins/MapIteratorPrototype/next/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `built-ins/MapIteratorPrototype/next/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### built-ins/Math

- [ ] `built-ins/Math/max/not-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}

```

- [ ] `built-ins/Math/max/not-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}

```

### built-ins/NaN

- [ ] `built-ins/NaN/15.1.1.1-0.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

- [ ] `built-ins/NaN/15.1.1.1-0.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

### built-ins/NativeErrors

- [ ] `built-ins/NativeErrors/ReferenceError/proto.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `built-ins/NativeErrors/ReferenceError/proto.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### built-ins/Number

- [ ] `built-ins/Number/prototype/toExponential/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `built-ins/Number/prototype/toExponential/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### built-ins/Object

- [ ] `built-ins/Object/hasOwn/toobject_null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/Object/hasOwn/toobject_null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/Promise

- [ ] `built-ins/Promise/all/iter-returns-symbol-reject.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/Promise/all/iter-returns-symbol-reject.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/Proxy

- [ ] `built-ins/Proxy/get/trap-is-not-callable-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `built-ins/Proxy/get/trap-is-not-callable-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

### built-ins/Reflect

- [ ] `built-ins/Reflect/getOwnPropertyDescriptor/return-abrupt-from-result.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Proxy"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}

```

- [ ] `built-ins/Reflect/getOwnPropertyDescriptor/return-abrupt-from-result.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Proxy"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}

```

### built-ins/RegExp

- [ ] `built-ins/RegExp/property-escapes/generated/Script_-_Cypro_Minoan.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `built-ins/RegExp/property-escapes/generated/Script_-_Cypro_Minoan.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### built-ins/RegExpStringIteratorPrototype

- [ ] `built-ins/RegExpStringIteratorPrototype/next/next-iteration.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `built-ins/RegExpStringIteratorPrototype/next/next-iteration.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### built-ins/Set

- [ ] `built-ins/Set/prototype/values/does-not-have-setdata-internal-slot-array.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/Set/prototype/values/does-not-have-setdata-internal-slot-array.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/SetIteratorPrototype

- [ ] `built-ins/SetIteratorPrototype/Symbol.toStringTag.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `built-ins/SetIteratorPrototype/Symbol.toStringTag.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### built-ins/ShadowRealm

- [ ] `built-ins/ShadowRealm/prototype/evaluate/wrapped-function-observing-their-scopes.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier ShadowRealm"}
{"diagnostic":"Go target: unsupported or unresolved identifier ShadowRealm"}

```

- [ ] `built-ins/ShadowRealm/prototype/evaluate/wrapped-function-observing-their-scopes.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier ShadowRealm"}
{"diagnostic":"Go target: unsupported or unresolved identifier ShadowRealm"}

```

### built-ins/SharedArrayBuffer

- [ ] `built-ins/SharedArrayBuffer/is-a-constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier SharedArrayBuffer"}
{"diagnostic":"Go target: unsupported or unresolved identifier SharedArrayBuffer"}

```

- [ ] `built-ins/SharedArrayBuffer/is-a-constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier SharedArrayBuffer"}
{"diagnostic":"Go target: unsupported or unresolved identifier SharedArrayBuffer"}

```

### built-ins/String

- [ ] `built-ins/String/prototype/localeCompare/S15.5.4.9_A8.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `built-ins/String/prototype/localeCompare/S15.5.4.9_A8.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### built-ins/StringIteratorPrototype

- [ ] `built-ins/StringIteratorPrototype/next/length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `built-ins/StringIteratorPrototype/next/length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### built-ins/SuppressedError

- [ ] `built-ins/SuppressedError/prototype/constructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier SuppressedError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SuppressedError"}

```

- [ ] `built-ins/SuppressedError/prototype/constructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier SuppressedError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SuppressedError"}

```

### built-ins/Symbol

- [ ] `built-ins/Symbol/symbol.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

- [ ] `built-ins/Symbol/symbol.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

### built-ins/Temporal

- [ ] `built-ins/Temporal/ZonedDateTime/prototype/equals/argument-propertybag-offset-not-agreeing-with-timezone.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Temporal"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}

```

- [ ] `built-ins/Temporal/ZonedDateTime/prototype/equals/argument-propertybag-offset-not-agreeing-with-timezone.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Temporal"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}

```

### built-ins/ThrowTypeError

- [ ] `built-ins/ThrowTypeError/forbidden-caller.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `built-ins/ThrowTypeError/forbidden-caller.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### built-ins/TypedArray

- [ ] `built-ins/TypedArray/prototype/some/return-abrupt-from-this-out-of-bounds.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/TypedArray/prototype/some/return-abrupt-from-this-out-of-bounds.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/TypedArrayConstructors

- [ ] `built-ins/TypedArrayConstructors/ctors/buffer-arg/length-is-symbol-throws-sab.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier SharedArrayBuffer"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"d
```

- [ ] `built-ins/TypedArrayConstructors/ctors/buffer-arg/length-is-symbol-throws-sab.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier SharedArrayBuffer"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"d
```

### built-ins/Uint8Array

- [ ] `built-ins/Uint8Array/prototype/setFromBase64/nonconstructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/Uint8Array/prototype/setFromBase64/nonconstructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/WeakMap

- [ ] `built-ins/WeakMap/prototype/getOrInsertComputed/this-not-object-throw.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/WeakMap/prototype/getOrInsertComputed/this-not-object-throw.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/WeakRef

- [ ] `built-ins/WeakRef/prototype/deref/this-does-not-have-internal-target-throws.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier WeakRef"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakRef"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakRef"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakRef"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakSet"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/WeakRef/prototype/deref/this-does-not-have-internal-target-throws.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier WeakRef"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakRef"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakRef"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakRef"}
{"diagnostic":"Go target: unsupported or unresolved identifier FinalizationRegistry"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakSet"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/WeakSet

- [ ] `built-ins/WeakSet/prototype/constructor/weakset-prototype-constructor-intrinsic.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier WeakSet"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakSet"}

```

- [ ] `built-ins/WeakSet/prototype/constructor/weakset-prototype-constructor-intrinsic.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier WeakSet"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakSet"}

```

### built-ins/decodeURI

- [ ] `built-ins/decodeURI/S15.1.3.1_A1.8_T2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier decodeURI"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}

```

- [ ] `built-ins/decodeURI/S15.1.3.1_A1.8_T2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier decodeURI"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}

```

### built-ins/decodeURIComponent

- [ ] `built-ins/decodeURIComponent/S15.1.3.2_A1.15_T6.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier decodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}

```

- [ ] `built-ins/decodeURIComponent/S15.1.3.2_A1.15_T6.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier decodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}

```

### built-ins/encodeURI

- [ ] `built-ins/encodeURI/S15.1.3.3_A4_T1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURI"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURI"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURI"}

```

- [ ] `built-ins/encodeURI/S15.1.3.3_A4_T1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURI"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURI"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURI"}

```

### built-ins/encodeURIComponent

- [ ] `built-ins/encodeURIComponent/S15.1.3.4_A6_T1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `built-ins/encodeURIComponent/S15.1.3.4_A6_T1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier encodeURIComponent"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### built-ins/global

- [ ] `built-ins/global/S10.2.3_A2.1_T2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

- [ ] `built-ins/global/S10.2.3_A2.1_T2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

### built-ins/isFinite

- [ ] `built-ins/isFinite/S15.1.2.5_A2.6.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier isFinite"}

```

- [ ] `built-ins/isFinite/S15.1.2.5_A2.6.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier isFinite"}

```

### built-ins/isNaN

- [ ] `built-ins/isNaN/return-true-nan.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}

```

- [ ] `built-ins/isNaN/return-true-nan.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}

```

### built-ins/parseFloat

- [ ] `built-ins/parseFloat/S15.1.2.3_A2_T1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}

```

- [ ] `built-ins/parseFloat/S15.1.2.3_A2_T1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseFloat"}

```

### built-ins/parseInt

- [ ] `built-ins/parseInt/S15.1.2.2_A4.1_T1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unres
```

- [ ] `built-ins/parseInt/S15.1.2.2_A4.1_T1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unresolved identifier parseInt"}
{"diagnostic":"Go target: unsupported or unres
```

### built-ins/undefined

- [ ] `built-ins/undefined/S15.1.1.3_A4.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

- [ ] `built-ins/undefined/S15.1.1.3_A4.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

### harness/assert-false.js

- [ ] `harness/assert-false.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-false.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-notsamevalue-nan.js

- [ ] `harness/assert-notsamevalue-nan.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-notsamevalue-nan.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-notsamevalue-notsame.js

- [ ] `harness/assert-notsamevalue-notsame.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/assert-notsamevalue-notsame.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/assert-notsamevalue-objects.js

- [ ] `harness/assert-notsamevalue-objects.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/assert-notsamevalue-objects.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/assert-notsamevalue-tostring.js

- [ ] `harness/assert-notsamevalue-tostring.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-notsamevalue-tostring.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-notsamevalue-zeros.js

- [ ] `harness/assert-notsamevalue-zeros.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/assert-notsamevalue-zeros.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/assert-obj.js

- [ ] `harness/assert-obj.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-obj.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-samevalue-nan.js

- [ ] `harness/assert-samevalue-nan.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/assert-samevalue-nan.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/assert-samevalue-objects.js

- [ ] `harness/assert-samevalue-objects.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-samevalue-objects.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-samevalue-same.js

- [ ] `harness/assert-samevalue-same.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/assert-samevalue-same.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/assert-samevalue-tostring.js

- [ ] `harness/assert-samevalue-tostring.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-samevalue-tostring.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-samevalue-zeros.js

- [ ] `harness/assert-samevalue-zeros.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-samevalue-zeros.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-throws-custom-typeerror.js

- [ ] `harness/assert-throws-custom-typeerror.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-throws-custom-typeerror.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-throws-custom.js

- [ ] `harness/assert-throws-custom.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/assert-throws-custom.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/assert-throws-incorrect-ctor.js

- [ ] `harness/assert-throws-incorrect-ctor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-throws-incorrect-ctor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-throws-native.js

- [ ] `harness/assert-throws-native.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier EvalError"}
{"diagnostic":"Go target: unsupported or unresolved identifier EvalError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}

```

- [ ] `harness/assert-throws-native.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier EvalError"}
{"diagnostic":"Go target: unsupported or unresolved identifier EvalError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}

```

### harness/assert-throws-no-arg.js

- [ ] `harness/assert-throws-no-arg.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-throws-no-arg.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-throws-no-error.js

- [ ] `harness/assert-throws-no-error.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-throws-no-error.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-throws-null-fn.js

- [ ] `harness/assert-throws-null-fn.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-throws-null-fn.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-throws-null.js

- [ ] `harness/assert-throws-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-throws-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-throws-primitive.js

- [ ] `harness/assert-throws-primitive.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-throws-primitive.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-throws-same-realm.js

- [ ] `harness/assert-throws-same-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `harness/assert-throws-same-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

### harness/assert-throws-single-arg.js

- [ ] `harness/assert-throws-single-arg.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-throws-single-arg.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-tostring.js

- [ ] `harness/assert-tostring.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assert-tostring.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/assert-true.js

- [ ] `harness/assert-true.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/assert-true.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/assertRelativeDateMs.js

- [ ] `harness/assertRelativeDateMs.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/assertRelativeDateMs.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/asyncHelpers-asyncTest-func-throws-sync.js

- [ ] `harness/asyncHelpers-asyncTest-func-throws-sync.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

- [ ] `harness/asyncHelpers-asyncTest-func-throws-sync.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

### harness/asyncHelpers-asyncTest-rejects-non-callable.js

- [ ] `harness/asyncHelpers-asyncTest-rejects-non-callable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

- [ ] `harness/asyncHelpers-asyncTest-rejects-non-callable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

### harness/asyncHelpers-asyncTest-return-not-thenable.js

- [ ] `harness/asyncHelpers-asyncTest-return-not-thenable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `harness/asyncHelpers-asyncTest-return-not-thenable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### harness/asyncHelpers-asyncTest-returns-undefined.js

- [ ] `harness/asyncHelpers-asyncTest-returns-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

- [ ] `harness/asyncHelpers-asyncTest-returns-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

### harness/asyncHelpers-asyncTest-then-rejects.js

- [ ] `harness/asyncHelpers-asyncTest-then-rejects.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

- [ ] `harness/asyncHelpers-asyncTest-then-rejects.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

### harness/asyncHelpers-asyncTest-then-resolves.js

- [ ] `harness/asyncHelpers-asyncTest-then-resolves.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

- [ ] `harness/asyncHelpers-asyncTest-then-resolves.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

### harness/asyncHelpers-asyncTest-without-async-flag.js

- [ ] `harness/asyncHelpers-asyncTest-without-async-flag.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier $DONE"}
{"diagnostic":"Go target: unsupported or unresolved identifier $DONE"}
{"diagnostic":"Go target: unsupported or unresolved identifier $DONE"}
{"diagnostic":"Go target: unsupported or unresolved identifier $DONE"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

- [ ] `harness/asyncHelpers-asyncTest-without-async-flag.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier $DONE"}
{"diagnostic":"Go target: unsupported or unresolved identifier $DONE"}
{"diagnostic":"Go target: unsupported or unresolved identifier $DONE"}
{"diagnostic":"Go target: unsupported or unresolved identifier $DONE"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}

```

### harness/asyncHelpers-throwsAsync-custom-typeerror.js

- [ ] `harness/asyncHelpers-throwsAsync-custom-typeerror.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-custom-typeerror.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/asyncHelpers-throwsAsync-custom.js

- [ ] `harness/asyncHelpers-throwsAsync-custom.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-custom.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/asyncHelpers-throwsAsync-func-never-settles.js

- [ ] `harness/asyncHelpers-throwsAsync-func-never-settles.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-func-never-settles.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/asyncHelpers-throwsAsync-func-throws-sync.js

- [ ] `harness/asyncHelpers-throwsAsync-func-throws-sync.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/asyncHelpers-throwsAsync-func-throws-sync.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/asyncHelpers-throwsAsync-incorrect-ctor.js

- [ ] `harness/asyncHelpers-throwsAsync-incorrect-ctor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-incorrect-ctor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/asyncHelpers-throwsAsync-invalid-func.js

- [ ] `harness/asyncHelpers-throwsAsync-invalid-func.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-invalid-func.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/asyncHelpers-throwsAsync-native.js

- [ ] `harness/asyncHelpers-throwsAsync-native.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier EvalError"}
{"diagnostic":"Go target: unsupported or unresolved identifier EvalError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}
{"diagnostic":"Go target: unsu
```

- [ ] `harness/asyncHelpers-throwsAsync-native.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier EvalError"}
{"diagnostic":"Go target: unsupported or unresolved identifier EvalError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}
{"diagnostic":"Go target: unsupported or unresolved identifier ReferenceError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}
{"diagnostic":"Go target: unsupported or unresolved identifier URIError"}
{"diagnostic":"Go target: unsu
```

### harness/asyncHelpers-throwsAsync-no-arg.js

- [ ] `harness/asyncHelpers-throwsAsync-no-arg.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-no-arg.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/asyncHelpers-throwsAsync-no-error.js

- [ ] `harness/asyncHelpers-throwsAsync-no-error.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-no-error.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/asyncHelpers-throwsAsync-null.js

- [ ] `harness/asyncHelpers-throwsAsync-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/asyncHelpers-throwsAsync-primitive.js

- [ ] `harness/asyncHelpers-throwsAsync-primitive.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-primitive.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/asyncHelpers-throwsAsync-resolved-error.js

- [ ] `harness/asyncHelpers-throwsAsync-resolved-error.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-resolved-error.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/asyncHelpers-throwsAsync-same-realm.js

- [ ] `harness/asyncHelpers-throwsAsync-same-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `harness/asyncHelpers-throwsAsync-same-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

### harness/asyncHelpers-throwsAsync-single-arg.js

- [ ] `harness/asyncHelpers-throwsAsync-single-arg.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

- [ ] `harness/asyncHelpers-throwsAsync-single-arg.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier globalThis"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}

```

### harness/byteConversionValues.js

- [ ] `harness/byteConversionValues.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/byteConversionValues.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/compare-array-arguments.js

- [ ] `harness/compare-array-arguments.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `harness/compare-array-arguments.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### harness/compare-array-arraylike.js

- [ ] `harness/compare-array-arraylike.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/compare-array-arraylike.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/compare-array-different-elements.js

- [ ] `harness/compare-array-different-elements.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/compare-array-different-elements.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/compare-array-different-length.js

- [ ] `harness/compare-array-different-length.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/compare-array-different-length.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/compare-array-empty.js

- [ ] `harness/compare-array-empty.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/compare-array-empty.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/compare-array-falsy-arguments.js

- [ ] `harness/compare-array-falsy-arguments.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/compare-array-falsy-arguments.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/compare-array-message.js

- [ ] `harness/compare-array-message.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `harness/compare-array-message.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### harness/compare-array-same-elements-different-order.js

- [ ] `harness/compare-array-same-elements-different-order.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/compare-array-same-elements-different-order.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/compare-array-same-elements-same-order.js

- [ ] `harness/compare-array-same-elements-same-order.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/compare-array-same-elements-same-order.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/compare-array-samevalue.js

- [ ] `harness/compare-array-samevalue.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/compare-array-samevalue.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/compare-array-sparse.js

- [ ] `harness/compare-array-sparse.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/compare-array-sparse.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/compare-array-symbol.js

- [ ] `harness/compare-array-symbol.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `harness/compare-array-symbol.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### harness/dateConstants.js

- [ ] `harness/dateConstants.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/dateConstants.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/decimalToHexString.js

- [ ] `harness/decimalToHexString.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/decimalToHexString.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/deepEqual-array.js

- [ ] `harness/deepEqual-array.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

- [ ] `harness/deepEqual-array.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

### harness/deepEqual-circular.js

- [ ] `harness/deepEqual-circular.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

- [ ] `harness/deepEqual-circular.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

### harness/deepEqual-deep.js

- [ ] `harness/deepEqual-deep.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

- [ ] `harness/deepEqual-deep.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

### harness/deepEqual-mapset.js

- [ ] `harness/deepEqual-mapset.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

- [ ] `harness/deepEqual-mapset.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

### harness/deepEqual-object.js

- [ ] `harness/deepEqual-object.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

- [ ] `harness/deepEqual-object.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

### harness/deepEqual-primitives-bigint.js

- [ ] `harness/deepEqual-primitives-bigint.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

- [ ] `harness/deepEqual-primitives-bigint.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

### harness/deepEqual-primitives.js

- [ ] `harness/deepEqual-primitives.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

- [ ] `harness/deepEqual-primitives.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier isNaN"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier Promise"}
{"diagnostic":"Go target: unsupported or unresolved identifier WeakMap"}
{"diagnostic":"Go target: unsupported or unresolved iden
```

### harness/detachArrayBuffer-host-detachArrayBuffer.js

- [ ] `harness/detachArrayBuffer-host-detachArrayBuffer.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `harness/detachArrayBuffer-host-detachArrayBuffer.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

### harness/detachArrayBuffer.js

- [ ] `harness/detachArrayBuffer.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `harness/detachArrayBuffer.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

### harness/fnGlobalObject.js

- [ ] `harness/fnGlobalObject.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

- [ ] `harness/fnGlobalObject.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

### harness/isConstructor.js

- [ ] `harness/isConstructor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `harness/isConstructor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### harness/nans.js

- [ ] `harness/nans.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}

```

- [ ] `harness/nans.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}

```

### harness/nativeFunctionMatcher.js

- [ ] `harness/nativeFunctionMatcher.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/nativeFunctionMatcher.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier SyntaxError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/promiseHelper.js

- [ ] `harness/promiseHelper.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/promiseHelper.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/propertyhelper-verifyconfigurable-configurable-object.js

- [ ] `harness/propertyhelper-verifyconfigurable-configurable-object.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

- [ ] `harness/propertyhelper-verifyconfigurable-configurable-object.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

### harness/propertyhelper-verifyconfigurable-configurable.js

- [ ] `harness/propertyhelper-verifyconfigurable-configurable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/propertyhelper-verifyconfigurable-configurable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/propertyhelper-verifyconfigurable-not-configurable.js

- [ ] `harness/propertyhelper-verifyconfigurable-not-configurable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/propertyhelper-verifyconfigurable-not-configurable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/propertyhelper-verifyenumerable-enumerable-symbol.js

- [ ] `harness/propertyhelper-verifyenumerable-enumerable-symbol.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `harness/propertyhelper-verifyenumerable-enumerable-symbol.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### harness/propertyhelper-verifyenumerable-enumerable.js

- [ ] `harness/propertyhelper-verifyenumerable-enumerable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/propertyhelper-verifyenumerable-enumerable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/propertyhelper-verifyenumerable-not-enumerable-symbol.js

- [ ] `harness/propertyhelper-verifyenumerable-not-enumerable-symbol.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/propertyhelper-verifyenumerable-not-enumerable-symbol.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/propertyhelper-verifyenumerable-not-enumerable.js

- [ ] `harness/propertyhelper-verifyenumerable-not-enumerable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/propertyhelper-verifyenumerable-not-enumerable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/propertyhelper-verifynotconfigurable-configurable.js

- [ ] `harness/propertyhelper-verifynotconfigurable-configurable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/propertyhelper-verifynotconfigurable-configurable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/propertyhelper-verifynotconfigurable-not-configurable.js

- [ ] `harness/propertyhelper-verifynotconfigurable-not-configurable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/propertyhelper-verifynotconfigurable-not-configurable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/propertyhelper-verifynotenumerable-enumerable-symbol.js

- [ ] `harness/propertyhelper-verifynotenumerable-enumerable-symbol.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `harness/propertyhelper-verifynotenumerable-enumerable-symbol.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### harness/propertyhelper-verifynotenumerable-enumerable.js

- [ ] `harness/propertyhelper-verifynotenumerable-enumerable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/propertyhelper-verifynotenumerable-enumerable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/propertyhelper-verifynotenumerable-not-enumerable-symbol.js

- [ ] `harness/propertyhelper-verifynotenumerable-not-enumerable-symbol.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `harness/propertyhelper-verifynotenumerable-not-enumerable-symbol.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### harness/propertyhelper-verifynotenumerable-not-enumerable.js

- [ ] `harness/propertyhelper-verifynotenumerable-not-enumerable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/propertyhelper-verifynotenumerable-not-enumerable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/propertyhelper-verifynotwritable-not-writable-strict.js

- [ ] `harness/propertyhelper-verifynotwritable-not-writable-strict.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/propertyhelper-verifynotwritable-not-writable-strict.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/propertyhelper-verifynotwritable-writable.js

- [ ] `harness/propertyhelper-verifynotwritable-writable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/propertyhelper-verifynotwritable-writable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/propertyhelper-verifywritable-array-length.js

- [ ] `harness/propertyhelper-verifywritable-array-length.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/propertyhelper-verifywritable-array-length.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/propertyhelper-verifywritable-not-writable.js

- [ ] `harness/propertyhelper-verifywritable-not-writable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/propertyhelper-verifywritable-not-writable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/propertyhelper-verifywritable-writable.js

- [ ] `harness/propertyhelper-verifywritable-writable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/propertyhelper-verifywritable-writable.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/proxytrapshelper-default.js

- [ ] `harness/proxytrapshelper-default.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/proxytrapshelper-default.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/proxytrapshelper-overrides.js

- [ ] `harness/proxytrapshelper-overrides.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/proxytrapshelper-overrides.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/sta.js

- [ ] `harness/sta.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/sta.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/tcoHelper.js

- [ ] `harness/tcoHelper.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `harness/tcoHelper.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### harness/testTypedArray-conversions-call-error.js

- [ ] `harness/testTypedArray-conversions-call-error.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic
```

- [ ] `harness/testTypedArray-conversions-call-error.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic
```

### harness/testTypedArray-conversions.js

- [ ] `harness/testTypedArray-conversions.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/testTypedArray-conversions.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/testTypedArray.js

- [ ] `harness/testTypedArray.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"
```

- [ ] `harness/testTypedArray.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"
```

### harness/verifyProperty-arguments.js

- [ ] `harness/verifyProperty-arguments.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}

```

- [ ] `harness/verifyProperty-arguments.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}

```

### harness/verifyProperty-configurable-object.js

- [ ] `harness/verifyProperty-configurable-object.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

- [ ] `harness/verifyProperty-configurable-object.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}
{"diagnostic":"Go target: this is supported only in class methods and their arrow closures"}

```

### harness/verifyProperty-desc-is-not-object.js

- [ ] `harness/verifyProperty-desc-is-not-object.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `harness/verifyProperty-desc-is-not-object.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### harness/verifyProperty-noproperty.js

- [ ] `harness/verifyProperty-noproperty.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}

```

- [ ] `harness/verifyProperty-noproperty.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}

```

### harness/verifyProperty-restore-accessor-symbol.js

- [ ] `harness/verifyProperty-restore-accessor-symbol.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `harness/verifyProperty-restore-accessor-symbol.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### harness/verifyProperty-restore-accessor.js

- [ ] `harness/verifyProperty-restore-accessor.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/verifyProperty-restore-accessor.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/verifyProperty-restore-symbol.js

- [ ] `harness/verifyProperty-restore-symbol.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `harness/verifyProperty-restore-symbol.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### harness/verifyProperty-restore.js

- [ ] `harness/verifyProperty-restore.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/verifyProperty-restore.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/verifyProperty-same-value.js

- [ ] `harness/verifyProperty-same-value.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/verifyProperty-same-value.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/verifyProperty-string-prop.js

- [ ] `harness/verifyProperty-string-prop.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/verifyProperty-string-prop.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/verifyProperty-symbol-prop.js

- [ ] `harness/verifyProperty-symbol-prop.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `harness/verifyProperty-symbol-prop.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### harness/verifyProperty-undefined-desc.js

- [ ] `harness/verifyProperty-undefined-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/verifyProperty-undefined-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/verifyProperty-value-error.js

- [ ] `harness/verifyProperty-value-error.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `harness/verifyProperty-value-error.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### harness/verifyProperty-value.js

- [ ] `harness/verifyProperty-value.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `harness/verifyProperty-value.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### harness/wellKnownIntrinsicObjects.js

- [ ] `harness/wellKnownIntrinsicObjects.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}

```

- [ ] `harness/wellKnownIntrinsicObjects.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}

```

### intl402/Array

- [ ] `intl402/Array/prototype/toLocaleString/invoke-element-tolocalestring.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `intl402/Array/prototype/toLocaleString/invoke-element-tolocalestring.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### intl402/BigInt

- [ ] `intl402/BigInt/prototype/toLocaleString/default-options-object-prototype.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `intl402/BigInt/prototype/toLocaleString/default-options-object-prototype.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### intl402/Collator

- [ ] `intl402/Collator/prototype/compare/builtin.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}

```

- [ ] `intl402/Collator/prototype/compare/builtin.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}

```

### intl402/Date

- [ ] `intl402/Date/prototype/returns-same-results-as-DateTimeFormat.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

- [ ] `intl402/Date/prototype/returns-same-results-as-DateTimeFormat.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

### intl402/DateTimeFormat

- [ ] `intl402/DateTimeFormat/prototype/formatRange/en-US.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

- [ ] `intl402/DateTimeFormat/prototype/formatRange/en-US.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

### intl402/DisplayNames

- [ ] `intl402/DisplayNames/proto.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}

```

- [ ] `intl402/DisplayNames/proto.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}

```

### intl402/DurationFormat

- [ ] `intl402/DurationFormat/prototype/formatToParts/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

- [ ] `intl402/DurationFormat/prototype/formatToParts/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

### intl402/FallbackSymbol

- [ ] `intl402/FallbackSymbol/per-realm.js` (sloppy): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

- [ ] `intl402/FallbackSymbol/per-realm.js` (strict): **runner-host-unimplemented**

```text
Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.
```

### intl402/Intl

- [ ] `intl402/Intl/getCanonicalLocales/locales-is-not-a-string.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

- [ ] `intl402/Intl/getCanonicalLocales/locales-is-not-a-string.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}

```

### intl402/ListFormat

- [ ] `intl402/ListFormat/prototype/format/iterable-invalid.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

- [ ] `intl402/ListFormat/prototype/format/iterable-invalid.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### intl402/Locale

- [ ] `intl402/Locale/prototype/minimize/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

- [ ] `intl402/Locale/prototype/minimize/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

### intl402/Number

- [ ] `intl402/Number/prototype/toLocaleString/builtin.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}

```

- [ ] `intl402/Number/prototype/toLocaleString/builtin.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Reflect"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}

```

### intl402/NumberFormat

- [ ] `intl402/NumberFormat/prototype/formatToParts/prop-desc.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

- [ ] `intl402/NumberFormat/prototype/formatToParts/prop-desc.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

### intl402/PluralRules

- [ ] `intl402/PluralRules/name.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

- [ ] `intl402/PluralRules/name.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

### intl402/RelativeTimeFormat

- [ ] `intl402/RelativeTimeFormat/prototype/format/en-us-numeric-always.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

- [ ] `intl402/RelativeTimeFormat/prototype/format/en-us-numeric-always.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

### intl402/Segmenter

- [ ] `intl402/Segmenter/constructor/constructor/newtarget-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

- [ ] `intl402/Segmenter/constructor/constructor/newtarget-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}

```

### intl402/String

- [ ] `intl402/String/prototype/localeCompare/missing-arguments-coerced-to-undefined.js` (sloppy): **build-error**

```text
# command-line-arguments
<case>/main.go:18238:29: tsTemp581.read undefined (type *tsCell has no field or method read)
<case>/main.go:18240:19: cannot use tsTemp581.set(loop, tsNumberValue(tsTemp583)) (value of struct type tsValue) as float64 value in assignment

```

- [ ] `intl402/String/prototype/localeCompare/missing-arguments-coerced-to-undefined.js` (strict): **build-error**

```text
# command-line-arguments
<case>/main.go:18241:29: tsTemp582.read undefined (type *tsCell has no field or method read)
<case>/main.go:18243:19: cannot use tsTemp582.set(loop, tsNumberValue(tsTemp584)) (value of struct type tsValue) as float64 value in assignment

```

### intl402/Temporal

- [ ] `intl402/Temporal/PlainDate/prototype/until/basic-indian.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

- [ ] `intl402/Temporal/PlainDate/prototype/until/basic-indian.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}
{"diagnostic":"Go target: base class must be a statically resolved class in the same source file"}

```

### intl402/TypedArray

- [ ] `intl402/TypedArray/prototype/toLocaleString/calls-toLocaleString-number-elements.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

- [ ] `intl402/TypedArray/prototype/toLocaleString/calls-toLocaleString-number-elements.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported statement KindJSTypeAliasDeclaration"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigInt64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier BigUint64Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Float16Array"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}

```

### intl402/constructors-string-and-single-element-array.js

- [ ] `intl402/constructors-string-and-single-element-array.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/constructors-string-and-single-element-array.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/constructors-taint-Object-prototype-2.js

- [ ] `intl402/constructors-taint-Object-prototype-2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/constructors-taint-Object-prototype-2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/constructors-taint-Object-prototype.js

- [ ] `intl402/constructors-taint-Object-prototype.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/constructors-taint-Object-prototype.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/default-locale-is-canonicalized.js

- [ ] `intl402/default-locale-is-canonicalized.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/default-locale-is-canonicalized.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/default-locale-is-supported.js

- [ ] `intl402/default-locale-is-supported.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/default-locale-is-supported.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/fallback-locales-are-supported.js

- [ ] `intl402/fallback-locales-are-supported.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/fallback-locales-are-supported.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/language-tags-canonicalized.js

- [ ] `intl402/language-tags-canonicalized.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/language-tags-canonicalized.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/language-tags-invalid.js

- [ ] `intl402/language-tags-invalid.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/language-tags-invalid.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/language-tags-valid.js

- [ ] `intl402/language-tags-valid.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/language-tags-valid.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/language-tags-with-underscore.js

- [ ] `intl402/language-tags-with-underscore.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/language-tags-with-underscore.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-consistent-with-resolvedOptions.js

- [ ] `intl402/supportedLocalesOf-consistent-with-resolvedOptions.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-consistent-with-resolvedOptions.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-default-locale-and-zxx-locale.js

- [ ] `intl402/supportedLocalesOf-default-locale-and-zxx-locale.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-default-locale-and-zxx-locale.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-duplicate-elements-removed.js

- [ ] `intl402/supportedLocalesOf-duplicate-elements-removed.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-duplicate-elements-removed.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-empty-and-undefined.js

- [ ] `intl402/supportedLocalesOf-empty-and-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-empty-and-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-locales-arg-coered-to-object.js

- [ ] `intl402/supportedLocalesOf-locales-arg-coered-to-object.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-locales-arg-coered-to-object.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-locales-arg-empty-array.js

- [ ] `intl402/supportedLocalesOf-locales-arg-empty-array.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-locales-arg-empty-array.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-returned-array-elements-are-not-frozen.js

- [ ] `intl402/supportedLocalesOf-returned-array-elements-are-not-frozen.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-returned-array-elements-are-not-frozen.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-taint-Array-2.js

- [ ] `intl402/supportedLocalesOf-taint-Array-2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-taint-Array-2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-taint-Array.js

- [ ] `intl402/supportedLocalesOf-taint-Array.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-taint-Array.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-test-option-localeMatcher.js

- [ ] `intl402/supportedLocalesOf-test-option-localeMatcher.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-test-option-localeMatcher.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-throws-if-element-not-string-or-object.js

- [ ] `intl402/supportedLocalesOf-throws-if-element-not-string-or-object.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-throws-if-element-not-string-or-object.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### intl402/supportedLocalesOf-unicode-extensions-ignored.js

- [ ] `intl402/supportedLocalesOf-unicode-extensions-ignored.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

- [ ] `intl402/supportedLocalesOf-unicode-extensions-ignored.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved identifier Error"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Symbol"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier RangeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier Boolean"}
{"diagnostic":"Go target: unsupported or unresolved identifier Intl"}
{"diagnostic":"Go target: unsupported or unresolved ide
```

### language/arguments-object/10.5-1gs.js

- [ ] `language/arguments-object/10.5-1gs.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier $DONOTEVALUATE"}
{"diagnostic":"Go target: unsupported assignment target"}

```

### language/arguments-object/10.5-7-b-2-s.js

- [ ] `language/arguments-object/10.5-7-b-2-s.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.5-7-b-2-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.5-7-b-3-s.js

- [ ] `language/arguments-object/10.5-7-b-3-s.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.5-7-b-3-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.5-7-b-4-s.js

- [ ] `language/arguments-object/10.5-7-b-4-s.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.5-7-b-4-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-10-c-ii-1-s.js

- [ ] `language/arguments-object/10.6-10-c-ii-1-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-10-c-ii-1.js

- [ ] `language/arguments-object/10.6-10-c-ii-1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-10-c-ii-2.js

- [ ] `language/arguments-object/10.6-10-c-ii-2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-11-b-1.js

- [ ] `language/arguments-object/10.6-11-b-1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.6-11-b-1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-12-1.js

- [ ] `language/arguments-object/10.6-12-1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-12-2.js

- [ ] `language/arguments-object/10.6-12-2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-13-a-1.js

- [ ] `language/arguments-object/10.6-13-a-1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-13-a-2.js

- [ ] `language/arguments-object/10.6-13-a-2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-13-a-3.js

- [ ] `language/arguments-object/10.6-13-a-3.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-13-c-1-s.js

- [ ] `language/arguments-object/10.6-13-c-1-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-13-c-2-s.js

- [ ] `language/arguments-object/10.6-13-c-2-s.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.6-13-c-2-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-13-c-3-s.js

- [ ] `language/arguments-object/10.6-13-c-3-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-14-c-1-s.js

- [ ] `language/arguments-object/10.6-14-c-1-s.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.6-14-c-1-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-14-c-4-s.js

- [ ] `language/arguments-object/10.6-14-c-4-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### language/arguments-object/10.6-2gs.js

- [ ] `language/arguments-object/10.6-2gs.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}

```

### language/arguments-object/10.6-5-1.js

- [ ] `language/arguments-object/10.6-5-1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.6-5-1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-6-1.js

- [ ] `language/arguments-object/10.6-6-1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.6-6-1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-6-2.js

- [ ] `language/arguments-object/10.6-6-2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.6-6-2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-6-3-s.js

- [ ] `language/arguments-object/10.6-6-3-s.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.6-6-3-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-6-3.js

- [ ] `language/arguments-object/10.6-6-3.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-6-4-s.js

- [ ] `language/arguments-object/10.6-6-4-s.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.6-6-4-s.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-6-4.js

- [ ] `language/arguments-object/10.6-6-4.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/10.6-7-1.js

- [ ] `language/arguments-object/10.6-7-1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/10.6-7-1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.1.6_A1_T2.js

- [ ] `language/arguments-object/S10.1.6_A1_T2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `language/arguments-object/S10.1.6_A1_T2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### language/arguments-object/S10.6_A1.js

- [ ] `language/arguments-object/S10.6_A1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/S10.6_A1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.6_A2.js

- [ ] `language/arguments-object/S10.6_A2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/S10.6_A2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.6_A3_T1.js

- [ ] `language/arguments-object/S10.6_A3_T1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/S10.6_A3_T1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.6_A3_T2.js

- [ ] `language/arguments-object/S10.6_A3_T2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/S10.6_A3_T2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.6_A3_T3.js

- [ ] `language/arguments-object/S10.6_A3_T3.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### language/arguments-object/S10.6_A3_T4.js

- [ ] `language/arguments-object/S10.6_A3_T4.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.6_A4.js

- [ ] `language/arguments-object/S10.6_A4.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.6_A5_T1.js

- [ ] `language/arguments-object/S10.6_A5_T1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/S10.6_A5_T1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.6_A5_T2.js

- [ ] `language/arguments-object/S10.6_A5_T2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/S10.6_A5_T2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.6_A5_T3.js

- [ ] `language/arguments-object/S10.6_A5_T3.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

- [ ] `language/arguments-object/S10.6_A5_T3.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### language/arguments-object/S10.6_A5_T4.js

- [ ] `language/arguments-object/S10.6_A5_T4.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/S10.6_A5_T4.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.6_A6.js

- [ ] `language/arguments-object/S10.6_A6.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/S10.6_A6.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/S10.6_A7.js

- [ ] `language/arguments-object/S10.6_A7.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/S10.6_A7.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/arguments-caller.js

- [ ] `language/arguments-object/arguments-caller.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/arguments-caller.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/async-gen-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/async-gen-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/async-gen-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/async-gen-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/async-gen-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/async-gen-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/async-gen-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/async-gen-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/async-gen-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/async-gen-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/async-gen-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/async-gen-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/async-gen-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/async-gen-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/async-gen-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/async-gen-named-func-expr-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/async-gen-named-func-expr-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/async-gen-named-func-expr-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/async-gen-named-func-expr-args-trailing-comma-null.js

- [ ] `language/arguments-object/async-gen-named-func-expr-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/async-gen-named-func-expr-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/async-gen-named-func-expr-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/async-gen-named-func-expr-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/async-gen-named-func-expr-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/async-gen-named-func-expr-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/async-gen-named-func-expr-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/async-gen-named-func-expr-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/async-gen-named-func-expr-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/async-gen-named-func-expr-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/async-gen-named-func-expr-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-func-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-async-gen-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-async-private-gen-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-gen-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-gen-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-gen-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-gen-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-gen-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-gen-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-gen-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-gen-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-gen-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-gen-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-gen-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-gen-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-gen-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-gen-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-gen-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

- [ ] `language/arguments-object/cls-decl-gen-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}

```

### language/arguments-object/cls-decl-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unknown class member \"prototype\""}
{"diagnostic":"Go target: unknown class member \"prototype\""}

```

- [ ] `language/arguments-object/cls-decl-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unknown class member \"prototype\""}
{"diagnostic":"Go target: unknown class member \"prototype\""}

```

### language/arguments-object/cls-decl-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unknown class member \"prototype\""}
{"diagnostic":"Go target: unknown class member \"prototype\""}

```

- [ ] `language/arguments-object/cls-decl-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unknown class member \"prototype\""}
{"diagnostic":"Go target: unknown class member \"prototype\""}

```

### language/arguments-object/cls-decl-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unknown class member \"prototype\""}
{"diagnostic":"Go target: unknown class member \"prototype\""}

```

- [ ] `language/arguments-object/cls-decl-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unknown class member \"prototype\""}
{"diagnostic":"Go target: unknown class member \"prototype\""}

```

### language/arguments-object/cls-decl-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unknown class member \"prototype\""}
{"diagnostic":"Go target: unknown class member \"prototype\""}

```

- [ ] `language/arguments-object/cls-decl-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unknown class member \"prototype\""}
{"diagnostic":"Go target: unknown class member \"prototype\""}

```

### language/arguments-object/cls-decl-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unknown class member \"prototype\""}
{"diagnostic":"Go target: unknown class member \"prototype\""}

```

- [ ] `language/arguments-object/cls-decl-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unknown class member \"prototype\""}
{"diagnostic":"Go target: unknown class member \"prototype\""}

```

### language/arguments-object/cls-decl-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/cls-decl-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/cls-decl-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/cls-decl-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/cls-decl-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/cls-decl-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/cls-decl-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/cls-decl-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/cls-decl-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/cls-decl-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-gen-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator methods are not supported"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-gen-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-private-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: methods require ordinary literal names"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: methods require ordinary literal names"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-private-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: methods require ordinary literal names"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: methods require ordinary literal names"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-private-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: methods require ordinary literal names"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: methods require ordinary literal names"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-private-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: methods require ordinary literal names"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: methods require ordinary literal names"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-private-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: methods require ordinary literal names"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: methods require ordinary literal names"}
{"diagnostic":"Go target: unsupported class member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

- [ ] `language/arguments-object/cls-decl-private-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: static members require ordinary literal names"}
{"diagnostic":"Go target: unsupported static member KindGetAccessor"}

```

### language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-func-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-gen-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-async-private-gen-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-gen-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-gen-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-gen-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-gen-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-gen-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-gen-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-gen-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-gen-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-gen-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-gen-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-gen-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-gen-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-gen-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-gen-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-gen-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-gen-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-private-gen-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/cls-expr-private-gen-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/cls-expr-private-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-private-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-private-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-private-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-private-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-private-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-private-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-private-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-private-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-private-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-private-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-private-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-private-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-private-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-private-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-null.js

- [ ] `language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

- [ ] `language/arguments-object/cls-expr-private-meth-static-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindClassExpression"}

```

### language/arguments-object/func-decl-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/func-decl-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/func-decl-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/func-decl-args-trailing-comma-null.js

- [ ] `language/arguments-object/func-decl-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/func-decl-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/func-decl-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/func-decl-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/func-decl-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/func-decl-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/func-decl-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/func-decl-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/func-decl-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/func-decl-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/func-decl-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/func-expr-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/func-expr-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/func-expr-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/func-expr-args-trailing-comma-null.js

- [ ] `language/arguments-object/func-expr-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/func-expr-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/func-expr-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/func-expr-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/func-expr-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/func-expr-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/func-expr-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/func-expr-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/func-expr-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/func-expr-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/func-expr-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/gen-func-decl-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/gen-func-decl-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-func-decl-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-func-decl-args-trailing-comma-null.js

- [ ] `language/arguments-object/gen-func-decl-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-func-decl-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-func-decl-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/gen-func-decl-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-func-decl-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-func-decl-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/gen-func-decl-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-func-decl-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-func-decl-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/gen-func-decl-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-func-decl-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-func-expr-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/gen-func-expr-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-func-expr-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-func-expr-args-trailing-comma-null.js

- [ ] `language/arguments-object/gen-func-expr-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-func-expr-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-func-expr-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/gen-func-expr-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-func-expr-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-func-expr-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/gen-func-expr-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-func-expr-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-func-expr-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/gen-func-expr-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-func-expr-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/gen-meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/gen-meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/gen-meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/gen-meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/gen-meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/gen-meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

- [ ] `language/arguments-object/gen-meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: generator functions are not supported"}

```

### language/arguments-object/mapped

- [ ] `language/arguments-object/mapped/nonconfigurable-nonwritable-descriptors-set-by-arguments.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Function"}
{"diagnostic":"Go target: unsupported or unresolved identifier Math"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier TypeError"}
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/meth-args-trailing-comma-multiple.js

- [ ] `language/arguments-object/meth-args-trailing-comma-multiple.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/meth-args-trailing-comma-multiple.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/meth-args-trailing-comma-null.js

- [ ] `language/arguments-object/meth-args-trailing-comma-null.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/meth-args-trailing-comma-null.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/meth-args-trailing-comma-single-args.js

- [ ] `language/arguments-object/meth-args-trailing-comma-single-args.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/meth-args-trailing-comma-single-args.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/meth-args-trailing-comma-spread-operator.js

- [ ] `language/arguments-object/meth-args-trailing-comma-spread-operator.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/meth-args-trailing-comma-spread-operator.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/meth-args-trailing-comma-undefined.js

- [ ] `language/arguments-object/meth-args-trailing-comma-undefined.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

- [ ] `language/arguments-object/meth-args-trailing-comma-undefined.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/arguments-object/non-strict-arguments-object-is-immutable.js

- [ ] `language/arguments-object/non-strict-arguments-object-is-immutable.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported expression KindDeleteExpression"}

```

### language/arguments-object/unmapped

- [ ] `language/arguments-object/unmapped/via-strict.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}
{"diagnostic":"Go target: unsupported or unresolved identifier arguments"}

```

### language/asi/S7.9.2_A1_T2.js

- [ ] `language/asi/S7.9.2_A1_T2.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9.2_A1_T2.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9.2_A1_T4.js

- [ ] `language/asi/S7.9.2_A1_T4.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9.2_A1_T4.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9.2_A1_T5.js

- [ ] `language/asi/S7.9.2_A1_T5.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9.2_A1_T5.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9.2_A1_T7.js

- [ ] `language/asi/S7.9.2_A1_T7.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9.2_A1_T7.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A1.js

- [ ] `language/asi/S7.9_A1.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier label2"}

```

- [ ] `language/asi/S7.9_A1.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier label2"}

```

### language/asi/S7.9_A10_T1.js

- [ ] `language/asi/S7.9_A10_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A10_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A10_T10.js

- [ ] `language/asi/S7.9_A10_T10.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A10_T10.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A10_T11.js

- [ ] `language/asi/S7.9_A10_T11.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A10_T11.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A10_T12.js

- [ ] `language/asi/S7.9_A10_T12.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A10_T12.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A10_T3.js

- [ ] `language/asi/S7.9_A10_T3.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A10_T3.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A10_T5.js

- [ ] `language/asi/S7.9_A10_T5.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A10_T5.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A10_T7.js

- [ ] `language/asi/S7.9_A10_T7.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A10_T7.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A10_T9.js

- [ ] `language/asi/S7.9_A10_T9.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A10_T9.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A11_T1.js

- [ ] `language/asi/S7.9_A11_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A11_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A11_T10.js

- [ ] `language/asi/S7.9_A11_T10.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A11_T10.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A11_T11.js

- [ ] `language/asi/S7.9_A11_T11.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A11_T11.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A11_T2.js

- [ ] `language/asi/S7.9_A11_T2.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A11_T2.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A11_T3.js

- [ ] `language/asi/S7.9_A11_T3.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A11_T3.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A11_T5.js

- [ ] `language/asi/S7.9_A11_T5.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A11_T5.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A11_T6.js

- [ ] `language/asi/S7.9_A11_T6.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A11_T6.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A11_T7.js

- [ ] `language/asi/S7.9_A11_T7.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A11_T7.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A11_T9.js

- [ ] `language/asi/S7.9_A11_T9.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A11_T9.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A2.js

- [ ] `language/asi/S7.9_A2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier label2"}

```

- [ ] `language/asi/S7.9_A2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier label2"}

```

### language/asi/S7.9_A3.js

- [ ] `language/asi/S7.9_A3.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A3.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A4.js

- [ ] `language/asi/S7.9_A4.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier $DONOTEVALUATE"}
{"diagnostic":"Go target: unsupported or unresolved identifier "}

```

- [ ] `language/asi/S7.9_A4.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier $DONOTEVALUATE"}
{"diagnostic":"Go target: unsupported or unresolved identifier "}

```

### language/asi/S7.9_A5.2_T1.js

- [ ] `language/asi/S7.9_A5.2_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A5.2_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A5.4_T1.js

- [ ] `language/asi/S7.9_A5.4_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A5.4_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A5.5_T1.js

- [ ] `language/asi/S7.9_A5.5_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A5.5_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A5.5_T2.js

- [ ] `language/asi/S7.9_A5.5_T2.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}

```

- [ ] `language/asi/S7.9_A5.5_T2.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}

```

### language/asi/S7.9_A5.5_T3.js

- [ ] `language/asi/S7.9_A5.5_T3.js` (sloppy): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}

```

- [ ] `language/asi/S7.9_A5.5_T3.js` (strict): **emit-error**

```text
{"diagnostic":"Go target: unsupported or unresolved identifier Object"}

```

### language/asi/S7.9_A5.5_T4.js

- [ ] `language/asi/S7.9_A5.5_T4.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A5.5_T4.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A5.5_T5.js

- [ ] `language/asi/S7.9_A5.5_T5.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A5.5_T5.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A5.6_T1.js

- [ ] `language/asi/S7.9_A5.6_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A5.6_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A5.6_T2.js

- [ ] `language/asi/S7.9_A5.6_T2.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A5.6_T2.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A5.8_T1.js

- [ ] `language/asi/S7.9_A5.8_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A5.8_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A5.9_T1.js

- [ ] `language/asi/S7.9_A5.9_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A5.9_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T1.js

- [ ] `language/asi/S7.9_A6.1_T1.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T1.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T10.js

- [ ] `language/asi/S7.9_A6.1_T10.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T10.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T11.js

- [ ] `language/asi/S7.9_A6.1_T11.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T11.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T12.js

- [ ] `language/asi/S7.9_A6.1_T12.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T12.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T13.js

- [ ] `language/asi/S7.9_A6.1_T13.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T13.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T2.js

- [ ] `language/asi/S7.9_A6.1_T2.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T2.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T3.js

- [ ] `language/asi/S7.9_A6.1_T3.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T3.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T4.js

- [ ] `language/asi/S7.9_A6.1_T4.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T4.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T5.js

- [ ] `language/asi/S7.9_A6.1_T5.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T5.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T6.js

- [ ] `language/asi/S7.9_A6.1_T6.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T6.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T7.js

- [ ] `language/asi/S7.9_A6.1_T7.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T7.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T8.js

- [ ] `language/asi/S7.9_A6.1_T8.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T8.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

### language/asi/S7.9_A6.1_T9.js

- [ ] `language/asi/S7.9_A6.1_T9.js` (sloppy): **runtime-fail**

```text
Async runtime errors: Expected an array

```

- [ ] `language/asi/S7.9_A6.1_T9.js` (strict): **runtime-fail**

```text
Async runtime errors: Expected an array

```

