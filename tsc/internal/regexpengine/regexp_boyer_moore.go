package regexpengine

import "math/bits"

type v8NodesContainedInLattice uint8

const (
	v8NodesLatticeNotYet v8NodesContainedInLattice = iota
	v8NodesLatticeIn
	v8NodesLatticeOut
	v8NodesLatticeUnknown
)

func v8NodesCombineLattice(left, right v8NodesContainedInLattice) v8NodesContainedInLattice {
	return left | right
}

type v8NodesBoyerMoorePositionInfo struct {
	mapBits  [2]uint64
	mapCount int
	word     v8NodesContainedInLattice
}

const (
	v8NodesBoyerMooreMapSize = 128
	v8NodesBoyerMooreMapMask = v8NodesBoyerMooreMapSize - 1
)

func (position *v8NodesBoyerMoorePositionInfo) At(index int) bool {
	if index < 0 || index >= v8NodesBoyerMooreMapSize {
		panic("V8 Boyer-Moore character index is out of range")
	}
	return position.mapBits[index>>6]&(uint64(1)<<uint(index&63)) != 0
}

func (position *v8NodesBoyerMoorePositionInfo) MapCount() int { return position.mapCount }
func (position *v8NodesBoyerMoorePositionInfo) IsNonWord() bool {
	return position.word == v8NodesLatticeOut
}
func (position *v8NodesBoyerMoorePositionInfo) IsWord() bool {
	return position.word == v8NodesLatticeIn
}
func (position *v8NodesBoyerMoorePositionInfo) RawBitset() [2]uint64 {
	return position.mapBits
}

func (position *v8NodesBoyerMoorePositionInfo) Set(character int) {
	position.SetInterval(v8ASTNewInterval(character, character))
}

func (position *v8NodesBoyerMoorePositionInfo) SetInterval(interval v8NodesInterval) {
	position.word = v8NodesBoyerMooreAddRange(position.word, interval)
	if interval.size() >= v8NodesBoyerMooreMapSize {
		position.mapCount = v8NodesBoyerMooreMapSize
		position.mapBits = [2]uint64{^uint64(0), ^uint64(0)}
		return
	}
	for character := interval.from; character <= interval.to; character++ {
		modCharacter := character & v8NodesBoyerMooreMapMask
		wordIndex := modCharacter >> 6
		bit := uint64(1) << uint(modCharacter&63)
		if position.mapBits[wordIndex]&bit == 0 {
			position.mapCount++
			position.mapBits[wordIndex] |= bit
		}
		if position.mapCount == v8NodesBoyerMooreMapSize {
			return
		}
	}
}

func (position *v8NodesBoyerMoorePositionInfo) SetAll() {
	position.word = v8NodesLatticeUnknown
	if position.mapCount != v8NodesBoyerMooreMapSize {
		position.mapCount = v8NodesBoyerMooreMapSize
		position.mapBits = [2]uint64{^uint64(0), ^uint64(0)}
	}
}

func v8NodesBoyerMooreAddRange(containment v8NodesContainedInLattice, interval v8NodesInterval) v8NodesContainedInLattice {
	// regexp-compiler.h's half-open word ranges, including the terminal marker.
	ranges := v8CompilerWordRanges[:]
	if containment == v8NodesLatticeUnknown {
		return containment
	}
	inside := false
	last := 0
	for _, edge := range ranges {
		if edge <= interval.from {
			inside = !inside
			last = edge
			continue
		}
		if last <= interval.from && interval.to < edge {
			membership := v8NodesLatticeOut
			if inside {
				membership = v8NodesLatticeIn
			}
			return v8NodesCombineLattice(containment, membership)
		}
		return v8NodesLatticeUnknown
	}
	return containment
}

func v8NodesBitsetFirstSetBit(set [2]uint64) int {
	if set[0] != 0 {
		return bits.TrailingZeros64(set[0])
	}
	if set[1] != 0 {
		return 64 + bits.TrailingZeros64(set[1])
	}
	return -1
}

type v8NodesBoyerMooreLookahead struct {
	length         int
	compiler       *v8CompilerCompiler
	maxChar        int
	bitmaps        []*v8NodesBoyerMoorePositionInfo
	cachesNodeInfo bool
}

