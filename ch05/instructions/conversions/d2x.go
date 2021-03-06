package conversions

import (
	"jvmgo/ch05/instructions/base"
	"jvmgo/ch05/rtda"
)

type D2F struct {base.NoOperandsInstruction}


func (d2f *D2F) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopDouble()
	f := float32(d)
	stack.PushFloat(f)
}

type D2I struct {base.NoOperandsInstruction}

func (d2i *D2I) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopDouble()
	i := int32(d)
	stack.PushInt(i)
}

type D2L struct {base.NoOperandsInstruction}

func (d2l *D2L) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	d := stack.PopDouble()
	l := int64(d)
	stack.PushLong(l)
}