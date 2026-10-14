package regexpengine

import (
	"fmt"
	"go/format"
	"go/token"
	"math/bits"
	"sort"
	"strconv"
	"strings"
	"unicode/utf16"
)

// CompileToGo parses a pattern with the V8-shaped frontend, lowers it through
// the V8 graph builder, and emits a standalone Go matcher function. The
// function accepts UTF-16 code units so its indices match JavaScript strings.
// Unsupported graph nodes are reported instead of being sent to another
// regular-expression engine.
func CompileToGo(source, flags, functionName string) (string, error) {
	return CompileUTF16ToGo(utf16.Encode([]rune(source)), flags, functionName)
}

// CompileUTF16ToGo is the UTF-16-preserving form of CompileToGo.
func CompileUTF16ToGo(source []uint16, flagText, functionName string) (string, error) {
	if !token.IsIdentifier(functionName) || token.Lookup(functionName).IsKeyword() {
		return "", fmt.Errorf("invalid Go matcher function name %q", functionName)
	}
	flags, err := ParseFlags(flagText)
	if err != nil {
		return "", err
	}
	parsed, err := v8ASTParseRegExpUTF16(source, flags, nil)
	if err != nil {
		return "", err
	}
	var generated string
	if atom, ok := parsed.Tree.(*v8ASTAtom); ok && len(atom.data) > 0 &&
		parsed.CaptureCount == 0 && !flags.IgnoreCase && !flags.Linear && !flags.Sticky {
		generated = v8EmitGoAtomMatcher(atom.data, flags, functionName)
	} else {
		graph, graphErr := v8BuildGoMatcherGraph(parsed, flags)
		if graphErr != nil {
			return "", graphErr
		}
		generated, err = v8EmitGoMatcher(graph, functionName)
		if err != nil {
			return "", err
		}
	}
	formatted, err := format.Source([]byte(generated))
	if err != nil {
		return "", fmt.Errorf("format generated Go matcher: %w\n%s", err, generated)
	}
	return string(formatted), nil
}

type v8GoMatcherGraph struct {
	root                v8NodesNode
	searchBody          v8NodesNode
	compiler            *v8CompilerCompiler
	flags               Flags
	captureCount        int
	registerCount       int
	metadata            v8ASTParserCompileData
	firstSet            []v8ASTCharacterRange
	hasFirstSet         bool
	searchFirstSet      []v8ASTCharacterRange
	hasSearchFirstSet   bool
	searchMayMatchEmpty bool
	searchPlan          v8GoBoyerMooreSearchPlan
}

// v8GoPreprocessGraphFactory supplies the concrete graph constructors needed
// by the existing port of Compiler::PreprocessRegExp. It intentionally does
// not implement any assembler or emission interface.
type v8GoPreprocessGraphFactory struct {
	lowerer *v8TonodeGraphLowerer
	nodes   []v8NodesNode
	body    v8NodesNode
}

func (factory *v8GoPreprocessGraphFactory) RegisterNode(node v8NodesNode) {
	if node != nil {
		factory.nodes = append(factory.nodes, node)
	}
}

func (factory *v8GoPreprocessGraphFactory) CreateForCharacterRanges(zone *v8CompilerZone, ranges *v8NodesZoneList[v8ASTCharacterRange], readBackward bool, onSuccess v8NodesNode, flags Flags) *v8NodesTextNode {
	return (v8CompilerTextNodeFactory{}).CreateForCharacterRanges(zone, ranges, readBackward, onSuccess, flags)
}

func (factory *v8GoPreprocessGraphFactory) NewChoiceNode(expectedAlternatives int, flags Flags, zone *v8CompilerZone) *v8NodesChoiceNode {
	node := v8NodesNewChoiceNode(expectedAlternatives, flags, zone)
	factory.RegisterNode(node)
	return node
}

func (factory *v8GoPreprocessGraphFactory) NewUnicodeLookaroundBuilder(forward bool, stepBack v8NodesNode, compiler *v8CompilerCompiler, stackRegister, positionRegister int) v8CompilerSurrogateLookaroundBuilder {
	builder := factory.lowerer.newLookaroundBuilder(forward, stepBack, compiler, stackRegister, positionRegister, 0, -1)
	return v8GoSurrogateLookaroundBuilder{factory: factory, builder: builder}
}

func (factory *v8GoPreprocessGraphFactory) CaptureToNode(tree any, compiler *v8CompilerCompiler, accept v8NodesNode) v8NodesNode {
	factory.body = factory.lowerer.captureToNode(tree, compiler, accept, nil)
	return factory.body
}

func (*v8GoPreprocessGraphFactory) EverythingClassRanges(*v8CompilerZone) *v8ASTClassRanges {
	return v8ASTFactory{}.newClassRangesFromStandardSet(v8ASTEverything)
}

func (factory *v8GoPreprocessGraphFactory) QuantifierToNode(minimum, maximum int, greedy bool, body *v8ASTClassRanges, compiler *v8CompilerCompiler, success v8NodesNode, containsAnchor bool) v8NodesNode {
	return factory.lowerer.quantifierToNode(minimum, maximum, greedy, body, compiler, success, containsAnchor)
}

func (factory *v8GoPreprocessGraphFactory) NewTextNodeForClassRanges(class *v8ASTClassRanges, readBackward bool, success v8NodesNode, flags Flags) *v8NodesTextNode {
	node := (v8CompilerTextNodeFactory{}).NewTextNodeForClassRanges(class, readBackward, success, flags)
	factory.RegisterNode(node)
	return node
}

type v8GoSurrogateLookaroundBuilder struct {
	factory *v8GoPreprocessGraphFactory
	builder *v8ASTLookaroundBuilder
}

func (builder v8GoSurrogateLookaroundBuilder) OnMatchSuccess() v8NodesNode {
	return builder.builder.onMatchSuccessNode()
}

func (builder v8GoSurrogateLookaroundBuilder) ForMatch(compiler *v8CompilerCompiler, matchTrail v8NodesNode) v8NodesGuardedAlternative {
	node := builder.factory.lowerer.lookaroundBuilderForMatch(builder.builder, compiler, matchTrail)
	return v8NodesNewGuardedAlternative(node)
}

func (builder v8GoSurrogateLookaroundBuilder) RegisterNode(node v8NodesNode) {
	builder.factory.RegisterNode(node)
}

func v8BuildGoMatcherGraph(parsed v8ASTParserCompileData, flags Flags) (*v8GoMatcherGraph, error) {
	tree := parsed.Tree
	if tree == nil {
		return nil, fmt.Errorf("V8 parser produced no regexp AST")
	}
	zone := &v8CompilerZone{}
	// This port-side limit replaces the architecture-specific assembler limit;
	// capture-count validation itself remains in the V8-shaped parser.
	const maxGraphRegisters = 1 << 20
	compiler := v8CompilerNewCompiler(zone, parsed.CaptureCount, flags, false, maxGraphRegisters)
	lowerer := v8TonodeNewGraphLowerer(nil)
	factory := &v8GoPreprocessGraphFactory{lowerer: lowerer}
	request := &v8CompilerPreprocessRequest{tree: tree, containsAnchor: parsed.ContainsAnchor}
	root := v8CompilerPreprocessRegExp(compiler, request, false, factory)
	if request.compileError != v8CompilerNoCompilationError || compiler.IsRegExpTooBig() {
		return nil, fmt.Errorf("regular expression graph exceeds the compiler limit")
	}
	if root == nil {
		return nil, fmt.Errorf("V8 graph construction produced no root node")
	}
	if factory.body == nil {
		return nil, fmt.Errorf("V8 graph construction produced no search body")
	}
	if analysisError := v8CompilerAnalyzeRegExp(root, false, nil); analysisError != errorNone {
		return nil, fmt.Errorf("V8 regexp graph analysis failed: %s", errorString(analysisError))
	}
	searchMayMatchEmpty := v8TonodeTreeMinMatch(tree) == 0
	var firstSet, searchFirstSet []v8ASTCharacterRange
	hasFirstSet, hasSearchFirstSet := false, false
	if compiler.HasSearchPrefix() && !searchMayMatchEmpty {
		searchFirstSet, hasSearchFirstSet = v8CompilerComputeFirstCharacterSet(factory.body, compiler)
		if hasSearchFirstSet {
			searchFirstSet = v8GoFirstCharacterUnitRanges(searchFirstSet)
		}
		if v8GoFirstCharacterSetCoversUTF16(searchFirstSet) {
			hasSearchFirstSet = false
		}
	} else if !searchMayMatchEmpty {
		firstSet, hasFirstSet = v8CompilerComputeFirstCharacterSet(root, compiler)
		if hasFirstSet {
			firstSet = v8GoFirstCharacterUnitRanges(firstSet)
		}
		if v8GoFirstCharacterSetCoversUTF16(firstSet) {
			hasFirstSet = false
		}
	}
	// FIRST sets are computed in forward code-unit offsets. The parser graph
	// also contains the contents of lookarounds, which can be read backward and
	// must not constrain the candidate's first input unit.
	if v8GoGraphContainsLookaround(root) {
		hasFirstSet = false
		hasSearchFirstSet = false
	}
	return &v8GoMatcherGraph{
		root:                root,
		searchBody:          factory.body,
		compiler:            compiler,
		flags:               flags,
		captureCount:        parsed.CaptureCount,
		registerCount:       compiler.nextRegister,
		metadata:            parsed,
		firstSet:            firstSet,
		hasFirstSet:         hasFirstSet,
		searchFirstSet:      searchFirstSet,
		hasSearchFirstSet:   hasSearchFirstSet,
		searchMayMatchEmpty: searchMayMatchEmpty,
	}, nil
}

