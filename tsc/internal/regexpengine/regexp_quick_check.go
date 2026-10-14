package regexpengine

// QuickCheckDetails is V8's fixed-width summary of the next input code units
// that may match. A quick check may admit false positives, but must not reject
// an input that the full graph can match.
const v8NodesQuickCheckMaxPositions = 4

type v8NodesQuickCheckPosition struct {
	mask                uint32
	value               uint32
	determinesPerfectly bool
	cannotMatch         bool
}

func (position *v8NodesQuickCheckPosition) Clear() {
	position.mask = 0
	position.value = 0
	position.determinesPerfectly = false
	position.cannotMatch = false
}

type v8NodesQuickCheckDetails struct {
	characters int
	positions  [v8NodesQuickCheckMaxPositions]v8NodesQuickCheckPosition
	mask       uint32
	value      uint32
}

type v8NodesQuickCheckDetailsAPI interface {
	Rationalize(oneByte bool) bool
	Merge(other *v8NodesQuickCheckDetails, fromIndex int)
	Advance(by int, oneByte bool)
	Clear()
	CannotMatch() bool
	SetCannotMatchFrom(index int)
	Characters() int
	SetCharacters(characters int)
	Position(index int) *v8NodesQuickCheckPosition
	Mask() uint32
	Value() uint32
	SetMask(mask uint32)
	SetValue(value uint32)
}

func v8NodesNewQuickCheckDetails() *v8NodesQuickCheckDetails {
	return &v8NodesQuickCheckDetails{}
}

func v8NodesNewQuickCheckDetailsWithCharacters(characters int) *v8NodesQuickCheckDetails {
	details := v8NodesNewQuickCheckDetails()
	details.SetCharacters(characters)
	return details
}

func (details *v8NodesQuickCheckDetails) Rationalize(oneByte bool) bool {
	charMask := v8CompilerCharMask(oneByte)
	charShift := uint(16)
	if oneByte {
		charShift = 8
	}
	details.mask = 0
	details.value = 0
	foundUsefulOperation := false
	shift := uint(0)
	for index := 0; index < details.characters; index++ {
		position := &details.positions[index]
		if position.mask&0xff != 0 {
			foundUsefulOperation = true
		}
		details.mask |= (position.mask & charMask) << shift
		details.value |= (position.value & charMask) << shift
		shift += charShift
	}
	return foundUsefulOperation
}

func (details *v8NodesQuickCheckDetails) Merge(other *v8NodesQuickCheckDetails, fromIndex int) {
	if other == nil || details.characters != other.characters {
		panic("V8 QuickCheckDetails.Merge requires equal character counts")
	}
	if fromIndex < 0 || fromIndex > details.characters {
		panic("V8 QuickCheckDetails.Merge index is out of range")
	}
	for index := fromIndex; index < details.characters; index++ {
		position := &details.positions[index]
		otherPosition := &other.positions[index]
		if position.cannotMatch {
			*position = *otherPosition
		} else if !otherPosition.cannotMatch {
			if position.mask != otherPosition.mask || position.value != otherPosition.value || !otherPosition.determinesPerfectly {
				position.determinesPerfectly = false
			}
			position.mask &= otherPosition.mask
			position.value &= position.mask
			otherPosition.value &= position.mask
			differingBits := position.value ^ otherPosition.value
			position.mask &= ^differingBits
			position.value &= position.mask
		}
	}
}

func (details *v8NodesQuickCheckDetails) Advance(by int, _ bool) {
	if by >= details.characters || by < 0 {
		details.Clear()
		return
	}
	remaining := details.characters - by
	for index := 0; index < remaining; index++ {
		details.positions[index] = details.positions[by+index]
	}
	for index := remaining; index < details.characters; index++ {
		details.positions[index].Clear()
	}
	details.characters = remaining
	// V8 intentionally leaves mask/value unchanged: after an Advance they have
	// already served their emitted check and will not be consumed a second time.
}

func (details *v8NodesQuickCheckDetails) Clear() {
	for index := 0; index < details.characters; index++ {
		details.positions[index].Clear()
	}
	details.characters = 0
}

func (details *v8NodesQuickCheckDetails) CannotMatch() bool {
	for index := 0; index < details.characters; index++ {
		if details.positions[index].cannotMatch {
			return true
		}
	}
	return false
}

func (details *v8NodesQuickCheckDetails) SetCannotMatchFrom(index int) {
	if index < 0 {
		panic("V8 QuickCheckDetails index cannot be negative")
	}
	for ; index < details.characters; index++ {
		details.positions[index].cannotMatch = true
	}
}

func (details *v8NodesQuickCheckDetails) Characters() int { return details.characters }

func (details *v8NodesQuickCheckDetails) SetCharacters(characters int) {
	if characters < 0 || characters > v8NodesQuickCheckMaxPositions {
		panic("V8 QuickCheckDetails character count is out of range")
	}
	details.characters = characters
}

func (details *v8NodesQuickCheckDetails) Position(index int) *v8NodesQuickCheckPosition {
	if index < 0 || index >= details.characters {
		panic("V8 QuickCheckDetails position is out of range")
	}
	return &details.positions[index]
}

func (details *v8NodesQuickCheckDetails) Mask() uint32          { return details.mask }
func (details *v8NodesQuickCheckDetails) Value() uint32         { return details.value }
func (details *v8NodesQuickCheckDetails) SetMask(mask uint32)   { details.mask = mask }
func (details *v8NodesQuickCheckDetails) SetValue(value uint32) { details.value = value }
