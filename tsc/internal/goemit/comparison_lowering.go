package goemit

import (
	"fmt"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/core"
)

// nativeNumericShape describes a comparison operand after expression
// evaluation. Nullable operands retain their native payload and a presence
// tag; dynamic values never enter this plan.
type nativeNumericShape struct {
	goType   string
	optional bool
	value    string
	tag      string
}

type nativePrimitiveShape struct {
	family   string
	optional bool
	value    string
	tag      string
}

func makeNativePrimitiveShape(kind, expression string) (nativePrimitiveShape, bool) {
	optional := optionalValueType(kind)
	if optional != "" {
		kind = optional
	}
	family := ""
	switch {
	case nativeNumericType(kind):
		family = "numeric"
	case kind == "bool":
		family = "boolean"
	case kind == "*tsString":
		family = "string"
	default:
		return nativePrimitiveShape{}, false
	}
	tag := "0"
	value := expression
	if optional != "" {
		tag = expression + ".tag"
		value = expression + ".value"
	}
	return nativePrimitiveShape{family: family, optional: optional != "", value: value, tag: tag}, true
}

// nativePrimitiveComparison handles nonnumeric same-family operations and
// strict mismatches that are statically impossible. Dynamic or coercive mixed
// pairs stay on the policy/runtime path.
func (b *machineBuilder) nativePrimitiveComparison(op, left, right string) (string, bool) {
	leftShape, ok := makeNativePrimitiveShape(b.tempType(left), left)
	if !ok {
		return "", false
	}
	rightShape, ok := makeNativePrimitiveShape(b.tempType(right), right)
	if !ok {
		return "", false
	}

	if leftShape.family != rightShape.family {
		if (leftShape.family == "boolean" && rightShape.family == "numeric") || (leftShape.family == "numeric" && rightShape.family == "boolean") {
			booleanShape, numericShape := leftShape, rightShape
			booleanOnLeft := true
			if leftShape.family == "numeric" {
				booleanShape, numericShape = rightShape, leftShape
				booleanOnLeft = false
			}
			if op == "==" || op == "!=" {
				booleanValue := booleanShape.value
				if booleanShape.optional {
					booleanValue = "(" + booleanShape.tag + "==0 && " + booleanShape.value + ")"
				}
				payloadEqual := "((" + booleanValue + " && " + numericShape.value + "==1) || (!" + booleanValue + " && " + numericShape.value + "==0))"
				equal := payloadEqual
				if booleanShape.optional && numericShape.optional {
					equal = fmt.Sprintf("((%s!=0 && %s!=0) || (%s==0 && %s==0 && %s))", booleanShape.tag, numericShape.tag, booleanShape.tag, numericShape.tag, payloadEqual)
				} else {
					if booleanShape.optional {
						equal = "(" + booleanShape.tag + "==0 && " + equal + ")"
					}
					if numericShape.optional {
						equal = "(" + numericShape.tag + "==0 && " + equal + ")"
					}
				}
				if op == "!=" {
					equal = "!(" + equal + ")"
				}
				return b.typedTemp(equal, "bool"), true
			}
			if op == "<" || op == "<=" || op == ">" || op == ">=" {
				booleanValue := booleanShape.value
				guards := []string{}
				if booleanShape.optional {
					booleanValue = "(" + booleanShape.tag + "==0 && " + booleanShape.value + ")"
					guards = append(guards, booleanShape.tag+"!=2")
				}
				compare := func(number string) string {
					zero, one := "0"+op+number, "1"+op+number
					if !booleanOnLeft {
						zero, one = number+op+"0", number+op+"1"
					}
					return "((" + booleanValue + " && " + one + ") || (!" + booleanValue + " && " + zero + "))"
				}
				terms := []string{}
				if numericShape.optional {
					terms = append(terms, "("+numericShape.tag+"==0 && "+compare(numericShape.value)+")")
					terms = append(terms, "("+numericShape.tag+"==1 && "+compare("0")+")")
				} else {
					terms = append(terms, compare(numericShape.value))
				}
				result := "(" + strings.Join(terms, " || ") + ")"
				if len(guards) != 0 {
					result = "(" + strings.Join(guards, " && ") + " && " + result + ")"
				}
				if numericShape.optional {
					result = "(" + numericShape.tag + "!=2 && " + result + ")"
				}
				return b.typedTemp(result, "bool"), true
			}
		}
		if op != "===" && op != "!==" {
			return "", false
		}
		equal := "false"
		if leftShape.optional && rightShape.optional {
			equal = fmt.Sprintf("(%s!=0 && %s==%s)", leftShape.tag, leftShape.tag, rightShape.tag)
		}
		if op == "!==" {
			equal = "!(" + equal + ")"
		}
		return b.typedTemp(equal, "bool"), true
	}

	if op == "==" || op == "!=" || op == "===" || op == "!==" {
		var payloadEqual string
		switch leftShape.family {
		case "boolean":
			payloadEqual = leftShape.value + "==" + rightShape.value
		case "string":
			payloadEqual = "tsStringEqual(" + leftShape.value + "," + rightShape.value + ")"
		default:
			return "", false
		}

		equal := payloadEqual
		if leftShape.optional && rightShape.optional {
			if op == "==" || op == "!=" {
				equal = fmt.Sprintf("((%s!=0 && %s!=0) || (%s==0 && %s==0 && %s))", leftShape.tag, rightShape.tag, leftShape.tag, rightShape.tag, payloadEqual)
			} else {
				equal = fmt.Sprintf("((%s!=0 && %s==%s) || (%s==0 && %s==0 && %s))", leftShape.tag, leftShape.tag, rightShape.tag, leftShape.tag, rightShape.tag, payloadEqual)
			}
		} else if leftShape.optional {
			equal = "(" + leftShape.tag + "==0 && " + payloadEqual + ")"
		} else if rightShape.optional {
			equal = "(" + rightShape.tag + "==0 && " + payloadEqual + ")"
		}
		if op == "!=" || op == "!==" {
			equal = "!(" + equal + ")"
		}
		return b.typedTemp(equal, "bool"), true
	}
	if leftShape.family == "string" && (op == "<" || op == "<=" || op == ">" || op == ">=") && (leftShape.optional || rightShape.optional) {
		leftTag, rightTag := leftShape.tag, rightShape.tag
		terms := []string{"(" + leftTag + "==0 && " + rightTag + "==0 && tsStringCompare(" + leftShape.value + "," + rightShape.value + ")" + op + "0)"}
		zeroCompare := "0" + op + "0"
		if leftShape.optional && rightShape.optional {
			terms = append(terms, "("+leftTag+"==1 && "+rightTag+"==1 && "+zeroCompare+")")
		}
		if leftShape.optional {
			terms = append(terms, "("+leftTag+"==1 && "+rightTag+"==0 && 0"+op+"tsStringNumber("+rightShape.value+"))")
		}
		if rightShape.optional {
			terms = append(terms, "("+leftTag+"==0 && "+rightTag+"==1 && tsStringNumber("+leftShape.value+")"+op+"0)")
		}
		return b.typedTemp("("+strings.Join(terms, " || ")+")", "bool"), true
	}

	if leftShape.family == "boolean" {
		var leftTrue, rightTrue string
		guards := []string{}
		if leftShape.optional {
			leftTrue = "(" + leftShape.tag + "==0 && " + leftShape.value + ")"
			guards = append(guards, leftShape.tag+"!=2")
		} else {
			leftTrue = leftShape.value
		}
		if rightShape.optional {
			rightTrue = "(" + rightShape.tag + "==0 && " + rightShape.value + ")"
			guards = append(guards, rightShape.tag+"!=2")
		} else {
			rightTrue = rightShape.value
		}
		var ordered string
		switch op {
		case "<":
			ordered = "(!" + leftTrue + " && " + rightTrue + ")"
		case "<=":
			ordered = "(!" + leftTrue + " || " + rightTrue + ")"
		case ">":
			ordered = "(" + leftTrue + " && !" + rightTrue + ")"
		case ">=":
			ordered = "(" + leftTrue + " || !" + rightTrue + ")"
		default:
			return "", false
		}
		if len(guards) != 0 {
			ordered = "(" + strings.Join(guards, " && ") + " && " + ordered + ")"
		}
		return b.typedTemp(ordered, "bool"), true
	}
	return "", false
}

