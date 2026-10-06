package goemit

import (
	"fmt"
	"go/constant"
	"go/token"
	"math"
	"strconv"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/scanner"
)

func (b *machineBuilder) identifier(node *ast.Node) string {
	if member := b.e.reference(node); member != nil && member.Kind == ast.KindEnumMember {
		if cell := b.e.bindings[member.Parent]; cell != nil {
			return "tsGet(" + cell.name + ".get()," + strconv.Quote(member.Name().Text()) + ")"
		}
	}
	if cell := b.e.binding(node); cell != nil {
		return cell.name + ".get()"
	}
	if imported := b.e.imports[b.e.reference(node)]; imported != "" {
		return b.builtin(imported)
	}
	switch node.Text() {
	case "String":
		return "tsFunc(func(args ...tsValue)tsValue {if len(args)==0 {return tsStringUnits(nil)};return tsStringValue(args[0])})"
	case "RegExp":
		return "tsFunc(tsNewRegExp)"
	case "NaN":
		return "math.NaN()"
	case "Infinity":
		return "math.Inf(1)"
	case "undefined":
		return "tsU"
	case "readFile", "delay", "setTimeout", "clearTimeout":
		// Ambient declarations describe runtime intrinsics; ordinary local bindings
		// always take precedence through the compiler reference resolver above.
		return b.builtin(node.Text())
	}
	b.e.fail(node, "unsupported or unresolved identifier "+node.Text())
	return "tsU"
}
func (b *machineBuilder) builtin(name string) string {
	switch name {
	case "readFile":
		return "tsFunc(func(args ...tsValue) tsValue { return loop.readFile(tsArg(args,0)) })"
	case "setTimeout":
		return "tsFunc(func(args ...tsValue) tsValue { extra:=[]tsValue{};if len(args)>2 {extra=args[2:]};return loop.setTimeout(tsArg(args,0),tsArg(args,1),extra...) })"
	case "clearTimeout":
		return "tsFunc(func(args ...tsValue) tsValue {return loop.clearTimeout(tsArg(args,0))})"
	case "delay":
		return "tsFunc(func(args ...tsValue) tsValue { return loop.delay(tsArg(args,0)) })"
	default:
		return "tsU"
	}
}
func (b *machineBuilder) expression(node *ast.Node) string {
	if node == nil {
		return "tsU"
	}
	if ast.IsOptionalChain(node) {
		return b.optionalChain(node)
	}
	switch node.Kind {
	case ast.KindNumericLiteral:
		text := strings.ReplaceAll(node.Text(), "_", "")
		number, err := strconv.ParseFloat(text, 64)
		if err != nil {
			b.e.fail(node, "unsupported numeric literal")
			return "tsU"
		}
		return b.typedTemp("float64("+strconv.FormatFloat(number, 'g', -1, 64)+")", "float64")
	case ast.KindObjectLiteralExpression:
		return b.objectLiteral(node)
	case ast.KindTaggedTemplateExpression:
		return b.taggedTemplate(node)
	case ast.KindRegularExpressionLiteral:
		text := node.Text()
		last := strings.LastIndex(text, "/")
		if last < 1 {
			b.e.fail(node, "invalid regular expression literal")
			return "tsU"
		}
		return b.temp("tsNewRegExp(" + stringLiteral(text[1:last]) + "," + stringLiteral(text[last+1:]) + ")")
	case ast.KindTemplateExpression:
		n := node.AsTemplateExpression()
		parts := []string{stringLiteral(n.Head.Text())}
		for _, span := range n.TemplateSpans.Nodes {
			item := span.AsTemplateSpan()
			value := b.expression(item.Expression)
			parts = append(parts, "tsStringValue("+value+")", stringLiteral(item.Literal.Text()))
		}
		return b.typedTemp("tsStringConcat("+strings.Join(parts, ",")+")", "*tsString")
	case ast.KindTypeOfExpression:
		value := b.expression(node.AsTypeOfExpression().Expression)
		return b.typedTemp("tsTypeOf("+value+")", "*tsString")
	case ast.KindStringLiteral, ast.KindNoSubstitutionTemplateLiteral:
		return b.typedTemp(stringLiteral(node.Text()), "*tsString")
	case ast.KindTrueKeyword:
		return b.typedTemp("true", "bool")
	case ast.KindFalseKeyword:
		return b.typedTemp("false", "bool")
	case ast.KindNullKeyword:
		return b.temp("nil")
	case ast.KindThisKeyword:
		if b.receiver == "" {
			b.e.fail(node, "this is supported only in class methods and their arrow closures")
			return "tsU"
		}
		if b.concrete != nil && b.concrete.static {
			return b.temp(b.receiver + ".properties.class")
		}
		return b.temp("tsRequireThis(" + b.receiver + ")")
	case ast.KindSuperKeyword:
		return b.receiver
	case ast.KindIdentifier:
		if cell := b.e.binding(node); cell != nil && b.e.strictNulls && (cell.primitive.kind == "" || cell.primitive.nulls != 0) && (!cell.maybeUndefined || b.definitelyInitialized(cell, node)) {
			if p := b.e.primitive(node); p.kind != "" && p.nulls == 0 {
				return b.typedTemp(fmt.Sprintf("tsNative[%s](tsBoundary(%s.get(),%q,0,false))", p.goType(), cell.name, p.kind), p.goType())
			}
		}
		if cell := b.e.binding(node); cell != nil && cell.primitive.kind != "" && (cell.primitive.nulls != 0 || (cell.maybeUndefined && !b.definitelyInitialized(cell, node))) {
			return b.typedTemp(cell.name+".optional()", "tsOptional["+cell.primitive.goType()+"]")
		}
		if cell := b.e.binding(node); cell != nil && cell.primitive.kind != "" && cell.primitive.nulls == 0 {
			safe := cell.declaration.Kind == ast.KindParameter || (cell.declaration.Kind == ast.KindVariableDeclaration && cell.declaration.AsVariableDeclaration().Initializer != nil)
			if safe {
				return b.typedTemp(cell.name+".read()", cell.primitive.goType())
			}
		}
		return b.temp(b.identifier(node))
	case ast.KindParenthesizedExpression:
		return b.expression(node.AsParenthesizedExpression().Expression)
	case ast.KindAsExpression:
		return b.numericAssertion(node.AsAsExpression().Expression, node.AsAsExpression().Type)
	case ast.KindTypeAssertionExpression:
		return b.numericAssertion(node.AsTypeAssertion().Expression, node.AsTypeAssertion().Type)
	case ast.KindSatisfiesExpression:
		return b.expression(node.AsSatisfiesExpression().Expression)
	case ast.KindNonNullExpression:
		return b.expression(node.AsNonNullExpression().Expression)
	case ast.KindGoExpression:
		callNode := node.AsGoExpression().Expression
		if callNode.Kind != ast.KindCallExpression {
			b.e.fail(node, "go requires a function call")
			return "tsU"
		}
		call := callNode.AsCallExpression()
		callee := b.expression(call.Expression)
		args := []string{}
		spread := false
		for _, arg := range call.Arguments.Nodes {
			if arg.Kind == ast.KindSpreadElement {
				spread = true
			}
		}
		if spread {
			return b.temp("loop.spawn(" + callee + "," + b.spreadArguments(call.Arguments.Nodes) + "...)")
		}
		for i := range call.Arguments.Nodes {
			args = append(args, b.callArgument(call, i))
		}
		return b.temp("loop.spawn(" + callee + ",[]tsValue{" + strings.Join(args, ",") + "}...)")
	case ast.KindAwaitExpression:
		if !b.async {
			b.e.fail(node, "await is supported only inside async functions")
			return "tsU"
		}
		value := b.expression(node.AsAwaitExpression().Expression)
		next := b.block()
		result := b.e.unique("Await")
		b.temps = append(b.temps, result)
		b.emit(fmt.Sprintf("m.await(%s,%d); return", value, next))
		b.current = next
		b.emit(result + " = m.result")
		return result
	case ast.KindArrowFunction, ast.KindFunctionExpression:
		return b.temp(b.function(node))
	case ast.KindNewExpression:
		n := node.AsNewExpression()
		if n.Expression.Kind == ast.KindIdentifier && n.Expression.Text() == "RegExp" && b.e.binding(n.Expression) == nil {
			args := []string{}
			if n.Arguments != nil {
				for _, arg := range n.Arguments.Nodes {
					args = append(args, b.expression(arg))
				}
			}
			return b.temp("tsNewRegExp(" + strings.Join(args, ",") + ")")
		}

		if b.e.classReference(n.Expression) != nil || (n.Expression.Kind == ast.KindIdentifier && b.e.binding(n.Expression) != nil) {
			class := b.expression(n.Expression)
			args := []string{}
			if n.Arguments != nil {
				spread := false
				for _, arg := range n.Arguments.Nodes {
					if arg.Kind == ast.KindSpreadElement {
						spread = true
					}
				}
				if spread {
					args = append(args, b.spreadArguments(n.Arguments.Nodes)+"...")
				} else {
					for i, arg := range n.Arguments.Nodes {
						args = append(args, b.constructorArgument(node, b.e.classReference(n.Expression), i, arg))
					}
				}
			}

			suffix := ""
			if len(args) > 0 {
				suffix = "," + strings.Join(args, ",")
			}
			return b.temp("tsConstruct(" + class + suffix + ")")
		}
		if n.Expression.Kind != ast.KindIdentifier || n.Expression.Text() != "Promise" || b.e.binding(n.Expression) != nil || n.Arguments == nil || len(n.Arguments.Nodes) != 1 {
			b.e.fail(node, "only new Promise(executor) is supported")
			return "tsU"
		}
		executor := b.expression(n.Arguments.Nodes[0])
		return b.temp("loop.constructPromise(" + executor + ")")
	case ast.KindCallExpression:
		call := node.AsCallExpression()
		if call.Expression.Kind == ast.KindImportKeyword {
			target := b.e.targets[node]
			if target == "" {
				b.e.fail(node, "dynamic import requires a resolved literal module")
				return "tsU"
			}
			return b.temp("modules[" + strconv.Quote(target) + "].load(loop)")
		}

		hasSpread := false
		for _, arg := range call.Arguments.Nodes {
			if arg.Kind == ast.KindSpreadElement {
				hasSpread = true
			}
		}
		if call.Expression.Kind == ast.KindIdentifier && !hasSpread {
			if fn := b.e.genericFunctions[b.e.reference(call.Expression)]; fn != nil {
				if value, ok := b.genericCall(node, fn); ok {
					return value
				}
			}
			if fn := b.e.nativeFunctions[b.e.reference(call.Expression)]; fn != nil && b.concrete == nil {
				return b.nativeCall(node, fn)
			}
		}
		if call.Expression.Kind == ast.KindSuperKeyword {
			if !b.constructor || b.declaring == nil || b.declaring.parent == nil {
				b.e.fail(node, "super() requires a derived constructor")
				return "tsU"
			}
			args := []string{}
			if hasSpread {
				args = append(args, b.spreadArguments(call.Arguments.Nodes)+"...")
			} else {
				for i, arg := range call.Arguments.Nodes {
					args = append(args, b.constructorArgument(node, b.declaring.parent, i, arg))
				}
			}

			b.emit("if " + b.receiver + ".properties.initialized {other:=" + b.receiver + ".newBlank();other." + ctorName(b.declaring.parent) + "(" + strings.Join(args, ",") + ");panic(\"Super constructor may only be called once\")}")
			b.emit(b.receiver + "." + ctorName(b.declaring.parent) + "(" + strings.Join(args, ",") + ")")
			b.initializeFields()
			b.initializeParameterProperties()
			return b.receiver
		}
		if call.QuestionDotToken != nil {
			b.e.fail(node, "optional calls are not supported")
		}
		if call.Expression.Kind == ast.KindIdentifier && b.e.imports[b.e.reference(call.Expression)] == "readFile" {
			if len(call.Arguments.Nodes) != 2 || call.Arguments.Nodes[1].Kind != ast.KindStringLiteral || (call.Arguments.Nodes[1].Text() != "utf8" && call.Arguments.Nodes[1].Text() != "utf-8") {
				b.e.fail(node, "imported readFile currently requires an explicit utf8 encoding")
			}
		}
		if call.Expression.Kind == ast.KindPropertyAccessExpression && !hasSpread {
			property := call.Expression.AsPropertyAccessExpression()
			name := property.Name().Text()
			var target string
			if property.Expression.Kind == ast.KindSuperKeyword && b.declaring != nil && b.declaring.parent != nil {
				method := b.declaring.parent.methods[name]
				if method != nil {
					target = b.receiver + "." + implName(b.declaring.parent.ownerOf(method), name)
				}
			} else if c := b.classOf(property.Expression); c != nil && c.methods[name] != nil {
				receiver := b.expression(property.Expression)
				target = b.classTarget(property.Expression, receiver, c) + ".Call" + memberName(name)
				if property.Expression.Kind == ast.KindThisKeyword {
					target = b.receiver + ".Call" + memberName(name)
				}
			}
			if target != "" {
				args := []string{}
				for i := range call.Arguments.Nodes {
					args = append(args, b.callArgument(call, i))
				}
				return b.temp(target + "(" + strings.Join(args, ",") + ")")
			}
		}
		callee := b.expression(call.Expression)
		if hasSpread {
			return b.temp("tsCall(" + callee + "," + b.spreadArguments(call.Arguments.Nodes) + "...)")
		}
		args := []string{}
		for i, arg := range call.Arguments.Nodes {
			if arg.Kind == ast.KindSpreadElement {
				b.e.fail(arg, "spread arguments are not supported")
			}
			args = append(args, b.callArgument(call, i))
		}
		suffix := ""
		if len(args) > 0 {
			suffix = "," + strings.Join(args, ",")
		}
		return b.temp("tsCall(" + callee + suffix + ")")
	case ast.KindPropertyAccessExpression:
		n := node.AsPropertyAccessExpression()
		if n.QuestionDotToken != nil {
			b.e.fail(node, "optional property access is not supported")
		}
		name := node.Name().Text()
		if n.Expression.Kind == ast.KindIdentifier && b.e.binding(n.Expression) == nil {
			if n.Expression.Text() == "Object" {
				return b.temp(b.objectBuiltin(name))
			}
			if n.Expression.Text() == "String" {
				return b.temp("tsFunc(func(args ...tsValue)tsValue{return tsStringStatic(" + strconv.Quote(name) + ",args)})")
			}
			switch n.Expression.Text() + "." + name {
			case "console.log":
				return b.temp("tsFunc(tsConsoleLog)")
			case "Promise.resolve":
				return b.temp("tsFunc(func(args ...tsValue) tsValue {return loop.resolved(tsArg(args,0),false)})")
			case "Promise.all":
				return b.temp("tsFunc(func(args ...tsValue) tsValue {return loop.all(tsArg(args,0))})")
			case "Promise.reject":
				return b.temp("tsFunc(func(args ...tsValue) tsValue {return loop.resolved(tsArg(args,0),true)})")
			}
		}
		receiver := b.expression(n.Expression)
		if value, ok := b.classProperty(n.Expression, receiver, name); ok {
			return b.temp(value)
		}
		return b.temp("tsGet(" + receiver + "," + strconv.Quote(name) + ")")
	case ast.KindElementAccessExpression:
		n := node.AsElementAccessExpression()
		receiver := b.expression(n.Expression)
		index := b.expression(n.ArgumentExpression)
		return b.temp("tsGet(" + receiver + "," + index + ")")
	case ast.KindArrayLiteralExpression:
		return b.arrayLiteral(node.AsArrayLiteralExpression().Elements.Nodes)
	case ast.KindBinaryExpression:
		return b.binary(node)
	case ast.KindConditionalExpression:
		n := node.AsConditionalExpression()
		condition := b.expression(n.Condition)
		if b.direct {
			result := b.temp("tsU")
			b.emit("if tsTruthy(" + condition + ") {")
			value := b.expression(n.WhenTrue)
			b.emit(result + "=" + value)
			b.emit("} else {")
			value = b.expression(n.WhenFalse)
			b.emit(result + "=" + value)
			b.emit("}")
			return result
		}
		yes, no, end := b.block(), b.block(), b.block()
		result := b.e.unique("Conditional")
		b.temps = append(b.temps, result)
		b.emit(fmt.Sprintf("if tsTruthy(%s) {m.pc=%d} else {m.pc=%d}; return", condition, yes, no))
		b.current = yes
		value := b.expression(n.WhenTrue)
		b.emit(result + " = " + value)
		b.jump(end)
		b.current = no
		value = b.expression(n.WhenFalse)
		b.emit(result + " = " + value)
		b.jump(end)
		b.current = end
		return result
	case ast.KindPrefixUnaryExpression:
		n := node.AsPrefixUnaryExpression()
		if n.Operator == ast.KindPlusPlusToken || n.Operator == ast.KindMinusMinusToken {
			return b.update(n.Operand, n.Operator, false)
		}
		operand := b.expression(n.Operand)
		switch n.Operator {
		case ast.KindMinusToken:
			if p := (primitive{kind: b.tempType(operand)}); p.numeric() {
				return b.typedTemp("-"+operand, p.goType())
			}
			return b.temp("-tsNumber(" + operand + ")")
		case ast.KindPlusToken:
			if (primitive{kind: b.tempType(operand)}).numeric() {
				return operand
			}
			return b.temp("tsNumber(" + operand + ")")
		case ast.KindTildeToken:
			if b.tempType(operand) == "float64" {
				return b.typedTemp("tsNumberBitwiseNot("+operand+")", "float64")
			}
			return b.temp("tsBitwiseNot(" + operand + ")")
		case ast.KindExclamationToken:
			return b.temp("!tsTruthy(" + operand + ")")
		}
	case ast.KindPostfixUnaryExpression:
		n := node.AsPostfixUnaryExpression()
		return b.update(n.Operand, n.Operator, true)
	case ast.KindVoidExpression:
		b.expression(node.AsVoidExpression().Expression)
		return "tsU"
	}
	b.e.fail(node, "unsupported expression "+node.Kind.String())
	return "tsU"
}

