// Package regexpengine AST node definitions for ECMAScript regular expressions.
package regexpengine

import "sort"

// classRange is the parser's compact property-table interchange format.
// The graph AST uses v8ASTCharacterRange after the property is expanded.
type classRange struct{ lo, hi rune }
type classTerm struct {
	ranges   []classRange
	property string
	negated  bool
}

// Ported V8 source declarations from regexp-ast.
// Copyright 2016 the V8 project authors. All rights reserved.
// V8 RegExp AST declarations adapted from regexp-ast.h.
// See LICENSE.v8 for attribution and license terms.
const (
	v8ASTMaxCodePoint = rune(0x10ffff)
	v8ASTInfinity     = quantifierInfinity
	v8ASTIntervalNone = -1
	v8ASTIntervalBad  = -2
)

// v8ASTStandardCharacterSet preserves the exact character tags from
// StandardCharacterSet in regexp-ast.h.
type v8ASTStandardCharacterSet rune

const (
	v8ASTWhitespace        v8ASTStandardCharacterSet = 's'
	v8ASTNotWhitespace     v8ASTStandardCharacterSet = 'S'
	v8ASTWord              v8ASTStandardCharacterSet = 'w'
	v8ASTNotWord           v8ASTStandardCharacterSet = 'W'
	v8ASTDigit             v8ASTStandardCharacterSet = 'd'
	v8ASTNotDigit          v8ASTStandardCharacterSet = 'D'
	v8ASTLineTerminator    v8ASTStandardCharacterSet = 'n'
	v8ASTNotLineTerminator v8ASTStandardCharacterSet = '.'
	v8ASTEverything        v8ASTStandardCharacterSet = '*'
)

// v8ASTInterval is the same inclusive interval representation already ported
// from regexp-ast.cc. The wrappers keep the header API names visible here.
type v8ASTInterval = registerInterval

func v8ASTNewInterval(from, to int) v8ASTInterval {
	return v8ASTInterval{from: from, to: to}
}

func v8ASTEmptyInterval() v8ASTInterval   { return emptyRegisterInterval() }
func v8ASTInvalidInterval() v8ASTInterval { return invalidRegisterInterval() }
func v8ASTUnionInterval(a, b v8ASTInterval) v8ASTInterval {
	return a.union(b)
}
func v8ASTIntervalContains(i v8ASTInterval, value int) bool { return i.contains(value) }
func v8ASTIntervalIsEmpty(i v8ASTInterval) bool             { return i.empty() }
func v8ASTIntervalIsValid(i v8ASTInterval) bool             { return i.valid() }
func v8ASTIntervalFrom(i v8ASTInterval) int                 { return i.from }
func v8ASTIntervalTo(i v8ASTInterval) int                   { return i.to }
func v8ASTIntervalSize(i v8ASTInterval) int                 { return i.size() }

// v8ASTCharacterRange is an inclusive Unicode code-point range. Go's zero
// value matches the C++ default constructor's [0,0] range.
type v8ASTCharacterRange struct {
	from rune
	to   rune
}

func v8ASTRangeSingleton(value rune) v8ASTCharacterRange {
	return v8ASTCharacterRange{from: value, to: value}
}

func v8ASTRange(from, to rune) v8ASTCharacterRange {
	return v8ASTCharacterRange{from: from, to: to}
}

func v8ASTRangeEverything() v8ASTCharacterRange {
	return v8ASTCharacterRange{from: 0, to: v8ASTMaxCodePoint}
}

func v8ASTRangeList(r v8ASTCharacterRange) []v8ASTCharacterRange {
	return []v8ASTCharacterRange{r}
}

func (r v8ASTCharacterRange) contains(value rune) bool {
	return r.from <= value && value <= r.to
}

func (r v8ASTCharacterRange) fromCodePoint() rune { return r.from }

func (r v8ASTCharacterRange) toCodePoint() rune { return r.to }

func (r v8ASTCharacterRange) isEverything(max rune) bool {
	return r.from == 0 && r.to >= max
}

func (r v8ASTCharacterRange) isSingleton() bool { return r.from == r.to }

