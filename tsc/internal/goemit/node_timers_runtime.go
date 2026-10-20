package goemit

// nodeTimersRuntime and nodeTimersPromisesRuntime model Bun's timer wrappers
// on the compiler runtime's event loop rather than introducing a host library.
const nodeTimerCoreRuntime = `
func tsNodeTimerSetInterval(loop *tsLoop, callback, milliseconds tsValue, args ...tsValue) *tsTask {
	if callback.kind != tsFunctionKind { tsPropertyFailure("setInterval expects a function") }
	handle := &tsTask{timer:true}
	delay := milliseconds
	if tsIsUndefined(delay) { delay = tsNumberValue(1) }
	var active *tsTask
	var schedule func()
	invoke := tsFunctionValue(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
		if handle.canceled { return tsU }
		tsCall(loop, callback, args...)
		if !handle.canceled { schedule() }
		return tsU
	}))
	schedule = func() {
		if handle.canceled { return }
		active = loop.setTimeout(invoke, delay)
		if handle.unref { active.unref = true }
	}
	handle.cancel = func() {
		if active != nil { loop.clearTimeout(tsTaskValue(active)) }
	}
	handle.refHook = func(referenced bool) {
		if active != nil { active.unref = !referenced }
	}
	schedule()
	return handle
}

func tsNodeTimerClearInterval(loop *tsLoop, handle tsValue) tsValue {
	return loop.clearTimeout(handle)
}
`

const nodeTimersRuntime = nodeTimerCoreRuntime + `
type tsNodeTimersState struct { sync.Once; value *tsObject; promises tsValue }
var tsNodeTimersSingleton tsNodeTimersState

func tsNodeTimersModule() tsValue {
	tsNodeTimersSingleton.Do(func() {
		out := tsNewObject()
		out.set("setTimeout", tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
			callback := tsArg(args, 0)
			extra := []tsValue{}
			if len(args) > 2 { extra = args[2:] }
			return tsTaskValue(loop.setTimeout(callback, tsArg(args, 1), extra...))
		}))
		out.set("clearTimeout", tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return loop.clearTimeout(tsArg(args, 0)) }))
		out.set("setInterval", tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
			callback := tsArg(args, 0)
			extra := []tsValue{}
			if len(args) > 2 { extra = args[2:] }
			return tsTaskValue(tsNodeTimerSetInterval(loop, callback, tsArg(args, 1), extra...))
		}))
		out.set("clearInterval", tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return tsNodeTimerClearInterval(loop, tsArg(args, 0)) }))
		out.set("setImmediate", tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
			callback := tsArg(args, 0)
			extra := []tsValue{}
			if len(args) > 1 { extra = args[1:] }
			return tsTaskValue(loop.setTimeout(callback, tsNumberValue(0), extra...))
		}))
		out.set("clearImmediate", tsFunc(func(loop *tsLoop, args ...tsValue) tsValue { return loop.clearTimeout(tsArg(args, 0)) }))
		out.set("active", tsNamedNativeMethod("active", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			handle := tsArg(args, 0)
			if handle.kind == tsTaskKind { task := (*tsTask)(handle.ref); if task.cancel == nil && !task.completed { task.canceled = false } }
			return tsU
		}))
		out.set("unenroll", tsNamedNativeMethod("unenroll", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return loop.clearTimeout(tsArg(args, 0)) }))
		out.set("enroll", tsNamedNativeMethod("enroll", 2, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue {
			panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"Error",message:"timers.enroll is not implemented for native handles"})})
		}))
		out.order = append(out.order, "promises")
		out.descriptors = map[string]*tsDescriptor{"promises": {
			accessor:true, enumerable:true, configurable:true,
			getter:tsNamedNativeMethod("get promises", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
				if tsIsUndefined(tsNodeTimersSingleton.promises) { tsNodeTimersSingleton.promises = tsNodeTimersPromisesModule() }
				return tsNodeTimersSingleton.promises
			}),
			setter:tsNamedNativeMethod("set promises", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
				tsNodeTimersSingleton.promises = tsArg(args, 0)
				return tsU
			}),
		}}
		tsNodeTimersSingleton.value = out
	})
	return tsObjectValue(tsNodeTimersSingleton.value)
}
`

