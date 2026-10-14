// Package regexpengine V8-shaped regular-expression graph construction.
package regexpengine

import (
	"sort"
)

// --- Opaque objects owned by neighboring V8 subsystems --------------------

type v8NodesZone = v8CompilerZone
type v8NodesNodePrinter[T any] struct{}
type v8NodesDotPrinterImpl struct{}
type v8NodesCapture = v8ASTCapture
type v8NodesClassRanges = v8ASTClassRanges
type v8NodesCharacterRange = v8ASTCharacterRange
type v8NodesTextElement = v8ASTTextElement
type v8NodesInterval = v8ASTInterval
type v8NodesFlags = Flags

// v8NodesZoneList mirrors the owning, arena-backed vectors used by the graph
// IR. The slice is only its Go storage representation; allocation lifetime
// and arena semantics remain an integration task.
type v8NodesZoneList[T any] struct{ items []T }

func (list *v8NodesZoneList[T]) Len() int { return len(list.items) }
func (list *v8NodesZoneList[T]) At(index int) *T {
	return &list.items[index]
}
func (list *v8NodesZoneList[T]) Add(value T) { list.items = append(list.items, value) }

// Forward declarations from regexp-nodes.h. The assembler-specific owners are
// opaque in the Go graph compiler; their data contracts are represented below
// where they affect analysis or source generation.
type v8NodesAlternativeGenerationList struct{}

type v8NodesPreloadState struct{}
type v8NodesSpecialLoopState struct{}
type v8NodesTrace struct{}

type v8NodesNodeInfo struct {
	beingAnalyzed       bool
	beenAnalyzed        bool
	followsWordInterest bool
	followsNewline      bool
	followsStart        bool
	atEnd               bool
	visited             bool
	replacementComputed bool
}

func v8NodesNewNodeInfo() v8NodesNodeInfo { return v8NodesNodeInfo{} }

func (info v8NodesNodeInfo) Matches(other v8NodesNodeInfo) bool {
	return info.atEnd == other.atEnd &&
		info.followsWordInterest == other.followsWordInterest &&
		info.followsNewline == other.followsNewline &&
		info.followsStart == other.followsStart
}

func (info *v8NodesNodeInfo) AddFromPreceding(other v8NodesNodeInfo) {
	info.atEnd = info.atEnd || other.atEnd
	info.followsWordInterest = info.followsWordInterest || other.followsWordInterest
	info.followsNewline = info.followsNewline || other.followsNewline
	info.followsStart = info.followsStart || other.followsStart
}

func (info v8NodesNodeInfo) HasLookbehind() bool {
	return info.followsWordInterest || info.followsNewline || info.followsStart
}

func (info *v8NodesNodeInfo) AddFromFollowing(other v8NodesNodeInfo) {
	info.followsWordInterest = info.followsWordInterest || other.followsWordInterest
	info.followsNewline = info.followsNewline || other.followsNewline
	info.followsStart = info.followsStart || other.followsStart
}

func (info *v8NodesNodeInfo) ResetCompilationState() {
	info.beingAnalyzed = false
	info.beenAnalyzed = false
}

type v8NodesEatsAtLeastInfo struct {
	fromPossiblyStart uint8
	fromNotStart      uint8
}

func v8NodesNewEatsAtLeastInfo(eats uint8) v8NodesEatsAtLeastInfo {
	return v8NodesEatsAtLeastInfo{fromPossiblyStart: eats, fromNotStart: eats}
}

func (info *v8NodesEatsAtLeastInfo) SetMin(other v8NodesEatsAtLeastInfo) {
	if other.fromPossiblyStart < info.fromPossiblyStart {
		info.fromPossiblyStart = other.fromPossiblyStart
	}
	if other.fromNotStart < info.fromNotStart {
		info.fromNotStart = other.fromNotStart
	}
}

func (info *v8NodesEatsAtLeastInfo) SetMax(other int) {
	maximum := uint8(other)
	if other < 0 {
		maximum = 0
	} else if other > 0xff {
		maximum = 0xff
	}
	if maximum > info.fromPossiblyStart {
		info.fromPossiblyStart = maximum
	}
	if maximum > info.fromNotStart {
		info.fromNotStart = maximum
	}
}

func (info v8NodesEatsAtLeastInfo) IsZero() bool {
	return info.fromPossiblyStart == 0 && info.fromNotStart == 0
}

type v8NodesEmitResult uint8

const (
	v8NodesEmitSuccess v8NodesEmitResult = iota
	v8NodesEmitError
)

func v8NodesSuccess() v8NodesEmitResult { return v8NodesEmitSuccess }
func v8NodesError() v8NodesEmitResult   { return v8NodesEmitError }

func (result v8NodesEmitResult) IsSuccess() bool { return result == v8NodesEmitSuccess }
func (result v8NodesEmitResult) IsError() bool   { return result == v8NodesEmitError }

type v8NodesDrainMode uint8

const (
	v8NodesDrainFull v8NodesDrainMode = iota
	v8NodesDrainRetryAtEntry
	v8NodesDrainRestoreOnly
	v8NodesDrainOmit
)

type v8NodesAtomicLoopKind uint8

const (
	v8NodesAtomicLoopNone v8NodesAtomicLoopKind = iota
	v8NodesAtomicLoopAtEnd
	v8NodesAtomicLoopTotal
	v8NodesAtomicLoopBoundary
	v8NodesAtomicLoopDisjoint
)

type v8NodesParkedGrant uint8

const (
	v8NodesParkedGrantNone v8NodesParkedGrant = iota
	v8NodesParkedGrantParked
	v8NodesParkedGrantUniformPrefix
	v8NodesParkedGrantNonEmptyUniformPrefix
)

// v8NodesNode is the graph-only node contract consumed by AST lowering and
// the generated-Go backend. It deliberately contains no assembler entry point.
type v8NodesNode interface {
	Accept(v8NodesNodeVisitor)
	AsEndNode() *v8NodesEndNode
	AsActionNode() *v8NodesActionNode
	AsChoiceNode() *v8NodesChoiceNode
	AsLoopChoiceNode() *v8NodesLoopChoiceNode
	AsNegativeLookaroundChoiceNode() *v8NodesNegativeLookaroundChoiceNode
	AsBackReferenceNode() *v8NodesBackReferenceNode
	AsAssertionNode() *v8NodesAssertionNode
	AsTextNode() *v8NodesTextNode
	AsUnanchoredAdvanceNode() *v8NodesUnanchoredAdvanceNode
	AsNegativeSubmatchSuccess() *v8NodesNegativeSubmatchSuccess
	AsSeqNode() *v8NodesSeqNode
	Zone() *v8NodesZone
	Flags() v8NodesFlags
	IsBacktrack() bool
	EatsAtLeast(notAtStart bool) uint32
	Info() *v8NodesNodeInfo
	EatsAtLeastInfo() *v8NodesEatsAtLeastInfo
	SetEatsAtLeastInfo(v8NodesEatsAtLeastInfo)
	OnWorkList() bool
	SetOnWorkList(bool)
	KeepRecursing(*v8CompilerCompiler) bool
	FixedLengthLoopLength() int
}

// v8NodesNodeBase ports the graph-owned state shared by every V8 Node. Label
// binding and assembler emission are deliberately represented by graph IDs in
// regexp_go_codegen.go; all node-local analysis and rewrite metadata lives here.
type v8NodesNodeBase struct {
	replacement v8NodesNode
	flags       v8NodesFlags
	onWorkList  bool
	info        v8NodesNodeInfo
	eatsAtLeast v8NodesEatsAtLeastInfo
	traceCount  int
	bmInfo      [2]*v8NodesBoyerMooreLookahead
	zone        *v8NodesZone
}

const (
	v8NodesLargeEatsAtLeastValue  = 255
	v8NodesMaxCopiesCodeGenerated = 10
	v8NodesFirstCharBudget        = 10
	v8NodesMaxCPOffset            = (1 << 15) - 1
	v8NodesMinCPOffset            = -v8NodesMaxCPOffset
)

var v8NodesNodeIsTooComplexForFixedLengthLoops = -int(^uint(0)>>1) - 1

func (node *v8NodesNodeBase) Zone() *v8NodesZone  { return node.zone }
func (node *v8NodesNodeBase) Flags() v8NodesFlags { return node.flags }
func (node *v8NodesNodeBase) IsBacktrack() bool   { return false }
func (node *v8NodesNodeBase) EatsAtLeast(notAtStart bool) uint32 {
	if notAtStart {
		return uint32(node.eatsAtLeast.fromNotStart)
	}
	return uint32(node.eatsAtLeast.fromPossiblyStart)
}
func (node *v8NodesNodeBase) Info() *v8NodesNodeInfo { return &node.info }
func (node *v8NodesNodeBase) EatsAtLeastInfo() *v8NodesEatsAtLeastInfo {
	return &node.eatsAtLeast
}
func (node *v8NodesNodeBase) SetEatsAtLeastInfo(info v8NodesEatsAtLeastInfo) {
	node.eatsAtLeast = info
}
func (node *v8NodesNodeBase) OnWorkList() bool         { return node.onWorkList }
func (node *v8NodesNodeBase) SetOnWorkList(value bool) { node.onWorkList = value }
func (node *v8NodesNodeBase) SetDoNotInline() {
	node.traceCount = v8NodesMaxCopiesCodeGenerated
}
func (node *v8NodesNodeBase) TraceCount() int { return node.traceCount }
func (node *v8NodesNodeBase) IncrementTraceCount() int {
	node.traceCount++
	return node.traceCount
}
func (node *v8NodesNodeBase) Replacement() v8NodesNode { return node.replacement }
func (node *v8NodesNodeBase) SetReplacement(replacement v8NodesNode) {
	node.replacement = replacement
}
func (node *v8NodesNodeBase) BMInfo(notAtStart bool) *v8NodesBoyerMooreLookahead {
	index := v8NodesBoolIndex(notAtStart)
	return node.bmInfo[index]
}
func (node *v8NodesNodeBase) SetBMInfo(notAtStart bool, lookahead *v8NodesBoyerMooreLookahead) {
	if lookahead == nil || !lookahead.CachesNodeInfo() {
		return
	}
	index := v8NodesBoolIndex(notAtStart)
	node.bmInfo[index] = lookahead
}
func (node *v8NodesNodeBase) KeepRecursing(compiler *v8CompilerCompiler) bool {
	return compiler != nil && !compiler.limitingRecursion && compiler.recursionDepth <= v8CompilerMaxRecursion
}
func (*v8NodesNodeBase) FixedLengthLoopLength() int {
	return v8NodesNodeIsTooComplexForFixedLengthLoops
}
func (*v8NodesNodeBase) AsEndNode() *v8NodesEndNode               { return nil }
func (*v8NodesNodeBase) AsActionNode() *v8NodesActionNode         { return nil }
func (*v8NodesNodeBase) AsChoiceNode() *v8NodesChoiceNode         { return nil }
func (*v8NodesNodeBase) AsLoopChoiceNode() *v8NodesLoopChoiceNode { return nil }
func (*v8NodesNodeBase) AsNegativeLookaroundChoiceNode() *v8NodesNegativeLookaroundChoiceNode {
	return nil
}
func (*v8NodesNodeBase) AsBackReferenceNode() *v8NodesBackReferenceNode             { return nil }
func (*v8NodesNodeBase) AsAssertionNode() *v8NodesAssertionNode                     { return nil }
func (*v8NodesNodeBase) AsTextNode() *v8NodesTextNode                               { return nil }
func (*v8NodesNodeBase) AsUnanchoredAdvanceNode() *v8NodesUnanchoredAdvanceNode     { return nil }
func (*v8NodesNodeBase) AsNegativeSubmatchSuccess() *v8NodesNegativeSubmatchSuccess { return nil }
func (*v8NodesNodeBase) AsSeqNode() *v8NodesSeqNode                                 { return nil }

const (
	v8NodesRecursionBudget = 200
)

// v8NodesSeqNode adds the mandatory success edge shared by consuming and
// deferred-action graph nodes.
type v8NodesSeqNode struct {
	v8NodesNodeBase
	onSuccess v8NodesNode
}

func v8NodesNewSeqNode(onSuccess v8NodesNode, flags v8NodesFlags) v8NodesSeqNode {
	return v8NodesSeqNode{
		v8NodesNodeBase: v8NodesNodeBase{flags: flags, zone: onSuccess.Zone()},
		onSuccess:       onSuccess,
	}
}

func (node *v8NodesSeqNode) OnSuccess() v8NodesNode        { return node.onSuccess }
func (node *v8NodesSeqNode) SetOnSuccess(next v8NodesNode) { node.onSuccess = next }
func (node *v8NodesSeqNode) AsSeqNode() *v8NodesSeqNode    { return node }

func v8NodesBoolIndex(value bool) int {
	if value {
		return 1
	}
	return 0
}

type v8NodesSeqNodeFactories interface {
	NewSeqNode(onSuccess v8NodesNode, flags v8NodesFlags) v8NodesSeqNode
}

type v8NodesActionType uint8

const (
	v8NodesSetRegisterForLoop v8NodesActionType = iota
	v8NodesIncrementRegister
	v8NodesStorePosition
	v8NodesRestorePosition
	v8NodesBeginPositiveSubmatch
	v8NodesBeginNegativeSubmatch
	v8NodesPositiveSubmatchSuccess
	v8NodesEmptyMatchCheck
	v8NodesClearCaptures
	v8NodesEatsAtLeastAction
)

// ActionNode's C++ union is expressed as named Go payload structs. Only the
// payload selected by actionType is meaningful for a given node.
type v8NodesActionSimpleData struct {
	registerFrom int
	registerTo   int
	value        int
	capture      *v8NodesCapture
}

type v8NodesActionSubmatchData struct {
	stackPointerRegister    int
	currentPositionRegister int
	clearRegisterCount      int
	clearRegisterFrom       int
	successNode             *v8NodesActionNode
}

type v8NodesActionEmptyMatchData struct {
	startRegister      int
	repetitionRegister int
	repetitionLimit    int
}

type v8NodesActionEatsAtLeastData struct{ characters int }

type v8NodesActionNode struct {
	v8NodesSeqNode
	actionType  v8NodesActionType
	simple      v8NodesActionSimpleData
	submatch    v8NodesActionSubmatchData
	emptyMatch  v8NodesActionEmptyMatchData
	eatsAtLeast v8NodesActionEatsAtLeastData
}

