package base

import (
	"jvmgo/ch11/rtda"
	"jvmgo/ch11/rtda/heap"
)

func InvokeMethod(frame *rtda.Frame, method *heap.Method) {
	thread := frame.Thread()
	newFrame := thread.NewFrame(method)
	thread.PushFrame(newFrame)

	argSlotSlot := int(method.ArgSlotCount())
	if argSlotSlot > 0 {
		for i := argSlotSlot-1; i >= 0; i-- {
			slot := frame.OperandStack().PopSlot()
			newFrame.LocalVars().SetSlot(uint(i),slot)
		}
	}
}