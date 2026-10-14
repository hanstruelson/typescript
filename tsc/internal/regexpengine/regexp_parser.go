package regexpengine

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"unicode/utf16"
)

type v8ASTParserTextBuilder struct {
	flags      Flags
	terms      *[]any
	characters []uint16
	text       []any
}

func v8ASTNewParserTextBuilder(flags Flags, terms *[]any) *v8ASTParserTextBuilder {
	if terms == nil {
		panic("V8 TextBuilder requires term storage")
	}
	return &v8ASTParserTextBuilder{flags: flags, terms: terms}
}

func (builder *v8ASTParserTextBuilder) flushCharacters() {
	if len(builder.characters) == 0 {
		return
	}
	atom := (v8ASTFactory{}).newAtom(builder.characters)
	builder.characters = nil
	builder.text = append(builder.text, atom)
}

func (builder *v8ASTParserTextBuilder) flushText() {
	builder.flushCharacters()
	switch len(builder.text) {
	case 0:
		return
	case 1:
		*builder.terms = append(*builder.terms, builder.text[0])
	default:
		text := (v8ASTFactory{}).newText()
		for _, element := range builder.text {
			v8ASTAppendToText(element, text)
		}
		*builder.terms = append(*builder.terms, text)
	}
	builder.text = builder.text[:0]
}

func (builder *v8ASTParserTextBuilder) addCharacter(character uint16) {
	if builder.needsDesugaringForIgnoreCase(rune(character)) {
		builder.addClassRangesForDesugaring(rune(character))
		return
	}
	builder.characters = append(builder.characters, character)
}

func (builder *v8ASTParserTextBuilder) addUnicodeCharacter(character rune) {
	if character < 0 || character > v8ASTMaxCodePoint {
		panic("V8 parser produced an invalid Unicode code point")
	}
	if character > 0xffff {
		if builder.needsDesugaringForIgnoreCase(character) {
			builder.addClassRangesForDesugaring(character)
			return
		}
		value := uint32(character - 0x10000)
		atom := (v8ASTFactory{}).newAtom([]uint16{
			uint16(0xd800 + (value >> 10)),
			uint16(0xdc00 + (value & 0x3ff)),
		})
		builder.addAtom(atom)
		return
	}
	if IsEitherUnicode(builder.flags) && character >= 0xd800 && character <= 0xdfff {
		builder.addClassRangesForDesugaring(character)
		return
	}
	builder.addCharacter(uint16(character))
}

func (builder *v8ASTParserTextBuilder) addClassRanges(class *v8ASTClassRanges) {
	if class == nil {
		panic("V8 TextBuilder cannot add nil class ranges")
	}
	if builder.needsDesugaringForUnicode(class) {
		builder.addTerm(class)
	} else {
		builder.addAtom(class)
	}
}

func (builder *v8ASTParserTextBuilder) addClassRangesForDesugaring(character rune) {
	class := (v8ASTFactory{}).newClassRanges(v8ASTRangeList(v8ASTRangeSingleton(character)), 0)
	builder.addTerm(class)
}

func (builder *v8ASTParserTextBuilder) addAtom(atom any) {
	if atom == nil || !v8ASTTreeIsTextElement(atom) {
		panic("V8 TextBuilder atom must be a text element")
	}
	builder.flushCharacters()
	builder.text = append(builder.text, atom)
}

func (builder *v8ASTParserTextBuilder) addTerm(term any) {
	if term == nil || !v8ASTTreeIsTextElement(term) {
		panic("V8 TextBuilder term must be a text element")
	}
	builder.flushText()
	*builder.terms = append(*builder.terms, term)
}

func (builder *v8ASTParserTextBuilder) needsDesugaringForUnicode(class *v8ASTClassRanges) bool {
	if !IsEitherUnicode(builder.flags) {
		return false
	}
	if builder.flags.IgnoreCase {
		return true
	}
	ranges := v8ASTCanonicalizeRanges(v8ASTCharacterSetRanges(&class.set))
	class.set.ranges = ranges
	if class.flags&v8ASTClassNegated != 0 {
		ranges = v8ASTNegateRanges(ranges)
	}
	for _, current := range ranges {
		if current.to >= 0x10000 || (current.from <= 0xdfff && current.to >= 0xd800) {
			return true
		}
	}
	return false
}

func (builder *v8ASTParserTextBuilder) needsDesugaringForIgnoreCase(character rune) bool {
	if !IsEitherUnicode(builder.flags) || !builder.flags.IgnoreCase {
		return false
	}
	// V8 consults SupplementaryEquivalentsSet for BMP characters, then closes
	// over the simple-fold class for supplementary characters.
	if character <= 0xffff {
		return v8ASTCaseFoldingHasSupplementaryEquivalents(uint16(character))
	}
	return v8ASTCaseFoldingHasOtherUnicodeEquivalent(character)
}

func (builder *v8ASTParserTextBuilder) popLastAtom() any {
	if len(builder.characters) != 0 {
		last := builder.characters[len(builder.characters)-1]
		if len(builder.characters) > 1 {
			prefix := (v8ASTFactory{}).newAtom(builder.characters[:len(builder.characters)-1])
			builder.text = append(builder.text, prefix)
		}
		builder.characters = nil
		return (v8ASTFactory{}).newAtom([]uint16{last})
	}
	if len(builder.text) == 0 {
		return nil
	}
	atom := builder.text[len(builder.text)-1]
	builder.text = builder.text[:len(builder.text)-1]
	return atom
}

func (builder *v8ASTParserTextBuilder) toRegExp() any {
	builder.flushText()
	switch len(*builder.terms) {
	case 0:
		return (v8ASTFactory{}).newEmpty()
	case 1:
		return (*builder.terms)[0]
	default:
		return (v8ASTFactory{}).newAlternative(append([]any(nil), (*builder.terms)...))
	}
}

// v8ASTParserBuilder mirrors Builder<kBuildAST>, including its pending-empty
// state and the distinction between text elements, terms, and assertions.
type v8ASTParserBuilder struct {
	flags        Flags
	pendingEmpty bool
	terms        []any
	alternatives []any
	textBuilder  *v8ASTParserTextBuilder
}

func v8ASTNewParserBuilder(flags Flags) *v8ASTParserBuilder {
	builder := &v8ASTParserBuilder{flags: flags}
	builder.textBuilder = v8ASTNewParserTextBuilder(flags, &builder.terms)
	return builder
}

func (builder *v8ASTParserBuilder) flushText() {
	builder.textBuilder.flushText()
}

func (builder *v8ASTParserBuilder) getFlags() Flags     { return builder.flags }
func (builder *v8ASTParserBuilder) ignoreCase() bool    { return builder.flags.IgnoreCase }
func (builder *v8ASTParserBuilder) multiline() bool     { return builder.flags.Multiline }
func (builder *v8ASTParserBuilder) dotAll() bool        { return builder.flags.DotAll }
func (builder *v8ASTParserBuilder) isUnicodeMode() bool { return IsEitherUnicode(builder.flags) }

func (builder *v8ASTParserBuilder) addCharacter(character uint16) {
	builder.pendingEmpty = false
	builder.textBuilder.addCharacter(character)
}

func (builder *v8ASTParserBuilder) addUnicodeCharacter(character rune) {
	builder.pendingEmpty = false
	builder.textBuilder.addUnicodeCharacter(character)
}

func (builder *v8ASTParserBuilder) addEmpty() {
	builder.pendingEmpty = true
}

func (builder *v8ASTParserBuilder) addClassRanges(class *v8ASTClassRanges) {
	builder.pendingEmpty = false
	builder.textBuilder.addClassRanges(class)
}

func (builder *v8ASTParserBuilder) addAtom(atom any) {
	if atom == nil {
		panic("V8 AST Builder cannot add a nil atom")
	}
	if v8ASTIsEmpty(atom) {
		builder.addEmpty()
		return
	}
	builder.pendingEmpty = false
	if v8ASTTreeIsTextElement(atom) {
		builder.textBuilder.addAtom(atom)
		return
	}
	builder.flushText()
	builder.terms = append(builder.terms, atom)
}

func (builder *v8ASTParserBuilder) addTerm(term any) {
	if term == nil || v8ASTIsEmpty(term) {
		panic("V8 AST Builder term must be non-empty")
	}
	builder.pendingEmpty = false
	if v8ASTTreeIsTextElement(term) {
		builder.textBuilder.addTerm(term)
		return
	}
	builder.flushText()
	builder.terms = append(builder.terms, term)
}

func (builder *v8ASTParserBuilder) addAssertion(assertion any) {
	if assertion == nil || v8ASTIsEmpty(assertion) {
		panic("V8 AST Builder assertion must be non-empty")
	}
	builder.pendingEmpty = false
	builder.flushText()
	builder.terms = append(builder.terms, assertion)
}

func (builder *v8ASTParserBuilder) newAlternative() {
	builder.flushTerms()
}

func (builder *v8ASTParserBuilder) flushTerms() {
	builder.flushText()
	var alternative any
	switch len(builder.terms) {
	case 0:
		alternative = (v8ASTFactory{}).newEmpty()
	case 1:
		alternative = builder.terms[0]
	default:
		alternative = (v8ASTFactory{}).newAlternative(append([]any(nil), builder.terms...))
	}
	builder.alternatives = append(builder.alternatives, alternative)
	builder.terms = builder.terms[:0]
	builder.textBuilder.terms = &builder.terms
}

func (builder *v8ASTParserBuilder) toRegExp() any {
	builder.flushTerms()
	switch len(builder.alternatives) {
	case 0:
		return (v8ASTFactory{}).newEmpty()
	case 1:
		return builder.alternatives[0]
	default:
		return (v8ASTFactory{}).newDisjunction(append([]any(nil), builder.alternatives...))
	}
}

func (builder *v8ASTParserBuilder) addQuantifier(minimum, maximum, index int, kind v8ASTQuantifierType) bool {
	if builder.pendingEmpty {
		builder.pendingEmpty = false
		return true
	}
	atom := builder.textBuilder.popLastAtom()
	if atom != nil {
		builder.flushText()
	} else if len(builder.terms) != 0 {
		atom = builder.terms[len(builder.terms)-1]
		builder.terms = builder.terms[:len(builder.terms)-1]
		if lookaround, ok := atom.(*v8ASTLookaround); ok {
			if IsEitherUnicode(builder.flags) || lookaround.typeOf == v8ASTLookbehind {
				return false
			}
		}
		if v8TonodeTreeMaxMatch(atom) == 0 {
			if minimum == 0 {
				return true
			}
			builder.terms = append(builder.terms, atom)
			return true
		}
	} else {
		panic("V8 Builder quantifier has no preceding atom")
	}
	quantifier := (v8ASTFactory{}).newQuantifier(minimum, maximum, kind, index, atom)
	builder.terms = append(builder.terms, quantifier)
	return true
}

// v8ASTSyntaxBuilder is the zero-AST Builder<kVerifySyntax> specialization.
// It retains the lookaround state needed to reject invalid quantifiers while
// discarding every node passed through the parser.
type v8ASTSyntaxBuilder struct {
	flags              Flags
	lastIsLookaround   bool
	lastLookaroundKind v8ASTLookaroundType
}

func v8ASTNewSyntaxBuilder(flags Flags) *v8ASTSyntaxBuilder {
	return &v8ASTSyntaxBuilder{flags: flags, lastLookaroundKind: v8ASTLookahead}
}

func (builder *v8ASTSyntaxBuilder) addCharacter(uint16) { builder.lastIsLookaround = false }
func (builder *v8ASTSyntaxBuilder) addUnicodeCharacter(rune) {
	builder.lastIsLookaround = false
}
func (builder *v8ASTSyntaxBuilder) addEmpty() { builder.lastIsLookaround = false }
func (builder *v8ASTSyntaxBuilder) addClassRanges(*v8ASTClassRanges) {
	builder.lastIsLookaround = false
}
func (builder *v8ASTSyntaxBuilder) addAtom(atom any) {
	if atom != nil && v8ASTIsEmpty(atom) {
		builder.addEmpty()
		return
	}
	if _, ok := atom.(*v8ASTLookaround); ok {
		panic("V8 syntax Builder requires AddLookaround for lookaround nodes")
	}
	builder.lastIsLookaround = false
}
func (builder *v8ASTSyntaxBuilder) addLookaround(kind v8ASTLookaroundType) {
	builder.lastIsLookaround = true
	builder.lastLookaroundKind = kind
}
func (builder *v8ASTSyntaxBuilder) flushText()          {}
func (builder *v8ASTSyntaxBuilder) toRegExp() any       { return nil }
func (builder *v8ASTSyntaxBuilder) getFlags() Flags     { return builder.flags }
func (builder *v8ASTSyntaxBuilder) ignoreCase() bool    { return builder.flags.IgnoreCase }
func (builder *v8ASTSyntaxBuilder) multiline() bool     { return builder.flags.Multiline }
func (builder *v8ASTSyntaxBuilder) dotAll() bool        { return builder.flags.DotAll }
func (builder *v8ASTSyntaxBuilder) isUnicodeMode() bool { return IsEitherUnicode(builder.flags) }

type v8ASTParseMode uint8

const (
	v8ASTBuildAST v8ASTParseMode = iota
	v8ASTVerifySyntax
)

type v8ASTClassSetOperandType uint8

const (
	v8ASTClassSetCharacter v8ASTClassSetOperandType = iota
	v8ASTClassStringDisjunction
	v8ASTNestedClass
	v8ASTCharacterClassEscape
	v8ASTClassSetRange
)

type v8ASTInClassEscapeState uint8

const (
	v8ASTNotInClass v8ASTInClassEscapeState = iota
	v8ASTInClass
)

type v8ASTClassSetResult struct {
	tree              any
	mayContainStrings bool
	success           bool
}