type v8GoGraphEmitter struct {
	graph                       *v8GoMatcherGraph
	name                        string
	entryNode                   v8NodesNode
	nodes                       []v8NodesNode
	nodeIDs                     map[v8NodesNode]int
	choiceCases                 []v8GoChoiceCase
	text                        strings.Builder
	rangeTables                 []string
	rangeTableNames             map[string]string
	stringSetTables             []string
	stringSetNames              map[*v8NodesFiniteStringSet]v8GoStringSetTableNames
	loopFrameNames              map[int]string
	stringSetDecoderEmitted     bool
	unsafeFirstCharacterFilters bool
}

type v8GoChoiceCase struct {
	key                     int
	target                  int
	guards                  *v8NodesZoneList[*v8NodesGuard]
	choiceLabel             int
	loopFrameName           string
	firstCharacterSet       []v8ASTCharacterRange
	hasFirstCharacterSet    bool
	firstCharacterMap       [2]uint64
	hasFirstCharacterFilter bool
}

func v8EmitGoMatcher(graph *v8GoMatcherGraph, name string) (string, error) {
	emitter := &v8GoGraphEmitter{graph: graph, name: name + "Into", entryNode: graph.root, nodeIDs: make(map[v8NodesNode]int), unsafeFirstCharacterFilters: v8GoGraphContainsLookaround(graph.root) || graph.searchMayMatchEmpty}
	if err := emitter.collectGraph(); err != nil {
		return "", err
	}
	if err := emitter.validateGraph(); err != nil {
		return "", err
	}
	if !emitter.graph.searchMayMatchEmpty {
		emitter.graph.searchPlan = v8GoBuildBoyerMooreSearchPlan(emitter.graph.searchBody, emitter.graph.compiler)
	} else {
		// A nullable expression can match at any position, including one whose
		// first input unit is outside the non-empty alternatives' first set.
		// Try the body at each legal start position rather than allowing V8's
		// search prefix or a first-character filter to skip an empty match.
		emitter.graph.searchPlan = v8GoBoyerMooreSearchPlan{
			enabled: true, advanceBy: 1,
			characterMap:   [2]uint64{^uint64(0), ^uint64(0)},
			unicodeAdvance: IsEitherUnicode(emitter.graph.flags),
		}
	}
	if emitter.graph.searchPlan.enabled {
		// The generated search loop tries each surviving candidate directly;
		// the V8 lazy wildcard prefix is therefore not part of this Go graph.
		emitter.entryNode = emitter.graph.searchBody
		emitter.nodes = nil
		emitter.nodeIDs = make(map[v8NodesNode]int)
		emitter.choiceCases = nil
		if err := emitter.collectGraph(); err != nil {
			return "", err
		}
		if err := emitter.validateGraph(); err != nil {
			return "", err
		}
	}
	emitter.prepareLoopFrames()
	emitter.writeHeader()
	for id, node := range emitter.nodes {
		if err := emitter.emitNode(id, node); err != nil {
			return "", err
		}
	}
	emitter.writeBacktrack()
	var source strings.Builder
	for _, table := range emitter.rangeTables {
		source.WriteString(table)
		source.WriteByte('\n')
	}
	for _, table := range emitter.stringSetTables {
		source.WriteString(table)
		source.WriteByte('\n')
	}
	captureSlots := (graph.captureCount + 1) * 2
	fmt.Fprintf(&source, "const %sCaptureSlots = %d\n", name, captureSlots)
	fmt.Fprintf(&source, "func %s(input []uint16, start int) (matchStart, matchEnd int, captures []int, ok bool) {\n", name)
	fmt.Fprintf(&source, "offsets := make([]int32, %sCaptureSlots)\n", name)
	fmt.Fprintf(&source, "matchStart, matchEnd, ok = %s(input, start, offsets)\n", emitter.name)
	source.WriteString("if !ok { return 0, 0, nil, false }; captures = make([]int, len(offsets)); for index, offset := range offsets { captures[index] = int(offset) }; return matchStart, matchEnd, captures, true\n}\n")
	source.WriteString(emitter.text.String())
	return source.String(), nil
}

func (emitter *v8GoGraphEmitter) prepareLoopFrames() {
	emitter.loopFrameNames = make(map[int]string)
	for id, node := range emitter.nodes {
		loop, ok := node.(*v8NodesLoopChoiceNode)
		if !ok || loop.BodyCanBeZeroLength() {
			continue
		}
		if loop.FixedLengthBodyIterationLength() > 0 || loop.SingleCodePointClassBody() {
			if loop.AtomicLoopKind() == v8NodesAtomicLoopAtEnd {
				emitter.loopFrameNames[id] = "rxLoopFrame" + strconv.Itoa(id)
			}
		}
	}
}

func (emitter *v8GoGraphEmitter) addNode(node v8NodesNode) {
	if node == nil {
		return
	}
	if _, exists := emitter.nodeIDs[node]; exists {
		return
	}
	emitter.nodeIDs[node] = len(emitter.nodes)
	emitter.nodes = append(emitter.nodes, node)
}

func (emitter *v8GoGraphEmitter) collectGraph() error {
	emitter.addNode(emitter.entryNode)
	for cursor := 0; cursor < len(emitter.nodes); cursor++ {
		node := emitter.nodes[cursor]
		switch value := node.(type) {
		case *v8NodesNegativeSubmatchSuccess:
			// This terminal owns the success/failure boundary of a negative
			// submatch and has no outgoing edge.
		case *v8NodesNegativeLookaroundChoiceNode:
			for i := 0; i < value.Alternatives().Len(); i++ {
				emitter.addNode(value.Alternatives().At(i).Node())
			}
		case *v8NodesLoopChoiceNode:
			for i := 0; i < value.Alternatives().Len(); i++ {
				emitter.addNode(value.Alternatives().At(i).Node())
			}
		case *v8NodesChoiceNode:
			for i := 0; i < value.Alternatives().Len(); i++ {
				emitter.addNode(value.Alternatives().At(i).Node())
			}
		case *v8NodesEndNode:
		case v8NodesNode:
			if seq := value.AsSeqNode(); seq != nil {
				emitter.addNode(seq.OnSuccess())
				continue
			}
			return fmt.Errorf("Go regex codegen does not yet support graph node %T", node)
		default:
			return fmt.Errorf("Go regex codegen does not yet support graph node %T", node)
		}
	}
	return nil
}

