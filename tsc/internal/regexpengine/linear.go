// Copyright 2020 the V8 project authors. All rights reserved.
// This file is an algorithmic Go port of V8's experimental RegExp engine;
// see LICENSE.v8 for the applicable license and attribution.
package regexpengine

import (
	"fmt"
)

// linearProgram is a Thompson-style prioritized NFA. It has no backtracking
// stack: each instruction is visited at most once at each UTF-16 input
// position. This is the execution model used by V8's experimental engine.
type linearProgram struct {
	code        []linearInstruction
	flags       Flags
	captures    int
	start       int
	lookbehinds []*linearProgram
}

type linearOp uint8

const (
	linearNever linearOp = iota
	linearRune
	linearAny
	linearClass
	linearSplit
	linearJump
	linearSave
	linearClearCapture
	linearAssert
	linearLookbehind
	linearBeginLoop
	linearEndLoop
	linearMatch
)

type linearInstruction struct {
	op         linearOp
	out        int
	alt        int
	rune       rune
	class      []classTerm
	negated    bool
	dotAll     bool
	multiline  bool
	capture    int
	assertion  nodeKind
	lookbehind int
}

type linearPatch struct {
	pc  int
	alt bool
}

type linearFragment struct {
	start int
	outs  []linearPatch
}

type linearCompiler struct {
	code        []linearInstruction
	flags       Flags
	lookbehinds []*linearProgram
}

// compileLinear lowers the AST using ordered NFA threads. It accepts the
// regular subset supported by V8's experimental engine and refuses patterns
// that would require backtracking or unsupported Unicode behavior.
func compileLinear(root *node, flags Flags, captures int) (*linearProgram, error) {
	if reason := linearUnsupportedReason(root, flags); reason != "" {
		return nil, fmt.Errorf("regular expression cannot use the linear engine: %s", reason)
	}
	c := &linearCompiler{flags: flags}
	begin := c.emit(linearInstruction{op: linearSave, capture: 0})
	fragment, err := c.compile(root)
	if err != nil {
		return nil, err
	}
	c.code[begin].out = fragment.start
	end := c.emit(linearInstruction{op: linearSave, capture: 1})
	c.patch(fragment.outs, end)
	match := c.emit(linearInstruction{op: linearMatch})
	c.code[end].out = match
	return &linearProgram{
		code:        c.code,
		flags:       flags,
		captures:    captures,
		start:       begin,
		lookbehinds: c.lookbehinds,
	}, nil
}

// linearUnsupportedReason mirrors the material checks in V8's
// ExperimentalCompiler::CanBeHandled for flags, captures, and quantifiers.
func linearUnsupportedReason(root *node, flags Flags) string {
	if !flags.Linear {
		return "the linear flag is not set"
	}
	if flags.HasIndices || flags.IgnoreCase || IsEitherUnicode(flags) {
		return "the linear engine only supports g, y, m, s, and l flags"
	}
	if reason := linearVisit(root, 1, flags); reason != "" {
		return reason
	}
	return ""
}

