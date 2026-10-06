package goemit

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/core"
)

// ModuleInput uses the compiler's resolved import graph rather than resolving
// paths independently. The resolver is borrowed only while this fragment emits.
type ModuleInput struct {
	File     *ast.SourceFile
	Resolver ReferenceResolver
	Targets  map[*ast.Node]string
}

func EmitModule(input ModuleInput, options *core.CompilerOptions) (string, []*ast.Diagnostic) {
	e := &emitter{file: input.File, resolver: input.Resolver, bindings: make(map[*ast.Node]*binding), imports: make(map[*ast.Node]string), targets: input.Targets}
	e.coerce = !options.CoerceAny.IsFalse()
	e.strictNulls = options.GetStrictOptionValue(options.StrictNullChecks)
	if options.SourceMap.IsTrue() || options.InlineSourceMap.IsTrue() {
		e.fail(input.File.AsNode(), "source maps are not supported")
	}
	if options.Jsx != core.JsxEmitNone {
		e.fail(input.File.AsNode(), "JSX is not supported")
	}
	if len(input.File.Diagnostics()) != 0 {
		return "", input.File.Diagnostics()
	}
	e.collect(input.File.AsNode(), input.File.AsNode())
	e.planNativeFunctions()
	e.planGenericFunctions()
	e.prepareClasses()
	if len(e.diags) != 0 {
		return "", e.diags
	}
	builder := e.newMachine(input.File.AsNode(), nil, true)
	builder.module = true
	for _, node := range input.File.Statements.Nodes {
		if node.Kind == ast.KindImportDeclaration && input.Targets[node] == "" {
			builder.importDeclaration(node)
		}
	}
	builder.statements(input.File.Statements.Nodes)
	builder.emit("m.complete(tsU)")
	return e.classText.String() + builder.finish(), e.diags
}

