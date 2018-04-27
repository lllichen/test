package heap

import "jvmgo/ch06/classfile"

type Field struct {
	ClassMember
}

func newFields (class *Class, cfFields []*classfile.MemberInfo ) []*Field{
	fields := make ([]*Field ,len(cfFields))
	for i, cfFields := range cfFields {
		fields[i] = &Field{}
		fields[i].class = class
		fields[i].copyMemberInfo(cfFields)
	}
	return fields
}