// v8ASTCharacterRangeOps is the complete static-operation surface declared by
// CharacterRange. Case-equivalence operations remain contracts because V8's
// legacy mode depends on its isolate/ICU canonicalization context.
type v8ASTCharacterRangeOps interface {
	addClassEscape(v8ASTStandardCharacterSet, *[]v8ASTCharacterRange, bool)
	addCaseEquivalents(isolate any, ranges *[]v8ASTCharacterRange, oneByte bool)
	addUnicodeCaseEquivalents(ranges *[]v8ASTCharacterRange)
	isCanonical([]v8ASTCharacterRange) bool
	canonicalize([]v8ASTCharacterRange) []v8ASTCharacterRange
	negate(src []v8ASTCharacterRange) []v8ASTCharacterRange
	intersect(lhs, rhs []v8ASTCharacterRange) []v8ASTCharacterRange
	intersects(lhs, rhs []v8ASTCharacterRange) bool
	subtract(src, remove []v8ASTCharacterRange) []v8ASTCharacterRange
	clampToOneByte([]v8ASTCharacterRange) []v8ASTCharacterRange
	equals(lhs, rhs []v8ASTCharacterRange) bool
}

func v8ASTRangesAreCanonical(ranges []v8ASTCharacterRange) bool {
	for i, current := range ranges {
		if current.from < 0 || current.from > current.to || current.to > v8ASTMaxCodePoint {
			return false
		}
		if i > 0 && ranges[i-1].to+1 >= current.from {
			return false
		}
	}
	return true
}

func v8ASTCanonicalizeRanges(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	if len(ranges) < 2 {
		return append([]v8ASTCharacterRange(nil), ranges...)
	}
	result := append([]v8ASTCharacterRange(nil), ranges...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].from == result[j].from {
			return result[i].to < result[j].to
		}
		return result[i].from < result[j].from
	})
	merged := result[:0]
	for _, current := range result {
		if current.from > current.to || current.to < 0 || current.from > v8ASTMaxCodePoint {
			continue
		}
		if current.from < 0 {
			current.from = 0
		}
		if current.to > v8ASTMaxCodePoint {
			current.to = v8ASTMaxCodePoint
		}
		if len(merged) == 0 || current.from > merged[len(merged)-1].to+1 {
			merged = append(merged, current)
			continue
		}
		if current.to > merged[len(merged)-1].to {
			merged[len(merged)-1].to = current.to
		}
	}
	return merged
}

func v8ASTNegateRanges(src []v8ASTCharacterRange) []v8ASTCharacterRange {
	src = v8ASTCanonicalizeRanges(src)
	result := make([]v8ASTCharacterRange, 0, len(src)+1)
	next := rune(0)
	for _, current := range src {
		if next < current.from {
			result = append(result, v8ASTCharacterRange{from: next, to: current.from - 1})
		}
		if current.to >= v8ASTMaxCodePoint {
			return result
		}
		next = current.to + 1
	}
	if next <= v8ASTMaxCodePoint {
		result = append(result, v8ASTCharacterRange{from: next, to: v8ASTMaxCodePoint})
	}
	return result
}

func v8ASTIntersectRanges(lhs, rhs []v8ASTCharacterRange) []v8ASTCharacterRange {
	lhs, rhs = v8ASTCanonicalizeRanges(lhs), v8ASTCanonicalizeRanges(rhs)
	result := make([]v8ASTCharacterRange, 0)
	for i, j := 0, 0; i < len(lhs) && j < len(rhs); {
		from := lhs[i].from
		if rhs[j].from > from {
			from = rhs[j].from
		}
		to := lhs[i].to
		if rhs[j].to < to {
			to = rhs[j].to
		}
		if from <= to {
			result = append(result, v8ASTCharacterRange{from: from, to: to})
		}
		if lhs[i].to < rhs[j].to {
			i++
		} else {
			j++
		}
	}
	return result
}

func v8ASTRangesIntersect(lhs, rhs []v8ASTCharacterRange) bool {
	return len(v8ASTIntersectRanges(lhs, rhs)) != 0
}

func v8ASTSubtractRanges(src, remove []v8ASTCharacterRange) []v8ASTCharacterRange {
	src, remove = v8ASTCanonicalizeRanges(src), v8ASTCanonicalizeRanges(remove)
	result := make([]v8ASTCharacterRange, 0, len(src))
	for _, source := range src {
		cursor := source.from
		for _, cut := range remove {
			if cut.to < cursor {
				continue
			}
			if cut.from > source.to {
				break
			}
			if cut.from > cursor {
				result = append(result, v8ASTCharacterRange{from: cursor, to: cut.from - 1})
			}
			if cut.to >= source.to {
				cursor = source.to + 1
				break
			}
			if cut.to >= cursor {
				cursor = cut.to + 1
			}
		}
		if cursor <= source.to {
			result = append(result, v8ASTCharacterRange{from: cursor, to: source.to})
		}
	}
	return result
}

