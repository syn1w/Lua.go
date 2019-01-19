package api

//
//  +----------------------------+-----+             +-------------+
//  |          Core Lua          |     |             |             |
//  |                            |     |             |             |
//  |   +-------------------+    | Lua |  <--------- |     Host    |
//  |   |      Lua State    |    | API |             |             |
//  |   +-------------------+    |     |             |             |
//  +----------------------------+-----+             +-------------+
//

// LuaType is a enum type, value LuaTxxx in api/consts.go
type LuaType int

// ArithmeticOp is arithmetic operator
type ArithmeticOp int

// CompareOp is compare operator
type CompareOp int

// ILuaState LuaState interface
type ILuaState interface {
	// basic stack manipulation
	GetTop() int
	AbsIndex(idx int) int
	CheckStack(n int) bool
	Pop(n int)
	Copy(fromIdx, toIdx int)
	PushValue(idx int)
	Replace(idx int)
	Insert(idx int)
	Remove(idx int)
	Rotate(idx, n int)
	SetTop(top int)

	// access methods (stack -> go)
	TypeName(t LuaType) string
	Type(idx int) LuaType
	IsNone(idx int) bool
	IsNil(idx int) bool
	IsNoneOrNil(idx int) bool
	IsBoolean(idx int) bool
	IsInteger(idx int) bool
	IsNumber(idx int) bool
	IsString(idx int) bool
	ToBoolean(idx int) bool
	ToInteger(idx int) int64
	ToIntegerX(idx int) (int64, bool)
	ToNumber(idx int) float64
	ToNumberX(idx int) (float64, bool)
	ToString(idx int) string
	ToStringX(idx int) (string, bool)

	// push methods (go -> stack)
	PushNil()
	PushBoolean(b bool)
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(str string)

	// operator
	Arithmetic(op ArithmeticOp)
	Compare(idx1, idx2 int, op CompareOp) bool
	Len(idx int)
	Concat(n int)

	// table function
	NewTable()
	CreateTable(nArr, nRecord int)

	// table get function
	GetTable(idx int) LuaType
	GetField(idx int, k string) LuaType
	GetI(idx int, i int64) LuaType

	// table set function
	SetTable(idx int)
	SetField(idx int, k string)
	SetI(idx int, n int64)

	// function call
	Load(chunk []byte, chunkName, mode string) int // mode: b(binary), t(text file), bt
	Call(nArgs, nResults int)

	// Go function
	PushGoFunction(f GoFunction)
	IsGoFunction(idx int) bool
	ToGoFunction(idx int) GoFunction
	PushGoClosure(f GoFunction, n int)

	// global table
	PushGlobalTable()
	GetGlobal(name string) LuaType
	SetGlobal(name string)
	Register(name string, f GoFunction)

	// metaTable
	GetMetaTable(idx int) bool
	SetMetaTable(idx int)
	RawLen(idx int) uint
	RawEqual(idx1 int, idx2 int) bool
	RawGet(idx int) LuaType
	RawSet(idx int)
	RawGetI(idx int, i int64) LuaType
	RawSetI(idx int, i int64)

	// iterator
	Next(idx int) bool
}

// GoFunction is called by lua
type GoFunction func(ILuaState) int
