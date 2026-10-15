package goemit

// internalLinkedListRuntime ports Node/Bun's circular intrusive linked-list
// primitives. List and item records stay ordinary generated-runtime objects.
const internalLinkedListRuntime = `
func tsInternalLinkedListModule(loop *tsLoop) tsValue {
	module := tsNewObject()
	module.set("init", tsNamedNativeMethod("init", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		list := tsArg(args, 0)
		tsSet(loop, list, tsStringReference(tsStringUTF8("_idleNext")), list)
		tsSet(loop, list, tsStringReference(tsStringUTF8("_idlePrev")), list)
		return list
	}))
	module.set("peek", tsNamedNativeMethod("peek", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		list := tsArg(args, 0)
		previous := tsGet(loop, list, tsStringReference(tsStringUTF8("_idlePrev")))
		if tsStrictEqual(previous, list) { return tsNull }
		return previous
	}))
	module.set("remove", tsNamedNativeMethod("remove", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		item := tsArg(args, 0)
		next := tsGet(loop, item, tsStringReference(tsStringUTF8("_idleNext")))
		previous := tsGet(loop, item, tsStringReference(tsStringUTF8("_idlePrev")))
		if tsTruthy(next) { tsSet(loop, next, tsStringReference(tsStringUTF8("_idlePrev")), previous) }
		if tsTruthy(previous) { tsSet(loop, previous, tsStringReference(tsStringUTF8("_idleNext")), next) }
		tsSet(loop, item, tsStringReference(tsStringUTF8("_idleNext")), tsNull)
		tsSet(loop, item, tsStringReference(tsStringUTF8("_idlePrev")), tsNull)
		return tsU
	}))
	module.set("append", tsNamedNativeMethod("append", 2, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		list, item := tsArg(args, 0), tsArg(args, 1)
		nextKey, previousKey := tsStringReference(tsStringUTF8("_idleNext")), tsStringReference(tsStringUTF8("_idlePrev"))
		attachedNext, attachedPrevious := tsGet(loop, item, nextKey), tsGet(loop, item, previousKey)
		if tsTruthy(attachedNext) || tsTruthy(attachedPrevious) {
			if tsTruthy(attachedNext) { tsSet(loop, attachedNext, previousKey, attachedPrevious) }
			if tsTruthy(attachedPrevious) { tsSet(loop, attachedPrevious, nextKey, attachedNext) }
			tsSet(loop, item, nextKey, tsNull); tsSet(loop, item, previousKey, tsNull)
		}
		next := tsGet(loop, list, nextKey)
		tsSet(loop, item, nextKey, next)
		tsSet(loop, item, previousKey, list)
		tsSet(loop, next, previousKey, item)
		tsSet(loop, list, nextKey, item)
		return tsU
	}))
	module.set("isEmpty", tsNamedNativeMethod("isEmpty", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
		list := tsArg(args, 0)
		return tsBooleanValue(tsStrictEqual(tsGet(loop, list, tsStringReference(tsStringUTF8("_idleNext"))), list))
	}))
	module.set("default", tsObjectValue(module))
	return tsObjectValue(module)
}
`