// Bundle formats one executable containing a persistent initializer per source.
func Bundle(names, fragments, roots []string) (string, error) {
	var out strings.Builder
	out.WriteString("package main\n\n" + Runtime + ValueRuntime + ModuleRuntime + ClassRuntime + TypeRuntime + StringRuntime + RegexRuntime + EqualityRuntime + ObjectRuntime + "\n")
	for index, fragment := range fragments {
		signature := fmt.Sprintf("func tsSourceModule%d(loop *tsLoop,modules map[string]*tsModule,module *tsModule)", index)
		out.WriteString(strings.Replace(fragment, "func(module *tsModule)", signature, 1))
		out.WriteString("\n")
	}
	out.WriteString("func main(){loop:=tsNewLoop();modules:=make(map[string]*tsModule)\n")
	for _, name := range names {
		fmt.Fprintf(&out, "modules[%q]=tsNewModule(loop,%q)\n", name, name)
	}
	for index, name := range names {
		fmt.Fprintf(&out, "modules[%q].initialize=func(module *tsModule){tsSourceModule%d(loop,modules,module)}\n", name, index)
	}
	for _, name := range roots {
		fmt.Fprintf(&out, "modules[%q].request(nil,true,func(result tsResult){if result.rejected {loop.errors=append(loop.errors,result.value)}})\n", name)
	}
	out.WriteString("if err:=loop.run();err!=nil {fmt.Fprintln(os.Stderr,err);os.Exit(1)}\n}\n")
	formatted, err := formatValueSource([]byte(out.String()))
	return string(formatted), err
}
func (b *machineBuilder) publish(name, value, cell string) {
	if !b.module {
		b.e.fail(b.owner, "exports require module bundling")
		return
	}
	next := b.block()
	b.emit(fmt.Sprintf("m.exportValue(%q,%s,%s,%d);return", name, value, cell, next))
	b.current = next
}
func (b *machineBuilder) requestModule(node *ast.Node, names []string, full bool) string {
	target := b.e.targets[node]
	if target == "" {
		b.e.fail(node, "module was not resolved by the compiler")
		return "tsU"
	}
	wanted := []string{}
	for _, name := range names {
		wanted = append(wanted, strconv.Quote(name))
	}
	next := b.block()
	b.emit(fmt.Sprintf("m.importModule(modules[%q],[]string{%s},%t,%d);return", target, strings.Join(wanted, ","), full, next))
	b.current = next
	return "modules[" + strconv.Quote(target) + "]"
}
func (b *machineBuilder) localImport(node *ast.Node) {
	clause := node.AsImportDeclaration().ImportClause
	if clause == nil {
		b.requestModule(node, nil, true)
		return
	}
	c := clause.AsImportClause()
	if c.PhaseModifier == ast.KindTypeKeyword {
		return
	}
	names := []string{}
	bindings := []struct {
		cell *binding
		name string
	}{}
	if clause.Name() != nil {
		names = append(names, "default")
		bindings = append(bindings, struct {
			cell *binding
			name string
		}{b.e.bindings[clause], "default"})
	}
	namespace := false
	if c.NamedBindings != nil {
		if c.NamedBindings.Kind == ast.KindNamedImports {
			for _, specifier := range c.NamedBindings.AsNamedImports().Elements.Nodes {
				s := specifier.AsImportSpecifier()
				if s.IsTypeOnly {
					continue
				}
				name := specifier.Name().Text()
				if s.PropertyName != nil {
					name = s.PropertyName.Text()
				}
				names = append(names, name)
				bindings = append(bindings, struct {
					cell *binding
					name string
				}{b.e.bindings[specifier], name})
			}
		} else {
			namespace = true
		}
	}
	// A purely type-only list performs no runtime demand.
	if len(names) == 0 && !namespace && c.NamedBindings != nil && len(c.NamedBindings.AsNamedImports().Elements.Nodes) > 0 {
		return
	}
	module := b.requestModule(node, names, namespace || len(names) == 0)
	for _, entry := range bindings {
		if entry.cell != nil {
			b.emit(fmt.Sprintf("%s.init(tsImportRef{module:%s,name:%q})", entry.cell.name, module, entry.name))
		}
	}
	if namespace {
		if cell := b.e.bindings[c.NamedBindings]; cell != nil {
			b.emit(cell.name + ".init(tsNamespace{module:" + module + "})")
		}
	}
}
func (b *machineBuilder) exportDeclaration(node *ast.Node) {
	n := node.AsExportDeclaration()
	if n.IsTypeOnly {
		return
	}
	if n.ExportClause == nil || n.ExportClause.Kind != ast.KindNamedExports {
		b.e.fail(node, "export-star requires an explicit named export list in this backend")
		return
	}
	if !b.module {
		b.e.fail(node, "exports require module bundling")
		return
	}
	names := []string{}
	for _, specifier := range n.ExportClause.AsNamedExports().Elements.Nodes {
		s := specifier.AsExportSpecifier()
		if s.IsTypeOnly {
			continue
		}
		name := specifier.Name().Text()
		if s.PropertyName != nil {
			name = s.PropertyName.Text()
		}
		names = append(names, name)
	}
	target := ""
	if n.ModuleSpecifier != nil {
		target = b.requestModule(node, names, false)
	}
	for _, specifier := range n.ExportClause.AsNamedExports().Elements.Nodes {
		s := specifier.AsExportSpecifier()
		if s.IsTypeOnly {
			continue
		}
		local := specifier.Name()
		if s.PropertyName != nil {
			local = s.PropertyName
		}
		if target != "" {
			b.publish(specifier.Name().Text(), fmt.Sprintf("tsImportRef{module:%s,name:%q}", target, local.Text()), "nil")
			continue
		}
		cell := b.e.binding(local)
		if cell == nil {
			if symbol := b.e.file.AsNode().Locals()[local.Text()]; symbol != nil {
				cell = b.e.bindings[symbol.ValueDeclaration]
			}
		}
		if cell == nil {
			b.e.fail(specifier, "cannot resolve local export "+local.Text())
			continue
		}
		b.publish(specifier.Name().Text(), cell.name+".get()", cell.name)
	}
}
