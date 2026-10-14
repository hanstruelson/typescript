// Copyright 2016 the V8 project authors. All rights reserved.
// This file is an algorithmic Go adaptation of behavior declared in V8's
// regexp-ast.h and implemented in regexp-ast.cc.
// See LICENSE.v8 for the applicable license and attribution.
package regexpengine

type v8ASTRangeOperations struct{}

var _ v8ASTCharacterRangeOps = v8ASTRangeOperations{}

func (v8ASTRangeOperations) addClassEscape(set v8ASTStandardCharacterSet, ranges *[]v8ASTCharacterRange, unicodeCaseClosure bool) {
	*ranges = v8TonodeAddClassEscape(set, *ranges, unicodeCaseClosure)
}

func (v8ASTRangeOperations) addCaseEquivalents(_ any, ranges *[]v8ASTCharacterRange, oneByte bool) {
	*ranges = v8ASTAddLegacyCaseEquivalents(*ranges, oneByte)
}

func (v8ASTRangeOperations) addUnicodeCaseEquivalents(ranges *[]v8ASTCharacterRange) {
	*ranges = v8TonodeAddUnicodeCaseEquivalents(*ranges)
}

func (v8ASTRangeOperations) isCanonical(ranges []v8ASTCharacterRange) bool {
	return v8ASTRangesAreCanonical(ranges)
}

func (v8ASTRangeOperations) canonicalize(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	return v8ASTCanonicalizeRanges(ranges)
}

func (v8ASTRangeOperations) negate(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	return v8ASTNegateRanges(ranges)
}

func (v8ASTRangeOperations) intersect(lhs, rhs []v8ASTCharacterRange) []v8ASTCharacterRange {
	return v8ASTIntersectRanges(lhs, rhs)
}

func (v8ASTRangeOperations) intersects(lhs, rhs []v8ASTCharacterRange) bool {
	return v8ASTRangesIntersect(lhs, rhs)
}

func (v8ASTRangeOperations) subtract(src, remove []v8ASTCharacterRange) []v8ASTCharacterRange {
	return v8ASTSubtractRanges(src, remove)
}

func (v8ASTRangeOperations) clampToOneByte(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	return v8ASTClampRangesToOneByte(ranges)
}

func (v8ASTRangeOperations) equals(lhs, rhs []v8ASTCharacterRange) bool {
	return v8ASTRangesEqual(lhs, rhs)
}

func v8ASTAddLegacyCaseEquivalents(ranges []v8ASTCharacterRange, oneByte bool) []v8ASTCharacterRange {
	ranges = v8ASTCanonicalizeRanges(ranges)
	if len(ranges) == 0 {
		return ranges
	}
	latin1Exceptions := [...]rune{0x039c, 0x03bc, 0x0178}
	eligible := make([]v8ASTCharacterRange, 0, len(ranges))
	for _, current := range ranges {
		from, to := current.from, current.to
		if from > 0xffff {
			continue
		}
		if to > 0xffff {
			to = 0xffff
		}
		if from >= v8TonodeLeadSurrogateStart && to <= v8TonodeTrailSurrogateEnd {
			continue
		}
		containsLatin1Exception := false
		for _, exception := range latin1Exceptions {
			if from <= exception && exception <= to {
				containsLatin1Exception = true
				break
			}
		}
		if oneByte && !containsLatin1Exception {
			if from > 0xff {
				continue
			}
			if to > 0xff {
				to = 0xff
			}
		}
		if from <= to {
			eligible = append(eligible, v8ASTRange(from, to))
		}
	}
	if len(eligible) == 0 {
		return ranges
	}
	contains := func(value rune) bool {
		for _, current := range eligible {
			if current.contains(value) {
				return true
			}
			if current.from > value {
				return false
			}
		}
		return false
	}
	closed := append([]v8ASTCharacterRange(nil), ranges...)
	for _, foldClass := range regexpCaseFoldClasses {
		for _, source := range foldClass {
			if source > 0xffff || !contains(source) {
				continue
			}
			canonical := canonicalizeLegacy(source)
			for _, candidate := range foldClass {
				if candidate <= 0xffff && canonicalizeLegacy(candidate) == canonical && !contains(candidate) {
					closed = append(closed, v8ASTRangeSingleton(candidate))
				}
			}
		}
	}
	return v8ASTCanonicalizeRanges(closed)
}

