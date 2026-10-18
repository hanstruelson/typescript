//go:build !no_node_modules

package goemit

// nodeConsoleRuntime implements Bun's node:console default export using the
// host console behavior available in the generated runtime.
const nodeConsoleRuntime = `
var tsNodeConsoleSingleton struct {
	sync.Once
	value *tsObject
	counts map[string]uint64
	timers map[string]time.Time
}

func tsNodeConsoleLabel(args []tsValue) string {
	if len(args) == 0 || tsIsUndefined(args[0]) { return "default" }
	return tsStringValue(args[0]).String()
}

func tsNodeConsolePrint(writer io.Writer, args ...tsValue) {
	for index, value := range args {
		if index != 0 { fmt.Fprint(writer, " ") }
		fmt.Fprint(writer, tsText(value))
		if value.kind == tsBigIntKind { fmt.Fprint(writer, "n") }
	}
	fmt.Fprintln(writer)
}

func tsNodeConsoleModule(loop *tsLoop) tsValue {
	tsNodeConsoleSingleton.Do(func() {
		out := tsNewObject()
		tsNodeConsoleSingleton.counts = make(map[string]uint64)
		tsNodeConsoleSingleton.timers = make(map[string]time.Time)
		write := func(name string, writer io.Writer) tsValue {
			return tsNamedNativeMethod(name, 0, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
				tsNodeConsolePrint(writer, args...)
				return tsU
			})
		}
		for _, name := range []string{"log", "info", "debug", "dir", "dirxml", "table", "trace"} { out.set(name, write(name, os.Stdout)) }
		out.set("warn", write("warn", os.Stderr))
		out.set("error", write("error", os.Stderr))
		out.set("assert", tsNamedNativeMethod("assert", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
			if len(args) == 0 || tsTruthy(args[0]) { return tsU }
			message := []tsValue{tsStringReference(tsStringUTF8("Assertion failed"))}
			if len(args) > 1 { message = append(message, args[1:]...) }
			tsNodeConsolePrint(os.Stderr, message...)
			return tsU
		}))
		out.set("count", tsNamedNativeMethod("count", 0, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
			label := tsNodeConsoleLabel(args)
			tsNodeConsoleSingleton.counts[label]++
			tsNodeConsolePrint(os.Stdout, tsStringReference(tsStringUTF8(label+": "+strconv.FormatUint(tsNodeConsoleSingleton.counts[label], 10))))
			return tsU
		}))
		out.set("countReset", tsNamedNativeMethod("countReset", 0, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
			delete(tsNodeConsoleSingleton.counts, tsNodeConsoleLabel(args))
			return tsU
		}))
		out.set("time", tsNamedNativeMethod("time", 0, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
			tsNodeConsoleSingleton.timers[tsNodeConsoleLabel(args)] = time.Now()
			return tsU
		}))
		for _, method := range []string{"timeLog", "timeEnd"} {
			name := method
			out.set(name, tsNamedNativeMethod(name, 0, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
				label := tsNodeConsoleLabel(args)
				started, ok := tsNodeConsoleSingleton.timers[label]
				if !ok { return tsU }
				message := []tsValue{tsStringReference(tsStringUTF8(label+": "+strconv.FormatFloat(float64(time.Since(started))/float64(time.Millisecond), 'f', 3, 64)+"ms"))}
				if len(args) > 1 { message = append(message, args[1:]...) }
				tsNodeConsolePrint(os.Stdout, message...)
				if name == "timeEnd" { delete(tsNodeConsoleSingleton.timers, label) }
				return tsU
			}))
		}
		for _, name := range []string{"group", "groupCollapsed"} { out.set(name, write(name, os.Stdout)) }
		out.set("groupEnd", tsNamedNativeMethod("groupEnd", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsU }))
		out.set("clear", tsNamedNativeMethod("clear", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsU }))
		tsNodeConsoleSingleton.value = out
		if globalObject := tsGlobalNamespace.object; globalObject != nil {
			if globalObject.get("console").kind == tsUndefinedKind {
				globalObject.set("console", tsObjectValue(out))
			} else if existing := globalObject.get("console"); existing.kind == tsObjectKind {
				tsNodeConsoleSingleton.value = (*tsObject)(existing.ref)
			}
		}
	})
	return tsObjectValue(tsNodeConsoleSingleton.value)
}
`
