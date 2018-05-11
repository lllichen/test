package references

import (
	"jvmgo/ch06/instructions/base"
	"jvmgo/ch06/rtda"
	"jvmgo/ch06/rtda/heap"
)

type Get_FIELD struct {
	base.Index16Instruction
}

func (getField *Get_FIELD) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	fieldRef := cp.GetConstant(getField.Index).(*heap.FieldRef)
	field := fieldRef.ResolvedFieldRef()
	if field.IsStatic(){
		panic("java.lang.IncompatibleChangeError")
	}
	stack := frame.OperandStack()
	ref := stack.PopRef()
	if ref == nil {
		panic("java.lang.NullPointerException")
	}

	descriptor := field.Descriptor()
	slotId := field.SlotId()
	slots := ref.Fields()
	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		stack.PushInt(slots.GetInt(slotId))
	case 'F':
		stack.PushFloat(slots.GetFloat(slotId))
	case 'J':
		stack.PushLong(slots.GetLong(slotId))
	case 'D':
		stack.PushDouble(slots.GetDouble(slotId))
	case 'L':
		stack.PushRef(slots.GetRef(slotId))
	}
}

