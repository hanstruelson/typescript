package goemit

// internalReplNodeErrorsRuntime supplies the Node error constructors used by
// the vendored readline/repl modules.
const internalReplNodeErrorsRuntime = `
var tsInternalReplNodeErrors struct { sync.Once; module tsValue }

func tsInternalReplNodeErrorMessage(loop *tsLoop, code string, args []tsValue) string {
	arg := func(index int) tsValue { return tsArg(args, index) }
	text := func(index int) string { return tsText(arg(index)) }
	switch code {
	case "ERR_INVALID_ARG_TYPE":
		name, expected := text(0), text(1)
		if expected == "" { expected = "a valid value" }
		return "The \"" + name + "\" argument must be " + expected
	case "ERR_INVALID_ARG_VALUE":
		name := text(0)
		message := "The \"" + name + "\" argument is invalid"
		if len(args) > 1 && !tsIsUndefined(arg(1)) { message += ". Received " + tsNodeErrorValueString(loop, arg(1)) }
		if len(args) > 2 && tsArg(args, 2).kind == tsStringKind { message += ". " + text(2) }
		return message
	case "ERR_MISSING_ARGS":
		names := make([]string, 0, len(args))
		for _, value := range args { names = append(names, "\""+tsText(value)+"\"") }
		if len(names) == 0 { return "The required arguments are missing" }
		if len(names) == 1 { return "The " + names[0] + " argument must be specified" }
		return "The " + strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1] + " arguments must be specified"
	case "ERR_USE_AFTER_CLOSE":
		return "Cannot use " + text(0) + " after close"
	case "ERR_INVALID_CURSOR_POS":
		return "Cannot set cursor position to an invalid position"
	case "ERR_SCRIPT_EXECUTION_INTERRUPTED":
		return "Script execution was interrupted by SIGINT"
	case "ERR_INVALID_STATE":
		message := "Invalid state"
		if !tsIsUndefined(arg(0)) { message += ": " + text(0) }
		return message
	case "ERR_CANNOT_WATCH_SIGINT":
		return "Cannot watch for interruptions when running asynchronously"
	case "ERR_INSPECTOR_NOT_AVAILABLE":
		return "Inspector is not available"
	case "ERR_INVALID_REPL_EVAL_CONFIG":
		return "Cannot specify both \"breakEvalOnSigint\" and \"eval\" for REPL"
	case "ERR_INVALID_REPL_INPUT":
		return text(0)
	}
	return code
}

func tsNodeErrorValueString(loop *tsLoop, value tsValue) string {
	if value.kind == tsStringKind { return "'" + tsText(value) + "'" }
	if value.kind == tsUndefinedKind { return "undefined" }
	if value.kind == tsNullKind { return "null" }
	if value.kind == tsBooleanKind { return tsText(value) }
	if value.kind == tsNumberKind { return tsText(value) }
	return "an instance of " + tsText(tsGet(loop, value, tsStringReference(tsStringUTF8("constructor"))))
}

func tsInternalReplNodeErrorConstructor(loop *tsLoop, name, code string, abort bool) tsValue {
	makeError := func(loop *tsLoop, args []tsValue) tsValue {
		message := tsInternalReplNodeErrorMessage(loop, code, args)
		var errorValue tsValue
		errorName := "Error"
		switch code {
		case "ERR_INVALID_ARG_TYPE", "ERR_INVALID_ARG_VALUE", "ERR_MISSING_ARGS", "ERR_USE_AFTER_CLOSE", "ERR_INVALID_CURSOR_POS":
			errorName = "TypeError"
		}
		if abort {
			if len(args) > 1 && !tsIsPrimitiveValue(args[1]) && tsHas(args[1], tsStringReference(tsStringUTF8("cause"))) {
				errorValue = tsFSException("AbortError", "ABORT_ERR", message)
				(*tsRuntimeError)(errorValue.ref).fields["cause"] = tsGet(loop, args[1], tsStringReference(tsStringUTF8("cause")))
			} else {
				errorValue = tsFSException("AbortError", "ABORT_ERR", message)
			}
		} else {
			errorValue = tsFSException(errorName, code, message)
		}
		return errorValue
	}
	function := tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return makeError(loop, args) })
	function.name = name
	function.length = 1
	function.constructible = true
	function.construct = func(loop *tsLoop, args ...tsValue) tsValue { return makeError(loop, args) }
	function.receiverCall = func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return makeError(loop, args) }
	function.properties = tsNewObject()
	errorName := "Error"
	switch code { case "ERR_INVALID_ARG_TYPE", "ERR_INVALID_ARG_VALUE", "ERR_MISSING_ARGS", "ERR_USE_AFTER_CLOSE", "ERR_INVALID_CURSOR_POS": errorName = "TypeError" }
	function.properties.set("prototype", tsErrorPrototype(errorName))
	instanceKey := tsPropertyKey(tsWellKnownSymbol("hasInstance"))
	function.properties.set(instanceKey, tsNamedNativeMethod("[Symbol.hasInstance]", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		value := tsArg(args, 0)
		return tsBooleanValue((value.kind == tsErrorKind || value.kind == tsObjectKind) && tsStrictEqual(tsGet(loop, value, tsStringReference(tsStringUTF8("code"))), tsStringReference(tsStringUTF8(code))))
	}))
	return tsFunctionValue(function)
}

func tsInternalReplNodePrimordialsIsErrorStackTraceLimitWritable(loop *tsLoop) tsValue {
	globalError := tsInternalReplPrimordialsGlobal(loop, "Error")
	object := tsInternalReplPrimordialsGlobal(loop, "Object")
	descriptorFn := tsGet(loop, object, tsStringReference(tsStringUTF8("getOwnPropertyDescriptor")))
	descriptor := tsCallReceiver(loop, descriptorFn, object, globalError, tsStringReference(tsStringUTF8("stackTraceLimit")))
	if tsIsUndefined(descriptor) {
		extensible := tsGet(loop, object, tsStringReference(tsStringUTF8("isExtensible")))
		return tsCallReceiver(loop, extensible, object, globalError)
	}
	if tsHas(descriptor, tsStringReference(tsStringUTF8("writable"))) { return tsGet(loop, descriptor, tsStringReference(tsStringUTF8("writable"))) }
	return tsBooleanValue(!tsIsUndefined(tsGet(loop, descriptor, tsStringReference(tsStringUTF8("set")))))
}

func tsInternalReplNodeErrorsModule(loop *tsLoop) tsValue {
	tsInternalReplNodeErrors.Do(func() {
		module, codes := tsNewObject(), tsNewObject()
		codeNames := []string{"ERR_INVALID_ARG_TYPE", "ERR_INVALID_ARG_VALUE", "ERR_MISSING_ARGS", "ERR_USE_AFTER_CLOSE", "ERR_INVALID_CURSOR_POS", "ERR_SCRIPT_EXECUTION_INTERRUPTED", "ERR_INVALID_STATE", "ERR_CANNOT_WATCH_SIGINT", "ERR_INSPECTOR_NOT_AVAILABLE", "ERR_INVALID_REPL_EVAL_CONFIG", "ERR_INVALID_REPL_INPUT"}
		for _, name := range codeNames {
			constructor := tsInternalReplNodeErrorConstructor(loop, name, name, false)
			codes.set(name, constructor)
		}
		module.set("AbortError", tsInternalReplNodeErrorConstructor(loop, "AbortError", "ABORT_ERR", true))
		module.set("isErrorStackTraceLimitWritable", tsNamedNativeMethod("isErrorStackTraceLimitWritable", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsInternalReplNodePrimordialsIsErrorStackTraceLimitWritable(loop) }))
		module.set("codes", tsObjectValue(codes))
		module.set("default", tsObjectValue(module))
		tsInternalReplNodeErrors.module = tsObjectValue(module)
	})
	return tsInternalReplNodeErrors.module
}
`