func v8ASTNewClassSetResult(tree any, mayContainStrings bool) v8ASTClassSetResult {
	return v8ASTClassSetResult{tree: tree, mayContainStrings: mayContainStrings, success: true}
}

func v8ASTClassSetResultFailure() v8ASTClassSetResult {
	return v8ASTClassSetResult{success: false}
}

type v8ASTSyntaxClassSetResult struct {
	mayContainStrings bool
	success           bool
}

func v8ASTNewSyntaxClassSetResult(mayContainStrings bool) v8ASTSyntaxClassSetResult {
	return v8ASTSyntaxClassSetResult{mayContainStrings: mayContainStrings, success: true}
}

func v8ASTSyntaxClassSetResultFromBuild(result v8ASTClassSetResult) v8ASTSyntaxClassSetResult {
	return v8ASTSyntaxClassSetResult{mayContainStrings: result.mayContainStrings, success: result.success}
}

func v8ASTSyntaxClassSetResultFailure() v8ASTSyntaxClassSetResult {
	return v8ASTSyntaxClassSetResult{success: false}
}

const v8ASTParserEndMarker rune = 1 << 21

type v8ASTParserCursor struct {
	input         []uint16
	byteInput     []byte
	inputIsBytes  bool
	flags         Flags
	oneByte       bool
	current       rune
	nextPosition  int
	hasMore       bool
	forceUnicode  bool
	failed        bool
	errorCode     string
	errorPosition int
	stackOverflow func() bool
}

func v8ASTNewParserCursor(input []uint16, flags Flags, oneByte bool, stackOverflow func() bool) *v8ASTParserCursor {
	cursor := &v8ASTParserCursor{
		input: input, flags: flags, oneByte: oneByte, current: v8ASTParserEndMarker,
		hasMore: true, stackOverflow: stackOverflow,
	}
	cursor.advance()
	return cursor
}

func v8ASTNewParserCursorOneByte(input []byte, flags Flags, stackOverflow func() bool) *v8ASTParserCursor {
	cursor := &v8ASTParserCursor{
		byteInput: input, inputIsBytes: true, oneByte: true, flags: flags,
		current: v8ASTParserEndMarker, hasMore: true, stackOverflow: stackOverflow,
	}
	cursor.advance()
	return cursor
}

func (cursor *v8ASTParserCursor) currentValue() rune { return cursor.current }
func (cursor *v8ASTParserCursor) hasMoreInput() bool { return cursor.hasMore }
func (cursor *v8ASTParserCursor) hasNextInput() bool {
	return cursor.nextPosition < cursor.inputLength()
}
func (cursor *v8ASTParserCursor) inputLength() int {
	if cursor.inputIsBytes {
		return len(cursor.byteInput)
	}
	return len(cursor.input)
}
func (cursor *v8ASTParserCursor) flagsValue() Flags { return cursor.flags }
func (cursor *v8ASTParserCursor) failedParse() bool { return cursor.failed }
func (cursor *v8ASTParserCursor) errorState() (string, int) {
	return cursor.errorCode, cursor.errorPosition
}
func (cursor *v8ASTParserCursor) unicodeSets() bool { return cursor.flags.UnicodeSets }
func (cursor *v8ASTParserCursor) ignoreCase() bool  { return cursor.flags.IgnoreCase }
func (cursor *v8ASTParserCursor) isUnicodeMode() bool {
	return IsEitherUnicode(cursor.flags) || cursor.forceUnicode
}
func (cursor *v8ASTParserCursor) addUnicodeCaseEquivalents(mode v8ASTParseMode) bool {
	return mode == v8ASTBuildAST && cursor.isUnicodeMode() && cursor.flags.IgnoreCase
}

func (cursor *v8ASTParserCursor) inputAt(index int) uint16 {
	if index < 0 || index >= cursor.inputLength() {
		panic("V8 parser cursor index is out of range")
	}
	if cursor.inputIsBytes {
		return uint16(cursor.byteInput[index])
	}
	return cursor.input[index]
}

func (cursor *v8ASTParserCursor) inputUnits(start, end int) []uint16 {
	if start < 0 || end < start || end > cursor.inputLength() {
		panic("V8 parser input slice is out of range")
	}
	if !cursor.inputIsBytes {
		return cursor.input[start:end]
	}
	units := make([]uint16, end-start)
	for index, value := range cursor.byteInput[start:end] {
		units[index] = uint16(value)
	}
	return units
}

func (cursor *v8ASTParserCursor) readNext(updatePosition bool) rune {
	position := cursor.nextPosition
	first := cursor.inputAt(position)
	position++
	result := rune(first)
	if !cursor.oneByte && cursor.isUnicodeMode() && position < cursor.inputLength() && first >= 0xd800 && first <= 0xdbff {
		second := cursor.inputAt(position)
		if second >= 0xdc00 && second <= 0xdfff {
			result = 0x10000 + ((rune(first) - 0xd800) << 10) + rune(second) - 0xdc00
			position++
		}
	}
	if updatePosition {
		cursor.nextPosition = position
	}
	return result
}

func (cursor *v8ASTParserCursor) next() rune {
	if cursor.hasNextInput() {
		return cursor.readNext(false)
	}
	return v8ASTParserEndMarker
}

func (cursor *v8ASTParserCursor) position() int {
	rewind := 1
	if cursor.current != v8ASTParserEndMarker && cursor.current > 0xffff {
		rewind = 2
	}
	position := cursor.nextPosition - rewind
	if position < 0 {
		return 0
	}
	return position
}

func (cursor *v8ASTParserCursor) reportError(code string) {
	cursor.reportErrorAt(code, cursor.position())
}

func (cursor *v8ASTParserCursor) reportErrorAt(code string, position int) {
	if cursor.failed {
		return
	}
	cursor.failed = true
	cursor.errorCode = code
	cursor.errorPosition = position
	cursor.current = v8ASTParserEndMarker
	cursor.nextPosition = cursor.inputLength()
	cursor.hasMore = false
}

func (cursor *v8ASTParserCursor) advance() {
	if cursor.hasNextInput() {
		if cursor.stackOverflow != nil && cursor.stackOverflow() {
			cursor.reportError(errorString(errorStackOverflow))
		} else {
			cursor.current = cursor.readNext(true)
		}
		return
	}
	cursor.current = v8ASTParserEndMarker
	cursor.nextPosition = cursor.inputLength() + 1
	cursor.hasMore = false
}

func (cursor *v8ASTParserCursor) advanceBy(distance int) {
	cursor.nextPosition += distance - 1
	cursor.advance()
}

func (cursor *v8ASTParserCursor) rewindByOneCodePoint() {
	if !cursor.hasMore {
		return
	}
	distance := -1
	if cursor.current > 0xffff {
		distance = -2
	}
	cursor.advanceBy(distance)
}

func (cursor *v8ASTParserCursor) reset(position int) {
	if position < 0 || position > cursor.inputLength() {
		panic("V8 parser cursor reset position is out of range")
	}
	cursor.nextPosition = position
	cursor.hasMore = position < cursor.inputLength()
	cursor.advance()
}

func (cursor *v8ASTParserCursor) isSyntaxCharacterOrSlash(character rune) bool {
	switch character {
	case '^', '$', '\\', '.', '*', '+', '?', '(', ')', '[', ']', '{', '}', '|', '/':
		return true
	default:
		return false
	}
}

func (cursor *v8ASTParserCursor) isClassSetSyntaxCharacter(character rune) bool {
	switch character {
	case '(', ')', '[', ']', '{', '}', '/', '-', '\\', '|':
		return true
	default:
		return false
	}
}

func (cursor *v8ASTParserCursor) isClassSetReservedPunctuator(character rune) bool {
	switch character {
	case '&', '-', '!', '#', '%', ',', ':', ';', '<', '=', '>', '@', '`', '~':
		return true
	default:
		return false
	}
}

func (cursor *v8ASTParserCursor) isClassSetReservedDoublePunctuator(character rune) bool {
	switch character {
	case '&', '!', '#', '$', '%', '*', '+', ',', '.', ':', ';', '<', '=', '>', '?', '@', '^', '`', '~':
		return cursor.next() == character
	default:
		return false
	}
}

type v8ASTForceUnicodeScope struct {
	cursor *v8ASTParserCursor
	active bool
}

func v8ASTNewForceUnicodeScope(cursor *v8ASTParserCursor) *v8ASTForceUnicodeScope {
	if cursor == nil || cursor.forceUnicode {
		panic("V8 ForceUnicodeScope requires a cursor outside another forced scope")
	}
	cursor.forceUnicode = true
	return &v8ASTForceUnicodeScope{cursor: cursor, active: true}
}

func (scope *v8ASTForceUnicodeScope) close() {
	if scope == nil || !scope.active || scope.cursor == nil || !scope.cursor.forceUnicode {
		panic("unbalanced V8 ForceUnicodeScope")
	}
	scope.cursor.forceUnicode = false
	scope.active = false
}

func (scope *v8ASTForceUnicodeScope) with(body func()) {
	if body == nil {
		panic("V8 ForceUnicodeScope requires a body")
	}
	defer scope.close()
	body()
}

func (cursor *v8ASTParserCursor) errorf(format string, args ...any) error {
	return fmt.Errorf("regular expression at %d: %s", cursor.position(), fmt.Sprintf(format, args...))
}

func v8ASTIsDecimalDigit(character rune) bool { return character >= '0' && character <= '9' }

// parseIntervalQuantifier ports ParserImpl::ParseIntervalQuantifier. Decimal
// overflow saturates at Tree::kInfinity while the cursor still consumes all
// digits, matching the source's recovery position.
func (cursor *v8ASTParserCursor) parseIntervalQuantifier() (minimum, maximum int, ok bool) {
	if cursor.current != '{' {
		panic("V8 interval quantifier parser must start at '{'")
	}
	start := cursor.position()
	cursor.advance()
	if !v8ASTIsDecimalDigit(cursor.current) {
		cursor.reset(start)
		return 0, 0, false
	}
	minimum = 0
	for v8ASTIsDecimalDigit(cursor.current) {
		digit := int(cursor.current - '0')
		if minimum > (v8ASTInfinity-digit)/10 {
			for v8ASTIsDecimalDigit(cursor.current) {
				cursor.advance()
			}
			minimum = v8ASTInfinity
			break
		}
		minimum = 10*minimum + digit
		cursor.advance()
	}
	if cursor.current == '}' {
		maximum = minimum
		cursor.advance()
		return minimum, maximum, true
	}
	if cursor.current != ',' {
		cursor.reset(start)
		return 0, 0, false
	}
	cursor.advance()
	if cursor.current == '}' {
		cursor.advance()
		return minimum, v8ASTInfinity, true
	}
	maximum = 0
	for v8ASTIsDecimalDigit(cursor.current) {
		digit := int(cursor.current - '0')
		if maximum > (v8ASTInfinity-digit)/10 {
			for v8ASTIsDecimalDigit(cursor.current) {
				cursor.advance()
			}
			maximum = v8ASTInfinity
			break
		}
		maximum = 10*maximum + digit
		cursor.advance()
	}
	if cursor.current != '}' {
		cursor.reset(start)
		return 0, 0, false
	}
	cursor.advance()
	return minimum, maximum, true
}

func v8ASTHexValue(character rune) int {
	switch {
	case character >= '0' && character <= '9':
		return int(character - '0')
	case character >= 'a' && character <= 'f':
		return int(character-'a') + 10
	case character >= 'A' && character <= 'F':
		return int(character-'A') + 10
	default:
		return -1
	}
}

// parseHexEscape ports the fixed-width hexadecimal reader. On failure, the
// cursor returns to the first digit so the caller can apply legacy identity
// escape or Unicode-mode error rules.
func (cursor *v8ASTParserCursor) parseHexEscape(length int) (rune, bool) {
	if length < 0 {
		panic("negative V8 hex escape width")
	}
	start := cursor.position()
	value := rune(0)
	for index := 0; index < length; index++ {
		digit := v8ASTHexValue(cursor.current)
		if digit < 0 {
			cursor.reset(start)
			return 0, false
		}
		value = value*16 + rune(digit)
		cursor.advance()
	}
	return value, true
}

func (cursor *v8ASTParserCursor) parseUnlimitedLengthHexNumber(maximum rune) (rune, bool) {
	value := rune(0)
	digit := v8ASTHexValue(cursor.current)
	if digit < 0 {
		return 0, false
	}
	for digit >= 0 {
		value = value*16 + rune(digit)
		if value > maximum {
			return 0, false
		}
		cursor.advance()
		digit = v8ASTHexValue(cursor.current)
	}
	return value, true
}

// parseUnicodeEscape expects the backslash and 'u' to have already been
// consumed. It combines a fixed-width escaped surrogate pair only in Unicode
// mode and rolls back malformed braced forms.
func (cursor *v8ASTParserCursor) parseUnicodeEscape() (rune, bool) {
	if cursor.current == '{' && cursor.isUnicodeMode() {
		start := cursor.position()
		cursor.advance()
		value, ok := cursor.parseUnlimitedLengthHexNumber(v8ASTMaxCodePoint)
		if ok && cursor.current == '}' {
			cursor.advance()
			return value, true
		}
		cursor.reset(start)
		return 0, false
	}
	value, ok := cursor.parseHexEscape(4)
	if ok && cursor.isUnicodeMode() && value >= 0xd800 && value <= 0xdbff && cursor.current == '\\' {
		start := cursor.position()
		if cursor.next() == 'u' {
			cursor.advanceBy(2)
			trail, trailOK := cursor.parseHexEscape(4)
			if trailOK && trail >= 0xdc00 && trail <= 0xdfff {
				combined := 0x10000 + ((value - 0xd800) << 10) + trail - 0xdc00
				return combined, true
			}
		}
		cursor.reset(start)
	}
	return value, ok
}

