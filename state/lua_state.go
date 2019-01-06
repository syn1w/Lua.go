package state

// LuaState impl api.ILuaState
type LuaState struct {
	stack *LuaStack
}

// NewLuaState new a LuaState
func NewLuaState() *LuaState {
	return &LuaState{
		stack: newLuaStack(20),
	}
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
