package goemit

// internalClusterIsPrimaryRuntime ports the tiny early cluster predicate.
// It inspects process.env without loading node:cluster, matching Bun's module.
const internalClusterIsPrimaryRuntime = `
func tsInternalClusterIsPrimaryModule(loop *tsLoop) tsValue {
	process := tsNodeProcessModule(loop)
	env := tsGet(loop, process, tsStringReference(tsStringUTF8("env")))
	isPrimary := true
	if env.kind == tsObjectKind {
		_, hasUniqueID := (*tsObject)(env.ref).values["NODE_UNIQUE_ID"]
		isPrimary = !hasUniqueID
	}
	module := tsNewObject()
	module.set("isPrimary", tsBooleanValue(isPrimary))
	module.set("default", tsObjectValue(module))
	return tsObjectValue(module)
}
`
