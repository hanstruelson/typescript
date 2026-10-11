// Package regexpengine is the standalone semantic core for ECMAScript RegExp.
// It uses UTF-16 positions like JavaScript and is independent of Goja.
package regexpengine

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
)

type Flag uint16

const (
	FlagGlobal Flag = 1 << iota
	FlagIgnoreCase
	FlagMultiline
	FlagSticky
	FlagUnicode
	FlagDotAll
	FlagLinear
	FlagHasIndices
	FlagUnicodeSets
)

type Flags struct {
	Global, HasIndices, IgnoreCase, Linear, Multiline, DotAll bool
	Unicode, UnicodeSets, Sticky                              bool
}

func (f Flags) Has(flag Flag) bool {
	switch flag {
	case FlagGlobal:
		return f.Global
	case FlagHasIndices:
		return f.HasIndices
	case FlagIgnoreCase:
		return f.IgnoreCase
	case FlagLinear:
		return f.Linear
	case FlagMultiline:
		return f.Multiline
	case FlagDotAll:
		return f.DotAll
	case FlagUnicode:
		return IsUnicode(f)
	case FlagUnicodeSets:
		return f.UnicodeSets
	case FlagSticky:
		return f.Sticky
	default:
		return false
	}
}

func (f *Flags) Set(flag Flag, value bool) {
	switch flag {
	case FlagGlobal:
		f.Global = value
	case FlagHasIndices:
		f.HasIndices = value
	case FlagIgnoreCase:
		f.IgnoreCase = value
	case FlagLinear:
		f.Linear = value
	case FlagMultiline:
		f.Multiline = value
	case FlagDotAll:
		f.DotAll = value
	case FlagUnicode:
		f.Unicode = value || f.UnicodeSets
	case FlagUnicodeSets:
		f.UnicodeSets = value
		if value {
			f.Unicode = true
		}
	case FlagSticky:
		f.Sticky = value
	}
}

func IsHasIndices(f Flags) bool    { return f.HasIndices }
func IsGlobal(f Flags) bool        { return f.Global }
func IsIgnoreCase(f Flags) bool    { return f.IgnoreCase }
func IsLinear(f Flags) bool        { return f.Linear }
func IsMultiline(f Flags) bool     { return f.Multiline }
func IsDotAll(f Flags) bool        { return f.DotAll }
func IsUnicode(f Flags) bool       { return f.Unicode && !f.UnicodeSets }
func IsUnicodeSets(f Flags) bool   { return f.UnicodeSets }
func IsSticky(f Flags) bool        { return f.Sticky }
func IsEitherUnicode(f Flags) bool { return f.Unicode || f.UnicodeSets }
func ShouldOptionallyStepBackToLeadSurrogate(f Flags) bool {
	return IsEitherUnicode(f) && (f.Global || f.Sticky)
}

func TryFlagFromChar(c rune) (Flag, bool) {
	switch c {
	case 'd':
		return FlagHasIndices, true
	case 'g':
		return FlagGlobal, true
	case 'i':
		return FlagIgnoreCase, true
	case 'l':
		return FlagLinear, true
	case 'm':
		return FlagMultiline, true
	case 's':
		return FlagDotAll, true
	case 'u':
		return FlagUnicode, true
	case 'v':
		return FlagUnicodeSets, true
	case 'y':
		return FlagSticky, true
	default:
		return 0, false
	}
}