func (emitter *v8GoGraphEmitter) validateGraph() error {
	registerCount := max(emitter.graph.registerCount, 2)
	validRegister := func(register int) bool {
		return register >= 0 && register < registerCount
	}
	for _, node := range emitter.nodes {
		if seq := node.AsSeqNode(); seq != nil && seq.OnSuccess() == nil {
			return fmt.Errorf("V8 graph node %T has no success edge", node)
		}
		switch value := node.(type) {
		case *v8NodesNegativeSubmatchSuccess:
			if value.Action() != v8NodesNegativeSubmatchSuccessEndAction {
				return fmt.Errorf("V8 negative-submatch terminal has end action %d", value.Action())
			}
			if !validRegister(value.stackPointerRegister) || !validRegister(value.currentPositionRegister) || value.clearCaptureCount < 0 || (value.clearCaptureCount > 0 && (value.clearCaptureStart < 0 || value.clearCaptureStart+value.clearCaptureCount > registerCount)) {
				return fmt.Errorf("V8 negative-submatch terminal has invalid register range")
			}
		case *v8NodesNegativeLookaroundChoiceNode:
			if value.Alternatives() == nil || value.Alternatives().Len() != 2 {
				return fmt.Errorf("V8 negative-lookaround choice must have two alternatives")
			}
			if err := emitter.validateChoiceAlternatives(value.Alternatives(), validRegister); err != nil {
				return err
			}
		case *v8NodesEndNode:
			switch value.Action() {
			case v8NodesBacktrack, v8NodesAccept:
			default:
				return fmt.Errorf("Go regex codegen does not yet support end action %d", value.Action())
			}
		case *v8NodesTextNode:
			if value.finiteStringSet != nil {
				if value.ReadBackward() || len(value.finiteStringSet.nodes) == 0 || value.finiteStringSet.root < 0 || value.finiteStringSet.root >= len(value.finiteStringSet.nodes) {
					return fmt.Errorf("V8 finite string-set node has invalid automaton data")
				}
				break
			}
			if value.Elements() == nil {
				return fmt.Errorf("V8 text graph has no elements")
			}
			for i := 0; i < value.Elements().Len(); i++ {
				element := *value.Elements().At(i)
				switch element.textType {
				case v8ASTTextAtom:
					atom, ok := element.tree.(*v8ASTAtom)
					if !ok || atom == nil {
						return fmt.Errorf("V8 text graph contains an invalid atom payload")
					}
				case v8ASTTextClassRanges:
					class, ok := element.tree.(*v8ASTClassRanges)
					if !ok || class == nil {
						return fmt.Errorf("V8 text graph contains an invalid character-class payload")
					}
				default:
					return fmt.Errorf("Go regex codegen does not yet support text element %d", element.textType)
				}
			}
		case *v8NodesAssertionNode:
			switch value.AssertionType() {
			case v8NodesAssertAtStart, v8NodesAssertAtEnd, v8NodesAssertAtBoundary, v8NodesAssertAtNonBoundary, v8NodesAssertAfterNewline:
			default:
				return fmt.Errorf("Go regex codegen does not yet support assertion node %d", value.AssertionType())
			}
		case *v8NodesActionNode:
			switch value.ActionType() {
			case v8NodesStorePosition, v8NodesRestorePosition, v8NodesSetRegisterForLoop, v8NodesIncrementRegister:
				if !validRegister(value.RegisterFrom()) {
					return fmt.Errorf("V8 action node %d references register %d outside the graph", value.ActionType(), value.RegisterFrom())
				}
			case v8NodesClearCaptures:
				if value.RegisterFrom() < 0 || value.RegisterTo() < value.RegisterFrom() || value.RegisterTo() >= registerCount {
					return fmt.Errorf("V8 clear-captures action has invalid register range")
				}
			case v8NodesBeginPositiveSubmatch, v8NodesBeginNegativeSubmatch:
				if !validRegister(value.submatch.stackPointerRegister) || !validRegister(value.submatch.currentPositionRegister) {
					return fmt.Errorf("V8 submatch action has invalid state registers")
				}
				if value.ActionType() == v8NodesBeginPositiveSubmatch && value.SuccessNode() == nil {
					return fmt.Errorf("V8 positive-submatch action has no success node")
				}
			case v8NodesPositiveSubmatchSuccess:
				if !validRegister(value.submatch.stackPointerRegister) || !validRegister(value.submatch.currentPositionRegister) || value.submatch.clearRegisterCount < 0 || (value.submatch.clearRegisterCount > 0 && (value.submatch.clearRegisterFrom < 0 || value.submatch.clearRegisterFrom+value.submatch.clearRegisterCount > registerCount)) {
					return fmt.Errorf("V8 positive-submatch success has invalid register range")
				}
			case v8NodesEmptyMatchCheck:
				if !validRegister(value.emptyMatch.startRegister) || (value.emptyMatch.repetitionRegister != v8CompilerNoRegister && !validRegister(value.emptyMatch.repetitionRegister)) {
					return fmt.Errorf("V8 empty-match check has invalid register")
				}
			case v8NodesEatsAtLeastAction:
			default:
				return fmt.Errorf("Go regex codegen does not yet support action node %d", value.ActionType())
			}
		case *v8NodesLoopChoiceNode:
			if value.Alternatives() == nil || value.Alternatives().Len() != 2 {
				return fmt.Errorf("V8 loop choice must have two alternatives")
			}
			if err := emitter.validateChoiceAlternatives(value.Alternatives(), validRegister); err != nil {
				return err
			}
		case *v8NodesChoiceNode:
			if value.Alternatives() == nil || value.Alternatives().Len() == 0 {
				return fmt.Errorf("V8 choice node has no alternatives")
			}
			if err := emitter.validateChoiceAlternatives(value.Alternatives(), validRegister); err != nil {
				return err
			}
		case *v8NodesUnanchoredAdvanceNode:
		case *v8NodesBackReferenceNode:
			if !validRegister(value.StartRegister()) || !validRegister(value.EndRegister()) {
				return fmt.Errorf("V8 backreference has invalid capture registers")
			}
		default:
			return fmt.Errorf("Go regex codegen does not yet support graph node %T", node)
		}
	}
	return nil
}

func (emitter *v8GoGraphEmitter) validateChoiceAlternatives(alternatives *v8NodesZoneList[v8NodesGuardedAlternative], validRegister func(int) bool) error {
	for index := 0; index < alternatives.Len(); index++ {
		alternative := alternatives.At(index)
		if alternative.Node() == nil {
			return fmt.Errorf("V8 choice alternative %d has no target node", index)
		}
		guards := alternative.Guards()
		for guardIndex := 0; guards != nil && guardIndex < guards.Len(); guardIndex++ {
			guard := *guards.At(guardIndex)
			if !validRegister(guard.Register()) {
				return fmt.Errorf("V8 choice guard references register %d outside the graph", guard.Register())
			}
			switch guard.Relation() {
			case v8NodesGuardLT, v8NodesGuardGEQ:
			default:
				return fmt.Errorf("Go regex codegen does not yet support guard relation %d", guard.Relation())
			}
		}
	}
	return nil
}

