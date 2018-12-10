package state

import (
	"vczn/luago/binchunk"
)

// LuaState impl api.ILuaState
type LuaState struct {
	stack *LuaStack
	proto *binchunk.ProtoType
	pc    int
}

// NewLuaState new a LuaState
func NewLuaState(sizeStack int, proto *binchunk.ProtoType) *LuaState {
	return &LuaState{
		stack: newLuaStack(20), // TODO
		proto: proto,
		pc:    0,
	}
}
