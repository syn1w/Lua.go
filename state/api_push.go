package state

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