func (cursor *v8ASTParserCursor) parseOctalLiteral() rune {
	if !v8ASTIsDecimalDigit(cursor.current) || cursor.current > '7' {
		panic("V8 octal parser must start at an octal digit")
	}
	value := cursor.current - '0'
	cursor.advance()
	if cursor.current >= '0' && cursor.current <= '7' {
		value = value*8 + cursor.current - '0'
		cursor.advance()
		if value < 32 && cursor.current >= '0' && cursor.current <= '7' {
			value = value*8 + cursor.current - '0'
			cursor.advance()
		}
	}
	return value
}

// parseCharacterEscape ports ParserImpl::ParseCharacterEscape. The callback
// performs V8's lazy capture scan for legacy `\k` identity-escape handling.
func (cursor *v8ASTParserCursor) parseCharacterEscape(state v8ASTInClassEscapeState, hasNamedCaptures func(v8ASTInClassEscapeState) bool) rune {
	if cursor.current != '\\' || !cursor.hasNextInput() {
		cursor.reportError(errorString(errorEscapeAtEndOfPattern))
		return 0
	}
	cursor.advance()
	character := cursor.current
	switch character {
	case 'f':
		cursor.advance()
		return '\f'
	case 'n':
		cursor.advance()
		return '\n'
	case 'r':
		cursor.advance()
		return '\r'
	case 't':
		cursor.advance()
		return '\t'
	case 'v':
		cursor.advance()
		return '\v'
	case 'c':
		controlLetter := cursor.next()
		letter := controlLetter &^ rune('A'^'a')
		if letter >= 'A' && letter <= 'Z' {
			cursor.advanceBy(2)
			return controlLetter & 0x1f
		}
		if cursor.isUnicodeMode() {
			cursor.reportError(errorString(errorInvalidUnicodeEscape))
			return 0
		}
		if state == v8ASTInClass && ((controlLetter >= '0' && controlLetter <= '9') || controlLetter == '_') {
			cursor.advanceBy(2)
			return controlLetter & 0x1f
		}
		return '\\'
	case '0':
		next := cursor.next()
		if next < '0' || next > '9' {
			cursor.advance()
			return 0
		}
		fallthrough
	case '1', '2', '3', '4', '5', '6', '7':
		if cursor.isUnicodeMode() {
			cursor.reportError(errorString(errorInvalidDecimalEscape))
			return 0
		}
		return cursor.parseOctalLiteral()
	case 'x':
		cursor.advance()
		value, ok := cursor.parseHexEscape(2)
		if ok {
			return value
		}
		if cursor.isUnicodeMode() {
			cursor.reportError(errorString(errorInvalidEscape))
			return 0
		}
		return 'x'
	case 'u':
		cursor.advance()
		value, ok := cursor.parseUnicodeEscape()
		if ok {
			return value
		}
		if cursor.isUnicodeMode() {
			cursor.reportError(errorString(errorInvalidUnicodeEscape))
			return 0
		}
		return 'u'
	}
	if cursor.flags.UnicodeSets && state == v8ASTInClass && cursor.isClassSetReservedPunctuator(character) {
		cursor.advance()
		return character
	}
	if cursor.isUnicodeMode() {
		if !cursor.isSyntaxCharacterOrSlash(character) && !(state == v8ASTInClass && character == '-') {
			cursor.reportError(errorString(errorInvalidEscape))
			return 0
		}
		cursor.advance()
		return character
	}
	if character == 'c' {
		cursor.reportError(errorString(errorInvalidEscape))
		return 0
	}
	cursor.advance()
	if character == 'k' && hasNamedCaptures != nil && hasNamedCaptures(state) {
		cursor.reportError(errorString(errorInvalidEscape))
		return 0
	}
	return character
}

const v8ASTMaxCaptureGroups = 32767

type v8ASTCaptureScanState struct {
	capturesStarted int
	captureCount    int
	hasNamed        bool
	scanned         bool
}

// scanForCaptures is the parser's lazy mini-parser. It skips escaped code
// points and whole character classes, including nested /v classes, while
// counting capture groups and detecting named-capture syntax.
func (cursor *v8ASTParserCursor) scanForCaptures(state *v8ASTCaptureScanState, escapeState v8ASTInClassEscapeState) {
	if state == nil || state.scanned {
		panic("V8 capture scan requires an unscanned state")
	}
	savedPosition := cursor.position()
	captureCount := state.capturesStarted
	if escapeState == v8ASTInClass {
		if cursor.isUnicodeMode() {
			panic("V8 capture scanner must not start inside a Unicode class")
		}
		for cursor.current != v8ASTParserEndMarker {
			character := cursor.current
			cursor.advance()
			if character == '\\' {
				cursor.advance()
			} else if character == ']' {
				break
			}
		}
	}
	for cursor.current != v8ASTParserEndMarker {
		character := cursor.current
		cursor.advance()
		switch character {
		case '\\':
			cursor.advance()
		case '[':
			nesting := 0
			for cursor.current != v8ASTParserEndMarker {
				classCharacter := cursor.current
				cursor.advance()
				switch classCharacter {
				case '\\':
					cursor.advance()
				case '[':
					if cursor.flags.UnicodeSets {
						nesting++
					}
				case ']':
					if nesting == 0 {
						classCharacter = v8ASTParserEndMarker
					} else {
						nesting--
					}
				}
				if classCharacter == v8ASTParserEndMarker {
					break
				}
			}
		case '(':
			if cursor.current == '?' {
				cursor.advance()
				if cursor.current != '<' {
					break
				}
				cursor.advance()
				if cursor.current == '=' || cursor.current == '!' {
					break
				}
				state.hasNamed = true
			}
			captureCount++
		}
	}
	state.captureCount = captureCount
	state.scanned = true
	cursor.reset(savedPosition)
}

func (cursor *v8ASTParserCursor) hasNamedCaptures(state *v8ASTCaptureScanState, escapeState v8ASTInClassEscapeState) bool {
	if state == nil {
		panic("V8 named-capture query requires scan state")
	}
	if state.scanned {
		return state.hasNamed
	}
	cursor.scanForCaptures(state, escapeState)
	return state.hasNamed
}

// parseBackReferenceIndex performs V8's forward-capture scan only when the
// apparent decimal reference exceeds captures already parsed.
func (cursor *v8ASTParserCursor) parseBackReferenceIndex(state *v8ASTCaptureScanState) (int, bool) {
	if cursor.current != '\\' || cursor.next() < '1' || cursor.next() > '9' {
		panic("V8 backreference parser must start at backslash followed by 1-9")
	}
	start := cursor.position()
	value := int(cursor.next() - '0')
	cursor.advanceBy(2)
	for v8ASTIsDecimalDigit(cursor.current) {
		digit := int(cursor.current - '0')
		if value > (v8ASTMaxCaptureGroups-digit)/10 {
			cursor.reset(start)
			return 0, false
		}
		value = value*10 + digit
		if value > v8ASTMaxCaptureGroups {
			cursor.reset(start)
			return 0, false
		}
		cursor.advance()
	}
	if value > state.capturesStarted {
		if !state.scanned {
			cursor.scanForCaptures(state, v8ASTNotInClass)
		}
		if value > state.captureCount {
			cursor.reset(start)
			return 0, false
		}
	}
	return value, true
}

// parseCaptureGroupName ports ParseCaptureGroupName and
// ParseCaptureNameCharacter. It is entered immediately after the opening '<'
// and leaves the cursor immediately after the closing '>'. Names are retained
// as UTF-16 because that is the parser's source representation.
func (cursor *v8ASTParserCursor) parseCaptureGroupName() ([]uint16, bool) {
	// V8 temporarily forces Unicode decoding for RegExpIdentifierName, even
	// when the enclosing pattern is in legacy mode. Re-read the current code
	// unit under that mode, then leave the cursor on '>' until normal decoding
	// has been restored for the character after the name.
	nameStart := cursor.position()
	previousForceUnicode := cursor.forceUnicode
	cursor.forceUnicode = true
	cursor.reset(nameStart)
	name := make([]rune, 0, 16)
	atStart := true
	for cursor.current != '>' && cursor.current != v8ASTParserEndMarker {
		characterPosition := cursor.position()
		character, ok := cursor.parseCaptureNameCodePoint()
		if !ok {
			cursor.forceUnicode = previousForceUnicode
			return nil, false
		}
		if atStart {
			if !isIdentifierStart(character) {
				cursor.reportErrorAt(errorString(errorInvalidCaptureGroupName), characterPosition)
				cursor.forceUnicode = previousForceUnicode
				return nil, false
			}
			atStart = false
		} else if !isIdentifierContinue(character) {
			cursor.reportErrorAt(errorString(errorInvalidCaptureGroupName), characterPosition)
			cursor.forceUnicode = previousForceUnicode
			return nil, false
		}
		name = append(name, character)
	}
	if atStart || cursor.current != '>' {
		cursor.reportError(errorString(errorInvalidCaptureGroupName))
		cursor.forceUnicode = previousForceUnicode
		return nil, false
	}
	cursor.forceUnicode = previousForceUnicode
	cursor.advance()
	return utf16.Encode(name), true
}

// parseCaptureNameEscape consumes a Unicode escape in a capture identifier.
// Capture names admit braced Unicode escapes independently of the pattern's
// legacy/u mode; malformed escapes retain ParseUnicodeEscape's reset point.
func (cursor *v8ASTParserCursor) parseCaptureNameEscape() (rune, string, bool) {
	if cursor.current != '\\' {
		panic("V8 capture-name escape parser must start at backslash")
	}
	start := cursor.position()
	cursor.advance()
	if cursor.current != 'u' {
		cursor.reset(start)
		return 0, errorString(errorInvalidCaptureGroupName), false
	}
	cursor.advance()
	// parseUnicodeEscape provides the same transactional recovery as V8:
	// malformed fixed-width escapes reset to their first hex digit, and bad
	// braced escapes reset to the opening brace. Since capture names force
	// Unicode mode, this also accepts \u{...} and combines escaped surrogates.
	value, ok := cursor.parseUnicodeEscape()
	if !ok {
		return 0, errorString(errorInvalidUnicodeEscape), false
	}
	return value, "", true
}

func (cursor *v8ASTParserCursor) parseCaptureNameCodePoint() (rune, bool) {
	start := cursor.position()
	var character rune
	var ok bool
	if cursor.current == '\\' {
		var errorCode string
		character, errorCode, ok = cursor.parseCaptureNameEscape()
		if !ok {
			cursor.reportError(errorCode)
			return 0, false
		}
	} else {
		character = cursor.current
		if character == v8ASTParserEndMarker {
			return 0, false
		}
		cursor.advance()
	}
	if character > v8ASTMaxCodePoint || character < 0 {
		cursor.reset(start)
		cursor.reportError(errorString(errorInvalidUnicodeEscape))
		return 0, false
	}
	return character, true
}

type v8ASTCaptureBranch struct {
	id     int
	branch int
}

type v8ASTNamedCapture struct {
	capture *v8ASTCapture
	index   int
	path    []v8ASTCaptureBranch
}

type v8ASTParserCaptureTable struct {
	captures    []*v8ASTCapture
	named       map[string][]v8ASTNamedCapture
	backrefs    []*v8ASTBackReference
	namedRefs   [][]uint16
	numericRefs []v8ASTNumericBackReference
	captureSize int
}

type v8ASTNumericBackReference struct {
	index     int
	reference *v8ASTBackReference
}

func v8ASTNewParserCaptureTable(captureCount int) *v8ASTParserCaptureTable {
	if captureCount < 0 || captureCount > v8ASTMaxCaptureGroups {
		panic("invalid V8 parser capture count")
	}
	return &v8ASTParserCaptureTable{
		captures:    make([]*v8ASTCapture, captureCount+1),
		named:       make(map[string][]v8ASTNamedCapture),
		captureSize: captureCount,
	}
}

func (table *v8ASTParserCaptureTable) captureCount() int {
	if table == nil {
		return 0
	}
	return table.captureSize
}

func v8ASTCapturePathFromParserPath(path []choiceBranch) []v8ASTCaptureBranch {
	result := make([]v8ASTCaptureBranch, len(path))
	for index, branch := range path {
		result[index] = v8ASTCaptureBranch{id: branch.id, branch: branch.branch}
	}
	return result
}

func v8ASTCaptureNameKey(name []uint16) string {
	key := make([]byte, 2*len(name))
	for index, unit := range name {
		key[2*index] = byte(unit)
		key[2*index+1] = byte(unit >> 8)
	}
	return string(key)
}

// v8ASTCaptureNameLess mirrors CaptureNameLess by comparing UTF-16 code units,
// including the surrogate ordering used for non-BMP identifier characters.
func v8ASTCaptureNameLess(left, right []uint16) bool {
	limit := len(left)
	if len(right) < limit {
		limit = len(right)
	}
	for index := 0; index < limit; index++ {
		if left[index] != right[index] {
			return left[index] < right[index]
		}
	}
	return len(left) < len(right)
}

func v8ASTCapturePathsAreExclusive(left, right []v8ASTCaptureBranch) bool {
	for _, lhs := range left {
		for _, rhs := range right {
			if lhs.id == rhs.id && lhs.branch != rhs.branch {
				return true
			}
		}
	}
	return false
}

