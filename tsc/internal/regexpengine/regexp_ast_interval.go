// Copyright 2016 the V8 project authors. All rights reserved.
// Inclusive register intervals shared by the AST and graph compiler port.
package regexpengine

type registerInterval struct {
	from int
	to   int
}

const (
	intervalNone    = -1
	intervalInvalid = -2
)

func emptyRegisterInterval() registerInterval {
	return registerInterval{from: intervalNone, to: intervalNone - 1}
}

func invalidRegisterInterval() registerInterval {
	return registerInterval{from: intervalInvalid, to: intervalInvalid}
}

func (i registerInterval) valid() bool { return i.from != intervalInvalid }

func (i registerInterval) empty() bool { return i.valid() && i.from == intervalNone }

func (i registerInterval) contains(register int) bool {
	return i.valid() && !i.empty() && i.from <= register && register <= i.to
}

func (i registerInterval) size() int {
	if !i.valid() || i.empty() {
		return 0
	}
	return i.to - i.from + 1
}

func (i registerInterval) union(other registerInterval) registerInterval {
	if !i.valid() {
		return i
	}
	if !other.valid() {
		return other
	}
	if other.empty() {
		return i
	}
	if i.empty() {
		return other
	}
	if other.from < i.from {
		i.from = other.from
	}
	if other.to > i.to {
		i.to = other.to
	}
	return i
}