func (emitter *v8GoGraphEmitter) writeHeader() {
	emitter.text.WriteString("func ")
	emitter.text.WriteString(emitter.name)
	emitter.text.WriteString("(input []uint16, start int, captures []int32) (matchStart, matchEnd int, ok bool) {\n")
	emitter.text.WriteString("const captureSlots = ")
	emitter.text.WriteString(strconv.Itoa((emitter.graph.captureCount + 1) * 2))
	emitter.text.WriteString("\n")
	emitter.text.WriteString("if start < 0 { start = 0 }; if start > len(input) { return 0, 0, false }; if len(captures) < captureSlots { return 0, 0, false }\n")
	if emitter.graph.hasFirstSet {
		if len(emitter.graph.firstSet) == 0 {
			emitter.text.WriteString("return 0, 0, false\n")
		} else {
			emitter.text.WriteString("if start < len(input) && !")
			emitter.text.WriteString(emitter.firstCharacterSetConditionAt(emitter.graph.firstSet, "start"))
			emitter.text.WriteString(" { return 0, 0, false }\n")
		}
	}
	emitter.text.WriteString("regs := make([]int, ")
	emitter.text.WriteString(strconv.Itoa(max(emitter.graph.registerCount, 2)))
	emitter.text.WriteString(")\nfor i := range regs { regs[i] = -1 }\n")
	emitter.text.WriteString("pos := start\nattemptStart := start\n")
	// Some graph shapes emit an immediate return before the matcher body, leaving
	// these registers write-only. Keep the generated Go valid; the compiler
	// removes these no-op reads.
	emitter.text.WriteString("_ = pos; _ = attemptStart\n")
	emitter.text.WriteString("const rxStackFrameSlots = ")
	emitter.text.WriteString(strconv.Itoa(max(emitter.graph.registerCount, 2) + 2))
	// regexp::Stack has a 64 MiB maximum and uses 4-byte slots on common V8 targets.
	emitter.text.WriteString("\nconst rxMaxStackSlots = 64 * 1024 * 1024 / 4\n")
	// V8 gives each execution a small backtracking stack which grows as needed.
	// A flat Go slice preserves that per-execution lifetime while avoiding a
	// separate capture-slice allocation for every saved alternative.
	emitter.text.WriteString("stack := make([]int32, 0, 128)\nvar choiceTarget, rxFrame int\n")
	loopFrameIDs := make([]int, 0, len(emitter.loopFrameNames))
	for id := range emitter.loopFrameNames {
		loopFrameIDs = append(loopFrameIDs, id)
	}
	sort.Ints(loopFrameIDs)
	for _, id := range loopFrameIDs {
		name := emitter.loopFrameNames[id]
		emitter.text.WriteString("var ")
		emitter.text.WriteString(name)
		emitter.text.WriteString(" = -1\n")
	}
	if emitter.hasCaptureCleanup() {
		// Keep these names distinct from runtime fields. The generated body is
		// rewritten to export runtime selectors, so generic names can otherwise
		// be mistaken for shared-runtime members.
		emitter.text.WriteString("type rxCaptureCleanup struct { stackDepth, captureFrom, captureTo int }; var rxCaptureCleanups []rxCaptureCleanup; rxCaptureCleanupMarks := make([]int, len(regs))\n")
	}
	emitter.emitCaseEqualityHelper()
	if emitter.graph.searchPlan.enabled {
		plan := emitter.graph.searchPlan
		emitter.text.WriteString("searchLoop: for candidate := start; candidate <= len(input); {\n")
		emitter.text.WriteString("for candidate+")
		emitter.text.WriteString(strconv.Itoa(plan.offset))
		emitter.text.WriteString(" < len(input) && !")
		emitter.text.WriteString(v8GoFirstCharacterMapConditionAt(plan.characterMap, "candidate+"+strconv.Itoa(plan.offset)))
		emitter.text.WriteString(" { candidate += ")
		emitter.text.WriteString(strconv.Itoa(plan.advanceBy))
		emitter.text.WriteString(" }\nif candidate > len(input) { break searchLoop }; pos = candidate; attemptStart = candidate\n")
		emitter.text.WriteString("for i := range regs { regs[i] = -1 }; stack = stack[:0]\n")
	}
	// Empty-capable bodies can accept at the first candidate without emitting
	// any failing node, leaving rxBacktrack otherwise unreachable. A runtime
	// invariant check keeps the shared label reachable to Go's compiler; start
	// and every search candidate are non-negative, so this branch is never taken.
	emitter.text.WriteString("if pos < 0 { goto rxBacktrack }\n")
	emitter.text.WriteString("goto rxNode")
	emitter.text.WriteString(strconv.Itoa(emitter.nodeIDs[emitter.entryNode]))
	emitter.text.WriteString("\n")
}

func (emitter *v8GoGraphEmitter) emitNode(id int, node v8NodesNode) error {
	emitter.text.WriteString("rxNode")
	emitter.text.WriteString(strconv.Itoa(id))
	emitter.text.WriteString(":\n")
	switch value := node.(type) {
	case *v8NodesNegativeSubmatchSuccess:
		emitter.emitNegativeSubmatchSuccess(value)
	case *v8NodesNegativeLookaroundChoiceNode:
		emitter.emitChoice(&value.v8NodesChoiceNode, id, false, false)
	case *v8NodesEndNode:
		if value.Action() == v8NodesBacktrack {
			emitter.text.WriteString("goto rxBacktrack\n")
			return nil
		}
		if value.Action() != v8NodesAccept {
			return fmt.Errorf("Go regex codegen does not yet support end action %d", value.Action())
		}
		emitter.text.WriteString("matchStart, matchEnd = regs[0], regs[1]\n")
		emitter.text.WriteString("if matchStart < 0 { matchStart = attemptStart }; if matchEnd < 0 { matchEnd = pos }\n")
		emitter.text.WriteString("for captureIndex := 0; captureIndex < captureSlots; captureIndex++ { captures[captureIndex] = int32(regs[captureIndex]) }\nreturn matchStart, matchEnd, true\n")
	case *v8NodesTextNode:
		if err := emitter.emitTextNode(value); err != nil {
			return err
		}
		emitter.gotoNode(value.OnSuccess())
	case *v8NodesAssertionNode:
		if err := emitter.emitAssertion(value); err != nil {
			return err
		}
		emitter.gotoNode(value.OnSuccess())
	case *v8NodesActionNode:
		if err := emitter.emitAction(value); err != nil {
			return err
		}
		emitter.gotoNode(value.OnSuccess())
	case *v8NodesLoopChoiceNode:
		emitter.emitChoice(&value.v8NodesChoiceNode, id, false, true)
	case *v8NodesChoiceNode:
		emitter.emitChoice(value, id, true, false)
	case *v8NodesUnanchoredAdvanceNode:
		emitter.text.WriteString("if pos >= len(input) { goto rxBacktrack }; ")
		if IsEitherUnicode(value.Flags()) {
			emitter.text.WriteString("if input[pos] >= 0xd800 && input[pos] <= 0xdbff && pos+1 < len(input) && input[pos+1] >= 0xdc00 && input[pos+1] <= 0xdfff { pos += 2 } else { pos++ }\n")
		} else {
			emitter.text.WriteString("pos++\n")
		}
		if emitter.graph.hasSearchFirstSet && !emitter.graph.searchPlan.enabled {
			emitter.text.WriteString("for pos < len(input) && !")
			emitter.text.WriteString(emitter.firstCharacterSetConditionAt(emitter.graph.searchFirstSet, "pos"))
			emitter.text.WriteString(" { ")
			if IsEitherUnicode(value.Flags()) {
				emitter.text.WriteString("if input[pos] >= 0xd800 && input[pos] <= 0xdbff && pos+1 < len(input) && input[pos+1] >= 0xdc00 && input[pos+1] <= 0xdfff { pos += 2 } else { pos++ }")
			} else {
				emitter.text.WriteString("pos++")
			}
			emitter.text.WriteString(" }; if pos >= len(input) { goto rxBacktrack }\n")
		}
		emitter.gotoNode(value.OnSuccess())
	case *v8NodesBackReferenceNode:
		emitter.emitBackReference(value)
		emitter.gotoNode(value.OnSuccess())
	default:
		return fmt.Errorf("Go regex codegen does not yet support graph node %T", node)
	}
	return nil
}

func (emitter *v8GoGraphEmitter) emitAtom(atom *v8ASTAtom) error {
	if atom == nil {
		return fmt.Errorf("nil atom in V8 text graph")
	}
	emitter.text.WriteString("if pos < 0 || pos+")
	emitter.text.WriteString(strconv.Itoa(len(atom.data)))
	emitter.text.WriteString(" > len(input) { goto rxBacktrack }\n")
	for offset, unit := range atom.data {
		emitter.text.WriteString("if input[pos+")
		emitter.text.WriteString(strconv.Itoa(offset))
		emitter.text.WriteString("] != ")
		emitter.text.WriteString(strconv.Itoa(int(unit)))
		emitter.text.WriteString(" { goto rxBacktrack }\n")
	}
	emitter.text.WriteString("pos += ")
	emitter.text.WriteString(strconv.Itoa(len(atom.data)))
	emitter.text.WriteString("\n")
	return nil
}