func nativeNumericType(kind string) bool {
	switch kind {
	case "float32", "float64", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	default:
		return false
	}
}

func makeNativeNumericShape(kind, expression string) (nativeNumericShape, bool) {
	optional := optionalValueType(kind)
	if optional != "" {
		if !nativeNumericType(optional) {
			return nativeNumericShape{}, false
		}
		return nativeNumericShape{goType: optional, optional: true, value: expression + ".value", tag: expression + ".tag"}, true
	}
	if !nativeNumericType(kind) {
		return nativeNumericShape{}, false
	}
	return nativeNumericShape{goType: kind, value: expression, tag: "0"}, true
}

// nativeNumericComparison emits the fastest exact plan available for a
// statically represented numeric pair. It returns false for dynamic and
// nonnumeric operands, which retain the shared runtime comparison path.
func (b *machineBuilder) nativeNumericComparison(op, left, right string) (string, bool) {
	leftShape, ok := makeNativeNumericShape(b.tempType(left), left)
	if !ok {
		return "", false
	}
	rightShape, ok := makeNativeNumericShape(b.tempType(right), right)
	if !ok {
		return "", false
	}

	expression := ""
	switch op {
	case "==", "!=", "===", "!==":
		payloadEqual := nativeNumericPayloadCompare("==", leftShape.value, leftShape.goType, rightShape.value, rightShape.goType)
		loose := op == "==" || op == "!="
		switch {
		case !leftShape.optional && !rightShape.optional:
			expression = payloadEqual
		case leftShape.optional && !rightShape.optional:
			expression = leftShape.tag + "==0 && " + payloadEqual
		case !leftShape.optional && rightShape.optional:
			expression = rightShape.tag + "==0 && " + payloadEqual
		case loose:
			expression = fmt.Sprintf("((%s!=0 && %s!=0) || (%s==0 && %s==0 && %s))", leftShape.tag, rightShape.tag, leftShape.tag, rightShape.tag, payloadEqual)
		default:
			expression = fmt.Sprintf("(%s==%s && (%s!=0 || %s))", leftShape.tag, rightShape.tag, leftShape.tag, payloadEqual)
		}
		if op == "!=" || op == "!==" {
			expression = "!(" + expression + ")"
		}
	case "<", "<=", ">", ">=":
		terms := []string{}
		appendTerm := func(condition, leftValue, leftType, rightValue, rightType string) {
			comparison := nativeNumericPayloadCompare(op, leftValue, leftType, rightValue, rightType)
			if condition == "true" {
				terms = append(terms, comparison)
			} else {
				terms = append(terms, fmt.Sprintf("(%s && %s)", condition, comparison))
			}
		}
		present := []string{}
		if leftShape.optional {
			present = append(present, leftShape.tag+"==0")
		}
		if rightShape.optional {
			present = append(present, rightShape.tag+"==0")
		}
		presentCondition := strings.Join(present, " && ")
		if presentCondition == "" {
			presentCondition = "true"
		}
		appendTerm(presentCondition, leftShape.value, leftShape.goType, rightShape.value, rightShape.goType)
		if leftShape.optional {
			appendTerm(leftShape.tag+"==1 && "+rightShape.tag+"==0", leftShape.goType+"(0)", leftShape.goType, rightShape.value, rightShape.goType)
		}
		if rightShape.optional {
			appendTerm(leftShape.tag+"==0 && "+rightShape.tag+"==1", leftShape.value, leftShape.goType, rightShape.goType+"(0)", rightShape.goType)
		}
		if leftShape.optional && rightShape.optional {
			appendTerm(leftShape.tag+"==1 && "+rightShape.tag+"==1", leftShape.goType+"(0)", leftShape.goType, rightShape.goType+"(0)", rightShape.goType)
		}
		expression = "(" + strings.Join(terms, " || ") + ")"
	default:
		return "", false
	}
	return b.typedTemp(expression, "bool"), true
}