func (f Flags) String() string {
	var b strings.Builder
	for _, c := range "dgilmsuvy" {
		flag, _ := TryFlagFromChar(c)
		if f.Has(flag) {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func ParseFlags(text string) (Flags, error) {
	var f, seen Flags
	for _, c := range text {
		flag, ok := TryFlagFromChar(c)
		if !ok {
			return Flags{}, fmt.Errorf("invalid regular expression flag %q", c)
		}
		if seen.Has(flag) {
			return Flags{}, fmt.Errorf("duplicate regular expression flag %q", c)
		}
		seen.Set(flag, true)
		f.Set(flag, true)
	}
	if strings.ContainsRune(text, 'u') && strings.ContainsRune(text, 'v') {
		return Flags{}, fmt.Errorf("regular expression flags u and v cannot be used together")
	}
	return f, nil
}

type nodeKind uint8

const (
	nodeEmpty nodeKind = iota
	nodeLiteral
	nodeAny
	nodeClass
	nodeConcat
	nodeAlternate
	nodeRepeat
	nodeCapture
	nodeBackref
	nodeLookaround
	nodeBegin
	nodeEnd
	nodeBoundary
)

type node struct {
	kind                          nodeKind
	rune                          rune
	class                         []classTerm
	set                           *classExpr
	negated                       bool
	children                      []*node
	min, max                      int
	greedy                        bool
	capture                       int
	captureIndices                []int
	refName, refDigits            string
	negative, behind              bool
	ignoreCase, dotAll, multiline bool
}

type choiceBranch struct{ id, branch int }
type namedCaptureOccurrence struct {
	index int
	path  []choiceBranch
}

type classRange struct{ lo, hi rune }
type classTerm struct {
	ranges   []classRange
	property string
	negated  bool
}
type classExprKind uint8

const (
	classExprUnion classExprKind = iota
	classExprIntersection
	classExprSubtraction
	classExprComplement
	classExprStringSet
	classExprTerm
)

type classExpr struct {
	kind     classExprKind
	children []*classExpr
	strings  [][]rune
	term     classTerm
}

type parser struct {
	text                                   []uint16
	index                                  int
	flags                                  Flags
	group                                  int
	nesting                                int
	groupNames                             map[string][]int
	namedCapturePaths                      map[string][]namedCaptureOccurrence
	choicePath                             []choiceBranch
	nextChoiceID                           int
	namedCaptureScanDone, hasNamedCaptures bool
}

const maxParserNesting = 2048
const maxCaptureGroups = 32767

func (p *parser) peek() rune {
	if p.index >= len(p.text) {
		return -1
	}
	return rune(p.text[p.index])
}
func (p *parser) take() rune {
	c := p.peek()
	if c >= 0 {
		p.index++
	}
	return c
}
func (p *parser) takePatternCharacter() rune {
	c := p.take()
	if p.flags.Unicode && c >= 0xd800 && c <= 0xdbff && p.index < len(p.text) {
		n := rune(p.text[p.index])
		if n >= 0xdc00 && n <= 0xdfff {
			p.index++
			return 0x10000 + ((c - 0xd800) << 10) + n - 0xdc00
		}
	}
	return c
}
func utf16String(units []uint16) string { return string(utf16.Decode(units)) }
func asciiFromUnits(units []uint16) string {
	var b strings.Builder
	for _, c := range units {
		if c > 0x7f {
			return ""
		}
		b.WriteByte(byte(c))
	}
	return b.String()
}
func (p *parser) errorf(format string, args ...any) error {
	return fmt.Errorf("regular expression at %d: %s", p.index, fmt.Sprintf(format, args...))
}

func (p *parser) nextCapture() (int, error) {
	if p.group >= maxCaptureGroups {
		return 0, p.errorf("%s", errorString(errorTooManyCaptures))
	}
	p.group++
	return p.group, nil
}

func (p *parser) scanNamedCaptures() {
	p.namedCaptureScanDone = true
	depth := 0
	for i := 0; i < len(p.text); i++ {
		c := rune(p.text[i])
		if c == '\\' {
			i++
			continue
		}
		if depth > 0 {
			if p.flags.UnicodeSets && c == '[' {
				depth++
			} else if c == ']' {
				depth--
			}
			continue
		}
		if c == '[' {
			depth = 1
			continue
		}
		if c == '(' && i+3 < len(p.text) && p.text[i+1] == '?' && p.text[i+2] == '<' && p.text[i+3] != '=' && p.text[i+3] != '!' {
			p.hasNamedCaptures = true
			return
		}
	}
}

func (p *parser) parse() (*node, error) {
	if !p.namedCaptureScanDone {
		p.scanNamedCaptures()
	}
	root, err := p.alternation()
	if err != nil {
		return nil, err
	}
	if p.index != len(p.text) {
		return nil, p.errorf("unexpected %q", p.peek())
	}
	normalizeClassRanges(root)
	return root, nil
}
func (p *parser) alternation() (*node, error) {
	choiceID := p.nextChoiceID
	p.nextChoiceID++
	p.choicePath = append(p.choicePath, choiceBranch{id: choiceID})
	defer func() { p.choicePath = p.choicePath[:len(p.choicePath)-1] }()
	parts := []*node{}
	for {
		n, e := p.concatenation()
		if e != nil {
			return nil, e
		}
		parts = append(parts, n)
		if p.peek() != '|' {
			break
		}
		p.take()
		p.choicePath[len(p.choicePath)-1].branch++
	}
	if len(parts) == 1 {
		return parts[0], nil
	}
	return &node{kind: nodeAlternate, children: parts}, nil
}
func (p *parser) concatenation() (*node, error) {
	parts := []*node{}
	for {
		c := p.peek()
		if c < 0 || c == ')' || c == '|' {
			break
		}
		n, e := p.quantified()
		if e != nil {
			return nil, e
		}
		if p.flags.Unicode && len(parts) > 0 {
			a := parts[len(parts)-1]
			if a.kind == nodeLiteral && n.kind == nodeLiteral && a.rune >= 0xd800 && a.rune <= 0xdbff && n.rune >= 0xdc00 && n.rune <= 0xdfff {
				a.rune = 0x10000 + ((a.rune - 0xd800) << 10) + n.rune - 0xdc00
				continue
			}
		}
		parts = append(parts, n)
	}
	if len(parts) == 0 {
		return &node{kind: nodeEmpty}, nil
	}
	if len(parts) == 1 {
		return parts[0], nil
	}
	return &node{kind: nodeConcat, children: parts}, nil
}
func (p *parser) quantified() (*node, error) {
	base, e := p.atom()
	if e != nil {
		return nil, e
	}
	min, max := 1, 1
	have := true
	switch p.peek() {
	case '*':
		p.take()
		min, max = 0, -1
	case '+':
		p.take()
		min, max = 1, -1
	case '?':
		p.take()
		min, max = 0, 1
	case '{':
		start := p.index
		p.take()
		if p.peek() < '0' || p.peek() > '9' {
			have = false
		} else {
			min = 0
			for p.peek() >= '0' && p.peek() <= '9' {
				min = min*10 + int(p.take()-'0')
				if min > 100000 {
					return nil, p.errorf("quantifier is too large")
				}
			}
			max = min
			if p.peek() == ',' {
				p.take()
				max = -1
				if p.peek() >= '0' && p.peek() <= '9' {
					max = 0
					for p.peek() >= '0' && p.peek() <= '9' {
						max = max*10 + int(p.take()-'0')
						if max > 100000 {
							return nil, p.errorf("quantifier is too large")
						}
					}
				}
			}
			if p.peek() != '}' {
				have = false
			} else {
				p.take()
			}
			if have && max >= 0 && max < min {
				return nil, p.errorf("maximum quantifier is smaller than minimum")
			}
		}
		if !have {
			if p.flags.Unicode {
				return nil, p.errorf("invalid or unterminated quantifier")
			}
			p.index = start
		}
	default:
		have = false
	}
	if !have {
		return base, nil
	}
	greedy := true
	if p.peek() == '?' {
		p.take()
		greedy = false
	}
	if p.peek() == '*' || p.peek() == '+' || p.peek() == '?' {
		return nil, p.errorf("multiple quantifiers")
	}
	if base.kind == nodeLookaround && (base.behind || p.flags.Unicode) {
		return nil, p.errorf("lookaround cannot be quantified in this mode")
	}
	return &node{kind: nodeRepeat, children: []*node{base}, min: min, max: max, greedy: greedy}, nil
}

func (p *parser) scopedFlags() (Flags, error) {
	f := p.flags
	var seen Flags
	dash := false
	after := 0
	for {
		c := p.peek()
		if c == ':' {
			p.take()
			if (seen.IgnoreCase || seen.Multiline || seen.DotAll) == false || dash && after == 0 {
				return Flags{}, p.errorf("invalid scoped modifiers")
			}
			return f, nil
		}
		if c == '-' {
			if dash {
				return Flags{}, p.errorf("duplicate modifier separator")
			}
			dash = true
			p.take()
			continue
		}
		if c != 'i' && c != 'm' && c != 's' {
			return Flags{}, p.errorf("invalid scoped modifier %q", c)
		}
		p.take()
		flag, _ := TryFlagFromChar(c)
		if seen.Has(flag) {
			return Flags{}, p.errorf("duplicate scoped modifier %q", c)
		}
		seen.Set(flag, true)
		if dash {
			after++
		}
		f.Set(flag, !dash)
	}
}
func (p *parser) atom() (*node, error) {
	switch c := p.take(); c {
	case -1:
		return nil, p.errorf("unexpected end of pattern")
	case '*', '+', '?':
		return nil, p.errorf("nothing to repeat")
	case '(':
		if p.nesting >= maxParserNesting {
			return nil, p.errorf("regular expression nesting limit exceeded")
		}
		p.nesting++
		defer func() { p.nesting-- }()
		capture := 0
		var e error
		scoped := false
		old := p.flags
		if p.peek() == '?' {
			p.take()
			q := p.peek()
			if q == 'i' || q == 'm' || q == 's' || q == '-' {
				f, e := p.scopedFlags()
				if e != nil {
					return nil, e
				}
				p.flags = f
				scoped = true
			} else {
				switch p.take() {
				case ':':
				case '=':
					return p.parseLookaround(false, false)
				case '!':
					return p.parseLookaround(true, false)
				case '<':
					if p.peek() == '=' || p.peek() == '!' {
						negative := p.take() == '!'
						return p.parseLookaround(negative, true)
					}
					capture, e = p.nextCapture()
					if e != nil {
						return nil, e
					}
					name, e := p.captureGroupName()
					if e != nil {
						return nil, e
					}
					if !validCaptureName(name) {
						return nil, p.errorf("invalid capture group name %q", name)
					}
					if p.groupNames == nil {
						p.groupNames = map[string][]int{}
					}
					if p.namedCapturePaths == nil {
						p.namedCapturePaths = map[string][]namedCaptureOccurrence{}
					}
					for _, existing := range p.namedCapturePaths[name] {
						if !capturePathsAreExclusive(existing.path, p.choicePath) {
							return nil, p.errorf("duplicate capture group name %q", name)
						}
					}
					p.groupNames[name] = append(p.groupNames[name], capture)
					path := append([]choiceBranch(nil), p.choicePath...)
					p.namedCapturePaths[name] = append(p.namedCapturePaths[name], namedCaptureOccurrence{index: capture, path: path})
				default:
					return nil, p.errorf("unsupported group prefix")
				}
			}
		} else {
			capture, e = p.nextCapture()
			if e != nil {
				return nil, e
			}
		}
		inner, e := p.alternation()
		if scoped {
			p.flags = old
		}
		if e != nil {
			return nil, e
		}
		if p.take() != ')' {
			return nil, p.errorf("unterminated group")
		}
		if capture == 0 {
			return inner, nil
		}
		return &node{kind: nodeCapture, children: []*node{inner}, capture: capture}, nil
	case '[':
		return p.characterClass()
	case '.':
		return &node{kind: nodeAny, dotAll: p.flags.DotAll}, nil
	case '^':
		return &node{kind: nodeBegin, multiline: p.flags.Multiline}, nil
	case '$':
		return &node{kind: nodeEnd, multiline: p.flags.Multiline}, nil
	case '\\':
		return p.escape(false)
	default:
		if p.flags.Unicode && (c == '{' || c == '}' || c == ']') {
			return nil, p.errorf("unescaped syntax character")
		}
		p.index--
		return &node{kind: nodeLiteral, rune: p.takePatternCharacter(), ignoreCase: p.flags.IgnoreCase}, nil
	}
}
func (p *parser) parseLookaround(negative, behind bool) (*node, error) {
	inner, e := p.alternation()
	if e != nil {
		return nil, e
	}
	if p.take() != ')' {
		return nil, p.errorf("unterminated lookaround")
	}
	return &node{kind: nodeLookaround, children: []*node{inner}, negative: negative, behind: behind}, nil
}

func (p *parser) captureGroupName() (string, error) {
	name := []rune{}
	for p.peek() != '>' && p.peek() >= 0 {
		c, e := p.captureNameCodePoint()
		if e != nil {
			return "", e
		}
		name = append(name, c)
	}
	if p.take() != '>' {
		return "", p.errorf("unterminated capture group name")
	}
	return string(name), nil
}
func (p *parser) captureNameCodePoint() (rune, error) {
	var c rune
	if p.peek() == '\\' {
		p.take()
		if p.take() != 'u' {
			return 0, p.errorf("invalid escape in capture name")
		}
		if p.peek() == '{' {
			p.take()
			start := p.index
			for isHexDigit(p.peek()) {
				p.take()
			}
			if start == p.index || p.take() != '}' {
				return 0, p.errorf("invalid Unicode escape in capture name")
			}
			v, e := strconv.ParseInt(asciiFromUnits(p.text[start:p.index-1]), 16, 32)
			if e != nil || v > unicode.MaxRune {
				return 0, p.errorf("invalid Unicode escape in capture name")
			}
			c = rune(v)
		} else {
			start := p.index
			if len(p.text)-start < 4 {
				return 0, p.errorf("invalid Unicode escape in capture name")
			}
			for i := 0; i < 4; i++ {
				if !isHexDigit(rune(p.text[start+i])) {
					return 0, p.errorf("invalid Unicode escape in capture name")
				}
			}
			p.index += 4
			v, e := strconv.ParseUint(asciiFromUnits(p.text[start:p.index]), 16, 16)
			if e != nil {
				return 0, p.errorf("invalid Unicode escape in capture name")
			}
			c = rune(v)
		}
	} else {
		first := p.take()
		c = first
		if first >= 0xd800 && first <= 0xdbff && p.peek() >= 0xdc00 && p.peek() <= 0xdfff {
			second := p.take()
			c = 0x10000 + ((first - 0xd800) << 10) + second - 0xdc00
		}
	}
	if c >= 0xd800 && c <= 0xdbff && (p.peek() == '\\' || p.peek() >= 0xdc00 && p.peek() <= 0xdfff) {
		mark := p.index
		next, e := p.captureNameCodePoint()
		if e == nil && next >= 0xdc00 && next <= 0xdfff {
			return 0x10000 + ((c - 0xd800) << 10) + next - 0xdc00, nil
		}
		p.index = mark
	}
	return c, nil
}
func isHexDigit(c rune) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

func (p *parser) escape(inClass bool) (*node, error) {
	c := p.take()
	if c < 0 {
		return nil, p.errorf("trailing escape")
	}
	lit := func(v rune) *node { return &node{kind: nodeLiteral, rune: v, ignoreCase: p.flags.IgnoreCase} }
	class := func(term classTerm) *node {
		return &node{kind: nodeClass, class: []classTerm{term}, ignoreCase: p.flags.IgnoreCase}
	}
	switch c {
	case 'd', 'D':
		return class(classTerm{ranges: []classRange{{'0', '9'}}, negated: c == 'D'}), nil
	case 'w', 'W':
		return class(classTerm{ranges: []classRange{{'0', '9'}, {'A', 'Z'}, {'a', 'z'}, {'_', '_'}}, negated: c == 'W'}), nil
	case 's', 'S':
		return class(classTerm{ranges: ecmaWhitespaceRanges(), negated: c == 'S'}), nil
	case 'n':
		return lit('\n'), nil
	case 'r':
		return lit('\r'), nil
	case 't':
		return lit('\t'), nil
	case 'f':
		return lit('\f'), nil
	case 'v':
		return lit('\v'), nil
	case 'b':
		if inClass {
			return lit('\b'), nil
		}
		return &node{kind: nodeBoundary, ignoreCase: p.flags.IgnoreCase}, nil
	case 'B':
		if inClass {
			if p.flags.Unicode {
				return nil, p.errorf("invalid identity escape")
			}
			return lit('B'), nil
		}
		return &node{kind: nodeBoundary, negated: true, ignoreCase: p.flags.IgnoreCase}, nil
	case 'x':
		return p.hexEscape(2)
	case 'u':
		if p.peek() == '{' {
			if !p.flags.Unicode {
				return lit('u'), nil
			}
			p.take()
			start := p.index
			for isHexDigit(p.peek()) {
				p.take()
			}
			if start == p.index || p.take() != '}' {
				return nil, p.errorf("invalid Unicode escape")
			}
			v, e := strconv.ParseInt(asciiFromUnits(p.text[start:p.index-1]), 16, 32)
			if e != nil || v > unicode.MaxRune {
				return nil, p.errorf("invalid Unicode escape")
			}
			return lit(rune(v)), nil
		}
		return p.hexEscape(4)
	case 'c':
		if p.peek() >= 'a' && p.peek() <= 'z' || p.peek() >= 'A' && p.peek() <= 'Z' {
			return lit(p.take() % 32), nil
		}
		if p.flags.Unicode {
			return nil, p.errorf("invalid control escape")
		}
		return lit('c'), nil
	case 'p', 'P':
		if !p.flags.Unicode {
			return lit(c), nil
		}
		if p.take() != '{' {
			return nil, p.errorf("invalid Unicode property escape")
		}
		start := p.index
		for p.peek() != '}' && p.peek() >= 0 {
			p.take()
		}
		if p.take() != '}' {
			return nil, p.errorf("unterminated Unicode property escape")
		}
		name := asciiFromUnits(p.text[start : p.index-1])
		prop, ok := canonicalUnicodeProperty(name)
		if !ok {
			return nil, p.errorf("unsupported Unicode property %q", name)
		}
		return class(classTerm{property: prop, negated: c == 'P'}), nil
	case 'k':
		if p.peek() == '<' && (p.flags.Unicode || p.hasNamedCaptures) {
			p.take()
			name, e := p.captureGroupName()
			if e != nil {
				return nil, e
			}
			return &node{kind: nodeBackref, refName: name, ignoreCase: p.flags.IgnoreCase}, nil
		}
		if p.flags.Unicode {
			return nil, p.errorf("invalid named backreference")
		}
		return lit('k'), nil
	case 'q':
		if inClass && p.flags.UnicodeSets {
			return p.classStringDisjunction()
		}
		if p.flags.Unicode {
			return nil, p.errorf("invalid identity escape")
		}
		return lit('q'), nil
	default:
		if c >= '0' && c <= '9' {
			if inClass {
				if p.flags.Unicode {
					if c == '0' && !(p.peek() >= '0' && p.peek() <= '9') {
						return lit(0), nil
					}
					return nil, p.errorf("invalid decimal escape in class")
				}
				if c >= '8' {
					return lit(c), nil
				}
				value := int(c - '0')
				for j := 0; j < 2 && p.peek() >= '0' && p.peek() <= '7'; j++ {
					next := value*8 + int(p.peek()-'0')
					if next > 255 {
						break
					}
					value = next
					p.take()
				}
				return lit(rune(value)), nil
			}
			digits := []uint16{uint16(c)}
			value := int(c - '0')
			for p.peek() >= '0' && p.peek() <= '9' {
				d := p.take()
				digits = append(digits, uint16(d))
				value = value*10 + int(d-'0')
			}
			return &node{kind: nodeBackref, capture: value, refDigits: asciiFromUnits(digits), ignoreCase: p.flags.IgnoreCase}, nil
		}
		if p.flags.Unicode && !validIdentityEscape(c, inClass, p.flags.UnicodeSets) {
			return nil, p.errorf("invalid identity escape \\%c", c)
		}
		return lit(c), nil
	}
}

func validIdentityEscape(c rune, inClass, sets bool) bool {
	if strings.ContainsRune("^$\\.*+?()[]{}|/", c) {
		return true
	}
	if inClass && c == '-' {
		return true
	}
	return sets && inClass && strings.ContainsRune("!#$%&*+,:;<=>?@^`~", c)
}
func (p *parser) hexEscape(width int) (*node, error) {
	start := p.index
	if !p.flags.Unicode && (len(p.text)-start < width) {
		return &node{kind: nodeLiteral, rune: rune(p.text[start-1]), ignoreCase: p.flags.IgnoreCase}, nil
	}
	for i := 0; i < width; i++ {
		c := p.take()
		if !isHexDigit(c) {
			return nil, p.errorf("invalid hexadecimal escape")
		}
	}
	v, e := strconv.ParseInt(asciiFromUnits(p.text[start:p.index]), 16, 32)
	if e != nil {
		return nil, p.errorf("invalid hexadecimal escape")
	}
	return &node{kind: nodeLiteral, rune: rune(v), ignoreCase: p.flags.IgnoreCase}, nil
}

func (p *parser) classStringDisjunction() (*node, error) {
	if p.take() != '{' {
		return nil, p.errorf("class string escape requires braces")
	}
	stringsIn := [][]rune{}
	current := []rune{}
	for {
		if p.peek() < 0 {
			return nil, p.errorf("unterminated class string")
		}
		switch p.peek() {
		case '|':
			p.take()
			stringsIn = append(stringsIn, current)
			current = nil
		case '}':
			p.take()
			stringsIn = append(stringsIn, current)
			return &node{kind: nodeClass, set: &classExpr{kind: classExprStringSet, strings: stringsIn}, ignoreCase: p.flags.IgnoreCase}, nil
		case '\\':
			p.take()
			n, e := p.escape(true)
			if e != nil {
				return nil, e
			}
			if n.kind != nodeLiteral {
				return nil, p.errorf("invalid character in class string")
			}
			current = append(current, n.rune)
		default:
			current = append(current, p.takePatternCharacter())
		}
	}
}

func ecmaWhitespaceRanges() []classRange {
	return []classRange{{'\t', '\r'}, {' ', ' '}, {0xa0, 0xa0}, {0x1680, 0x1680}, {0x2000, 0x200a}, {0x2028, 0x2029}, {0x202f, 0x202f}, {0x205f, 0x205f}, {0x3000, 0x3000}, {0xfeff, 0xfeff}}
}
func normalizePropertyName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}
func canonicalUnicodeProperty(s string) (string, bool) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) == 1 {
		v := normalizePropertyName(parts[0])
		cats := map[string]string{"l": "L", "letter": "L", "lu": "Lu", "uppercaseletter": "Lu", "ll": "Ll", "lowercaseletter": "Ll", "lt": "Lt", "lm": "Lm", "lo": "Lo", "m": "M", "n": "N", "nd": "Nd", "decimalnumber": "Nd", "p": "P", "s": "S", "z": "Z", "c": "C", "white_space": "White_Space", "whitespace": "White_Space", "alphabetic": "Alphabetic", "ascii": "ASCII", "any": "Any", "assigned": "Assigned", "idstart": "ID_Start", "idcontinue": "ID_Continue"}
		if c, ok := cats[v]; ok {
			if _, exists := unicode.Categories[c]; exists {
				return "gc:" + c, true
			}
			return "binary:" + c, true
		}
		for name := range unicode.Properties {
			if normalizePropertyName(name) == v {
				return "binary:" + name, true
			}
		}
		return "", false
	}
	key, value := normalizePropertyName(parts[0]), normalizePropertyName(parts[1])
	if key == "gc" || key == "generalcategory" {
		return canonicalUnicodeProperty(value)
	}
	if key == "sc" || key == "script" || key == "scx" || key == "scriptextensions" {
		for name := range unicode.Scripts {
			if normalizePropertyName(name) == value {
				return "script:" + name, true
			}
		}
	}
	return "", false
}
func unicodePropertyContains(c rune, property string) bool {
	parts := strings.SplitN(property, ":", 2)
	if len(parts) != 2 {
		return false
	}
	switch parts[0] {
	case "gc":
		return unicode.Is(unicode.Categories[parts[1]], c)
	case "script":
		return unicode.Is(unicode.Scripts[parts[1]], c)
	case "binary":
		switch parts[1] {
		case "White_Space":
			for _, r := range ecmaWhitespaceRanges() {
				if c >= r.lo && c <= r.hi {
					return true
				}
			}
			return false
		case "ASCII":
			return c <= 0x7f
		case "Any":
			return c >= 0 && c <= unicode.MaxRune
		case "Assigned":
			return !unicode.Is(unicode.Categories["C"], c)
		case "Alphabetic":
			return unicode.IsLetter(c) || unicode.Is(unicode.Categories["Nl"], c)
		case "ID_Start":
			return isIdentifierStart(c)
		case "ID_Continue":
			return isIdentifierContinue(c)
		default:
			return unicode.Is(unicode.Properties[parts[1]], c)
		}
	}
	return false
}