func v8ASTClampRangesToOneByte(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	result := make([]v8ASTCharacterRange, 0, len(ranges))
	for _, current := range v8ASTCanonicalizeRanges(ranges) {
		if current.from > 0xff {
			break
		}
		if current.to > 0xff {
			current.to = 0xff
		}
		result = append(result, current)
	}
	return result
}

func v8ASTRangesEqual(lhs, rhs []v8ASTCharacterRange) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i := range lhs {
		if lhs[i] != rhs[i] {
			return false
		}
	}
	return true
}

// v8ASTStackLimiter records the source class's mutable depth budget. Go does
// not prohibit copying structs, so callers must pass it down by pointer when
// they need the C++ non-copying decrement semantics.
type v8ASTStackLimiter struct{ budget int }

func v8ASTNewStackLimiter(budget int) v8ASTStackLimiter {
	return v8ASTStackLimiter{budget: budget}
}

func (limiter *v8ASTStackLimiter) subtract(value int) v8ASTStackLimiter {
	return v8ASTStackLimiter{budget: limiter.budget - value}
}

type v8ASTStackLimiterAPI interface {
	isOverflowed() bool
}

func (limiter *v8ASTStackLimiter) isOverflowed() bool {
	return limiter == nil || limiter.budget <= 0
}

// Compiler and Node are deliberately opaque here. regexp-compiler-tonode.cc
// owns their behavior; this header only uses them as graph-lowering handles.
type v8ASTCompiler interface{}
type v8ASTNode interface{}

type v8ASTVisitor interface {
	visitDisjunction(*v8ASTDisjunction, any) any
	visitAlternative(*v8ASTAlternative, any) any
	visitAssertion(*v8ASTAssertion, any) any
	visitClassRanges(*v8ASTClassRanges, any) any
	visitClassSetOperand(*v8ASTClassSetOperand, any) any
	visitClassSetExpression(*v8ASTClassSetExpression, any) any
	visitAtom(*v8ASTAtom, any) any
	visitQuantifier(*v8ASTQuantifier, any) any
	visitCapture(*v8ASTCapture, any) any
	visitGroup(*v8ASTGroup, any) any
	visitLookaround(*v8ASTLookaround, any) any
	visitBackReference(*v8ASTBackReference, any) any
	visitEmpty(*v8ASTEmpty, any) any
	visitText(*v8ASTText, any) any
}

// v8ASTTree is the complete virtual surface from Tree, including the
// FOR_EACH_REG_EXP_TREE_TYPE-generated AsX/IsX methods. Go's active AST uses
// tagged nodes and switches; this interface records the port target without
// making the dormant state structs pretend to implement it.
type v8ASTTree interface {
	accept(v8ASTVisitor, any) any
	toNode(v8ASTCompiler, v8ASTNode) v8ASTNode
	toNodeImpl(v8ASTCompiler, v8ASTNode) v8ASTNode
	isTextElement() bool
	isCertainlyAnchoredAtStart(budget int) bool
	isCertainlyAnchoredAtEnd(budget int) bool
	minMatch() int
	maxMatch() int
	captureRegisters(*v8ASTStackLimiter) v8ASTInterval
	appendToText(*v8ASTText)
	asDisjunction() *v8ASTDisjunction
	isDisjunction() bool
	asAlternative() *v8ASTAlternative
	isAlternative() bool
	asAssertion() *v8ASTAssertion
	isAssertion() bool
	asClassRanges() *v8ASTClassRanges
	isClassRanges() bool
	asClassSetOperand() *v8ASTClassSetOperand
	isClassSetOperand() bool
	asClassSetExpression() *v8ASTClassSetExpression
	isClassSetExpression() bool
	asAtom() *v8ASTAtom
	isAtom() bool
	asQuantifier() *v8ASTQuantifier
	isQuantifier() bool
	asCapture() *v8ASTCapture
	isCapture() bool
	asGroup() *v8ASTGroup
	isGroup() bool
	asLookaround() *v8ASTLookaround
	isLookaround() bool
	asBackReference() *v8ASTBackReference
	isBackReference() bool
	asEmpty() *v8ASTEmpty
	isEmpty() bool
	asText() *v8ASTText
	isText() bool
}

type v8ASTTreeAPI interface{ v8ASTTree }

type v8ASTDisjunctionAPI interface {
	v8ASTTree
	alternatives() []v8ASTTree
	sortConsecutiveAtoms(v8ASTCompiler) bool
	rationalizeConsecutiveAtoms(v8ASTCompiler)
	fixSingleCharacterDisjunctions(v8ASTCompiler)
}

type v8ASTAlternativeAPI interface {
	v8ASTTree
	nodes() []v8ASTTree
}

