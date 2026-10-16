package goemit

// SliceRuntime represents an index or a Python-style range with a small,
// tagged Go value. User code receives these values in one ordinary TS array.
const SliceRuntime = `
type tsSlicePartKind uint8
const (
	tsSliceIndexPart tsSlicePartKind = iota
	tsSliceRangePart
)
type tsSlice struct {
	kind tsSlicePartKind
	index float64
	start, stop, step float64
	hasStart, hasStop, hasStep bool
}
func tsIndexSliceValue(value tsValue) tsValue {
	return tsSliceValue(&tsSlice{kind: tsSliceIndexPart, index: tsNumber(value)})
}
func tsRangeSliceValue(start, stop, step tsValue) tsValue {
	value := &tsSlice{kind: tsSliceRangePart}
	if !tsIsUndefined(start) { value.start = tsNumber(start); value.hasStart = true }
	if !tsIsUndefined(stop) { value.stop = tsNumber(stop); value.hasStop = true }
	if !tsIsUndefined(step) { value.step = tsNumber(step); value.hasStep = true }
	return tsSliceValue(value)
}
func tsSliceGet(value, key tsValue) tsValue {
	part := (*tsSlice)(value.ref)
	switch tsPropertyKey(key) {
	case "kind":
		if part.kind == tsSliceIndexPart { return tsStringReference(tsStringUTF8("index")) }
		return tsStringReference(tsStringUTF8("range"))
	case "index":
		if part.kind == tsSliceIndexPart { return tsNumberValue(part.index) }
	case "start":
		if part.kind == tsSliceRangePart && part.hasStart { return tsNumberValue(part.start) }
	case "stop":
		if part.kind == tsSliceRangePart && part.hasStop { return tsNumberValue(part.stop) }
	case "step":
		if part.kind == tsSliceRangePart && part.hasStep { return tsNumberValue(part.step) }
	}
	return tsU
}
`
