package math

import (
	"jvmgo/ch10/instructions/base"
	"jvmgo/ch10/rtda"
)

type IAND struct {
	base.NoOperandsInstruction
}

func (iand *IAND) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopInt()
	v1 := stack.PopInt()
	result := v1 & v2
	stack.PushInt(result)
}

type LAND struct {
	base.NoOperandsInstruction
}

func (lang *LAND) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	v2 := stack.PopLong()
	v1 := stack.PopLong()
	result := v1 & v2
	stack.PushLong(result)

}
