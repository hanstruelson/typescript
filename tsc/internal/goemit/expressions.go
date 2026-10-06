package goemit

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

func (b *machineBuilder) identifier(node *ast.Node) string {
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
		return b.temp("true")
	case ast.KindFalseKeyword:
		return b.temp("false")
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
		return b.expression(node.AsAsExpression().Expression)
	case ast.KindTypeAssertionExpression:
		return b.expression(node.AsTypeAssertion().Expression)
	case ast.KindNonNullExpression:
		return b.expression(node.AsNonNullExpression().Expression)
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
				for _, arg := range n.Arguments.Nodes {
					args = append(args, b.expression(arg))
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
		return b.temp("loop.construct(" + executor + ")")
	case ast.KindCallExpression:
		call := node.AsCallExpression()
		hasSpread := false
		for _, arg := range call.Arguments.Nodes {
			if arg.Kind == ast.KindSpreadElement {
				hasSpread = true
			}
		}
		if call.Expression.Kind == ast.KindIdentifier && b.concrete == nil && !hasSpread {
			if fn := b.e.nativeFunctions[b.e.reference(call.Expression)]; fn != nil {
				return b.nativeCall(node, fn)
			}
		}
		if call.Expression.Kind == ast.KindSuperKeyword {
			if !b.constructor || b.declaring == nil || b.declaring.parent == nil {
				b.e.fail(node, "super() requires a derived constructor")
				return "tsU"
			}
			args := []string{}
			for _, arg := range call.Arguments.Nodes {
				args = append(args, b.expression(arg))
			}
			b.emit("if " + b.receiver + ".properties.initialized {other:=" + b.receiver + ".newBlank();other." + ctorName(b.declaring.parent) + "(" + strings.Join(args, ",") + ");panic(\"Super constructor may only be called once\")}")
			b.emit(b.receiver + "." + ctorName(b.declaring.parent) + "(" + strings.Join(args, ",") + ")")
			b.initializeFields()
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
				for _, arg := range call.Arguments.Nodes {
					args = append(args, b.expression(arg))
				}
				return b.temp(target + "(" + strings.Join(args, ",") + ")")
			}
		}
		callee := b.expression(call.Expression)
		if hasSpread {
			return b.temp("tsCall(" + callee + "," + b.spreadArguments(call.Arguments.Nodes) + "...)")
		}
		args := []string{}
		for _, arg := range call.Arguments.Nodes {
			if arg.Kind == ast.KindSpreadElement {
				b.e.fail(arg, "spread arguments are not supported")
			}
			args = append(args, b.expression(arg))
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
		return b.temp("&tsArray{values:" + b.spreadArguments(node.AsArrayLiteralExpression().Elements.Nodes) + "}")
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
			return b.temp("-tsNumber(" + operand + ")")
		case ast.KindPlusToken:
			return b.temp("tsNumber(" + operand + ")")
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
	value := b.temp(write("tsNumber(" + old + ")" + op + "1"))
	if postfix {
		return old
	}
	return value
}
func (b *machineBuilder) binary(node *ast.Node) string {
	n := node.AsBinaryExpression()
	operator := n.OperatorToken.Kind
	if operator == ast.KindEqualsToken || operator == ast.KindPlusEqualsToken || operator == ast.KindMinusEqualsToken || operator == ast.KindAsteriskEqualsToken {
		read, write := b.lvalue(n.Left)
		old := ""
		if operator != ast.KindEqualsToken {
			old = b.temp(read)
		}
		value := b.expression(n.Right)
		if old != "" {
			op := map[ast.Kind]string{ast.KindPlusEqualsToken: "+", ast.KindMinusEqualsToken: "-", ast.KindAsteriskEqualsToken: "*"}[operator]
			value = b.temp("tsBinary(" + strconv.Quote(op) + "," + old + "," + value + ")")
		}
		if n.Left.Kind == ast.KindIdentifier {
			if cell := b.e.binding(n.Left); cell != nil && cell.primitive.kind != "" && b.tempType(value) == cell.primitive.goType() {
				return b.typedTemp(write(value), cell.primitive.goType())
			}
		}
		return b.temp(write(value))
	}
	left := b.expression(n.Left)
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
	right := b.expression(n.Right)
	if operator == ast.KindInstanceOfKeyword {
		return b.temp("tsInstanceOf(" + left + "," + right + ")")
	}
	operators := map[ast.Kind]string{ast.KindPlusToken: "+", ast.KindMinusToken: "-", ast.KindAsteriskToken: "*", ast.KindSlashToken: "/", ast.KindPercentToken: "%", ast.KindLessThanToken: "<", ast.KindLessThanEqualsToken: "<=", ast.KindGreaterThanToken: ">", ast.KindGreaterThanEqualsToken: ">=", ast.KindEqualsEqualsEqualsToken: "===", ast.KindExclamationEqualsEqualsToken: "!=="}
	op, ok := operators[operator]
	if !ok {
		b.e.fail(node, "unsupported binary operator "+operator.String())
		return "tsU"
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