func (table *v8ASTParserCaptureTable) createCaptureAtIndex(index int, name []uint16, path []v8ASTCaptureBranch) (*v8ASTCapture, error) {
	if table == nil || index < 1 || index > v8ASTMaxCaptureGroups {
		return nil, fmt.Errorf("%s", errorString(errorTooManyCaptures))
	}
	if index >= len(table.captures) {
		captures := make([]*v8ASTCapture, index+1)
		copy(captures, table.captures)
		table.captures = captures
	}
	if table.captures[index] != nil {
		return nil, fmt.Errorf("V8 parser capture index %d was registered twice", index)
	}
	capture := (v8ASTFactory{}).newCapture(index)
	if name != nil {
		capture.name = append([]uint16(nil), name...)
		if err := table.registerNamedCapture(index, name, path, capture); err != nil {
			return nil, err
		}
	}
	table.captures[index] = capture
	if index > table.captureSize {
		table.captureSize = index
	}
	return capture, nil
}

func (table *v8ASTParserCaptureTable) registerNamedCapture(index int, name []uint16, path []v8ASTCaptureBranch, capture *v8ASTCapture) error {
	if table == nil || index < 1 || len(name) == 0 || !validCaptureName(string(utf16.Decode(name))) {
		return fmt.Errorf("%s", errorString(errorInvalidCaptureGroupName))
	}
	key := v8ASTCaptureNameKey(name)
	for _, existing := range table.named[key] {
		if !v8ASTCapturePathsAreExclusive(existing.path, path) {
			return fmt.Errorf("%s", errorString(errorDuplicateCaptureGroupName))
		}
	}
	table.named[key] = append(table.named[key], v8ASTNamedCapture{
		capture: capture,
		index:   index,
		path:    append([]v8ASTCaptureBranch(nil), path...),
	})
	return nil
}

func (table *v8ASTParserCaptureTable) getCapture(index int) *v8ASTCapture {
	if table == nil || index < 0 || index >= len(table.captures) {
		return nil
	}
	return table.captures[index]
}

func (table *v8ASTParserCaptureTable) getNamedCaptures(name []uint16) []*v8ASTCapture {
	if table == nil {
		return nil
	}
	entries := table.named[v8ASTCaptureNameKey(name)]
	captures := make([]*v8ASTCapture, 0, len(entries))
	for _, entry := range entries {
		captures = append(captures, entry.capture)
	}
	return captures
}

// parseNamedBackReferenceName starts at '<' after the parser has consumed
// `\\k`. The syntax-only path uses this to validate without allocating AST.
func (cursor *v8ASTParserCursor) parseNamedBackReferenceName() ([]uint16, bool) {
	if cursor.current != '<' {
		panic("V8 named-backreference parser must start at '<'")
	}
	cursor.advance()
	return cursor.parseCaptureGroupName()
}

func (cursor *v8ASTParserCursor) parseNamedBackReference() (*v8ASTBackReference, bool) {
	name, ok := cursor.parseNamedBackReferenceName()
	if !ok {
		return nil, false
	}
	reference := (v8ASTFactory{}).newBackReference()
	reference.setName(name)
	return reference, true
}

func (table *v8ASTParserCaptureTable) addNamedBackReference(reference *v8ASTBackReference) {
	if table == nil || reference == nil || len(v8ASTBackReferenceName(reference)) == 0 {
		panic("V8 named-reference registry requires a named backreference")
	}
	table.backrefs = append(table.backrefs, reference)
	table.namedRefs = append(table.namedRefs, append([]uint16(nil), v8ASTBackReferenceName(reference)...))
}

func (table *v8ASTParserCaptureTable) addNamedReferenceName(name []uint16) {
	if table == nil || len(name) == 0 {
		panic("V8 named-reference registry requires a name")
	}
	table.namedRefs = append(table.namedRefs, append([]uint16(nil), name...))
}

func (table *v8ASTParserCaptureTable) newNumericBackReference(index int) (*v8ASTBackReference, error) {
	if table == nil || index < 1 || index > v8ASTMaxCaptureGroups {
		return nil, fmt.Errorf("%s", errorString(errorInvalidDecimalEscape))
	}
	reference := (v8ASTFactory{}).newBackReference()
	if capture := table.getCapture(index); capture != nil {
		capture.isBackreferenced = true
		reference.addCapture(capture)
	} else {
		table.numericRefs = append(table.numericRefs, v8ASTNumericBackReference{index: index, reference: reference})
	}
	return reference, nil
}

// patchNamedBackReferences resolves forward references after all captures have
// been registered. Each possible duplicate-name capture remains attached so
// the matcher can select the participating capture at execution time.
func (table *v8ASTParserCaptureTable) patchNamedBackReferences(hasNamedCaptures, unicodeMode bool) error {
	if table == nil {
		return fmt.Errorf("nil V8 parser capture table")
	}
	for _, name := range table.namedRefs {
		entries := table.named[v8ASTCaptureNameKey(name)]
		if len(entries) == 0 && (hasNamedCaptures || unicodeMode) {
			return fmt.Errorf("%s", errorString(errorInvalidNamedCaptureReference))
		}
	}
	for _, reference := range table.backrefs {
		reference.captures = reference.captures[:0]
		entries := table.named[v8ASTCaptureNameKey(v8ASTBackReferenceName(reference))]
		for _, entry := range entries {
			if entry.capture == nil {
				continue
			}
			entry.capture.isBackreferenced = true
			reference.addCapture(entry.capture)
		}
	}
	for _, pending := range table.numericRefs {
		capture := table.getCapture(pending.index)
		if capture == nil {
			return fmt.Errorf("%s", errorString(errorInvalidDecimalEscape))
		}
		capture.isBackreferenced = true
		pending.reference.addCapture(capture)
	}
	table.numericRefs = table.numericRefs[:0]
	return nil
}

type v8ASTParserBuilderAPI interface {
	addCharacter(uint16)
	addUnicodeCharacter(rune)
	addEmpty()
	addClassRanges(*v8ASTClassRanges)
	addAtom(any)
	addAssertion(any)
	newAlternative()
	toRegExp() any
	addQuantifier(int, int, int, v8ASTQuantifierType) bool
}

// v8ASTParser is the recursive-descent driver around the source-shaped cursor,
// builders, parser states, capture table, and character-class productions.
// Its AST is lowered by the V8-shaped graph builder in regexp_compiler.go.
type v8ASTParser struct {
	cursor                *v8ASTParserCursor
	flags                 Flags
	mode                  v8ASTParseMode
	captureScan           *v8ASTCaptureScanState
	captures              *v8ASTParserCaptureTable
	state                 *v8ASTParserState
	choicePath            []choiceBranch
	nextChoiceID          int
	nextLookID            int
	quantifierID          int
	nesting               int
	simple                bool
	containsAnchor        bool
	bufferBoundaries      bool
	possessiveQuantifiers bool
}

// v8ASTParserOptions captures V8 runtime feature flags that affect only
// internal RegExp syntax. Both flags default off, matching V8's defaults.
type v8ASTParserOptions struct {
	bufferBoundaries      bool
	possessiveQuantifiers bool
}

func v8ASTNewParser(input []uint16, flags Flags, oneByte bool, mode v8ASTParseMode, stackOverflow func() bool) *v8ASTParser {
	return v8ASTNewParserWithOptions(input, flags, oneByte, mode, stackOverflow, v8ASTParserOptions{})
}

func v8ASTNewParserWithOptions(input []uint16, flags Flags, oneByte bool, mode v8ASTParseMode, stackOverflow func() bool, options v8ASTParserOptions) *v8ASTParser {
	return v8ASTNewParserFromCursor(v8ASTNewParserCursor(input, flags, oneByte, stackOverflow), flags, mode, options)
}

func v8ASTNewParserOneByte(input []byte, flags Flags, mode v8ASTParseMode, stackOverflow func() bool, options v8ASTParserOptions) *v8ASTParser {
	return v8ASTNewParserFromCursor(v8ASTNewParserCursorOneByte(input, flags, stackOverflow), flags, mode, options)
}

func v8ASTNewParserFromCursor(cursor *v8ASTParserCursor, flags Flags, mode v8ASTParseMode, options v8ASTParserOptions) *v8ASTParser {
	return &v8ASTParser{
		cursor:                cursor,
		flags:                 flags,
		mode:                  mode,
		captureScan:           &v8ASTCaptureScanState{},
		bufferBoundaries:      options.bufferBoundaries,
		possessiveQuantifiers: options.possessiveQuantifiers,
	}
}

func (parser *v8ASTParser) newBuilder() v8ASTParserBuilderAPI {
	if parser.mode == v8ASTBuildAST {
		return v8ASTNewParserBuilder(parser.flags)
	}
	return v8ASTNewSyntaxBuilder(parser.flags)
}

func (parser *v8ASTParser) parse() (any, error) {
	result, err := parser.parseCompileData()
	return result.Tree, err
}

// v8ASTParserCompileData mirrors the semantic fields filled by V8's
// ParserImpl::Parse for AST construction. Syntax-only mode leaves AST metadata
// empty while still returning the same diagnostics.
type v8ASTParserCompileData struct {
	Tree           any
	Simple         bool
	ContainsAnchor bool
	CaptureCount   int
	NamedCaptures  []*v8ASTCapture
	Error          string
	ErrorPosition  int
}

func (parser *v8ASTParser) parserCompileData(tree any, err error) v8ASTParserCompileData {
	result := v8ASTParserCompileData{Tree: tree}
	if parser.mode == v8ASTBuildAST {
		result.ContainsAnchor = parser.containsAnchor
		result.CaptureCount = parser.captureScan.capturesStarted
		result.Simple = parser.simple && result.CaptureCount == 0
		if parser.captures != nil {
			names := make([][]uint16, 0, len(parser.captures.named))
			for _, entries := range parser.captures.named {
				if len(entries) != 0 {
					names = append(names, entries[0].capture.name)
				}
			}
			sort.Slice(names, func(left, right int) bool {
				return v8ASTCaptureNameLess(names[left], names[right])
			})
			for _, name := range names {
				for _, entry := range parser.captures.named[v8ASTCaptureNameKey(name)] {
					result.NamedCaptures = append(result.NamedCaptures, entry.capture)
				}
			}
		}
	}
	if err != nil {
		result.Error = parser.cursor.errorCode
		result.ErrorPosition = parser.cursor.errorPosition
		if result.Error == "" {
			result.Error = err.Error()
			result.ErrorPosition = parser.cursor.position()
		}
	}
	return result
}

func (parser *v8ASTParser) parseCompileData() (v8ASTParserCompileData, error) {
	// V8 counts captures lazily: ordinary patterns do not pay for a pre-scan.
	// A numeric forward reference or legacy \k escape triggers the mini-parser.
	parser.state = v8ASTNewParserState(nil, v8ASTInitial, v8ASTLookahead, 0, nil, parser.flags, parser.mode == v8ASTBuildAST)
	root, err := parser.parseDisjunction(parser.state)
	if err != nil {
		return parser.parserCompileData(nil, err), err
	}
	if parser.cursor.failedParse() {
		err = parser.parseError(parser.cursor.errorCode)
		return parser.parserCompileData(nil, err), err
	}
	if parser.cursor.current != v8ASTParserEndMarker {
		if parser.cursor.current == ')' {
			err = parser.parseError(errorString(errorUnmatchedParen))
		} else {
			err = parser.parseError("unexpected-token")
		}
		return parser.parserCompileData(nil, err), err
	}
	if atom, ok := root.(*v8ASTAtom); ok && len(atom.data) == parser.cursor.inputLength() {
		parser.simple = true
	}
	if parser.captures != nil {
		if err := parser.captures.patchNamedBackReferences(parser.captureScan.hasNamed, parser.cursor.isUnicodeMode()); err != nil {
			err = parser.parseError(err.Error())
			return parser.parserCompileData(nil, err), err
		}
	}
	return parser.parserCompileData(root, nil), nil
}

func (parser *v8ASTParser) captureTable() *v8ASTParserCaptureTable {
	if parser.captures == nil {
		parser.captures = v8ASTNewParserCaptureTable(0)
	}
	return parser.captures
}

func (parser *v8ASTParser) parseError(message string) error {
	position := parser.cursor.position()
	if parser.cursor.failed {
		position = parser.cursor.errorPosition
		message = parser.cursor.errorCode
	} else {
		parser.cursor.reportError(message)
	}
	return fmt.Errorf("regular expression at %d: %s", position, message)
}

func (parser *v8ASTParser) parseDisjunction(state *v8ASTParserState) (any, error) {
	if state == nil {
		return nil, fmt.Errorf("nil V8 parser state")
	}
	choiceID := parser.nextChoiceID
	parser.nextChoiceID++
	parser.choicePath = append(parser.choicePath, choiceBranch{id: choiceID})
	defer func() { parser.choicePath = parser.choicePath[:len(parser.choicePath)-1] }()
	for {
		if _, err := parser.parseAlternative(state); err != nil {
			return nil, err
		}
		if parser.cursor.current != '|' {
			break
		}
		parser.cursor.advance()
		last := len(parser.choicePath) - 1
		parser.choicePath[last].branch++
		state.newAlternative(parser.captureScan.capturesStarted)
		if state.buildAST {
			state.buildBuilder.newAlternative()
		} else {
			state.syntaxBuilder.newAlternative()
		}
	}
	if state.buildAST {
		return state.buildBuilder.toRegExp(), nil
	}
	return nil, nil
}

func (parser *v8ASTParser) parseAlternative(state *v8ASTParserState) (bool, error) {
	if parser.cursor.failedParse() {
		return false, parser.parseError(parser.cursor.errorCode)
	}
	parsed := false
	for parser.cursor.current != v8ASTParserEndMarker && parser.cursor.current != ')' && parser.cursor.current != '|' {
		if parser.cursor.failedParse() {
			return parsed, parser.parseError("parser-cursor-failed")
		}
		lookaround, kind, err := parser.parseAtom(state)
		if err != nil {
			return parsed, err
		}
		parsed = true
		quantifierIndex := parser.quantifierID
		if err := parser.parseQuantifier(state, quantifierIndex); err != nil {
			return parsed, err
		}
		if lookaround && state.buildAST {
			// Builder::AddAtom records the AST node; the syntax-only builder
			// stores this bit separately for AddQuantifierToAtom validation.
			_ = kind
		}
	}
	if !parsed {
		parser.addEmpty(state)
	}
	return parsed, nil
}

