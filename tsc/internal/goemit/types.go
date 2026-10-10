package goemit

import (
	"fmt"
	"strconv"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

type primitive struct {
	kind  string
	nulls uint8
}

func syntaxPrimitive(node *ast.Node) primitive {
	if node == nil {
		return primitive{}
	}
	switch node.Kind {
	case ast.KindNumberKeyword:
		return primitive{kind: "number"}
	case ast.KindFloat32Keyword:
		return primitive{kind: "float32"}
	case ast.KindFloat64Keyword:
		return primitive{kind: "float64"}
	case ast.KindIntKeyword:
		return primitive{kind: "int"}
	case ast.KindInt8Keyword:
		return primitive{kind: "int8"}
	case ast.KindInt16Keyword:
		return primitive{kind: "int16"}
	case ast.KindInt32Keyword:
		return primitive{kind: "int32"}
	case ast.KindInt64Keyword:
		return primitive{kind: "int64"}
	case ast.KindUintKeyword:
		return primitive{kind: "uint"}
	case ast.KindUint8Keyword:
		return primitive{kind: "uint8"}
	case ast.KindUint16Keyword:
		return primitive{kind: "uint16"}
	case ast.KindUint32Keyword:
		return primitive{kind: "uint32"}
	case ast.KindUint64Keyword:
		return primitive{kind: "uint64"}

	case ast.KindStringKeyword:
		return primitive{kind: "string"}
	case ast.KindBooleanKeyword:
		return primitive{kind: "boolean"}
	case ast.KindNullKeyword:
		return primitive{nulls: 1}
	case ast.KindUndefinedKeyword:
		return primitive{nulls: 2}
	case ast.KindTypeReference:
		if n := node.AsTypeReferenceNode().TypeName; n.Kind == ast.KindIdentifier && n.Text() == "undefined" {
			return primitive{nulls: 2}
		}
	case ast.KindLiteralType:
		return syntaxPrimitive(node.AsLiteralTypeNode().Literal)
	case ast.KindUnionType:
		result := primitive{}
		for _, part := range node.AsUnionTypeNode().Types.Nodes {
			p := syntaxPrimitive(part)
			if p.kind == "" && p.nulls == 0 {
				return primitive{}
			}
			if p.kind != "" {
				if result.kind != "" && result.kind != p.kind {
					return primitive{}
				}
				result.kind = p.kind
			}
			result.nulls |= p.nulls
		}
		return result
	}
	return primitive{}
}
func (e *emitter) primitive(node *ast.Node) primitive {
	if e.dynamicFieldAccess(node) {
		return primitive{}
	}
	if node == nil {
		return primitive{}
	}
	if len(e.specializationCalls) > 0 {
		if resolver, ok := e.resolver.(interface {
			GetEmitSpecializedPrimitiveType(*ast.Node, []*ast.Node) (string, uint8)
		}); ok {
			kind, mask := resolver.GetEmitSpecializedPrimitiveType(node, e.specializationCalls)
			return primitive{kind, mask}
		}
	}

	declaration := node
	if node.Kind == ast.KindIdentifier {
		declaration = e.reference(node)
	}
	if declaration != nil && declaration.Kind == ast.KindVariableDeclaration {
		if init := declaration.AsVariableDeclaration().Initializer; init != nil && init.Kind == ast.KindCallExpression {
			callee := init.AsCallExpression().Expression
			if callee.Kind == ast.KindIdentifier && callee.Text() == "setTimeout" {
				decl := e.reference(callee)
				if decl == nil || (decl.Kind == ast.KindFunctionDeclaration && decl.Body() == nil) {
					return primitive{}
				}
			}
		}
	}
	if p := e.annotationPrimitive(node.Type()); p.kind != "" {
		if node.Kind == ast.KindParameter && node.QuestionToken() != nil {
			p.nulls |= 2
		}
		return p
	}
	if resolver, ok := e.resolver.(interface {
		GetEmitPrimitiveType(*ast.Node) (string, uint8)
	}); ok {
		kind, mask := resolver.GetEmitPrimitiveType(node)
		return primitive{kind, mask}
	}
	return primitive{}
}
func (p primitive) goType() string {
	switch p.kind {
	case "number":
		return "float64"
	case "string":
		return "*tsString"
	case "boolean":
		return "bool"
	case "float32", "float64", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return p.kind
	}
	return "tsValue"
}
func (cell *binding) pointerType() string {
	if cell.arrayElement.kind != "" {
		return "*tsDenseArrayCell[" + cell.arrayElement.goType() + "]"
	}
	if cell.primitive.kind != "" {
		return "*tsTypedCell[" + cell.primitive.goType() + "]"
	}
	return "*tsCell"
}
func (b *machineBuilder) bindingFactory(cell *binding, initialized bool) string {
	if p := cell.arrayElement; p.kind != "" {
		return fmt.Sprintf("tsDenseArrayBinding[%s](%t,%t,%q,%t,%d)", p.goType(), initialized, cell.constant, p.kind, b.e.coerce, p.nulls)
	}
	if cell.primitive.kind == "" {
		if p := b.e.arrayElementPrimitive(cell.declaration); p.kind != "" && !ast.IsFunctionLike(cell.declaration) {
			return fmt.Sprintf("tsArrayBinding(%t,%t,%q,%d,%t)", initialized, cell.constant, p.kind, p.nulls, b.e.coerce)
		}
		return fmt.Sprintf("tsBinding(%t,%t)", initialized, cell.constant)
	}
	return fmt.Sprintf("tsTypedBinding[%s](%t,%t,%q,%d,%t)", cell.primitive.goType(), initialized, cell.constant, cell.primitive.kind, cell.primitive.nulls, b.e.coerce)
}
func (b *machineBuilder) typedTemp(value, kind string) string {
	name := b.temp(value)
	if b.tempTypes == nil {
		b.tempTypes = map[string]string{}
	}
	b.tempTypes[name] = kind
	return name
}
func (b *machineBuilder) tempType(name string) string {
	if kind := b.tempTypes[name]; kind != "" {
		return kind
	}
	return "tsValue"
}
func (b *machineBuilder) validateParameter(param *ast.Node) {
	if cell := b.e.bindings[param]; cell != nil && cell.primitive.kind != "" {
		b.emit(cell.name + ".require()")
	}
}
func (b *machineBuilder) returnValue(value string) string {
	p := b.e.returnPrimitive(b.owner)
	if p.kind == "" || b.constructor {
		if !b.constructor {
			return b.arrayBoundary(value, b.owner)
		}
		return value
	}
	return fmt.Sprintf("tsBoundary(%s,%s,%d,%t)", value, strconv.Quote(p.kind), p.nulls, b.e.coerce)
}
func (e *emitter) annotationPrimitive(node *ast.Node) primitive {
	if node == nil {
		return primitive{}
	}
	if len(e.specializationCalls) > 0 {
		if resolver, ok := e.resolver.(interface {
			GetEmitSpecializedPrimitiveType(*ast.Node, []*ast.Node) (string, uint8)
		}); ok {
			kind, mask := resolver.GetEmitSpecializedPrimitiveType(node, e.specializationCalls)
			return primitive{kind, mask}
		}
	}

	// Preserve explicitly written nullable primitives even without strict null checks.
	if p := syntaxPrimitive(node); p.kind != "" {
		return p
	}
	// The checker instantiates generic aliases and resolves lexical type names.
	// A type parameter is deliberately dynamic, even when an outer alias shares
	// its spelling. Textual expansion cannot safely make that decision.
	if resolver, ok := e.resolver.(interface {
		GetEmitPrimitiveType(*ast.Node) (string, uint8)
	}); ok {
		kind, mask := resolver.GetEmitPrimitiveType(node)
		return primitive{kind, mask}
	}
	return e.annotationSeen(node, map[*ast.Node]bool{})
}
func (e *emitter) annotationSeen(node *ast.Node, seen map[*ast.Node]bool) primitive {
	if node == nil || seen[node] {
		return primitive{}
	}
	seen[node] = true
	defer delete(seen, node)
	if node.Kind == ast.KindTypeReference {
		n := node.AsTypeReferenceNode().TypeName
		if n.Kind == ast.KindIdentifier {
			if symbol := e.file.AsNode().Locals()[n.Text()]; symbol != nil {
				for _, decl := range symbol.Declarations {
					if decl.Kind == ast.KindTypeAliasDeclaration {
						return e.annotationSeen(decl.AsTypeAliasDeclaration().Type, seen)
					}
				}
			}
		}
	}
	if node.Kind == ast.KindParenthesizedType {
		return e.annotationSeen(node.AsParenthesizedTypeNode().Type, seen)
	}
	if node.Kind == ast.KindUnionType {
		result := primitive{}
		for _, part := range node.AsUnionTypeNode().Types.Nodes {
			p := e.annotationSeen(part, seen)
			if p.kind == "" && p.nulls == 0 {
				return primitive{}
			}
			if p.kind != "" {
				if result.kind != "" && result.kind != p.kind {
					return primitive{}
				}
				result.kind = p.kind
			}
			result.nulls |= p.nulls
		}
		return result
	}
	return syntaxPrimitive(node)
}

// Reads of captured or not-yet-initialized var bindings retain undefined.
func (b *machineBuilder) definitelyInitialized(cell *binding, node *ast.Node) bool {
	if !cell.maybeUndefined {
		return true
	}
	if cell.owner != b.owner || cell.declaration.Kind != ast.KindVariableDeclaration {
		return false
	}
	decl := cell.declaration.AsVariableDeclaration()
	if decl.Initializer == nil || node.Loc.Pos() < cell.declaration.Loc.End() {
		return false
	}
	parent := cell.declaration.Parent
	if parent == nil {
		return false
	}
	parent = parent.Parent
	if parent == nil || parent.Kind != ast.KindVariableStatement {
		return false
	}
	container := parent.Parent
	return container == b.owner || (b.owner.Body() != nil && container == b.owner.Body())
}

// Parameter bindings remain uninitialized until their position in the parameter
// list is evaluated, including names nested inside destructuring patterns.
func (cell *binding) isParameter() bool {
	for node := cell.declaration; node != nil && node != cell.owner; node = node.Parent {
		if node.Kind == ast.KindParameter {
			return true
		}
	}
	return false
}

func (p primitive) numeric() bool {
	switch p.kind {
	case "number", "float32", "float64", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	}
	return false
}

func (e *emitter) returnPrimitive(owner *ast.Node) primitive {
	node := owner.Type()
	if node != nil && node.Kind == ast.KindTypeReference {
		ref := node.AsTypeReferenceNode()
		if ref.TypeName.Kind == ast.KindIdentifier && ref.TypeName.Text() == "Promise" && ref.TypeArguments != nil && len(ref.TypeArguments.Nodes) == 1 {
			return e.annotationPrimitive(ref.TypeArguments.Nodes[0])
		}
	}
	return e.annotationPrimitive(node)
}

func (p primitive) elementGoType() string {
	if p.kind == "" {
		return "tsValue"
	}
	if p.nulls != 0 {
		return "tsOptional[" + p.goType() + "]"
	}
	return p.goType()
}
