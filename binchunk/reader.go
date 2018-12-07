package binchunk

import (
	"encoding/binary"
	"math"
)

// Reader is for reading binary chunk file
type Reader struct {
	data []byte
}

func (r *Reader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

func (r *Reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return i
}

func (r *Reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return i
}

func (r *Reader) readLuaInteger() int64 {
	return int64(r.readUint64())
}

func (r *Reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}

// string
// NULL           0
// n<=0xFD(253)   n+1(1byte)  | bytes
// n>=0xFE(254)   0xFF | n+1(8bytes)    | bytes
func (r *Reader) readString() string {
	size := uint64(r.readByte())
	if size == 0x0 { // NULL
		return ""
	}
	if size == 0xFF { // long string
		size = r.readUint64()
	}
	str := r.readBytes(size - 1)
	return string(str)
}

func (r *Reader) readBytes(n uint64) []byte {
	bs := r.data[:n]
	r.data = r.data[n:]
	return bs
}

func (r *Reader) readCode() []uint32 {
	code := make([]uint32, r.readUint32())
	for i := range code {
		code[i] = r.readUint32()
	}
	return code
}

func (r *Reader) readConstant() interface{} {
	switch r.readByte() {
	case tagNil:
		return nil
	case tagBoolean:
		return r.readByte() != 0
	case tagInteger:
		return r.readLuaInteger()
	case tagNumber:
		return r.readLuaNumber()
	case tagShortString, tagLongString:
		return r.readString()
	default:
		panic("corrupted")
	}
}

func (r *Reader) readConstants() []interface{} {
	constants := make([]interface{}, r.readUint32())
	for i := range constants {
		constants[i] = r.readConstant()
	}
	return constants
}

func (r *Reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: r.readByte(),
			Idx:     r.readByte(),
		}
	}
	return upvalues
}

func (r *Reader) readProtos(parentSource string) []*ProtoType {
	protos := make([]*ProtoType, r.readUint32())
	for i := range protos {
		protos[i] = r.readProto(parentSource)
	}
	return protos
}

func (r *Reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, r.readUint32())
	for i := range lineInfo {
		lineInfo[i] = r.readUint32()
	}
	return lineInfo
}

func (r *Reader) readLocVars() []LocVar {
	locVars := make([]LocVar, r.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: r.readString(),
			StartPC: r.readUint32(),
			EndPC:   r.readUint32(),
		}
	}
	return locVars
}

func (r *Reader) readUpvalueNames() []string {
	names := make([]string, r.readUint32())
	for i := range names {
		names[i] = r.readString()
	}
	return names
}

func (r *Reader) readProto(parentSource string) *ProtoType {
	source := r.readString()
	if source == "" {
		source = parentSource
	}

	return &ProtoType{
		Source:          source,
		LineDefined:     r.readUint32(),
		LastLineDefined: r.readUint32(),
		NumParams:       r.readByte(),
		IsVararg:        r.readByte(),
		MaxStackSize:    r.readByte(),
		Code:            r.readCode(),
		Constants:       r.readConstants(),
		Upvalues:        r.readUpvalues(),
		Protos:          r.readProtos(source),
		LineInfo:        r.readLineInfo(),
		LocVars:         r.readLocVars(),
		UpvalueNames:    r.readUpvalueNames(),
	}
}

func (r *Reader) checkHeader() {
	if string(r.readBytes(4)) != luaSignature {
		panic("not a precompiled chunk!")
	}
	if r.readByte() != luacVersion {
		panic("version mismatch!")
	}
	if r.readByte() != luacFormat {
		panic("format mismatch!")
	}
	if string(r.readBytes(6)) != luacData {
		panic("corrupted!")
	}
	if r.readByte() != cintSize {
		panic("int size mismatch!")
	}
	if r.readByte() != csizetSize {
		panic("size_t size mismatch!")
	}
	if r.readByte() != instructionSize {
		panic("instruction size mismatch!")
	}
	if r.readByte() != luaIntegerSize {
		panic("lua integer mismatch!")
	}
	if r.readByte() != luaNumberSize {
		panic("lua number mismatch!")
	}
	if r.readLuaInteger() != luacInt {
		panic("endianness mismatch!")
	}
	if r.readLuaNumber() != luacNumber {
		panic("float format mismatch!")
	}
}
