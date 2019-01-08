package state

import (
	"vczn/luago/api"
)

// ------------------------------------
//      push methods(go -> stack)
// ------------------------------------

// PushNil pushes nil
func (s *LuaState) PushNil() {
	s.stack.push(nil)
}

// PushBoolean pushes boolean value
func (s *LuaState) PushBoolean(b bool) {
	s.stack.push(b)
}

// PushInteger pushes Lua Integer value
func (s *LuaState) PushInteger(n int64) {
	s.stack.push(n)
}

// PushNumber pushes Lua Number value
func (s *LuaState) PushNumber(n float64) {
	s.stack.push(n)
}

// PushString pushes string value
func (s *LuaState) PushString(str string) {
	s.stack.push(str)
}

// PushGoFunction pushes go function into stack
func (s *LuaState) PushGoFunction(goFunc api.GoFunction) {
	s.stack.push(newGoClosure(goFunc))
}

// IsGoFunction returns stack[idx] if is go function
func (s *LuaState) IsGoFunction(idx int) bool {
	val := s.stack.get(idx)
	if c, ok := val.(*luaClosure); ok {
		return c.goFunc != nil
	}
	return false
}

// ToGoFunction converts stack[idx] to GoFunction
func (s *LuaState) ToGoFunction(idx int) api.GoFunction {
	val := s.stack.get(idx)
	if c, ok := val.(*luaClosure); ok {
		return c.goFunc
	}

	return nil
}
