package goemit

// internalErrorsRuntime adapts Node's shared stream-error aggregation helper.
const internalErrorsRuntime = `
func tsInternalErrorsAggregateTwoErrors(loop *tsLoop, innerError, outerError tsValue) tsValue {
    if tsTruthy(innerError) && tsTruthy(outerError) && !tsStrictEqual(innerError, outerError) {
        outerErrors := tsGet(loop, outerError, tsStringReference(tsStringUTF8("errors")))
        if outerErrors.kind == tsArrayKind {
            outerErrorsArray := (*tsArray)(outerErrors.ref)
            outerErrorsArray.values = append(outerErrorsArray.values, innerError)
            return outerError
        }
        message := tsGet(loop, outerError, tsStringReference(tsStringUTF8("message")))
        aggregate := tsNewError(loop, "AggregateError", []tsValue{tsArrayValue(&tsArray{values: []tsValue{outerError, innerError}}), message})
        if code := tsGet(loop, outerError, tsStringReference(tsStringUTF8("code"))); !tsIsUndefined(code) {
            (*tsObject)(aggregate.ref).set("code", code)
        }
        return aggregate
    }
    if tsTruthy(innerError) {
        return innerError
    }
    return outerError
}

func tsInternalErrorsModule(loop *tsLoop) tsValue {
    module := tsNewObject()
    aggregate := tsNamedNativeMethod("aggregateTwoErrors", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
        return tsInternalErrorsAggregateTwoErrors(loop, tsArg(args, 0), tsArg(args, 1))
    })
    module.set("aggregateTwoErrors", aggregate)
    module.set("default", tsObjectValue(module))
    return tsObjectValue(module)
}
`