func (emitter *v8GoGraphEmitter) emitClass(class *v8ASTClassRanges) error {
	if class == nil {
		return fmt.Errorf("nil character class in V8 text graph")
	}
	condition, err := emitter.classConditionAt(class, "pos")
	if err != nil {
		return err
	}
	emitter.text.WriteString("if pos < 0 || pos >= len(input) || !(")
	emitter.text.WriteString(condition)
	emitter.text.WriteString(") { goto rxBacktrack }; pos++\n")
	return nil
}

func (emitter *v8GoGraphEmitter) classConditionAt(class *v8ASTClassRanges, index string) (string, error) {
	character := "input[" + index + "]"
	set := class.set
	condition := "false"
	if set.standardSet != nil {
		whitespace := "(" + character + " == 9 || " + character + " == 10 || " + character + " == 11 || " + character + " == 12 || " + character + " == 13 || " + character + " == 32 || " + character + " == 160 || " + character + " == 5760 || (" + character + " >= 8192 && " + character + " <= 8202) || " + character + " == 8232 || " + character + " == 8233 || " + character + " == 8239 || " + character + " == 8287 || " + character + " == 12288 || " + character + " == 65279)"
		word := "((" + character + " >= '0' && " + character + " <= '9') || (" + character + " >= 'A' && " + character + " <= 'Z') || " + character + " == '_' || (" + character + " >= 'a' && " + character + " <= 'z'))"
		digit := "(" + character + " >= '0' && " + character + " <= '9')"
		lineTerminator := "(" + character + " == 10 || " + character + " == 13 || " + character + " == 0x2028 || " + character + " == 0x2029)"
		switch *set.standardSet {
		case v8ASTWhitespace:
			condition = whitespace
		case v8ASTNotWhitespace:
			condition = "!(" + whitespace + ")"
		case v8ASTWord:
			condition = word
		case v8ASTNotWord:
			condition = "!(" + word + ")"
		case v8ASTDigit:
			condition = digit
		case v8ASTNotDigit:
			condition = "!(" + digit + ")"
		case v8ASTLineTerminator:
			condition = lineTerminator
		case v8ASTNotLineTerminator:
			condition = "!(" + lineTerminator + ")"
		case v8ASTEverything:
			condition = "true"
		default:
			return "", fmt.Errorf("Go regex codegen does not yet support standard character class %q", rune(*set.standardSet))
		}
	} else if len(set.ranges) != 0 {
		ranges := v8TonodeCanonicalizeRanges(set.ranges)
		if ranges[len(ranges)-1].to > 0xffff {
			// The class graph keeps scalar ranges above the BMP even when the
			// legacy matcher cursor is a UTF-16 code unit. Widen before emitting
			// constants and bitmap indexes so Go does not reject 0x10ffff as a
			// uint16 constant.
			character = "uint32(" + character + ")"
		}
		condition = emitter.rangeCondition(character, ranges)
	}
	if class.flags&v8ASTClassNegated != 0 {
		condition = "!(" + condition + ")"
	}
	return condition, nil
}

// rangeCondition keeps small character classes inline and stores large ones as
// a compact UTF-16 bitmap. Expanding thousands of Unicode intervals into a Go
// boolean expression makes `go build` parse and optimize a huge syntax tree;
// a read-only bitmap is both cheaper to compile and constant-time at runtime.
func (emitter *v8GoGraphEmitter) rangeCondition(character string, ranges []v8ASTCharacterRange) string {
	if len(ranges) <= 4 {
		parts := make([]string, 0, len(ranges))
		for _, current := range ranges {
			parts = append(parts, fmt.Sprintf("(%s >= %d && %s <= %d)", character, current.from, character, current.to))
		}
		return "(" + strings.Join(parts, " || ") + ")"
	}
	if len(ranges) <= 64 {
		return v8GoClassRangeTreeCondition(character, ranges)
	}
	key := v8GoRangeTableKey(ranges)
	if emitter.rangeTableNames == nil {
		emitter.rangeTableNames = make(map[string]string)
	}
	tableName, exists := emitter.rangeTableNames[key]
	if !exists {
		tableName = fmt.Sprintf("%sRangeBits%d", emitter.name, len(emitter.rangeTables))
		bitCount := 1 << 10 // 65,536 UTF-16 code units.
		if ranges[len(ranges)-1].to > 0xffff {
			bitCount = (0x110000 + 63) / 64 // Unicode scalar values.
		}
		bits := make([]uint64, bitCount)
		for _, current := range ranges {
			for value := current.from; value <= current.to; value++ {
				unit := uint32(value)
				bits[unit>>6] |= uint64(1) << (unit & 63)
			}
		}
		var declaration strings.Builder
		fmt.Fprintf(&declaration, "var %s = [%d]uint64{\n", tableName, bitCount)
		for index, value := range bits {
			if value != 0 {
				fmt.Fprintf(&declaration, "%d: 0x%x,\n", index, value)
			}
		}
		declaration.WriteString("}\n")
		emitter.rangeTables = append(emitter.rangeTables, declaration.String())
		emitter.rangeTableNames[key] = tableName
	}
	return fmt.Sprintf("((%s[%s>>6] & (uint64(1) << (%s&63))) != 0)", tableName, character, character)
}

func v8GoRangeTableKey(ranges []v8ASTCharacterRange) string {
	var key strings.Builder
	for _, current := range ranges {
		key.WriteString(strconv.Itoa(int(current.from)))
		key.WriteByte(':')
		key.WriteString(strconv.Itoa(int(current.to)))
		key.WriteByte(';')
	}
	return key.String()
}

// v8GoClassRangeTreeCondition emits the same ordered-range membership test as
// V8's IsCharacterInRangeArray, but encodes the binary search as a balanced
// short-circuit expression. Small classes stay on the straight-line path
// above; large classes avoid testing every interval for every input unit.
func v8GoClassRangeTreeCondition(character string, ranges []v8ASTCharacterRange) string {
	if len(ranges) == 0 {
		return "false"
	}
	middle := len(ranges) / 2
	current := ranges[middle]
	parts := []string{fmt.Sprintf("(%s >= %d && %s <= %d)", character, current.from, character, current.to)}
	if middle != 0 {
		parts = append(parts, fmt.Sprintf("(%s < %d && %s)", character, current.from, v8GoClassRangeTreeCondition(character, ranges[:middle])))
	}
	if middle+1 < len(ranges) {
		parts = append(parts, fmt.Sprintf("(%s > %d && %s)", character, current.to, v8GoClassRangeTreeCondition(character, ranges[middle+1:])))
	}
	return "(" + strings.Join(parts, " || ") + ")"
}

func (emitter *v8GoGraphEmitter) emitAssertion(node *v8NodesAssertionNode) error {
	condition := "false"
	switch node.AssertionType() {
	case v8NodesAssertAtStart:
		// Multiline anchors also match only between complete CRLF sequences.
		condition = "pos == 0"
		if node.Flags().Multiline {
			condition = "(pos == 0 || (pos > 0 && (input[pos-1] == 10 || input[pos-1] == 13 || input[pos-1] == 0x2028 || input[pos-1] == 0x2029) && !(input[pos-1] == 13 && pos < len(input) && input[pos] == 10)))"
		}
	case v8NodesAssertAtEnd:
		// Multiline $ matches before a line terminator, but not between CRLF.
		condition = "pos == len(input)"
		if node.Flags().Multiline {
			condition = "(pos == len(input) || (pos < len(input) && (input[pos] == 10 || input[pos] == 13 || input[pos] == 0x2028 || input[pos] == 0x2029) && !(input[pos] == 10 && pos > 0 && input[pos-1] == 13)))"
		}
	case v8NodesAssertAtBoundary, v8NodesAssertAtNonBoundary:
		left := "(pos > 0 && " + v8GoWordConditionAt("pos-1", node.Flags()) + ")"
		right := "(pos < len(input) && " + v8GoWordConditionAt("pos", node.Flags()) + ")"
		condition = "(" + left + " != " + right + ")"
		if node.AssertionType() == v8NodesAssertAtNonBoundary {
			condition = "!(" + condition + ")"
		}
	case v8NodesAssertAfterNewline:
		condition = "(pos == 0 || input[pos-1] == 10 || input[pos-1] == 13 || input[pos-1] == 0x2028 || input[pos-1] == 0x2029)"
	default:
		return fmt.Errorf("Go regex codegen does not yet support assertion node %d", node.AssertionType())
	}
	emitter.text.WriteString("if !(")
	emitter.text.WriteString(condition)
	emitter.text.WriteString(") { goto rxBacktrack }\n")
	return nil
}