func linearVisit(n *node, replicationFactor int, flags Flags) string {
	if n == nil {
		return "invalid empty AST node"
	}
	if n.ignoreCase {
		return "scoped ignore-case modifiers are unsupported by the linear engine"
	}
	switch n.kind {
	case nodeBackref:
		return "backreferences require backtracking"
	case nodeLookaround:
		if flags.Global || flags.Sticky {
			return "V8's linear engine does not support lookarounds with global or sticky flags"
		}
		if !n.behind {
			return "V8's default linear engine does not support lookahead"
		}
		if linearContainsCapture(n.children[0]) {
			return "V8's default linear engine only supports captureless lookbehind"
		}
	case nodeClass:
		if n.set != nil {
			return "Unicode set expressions and class strings are unsupported"
		}
	case nodeRepeat:
		const maxReplicationFactor = 16
		if n.min > maxReplicationFactor || (n.max >= 0 && n.max > maxReplicationFactor) {
			return "quantifier bounds exceed the linear bytecode replication limit"
		}
		localReplication := n.max
		if n.max < 0 {
			if n.min > 0 && linearMinimumWidth(n.children[0]) > 0 {
				localReplication = maxInt(n.min, 1)
			} else {
				localReplication = n.min + 1
			}
		}
		if localReplication != 0 && replicationFactor > maxReplicationFactor/localReplication {
			return "nested quantifiers exceed the linear bytecode replication limit"
		}
		if reason := linearVisit(n.children[0], replicationFactor*localReplication, flags); reason != "" {
			return reason
		}
		return ""
	case nodeEmpty, nodeLiteral, nodeAny, nodeConcat, nodeAlternate, nodeCapture, nodeBegin, nodeEnd, nodeBoundary:
	default:
		return "unsupported AST node"
	}
	for _, child := range n.children {
		if reason := linearVisit(child, replicationFactor, flags); reason != "" {
			return reason
		}
	}
	return ""
}

func linearContainsCapture(root *node) bool {
	if root == nil {
		return false
	}
	if root.kind == nodeCapture {
		return true
	}
	for _, child := range root.children {
		if linearContainsCapture(child) {
			return true
		}
	}
	return false
}

func linearMinimumWidth(n *node) int {
	const capWidth = 17
	if n == nil {
		return 0
	}
	switch n.kind {
	case nodeEmpty, nodeBegin, nodeEnd, nodeBoundary, nodeLookaround:
		return 0
	case nodeLiteral, nodeAny, nodeClass:
		return 1
	case nodeCapture:
		return linearMinimumWidth(n.children[0])
	case nodeConcat:
		width := 0
		for _, child := range n.children {
			width += linearMinimumWidth(child)
			if width >= capWidth {
				return capWidth
			}
		}
		return width
	case nodeAlternate:
		if len(n.children) == 0 {
			return 0
		}
		width := capWidth
		for _, child := range n.children {
			width = minInt(width, linearMinimumWidth(child))
		}
		return width
	case nodeRepeat:
		width := linearMinimumWidth(n.children[0])
		if width == 0 || n.min == 0 {
			return 0
		}
		if n.min > capWidth/width {
			return capWidth
		}
		return width * n.min
	default:
		return 0
	}
}

func (c *linearCompiler) emit(instruction linearInstruction) int {
	if instruction.op != linearSplit {
		instruction.out = -1
		instruction.alt = -1
	}
	c.code = append(c.code, instruction)
	return len(c.code) - 1
}

func (c *linearCompiler) patch(outs []linearPatch, target int) {
	for _, patch := range outs {
		if patch.alt {
			c.code[patch.pc].alt = target
		} else {
			c.code[patch.pc].out = target
		}
	}
}

func (c *linearCompiler) append(first, second linearFragment) linearFragment {
	c.patch(first.outs, second.start)
	return linearFragment{start: first.start, outs: second.outs}
}

func (c *linearCompiler) empty() linearFragment {
	pc := c.emit(linearInstruction{op: linearJump})
	return linearFragment{start: pc, outs: []linearPatch{{pc: pc}}}
}

