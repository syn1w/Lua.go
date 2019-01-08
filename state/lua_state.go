package state

import "vczn/luago/api"

// LuaState impl api.ILuaState
type LuaState struct {
	registry *LuaTable
	stack    *LuaStack
}

// NewLuaState new a LuaState
func NewLuaState() *LuaState {
	registry := NewLuaTable(0, 0)
	registry.put(api.LuaRidxGlobals, NewLuaTable(0, 0))
	luastate := &LuaState{
		registry: registry,
	}

	luastate.pushLuaStack(newLuaStack(api.LuaMinStack, luastate))

	return luastate
}

func (s *LuaState) pushLuaStack(stack *LuaStack) {
	stack.prev = s.stack
	s.stack = stack
}

func (s *LuaState) popLuaStack() {
	stack := s.stack
	s.stack = s.stack.prev
	stack.prev = nil
}
