package goemit

// These candidates are experimental benchmarks, not the production runtime.
// Each variant preserves distinct null/undefined and GC-visible references.

import (
	"math"
	"runtime"
	"testing"
	"unsafe"
)

type benchKind uint8

const (
	benchUndefinedKind benchKind = iota
	benchNullKind
	benchNumberKind
	benchBooleanKind
	benchStringKind
	benchObjectKind
)

type benchUndefined struct{}
type benchText struct{ units []uint16 }
type benchObject struct{ number float64 }

// 24 bytes on amd64: inline tag, full numeric payload, GC-visible pointer.
// The pointer never passes through uintptr or integer payload bits.
// A safe 32-byte alternative uses a Go interface only for references; numeric
// payloads remain unboxed and its tag controls JavaScript semantics.
type benchValue32 struct {
	number    float64
	reference any
	kind      benchKind
}

type benchValue24 struct {
	number    float64
	reference unsafe.Pointer
	kind      benchKind
}

// 16 bytes on amd64: nil reference means Number; other kinds use a header.
// Headers for null, undefined, and boolean are shared immutable sentinels.
// Object/string headers cost memory and an indirection, so measure both.
type benchHeader struct {
	reference unsafe.Pointer
	kind      benchKind
}
type benchValue16 struct {
	number    float64
	reference *benchHeader
}
type benchHeaderObject struct {
	header benchHeader
	number float64
}
type benchHeaderText struct {
	header benchHeader
	units  []uint16
}

var benchUndefinedHeader = benchHeader{kind: benchUndefinedKind}
var benchNullHeader = benchHeader{kind: benchNullKind}
var benchBooleanHeader = benchHeader{kind: benchBooleanKind}

var benchNumberSink float64
var benchBoolSink bool
var benchAnySink any
var bench24Sink benchValue24
var bench16Sink benchValue16
var bench32Sink benchValue32
var bench32SliceSink []benchValue32
var benchAnySliceSink []any
var bench24SliceSink []benchValue24
var bench16SliceSink []benchValue16
var benchNativeSliceSink []float64

type benchValues struct {
	dynamic []any
	inline  []benchValue24
	header  []benchValue16
	native  []float64
}

func makeBenchValues(mixed bool) benchValues {
	values := benchValues{make([]any, 4096), make([]benchValue24, 4096), make([]benchValue16, 4096), make([]float64, 4096)}
	for i := range values.dynamic {
		number := float64(i) + 0.125
		values.native[i] = number
		kind := benchNumberKind
		if mixed {
			kind = benchKind(i % 6)
		}
		value24 := benchValue24{kind: kind}
		value16 := benchValue16{}
		switch kind {
		case benchNumberKind:
			values.dynamic[i] = number
			value24.number = number
			value16.number = number
		case benchUndefinedKind:
			values.dynamic[i] = benchUndefined{}
			value16.reference = &benchUndefinedHeader
		case benchNullKind:
			values.dynamic[i] = nil
			value16.reference = &benchNullHeader
		case benchBooleanKind:
			values.dynamic[i] = true
			value24.number = 1
			value16.number = 1
			value16.reference = &benchBooleanHeader
		case benchStringKind:
			text := &benchText{units: []uint16{'a', 'b', 'c'}}
			values.dynamic[i] = text
			value24.reference = unsafe.Pointer(text)
			header := &benchHeaderText{units: text.units}
			header.header = benchHeader{unsafe.Pointer(header), kind}
			value16.reference = &header.header
		case benchObjectKind:
			object := &benchObject{number}
			values.dynamic[i] = object
			value24.reference = unsafe.Pointer(object)
			header := &benchHeaderObject{number: number}
			header.header = benchHeader{unsafe.Pointer(header), kind}
			value16.reference = &header.header
		}
		values.inline[i] = value24
		values.header[i] = value16
	}
	return values
}