func (c *linearCompiler) compile(n *node) (linearFragment, error) {
	switch n.kind {
	case nodeEmpty:
		return c.empty(), nil
	case nodeLiteral:
		pc := c.emit(linearInstruction{op: linearRune, rune: n.rune})
		return linearFragment{start: pc, outs: []linearPatch{{pc: pc}}}, nil
	case nodeAny:
		pc := c.emit(linearInstruction{op: linearAny, dotAll: n.dotAll})
		return linearFragment{start: pc, outs: []linearPatch{{pc: pc}}}, nil
	case nodeClass:
		if n.set != nil {
			return linearFragment{}, fmt.Errorf("Unicode set class is not linear-engine compatible")
		}
		pc := c.emit(linearInstruction{op: linearClass, class: n.class, negated: n.negated})
		return linearFragment{start: pc, outs: []linearPatch{{pc: pc}}}, nil
	case nodeConcat:
		if len(n.children) == 0 {
			return c.empty(), nil
		}
		fragment, err := c.compile(n.children[0])
		if err != nil {
			return linearFragment{}, err
		}
		for _, child := range n.children[1:] {
			next, err := c.compile(child)
			if err != nil {
				return linearFragment{}, err
			}
			fragment = c.append(fragment, next)
		}
		return fragment, nil
	case nodeAlternate:
		if len(n.children) == 0 {
			pc := c.emit(linearInstruction{op: linearNever})
			return linearFragment{start: pc}, nil
		}
		fragment, err := c.compile(n.children[len(n.children)-1])
		if err != nil {
			return linearFragment{}, err
		}
		for index := len(n.children) - 2; index >= 0; index-- {
			left, err := c.compile(n.children[index])
			if err != nil {
				return linearFragment{}, err
			}
			split := c.emit(linearInstruction{op: linearSplit, out: left.start, alt: fragment.start})
			outs := append(left.outs, fragment.outs...)
			fragment = linearFragment{start: split, outs: outs}
		}
		return fragment, nil
	case nodeCapture:
		begin := c.emit(linearInstruction{op: linearSave, capture: 2 * n.capture})
		body, err := c.compile(n.children[0])
		if err != nil {
			return linearFragment{}, err
		}
		c.code[begin].out = body.start
		end := c.emit(linearInstruction{op: linearSave, capture: 2*n.capture + 1})
		c.patch(body.outs, end)
		return linearFragment{start: begin, outs: []linearPatch{{pc: end}}}, nil
	case nodeRepeat:
		return c.compileRepeat(n)
	case nodeBegin, nodeEnd, nodeBoundary:
		pc := c.emit(linearInstruction{op: linearAssert, assertion: n.kind, negated: n.negated, multiline: n.multiline})
		return linearFragment{start: pc, outs: []linearPatch{{pc: pc}}}, nil
	case nodeLookaround:
		if !n.behind || c.flags.Global || c.flags.Sticky || linearContainsCapture(n.children[0]) {
			return linearFragment{}, fmt.Errorf("lookaround is not supported by the V8-compatible linear engine subset")
		}
		lookbehind, err := compileLinearLookbehind(n.children[0], c.flags)
		if err != nil {
			return linearFragment{}, err
		}
		index := len(c.lookbehinds)
		c.lookbehinds = append(c.lookbehinds, lookbehind)
		pc := c.emit(linearInstruction{op: linearLookbehind, lookbehind: index, negated: n.negative})
		return linearFragment{start: pc, outs: []linearPatch{{pc: pc}}}, nil
	case nodeBackref:
		return linearFragment{}, fmt.Errorf("node kind %d is not supported by the linear engine", n.kind)
	default:
		return linearFragment{}, fmt.Errorf("unknown regexp AST node kind %d", n.kind)
	}
}

// compileLinearLookbehind builds the captureless body automaton used to fill
// the lookbehind truth table. It shares the interpreter's NFA representation,
// but has no synthetic whole-match capture registers.
func compileLinearLookbehind(root *node, flags Flags) (*linearProgram, error) {
	c := &linearCompiler{flags: flags}
	fragment, err := c.compile(root)
	if err != nil {
		return nil, err
	}
	match := c.emit(linearInstruction{op: linearMatch})
	c.patch(fragment.outs, match)
	return &linearProgram{code: c.code, flags: flags, start: fragment.start, lookbehinds: c.lookbehinds}, nil
}