func (p *parser) characterClass() (*node, error) {
	if p.flags.UnicodeSets {
		set, e := p.unicodeSetClass()
		if e != nil {
			return nil, e
		}
		return &node{kind: nodeClass, set: set, ignoreCase: p.flags.IgnoreCase}, nil
	}
	neg := false
	if p.peek() == '^' {
		p.take()
		neg = true
	}
	if p.peek() == ']' {
		p.take()
		return &node{kind: nodeClass, negated: neg, ignoreCase: p.flags.IgnoreCase}, nil
	}
	ranges := []classRange{}
	terms := []classTerm{}
	for {
		if p.peek() < 0 {
			return nil, p.errorf("unterminated character class")
		}
		if p.peek() == ']' {
			p.take()
			if len(ranges) > 0 {
				terms = append(terms, classTerm{ranges: ranges})
			}
			return &node{kind: nodeClass, class: terms, negated: neg, ignoreCase: p.flags.IgnoreCase}, nil
		}
		var first *node
		var e error
		if p.peek() == '\\' {
			p.take()
			first, e = p.escape(true)
		} else {
			first = &node{kind: nodeLiteral, rune: p.takePatternCharacter(), ignoreCase: p.flags.IgnoreCase}
		}
		if e != nil {
			return nil, e
		}
		if first.kind != nodeLiteral {
			if first.kind == nodeClass && first.set == nil {
				terms = append(terms, first.class...)
				continue
			}
			return nil, p.errorf("invalid class atom")
		}
		lo, hi := first.rune, first.rune
		if p.peek() == '-' {
			p.take()
			if p.peek() == ']' {
				ranges = append(ranges, classRange{lo, lo}, classRange{'-', '-'})
				continue
			}
			var second *node
			if p.peek() == '\\' {
				p.take()
				second, e = p.escape(true)
			} else {
				second = &node{kind: nodeLiteral, rune: p.takePatternCharacter()}
			}
			if e != nil || second.kind != nodeLiteral {
				return nil, p.errorf("invalid character class range")
			}
			hi = second.rune
			if hi < lo {
				return nil, p.errorf("invalid character class range")
			}
		}
		ranges = append(ranges, classRange{lo, hi})
	}
}
func (p *parser) unicodeSetClass() (*classExpr, error) {
	if p.nesting >= maxParserNesting {
		return nil, p.errorf("regular expression nesting limit exceeded")
	}
	p.nesting++
	defer func() { p.nesting-- }()
	neg := false
	if p.peek() == '^' {
		p.take()
		neg = true
	}
	if p.peek() == ']' {
		p.take()
		empty := &classExpr{kind: classExprUnion}
		if neg {
			return &classExpr{kind: classExprComplement, children: []*classExpr{empty}}, nil
		}
		return empty, nil
	}
	first, _, e := p.unicodeSetOperand(true)
	if e != nil {
		return nil, e
	}
	items := []*classExpr{first}
	operator := rune(0)
	for p.peek() != ']' && p.peek() >= 0 {
		if p.peek() == '&' && p.index+1 < len(p.text) && p.text[p.index+1] == '&' || p.peek() == '-' && p.index+1 < len(p.text) && p.text[p.index+1] == '-' {
			op := p.take()
			p.take()
			if operator != 0 && operator != op {
				return nil, p.errorf("cannot mix set intersection and subtraction")
			}
			operator = op
			right, _, err := p.unicodeSetOperand(false)
			if err != nil {
				return nil, err
			}
			kind := classExprIntersection
			if op == '-' {
				kind = classExprSubtraction
			}
			if len(items) == 1 {
				items = []*classExpr{items[0], right}
			} else {
				items = []*classExpr{{kind: kind, children: items}}
				items = append(items, right)
			}
			continue
		}
		if operator != 0 {
			return nil, p.errorf("set operation operands must be separated")
		}
		right, _, err := p.unicodeSetOperand(true)
		if err != nil {
			return nil, err
		}
		items = append(items, right)
	}
	if p.take() != ']' {
		return nil, p.errorf("unterminated Unicode set class")
	}
	var result *classExpr
	if operator != 0 {
		kind := classExprIntersection
		if operator == '-' {
			kind = classExprSubtraction
		}
		result = &classExpr{kind: kind, children: items}
	} else {
		result = &classExpr{kind: classExprUnion, children: items}
	}
	if neg {
		if classExprContainsStrings(result) {
			return nil, p.errorf("negated Unicode set cannot contain strings")
		}
		result = &classExpr{kind: classExprComplement, children: []*classExpr{result}}
	}
	return result, nil
}
func (p *parser) unicodeSetOperand(allowRange bool) (*classExpr, bool, error) {
	if p.peek() == '[' {
		p.take()
		set, err := p.unicodeSetClass()
		return set, false, err
	}
	var n *node
	var e error
	if p.peek() == '\\' {
		p.take()
		n, e = p.escape(true)
	} else {
		c := p.takePatternCharacter()
		if strings.ContainsRune("-()[]{}|/\\", c) {
			return nil, false, p.errorf("reserved character %q in Unicode set", c)
		}
		n = &node{kind: nodeLiteral, rune: c, ignoreCase: p.flags.IgnoreCase}
	}
	if e != nil {
		return nil, false, e
	}
	if n.kind == nodeClass {
		if n.set != nil {
			return n.set, false, nil
		}
		children := []*classExpr{}
		for _, term := range n.class {
			children = append(children, &classExpr{kind: classExprTerm, term: term})
		}
		return &classExpr{kind: classExprUnion, children: children}, false, nil
	}
	if n.kind != nodeLiteral {
		return nil, false, p.errorf("invalid Unicode set operand")
	}
	lo, hi := n.rune, n.rune
	isRange := false
	if allowRange && p.peek() == '-' && !(p.index+1 < len(p.text) && p.text[p.index+1] == '-') && p.index+1 < len(p.text) && p.text[p.index+1] != ']' {
		p.take()
		isRange = true
		var last *node
		if p.peek() == '\\' {
			p.take()
			last, e = p.escape(true)
		} else {
			last = &node{kind: nodeLiteral, rune: p.takePatternCharacter()}
		}
		if e != nil || last.kind != nodeLiteral {
			return nil, false, p.errorf("invalid Unicode set range")
		}
		hi = last.rune
		if hi < lo {
			return nil, false, p.errorf("invalid Unicode set range")
		}
	}
	return &classExpr{kind: classExprTerm, term: classTerm{ranges: []classRange{{lo, hi}}}}, isRange, nil
}
func classExprContainsStrings(expr *classExpr) bool {
	if expr == nil {
		return false
	}
	if expr.kind == classExprStringSet && len(expr.strings) > 0 {
		return true
	}
	for _, c := range expr.children {
		if classExprContainsStrings(c) {
			return true
		}
	}
	return false
}

