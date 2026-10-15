package goemit

// internalWorkerMessagingRuntime supplies the Node worker-message routing
// helpers over the local worker_threads adapter's MessagePort implementation.
const internalWorkerMessagingRuntime = `
var tsInternalWorkerMessagingSingleton struct { sync.Once; module tsValue }

func tsInternalWorkerMessagingModule(loop *tsLoop) tsValue {
	tsInternalWorkerMessagingSingleton.Do(func() {
		module := tsNewObject()
		module.set("initThreadInfo", tsNamedNativeMethod("initThreadInfo", 3, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
			tsNodeWorkerThreads.threadID = int(tsNumber(tsArg(args, 0)))
			tsNodeWorkerThreads.isMainThread = tsTruthy(tsArg(args, 1))
			return tsU
		}))
		module.set("createMessagingChannel", tsNamedNativeMethod("createMessagingChannel", 0, func(loop *tsLoop, _ tsValue, _ ...tsValue) tsValue {
			channel := tsNodeWorkerMakeChannel(loop)
			return tsNodeWorkerObject(loop, map[string]tsValue{"portToMain": tsGet(loop, channel, tsNodeWorkerKey("port1")), "portToWorker": tsGet(loop, channel, tsNodeWorkerKey("port2"))})
		}))
		module.set("registerMainThreadPort", tsNamedNativeMethod("registerMainThreadPort", 2, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsU }))
		module.set("destroyMainThreadPort", tsNamedNativeMethod("destroyMainThreadPort", 1, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsU }))
		module.set("setupMainThreadPort", tsNamedNativeMethod("setupMainThreadPort", 2, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsU }))
		module.set("postMessageToThread", tsNamedNativeMethod("postMessageToThread", 4, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeWorkerPostMessageToThread(loop, args) }))
		module.set("default", tsObjectValue(module))
		tsInternalWorkerMessagingSingleton.module = tsObjectValue(module)
	})
	return tsInternalWorkerMessagingSingleton.module
}
`
