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

// LuaState constants
const (
	LuaMinStack = 20
	LuaMaxStack = 1000000

	//      -max                                       max
	//        |                                         |
	//     |__|____________________|____________________|
	//     |                       0
	// -max-1000
	LuaRegistryIndex = -LuaMaxStack - 1000

	LuaRidxGlobals int64 = 2 // index in global
)
