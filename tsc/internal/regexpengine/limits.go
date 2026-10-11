package regexpengine

import "fmt"

// ExecutionLimits bounds work for callers that run patterns against
// untrusted input. MaxBacktrackBytes is a conservative estimate of active and
// saved VM capture/map state, including Go map/slice overhead.
type ExecutionLimits struct {
	MaxSteps          uint64
	MaxBacktrackBytes uint64
	MaxAssertionDepth int
}

// DefaultExecutionLimits returns finite limits suitable for normal runtime
// requests. Callers with trusted workloads can explicitly choose larger
// limits with FindUTF16WithLimits or FindAllWithLimits.
func DefaultExecutionLimits() ExecutionLimits {
	return ExecutionLimits{
		MaxSteps:          10_000_000,
		MaxBacktrackBytes: 64 << 20,
		MaxAssertionDepth: 128,
	}
}

// ExecutionLimitError reports a resource ceiling reached during matching.
// This is distinct from a normal no-match result.
type ExecutionLimitError struct {
	Resource string
	Maximum  uint64
}

func (e *ExecutionLimitError) Error() string {
	return fmt.Sprintf("regular expression %s limit exceeded (maximum %d)", e.Resource, e.Maximum)
}

type executionBudget struct {
	limits         ExecutionLimits
	steps          uint64
	backtrackBytes uint64
	auxiliaryBytes uint64
	assertionDepth int
	err            error
}

func newExecutionBudget(limits ExecutionLimits) (*executionBudget, error) {
	if limits.MaxSteps == 0 || limits.MaxBacktrackBytes == 0 || limits.MaxAssertionDepth <= 0 {
		return nil, fmt.Errorf("regular expression execution limits must all be positive")
	}
	return &executionBudget{limits: limits}, nil
}

func (b *executionBudget) step() error {
	if b == nil {
		return nil
	}
	if b.err != nil {
		return b.err
	}
	if b.steps >= b.limits.MaxSteps {
		b.err = &ExecutionLimitError{Resource: "step", Maximum: b.limits.MaxSteps}
		return b.err
	}
	b.steps++
	return nil
}

func (b *executionBudget) enterAssertion() (func(), error) {
	if b == nil {
		return func() {}, nil
	}
	if b.assertionDepth >= b.limits.MaxAssertionDepth {
		b.err = &ExecutionLimitError{Resource: "assertion depth", Maximum: uint64(b.limits.MaxAssertionDepth)}
		return nil, b.err
	}
	b.assertionDepth++
	return func() { b.assertionDepth-- }, nil
}

func executionStateBytes(state executionState) uint64 {
	// The estimate includes the struct, slice backing storage, and generous
	// per-entry allowances for the two maps cloned with each backtrack state.
	return executionFrameBytes(len(state.captures), len(state.emptySplits), len(state.loopStarts))
}

func executionFrameBytes(captureCount, splitCount, loopCount int) uint64 {
	// Include map headers and a conservative allowance for each entry/bucket.
	return 128 + uint64(captureCount)*8 + uint64(splitCount)*32 + uint64(loopCount)*24
}

func (b *executionBudget) adjustAuxiliary(oldBytes, newBytes uint64) error {
	if b == nil {
		return nil
	}
	if b.err != nil {
		return b.err
	}
	retained := b.backtrackBytes + b.auxiliaryBytes
	if retained > b.limits.MaxBacktrackBytes || newBytes > oldBytes && newBytes-oldBytes > b.limits.MaxBacktrackBytes-retained {
		b.err = &ExecutionLimitError{Resource: "execution memory", Maximum: b.limits.MaxBacktrackBytes}
		return b.err
	}
	b.auxiliaryBytes = b.auxiliaryBytes - oldBytes + newBytes
	return nil
}

func (b *executionBudget) checkStateSize(captureCount, splitCount, loopCount int) error {
	if b == nil {
		return nil
	}
	if b.err != nil {
		return b.err
	}
	bytes := executionFrameBytes(captureCount, splitCount, loopCount)
	retained := b.backtrackBytes + b.auxiliaryBytes
	if retained > b.limits.MaxBacktrackBytes || bytes > b.limits.MaxBacktrackBytes-retained {
		b.err = &ExecutionLimitError{Resource: "backtracking memory", Maximum: b.limits.MaxBacktrackBytes}
		return b.err
	}
	return nil
}

func (b *executionBudget) pushState(state executionState) error {
	if b == nil {
		return nil
	}
	bytes := executionStateBytes(state)
	retained := b.backtrackBytes + b.auxiliaryBytes
	if retained > b.limits.MaxBacktrackBytes || bytes > b.limits.MaxBacktrackBytes-retained {
		b.err = &ExecutionLimitError{Resource: "backtracking memory", Maximum: b.limits.MaxBacktrackBytes}
		return b.err
	}
	b.backtrackBytes += bytes
	return nil
}

func (b *executionBudget) popState(state executionState) {
	if b == nil {
		return
	}
	b.backtrackBytes -= executionStateBytes(state)
}

func (b *executionBudget) replaceState(old, updated executionState) error {
	if b == nil {
		return nil
	}
	oldBytes, newBytes := executionStateBytes(old), executionStateBytes(updated)
	retained := b.backtrackBytes + b.auxiliaryBytes
	if retained > b.limits.MaxBacktrackBytes || newBytes > oldBytes && newBytes-oldBytes > b.limits.MaxBacktrackBytes-retained {
		b.err = &ExecutionLimitError{Resource: "backtracking memory", Maximum: b.limits.MaxBacktrackBytes}
		return b.err
	}
	b.backtrackBytes = b.backtrackBytes - oldBytes + newBytes
	return nil
}