func (c *linearCompiler) compileRepeat(n *node) (linearFragment, error) {
	body := n.children[0]
	result := linearFragment{start: -1}
	for count := 0; count < n.min; count++ {
		fragment, err := c.compileRepeatedBody(body, false)
		if err != nil {
			return linearFragment{}, err
		}
		if result.start < 0 {
			result = fragment
		} else {
			result = c.append(result, fragment)
		}
	}

	if n.max < 0 {
		star, err := c.compileStar(body, n.greedy)
		if err != nil {
			return linearFragment{}, err
		}
		if result.start < 0 {
			return star, nil
		}
		return c.append(result, star), nil
	}

	for count := n.min; count < n.max; count++ {
		optional, err := c.compileOptional(body, n.greedy)
		if err != nil {
			return linearFragment{}, err
		}
		if result.start < 0 {
			result = optional
		} else {
			result = c.append(result, optional)
		}
	}
	if result.start < 0 {
		return c.empty(), nil
	}
	return result, nil
}

func (c *linearCompiler) compileRepeatedBody(body *node, rejectEmptyIteration bool) (linearFragment, error) {
	fragment, err := c.compile(body)
	if err != nil {
		return linearFragment{}, err
	}
	fragment = c.prependCaptureClears(body, fragment)
	if !rejectEmptyIteration {
		return fragment, nil
	}
	begin := c.emit(linearInstruction{op: linearBeginLoop})
	c.code[begin].out = fragment.start
	end := c.emit(linearInstruction{op: linearEndLoop})
	c.patch(fragment.outs, end)
	return linearFragment{start: begin, outs: []linearPatch{{pc: end}}}, nil
}

func (c *linearCompiler) prependCaptureClears(body *node, fragment linearFragment) linearFragment {
	captures := linearCaptureIDs(body)
	if len(captures) == 0 {
		return fragment
	}
	var prefix linearFragment
	for index, capture := range captures {
		pc := c.emit(linearInstruction{op: linearClearCapture, capture: capture})
		clear := linearFragment{start: pc, outs: []linearPatch{{pc: pc}}}
		if index == 0 {
			prefix = clear
		} else {
			prefix = c.append(prefix, clear)
		}
	}
	return c.append(prefix, fragment)
}

func (c *linearCompiler) compileOptional(body *node, greedy bool) (linearFragment, error) {
	fragment, err := c.compileRepeatedBody(body, true)
	if err != nil {
		return linearFragment{}, err
	}
	pc := c.emit(linearInstruction{op: linearSplit, out: -1, alt: -1})
	if greedy {
		c.code[pc].out = fragment.start
		return linearFragment{start: pc, outs: append(fragment.outs, linearPatch{pc: pc, alt: true})}, nil
	}
	c.code[pc].alt = fragment.start
	return linearFragment{start: pc, outs: append(fragment.outs, linearPatch{pc: pc})}, nil
}

func (c *linearCompiler) compileStar(body *node, greedy bool) (linearFragment, error) {
	fragment, err := c.compileRepeatedBody(body, true)
	if err != nil {
		return linearFragment{}, err
	}
	pc := c.emit(linearInstruction{op: linearSplit, out: -1, alt: -1})
	if greedy {
		c.code[pc].out = fragment.start
		c.patch(fragment.outs, pc)
		return linearFragment{start: pc, outs: []linearPatch{{pc: pc, alt: true}}}, nil
	}
	c.code[pc].alt = fragment.start
	c.patch(fragment.outs, pc)
	return linearFragment{start: pc, outs: []linearPatch{{pc: pc}}}, nil
}