func (set *v8ASTCharacterSet) rangesForCompiler() []v8ASTCharacterRange {
	return v8ASTCharacterSetRanges(set)
}

func (set *v8ASTCharacterSet) standardSetTypeValue() (v8ASTStandardCharacterSet, bool) {
	if set == nil || set.standardSet == nil {
		return 0, false
	}
	return *set.standardSet, true
}

func v8ASTClassSetOperandRanges(operand *v8ASTClassSetOperand) []v8ASTCharacterRange {
	if operand == nil {
		return nil
	}
	return operand.ranges
}

func v8ASTClassSetOperandStrings(operand *v8ASTClassSetOperand) *v8ASTCharacterClassStrings {
	if operand == nil {
		return nil
	}
	return &operand.strings
}

func v8ASTClassSetExpressionOperands(expression *v8ASTClassSetExpression) []any {
	if expression == nil {
		return nil
	}
	return expression.operands
}

func v8ASTClassSetExpressionOperation(expression *v8ASTClassSetExpression) v8ASTClassSetOperation {
	if expression == nil {
		panic("nil V8 class-set expression")
	}
	return expression.operation
}

func v8ASTClassSetExpressionIsNegated(expression *v8ASTClassSetExpression) bool {
	return expression != nil && expression.isNegated
}

func v8ASTClassSetExpressionMayContainStrings(expression *v8ASTClassSetExpression) bool {
	return expression != nil && expression.mayContainStrings
}

func v8ASTAtomData(atom *v8ASTAtom) []uint16 {
	if atom == nil {
		return nil
	}
	return atom.data
}

func v8ASTTextElementCodePointOffset(element *v8ASTTextElement) int {
	if element == nil {
		return -1
	}
	return element.cpOffset
}

func v8ASTTextElementSetCodePointOffset(element *v8ASTTextElement, offset int) {
	if element == nil {
		panic("nil V8 text element")
	}
	element.cpOffset = offset
}

func v8ASTTextElementTree(element *v8ASTTextElement) any {
	if element == nil {
		return nil
	}
	return element.tree
}

func v8ASTTextElementAtom(element *v8ASTTextElement) *v8ASTAtom {
	if element == nil || element.textType != v8ASTTextAtom {
		panic("V8 TextElement::atom called for a non-atom")
	}
	atom, ok := element.tree.(*v8ASTAtom)
	if !ok || atom == nil {
		panic("V8 atom text element contains an invalid tree")
	}
	return atom
}

func v8ASTTextElementClassRanges(element *v8ASTTextElement) *v8ASTClassRanges {
	if element == nil || element.textType != v8ASTTextClassRanges {
		panic("V8 TextElement::class_ranges called for a non-class")
	}
	class, ok := element.tree.(*v8ASTClassRanges)
	if !ok || class == nil {
		panic("V8 class text element contains an invalid tree")
	}
	return class
}

func v8ASTTextElements(text *v8ASTText) []v8ASTTextElement {
	if text == nil {
		return nil
	}
	return text.elements
}

func v8ASTQuantifierBody(quantifier *v8ASTQuantifier) any {
	if quantifier == nil {
		return nil
	}
	return quantifier.body
}

func v8ASTQuantifierMin(quantifier *v8ASTQuantifier) int {
	if quantifier == nil {
		return 0
	}
	return quantifier.min
}

func v8ASTQuantifierMax(quantifier *v8ASTQuantifier) int {
	if quantifier == nil {
		return 0
	}
	return quantifier.max
}

func v8ASTQuantifierKind(quantifier *v8ASTQuantifier) v8ASTQuantifierType {
	if quantifier == nil {
		panic("nil V8 quantifier")
	}
	return quantifier.quantifier
}

func v8ASTQuantifierIndex(quantifier *v8ASTQuantifier) int {
	if quantifier == nil {
		return 0
	}
	return quantifier.index
}