// lvalue snapshots the receiver/key before an RHS or await can mutate them.
func (b *machineBuilder) lvalue(node *ast.Node) (string, func(string) string) {
	if node.Kind == ast.KindIdentifier {
		if cell := b.e.binding(node); cell != nil {
			pointer := b.typedTemp(cell.name, cell.pointerType())
			return pointer + ".get()", func(value string) string {
				method := "set"
				if cell.primitive.kind != "" && b.tempType(value) == cell.primitive.goType() {
					method = "setNative"
				}
				return pointer + "." + method + "(" + value + ")"
			}
		}
	}
	if node.Kind == ast.KindPropertyAccessExpression {
		n := node.AsPropertyAccessExpression()
		if n.Expression.Kind == ast.KindSuperKeyword {
			if b.declaring != nil && b.declaring.static && b.declaring.parent != nil && b.declaring.parent.fields[n.Name().Text()] != nil {
				read, _ := b.classProperty(n.Expression, b.receiver, n.Name().Text())
				return read, func(value string) string {
					return b.receiver + ".Set" + memberName(n.Name().Text()) + "(" + value + ")"
				}
			}
			b.e.fail(node, "super property assignment requires accessor lowering, which is not supported yet")
			return "tsU", func(string) string { return "tsU" }
		}
		receiver := b.expression(n.Expression)
		name := n.Name().Text()
		c := b.classOf(n.Expression)
		if c != nil && c.fields[name] != nil {
			target := b.classTarget(n.Expression, receiver, c)
			if n.Expression.Kind == ast.KindThisKeyword {
				target = b.receiver
			}
			id := memberName(name)
			return target + ".Get" + id + "()", func(value string) string {
				shape := b.e.primitive(c.fields[name])
				if n.Expression.Kind == ast.KindThisKeyword && !c.static && shape.kind != "" && shape.nulls == 0 && b.tempType(value) == shape.goType() {
					return "func() tsValue {" + target + "." + id + "=" + value + ";" + target + ".has" + id + "=true;return " + value + "}()"
				}
				return target + ".Set" + id + "(" + value + ")"
			}
		}
		return "tsGet(" + receiver + "," + strconv.Quote(name) + ")", func(value string) string { return "tsSet(" + receiver + "," + strconv.Quote(name) + "," + value + ")" }
	}
	if node.Kind == ast.KindElementAccessExpression {
		n := node.AsElementAccessExpression()
		receiver := b.expression(n.Expression)
		key := b.expression(n.ArgumentExpression)
		return "tsGet(" + receiver + "," + key + ")", func(value string) string { return "tsSet(" + receiver + "," + key + "," + value + ")" }
	}
	b.e.fail(node, "unsupported assignment target")
	return "tsU", func(string) string { return "tsU" }
}
func (b *machineBuilder) update(node *ast.Node, operator ast.Kind, postfix bool) string {
	read, write := b.lvalue(node)
	old := b.temp(read)
	op := "+"
	if operator == ast.KindMinusMinusToken {
		op = "-"
	}
	operand := "tsNumber(" + old + ")" + op + "1"
	if p := b.e.primitive(node); p.numeric() {
		one := b.typedTemp(p.goType()+"(1)", p.goType())
		operand = "tsBinary(" + strconv.Quote(op) + "," + old + "," + one + ")"
	}
	value := b.temp(write(operand))
	if postfix {
		return old
	}
	return value
}
func (b *machineBuilder) binary(node *ast.Node) string {
	n := node.AsBinaryExpression()
	operator := n.OperatorToken.Kind
	if operator == ast.KindAmpersandAmpersandEqualsToken || operator == ast.KindBarBarEqualsToken || operator == ast.KindQuestionQuestionEqualsToken {
		return b.logicalAssignment(n)
	}
	if operator == ast.KindEqualsToken || compoundOperator(operator) != "" {
		read, write := b.lvalue(n.Left)
		old := ""
		if operator != ast.KindEqualsToken {
			var cell *binding
			if n.Left.Kind == ast.KindIdentifier {
				cell = b.e.binding(n.Left)
			}
			if cell != nil && cell.primitive.goType() == "float64" && cell.primitive.nulls == 0 {
				old = b.typedTemp(strings.TrimSuffix(read, ".get()")+".read()", "float64")
			} else {
				old = b.temp(read)
			}
		}
		value := ""
		var assignment *binding
		if n.Left.Kind == ast.KindIdentifier {
			assignment = b.e.binding(n.Left)
		}
		if cell := assignment; cell != nil && cell.primitive.numeric() && cell.primitive.nulls == 0 {
			value = b.numericExpression(n.Right, cell.primitive)
		} else {
			value = b.expression(n.Right)
		}
		if old != "" {
			op := compoundOperator(operator)
			if op == "**" && b.tempType(old) == "float64" && b.tempType(value) == "float64" {
				value = b.typedTemp("tsPow("+old+","+value+")", "float64")
			} else if isBitwiseOperator(op) && b.tempType(old) == "float64" && b.tempType(value) == "float64" {
				value = b.typedTemp("tsNumberBitwise("+strconv.Quote(op)+","+old+","+value+")", "float64")
			} else {
				value = b.temp("tsBinary(" + strconv.Quote(op) + "," + old + "," + value + ")")
			}
		}
		if n.Left.Kind == ast.KindIdentifier {
			if cell := b.e.binding(n.Left); cell != nil && cell.primitive.kind != "" && b.tempType(value) == cell.primitive.goType() {
				return b.typedTemp(write(value), cell.primitive.goType())
			}
		}
		return b.temp(write(value))
	}
	left := ""
	if p := b.e.primitive(n.Right); p.numeric() && p.kind != "number" && n.Left.Kind == ast.KindNumericLiteral {
		left = b.numericExpression(n.Left, p)
	} else {
		left = b.expression(n.Left)
	}
	if operator == ast.KindCommaToken {
		return b.expression(n.Right)
	}
	if operator == ast.KindAmpersandAmpersandToken || operator == ast.KindBarBarToken || operator == ast.KindQuestionQuestionToken {
		result := b.temp(left)
		condition := "tsTruthy(" + left + ")"
		if operator == ast.KindBarBarToken {
			condition = "!" + condition
		}
		if operator == ast.KindQuestionQuestionToken {
			condition = "tsNullish(" + left + ")"
		}
		if b.direct {
			b.emit("if " + condition + " {")
			value := b.expression(n.Right)
			b.emit(result + "=" + value)
			b.emit("}")
			return result
		}
		right, end := b.block(), b.block()
		b.emit(fmt.Sprintf("if %s {m.pc=%d} else {m.pc=%d}; return", condition, right, end))
		b.current = right
		value := b.expression(n.Right)
		b.emit(result + " = " + value)
		b.jump(end)
		b.current = end
		return result
	}
	right := ""
	if p := (primitive{kind: b.tempType(left)}); p.numeric() && p.kind != "float64" && n.Right.Kind == ast.KindNumericLiteral {
		right = b.numericExpression(n.Right, p)
	} else {
		right = b.expression(n.Right)
	}
	if operator == ast.KindInstanceOfKeyword {
		return b.temp("tsInstanceOf(" + left + "," + right + ")")
	}
	operators := map[ast.Kind]string{ast.KindPlusToken: "+", ast.KindMinusToken: "-", ast.KindAsteriskToken: "*", ast.KindSlashToken: "/", ast.KindPercentToken: "%", ast.KindLessThanToken: "<", ast.KindLessThanEqualsToken: "<=", ast.KindGreaterThanToken: ">", ast.KindGreaterThanEqualsToken: ">=", ast.KindEqualsEqualsEqualsToken: "===", ast.KindExclamationEqualsEqualsToken: "!==", ast.KindEqualsEqualsToken: "==", ast.KindExclamationEqualsToken: "!=", ast.KindAsteriskAsteriskToken: "**", ast.KindAmpersandToken: "&", ast.KindBarToken: "|", ast.KindCaretToken: "^", ast.KindLessThanLessThanToken: "<<", ast.KindGreaterThanGreaterThanToken: ">>", ast.KindGreaterThanGreaterThanGreaterThanToken: ">>>", ast.KindInKeyword: "in"}
	op, ok := operators[operator]
	if !ok {
		b.e.fail(node, "unsupported binary operator "+operator.String())
		return "tsU"
	}
	// Operands have already been saved in left-to-right order. A dynamic value
	// compared with a typed primitive uses the same policy as typed assignments.
	switch op {
	case "==", "!=", "===", "!==", "<", "<=", ">", ">=":
		leftPlan, rightPlan := b.e.primitive(n.Left), b.e.primitive(n.Right)
		if leftPlan.kind != "" && rightPlan.kind == "" && !isNullishComparisonOperand(n.Right) {
			right = b.comparisonBoundary(right, leftPlan)
		} else if rightPlan.kind != "" && leftPlan.kind == "" && !isNullishComparisonOperand(n.Left) {
			left = b.comparisonBoundary(left, rightPlan)
		}
	}
	if op == "**" && b.tempType(left) == "float64" && b.tempType(right) == "float64" {
		return b.typedTemp("tsPow("+left+","+right+")", "float64")
	}
	if isBitwiseOperator(op) && b.tempType(left) == "float64" && b.tempType(right) == "float64" {
		return b.typedTemp("tsNumberBitwise("+strconv.Quote(op)+","+left+","+right+")", "float64")
	}
	if op == "in" || op == "&" || op == "|" || op == "^" || op == "<<" || op == ">>" || op == ">>>" || op == "**" || op == "==" || op == "!=" {
		return b.temp("tsBinary(" + strconv.Quote(op) + "," + left + "," + right + ")")
	}
	if b.tempType(left) == "bool" && b.tempType(right) == "bool" && (op == "===" || op == "!==") {
		goop := "=="
		if op == "!==" {
			goop = "!="
		}
		return b.typedTemp(left+goop+right, "bool")
	}
	if b.tempType(left) == "*tsString" && b.tempType(right) == "*tsString" {
		switch op {
		case "+":
			return b.typedTemp("tsStringConcat("+left+","+right+")", "*tsString")
		case "===":
			return b.typedTemp("tsStringEqual("+left+","+right+")", "bool")
		case "!==":
			return b.typedTemp("!tsStringEqual("+left+","+right+")", "bool")
		case "<", "<=", ">", ">=":
			return b.typedTemp("tsStringCompare("+left+","+right+")"+op+"0", "bool")
		}
	}
	numberLeft, numberRight := b.tempType(left), b.tempType(right)
	if numberLeft != numberRight && (primitive{kind: numberLeft}).numeric() && (primitive{kind: numberRight}).numeric() && (op == "+" || op == "-" || op == "*" || op == "/" || op == "%") {
		target := numericPromotion(numberLeft, numberRight)
		left = b.comparisonBoundary(left, primitive{kind: target})
		right = b.comparisonBoundary(right, primitive{kind: target})
		numberLeft = target
		numberRight = target
	}

	if numberLeft == numberRight && numberLeft != "float64" && (primitive{kind: numberLeft}).numeric() {
		goop := op
		kind := numberLeft
		switch op {
		case "===":
			goop = "=="
			kind = "bool"
		case "!==":
			goop = "!="
			kind = "bool"
		case "<", "<=", ">", ">=":
			kind = "bool"
		}
		if op == "%" && numberLeft == "float32" {
			return b.typedTemp("float32(math.Mod(float64("+left+"),float64("+right+")))", kind)
		}
		return b.typedTemp(left+goop+right, kind)
	}

	if (numberLeft == "float64" || numberLeft == "tsOptional[float64]") && (numberRight == "float64" || numberRight == "tsOptional[float64]") {
		if (op == "===" || op == "!==") && (numberLeft != "float64" || numberRight != "float64") {
			return b.temp("tsBinary(" + strconv.Quote(op) + "," + left + "," + right + ")")
		}
		if numberLeft != "float64" {
			left = "tsNullableNumber(" + left + ")"
		}
		if numberRight != "float64" {
			right = "tsNullableNumber(" + right + ")"
		}
		kind := "float64"
		switch op {
		case "<", "<=", ">", ">=", "===", "!==":
			kind = "bool"
		}
		goop := op
		if op == "===" {
			goop = "=="
		}
		if op == "!==" {
			goop = "!="
		}
		if op == "%" {
			return b.typedTemp("math.Mod("+left+","+right+")", kind)
		}
		return b.typedTemp(left+goop+right, kind)
	}
	return b.temp("tsBinary(" + strconv.Quote(op) + "," + left + "," + right + ")")
}