type v8ASTAssertionType uint8

const (
	v8ASTStartOfLine v8ASTAssertionType = iota
	v8ASTStartOfInput
	v8ASTEndOfLine
	v8ASTEndOfInput
	v8ASTEndOfBuffer
	v8ASTBoundary
	v8ASTNonBoundary
	v8ASTLastAssertionType = v8ASTNonBoundary
)

type v8ASTAssertionAPI interface {
	v8ASTTree
	assertionType() v8ASTAssertionType
}

type v8ASTCharacterSet struct {
	ranges      []v8ASTCharacterRange
	standardSet *v8ASTStandardCharacterSet
}

type v8ASTCharacterSetAPI interface {
	ranges() []v8ASTCharacterRange
	standardSetType() v8ASTStandardCharacterSet
	setStandardSetType(v8ASTStandardCharacterSet)
	isStandard() bool
	canonicalize()
}

type v8ASTClassRangesFlags uint8

const (
	v8ASTClassNegated v8ASTClassRangesFlags = 1 << iota
	v8ASTClassContainsSplitSurrogate
	v8ASTClassNoCaseFoldingNeeded
	v8ASTClassCertainlyOneCodePoint
	v8ASTClassCertainlyTwoCodePoints
)

type v8ASTClassRanges struct {
	set   v8ASTCharacterSet
	flags v8ASTClassRangesFlags
}

type v8ASTClassRangesAPI interface {
	v8ASTTree
	isStandard() bool
	standardType() v8ASTStandardCharacterSet
	characterSet() v8ASTCharacterSet
	ranges() []v8ASTCharacterRange
	isNegated() bool
	containsSplitSurrogate() bool
	noCaseFoldingNeeded() bool
	isCertainlyOneCodePoint() bool
	isCertainlyTwoCodePoints() bool
}

type v8ASTCharacterClassString struct {
	codePoints []rune
	// Keep the concrete AST payload constructible while v8ASTTree remains a
	// declaration of the upstream virtual API rather than a Go implementation.
	tree any
}

// v8ASTCharacterClassStringLess retains the V8 ordering: longer strings first,
// then ascending code-point order for equal lengths.
type v8ASTCharacterClassStringLess struct{}

func (v8ASTCharacterClassStringLess) less(lhs, rhs []rune) bool {
	if len(lhs) != len(rhs) {
		return len(lhs) > len(rhs)
	}
	for i := range lhs {
		if lhs[i] != rhs[i] {
			return lhs[i] < rhs[i]
		}
	}
	return false
}

// v8ASTCharacterClassStrings is kept as an ordered slice because Go slices,
// unlike V8's ZoneMap keys, cannot be used directly as map keys.
type v8ASTCharacterClassStrings []v8ASTCharacterClassString

type v8ASTClassSetOperand struct {
	ranges   []v8ASTCharacterRange
	strings  v8ASTCharacterClassStrings
	minMatch int
	maxMatch int
}

type v8ASTClassSetOperandAPI interface {
	v8ASTTree
	union(other *v8ASTClassSetOperand)
	intersect(other *v8ASTClassSetOperand, scratch *[]v8ASTCharacterRange)
	subtract(other *v8ASTClassSetOperand, scratch *[]v8ASTCharacterRange)
	hasStrings() bool
	ranges() []v8ASTCharacterRange
	strings() *v8ASTCharacterClassStrings
}

type v8ASTClassSetOperation uint8

const (
	v8ASTClassSetUnion v8ASTClassSetOperation = iota
	v8ASTClassSetIntersection
	v8ASTClassSetSubtraction
)

type v8ASTClassSetExpression struct {
	operation         v8ASTClassSetOperation
	isNegated         bool
	mayContainStrings bool
	operands          []any
	maxMatch          int
}

type v8ASTClassSetExpressionAPI interface {
	v8ASTTree
	operation() v8ASTClassSetOperation
	isNegated() bool
	mayContainStrings() bool
	operands() []v8ASTTree
	computeExpression(root v8ASTTree, scratch *[]v8ASTCharacterRange) *v8ASTClassSetOperand
	emptyClassSetExpression(negated bool) *v8ASTClassSetExpression
}

type v8ASTAtom struct{ data []uint16 }

type v8ASTAtomAPI interface {
	v8ASTTree
	data() []uint16
	length() int
}

type v8ASTTextType uint8

const (
	v8ASTTextAtom v8ASTTextType = iota
	v8ASTTextClassRanges
)