func v8ASTQuantifierIsPossessive(quantifier *v8ASTQuantifier) bool {
	return quantifier != nil && quantifier.quantifier == v8ASTPossessive
}

func v8ASTQuantifierIsNonGreedy(quantifier *v8ASTQuantifier) bool {
	return quantifier != nil && quantifier.quantifier == v8ASTNonGreedy
}

func v8ASTQuantifierIsGreedy(quantifier *v8ASTQuantifier) bool {
	return quantifier != nil && quantifier.quantifier == v8ASTGreedy
}

func v8ASTCaptureBody(capture *v8ASTCapture) any {
	if capture == nil {
		return nil
	}
	return capture.body
}

func v8ASTCaptureIndex(capture *v8ASTCapture) int {
	if capture == nil {
		return 0
	}
	return capture.index
}

func v8ASTCaptureIsBackreferenced(capture *v8ASTCapture) bool {
	return capture != nil && capture.isBackreferenced
}

func (capture *v8ASTCapture) setIsBackreferenced() {
	if capture == nil {
		panic("nil V8 capture")
	}
	capture.isBackreferenced = true
}

func (capture *v8ASTCapture) setName(name []uint16) {
	if capture == nil {
		panic("nil V8 capture")
	}
	capture.name = append(capture.name[:0], name...)
}

func v8ASTCaptureName(capture *v8ASTCapture) []uint16 {
	if capture == nil {
		return nil
	}
	return capture.name
}

func v8ASTGroupBody(group *v8ASTGroup) any {
	if group == nil {
		return nil
	}
	return group.body
}

func v8ASTGroupFlags(group *v8ASTGroup) Flags {
	if group == nil {
		return Flags{}
	}
	return group.flags
}

func v8ASTLookaroundBody(lookaround *v8ASTLookaround) any {
	if lookaround == nil {
		return nil
	}
	return lookaround.body
}

func v8ASTLookaroundIsPositive(lookaround *v8ASTLookaround) bool {
	return lookaround != nil && lookaround.isPositive
}

func v8ASTLookaroundCaptureCount(lookaround *v8ASTLookaround) int {
	if lookaround == nil {
		return 0
	}
	return lookaround.captureCount
}

func v8ASTLookaroundCaptureFrom(lookaround *v8ASTLookaround) int {
	if lookaround == nil {
		return 0
	}
	return lookaround.captureFrom
}

func v8ASTLookaroundTypeOf(lookaround *v8ASTLookaround) v8ASTLookaroundType {
	if lookaround == nil {
		panic("nil V8 lookaround")
	}
	return lookaround.typeOf
}

func v8ASTLookaroundIndex(lookaround *v8ASTLookaround) int {
	if lookaround == nil {
		return 0
	}
	return lookaround.index
}

func (builder *v8ASTLookaroundBuilder) onMatchSuccessValue() v8NodesNode {
	if builder == nil {
		panic("nil V8 lookaround builder")
	}
	node, ok := builder.onMatchSuccess.(v8NodesNode)
	if !ok || node == nil {
		panic("V8 lookaround builder has no match-success node")
	}
	return node
}

func (builder *v8ASTLookaroundBuilder) forMatch(compiler v8ASTCompiler, match v8ASTNode) v8ASTNode {
	graphCompiler, compilerOK := compiler.(*v8CompilerCompiler)
	graphMatch, matchOK := match.(v8NodesNode)
	if !compilerOK || !matchOK || builder == nil || graphMatch == nil {
		panic("V8 lookaround builder requires graph compiler and node values")
	}
	return v8TonodeNewGraphLowerer(nil).lookaroundBuilderForMatch(builder, graphCompiler, graphMatch)
}

func v8ASTBackReferenceCaptures(reference *v8ASTBackReference) []*v8ASTCapture {
	if reference == nil {
		return nil
	}
	return reference.captures
}

func (reference *v8ASTBackReference) setName(name []uint16) {
	if reference == nil {
		panic("nil V8 backreference")
	}
	reference.name = append(reference.name[:0], name...)
}

func v8ASTBackReferenceName(reference *v8ASTBackReference) []uint16 {
	if reference == nil {
		return nil
	}
	return reference.name
}