func (emitter *v8GoGraphEmitter) emitAction(node *v8NodesActionNode) error {
	switch node.ActionType() {
	case v8NodesStorePosition:
		emitter.text.WriteString("regs[")
		emitter.text.WriteString(strconv.Itoa(node.simple.registerFrom))
		emitter.text.WriteString("] = pos\n")
	case v8NodesRestorePosition:
		emitter.text.WriteString("if regs[")
		emitter.text.WriteString(strconv.Itoa(node.simple.registerFrom))
		emitter.text.WriteString("] < 0 { goto rxBacktrack }; pos = regs[")
		emitter.text.WriteString(strconv.Itoa(node.simple.registerFrom))
		emitter.text.WriteString("]\n")
	case v8NodesSetRegisterForLoop:
		emitter.text.WriteString("regs[")
		emitter.text.WriteString(strconv.Itoa(node.simple.registerFrom))
		emitter.text.WriteString("] = ")
		emitter.text.WriteString(strconv.Itoa(node.simple.value))
		emitter.text.WriteString("\n")
	case v8NodesIncrementRegister:
		emitter.text.WriteString("regs[")
		emitter.text.WriteString(strconv.Itoa(node.simple.registerFrom))
		emitter.text.WriteString("]++\n")
	case v8NodesClearCaptures:
		for register := node.simple.registerFrom; register <= node.simple.registerTo; register++ {
			emitter.text.WriteString("regs[")
			emitter.text.WriteString(strconv.Itoa(register))
			emitter.text.WriteString("] = -1\n")
		}
	case v8NodesBeginPositiveSubmatch, v8NodesBeginNegativeSubmatch:
		emitter.emitBeginSubmatch(node)
	case v8NodesPositiveSubmatchSuccess:
		emitter.emitPositiveSubmatchSuccess(node)
	case v8NodesEmptyMatchCheck:
		emitter.emitEmptyMatchCheck(node)
	case v8NodesEatsAtLeastAction:
	default:
		return fmt.Errorf("Go regex codegen does not yet support action node %d", node.ActionType())
	}
	return nil
}

// hasCaptureCleanup keeps the regular expression backtrack path free of
// lookaround bookkeeping unless a positive lookaround actually captures.
func (emitter *v8GoGraphEmitter) hasCaptureCleanup() bool {
	for _, node := range emitter.nodes {
		if action, ok := node.(*v8NodesActionNode); ok && action.ActionType() == v8NodesPositiveSubmatchSuccess && action.submatch.clearRegisterCount > 0 {
			return true
		}
	}
	return false
}

func (emitter *v8GoGraphEmitter) emitChoice(node *v8NodesChoiceNode, id int, allowDispatch, isLoopChoice bool) {
	alternatives := node.Alternatives()
	var precomputedFirstSets [][]v8ASTCharacterRange
	if allowDispatch && !emitter.unsafeFirstCharacterFilters {
		var dispatched bool
		precomputedFirstSets, dispatched = emitter.emitChoiceDispatch(node)
		if dispatched {
			return
		}
	}
	for index := alternatives.Len() - 1; index >= 1; index-- {
		alternative := alternatives.At(index)
		target := emitter.nodeIDs[alternative.Node()]
		key := len(emitter.choiceCases) + 1
		var firstCharacterSet []v8ASTCharacterRange
		hasFirstCharacterSet := false
		if emitter.unsafeFirstCharacterFilters {
			// The same forward-only quick-check tables used by the dispatch
			// optimizer are unsound for alternatives reached inside lookarounds.
		} else if len(precomputedFirstSets) == alternatives.Len() {
			firstCharacterSet = precomputedFirstSets[index]
			if IsEitherUnicode(node.Flags()) {
				firstCharacterSet = v8GoFirstCharacterUnitRanges(firstCharacterSet)
			}
			hasFirstCharacterSet = true
		} else {
			firstCharacterSet, hasFirstCharacterSet = v8CompilerComputeFirstCharacterSet(alternative.Node(), emitter.graph.compiler)
		}
		if hasFirstCharacterSet && v8GoFirstCharacterSetCoversUTF16(firstCharacterSet) {
			hasFirstCharacterSet = false
		}
		firstCharacterMap, filterFirstCharacter := [2]uint64{}, false
		if !hasFirstCharacterSet && !emitter.unsafeFirstCharacterFilters {
			firstCharacterMap, filterFirstCharacter = v8GoFirstCharacterMap(alternative.Node(), emitter.graph.compiler)
		}
		loopFrameName := ""
		if isLoopChoice {
			loopFrameName = emitter.loopFrameNames[id]
			// Loop continuation alternatives may be zero-width assertions such
			// as `$`. Their empty first-character set does not mean they are
			// impossible, so loop choices must keep the branch unconditionally.
			hasFirstCharacterSet = false
			filterFirstCharacter = false
		}
		emitter.choiceCases = append(emitter.choiceCases, v8GoChoiceCase{
			key: key, target: target, guards: alternative.Guards(), choiceLabel: id,
			loopFrameName:     loopFrameName,
			firstCharacterSet: firstCharacterSet, hasFirstCharacterSet: hasFirstCharacterSet,
			firstCharacterMap: firstCharacterMap, hasFirstCharacterFilter: filterFirstCharacter,
		})
		stackCapName := "rxStackCap" + strconv.Itoa(key)
		grownStackName := "rxGrownStack" + strconv.Itoa(key)
		if hasFirstCharacterSet && len(firstCharacterSet) == 0 {
			// The source walk proved that this branch cannot consume a first
			// code unit on any successful path, so it cannot match.
			continue
		} else if hasFirstCharacterSet {
			emitter.text.WriteString("if pos >= len(input) || ")
			emitter.text.WriteString(emitter.firstCharacterSetConditionAt(firstCharacterSet, "pos"))
			emitter.text.WriteString(" { if len(stack)+rxStackFrameSlots > rxMaxStackSlots { return -1, -1, false }; ")
		} else if filterFirstCharacter {
			emitter.text.WriteString("if pos >= len(input) || ")
			emitter.text.WriteString(v8GoFirstCharacterMapContains(firstCharacterMap))
			emitter.text.WriteString(" { if len(stack)+rxStackFrameSlots > rxMaxStackSlots { return -1, -1, false }; ")
		} else {
			emitter.text.WriteString("if len(stack)+rxStackFrameSlots > rxMaxStackSlots { return -1, -1, false }; ")
		}
		emitter.text.WriteString("if len(stack)+rxStackFrameSlots > cap(stack) { ")
		emitter.text.WriteString(stackCapName)
		emitter.text.WriteString(" := cap(stack)*2; if ")
		emitter.text.WriteString(stackCapName)
		emitter.text.WriteString(" < len(stack)+rxStackFrameSlots { ")
		emitter.text.WriteString(stackCapName)
		emitter.text.WriteString(" = len(stack)+rxStackFrameSlots }; if ")
		emitter.text.WriteString(stackCapName)
		emitter.text.WriteString(" > rxMaxStackSlots { ")
		emitter.text.WriteString(stackCapName)
		emitter.text.WriteString(" = rxMaxStackSlots }; ")
		emitter.text.WriteString(grownStackName)
		emitter.text.WriteString(" := make([]int32, len(stack), ")
		emitter.text.WriteString(stackCapName)
		emitter.text.WriteString("); copy(")
		emitter.text.WriteString(grownStackName)
		emitter.text.WriteString(", stack); stack = ")
		emitter.text.WriteString(grownStackName)
		emitter.text.WriteString(" }; ")
		if loopFrameName != "" {
			emitter.text.WriteString("if ")
			emitter.text.WriteString(loopFrameName)
			emitter.text.WriteString(" < 0 { stack = append(stack, int32(")
			emitter.text.WriteString(strconv.Itoa(key))
			emitter.text.WriteString("), int32(pos)")
			for register := 0; register < max(emitter.graph.registerCount, 2); register++ {
				emitter.text.WriteString(", int32(regs[")
				emitter.text.WriteString(strconv.Itoa(register))
				emitter.text.WriteString("]) ")
			}
			emitter.text.WriteString("); ")
			emitter.text.WriteString(loopFrameName)
			emitter.text.WriteString(" = len(stack)-rxStackFrameSlots } else { stack = stack[:")
			emitter.text.WriteString(loopFrameName)
			emitter.text.WriteString("+rxStackFrameSlots]; stack[")
			emitter.text.WriteString(loopFrameName)
			emitter.text.WriteString("] = int32(")
			emitter.text.WriteString(strconv.Itoa(key))
			emitter.text.WriteString("); stack[")
			emitter.text.WriteString(loopFrameName)
			emitter.text.WriteString("+1] = int32(pos)")
			for register := 0; register < max(emitter.graph.registerCount, 2); register++ {
				emitter.text.WriteString("; stack[")
				emitter.text.WriteString(loopFrameName)
				emitter.text.WriteString("+")
				emitter.text.WriteString(strconv.Itoa(register + 2))
				emitter.text.WriteString("] = int32(regs[")
				emitter.text.WriteString(strconv.Itoa(register))
				emitter.text.WriteString("])")
			}
			emitter.text.WriteString(" }\n")
		} else {
			emitter.text.WriteString("stack = append(stack, int32(")
			emitter.text.WriteString(strconv.Itoa(key))
			emitter.text.WriteString("), int32(pos)")
			for register := 0; register < max(emitter.graph.registerCount, 2); register++ {
				emitter.text.WriteString(", int32(regs[")
				emitter.text.WriteString(strconv.Itoa(register))
				emitter.text.WriteString("]) ")
			}
			emitter.text.WriteString(")")
		}
		if hasFirstCharacterSet && len(firstCharacterSet) != 0 || filterFirstCharacter {
			emitter.text.WriteString(" }")
		}
		emitter.text.WriteString("\n")
	}
	if alternatives.Len() == 0 {
		emitter.text.WriteString("goto rxBacktrack\n")
		return
	}
	first := alternatives.At(0)
	if condition := v8GoGuardCondition(first.Guards()); condition != "true" {
		emitter.text.WriteString("if !(")
		emitter.text.WriteString(condition)
		emitter.text.WriteString(") { goto rxBacktrack }\n")
	}
	emitter.gotoNode(first.Node())
}

