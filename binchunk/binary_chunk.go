package binchunk

type binaryChunk struct {
	header
	sizeUpvalues byte // number of upvalues
	mainFunc     *ProtoType
}

// xxd ./luac.out 查看 header

// windows 下的 header 格式不一样，下面使用 Linux Lua 5.3
type header struct {
	signature       [4]byte // magic number: ESC, L, u, a -> 0x1B4C7561 -> "\x1bLua"
	version         byte    // 5.3.4 -> 0x53
	format          byte    // 00
	luacData        [6]byte // 1993 0D(\r)0A(\n) 1A0A in Unix, 0104 0404 0800 Windows
	cintSize        byte    // 04
	sizetSize       byte    // 08
	instructionSize byte    // 04
	luaIntegerSize  byte    // 08
	luaNumberSize   byte    // 08
	luacInt         int64   // 78 56 00 00 00 00 00 00 -> 5678 little endian
	luacNum         float64 // 00 00 00 00 00 28 77 40 -> IEEE 754 float64 370.5
}

// ProtoType 函数原型：函数的基本信息、指令表、常量表、upvalue 表、子函数原型、调试信息, ...
type ProtoType struct {
	Source          string        // source file name
	LineDefined     uint32        // 起始行 main == 0, other > 0
	LastLineDefined uint32        // last line
	NumParams       byte          // 固定参数个数 main == 0
	IsVararg        byte          // 2: declared vararg; 1: uses vararg; 0 not vararg
	MaxStackSize    byte          // number of stacks
	Code            []uint32      // 指令列表，每条指令 4 bytes
	Constants       []interface{} // 常量表，字面量 nil, boolean, number, integer, string
	Upvalues        []Upvalue     // upvalues table, 2 bytes per element
	Protos          []*ProtoType  // sub funtion proto table
	LineInfo        []uint32
	LocVars         []LocVar
	UpvalueNames    []string
}

// header Linux Lua 5.3
const (
	luaSignature    = "\x1bLua"
	luacVersion     = 0x53
	luacFormat      = 0
	luacData        = "\x19\x93\x0d\x0a\x1a\x0a" // 1993 \r\n \x1a \n
	cintSize        = 0x04
	csizetSize      = 0x08
	instructionSize = 0x04
	luaIntegerSize  = 0x08
	luaNumberSize   = 0x08
	luacInt         = 0x5678
	luacNumber      = 370.5
)

// 0x01 number of upvalues

// Source
// 40 @
// 0A len+1  "@test.lua"  // src file @ 表明确实是从 Lua 源文件编译而来，其他类似有 =stdin

// 40 74 65 73 74 2e 6c 75 61 @test.lua
// 0000 0000 LineDefined
// 0000 0000 LastLineDefined
// 00        NumParams
// 02        IsVararg
// 02        MaxStackSize

// Code
// TODO
//                                04 00 00 00 06 00 40 00 41         ......@.A
// 00000040: 40 00 00 24 40 00 01 26 00 80 00 02 00 00 00 04  @..$@..&........
// 00000050: 06 70 72 69 6e 74 04 06                          .print..

// constants table
const (
	tagNil         = 0x00
	tagBoolean     = 0x01
	tagNumber      = 0x03
	tagInteger     = 0x13
	tagShortString = 0x04
	tagLongString  = 0x14
)

// 68 65 6c 6c 6f "hello"

// upvalues table
// 01 00

// Upvalue type
type Upvalue struct {
	Instack byte
	Idx     byte
}

// Protos
// 00 00 表示无 subfunction

// lineInfo
// 行号和指令表中的指令一一对应，分别记录每条在源代码中对应的 lineno
// 01 00 00 00 -> line 1, 00 00 04 00 00 00
// 01 00 00 00 -> line 1
// 01 00 00 00 -> line 1
// 01 00 00 00 -> line 1

// LocVar is local variables struct
// 01 00 00 00 00 00 00 00 01 00 00 00 80 00 00 00 05
type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

// UpvalueNames
// 5f 45 4e 56 _ENV

// Undump is for parsing binary chunk file to generate *ProtoType info
func Undump(data []byte) *ProtoType {
	reader := &Reader{data}
	reader.checkHeader() // check header
	reader.readByte()    // 跳过 Upvalue numbers
	return reader.readProto("")
}