// computeFirstCharacterFilter builds a conservative set of AST leaves that
// can consume the first input character. A nil result means the AST contains
// a form we cannot prove safe to filter (or can match without consuming).
// This is intentionally fail-open: it may miss optimization opportunities,
// but must never discard a valid starting position.
func computeFirstCharacterFilter(root *node) []*node {
	var first func(*node) (leaves []*node, nullable, known bool)
	first = func(n *node) ([]*node, bool, bool) {
		if n == nil {
			return nil, true, true
		}
		switch n.kind {
		case nodeLiteral, nodeClass:
			if n.kind == nodeClass && n.set != nil && !classExprCanFilter(n.set) {
				return nil, false, false
			}
			return []*node{n}, false, true
		case nodeAny, nodeBackref:
			// Any is not selective. A backreference's first character depends
			// on a runtime capture, so it cannot safely contribute a static set.
			return nil, n.kind == nodeBackref, false
		case nodeEmpty, nodeBegin, nodeEnd, nodeBoundary, nodeLookaround:
			return nil, true, true
		case nodeCapture:
			return first(n.children[0])
		case nodeAlternate:
			var out []*node
			nullable := false
			for _, child := range n.children {
				leaves, canBeEmpty, ok := first(child)
				if !ok {
					return nil, false, false
				}
				out = append(out, leaves...)
				nullable = nullable || canBeEmpty
			}
			return out, nullable, true
		case nodeConcat:
			var out []*node
			nullable := true
			for _, child := range n.children {
				if !nullable {
					break
				}
				leaves, canBeEmpty, ok := first(child)
				if !ok {
					return nil, false, false
				}
				out = append(out, leaves...)
				nullable = nullable && canBeEmpty
			}
			return out, nullable, true
		case nodeRepeat:
			if n.max == 0 {
				return nil, true, true
			}
			leaves, bodyNullable, ok := first(n.children[0])
			if !ok {
				return nil, false, false
			}
			return leaves, n.min == 0 || bodyNullable, true
		}
		return nil, false, false
	}
	leaves, nullable, known := first(root)
	if !known || nullable || len(leaves) == 0 {
		return nil
	}
	return leaves
}