func (parser *v8ASTParser) addCharacter(state *v8ASTParserState, character rune) {
	if !state.buildAST {
		state.syntaxBuilder.addUnicodeCharacter(character)
		return
	}
	if parser.cursor.isUnicodeMode() || character > 0xffff {
		state.buildBuilder.addUnicodeCharacter(character)
		return
	}
	state.buildBuilder.addCharacter(uint16(character))
}

func (parser *v8ASTParser) addEmpty(state *v8ASTParserState) {
	if state.buildAST {
		state.buildBuilder.addEmpty()
		return
	}
	state.syntaxBuilder.addEmpty()
}

func (parser *v8ASTParser) addAtom(state *v8ASTParserState, atom any, lookaround *v8ASTLookaroundType) {
	if state.buildAST {
		if _, isClassSetExpression := atom.(*v8ASTClassSetExpression); isClassSetExpression {
			state.buildBuilder.addTerm(atom)
			return
		}
		state.buildBuilder.addAtom(atom)
		return
	}
	if lookaround != nil {
		state.syntaxBuilder.addLookaround(*lookaround)
		return
	}
	state.syntaxBuilder.addAtom(atom)
}

func (parser *v8ASTParser) addTerm(state *v8ASTParserState, term any) {
	if state.buildAST {
		state.buildBuilder.addTerm(term)
		return
	}
	state.syntaxBuilder.addTerm(term)
}

func (parser *v8ASTParser) addClassRanges(state *v8ASTParserState, class *v8ASTClassRanges) {
	if state.buildAST {
		state.buildBuilder.addClassRanges(class)
		return
	}
	state.syntaxBuilder.addClassRanges(class)
}

func (parser *v8ASTParser) addAssertion(state *v8ASTParserState, assertion *v8ASTAssertion) {
	if state.buildAST {
		state.buildBuilder.addAssertion(assertion)
		return
	}
	state.syntaxBuilder.addAssertion(assertion)
}

func (parser *v8ASTParser) parseAtom(state *v8ASTParserState) (bool, v8ASTLookaroundType, error) {
	cursor := parser.cursor
	character := cursor.current
	switch character {
	case '*', '+', '?':
		return false, 0, parser.parseError(errorString(errorNothingToRepeat))
	case '(':
		lookaround, kind, err := parser.parseOpenParenthesis(state)
		return lookaround, kind, err
	case ')', '|', v8ASTParserEndMarker:
		return false, 0, parser.parseError("missing-atom")
	case '[':
		return parser.parseCharacterClass(state)
	case '.':
		cursor.advance()
		set := v8ASTNotLineTerminator
		if parser.flags.DotAll {
			set = v8ASTEverything
		}
		if state.buildAST {
			parser.addClassRanges(state, (v8ASTFactory{}).newClassRangesFromStandardSet(set))
		} else {
			parser.addClassRanges(state, nil)
		}
		return false, 0, nil
	case '^':
		cursor.advance()
		parser.containsAnchor = true
		assertion := v8ASTStartOfInput
		if parser.flags.Multiline {
			assertion = v8ASTStartOfLine
		}
		if state.buildAST {
			parser.addAssertion(state, (v8ASTFactory{}).newAssertion(assertion))
		} else {
			parser.addAssertion(state, nil)
		}
		return false, 0, nil
	case '{':
		if _, _, ok := cursor.parseIntervalQuantifier(); ok {
			return false, 0, parser.parseError(errorString(errorNothingToRepeat))
		}
		if cursor.isUnicodeMode() {
			return false, 0, parser.parseError(errorString(errorLoneQuantifierBrackets))
		}
		cursor.advance()
		parser.addCharacter(state, '{')
		return false, 0, nil
	case '}', ']':
		if cursor.isUnicodeMode() {
			return false, 0, parser.parseError(errorString(errorLoneQuantifierBrackets))
		}
		cursor.advance()
		parser.addCharacter(state, character)
		return false, 0, nil
	case '$':
		cursor.advance()
		assertion := v8ASTEndOfInput
		if parser.flags.Multiline {
			assertion = v8ASTEndOfLine
		}
		if state.buildAST {
			parser.addAssertion(state, (v8ASTFactory{}).newAssertion(assertion))
		} else {
			parser.addAssertion(state, nil)
		}
		return false, 0, nil
	case '\\':
		return parser.parseEscapeAtom(state)
	default:
		cursor.advance()
		parser.addCharacter(state, character)
		return false, 0, nil
	}
}

func (parser *v8ASTParser) parseEscapeAtom(state *v8ASTParserState) (bool, v8ASTLookaroundType, error) {
	cursor := parser.cursor
	if !cursor.hasNextInput() {
		return false, 0, parser.parseError(errorString(errorEscapeAtEndOfPattern))
	}
	next := cursor.next()
	if parser.bufferBoundaries && cursor.isUnicodeMode() {
		var assertion v8ASTAssertionType
		recognized := true
		switch next {
		case 'A':
			assertion = v8ASTStartOfInput
			parser.containsAnchor = true
		case 'z':
			assertion = v8ASTEndOfInput
		case 'Z':
			assertion = v8ASTEndOfBuffer
		default:
			recognized = false
		}
		if recognized {
			cursor.advanceBy(2)
			if state.buildAST {
				parser.addAssertion(state, (v8ASTFactory{}).newAssertion(assertion))
			} else {
				parser.addAssertion(state, nil)
			}
			return false, 0, nil
		}
	}
	if (next == 'p' || next == 'P') && cursor.isUnicodeMode() {
		cursor.advanceBy(2)
		if cursor.current != '{' {
			return false, 0, parser.parseError(errorString(errorInvalidPropertyName))
		}
		propertyName, err := parser.parseUnicodePropertyName(errorInvalidPropertyName)
		if err != nil {
			return false, 0, err
		}
		property, ok := canonicalUnicodeProperty(propertyName)
		if !ok {
			return false, 0, parser.parseError(errorString(errorInvalidPropertyName))
		}
		if strings.HasPrefix(property, "string-property:") {
			if !parser.flags.UnicodeSets || next == 'P' {
				return false, 0, parser.parseError(errorString(errorInvalidPropertyName))
			}
			data, exists := regexpUnicodeStringProperties[strings.TrimPrefix(property, "string-property:")]
			if !exists {
				return false, 0, parser.parseError(errorString(errorInvalidPropertyName))
			}
			operand := parser.parserClassOperandFromProperty(data)
			if state.buildAST {
				parser.addTerm(state, (v8ASTFactory{}).newClassSetOperand(operand.ranges, operand.strings))
			} else {
				parser.addTerm(state, nil)
			}
			return false, 0, nil
		}
		var ranges []v8ASTCharacterRange
		var rangesOK bool
		if state.buildAST {
			ranges, rangesOK = v8ASTUnicodePropertyRanges(property)
		} else {
			rangesOK = v8ASTUnicodePropertyAvailable(property)
		}
		if !rangesOK {
			return false, 0, parser.parseError(errorString(errorInvalidPropertyName))
		}
		if state.buildAST && parser.flags.UnicodeSets {
			ranges = parser.closeUnicodeSetPropertyRanges(ranges)
		}
		if state.buildAST && next == 'P' {
			ranges = v8ASTNegateRanges(ranges)
		}
		if parser.flags.UnicodeSets {
			if state.buildAST {
				parser.addTerm(state, (v8ASTFactory{}).newClassSetOperand(ranges, nil))
			} else {
				parser.addTerm(state, nil)
			}
		} else {
			if state.buildAST {
				parser.addClassRanges(state, (v8ASTFactory{}).newClassRanges(ranges, 0))
			} else {
				parser.addClassRanges(state, nil)
			}
		}
		return false, 0, nil
	}
	if next == 'd' || next == 'D' || next == 's' || next == 'S' || next == 'w' || next == 'W' {
		cursor.advanceBy(2)
		if state.buildAST {
			class := (v8ASTFactory{}).newClassRangesFromStandardSet(v8ASTStandardCharacterSet(next))
			parser.addClassRanges(state, class)
		} else {
			parser.addClassRanges(state, nil)
		}
		return false, 0, nil
	}
	if next == 'b' || next == 'B' {
		cursor.advanceBy(2)
		assertion := v8ASTBoundary
		if next == 'B' {
			assertion = v8ASTNonBoundary
		}
		if state.buildAST {
			parser.addAssertion(state, (v8ASTFactory{}).newAssertion(assertion))
		} else {
			parser.addAssertion(state, nil)
		}
		return false, 0, nil
	}
	if next >= '1' && next <= '9' {
		if _, ok := cursor.parseBackReferenceIndex(parser.captureScan); ok {
			index := 0
			// The numeric parser leaves the cursor after the complete decimal
			// sequence; recover its value from the consumed UTF-16 source.
			start := cursor.position()
			for start > 0 && cursor.inputAt(start-1) >= '0' && cursor.inputAt(start-1) <= '9' {
				start--
			}
			for position := start; position < cursor.position(); position++ {
				index = index*10 + int(cursor.inputAt(position)-'0')
			}
			if state.isInsideCaptureGroup(index) {
				parser.addEmpty(state)
				return false, 0, nil
			}
			if state.buildAST {
				reference, err := parser.captureTable().newNumericBackReference(index)
				if err != nil {
					return false, 0, parser.parseError(errorString(errorInvalidDecimalEscape))
				}
				parser.addAtom(state, reference, nil)
			} else {
				parser.addAtom(state, nil, nil)
			}
			return false, 0, nil
		}
		if cursor.isUnicodeMode() {
			return false, 0, parser.parseError(errorString(errorInvalidEscape))
		}
	}
	if next == 'k' {
		hasNamed := cursor.hasNamedCaptures(parser.captureScan, v8ASTNotInClass)
		if next2 := cursor.inputAtSafe(cursor.position() + 2); next2 == '<' && (cursor.isUnicodeMode() || hasNamed) {
			cursor.advanceBy(2)
			name, ok := cursor.parseNamedBackReferenceName()
			if !ok {
				return false, 0, parser.parseError(errorString(errorInvalidCaptureGroupName))
			}
			if state.isInsideNamedCaptureGroup(name) {
				parser.addEmpty(state)
				return false, 0, nil
			}
			if state.buildAST {
				reference := (v8ASTFactory{}).newBackReference()
				reference.setName(name)
				parser.captureTable().addNamedBackReference(reference)
				parser.addAtom(state, reference, nil)
			} else {
				parser.captureTable().addNamedReferenceName(name)
				parser.addAtom(state, nil, nil)
			}
			return false, 0, nil
		}
		if cursor.isUnicodeMode() || hasNamed {
			return false, 0, parser.parseError(errorString(errorInvalidNamedReference))
		}
	}
	character := cursor.parseCharacterEscape(v8ASTNotInClass, func(v8ASTInClassEscapeState) bool {
		return parser.captureScan.hasNamed
	})
	if cursor.failedParse() {
		return false, 0, parser.parseError(cursor.errorCode)
	}
	parser.addCharacter(state, character)
	return false, 0, nil
}

func (cursor *v8ASTParserCursor) inputAtSafe(index int) rune {
	if index < 0 || index >= cursor.inputLength() {
		return v8ASTParserEndMarker
	}
	return rune(cursor.inputAt(index))
}

func (parser *v8ASTParser) parseQuantifier(state *v8ASTParserState, index int) error {
	cursor := parser.cursor
	minimum, maximum := 0, 0
	kind := v8ASTGreedy
	switch cursor.current {
	case '*':
		minimum, maximum = 0, v8ASTInfinity
		cursor.advance()
	case '+':
		minimum, maximum = 1, v8ASTInfinity
		cursor.advance()
	case '?':
		minimum, maximum = 0, 1
		cursor.advance()
	case '{':
		var ok bool
		minimum, maximum, ok = cursor.parseIntervalQuantifier()
		if !ok {
			if cursor.isUnicodeMode() {
				return parser.parseError(errorString(errorIncompleteQuantifier))
			}
			return nil
		}
		if maximum < minimum {
			return parser.parseError(errorString(errorRangeOutOfOrder))
		}
	default:
		return nil
	}
	if cursor.current == '?' {
		kind = v8ASTNonGreedy
		cursor.advance()
	} else if parser.possessiveQuantifiers && cursor.current == '+' {
		kind = v8ASTPossessive
		cursor.advance()
	}
	valid := false
	if state.buildAST {
		valid = state.buildBuilder.addQuantifier(minimum, maximum, index, kind)
	} else {
		valid = state.syntaxBuilder.addQuantifier(minimum, maximum, index, kind)
	}
	if !valid {
		return parser.parseError(errorString(errorInvalidQuantifier))
	}
	parser.quantifierID++
	return nil
}

func (parser *v8ASTParser) parseScopedFlags() (Flags, error) {
	flags := parser.flags
	var seen Flags
	dash := false
	for {
		character := parser.cursor.current
		if character == ':' {
			if !seen.IgnoreCase && !seen.Multiline && !seen.DotAll {
				return Flags{}, parser.parseError(errorString(errorInvalidFlagGroup))
			}
			parser.cursor.advance()
			return flags, nil
		}
		if character == '-' {
			if dash {
				return Flags{}, parser.parseError(errorString(errorMultipleFlagDashes))
			}
			dash = true
			parser.cursor.advance()
			continue
		}
		flag, ok := TryFlagFromChar(character)
		if !ok || (flag != FlagIgnoreCase && flag != FlagMultiline && flag != FlagDotAll) {
			return Flags{}, parser.parseError(errorString(errorInvalidGroup))
		}
		if seen.Has(flag) {
			return Flags{}, parser.parseError(errorString(errorRepeatedFlag))
		}
		seen.Set(flag, true)
		flags.Set(flag, !dash)
		parser.cursor.advance()
	}
}

