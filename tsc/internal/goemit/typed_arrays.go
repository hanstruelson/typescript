package goemit

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"strconv"
)

func (b *machineBuilder) nativeArrayShape(node *ast.Node, depth int) (string, string) {
	if node == nil || depth > 8 {
		return "", ""
	}
	name := ""
	if node.Kind == ast.KindNewExpression {
		expression := node.AsNewExpression().Expression
		if expression.Kind == ast.KindIdentifier && b.e.binding(expression) == nil {
			name = expression.Text()
		}
	} else if node.Kind == ast.KindIdentifier {
		decl := b.e.reference(node)
		if decl != nil {
			if typ := decl.Type(); typ != nil && typ.Kind == ast.KindTypeReference {
				typeName := typ.AsTypeReferenceNode().TypeName
				if typeName.Kind == ast.KindIdentifier {
					name = typeName.Text()
				}
			}
			if name == "" && decl.Kind == ast.KindVariableDeclaration {
				return b.nativeArrayShape(decl.AsVariableDeclaration().Initializer, depth+1)
			}
		}
	}
	switch name {
	case "Float64Array":
		return name, "float64"
	case "Float32Array":
		return name, "float32"
	case "Int8Array":
		return name, "int8"
	case "Uint8Array", "Uint8ClampedArray":
		return name, "uint8"
	case "Int16Array":
		return name, "int16"
	case "Uint16Array":
		return name, "uint16"
	case "Int32Array":
		return name, "int32"
	case "Uint32Array":
		return name, "uint32"
	}
	return "", ""
}
func (b *machineBuilder) nativeArrayRead(node *ast.Node, receiver, key string) string {
	if p := b.e.arrayElementPrimitive(node); p.kind != "" && p.nulls == 0 {
		return "tsGrowableArrayRead[" + p.goType() + "](" + receiver + "," + key + "," + strconv.Quote(p.kind) + ")"
	}
	if name, shape := b.nativeArrayShape(node, 0); shape != "" {
		return "tsNativeArrayRead[" + shape + "](" + receiver + "," + key + "," + strconv.Quote(name) + ")"
	}
	return "tsGet(" + receiver + "," + key + ")"
}
func (b *machineBuilder) nativeArrayWrite(node *ast.Node, receiver, key, value string) string {
	if p := b.e.arrayElementPrimitive(node); p.kind != "" && p.nulls == 0 {
		return "tsGrowableArrayWrite[" + p.goType() + "](" + receiver + "," + key + "," + value + "," + strconv.Quote(p.kind) + ")"
	}
	if name, shape := b.nativeArrayShape(node, 0); shape != "" {
		return "tsNativeArrayWrite[" + shape + "](" + receiver + "," + key + "," + value + "," + strconv.Quote(name) + ")"
	}
	return "tsSet(" + receiver + "," + key + "," + value + ")"
}

const NativeArrayAccessRuntime = `
type tsNativeArrayStorage[T tsNativeNumber]struct {values []T}
func tsNativeArrayRead[T tsNativeNumber](value,key tsValue,name string)tsValue{if value.kind==tsTypedArrayKind{array:=(*tsTypedArray)(value.ref);if array.name==name&&key.kind==tsNumberKind{index:=key.number;if index>=0&&index<float64(array.length)&&math.Trunc(index)==index{return tsNumberValue(float64((*tsNativeArrayStorage[T])(array.storage).values[int(index)]))}}};return tsGet(value,key)}
func tsNativeArrayWrite[T tsNativeNumber](value,key,item tsValue,name string)tsValue{if value.kind==tsTypedArrayKind{array:=(*tsTypedArray)(value.ref);if array.name==name&&key.kind==tsNumberKind{index:=key.number;if index>=0&&index<float64(array.length)&&math.Trunc(index)==index{number:=tsNumber(tsToPrimitive(item));var result T;switch name{case "Float64Array","Float32Array":result=T(number);case "Uint8ClampedArray":if math.IsNaN(number)||number<=0{number=0}else if number>=255{number=255}else{number=math.RoundToEven(number)};result=T(number);default:result=T(tsNumberUint32(number))};(*tsNativeArrayStorage[T])(array.storage).values[int(index)]=result;return item}}};return tsSet(value,key,item)}
`
