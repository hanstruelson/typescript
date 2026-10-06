package goemit

import (
	"bytes"
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
		return nil, err
	}
	pass := valueEmission{fields: map[string]map[string]string{}, functions: map[string][]valueSignature{}, globals: map[string]string{}}
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
			if d.Name.Name == "tsValueOf" || d.Name.Name == "tsUnbox" || d.Name.Name == "tsPrimitiveValue" || d.Name.Name == "tsNative" {
				continue
			}
			env := copyValueEnvironment(pass.globals)
			pass.parameters(d.Type, env)
			if d.Recv != nil {
				pass.fieldNames(d.Recv, env)
			}
			pass.block(d.Body, env, valueFunctionSignature(d.Type).result)
		}
	}
	var out bytes.Buffer
	err = format.Node(&out, set, file)
	return out.Bytes(), err
}

type valueSignature struct {
	parameters []string
	variadic   bool
	result     string
}
type valueEmission struct {
	fields    map[string]map[string]string
	functions map[string][]valueSignature
	globals   map[string]string
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
	if fn.Results != nil && len(fn.Results.List) > 0 {
		sig.result = valueType(fn.Results.List[0].Type)
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
	if at := strings.IndexByte(kind, '['); at >= 0 {
		kind = kind[:at]
	}
	return kind
}
func valueElement(kind string) string {
	if strings.HasPrefix(kind, "[]") {
		return kind[2:]
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
	}
	return ""
}
func (p *valueEmission) signature(call *ast.CallExpr, env map[string]string) valueSignature {
	// A local callback or function-valued field takes precedence over a method
	// with the same name elsewhere in the generated runtime.
	if ident, ok := call.Fun.(*ast.Ident); ok {
		if typ := env[ident.Name]; strings.HasPrefix(typ, "func(") {
			if parsed, err := parser.ParseExpr(typ); err == nil {
				if fn, ok := parsed.(*ast.FuncType); ok {
					return valueFunctionSignature(fn)
				}
			}
		}
	}
	candidates := p.functions[valueCallName(call.Fun)]
	for _, sig := range candidates {
		if (!sig.variadic && len(sig.parameters) == len(call.Args)) || (sig.variadic && len(call.Args) >= len(sig.parameters)-1) {
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
		return p.fields[valueBase(p.kind(e.X, env))][e.Sel.Name]
	case *ast.FuncLit:
		return valueType(e.Type)
	case *ast.CallExpr:
		name := valueCallName(e.Fun)
		switch name {
		case "tsNumberValue", "tsBooleanValue", "tsStringReference", "tsValueOf", "tsPrimitiveValue":
			return "tsValue"
		case "float64", "float32", "int", "int64", "uint16", "uint32", "uint64", "bool", "string":
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
	if expected != "tsValue" {
		return expr
	}
	kind := p.kind(expr, env)
	switch kind {
	case "tsValue":
		return expr
	case "nil":
		return ast.NewIdent("tsNull")
	case "float64", "float32", "int", "int64", "uint16", "uint32", "uint64":
		if kind != "float64" {
			expr = valueWrapper("float64", expr)
		}
		return valueWrapper("tsNumberValue", expr)
	case "bool":
		return valueWrapper("tsBooleanValue", expr)
	case "*tsString":
		return valueWrapper("tsStringReference", expr)
	case "string":
		return valueWrapper("tsStringReference", valueWrapper("tsStringUTF8", expr))
	case "T":
		return valueWrapper("tsPrimitiveValue", expr)
	}
	return valueWrapper("tsValueOf", expr)
}
func (p *valueEmission) expression(expr ast.Expr, env map[string]string) ast.Expr {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *ast.CallExpr:
		e.Fun = p.expression(e.Fun, env)
		for i, arg := range e.Args {
			e.Args[i] = p.expression(arg, env)
		}
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
					e.Args[i] = p.coerce(arg, sig.parameters[at], env)
				}
			}
		}
	case *ast.TypeAssertExpr:
		e.X = p.expression(e.X, env)
		if p.kind(e.X, env) == "tsValue" {
			// Preserve Go's comma-ok form (used by the runtime's type switches and
			// checked assertions). Unboxing keeps the original assertion semantics.
			e.X = valueWrapper("tsUnbox", e.X)
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
		p.block(e.Body, child, valueFunctionSignature(e.Type).result)
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
		for i, left := range s.Lhs {
			if i >= len(s.Rhs) {
				if name, ok := left.(*ast.Ident); ok {
					env[name.Name] = "bool"
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
			s.Results[i] = p.coerce(p.expression(expr, env), result, env)
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