const nodeTimersPromisesRuntime = `
func tsNodeTimerRejectedIterator(failure tsValue) tsValue {
	iterator := &tsAsyncIterator{
		next: func(loop *tsLoop) *tsPromise { return loop.resolved(failure, true) },
		close: func(loop *tsLoop) *tsPromise { return loop.resolved(tsIteratorResult(tsU, true), false) },
	}
	return tsAsyncIteratorObject(iterator)
}

func tsNodeTimerValidateOptions(options tsValue) (*tsAbortState, bool) {
	if tsIsUndefined(options) { return nil, true }
	if options.kind != tsObjectKind && options.kind != tsInstanceKind { tsFSInvalid("options", "an object") }
	signal := tsFSAbortSignal(options)
	if !tsIsUndefined(tsFSProperty(options, "ref")) && tsFSProperty(options, "ref").kind != tsBooleanKind { tsFSInvalid("options.ref", "a boolean") }
	return signal, tsIsUndefined(tsFSProperty(options, "ref")) || tsTruthy(tsFSProperty(options, "ref"))
}

func tsNodeTimerPromise(loop *tsLoop, milliseconds, value, options tsValue, immediate bool) tsValue {
	promise := loop.promise()
	func() {
		defer func() { if failure := recover(); failure != nil { promise.settle(tsResult{tsUnwrap(failure), true}) } }()
		if tsIsUndefined(milliseconds) { milliseconds = tsNumberValue(1) }
		if !tsIsNumeric(milliseconds) { tsFSInvalid("delay", "a number") }
		signal, reference := tsNodeTimerValidateOptions(options)
		if signal != nil { if _, aborted := tsAbortReason(signal); aborted { promise.settle(tsResult{tsFSAbortError(signal), true}); return } }
		delay := milliseconds
		if immediate { delay = tsNumberValue(0) }
		var task *tsTask
		task = loop.setTimeout(tsFunc(func(_ *tsLoop, _ ...tsValue) tsValue { promise.settle(tsResult{value, false}); return tsU }), delay)
		if !reference { task.unref = true }
		if signal != nil {
			signal.followers = append(signal.followers, func(_ tsValue) {
				loop.clearTimeout(tsTaskValue(task))
				promise.settle(tsResult{tsFSAbortError(signal), true})
			})
		}
	}()
	return tsPromiseValue(promise)
}

func tsNodeTimerInterval(loop *tsLoop, milliseconds, value, options tsValue) tsValue {
	var validationFailure tsValue
	var signal *tsAbortState
	var reference bool
	func() {
		defer func() { if failure := recover(); failure != nil { validationFailure = tsUnwrap(failure) } }()
		if tsIsUndefined(milliseconds) { milliseconds = tsNumberValue(1) }
		if !tsIsNumeric(milliseconds) { tsFSInvalid("delay", "a number") }
		signal, reference = tsNodeTimerValidateOptions(options)
	}()
	if !tsIsUndefined(validationFailure) { return tsNodeTimerRejectedIterator(validationFailure) }
	if signal != nil { if _, aborted := tsAbortReason(signal); aborted { return tsNodeTimerRejectedIterator(tsFSAbortError(signal)) } }

	var active *tsTask
	var waiters []*tsPromise
	var queued uint64
	closed, aborted := false, false
	var iterator *tsAsyncIterator
	var schedule func()
	settleNext := func() {
		for len(waiters) > 0 && queued > 0 {
			waiter := waiters[0]; waiters = waiters[1:]; queued--
			waiter.settle(tsResult{tsIteratorResult(value, false), false})
		}
		if len(waiters) > 0 && aborted && queued == 0 {
			failure := tsFSAbortError(signal)
			for _, waiter := range waiters { waiter.settle(tsResult{failure, true}) }
			waiters = nil
		}
	}
	schedule = func() {
		if closed || aborted { return }
		active = loop.setTimeout(tsFunc(func(loop *tsLoop, _ ...tsValue) tsValue {
			if closed || aborted { return tsU }
			queued++
			settleNext()
			schedule()
			return tsU
		}), milliseconds)
		if !reference { active.unref = true }
	}
	iterator = &tsAsyncIterator{
		next: func(loop *tsLoop) *tsPromise {
			if closed { return loop.resolved(tsIteratorResult(tsU, true), false) }
			if queued > 0 { queued--; return loop.resolved(tsIteratorResult(value, false), false) }
			if aborted { return loop.resolved(tsFSAbortError(signal), true) }
			promise := loop.promise(); waiters = append(waiters, promise); return promise
		},
		close: func(loop *tsLoop) *tsPromise {
			if !closed {
				closed = true
				if active != nil { loop.clearTimeout(tsTaskValue(active)) }
				for _, waiter := range waiters { waiter.settle(tsResult{tsIteratorResult(tsU, true), false}) }
				waiters = nil
			}
			return loop.resolved(tsIteratorResult(tsU, true), false)
		},
	}
	_ = iterator
	schedule()
	if signal != nil {
		signal.followers = append(signal.followers, func(_ tsValue) {
			aborted = true
			if active != nil { loop.clearTimeout(tsTaskValue(active)) }
			settleNext()
		})
	}
	return tsAsyncIteratorObject(iterator)
}

func tsNodeTimersPromisesModule() tsValue {
	out := tsNewObject()
	out.set("setTimeout", tsNamedNativeMethod("setTimeout", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeTimerPromise(loop, tsArg(args, 0), tsArg(args, 1), tsArg(args, 2), false) }))
	out.set("setImmediate", tsNamedNativeMethod("setImmediate", 0, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeTimerPromise(loop, tsNumberValue(0), tsArg(args, 0), tsArg(args, 1), true) }))
	out.set("setInterval", tsNamedNativeMethod("setInterval", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeTimerInterval(loop, tsArg(args, 0), tsArg(args, 1), tsArg(args, 2)) }))
	scheduler := tsNewObject()
	scheduler.set("wait", tsNamedNativeMethod("wait", 1, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeTimerPromise(loop, tsArg(args, 0), tsU, tsArg(args, 1), false) }))
	scheduler.set("yield", tsNamedNativeMethod("yield", 0, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeTimerPromise(loop, tsNumberValue(0), tsU, tsArg(args, 0), true) }))
	out.set("scheduler", tsObjectValue(scheduler))
	return tsObjectValue(out)
}
`