// v8ASTCharacterSetRanges implements CharacterSet::ranges. Standard sets are
// materialized only when a caller asks for their ranges, preserving the
// upstream lazy representation.
func v8ASTCharacterSetRanges(set *v8ASTCharacterSet) []v8ASTCharacterRange {
	if set == nil {
		panic("nil V8 character set")
	}
	if set.ranges == nil {
		if set.standardSet == nil {
			panic("V8 character set has neither ranges nor a standard-set tag")
		}
		set.ranges = v8TonodeAddClassEscape(*set.standardSet, nil, false)
	}
	return set.ranges
}

func v8ASTCharacterSetCanonicalize(set *v8ASTCharacterSet) {
	if set != nil && set.ranges != nil {
		set.ranges = v8ASTCanonicalizeRanges(set.ranges)
	}
}

func (set *v8ASTCharacterSet) setStandardSetType(kind v8ASTStandardCharacterSet) {
	set.standardSet = &kind
}

func (set *v8ASTCharacterSet) isStandard() bool {
	return set != nil && set.standardSet != nil
}

func (class *v8ASTClassRanges) isStandard() bool {
	if class == nil || class.flags&v8ASTClassNegated != 0 {
		return false
	}
	_, ok := v8TonodeIsStandardClass(&class.set, false)
	return ok
}

func (class *v8ASTClassRanges) standardType() (v8ASTStandardCharacterSet, bool) {
	if class == nil || !class.isStandard() || class.set.standardSet == nil {
		return 0, false
	}
	return *class.set.standardSet, true
}

func (class *v8ASTClassRanges) characterSet() v8ASTCharacterSet {
	if class == nil {
		return v8ASTCharacterSet{}
	}
	return class.set
}

func (class *v8ASTClassRanges) ranges() []v8ASTCharacterRange {
	if class == nil {
		return nil
	}
	return v8ASTCharacterSetRanges(&class.set)
}

func (class *v8ASTClassRanges) isNegated() bool {
	return class != nil && class.flags&v8ASTClassNegated != 0
}

func (class *v8ASTClassRanges) containsSplitSurrogate() bool {
	return class != nil && class.flags&v8ASTClassContainsSplitSurrogate != 0
}

func (class *v8ASTClassRanges) noCaseFoldingNeeded() bool {
	return class != nil && class.flags&v8ASTClassNoCaseFoldingNeeded != 0
}

func (class *v8ASTClassRanges) isCertainlyOneCodePoint() bool {
	return class != nil && class.flags&v8ASTClassCertainlyOneCodePoint != 0
}

func (class *v8ASTClassRanges) isCertainlyTwoCodePoints() bool {
	return class != nil && class.flags&v8ASTClassCertainlyTwoCodePoints != 0
}

func (class *v8ASTClassRanges) appendToText(text *v8ASTText) {
	if class == nil || text == nil {
		panic("V8 class-range text append requires non-nil nodes")
	}
	text.addElement(v8ASTTextElement{textType: v8ASTTextClassRanges, tree: class})
}

func (atom *v8ASTAtom) appendToText(text *v8ASTText) {
	if atom == nil || text == nil {
		panic("V8 atom text append requires non-nil nodes")
	}
	text.addElement(v8ASTTextElement{textType: v8ASTTextAtom, tree: atom})
}

func (text *v8ASTText) appendToText(destination *v8ASTText) {
	if text == nil || destination == nil {
		panic("V8 text append requires non-nil nodes")
	}
	for _, element := range text.elements {
		destination.addElement(element)
	}
}

// v8ASTAppendToText is the Go dispatch equivalent of Tree::AppendToText and
// its Atom/ClassRanges/Text overrides. Unsupported tree kinds are unreachable
// in the parser's text-coalescing path, as they are in V8.
func v8ASTAppendToText(tree any, text *v8ASTText) {
	switch value := tree.(type) {
	case *v8ASTAtom:
		value.appendToText(text)
	case *v8ASTClassRanges:
		value.appendToText(text)
	case *v8ASTText:
		value.appendToText(text)
	default:
		panic("V8 AST node cannot be appended to Text")
	}
}