func v8NodesNewBoyerMooreLookahead(length int, compiler *v8CompilerCompiler) *v8NodesBoyerMooreLookahead {
	if length < 0 || compiler == nil {
		panic("V8 Boyer-Moore lookahead requires a nonnegative length and compiler")
	}
	maxChar := 0xffff
	if compiler.oneByte {
		maxChar = 0xff
	}
	lookahead := &v8NodesBoyerMooreLookahead{
		length:         length,
		compiler:       compiler,
		maxChar:        maxChar,
		bitmaps:        make([]*v8NodesBoyerMoorePositionInfo, length),
		cachesNodeInfo: true,
	}
	for index := range lookahead.bitmaps {
		lookahead.bitmaps[index] = &v8NodesBoyerMoorePositionInfo{}
	}
	return lookahead
}

func (lookahead *v8NodesBoyerMooreLookahead) Length() int  { return lookahead.length }
func (lookahead *v8NodesBoyerMooreLookahead) MaxChar() int { return lookahead.maxChar }
func (lookahead *v8NodesBoyerMooreLookahead) Compiler() *v8CompilerCompiler {
	return lookahead.compiler
}
func (lookahead *v8NodesBoyerMooreLookahead) Count(mapNumber int) int {
	return lookahead.At(mapNumber).MapCount()
}
func (lookahead *v8NodesBoyerMooreLookahead) At(index int) *v8NodesBoyerMoorePositionInfo {
	if index < 0 || index >= len(lookahead.bitmaps) {
		panic("V8 Boyer-Moore lookahead index is out of range")
	}
	return lookahead.bitmaps[index]
}
func (lookahead *v8NodesBoyerMooreLookahead) Set(mapNumber, character int) {
	if character > lookahead.maxChar {
		return
	}
	lookahead.At(mapNumber).Set(character)
}
func (lookahead *v8NodesBoyerMooreLookahead) SetInterval(mapNumber int, interval v8NodesInterval) {
	if interval.from > lookahead.maxChar {
		return
	}
	if interval.to > lookahead.maxChar {
		interval.to = lookahead.maxChar
	}
	lookahead.At(mapNumber).SetInterval(interval)
}
func (lookahead *v8NodesBoyerMooreLookahead) SetAll(mapNumber int) {
	lookahead.At(mapNumber).SetAll()
}
func (lookahead *v8NodesBoyerMooreLookahead) SetRest(fromMap int) {
	if fromMap < 0 || fromMap > lookahead.length {
		panic("V8 Boyer-Moore SetRest index is out of range")
	}
	for index := fromMap; index < lookahead.length; index++ {
		lookahead.SetAll(index)
	}
}
func (lookahead *v8NodesBoyerMooreLookahead) CachesNodeInfo() bool {
	return lookahead.cachesNodeInfo
}
func (lookahead *v8NodesBoyerMooreLookahead) SetCachesNodeInfo(value bool) {
	lookahead.cachesNodeInfo = value
}

func (lookahead *v8NodesBoyerMooreLookahead) findWorthwhileInterval() (from, to int, found bool) {
	biggestPoints := 0
	for maxCharacters := 4; maxCharacters < 32; maxCharacters *= 2 {
		biggestPoints = lookahead.findBestInterval(maxCharacters, biggestPoints, &from, &to)
	}
	return from, to, biggestPoints != 0
}