func (parser *v8ASTParser) parseOpenParenthesis(parent *v8ASTParserState) (bool, v8ASTLookaroundType, error) {
	cursor := parser.cursor
	if cursor.current != '(' {
		panic("V8 group parser must start at '('")
	}
	if parser.nesting >= maxParserNesting {
		return false, 0, parser.parseError(errorString(errorStackOverflow))
	}
	parser.nesting++
	defer func() { parser.nesting-- }()
	cursor.advance()
	oldFlags := parser.flags
	flagsChanged := false
	defer func() {
		if flagsChanged {
			parser.flags = oldFlags
			cursor.flags = oldFlags
		}
	}()
	groupType := v8ASTCaptureGroup
	lookaroundType := v8ASTLookahead
	positive := true
	captureName := []uint16(nil)
	if cursor.current == '?' {
		cursor.advance()
		if cursor.current == 'i' || cursor.current == 'm' || cursor.current == 's' || cursor.current == '-' {
			flags, err := parser.parseScopedFlags()
			if err != nil {
				return false, 0, err
			}
			parser.flags = flags
			cursor.flags = flags
			flagsChanged = true
			groupType = v8ASTGrouping
		} else {
			switch cursor.current {
			case ':':
				groupType = v8ASTGrouping
				cursor.advance()
			case '=', '!':
				groupType = v8ASTPositiveLookaround
				positive = cursor.current == '='
				if !positive {
					groupType = v8ASTNegativeLookaround
				}
				cursor.advance()
			case '<':
				cursor.advance()
				if cursor.current == '=' || cursor.current == '!' {
					groupType = v8ASTPositiveLookaround
					positive = cursor.current == '='
					if !positive {
						groupType = v8ASTNegativeLookaround
					}
					lookaroundType = v8ASTLookbehind
					cursor.advance()
				} else {
					name, ok := cursor.parseCaptureGroupName()
					if !ok {
						return false, 0, parser.parseError(errorString(errorInvalidCaptureGroupName))
					}
					captureName = name
					parser.captureScan.hasNamed = true
				}
			default:
				return false, 0, parser.parseError(errorString(errorInvalidGroup))
			}
		}
	}
	captureIndex := parser.captureScan.capturesStarted
	if groupType == v8ASTCaptureGroup {
		if parser.captureScan.capturesStarted >= v8ASTMaxCaptureGroups {
			return false, 0, parser.parseError(errorString(errorTooManyCaptures))
		}
		parser.captureScan.capturesStarted++
		captureIndex = parser.captureScan.capturesStarted
		if captureName != nil {
			parser.captureScan.hasNamed = true
		}
	}
	childState := v8ASTNewParserState(parent, groupType, lookaroundType, captureIndex, captureName, parser.flags, parser.mode == v8ASTBuildAST)
	body, err := parser.parseDisjunction(childState)
	if err != nil {
		return false, 0, err
	}
	if cursor.current != ')' {
		return false, 0, parser.parseError(errorString(errorUnterminatedGroup))
	}
	cursor.advance()
	if groupType == v8ASTCaptureGroup && captureName != nil && parent.isInsideNamedCaptureGroup(captureName) {
		return false, 0, parser.parseError(errorString(errorDuplicateCaptureGroupName))
	}
	if parser.mode == v8ASTVerifySyntax {
		if groupType == v8ASTCaptureGroup && captureName != nil {
			if err := parser.captureTable().registerNamedCapture(captureIndex, captureName, v8ASTCapturePathFromParserPath(parser.choicePath), nil); err != nil {
				return false, 0, parser.parseError(err.Error())
			}
		}
		if groupType == v8ASTPositiveLookaround || groupType == v8ASTNegativeLookaround {
			parser.addAtom(parent, nil, &lookaroundType)
		} else {
			parser.addAtom(parent, nil, nil)
		}
		return groupType == v8ASTPositiveLookaround || groupType == v8ASTNegativeLookaround, lookaroundType, nil
	}
	var atom any = body
	switch groupType {
	case v8ASTCaptureGroup:
		capture, captureErr := parser.captureTable().createCaptureAtIndex(captureIndex, captureName, v8ASTCapturePathFromParserPath(parser.choicePath))
		if captureErr != nil {
			return false, 0, parser.parseError(captureErr.Error())
		}
		capture.setBody(body)
		atom = capture
	case v8ASTPositiveLookaround, v8ASTNegativeLookaround:
		captureFrom := captureIndex
		captureCount := parser.captureScan.capturesStarted - captureIndex
		atom = (v8ASTFactory{}).newLookaround(body, positive, captureCount, captureFrom, lookaroundType, parser.nextLookID)
		parser.nextLookID++
	case v8ASTGrouping:
		atom = (v8ASTFactory{}).newGroup(body, parser.flags)
	}
	parser.addAtom(parent, atom, nil)
	return groupType == v8ASTPositiveLookaround || groupType == v8ASTNegativeLookaround, lookaroundType, nil
}

func (builder *v8ASTSyntaxBuilder) addTerm(any) { builder.lastIsLookaround = false }
func (builder *v8ASTSyntaxBuilder) addAssertion(any) {
	builder.lastIsLookaround = false
}
func (builder *v8ASTSyntaxBuilder) newAlternative() { builder.lastIsLookaround = false }
func (builder *v8ASTSyntaxBuilder) addQuantifier(int, int, int, v8ASTQuantifierType) bool {
	if builder.lastIsLookaround && (IsEitherUnicode(builder.flags) || builder.lastLookaroundKind == v8ASTLookbehind) {
		return false
	}
	builder.lastIsLookaround = false
	return true
}

type v8ASTSubexpressionType int8

const (
	v8ASTInitial v8ASTSubexpressionType = iota
	v8ASTCaptureGroup
	v8ASTPositiveLookaround
	v8ASTNegativeLookaround
	v8ASTGrouping
)

// v8ASTParserState is the linked ParserState stack used while parsing one
// disjunction. It records captures that cannot participate in the current
// alternative, which is required when duplicate named captures are resolved.
type v8ASTParserState struct {
	previous         *v8ASTParserState
	buildAST         bool
	buildBuilder     *v8ASTParserBuilder
	syntaxBuilder    *v8ASTSyntaxBuilder
	groupType        v8ASTSubexpressionType
	lookaroundType   v8ASTLookaroundType
	captureIndex     int
	captureName      []uint16
	nonParticipating []v8ASTInterval
}

func v8ASTNewParserState(previous *v8ASTParserState, groupType v8ASTSubexpressionType, lookaroundType v8ASTLookaroundType, captureIndex int, captureName []uint16, flags Flags, buildAST bool) *v8ASTParserState {
	state := &v8ASTParserState{
		previous: previous, buildAST: buildAST, groupType: groupType,
		lookaroundType: lookaroundType, captureIndex: captureIndex,
		captureName: append([]uint16(nil), captureName...),
	}
	if previous != nil {
		state.nonParticipating = append(state.nonParticipating, previous.nonParticipating...)
	}
	if buildAST {
		state.buildBuilder = v8ASTNewParserBuilder(flags)
	} else {
		state.syntaxBuilder = v8ASTNewSyntaxBuilder(flags)
	}
	return state
}

func (state *v8ASTParserState) isSubexpression() bool {
	return state != nil && state.previous != nil
}

func (state *v8ASTParserState) isNamedCapture() bool {
	return state != nil && state.captureName != nil
}

func (state *v8ASTParserState) isInsideCaptureGroup(index int) bool {
	for current := state; current != nil; current = current.previous {
		if current.groupType != v8ASTCaptureGroup {
			continue
		}
		if index == current.captureIndex {
			return true
		}
		if index > current.captureIndex {
			return false
		}
	}
	return false
}

func v8ASTCaptureNamesEqual(lhs, rhs []uint16) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for index := range lhs {
		if lhs[index] != rhs[index] {
			return false
		}
	}
	return true
}

func (state *v8ASTParserState) isInsideNamedCaptureGroup(name []uint16) bool {
	if name == nil {
		panic("V8 named-capture lookup requires a name")
	}
	for current := state; current != nil; current = current.previous {
		if current.captureName != nil && v8ASTCaptureNamesEqual(current.captureName, name) {
			return true
		}
	}
	return false
}

func (state *v8ASTParserState) newAlternative(capturesStarted int) {
	if state == nil {
		panic("nil V8 parser state")
	}
	if state.captureIndex == capturesStarted {
		return
	}
	from, to := state.captureIndex+1, capturesStarted
	if from > to {
		panic("V8 parser capture interval is reversed")
	}
	if count := len(state.nonParticipating); count > 0 && state.nonParticipating[count-1].to+1 == to {
		last := &state.nonParticipating[count-1]
		if from < last.from {
			panic("V8 parser capture interval overlaps out of order")
		}
		*last = v8ASTUnionInterval(*last, v8ASTNewInterval(from, to))
		return
	}
	state.nonParticipating = append(state.nonParticipating, v8ASTNewInterval(from, to))
}

// Ported V8 source declarations from regexp-parser-class.
// Copyright 2016 the V8 project authors. All rights reserved.
// V8 RegExp character-class parsing adapted from regexp-parser.cc.
// See LICENSE.v8 for attribution and license terms.
type v8ASTParserClassOperand struct {
	tree              any
	ranges            []v8ASTCharacterRange
	strings           v8ASTCharacterClassStrings
	character         rune
	kind              v8ASTClassSetOperandType
	isCharacter       bool
	mayContainStrings bool
}

type v8ASTParsedClassSet struct {
	tree              *v8ASTClassSetExpression
	mayContainStrings bool
}

var v8ASTUnicodePropertyRangeCache sync.Map

func v8ASTFromClassRanges(ranges []classRange) []v8ASTCharacterRange {
	converted := make([]v8ASTCharacterRange, len(ranges))
	for index, current := range ranges {
		converted[index] = v8ASTCharacterRange{from: current.lo, to: current.hi}
	}
	return converted
}

func v8ASTToClassRanges(ranges []v8ASTCharacterRange) []classRange {
	converted := make([]classRange, len(ranges))
	for index, current := range ranges {
		converted[index] = classRange{lo: current.from, hi: current.to}
	}
	return converted
}

func v8ASTUnicodePropertyRanges(property string) ([]v8ASTCharacterRange, bool) {
	if cached, ok := v8ASTUnicodePropertyRangeCache.Load(property); ok {
		return append([]v8ASTCharacterRange(nil), cached.([]v8ASTCharacterRange)...), true
	}
	var ranges []classRange
	switch {
	case strings.HasPrefix(property, "gc:"):
		ranges = regexpUnicodeGeneralCategoryRanges[strings.TrimPrefix(property, "gc:")]
	case strings.HasPrefix(property, "script:"):
		ranges = regexpUnicodeScriptRanges[strings.TrimPrefix(property, "script:")]
	case strings.HasPrefix(property, "binary:"):
		name := strings.TrimPrefix(property, "binary:")
		switch name {
		case "Any":
			ranges = []classRange{{0, 0x10ffff}}
		case "ASCII":
			ranges = []classRange{{0, 0x7f}}
		case "Assigned":
			ranges = v8ASTToClassRanges(v8ASTNegateRanges(v8ASTFromClassRanges(regexpUnicodeUnassignedRanges)))
		default:
			ranges = unicodeBinaryPropertyRanges(name)
		}
	case strings.HasPrefix(property, "script-extension:"):
		values := strings.SplitN(strings.TrimPrefix(property, "script-extension:"), "|", 2)
		if len(values) != 2 {
			return nil, false
		}
		// Script_Extensions equals Script except on the sparse UCD override
		// intervals. Build the cached set by applying those interval additions
		// and removals instead of scanning all 1.1 million code points.
		base, ok := regexpUnicodeScriptRanges[values[1]]
		if !ok && !regexpUnicodeScriptExtensionNonEmpty[values[0]] {
			return nil, false
		}
		included := v8ASTFromClassRanges(base)
		excluded := make([]v8ASTCharacterRange, 0, 16)
		needle := "|" + values[0] + "|"
		for _, override := range regexpScriptExtensionOverrides {
			rangeValue := v8ASTCharacterRange{from: override.lo, to: override.hi}
			if strings.Contains(override.values, needle) {
				included = append(included, rangeValue)
			} else {
				excluded = append(excluded, rangeValue)
			}
		}
		ranges = v8ASTToClassRanges(v8ASTSubtractRanges(included, excluded))
	default:
		return nil, false
	}
	if ranges == nil {
		return nil, false
	}
	converted := v8ASTFromClassRanges(ranges)
	v8ASTUnicodePropertyRangeCache.Store(property, append([]v8ASTCharacterRange(nil), converted...))
	return converted, true
}

func v8ASTUnicodePropertyAvailable(property string) bool {
	switch {
	case strings.HasPrefix(property, "gc:"):
		_, ok := regexpUnicodeGeneralCategoryRanges[strings.TrimPrefix(property, "gc:")]
		return ok
	case strings.HasPrefix(property, "script:"):
		ranges, ok := regexpUnicodeScriptRanges[strings.TrimPrefix(property, "script:")]
		return ok && len(ranges) != 0
	case strings.HasPrefix(property, "binary:"):
		name := strings.TrimPrefix(property, "binary:")
		if name == "Any" || name == "ASCII" || name == "Assigned" {
			return true
		}
		return unicodeBinaryPropertyRanges(name) != nil
	case strings.HasPrefix(property, "script-extension:"):
		values := strings.SplitN(strings.TrimPrefix(property, "script-extension:"), "|", 2)
		if len(values) != 2 {
			return false
		}
		code, ok := regexpUnicodeScriptCodeByName[values[1]]
		return ok && code == values[0] && regexpUnicodeScriptExtensionNonEmpty[code]
	default:
		return false
	}
}