func BenchmarkValueTypeCheck(b *testing.B) {
	values := makeBenchValues(true)
	b.Run("Any", func(b *testing.B) {
		i := 0
		result := false
		for b.Loop() {
			_, result = values.dynamic[i&4095].(float64)
			i++
		}
		benchBoolSink = result
	})
	b.Run("Inline24", func(b *testing.B) {
		i := 0
		result := false
		for b.Loop() {
			result = values.inline[i&4095].kind == benchNumberKind
			i++
		}
		benchBoolSink = result
	})
	b.Run("Header16", func(b *testing.B) {
		i := 0
		result := false
		for b.Loop() {
			result = values.header[i&4095].reference == nil
			i++
		}
		benchBoolSink = result
	})
}
func BenchmarkValueNumericRead(b *testing.B) {
	values := makeBenchValues(false)
	b.Run("Native", func(b *testing.B) {
		i := 0
		sum := 0.0
		for b.Loop() {
			sum += values.native[i&4095]
			i++
		}
		benchNumberSink = sum
	})
	b.Run("Any", func(b *testing.B) {
		i := 0
		sum := 0.0
		for b.Loop() {
			sum += values.dynamic[i&4095].(float64)
			i++
		}
		benchNumberSink = sum
	})
	b.Run("Inline24", func(b *testing.B) {
		i := 0
		sum := 0.0
		for b.Loop() {
			v := values.inline[i&4095]
			if v.kind != benchNumberKind {
				b.Fatal("wrong type")
			}
			sum += v.number
			i++
		}
		benchNumberSink = sum
	})
	b.Run("Header16", func(b *testing.B) {
		i := 0
		sum := 0.0
		for b.Loop() {
			v := values.header[i&4095]
			if v.reference != nil {
				b.Fatal("wrong type")
			}
			sum += v.number
			i++
		}
		benchNumberSink = sum
	})
}
func BenchmarkValueNumericUpdate(b *testing.B) {
	b.Run("Native", func(b *testing.B) {
		cell := new(float64)
		*cell = 1024.125
		benchAnySink = cell
		b.ResetTimer()
		for b.Loop() {
			*cell += 1.125
		}
		benchNumberSink = *cell
	})
	b.Run("Any", func(b *testing.B) {
		cell := new(any)
		*cell = float64(1024.125)
		benchAnySink = cell
		b.ResetTimer()
		for b.Loop() {
			*cell = (*cell).(float64) + 1.125
		}
		benchNumberSink = (*cell).(float64)
	})
	b.Run("Inline24", func(b *testing.B) {
		cell := &benchValue24{number: 1024.125, kind: benchNumberKind}
		benchAnySink = cell
		b.ResetTimer()
		for b.Loop() {
			if cell.kind != benchNumberKind {
				b.Fatal("wrong type")
			}
			cell.number += 1.125
		}
		benchNumberSink = cell.number
	})
	b.Run("Header16", func(b *testing.B) {
		cell := &benchValue16{number: 1024.125}
		benchAnySink = cell
		b.ResetTimer()
		for b.Loop() {
			if cell.reference != nil {
				b.Fatal("wrong type")
			}
			cell.number += 1.125
		}
		benchNumberSink = cell.number
	})
}
func benchAnyContribution(value any) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case bool:
		if v {
			return 1
		}
	case *benchText:
		return float64(len(v.units))
	case *benchObject:
		return v.number
	}
	return 0
}
func bench24Contribution(value benchValue24) float64 {
	switch value.kind {
	case benchNumberKind, benchBooleanKind:
		return value.number
	case benchStringKind:
		return float64(len((*benchText)(value.reference).units))
	case benchObjectKind:
		return (*benchObject)(value.reference).number
	}
	return 0
}
func bench16Contribution(value benchValue16) float64 {
	if value.reference == nil {
		return value.number
	}
	switch value.reference.kind {
	case benchBooleanKind:
		return value.number
	case benchStringKind:
		return float64(len((*benchHeaderText)(value.reference.reference).units))
	case benchObjectKind:
		return (*benchHeaderObject)(value.reference.reference).number
	}
	return 0
}
func BenchmarkValueMixedDispatch(b *testing.B) {
	values := makeBenchValues(true)
	b.Run("Any", func(b *testing.B) {
		i := 0
		sum := 0.0
		for b.Loop() {
			sum += benchAnyContribution(values.dynamic[i&4095])
			i++
		}
		benchNumberSink = sum
	})
	b.Run("Inline24", func(b *testing.B) {
		i := 0
		sum := 0.0
		for b.Loop() {
			sum += bench24Contribution(values.inline[i&4095])
			i++
		}
		benchNumberSink = sum
	})
	b.Run("Header16", func(b *testing.B) {
		i := 0
		sum := 0.0
		for b.Loop() {
			sum += bench16Contribution(values.header[i&4095])
			i++
		}
		benchNumberSink = sum
	})
}