func isNullishComparisonOperand(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindNullKeyword:
		return true
	case ast.KindIdentifier:
		return node.Text() == "undefined"
	case ast.KindParenthesizedExpression:
		return isNullishComparisonOperand(node.AsParenthesizedExpression().Expression)
	}
	return false
}
func (b *machineBuilder) comparisonBoundary(value string, p primitive) string {
	checked := fmt.Sprintf("tsBoundary(%s,%q,%d,%t)", value, p.kind, p.nulls, b.e.coerce)
	if p.nulls != 0 {
		return b.typedTemp("tsOptionalFrom["+p.goType()+"]("+checked+")", "tsOptional["+p.goType()+"]")
	}
	return b.typedTemp("tsNative["+p.goType()+"]("+checked+")", p.goType())
}

func (b *machineBuilder) numericAssertion(expression, typeNode *ast.Node) string {
	p := b.e.annotationPrimitive(typeNode)
	if p.numeric() {
		return b.numericExpression(expression, p)
	}
	return b.expression(expression)
}
func (b *machineBuilder) numericExpression(node *ast.Node, p primitive) string {
	if node.Kind == ast.KindParenthesizedExpression {
		return b.numericExpression(node.Expression(), p)
	}

	if p.kind != "number" && p.kind != "float64" {
		if c, known := b.nativeConstant(node); known {
			if p.kind == "float32" {
				v, _ := constant.Float64Val(constant.ToFloat(c))
				rounded := float32(v)
				if !math.IsInf(float64(rounded), 0) {
					return b.typedTemp("float32("+strconv.FormatFloat(float64(rounded), 'g', -1, 32)+")", "float32")
				}
			}
			if p.kind != "float32" {
				bits := 64
				switch p.kind {
				case "int", "uint":
					bits = strconv.IntSize
				case "int8", "uint8":
					bits = 8
				case "int16", "uint16":
					bits = 16
				case "int32", "uint32":
					bits = 32
				}
				integer := constant.ToInt(c)
				if integer.Kind() == constant.Int {
					if strings.HasPrefix(p.kind, "uint") {
						v, ok := constant.Uint64Val(integer)
						if ok && (bits == 64 || v < uint64(1)<<bits) {
							return b.typedTemp(p.goType()+"("+strconv.FormatUint(v, 10)+")", p.goType())
						}
					} else {
						v, ok := constant.Int64Val(integer)
						if ok && (bits == 64 || (v >= -(int64(1)<<(bits-1)) && v < int64(1)<<(bits-1))) {
							return b.typedTemp(p.goType()+"("+strconv.FormatInt(v, 10)+")", p.goType())
						}
					}
				}
			}
			b.e.fail(node, "numeric constant cannot be represented as "+p.kind)
			return "tsU"
		}
	}

	value := b.expression(node)
	if b.tempType(value) == p.goType() {
		return value
	}
	return b.comparisonBoundary(value, p)
}