func v8ASTMergeCharacterRanges(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	return v8ASTCanonicalizeRanges(append([]v8ASTCharacterRange(nil), ranges...))
}

func (parser *v8ASTParser) addMaybeSimpleCaseFoldedRange(ranges []v8ASTCharacterRange, value v8ASTCharacterRange) []v8ASTCharacterRange {
	if parser.mode != v8ASTBuildAST {
		return nil
	}
	if parser.flags.IgnoreCase {
		ranges = append(ranges, value)
		return v8TonodeAddUnicodeCaseEquivalents(ranges)
	}
	return v8ASTMergeCharacterRanges(append(ranges, value))
}

func (parser *v8ASTParser) closeUnicodeSetPropertyRanges(ranges []v8ASTCharacterRange) []v8ASTCharacterRange {
	if parser.mode != v8ASTBuildAST {
		return nil
	}
	if parser.flags.UnicodeSets && parser.flags.IgnoreCase {
		return v8TonodeAddUnicodeCaseEquivalents(ranges)
	}
	return v8ASTMergeCharacterRanges(ranges)
}

// parseUnicodePropertyName ports ParsePropertyClassName's ASCII grammar and
// cursor recovery. Property aliases are exact: only ASCII letters, digits,
// underscores, and the single name/value separator are accepted. Keep the
// cursor on the offending code unit when reporting an invalid character so
// Compile and VerifySyntax retain V8's error position.
func (parser *v8ASTParser) parseUnicodePropertyName(errorCode regexpError) (string, error) {
	cursor := parser.cursor
	if cursor.current != '{' {
		return "", parser.parseError(errorString(errorCode))
	}
	cursor.advance()
	var name strings.Builder
	for cursor.current != '}' && cursor.current != '=' {
		if !v8ASTIsUnicodePropertyValueCharacter(cursor.current) || !cursor.hasNextInput() {
			cursor.reportErrorAt(errorString(errorCode), cursor.position())
			return "", parser.parseError(errorString(errorCode))
		}
		name.WriteByte(byte(cursor.current))
		cursor.advance()
	}
	if cursor.current == '=' {
		name.WriteRune('=')
		cursor.advance()
		for cursor.current != '}' {
			if !v8ASTIsUnicodePropertyValueCharacter(cursor.current) || !cursor.hasNextInput() {
				cursor.reportErrorAt(errorString(errorCode), cursor.position())
				return "", parser.parseError(errorString(errorCode))
			}
			name.WriteByte(byte(cursor.current))
			cursor.advance()
		}
	}
	if cursor.current != '}' {
		cursor.reportErrorAt(errorString(errorCode), cursor.position())
		return "", parser.parseError(errorString(errorCode))
	}
	cursor.advance()
	return name.String(), nil
}

func v8ASTIsUnicodePropertyValueCharacter(character rune) bool {
	return character >= 'a' && character <= 'z' ||
		character >= 'A' && character <= 'Z' ||
		character >= '0' && character <= '9' || character == '_'
}

func (parser *v8ASTParser) parseCharacterClass(state *v8ASTParserState) (bool, v8ASTLookaroundType, error) {
	if parser.cursor.current != '[' {
		panic("V8 character-class parser must start at '['")
	}
	parser.cursor.advance()
	if parser.flags.UnicodeSets {
		expression, err := parser.parseUnicodeSetClass()
		if err != nil {
			return false, 0, err
		}
		parser.addAtom(state, expression.tree, nil)
		return false, 0, nil
	}
	negated := false
	if parser.cursor.current == '^' {
		negated = true
		parser.cursor.advance()
	}
	var ranges []v8ASTCharacterRange
	if parser.mode == v8ASTBuildAST {
		ranges = make([]v8ASTCharacterRange, 0, 8)
	}
	if parser.cursor.current == ']' {
		parser.cursor.advance()
		if parser.mode == v8ASTBuildAST {
			flags := v8ASTClassRangesFlags(0)
			if negated {
				flags |= v8ASTClassNegated
			}
			parser.addClassRanges(state, (v8ASTFactory{}).newClassRanges(nil, flags))
		} else {
			parser.addClassRanges(state, nil)
		}
		return false, 0, nil
	}
	for parser.cursor.current != v8ASTParserEndMarker && parser.cursor.current != ']' {
		first, err := parser.parseClassOperand(false)
		if err != nil {
			return false, 0, err
		}
		if parser.cursor.current != '-' {
			ranges = append(ranges, first.ranges...)
			continue
		}
		parser.cursor.advance()
		if parser.cursor.current == ']' {
			ranges = append(ranges, first.ranges...)
			ranges = append(ranges, v8ASTCharacterRange{from: '-', to: '-'})
			break
		}
		second, err := parser.parseClassOperand(false)
		if err != nil {
			return false, 0, err
		}
		if !first.isCharacter || !second.isCharacter {
			if parser.cursor.isUnicodeMode() {
				return false, 0, parser.parseError(errorString(errorInvalidCharacterClass))
			}
			ranges = append(ranges, first.ranges...)
			ranges = append(ranges, v8ASTCharacterRange{from: '-', to: '-'})
			ranges = append(ranges, second.ranges...)
			continue
		}
		if second.character < first.character {
			return false, 0, parser.parseError(errorString(errorOutOfOrderCharacterClass))
		}
		ranges = append(ranges, v8ASTCharacterRange{from: first.character, to: second.character})
	}
	if parser.cursor.current != ']' {
		return false, 0, parser.parseError(errorString(errorUnterminatedCharacterClass))
	}
	parser.cursor.advance()
	if parser.mode == v8ASTBuildAST {
		flags := v8ASTClassRangesFlags(0)
		if negated {
			flags |= v8ASTClassNegated
		}
		if !parser.flags.IgnoreCase {
			flags |= v8ASTClassNoCaseFoldingNeeded
		}
		if parser.cursor.oneByte && !negated {
			flags |= v8ASTClassCertainlyOneCodePoint
		}
		ranges = v8ASTMergeCharacterRanges(ranges)
		parser.addClassRanges(state, (v8ASTFactory{}).newClassRanges(ranges, flags))
	} else {
		parser.addClassRanges(state, nil)
	}
	return false, 0, nil
}

func (parser *v8ASTParser) parseClassOperand(unicodeSets bool) (v8ASTParserClassOperand, error) {
	cursor := parser.cursor
	if cursor.current == '[' && unicodeSets {
		cursor.advance()
		expression, err := parser.parseUnicodeSetClass()
		if err != nil {
			return v8ASTParserClassOperand{}, err
		}
		return v8ASTParserClassOperand{
			tree:              expression.tree,
			kind:              v8ASTNestedClass,
			mayContainStrings: expression.mayContainStrings,
		}, nil
	}
	if cursor.current == '\\' {
		return parser.parseClassEscape(unicodeSets)
	}
	character := cursor.current
	if character == v8ASTParserEndMarker || character == ']' {
		return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidCharacterInClass))
	}
	if unicodeSets {
		if cursor.isClassSetSyntaxCharacter(character) {
			return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidCharacterInClass))
		}
		if cursor.isClassSetReservedDoublePunctuator(character) && cursor.next() == character {
			return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidClassSetOperation))
		}
	}
	cursor.advance()
	var ranges []v8ASTCharacterRange
	if parser.mode == v8ASTBuildAST {
		ranges = []v8ASTCharacterRange{v8ASTRangeSingleton(character)}
	}
	return v8ASTParserClassOperand{
		ranges:      ranges,
		character:   character,
		kind:        v8ASTClassSetCharacter,
		isCharacter: true,
	}, nil
}

func (parser *v8ASTParser) parseClassSetItem() (v8ASTParserClassOperand, error) {
	return parser.parseClassOperand(true)
}

func (parser *v8ASTParser) classSetCharacterRange(character rune) []v8ASTCharacterRange {
	if parser.mode != v8ASTBuildAST {
		return nil
	}
	return []v8ASTCharacterRange{v8ASTRangeSingleton(character)}
}

func (parser *v8ASTParser) parseClassEscape(unicodeSets bool) (v8ASTParserClassOperand, error) {
	cursor := parser.cursor
	if cursor.current != '\\' {
		panic("V8 class-escape parser must start at backslash")
	}
	start := cursor.position()
	cursor.advance()
	character := cursor.current
	if character == v8ASTParserEndMarker {
		return v8ASTParserClassOperand{}, parser.parseError(errorString(errorEscapeAtEndOfPattern))
	}
	if character == 'd' || character == 'D' || character == 's' || character == 'S' || character == 'w' || character == 'W' {
		cursor.advance()
		var ranges []v8ASTCharacterRange
		if parser.mode == v8ASTBuildAST {
			kind := v8ASTStandardCharacterSet(character)
			addUnicodeCaseEquivalents := cursor.isUnicodeMode() && parser.flags.IgnoreCase
			ranges = v8TonodeAddClassEscape(kind, nil, addUnicodeCaseEquivalents)
		}
		return v8ASTParserClassOperand{
			ranges: ranges,
			kind:   v8ASTCharacterClassEscape,
		}, nil
	}
	if character == 'q' && unicodeSets {
		cursor.advance()
		return parser.parseClassStringDisjunction()
	}
	if character == 'p' || character == 'P' {
		if !cursor.isUnicodeMode() {
			cursor.advance()
			return v8ASTParserClassOperand{ranges: parser.classSetCharacterRange(character), character: character, kind: v8ASTClassSetCharacter, isCharacter: true}, nil
		}
		cursor.advance()
		propertyName, err := parser.parseUnicodePropertyName(errorInvalidClassPropertyName)
		if err != nil {
			return v8ASTParserClassOperand{}, err
		}
		property, ok := canonicalUnicodeProperty(propertyName)
		if !ok {
			return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidClassPropertyName))
		}
		if strings.HasPrefix(property, "string-property:") {
			if !unicodeSets || character == 'P' {
				return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidClassPropertyName))
			}
			data, ok := regexpUnicodeStringProperties[strings.TrimPrefix(property, "string-property:")]
			if !ok {
				return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidClassPropertyName))
			}
			return parser.parserClassOperandFromProperty(data), nil
		}
		var ranges []v8ASTCharacterRange
		var rangesOK bool
		if parser.mode == v8ASTBuildAST {
			ranges, rangesOK = v8ASTUnicodePropertyRanges(property)
		} else {
			rangesOK = v8ASTUnicodePropertyAvailable(property)
		}
		if !rangesOK {
			return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidClassPropertyName))
		}
		if parser.mode == v8ASTBuildAST {
			ranges = parser.closeUnicodeSetPropertyRanges(ranges)
			if character == 'P' {
				ranges = v8ASTNegateRanges(ranges)
			}
		}
		return v8ASTParserClassOperand{ranges: ranges, kind: v8ASTCharacterClassEscape}, nil
	}
	if character == 'b' {
		cursor.advance()
		return v8ASTParserClassOperand{ranges: parser.classSetCharacterRange('\b'), character: '\b', kind: v8ASTClassSetCharacter, isCharacter: true}, nil
	}
	if character == 'B' && cursor.isUnicodeMode() {
		return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidEscape))
	}
	if character == 'k' && !cursor.isUnicodeMode() {
		cursor.advance()
		if cursor.hasNamedCaptures(parser.captureScan, v8ASTInClass) {
			return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidEscape))
		}
		return v8ASTParserClassOperand{ranges: parser.classSetCharacterRange('k'), character: 'k', kind: v8ASTClassSetCharacter, isCharacter: true}, nil
	}
	if character == 'k' {
		return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidEscape))
	}
	cursor.reset(start)
	value := cursor.parseCharacterEscape(v8ASTInClass, func(v8ASTInClassEscapeState) bool {
		return parser.captureScan.hasNamed
	})
	if cursor.failedParse() {
		return v8ASTParserClassOperand{}, parser.parseError(cursor.errorCode)
	}
	return v8ASTParserClassOperand{ranges: parser.classSetCharacterRange(value), character: value, kind: v8ASTClassSetCharacter, isCharacter: true}, nil
}

func (parser *v8ASTParser) parserClassOperandFromProperty(data unicodeStringPropertyData) v8ASTParserClassOperand {
	if parser.mode != v8ASTBuildAST {
		return v8ASTParserClassOperand{kind: v8ASTCharacterClassEscape, mayContainStrings: len(data.strings) != 0}
	}
	ranges := v8ASTFromClassRanges(data.ranges)
	ranges = parser.closeUnicodeSetPropertyRanges(ranges)
	stringsIn := make(v8ASTCharacterClassStrings, 0, len(data.strings))
	for _, codePoints := range data.strings {
		// Unicode string-property tables are immutable generated data. Keep
		// their code-point slices by reference in the common case; copying every
		// sequence here doubled the table's transient memory before lowering.
		foldedCodePoints := codePoints
		if parser.flags.IgnoreCase {
			foldedCodePoints = append([]rune(nil), codePoints...)
			for index, codePoint := range foldedCodePoints {
				foldedCodePoints[index] = regexpSimpleCaseFold(codePoint)
			}
		}
		var tree any
		if parser.mode == v8ASTBuildAST {
			tree = (v8ASTFactory{}).newEmpty()
			if len(codePoints) != 0 {
				tree = (v8ASTFactory{}).newAtom(utf16.Encode(codePoints))
			}
		}
		stringsIn = append(stringsIn, v8ASTCharacterClassString{codePoints: foldedCodePoints, tree: tree})
	}
	return v8ASTParserClassOperand{ranges: ranges, strings: stringsIn, kind: v8ASTCharacterClassEscape, mayContainStrings: len(stringsIn) != 0}
}

