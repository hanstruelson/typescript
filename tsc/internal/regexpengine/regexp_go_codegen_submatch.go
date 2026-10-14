package regexpengine

import "strconv"

// These routines lower V8's submatch bookkeeping to the generated matcher's
// own backtracking stack. A lookaround records the stack depth and input
// position at entry. Positive success commits the body and resumes at the
// saved position; negative success discards the body and fails the enclosing
// choice.
func (emitter *v8GoGraphEmitter) emitBeginSubmatch(node *v8NodesActionNode) {
	emitter.text.WriteString("regs[")
	emitter.text.WriteString(strconv.Itoa(node.submatch.stackPointerRegister))
	emitter.text.WriteString("] = len(stack); regs[")
	emitter.text.WriteString(strconv.Itoa(node.submatch.currentPositionRegister))
	emitter.text.WriteString("] = pos")
	if emitter.hasCaptureCleanup() {
		emitter.text.WriteString("; rxCaptureCleanupMarks[")
		emitter.text.WriteString(strconv.Itoa(node.submatch.stackPointerRegister))
		emitter.text.WriteString("] = len(rxCaptureCleanups)")
	}
	emitter.text.WriteString("\n")
}

func (emitter *v8GoGraphEmitter) emitPositiveSubmatchSuccess(node *v8NodesActionNode) {
	stackRegister := strconv.Itoa(node.submatch.stackPointerRegister)
	positionRegister := strconv.Itoa(node.submatch.currentPositionRegister)
	emitter.text.WriteString("if regs[")
	emitter.text.WriteString(stackRegister)
	emitter.text.WriteString("] < 0 || regs[")
	emitter.text.WriteString(stackRegister)
	emitter.text.WriteString("] > len(stack) { goto rxBacktrack }; stack = stack[:regs[")
	emitter.text.WriteString(stackRegister)
	emitter.text.WriteString("]]; pos = regs[")
	emitter.text.WriteString(positionRegister)
	emitter.text.WriteString("]; if pos < 0 || pos > len(input) { goto rxBacktrack }\n")
	if emitter.hasCaptureCleanup() {
		emitter.text.WriteString("rxCaptureCleanups = rxCaptureCleanups[:rxCaptureCleanupMarks[")
		emitter.text.WriteString(stackRegister)
		emitter.text.WriteString("]]\n")
		if node.submatch.clearRegisterCount > 0 {
			emitter.text.WriteString("rxCaptureCleanups = append(rxCaptureCleanups, rxCaptureCleanup{len(stack), ")
			emitter.text.WriteString(strconv.Itoa(node.submatch.clearRegisterFrom))
			emitter.text.WriteString(", ")
			emitter.text.WriteString(strconv.Itoa(node.submatch.clearRegisterFrom + node.submatch.clearRegisterCount))
			emitter.text.WriteString("})\n")
		}
	}
}

func (emitter *v8GoGraphEmitter) emitNegativeSubmatchSuccess(node *v8NodesNegativeSubmatchSuccess) {
	stackRegister := strconv.Itoa(node.stackPointerRegister)
	positionRegister := strconv.Itoa(node.currentPositionRegister)
	emitter.text.WriteString("if regs[")
	emitter.text.WriteString(stackRegister)
	emitter.text.WriteString("] < 0 || regs[")
	emitter.text.WriteString(stackRegister)
	emitter.text.WriteString("] > len(stack) { goto rxBacktrack }; stack = stack[:regs[")
	emitter.text.WriteString(stackRegister)
	emitter.text.WriteString("]]; pos = regs[")
	emitter.text.WriteString(positionRegister)
	emitter.text.WriteString("]; if pos < 0 || pos > len(input) { goto rxBacktrack }\n")
	if emitter.hasCaptureCleanup() {
		emitter.text.WriteString("rxCaptureCleanups = rxCaptureCleanups[:rxCaptureCleanupMarks[")
		emitter.text.WriteString(stackRegister)
		emitter.text.WriteString("]]\n")
	}
	for register := node.clearCaptureStart; register < node.clearCaptureStart+node.clearCaptureCount; register++ {
		emitter.text.WriteString("regs[")
		emitter.text.WriteString(strconv.Itoa(register))
		emitter.text.WriteString("] = -1\n")
	}
	emitter.text.WriteString("goto rxBacktrack\n")
}

func (emitter *v8GoGraphEmitter) emitEmptyMatchCheck(node *v8NodesActionNode) {
	startRegister := strconv.Itoa(node.emptyMatch.startRegister)
	repetitionRegister := node.emptyMatch.repetitionRegister
	if repetitionRegister != v8CompilerNoRegister {
		emitter.text.WriteString("if regs[")
		emitter.text.WriteString(strconv.Itoa(repetitionRegister))
		emitter.text.WriteString("] > ")
		emitter.text.WriteString(strconv.Itoa(node.emptyMatch.repetitionLimit))
		emitter.text.WriteString(" && regs[")
		emitter.text.WriteString(startRegister)
		emitter.text.WriteString("] == pos { goto rxBacktrack }\n")
	} else {
		emitter.text.WriteString("if regs[")
		emitter.text.WriteString(startRegister)
		emitter.text.WriteString("] == pos { goto rxBacktrack }\n")
	}
}
