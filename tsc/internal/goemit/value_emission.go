package goemit

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

// formatValueSource marshals native Go expressions into the value ABI at the
// generated Go's explicit type boundaries. It leaves native arithmetic alone.
// The source builder supplies all local types; no Go toolchain subprocess or
// TypeScript rechecking is needed during this pass.
func formatValueSource(source []byte) ([]byte, error) {
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, "output.go", source, parser.ParseComments)
	if err != nil {
		line := 0
		_, _ = fmt.Sscanf(err.Error(), "output.go:%d:", &line)
		lines := bytes.Split(source, []byte("\n"))
		start, end := line-2, line+2
		if start < 1 {
			start = 1
		}
		if end > len(lines) {
			end = len(lines)
		}
		var excerpt bytes.Buffer
		for index := start; index <= end; index++ {
			fmt.Fprintf(&excerpt, "%d: %s\n", index, lines[index-1])
		}
		return nil, fmt.Errorf("%w\n%s", err, excerpt.String())
	}
	threadLoopContext(file)
	pass := valueEmission{fields: map[string]map[string]string{}, functions: map[string][]valueSignature{}, globals: map[string]string{}}
	ast.Inspect(file, func(node ast.Node) bool {
		structure, ok := node.(*ast.StructType)
		if !ok {
			return true
		}
		fields := map[string]string{}
		for _, field := range structure.Fields.List {
			for _, name := range field.Names {
				fields[name.Name] = valueType(field.Type)
			}
		}
		pass.fields[valueType(structure)] = fields
		return true
	})
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				if typ, ok := spec.(*ast.TypeSpec); ok {
					if structure, ok := typ.Type.(*ast.StructType); ok {
						fields := map[string]string{}
						for _, field := range structure.Fields.List {
							for _, name := range field.Names {
								fields[name.Name] = valueType(field.Type)
							}
						}
						pass.fields[typ.Name.Name] = fields
					}
				} else if value, ok := spec.(*ast.ValueSpec); ok {
					if structure, ok := value.Type.(*ast.StructType); ok {
						fields := map[string]string{}
						for _, field := range structure.Fields.List {
							for _, name := range field.Names {
								fields[name.Name] = valueType(field.Type)
							}
						}
						pass.fields[valueType(value.Type)] = fields
						for _, name := range value.Names {
							pass.fields[name.Name] = fields
						}
					}
				}
			}
		case *ast.FuncDecl:
			pass.functions[d.Name.Name] = append(pass.functions[d.Name.Name], valueFunctionSignature(d.Type))
		}
	}
	for _, decl := range file.Decls {
		if d, ok := decl.(*ast.GenDecl); ok {
			pass.declaration(d, pass.globals)
		}
	}
	for _, decl := range file.Decls {
		if d, ok := decl.(*ast.FuncDecl); ok {
			// Host adapters deliberately accept/return Go interfaces. Their body is not
			// JavaScript code and must not recursively call itself through this pass.
			if d.Name.Name == "tsUnwrap" || d.Name.Name == "tsPrimitiveValue" || d.Name.Name == "tsNative" || d.Name.Name == "tsNodeVMInstallModuleMethods" || d.Name.Name == "tsInternalClusterSharedCreate" {
				continue
			}
			pass.scope = d.Name.Name
			env := copyValueEnvironment(pass.globals)
			pass.parameters(d.Type, env)
			if d.Recv != nil {
				pass.fieldNames(d.Recv, env)
			}
			pass.block(d.Body, env, valueReturnTypes(d.Type))
		}
	}
	if pass.err != nil {
		return nil, pass.err
	}
	var out bytes.Buffer
	err = format.Node(&out, set, file)
	return out.Bytes(), err
}

type valueSignature struct {
	typeParameters []string
	parameters     []string
	variadic       bool
	result         string
	results        []string
}
type valueEmission struct {
	scope       string
	callContext string
	err         error
	fields      map[string]map[string]string
	functions   map[string][]valueSignature
	globals     map[string]string
}

