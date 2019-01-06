package api

// ILuaVM is Lua virtual machine interface
type ILuaVM interface {
	ILuaState
	PC() int
	AddPC(n int)   // pc=pc+n
	Fetch() uint32 // fetch current instruction, pc=pc+1
	GetConst(idx int)
	GetRK(rk int)

	RegisterCount() int
	LoadVararg(n int)
	LoadProto(idx int)
}
