package goemit

// nodeUtilTypesRuntime translates Bun's native predicate table onto the
// concrete value tags used by the generated runtime.
const nodeUtilTypesRuntime = `
type tsNodeUtilTypesState struct { sync.Once; value *tsObject }
var tsNodeUtilTypesSingleton tsNodeUtilTypesState

func tsNodeUtilBoxedKind(value tsValue) tsKind {
	if value.kind != tsObjectKind { return tsUndefinedKind }
	return (*tsObject)(value.ref).boxedPrimitive.kind
}

func tsNodeUtilType(name string, value tsValue) bool {
	boxed := tsNodeUtilBoxedKind(value)
	object := (*tsObject)(nil)
	if value.kind == tsObjectKind { object = (*tsObject)(value.ref) }
	switch name {
	case "isExternal", "isWeakMap", "isWeakSet", "isSharedArrayBuffer", "isCryptoKey":
		return false
	case "isKeyObject":
		return tsNodeCryptoKeyFor(value) != nil
	case "isDate":
		return object != nil && object.nativeDate
	case "isArgumentsObject":
		return object != nil && object.argumentMap != nil
	case "isBigIntObject":
		return boxed == tsBigIntKind
	case "isBooleanObject":
		return boxed == tsBooleanKind
	case "isNumberObject":
		return tsIsNumeric(tsValue{kind:boxed})
	case "isStringObject":
		return boxed == tsStringKind
	case "isSymbolObject":
		return boxed == tsSymbolKind
	case "isNativeError":
		return value.kind == tsErrorKind || object != nil && object.nativeError
	case "isRegExp":
		return value.kind == tsRegExpKind
	case "isAsyncFunction":
		return value.kind == tsFunctionKind && ((*tsFunction)(value.ref).async || (*tsFunction)(value.ref).asyncGeneratorFunction)
	case "isGeneratorFunction":
		return value.kind == tsFunctionKind && (*tsFunction)(value.ref).generator
	case "isGeneratorObject":
		if value.kind != tsIteratorKind { return false }
		iterator := (*tsIterator)(value.ref)
		return iterator.machine != nil && (iterator.machine.generator || iterator.machine.asyncGenerator != nil)
	case "isPromise":
		return value.kind == tsPromiseKind
	case "isMap":
		return value.kind == tsCollectionKind && !(*tsCollection)(value.ref).set
	case "isSet":
		return value.kind == tsCollectionKind && (*tsCollection)(value.ref).set
	case "isMapIterator":
		return value.kind == tsIteratorKind && (*tsIterator)(value.ref).collectionIterator && !(*tsIterator)(value.ref).collectionSet
	case "isSetIterator":
		return value.kind == tsIteratorKind && (*tsIterator)(value.ref).collectionIterator && (*tsIterator)(value.ref).collectionSet
	case "isArrayBuffer", "isAnyArrayBuffer":
		return value.kind == tsArrayBufferKind
	case "isDataView":
		return value.kind == tsDataViewKind
	case "isProxy":
		return object != nil && object.proxy
	case "isModuleNamespaceObject":
		return value.kind == tsNamespaceKind
	case "isBoxedPrimitive":
		return boxed == tsBooleanKind || tsIsNumeric(tsValue{kind:boxed}) || boxed == tsStringKind || boxed == tsSymbolKind || boxed == tsBigIntKind
	case "isArrayBufferView", "isTypedArray":
		return value.kind == tsTypedArrayKind || name == "isArrayBufferView" && value.kind == tsDataViewKind
	case "isUint8Array", "isUint8ClampedArray", "isUint16Array", "isUint32Array", "isInt8Array", "isInt16Array", "isInt32Array", "isFloat16Array", "isFloat32Array", "isFloat64Array", "isBigInt64Array", "isBigUint64Array":
		if value.kind != tsTypedArrayKind { return false }
		return (*tsTypedArray)(value.ref).name == name[2:]
	case "isEventTarget":
		return object != nil && (object.abort != nil || object.nativeClass != nil && object.nativeClass.builtin == "EventTarget")
	default:
		return false
	}
}

func tsNodeUtilTypesModule() tsValue {
	tsNodeUtilTypesSingleton.Do(func() {
		out := tsNewObject()
		for _, name := range []string{
			"isExternal", "isDate", "isArgumentsObject", "isBigIntObject", "isBooleanObject", "isNumberObject", "isStringObject", "isSymbolObject",
			"isNativeError", "isRegExp", "isAsyncFunction", "isGeneratorFunction", "isGeneratorObject", "isPromise", "isMap", "isSet", "isMapIterator", "isSetIterator",
			"isWeakMap", "isWeakSet", "isArrayBuffer", "isDataView", "isSharedArrayBuffer", "isProxy", "isModuleNamespaceObject", "isAnyArrayBuffer", "isBoxedPrimitive", "isArrayBufferView", "isTypedArray",
			"isUint8Array", "isUint8ClampedArray", "isUint16Array", "isUint32Array", "isInt8Array", "isInt16Array", "isInt32Array", "isFloat16Array", "isFloat32Array", "isFloat64Array", "isBigInt64Array", "isBigUint64Array", "isKeyObject", "isCryptoKey", "isEventTarget",
		} {
			methodName := name
			out.set(methodName, tsNamedNativeMethod(methodName, 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
				return tsBooleanValue(tsNodeUtilType(methodName, tsArg(args, 0)))
			}))
		}
		tsNodeUtilTypesSingleton.value = out
	})
	return tsObjectValue(tsNodeUtilTypesSingleton.value)
}
`