type v8GoChoiceDispatchSegment struct {
	from, to uint32
	target   int // -1 means that no alternative can start in this interval.
}

// emitChoiceDispatch replaces a linear chain of ordered first-character
// probes with a balanced decision tree only when the alternatives are
// provably disjoint at their first character and the tree has a lower
// worst-case comparison count. Since only one branch can then match, no
// backtracking frame is needed for the skipped alternatives.
func (emitter *v8GoGraphEmitter) emitChoiceDispatch(node *v8NodesChoiceNode) ([][]v8ASTCharacterRange, bool) {
	if node == nil || node.Alternatives() == nil || node.Alternatives().Len() < 3 || node.Flags().IgnoreCase {
		return nil, false
	}
	linearChecks := node.Alternatives().Len() - 1
	unicode := IsEitherUnicode(node.Flags())
	setupChecks := 1
	if unicode {
		setupChecks = 6 // bounds check plus worst-case surrogate-pair decoding predicates
	}
	// Even perfectly packed disjoint sets need at least one interval per
	// alternative. Avoid first-set analysis when that lower bound already makes
	// the tree no cheaper than the existing chain.
	if bits.Len(uint(node.Alternatives().Len()-1))+setupChecks >= linearChecks {
		return nil, false
	}
	maxCharacter := uint32(0xffff)
	if unicode {
		maxCharacter = 0x10ffff
	}
	intervals := make([]v8GoChoiceDispatchSegment, 0, node.Alternatives().Len())
	firstSets := make([][]v8ASTCharacterRange, 0, node.Alternatives().Len())
	for index := 0; index < node.Alternatives().Len(); index++ {
		alternative := node.Alternatives().At(index)
		if alternative.Node() == nil || alternative.Guards() != nil && alternative.Guards().Len() != 0 {
			return nil, false
		}
		firstSet, known := v8CompilerComputeFirstCodePointSet(alternative.Node(), emitter.graph.compiler)
		if !known || len(firstSet) == 0 {
			return nil, false
		}
		firstSets = append(firstSets, firstSet)
		target, exists := emitter.nodeIDs[alternative.Node()]
		if !exists {
			return nil, false
		}
		for _, current := range firstSet {
			if current.from < 0 || current.to < current.from || uint32(current.to) > maxCharacter {
				return nil, false
			}
			intervals = append(intervals, v8GoChoiceDispatchSegment{
				from: uint32(current.from), to: uint32(current.to), target: target,
			})
		}
	}
	sort.Slice(intervals, func(left, right int) bool {
		if intervals[left].from != intervals[right].from {
			return intervals[left].from < intervals[right].from
		}
		return intervals[left].to < intervals[right].to
	})
	for index := 1; index < len(intervals); index++ {
		if intervals[index-1].to >= intervals[index].from {
			// A character accepted by two branches must retain ordered
			// backtracking semantics, so overlapping sets use the old path.
			return firstSets, false
		}
	}

	segments := make([]v8GoChoiceDispatchSegment, 0, len(intervals)*2+1)
	appendSegment := func(segment v8GoChoiceDispatchSegment) {
		if len(segments) != 0 {
			last := &segments[len(segments)-1]
			if last.target == segment.target && last.to+1 == segment.from {
				last.to = segment.to
				return
			}
		}
		segments = append(segments, segment)
	}
	cursor := uint32(0)
	for _, interval := range intervals {
		if cursor < interval.from {
			appendSegment(v8GoChoiceDispatchSegment{from: cursor, to: interval.from - 1, target: -1})
		}
		appendSegment(interval)
		cursor = interval.to + 1
	}
	if cursor <= maxCharacter {
		appendSegment(v8GoChoiceDispatchSegment{from: cursor, to: maxCharacter, target: -1})
	}
	if len(segments) < 2 {
		return nil, false
	}

	// The linear emitter checks every later alternative. A balanced tree over
	// the exhaustive intervals needs ceil(log2(segment count)) comparisons;
	// include a conservative setup charge for checking/decoding the input unit.
	treeChecks := bits.Len(uint(len(segments) - 1))
	if treeChecks+setupChecks >= linearChecks {
		return firstSets, false
	}

	emitter.text.WriteString("{\nif pos >= len(input) { goto rxBacktrack }; rxDispatchChar := uint32(input[pos])\n")
	if unicode {
		emitter.text.WriteString("if rxDispatchChar >= 0xd800 && rxDispatchChar <= 0xdbff && pos+1 < len(input) && input[pos+1] >= 0xdc00 && input[pos+1] <= 0xdfff { rxDispatchChar = 0x10000 + ((rxDispatchChar-0xd800)<<10) + (uint32(input[pos+1])-0xdc00) }\n")
	}
	emitter.emitChoiceDispatchTree(segments, 0, len(segments)-1)
	emitter.text.WriteString("}\n")
	return nil, true
}

func (emitter *v8GoGraphEmitter) emitChoiceDispatchTree(segments []v8GoChoiceDispatchSegment, first, last int) {
	if first == last {
		if segments[first].target < 0 {
			emitter.text.WriteString("goto rxBacktrack\n")
			return
		}
		emitter.text.WriteString("goto rxNode")
		emitter.text.WriteString(strconv.Itoa(segments[first].target))
		emitter.text.WriteByte('\n')
		return
	}
	middle := first + (last-first)/2
	emitter.text.WriteString("if rxDispatchChar <= ")
	emitter.text.WriteString(strconv.FormatUint(uint64(segments[middle].to), 10))
	emitter.text.WriteString(" {\n")
	emitter.emitChoiceDispatchTree(segments, first, middle)
	emitter.text.WriteString("} else {\n")
	emitter.emitChoiceDispatchTree(segments, middle+1, last)
	emitter.text.WriteString("}\n")
}