func classExprCanFilter(expr *classExpr) bool {
	if expr == nil || classExprContainsStrings(expr) {
		return false
	}
	switch expr.kind {
	case classExprTerm:
		return true
	case classExprUnion, classExprIntersection, classExprSubtraction, classExprComplement:
		for _, child := range expr.children {
			if !classExprCanFilter(child) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func literalAlternativeSequence(n *node) ([]*node, bool) {
	switch n.kind {
	case nodeLiteral:
		return []*node{n}, true
	case nodeConcat:
		sequence := make([]*node, 0, len(n.children))
		for _, child := range n.children {
			if child.kind != nodeLiteral {
				return nil, false
			}
			sequence = append(sequence, child)
		}
		return sequence, true
	default:
		return nil, false
	}
}

func literalSequenceNode(sequence []*node) *node {
	switch len(sequence) {
	case 0:
		return &node{kind: nodeEmpty}
	case 1:
		return sequence[0]
	default:
		return &node{kind: nodeConcat, children: sequence}
	}
}

// optimizeLiteralAlternatives applies two semantics-preserving graph-building
// rewrites from V8's Disjunction lowering: single-code-point alternatives
// become a class, and a shared literal prefix is emitted once before the
// remaining alternatives. It deliberately keeps source order and only handles
// literal-only branches, avoiding capture/flag-boundary changes.
func optimizeLiteralAlternatives(root *node) *node {
	if root == nil {
		return nil
	}
	for i, child := range root.children {
		root.children[i] = optimizeLiteralAlternatives(child)
	}
	if root.kind != nodeAlternate || len(root.children) < 2 {
		return root
	}

	allSingle := true
	ignoreCase := root.children[0].ignoreCase
	ranges := make([]classRange, 0, len(root.children))
	for _, child := range root.children {
		if child.kind != nodeLiteral || child.ignoreCase != ignoreCase {
			allSingle = false
			break
		}
		ranges = append(ranges, classRange{child.rune, child.rune})
	}
	if allSingle {
		return &node{kind: nodeClass, class: []classTerm{{ranges: ranges}}, ignoreCase: ignoreCase}
	}

	sequences := make([][]*node, len(root.children))
	minLength := -1
	for i, child := range root.children {
		sequence, ok := literalAlternativeSequence(child)
		if !ok || len(sequence) == 0 {
			return root
		}
		sequences[i] = sequence
		if minLength < 0 || len(sequence) < minLength {
			minLength = len(sequence)
		}
	}
	prefixLength := 0
	for prefixLength < minLength {
		first := sequences[0][prefixLength]
		common := true
		for i := 1; i < len(sequences); i++ {
			other := sequences[i][prefixLength]
			if first.rune != other.rune || first.ignoreCase != other.ignoreCase {
				common = false
				break
			}
		}
		if !common {
			break
		}
		prefixLength++
	}
	if prefixLength == 0 {
		return root
	}

	suffixes := make([]*node, len(sequences))
	allSuffixesEmpty := true
	for i, sequence := range sequences {
		suffixes[i] = literalSequenceNode(sequence[prefixLength:])
		allSuffixesEmpty = allSuffixesEmpty && len(sequence) == prefixLength
	}
	prefix := literalSequenceNode(sequences[0][:prefixLength])
	if allSuffixesEmpty {
		return prefix
	}
	suffix := optimizeLiteralAlternatives(&node{kind: nodeAlternate, children: suffixes})
	if suffix.kind == nodeEmpty {
		return prefix
	}
	if prefix.kind == nodeConcat {
		return &node{kind: nodeConcat, children: append(prefix.children, suffix)}
	}
	return &node{kind: nodeConcat, children: []*node{prefix, suffix}}
}

func (p *Program) firstCharacterMatches(text []uint16, pos int) bool {
	if len(p.firstChars) == 0 {
		return true
	}
	if pos < 0 || pos >= len(text) {
		return false
	}
	value := inputCharacter(text, pos, p.Flags.Unicode).value
	if value <= 0xff {
		return p.firstASCII[value>>6]&(uint64(1)<<(value&63)) != 0
	}
	return p.firstCharacterValueMatches(value)
}

func (p *Program) firstCharacterValueMatches(value rune) bool {
	for _, n := range p.firstChars {
		switch n.kind {
		case nodeLiteral:
			if equalRune(value, n.rune, n.ignoreCase, p.Flags.Unicode) {
				return true
			}
		case nodeClass:
			if n.set != nil {
				if classExprContains(value, n.set, n.ignoreCase, p.Flags.Unicode, p.Flags.UnicodeSets) {
					return true
				}
			} else if classContains(value, n.class, n.negated, n.ignoreCase, p.Flags.Unicode, p.Flags.UnicodeSets) {
				return true
			}
		}
	}
	return false
}

func isIdentifierStart(c rune) bool {
	return c == '$' || c == '_' || unicode.IsLetter(c) || unicode.Is(unicode.Categories["Nl"], c) || unicode.Is(unicode.Properties["Other_ID_Start"], c)
}
func isIdentifierContinue(c rune) bool {
	return isIdentifierStart(c) || unicode.Is(unicode.Categories["Mn"], c) || unicode.Is(unicode.Categories["Mc"], c) || unicode.IsDigit(c) || unicode.Is(unicode.Categories["Pc"], c) || unicode.Is(unicode.Properties["Other_ID_Continue"], c) || c == 0x200c || c == 0x200d
}
func validCaptureName(s string) bool {
	for i, c := range s {
		if i == 0 {
			if !isIdentifierStart(c) {
				return false
			}
		} else if !isIdentifierContinue(c) {
			return false
		}
	}
	return s != ""
}

func capturePathsAreExclusive(left, right []choiceBranch) bool {
	for _, a := range left {
		for _, b := range right {
			if a.id == b.id && a.branch != b.branch {
				return true
			}
		}
	}
	return false
}

func literalSequence(s string, ignore bool) *node {
	children := []*node{}
	for _, c := range s {
		children = append(children, &node{kind: nodeLiteral, rune: c, ignoreCase: ignore})
	}
	return &node{kind: nodeConcat, children: children}
}
func resolveReferences(root *node, p *parser) (*node, error) {
	for i, c := range root.children {
		n, e := resolveReferences(c, p)
		if e != nil {
			return nil, e
		}
		root.children[i] = n
	}
	if root.kind != nodeBackref {
		return root, nil
	}
	if root.refName != "" {
		if ids, ok := p.groupNames[root.refName]; ok {
			root.captureIndices = append([]int(nil), ids...)
			if len(ids) != 0 {
				root.capture = ids[0]
			}
			return root, nil
		}
		if p.flags.Unicode || p.hasNamedCaptures {
			return nil, fmt.Errorf("regular expression: invalid named backreference %q", root.refName)
		}
		return literalSequence("k<"+root.refName+">", root.ignoreCase), nil
	}
	if root.capture > 0 && root.capture <= p.group {
		return root, nil
	}
	if p.flags.Unicode {
		return nil, fmt.Errorf("regular expression: invalid numeric backreference \\%s", root.refDigits)
	}
	if root.refDigits == "" {
		return root, nil
	}
	digits := root.refDigits
	for length := len(digits); length > 0; length-- {
		id, _ := strconv.Atoi(digits[:length])
		if id > 0 && id <= p.group {
			root.capture = id
			nodes := []*node{root}
			for _, d := range digits[length:] {
				nodes = append(nodes, &node{kind: nodeLiteral, rune: d, ignoreCase: root.ignoreCase})
			}
			if len(nodes) == 1 {
				return root, nil
			}
			return &node{kind: nodeConcat, children: nodes}, nil
		}
	}
	value := 0
	length := 0
	for length < len(digits) && length < 3 && digits[length] >= '0' && digits[length] <= '7' {
		next := value*8 + int(digits[length]-'0')
		if next > 255 {
			break
		}
		value = next
		length++
	}
	nodes := []*node{{kind: nodeLiteral, rune: rune(value), ignoreCase: root.ignoreCase}}
	for _, d := range digits[length:] {
		nodes = append(nodes, &node{kind: nodeLiteral, rune: d, ignoreCase: root.ignoreCase})
	}
	return &node{kind: nodeConcat, children: nodes}, nil
}

func captureIDs(root *node) []int {
	out := []int{}
	var walk func(*node)
	walk = func(n *node) {
		if n.kind == nodeCapture && n.capture > 0 {
			out = append(out, n.capture)
		}
		for _, c := range n.children {
			walk(c)
		}
	}
	walk(root)
	return out
}

type op uint8

const (
	opChar op = iota
	opAny
	opClass
	opBackref
	opSplit
	opJump
	opSaveStart
	opSaveEnd
	opClearCaptures
	opBegin
	opEnd
	opBoundary
	opAssert
	opMatch
	opGreedyLoop
)

type instruction struct {
	op                                     op
	rune                                   rune
	class                                  []classTerm
	set                                    *classExpr
	negated, ignoreCase, dotAll, multiline bool
	a, b                                   int
	capture                                int
	guard                                  bool
	sub                                    *Program
	groups                                 []int
	captureIndices                         []int
	negative, behind                       bool
	loop                                   *instruction
}
type Span struct{ Start, End int }
type Match struct {
	Start, End int
	Captures   []Span
}

// CaptureUTF16 returns a view of a participating capture in the original
// UTF-16 subject. The returned slice aliases subject and remains valid only
// while the caller keeps that subject unchanged. The boolean distinguishes an
// unmatched or missing capture from a participating empty capture.
func (m Match) CaptureUTF16(subject []uint16, index int) ([]uint16, bool) {
	if index < 0 || index >= len(m.Captures) {
		return nil, false
	}
	span := m.Captures[index]
	if span.Start < 0 || span.End < span.Start || span.End > len(subject) {
		return nil, false
	}
	return subject[span.Start:span.End], true
}

type Program struct {
	Source       string
	Flags        Flags
	source       []uint16
	code         []instruction
	atom         []uint16
	atomFailure  []int
	firstChars   []*node
	firstASCII   [4]uint64
	captures     int
	GroupIndices map[string]int
	// GroupCaptureIndices includes every statically distinct capture that may
	// provide a duplicate named group. GroupIndices keeps the first index for
	// callers that only need a stable representative.
	GroupCaptureIndices map[string][]int
	reverse             bool
	linear              *linearProgram
}

func Compile(source, flags string) (*Program, error) {
	return CompileUTF16(utf16.Encode([]rune(source)), flags)
}
func CompileUTF16(source []uint16, flags string) (*Program, error) {
	f, e := ParseFlags(flags)
	if e != nil {
		return nil, e
	}
	p := &parser{text: append([]uint16(nil), source...), flags: f, groupNames: map[string][]int{}}
	root, e := p.parse()
	if e != nil {
		return nil, e
	}
	root, e = resolveReferences(root, p)
	if e != nil {
		return nil, e
	}
	root = optimizeLiteralAlternatives(root)
	normalizeClassRanges(root)
	program := &Program{Source: utf16String(source), Flags: f, source: append([]uint16(nil), source...), captures: p.group, GroupIndices: map[string]int{}, GroupCaptureIndices: map[string][]int{}}
	for name, indices := range p.groupNames {
		if len(indices) > 0 {
			program.GroupIndices[name] = indices[0]
		}
		program.GroupCaptureIndices[name] = append([]int(nil), indices...)
	}
	if f.Linear {
		program.linear, e = compileLinear(root, f, p.group)
		if e != nil {
			return nil, e
		}
		return program, nil
	}
	if p.group == 0 && !f.IgnoreCase && !f.Sticky && !IsEitherUnicode(f) {
		if atom, failure, ok := compileLiteralAtom(root); ok {
			program.atom, program.atomFailure = atom, failure
			return program, nil
		}
	}
	program.firstChars = computeFirstCharacterFilter(root)
	for c := rune(0); c < 256; c++ {
		if program.firstCharacterValueMatches(c) {
			program.firstASCII[c>>6] |= uint64(1) << (c & 63)
		}
	}
	program.code = compileNode(root, p.group, f, false)
	return program, nil
}

func compileNode(root *node, captures int, flags Flags, reverse bool) []instruction {
	code := []instruction{}
	var emit func(*node) (int, int)
	emit = func(n *node) (int, int) {
		start := len(code)
		switch n.kind {
		case nodeEmpty:
		case nodeLiteral:
			code = append(code, instruction{op: opChar, rune: n.rune, ignoreCase: n.ignoreCase})
		case nodeAny:
			code = append(code, instruction{op: opAny, dotAll: n.dotAll})
		case nodeClass:
			code = append(code, instruction{op: opClass, class: n.class, set: n.set, negated: n.negated, ignoreCase: n.ignoreCase})
		case nodeBackref:
			code = append(code, instruction{op: opBackref, capture: n.capture, captureIndices: append([]int(nil), n.captureIndices...), ignoreCase: n.ignoreCase})
		case nodeLookaround:
			body := n.children[0]
			sub := &Program{Flags: flags, captures: captures, reverse: n.behind, code: compileNode(body, captures, flags, n.behind)}
			code = append(code, instruction{op: opAssert, sub: sub, groups: captureIDs(body), negative: n.negative, behind: n.behind})
		case nodeBegin:
			code = append(code, instruction{op: opBegin, multiline: n.multiline})
		case nodeEnd:
			code = append(code, instruction{op: opEnd, multiline: n.multiline})
		case nodeBoundary:
			code = append(code, instruction{op: opBoundary, negated: n.negated, ignoreCase: n.ignoreCase})
		case nodeCapture:
			code = append(code, instruction{op: opSaveStart, capture: n.capture})
			emit(n.children[0])
			code = append(code, instruction{op: opSaveEnd, capture: n.capture})
		case nodeConcat:
			if reverse {
				for i := len(n.children) - 1; i >= 0; i-- {
					emit(n.children[i])
				}
			} else {
				for _, child := range n.children {
					emit(child)
				}
			}
		case nodeAlternate:
			if len(n.children) == 0 {
				break
			}
			split := len(code)
			code = append(code, instruction{op: opSplit})
			left := len(code)
			emit(n.children[0])
			jump := len(code)
			code = append(code, instruction{op: opJump})
			right := len(code)
			if len(n.children) == 2 {
				emit(n.children[1])
			} else {
				emit(&node{kind: nodeAlternate, children: n.children[1:]})
			}
			code[split].a, code[split].b = left, jump+1
			code[jump].a = len(code)
			_ = right
		case nodeRepeat:
			groups := captureIDs(n.children[0])
			for i := 0; i < n.min; i++ {
				if i > 0 && n.max != 1 && len(groups) > 0 {
					code = append(code, instruction{op: opClearCaptures, groups: groups})
				}
				emit(n.children[0])
			}
			if n.max < 0 && n.greedy {
				if body, ok := repeatBodyInstruction(n.children[0]); ok {
					code = append(code, instruction{op: opGreedyLoop, loop: body})
					break
				}
			}
			if n.max < 0 {
				loop := len(code)
				code = append(code, instruction{op: opSplit, guard: true})
				body := len(code)
				if len(groups) > 0 {
					code = append(code, instruction{op: opClearCaptures, groups: groups})
				}
				emit(n.children[0])
				code = append(code, instruction{op: opJump, a: loop})
				exit := len(code)
				if n.greedy {
					code[loop].a, code[loop].b = body, exit
				} else {
					code[loop].a, code[loop].b = exit, body
				}
			} else {
				for i := n.min; i < n.max; i++ {
					split := len(code)
					code = append(code, instruction{op: opSplit})
					body := len(code)
					if i > 0 && len(groups) > 0 {
						code = append(code, instruction{op: opClearCaptures, groups: groups})
					}
					emit(n.children[0])
					if n.greedy {
						code[split].a, code[split].b = body, len(code)
					} else {
						code[split].a, code[split].b = len(code), body
					}
				}
			}
		}
		return start, len(code)
	}
	emit(root)
	code = append(code, instruction{op: opMatch})
	return code
}

// repeatBodyInstruction recognizes deterministic, single-code-point bodies
// whose greedy repetition can be scanned directly without one VM frame per
// iteration. Captures and Unicode-set string alternatives remain on the
// general backtracking path.
func repeatBodyInstruction(n *node) (*instruction, bool) {
	var body instruction
	switch n.kind {
	case nodeLiteral:
		body = instruction{op: opChar, rune: n.rune, ignoreCase: n.ignoreCase}
	case nodeAny:
		body = instruction{op: opAny, dotAll: n.dotAll}
	case nodeClass:
		if n.set != nil {
			return nil, false
		}
		body = instruction{op: opClass, class: n.class, negated: n.negated, ignoreCase: n.ignoreCase}
	default:
		return nil, false
	}
	return &body, true
}

type inputChar struct {
	value rune
	width int
}

func inputCharacter(text []uint16, pos int, unicodeMode bool) inputChar {
	c := rune(text[pos])
	if unicodeMode && c >= 0xd800 && c <= 0xdbff && pos+1 < len(text) {
		d := rune(text[pos+1])
		if d >= 0xdc00 && d <= 0xdfff {
			return inputChar{0x10000 + ((c - 0xd800) << 10) + d - 0xdc00, 2}
		}
	}
	return inputChar{c, 1}
}
func previousInputCharacter(text []uint16, pos int, unicodeMode bool) inputChar {
	c := rune(text[pos-1])
	if unicodeMode && c >= 0xdc00 && c <= 0xdfff && pos > 1 {
		d := rune(text[pos-2])
		if d >= 0xd800 && d <= 0xdbff {
			return inputChar{0x10000 + ((d - 0xd800) << 10) + c - 0xdc00, 2}
		}
	}
	return inputChar{c, 1}
}
func previousCharacter(text []uint16, pos int, unicodeMode bool) rune {
	return previousInputCharacter(text, pos, unicodeMode).value
}
func advanceStringIndex(text []uint16, pos int, unicodeMode bool) int {
	if pos >= len(text) {
		return pos + 1
	}
	return pos + inputCharacter(text, pos, unicodeMode).width
}

func equalRune(a, b rune, ignoreCase, unicodeMode bool) bool {
	if a == b {
		return true
	}
	if !ignoreCase {
		return false
	}
	if unicodeMode {
		return forEachUnicodeCaseEquivalent(a, func(c rune) bool { return c == b })
	}
	return canonicalizeLegacy(a) == canonicalizeLegacy(b)
}
func canonicalizeLegacy(c rune) rune {
	upper := unicode.ToUpper(c)
	if upper > 0xff || c >= 0x80 && upper < 0x80 {
		return c
	}
	return upper
}
func classTermContains(value rune, term classTerm, ignoreCase, unicodeMode, unicodeSets bool) bool {
	if term.property != "" {
		contains := func(c rune) bool { return unicodePropertyContains(c, term.property) }
		matched := contains(value)
		if ignoreCase && unicodeMode {
			if term.negated && !unicodeSets {
				return forEachUnicodeCaseEquivalent(value, func(c rune) bool { return !contains(c) })
			}
			matched = forEachUnicodeCaseEquivalent(value, contains)
		}
		return matched != term.negated
	}
	contains := func(c rune) bool {
		index := sort.Search(len(term.ranges), func(i int) bool { return term.ranges[i].hi >= c })
		return index < len(term.ranges) && term.ranges[index].lo <= c
	}
	matched := contains(value)
	if ignoreCase && unicodeMode {
		matched = forEachUnicodeCaseEquivalent(value, contains)
	} else if ignoreCase {
		matched = forEachLegacyCaseEquivalent(value, contains)
	}
	if term.negated {
		return !matched
	}
	return matched
}

func normalizeClassRanges(root *node) {
	var normalizeTerm func(*classTerm)
	normalizeTerm = func(term *classTerm) {
		if len(term.ranges) < 2 {
			return
		}
		sort.Slice(term.ranges, func(i, j int) bool {
			if term.ranges[i].lo == term.ranges[j].lo {
				return term.ranges[i].hi < term.ranges[j].hi
			}
			return term.ranges[i].lo < term.ranges[j].lo
		})
		merged := term.ranges[:0]
		for _, current := range term.ranges {
			if len(merged) == 0 || current.lo > merged[len(merged)-1].hi+1 {
				merged = append(merged, current)
				continue
			}
			if current.hi > merged[len(merged)-1].hi {
				merged[len(merged)-1].hi = current.hi
			}
		}
		term.ranges = merged
	}
	var visitExpr func(*classExpr)
	visitExpr = func(expr *classExpr) {
		if expr == nil {
			return
		}
		if expr.kind == classExprTerm {
			normalizeTerm(&expr.term)
		}
		for _, child := range expr.children {
			visitExpr(child)
		}
	}
	var visitNode func(*node)
	visitNode = func(current *node) {
		if current == nil {
			return
		}
		for i := range current.class {
			normalizeTerm(&current.class[i])
		}
		visitExpr(current.set)
		for _, child := range current.children {
			visitNode(child)
		}
	}
	visitNode(root)
}
func classContains(value rune, terms []classTerm, negated, ignoreCase, unicodeMode, unicodeSets bool) bool {
	matched := false
	for _, t := range terms {
		if classTermContains(value, t, ignoreCase, unicodeMode, unicodeSets) {
			matched = true
			break
		}
	}
	return matched != negated
}
func classExprContains(value rune, e *classExpr, ignoreCase, unicodeMode, unicodeSets bool) bool {
	switch e.kind {
	case classExprUnion:
		for _, c := range e.children {
			if classExprContains(value, c, ignoreCase, unicodeMode, unicodeSets) {
				return true
			}
		}
		return false
	case classExprIntersection:
		return classExprContains(value, e.children[0], ignoreCase, unicodeMode, unicodeSets) && classExprContains(value, e.children[1], ignoreCase, unicodeMode, unicodeSets)
	case classExprSubtraction:
		return classExprContains(value, e.children[0], ignoreCase, unicodeMode, unicodeSets) && !classExprContains(value, e.children[1], ignoreCase, unicodeMode, unicodeSets)
	case classExprComplement:
		return !classExprContains(value, e.children[0], ignoreCase, unicodeMode, unicodeSets)
	case classExprStringSet:
		return false
	case classExprTerm:
		return classTermContains(value, e.term, ignoreCase, unicodeMode, unicodeSets)
	}
	return false
}
func (p *Program) classEndsAt(text []uint16, pos int, ins instruction) []int {
	if ins.set != nil {
		return classExprEnds(text, pos, ins.set, p.Flags, p.reverse)
	}
	end, ok := p.classEndAt(text, pos, ins)
	if !ok {
		return nil
	}
	return []int{end}
}

// classEndAt handles the overwhelmingly common one-code-point class case
// without allocating the one-element slice needed by Unicode-set alternatives.
func (p *Program) classEndAt(text []uint16, pos int, ins instruction) (int, bool) {
	if pos < 0 || pos > len(text) {
		return 0, false
	}
	var ch inputChar
	if p.reverse {
		if pos == 0 {
			return 0, false
		}
		ch = previousInputCharacter(text, pos, p.Flags.Unicode)
	} else {
		if pos == len(text) {
			return 0, false
		}
		ch = inputCharacter(text, pos, p.Flags.Unicode)
	}
	if !classContains(ch.value, ins.class, ins.negated, ins.ignoreCase, p.Flags.Unicode, p.Flags.UnicodeSets) {
		return 0, false
	}
	if p.reverse {
		return pos - ch.width, true
	}
	return pos + ch.width, true
}

func classExprEnds(text []uint16, pos int, e *classExpr, flags Flags, reverse bool) []int {
	switch e.kind {
	case classExprUnion:
		ends := []int{}
		for _, c := range e.children {
			ends = append(ends, classExprEnds(text, pos, c, flags, reverse)...)
		}
		return uniqueOrderedEnds(ends, reverse)
	case classExprIntersection, classExprSubtraction:
		left := classExprEnds(text, pos, e.children[0], flags, reverse)
		right := map[int]bool{}
		for _, n := range classExprEnds(text, pos, e.children[1], flags, reverse) {
			right[n] = true
		}
		out := []int{}
		for _, n := range left {
			if right[n] == (e.kind == classExprIntersection) {
				out = append(out, n)
			}
		}
		return uniqueOrderedEnds(out, reverse)
	case classExprComplement:
		var ch inputChar
		if reverse {
			if pos <= 0 {
				return nil
			}
			ch = previousInputCharacter(text, pos, flags.Unicode)
		} else {
			if pos >= len(text) {
				return nil
			}
			ch = inputCharacter(text, pos, flags.Unicode)
		}
		end := pos + ch.width
		if reverse {
			end = pos - ch.width
		}
		for _, n := range classExprEnds(text, pos, e.children[0], flags, reverse) {
			if n == end {
				return nil
			}
		}
		return []int{end}
	case classExprStringSet:
		out := []int{}
		for _, s := range e.strings {
			if end, ok := matchClassString(text, pos, s, flags, reverse); ok {
				out = append(out, end)
			}
		}
		return uniqueOrderedEnds(out, reverse)
	case classExprTerm:
		var ch inputChar
		if reverse {
			if pos <= 0 {
				return nil
			}
			ch = previousInputCharacter(text, pos, flags.Unicode)
		} else {
			if pos >= len(text) {
				return nil
			}
			ch = inputCharacter(text, pos, flags.Unicode)
		}
		if !classTermContains(ch.value, e.term, flags.IgnoreCase, flags.Unicode, flags.UnicodeSets) {
			return nil
		}
		end := pos + ch.width
		if reverse {
			end = pos - ch.width
		}
		return []int{end}
	}
	return nil
}
func matchClassString(text []uint16, pos int, s []rune, flags Flags, reverse bool) (int, bool) {
	for i := 0; i < len(s); i++ {
		if reverse {
			if pos <= 0 {
				return pos, false
			}
			c := previousInputCharacter(text, pos, flags.Unicode)
			if !equalRune(c.value, s[len(s)-i-1], flags.IgnoreCase, flags.Unicode) {
				return pos, false
			}
			pos -= c.width
		} else {
			if pos >= len(text) {
				return pos, false
			}
			c := inputCharacter(text, pos, flags.Unicode)
			if !equalRune(c.value, s[i], flags.IgnoreCase, flags.Unicode) {
				return pos, false
			}
			pos += c.width
		}
	}
	return pos, true
}
func uniqueOrderedEnds(ends []int, reverse bool) []int {
	seen := map[int]bool{}
	out := []int{}
	for _, n := range ends {
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if reverse {
			return out[i] < out[j]
		}
		return out[i] > out[j]
	})
	return out
}

type splitVisit struct{ pc, position int }
type executionState struct {
	pc, position int
	captures     []int
	emptySplits  map[splitVisit]struct{}
	loopStarts   map[int]int
	guardExit    splitVisit
	hasGuardExit bool
	greedyRetry  bool
	loopStart    int
}

func cloneIntSlice(v []int) []int { return append([]int(nil), v...) }
func cloneSplitVisits(v map[splitVisit]struct{}) map[splitVisit]struct{} {
	out := make(map[splitVisit]struct{}, len(v))
	for k := range v {
		out[k] = struct{}{}
	}
	return out
}
func cloneLoopStarts(v map[int]int) map[int]int {
	out := make(map[int]int, len(v))
	for pc, position := range v {
		out[pc] = position
	}
	return out
}

func (p *Program) consumeAt(text []uint16, pos, pc int, captures []int) (int, bool) {
	ins := p.code[pc]
	if ins.op == opBackref {
		if len(ins.captureIndices) != 0 {
			// Duplicate named groups are permitted only on mutually exclusive
			// paths, so at most one of their capture slots can be set here.
			// JavaScript treats a named backreference as empty when none of the
			// alternatives containing that name participated.
			for _, capture := range ins.captureIndices {
				a, b := captures[capture*2], captures[capture*2+1]
				if a >= 0 && b >= 0 {
					return p.consumeBackreferenceAt(text, pos, ins, a, b)
				}
			}
			return pos, true
		}
		a, b := captures[ins.capture*2], captures[ins.capture*2+1]
		if a < 0 || b < 0 {
			return pos, true
		}
		return p.consumeBackreferenceAt(text, pos, ins, a, b)
	}
	return p.consumeSimpleAt(text, pos, ins)
}

func (p *Program) consumeBackreferenceAt(text []uint16, pos int, ins instruction, a, b int) (int, bool) {
	if a > b {
		a, b = b, a
	}
	at := pos
	if p.reverse {
		for q := b; q > a; {
			if at <= 0 {
				return pos, false
			}
			left := previousInputCharacter(text, q, p.Flags.Unicode)
			right := previousInputCharacter(text, at, p.Flags.Unicode)
			if !equalRune(left.value, right.value, ins.ignoreCase, p.Flags.Unicode) {
				return pos, false
			}
			if p.Flags.Unicode && (b-q)+left.width != (pos-at)+right.width {
				return pos, false
			}
			q -= left.width
			at -= right.width
		}
	} else {
		for q := a; q < b; {
			if at >= len(text) {
				return pos, false
			}
			left := inputCharacter(text, q, p.Flags.Unicode)
			right := inputCharacter(text, at, p.Flags.Unicode)
			if !equalRune(left.value, right.value, ins.ignoreCase, p.Flags.Unicode) {
				return pos, false
			}
			if p.Flags.Unicode && (q-a)+left.width != (at-pos)+right.width {
				return pos, false
			}
			q += left.width
			at += right.width
		}
	}
	return at, true
}

func (p *Program) consumeSimpleAt(text []uint16, pos int, ins instruction) (int, bool) {
	if ins.op == opClass {
		return p.classEndAt(text, pos, ins)
	}
	var c inputChar
	if p.reverse {
		if pos <= 0 {
			return pos, false
		}
		c = previousInputCharacter(text, pos, p.Flags.Unicode)
	} else {
		if pos >= len(text) {
			return pos, false
		}
		c = inputCharacter(text, pos, p.Flags.Unicode)
	}
	next := pos + c.width
	if p.reverse {
		next = pos - c.width
	}
	switch ins.op {
	case opChar:
		if equalRune(c.value, ins.rune, ins.ignoreCase, p.Flags.Unicode) {
			return next, true
		}
	case opAny:
		if ins.dotAll || !isLineTerminator(c.value) {
			return next, true
		}
	}
	return pos, false
}

func (p *Program) previousGreedyPosition(text []uint16, pos, start int) (int, bool) {
	if p.reverse {
		if pos >= start || pos < 0 || pos >= len(text) {
			return 0, false
		}
		previous := pos + inputCharacter(text, pos, p.Flags.Unicode).width
		return previous, previous <= start
	}
	if pos <= start || pos > len(text) || pos == 0 {
		return 0, false
	}
	previous := pos - previousInputCharacter(text, pos, p.Flags.Unicode).width
	return previous, previous >= start
}

func (p *Program) runAt(text []uint16, start int) (Match, bool) {
	return p.runAtWithCaptures(text, start, nil, -1)
}
func (p *Program) runAtWithCaptures(text []uint16, start int, initial []int, requiredEnd int) (Match, bool) {
	return p.runAtWithExecutionBudget(text, start, initial, requiredEnd, nil)
}

func (p *Program) runAtWithExecutionBudget(text []uint16, start int, initial []int, requiredEnd int, budget *executionBudget) (Match, bool) {
	baseBacktrackBytes := uint64(0)
	activeBytes := uint64(0)
	if budget != nil {
		baseBacktrackBytes = budget.backtrackBytes
		activeBytes = executionFrameBytes(2*(p.captures+1), 0, 0)
		if budget.adjustAuxiliary(0, activeBytes) != nil {
			return Match{}, false
		}
		defer func() {
			budget.adjustAuxiliary(activeBytes, 0)
			budget.backtrackBytes = baseBacktrackBytes
		}()
	}
	captures := make([]int, 2*(p.captures+1))
	if initial == nil {
		for i := range captures {
			captures[i] = -1
		}
		captures[0] = start
	} else {
		copy(captures, initial)
	}
	pc, pos := 0, start
	empty := map[splitVisit]struct{}{}
	loopStarts := map[int]int{}
	backtrack := []executionState{}
	fail := func() bool {
		for len(backtrack) != 0 && (budget == nil || budget.err == nil) {
			last := len(backtrack) - 1
			s := backtrack[last]
			backtrack = backtrack[:last]
			budget.popState(s)
			if s.greedyRetry {
				previous, ok := p.previousGreedyPosition(text, s.position, s.loopStart)
				if !ok {
					continue
				}
				if budget.adjustAuxiliary(activeBytes, executionFrameBytes(len(s.captures), len(s.emptySplits), len(s.loopStarts))) != nil {
					return false
				}
				activeCaptures := cloneIntSlice(s.captures)
				activeEmpty := cloneSplitVisits(s.emptySplits)
				activeLoopStarts := cloneLoopStarts(s.loopStarts)
				activeBytes = executionFrameBytes(len(activeCaptures), len(activeEmpty), len(activeLoopStarts))
				s.position = previous
				if budget.pushState(s) != nil {
					return false
				}
				backtrack = append(backtrack, s)
				pc, pos = s.pc, s.position
				captures, empty, loopStarts = activeCaptures, activeEmpty, activeLoopStarts
				return true
			}
			if budget.adjustAuxiliary(activeBytes, executionStateBytes(s)) != nil {
				return false
			}
			activeBytes = executionStateBytes(s)
			pc, pos, captures, empty, loopStarts = s.pc, s.position, s.captures, s.emptySplits, s.loopStarts
			return true
		}
		return false
	}
	for {
		if budget.step() != nil {
			return Match{}, false
		}
		if pc < 0 || pc >= len(p.code) {
			if !fail() {
				return Match{}, false
			}
			continue
		}
		ins := p.code[pc]
		switch ins.op {
		case opClass:
			if ins.set == nil {
				next, ok := p.classEndAt(text, pos, ins)
				if !ok {
					if !fail() {
						return Match{}, false
					}
					continue
				}
				pc++
				pos = next
				continue
			}
			ends := p.classEndsAt(text, pos, ins)
			if len(ends) == 0 {
				if !fail() {
					return Match{}, false
				}
				continue
			}
			for i := len(ends) - 1; i >= 1; i-- {
				if budget.checkStateSize(len(captures), len(empty), len(loopStarts)) != nil {
					return Match{}, false
				}
				state := executionState{pc: pc + 1, position: ends[i], captures: cloneIntSlice(captures), emptySplits: cloneSplitVisits(empty), loopStarts: cloneLoopStarts(loopStarts)}
				if budget.pushState(state) != nil {
					return Match{}, false
				}
				backtrack = append(backtrack, state)
			}
			pc++
			pos = ends[0]
		case opChar, opAny, opBackref:
			next, ok := p.consumeAt(text, pos, pc, captures)
			if !ok {
				if !fail() {
					return Match{}, false
				}
				continue
			}
			pc++
			pos = next
		case opGreedyLoop:
			if ins.loop == nil {
				panic("greedy loop missing body instruction")
			}
			loopStart := pos
			for {
				if budget.step() != nil {
					return Match{}, false
				}
				next, ok := p.consumeSimpleAt(text, pos, *ins.loop)
				if !ok {
					break
				}
				pos = next
			}
			if pos != loopStart {
				if budget.checkStateSize(len(captures), len(empty), len(loopStarts)) != nil {
					return Match{}, false
				}
				state := executionState{
					pc: pc + 1, position: pos, captures: cloneIntSlice(captures),
					emptySplits: cloneSplitVisits(empty), loopStarts: cloneLoopStarts(loopStarts),
					greedyRetry: true, loopStart: loopStart,
				}
				if budget.pushState(state) != nil {
					return Match{}, false
				}
				backtrack = append(backtrack, state)
			}
			pc++
			continue
		case opSplit:
			var guardKey splitVisit
			if ins.guard {
				guardKey = splitVisit{pc, pos}
				if _, ok := empty[guardKey]; ok {
					// If the first iteration of this loop was itself empty, retain
					// its captures on the loop-exit state. If earlier iterations
					// consumed input, the already-saved exit state has the correct
					// captures and must not be overwritten by this empty attempt.
					if loopStart, seen := loopStarts[pc]; seen && loopStart == pos {
						for i := len(backtrack) - 1; i >= 0; i-- {
							if backtrack[i].hasGuardExit && backtrack[i].guardExit == guardKey {
								updated := backtrack[i]
								updated.captures = cloneIntSlice(captures)
								if budget.replaceState(backtrack[i], updated) != nil {
									return Match{}, false
								}
								backtrack[i] = updated
								break
							}
						}
					}
					if !fail() {
						return Match{}, false
					}
					continue
				}
				newSplitCount, newLoopCount := len(empty), len(loopStarts)
				if _, exists := empty[guardKey]; !exists {
					newSplitCount++
				}
				if _, exists := loopStarts[pc]; !exists {
					newLoopCount++
				}
				newActiveBytes := executionFrameBytes(len(captures), newSplitCount, newLoopCount)
				if budget.adjustAuxiliary(activeBytes, newActiveBytes) != nil {
					return Match{}, false
				}
				activeBytes = newActiveBytes
				empty[guardKey] = struct{}{}
				if _, seen := loopStarts[pc]; !seen {
					loopStarts[pc] = pos
				}
			}
			if budget.checkStateSize(len(captures), len(empty), len(loopStarts)) != nil {
				return Match{}, false
			}
			state := executionState{pc: ins.b, position: pos, captures: cloneIntSlice(captures), emptySplits: cloneSplitVisits(empty), loopStarts: cloneLoopStarts(loopStarts)}
			if ins.guard {
				state.guardExit, state.hasGuardExit = guardKey, true
			}
			if budget.pushState(state) != nil {
				return Match{}, false
			}
			backtrack = append(backtrack, state)
			pc = ins.a
		case opJump:
			pc = ins.a
		case opSaveStart:
			if p.reverse {
				captures[ins.capture*2+1] = pos
			} else {
				captures[ins.capture*2] = pos
			}
			pc++
		case opSaveEnd:
			if p.reverse {
				captures[ins.capture*2] = pos
			} else {
				captures[ins.capture*2+1] = pos
			}
			pc++
		case opClearCaptures:
			for _, g := range ins.groups {
				captures[2*g], captures[2*g+1] = -1, -1
			}
			pc++
		case opBegin:
			if lineStartAssertion(text, pos, ins.multiline) {
				pc++
				continue
			}
			if !fail() {
				return Match{}, false
			}
		case opEnd:
			if lineEndAssertion(text, pos, ins.multiline) {
				pc++
				continue
			}
			if !fail() {
				return Match{}, false
			}
		case opBoundary:
			f := p.Flags
			f.IgnoreCase = ins.ignoreCase
			before := pos > 0 && isWordRune(previousCharacter(text, pos, p.Flags.Unicode), f)
			after := pos < len(text) && isWordRune(inputCharacter(text, pos, p.Flags.Unicode).value, f)
			if (before != after) != ins.negated {
				pc++
				continue
			}
			if !fail() {
				return Match{}, false
			}
		case opAssert:
			subcaps := cloneIntSlice(captures)
			if !ins.negative {
				for _, g := range ins.groups {
					subcaps[2*g], subcaps[2*g+1] = -1, -1
				}
			}
			leaveAssertion, err := budget.enterAssertion()
			if err != nil {
				return Match{}, false
			}
			nested, ok := ins.sub.runAtWithExecutionBudget(text, pos, subcaps, -1, budget)
			leaveAssertion()
			if budget != nil && budget.err != nil {
				return Match{}, false
			}
			if ok == ins.negative {
				if !fail() {
					return Match{}, false
				}
				continue
			}
			if !ins.negative {
				for _, g := range ins.groups {
					if g < len(nested.Captures) {
						captures[2*g], captures[2*g+1] = nested.Captures[g].Start, nested.Captures[g].End
					}
				}
			}
			pc++
		case opMatch:
			if requiredEnd >= 0 && pos != requiredEnd {
				if !fail() {
					return Match{}, false
				}
				continue
			}
			captures[1] = pos
			spans := make([]Span, p.captures+1)
			for i := range spans {
				a, b := captures[2*i], captures[2*i+1]
				if a > b && a >= 0 && b >= 0 {
					a, b = b, a
				}
				spans[i] = Span{a, b}
			}
			startMatch, endMatch := spans[0].Start, spans[0].End
			if startMatch < 0 {
				startMatch, endMatch = start, pos
			}
			return Match{Start: startMatch, End: endMatch, Captures: spans}, true
		default:
			panic("invalid regexp opcode")
		}
	}
}

func (p *Program) Find(input string, start int) (Match, bool) {
	return p.FindUTF16(utf16.Encode([]rune(input)), start)
}
func (p *Program) FindUTF16(input []uint16, start int) (Match, bool) {
	if p.linear != nil {
		return p.linear.FindUTF16(input, start)
	}
	if start < 0 {
		start = 0
	}
	if start > len(input) {
		return Match{}, false
	}
	if ShouldOptionallyStepBackToLeadSurrogate(p.Flags) && start > 0 && start < len(input) && input[start] >= 0xdc00 && input[start] <= 0xdfff && input[start-1] >= 0xd800 && input[start-1] <= 0xdbff {
		start--
	}
	if len(p.atom) != 0 {
		if at := p.findLiteralAtom(input, start); at >= 0 {
			end := at + len(p.atom)
			return Match{Start: at, End: end, Captures: []Span{{Start: at, End: end}}}, true
		}
		return Match{}, false
	}
	for at := start; at <= len(input); {
		if len(p.firstChars) != 0 && !p.firstCharacterMatches(input, at) {
			if p.Flags.Sticky {
				return Match{}, false
			}
			next := AdvanceStringIndexUTF16(input, at, p.Flags.Unicode)
			if next <= at {
				break
			}
			at = next
			continue
		}
		if m, ok := p.runAt(input, at); ok {
			return m, true
		}
		if p.Flags.Sticky {
			return Match{}, false
		}
		next := AdvanceStringIndexUTF16(input, at, p.Flags.Unicode)
		if next <= at {
			break
		}
		at = next
	}
	return Match{}, false
}
func (p *Program) FindAll(input string) []Match {
	units := utf16.Encode([]rune(input))
	out := []Match{}
	start := 0
	for start <= len(units) {
		m, ok := p.FindUTF16(units, start)
		if !ok {
			break
		}
		out = append(out, m)
		if !p.Flags.Global {
			break
		}
		if m.End == m.Start {
			start = AdvanceStringIndexUTF16(units, m.End, p.Flags.Unicode)
		} else {
			start = m.End
		}
	}
	return out
}

// FindWithLimits runs a single search with explicit work and memory ceilings.
// The older Find methods remain unbounded for source compatibility; callers
// accepting user-controlled patterns or input should use this API.
func (p *Program) FindWithLimits(input string, start int, limits ExecutionLimits) (Match, bool, error) {
	return p.FindUTF16WithLimits(utf16.Encode([]rune(input)), start, limits)
}

// FindWithLinearFallback retries a resource-limited backtracking search using
// the linear engine when the pattern belongs to its supported subset. Each
// engine attempt is bounded by limits. Patterns outside that subset preserve
// and return the original backtracking limit error. The retry is opt-in so
// FindWithLimits retains its explicit error behavior.
func (p *Program) FindWithLinearFallback(input string, start int, limits ExecutionLimits) (Match, bool, error) {
	return p.FindUTF16WithLinearFallback(utf16.Encode([]rune(input)), start, limits)
}

// FindUTF16WithLinearFallback is the UTF-16 counterpart of
// FindWithLinearFallback.
func (p *Program) FindUTF16WithLinearFallback(input []uint16, start int, limits ExecutionLimits) (Match, bool, error) {
	match, matched, err := p.FindUTF16WithLimits(input, start, limits)
	if err == nil {
		return match, matched, nil
	}
	var limitErr *ExecutionLimitError
	if !errors.As(err, &limitErr) || p.linear != nil {
		return Match{}, false, err
	}

	fallbackFlags := p.Flags
	fallbackFlags.Linear = true
	fallback, compileErr := CompileUTF16(p.source, fallbackFlags.String())
	if compileErr != nil {
		return Match{}, false, err
	}
	return fallback.FindUTF16WithLimits(input, start, limits)
}

// FindUTF16WithLimits is the bounded counterpart of FindUTF16. A resource
// limit is returned as an error, never disguised as a normal no-match result.
func (p *Program) FindUTF16WithLimits(input []uint16, start int, limits ExecutionLimits) (Match, bool, error) {
	budget, err := newExecutionBudget(limits)
	if err != nil {
		return Match{}, false, err
	}
	return p.findUTF16WithBudget(input, start, budget)
}

func (p *Program) findUTF16WithBudget(input []uint16, start int, budget *executionBudget) (Match, bool, error) {
	if start < 0 {
		start = 0
	}
	if start > len(input) {
		return Match{}, false, nil
	}
	if p.linear != nil {
		work := linearWorkUpperBound(p.linear, len(input))
		if work > budget.limits.MaxSteps-budget.steps {
			return Match{}, false, &ExecutionLimitError{Resource: "step", Maximum: budget.limits.MaxSteps}
		}
		budget.steps += work
		match, ok := p.linear.FindUTF16(input, start)
		return match, ok, nil
	}
	if ShouldOptionallyStepBackToLeadSurrogate(p.Flags) && start > 0 && start < len(input) && input[start] >= 0xdc00 && input[start] <= 0xdfff && input[start-1] >= 0xd800 && input[start-1] <= 0xdbff {
		start--
	}
	if len(p.atom) != 0 {
		at := p.findLiteralAtomWithBudget(input, start, budget)
		if budget.err != nil {
			return Match{}, false, budget.err
		}
		if at >= 0 {
			end := at + len(p.atom)
			return Match{Start: at, End: end, Captures: []Span{{Start: at, End: end}}}, true, nil
		}
		return Match{}, false, nil
	}
	for at := start; at <= len(input); {
		if len(p.firstChars) != 0 && !p.firstCharacterMatches(input, at) {
			if p.Flags.Sticky {
				return Match{}, false, nil
			}
			next := AdvanceStringIndexUTF16(input, at, p.Flags.Unicode)
			if next <= at {
				break
			}
			at = next
			continue
		}
		match, ok := p.runAtWithExecutionBudget(input, at, nil, -1, budget)
		if budget.err != nil {
			return Match{}, false, budget.err
		}
		if ok {
			return match, true, nil
		}
		if p.Flags.Sticky {
			return Match{}, false, nil
		}
		next := AdvanceStringIndexUTF16(input, at, p.Flags.Unicode)
		if next <= at {
			break
		}
		at = next
	}
	return Match{}, false, nil
}

// FindAllWithLimits shares one work budget across the complete global search,
// so repeated empty matches cannot reset the caller's execution ceiling.
func (p *Program) FindAllWithLimits(input string, limits ExecutionLimits) ([]Match, error) {
	budget, err := newExecutionBudget(limits)
	if err != nil {
		return nil, err
	}
	units := utf16.Encode([]rune(input))
	matches := make([]Match, 0)
	start := 0
	for start <= len(units) {
		match, ok, err := p.findUTF16WithBudget(units, start, budget)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		matches = append(matches, match)
		if !p.Flags.Global {
			break
		}
		if match.End == match.Start {
			start = AdvanceStringIndexUTF16(units, match.End, p.Flags.Unicode)
		} else {
			start = match.End
		}
	}
	return matches, nil
}

func linearWorkUpperBound(program *linearProgram, inputLength int) uint64 {
	var instructionCount uint64
	var count func(*linearProgram)
	count = func(current *linearProgram) {
		instructionCount += uint64(len(current.code))
		for _, lookbehind := range current.lookbehinds {
			count(lookbehind)
		}
	}
	count(program)
	positions := uint64(inputLength) + 1
	if instructionCount != 0 && positions > ^uint64(0)/instructionCount {
		return ^uint64(0)
	}
	return positions * instructionCount
}

func isASCIIWordRune(c rune) bool {
	return c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c == '_'
}
func isWordRune(c rune, flags Flags) bool {
	if isASCIIWordRune(c) {
		return true
	}
	if flags.IgnoreCase && flags.Unicode {
		for folded := unicode.SimpleFold(c); folded != c; folded = unicode.SimpleFold(folded) {
			if isASCIIWordRune(folded) {
				return true
			}
		}
	}
	return false
}
func isLineTerminator(c rune) bool { return c == '\n' || c == '\r' || c == 0x2028 || c == 0x2029 }
func lineStartAssertion(text []uint16, pos int, multiline bool) bool {
	if pos == 0 {
		return true
	}
	return multiline && isLineTerminator(rune(text[pos-1]))
}
func lineEndAssertion(text []uint16, pos int, multiline bool) bool {
	if pos == len(text) {
		return true
	}
	return multiline && isLineTerminator(rune(text[pos]))
}