func (b *machineBuilder) callArgument(call *ast.CallExpression, index int) string {
	var declaration *ast.Node
	if call.Expression.Kind == ast.KindIdentifier {
		declaration = b.e.reference(call.Expression)
	}
	if call.Expression.Kind == ast.KindPropertyAccessExpression {
		property := call.Expression.AsPropertyAccessExpression()
		if c := b.classOf(property.Expression); c != nil {
			declaration = c.methods[property.Name().Text()]
		}
	}
	if declaration != nil && declaration.Kind == ast.KindVariableDeclaration {
		declaration = declaration.AsVariableDeclaration().Initializer
	}
	if declaration != nil && ast.IsFunctionLike(declaration) {
		p := b.argumentPrimitive(declaration, index)
		if p.numeric() && p.nulls == 0 {
			return b.numericExpression(call.Arguments.Nodes[index], p)
		}
	}
	return b.expression(call.Arguments.Nodes[index])
}

func numericPromotion(a, b string) string { return core.NativeNumericPromotion(a, b) }

func (b *machineBuilder) nativeConstant(node *ast.Node) (constant.Value, bool) {
	switch node.Kind {
	case ast.KindNumericLiteral:
		text := strings.TrimSpace(b.e.file.Text()[scanner.SkipTrivia(b.e.file.Text(), node.Pos()):node.End()])
		text = strings.ReplaceAll(text, "_", "")
		kind := token.INT
		if strings.ContainsAny(text, ".eEpP") && !strings.HasPrefix(strings.ToLower(text), "0x") {
			kind = token.FLOAT
		}
		value := constant.MakeFromLiteral(text, kind, 0)
		return value, value.Kind() != constant.Unknown
	case ast.KindParenthesizedExpression:
		return b.nativeConstant(node.Expression())
	case ast.KindPrefixUnaryExpression:
		unary := node.AsPrefixUnaryExpression()
		if unary.Operator != ast.KindPlusToken && unary.Operator != ast.KindMinusToken {
			return nil, false
		}
		value, known := b.nativeConstant(unary.Operand)
		if !known {
			return nil, false
		}
		op := token.ADD
		if unary.Operator == ast.KindMinusToken {
			op = token.SUB
		}
		return constant.UnaryOp(op, value, 0), true
	case ast.KindBinaryExpression:
		binary := node.AsBinaryExpression()
		op := map[ast.Kind]token.Token{ast.KindPlusToken: token.ADD, ast.KindMinusToken: token.SUB, ast.KindAsteriskToken: token.MUL, ast.KindSlashToken: token.QUO, ast.KindPercentToken: token.REM}[binary.OperatorToken.Kind]
		if op == token.ILLEGAL {
			return nil, false
		}
		left, lk := b.nativeConstant(binary.Left)
		right, rk := b.nativeConstant(binary.Right)
		if !lk || !rk {
			return nil, false
		}
		if (op == token.QUO || op == token.REM) && constant.Sign(right) == 0 {
			return nil, false
		}
		if op == token.QUO {
			left = constant.ToFloat(left)
			right = constant.ToFloat(right)
		}
		if op == token.REM {
			left = constant.ToInt(left)
			right = constant.ToInt(right)
			if left.Kind() != constant.Int || right.Kind() != constant.Int {
				return nil, false
			}
		}
		return constant.BinaryOp(left, op, right), true
	}
	return nil, false
}

