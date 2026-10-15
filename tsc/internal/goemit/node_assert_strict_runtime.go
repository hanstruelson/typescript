package goemit

// nodeAssertStrictRuntime is the one-line node:assert/strict facade from Bun.
const nodeAssertStrictRuntime = `
func tsNodeAssertStrictModule() tsValue { tsNodeAssertModule();return tsNodeAssertSingleton.strict }
`