func v8NodesNewActionNode(actionType v8NodesActionType, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesActionNode {
	return &v8NodesActionNode{
		v8NodesSeqNode: v8NodesNewSeqNode(onSuccess, flags),
		actionType:     actionType,
	}
}

func v8NodesNewSimpleActionNode(actionType v8NodesActionType, onSuccess v8NodesNode, flags v8NodesFlags, registerFrom, registerTo, value int) *v8NodesActionNode {
	if registerTo == -1 {
		registerTo = registerFrom
	}
	node := v8NodesNewActionNode(actionType, onSuccess, flags)
	node.simple = v8NodesActionSimpleData{
		registerFrom: registerFrom,
		registerTo:   registerTo,
		value:        value,
	}
	return node
}

func (node *v8NodesActionNode) ActionType() v8NodesActionType     { return node.actionType }
func (node *v8NodesActionNode) AsActionNode() *v8NodesActionNode  { return node }
func (node *v8NodesActionNode) Accept(visitor v8NodesNodeVisitor) { visitor.VisitAction(node) }
func (node *v8NodesActionNode) SuccessNode() *v8NodesActionNode {
	if node.actionType != v8NodesBeginPositiveSubmatch {
		panic("success node is only defined for begin-positive-submatch")
	}
	return node.submatch.successNode
}
func (node *v8NodesActionNode) StoredEatsAtLeast() int {
	if node.actionType != v8NodesEatsAtLeastAction {
		panic("stored eats-at-least is only defined for eats-at-least action")
	}
	return node.eatsAtLeast.characters
}
func (node *v8NodesActionNode) Mentions(register int) bool {
	return register >= node.RegisterFrom() && register <= node.RegisterTo()
}
func (node *v8NodesActionNode) Value() int {
	if node.actionType != v8NodesSetRegisterForLoop {
		panic("value is only defined for set-register-for-loop")
	}
	return node.simple.value
}
func (node *v8NodesActionNode) IsSimpleAction() bool {
	switch node.actionType {
	case v8NodesStorePosition, v8NodesRestorePosition, v8NodesIncrementRegister,
		v8NodesSetRegisterForLoop, v8NodesClearCaptures:
		return true
	default:
		return false
	}
}
func (node *v8NodesActionNode) IsRegisterOnlyAction() bool {
	switch node.actionType {
	case v8NodesEatsAtLeastAction, v8NodesStorePosition,
		v8NodesIncrementRegister, v8NodesSetRegisterForLoop,
		v8NodesClearCaptures:
		return true
	default:
		return false
	}
}
func (node *v8NodesActionNode) RegisterFrom() int {
	if !node.IsSimpleAction() {
		panic("register range is only defined for simple actions")
	}
	return node.simple.registerFrom
}
func (node *v8NodesActionNode) RegisterTo() int { return node.simple.registerTo }

func (node *v8NodesActionNode) StoresBackreferencedCapture() bool {
	return node.actionType == v8NodesStorePosition && node.simple.capture != nil &&
		node.simple.capture.isBackreferenced
}

type v8NodesActionNodeFactories interface {
	NewActionNode(actionType v8NodesActionType, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesActionNode
	NewSimpleActionNode(actionType v8NodesActionType, onSuccess v8NodesNode, flags v8NodesFlags, registerFrom, registerTo, value int) *v8NodesActionNode
	SetRegisterForLoop(reg, value int, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesActionNode
	IncrementRegister(reg int, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesActionNode
	StorePosition(reg int, onSuccess v8NodesNode, flags v8NodesFlags, capture *v8NodesCapture) *v8NodesActionNode
	RestorePosition(reg int, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesActionNode
	ClearCaptures(interval v8NodesInterval, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesActionNode
	BeginPositiveSubmatch(stackPointerReg, positionReg int, body v8NodesNode, success *v8NodesActionNode, flags v8NodesFlags) *v8NodesActionNode
	BeginNegativeSubmatch(stackPointerReg, positionReg int, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesActionNode
	PositiveSubmatchSuccess(stackPointerReg, restoreReg, clearCaptureCount, clearCaptureFrom int, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesActionNode
	EmptyMatchCheck(startReg, repetitionReg, repetitionLimit int, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesActionNode
	EatsAtLeast(characters int, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesActionNode
}

// TextNode is a sequence of literals and/or character classes, read in the
// normal or reverse direction. Offsets are cached for fixed-loop and BM scans.
type v8NodesTextNode struct {
	v8NodesSeqNode
	elements          *v8NodesZoneList[v8NodesTextElement]
	finiteStringSet   *v8NodesFiniteStringSet
	readBackward      bool
	offsetsCalculated bool
}

func v8NodesNewTextNode(elements *v8NodesZoneList[v8NodesTextElement], readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode {
	return &v8NodesTextNode{
		v8NodesSeqNode: v8NodesNewSeqNode(onSuccess, flags),
		elements:       elements,
		readBackward:   readBackward,
	}
}

func (node *v8NodesTextNode) Elements() *v8NodesZoneList[v8NodesTextElement] { return node.elements }
func (node *v8NodesTextNode) ReadBackward() bool                             { return node.readBackward }
func (node *v8NodesTextNode) AsTextNode() *v8NodesTextNode                   { return node }
func (node *v8NodesTextNode) Accept(visitor v8NodesNodeVisitor)              { visitor.VisitText(node) }

// MakeCaseIndependent applies the legacy (non-u/non-v) case-equivalence
// closure to ordinary class ranges. Unicode case closure is already performed
// while lowering AST classes, and ECMAScript standard classes are invariant.
func (node *v8NodesTextNode) MakeCaseIndependent(oneByte bool) {
	if !node.flags.IgnoreCase || v8CompilerNeedsUnicodeCaseEquivalents(node.flags) {
		return
	}
	for index := range node.elements.items {
		element := &node.elements.items[index]
		if element.textType != v8ASTTextClassRanges {
			continue
		}
		class, ok := element.tree.(*v8ASTClassRanges)
		if !ok || class == nil {
			panic("V8 TextNode::MakeCaseIndependent found invalid class ranges")
		}
		if class.set.standardSet != nil {
			continue
		}
		v8ASTCharacterSetCanonicalize(&class.set)
		class.set.ranges = v8ASTAddLegacyCaseEquivalents(class.set.ranges, oneByte)
	}
}

// CalculateOffsets assigns each text element its fixed code-point offset from
// the start of this TextNode. Atom payloads are UTF-16 units; every class
// element consumes one code point, matching TextElement::length().
func (node *v8NodesTextNode) CalculateOffsets() {
	if node.elements == nil {
		panic("V8 TextNode has no elements")
	}
	offset := 0
	for index := range node.elements.items {
		element := &node.elements.items[index]
		element.cpOffset = offset
		offset += element.length()
	}
	node.offsetsCalculated = true
}

func (node *v8NodesTextNode) Length() int {
	if node.finiteStringSet != nil {
		// A finite string set is variable width, so its fixed text length is
		// intentionally unknown to the ordinary TextNode optimizations.
		return 0
	}
	if node.elements == nil || len(node.elements.items) == 0 {
		return 0
	}
	if !node.offsetsCalculated {
		node.CalculateOffsets()
	}
	last := node.elements.items[len(node.elements.items)-1]
	return last.cpOffset + last.length()
}

func (node *v8NodesTextNode) FixedLengthLoopLength() int {
	if node.finiteStringSet != nil {
		return v8NodesNodeIsTooComplexForFixedLengthLoops
	}
	return node.Length()
}

// GetSuccessorOfOmnivorousTextNode recognizes the exact one-class TextNode
// shape used by V8's implicit-search-loop analysis. It returns the successor
// only if this node consumes every character in the compiler's input width.
func (node *v8NodesTextNode) GetSuccessorOfOmnivorousTextNode(compiler *v8CompilerCompiler) v8NodesNode {
	if node.finiteStringSet != nil || node.readBackward || node.elements == nil || node.elements.Len() != 1 || compiler == nil {
		return nil
	}
	element := node.elements.items[0]
	if element.textType != v8ASTTextClassRanges {
		return nil
	}
	class, ok := element.tree.(*v8ASTClassRanges)
	if !ok || class == nil {
		return nil
	}
	ranges := v8ASTCanonicalizeRanges(class.set.ranges)
	if class.flags&v8ASTClassNegated != 0 {
		if len(ranges) == 0 {
			return node.onSuccess
		}
		return nil
	}
	if len(ranges) != 1 {
		return nil
	}
	maxCodeUnit := rune(0xffff)
	if compiler.oneByte {
		maxCodeUnit = 0xff
	}
	if ranges[0].isEverything(maxCodeUnit) {
		return node.onSuccess
	}
	return nil
}

type v8NodesTextNodeFactories interface {
	NewTextNode(elements *v8NodesZoneList[v8NodesTextElement], readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode
	NewTextNodeForClassRanges(ranges *v8NodesClassRanges, readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode
	CreateForCharacterRanges(zone *v8NodesZone, ranges *v8NodesZoneList[v8NodesCharacterRange], readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode
	CreateForSurrogatePair(zone *v8NodesZone, lead v8NodesCharacterRange, trailRanges *v8NodesZoneList[v8NodesCharacterRange], readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode
	CreateForSurrogatePairRanges(zone *v8NodesZone, leadRanges *v8NodesZoneList[v8NodesCharacterRange], trail v8NodesCharacterRange, readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode
}

type v8NodesAssertionType uint8

const (
	v8NodesAssertAtEnd v8NodesAssertionType = iota
	v8NodesAssertAtStart
	v8NodesAssertAtBoundary
	v8NodesAssertAtNonBoundary
	v8NodesAssertAfterNewline
)

type v8NodesAssertionNode struct {
	v8NodesSeqNode
	assertionType v8NodesAssertionType
}

func v8NodesNewAssertionNode(assertionType v8NodesAssertionType, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode {
	return &v8NodesAssertionNode{
		v8NodesSeqNode: v8NodesNewSeqNode(onSuccess, flags),
		assertionType:  assertionType,
	}
}

func (node *v8NodesAssertionNode) AssertionType() v8NodesAssertionType    { return node.assertionType }
func (node *v8NodesAssertionNode) AsAssertionNode() *v8NodesAssertionNode { return node }
func (node *v8NodesAssertionNode) Accept(visitor v8NodesNodeVisitor) {
	visitor.VisitAssertion(node)
}

type v8NodesAssertionNodeFactories interface {
	AtEnd(onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode
	AtStart(onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode
	AtBoundary(onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode
	AtNonBoundary(onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode
	AfterNewline(onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode
}

type v8NodesIfPrevious uint8

const (
	v8NodesPreviousIsNonWord v8NodesIfPrevious = iota
	v8NodesPreviousIsWord
)

type v8NodesBackReferenceNode struct {
	v8NodesSeqNode
	startRegister int
	endRegister   int
	readBackward  bool
}

func v8NodesNewBackReferenceNode(startRegister, endRegister int, readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesBackReferenceNode {
	return &v8NodesBackReferenceNode{
		v8NodesSeqNode: v8NodesNewSeqNode(onSuccess, flags),
		startRegister:  startRegister,
		endRegister:    endRegister,
		readBackward:   readBackward,
	}
}

func (node *v8NodesBackReferenceNode) StartRegister() int { return node.startRegister }
func (node *v8NodesBackReferenceNode) EndRegister() int   { return node.endRegister }
func (node *v8NodesBackReferenceNode) ReadBackward() bool { return node.readBackward }
func (node *v8NodesBackReferenceNode) AsBackReferenceNode() *v8NodesBackReferenceNode {
	return node
}
func (node *v8NodesBackReferenceNode) Accept(visitor v8NodesNodeVisitor) {
	visitor.VisitBackReference(node)
}

type v8NodesUnanchoredAdvanceNode struct{ v8NodesSeqNode }

func v8NodesNewUnanchoredAdvanceNode(onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesUnanchoredAdvanceNode {
	return &v8NodesUnanchoredAdvanceNode{v8NodesSeqNode: v8NodesNewSeqNode(onSuccess, flags)}
}
func (node *v8NodesUnanchoredAdvanceNode) AsUnanchoredAdvanceNode() *v8NodesUnanchoredAdvanceNode {
	return node
}
func (node *v8NodesUnanchoredAdvanceNode) Accept(visitor v8NodesNodeVisitor) {
	visitor.VisitUnanchoredAdvance(node)
}

type v8NodesEndAction uint8

const (
	v8NodesAccept v8NodesEndAction = iota
	v8NodesBacktrack
	v8NodesNegativeSubmatchSuccessEndAction
)

type v8NodesEndNode struct {
	v8NodesNodeBase
	action v8NodesEndAction
}

func v8NodesNewEndNode(action v8NodesEndAction, flags v8NodesFlags, zone *v8NodesZone) *v8NodesEndNode {
	node := &v8NodesEndNode{v8NodesNodeBase: v8NodesNodeBase{flags: flags, zone: zone}, action: action}
	if action == v8NodesBacktrack {
		node.SetEatsAtLeastInfo(v8NodesNewEatsAtLeastInfo(v8NodesLargeEatsAtLeastValue))
	}
	return node
}

func (node *v8NodesEndNode) Action() v8NodesEndAction          { return node.action }
func (node *v8NodesEndNode) IsBacktrack() bool                 { return node.action == v8NodesBacktrack }
func (node *v8NodesEndNode) AsEndNode() *v8NodesEndNode        { return node }
func (node *v8NodesEndNode) Accept(visitor v8NodesNodeVisitor) { visitor.VisitEnd(node) }

type v8NodesNegativeSubmatchSuccess struct {
	v8NodesEndNode
	stackPointerRegister    int
	currentPositionRegister int
	clearCaptureCount       int
	clearCaptureStart       int
}

func v8NodesNewNegativeSubmatchSuccess(stackPointerRegister, positionRegister, clearCaptureCount, clearCaptureStart int, flags v8NodesFlags, zone *v8NodesZone) *v8NodesNegativeSubmatchSuccess {
	return &v8NodesNegativeSubmatchSuccess{
		v8NodesEndNode:          *v8NodesNewEndNode(v8NodesNegativeSubmatchSuccessEndAction, flags, zone),
		stackPointerRegister:    stackPointerRegister,
		currentPositionRegister: positionRegister,
		clearCaptureCount:       clearCaptureCount,
		clearCaptureStart:       clearCaptureStart,
	}
}
func (node *v8NodesNegativeSubmatchSuccess) AsNegativeSubmatchSuccess() *v8NodesNegativeSubmatchSuccess {
	return node
}

type v8NodesEndNodeFactories interface {
	NewEndNode(action v8NodesEndAction, flags v8NodesFlags, zone *v8NodesZone) *v8NodesEndNode
}

// v8NodesGraphNodeFactory ports the small inline/static constructors in
// regexp-nodes.h. Node semantics live on the concrete types; this factory
// centralizes the constructor choices used by AST-to-graph lowering.
type v8NodesGraphNodeFactory struct{}

var (
	_ v8NodesSeqNodeFactories        = v8NodesGraphNodeFactory{}
	_ v8NodesChoiceNodeFactories     = v8NodesGraphNodeFactory{}
	_ v8NodesLoopChoiceNodeFactories = v8NodesGraphNodeFactory{}
	_ v8NodesEndNodeFactories        = v8NodesGraphNodeFactory{}
	_ v8NodesAssertionNodeFactories  = v8NodesGraphNodeFactory{}
)

func (v8NodesGraphNodeFactory) NewChoiceNode(expectedSize int, flags v8NodesFlags, zone *v8NodesZone) *v8NodesChoiceNode {
	return v8NodesNewChoiceNode(expectedSize, flags, zone)
}

func (v8NodesGraphNodeFactory) NewSeqNode(onSuccess v8NodesNode, flags v8NodesFlags) v8NodesSeqNode {
	return v8NodesNewSeqNode(onSuccess, flags)
}

func (v8NodesGraphNodeFactory) NewLoopChoiceNode(bodyCanBeZeroLength, readBackward bool, flags v8NodesFlags, zone *v8NodesZone) *v8NodesLoopChoiceNode {
	return v8NodesNewLoopChoiceNode(bodyCanBeZeroLength, readBackward, flags, zone)
}

func (v8NodesGraphNodeFactory) NewEndNode(action v8NodesEndAction, flags v8NodesFlags, zone *v8NodesZone) *v8NodesEndNode {
	return v8NodesNewEndNode(action, flags, zone)
}

func (v8NodesGraphNodeFactory) AtEnd(success v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode {
	return v8NodesNewAssertionNode(v8NodesAssertAtEnd, success, flags)
}

func (v8NodesGraphNodeFactory) AtStart(success v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode {
	return v8NodesNewAssertionNode(v8NodesAssertAtStart, success, flags)
}

func (v8NodesGraphNodeFactory) AtBoundary(success v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode {
	return v8NodesNewAssertionNode(v8NodesAssertAtBoundary, success, flags)
}

func (v8NodesGraphNodeFactory) AtNonBoundary(success v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode {
	return v8NodesNewAssertionNode(v8NodesAssertAtNonBoundary, success, flags)
}

func (v8NodesGraphNodeFactory) AfterNewline(success v8NodesNode, flags v8NodesFlags) *v8NodesAssertionNode {
	return v8NodesNewAssertionNode(v8NodesAssertAfterNewline, success, flags)
}

func (v8NodesGraphNodeFactory) NewNegativeLookaroundChoiceNode(thisMustFail, thenDoThis v8NodesGuardedAlternative, flags v8NodesFlags, zone *v8NodesZone) *v8NodesNegativeLookaroundChoiceNode {
	return v8NodesNewNegativeLookaroundChoiceNode(thisMustFail, thenDoThis, flags, zone)
}

type v8NodesGuardRelation uint8

const (
	v8NodesGuardLT v8NodesGuardRelation = iota
	v8NodesGuardGEQ
)

type v8NodesGuard struct {
	register int
	relation v8NodesGuardRelation
	value    int
}

func v8NodesNewGuard(register int, relation v8NodesGuardRelation, value int) *v8NodesGuard {
	return &v8NodesGuard{register: register, relation: relation, value: value}
}
func (guard *v8NodesGuard) Register() int                  { return guard.register }
func (guard *v8NodesGuard) Relation() v8NodesGuardRelation { return guard.relation }
func (guard *v8NodesGuard) Value() int                     { return guard.value }

type v8NodesGuardAPI interface {
	Register() int
	Relation() v8NodesGuardRelation
	Value() int
}

type v8NodesGuardedAlternative struct {
	node   v8NodesNode
	guards *v8NodesZoneList[*v8NodesGuard]
}

func v8NodesNewGuardedAlternative(node v8NodesNode) v8NodesGuardedAlternative {
	return v8NodesGuardedAlternative{node: node}
}
func (alternative *v8NodesGuardedAlternative) AddGuard(guard *v8NodesGuard, _ *v8NodesZone) {
	if alternative.guards == nil {
		alternative.guards = &v8NodesZoneList[*v8NodesGuard]{items: make([]*v8NodesGuard, 0, 1)}
	}
	alternative.guards.Add(guard)
}
func (alternative *v8NodesGuardedAlternative) Node() v8NodesNode        { return alternative.node }
func (alternative *v8NodesGuardedAlternative) SetNode(node v8NodesNode) { alternative.node = node }
func (alternative *v8NodesGuardedAlternative) Guards() *v8NodesZoneList[*v8NodesGuard] {
	return alternative.guards
}

type v8NodesGuardedAlternativeAPI interface {
	AddGuard(*v8NodesGuard, *v8NodesZone)
	Node() v8NodesNode
	SetNode(v8NodesNode)
	Guards() *v8NodesZoneList[*v8NodesGuard]
}

// ChoiceNode owns ordered guarded alternatives and the compiler's search,
// preload, quick-check, and fixed-length-loop emission contracts.
type v8NodesChoiceNode struct {
	v8NodesNodeBase
	alternatives    *v8NodesZoneList[v8NodesGuardedAlternative]
	notAtStart      bool
	beingCalculated bool
}

func v8NodesNewChoiceNode(expectedSize int, flags v8NodesFlags, zone *v8NodesZone) *v8NodesChoiceNode {
	return &v8NodesChoiceNode{
		v8NodesNodeBase: v8NodesNodeBase{flags: flags, zone: zone},
		alternatives:    &v8NodesZoneList[v8NodesGuardedAlternative]{items: make([]v8NodesGuardedAlternative, 0, expectedSize)},
	}
}

func (node *v8NodesChoiceNode) AddAlternative(alternative v8NodesGuardedAlternative) {
	node.alternatives.Add(alternative)
}
func (node *v8NodesChoiceNode) Alternatives() *v8NodesZoneList[v8NodesGuardedAlternative] {
	return node.alternatives
}
func (node *v8NodesChoiceNode) BeingCalculated() bool                  { return node.beingCalculated }
func (node *v8NodesChoiceNode) NotAtStart() bool                       { return node.notAtStart }
func (node *v8NodesChoiceNode) SetNotAtStart()                         { node.notAtStart = true }
func (node *v8NodesChoiceNode) SetBeingCalculated(value bool)          { node.beingCalculated = value }
func (node *v8NodesChoiceNode) ReadBackward() bool                     { return false }
func (node *v8NodesChoiceNode) AsChoiceNode() *v8NodesChoiceNode       { return node }
func (node *v8NodesChoiceNode) Accept(visitor v8NodesNodeVisitor)      { visitor.VisitChoice(node) }
func (*v8NodesChoiceNode) TryToEmitQuickCheckForAlternative(bool) bool { return true }

// fixedLengthLoopLengthForAlternative recognizes the fixed-width sequence
// from one alternative back to this choice. It preserves V8's recursion and
// signed code-position bounds because these values feed the fixed-loop graph
// optimization even when the final backend emits Go source.
func (node *v8NodesChoiceNode) fixedLengthLoopLengthForAlternative(alternative *v8NodesGuardedAlternative) int {
	if alternative == nil || alternative.Node() == nil {
		return v8NodesNodeIsTooComplexForFixedLengthLoops
	}
	length := 0
	current := alternative.Node()
	depth := 0
	for current != node {
		if depth > v8CompilerMaxRecursion {
			return v8NodesNodeIsTooComplexForFixedLengthLoops
		}
		depth++
		nodeLength := current.FixedLengthLoopLength()
		if nodeLength == v8NodesNodeIsTooComplexForFixedLengthLoops {
			return v8NodesNodeIsTooComplexForFixedLengthLoops
		}
		length += nodeLength
		seq := current.AsSeqNode()
		if seq == nil || seq.OnSuccess() == nil {
			return v8NodesNodeIsTooComplexForFixedLengthLoops
		}
		current = seq.OnSuccess()
	}
	if node.ReadBackward() {
		length = -length
	}
	if length < v8NodesMinCPOffset || length > v8NodesMaxCPOffset {
		return v8NodesNodeIsTooComplexForFixedLengthLoops
	}
	return length
}

type v8NodesNegativeLookaroundChoiceNode struct{ v8NodesChoiceNode }

func v8NodesNewNegativeLookaroundChoiceNode(thisMustFail, thenDoThis v8NodesGuardedAlternative, flags v8NodesFlags, zone *v8NodesZone) *v8NodesNegativeLookaroundChoiceNode {
	node := &v8NodesNegativeLookaroundChoiceNode{v8NodesChoiceNode: *v8NodesNewChoiceNode(2, flags, zone)}
	node.AddAlternative(thisMustFail)
	node.AddAlternative(thenDoThis)
	return node
}

func (node *v8NodesNegativeLookaroundChoiceNode) LookaroundNode() v8NodesNode {
	return node.alternatives.At(v8NodesLookaroundIndex).node
}
func (node *v8NodesNegativeLookaroundChoiceNode) ContinueNode() v8NodesNode {
	return node.alternatives.At(v8NodesContinueIndex).node
}
func (*v8NodesNegativeLookaroundChoiceNode) TryToEmitQuickCheckForAlternative(isFirst bool) bool {
	return !isFirst
}
func (node *v8NodesNegativeLookaroundChoiceNode) AsNegativeLookaroundChoiceNode() *v8NodesNegativeLookaroundChoiceNode {
	return node
}
func (node *v8NodesNegativeLookaroundChoiceNode) Accept(visitor v8NodesNodeVisitor) {
	visitor.VisitNegativeLookaroundChoice(node)
}

type v8NodesChoiceNodeFactories interface {
	NewChoiceNode(expectedSize int, flags v8NodesFlags, zone *v8NodesZone) *v8NodesChoiceNode
}

type v8NodesLoopChoiceNode struct {
	v8NodesChoiceNode
	loopNode            v8NodesNode
	continueNode        v8NodesNode
	bodyCanBeZeroLength bool
	// A repeated single-code-point class has no capture side effects and
	// chooses its match solely from the current input code point. Its completed
	// iteration's internal alternatives can therefore be discarded when an
	// end-anchored continuation is retried.
	singleCodePointClassBody bool
	readBackward             bool
	atomicLoopKindValid      bool
	atomicLoopKind           v8NodesAtomicLoopKind
}

func v8NodesNewLoopChoiceNode(bodyCanBeZeroLength, readBackward bool, flags v8NodesFlags, zone *v8NodesZone) *v8NodesLoopChoiceNode {
	return &v8NodesLoopChoiceNode{
		v8NodesChoiceNode:   *v8NodesNewChoiceNode(2, flags, zone),
		bodyCanBeZeroLength: bodyCanBeZeroLength,
		readBackward:        readBackward,
	}
}

func (node *v8NodesLoopChoiceNode) AddLoopAlternative(alternative v8NodesGuardedAlternative) {
	node.loopNode = alternative.node
	node.AddAlternative(alternative)
}
func (node *v8NodesLoopChoiceNode) AddContinueAlternative(alternative v8NodesGuardedAlternative) {
	node.continueNode = alternative.node
	node.AddAlternative(alternative)
}
func (node *v8NodesLoopChoiceNode) LoopNode() v8NodesNode     { return node.loopNode }
func (node *v8NodesLoopChoiceNode) ContinueNode() v8NodesNode { return node.continueNode }
func (node *v8NodesLoopChoiceNode) BodyCanBeZeroLength() bool {
	return node.bodyCanBeZeroLength
}
func (node *v8NodesLoopChoiceNode) SingleCodePointClassBody() bool {
	return node.singleCodePointClassBody
}
func (node *v8NodesLoopChoiceNode) SetSingleCodePointClassBody(value bool) {
	node.singleCodePointClassBody = value
}
func (node *v8NodesLoopChoiceNode) ReadBackward() bool                       { return node.readBackward }
func (node *v8NodesLoopChoiceNode) AsLoopChoiceNode() *v8NodesLoopChoiceNode { return node }
func (node *v8NodesLoopChoiceNode) FixedLengthBodyIterationLength() int {
	if node.alternatives == nil || node.alternatives.Len() == 0 {
		return v8NodesNodeIsTooComplexForFixedLengthLoops
	}
	return node.fixedLengthLoopLengthForAlternative(node.alternatives.At(0))
}
func (node *v8NodesLoopChoiceNode) AtomicLoopKindCache() (v8NodesAtomicLoopKind, bool) {
	return node.atomicLoopKind, node.atomicLoopKindValid
}
func (node *v8NodesLoopChoiceNode) AtomicLoopKind() v8NodesAtomicLoopKind {
	return v8NodesLoopChoiceAtomicLoopKind(node)
}
func (node *v8NodesLoopChoiceNode) SetAtomicLoopKindCache(kind v8NodesAtomicLoopKind) {
	node.atomicLoopKind = kind
	node.atomicLoopKindValid = true
}
func (node *v8NodesLoopChoiceNode) Accept(visitor v8NodesNodeVisitor) {
	visitor.VisitLoopChoice(node)
}

type v8NodesLoopChoiceNodeFactories interface {
	NewLoopChoiceNode(bodyCanBeZeroLength, readBackward bool, flags v8NodesFlags, zone *v8NodesZone) *v8NodesLoopChoiceNode
}

const (
	v8NodesLookaroundIndex = 0
	v8NodesContinueIndex   = 1
)

// v8NodesNodeVisitor preserves FOR_EACH_NODE_TYPE's complete visitor surface.
type v8NodesNodeVisitor interface {
	VisitEnd(*v8NodesEndNode)
	VisitAction(*v8NodesActionNode)
	VisitChoice(*v8NodesChoiceNode)
	VisitLoopChoice(*v8NodesLoopChoiceNode)
	VisitNegativeLookaroundChoice(*v8NodesNegativeLookaroundChoiceNode)
	VisitBackReference(*v8NodesBackReferenceNode)
	VisitAssertion(*v8NodesAssertionNode)
	VisitText(*v8NodesTextNode)
	VisitUnanchoredAdvance(*v8NodesUnanchoredAdvanceNode)
}

// v8NodesNodeCasts preserves the default-null virtual As* casts. Concrete node
// implementations override exactly one ordinary node-kind cast; seq nodes
// also override AsSeqNode, and NegativeSubmatchSuccess has its own cast.
type v8NodesNodeCasts interface {
	AsEndNode() *v8NodesEndNode
	AsActionNode() *v8NodesActionNode
	AsChoiceNode() *v8NodesChoiceNode
	AsLoopChoiceNode() *v8NodesLoopChoiceNode
	AsNegativeLookaroundChoiceNode() *v8NodesNegativeLookaroundChoiceNode
	AsBackReferenceNode() *v8NodesBackReferenceNode
	AsAssertionNode() *v8NodesAssertionNode
	AsTextNode() *v8NodesTextNode
	AsUnanchoredAdvanceNode() *v8NodesUnanchoredAdvanceNode
	AsNegativeSubmatchSuccess() *v8NodesNegativeSubmatchSuccess
	AsSeqNode() *v8NodesSeqNode
}

// Ported V8 source declarations from regexp-compiler-tonode.
// Copyright 2019 the V8 project authors. All rights reserved.
// V8 RegExp AST-to-graph lowering adapted from regexp-compiler-tonode.cc.
// See LICENSE.v8 for attribution and license terms.
const (
	v8TonodeMaxCodePoint        = rune(0x10ffff)
	v8TonodeMaxUTF16CodeUnit    = rune(0xffff)
	v8TonodeMaxUTF16CodeUnitU   = uint32(0xffff)
	v8TonodeRangeEndMarker      = v8CompilerRangeEndMarker
	v8TonodeMaxRangesToInline   = 32
	v8TonodeMaxUnrolledMinMatch = 3
	v8TonodeMaxUnrolledMaxMatch = 3
	v8TonodeMaxExpansionFactor  = 6
	v8TonodeExpansionOverflow   = v8TonodeMaxExpansionFactor + 1
)

const (
	v8TonodeLeadSurrogateStart  = rune(0xd800)
	v8TonodeLeadSurrogateEnd    = rune(0xdbff)
	v8TonodeTrailSurrogateStart = rune(0xdc00)
	v8TonodeTrailSurrogateEnd   = rune(0xdfff)
	v8TonodeNonBMPStart         = rune(0x10000)
	v8TonodeNonBMPEnd           = rune(0x10ffff)
)

// The source keeps compiler_constants.h's sorted, half-open tables. The
// marker and endpoint convention are retained here so AddClass and
// AddClassNegated can be compared directly with the V8 routines.
var (
	v8TonodeSpaceRanges          = v8CompilerSpaceRanges[:]
	v8TonodeWordRanges           = v8CompilerWordRanges[:]
	v8TonodeDigitRanges          = v8CompilerDigitRanges[:]
	v8TonodeLineTerminatorRanges = v8CompilerLineTerminatorRanges[:]
)

type v8TonodeCompiler = v8ASTCompiler
type v8TonodeNode = v8ASTNode
type v8TonodeTree = v8ASTTree

// v8TonodeUnicodeRangeSplitter preserves the four output buckets from the
// V8 class. The two ordinary BMP regions share bmp, while lone surrogates
// remain separate so lowering can add the required lookaround guards.
type v8TonodeUnicodeRangeSplitter struct {
	bmp             []v8ASTCharacterRange
	leadSurrogates  []v8ASTCharacterRange
	trailSurrogates []v8ASTCharacterRange
	nonBMP          []v8ASTCharacterRange
}

func v8TonodeNewUnicodeRangeSplitter(base []v8ASTCharacterRange) *v8TonodeUnicodeRangeSplitter {
	splitter := &v8TonodeUnicodeRangeSplitter{}
	for _, current := range base {
		splitter.addRange(current)
	}
	return splitter
}

func (splitter *v8TonodeUnicodeRangeSplitter) addRange(current v8ASTCharacterRange) {
	type partition struct {
		from, to rune
		bucket   *[]v8ASTCharacterRange
	}
	parts := [...]partition{
		{0, v8TonodeLeadSurrogateStart - 1, &splitter.bmp},
		{v8TonodeLeadSurrogateStart, v8TonodeLeadSurrogateEnd, &splitter.leadSurrogates},
		{v8TonodeTrailSurrogateStart, v8TonodeTrailSurrogateEnd, &splitter.trailSurrogates},
		{v8TonodeTrailSurrogateEnd + 1, v8TonodeNonBMPStart - 1, &splitter.bmp},
		{v8TonodeNonBMPStart, v8TonodeNonBMPEnd, &splitter.nonBMP},
	}
	for _, part := range parts {
		if part.from > current.to {
			break
		}
		from, to := part.from, part.to
		if current.from > from {
			from = current.from
		}
		if current.to < to {
			to = current.to
		}
		if from <= to {
			*part.bucket = append(*part.bucket, v8ASTRange(from, to))
		}
	}
}

func (splitter *v8TonodeUnicodeRangeSplitter) bmpRanges() []v8ASTCharacterRange {
	return splitter.bmp
}

func (splitter *v8TonodeUnicodeRangeSplitter) leadSurrogateRanges() []v8ASTCharacterRange {
	return splitter.leadSurrogates
}

func (splitter *v8TonodeUnicodeRangeSplitter) trailSurrogateRanges() []v8ASTCharacterRange {
	return splitter.trailSurrogates
}

func (splitter *v8TonodeUnicodeRangeSplitter) nonBMPRanges() []v8ASTCharacterRange {
	return splitter.nonBMP
}

// v8TonodeRangeVector is the Go slice counterpart of V8's small-vector range
// storage. Graph construction consumes these ranges but is separately scoped.
type v8TonodeRangeVector []v8ASTCharacterRange

type v8TonodeUC16Range uint32

func v8TonodePackUC16Range(from, to uint16) v8TonodeUC16Range {
	return v8TonodeUC16Range(uint32(from)<<16 | uint32(to))
}

func v8TonodeUC16RangeFrom(packed v8TonodeUC16Range) uint16 {
	return uint16(uint32(packed) >> 16)
}

func v8TonodeUC16RangeTo(packed v8TonodeUC16Range) uint16 {
	return uint16(packed)
}

type v8TonodeSurrogatePairGroup struct {
	leading  v8ASTCharacterRange
	trailing []v8ASTCharacterRange
}

type v8TonodeSurrogatePairPlan struct {
	leadingWithFullTrailing []v8ASTCharacterRange
	groupedByLeading        map[v8TonodeUC16Range][]v8ASTCharacterRange
}

// v8TonodeSurrogateParts maps a supplementary code point to its UTF-16 pair.
func v8TonodeSurrogateParts(codePoint rune) (lead, trail uint16) {
	value := uint32(codePoint - v8TonodeNonBMPStart)
	return uint16(uint32(v8TonodeLeadSurrogateStart) + (value >> 10)),
		uint16(uint32(v8TonodeTrailSurrogateStart) + (value & 0x3ff))
}

// v8TonodeGroupNonBMPRanges preserves V8's grouping rule: ranges with the
// same lead-surrogate interval share one trailing-range list, and full
// trailing spans are collected in a separate group.
func v8TonodeGroupNonBMPRanges(ranges []v8ASTCharacterRange) v8TonodeSurrogatePairPlan {
	plan := v8TonodeSurrogatePairPlan{
		groupedByLeading: make(map[v8TonodeUC16Range][]v8ASTCharacterRange),
	}
	add := func(fromLead, toLead, fromTrail, toTrail uint16) {
		key := v8TonodePackUC16Range(fromLead, toLead)
		trail := v8ASTRange(rune(fromTrail), rune(toTrail))
		if fromTrail == uint16(v8TonodeTrailSurrogateStart) && toTrail == uint16(v8TonodeTrailSurrogateEnd) {
			plan.leadingWithFullTrailing = append(plan.leadingWithFullTrailing, v8ASTRange(rune(fromLead), rune(toLead)))
			return
		}
		plan.groupedByLeading[key] = append(plan.groupedByLeading[key], trail)
	}

	for _, interval := range v8TonodeCanonicalizeRanges(ranges) {
		if interval.from < v8TonodeNonBMPStart || interval.to > v8TonodeNonBMPEnd {
			continue
		}
		fromLead, fromTrail := v8TonodeSurrogateParts(interval.from)
		toLead, toTrail := v8TonodeSurrogateParts(interval.to)
		if fromLead == toLead {
			add(fromLead, toLead, fromTrail, toTrail)
			continue
		}
		if fromTrail != uint16(v8TonodeTrailSurrogateStart) {
			add(fromLead, fromLead, fromTrail, uint16(v8TonodeTrailSurrogateEnd))
			fromLead++
		}
		if toTrail != uint16(v8TonodeTrailSurrogateEnd) {
			add(toLead, toLead, uint16(v8TonodeTrailSurrogateStart), toTrail)
			toLead--
		}
		if fromLead <= toLead {
			add(fromLead, toLead, uint16(v8TonodeTrailSurrogateStart), uint16(v8TonodeTrailSurrogateEnd))
		}
	}
	plan.leadingWithFullTrailing = v8TonodeCanonicalizeRanges(plan.leadingWithFullTrailing)
	for key, trailing := range plan.groupedByLeading {
		plan.groupedByLeading[key] = v8TonodeCanonicalizeRanges(trailing)
	}
	return plan
}

func v8TonodeCompareRanges(ranges []v8ASTCharacterRange, special []int) bool {
	if len(special) == 0 || special[len(special)-1] != v8TonodeRangeEndMarker {
		return false
	}
	count := len(special) - 1
	if len(ranges)*2 != count {
		return false
	}
	for i := 0; i < count; i += 2 {
		rangeAt := ranges[i/2]
		if rangeAt.from != rune(special[i]) || rangeAt.to != rune(special[i+1]-1) {
			return false
		}
	}
	return true
}

func v8TonodeCompareInverseRanges(ranges []v8ASTCharacterRange, special []int) bool {
	if len(special) < 3 || special[len(special)-1] != v8TonodeRangeEndMarker || special[0] == 0 {
		return false
	}
	count := len(special) - 1
	if len(ranges) != (count>>1)+1 || ranges[0].from != 0 {
		return false
	}
	current := ranges[0]
	for i := 0; i < count; i += 2 {
		if rune(special[i]) != current.to+1 {
			return false
		}
		current = ranges[(i>>1)+1]
		if rune(special[i+1]) != current.from {
			return false
		}
	}
	return current.to == v8TonodeMaxCodePoint
}

func v8TonodeAddClass(special []int, ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	if len(special) == 0 || special[len(special)-1] != v8TonodeRangeEndMarker {
		return ranges
	}
	for i := 0; i+1 < len(special)-1; i += 2 {
		if special[i] < special[i+1] {
			ranges = append(ranges, v8ASTRange(rune(special[i]), rune(special[i+1]-1)))
		}
	}
	return ranges
}

func v8TonodeAddClassNegated(special []int, ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	if len(special) == 0 || special[len(special)-1] != v8TonodeRangeEndMarker || special[0] == 0 {
		return ranges
	}
	last := 0
	for i := 0; i+1 < len(special)-1; i += 2 {
		if last <= special[i]-1 && special[i] < special[i+1] {
			ranges = append(ranges, v8ASTRange(rune(last), rune(special[i]-1)))
		}
		last = special[i+1]
	}
	return append(ranges, v8ASTRange(rune(last), v8TonodeMaxCodePoint))
}

func v8TonodeAddClassEscape(set v8ASTStandardCharacterSet, ranges []v8ASTCharacterRange, unicodeCaseClosure bool) []v8ASTCharacterRange {
	if unicodeCaseClosure && (set == v8ASTWord || set == v8ASTNotWord) {
		closed := v8ASTCaseFoldingCloseOver(v8TonodeAddClass(v8TonodeWordRanges, nil), v8TonodeUnicodeCaseFolding)
		if set == v8ASTNotWord {
			closed = v8ASTNegateRanges(closed)
		}
		return append(ranges, closed...)
	}
	switch set {
	case v8ASTWhitespace:
		return v8TonodeAddClass(v8TonodeSpaceRanges, ranges)
	case v8ASTNotWhitespace:
		return v8TonodeAddClassNegated(v8TonodeSpaceRanges, ranges)
	case v8ASTWord:
		return v8TonodeAddClass(v8TonodeWordRanges, ranges)
	case v8ASTNotWord:
		return v8TonodeAddClassNegated(v8TonodeWordRanges, ranges)
	case v8ASTDigit:
		return v8TonodeAddClass(v8TonodeDigitRanges, ranges)
	case v8ASTNotDigit:
		return v8TonodeAddClassNegated(v8TonodeDigitRanges, ranges)
	case v8ASTLineTerminator:
		return v8TonodeAddClass(v8TonodeLineTerminatorRanges, ranges)
	case v8ASTNotLineTerminator:
		return v8TonodeAddClassNegated(v8TonodeLineTerminatorRanges, ranges)
	case v8ASTEverything:
		return append(ranges, v8ASTRangeEverything())
	default:
		return ranges
	}
}

func v8TonodeClassRangesIsStandard(classRanges *v8ASTClassRanges) bool {
	if classRanges == nil || classRanges.flags&v8ASTClassNegated != 0 {
		return false
	}
	if classRanges.set.standardSet != nil {
		return true
	}
	v8TonodeCharacterSetRanges(&classRanges.set)
	_, ok := v8TonodeIsStandardClass(&classRanges.set, false)
	return ok
}

func v8TonodeIsStandardClass(set *v8ASTCharacterSet, negated bool) (v8ASTStandardCharacterSet, bool) {
	if set == nil {
		return 0, false
	}
	if negated {
		return 0, false
	}
	if set.standardSet != nil {
		return *set.standardSet, true
	}
	candidates := []struct {
		set    v8ASTStandardCharacterSet
		ranges []int
		invert bool
	}{
		{v8ASTWhitespace, v8TonodeSpaceRanges, false},
		{v8ASTNotWhitespace, v8TonodeSpaceRanges, true},
		{v8ASTNotLineTerminator, v8TonodeLineTerminatorRanges, true},
		{v8ASTLineTerminator, v8TonodeLineTerminatorRanges, false},
		{v8ASTWord, v8TonodeWordRanges, false},
		{v8ASTNotWord, v8TonodeWordRanges, true},
	}
	for _, candidate := range candidates {
		matches := v8TonodeCompareRanges(set.ranges, candidate.ranges)
		if candidate.invert {
			matches = v8TonodeCompareInverseRanges(set.ranges, candidate.ranges)
		}
		if matches {
			set.standardSet = &candidate.set
			return candidate.set, true
		}
	}
	return 0, false
}

func v8TonodeCharacterSetRanges(set *v8ASTCharacterSet) []v8ASTCharacterRange {
	if set.ranges == nil {
		if set.standardSet == nil {
			panic("standard character set is not initialized")
		}
		set.ranges = v8TonodeAddClassEscape(*set.standardSet, nil, false)
	}
	return set.ranges
}

func v8TonodeCanonicalizeCharacterSet(set *v8ASTCharacterSet) {
	if set != nil && set.ranges != nil {
		set.ranges = v8TonodeCanonicalizeRanges(set.ranges)
	}
}

func v8TonodeRangesAreCanonical(ranges []v8ASTCharacterRange) bool {
	if len(ranges) <= 1 {
		return true
	}
	max := ranges[0].to
	for i := 1; i < len(ranges); i++ {
		if ranges[i].from <= max+1 {
			return false
		}
		max = ranges[i].to
	}
	return true
}

func v8TonodeMoveRanges(ranges []v8ASTCharacterRange, from, to, count int) {
	if from < to {
		for i := count - 1; i >= 0; i-- {
			ranges[to+i] = ranges[from+i]
		}
	} else {
		for i := 0; i < count; i++ {
			ranges[to+i] = ranges[from+i]
		}
	}
}

func v8TonodeInsertRangeInCanonicalList(ranges []v8ASTCharacterRange, count int, insert v8ASTCharacterRange) ([]v8ASTCharacterRange, int) {
	from, to := insert.from, insert.to
	startPos, endPos := 0, count
	for i := count - 1; i >= 0; i-- {
		current := ranges[i]
		if current.from > to+1 {
			endPos = i
		} else if current.to+1 < from {
			startPos = i + 1
			break
		}
	}
	if startPos == endPos {
		if count >= len(ranges) {
			ranges = append(ranges, v8ASTCharacterRange{})
		}
		if startPos < count {
			v8TonodeMoveRanges(ranges, startPos, startPos+1, count-startPos)
		}
		ranges[startPos] = insert
		return ranges, count + 1
	}
	if startPos+1 == endPos {
		existing := ranges[startPos]
		if existing.from < from {
			from = existing.from
		}
		if existing.to > to {
			to = existing.to
		}
		ranges[startPos] = v8ASTRange(from, to)
		return ranges, count
	}
	newFrom := ranges[startPos].from
	if from < newFrom {
		newFrom = from
	}
	newTo := ranges[endPos-1].to
	if to > newTo {
		newTo = to
	}
	if endPos < count {
		v8TonodeMoveRanges(ranges, endPos, startPos+1, count-endPos)
	}
	ranges[startPos] = v8ASTRange(newFrom, newTo)
	return ranges, count - (endPos - startPos) + 1
}

func v8TonodeCanonicalizeRanges(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	if len(ranges) <= 1 {
		return ranges
	}
	n := len(ranges)
	max, i := ranges[0].to, 1
	for i < n {
		current := ranges[i]
		if current.from <= max+1 {
			break
		}
		max = current.to
		i++
	}
	if i == n {
		return ranges
	}
	read, canonicalCount := i, i
	for read < n {
		ranges, canonicalCount = v8TonodeInsertRangeInCanonicalList(ranges, canonicalCount, ranges[read])
		read++
	}
	return ranges[:canonicalCount]
}

func v8TonodeNegateRanges(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	if !v8TonodeRangesAreCanonical(ranges) {
		panic("negate requires canonical ranges")
	}
	result := make([]v8ASTCharacterRange, 0, len(ranges)+1)
	from, index := rune(0), 0
	if len(ranges) > 0 && ranges[0].from == 0 {
		from = ranges[0].to + 1
		index = 1
	}
	for ; index < len(ranges); index++ {
		current := ranges[index]
		result = append(result, v8ASTRange(from, current.from-1))
		from = current.to + 1
	}
	if from < v8TonodeMaxCodePoint {
		result = append(result, v8ASTRange(from, v8TonodeMaxCodePoint))
	}
	return result
}

func v8TonodeIntersectRanges(lhs, rhs []v8ASTCharacterRange) []v8ASTCharacterRange {
	if !v8TonodeRangesAreCanonical(lhs) || !v8TonodeRangesAreCanonical(rhs) {
		panic("intersection requires canonical ranges")
	}
	result := make([]v8ASTCharacterRange, 0)
	for i, j := 0, 0; i < len(lhs) && j < len(rhs); {
		if lhs[i].to < rhs[j].from {
			i++
			continue
		}
		if rhs[j].to < lhs[i].from {
			j++
			continue
		}
		from, to := lhs[i].from, lhs[i].to
		if rhs[j].from > from {
			from = rhs[j].from
		}
		if rhs[j].to < to {
			to = rhs[j].to
		}
		result = append(result, v8ASTRange(from, to))
		if to == lhs[i].to {
			i++
		} else {
			j++
		}
	}
	return result
}

func v8TonodeRangesIntersect(lhs, rhs []v8ASTCharacterRange) bool {
	if !v8TonodeRangesAreCanonical(lhs) || !v8TonodeRangesAreCanonical(rhs) {
		panic("intersection requires canonical ranges")
	}
	for i, j := 0, 0; i < len(lhs) && j < len(rhs); {
		if lhs[i].to < rhs[j].from {
			i++
		} else if rhs[j].to < lhs[i].from {
			j++
		} else {
			return true
		}
	}
	return false
}

func v8TonodeSafeAdvanceRange(ranges []v8ASTCharacterRange, index int, from, to rune) (int, rune, rune) {
	index++
	if index < len(ranges) {
		return index, ranges[index].from, ranges[index].to
	}
	return index, v8TonodeMaxCodePoint + 1, to
}

func v8TonodeSubtractRanges(src, toRemove []v8ASTCharacterRange) []v8ASTCharacterRange {
	if !v8TonodeRangesAreCanonical(src) || !v8TonodeRangesAreCanonical(toRemove) {
		panic("subtract requires canonical ranges")
	}
	if len(src) == 0 {
		return nil
	}
	result := make([]v8ASTCharacterRange, 0, len(src))
	srcIndex, removeIndex := 0, 0
	from, to := src[0].from, src[0].to
	for srcIndex < len(src) && removeIndex < len(toRemove) {
		remove := toRemove[removeIndex]
		switch {
		case remove.to < from:
			removeIndex++
		case to < remove.from:
			result = append(result, v8ASTRange(from, to))
			srcIndex, from, to = v8TonodeSafeAdvanceRange(src, srcIndex, from, to)
		case from >= remove.from && to <= remove.to:
			srcIndex, from, to = v8TonodeSafeAdvanceRange(src, srcIndex, from, to)
		case from < remove.from && to > remove.to:
			result = append(result, v8ASTRange(from, remove.from-1))
			from = remove.to + 1
			removeIndex++
		case from < remove.from:
			result = append(result, v8ASTRange(from, remove.from-1))
			srcIndex, from, to = v8TonodeSafeAdvanceRange(src, srcIndex, from, to)
		case to > remove.to:
			from = remove.to + 1
			removeIndex++
		default:
			panic("unreachable range subtraction state")
		}
	}
	if from <= to {
		result = append(result, v8ASTRange(from, to))
	}
	srcIndex++
	for ; srcIndex < len(src); srcIndex++ {
		result = append(result, src[srcIndex])
	}
	return result
}

func v8TonodeClampRangesToOneByte(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	if !v8TonodeRangesAreCanonical(ranges) {
		panic("clamping requires canonical ranges")
	}
	n := len(ranges)
	for ; n > 0; n-- {
		current := &ranges[n-1]
		if current.from <= 0xff {
			if current.to > 0xff {
				current.to = 0xff
			}
			break
		}
	}
	return ranges[:n]
}

func v8TonodeRangesEqual(lhs, rhs []v8ASTCharacterRange) bool {
	if !v8TonodeRangesAreCanonical(lhs) || !v8TonodeRangesAreCanonical(rhs) {
		panic("range equality requires canonical ranges")
	}
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

func v8TonodeToCanonicalRangeList(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	if len(ranges) == 0 {
		return nil
	}
	copyOfRanges := append([]v8ASTCharacterRange(nil), ranges...)
	return v8TonodeCanonicalizeRanges(copyOfRanges)
}

func v8TonodeUnionClassSetOperand(dst, src *v8ASTClassSetOperand) {
	dst.ranges = append(dst.ranges, src.ranges...)
	for _, candidate := range src.strings {
		found := false
		for _, existing := range dst.strings {
			if equalRunes(existing.codePoints, candidate.codePoints) {
				found = true
				break
			}
		}
		if !found {
			dst.strings = append(dst.strings, candidate)
		}
	}
	v8TonodeSortClassSetStrings(dst.strings)
}

func v8TonodeIntersectClassSetOperand(dst *v8ASTClassSetOperand, other *v8ASTClassSetOperand) {
	dst.ranges = v8TonodeIntersectRanges(dst.ranges, other.ranges)
	kept := dst.strings[:0]
	for _, current := range dst.strings {
		for _, candidate := range other.strings {
			if equalRunes(current.codePoints, candidate.codePoints) {
				kept = append(kept, current)
				break
			}
		}
	}
	dst.strings = kept
}

func v8TonodeSubtractClassSetOperand(dst *v8ASTClassSetOperand, other *v8ASTClassSetOperand) {
	dst.ranges = v8TonodeSubtractRanges(dst.ranges, other.ranges)
	kept := dst.strings[:0]
	for _, current := range dst.strings {
		remove := false
		for _, candidate := range other.strings {
			if equalRunes(current.codePoints, candidate.codePoints) {
				remove = true
				break
			}
		}
		if !remove {
			kept = append(kept, current)
		}
	}
	dst.strings = kept
}

func v8TonodeSortClassSetStrings(strings v8ASTCharacterClassStrings) {
	sort.SliceStable(strings, func(i, j int) bool {
		return (v8ASTCharacterClassStringLess{}).less(strings[i].codePoints, strings[j].codePoints)
	})
}

func equalRunes(lhs, rhs []rune) bool {
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

// v8TonodeExpansionCompiler is the small part of Compiler state borrowed by
// ExpansionLimiter. close() is the Go equivalent of its destructor.
type v8TonodeExpansionCompiler interface {
	CurrentExpansionFactor() int
	SetCurrentExpansionFactor(int)
}

type v8TonodeExpansionLimiter struct {
	compiler             v8TonodeExpansionCompiler
	savedExpansionFactor int
	okToExpand           bool
	closed               bool
}

func v8TonodeNewExpansionLimiter(compiler v8TonodeExpansionCompiler, factor int) *v8TonodeExpansionLimiter {
	if factor <= 0 {
		panic("regexp graph expansion factor must be positive")
	}
	saved := compiler.CurrentExpansionFactor()
	limiter := &v8TonodeExpansionLimiter{
		compiler:             compiler,
		savedExpansionFactor: saved,
		okToExpand:           saved <= v8TonodeMaxExpansionFactor,
	}
	if limiter.okToExpand {
		if factor > v8TonodeMaxExpansionFactor {
			limiter.okToExpand = false
			compiler.SetCurrentExpansionFactor(v8TonodeExpansionOverflow)
		} else {
			newFactor := saved * factor
			limiter.okToExpand = newFactor <= v8TonodeMaxExpansionFactor
			compiler.SetCurrentExpansionFactor(newFactor)
		}
	}
	return limiter
}

func (limiter *v8TonodeExpansionLimiter) okToExpandNow() bool { return limiter.okToExpand }

func (limiter *v8TonodeExpansionLimiter) close() {
	if limiter == nil || limiter.closed {
		return
	}
	limiter.compiler.SetCurrentExpansionFactor(limiter.savedExpansionFactor)
	limiter.closed = true
}

type v8TonodeFlagCompiler interface {
	Flags() Flags
	SetFlags(Flags)
}

type v8TonodeModifiersScope struct {
	compiler      v8TonodeFlagCompiler
	previousFlags Flags
	closed        bool
}

func v8TonodeNewModifiersScope(compiler v8TonodeFlagCompiler, flags Flags) *v8TonodeModifiersScope {
	scope := &v8TonodeModifiersScope{compiler: compiler, previousFlags: compiler.Flags()}
	compiler.SetFlags(flags)
	return scope
}

func (scope *v8TonodeModifiersScope) close() {
	if scope == nil || scope.closed {
		return
	}
	scope.compiler.SetFlags(scope.previousFlags)
	scope.closed = true
}

type v8TonodeAssertionSequenceState struct {
	seenAssertions   uint32
	sequenceStartIdx int
}

type v8TonodeCaseFoldingMode uint8

const (
	v8TonodeNonUnicodeCaseFolding v8TonodeCaseFoldingMode = iota
	v8TonodeUnicodeCaseFolding
)

func v8TonodeCaseFoldingModeForFlags(flags Flags) v8TonodeCaseFoldingMode {
	if IsEitherUnicode(flags) {
		return v8TonodeUnicodeCaseFolding
	}
	return v8TonodeNonUnicodeCaseFolding
}

func v8TonodeCanonicalLegacyRune(canonicalize func(rune) rune, value rune) rune {
	if canonicalize == nil {
		return canonicalizeLegacy(value)
	}
	return canonicalize(value)
}

func v8TonodeCompareCaseInsensitiveUnicode(_ v8TonodeCaseFoldingMode, lhs, rhs uint16) int {
	if lhs == rhs {
		return 0
	}
	left := v8ASTCaseFoldingEquivalenceKey(rune(lhs), v8TonodeUnicodeCaseFolding)
	right := v8ASTCaseFoldingEquivalenceKey(rune(rhs), v8TonodeUnicodeCaseFolding)
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

func v8TonodeCompareCaseInsensitiveLegacy(canonicalize func(rune) rune, lhs, rhs rune) int {
	if lhs == rhs {
		return 0
	}
	if lhs >= 'a' || rhs >= 'a' {
		lhs = v8TonodeCanonicalLegacyRune(canonicalize, lhs)
		rhs = v8TonodeCanonicalLegacyRune(canonicalize, rhs)
	}
	if lhs < rhs {
		return -1
	}
	if lhs > rhs {
		return 1
	}
	return 0
}

func v8TonodeCharactersEqualWithMode(ignoreCase bool, lhs, rhs uint16, mode v8TonodeCaseFoldingMode, canonicalize func(rune) rune) bool {
	if lhs == rhs {
		return true
	}
	if !ignoreCase {
		return false
	}
	if mode == v8TonodeUnicodeCaseFolding {
		return v8TonodeCompareCaseInsensitiveUnicode(mode, lhs, rhs) == 0
	}
	return v8TonodeCompareCaseInsensitiveLegacy(canonicalize, rune(lhs), rune(rhs)) == 0
}

func v8TonodeCharAtEquals(ignoreCase bool, index int, lhs, rhs *v8ASTAtom, mode v8TonodeCaseFoldingMode, canonicalize func(rune) rune) bool {
	if lhs == nil || rhs == nil || index < 0 || index >= len(lhs.data) || index >= len(rhs.data) {
		return false
	}
	return v8TonodeCharactersEqualWithMode(ignoreCase, lhs.data[index], rhs.data[index], mode, canonicalize)
}

// These interfaces enumerate the source methods in regexp-compiler-tonode.cc.
// v8TonodeGraphLowerer below supplies the AST-to-graph behavior using the
// corresponding Go node and compiler state records.
type v8TonodeLoweringMethods interface {
	treeToNode(v8TonodeTree, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	atomToNodeImpl(*v8ASTAtom, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	textToNodeImpl(*v8ASTText, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	classRangesIsStandard(*v8ASTClassRanges) bool
	classRangesToNodeImpl(*v8ASTClassRanges, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	classSetOperandToNodeImpl(*v8ASTClassSetOperand, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	classSetExpressionToNodeImpl(*v8ASTClassSetExpression, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	disjunctionSortConsecutiveAtoms(*v8ASTDisjunction, v8TonodeCompiler) bool
	disjunctionRationalizeConsecutiveAtoms(*v8ASTDisjunction, v8TonodeCompiler)
	disjunctionFixSingleCharacterDisjunctions(*v8ASTDisjunction, v8TonodeCompiler)
	disjunctionToNodeImpl(*v8ASTDisjunction, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	quantifierToNode(min, max int, greedy bool, body v8TonodeTree, compiler v8TonodeCompiler, success v8TonodeNode, notAtStart bool) v8TonodeNode
	quantifierToNodeImpl(*v8ASTQuantifier, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	assertionToNodeImpl(*v8ASTAssertion, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	backReferenceToNodeImpl(*v8ASTBackReference, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	emptyToNodeImpl(*v8ASTEmpty, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	groupToNodeImpl(*v8ASTGroup, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	lookaroundBuilderForMatch(*v8ASTLookaroundBuilder, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	lookaroundToNodeImpl(*v8ASTLookaround, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	captureToNodeImpl(*v8ASTCapture, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	captureToNode(v8TonodeTree, v8TonodeCompiler, v8TonodeNode, *v8ASTCapture) v8TonodeNode
	alternativeToNodeImpl(*v8ASTAlternative, v8TonodeCompiler, v8TonodeNode) v8TonodeNode
}

type v8TonodeHelperMethods interface {
	compareInverseRanges([]v8ASTCharacterRange, []int) bool
	compareRanges([]v8ASTCharacterRange, []int) bool
	toCanonicalRangeList([]v8ASTCharacterRange) []v8ASTCharacterRange
	moveRanges(ranges []v8ASTCharacterRange, from, to, count int)
	insertRangeInCanonicalList(ranges []v8ASTCharacterRange, count int, insert v8ASTCharacterRange) ([]v8ASTCharacterRange, int)
	characterSetRanges(*v8ASTCharacterSet) []v8ASTCharacterRange
	characterSetCanonicalize(*v8ASTCharacterSet)
	addBMPCharacters(v8TonodeCompiler, v8TonodeNode, v8TonodeNode, *v8TonodeUnicodeRangeSplitter)
	addNonBMPSurrogatePairs(v8TonodeCompiler, v8TonodeNode, v8TonodeNode, *v8TonodeUnicodeRangeSplitter)
	negativeLookaroundAgainstReadDirectionAndMatch(v8TonodeCompiler, []v8ASTCharacterRange, []v8ASTCharacterRange, v8TonodeNode, bool) v8TonodeNode
	matchAndNegativeLookaroundInReadDirection(v8TonodeCompiler, []v8ASTCharacterRange, []v8ASTCharacterRange, v8TonodeNode, bool) v8TonodeNode
	addLoneLeadSurrogates(v8TonodeCompiler, v8TonodeNode, v8TonodeNode, *v8TonodeUnicodeRangeSplitter)
	addLoneTrailSurrogates(v8TonodeCompiler, v8TonodeNode, v8TonodeNode, *v8TonodeUnicodeRangeSplitter)
	unanchoredAdvance(v8TonodeCompiler, v8TonodeNode) v8TonodeNode
	startsWithAtom(v8TonodeTree) bool
	firstAtom(v8TonodeTree) *v8ASTAtom
	compareFirstChar(v8TonodeTree, v8TonodeTree) int
	caseFoldingMode(Flags) v8TonodeCaseFoldingMode
	boundaryAssertionAsLookaround(v8TonodeCompiler, v8TonodeNode, v8ASTAssertionType) v8TonodeNode
	addClass([]int, []v8ASTCharacterRange) []v8ASTCharacterRange
	addClassNegated([]int, []v8ASTCharacterRange) []v8ASTCharacterRange
	canonicalLegacyRune(func(rune) rune, rune) rune
	compareCaseInsensitiveUnicode(v8TonodeCaseFoldingMode, uint16, uint16) int
	compareCaseInsensitiveLegacy(func(rune) rune, rune, rune) int
	compareFirstCharCaseInsensitive(v8TonodeTree, v8TonodeTree, v8TonodeCaseFoldingMode, func(rune) rune) int
	charactersEqual(bool, uint16, uint16, v8TonodeCaseFoldingMode, func(rune) rune) bool
	charAtEquals(bool, int, *v8ASTAtom, *v8ASTAtom, v8TonodeCaseFoldingMode, func(rune) rune) bool
}

type v8TonodeCharacterRangeMethods interface {
	addUnicodeCaseEquivalents(*[]v8ASTCharacterRange)
	addClassEscape(v8ASTStandardCharacterSet, *[]v8ASTCharacterRange, bool)
	addCaseEquivalents(isolate any, ranges *[]v8ASTCharacterRange, oneByte bool)
	isCanonical([]v8ASTCharacterRange) bool
	canonicalize([]v8ASTCharacterRange) []v8ASTCharacterRange
	negate([]v8ASTCharacterRange) []v8ASTCharacterRange
	intersect([]v8ASTCharacterRange, []v8ASTCharacterRange) []v8ASTCharacterRange
	intersects([]v8ASTCharacterRange, []v8ASTCharacterRange) bool
	safeAdvanceRange([]v8ASTCharacterRange, int) (int, rune, rune)
	subtract([]v8ASTCharacterRange, []v8ASTCharacterRange) []v8ASTCharacterRange
	clampToOneByte([]v8ASTCharacterRange) []v8ASTCharacterRange
	equals([]v8ASTCharacterRange, []v8ASTCharacterRange) bool
}

type v8TonodeClassSetMethods interface {
	classSetOperandUnion(*v8ASTClassSetOperand, *v8ASTClassSetOperand)
	classSetOperandIntersect(*v8ASTClassSetOperand, *v8ASTClassSetOperand)
	classSetOperandSubtract(*v8ASTClassSetOperand, *v8ASTClassSetOperand)
	classSetExpressionCompute(any) *v8ASTClassSetOperand
}

type v8TonodeAstRewriteMethods interface {
	assertionSequenceMaybeRewrite([]v8TonodeTree, v8TonodeCompiler) []v8TonodeTree
	appendFlattenedAndFold([]v8TonodeTree, v8TonodeTree, Flags, *v8TonodeAssertionSequenceState)
}

type v8TonodeDiagnostics interface {
	trace(message string)
	traceWithNode(message string, tree v8TonodeTree)
	registerNode(node v8TonodeNode)
}

// v8TonodeGraphLowerer is the concrete AST-to-node pass from
// regexp-compiler-tonode.cc. The source AST declaration interfaces remain
// useful as a reference surface, while this implementation dispatches over
// their field-for-field Go state records. It builds the same ordered graph
// nodes that the compiler core consumes.
type v8TonodeGraphLowerer struct {
	stack v8CompilerAnalysisStackChecker
}

func v8TonodeNewGraphLowerer(stack v8CompilerAnalysisStackChecker) *v8TonodeGraphLowerer {
	return &v8TonodeGraphLowerer{stack: stack}
}

func (lowerer *v8TonodeGraphLowerer) backtrack(compiler *v8CompilerCompiler) v8NodesNode {
	return v8NodesNewEndNode(v8NodesBacktrack, compiler.Flags(), compiler.Zone())
}

func (lowerer *v8TonodeGraphLowerer) track(node v8NodesNode) v8NodesNode {
	return node
}

// treeToNode performs the source's early pruning, periodic stack check, and
// regexp-too-large handling before dispatching to a concrete AST node.
func (lowerer *v8TonodeGraphLowerer) treeToNode(tree any, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if compiler == nil || success == nil {
		panic("V8 Tree::ToNode requires a compiler and successor")
	}
	if success.IsBacktrack() {
		return success
	}
	compiler.ToNodeMaybeCheckForStackOverflow(lowerer.stack)
	if compiler.IsRegExpTooBig() {
		return lowerer.backtrack(compiler)
	}
	return lowerer.treeToNodeImpl(tree, compiler, success)
}

// treeToNodeImpl is the virtual ToNodeImpl dispatch generated by the AST
// subclasses' DECL_BOILERPLATE declarations.
func (lowerer *v8TonodeGraphLowerer) treeToNodeImpl(tree any, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if compiler == nil || success == nil {
		panic("V8 Tree::ToNodeImpl requires a compiler and successor")
	}
	switch value := tree.(type) {
	case *v8ASTAtom:
		return lowerer.atomToNodeImpl(value, compiler, success)
	case *v8ASTText:
		return lowerer.textToNodeImpl(value, compiler, success)
	case *v8ASTClassRanges:
		return lowerer.classRangesToNodeImpl(value, compiler, success)
	case *v8ASTClassSetOperand:
		return lowerer.classSetOperandToNodeImpl(value, compiler, success)
	case *v8ASTClassSetExpression:
		return lowerer.classSetExpressionToNodeImpl(value, compiler, success)
	case *v8ASTDisjunction:
		return lowerer.disjunctionToNodeImpl(value, compiler, success)
	case *v8ASTAlternative:
		return lowerer.alternativeToNodeImpl(value, compiler, success)
	case *v8ASTQuantifier:
		return lowerer.quantifierToNodeImpl(value, compiler, success)
	case *v8ASTAssertion:
		return lowerer.assertionToNodeImpl(value, compiler, success)
	case *v8ASTBackReference:
		return lowerer.backReferenceToNodeImpl(value, compiler, success)
	case *v8ASTEmpty:
		return lowerer.emptyToNodeImpl(value, compiler, success)
	case *v8ASTGroup:
		return lowerer.groupToNodeImpl(value, compiler, success)
	case *v8ASTLookaround:
		return lowerer.lookaroundToNodeImpl(value, compiler, success)
	case *v8ASTCapture:
		return lowerer.captureToNodeImpl(value, compiler, success)
	default:
		panic("unknown V8 regular-expression AST node in graph lowering")
	}
}

func (lowerer *v8TonodeGraphLowerer) atomToNodeImpl(atom *v8ASTAtom, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if atom == nil {
		panic("V8 Atom::ToNodeImpl requires an atom")
	}
	if len(atom.data) == 0 {
		return success
	}
	elements := &v8NodesZoneList[v8NodesTextElement]{}
	elements.Add(v8CompilerTextElementFromAtom(atom))
	node := (v8CompilerTextNodeFactory{}).NewTextNode(elements, compiler.ReadBackward(), success, compiler.Flags())
	if compiler.OneByte() && !node.CanMatchLatin1(compiler) {
		return lowerer.backtrack(compiler)
	}
	return lowerer.track(node)
}

func (lowerer *v8TonodeGraphLowerer) textToNodeImpl(text *v8ASTText, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if text == nil {
		panic("V8 Text::ToNodeImpl requires text")
	}
	if len(text.elements) == 0 {
		return success
	}
	elements := &v8NodesZoneList[v8NodesTextElement]{items: append([]v8NodesTextElement(nil), text.elements...)}
	node := (v8CompilerTextNodeFactory{}).NewTextNode(elements, compiler.ReadBackward(), success, compiler.Flags())
	if compiler.OneByte() && !node.CanMatchLatin1(compiler) {
		return lowerer.backtrack(compiler)
	}
	return lowerer.track(node)
}

func v8TonodeClassFromRanges(ranges []v8ASTCharacterRange, flags v8ASTClassRangesFlags) *v8ASTClassRanges {
	return &v8ASTClassRanges{set: v8ASTCharacterSet{ranges: append([]v8ASTCharacterRange(nil), ranges...)}, flags: flags}
}

func (lowerer *v8TonodeGraphLowerer) textForRanges(ranges []v8ASTCharacterRange, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if len(ranges) == 0 {
		return lowerer.backtrack(compiler)
	}
	list := &v8NodesZoneList[v8NodesCharacterRange]{items: append([]v8NodesCharacterRange(nil), ranges...)}
	return (v8CompilerTextNodeFactory{}).CreateForCharacterRanges(compiler.Zone(), list, compiler.ReadBackward(), success, compiler.Flags())
}

func (lowerer *v8TonodeGraphLowerer) textForSurrogatePair(lead v8ASTCharacterRange, trails []v8ASTCharacterRange, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	trailList := &v8NodesZoneList[v8NodesCharacterRange]{items: append([]v8NodesCharacterRange(nil), trails...)}
	return (v8CompilerTextNodeFactory{}).CreateForSurrogatePair(compiler.Zone(), lead, trailList, compiler.ReadBackward(), success, compiler.Flags())
}

func (lowerer *v8TonodeGraphLowerer) textForSurrogatePairRanges(leads []v8ASTCharacterRange, trail v8ASTCharacterRange, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	leadList := &v8NodesZoneList[v8NodesCharacterRange]{items: append([]v8NodesCharacterRange(nil), leads...)}
	return (v8CompilerTextNodeFactory{}).CreateForSurrogatePairRanges(compiler.Zone(), leadList, trail, compiler.ReadBackward(), success, compiler.Flags())
}

func (lowerer *v8TonodeGraphLowerer) addBMPCharacters(result *v8NodesChoiceNode, success v8NodesNode, splitter *v8TonodeUnicodeRangeSplitter, compiler *v8CompilerCompiler) {
	ranges := v8TonodeToCanonicalRangeList(splitter.bmpRanges())
	if len(ranges) == 0 {
		return
	}
	result.AddAlternative(v8NodesNewGuardedAlternative(lowerer.textForRanges(ranges, compiler, success)))
}

func (lowerer *v8TonodeGraphLowerer) addNonBMPSurrogatePairs(result *v8NodesChoiceNode, success v8NodesNode, splitter *v8TonodeUnicodeRangeSplitter, compiler *v8CompilerCompiler) {
	if compiler.OneByte() {
		panic("V8 non-BMP surrogate lowering is unavailable in one-byte mode")
	}
	plan := v8TonodeGroupNonBMPRanges(splitter.nonBMPRanges())
	if len(plan.leadingWithFullTrailing) != 0 {
		result.AddAlternative(v8NodesNewGuardedAlternative(lowerer.textForSurrogatePairRanges(
			plan.leadingWithFullTrailing,
			v8ASTRange(v8TonodeTrailSurrogateStart, v8TonodeTrailSurrogateEnd), compiler, success)))
	}
	for packed, trailing := range plan.groupedByLeading {
		lead := v8ASTRange(rune(v8TonodeUC16RangeFrom(packed)), rune(v8TonodeUC16RangeTo(packed)))
		result.AddAlternative(v8NodesNewGuardedAlternative(lowerer.textForSurrogatePair(lead, trailing, compiler, success)))
	}
}

func (lowerer *v8TonodeGraphLowerer) negativeLookaroundAgainstReadDirectionAndMatch(lookbehind, match []v8ASTCharacterRange, success v8NodesNode, readBackward bool, compiler *v8CompilerCompiler) v8NodesNode {
	oldDirection := compiler.ReadBackward()
	compiler.SetReadBackward(readBackward)
	matchNode := lowerer.textForRanges(match, compiler, success)
	stackReg, positionReg := compiler.UnicodeLookaroundStackRegister(), compiler.UnicodeLookaroundPositionRegister()
	builder := lowerer.newLookaroundBuilder(false, matchNode, compiler, stackReg, positionReg, 0, -1)
	compiler.SetReadBackward(!readBackward)
	lookbehindNode := lowerer.textForRanges(lookbehind, compiler, builder.onMatchSuccess.(v8NodesNode))
	compiler.SetReadBackward(oldDirection)
	return builder.forMatch(compiler, lookbehindNode).(v8NodesNode)
}

func (lowerer *v8TonodeGraphLowerer) matchAndNegativeLookaroundInReadDirection(match, lookahead []v8ASTCharacterRange, success v8NodesNode, readBackward bool, compiler *v8CompilerCompiler) v8NodesNode {
	stackReg, positionReg := compiler.UnicodeLookaroundStackRegister(), compiler.UnicodeLookaroundPositionRegister()
	builder := lowerer.newLookaroundBuilder(false, success, compiler, stackReg, positionReg, 0, -1)
	oldDirection := compiler.ReadBackward()
	compiler.SetReadBackward(readBackward)
	negativeMatch := lowerer.textForRanges(lookahead, compiler, builder.onMatchSuccess.(v8NodesNode))
	matchNode := lowerer.textForRanges(match, compiler, builder.forMatch(compiler, negativeMatch).(v8NodesNode))
	compiler.SetReadBackward(oldDirection)
	return matchNode
}

func (lowerer *v8TonodeGraphLowerer) addLoneLeadSurrogates(result *v8NodesChoiceNode, success v8NodesNode, splitter *v8TonodeUnicodeRangeSplitter, compiler *v8CompilerCompiler) {
	leads := v8TonodeToCanonicalRangeList(splitter.leadSurrogateRanges())
	if len(leads) == 0 {
		return
	}
	trails := v8ASTRangeList(v8ASTRange(v8TonodeTrailSurrogateStart, v8TonodeTrailSurrogateEnd))
	var match v8NodesNode
	if compiler.ReadBackward() {
		match = lowerer.negativeLookaroundAgainstReadDirectionAndMatch(trails, leads, success, true, compiler)
	} else {
		match = lowerer.matchAndNegativeLookaroundInReadDirection(leads, trails, success, false, compiler)
	}
	result.AddAlternative(v8NodesNewGuardedAlternative(match))
}

func (lowerer *v8TonodeGraphLowerer) addLoneTrailSurrogates(result *v8NodesChoiceNode, success v8NodesNode, splitter *v8TonodeUnicodeRangeSplitter, compiler *v8CompilerCompiler) {
	trails := v8TonodeToCanonicalRangeList(splitter.trailSurrogateRanges())
	if len(trails) == 0 {
		return
	}
	leads := v8ASTRangeList(v8ASTRange(v8TonodeLeadSurrogateStart, v8TonodeLeadSurrogateEnd))
	var match v8NodesNode
	if compiler.ReadBackward() {
		match = lowerer.matchAndNegativeLookaroundInReadDirection(trails, leads, success, true, compiler)
	} else {
		match = lowerer.negativeLookaroundAgainstReadDirectionAndMatch(leads, trails, success, false, compiler)
	}
	result.AddAlternative(v8NodesNewGuardedAlternative(match))
}

func (lowerer *v8TonodeGraphLowerer) unanchoredAdvance(compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if compiler.ReadBackward() {
		panic("V8 unanchored advance cannot run while reading backward")
	}
	return v8NodesNewUnanchoredAdvanceNode(success, compiler.Flags())
}

func v8TonodeAddUnicodeCaseEquivalents(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	ranges = v8TonodeCanonicalizeRanges(ranges)
	if len(ranges) == 1 && ranges[0].isEverything(v8TonodeMaxCodePoint) {
		return ranges
	}
	return v8ASTCaseFoldingCloseOver(ranges, v8TonodeUnicodeCaseFolding)
}

func (lowerer *v8TonodeGraphLowerer) classRangesToNodeImpl(class *v8ASTClassRanges, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if class == nil {
		panic("V8 ClassRanges::ToNodeImpl requires class ranges")
	}
	v8ASTCharacterSetCanonicalize(&class.set)
	ranges := append([]v8ASTCharacterRange(nil), v8TonodeCharacterSetRanges(&class.set)...)
	if v8CompilerNeedsUnicodeCaseEquivalents(compiler.Flags()) && class.flags&v8ASTClassNoCaseFoldingNeeded == 0 {
		// Case-insensitive Unicode semantics close \w before complementing it.
		// Closing the already-negated \W range incorrectly adds fold-equivalent
		// non-word characters such as long s back into the set.
		if standard, ok := class.set.standardSetTypeValue(); ok && standard == v8ASTNotWord {
			ranges = v8TonodeAddClassEscape(standard, nil, true)
		} else {
			ranges = v8TonodeAddUnicodeCaseEquivalents(ranges)
		}
		class.set.ranges = ranges
	}

	if !IsEitherUnicode(compiler.Flags()) || compiler.OneByte() || class.flags&v8ASTClassContainsSplitSurrogate != 0 {
		node := (v8CompilerTextNodeFactory{}).NewTextNodeForClassRanges(class, compiler.ReadBackward(), success, compiler.Flags())
		if compiler.OneByte() && !node.CanMatchLatin1(compiler) {
			return lowerer.backtrack(compiler)
		}
		return lowerer.track(node)
	}

	if class.flags&v8ASTClassNegated != 0 {
		ranges = v8TonodeNegateRanges(v8TonodeCanonicalizeRanges(ranges))
	}
	if len(ranges) == 0 {
		return lowerer.backtrack(compiler)
	}
	if class.set.standardSet != nil && *class.set.standardSet == v8ASTEverything {
		return lowerer.unanchoredAdvance(compiler, success)
	}
	choice := v8NodesNewChoiceNode(2, compiler.Flags(), compiler.Zone())
	splitter := v8TonodeNewUnicodeRangeSplitter(ranges)
	lowerer.addBMPCharacters(choice, success, splitter, compiler)
	lowerer.addNonBMPSurrogatePairs(choice, success, splitter, compiler)
	lowerer.addLoneLeadSurrogates(choice, success, splitter, compiler)
	lowerer.addLoneTrailSurrogates(choice, success, splitter, compiler)
	if choice.Alternatives().Len() == 0 {
		return lowerer.backtrack(compiler)
	}
	if choice.Alternatives().Len() == 1 {
		return choice.Alternatives().At(0).Node()
	}
	return lowerer.track(choice)
}

func (lowerer *v8TonodeGraphLowerer) classSetOperandToNodeImpl(operand *v8ASTClassSetOperand, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if operand == nil {
		panic("V8 ClassSetOperand::ToNodeImpl requires an operand")
	}
	if compiler.Flags().UnicodeSets && !compiler.OneByte() && !compiler.ReadBackward() && v8TonodeShouldUseFiniteStringSet(operand.strings) {
		set := v8TonodeBuildFiniteStringSet(operand.strings, operand.ranges)
		node := v8NodesNewTextNode(&v8NodesZoneList[v8NodesTextElement]{}, false, success, compiler.Flags())
		node.finiteStringSet = set
		return lowerer.track(node)
	}
	size := len(operand.strings)
	if len(operand.ranges) != 0 {
		size++
	}
	if size == 0 {
		return lowerer.backtrack(compiler)
	}
	v8TonodeSortClassSetStrings(operand.strings)
	alternatives := make([]any, 0, size)
	var emptyString any
	for _, item := range operand.strings {
		if _, empty := item.tree.(*v8ASTEmpty); empty {
			emptyString = item.tree
		} else {
			alternatives = append(alternatives, item.tree)
		}
	}
	if len(operand.ranges) != 0 {
		alternatives = append(alternatives, v8TonodeClassFromRanges(operand.ranges, v8ASTClassNoCaseFoldingNeeded))
	}
	if emptyString != nil {
		alternatives = append(alternatives, emptyString)
	}
	var tree any
	if len(alternatives) == 1 {
		tree = alternatives[0]
	} else {
		tree = &v8ASTDisjunction{alternatives: alternatives}
	}
	return lowerer.treeToNode(tree, compiler, success)
}

func (lowerer *v8TonodeGraphLowerer) classSetExpressionCompute(root any) *v8ASTClassSetOperand {
	switch value := root.(type) {
	case *v8ASTClassSetOperand:
		return value
	case *v8ASTClassSetExpression:
		if value == nil {
			panic("nil V8 class-set expression")
		}
		if len(value.operands) == 0 {
			if value.isNegated {
				return &v8ASTClassSetOperand{ranges: v8ASTRangeList(v8ASTRangeEverything())}
			}
			return &v8ASTClassSetOperand{}
		}
		result := lowerer.classSetExpressionCompute(value.operands[0])
		for _, operandTree := range value.operands[1:] {
			other := lowerer.classSetExpressionCompute(operandTree)
			switch value.operation {
			case v8ASTClassSetUnion:
				v8TonodeUnionClassSetOperand(result, other)
			case v8ASTClassSetIntersection:
				v8TonodeIntersectClassSetOperand(result, other)
			case v8ASTClassSetSubtraction:
				v8TonodeSubtractClassSetOperand(result, other)
			default:
				panic("unknown V8 class-set operation")
			}
		}
		result.ranges = v8TonodeCanonicalizeRanges(result.ranges)
		if value.isNegated {
			if len(result.strings) != 0 {
				panic("V8 cannot complement a class set containing strings")
			}
			result.ranges = v8TonodeNegateRanges(v8TonodeCanonicalizeRanges(result.ranges))
			value.isNegated = false
		}
		value.operands = []any{result}
		return result
	default:
		panic("unknown V8 class-set expression operand")
	}
}

func (lowerer *v8TonodeGraphLowerer) classSetExpressionToNodeImpl(expression *v8ASTClassSetExpression, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	root := lowerer.classSetExpressionCompute(expression)
	return lowerer.classSetOperandToNodeImpl(root, compiler, success)
}

func v8TonodeTreeMinMatch(tree any) int {
	switch value := tree.(type) {
	case *v8ASTAtom:
		return len(value.data)
	case *v8ASTText:
		return value.length
	case *v8ASTClassRanges:
		return value.minMatch()
	case *v8ASTClassSetOperand:
		return value.minMatch
	case *v8ASTClassSetExpression:
		// A UnicodeSets class expression can match a character or a non-empty
		// string, but never the empty string. The expression may describe an
		// empty language (for example, an intersection of disjoint sets); the
		// lower bound still holds vacuously because it then has no successful
		// match. Returning zero here makes quantified class expressions enter
		// the empty-body loop path, which restores a position register that was
		// never saved. Keep this fact explicit for quantifier lowering.
		return 1
	case *v8ASTDisjunction:
		return value.minMatch
	case *v8ASTAlternative:
		return value.minMatch
	case *v8ASTQuantifier:
		return value.minMatch
	case *v8ASTCapture:
		return value.minMatch
	case *v8ASTGroup:
		return value.minMatch
	case *v8ASTAssertion, *v8ASTBackReference, *v8ASTEmpty, *v8ASTLookaround:
		return 0
	default:
		panic("unknown V8 AST minimum-match query")
	}
}

func v8TonodeTreeMaxMatch(tree any) int {
	switch value := tree.(type) {
	case *v8ASTAtom:
		return len(value.data)
	case *v8ASTText:
		return value.length
	case *v8ASTClassRanges:
		return value.maxMatch()
	case *v8ASTClassSetOperand:
		return value.maxMatch
	case *v8ASTClassSetExpression:
		return value.maxMatch
	case *v8ASTDisjunction:
		return value.maxMatch
	case *v8ASTAlternative:
		return value.maxMatch
	case *v8ASTQuantifier:
		return value.maxMatch
	case *v8ASTCapture:
		return value.maxMatch
	case *v8ASTGroup:
		return value.maxMatch
	case *v8ASTAssertion, *v8ASTEmpty, *v8ASTLookaround:
		return 0
	case *v8ASTBackReference:
		return v8ASTInfinity
	default:
		panic("unknown V8 AST maximum-match query")
	}
}

func v8TonodeCaptureRegisters(tree any, budget int) v8ASTInterval {
	return v8ASTCaptureRegisters(tree, v8ASTNewStackLimiter(budget))
}

func (lowerer *v8TonodeGraphLowerer) quantifierToNodeImpl(quantifier *v8ASTQuantifier, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if quantifier == nil {
		panic("V8 Quantifier::ToNodeImpl requires a quantifier")
	}
	return lowerer.quantifierToNode(quantifier.min, quantifier.max,
		quantifier.quantifier == v8ASTGreedy,
		quantifier.body, compiler, success, false)
}

func (lowerer *v8TonodeGraphLowerer) quantifierToNode(minimum, maximum int, greedy bool, body any, compiler *v8CompilerCompiler, success v8NodesNode, notAtStart bool) v8NodesNode {
	if maximum == 0 {
		return success
	}
	bodyMin := v8TonodeTreeMinMatch(body)
	bodyCanBeEmpty := bodyMin == 0
	captureRegisters := v8TonodeCaptureRegisters(body, v8NodesRecursionBudget)
	if !v8ASTIntervalIsValid(captureRegisters) {
		compiler.SetRegExpTooBig()
		return lowerer.backtrack(compiler)
	}
	needsCaptureClearing := !v8ASTIntervalIsEmpty(captureRegisters) && maximum != 1
	bodyStartRegister := v8CompilerNoRegister
	if bodyCanBeEmpty {
		// The source reserves this register before constructing the loop.
		bodyStartRegister = compiler.AllocateRegister()
	}
	wantUnroll := compiler.Optimize()
	if !bodyCanBeEmpty {
		if wantUnroll && !needsCaptureClearing {
			factor := minimum
			if maximum != minimum {
				factor++
			}
			limiter := v8TonodeNewExpansionLimiter(compiler, factor)
			if minimum > 0 && minimum <= v8TonodeMaxUnrolledMinMatch && limiter.okToExpandNow() {
				newMaximum := maximum
				if maximum != v8ASTInfinity {
					newMaximum -= minimum
				}
				answer := lowerer.quantifierToNode(0, newMaximum, greedy, body, compiler, success, true)
				for index := 0; index < minimum; index++ {
					answer = lowerer.treeToNode(body, compiler, answer)
				}
				limiter.close()
				return answer
			}
			limiter.close()
			if maximum <= v8TonodeMaxUnrolledMaxMatch && minimum == 0 {
				limiter = v8TonodeNewExpansionLimiter(compiler, maximum)
				if limiter.okToExpandNow() {
					answer := success
					for index := 0; index < maximum; index++ {
						choice := v8NodesNewChoiceNode(2, compiler.Flags(), compiler.Zone())
						bodyNode := lowerer.treeToNode(body, compiler, answer)
						if greedy {
							choice.AddAlternative(v8NodesNewGuardedAlternative(bodyNode))
							choice.AddAlternative(v8NodesNewGuardedAlternative(success))
						} else {
							choice.AddAlternative(v8NodesNewGuardedAlternative(success))
							choice.AddAlternative(v8NodesNewGuardedAlternative(bodyNode))
						}
						if notAtStart && !compiler.ReadBackward() {
							choice.SetNotAtStart()
						}
						answer = choice
					}
					limiter.close()
					return answer
				}
				limiter.close()
			}
		}
	}

	hasMinimum, hasMaximum := minimum > 0, maximum < v8ASTInfinity
	needsCounter := hasMinimum || hasMaximum
	counterRegister := v8CompilerNoRegister
	if needsCounter {
		counterRegister = compiler.AllocateRegister()
	}
	center := v8NodesNewLoopChoiceNode(bodyCanBeEmpty, compiler.ReadBackward(), compiler.Flags(), compiler.Zone())
	// Unicode property escapes can lower to a small choice graph for UTF-16
	// surrogate handling even though the source atom is one deterministic
	// code-point class. Preserve that fact so the Go backend can keep only the
	// current greedy continuation for an end-anchored repetition.
	switch class := body.(type) {
	case *v8ASTClassRanges:
		center.SetSingleCodePointClassBody(true)
	case *v8ASTClassSetExpression:
		center.SetSingleCodePointClassBody(class != nil && !class.mayContainStrings)
	}
	if notAtStart && !compiler.ReadBackward() {
		center.SetNotAtStart()
	}
	loopReturn := v8NodesNode(center)
	if needsCounter {
		loopReturn = v8CompilerIncrementRegister(counterRegister, loopReturn, compiler.Flags())
	}
	if bodyCanBeEmpty {
		loopReturn = v8CompilerEmptyMatchCheck(bodyStartRegister, counterRegister, minimum, loopReturn, compiler.Flags())
	}
	bodyNode := lowerer.treeToNode(body, compiler, loopReturn)
	if bodyNode.IsBacktrack() {
		if hasMinimum {
			return bodyNode
		}
		return success
	}
	if bodyCanBeEmpty {
		bodyNode = v8CompilerRestorePosition(bodyStartRegister, bodyNode, compiler.Flags())
		bodyNode = v8CompilerStorePosition(bodyStartRegister, bodyNode, compiler.Flags(), nil)
	}
	if needsCaptureClearing {
		bodyNode = v8CompilerClearCaptures(captureRegisters, bodyNode, compiler.Flags())
	}
	bodyAlternative := v8NodesNewGuardedAlternative(bodyNode)
	if hasMaximum {
		bodyAlternative.AddGuard(v8NodesNewGuard(counterRegister, v8NodesGuardLT, maximum), compiler.Zone())
	}
	restAlternative := v8NodesNewGuardedAlternative(success)
	if hasMinimum {
		restAlternative.AddGuard(v8NodesNewGuard(counterRegister, v8NodesGuardGEQ, minimum), compiler.Zone())
	}
	if greedy {
		center.AddLoopAlternative(bodyAlternative)
		center.AddContinueAlternative(restAlternative)
	} else {
		center.AddContinueAlternative(restAlternative)
		center.AddLoopAlternative(bodyAlternative)
	}
	var result v8NodesNode = center
	if minimum > 0 && bodyMin > 0 && !compiler.ReadBackward() {
		eats := minimum * bodyMin
		if minimum > 256 {
			eats = 256 * min(bodyMin, 256)
		} else if bodyMin > 256 {
			eats = minimum * 256
		}
		if eats > 255 {
			eats = 255
		}
		result = v8CompilerEatsAtLeast(eats, result, compiler.Flags())
	}
	if needsCounter {
		result = v8CompilerSetRegisterForLoop(counterRegister, 0, result, compiler.Flags())
	}
	return result
}

func (lowerer *v8TonodeGraphLowerer) newLookaroundBuilder(positive bool, success v8NodesNode, compiler *v8CompilerCompiler, stackRegister, positionRegister, captureRegisterCount, captureRegisterStart int) *v8ASTLookaroundBuilder {
	var matchSuccess v8NodesNode
	if positive {
		matchSuccess = v8CompilerPositiveSubmatchSuccess(stackRegister, positionRegister,
			captureRegisterCount, captureRegisterStart, success, compiler.Flags())
	} else {
		matchSuccess = v8NodesNewNegativeSubmatchSuccess(stackRegister, positionRegister,
			captureRegisterCount, captureRegisterStart, compiler.Flags(), compiler.Zone())
	}
	return &v8ASTLookaroundBuilder{
		isPositive:           positive,
		onMatchSuccess:       matchSuccess,
		onSuccess:            success,
		stackPointerRegister: stackRegister,
		positionRegister:     positionRegister,
	}
}

func (builder *v8ASTLookaroundBuilder) onMatchSuccessNode() v8NodesNode {
	node, ok := builder.onMatchSuccess.(v8NodesNode)
	if !ok || node == nil {
		panic("V8 lookaround builder has no match-success node")
	}
	return node
}

func (lowerer *v8TonodeGraphLowerer) lookaroundBuilderForMatch(builder *v8ASTLookaroundBuilder, compiler *v8CompilerCompiler, match v8NodesNode) v8NodesNode {
	if builder.isPositive {
		success, ok := builder.onMatchSuccess.(*v8NodesActionNode)
		if !ok {
			panic("positive V8 lookaround success is not an action node")
		}
		return v8CompilerBeginPositiveSubmatch(builder.stackPointerRegister,
			builder.positionRegister, match, success, compiler.Flags())
	}
	choice := v8NodesNewNegativeLookaroundChoiceNode(
		v8NodesNewGuardedAlternative(match),
		v8NodesNewGuardedAlternative(builder.onSuccess.(v8NodesNode)),
		compiler.Flags(), compiler.Zone())
	return v8CompilerBeginNegativeSubmatch(builder.stackPointerRegister,
		builder.positionRegister, choice, compiler.Flags())
}

func (lowerer *v8TonodeGraphLowerer) lookaroundToNodeImpl(lookaround *v8ASTLookaround, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if lookaround == nil {
		panic("V8 Lookaround::ToNodeImpl requires a lookaround")
	}
	stackRegister, positionRegister := compiler.AllocateRegister(), compiler.AllocateRegister()
	registerCount := lookaround.captureCount * 2
	registerStart := 2 + lookaround.captureFrom*2
	wasReadingBackward := compiler.ReadBackward()
	compiler.SetReadBackward(lookaround.typeOf == v8ASTLookbehind)
	defer compiler.SetReadBackward(wasReadingBackward)
	builder := lowerer.newLookaroundBuilder(lookaround.isPositive, success, compiler,
		stackRegister, positionRegister, registerCount, registerStart)
	match := lowerer.treeToNode(lookaround.body, compiler, builder.onMatchSuccessNode())
	if match.IsBacktrack() && (lookaround.isPositive || compiler.IsRegExpTooBig()) {
		return match
	}
	return lowerer.lookaroundBuilderForMatch(builder, compiler, match)
}

func (lowerer *v8TonodeGraphLowerer) captureToNodeImpl(capture *v8ASTCapture, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if capture == nil {
		panic("V8 Capture::ToNodeImpl requires a capture")
	}
	return lowerer.captureToNode(capture.body, compiler, success, capture)
}

func (lowerer *v8TonodeGraphLowerer) captureToNode(body any, compiler *v8CompilerCompiler, success v8NodesNode, capture *v8ASTCapture) v8NodesNode {
	if body == nil {
		panic("V8 Capture::ToNode requires a body")
	}
	index := 0
	if capture != nil {
		index = capture.index
	}
	startRegister, endRegister := v8ASTCaptureStartRegister(index), v8ASTCaptureEndRegister(index)
	if compiler.ReadBackward() {
		startRegister, endRegister = endRegister, startRegister
	}
	storeEnd := v8CompilerStorePosition(endRegister, success, compiler.Flags(), capture)
	bodyNode := lowerer.treeToNode(body, compiler, storeEnd)
	if bodyNode.IsBacktrack() {
		return bodyNode
	}
	return v8CompilerStorePosition(startRegister, bodyNode, compiler.Flags(), capture)
}

func (lowerer *v8TonodeGraphLowerer) backReferenceToNodeImpl(backReference *v8ASTBackReference, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if backReference == nil {
		panic("V8 BackReference::ToNodeImpl requires a backreference")
	}
	result := success
	for _, capture := range backReference.captures {
		if capture == nil {
			continue
		}
		capture.isBackreferenced = true
		result = v8NodesNewBackReferenceNode(v8ASTCaptureStartRegister(capture.index),
			v8ASTCaptureEndRegister(capture.index), compiler.ReadBackward(), result, compiler.Flags())
	}
	return result
}

func (lowerer *v8TonodeGraphLowerer) emptyToNodeImpl(_ *v8ASTEmpty, _ *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	return success
}

func (lowerer *v8TonodeGraphLowerer) groupToNodeImpl(group *v8ASTGroup, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if group == nil {
		panic("V8 Group::ToNodeImpl requires a group")
	}
	scope := v8TonodeNewModifiersScope(compiler, group.flags)
	defer scope.close()
	return lowerer.treeToNode(group.body, compiler, success)
}

func (lowerer *v8TonodeGraphLowerer) boundaryAssertionAsLookaround(compiler *v8CompilerCompiler, success v8NodesNode, assertionType v8ASTAssertionType) v8NodesNode {
	if !v8CompilerNeedsUnicodeCaseEquivalents(compiler.Flags()) {
		panic("Unicode boundary lookaround requested without Unicode ignore-case")
	}
	wordRanges := v8TonodeAddClassEscape(v8ASTWord, nil, true)
	stackRegister := compiler.UnicodeLookaroundStackRegister()
	positionRegister := compiler.UnicodeLookaroundPositionRegister()
	choice := v8NodesNewChoiceNode(2, compiler.Flags(), compiler.Zone())
	for index := 0; index < 2; index++ {
		lookbehindForWord := index == 0
		lookaheadForWord := assertionType == v8ASTBoundary != lookbehindForWord
		lookbehind := lowerer.newLookaroundBuilder(lookbehindForWord, success, compiler, stackRegister, positionRegister, 0, -1)
		oldDirection := compiler.ReadBackward()
		compiler.SetReadBackward(true)
		backward := lowerer.textForRanges(wordRanges, compiler, lookbehind.onMatchSuccessNode())
		compiler.SetReadBackward(oldDirection)
		prior := lowerer.lookaroundBuilderForMatch(lookbehind, compiler, backward)
		lookahead := lowerer.newLookaroundBuilder(lookaheadForWord, prior, compiler, stackRegister, positionRegister, 0, -1)
		compiler.SetReadBackward(false)
		forward := lowerer.textForRanges(wordRanges, compiler, lookahead.onMatchSuccessNode())
		compiler.SetReadBackward(oldDirection)
		choice.AddAlternative(v8NodesNewGuardedAlternative(lowerer.lookaroundBuilderForMatch(lookahead, compiler, forward)))
	}
	return choice
}

func (lowerer *v8TonodeGraphLowerer) assertionToNodeImpl(assertion *v8ASTAssertion, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if assertion == nil {
		panic("V8 Assertion::ToNodeImpl requires an assertion")
	}
	flags := compiler.Flags()
	switch assertion.assertionType {
	case v8ASTStartOfLine:
		return v8NodesNewAssertionNode(v8NodesAssertAfterNewline, success, flags)
	case v8ASTStartOfInput:
		return v8NodesNewAssertionNode(v8NodesAssertAtStart, success, flags)
	case v8ASTBoundary, v8ASTNonBoundary:
		if v8CompilerNeedsUnicodeCaseEquivalents(flags) {
			return lowerer.boundaryAssertionAsLookaround(compiler, success, assertion.assertionType)
		}
		kind := v8NodesAssertAtBoundary
		if assertion.assertionType == v8ASTNonBoundary {
			kind = v8NodesAssertAtNonBoundary
		}
		return v8NodesNewAssertionNode(kind, success, flags)
	case v8ASTEndOfInput:
		return v8NodesNewAssertionNode(v8NodesAssertAtEnd, success, flags)
	case v8ASTEndOfBuffer:
		stackRegister, positionRegister := compiler.AllocateRegister(), compiler.AllocateRegister()
		lookahead := lowerer.newLookaroundBuilder(true, success, compiler, stackRegister, positionRegister, 0, -1)
		innerSuccess := lookahead.onMatchSuccessNode()
		crlfEnd := v8NodesNewAssertionNode(v8NodesAssertAtEnd, innerSuccess, flags)
		crlf := &v8ASTAtom{data: []uint16{'\r', '\n'}}
		crlfNode := lowerer.atomToNodeImpl(crlf, compiler, crlfEnd)
		lineClass := &v8ASTClassRanges{set: v8ASTCharacterSet{standardSet: func() *v8ASTStandardCharacterSet { set := v8ASTLineTerminator; return &set }()}}
		lineEnd := v8NodesNewAssertionNode(v8NodesAssertAtEnd, innerSuccess, flags)
		lineNode := lowerer.classRangesToNodeImpl(lineClass, compiler, lineEnd)
		innerChoice := v8NodesNewChoiceNode(2, flags, compiler.Zone())
		innerChoice.AddAlternative(v8NodesNewGuardedAlternative(crlfNode))
		innerChoice.AddAlternative(v8NodesNewGuardedAlternative(lineNode))
		lookaheadNode := lowerer.lookaroundBuilderForMatch(lookahead, compiler, innerChoice)
		outerChoice := v8NodesNewChoiceNode(2, flags, compiler.Zone())
		outerChoice.AddAlternative(v8NodesNewGuardedAlternative(lookaheadNode))
		outerChoice.AddAlternative(v8NodesNewGuardedAlternative(v8NodesNewAssertionNode(v8NodesAssertAtEnd, success, flags)))
		return outerChoice
	case v8ASTEndOfLine:
		stackRegister, positionRegister := compiler.AllocateRegister(), compiler.AllocateRegister()
		choice := v8NodesNewChoiceNode(2, flags, compiler.Zone())
		submatchSuccess := v8CompilerPositiveSubmatchSuccess(stackRegister, positionRegister, 0, -1, success, flags)
		lineClass := &v8ASTClassRanges{set: v8ASTCharacterSet{standardSet: func() *v8ASTStandardCharacterSet { set := v8ASTLineTerminator; return &set }()}}
		lineMatcher := lowerer.classRangesToNodeImpl(lineClass, compiler, submatchSuccess)
		lookahead := v8CompilerBeginPositiveSubmatch(stackRegister, positionRegister, lineMatcher, submatchSuccess, flags)
		choice.AddAlternative(v8NodesNewGuardedAlternative(lookahead))
		choice.AddAlternative(v8NodesNewGuardedAlternative(v8NodesNewAssertionNode(v8NodesAssertAtEnd, success, flags)))
		return choice
	default:
		panic("unknown V8 assertion kind")
	}
}

func v8TonodeStartsWithAtom(tree any) bool {
	switch value := tree.(type) {
	case *v8ASTAtom:
		return true
	case *v8ASTText:
		return len(value.elements) > 0 && value.elements[0].textType == v8ASTTextAtom
	default:
		return false
	}
}

func v8TonodeFirstAtom(tree any) *v8ASTAtom {
	switch value := tree.(type) {
	case *v8ASTAtom:
		return value
	case *v8ASTText:
		if len(value.elements) > 0 && value.elements[0].textType == v8ASTTextAtom {
			atom, _ := value.elements[0].tree.(*v8ASTAtom)
			return atom
		}
	}
	return nil
}

func v8TonodeCompareCaseInsensitive(flags Flags, lhs, rhs uint16) int {
	mode := v8TonodeCaseFoldingModeForFlags(flags)
	if lhs == rhs {
		return 0
	}
	left := v8ASTCaseFoldingEquivalenceKey(rune(lhs), mode)
	right := v8ASTCaseFoldingEquivalenceKey(rune(rhs), mode)
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

func v8TonodeCompareFirstChar(lhs, rhs any, flags Flags) int {
	left, right := v8TonodeFirstAtom(lhs), v8TonodeFirstAtom(rhs)
	if left == nil || right == nil || len(left.data) == 0 || len(right.data) == 0 {
		panic("V8 atom-run comparison requires nonempty first atoms")
	}
	a, b := left.data[0], right.data[0]
	if flags.IgnoreCase {
		return v8TonodeCompareCaseInsensitive(flags, a, b)
	}
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func v8TonodeCompareFirstCharCaseInsensitive(lhs, rhs any, mode v8TonodeCaseFoldingMode, canonicalize func(rune) rune) int {
	left, right := v8TonodeFirstAtom(lhs), v8TonodeFirstAtom(rhs)
	if left == nil || right == nil || len(left.data) == 0 || len(right.data) == 0 {
		panic("V8 case-insensitive atom comparison requires nonempty first atoms")
	}
	if mode == v8TonodeUnicodeCaseFolding {
		return v8TonodeCompareCaseInsensitiveUnicode(mode, left.data[0], right.data[0])
	}
	return v8TonodeCompareCaseInsensitiveLegacy(canonicalize, rune(left.data[0]), rune(right.data[0]))
}

func (lowerer *v8TonodeGraphLowerer) disjunctionSortConsecutiveAtoms(disjunction *v8ASTDisjunction, compiler *v8CompilerCompiler) bool {
	found := false
	items := disjunction.alternatives
	for index := 0; index < len(items); {
		for index < len(items) && !v8TonodeStartsWithAtom(items[index]) {
			index++
		}
		start := index
		for index < len(items) && v8TonodeStartsWithAtom(items[index]) {
			index++
		}
		if index-start < 2 {
			continue
		}
		found = true
		flags := compiler.Flags()
		sort.SliceStable(items[start:index], func(i, j int) bool {
			return v8TonodeCompareFirstChar(items[start+i], items[start+j], flags) < 0
		})
	}
	return found
}

func v8TonodeCharactersEqual(flags Flags, lhs, rhs uint16) bool {
	if lhs == rhs {
		return true
	}
	return flags.IgnoreCase && v8TonodeCompareCaseInsensitive(flags, lhs, rhs) == 0
}

func v8TonodeAppendTextElement(text *v8ASTText, element v8ASTTextElement) {
	text.elements = append(text.elements, element)
	text.length += v8CompilerTextElementLength(element)
}

func (lowerer *v8TonodeGraphLowerer) disjunctionRationalizeConsecutiveAtoms(disjunction *v8ASTDisjunction, compiler *v8CompilerCompiler) {
	items := disjunction.alternatives
	write, index := 0, 0
	flags := compiler.Flags()
	for index < len(items) {
		if !v8TonodeStartsWithAtom(items[index]) {
			items[write] = items[index]
			write++
			index++
			continue
		}
		first := index
		firstAtom := v8TonodeFirstAtom(items[first])
		if len(firstAtom.data) == 0 {
			panic("V8 atom rationalization found an empty prefix")
		}
		prefixLength := len(firstAtom.data)
		index++
		for index < len(items) && v8TonodeStartsWithAtom(items[index]) {
			atom := v8TonodeFirstAtom(items[index])
			if !v8TonodeCharactersEqual(flags, atom.data[0], firstAtom.data[0]) {
				break
			}
			if len(atom.data) < prefixLength {
				prefixLength = len(atom.data)
			}
			index++
		}
		// Two alternatives are enough to justify factoring a shared prefix.
		// Besides reducing graph nodes for ordinary `abc|abd` patterns, this is
		// the step that turns a sorted finite string set into a prefix trie.
		// Requiring three branches left every binary fork duplicated in the
		// generated matcher, which compounds badly for large Unicode string
		// properties.
		if index-first < 2 {
			for cursor := first; cursor < index; cursor++ {
				items[write] = items[cursor]
				write++
			}
			continue
		}
		for cursor := first + 1; cursor < index && prefixLength > 1; cursor++ {
			other := v8TonodeFirstAtom(items[cursor])
			for offset := 1; offset < prefixLength; offset++ {
				if !v8TonodeCharactersEqual(flags, firstAtom.data[offset], other.data[offset]) {
					prefixLength = offset
					break
				}
			}
		}
		prefix := &v8ASTAtom{data: append([]uint16(nil), firstAtom.data[:prefixLength]...)}
		suffixes := make([]any, 0, index-first)
		for cursor := first; cursor < index; cursor++ {
			switch alternative := items[cursor].(type) {
			case *v8ASTAtom:
				if len(alternative.data) == prefixLength {
					suffixes = append(suffixes, &v8ASTEmpty{})
				} else {
					suffixes = append(suffixes, &v8ASTAtom{data: append([]uint16(nil), alternative.data[prefixLength:]...)})
				}
			case *v8ASTText:
				oldFirst := v8TonodeFirstAtom(alternative)
				text := &v8ASTText{}
				if len(oldFirst.data) > prefixLength {
					v8TonodeAppendTextElement(text, v8CompilerTextElementFromAtom(&v8ASTAtom{data: append([]uint16(nil), oldFirst.data[prefixLength:]...)}))
				}
				for _, element := range alternative.elements[1:] {
					v8TonodeAppendTextElement(text, element)
				}
				if len(text.elements) == 0 {
					suffixes = append(suffixes, &v8ASTEmpty{})
				} else {
					suffixes = append(suffixes, text)
				}
			default:
				panic("V8 atom factoring encountered a non-atom alternative")
			}
		}
		suffixTree := any(&v8ASTDisjunction{alternatives: suffixes})
		items[write] = &v8ASTAlternative{nodes: []any{prefix, suffixTree}}
		write++
	}
	disjunction.alternatives = items[:write]
}

func (lowerer *v8TonodeGraphLowerer) disjunctionFixSingleCharacterDisjunctions(disjunction *v8ASTDisjunction, compiler *v8CompilerCompiler) {
	items := disjunction.alternatives
	write, index := 0, 0
	flags := compiler.Flags()
	for index < len(items) {
		atom, ok := items[index].(*v8ASTAtom)
		if !ok || len(atom.data) != 1 {
			items[write] = items[index]
			write++
			index++
			continue
		}
		start := index
		containsTrail := uint16(v8TonodeTrailSurrogateStart) <= atom.data[0] && atom.data[0] <= uint16(v8TonodeTrailSurrogateEnd)
		for index < len(items) {
			current, ok := items[index].(*v8ASTAtom)
			if !ok || len(current.data) != 1 {
				break
			}
			character := current.data[0]
			containsTrail = containsTrail || uint16(v8TonodeTrailSurrogateStart) <= character && character <= uint16(v8TonodeTrailSurrogateEnd)
			index++
		}
		if index-start > 1 {
			ranges := make([]v8ASTCharacterRange, 0, index-start)
			for cursor := start; cursor < index; cursor++ {
				character := items[cursor].(*v8ASTAtom).data[0]
				ranges = append(ranges, v8ASTRangeSingleton(rune(character)))
			}
			classFlags := v8ASTClassRangesFlags(0)
			if IsEitherUnicode(flags) && containsTrail {
				classFlags |= v8ASTClassContainsSplitSurrogate
			}
			items[write] = v8TonodeClassFromRanges(ranges, classFlags)
			write++
		} else {
			items[write] = items[start]
			write++
		}
	}
	disjunction.alternatives = items[:write]
}

func (lowerer *v8TonodeGraphLowerer) disjunctionToNodeImpl(disjunction *v8ASTDisjunction, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if disjunction == nil {
		panic("V8 Disjunction::ToNodeImpl requires a disjunction")
	}
	if len(disjunction.alternatives) > 1 {
		if !compiler.ReadBackward() && compiler.Optimize() {
			if lowerer.disjunctionSortConsecutiveAtoms(disjunction, compiler) {
				lowerer.disjunctionRationalizeConsecutiveAtoms(disjunction, compiler)
			}
		}
		lowerer.disjunctionFixSingleCharacterDisjunctions(disjunction, compiler)
		if len(disjunction.alternatives) == 1 {
			return lowerer.treeToNode(disjunction.alternatives[0], compiler, success)
		}
	}
	choice := v8NodesNewChoiceNode(len(disjunction.alternatives), compiler.Flags(), compiler.Zone())
	for _, alternative := range disjunction.alternatives {
		node := lowerer.treeToNode(alternative, compiler, success)
		if !node.IsBacktrack() {
			choice.AddAlternative(v8NodesNewGuardedAlternative(node))
		}
	}
	if choice.Alternatives().Len() == 0 {
		return lowerer.backtrack(compiler)
	}
	if choice.Alternatives().Len() == 1 {
		return choice.Alternatives().At(0).Node()
	}
	return choice
}

func v8TonodeRewriteAssertionSequence(terms []any, flags Flags) []any {
	hasAssertion, needsFlattening := false, false
	for _, term := range terms {
		switch value := term.(type) {
		case *v8ASTAssertion:
			hasAssertion = true
		case *v8ASTAlternative:
			needsFlattening = true
		case *v8ASTGroup:
			if value.flags == flags {
				needsFlattening = true
			}
		}
	}
	if !hasAssertion && !needsFlattening {
		return terms
	}
	flattened := make([]any, 0, len(terms))
	seen := uint32(0)
	sequenceStart := -1
	var appendTerm func(any)
	appendTerm = func(term any) {
		switch value := term.(type) {
		case *v8ASTAssertion:
			bit := uint32(1) << value.assertionType
			if seen&bit != 0 {
				return
			}
			seen |= bit
			if seen&(uint32(1)<<v8ASTBoundary) != 0 && seen&(uint32(1)<<v8ASTNonBoundary) != 0 {
				if sequenceStart < 0 {
					panic("V8 assertion rewrite lost its sequence start")
				}
				flattened = flattened[:sequenceStart]
				flattened = append(flattened, &v8ASTClassRanges{})
				return
			}
			if sequenceStart < 0 {
				sequenceStart = len(flattened)
			}
			flattened = append(flattened, term)
		case *v8ASTAlternative:
			for _, child := range value.nodes {
				appendTerm(child)
			}
		case *v8ASTGroup:
			if value.flags == flags {
				appendTerm(value.body)
			} else {
				seen, sequenceStart = 0, -1
				flattened = append(flattened, term)
			}
		default:
			seen, sequenceStart = 0, -1
			flattened = append(flattened, term)
		}
	}
	for _, term := range terms {
		appendTerm(term)
	}
	return flattened
}

func (lowerer *v8TonodeGraphLowerer) alternativeToNodeImpl(alternative *v8ASTAlternative, compiler *v8CompilerCompiler, success v8NodesNode) v8NodesNode {
	if alternative == nil {
		panic("V8 Alternative::ToNodeImpl requires an alternative")
	}
	children := v8TonodeRewriteAssertionSequence(alternative.nodes, compiler.Flags())
	current := success
	if compiler.ReadBackward() {
		for _, child := range children {
			current = lowerer.treeToNode(child, compiler, current)
		}
	} else {
		for index := len(children) - 1; index >= 0; index-- {
			current = lowerer.treeToNode(children[index], compiler, current)
		}
	}
	return current
}

// Ported V8 source declarations from regexp-compiler-core.
// Copyright 2019 the V8 project authors. All rights reserved.
// V8 RegExp compiler logic adapted from regexp-compiler.cc and
// regexp-compiler.h. See LICENSE.v8 for attribution and license terms.
const (
	v8CompilerNoRegister               = -1
	v8CompilerToNodeStackCheckInterval = 64
)

// v8CompilerNeedsUnicodeCaseEquivalents ports the compiler-header predicate
// that selects full Unicode case closure only when both flags are active.
func v8CompilerNeedsUnicodeCaseEquivalents(flags Flags) bool {
	return IsEitherUnicode(flags) && flags.IgnoreCase
}

// v8CompilerAppendTreeToText follows Tree::AppendToText's unreachable base
// case. The concrete atom and class-range cases append one TextElement.
func v8CompilerAppendTreeToText(tree any, text *v8ASTText) {
	switch value := tree.(type) {
	case *v8ASTAtom:
		v8CompilerAppendTextElement(text, v8CompilerTextElementFromAtom(value))
	case *v8ASTClassRanges:
		v8CompilerAppendTextElement(text, v8CompilerTextElementFromClassRanges(value))
	case *v8ASTText:
		for _, element := range value.elements {
			v8CompilerAppendTextElement(text, element)
		}
	default:
		panic("V8 Tree::AppendToText is unreachable for this AST node")
	}
}

func v8CompilerTextElementFromAtom(atom *v8ASTAtom) v8ASTTextElement {
	return v8ASTTextElement{textType: v8ASTTextAtom, tree: atom}
}

func v8CompilerTextElementFromClassRanges(ranges *v8ASTClassRanges) v8ASTTextElement {
	return v8ASTTextElement{textType: v8ASTTextClassRanges, tree: ranges}
}

// v8CompilerNewClassRangesFromList ports the zone-owned ClassRanges wrapper
// used by TextNode's range constructors. Go owns a copy of the range slice;
// the source node retains the class-range object and its exact input ranges.
func v8CompilerNewClassRangesFromList(zone *v8CompilerZone, ranges *v8NodesZoneList[v8NodesCharacterRange]) *v8ASTClassRanges {
	if zone == nil || ranges == nil {
		panic("V8 text class-range construction requires a zone and ranges")
	}
	return &v8ASTClassRanges{
		set: v8ASTCharacterSet{ranges: append([]v8ASTCharacterRange(nil), ranges.items...)},
	}
}

// v8CompilerTextNodeFactory ports the concrete TextNode constructors declared
// in regexp-nodes.h and defined in regexp-compiler.cc. The factory makes the
// source's one-element atom-vs-class distinction for surrogate-pair prefixes.
type v8CompilerTextNodeFactory struct{}

var _ v8NodesTextNodeFactories = v8CompilerTextNodeFactory{}

func (v8CompilerTextNodeFactory) NewTextNode(elements *v8NodesZoneList[v8NodesTextElement], readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode {
	if elements == nil || elements.Len() == 0 || onSuccess == nil {
		panic("V8 TextNode construction requires elements and a successor")
	}
	return v8NodesNewTextNode(elements, readBackward, onSuccess, flags)
}

func (factory v8CompilerTextNodeFactory) NewTextNodeForClassRanges(ranges *v8NodesClassRanges, readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode {
	if ranges == nil {
		panic("V8 single-class TextNode requires class ranges")
	}
	elements := &v8NodesZoneList[v8NodesTextElement]{}
	elements.Add(v8CompilerTextElementFromClassRanges(ranges))
	return factory.NewTextNode(elements, readBackward, onSuccess, flags)
}

func (factory v8CompilerTextNodeFactory) CreateForCharacterRanges(zone *v8NodesZone, ranges *v8NodesZoneList[v8NodesCharacterRange], readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode {
	class := v8CompilerNewClassRangesFromList(zone, ranges)
	return factory.NewTextNodeForClassRanges(class, readBackward, onSuccess, flags)
}

func (factory v8CompilerTextNodeFactory) CreateForSurrogatePair(zone *v8NodesZone, lead v8NodesCharacterRange, trailRanges *v8NodesZoneList[v8NodesCharacterRange], readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode {
	if zone == nil || trailRanges == nil {
		panic("V8 surrogate-pair TextNode requires a zone and trail ranges")
	}
	elements := &v8NodesZoneList[v8NodesTextElement]{}
	if lead.from == lead.to {
		if lead.from < 0 || lead.from > 0xffff {
			panic("V8 surrogate-pair lead atom is outside UTF-16")
		}
		atom := &v8ASTAtom{data: []uint16{uint16(lead.from)}}
		elements.Add(v8CompilerTextElementFromAtom(atom))
	} else {
		elements.Add(v8CompilerTextElementFromClassRanges(v8CompilerNewClassRangesFromList(zone, &v8NodesZoneList[v8NodesCharacterRange]{items: []v8NodesCharacterRange{lead}})))
	}
	elements.Add(v8CompilerTextElementFromClassRanges(v8CompilerNewClassRangesFromList(zone, trailRanges)))
	return factory.NewTextNode(elements, readBackward, onSuccess, flags)
}

func (factory v8CompilerTextNodeFactory) CreateForSurrogatePairRanges(zone *v8NodesZone, leadRanges *v8NodesZoneList[v8NodesCharacterRange], trail v8NodesCharacterRange, readBackward bool, onSuccess v8NodesNode, flags v8NodesFlags) *v8NodesTextNode {
	if zone == nil || leadRanges == nil {
		panic("V8 surrogate-pair TextNode requires a zone and lead ranges")
	}
	elements := &v8NodesZoneList[v8NodesTextElement]{}
	elements.Add(v8CompilerTextElementFromClassRanges(v8CompilerNewClassRangesFromList(zone, leadRanges)))
	trailRanges := &v8NodesZoneList[v8NodesCharacterRange]{items: []v8NodesCharacterRange{trail}}
	elements.Add(v8CompilerTextElementFromClassRanges(v8CompilerNewClassRangesFromList(zone, trailRanges)))
	return factory.NewTextNode(elements, readBackward, onSuccess, flags)
}

func v8CompilerTextElementLength(element v8ASTTextElement) int {
	switch element.textType {
	case v8ASTTextAtom:
		atom, ok := element.tree.(*v8ASTAtom)
		if !ok || atom == nil {
			panic("V8 TextElement atom payload is missing")
		}
		return len(atom.data)
	case v8ASTTextClassRanges:
		return 1
	default:
		panic("unknown V8 TextElement kind")
	}
}

func v8CompilerAppendTextElement(text *v8ASTText, element v8ASTTextElement) {
	text.elements = append(text.elements, element)
	text.length += v8CompilerTextElementLength(element)
}

// CanMatchLatin1 ports TextNode::CanMatchLatin1. The graph builder uses this
// when constructing a one-byte matcher to discard nodes which cannot match
// any Latin-1 input.
func (node *v8NodesTextNode) CanMatchLatin1(compiler *v8CompilerCompiler) bool {
	if node == nil || compiler == nil {
		panic("V8 TextNode::CanMatchLatin1 requires node and compiler")
	}
	if node.finiteStringSet != nil {
		// The finite set may contain Latin-1 alternatives; keeping it is the
		// conservative answer for the one-byte graph's dead-node pruning.
		return true
	}
	for index := 0; index < node.elements.Len(); index++ {
		element := *node.elements.At(index)
		if element.textType == v8ASTTextAtom {
			atom, ok := element.tree.(*v8ASTAtom)
			if !ok || atom == nil {
				panic("V8 TextNode::CanMatchLatin1 found invalid atom")
			}
			for _, character := range atom.data {
				if !node.flags.IgnoreCase {
					if character > 0xff {
						return false
					}
					continue
				}
				letters := v8CompilerCaseIndependentLetters(node.flags, character, false)
				if len(letters) == 0 || letters[0] > 0xff {
					return false
				}
			}
			continue
		}
		if element.textType != v8ASTTextClassRanges {
			panic("unknown V8 TextElement kind in CanMatchLatin1")
		}
		class, ok := element.tree.(*v8ASTClassRanges)
		if !ok || class == nil {
			panic("V8 TextNode::CanMatchLatin1 found invalid class ranges")
		}
		ranges := v8TonodeCanonicalizeRanges(v8TonodeCharacterSetRanges(&class.set))
		if len(ranges) == 0 && class.flags&v8ASTClassNegated == 0 {
			return false
		}
		caseComplications := !IsEitherUnicode(node.flags) && node.flags.IgnoreCase && v8CompilerRangesContainLatin1Equivalents(ranges)
		if class.flags&v8ASTClassNegated != 0 {
			if len(ranges) != 0 && ranges[0].from == 0 && ranges[0].to >= 0xff && !caseComplications {
				return false
			}
		} else if (len(ranges) == 0 || ranges[0].from > 0xff) && !caseComplications {
			return false
		}
	}
	return true
}

func v8CompilerCaseIndependentLetters(flags Flags, character uint16, oneByte bool) []uint16 {
	if !flags.IgnoreCase {
		if oneByte && character > 0xff {
			return nil
		}
		return []uint16{character}
	}
	for _, class := range regexpCaseFoldClasses {
		contains := false
		for _, member := range class {
			if member == rune(character) {
				contains = true
				break
			}
		}
		if !contains {
			continue
		}
		letters := make([]uint16, 0, len(class))
		for _, member := range class {
			if member < 0 || member > 0xffff {
				continue
			}
			if oneByte && member > 0xff {
				continue
			}
			candidate := uint16(member)
			if v8TonodeCompareCaseInsensitive(flags, character, candidate) == 0 {
				if len(letters) == 4 {
					panic("V8 case-equivalence class exceeds its four-character buffer")
				}
				letters = append(letters, candidate)
			}
		}
		if len(letters) != 0 {
			return letters
		}
	}
	if oneByte && character > 0xff {
		return nil
	}
	return []uint16{character}
}

func v8CompilerRangesContainLatin1Equivalents(ranges []v8ASTCharacterRange) bool {
	for _, current := range ranges {
		if current.contains(0x039c) || current.contains(0x03bc) || current.contains(0x0178) {
			return true
		}
	}
	return false
}

// v8CompilerZone preserves the ownership boundary used by V8 graph constructors.
// Go's garbage collector owns the actual graph allocations.
type v8CompilerZone struct{}

type v8CompilerCompilationError uint8

const (
	v8CompilerNoCompilationError v8CompilerCompilationError = iota
	v8CompilerTooLargeError
)

type v8CompilerAnalysisStackChecker interface {
	HasOverflowed() bool
}

type v8CompilerCompiler struct {
	accept                            *v8NodesEndNode
	nextRegister                      int
	unicodeLookaroundStackRegister    int
	unicodeLookaroundPositionRegister int
	flags                             Flags
	oneByte                           bool
	regExpTooBig                      bool
	toNodeOverflowCheckTicks          uint64
	optimize                          bool
	readBackward                      bool
	hasSearchPrefix                   bool
	currentExpansionFactor            int
	recursionDepth                    int
	limitingRecursion                 bool
	frequencyCollator                 v8CompilerFrequencyCollator
	zone                              *v8CompilerZone
	maxRegister                       int
}

// The non-macOS limit is the one used by V8's default compiler build. V8 uses
// 50 on macOS because its secondary-thread stacks are smaller; the Go emitter
// does not run recursive assembler emission, but retains this source limit for
// graph analyses that ask Node.KeepRecursing.
const v8CompilerMaxRecursion = 100

// NewCompiler initializes the state needed by V8's AST-to-graph lowering.
// The generated-Go emitter owns execution; this compiler has no assembler,
// bytecode, work list, or interpreter state.
func v8CompilerNewCompiler(zone *v8CompilerZone, captureCount int, flags Flags, oneByte bool, maxRegister int) *v8CompilerCompiler {
	if captureCount < 0 {
		panic("negative V8 capture count")
	}
	nextRegister := 2 * (captureCount + 1)
	if maxRegister < nextRegister {
		panic("V8 capture registers exceed the graph register limit")
	}
	compiler := &v8CompilerCompiler{
		nextRegister:                      nextRegister,
		unicodeLookaroundStackRegister:    v8CompilerNoRegister,
		unicodeLookaroundPositionRegister: v8CompilerNoRegister,
		flags:                             flags,
		oneByte:                           oneByte,
		optimize:                          true,
		currentExpansionFactor:            1,
		frequencyCollator:                 v8CompilerNewFrequencyCollator(),
		zone:                              zone,
		maxRegister:                       maxRegister,
	}
	compiler.accept = v8NodesNewEndNode(v8NodesAccept, flags, zone)
	return compiler
}

func (compiler *v8CompilerCompiler) AllocateRegister() int {
	if compiler.nextRegister >= compiler.maxRegister {
		compiler.SetRegExpTooBig()
		return compiler.nextRegister
	}
	register := compiler.nextRegister
	compiler.nextRegister++
	return register
}

func (compiler *v8CompilerCompiler) UnicodeLookaroundStackRegister() int {
	if compiler.unicodeLookaroundStackRegister == v8CompilerNoRegister {
		compiler.unicodeLookaroundStackRegister = compiler.AllocateRegister()
	}
	return compiler.unicodeLookaroundStackRegister
}

func (compiler *v8CompilerCompiler) UnicodeLookaroundPositionRegister() int {
	if compiler.unicodeLookaroundPositionRegister == v8CompilerNoRegister {
		compiler.unicodeLookaroundPositionRegister = compiler.AllocateRegister()
	}
	return compiler.unicodeLookaroundPositionRegister
}

func (compiler *v8CompilerCompiler) SetRegExpTooBig()           { compiler.regExpTooBig = true }
func (compiler *v8CompilerCompiler) IsRegExpTooBig() bool       { return compiler.regExpTooBig }
func (compiler *v8CompilerCompiler) Zone() *v8CompilerZone      { return compiler.zone }
func (compiler *v8CompilerCompiler) Flags() Flags               { return compiler.flags }
func (compiler *v8CompilerCompiler) SetFlags(flags Flags)       { compiler.flags = flags }
func (compiler *v8CompilerCompiler) OneByte() bool              { return compiler.oneByte }
func (compiler *v8CompilerCompiler) Optimize() bool             { return compiler.optimize }
func (compiler *v8CompilerCompiler) SetOptimize(value bool)     { compiler.optimize = value }
func (compiler *v8CompilerCompiler) ReadBackward() bool         { return compiler.readBackward }
func (compiler *v8CompilerCompiler) SetReadBackward(value bool) { compiler.readBackward = value }
func (compiler *v8CompilerCompiler) HasSearchPrefix() bool      { return compiler.hasSearchPrefix }
func (compiler *v8CompilerCompiler) CurrentExpansionFactor() int {
	return compiler.currentExpansionFactor
}
func (compiler *v8CompilerCompiler) SetCurrentExpansionFactor(value int) {
	compiler.currentExpansionFactor = value
}
func (compiler *v8CompilerCompiler) RecursionDepth() int { return compiler.recursionDepth }
func (compiler *v8CompilerCompiler) IncrementRecursionDepth() {
	compiler.recursionDepth++
}
func (compiler *v8CompilerCompiler) DecrementRecursionDepth() {
	compiler.recursionDepth--
}
func (compiler *v8CompilerCompiler) LimitingRecursion() bool {
	return compiler.limitingRecursion
}
func (compiler *v8CompilerCompiler) SetLimitingRecursion(value bool) {
	compiler.limitingRecursion = value
}
func (compiler *v8CompilerCompiler) FrequencyCollator() *v8CompilerFrequencyCollator {
	return &compiler.frequencyCollator
}

type v8CompilerSurrogateLookaroundBuilder interface {
	OnMatchSuccess() v8NodesNode
	ForMatch(compiler *v8CompilerCompiler, matchTrail v8NodesNode) v8NodesGuardedAlternative
	RegisterNode(v8NodesNode)
}

type v8CompilerSurrogateGraphFactory interface {
	CreateForCharacterRanges(zone *v8CompilerZone, ranges *v8NodesZoneList[v8ASTCharacterRange], readBackward bool, onSuccess v8NodesNode, flags Flags) *v8NodesTextNode
	NewChoiceNode(expectedAlternatives int, flags Flags, zone *v8CompilerZone) *v8NodesChoiceNode
	NewUnicodeLookaroundBuilder(forward bool, stepBack v8NodesNode, compiler *v8CompilerCompiler, stackRegister, positionRegister int) v8CompilerSurrogateLookaroundBuilder
	RegisterNode(v8NodesNode)
}

// v8CompilerOptionallyStepBackToLeadSurrogate ports the graph construction
// from Compiler::OptionallyStepBackToLeadSurrogate. The lookaround builder is
// an explicit contract because its implementation belongs to the separate
// lookaround graph-construction source file.
func v8CompilerOptionallyStepBackToLeadSurrogate(compiler *v8CompilerCompiler, onSuccess v8NodesNode, factory v8CompilerSurrogateGraphFactory) v8NodesNode {
	if compiler == nil || onSuccess == nil || factory == nil {
		panic("V8 optional surrogate rewind requires compiler, successor, and graph factory")
	}
	if compiler.readBackward {
		panic("V8 optional surrogate rewind is only valid for forward compilation")
	}
	leadRanges := &v8NodesZoneList[v8ASTCharacterRange]{items: []v8ASTCharacterRange{{from: 0xd800, to: 0xdbff}}}
	trailRanges := &v8NodesZoneList[v8ASTCharacterRange]{items: []v8ASTCharacterRange{{from: 0xdc00, to: 0xdfff}}}
	stackRegister := compiler.UnicodeLookaroundStackRegister()
	positionRegister := compiler.UnicodeLookaroundPositionRegister()
	stepBack := factory.CreateForCharacterRanges(compiler.zone, leadRanges, true, onSuccess, compiler.flags)
	builder := factory.NewUnicodeLookaroundBuilder(true, stepBack, compiler, stackRegister, positionRegister)
	if builder == nil {
		panic("V8 Unicode lookaround factory returned a nil builder")
	}
	matchTrail := factory.CreateForCharacterRanges(compiler.zone, trailRanges, false, builder.OnMatchSuccess(), compiler.flags)
	builder.RegisterNode(stepBack)
	builder.RegisterNode(matchTrail)
	choice := factory.NewChoiceNode(2, compiler.flags, compiler.zone)
	choice.AddAlternative(builder.ForMatch(compiler, matchTrail))
	choice.AddAlternative(v8NodesNewGuardedAlternative(onSuccess))
	factory.RegisterNode(choice)
	return choice
}

type v8CompilerPreprocessRequest struct {
	tree           any
	containsAnchor bool
	compileError   v8CompilerCompilationError
}

type v8CompilerPreprocessGraphFactory interface {
	v8CompilerSurrogateGraphFactory
	CaptureToNode(tree any, compiler *v8CompilerCompiler, accept v8NodesNode) v8NodesNode
	EverythingClassRanges(zone *v8CompilerZone) *v8ASTClassRanges
	QuantifierToNode(minimum, maximum int, greedy bool, body *v8ASTClassRanges, compiler *v8CompilerCompiler, success v8NodesNode, containsAnchor bool) v8NodesNode
	NewTextNodeForClassRanges(class *v8ASTClassRanges, readBackward bool, success v8NodesNode, flags Flags) *v8NodesTextNode
}

// v8CompilerPreprocessRegExp ports Compiler::PreprocessRegExp's ordering:
// accept registration, capture-zero wrapper, optional lazy search prefix,
// anchor-sensitive one-step unrolling, then Unicode lead-surrogate rewind.
// AST lowering and quantifier/lookaround construction remain narrow factory
// contracts owned by the corresponding source modules.
func v8CompilerPreprocessRegExp(compiler *v8CompilerCompiler, request *v8CompilerPreprocessRequest, oneByte bool, factory v8CompilerPreprocessGraphFactory) v8NodesNode {
	if compiler == nil || request == nil || request.tree == nil || factory == nil {
		panic("V8 PreprocessRegExp requires compiler, request, AST, and graph factory")
	}
	factory.RegisterNode(compiler.accept)
	capturedBody := factory.CaptureToNode(request.tree, compiler, compiler.accept)
	if capturedBody == nil {
		panic("V8 capture-zero lowering returned a nil graph")
	}
	node := capturedBody
	if !v8ASTCertainlyAnchoredAtStart(request.tree, v8NodesRecursionBudget) && !IsSticky(compiler.flags) {
		compiler.hasSearchPrefix = true
		loopBodyClass := factory.EverythingClassRanges(compiler.zone)
		loop := factory.QuantifierToNode(0, v8ASTInfinity, false, loopBodyClass, compiler, capturedBody, request.containsAnchor)
		if loop == nil {
			panic("V8 unanchored search quantifier lowering returned nil")
		}
		if request.containsAnchor {
			firstStep := factory.NewChoiceNode(2, compiler.flags, compiler.zone)
			firstStep.AddAlternative(v8NodesNewGuardedAlternative(capturedBody))
			advanceClass := factory.EverythingClassRanges(compiler.zone)
			firstStep.AddAlternative(v8NodesNewGuardedAlternative(factory.NewTextNodeForClassRanges(advanceClass, false, loop, compiler.flags)))
			factory.RegisterNode(firstStep)
			node = firstStep
		} else {
			node = loop
		}
	}
	if !oneByte && IsEitherUnicode(compiler.flags) && (IsGlobal(compiler.flags) || IsSticky(compiler.flags)) {
		node = v8CompilerOptionallyStepBackToLeadSurrogate(compiler, node, factory)
	}
	if compiler.regExpTooBig {
		request.compileError = v8CompilerTooLargeError
	}
	if node == nil {
		panic("V8 PreprocessRegExp produced a nil graph")
	}
	return node
}

func (compiler *v8CompilerCompiler) ToNodeCheckForStackOverflow(stack v8CompilerAnalysisStackChecker) {
	if compiler == nil {
		panic("V8 ToNode stack check requires a compiler")
	}
	if stack != nil && stack.HasOverflowed() {
		compiler.SetRegExpTooBig()
	}
}

// ToNodeMaybeCheckForStackOverflow preserves Compiler's periodic check policy.
// The first recursive-construction call checks immediately, then every 64th
// call; the checker itself remains injected because the source uses Isolate's
// stack limit.
func (compiler *v8CompilerCompiler) ToNodeMaybeCheckForStackOverflow(stack v8CompilerAnalysisStackChecker) {
	if compiler == nil {
		panic("V8 periodic ToNode stack check requires a compiler")
	}
	if compiler.toNodeOverflowCheckTicks%v8CompilerToNodeStackCheckInterval == 0 {
		compiler.ToNodeCheckForStackOverflow(stack)
	}
	compiler.toNodeOverflowCheckTicks++
}

// Action factories translate the union payloads in regexp-nodes.h and
// regexp-compiler.cc. They preserve the source register ranges and payloads.
func v8CompilerSetRegisterForLoop(register, value int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8NodesNewSimpleActionNode(v8NodesSetRegisterForLoop, onSuccess, flags, register, register, value)
}
func v8CompilerIncrementRegister(register int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8NodesNewSimpleActionNode(v8NodesIncrementRegister, onSuccess, flags, register, register, 0)
}
func v8CompilerStorePosition(register int, onSuccess v8NodesNode, flags Flags, capture *v8NodesCapture) *v8NodesActionNode {
	node := v8NodesNewSimpleActionNode(v8NodesStorePosition, onSuccess, flags, register, register, 0)
	node.simple.capture = capture
	return node
}
func v8CompilerRestorePosition(register int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8NodesNewSimpleActionNode(v8NodesRestorePosition, onSuccess, flags, register, register, 0)
}
func v8CompilerClearCaptures(interval v8NodesInterval, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8NodesNewSimpleActionNode(v8NodesClearCaptures, onSuccess, flags, interval.from, interval.to, 0)
}
func v8CompilerBeginPositiveSubmatch(stackRegister, positionRegister int, body v8NodesNode, success *v8NodesActionNode, flags Flags) *v8NodesActionNode {
	node := v8NodesNewActionNode(v8NodesBeginPositiveSubmatch, body, flags)
	node.submatch.stackPointerRegister = stackRegister
	node.submatch.currentPositionRegister = positionRegister
	node.submatch.successNode = success
	return node
}
func v8CompilerBeginNegativeSubmatch(stackRegister, positionRegister int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	node := v8NodesNewActionNode(v8NodesBeginNegativeSubmatch, onSuccess, flags)
	node.submatch.stackPointerRegister = stackRegister
	node.submatch.currentPositionRegister = positionRegister
	return node
}
func v8CompilerPositiveSubmatchSuccess(stackRegister, positionRegister, clearRegisterCount, clearRegisterFrom int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	node := v8NodesNewActionNode(v8NodesPositiveSubmatchSuccess, onSuccess, flags)
	node.submatch.stackPointerRegister = stackRegister
	node.submatch.currentPositionRegister = positionRegister
	node.submatch.clearRegisterCount = clearRegisterCount
	node.submatch.clearRegisterFrom = clearRegisterFrom
	return node
}
func v8CompilerEmptyMatchCheck(startRegister, repetitionRegister, repetitionLimit int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	node := v8NodesNewActionNode(v8NodesEmptyMatchCheck, onSuccess, flags)
	node.emptyMatch = v8NodesActionEmptyMatchData{startRegister: startRegister, repetitionRegister: repetitionRegister, repetitionLimit: repetitionLimit}
	return node
}
func v8CompilerEatsAtLeast(characters int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	node := v8NodesNewActionNode(v8NodesEatsAtLeastAction, onSuccess, flags)
	node.eatsAtLeast.characters = characters
	return node
}

// v8CompilerActionNodeFactory implements the complete static ActionNode
// factory surface used by the source graph builder.
type v8CompilerActionNodeFactory struct{}

var _ v8NodesActionNodeFactories = v8CompilerActionNodeFactory{}

func (v8CompilerActionNodeFactory) NewActionNode(actionType v8NodesActionType, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8NodesNewActionNode(actionType, onSuccess, flags)
}
func (v8CompilerActionNodeFactory) NewSimpleActionNode(actionType v8NodesActionType, onSuccess v8NodesNode, flags Flags, registerFrom, registerTo, value int) *v8NodesActionNode {
	return v8NodesNewSimpleActionNode(actionType, onSuccess, flags, registerFrom, registerTo, value)
}
func (v8CompilerActionNodeFactory) SetRegisterForLoop(register, value int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8CompilerSetRegisterForLoop(register, value, onSuccess, flags)
}
func (v8CompilerActionNodeFactory) IncrementRegister(register int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8CompilerIncrementRegister(register, onSuccess, flags)
}
func (v8CompilerActionNodeFactory) StorePosition(register int, onSuccess v8NodesNode, flags Flags, capture *v8NodesCapture) *v8NodesActionNode {
	return v8CompilerStorePosition(register, onSuccess, flags, capture)
}
func (v8CompilerActionNodeFactory) RestorePosition(register int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8CompilerRestorePosition(register, onSuccess, flags)
}
func (v8CompilerActionNodeFactory) ClearCaptures(interval v8NodesInterval, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8CompilerClearCaptures(interval, onSuccess, flags)
}
func (v8CompilerActionNodeFactory) BeginPositiveSubmatch(stackRegister, positionRegister int, body v8NodesNode, success *v8NodesActionNode, flags Flags) *v8NodesActionNode {
	return v8CompilerBeginPositiveSubmatch(stackRegister, positionRegister, body, success, flags)
}
func (v8CompilerActionNodeFactory) BeginNegativeSubmatch(stackRegister, positionRegister int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8CompilerBeginNegativeSubmatch(stackRegister, positionRegister, onSuccess, flags)
}
func (v8CompilerActionNodeFactory) PositiveSubmatchSuccess(stackRegister, positionRegister, clearRegisterCount, clearRegisterFrom int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8CompilerPositiveSubmatchSuccess(stackRegister, positionRegister, clearRegisterCount, clearRegisterFrom, onSuccess, flags)
}
func (v8CompilerActionNodeFactory) EmptyMatchCheck(startRegister, repetitionRegister, repetitionLimit int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8CompilerEmptyMatchCheck(startRegister, repetitionRegister, repetitionLimit, onSuccess, flags)
}
func (v8CompilerActionNodeFactory) EatsAtLeast(characters int, onSuccess v8NodesNode, flags Flags) *v8NodesActionNode {
	return v8CompilerEatsAtLeast(characters, onSuccess, flags)
}
