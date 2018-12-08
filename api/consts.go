package api

// Lua Type constants
const (
	LuaTNone = iota - 1 // -1
	LuaTNil
	LuaTBoolean
	LuaTLightUserData
	LuaTNumber
	LuaTString
	LuaTTable
	LuaTFunction
	LuaTUserData
	LuaTThread
)

// Lua arithmetic operator
const (
	LuaOpAdd  = iota // +
	LuaOpSub         // -(binary)
	LuaOpMul         // *
	LuaOpMod         // %
	LuaOpPow         // ^
	LuaOpDiv         // /
	LuaOpIDiv        // //
	LuaOpBAnd        // &
	LuaOpBOr         // |
	LuaOpBXor        // ~
	LuaOpBShl        // <<
	LuaOpBShr        // >>
	LuaOpUnm         // -(unary)
	LuaOpBNot        // ~
)

// Lua compare operator
const (
	LuaOpEq = iota // ==
	LuaOpLt        // <
	LuaOpLe        // <=
)