func v8GoFirstCharacterSetCoversUTF16(ranges []v8ASTCharacterRange) bool {
	return len(ranges) == 1 && ranges[0].from <= 0 && ranges[0].to >= 0xffff
}

// V8 graph ranges can denote Unicode scalar values above the BMP. The Go
// matcher cursor is a UTF-16 code-unit index, so a necessary first-code-unit
// filter for a supplementary value is its lead-surrogate range.
func v8GoFirstCharacterUnitRanges(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	units := make([]v8ASTCharacterRange, 0, len(ranges)+1)
	for _, current := range ranges {
		if current.from <= 0xffff {
			end := min(current.to, rune(0xffff))
			if current.from <= end {
				units = append(units, v8ASTCharacterRange{from: current.from, to: end})
			}
		}
		if current.to >= 0x10000 {
			start := max(current.from, rune(0x10000))
			leadStart := rune(0xd800) + ((start - 0x10000) >> 10)
			leadEnd := rune(0xd800) + ((current.to - 0x10000) >> 10)
			units = append(units, v8ASTCharacterRange{from: leadStart, to: leadEnd})
		}
	}
	return v8TonodeCanonicalizeRanges(units)
}

func v8GoFirstCharacterSetCondition(ranges []v8ASTCharacterRange) string {
	return v8GoFirstCharacterSetConditionAt(ranges, "pos")
}

func v8GoFirstCharacterSetConditionAt(ranges []v8ASTCharacterRange, index string) string {
	if len(ranges) == 0 {
		return "false"
	}
	ranges = v8TonodeCanonicalizeRanges(ranges)
	character := "uint32(input[" + index + "])"
	if len(ranges) > 4 {
		return v8GoClassRangeTreeCondition(character, ranges)
	}
	conditions := make([]string, 0, len(ranges))
	for _, current := range ranges {
		if current.from == current.to {
			conditions = append(conditions, character+" == "+strconv.FormatInt(int64(current.from), 10))
		} else {
			conditions = append(conditions, "("+character+" >= "+strconv.FormatInt(int64(current.from), 10)+" && "+character+" <= "+strconv.FormatInt(int64(current.to), 10)+")")
		}
	}
	return "(" + strings.Join(conditions, " || ") + ")"
}

func (emitter *v8GoGraphEmitter) firstCharacterSetConditionAt(ranges []v8ASTCharacterRange, index string) string {
	if len(ranges) == 0 {
		return "false"
	}
	ranges = v8TonodeCanonicalizeRanges(ranges)
	character := "uint32(input[" + index + "])"
	if len(ranges) > 4 {
		return emitter.rangeCondition(character, ranges)
	}
	return v8GoFirstCharacterSetConditionAt(ranges, index)
}

// v8GoFirstCharacterMap returns a conservative first-code-unit map for one
// alternative. A set bit means the branch may match that code unit; an absent
// bit proves that the branch cannot succeed at the current input position.
func v8GoFirstCharacterMap(node v8NodesNode, compiler *v8CompilerCompiler) ([2]uint64, bool) {
	if node == nil || compiler == nil {
		return [2]uint64{}, false
	}
	lookahead := v8NodesNewBoyerMooreLookahead(1, compiler)
	lookahead.SetCachesNodeInfo(false)
	v8CompilerFillInBMInfo(node, 0, v8NodesRecursionBudget, lookahead, false)
	position := lookahead.At(0)
	if position.MapCount() == 0 || position.MapCount() == v8NodesBoyerMooreMapSize {
		return [2]uint64{}, false
	}
	return position.RawBitset(), true
}

func v8GoFirstCharacterMapContains(characterMap [2]uint64) string {
	return v8GoFirstCharacterMapConditionAt(characterMap, "pos")
}

func v8GoFirstCharacterMapConditionAt(characterMap [2]uint64, index string) string {
	return "((uint16(input[" + index + "])&64 == 0 && (uint64(" + strconv.FormatUint(characterMap[0], 10) + ")>>uint(input[" + index + "]&63))&1 != 0) || (uint16(input[" + index + "])&64 != 0 && (uint64(" + strconv.FormatUint(characterMap[1], 10) + ")>>uint(input[" + index + "]&63))&1 != 0))"
}

func v8GoGuardCondition(guards *v8NodesZoneList[*v8NodesGuard]) string {
	if guards == nil || guards.Len() == 0 {
		return "true"
	}
	parts := make([]string, 0, guards.Len())
	for index := 0; index < guards.Len(); index++ {
		guard := *guards.At(index)
		operator := "<"
		if guard.Relation() == v8NodesGuardGEQ {
			operator = ">="
		}
		parts = append(parts, fmt.Sprintf("regs[%d] %s %d", guard.Register(), operator, guard.Value()))
	}
	return strings.Join(parts, " && ")
}

func (emitter *v8GoGraphEmitter) emitBackReference(node *v8NodesBackReferenceNode) {
	emitter.emitV8BackReference(node)
}

func (emitter *v8GoGraphEmitter) gotoNode(node v8NodesNode) {
	emitter.text.WriteString("goto rxNode")
	emitter.text.WriteString(strconv.Itoa(emitter.nodeIDs[node]))
	emitter.text.WriteString("\n")
}

func (emitter *v8GoGraphEmitter) writeBacktrack() {
	emitter.text.WriteString("rxBacktrack:\n")
	if emitter.hasCaptureCleanup() {
		emitter.text.WriteString("for len(rxCaptureCleanups) > 0 && len(stack) <= rxCaptureCleanups[len(rxCaptureCleanups)-1].stackDepth { rxCleanup := rxCaptureCleanups[len(rxCaptureCleanups)-1]; rxCaptureCleanups = rxCaptureCleanups[:len(rxCaptureCleanups)-1]; for rxCleanupRegister := rxCleanup.captureFrom; rxCleanupRegister < rxCleanup.captureTo; rxCleanupRegister++ { regs[rxCleanupRegister] = -1 } }\n")
	}
	emitter.text.WriteString("if len(stack) < rxStackFrameSlots {")
	if emitter.graph.searchPlan.enabled {
		if emitter.graph.searchPlan.unicodeAdvance {
			emitter.text.WriteString(" candidate++; if candidate < len(input) && candidate > 0 && input[candidate-1] >= 0xd800 && input[candidate-1] <= 0xdbff && input[candidate] >= 0xdc00 && input[candidate] <= 0xdfff { candidate++ }; continue searchLoop")
		} else {
			emitter.text.WriteString(" candidate++; continue searchLoop")
		}
	} else {
		emitter.text.WriteString(" return 0, 0, false")
	}
	emitter.text.WriteString(" }\n")
	emitter.text.WriteString("rxFrame = len(stack)-rxStackFrameSlots; choiceTarget = int(stack[rxFrame]); pos = int(stack[rxFrame+1]); for rxIndex := range regs { regs[rxIndex] = int(stack[rxFrame+2+rxIndex]) }; stack = stack[:rxFrame]\n")
	emitter.text.WriteString("switch choiceTarget {\n")
	for _, branch := range emitter.choiceCases {
		emitter.text.WriteString("case ")
		emitter.text.WriteString(strconv.Itoa(branch.key))
		emitter.text.WriteString(":\n")
		if branch.loopFrameName != "" {
			emitter.text.WriteString(branch.loopFrameName)
			emitter.text.WriteString(" = -1\n")
		}
		if condition := v8GoGuardCondition(branch.guards); condition != "true" {
			emitter.text.WriteString("if !(")
			emitter.text.WriteString(condition)
			emitter.text.WriteString(") { goto rxBacktrack }\n")
		}
		emitter.text.WriteString("goto rxNode")
		emitter.text.WriteString(strconv.Itoa(branch.target))
		emitter.text.WriteString("\n")
	}
	emitter.text.WriteString("default: return 0, 0, false\n}\n")
	if emitter.graph.searchPlan.enabled {
		emitter.text.WriteString("}\nreturn 0, 0, false\n")
	}
	emitter.text.WriteString("}\n")
}
