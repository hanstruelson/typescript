package goemit

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

type genericFunction struct {
	machine                    bool
	captures                   []*binding
	instances                  map[string]string
	text                       string
	node                       *ast.Node
	parameters, typeParameters []string
	result                     string
	constraints                map[string]string
}

type genericBody struct {
	e         *emitter
	fn        *genericFunction
	names     map[*ast.Node]string
	types     map[*ast.Node]string
	typeNames map[string]string
	next      int
}

// Each checked representation gets one concrete implementation. Only bodies
// proven independent of their lexical environment are lifted; ordinary closures
// keep their existing ABI until closure specialization is available.
func (e *emitter) planGenericFunctions() {
	e.genericFunctions = map[*ast.Node]*genericFunction{}
	for _, node := range e.file.Statements.Nodes {
		if node.Kind != ast.KindFunctionDeclaration || node.Body() == nil || ast.HasSyntacticModifier(node, ast.ModifierFlagsAsync) || len(node.TypeParameters()) == 0 {
			continue
		}
		f := &genericFunction{node: node, constraints: map[string]string{}}
		g := &genericBody{e: e, fn: f, names: map[*ast.Node]string{}, types: map[*ast.Node]string{}, typeNames: map[string]string{}}
		valid := true
		for i, param := range node.TypeParameters() {
			name := fmt.Sprintf("G%d", i)
			g.typeNames[param.Name().Text()] = name
			f.typeParameters = append(f.typeParameters, name)
			if c := param.AsTypeParameterDeclaration().Constraint; c != nil && c.Kind != ast.KindAnyKeyword && c.Kind != ast.KindUnknownKeyword {
				p := e.annotationPrimitive(c)
				if p.kind == "" || p.nulls != 0 {
					valid = false
					break
				}
				f.constraints[name] = p.goType()
			}
		}
		for i, param := range node.Parameters() {
			decl := param.AsParameterDeclaration()
			typ := g.annotation(param.Type())
			if param.Name() == nil || param.Name().Kind != ast.KindIdentifier || typ == "" || decl.Initializer != nil || decl.DotDotDotToken != nil || param.QuestionToken() != nil {
				valid = false
				break
			}
			name := fmt.Sprintf("a%d", i)
			g.names[param] = name
			g.types[param] = typ
			f.parameters = append(f.parameters, typ)
		}
		f.result = g.annotation(node.Type())
		if f.result == "" {
			valid = false
		}
		body, ok := g.statement(node.Body())
		if !valid || !ok {
			if !directClassBody(node) {
				continue
			}
			eligible := true
			for _, param := range node.Parameters() {
				d := param.AsParameterDeclaration()
				if param.Name() == nil || param.Name().Kind != ast.KindIdentifier || d.Initializer != nil || d.DotDotDotToken != nil || param.QuestionToken() != nil {
					eligible = false
				}
			}
			if !eligible {
				continue
			}
			f.machine = true
			f.captures = nil
			f.parameters = make([]string, len(node.Parameters()))
			f.constraints = map[string]string{}
			f.typeParameters = make([]string, len(node.TypeParameters()))
			for i := range f.typeParameters {
				f.typeParameters[i] = fmt.Sprintf("G%d", i)
			}
			var free func(*ast.Node)
			seen := map[*binding]bool{}
			free = func(n *ast.Node) {
				if n.Kind == ast.KindIdentifier {
					if cell := e.binding(n); cell != nil && cell.owner == e.file.AsNode() && !seen[cell] {
						seen[cell] = true
						f.captures = append(f.captures, cell)
					}
				}
				n.ForEachChild(func(child *ast.Node) bool { free(child); return false })
			}
			free(node.Body())
		}
		e.genericFunctions[node] = f
		f.text = body
		f.instances = map[string]string{}
	}
	var visit func(*ast.Node)
	visit = func(node *ast.Node) {
		if node.Kind == ast.KindBinaryExpression {
			n := node.AsBinaryExpression()
			if ast.IsAssignmentOperator(n.OperatorToken.Kind) && n.Left.Kind == ast.KindIdentifier {
				delete(e.genericFunctions, e.reference(n.Left))
			}
		}
		node.ForEachChild(func(child *ast.Node) bool { visit(child); return false })
	}
	visit(e.file.AsNode())
}
func (g *genericBody) annotation(node *ast.Node) string {
	if node == nil {
		return ""
	}
	if node.Kind == ast.KindTypeReference && node.AsTypeReferenceNode().TypeName.Kind == ast.KindIdentifier {
		if name := g.typeNames[node.AsTypeReferenceNode().TypeName.Text()]; name != "" {
			return name
		}
	}
	if node.Kind == ast.KindAnyKeyword || node.Kind == ast.KindUnknownKeyword {
		return "tsValue"
	}
	p := g.e.annotationPrimitive(node)
	if p.kind != "" && p.nulls == 0 {
		return p.goType()
	}
	return ""
}
func (g *genericBody) numeric(typ string) bool {
	if native := g.fn.constraints[typ]; native != "" {
		typ = native
	}
	switch typ {
	case "float64", "float32", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	}
	return false
}
func (g *genericBody) expression(node *ast.Node) (string, string, bool) {
	if node == nil {
		return "", "", false
	}
	switch node.Kind {
	case ast.KindIdentifier:
		decl := g.e.reference(node)
		name, ok := g.names[decl]
		if ok {
			return name, g.types[decl], true
		}
		cell := g.e.bindings[decl]
		if cell == nil || cell.owner == g.fn.node || cell.primitive.nulls != 0 || cell.maybeUndefined {
			return "", "", false
		}
		capture := ""
		for i, existing := range g.fn.captures {
			if existing == cell {
				capture = fmt.Sprintf("environment%d", i)
				break
			}
		}
		if capture == "" {
			capture = fmt.Sprintf("environment%d", len(g.fn.captures))
			g.fn.captures = append(g.fn.captures, cell)
		}
		if cell.primitive.kind != "" {
			return capture + ".read()", cell.primitive.goType(), true
		}
		return capture + ".get()", "tsValue", true
	case ast.KindParenthesizedExpression:
		value, typ, ok := g.expression(node.Expression())
		return "(" + value + ")", typ, ok
	case ast.KindAsExpression, ast.KindTypeAssertionExpression:
		value, typ, ok := g.expression(node.Expression())
		target := g.annotation(node.Type())
		if !ok || target != typ {
			return "", "", false
		}
		return value, typ, true
	case ast.KindTrueKeyword:
		return "true", "bool", true
	case ast.KindFalseKeyword:
		return "false", "bool", true
	case ast.KindNumericLiteral:
		return node.Text(), "constant", true
	case ast.KindBinaryExpression:
		n := node.AsBinaryExpression()
		left, lt, lok := g.expression(n.Left)
		right, rt, rok := g.expression(n.Right)
		if !lok || !rok {
			return "", "", false
		}
		typ := lt
		if typ == "constant" {
			typ = rt
		}
		if lt != rt && lt != "constant" && rt != "constant" {
			if g.fn.constraints[lt] == rt {
				right = lt + "(" + right + ")"
			} else if g.fn.constraints[rt] == lt {
				left = rt + "(" + left + ")"
				typ = rt
			} else {
				return "", "", false
			}
		}
		// Let the checked body emitter handle native-width constants and their
		// range checks rather than relying on Go's untyped-constant rules.
		if lt == "constant" || rt == "constant" {
			native := typ
			if bound := g.fn.constraints[typ]; bound != "" {
				native = bound
			}
			if native != "float64" {
				return "", "", false
			}
		}
		ops := map[ast.Kind]string{ast.KindPlusToken: "+", ast.KindMinusToken: "-", ast.KindAsteriskToken: "*", ast.KindSlashToken: "/", ast.KindLessThanToken: "<", ast.KindLessThanEqualsToken: "<=", ast.KindGreaterThanToken: ">", ast.KindGreaterThanEqualsToken: ">=", ast.KindEqualsEqualsEqualsToken: "==", ast.KindExclamationEqualsEqualsToken: "!="}
		op := ops[n.OperatorToken.Kind]
		if op == "" || !g.numeric(typ) {
			return "", "", false
		}
		result := typ
		if op == "<" || op == "<=" || op == ">" || op == ">=" || op == "==" || op == "!=" {
			result = "bool"
		}
		return "(" + left + op + right + ")", result, true
	}
	return "", "", false
}
func (g *genericBody) statement(node *ast.Node) (string, bool) {
	if node == nil {
		return "", true
	}
	switch node.Kind {
	case ast.KindBlock:
		lines := []string{}
		for _, stmt := range node.AsBlock().Statements.Nodes {
			line, ok := g.statement(stmt)
			if !ok {
				return "", false
			}
			lines = append(lines, line)
		}
		return strings.Join(lines, "\n"), true
	case ast.KindReturnStatement:
		value, typ, ok := g.expression(node.AsReturnStatement().Expression)
		if !ok {
			return "", false
		}
		if typ != g.fn.result {
			if g.fn.constraints[typ] == g.fn.result {
				value = "(" + g.fn.result + ")(" + value + ")"
			} else {
				return "", false
			}
		}
		return "return " + value, true
	case ast.KindIfStatement:
		n := node.AsIfStatement()
		condition, typ, ok := g.expression(n.Expression)
		if !ok || typ != "bool" {
			return "", false
		}
		yes, ok := g.statement(n.ThenStatement)
		if !ok {
			return "", false
		}
		no, ok := g.statement(n.ElseStatement)
		if !ok {
			return "", false
		}
		result := "if " + condition + " {\n" + yes + "\n}"
		if n.ElseStatement != nil {
			result += " else {\n" + no + "\n}"
		}
		return result, true
	case ast.KindVariableStatement:
		lines := []string{}
		for _, decl := range node.AsVariableStatement().DeclarationList.AsVariableDeclarationList().Declarations.Nodes {
			if decl.Name().Kind != ast.KindIdentifier {
				return "", false
			}
			value, typ, ok := g.expression(decl.AsVariableDeclaration().Initializer)
			if !ok || typ == "constant" {
				return "", false
			}
			if decl.Type() != nil && g.annotation(decl.Type()) != typ {
				return "", false
			}
			name := fmt.Sprintf("v%d", g.next)
			g.next++
			g.names[decl] = name
			g.types[decl] = typ
			lines = append(lines, name+":="+value+";_="+name)
		}
		return strings.Join(lines, "\n"), true
	}
	return "", false
}