func v8ASTClassSetStringEqual(lhs, rhs v8ASTCharacterClassString) bool {
	if len(lhs.codePoints) != len(rhs.codePoints) {
		return false
	}
	for i := range lhs.codePoints {
		if lhs.codePoints[i] != rhs.codePoints[i] {
			return false
		}
	}
	return true
}

func v8ASTClassSetHasString(strings v8ASTCharacterClassStrings, candidate v8ASTCharacterClassString) bool {
	for _, existing := range strings {
		if v8ASTClassSetStringEqual(existing, candidate) {
			return true
		}
	}
	return false
}

func (operand *v8ASTClassSetOperand) hasStrings() bool {
	return operand != nil && len(operand.strings) != 0
}

func (operand *v8ASTClassSetOperand) union(other *v8ASTClassSetOperand) {
	if operand == nil || other == nil {
		panic("V8 class-set union requires operands")
	}
	operand.ranges = append(operand.ranges, other.ranges...)
	for _, candidate := range other.strings {
		if !v8ASTClassSetHasString(operand.strings, candidate) {
			candidate.codePoints = append([]rune(nil), candidate.codePoints...)
			operand.strings = append(operand.strings, candidate)
		}
	}
	v8TonodeSortClassSetStrings(operand.strings)
}

func (operand *v8ASTClassSetOperand) intersect(other *v8ASTClassSetOperand, scratch *[]v8ASTCharacterRange) {
	if operand == nil || other == nil || scratch == nil {
		panic("V8 class-set intersection requires operands and scratch ranges")
	}
	*scratch = v8ASTIntersectRanges(operand.ranges, other.ranges)
	operand.ranges, *scratch = *scratch, operand.ranges[:0]
	if len(operand.strings) == 0 {
		return
	}
	if !other.hasStrings() {
		operand.strings = operand.strings[:0]
		return
	}
	kept := operand.strings[:0]
	for _, current := range operand.strings {
		if v8ASTClassSetHasString(other.strings, current) {
			kept = append(kept, current)
		}
	}
	operand.strings = kept
}

func (operand *v8ASTClassSetOperand) subtract(other *v8ASTClassSetOperand, scratch *[]v8ASTCharacterRange) {
	if operand == nil || other == nil || scratch == nil {
		panic("V8 class-set subtraction requires operands and scratch ranges")
	}
	*scratch = v8ASTSubtractRanges(operand.ranges, other.ranges)
	operand.ranges, *scratch = *scratch, operand.ranges[:0]
	if !operand.hasStrings() || !other.hasStrings() {
		return
	}
	kept := operand.strings[:0]
	for _, current := range operand.strings {
		if !v8ASTClassSetHasString(other.strings, current) {
			kept = append(kept, current)
		}
	}
	operand.strings = kept
}

// v8ASTComputeClassSetExpression ports the in-place evaluation in
// ClassSetExpression::ComputeExpression. The evaluated result is installed as
// the expression's sole operand and negation is consumed exactly once.
func v8ASTComputeClassSetExpression(root any, scratch *[]v8ASTCharacterRange) *v8ASTClassSetOperand {
	if scratch == nil {
		panic("V8 class-set expression requires scratch storage")
	}
	if len(*scratch) != 0 {
		panic("V8 class-set expression scratch storage must be empty")
	}
	switch value := root.(type) {
	case *v8ASTClassSetOperand:
		return value
	case *v8ASTClassSetExpression:
		if value == nil || len(value.operands) == 0 {
			panic("V8 class-set expression has no operands")
		}
		result := v8ASTComputeClassSetExpression(value.operands[0], scratch)
		for _, child := range value.operands[1:] {
			operand := v8ASTComputeClassSetExpression(child, scratch)
			switch value.operation {
			case v8ASTClassSetUnion:
				result.union(operand)
			case v8ASTClassSetIntersection:
				result.intersect(operand, scratch)
			case v8ASTClassSetSubtraction:
				result.subtract(operand, scratch)
			default:
				panic("unknown V8 class-set operation")
			}
		}
		if value.operation == v8ASTClassSetUnion {
			result.ranges = v8ASTCanonicalizeRanges(result.ranges)
		}
		if value.isNegated {
			if result.hasStrings() {
				panic("V8 negated class-set expression contains strings")
			}
			*scratch = v8ASTNegateRanges(result.ranges)
			result.ranges, *scratch = *scratch, result.ranges[:0]
			value.isNegated = false
		}
		value.operands = []any{result}
		return result
	default:
		panic("V8 class-set expression contains an invalid node")
	}
}

