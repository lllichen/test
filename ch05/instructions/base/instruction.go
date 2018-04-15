package base

import "jvmgo/ch05/rtda"

type Instruction interface {
	FetchOperands(reader *ByteCodeReader)
	Execute(frame *rtda.Frame)
}

type NoOperandsInstruction struct {

}

func (noi *NoOperandsInstruction) FetchOperands(reader *ByteCodeReader){

}


type BranchInstruction struct {
	offset int
}


func (bi *BranchInstruction) FetchOperands(reader *ByteCodeReader){
bi.offset = int(reader.ReadInt16())
}

type Index8Instruction struct {
	Index uint
}

func (index8Instruction *Index8Instruction) FetchOperands(reader *ByteCodeReader){
	index8Instruction.Index = uint(reader.readUint8())
}

type Index16Instruction struct {
	Index uint
}

func (index16Instruction *Index16Instruction) FetchOperands(reader *ByteCodeReader){
	index16Instruction.Index = uint(reader.readUint16())
}