package classfile

import "fmt"

type ClassFile struct {
	//magic      uint32
	minorVersion uint16
	majorVersion uint16
	constantPool ConstantPool
	accessFlags uint16
	thisClass uint16
	superClass uint16
	interfaces []uint16
	fields  []*MemberInfo
	methods []*MemberInfo
	attributes []AttributeInfo

}


func Parse(classData []byte)(cf *ClassFile,err error) {
	//interface assertion
	defer func(){
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			// check err ,not err ,format to err
			if !ok {
				err = fmt.Errorf("%v",r)
			}
		}
	}()
	cr := &ClassReader{classData}
	cf = &ClassFile{}
	cf.read(cr)
	return
}

func (cf *ClassFile) read(reader *ClassReader) {
	cf.readAndCheckMagic(reader)
	cf.readAndCheckVersion(reader)
}

func (cf *ClassFile) MajorVersion() uint16{
	return cf.majorVersion
}
func (cf *ClassFile) readAndCheckMagic(reader *ClassReader) {
	magic := reader.readUint32()
	if magic != 0XCAFEBABE {
		panic("java.lang.ClassFormatError: magic!")
	}
}
func (cf *ClassFile) readAndCheckVersion(reader *ClassReader) {
	cf.minorVersion = reader.readUint16()
	cf.majorVersion = reader.readUint16()
	fmt.Printf("minorVersion is %d\n",cf.minorVersion)
	fmt.Printf("majorVersion is %d\n",cf.majorVersion)
	switch cf.majorVersion {
	case 45:
		return
	case 46,47,48,49,50,51,52:
		if cf.minorVersion == 0 {
			return
		}
	}
	panic("java.lang.UnsupportedClassVersionError!")
}







