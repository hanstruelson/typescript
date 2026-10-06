package goemit

import (
	"fmt"
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"strings"
)

type nativeCallbackArg struct {
	name, typ string
	p         primitive
	array     bool
}

func (e *emitter) primitiveFunctionReturn(node *ast.Node) primitive {
	if resolver, ok := e.resolver.(interface {
		GetEmitFunctionReturnPrimitiveType(*ast.Node, []*ast.Node) (string, uint8)
	}); ok {
		kind, nulls := resolver.GetEmitFunctionReturnPrimitiveType(node, e.specializationCalls)
		return primitive{kind, nulls}
	}
	return e.returnPrimitive(node)
}
func (b *machineBuilder) callbackReturn(node *ast.Node) primitive {
	return b.e.primitiveFunctionReturn(node)
}

// Literal synchronous callbacks have an entirely native ABI. Actual dynamic
// callable values continue to use the TS function ABI, since their parameter
// types cannot be assumed at compile time.
func (b *machineBuilder) nativeCallback(node *ast.Node, args []nativeCallbackArg) (string, primitive, bool) {
	for node.Kind == ast.KindParenthesizedExpression {
		node = node.AsParenthesizedExpression().Expression
	}
	if node.Kind == ast.KindIdentifier {
		if fn := b.e.nativeFunctions[b.e.reference(node)]; fn != nil {
			if cell := b.e.binding(node); cell != nil && cell.declaration.Kind == ast.KindVariableDeclaration {
				b.emit("_=" + cell.name + ".get()")
			}
			if len(fn.parameters) > len(args) {
				return "", primitive{}, false
			}
			signature := []string{"loop *tsLoop"}
			values := []string{"loop"}
			for i, arg := range args {
				signature = append(signature, arg.name+" "+arg.typ)
				if i >= len(fn.parameters) {
					continue
				}
				b.tempTypes[arg.name] = arg.typ
				if fn.arrayParameters[i].kind != "" {
					if !arg.array {
						return "", primitive{}, false
					}
					value := arg.name
					target := fn.arrayParameters[i]
					if target.kind != arg.p.kind || target.nulls != arg.p.nulls {
						value = fmt.Sprintf("tsDenseStorage[%s](tsArrayBoundary(tsArrayValue(%s.owner),%q,%d,%t),%q,%d)", target.goType(), arg.name, target.kind, target.nulls, b.e.coerce, target.kind, target.nulls)
					}
					values = append(values, value)
				} else {
					values = append(values, b.nativeCallbackConversion(arg.name, arg.p, fn.parameters[i]))
				}
			}
			f := "func(" + strings.Join(signature, ",") + ") " + fn.result.goType() + " {return " + fn.name + "(" + strings.Join(values, ",") + ")}"
			return b.typedTemp(f, "func("+strings.Join(callbackTypes(args), ",")+") "+fn.result.goType()), fn.result, true
		}
	}

	if node.Kind == ast.KindIdentifier {
		if definition := b.e.stableCallbacks[b.e.reference(node)]; definition != nil {
			if cell := b.e.binding(node); cell != nil {
				b.emit("_=" + cell.name + ".get()")
			}
			node = definition
		}
	}

	if node.Kind != ast.KindArrowFunction && node.Kind != ast.KindFunctionExpression && node.Kind != ast.KindFunctionDeclaration {
		return "", primitive{}, false
	}
	if node.Name() != nil && node.Kind != ast.KindFunctionDeclaration || !nativeCallbackBody(node) || len(node.TypeParameters()) != 0 {
		return "", primitive{}, false
	}
	if len(node.Parameters()) > len(args) {
		return "", primitive{}, false
	}
	for i, param := range node.Parameters() {
		decl := param.AsParameterDeclaration()
		if param.Name().Kind != ast.KindIdentifier || decl.DotDotDotToken != nil {
			return "", primitive{}, false
		}
		cell := b.e.bindings[param]
		if cell == nil {
			return "", primitive{}, false
		}
		if args[i].array {
			if cell.arrayElement.kind == "" {
				return "", primitive{}, false
			}
		} else if cell.primitive.kind == "" {
			return "", primitive{}, false
		}
	}
	r := b.callbackReturn(node)
	child := b.e.newMachine(node, b, false)
	child.direct = directClassBody(node)
	child.tempTypes = map[string]string{}
	child.nativeResult = r
	arrayResult := b.e.arrayElementPrimitive(node)
	if r.kind == "" && arrayResult.kind != "" {
		child.nativeArrayResult = arrayResult
	}
	if !child.direct {
		child.nativeOutput = child.e.unique("NativeReturnVector")
		child.temps = append(child.temps, child.nativeOutput)
		child.tempTypes[child.nativeOutput] = "tsValue"
		if r.kind != "" {
			child.tempTypes[child.nativeOutput] = r.elementGoType()
			child.blocks[0].lines = nil
		}
	}
	if r.kind != "" {
		child.nativeReturn = r.elementGoType()
	}
	if arrayResult.kind != "" {
		child.nativeReturn = "*tsGrowableStorage[" + arrayResult.goType() + "]"
		if !child.direct {
			child.tempTypes[child.nativeOutput] = child.nativeReturn
			child.blocks[0].lines = nil
		}
	}
	if node.Kind == ast.KindArrowFunction {
		child.concrete = b.concrete
		child.declaring = b.declaring
		child.receiver = b.receiver
	}
	signature := []string{"loop *tsLoop"}
	for i, arg := range args {
		signature = append(signature, arg.name+" "+arg.typ)
		child.tempTypes[arg.name] = arg.typ
		if i >= len(node.Parameters()) {
			continue
		}
		cell := b.e.bindings[node.Parameters()[i]]
		if arg.array {
			value := arg.name
			if cell.arrayElement.kind != arg.p.kind || cell.arrayElement.nulls != arg.p.nulls {
				target := cell.arrayElement
				value = fmt.Sprintf("tsDenseStorage[%s](tsArrayBoundary(tsArrayValue(%s.owner),%q,%d,%t),%q,%d)", target.goType(), arg.name, target.kind, target.nulls, b.e.coerce, target.kind, target.nulls)
			}
			child.emit(cell.name + ".initStorage(" + value + ")")
			continue
		}
		method := "initNative"
		if cell.primitive.nulls != 0 {
			method = "initOptional"
		}

		initializer := node.Parameters()[i].AsParameterDeclaration().Initializer
		if initializer != nil && arg.p.nulls&2 != 0 {
			if child.direct {
				child.emit("if " + arg.name + ".tag==2 {")
				fallback := child.denseElement(initializer, cell.primitive)
				child.emit(cell.name + "." + method + "(" + fallback + ")")
				child.emit("} else {")
				value := child.comparisonBoundary(arg.name, cell.primitive)
				child.emit(cell.name + "." + method + "(" + value + ")")
				child.emit("}")
			} else {
				missing, present, end := child.block(), child.block(), child.block()
				child.emit(fmt.Sprintf("if %s.tag==2{m.pc=%d}else{m.pc=%d};return", arg.name, missing, present))
				child.current = missing
				fallback := child.denseElement(initializer, cell.primitive)
				child.emit(cell.name + "." + method + "(" + fallback + ")")
				child.jump(end)
				child.current = present
				value := child.comparisonBoundary(arg.name, cell.primitive)
				child.emit(cell.name + "." + method + "(" + value + ")")
				child.jump(end)
				child.current = end
			}
		} else {
			value := child.comparisonBoundary(arg.name, cell.primitive)
			child.emit(cell.name + "." + method + "(" + value + ")")
		}

	}
	if node.Body().Kind == ast.KindBlock {
		if child.direct {
			child.directStatements(node.Body().AsBlock().Statements.Nodes)
		} else {
			child.statements(node.Body().AsBlock().Statements.Nodes)
			child.endNativeMachine()
		}
	} else {
		value := ""
		if arrayResult.kind != "" {
			value = child.denseStorage(node.Body(), arrayResult)
		} else if r.kind != "" {
			value = child.denseElement(node.Body(), r)
		} else {
			value = child.expression(node.Body())
		}
		if arrayResult.kind != "" {
		} else if r.kind != "" {
			value = child.comparisonBoundary(value, r)
		} else {
			value = child.returnValue(value)
		}
		child.emit("return " + value)
	}
	resultType := "tsValue"
	if r.kind != "" {
		resultType = r.elementGoType()
	}
	if arrayResult.kind != "" {
		resultType = "*tsGrowableStorage[" + arrayResult.goType() + "]"
	}
	body := child.finishDirect()
	if !child.direct {
		body = child.finishNativeMachine()
	}
	function := "func(" + strings.Join(signature, ",") + ") " + resultType + " {\n" + body + "\n}"
	captures, values := []string{}, []string{}
	for _, cell := range child.captures {
		captures = append(captures, cell.name+" "+cell.pointerType())
		values = append(values, cell.name)
	}
	function = "func(" + strings.Join(captures, ",") + ") func(" + strings.Join(callbackTypes(args), ",") + ") " + resultType + " {return " + function + "}(" + strings.Join(values, ",") + ")"
	return b.typedTemp(function, "func("+strings.Join(callbackTypes(args), ",")+") "+resultType), r, true
}
func (b *machineBuilder) nativeCallbackConversion(value string, source, target primitive) string {
	if source.nulls != 0 && target.nulls == 0 && source.kind != "" {
		return b.nativeCallbackConversion("tsDensePresent("+value+")", primitive{kind: source.kind}, target)
	}
	if source.goType() == target.goType() && source.nulls == 0 {
		return value
	}
	if source.nulls == 0 && target.nulls == 0 && source.kind != "" {
		return "tsConvert" + conversionName(target.kind) + "From" + conversionName(source.kind) + "(" + value + "," + fmt.Sprint(b.e.coerce) + ")"
	}
	return fmt.Sprintf("tsNative[%s](tsBoundary(%s,%q,%d,%t))", target.goType(), value, target.kind, target.nulls, b.e.coerce)
}
func callbackTypes(args []nativeCallbackArg) []string {
	types := []string{"*tsLoop"}
	for _, arg := range args {
		types = append(types, arg.typ)
	}
	return types
}
func (b *machineBuilder) nativeCondition(value string, p primitive) string {
	if b.tempType(value) == "tsValue" {
		return "tsTruthy(" + value + ")"
	}
	return b.nativeTruth(value, p)
}
func (b *machineBuilder) nativeTruth(value string, p primitive) string {
	if p.kind == "" {
		return "tsTruthy(" + value + ")"
	}
	if p.nulls != 0 {
		return value + ".tag==0&&(" + b.nativeTruth(value+".value", primitive{kind: p.kind}) + ")"
	}
	switch p.kind {
	case "boolean":
		return value
	case "string":
		return "len(" + value + ".units)!=0"
	case "number", "float64", "float32":
		return value + "!=0&&!math.IsNaN(float64(" + value + "))"
	default:
		return value + "!=0"
	}
}
func (b *machineBuilder) nativeCallbackCall(call *ast.CallExpression, p primitive) (string, bool) {
	if p.kind == "" {
		return "", false
	}
	name := call.Expression.AsPropertyAccessExpression().Name().Text()
	switch name {
	case "map", "flatMap", "filter", "forEach", "every", "some", "find", "findLast", "findIndex", "findLastIndex", "reduce", "reduceRight", "sort", "toSorted":
	default:
		return "", false
	}
	nodes := call.Arguments.Nodes
	if len(nodes) == 0 {
		return "", false
	}
	for _, node := range nodes {
		if node.Kind == ast.KindSpreadElement {
			return "", false
		}
	}
	// First prove eligibility, before evaluating the receiver or any arguments.
	callbackValue := nodes[0]
	callback := nodes[0]
	for callback.Kind == ast.KindParenthesizedExpression {
		callback = callback.AsParenthesizedExpression().Expression
	}
	if callback.Kind == ast.KindIdentifier {
		if fn := b.e.nativeFunctions[b.e.reference(callback)]; fn != nil {
			callback = fn.node
		}
	}
	if callback.Kind == ast.KindIdentifier {
		decl := b.e.reference(callback)
		if decl != nil && decl.Kind == ast.KindVariableDeclaration {
			cell := b.e.bindings[decl]
			init := decl.AsVariableDeclaration().Initializer
			if cell != nil && cell.constant && nativeStableCallbackBinding(decl, cell.owner) && init != nil && (init.Kind == ast.KindArrowFunction || init.Kind == ast.KindFunctionExpression) {
				callback = init
				callbackValue = init
			}
		}
	}
	if callback.Kind == ast.KindIdentifier {
		if definition := b.e.stableCallbacks[b.e.reference(callback)]; definition != nil {
			callback = definition
		}
	}
	if (callback.Kind != ast.KindArrowFunction && callback.Kind != ast.KindFunctionExpression && callback.Kind != ast.KindFunctionDeclaration) || !nativeCallbackBody(callback) {
		return "", false
	}
	for _, param := range callback.Parameters() {
		decl := param.AsParameterDeclaration()
		if param.Name().Kind != ast.KindIdentifier || decl.DotDotDotToken != nil {
			return "", false
		}
		cell := b.e.bindings[param]
		if cell == nil || cell.primitive.kind == "" && cell.arrayElement.kind == "" {
			return "", false
		}
	}
	if callback.Name() != nil && callback.Kind != ast.KindFunctionDeclaration || len(callback.TypeParameters()) != 0 {
		return "", false
	}
	expected := 3
	if name == "reduce" || name == "reduceRight" {
		expected = 4
	}
	if name == "sort" || name == "toSorted" {
		expected = 2
	}
	if len(callback.Parameters()) > expected {
		return "", false
	}
	for i, param := range callback.Parameters() {
		arrayPosition := expected - 1
		receivesArray := name != "sort" && name != "toSorted" && i == arrayPosition
		cell := b.e.bindings[param]
		if receivesArray != (cell.arrayElement.kind != "") {
			return "", false
		}
	}
	inferred := b.callbackReturn(callback)
	if (name == "map" || name == "reduce" || name == "reduceRight" || name == "flatMap" && b.e.arrayElementPrimitive(callback).kind == "") && inferred.kind == "" {
		return "", false
	}
	storage := b.denseStorage(call.Expression.AsPropertyAccessExpression().Expression, p)
	source := nativeCallbackArg{"nativeValue", p.elementGoType(), p, false}
	index := nativeCallbackArg{"nativeIndex", "float64", primitive{kind: "number"}, false}
	array := nativeCallbackArg{"nativeArray", "*tsGrowableStorage[" + p.goType() + "]", p, true}
	args := []nativeCallbackArg{source, index, array}
	result := b.callbackReturn(callback)
	reduce := name == "reduce" || name == "reduceRight"
	if reduce {
		result = b.e.primitive(call.Parent)
		if r := b.callbackReturn(callback); r.kind != "" {
			result = r
		}
		if result.kind == "" {
			return "", false
		}
		args = append([]nativeCallbackArg{{"nativeAccumulator", result.elementGoType(), result, false}}, args...)
	}
	if name == "sort" || name == "toSorted" {
		args = []nativeCallbackArg{source, {"nativeOther", p.elementGoType(), p, false}}
	}
	if len(callback.Parameters()) > len(args) {
		return "", false
	}
	if callbackValue != nodes[0] && nodes[0].Kind == ast.KindIdentifier {
		if cell := b.e.binding(nodes[0]); cell != nil {
			b.emit("_=" + cell.name + ".get()")
		}
	}
	fn, r, ok := b.nativeCallback(callbackValue, args)
	if !ok {
		return "", false
	}
	initial := ""
	if reduce && len(nodes) > 1 {
		initial = b.comparisonBoundary(b.expression(nodes[1]), result)
	}
	extras := 1
	if reduce {
		extras = 2
	}
	for _, node := range nodes[min(extras, len(nodes)):] {
		b.emit("_=" + b.expression(node))
	}
	item := "tsDenseRead[" + p.goType() + "](" + storage + ",float64(i))"
	if p.nulls != 0 {
		item = "tsNullableAt[" + p.goType() + "](" + storage + ",i)"
	}
	invoke := fn + "(loop," + item + ",float64(i)," + storage + ")"
	switch name {
	case "map", "flatMap":
		if name == "flatMap" {
			if out := b.e.arrayElementPrimitive(callback); out.kind != "" {
				return b.typedTemp(fmt.Sprintf("tsDenseFlatMap[%s,%s](%s,%q,%d,%t,func(i int)*tsGrowableStorage[%s]{return %s})", p.goType(), out.goType(), storage, out.kind, out.nulls, b.e.coerce, out.goType(), invoke), "*tsArray"), true
			}
		}
		out := r
		if target := b.e.arrayElementPrimitive(call.Parent); target.kind == out.kind {
			out.nulls |= target.nulls
		}
		if out.kind == "" {
			return "", false
		}
		helper := "tsDenseMap"
		ret := out.goType()
		if r.nulls != 0 {
			helper = "tsNullableMap"
			ret = out.elementGoType()
		}
		cb := "func(i int) " + ret + " {return " + invoke + "}"
		return b.typedTemp(fmt.Sprintf("%s[%s,%s](%s,%q,%d,%t,%s)", helper, p.goType(), out.goType(), storage, out.kind, out.nulls, b.e.coerce, cb), "*tsArray"), true
	case "forEach":
		return b.temp("tsDenseForEach[" + p.goType() + "](" + storage + ",func(i int){_=" + invoke + "})"), true
	case "filter", "every", "some", "find", "findLast", "findIndex", "findLastIndex":
		pred := "func(i int) bool {result:=" + invoke + ";return " + b.nativeTruth("result", r) + "}"
		if name == "filter" {
			return b.typedTemp(fmt.Sprintf("tsDenseFilter[%s](%s,%d,%t,%s)", p.goType(), storage, b.e.arrayElementPrimitive(call.AsNode()).nulls, b.e.coerce, pred), "*tsArray"), true
		}
		if name == "every" || name == "some" {
			return b.typedTemp(fmt.Sprintf("tsDenseTest[%s](%s,%t,%s)", p.goType(), storage, name == "every", pred), "bool"), true
		}
		if name == "find" || name == "findLast" {
			return b.typedTemp(fmt.Sprintf("tsDenseFind[%s](%s,%t,%s)", p.goType(), storage, name == "findLast", pred), "tsOptional["+p.goType()+"]"), true
		}
		found := b.typedTemp(fmt.Sprintf("tsDenseFindIndex[%s](%s,%t,%s)", p.goType(), storage, name == "findLast" || name == "findLastIndex", pred), "int")
		if name == "findIndex" || name == "findLastIndex" {
			return b.typedTemp("float64("+found+")", "float64"), true
		}
		return b.typedTemp("tsNullableAt["+p.goType()+"]("+storage+","+found+")", "tsOptional["+p.goType()+"]"), true
	case "reduce", "reduceRight":
		if initial == "" {
			initial = result.elementGoType() + "{}"
			if result.nulls == 0 {
				initial = result.goType() + "(0)"
				if result.kind == "string" {
					initial = `tsStringUTF8("")`
				} else if result.kind == "boolean" {
					initial = "false"
				}
			}
		}
		seed := b.e.unique("ReduceSeed")
		b.tempTypes[seed] = p.elementGoType()
		seedBuilder := b.e.newMachine(b.owner, b, false)
		seedBuilder.direct = true
		seedBuilder.tempTypes = map[string]string{}
		seedBuilder.tempTypes[seed] = p.elementGoType()
		seedValue := seedBuilder.comparisonBoundary(seed, result)
		seedCode := strings.Join(seedBuilder.blocks[0].lines, "\n")
		for _, temp := range seedBuilder.temps {
			seedCode = "var " + temp + " " + seedBuilder.tempType(temp) + ";_=" + temp + "\n" + seedCode
		}
		seedFn := "func(i int) " + result.elementGoType() + " {" + seed + ":=" + item + ";" + seedCode + ";return " + seedValue + "}"
		cb := "func(acc " + result.elementGoType() + ",i int) " + result.elementGoType() + " {return " + fn + "(loop,acc," + item + ",float64(i)," + storage + ")}"
		return b.typedTemp(fmt.Sprintf("tsDenseReduce[%s,%s](%s,%s,%t,%t,%s,%s)", p.goType(), result.elementGoType(), storage, initial, len(nodes) > 1, name == "reduceRight", seedFn, cb), result.elementGoType()), true
	case "sort", "toSorted":
		if name == "toSorted" {
			storage = b.typedTemp(fmt.Sprintf("tsDenseStorage[%s](tsDenseSlice(%s,%q,%t,0,math.Inf(1)),%q,%d)", p.goType(), storage, p.kind, b.e.coerce, p.kind, p.nulls), "*tsGrowableStorage["+p.goType()+"]")
		}
		first := "tsDenseRead[" + p.goType() + "](" + storage + ",float64(i))"
		second := "tsDenseRead[" + p.goType() + "](" + storage + ",float64(j))"
		if p.nulls != 0 {
			first = "tsNullableAt[" + p.goType() + "](" + storage + ",i)"
			second = "tsNullableAt[" + p.goType() + "](" + storage + ",j)"
		}
		comparison := fn + "(loop," + first + "," + second + ")"
		number := "result"
		if r.numeric() && r.nulls == 0 {
			number = "float64(result)"
		} else if r.kind != "number" && r.kind != "float64" {
			number = "tsNumber(result)"
		}
		less := "func(i,j int)bool{result:=" + comparison + ";return " + number + "<0}"
		b.emit("tsDenseSort[" + p.goType() + "](" + storage + "," + less + ")")
		return b.temp("tsArrayValue(" + storage + ".owner)"), true
	}
	return "", false
}