// v8ASTCaptureRegisters mirrors the overrides in regexp-ast.cc. V8 checks the
// native stack after a soft recursion budget is exhausted. Go grows goroutine
// stacks and exposes no equivalent stack-limit query, so use an explicit work
// stack here; this preserves the tree walk without rejecting a valid tree at
// the soft-budget boundary or risking recursive stack exhaustion.
func v8ASTCaptureRegisters(tree any, _ v8ASTStackLimiter) v8ASTInterval {
	type frame struct {
		tree        any
		children    []any
		single      any
		hasSingle   bool
		next        int
		result      v8ASTInterval
		own         v8ASTInterval
		initialized bool
	}
	stack := []frame{{tree: tree, result: v8ASTEmptyInterval(), own: v8ASTEmptyInterval()}}
	for len(stack) != 0 {
		current := &stack[len(stack)-1]
		if !current.initialized {
			current.initialized = true
			switch value := current.tree.(type) {
			case *v8ASTDisjunction:
				if value == nil {
					return v8ASTInvalidInterval()
				}
				current.children = value.alternatives
			case *v8ASTAlternative:
				if value == nil {
					return v8ASTInvalidInterval()
				}
				current.children = value.nodes
			case *v8ASTCapture:
				if value == nil || value.body == nil {
					return v8ASTInvalidInterval()
				}
				current.single, current.hasSingle = value.body, true
				current.own = v8ASTNewInterval(v8ASTCaptureStartRegister(value.index), v8ASTCaptureEndRegister(value.index))
			case *v8ASTGroup:
				if value == nil || value.body == nil {
					return v8ASTInvalidInterval()
				}
				current.single, current.hasSingle = value.body, true
			case *v8ASTLookaround:
				if value == nil || value.body == nil {
					return v8ASTInvalidInterval()
				}
				current.single, current.hasSingle = value.body, true
			case *v8ASTQuantifier:
				if value == nil || value.body == nil {
					return v8ASTInvalidInterval()
				}
				current.single, current.hasSingle = value.body, true
			case *v8ASTClassSetExpression, *v8ASTText:
				// Neither type overrides Tree::CaptureRegisters.
			case *v8ASTAssertion, *v8ASTClassRanges, *v8ASTClassSetOperand, *v8ASTAtom, *v8ASTBackReference, *v8ASTEmpty:
			default:
				return v8ASTInvalidInterval()
			}
		}

		childCount := len(current.children)
		if current.hasSingle {
			childCount = 1
		}
		if current.next < childCount {
			var child any
			if current.hasSingle {
				child = current.single
			} else {
				child = current.children[current.next]
			}
			current.next++
			stack = append(stack, frame{tree: child, result: v8ASTEmptyInterval(), own: v8ASTEmptyInterval()})
			continue
		}

		result := v8ASTUnionInterval(current.own, current.result)
		stack = stack[:len(stack)-1]
		if len(stack) == 0 {
			return result
		}
		parent := &stack[len(stack)-1]
		parent.result = v8ASTUnionInterval(parent.result, result)
	}
	return v8ASTEmptyInterval()
}