func (b *machineBuilder) genericCall(node *ast.Node, f *genericFunction) (string, bool) {
	resolver, ok := b.e.resolver.(interface {
		GetEmitGenericCallTypes(*ast.Node, []*ast.Node) (*ast.Node, []string, []string, string, string)
	})
	if !ok {
		return "", false
	}
	declaration, types, parameters, result, semanticKey := resolver.GetEmitGenericCallTypes(node, b.e.specializationCalls)
	if declaration != f.node || len(types) != len(f.typeParameters) || len(parameters) != len(f.parameters) {
		return "", false
	}
	for i, typ := range types {
		constraint := f.constraints[f.typeParameters[i]]
		if constraint != "" && constraint != typ {
			return "", false
		}
	}
	// Missing parameters preserve JavaScript's undefined behavior through the
	// ordinary wrapper. Extra arguments are still evaluated for their effects.
	call := node.AsCallExpression()
	if len(call.Arguments.Nodes) < len(parameters) {
		return "", false
	}
	key := strings.Join(types, ",") + ";" + strings.Join(parameters, ",") + ";" + result
	if f.machine {
		key += ";" + semanticKey
	}
	if f.instances[key] == "" && (len(f.instances) >= 64 || len(b.e.specializationCalls) >= 16) {
		return "", false
	}
	args := []string{}
	for _, cell := range f.captures {
		args = append(args, cell.name)
	}
	for i, arg := range call.Arguments.Nodes {
		if i >= len(parameters) {
			b.expression(arg)
			continue
		}
		kind := parameters[i]
		p := goPrimitive(kind)
		value := ""
		if p.numeric() {
			value = b.numericExpression(arg, p)
		} else {
			value = b.expression(arg)
		}
		if kind != "tsValue" && b.tempType(value) != kind {
			value = fmt.Sprintf("tsNative[%s](tsBoundary(%s,%q,0,%t))", kind, value, p.kind, b.e.coerce)
		}
		args = append(args, value)
	}
	name := f.instances[key]
	if name == "" {
		h := fnv.New64a()
		h.Write([]byte(b.e.file.FileName()))
		name = fmt.Sprintf("%s_%x", b.e.unique("Specialized"), h.Sum64())
		f.instances[key] = name
		replacements := map[string]string{}
		for i, param := range f.typeParameters {
			replacements[param] = types[i]
		}
		body := genericPlaceholders.ReplaceAllStringFunc(f.text, func(name string) string { return replacements[name] })
		if f.machine {
			body = b.specializedBody(node, f, parameters, result)
		}
		signature := []string{}
		for i, cell := range f.captures {
			if f.machine {
				signature = append(signature, cell.name+" "+cell.pointerType())
			} else {
				signature = append(signature, fmt.Sprintf("environment%d %s", i, cell.pointerType()))
			}
		}
		for i, typ := range parameters {
			signature = append(signature, fmt.Sprintf("a%d %s", i, typ))
		}
		b.e.classText.WriteString("func " + name + "(" + strings.Join(signature, ",") + ") " + result + " {\n" + body + "\n}\n")
	}
	return b.typedTemp(name+"("+strings.Join(args, ",")+")", result), true
}
func goPrimitive(typ string) primitive {
	switch typ {
	case "float64":
		return primitive{kind: "number"}
	case "bool":
		return primitive{kind: "boolean"}
	case "*tsString":
		return primitive{kind: "string"}
	case "tsValue":
		return primitive{}
	}
	return primitive{kind: typ}
}