// Prevent folding the calling convention into its caller. Each variant pays the
// same one call, with representation-dependent argument/return costs.
//
//go:noinline
func benchNativeAdd(a, b float64) float64 { return a + b }

//go:noinline
func benchAnyAdd(a, b any) any { return a.(float64) + b.(float64) }

//go:noinline
func bench24Add(a, b benchValue24) benchValue24 {
	if a.kind != benchNumberKind || b.kind != benchNumberKind {
		panic("wrong type")
	}
	return benchValue24{number: a.number + b.number, kind: benchNumberKind}
}

//go:noinline
func bench16Add(a, b benchValue16) benchValue16 {
	if a.reference != nil || b.reference != nil {
		panic("wrong type")
	}
	return benchValue16{number: a.number + b.number}
}
func BenchmarkValueCall(b *testing.B) {
	b.Run("Native", func(b *testing.B) {
		result := 0.0
		for b.Loop() {
			result = benchNativeAdd(1024.125, 2048.25)
		}
		benchNumberSink = result
	})
	b.Run("Any", func(b *testing.B) {
		var a, bValue any = float64(1024.125), float64(2048.25)
		var result any
		for b.Loop() {
			result = benchAnyAdd(a, bValue)
		}
		benchAnySink = result
	})
	b.Run("Inline24", func(b *testing.B) {
		a := benchValue24{number: 1024.125, kind: benchNumberKind}
		c := benchValue24{number: 2048.25, kind: benchNumberKind}
		result := benchValue24{}
		for b.Loop() {
			result = bench24Add(a, c)
		}
		bench24Sink = result
	})
	b.Run("Header16", func(b *testing.B) {
		a := benchValue16{number: 1024.125}
		c := benchValue16{number: 2048.25}
		result := benchValue16{}
		for b.Loop() {
			result = bench16Add(a, c)
		}
		bench16Sink = result
	})
}
func BenchmarkValueObjectRead(b *testing.B) {
	object := &benchObject{1024.125}
	header := &benchHeaderObject{number: 1024.125}
	header.header = benchHeader{unsafe.Pointer(header), benchObjectKind}
	b.Run("Any", func(b *testing.B) {
		value := benchIdentityAny(object)
		sum := 0.0
		for b.Loop() {
			v, ok := value.(*benchObject)
			if !ok {
				b.Fatal("wrong type")
			}
			sum += v.number
		}
		benchNumberSink = sum
	})
	b.Run("Inline24", func(b *testing.B) {
		value := benchIdentity24(benchValue24{kind: benchObjectKind, reference: unsafe.Pointer(object)})
		sum := 0.0
		for b.Loop() {
			if value.kind != benchObjectKind {
				b.Fatal("wrong type")
			}
			sum += (*benchObject)(value.reference).number
		}
		benchNumberSink = sum
	})
	b.Run("Header16", func(b *testing.B) {
		value := benchIdentity16(benchValue16{reference: &header.header})
		sum := 0.0
		for b.Loop() {
			if value.reference == nil || value.reference.kind != benchObjectKind {
				b.Fatal("wrong type")
			}
			sum += (*benchHeaderObject)(value.reference.reference).number
		}
		benchNumberSink = sum
	})
}
func BenchmarkValueBuildNumbers(b *testing.B) {
	const count = 1024
	b.Run("Native", func(b *testing.B) {
		for b.Loop() {
			values := make([]float64, count)
			for i := range values {
				values[i] = float64(i) + 0.125
			}
			benchNativeSliceSink = values
		}
	})
	b.Run("Any", func(b *testing.B) {
		for b.Loop() {
			values := make([]any, count)
			for i := range values {
				values[i] = float64(i) + 0.125
			}
			benchAnySliceSink = values
		}
	})
	b.Run("Inline24", func(b *testing.B) {
		for b.Loop() {
			values := make([]benchValue24, count)
			for i := range values {
				values[i] = benchValue24{number: float64(i) + 0.125, kind: benchNumberKind}
			}
			bench24SliceSink = values
		}
	})
	b.Run("Header16", func(b *testing.B) {
		for b.Loop() {
			values := make([]benchValue16, count)
			for i := range values {
				values[i].number = float64(i) + 0.125
			}
			bench16SliceSink = values
		}
	})
}
func BenchmarkValueBuildObjects(b *testing.B) {
	const count = 256
	b.Run("Any", func(b *testing.B) {
		for b.Loop() {
			values := make([]any, count)
			for i := range values {
				values[i] = &benchObject{float64(i) + 0.125}
			}
			benchAnySliceSink = values
		}
	})
	b.Run("Inline24", func(b *testing.B) {
		for b.Loop() {
			values := make([]benchValue24, count)
			for i := range values {
				values[i] = benchValue24{kind: benchObjectKind, reference: unsafe.Pointer(&benchObject{float64(i) + 0.125})}
			}
			bench24SliceSink = values
		}
	})
	b.Run("Header16", func(b *testing.B) {
		for b.Loop() {
			values := make([]benchValue16, count)
			for i := range values {
				object := &benchHeaderObject{number: float64(i) + 0.125}
				object.header = benchHeader{unsafe.Pointer(object), benchObjectKind}
				values[i].reference = &object.header
			}
			bench16SliceSink = values
		}
	})
}
func TestValueCandidateSemantics(t *testing.T) {
	t.Logf("sizes: any=%d, inline=%d, header=%d, header metadata=%d", unsafe.Sizeof(any(nil)), unsafe.Sizeof(benchValue24{}), unsafe.Sizeof(benchValue16{}), unsafe.Sizeof(benchHeader{}))
	values := makeBenchValues(true)
	for i := range values.dynamic {
		expected := benchAnyContribution(values.dynamic[i])
		if bench24Contribution(values.inline[i]) != expected || bench16Contribution(values.header[i]) != expected {
			t.Fatalf("mismatch at %d", i)
		}
	}
	if benchUndefinedHeader.kind == benchNullHeader.kind {
		t.Fatal("null and undefined merged")
	}
	// Preserve all float64 bit patterns, including NaN and negative zero.
	for _, number := range []float64{math.NaN(), math.Inf(1), math.Copysign(0, -1), math.SmallestNonzeroFloat64} {
		inline := benchValue24{number: number, kind: benchNumberKind}
		header := benchValue16{number: number}
		if math.Float64bits(inline.number) != math.Float64bits(number) || math.Float64bits(header.number) != math.Float64bits(number) {
			t.Fatal("numeric bits changed")
		}
	}
}
func TestValueCandidateGCReferences(t *testing.T) {
	// Candidate references are the sole roots after the factory returns. Both
	// pointer fields are visible to Go's collector; no uintptr is involved.
	values := makeBenchValues(true)
	values.dynamic = nil
	for range 3 {
		runtime.GC()
	}
	for i := range values.inline {
		a, c := bench24Contribution(values.inline[i]), bench16Contribution(values.header[i])
		if a != c {
			t.Fatalf("reference mismatch at %d", i)
		}
	}
	runtime.KeepAlive(values)
}