func valueType(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	var out bytes.Buffer
	_ = format.Node(&out, token.NewFileSet(), expr)
	return out.String()
}
func valueFunctionSignature(fn *ast.FuncType) valueSignature {
	sig := valueSignature{}
	if fn.TypeParams != nil {
		for _, field := range fn.TypeParams.List {
			for _, name := range field.Names {
				sig.typeParameters = append(sig.typeParameters, name.Name)
			}
		}
	}
	if fn.Params != nil {
		for _, field := range fn.Params.List {
			kind := valueType(field.Type)
			if _, ok := field.Type.(*ast.Ellipsis); ok {
				sig.variadic = true
				kind = strings.TrimPrefix(kind, "...")
			}
			count := len(field.Names)
			if count == 0 {
				count = 1
			}
			for range count {
				sig.parameters = append(sig.parameters, kind)
			}
		}
	}
	if fn.Results != nil {
		for _, field := range fn.Results.List {
			count := len(field.Names)
			if count == 0 {
				count = 1
			}
			for range count {
				sig.results = append(sig.results, valueType(field.Type))
			}
		}
	}
	if len(sig.results) > 0 {
		sig.result = sig.results[0]
	}
	return sig
}
func copyValueEnvironment(env map[string]string) map[string]string {
	copy := map[string]string{}
	for name, kind := range env {
		copy[name] = kind
	}
	return copy
}
func (p *valueEmission) fieldNames(fields *ast.FieldList, env map[string]string) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		kind := valueType(field.Type)
		kind = strings.ReplaceAll(kind, "...", "[]")
		for _, name := range field.Names {
			env[name.Name] = kind
		}
	}
}
func (p *valueEmission) parameters(fn *ast.FuncType, env map[string]string) {
	p.fieldNames(fn.Params, env)
	p.fieldNames(fn.Results, env)
}
func valueBase(kind string) string {
	kind = strings.TrimPrefix(kind, "*")
	if strings.HasPrefix(kind, "struct{") {
		return kind
	}
	if strings.HasPrefix(kind, "[]") {
		return strings.TrimPrefix(kind, "[]")
	}
	if strings.HasPrefix(kind, "[") {
		if end := strings.IndexByte(kind, ']'); end >= 0 {
			return kind[end+1:]
		}
	}
	return kind
}
func valueElement(kind string) string {
	if strings.HasPrefix(kind, "[") {
		if at := strings.IndexByte(kind, ']'); at >= 0 {
			return kind[at+1:]
		}
	}
	if strings.HasPrefix(kind, "map[") {
		if at := strings.IndexByte(kind, ']'); at >= 0 {
			return kind[at+1:]
		}
	}
	return ""
}
func valueCallName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return e.Sel.Name
	case *ast.IndexExpr:
		return valueCallName(e.X)
	case *ast.IndexListExpr:
		return valueCallName(e.X)
	}
	return ""
}
func (p *valueEmission) signature(call *ast.CallExpr, env map[string]string) valueSignature {
	// Callback signatures live on locals and struct fields, including the
	// concrete class views. Resolve these before same-named runtime methods.
	if typ := p.kind(call.Fun, env); strings.HasPrefix(typ, "func(") {
		if parsed, err := parser.ParseExpr(typ); err == nil {
			if fn, ok := parsed.(*ast.FuncType); ok {
				return valueFunctionSignature(fn)
			}
		}
	}
	// Binding cell methods are emitted as selectors. Their generic receiver
	// types are intentionally not expanded into the package-wide overload table,
	// but the one-value set/init ABI is stable across the cell implementations.
	if selector, ok := call.Fun.(*ast.SelectorExpr); ok && len(call.Args) == 1 &&
		(selector.Sel.Name == "set" || selector.Sel.Name == "init") {
		return valueSignature{parameters: []string{"tsValue"}, result: "tsValue", results: []string{"tsValue"}}
	}
	candidates := p.functions[valueCallName(call.Fun)]
	for _, sig := range candidates {
		if (!sig.variadic && len(sig.parameters) == len(call.Args)) || (sig.variadic && len(call.Args) >= len(sig.parameters)-1) {
			types := []ast.Expr{}
			switch indexed := call.Fun.(type) {
			case *ast.IndexExpr:
				types = append(types, indexed.Index)
			case *ast.IndexListExpr:
				types = indexed.Indices
			}
			if len(types) == len(sig.typeParameters) && len(types) > 0 {
				replace := func(typ string) string {
					for i, name := range sig.typeParameters {
						if typ == name {
							return valueType(types[i])
						}
						if typ == "[]"+name {
							return "[]" + valueType(types[i])
						}
					}
					return typ
				}
				sig.parameters = append([]string(nil), sig.parameters...)
				for i, typ := range sig.parameters {
					sig.parameters[i] = replace(typ)
				}
				sig.results = append([]string(nil), sig.results...)
				for i, typ := range sig.results {
					sig.results[i] = replace(typ)
				}
				if len(sig.results) > 0 {
					sig.result = sig.results[0]
				} else {
					sig.result = replace(sig.result)
				}
			}
			return sig
		}
	}
	return valueSignature{}
}
func (p *valueEmission) kind(expr ast.Expr, env map[string]string) string {
	switch e := expr.(type) {
	case *ast.Ident:
		if kind := env[e.Name]; kind != "" {
			return kind
		}
		switch e.Name {
		case "true", "false":
			return "bool"
		case "nil":
			return "nil"
		}
	case *ast.BasicLit:
		switch e.Kind {
		case token.INT, token.FLOAT:
			return "float64"
		case token.STRING:
			return "string"
		}
	case *ast.CompositeLit:
		return valueType(e.Type)
	case *ast.ParenExpr:
		return p.kind(e.X, env)
	case *ast.UnaryExpr:
		kind := p.kind(e.X, env)
		if e.Op == token.AND {
			return "*" + kind
		}
		if e.Op == token.MUL {
			return strings.TrimPrefix(kind, "*")
		}
		if e.Op == token.NOT {
			return "bool"
		}
		return kind
	case *ast.StarExpr:
		return strings.TrimPrefix(p.kind(e.X, env), "*")
	case *ast.TypeAssertExpr:
		return valueType(e.Type)
	case *ast.IndexExpr:
		return valueElement(p.kind(e.X, env))
	case *ast.SliceExpr:
		return p.kind(e.X, env)
	case *ast.SelectorExpr:
		if field := p.fields[valueBase(p.kind(e.X, env))][e.Sel.Name]; field != "" {
			return field
		}
		if name, ok := e.X.(*ast.Ident); ok {
			return p.fields[name.Name][e.Sel.Name]
		}
	case *ast.FuncLit:
		return valueType(e.Type)
	case *ast.CallExpr:
		// Go conversions such as []tsValue(nil) are represented as calls in the
		// syntax tree. Preserve the target type so later index assignments keep
		// their element type through append/copy initialization patterns.
		switch e.Fun.(type) {
		case *ast.ArrayType, *ast.MapType, *ast.ChanType:
			return valueType(e.Fun)
		}
		if selector, ok := e.Fun.(*ast.SelectorExpr); ok {
			if selector.Sel.Name == "init" {
				return "tsValue"
			}
			if pkg, ok := selector.X.(*ast.Ident); ok {
				switch pkg.Name {
				case "math":
					if selector.Sel.Name == "IsNaN" || selector.Sel.Name == "IsInf" {
						return "bool"
					}
					return "float64"
				case "fmt":
					if selector.Sel.Name == "Sprint" || selector.Sel.Name == "Sprintf" {
						return "string"
					}
				case "strings":
					if selector.Sel.Name == "Join" {
						return "string"
					}
				}
			}
			if selector.Sel.Name == "Error" {
				return "string"
			}
		}
		name := valueCallName(e.Fun)
		if strings.HasPrefix(name, "ts") && strings.HasSuffix(name, "Module") {
			return "tsValue"
		}
		switch name {
		case "tsNativeClass":
			return "*tsClass"
		case "tsInstanceProperties":
			return "*tsProperties"
		case "tsDescriptorObject":
			return "*tsObject"
		case "metadata":
			return "*tsObject"
		case "node":
			return "*tsJSONNode"
		case "tsIterate", "tsArrayLikeIterator":
			return "*tsIterator"
		case "tsObjectKeys", "tsOwnKeys":
			return "*tsArray"
		}
		switch name {
		case "tsValue", "tsNumberValue", "tsBooleanValue", "tsStringReference", "tsPrimitiveValue", "tsFunctionValue", "tsObjectValue", "tsArrayValue", "tsArrayBuiltin", "tsClassValue", "tsErrorValue", "tsPromiseValue", "tsCollectionValue", "tsIteratorValueReference", "tsBigIntBinary", "tsBigIntSigned", "tsBigIntFrom", "tsConvertBigInt", "tsDateConstruct", "tsFSConstructStats", "tsFSConstructDirent", "tsFSDateMillis", "tsFSDate", "tsStringStatic", "tsStringConstructor", "tsBooleanConstructor", "tsNumberConstructor", "tsBigIntFunction", "tsUnicodeString", "tsRegExpLegacyGetter", "tsRegExpConstructor", "tsNewCompiledRegExp", "tsRegExpStringIteratorPrototype", "tsBinary", "tsPolicyCompare", "tsBoundary", "tsArrayBoundary", "tsComparisonOperand", "tsConvertValue", "tsToPrimitive", "tsPrimitiveHint":
			return "tsValue"
		case "tsNodeQUICMakeConstructor", "tsNodeQUICMakeStatsClass", "tsNodeQUICStateConstructor", "tsNodeSQLiteConstructor", "tsInternalMockTimersConstructor", "tsNodeTestNewMockTracker", "tsNodeTestMockFunctionContext":
			return "tsValue"
		case "tsNodeVMValidateModule":
			return "*tsNodeVMModuleState"
		case "tsObjectPrototype", "tsObjectConstructor", "tsProxyConstructor", "tsFunctionConstructor", "tsFunctionPrototype", "tsClassPrototype", "tsPrototypeOf", "tsRegExpPrototype", "tsErrorPrototype", "tsErrorConstructor", "tsGlobalObject", "tsGeneratorPrototype", "tsGeneratorFunctionPrototype", "tsWellKnownSymbol", "tsSymbolFunction", "tsPathModuleFor", "tsNodeURLPathToFileURL", "tsNodeURLFileURLToPath", "tsNodeOSConstants", "tsNodeOSString", "tsNodeOSNumber", "tsNodeOSCPUs", "tsNodeOSNetworkInterfaces", "tsNodeTimerRejectedIterator", "tsBuiltinClass", "tsInternalPrimordialsGetter", "tsNodeModuleConstructorValue", "tsNodeModuleCompileCacheStatus", "tsNodeParseArgsFunction", "tsNodeUtilSystemErrorMap", "tsNodePerfEntry", "tsNodePerfEntryList", "tsNodePerfGetEntries", "tsNodePerfUtilization", "tsNodePerfObserverConstructor", "tsNodeHistogramValue", "tsInternalSharedPerformanceEntry", "tsInternalSharedNodeEntryList", "tsInternalSharedPerformanceNodeEntry", "tsInternalSharedNodeObserverConstructor", "tsInternalSharedNodeObserverDisconnect", "tsInternalSharedErrorClass", "tsInternalSharedHideFromStack", "tsInternalValidatorCheckRanges", "tsInternalValidatorPath", "tsInternalValidatorFsPath", "tsInternalValidatorCheckLink", "tsInternalValidatorPort", "tsNodeDuplexifyOptions", "tsTypedArrayBufferExtended", "tsTypedArrayBuffer", "tsTypedArrayBigIntElement", "tsBigIntElement", "tsTypedArrayElement", "tsArrayMutate", "tsNodeIterFromIteratorObject", "tsNodeIterFromIterable", "tsNodeIterFromEmpty", "tsNodeIterFromArrayBatch", "tsNodeIterFromAsyncIterator", "tsNodeIterFromNormalizeAsyncValue", "tsNodeIterFromDeferredAsyncResult", "tsNodeIterFromArray", "tsNodeClassicBatchedReadable", "tsNodeCryptoKeyValue", "tsNodeCryptoCallable", "tsNodeCryptoConstructorPrototype", "tsNodeProcessBindingConstants", "tsNodeProcessReallyExit", "tsNodeProcessResourceUsage", "tsNodeProcessConfig", "tsNodeProcessStdin", "tsNodeProcessWriteStream", "tsNodeProcessMemoryUsage", "tsNodeProcessSetUncaughtExceptionCapture", "tsNodeDomainEnter", "tsNodeDomainExit", "tsNodeDiagnosticsChannelConstructor", "tsNodeALSMethod", "tsNodeAsyncResourceMethod", "tsNodeAsyncLocalStorageConstructor", "tsNodeAsyncResourceConstructor", "tsNodeAsyncHooksCreateHook", "tsNodeAssertDecorate", "tsGet", "tsSet", "tsSetPrototype", "tsObjectBuiltin", "tsConstruct", "tsFSProperty", "tsFSOption", "tsFSStats", "tsFSZeroStats", "tsFSDir", "tsFSDirent", "tsFSInvoke", "tsIteratorResult", "tsNodeBuffer", "tsBYOBTransfer", "tsWebStreamValue", "tsFSReadLines", "tsAsyncIteratorObject", "tsBlobObject", "tsNamedNativeMethod", "tsNodeURLParamsObject", "tsNodeURLObject", "tsNodeNetMakeAddr", "tsNodeNetNormalizeArgs", "tsNodeDNSObject", "tsNodeDNSStrings", "tsNodeDNSAsync", "tsNodeDNSLookupAsync", "tsNodeDNSServiceAsync", "tsNodeDNSError", "tsInternalTLSKey", "tsInternalTLSError", "tsInternalTLSSecureRange", "tsInternalTLSStripCipherNames", "tsInternalTLSNormalizePEM", "tsInternalTLSClientOptions", "tsInternalTLSHandshakeError", "tsInternalTLSProcessPFX", "tsInternalTLSUnsealPFX":
			return "tsValue"
		case "float64", "float32", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "bool", "string":
			return name
		case "make":
			if len(e.Args) > 0 {
				return valueType(e.Args[0])
			}
		case "append":
			if len(e.Args) > 0 {
				return p.kind(e.Args[0], env)
			}
		case "len", "cap":
			return "int"
		case "new":
			if len(e.Args) > 0 {
				return "*" + valueType(e.Args[0])
			}
		}
		if _, ok := e.Fun.(*ast.ParenExpr); ok {
			return strings.Trim(valueType(e.Fun), "()")
		}
		if fn, ok := e.Fun.(*ast.FuncLit); ok {
			return valueFunctionSignature(fn.Type).result
		}
		return p.signature(e, env).result
	case *ast.BinaryExpr:
		switch e.Op {
		case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ, token.LAND, token.LOR:
			return "bool"
		}
		left, right := p.kind(e.X, env), p.kind(e.Y, env)
		if left == "string" || right == "string" {
			return "string"
		}
		if left == "float64" || right == "float64" {
			return "float64"
		}
		return left
	}
	return ""
}
func valueWrapper(name string, expr ast.Expr) ast.Expr {
	return &ast.CallExpr{Fun: ast.NewIdent(name), Args: []ast.Expr{expr}}
}
func (p *valueEmission) coerce(expr ast.Expr, expected string, env map[string]string) ast.Expr {
	if expected == "tsBindingCell" {
		kind := p.kind(expr, env)
		if kind == "nil" {
			return &ast.CompositeLit{Type: ast.NewIdent("tsBindingCell")}
		}
		if kind == "*tsCell" || (valueBase(kind) == "tsTypedCell" || valueBase(kind) == "tsDenseArrayCell") {
			return &ast.CallExpr{Fun: &ast.SelectorExpr{X: expr, Sel: ast.NewIdent("bindingRef")}}
		}
	}
	if expected != "tsValue" {
		return expr
	}
	kind := p.kind(expr, env)
	switch kind {
	case "tsValue":
		return expr
	case "nil":
		return ast.NewIdent("tsNull")
	case "float64":
		return valueWrapper("tsNumberValue", expr)
	case "float32", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return valueWrapper("ts"+strings.ToUpper(kind[:1])+kind[1:]+"Value", expr)
	case "bool":
		return valueWrapper("tsBooleanValue", expr)
	case "*tsString":
		return valueWrapper("tsStringReference", expr)
	case "string":
		return valueWrapper("tsStringReference", valueWrapper("tsStringUTF8", expr))
	case "T":
		return valueWrapper("tsPrimitiveValue", expr)
	}
	constructors := map[string]string{"*tsObject": "tsObjectValue", "*tsArray": "tsArrayValue", "*tsFunction": "tsFunctionValue", "*tsClass": "tsClassValue", "*tsProperties": "tsInstanceValue", "*tsRegExp": "tsRegExpValue", "*tsRuntimeError": "tsErrorValue", "*tsPromise": "tsPromiseValue", "*tsTask": "tsTaskValue", "*tsModule": "tsModuleValue", "*tsIterator": "tsIteratorValueReference", "tsNamespace": "tsNamespaceValue", "tsImportRef": "tsImportValue"}
	if constructor := constructors[kind]; constructor != "" {
		return valueWrapper(constructor, expr)
	}
	if strings.HasPrefix(valueBase(kind), "tsOptional[") {
		return &ast.CallExpr{Fun: &ast.SelectorExpr{X: expr, Sel: ast.NewIdent("raw")}}
	}
	if strings.HasPrefix(kind, "*") && p.fields[valueBase(kind)]["properties"] == "*tsProperties" {
		return valueWrapper("tsInstanceValue", &ast.SelectorExpr{X: expr, Sel: ast.NewIdent("properties")})
	}
	if p.err == nil {
		var detail bytes.Buffer
		_ = format.Node(&detail, token.NewFileSet(), expr)
		candidates := ""
		if call, ok := expr.(*ast.CallExpr); ok {
			candidates = fmt.Sprintf("; call %s signatures: %#v", valueCallName(call.Fun), p.functions[valueCallName(call.Fun)])
		}
		context := ""
		if p.callContext != "" {
			context = "; within " + p.callContext
		}
		p.err = fmt.Errorf("cannot marshal %s (%s) into tsValue in %s%s%s", detail.String(), kind, p.scope, candidates, context)
	}
	return expr
}
func (p *valueEmission) expression(expr ast.Expr, env map[string]string) ast.Expr {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *ast.CallExpr:
		previousContext := p.callContext
		callName := valueCallName(e.Fun)
		if previousContext == "" {
			p.callContext = callName
		} else {
			p.callContext = previousContext + " -> " + callName
		}
		currentContext := p.callContext
		defer func() { p.callContext = previousContext }()
		e.Fun = p.expression(e.Fun, env)
		for i, arg := range e.Args {
			p.callContext = fmt.Sprintf("%s arg[%d]", currentContext, i)
			e.Args[i] = p.expression(arg, env)
		}
		p.callContext = currentContext
		if valueCallName(e.Fun) == "append" && len(e.Args) > 1 {
			element := valueElement(p.kind(e.Args[0], env))
			for i := 1; i < len(e.Args); i++ {
				if e.Ellipsis.IsValid() && i == len(e.Args)-1 {
					continue
				}
				e.Args[i] = p.coerce(e.Args[i], element, env)
			}
		} else {
			sig := p.signature(e, env)
			for i, arg := range e.Args {
				if e.Ellipsis.IsValid() && i == len(e.Args)-1 {
					continue
				}
				at := i
				if sig.variadic && at >= len(sig.parameters) {
					at = len(sig.parameters) - 1
				}
				if at >= 0 && at < len(sig.parameters) {
					p.callContext = fmt.Sprintf("%s arg[%d]", currentContext, i)
					e.Args[i] = p.coerce(arg, sig.parameters[at], env)
				}
			}
		}
	case *ast.TypeAssertExpr:
		e.X = p.expression(e.X, env)
		if p.kind(e.X, env) == "tsValue" {
			if p.err == nil {
				p.err = fmt.Errorf("interface assertion on JavaScript value: %s", valueType(e))
			}
		}
	case *ast.CompositeLit:
		kind := valueType(e.Type)
		fields := p.fields[valueBase(kind)]
		order := []string{}
		if structure, ok := e.Type.(*ast.StructType); ok {
			_ = structure
		}
		for _, field := range p.fieldOrder(kind) {
			order = append(order, field)
		}
		for i, item := range e.Elts {
			if pair, ok := item.(*ast.KeyValueExpr); ok {
				pair.Value = p.expression(pair.Value, env)
				expected := valueElement(kind)
				if key, ok := pair.Key.(*ast.Ident); ok && fields != nil {
					expected = fields[key.Name]
				}
				pair.Value = p.coerce(pair.Value, expected, env)
			} else {
				item = p.expression(item, env)
				expected := valueElement(kind)
				if expected == "" && i < len(order) {
					expected = fields[order[i]]
				}
				e.Elts[i] = p.coerce(item, expected, env)
			}
		}
	case *ast.FuncLit:
		child := copyValueEnvironment(env)
		p.parameters(e.Type, child)
		p.block(e.Body, child, valueReturnTypes(e.Type))
	case *ast.BinaryExpr:
		e.X = p.expression(e.X, env)
		e.Y = p.expression(e.Y, env)
		if e.Op == token.EQL || e.Op == token.NEQ {
			if p.kind(e.X, env) == "tsValue" && p.kind(e.Y, env) == "nil" {
				e.Y = ast.NewIdent("tsNull")
			}
			if p.kind(e.Y, env) == "tsValue" && p.kind(e.X, env) == "nil" {
				e.X = ast.NewIdent("tsNull")
			}
		}
	case *ast.UnaryExpr:
		e.X = p.expression(e.X, env)
	case *ast.StarExpr:
		e.X = p.expression(e.X, env)
	case *ast.ParenExpr:
		e.X = p.expression(e.X, env)
	case *ast.SelectorExpr:
		e.X = p.expression(e.X, env)
	case *ast.IndexListExpr:
		e.X = p.expression(e.X, env)
	case *ast.IndexExpr:
		e.X = p.expression(e.X, env)
		e.Index = p.expression(e.Index, env)
	case *ast.SliceExpr:
		e.X = p.expression(e.X, env)
		e.Low = p.expression(e.Low, env)
		e.High = p.expression(e.High, env)
		e.Max = p.expression(e.Max, env)
	}
	return expr
}