type v8ASTTextElement struct {
	cpOffset int
	textType v8ASTTextType
	// tree holds the concrete atom or class-range payload. The surrounding
	// v8ASTTree provides the common interface for concrete AST node types.
	tree any
}

type v8ASTTextElementAPI interface {
	textType() v8ASTTextType
	codePointOffset() int
	setCodePointOffset(int)
	length() int
	tree() v8ASTTree
	atom() *v8ASTAtom
	classRanges() *v8ASTClassRanges
}

type v8ASTText struct {
	elements []v8ASTTextElement
	length   int
}

type v8ASTTextAPI interface {
	v8ASTTree
	addElement(v8ASTTextElement)
	elements() []v8ASTTextElement
	startsWithAtom() bool
	firstAtom() *v8ASTAtom
}

type v8ASTQuantifierType uint8

const (
	v8ASTGreedy v8ASTQuantifierType = iota
	v8ASTNonGreedy
	v8ASTPossessive
)

type v8ASTQuantifier struct {
	body       any
	min        int
	max        int
	minMatch   int
	maxMatch   int
	quantifier v8ASTQuantifierType
	index      int
}

type v8ASTQuantifierAPI interface {
	toNode(min, max int, greedy bool, body v8ASTTree, compiler v8ASTCompiler, success v8ASTNode, notAtStart bool) v8ASTNode
	min() int
	max() int
	quantifierType() v8ASTQuantifierType
	index() int
	isPossessive() bool
	isNonGreedy() bool
	isGreedy() bool
	body() v8ASTTree
}

type v8ASTCapture struct {
	body             any
	index            int
	minMatch         int
	maxMatch         int
	name             []uint16
	isBackreferenced bool
}

type v8ASTCaptureAPI interface {
	toNode(body v8ASTTree, compiler v8ASTCompiler, success v8ASTNode, capture *v8ASTCapture) v8ASTNode
	body() v8ASTTree
	setBody(v8ASTTree)
	index() int
	isBackreferenced() bool
	setIsBackreferenced()
	name() []uint16
	setName([]uint16)
}

func v8ASTCaptureStartRegister(index int) int { return index * 2 }
func v8ASTCaptureEndRegister(index int) int   { return index*2 + 1 }

type v8ASTGroup struct {
	body     any
	flags    Flags
	minMatch int
	maxMatch int
}

type v8ASTGroupAPI interface {
	v8ASTTree
	body() v8ASTTree
	flags() Flags
}

type v8ASTLookaroundType uint8

const (
	v8ASTLookahead v8ASTLookaroundType = iota
	v8ASTLookbehind
)

type v8ASTLookaround struct {
	body         any
	isPositive   bool
	captureCount int
	captureFrom  int
	typeOf       v8ASTLookaroundType
	index        int
}

type v8ASTLookaroundAPI interface {
	v8ASTTree
	body() v8ASTTree
	isPositive() bool
	captureCount() int
	captureFrom() int
	typeOf() v8ASTLookaroundType
	index() int
}

type v8ASTLookaroundBuilder struct {
	isPositive           bool
	onMatchSuccess       v8ASTNode
	onSuccess            v8ASTNode
	stackPointerRegister int
	positionRegister     int
}

type v8ASTLookaroundBuilderAPI interface {
	onMatchSuccess() v8ASTNode
	forMatch(compiler v8ASTCompiler, match v8ASTNode) v8ASTNode
}

type v8ASTBackReference struct {
	captures []*v8ASTCapture
	name     []uint16
}

type v8ASTBackReferenceAPI interface {
	v8ASTTree
	captures() []*v8ASTCapture
	addCapture(*v8ASTCapture)
	name() []uint16
	setName([]uint16)
}

type v8ASTEmpty struct{}

type v8ASTEmptyAPI interface{ v8ASTTree }