// The general specialization path reuses the established lowering machinery,
// with checker substitutions and concrete cells throughout the function body.
func (b *machineBuilder) specializedBody(call *ast.Node, f *genericFunction, parameters []string, result string) string {
	e := b.e
	savedBindings := e.bindings
	copied := make(map[*ast.Node]*binding, len(savedBindings))
	for node, cell := range savedBindings {
		copied[node] = cell
	}
	e.bindings = copied
	e.specializationCalls = append(e.specializationCalls, call)
	defer func() {
		e.bindings = savedBindings
		e.specializationCalls = e.specializationCalls[:len(e.specializationCalls)-1]
	}()
	for node, cell := range savedBindings {
		owner := cell.owner
		local := false
		for owner != nil {
			if owner == f.node {
				local = true
				break
			}
			owner = owner.Parent
		}
		if !local {
			continue
		}
		clone := *cell
		clone.primitive = e.primitive(node)
		e.bindings[node] = &clone
	}
	child := e.newMachine(f.node, b, false)
	child.direct = true
	if result != "tsValue" {
		child.nativeReturn = result
	}
	for i, param := range f.node.Parameters() {
		cell := e.bindings[param]
		method := "init"
		if cell.primitive.kind != "" && cell.primitive.nulls == 0 {
			method = "initNative"
		}
		child.emit(fmt.Sprintf("%s.%s(a%d)", cell.name, method, i))
	}
	child.directStatements(f.node.Body().AsBlock().Statements.Nodes)
	return child.finishDirect()
}

var genericPlaceholders = regexp.MustCompile(`\bG[0-9]+\b`)