// nativeDynamicComparison preserves the known operand's native shape until
// runtime dispatch. The generic helpers centralize policy conversion and keep
// emitted code independent of the number of source type pairs.
func (b *machineBuilder) nativeDynamicComparison(op, known, dynamic string, p primitive, knownLeft bool) (string, bool) {
	if p.kind == "" {
		switch optionalValueType(b.tempType(known)) {
		case "*tsString":
			p.kind = "string"
		case "bool":
			p.kind = "boolean"
		}
	}
	if p.kind == "" || (p.kind != "string" && p.kind != "boolean" && !p.numeric()) {
		return "", false
	}
	goType := p.goType()
	if optionalValueType(b.tempType(known)) != "" {
		return b.typedTemp(fmt.Sprintf("tsPolicyCompareOptionalNative[%s](loop,%q,%s,%s,%q,%d,%t,%t)", goType, op, known, dynamic, p.kind, p.nulls, b.e.coerce, knownLeft), "bool"), true
	}
	if b.tempType(known) != goType {
		return "", false
	}
	return b.typedTemp(fmt.Sprintf("tsPolicyCompareNative[%s](loop,%q,%s,%s,%q,%d,%t,%t)", goType, op, known, dynamic, p.kind, p.nulls, b.e.coerce, knownLeft), "bool"), true
}