func linearCaptureIDs(root *node) []int {
	seen := make(map[int]struct{})
	var captures []int
	var visit func(*node)
	visit = func(n *node) {
		if n == nil {
			return
		}
		if n.kind == nodeCapture {
			if _, ok := seen[n.capture]; !ok {
				seen[n.capture] = struct{}{}
				captures = append(captures, n.capture)
			}
		}
		for _, child := range n.children {
			visit(child)
		}
	}
	visit(root)
	return captures
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type linearThread struct {
	pc                          int
	captures                    []int
	consumedSinceLastQuantifier bool
}

func (p *linearProgram) FindUTF16(input []uint16, start int) (Match, bool) {
	if start < 0 {
		start = 0
	}
	if start > len(input) || len(p.code) == 0 {
		return Match{}, false
	}
	lookbehindTables := p.computeLookbehindTables(input)
	stateCapacity := 2 * len(p.code)
	active := make([]linearThread, 0, stateCapacity)
	closed := make([]linearThread, 0, stateCapacity)
	next := make([]linearThread, 0, stateCapacity)
	var best *Match
	seen := make([]uint32, 2*len(p.code))
	var generation uint32
	for position := start; position <= len(input); position++ {
		generation = nextLinearGeneration(seen, generation)
		closed = closed[:0]
		for _, thread := range active {
			p.addThread(&closed, thread.pc, thread.captures, thread.consumedSinceLastQuantifier, input, position, seen, generation, lookbehindTables)
		}
		if best == nil && (!p.flags.Sticky || position == start) {
			captures := make([]int, 2*(p.captures+1))
			for index := range captures {
				captures[index] = -1
			}
			p.addThread(&closed, p.start, captures, true, input, position, seen, generation, lookbehindTables)
		}

		matchIndex := -1
		for index, thread := range closed {
			if p.code[thread.pc].op == linearMatch {
				matchIndex = index
				candidate := p.matchFromThread(thread, position)
				best = &candidate
				closed = closed[:index]
				break
			}
		}
		if matchIndex == 0 {
			return *best, true
		}
		if position == len(input) {
			if best != nil {
				return *best, true
			}
			return Match{}, false
		}

		next = next[:0]
		character := rune(input[position])
		for _, thread := range closed {
			instruction := p.code[thread.pc]
			if p.consumes(instruction, character) {
				next = append(next, linearThread{pc: instruction.out, captures: thread.captures, consumedSinceLastQuantifier: true})
			}
		}
		if len(next) == 0 {
			if best != nil {
				return *best, true
			}
			if p.flags.Sticky {
				return Match{}, false
			}
		}
		active, next = next, active[:0]
		if best != nil && len(active) == 0 {
			return *best, true
		}
	}
	if best != nil {
		return *best, true
	}
	return Match{}, false
}

func (p *linearProgram) computeLookbehindTables(input []uint16) [][]bool {
	tables := make([][]bool, len(p.lookbehinds))
	for index, subprogram := range p.lookbehinds {
		nestedTables := subprogram.computeLookbehindTables(input)
		tables[index] = subprogram.matchEnds(input, nestedTables)
	}
	return tables
}

// matchEnds computes every UTF-16 position at which this automaton accepts.
// New start states are injected as the input advances, just as in the V8
// breadth-first interpreter; the work is bounded by input length times the
// automaton size rather than re-running a lookbehind from every position.
func (p *linearProgram) matchEnds(input []uint16, lookbehindTables [][]bool) []bool {
	ends := make([]bool, len(input)+1)
	stateCapacity := 2 * len(p.code)
	active := make([]linearThread, 0, stateCapacity)
	closed := make([]linearThread, 0, stateCapacity)
	next := make([]linearThread, 0, stateCapacity)
	seen := make([]uint32, 2*len(p.code))
	var generation uint32
	for position := 0; position <= len(input); position++ {
		generation = nextLinearGeneration(seen, generation)
		closed = closed[:0]
		for _, thread := range active {
			p.addThread(&closed, thread.pc, nil, thread.consumedSinceLastQuantifier, input, position, seen, generation, lookbehindTables)
		}
		p.addThread(&closed, p.start, nil, true, input, position, seen, generation, lookbehindTables)
		for _, thread := range closed {
			if p.code[thread.pc].op == linearMatch {
				ends[position] = true
				break
			}
		}
		if position == len(input) {
			break
		}
		next = next[:0]
		character := rune(input[position])
		for _, thread := range closed {
			instruction := p.code[thread.pc]
			if p.consumes(instruction, character) {
				next = append(next, linearThread{pc: instruction.out, consumedSinceLastQuantifier: true})
			}
		}
		active, next = next, active[:0]
	}
	return ends
}

func nextLinearGeneration(seen []uint32, generation uint32) uint32 {
	generation++
	if generation == 0 {
		clear(seen)
		return 1
	}
	return generation
}

func (p *linearProgram) addThread(out *[]linearThread, pc int, captures []int, consumedSinceLastQuantifier bool, input []uint16, position int, seen []uint32, generation uint32, lookbehindTables [][]bool) {
	for pc >= 0 && pc < len(p.code) {
		state := 2 * pc
		if consumedSinceLastQuantifier {
			state++
		}
		if seen[state] == generation {
			return
		}
		seen[state] = generation
		instruction := p.code[pc]
		switch instruction.op {
		case linearJump:
			pc = instruction.out
		case linearSplit:
			p.addThread(out, instruction.out, captures, consumedSinceLastQuantifier, input, position, seen, generation, lookbehindTables)
			pc = instruction.alt
		case linearSave:
			captures = cloneInts(captures)
			if instruction.capture >= 0 && instruction.capture < len(captures) {
				captures[instruction.capture] = position
			}
			pc = instruction.out
		case linearClearCapture:
			captures = cloneInts(captures)
			register := instruction.capture * 2
			if register >= 0 && register+1 < len(captures) {
				captures[register], captures[register+1] = -1, -1
			}
			pc = instruction.out
		case linearAssert:
			if !p.assertionMatches(instruction, input, position) {
				return
			}
			pc = instruction.out
		case linearLookbehind:
			matched := instruction.lookbehind >= 0 && instruction.lookbehind < len(lookbehindTables) && position < len(lookbehindTables[instruction.lookbehind]) && lookbehindTables[instruction.lookbehind][position]
			if matched == instruction.negated {
				return
			}
			pc = instruction.out
		case linearBeginLoop:
			consumedSinceLastQuantifier = false
			pc = instruction.out
		case linearEndLoop:
			if !consumedSinceLastQuantifier {
				return
			}
			pc = instruction.out
		case linearNever:
			return
		default:
			*out = append(*out, linearThread{pc: pc, captures: captures, consumedSinceLastQuantifier: consumedSinceLastQuantifier})
			return
		}
	}
}

func cloneInts(values []int) []int {
	return append([]int(nil), values...)
}

func (p *linearProgram) assertionMatches(instruction linearInstruction, input []uint16, position int) bool {
	switch instruction.assertion {
	case nodeBegin:
		return lineStartAssertion(input, position, instruction.multiline)
	case nodeEnd:
		return lineEndAssertion(input, position, instruction.multiline)
	case nodeBoundary:
		before := position > 0 && isWordRune(previousCharacter(input, position, false), p.flags)
		after := position < len(input) && isWordRune(rune(input[position]), p.flags)
		matched := before != after
		if instruction.negated {
			return !matched
		}
		return matched
	default:
		return false
	}
}

func (p *linearProgram) consumes(instruction linearInstruction, character rune) bool {
	switch instruction.op {
	case linearRune:
		return character == instruction.rune
	case linearAny:
		return instruction.dotAll || !isLineTerminator(character)
	case linearClass:
		return classContains(character, instruction.class, instruction.negated, false, false, false)
	default:
		return false
	}
}

func (p *linearProgram) matchFromThread(thread linearThread, end int) Match {
	captures := thread.captures
	spans := make([]Span, p.captures+1)
	for index := range spans {
		start, finish := -1, -1
		if 2*index+1 < len(captures) {
			start, finish = captures[2*index], captures[2*index+1]
		}
		spans[index] = Span{Start: start, End: finish}
	}
	start := 0
	if len(captures) != 0 && captures[0] >= 0 {
		start = captures[0]
	}
	if len(spans) != 0 {
		spans[0] = Span{Start: start, End: end}
	}
	return Match{Start: start, End: end, Captures: spans}
}