func compoundOperator(kind ast.Kind) string {
	return map[ast.Kind]string{ast.KindPlusEqualsToken: "+", ast.KindMinusEqualsToken: "-", ast.KindAsteriskEqualsToken: "*", ast.KindSlashEqualsToken: "/", ast.KindPercentEqualsToken: "%", ast.KindAsteriskAsteriskEqualsToken: "**", ast.KindAmpersandEqualsToken: "&", ast.KindBarEqualsToken: "|", ast.KindCaretEqualsToken: "^", ast.KindLessThanLessThanEqualsToken: "<<", ast.KindGreaterThanGreaterThanEqualsToken: ">>", ast.KindGreaterThanGreaterThanGreaterThanEqualsToken: ">>>"}[kind]
}
func (b *machineBuilder) logicalAssignment(node *ast.BinaryExpression) string {
	read, write := b.lvalue(node.Left)
	result := b.temp(read)
	condition := "tsTruthy(" + result + ")"
	switch node.OperatorToken.Kind {
	case ast.KindBarBarEqualsToken:
		condition = "!" + condition
	case ast.KindQuestionQuestionEqualsToken:
		condition = "tsNullish(" + result + ")"
	}
	if b.direct {
		b.emit("if " + condition + " {")
		value := b.expression(node.Right)
		b.emit(result + "=" + write(value))
		b.emit("}")
		return result
	}
	yes, end := b.block(), b.block()
	b.emit(fmt.Sprintf("if %s{m.pc=%d}else{m.pc=%d};return", condition, yes, end))
	b.current = yes
	value := b.expression(node.Right)
	b.emit(result + "=" + write(value))
	b.jump(end)
	b.current = end
	return result
}

func isBitwiseOperator(op string) bool {
	switch op {
	case "&", "|", "^", "<<", ">>", ">>>":
		return true
	}
	return false
}