// v8ASTFactories records every constructor overload from the header. Zone
// ownership is replaced by Go slices; v8ASTFactory below computes the same
// cached match bounds and empty-class representation as regexp-ast.cc.
type v8ASTFactories interface {
	newDisjunction([]any) *v8ASTDisjunction
	newAlternative([]any) *v8ASTAlternative
	newAssertion(v8ASTAssertionType) *v8ASTAssertion
	newCharacterSet(v8ASTStandardCharacterSet) v8ASTCharacterSet
	newCharacterSetFromRanges([]v8ASTCharacterRange) v8ASTCharacterSet
	newClassRanges([]v8ASTCharacterRange, v8ASTClassRangesFlags) *v8ASTClassRanges
	newClassRangesFromStandardSet(v8ASTStandardCharacterSet) *v8ASTClassRanges
	newClassSetOperand([]v8ASTCharacterRange, v8ASTCharacterClassStrings) *v8ASTClassSetOperand
	newClassSetExpression(v8ASTClassSetOperation, bool, bool, []any) *v8ASTClassSetExpression
	newAtom([]uint16) *v8ASTAtom
	newText() *v8ASTText
	textElementFromAtom(*v8ASTAtom) v8ASTTextElement
	textElementFromClassRanges(*v8ASTClassRanges) v8ASTTextElement
	newQuantifier(min, max int, kind v8ASTQuantifierType, index int, body any) *v8ASTQuantifier
	newCapture(index int) *v8ASTCapture
	newGroup(body any, flags Flags) *v8ASTGroup
	newLookaround(body any, positive bool, captureCount, captureFrom int, kind v8ASTLookaroundType, index int) *v8ASTLookaround
	newLookaroundBuilder(positive bool, success v8ASTNode, compiler v8ASTCompiler, stackPointerRegister, positionRegister, captureRegisterCount, captureRegisterStart int) *v8ASTLookaroundBuilder
	newBackReference() *v8ASTBackReference
	newBackReferenceTo(*v8ASTCapture) *v8ASTBackReference
	newEmpty() *v8ASTEmpty
}

type v8ASTFactory struct{}

var _ v8ASTFactories = v8ASTFactory{}

func (v8ASTFactory) newDisjunction(alternatives []any) *v8ASTDisjunction {
	if len(alternatives) < 2 {
		panic("V8 Disjunction requires at least two alternatives")
	}
	minimum, maximum := v8ASTInfinity, 0
	for _, alternative := range alternatives {
		childMin, childMax := v8TonodeTreeMinMatch(alternative), v8TonodeTreeMaxMatch(alternative)
		if childMin < minimum {
			minimum = childMin
		}
		if childMax > maximum {
			maximum = childMax
		}
	}
	return &v8ASTDisjunction{alternatives: append([]any(nil), alternatives...), minMatch: minimum, maxMatch: maximum}
}

func (v8ASTFactory) newAlternative(nodes []any) *v8ASTAlternative {
	if len(nodes) < 2 {
		panic("V8 Alternative requires at least two terms")
	}
	minimum, maximum := 0, 0
	for _, child := range nodes {
		minimum = v8ASTSaturatingMatchSum(minimum, v8TonodeTreeMinMatch(child))
		maximum = v8ASTSaturatingMatchSum(maximum, v8TonodeTreeMaxMatch(child))
	}
	return &v8ASTAlternative{nodes: append([]any(nil), nodes...), minMatch: minimum, maxMatch: maximum}
}

func (v8ASTFactory) newAssertion(kind v8ASTAssertionType) *v8ASTAssertion {
	if kind > v8ASTLastAssertionType {
		panic("unknown V8 assertion kind")
	}
	return &v8ASTAssertion{assertionType: kind}
}

func (v8ASTFactory) newCharacterSet(kind v8ASTStandardCharacterSet) v8ASTCharacterSet {
	return v8ASTCharacterSet{standardSet: &kind}
}

func (v8ASTFactory) newCharacterSetFromRanges(ranges []v8ASTCharacterRange) v8ASTCharacterSet {
	return v8ASTCharacterSet{ranges: append([]v8ASTCharacterRange(nil), ranges...)}
}

func (v8ASTFactory) newClassRanges(ranges []v8ASTCharacterRange, flags v8ASTClassRangesFlags) *v8ASTClassRanges {
	copyOfRanges := append([]v8ASTCharacterRange(nil), ranges...)
	if len(copyOfRanges) == 0 {
		copyOfRanges = v8ASTRangeList(v8ASTRangeEverything())
		flags ^= v8ASTClassNegated
	}
	class := &v8ASTClassRanges{set: v8ASTCharacterSet{ranges: copyOfRanges}, flags: flags}
	if class.flags&v8ASTClassNegated == 0 && class.flags&v8ASTClassCertainlyTwoCodePoints == 0 && class.flags&v8ASTClassNoCaseFoldingNeeded != 0 {
		foundBasicPlane := false
		for _, interval := range copyOfRanges {
			if interval.from < 0x10000 {
				foundBasicPlane = true
				break
			}
		}
		if !foundBasicPlane {
			class.flags |= v8ASTClassCertainlyTwoCodePoints
		}
	}
	return class
}

func (factory v8ASTFactory) newClassRangesFromStandardSet(kind v8ASTStandardCharacterSet) *v8ASTClassRanges {
	return &v8ASTClassRanges{set: factory.newCharacterSet(kind)}
}