//go:noinline
func bench32Add(a, b benchValue32) benchValue32 {
	if a.kind != benchNumberKind || b.kind != benchNumberKind {
		panic("wrong type")
	}
	return benchValue32{number: a.number + b.number, kind: benchNumberKind}
}
func bench32Contribution(value benchValue32) float64 {
	switch value.kind {
	case benchNumberKind, benchBooleanKind:
		return value.number
	case benchStringKind:
		return float64(len(value.reference.(*benchText).units))
	case benchObjectKind:
		return value.reference.(*benchObject).number
	}
	return 0
}
func makeBench32Values(mixed bool) []benchValue32 {
	original := makeBenchValues(mixed)
	out := make([]benchValue32, len(original.dynamic))
	for i, value := range original.dynamic {
		out[i].kind = original.inline[i].kind
		switch v := value.(type) {
		case float64:
			out[i].number = v
		case bool:
			if v {
				out[i].number = 1
			}
		case *benchText, *benchObject:
			out[i].reference = v
		}
	}
	return out
}
func BenchmarkValueSafe32(b *testing.B) {
	b.Run("TypeCheck", func(b *testing.B) {
		values := makeBench32Values(true)
		i := 0
		result := false
		b.ResetTimer()
		for b.Loop() {
			result = values[i&4095].kind == benchNumberKind
			i++
		}
		benchBoolSink = result
	})
	b.Run("NumericRead", func(b *testing.B) {
		values := makeBench32Values(false)
		i := 0
		sum := 0.0
		b.ResetTimer()
		for b.Loop() {
			v := values[i&4095]
			if v.kind != benchNumberKind {
				b.Fatal("wrong type")
			}
			sum += v.number
			i++
		}
		benchNumberSink = sum
	})
	b.Run("NumericUpdate", func(b *testing.B) {
		cell := &benchValue32{number: 1024.125, kind: benchNumberKind}
		benchAnySink = cell
		b.ResetTimer()
		for b.Loop() {
			if cell.kind != benchNumberKind {
				b.Fatal("wrong type")
			}
			cell.number += 1.125
		}
		benchNumberSink = cell.number
	})
	b.Run("MixedDispatch", func(b *testing.B) {
		values := makeBench32Values(true)
		i := 0
		sum := 0.0
		b.ResetTimer()
		for b.Loop() {
			sum += bench32Contribution(values[i&4095])
			i++
		}
		benchNumberSink = sum
	})
	b.Run("Call", func(b *testing.B) {
		a := benchValue32{number: 1024.125, kind: benchNumberKind}
		c := benchValue32{number: 2048.25, kind: benchNumberKind}
		result := benchValue32{}
		for b.Loop() {
			result = bench32Add(a, c)
		}
		bench32Sink = result
	})
	b.Run("ObjectRead", func(b *testing.B) {
		value := benchIdentity32(benchValue32{kind: benchObjectKind, reference: &benchObject{1024.125}})
		sum := 0.0
		for b.Loop() {
			if value.kind != benchObjectKind {
				b.Fatal("wrong type")
			}
			sum += value.reference.(*benchObject).number
		}
		benchNumberSink = sum
	})
	b.Run("BuildNumbers", func(b *testing.B) {
		for b.Loop() {
			values := make([]benchValue32, 1024)
			for i := range values {
				values[i] = benchValue32{number: float64(i) + 0.125, kind: benchNumberKind}
			}
			bench32SliceSink = values
		}
	})
	b.Run("BuildObjects", func(b *testing.B) {
		for b.Loop() {
			values := make([]benchValue32, 256)
			for i := range values {
				values[i] = benchValue32{kind: benchObjectKind, reference: &benchObject{float64(i) + 0.125}}
			}
			bench32SliceSink = values
		}
	})
}

// Identity functions stop the object-read tests from proving the type/tag at
// compile time. Construction is outside the timed loop for every candidate.
//
//go:noinline
func benchIdentityAny(value any) any { return value }

//go:noinline
func benchIdentity24(value benchValue24) benchValue24 { return value }

//go:noinline
func benchIdentity16(value benchValue16) benchValue16 { return value }

//go:noinline
func benchIdentity32(value benchValue32) benchValue32 { return value }
func TestValueSafe32(t *testing.T) {
	t.Logf("safe tagged value=%d bytes", unsafe.Sizeof(benchValue32{}))
	values := makeBench32Values(true)
	original := makeBenchValues(true)
	for range 3 {
		runtime.GC()
	}
	for i, value := range values {
		if bench32Contribution(value) != benchAnyContribution(original.dynamic[i]) {
			t.Fatalf("mismatch at %d", i)
		}
	}
	runtime.KeepAlive(values)
}