// truthCondition and nullishCondition keep known native values out of tsValue
// merely to decide a branch. The caller has already evaluated value once.
func (b *machineBuilder) truthCondition(value string) string {
	kind := b.tempType(value)
	optional := optionalValueType(kind)
	payload := value
	if optional != "" {
		kind = optional
		payload += ".value"
	}
	condition := ""
	switch kind {
	case "bool":
		condition = payload
	case "*tsString":
		condition = "len(" + payload + ".units)!=0"
	case "float32", "float64":
		condition = fmt.Sprintf("(%s!=0 && !math.IsNaN(float64(%s)))", payload, payload)
	default:
		if nativeNumericType(kind) {
			condition = payload + "!=0"
		}
	}
	if condition == "" {
		return "tsTruthy(" + value + ")"
	}
	if optional != "" {
		return "(" + value + ".tag==0 && " + condition + ")"
	}
	return condition
}

func (b *machineBuilder) nullishCondition(value string) string {
	kind := b.tempType(value)
	if optionalValueType(kind) != "" {
		return value + ".tag!=0"
	}
	if kind == "tsValue" {
		return "tsNullish(" + value + ")"
	}
	if kind == "bool" || kind == "*tsString" || nativeNumericType(kind) {
		return "false"
	}
	return "tsNullish(" + value + ")"
}

// nativeBitwiseComparison follows the runtime's existing families: an exact
// same-type integer pair uses that native width, while mixed numeric types and
// booleans use JavaScript's Number 32-bit conversion. Strings are primitive
// Number inputs here; dynamic operands retain runtime ToPrimitive dispatch.
func (b *machineBuilder) nativeBitwiseComparison(op, left, right string) (string, bool) {
	leftType, rightType := b.tempType(left), b.tempType(right)
	if optionalValueType(leftType) != "" || optionalValueType(rightType) != "" {
		return "", false
	}
	leftInteger := nativeIntegerType(leftType)
	rightInteger := nativeIntegerType(rightType)
	if leftInteger && rightInteger && leftType == rightType {
		if op == "&" || op == "|" || op == "^" {
			return b.typedTemp(left+op+right, leftType), true
		}
		return b.typedTemp(fmt.Sprintf("tsNativeIntegerBitwise[%s](%q,%s,%s)", leftType, op, left, right), leftType), true
	}
	leftNumber, ok := bitwiseNumberOperand(left, leftType)
	if !ok {
		return "", false
	}
	rightNumber, ok := bitwiseNumberOperand(right, rightType)
	if !ok {
		return "", false
	}
	return b.typedTemp(fmt.Sprintf("tsNumberBitwise(%q,%s,%s)", op, leftNumber, rightNumber), "float64"), true
}

func nativeIntegerType(kind string) bool {
	switch kind {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	default:
		return false
	}
}

func bitwiseNumberOperand(value, kind string) (string, bool) {
	switch {
	case kind == "bool":
		return "tsBooleanNumber(" + value + ")", true
	case kind == "*tsString":
		return "tsStringNumber(" + value + ")", true
	case nativeNumericType(kind):
		return "float64(" + value + ")", true
	default:
		return "", false
	}
}

func nativeNumericPayloadCompare(op, left, leftType, right, rightType string) string {
	target := core.NativeNumericCommonType(leftType, rightType)
	if target != "" {
		if target != leftType {
			left = target + "(" + left + ")"
		}
		if target != rightType {
			right = target + "(" + right + ")"
		}
		if op == "==" || op == "===" {
			return left + "==" + right
		}
		if op == "!=" || op == "!==" {
			return left + "!=" + right
		}
		return left + op + right
	}
	return fmt.Sprintf("tsNumericComparison(%q,tsPrimitiveValue(%s),tsPrimitiveValue(%s))", op, left, right)
}