func (v8ASTFactory) newClassSetOperand(ranges []v8ASTCharacterRange, strings v8ASTCharacterClassStrings) *v8ASTClassSetOperand {
	operand := &v8ASTClassSetOperand{ranges: append([]v8ASTCharacterRange(nil), ranges...), strings: append(v8ASTCharacterClassStrings(nil), strings...)}
	// Class-set alternatives can be strings with more than one code unit. Start
	// the minimum at infinity so a string-only operand gets the minimum of its
	// alternatives instead of incorrectly retaining the zero value. A range
	// always consumes at least one UTF-16 code unit and at most two.
	operand.minMatch = v8ASTInfinity
	if len(operand.ranges) != 0 {
		operand.minMatch, operand.maxMatch = 1, 2
	}
	for _, stringAlternative := range operand.strings {
		minimum := v8TonodeTreeMinMatch(stringAlternative.tree)
		maximum := v8TonodeTreeMaxMatch(stringAlternative.tree)
		if minimum < operand.minMatch {
			operand.minMatch = minimum
		}
		if maximum > operand.maxMatch {
			operand.maxMatch = maximum
		}
	}
	if operand.minMatch == v8ASTInfinity {
		// An empty operand has no successful match. Zero is the conservative
		// lower bound used by the quantifier lowering path for that case.
		operand.minMatch = 0
	}
	v8TonodeSortClassSetStrings(operand.strings)
	return operand
}

func (v8ASTFactory) newClassSetExpression(operation v8ASTClassSetOperation, negated, mayContainStrings bool, operands []any) *v8ASTClassSetExpression {
	if negated && mayContainStrings {
		panic("V8 negated class-set expression cannot contain strings")
	}
	maximum := 0
	if negated {
		maximum = 2
	} else {
		for _, operand := range operands {
			if candidate := v8TonodeTreeMaxMatch(operand); candidate > maximum {
				maximum = candidate
			}
		}
	}
	return &v8ASTClassSetExpression{operation: operation, isNegated: negated, mayContainStrings: mayContainStrings, operands: append([]any(nil), operands...), maxMatch: maximum}
}

func (v8ASTFactory) newAtom(data []uint16) *v8ASTAtom {
	return &v8ASTAtom{data: append([]uint16(nil), data...)}
}

func (v8ASTFactory) newText() *v8ASTText {
	return &v8ASTText{elements: make([]v8ASTTextElement, 0, 2)}
}

func (v8ASTFactory) textElementFromAtom(atom *v8ASTAtom) v8ASTTextElement {
	if atom == nil {
		panic("V8 atom text element requires an atom")
	}
	return v8ASTTextElement{cpOffset: -1, textType: v8ASTTextAtom, tree: atom}
}

func (v8ASTFactory) textElementFromClassRanges(class *v8ASTClassRanges) v8ASTTextElement {
	if class == nil {
		panic("V8 class-range text element requires class ranges")
	}
	return v8ASTTextElement{cpOffset: -1, textType: v8ASTTextClassRanges, tree: class}
}

func (v8ASTFactory) newQuantifier(minimum, maximum int, kind v8ASTQuantifierType, index int, body any) *v8ASTQuantifier {
	if body == nil || minimum < 0 || maximum < minimum || kind > v8ASTPossessive {
		panic("invalid V8 quantifier state")
	}
	bodyMin, bodyMax := v8TonodeTreeMinMatch(body), v8TonodeTreeMaxMatch(body)
	return &v8ASTQuantifier{
		body: body, min: minimum, max: maximum, quantifier: kind, index: index,
		minMatch: v8ASTSaturatingMatchProduct(minimum, bodyMin),
		maxMatch: v8ASTSaturatingMatchProduct(maximum, bodyMax),
	}
}

func (v8ASTFactory) newCapture(index int) *v8ASTCapture {
	if index < 0 {
		panic("negative V8 capture index")
	}
	return &v8ASTCapture{index: index}
}

func (v8ASTFactory) newGroup(body any, flags Flags) *v8ASTGroup {
	if body == nil {
		panic("V8 group requires a body")
	}
	return &v8ASTGroup{body: body, flags: flags, minMatch: v8TonodeTreeMinMatch(body), maxMatch: v8TonodeTreeMaxMatch(body)}
}

func (v8ASTFactory) newLookaround(body any, positive bool, captureCount, captureFrom int, kind v8ASTLookaroundType, index int) *v8ASTLookaround {
	if body == nil || captureCount < 0 || captureFrom < 0 || kind > v8ASTLookbehind {
		panic("invalid V8 lookaround state")
	}
	return &v8ASTLookaround{body: body, isPositive: positive, captureCount: captureCount, captureFrom: captureFrom, typeOf: kind, index: index}
}