func (lookahead *v8NodesBoyerMooreLookahead) findBestInterval(maxCharacters, oldBiggestPoints int, from, to *int) int {
	biggestPoints := oldBiggestPoints
	for index := 0; index < lookahead.length; {
		for index < lookahead.length && lookahead.Count(index) > maxCharacters {
			index++
		}
		if index == lookahead.length {
			break
		}
		rememberedFrom := index
		var union [2]uint64
		for index < lookahead.length && lookahead.Count(index) <= maxCharacters {
			bitmap := lookahead.At(index).RawBitset()
			union[0] |= bitmap[0]
			union[1] |= bitmap[1]
			index++
		}

		frequency := 0
		for character := v8NodesBitsetFirstSetBit(union); character != -1; character = v8NodesBitsetFirstSetBit(union) {
			frequency += lookahead.compiler.FrequencyCollator().Frequency(character) + 1
			union[character>>6] &^= uint64(1) << uint(character&63)
		}

		inQuickCheckRange := index-rememberedFrom < 4 ||
			(lookahead.compiler.oneByte && rememberedFrom <= 4) ||
			(!lookahead.compiler.oneByte && rememberedFrom <= 2)
		probability := v8NodesBoyerMooreMapSize - frequency
		if inQuickCheckRange {
			probability = v8NodesBoyerMooreMapSize/2 - frequency
		}
		points := (index - rememberedFrom) * probability
		if points > biggestPoints {
			*from = rememberedFrom
			*to = index - 1
			biggestPoints = points
		}
	}
	return biggestPoints
}

func (lookahead *v8NodesBoyerMooreLookahead) GetSkipTable(minLookahead, maxLookahead int, withNibbleTable bool) (booleanTable, nibbleTable []byte, skip int) {
	if minLookahead < 0 || maxLookahead < minLookahead || maxLookahead >= lookahead.length {
		panic("V8 Boyer-Moore skip-table range is out of bounds")
	}
	booleanTable = make([]byte, v8NodesBoyerMooreMapSize)
	if withNibbleTable {
		nibbleTable = make([]byte, v8NodesBoyerMooreMapSize/8)
	}
	for index := maxLookahead; index >= minLookahead; index-- {
		set := lookahead.At(index).RawBitset()
		for character := v8NodesBitsetFirstSetBit(set); character != -1; character = v8NodesBitsetFirstSetBit(set) {
			booleanTable[character] = 1
			if withNibbleTable {
				lowNibble := character & 0x0f
				highNibble := (character >> 4) & 0x07
				nibbleTable[lowNibble] |= byte(1 << uint(highNibble))
			}
			set[character>>6] &^= uint64(1) << uint(character&63)
		}
	}
	return booleanTable, nibbleTable, maxLookahead + 1 - minLookahead
}

// BuildSkipTable exposes the scan inputs independently of V8's byte-array
// allocator and macro assembler. The Go emitter owns the actual scan code.
func (lookahead *v8NodesBoyerMooreLookahead) BuildSkipTable() (offset, advanceBy int, booleanTable, nibbleTable []byte, found bool) {
	minLookahead, maxLookahead, worthwhile := lookahead.findWorthwhileInterval()
	if !worthwhile {
		return 0, 0, nil, nil, false
	}
	booleanTable, nibbleTable, advanceBy = lookahead.GetSkipTable(minLookahead, maxLookahead, true)
	return maxLookahead, advanceBy, booleanTable, nibbleTable, true
}

type v8NodesCharacterFrequency struct {
	counter   int
	character int
}

func (frequency *v8NodesCharacterFrequency) Increment()     { frequency.counter++ }
func (frequency *v8NodesCharacterFrequency) Counter() int   { return frequency.counter }
func (frequency *v8NodesCharacterFrequency) Character() int { return frequency.character }

type v8CompilerFrequencyCollator struct {
	frequencies  [v8NodesBoyerMooreMapSize]v8NodesCharacterFrequency
	totalSamples int
}

func v8CompilerNewFrequencyCollator() v8CompilerFrequencyCollator {
	var collator v8CompilerFrequencyCollator
	for index := range collator.frequencies {
		collator.frequencies[index].character = index
	}
	return collator
}

func (collator *v8CompilerFrequencyCollator) CountCharacter(character int) {
	index := character & v8NodesBoyerMooreMapMask
	collator.frequencies[index].Increment()
	collator.totalSamples++
}

func (collator *v8CompilerFrequencyCollator) Frequency(character int) int {
	if character&v8NodesBoyerMooreMapMask != character {
		panic("V8 frequency character is outside the 128-entry table")
	}
	if collator.totalSamples < 1 {
		return 1
	}
	return (collator.frequencies[character].Counter() * v8NodesBoyerMooreMapSize) / collator.totalSamples
}