func (parser *v8ASTParser) parseClassStringDisjunction() (v8ASTParserClassOperand, error) {
	if parser.cursor.current != '{' {
		return v8ASTParserClassOperand{}, parser.parseError(errorString(errorInvalidEscape))
	}
	parser.cursor.advance()
	var ranges []v8ASTCharacterRange
	var stringsIn v8ASTCharacterClassStrings
	var current []rune
	currentLength := 0
	containsStrings := false
	if parser.mode == v8ASTBuildAST {
		ranges = make([]v8ASTCharacterRange, 0, 4)
		stringsIn = make(v8ASTCharacterClassStrings, 0, 4)
		current = make([]rune, 0, 8)
	}
	addAlternative := func(value []rune, length int) {
		if parser.flags.IgnoreCase {
			for index, character := range value {
				value[index] = regexpSimpleCaseFold(character)
			}
		}
		if length == 1 {
			if parser.mode == v8ASTBuildAST {
				ranges = append(ranges, v8ASTRangeSingleton(value[0]))
			}
			return
		}
		containsStrings = true
		if parser.mode != v8ASTBuildAST {
			return
		}
		points := append([]rune(nil), value...)
		var tree any = (v8ASTFactory{}).newEmpty()
		if len(points) > 0 {
			tree = (v8ASTFactory{}).newAtom(utf16.Encode(points))
		}
		stringsIn = append(stringsIn, v8ASTCharacterClassString{codePoints: points, tree: tree})
	}
	for parser.cursor.current != v8ASTParserEndMarker {
		switch parser.cursor.current {
		case '|':
			addAlternative(current, currentLength)
			current = nil
			currentLength = 0
			parser.cursor.advance()
		case '}':
			addAlternative(current, currentLength)
			parser.cursor.advance()
			if parser.mode == v8ASTBuildAST {
				ranges = v8ASTMergeCharacterRanges(ranges)
			}
			return v8ASTParserClassOperand{ranges: ranges, strings: stringsIn, kind: v8ASTClassStringDisjunction, mayContainStrings: containsStrings}, nil
		default:
			character, err := parser.parseClassSetCharacter()
			if err != nil {
				return v8ASTParserClassOperand{}, err
			}
			if parser.mode == v8ASTBuildAST {
				current = append(current, character)
			}
			currentLength++
		}
	}
	return v8ASTParserClassOperand{}, parser.parseError(errorString(errorUnterminatedCharacterClass))
}

func (parser *v8ASTParser) parseUnicodeSetClass() (v8ASTParsedClassSet, error) {
	if parser.nesting >= maxParserNesting {
		return v8ASTParsedClassSet{}, parser.parseError(errorString(errorStackOverflow))
	}
	parser.nesting++
	defer func() { parser.nesting-- }()
	negated := false
	if parser.cursor.current == '^' {
		negated = true
		parser.cursor.advance()
	}
	if parser.cursor.current == ']' {
		parser.cursor.advance()
		return parser.newClassSetExpression(v8ASTClassSetUnion, negated, false, nil), nil
	}
	if parser.cursor.current == v8ASTParserEndMarker {
		return v8ASTParsedClassSet{}, parser.parseError(errorString(errorUnterminatedCharacterClass))
	}
	first, err := parser.parseClassOperand(true)
	if err != nil {
		return v8ASTParsedClassSet{}, err
	}
	switch {
	case parser.hasUnicodeSetOperator('-'):
		first = parser.classSetOperandWithTree(first)
		return parser.parseClassSubtraction(negated, first)
	case parser.hasUnicodeSetOperator('&'):
		first = parser.classSetOperandWithTree(first)
		return parser.parseClassIntersection(negated, first)
	default:
		return parser.parseClassUnion(negated, first)
	}
}

func (parser *v8ASTParser) parseClassSetCharacter() (rune, error) {
	cursor := parser.cursor
	character := cursor.current
	if character == '\\' {
		next := cursor.next()
		if next == 'b' {
			cursor.advanceBy(2)
			return '\b', nil
		}
		if next == v8ASTParserEndMarker {
			return 0, parser.parseError(errorString(errorEscapeAtEndOfPattern))
		}
		value := cursor.parseCharacterEscape(v8ASTInClass, func(v8ASTInClassEscapeState) bool {
			return parser.cursor.hasNamedCaptures(parser.captureScan, v8ASTInClass)
		})
		if cursor.failedParse() {
			return 0, parser.parseError(cursor.errorCode)
		}
		return value, nil
	}
	if cursor.isClassSetSyntaxCharacter(character) {
		return 0, parser.parseError(errorString(errorInvalidCharacterInClass))
	}
	if cursor.isClassSetReservedDoublePunctuator(character) {
		return 0, parser.parseError(errorString(errorInvalidClassSetOperation))
	}
	cursor.advance()
	return character, nil
}

func (parser *v8ASTParser) classSetOperandWithTree(operand v8ASTParserClassOperand) v8ASTParserClassOperand {
	if operand.tree != nil || parser.mode != v8ASTBuildAST {
		return operand
	}
	ranges := append([]v8ASTCharacterRange(nil), operand.ranges...)
	if operand.isCharacter {
		ranges = parser.addMaybeSimpleCaseFoldedRange(nil, v8ASTRangeSingleton(operand.character))
	}
	operand.tree = (v8ASTFactory{}).newClassSetOperand(ranges, operand.strings)
	return operand
}

func (parser *v8ASTParser) newClassSetExpression(operation v8ASTClassSetOperation, negated, mayContainStrings bool, operands []any) v8ASTParsedClassSet {
	if parser.mode == v8ASTBuildAST {
		return v8ASTParsedClassSet{tree: (v8ASTFactory{}).newClassSetExpression(operation, negated, mayContainStrings, operands), mayContainStrings: mayContainStrings}
	}
	return v8ASTParsedClassSet{mayContainStrings: mayContainStrings}
}

func (parser *v8ASTParser) parseClassUnion(negated bool, first v8ASTParserClassOperand) (v8ASTParsedClassSet, error) {
	var operands []any
	var ranges []v8ASTCharacterRange
	var stringsIn v8ASTCharacterClassStrings
	if parser.mode == v8ASTBuildAST {
		operands = make([]any, 0, 4)
		ranges = make([]v8ASTCharacterRange, 0, 4)
		stringsIn = make(v8ASTCharacterClassStrings, 0, 2)
	}
	mayContainStrings := first.mayContainStrings
	lastKind, character := first.kind, first.character
	if first.tree != nil {
		operands = append(operands, first.tree)
	} else if !first.isCharacter {
		ranges = append(ranges, first.ranges...)
		stringsIn = append(stringsIn, first.strings...)
	}
	for parser.cursor.current != ']' && parser.cursor.current != v8ASTParserEndMarker {
		if parser.cursor.current == '-' {
			if parser.cursor.next() == '-' {
				return v8ASTParsedClassSet{}, parser.parseError(errorString(errorInvalidClassSetOperation))
			}
			parser.cursor.advance()
			if parser.cursor.current == v8ASTParserEndMarker {
				break
			}
			if lastKind != v8ASTClassSetCharacter {
				return v8ASTParsedClassSet{}, parser.parseError(errorString(errorInvalidCharacterClass))
			}
			from := character
			right, err := parser.parseClassOperand(true)
			if err != nil {
				return v8ASTParsedClassSet{}, err
			}
			if !right.isCharacter {
				return v8ASTParsedClassSet{}, parser.parseError(errorString(errorInvalidCharacterClass))
			}
			if from > right.character {
				return v8ASTParsedClassSet{}, parser.parseError(errorString(errorOutOfOrderCharacterClass))
			}
			ranges = parser.addMaybeSimpleCaseFoldedRange(ranges, v8ASTCharacterRange{from: from, to: right.character})
			lastKind, character = v8ASTClassSetRange, right.character
			continue
		}
		if lastKind == v8ASTClassSetCharacter {
			ranges = parser.addMaybeSimpleCaseFoldedRange(ranges, v8ASTRangeSingleton(character))
		}
		right, err := parser.parseClassOperand(true)
		if err != nil {
			return v8ASTParsedClassSet{}, err
		}
		lastKind, character = right.kind, right.character
		mayContainStrings = mayContainStrings || right.mayContainStrings
		if right.tree != nil {
			if parser.mode == v8ASTBuildAST {
				if len(ranges) != 0 || len(stringsIn) != 0 {
					operands = append(operands, (v8ASTFactory{}).newClassSetOperand(ranges, stringsIn))
					ranges = make([]v8ASTCharacterRange, 0, 4)
					stringsIn = make(v8ASTCharacterClassStrings, 0, 2)
				}
				operands = append(operands, right.tree)
			}
		} else if !right.isCharacter {
			ranges = append(ranges, right.ranges...)
			stringsIn = append(stringsIn, right.strings...)
		}
	}
	if parser.cursor.current == v8ASTParserEndMarker {
		return v8ASTParsedClassSet{}, parser.parseError(errorString(errorUnterminatedCharacterClass))
	}
	if lastKind == v8ASTClassSetCharacter {
		ranges = parser.addMaybeSimpleCaseFoldedRange(ranges, v8ASTRangeSingleton(character))
	}
	if parser.mode == v8ASTBuildAST && (len(ranges) != 0 || len(stringsIn) != 0) {
		operands = append(operands, (v8ASTFactory{}).newClassSetOperand(ranges, stringsIn))
	}
	parser.cursor.advance()
	if negated && mayContainStrings {
		return v8ASTParsedClassSet{}, parser.parseError(errorString(errorNegatedCharacterClassWithStrings))
	}
	return parser.newClassSetExpression(v8ASTClassSetUnion, negated, mayContainStrings, operands), nil
}

func (parser *v8ASTParser) parseClassIntersection(negated bool, first v8ASTParserClassOperand) (v8ASTParsedClassSet, error) {
	first = parser.classSetOperandWithTree(first)
	operands := []any(nil)
	if parser.mode == v8ASTBuildAST {
		operands = append(operands, first.tree)
	}
	mayContainStrings := first.mayContainStrings
	for parser.cursor.current != ']' && parser.cursor.current != v8ASTParserEndMarker {
		if !parser.hasUnicodeSetOperator('&') {
			return v8ASTParsedClassSet{}, parser.parseError(errorString(errorInvalidClassSetOperation))
		}
		parser.cursor.advanceBy(2)
		if parser.cursor.current == '&' {
			return v8ASTParsedClassSet{}, parser.parseError(errorString(errorInvalidCharacterInClass))
		}
		if parser.cursor.current == v8ASTParserEndMarker {
			return v8ASTParsedClassSet{}, parser.parseError(errorString(errorUnterminatedCharacterClass))
		}
		right, err := parser.parseClassOperand(true)
		if err != nil {
			return v8ASTParsedClassSet{}, err
		}
		right = parser.classSetOperandWithTree(right)
		mayContainStrings = mayContainStrings && right.mayContainStrings
		if parser.mode == v8ASTBuildAST {
			operands = append(operands, right.tree)
		}
	}
	if parser.cursor.current == v8ASTParserEndMarker {
		return v8ASTParsedClassSet{}, parser.parseError(errorString(errorUnterminatedCharacterClass))
	}
	if negated && mayContainStrings {
		return v8ASTParsedClassSet{}, parser.parseError(errorString(errorNegatedCharacterClassWithStrings))
	}
	parser.cursor.advance()
	return parser.newClassSetExpression(v8ASTClassSetIntersection, negated, mayContainStrings, operands), nil
}

func (parser *v8ASTParser) parseClassSubtraction(negated bool, first v8ASTParserClassOperand) (v8ASTParsedClassSet, error) {
	first = parser.classSetOperandWithTree(first)
	mayContainStrings := first.mayContainStrings
	if negated && mayContainStrings {
		return v8ASTParsedClassSet{}, parser.parseError(errorString(errorNegatedCharacterClassWithStrings))
	}
	operands := []any(nil)
	if parser.mode == v8ASTBuildAST {
		operands = append(operands, first.tree)
	}
	for parser.cursor.current != ']' && parser.cursor.current != v8ASTParserEndMarker {
		if !parser.hasUnicodeSetOperator('-') {
			return v8ASTParsedClassSet{}, parser.parseError(errorString(errorInvalidClassSetOperation))
		}
		parser.cursor.advanceBy(2)
		if parser.cursor.current == v8ASTParserEndMarker {
			return v8ASTParsedClassSet{}, parser.parseError(errorString(errorUnterminatedCharacterClass))
		}
		right, err := parser.parseClassOperand(true)
		if err != nil {
			return v8ASTParsedClassSet{}, err
		}
		right = parser.classSetOperandWithTree(right)
		if parser.mode == v8ASTBuildAST {
			operands = append(operands, right.tree)
		}
	}
	if parser.cursor.current == v8ASTParserEndMarker {
		return v8ASTParsedClassSet{}, parser.parseError(errorString(errorUnterminatedCharacterClass))
	}
	parser.cursor.advance()
	return parser.newClassSetExpression(v8ASTClassSetSubtraction, negated, mayContainStrings, operands), nil
}

func v8ASTClassSetMayContainStrings(operand any) bool {
	switch value := operand.(type) {
	case *v8ASTClassSetOperand:
		return value != nil && value.hasStrings()
	case *v8ASTClassSetExpression:
		return value != nil && value.mayContainStrings
	default:
		return false
	}
}

func (parser *v8ASTParser) hasUnicodeSetOperator(operator rune) bool {
	return parser.cursor.current == operator && parser.cursor.next() == operator
}