func v8ASTCertainlyAnchoredAtStart(tree any, budget int) bool {
	if assertion, ok := tree.(*v8ASTAssertion); ok {
		return assertion.assertionType == v8ASTStartOfInput
	}
	if budget < 0 {
		return false
	}
	switch value := tree.(type) {
	case *v8ASTAlternative:
		for _, child := range value.nodes {
			if v8ASTCertainlyAnchoredAtStart(child, budget-1) {
				return true
			}
			if v8TonodeTreeMaxMatch(child) > 0 {
				return false
			}
		}
		return false
	case *v8ASTDisjunction:
		for _, child := range value.alternatives {
			if !v8ASTCertainlyAnchoredAtStart(child, budget-1) {
				return false
			}
		}
		return true
	case *v8ASTCapture:
		return v8ASTCertainlyAnchoredAtStart(value.body, budget-1)
	case *v8ASTGroup:
		return v8ASTCertainlyAnchoredAtStart(value.body, budget-1)
	case *v8ASTLookaround:
		return value.isPositive && value.typeOf == v8ASTLookahead && v8ASTCertainlyAnchoredAtStart(value.body, budget-1)
	default:
		return false
	}
}

func v8ASTCertainlyAnchoredAtEnd(tree any, budget int) bool {
	if assertion, ok := tree.(*v8ASTAssertion); ok {
		return assertion.assertionType == v8ASTEndOfInput
	}
	if budget < 0 {
		return false
	}
	switch value := tree.(type) {
	case *v8ASTAlternative:
		for i := len(value.nodes) - 1; i >= 0; i-- {
			child := value.nodes[i]
			if v8ASTCertainlyAnchoredAtEnd(child, budget-1) {
				return true
			}
			if v8TonodeTreeMaxMatch(child) > 0 {
				return false
			}
		}
		return false
	case *v8ASTDisjunction:
		for _, child := range value.alternatives {
			if !v8ASTCertainlyAnchoredAtEnd(child, budget-1) {
				return false
			}
		}
		return true
	case *v8ASTCapture:
		return v8ASTCertainlyAnchoredAtEnd(value.body, budget-1)
	case *v8ASTGroup:
		return v8ASTCertainlyAnchoredAtEnd(value.body, budget-1)
	default:
		return false
	}
}

func v8ASTAccept(tree any, visitor v8ASTVisitor, data any) any {
	if visitor == nil {
		panic("V8 AST visitor is nil")
	}
	switch value := tree.(type) {
	case *v8ASTDisjunction:
		return visitor.visitDisjunction(value, data)
	case *v8ASTAlternative:
		return visitor.visitAlternative(value, data)
	case *v8ASTAssertion:
		return visitor.visitAssertion(value, data)
	case *v8ASTClassRanges:
		return visitor.visitClassRanges(value, data)
	case *v8ASTClassSetOperand:
		return visitor.visitClassSetOperand(value, data)
	case *v8ASTClassSetExpression:
		return visitor.visitClassSetExpression(value, data)
	case *v8ASTAtom:
		return visitor.visitAtom(value, data)
	case *v8ASTQuantifier:
		return visitor.visitQuantifier(value, data)
	case *v8ASTCapture:
		return visitor.visitCapture(value, data)
	case *v8ASTGroup:
		return visitor.visitGroup(value, data)
	case *v8ASTLookaround:
		return visitor.visitLookaround(value, data)
	case *v8ASTBackReference:
		return visitor.visitBackReference(value, data)
	case *v8ASTEmpty:
		return visitor.visitEmpty(value, data)
	case *v8ASTText:
		return visitor.visitText(value, data)
	default:
		panic("V8 AST visitor received an unknown tree kind")
	}
}

func v8ASTTreeIsTextElement(tree any) bool {
	switch tree.(type) {
	case *v8ASTClassRanges, *v8ASTClassSetOperand, *v8ASTClassSetExpression, *v8ASTAtom, *v8ASTText:
		return true
	default:
		return false
	}
}