func (v8ASTFactory) newLookaroundBuilder(positive bool, success v8ASTNode, compiler v8ASTCompiler, stackRegister, positionRegister, captureRegisterCount, captureRegisterStart int) *v8ASTLookaroundBuilder {
	graphCompiler, ok := compiler.(*v8CompilerCompiler)
	graphSuccess, nodeOK := success.(v8NodesNode)
	if !ok || !nodeOK {
		panic("V8 lookaround builder requires the graph compiler and node types")
	}
	return v8TonodeNewGraphLowerer(nil).newLookaroundBuilder(positive, graphSuccess, graphCompiler, stackRegister, positionRegister, captureRegisterCount, captureRegisterStart)
}

func (v8ASTFactory) newBackReference() *v8ASTBackReference {
	return &v8ASTBackReference{}
}

func (factory v8ASTFactory) newBackReferenceTo(capture *v8ASTCapture) *v8ASTBackReference {
	if capture == nil {
		panic("V8 backreference requires a capture")
	}
	return &v8ASTBackReference{captures: []*v8ASTCapture{capture}}
}

func (v8ASTFactory) newEmpty() *v8ASTEmpty { return &v8ASTEmpty{} }

func (v8ASTFactory) emptyClassSetExpression(negated bool) *v8ASTClassSetExpression {
	operand := &v8ASTClassSetOperand{}
	return (v8ASTFactory{}).newClassSetExpression(v8ASTClassSetUnion, negated, false, []any{operand})
}

func v8ASTSaturatingMatchSum(previous, increase int) int {
	if previous >= v8ASTInfinity || increase >= v8ASTInfinity || increase > v8ASTInfinity-previous {
		return v8ASTInfinity
	}
	return previous + increase
}

func v8ASTSaturatingMatchProduct(count, width int) int {
	if count <= 0 || width <= 0 {
		return 0
	}
	if count >= v8ASTInfinity || width > v8ASTInfinity/count {
		return v8ASTInfinity
	}
	return count * width
}

func (class *v8ASTClassRanges) minMatch() int {
	if class.flags&v8ASTClassCertainlyTwoCodePoints != 0 {
		return 2
	}
	return 1
}

func (class *v8ASTClassRanges) maxMatch() int {
	if class.flags&v8ASTClassCertainlyOneCodePoint != 0 {
		return 1
	}
	return 2
}

func (element v8ASTTextElement) length() int {
	switch element.textType {
	case v8ASTTextAtom:
		atom, ok := element.tree.(*v8ASTAtom)
		if !ok || atom == nil {
			panic("V8 text atom element has an invalid payload")
		}
		return len(atom.data)
	case v8ASTTextClassRanges:
		return 1
	default:
		panic("unknown V8 text element kind")
	}
}

func (text *v8ASTText) addElement(element v8ASTTextElement) {
	text.elements = append(text.elements, element)
	text.length += element.length()
}

func (text *v8ASTText) startsWithAtom() bool {
	return text != nil && len(text.elements) != 0 && text.elements[0].textType == v8ASTTextAtom
}

func (text *v8ASTText) firstAtom() *v8ASTAtom {
	if !text.startsWithAtom() {
		return nil
	}
	atom, _ := text.elements[0].tree.(*v8ASTAtom)
	return atom
}

func (capture *v8ASTCapture) setBody(body any) {
	if capture == nil || body == nil {
		panic("V8 capture body cannot be nil")
	}
	capture.body = body
	capture.minMatch = v8TonodeTreeMinMatch(body)
	capture.maxMatch = v8TonodeTreeMaxMatch(body)
}

func (backReference *v8ASTBackReference) addCapture(capture *v8ASTCapture) {
	if capture == nil {
		panic("V8 backreference capture cannot be nil")
	}
	backReference.captures = append(backReference.captures, capture)
}

func (expression *v8ASTClassSetExpression) emptyClassSetExpression(negated bool) *v8ASTClassSetExpression {
	return (v8ASTFactory{}).emptyClassSetExpression(negated)
}

// v8ASTDisjunction and its siblings below are field-for-field state records
// for the corresponding final classes in regexp-ast.h.
type v8ASTDisjunction struct {
	alternatives []any
	minMatch     int
	maxMatch     int
}

type v8ASTAlternative struct {
	nodes    []any
	minMatch int
	maxMatch int
}

type v8ASTAssertion struct{ assertionType v8ASTAssertionType }