// Positional runtime literals are deliberately limited to these stable records.
func (p *valueEmission) fieldOrder(kind string) []string {
	switch valueBase(kind) {
	case "tsCell":
		return []string{"value", "initialized", "constant"}
	case "tsResult":
		return []string{"value", "rejected"}
	case "tsThrown":
		return []string{"value"}
	case "tsOptional":
		return []string{"value", "tag"}
	case "tsRuntimeError":
		return []string{"name", "message"}
	case "tsImportRef":
		return []string{"module", "name"}
	case "tsNamespace":
		return []string{"module"}
	}
	return nil
}
func (p *valueEmission) multipleAssignmentResults(expr ast.Expr, env map[string]string) []string {
	switch value := expr.(type) {
	case *ast.CallExpr:
		switch valueCallName(value.Fun) {
		case "tsFSReadDirectory":
			return []string{"tsValue", "error"}
		case "tsRegExpLiteralSearch":
			return []string{"tsValue", "bool"}
		case "tsNodeDNSResolveValues":
			return []string{"tsValue", "error"}
		}
		return p.signature(value, env).results
	case *ast.IndexExpr:
		if mapType := p.kind(value.X, env); strings.HasPrefix(mapType, "map[") {
			return []string{p.kind(value, env), "bool"}
		}
	case *ast.TypeAssertExpr:
		return []string{valueType(value.Type), "bool"}
	case *ast.UnaryExpr:
		if value.Op == token.ARROW {
			channelType := p.kind(value.X, env)
			channelType = strings.TrimPrefix(channelType, "<-chan ")
			channelType = strings.TrimPrefix(channelType, "chan ")
			return []string{channelType, "bool"}
		}
	}
	return nil
}
func (p *valueEmission) declaration(decl *ast.GenDecl, env map[string]string) {
	for _, spec := range decl.Specs {
		if v, ok := spec.(*ast.ValueSpec); ok {
			kind := valueType(v.Type)
			for i, expr := range v.Values {
				v.Values[i] = p.coerce(p.expression(expr, env), kind, env)
			}
			for i, name := range v.Names {
				actual := kind
				if actual == "" && i < len(v.Values) {
					actual = p.kind(v.Values[i], env)
				}
				env[name.Name] = actual
			}
		}
	}
}
func (p *valueEmission) block(block *ast.BlockStmt, env map[string]string, result string) {
	if block == nil {
		return
	}
	for _, stmt := range block.List {
		p.statement(stmt, env, result)
	}
}
func (p *valueEmission) statement(stmt ast.Stmt, env map[string]string, result string) {
	if stmt == nil {
		return
	}
	switch s := stmt.(type) {
	case *ast.BlockStmt:
		p.block(s, copyValueEnvironment(env), result)
	case *ast.DeclStmt:
		if decl, ok := s.Decl.(*ast.GenDecl); ok {
			p.declaration(decl, env)
		}
	case *ast.AssignStmt:
		for i, expr := range s.Rhs {
			s.Rhs[i] = p.expression(expr, env)
		}
		var tupleResults []string
		if len(s.Rhs) == 1 {
			tupleResults = p.multipleAssignmentResults(s.Rhs[0], env)
		}
		multipleResultAssignment := len(s.Rhs) == 1 && len(tupleResults) == len(s.Lhs) && len(tupleResults) > 1
		for i, left := range s.Lhs {
			left = p.expression(left, env)
			s.Lhs[i] = left
			if multipleResultAssignment {
				if name, ok := left.(*ast.Ident); ok && s.Tok == token.DEFINE {
					env[name.Name] = tupleResults[i]
				}
				continue
			}
			if i >= len(s.Rhs) {
				if name, ok := left.(*ast.Ident); ok && s.Tok == token.DEFINE {
					kind := ""
					if i < len(tupleResults) {
						kind = tupleResults[i]
					}
					env[name.Name] = kind
				}
				continue
			}
			expected := p.kind(left, env)
			if s.Tok == token.DEFINE {
				expected = ""
			}
			s.Rhs[i] = p.coerce(s.Rhs[i], expected, env)
			if name, ok := left.(*ast.Ident); ok && s.Tok == token.DEFINE {
				env[name.Name] = p.kind(s.Rhs[i], env)
			}
		}
	case *ast.ExprStmt:
		s.X = p.expression(s.X, env)
	case *ast.ReturnStmt:
		for i, expr := range s.Results {
			target := result
			if results := strings.Split(result, "|"); len(results) > 1 {
				target = ""
				if i < len(results) {
					target = results[i]
				}
			}
			s.Results[i] = p.coerce(p.expression(expr, env), target, env)
		}
	case *ast.IfStmt:
		child := copyValueEnvironment(env)
		p.statement(s.Init, child, result)
		s.Cond = p.expression(s.Cond, child)
		p.block(s.Body, copyValueEnvironment(child), result)
		p.statement(s.Else, copyValueEnvironment(child), result)
	case *ast.ForStmt:
		child := copyValueEnvironment(env)
		p.statement(s.Init, child, result)
		s.Cond = p.expression(s.Cond, child)
		p.statement(s.Post, child, result)
		p.block(s.Body, child, result)
	case *ast.RangeStmt:
		s.X = p.expression(s.X, env)
		child := copyValueEnvironment(env)
		if name, ok := s.Key.(*ast.Ident); ok {
			child[name.Name] = "int"
		}
		if name, ok := s.Value.(*ast.Ident); ok {
			child[name.Name] = valueElement(p.kind(s.X, env))
		}
		p.block(s.Body, child, result)
	case *ast.SwitchStmt:
		child := copyValueEnvironment(env)
		p.statement(s.Init, child, result)
		s.Tag = p.expression(s.Tag, child)
		for _, item := range s.Body.List {
			clause := item.(*ast.CaseClause)
			for i, expr := range clause.List {
				clause.List[i] = p.expression(expr, child)
			}
			local := copyValueEnvironment(child)
			for _, stmt := range clause.Body {
				p.statement(stmt, local, result)
			}
		}
	case *ast.TypeSwitchStmt:
		child := copyValueEnvironment(env)
		p.statement(s.Init, child, result)
		p.statement(s.Assign, child, result)
		variable := ""
		if assign, ok := s.Assign.(*ast.AssignStmt); ok {
			if name, ok := assign.Lhs[0].(*ast.Ident); ok {
				variable = name.Name
			}
		}
		for _, item := range s.Body.List {
			clause := item.(*ast.CaseClause)
			local := copyValueEnvironment(child)
			if len(clause.List) == 1 {
				local[variable] = valueType(clause.List[0])
			}
			for _, stmt := range clause.Body {
				p.statement(stmt, local, result)
			}
		}
	case *ast.DeferStmt:
		s.Call = p.expression(s.Call, env).(*ast.CallExpr)
	case *ast.GoStmt:
		s.Call = p.expression(s.Call, env).(*ast.CallExpr)
	case *ast.SendStmt:
		s.Value = p.expression(s.Value, env)
		s.Chan = p.expression(s.Chan, env)
	case *ast.SelectStmt:
		for _, item := range s.Body.List {
			clause := item.(*ast.CommClause)
			child := copyValueEnvironment(env)
			p.statement(clause.Comm, child, result)
			for _, stmt := range clause.Body {
				p.statement(stmt, child, result)
			}
		}
	case *ast.LabeledStmt:
		p.statement(s.Stmt, env, result)
	}
}

func valueReturnTypes(fn *ast.FuncType) string {
	results := []string{}
	if fn.Results != nil {
		for _, field := range fn.Results.List {
			count := len(field.Names)
			if count == 0 {
				count = 1
			}
			for i := 0; i < count; i++ {
				results = append(results, valueType(field.Type))
			}
		}
	}
	return strings.Join(results, "|")
}