const NativeCallbackRuntime = `
func tsDenseFlatMap[T,O tsPrimitive](storage *tsGrowableStorage[T],kind string,nulls uint8,coerce bool,callback func(int)*tsGrowableStorage[O])*tsArray{length:=len(storage.values);result:=tsNewGrowableArray(kind,nulls,coerce);out:=(*tsGrowableStorage[O])(result.native);for i:=0;i<length;i++{if i>=len(storage.values){continue};items:=callback(i);if items.owner.elementNulls!=0{for j,item:=range items.values{checked:=tsNullableCheck(tsOptional[O]{value:item,tag:items.tags[j]},nulls);out.values=append(out.values,checked.value);if nulls!=0{out.tags=append(out.tags,checked.tag)}}}else{out.values=append(out.values,items.values...);if nulls!=0{out.tags=append(out.tags,make([]uint8,len(items.values))...)}}};return result}
func tsDenseForEach[T tsPrimitive](storage *tsGrowableStorage[T],callback func(int))tsValue{length:=len(storage.values);for i:=0;i<length;i++{if i<len(storage.values){callback(i)}};return tsU}
func tsDenseMap[T,O tsPrimitive](storage *tsGrowableStorage[T],kind string,nulls uint8,coerce bool,callback func(int)O)*tsArray{length:=len(storage.values);result:=tsNewGrowableArray(kind,nulls,coerce);out:=(*tsGrowableStorage[O])(result.native);out.values=make([]O,length);if nulls!=0{out.tags=make([]uint8,length)};for i:=0;i<length;i++{if i>=len(storage.values){tsArrayRangeFailure("Dense map cannot create holes after the source shrinks")};out.values[i]=callback(i)};return result}
func tsNullableMap[T,O tsPrimitive](storage *tsGrowableStorage[T],kind string,nulls uint8,coerce bool,callback func(int)tsOptional[O])*tsArray{length:=len(storage.values);result:=tsNewGrowableArray(kind,nulls,coerce);out:=(*tsGrowableStorage[O])(result.native);out.values=make([]O,length);out.tags=make([]uint8,length);for i:=0;i<length;i++{if i>=len(storage.values){tsArrayRangeFailure("Dense map cannot create holes after the source shrinks")};item:=tsNullableCheck(callback(i),nulls);out.values[i]=item.value;out.tags[i]=item.tag};return result}
func tsDenseFilter[T tsPrimitive](storage *tsGrowableStorage[T],nulls uint8,coerce bool,callback func(int)bool)*tsArray{length:=len(storage.values);result:=tsNewGrowableArray(storage.owner.elementKind,nulls,coerce);out:=(*tsGrowableStorage[T])(result.native);for i:=0;i<length;i++{if i>=len(storage.values){continue};item:=storage.values[i];tag:=uint8(0);if storage.owner.elementNulls!=0{tag=storage.tags[i]};if callback(i){checked:=tsNullableCheck(tsOptional[T]{value:item,tag:tag},nulls);out.values=append(out.values,checked.value);if nulls!=0{out.tags=append(out.tags,checked.tag)}}};return result}
func tsDenseTest[T tsPrimitive](storage *tsGrowableStorage[T],every bool,callback func(int)bool)bool{length:=len(storage.values);for i:=0;i<length;i++{if i>=len(storage.values){continue};if callback(i)!=every{return !every}};return every}
func tsDenseFind[T tsPrimitive](storage *tsGrowableStorage[T],reverse bool,callback func(int)bool)tsOptional[T]{length:=len(storage.values);start,step:=0,1;if reverse{start,step=length-1,-1};for i:=start;i>=0&&i<length;i+=step{if i>=len(storage.values){tsArrayRangeFailure("Dense find cannot supply a missing element to a typed callback")};item:=tsNullableAt(storage,i);if callback(i){return item}};return tsOptional[T]{tag:2}}
func tsDenseFindIndex[T tsPrimitive](storage *tsGrowableStorage[T],reverse bool,callback func(int)bool)int{length:=len(storage.values);start,step:=0,1;if reverse{start,step=length-1,-1};for i:=start;i>=0&&i<length;i+=step{if i>=len(storage.values){tsArrayRangeFailure("Dense find cannot supply a missing element to a typed callback")};if callback(i){return i}};return -1}
func tsDenseReduce[T tsPrimitive,R any](storage *tsGrowableStorage[T],acc R,has,reverse bool,seed func(int)R,callback func(R,int)R)R{length:=len(storage.values);start,step:=0,1;if reverse{start,step=length-1,-1};for i:=start;i>=0&&i<length;i+=step{if i>=len(storage.values){continue};if !has{acc=seed(i);has=true}else{acc=callback(acc,i)}};if !has{tsArrayTypeFailure("Reduce of empty array without initial value")};return acc}
type tsDenseSorter[T tsPrimitive]struct{storage *tsGrowableStorage[T];length int;less func(int,int)bool}
func(s tsDenseSorter[T])Len()int{return s.length}
func(s tsDenseSorter[T])Less(i,j int)bool{if len(s.storage.values)!=s.length{tsArrayRangeFailure("Cannot resize an array from its sort comparator")};if s.storage.owner.elementNulls!=0{a,b:=s.storage.tags[i],s.storage.tags[j];if a==2||b==2{return a!=2}};return s.less(i,j)}
func(s tsDenseSorter[T])Swap(i,j int){if len(s.storage.values)!=s.length{tsArrayRangeFailure("Cannot resize an array from its sort comparator")};s.storage.values[i],s.storage.values[j]=s.storage.values[j],s.storage.values[i];if s.storage.owner.elementNulls!=0{s.storage.tags[i],s.storage.tags[j]=s.storage.tags[j],s.storage.tags[i]}}
func tsDenseSort[T tsPrimitive](storage *tsGrowableStorage[T],less func(int,int)bool){sort.Stable(tsDenseSorter[T]{storage:storage,length:len(storage.values),less:less})}
`