// The AsX/IsX pairs are the Go equivalents of the DECL_BOILERPLATE macro's
// generated Tree downcasts. A failed AsX returns nil, as in the C++ virtual
// defaults; a concrete match returns the original node pointer.
func v8ASTAsDisjunction(tree any) *v8ASTDisjunction {
	value, _ := tree.(*v8ASTDisjunction)
	return value
}
func v8ASTIsDisjunction(tree any) bool { return v8ASTAsDisjunction(tree) != nil }
func v8ASTAsAlternative(tree any) *v8ASTAlternative {
	value, _ := tree.(*v8ASTAlternative)
	return value
}
func v8ASTIsAlternative(tree any) bool          { return v8ASTAsAlternative(tree) != nil }
func v8ASTAsAssertion(tree any) *v8ASTAssertion { value, _ := tree.(*v8ASTAssertion); return value }
func v8ASTIsAssertion(tree any) bool            { return v8ASTAsAssertion(tree) != nil }
func v8ASTAsClassRanges(tree any) *v8ASTClassRanges {
	value, _ := tree.(*v8ASTClassRanges)
	return value
}
func v8ASTIsClassRanges(tree any) bool { return v8ASTAsClassRanges(tree) != nil }
func v8ASTAsClassSetOperand(tree any) *v8ASTClassSetOperand {
	value, _ := tree.(*v8ASTClassSetOperand)
	return value
}
func v8ASTIsClassSetOperand(tree any) bool { return v8ASTAsClassSetOperand(tree) != nil }
func v8ASTAsClassSetExpression(tree any) *v8ASTClassSetExpression {
	value, _ := tree.(*v8ASTClassSetExpression)
	return value
}
func v8ASTIsClassSetExpression(tree any) bool { return v8ASTAsClassSetExpression(tree) != nil }
func v8ASTAsAtom(tree any) *v8ASTAtom         { value, _ := tree.(*v8ASTAtom); return value }
func v8ASTIsAtom(tree any) bool               { return v8ASTAsAtom(tree) != nil }
func v8ASTAsQuantifier(tree any) *v8ASTQuantifier {
	value, _ := tree.(*v8ASTQuantifier)
	return value
}
func v8ASTIsQuantifier(tree any) bool       { return v8ASTAsQuantifier(tree) != nil }
func v8ASTAsCapture(tree any) *v8ASTCapture { value, _ := tree.(*v8ASTCapture); return value }
func v8ASTIsCapture(tree any) bool          { return v8ASTAsCapture(tree) != nil }
func v8ASTAsGroup(tree any) *v8ASTGroup     { value, _ := tree.(*v8ASTGroup); return value }
func v8ASTIsGroup(tree any) bool            { return v8ASTAsGroup(tree) != nil }
func v8ASTAsLookaround(tree any) *v8ASTLookaround {
	value, _ := tree.(*v8ASTLookaround)
	return value
}
func v8ASTIsLookaround(tree any) bool { return v8ASTAsLookaround(tree) != nil }
func v8ASTAsBackReference(tree any) *v8ASTBackReference {
	value, _ := tree.(*v8ASTBackReference)
	return value
}
func v8ASTIsBackReference(tree any) bool { return v8ASTAsBackReference(tree) != nil }
func v8ASTAsEmpty(tree any) *v8ASTEmpty  { value, _ := tree.(*v8ASTEmpty); return value }
func v8ASTIsEmpty(tree any) bool         { return v8ASTAsEmpty(tree) != nil }
func v8ASTAsText(tree any) *v8ASTText    { value, _ := tree.(*v8ASTText); return value }
func v8ASTIsText(tree any) bool          { return v8ASTAsText(tree) != nil }

// v8ASTTreeToNode is the Tree::ToNode entry point over the Go graph compiler.
// It applies successor pruning and the compiler's periodic overflow/size
// checks before dispatching to the source-mapped ToNodeImpl routines.
func v8ASTTreeToNode(tree any, compiler *v8CompilerCompiler, success v8NodesNode, stack v8CompilerAnalysisStackChecker) v8NodesNode {
	if compiler == nil || success == nil {
		panic("V8 Tree::ToNode requires a compiler and successor")
	}
	if success.IsBacktrack() {
		return success
	}
	compiler.ToNodeMaybeCheckForStackOverflow(stack)
	if compiler.IsRegExpTooBig() {
		return v8TonodeNewGraphLowerer(stack).backtrack(compiler)
	}
	return v8TonodeNewGraphLowerer(stack).treeToNodeImpl(tree, compiler, success)
}

func v8ASTTreeToNodeImpl(tree any, compiler *v8CompilerCompiler, success v8NodesNode, stack v8CompilerAnalysisStackChecker) v8NodesNode {
	return v8TonodeNewGraphLowerer(stack).treeToNodeImpl(tree, compiler, success)
}
