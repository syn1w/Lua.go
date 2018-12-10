package state

// PC returns current pc
func (s *LuaState) PC() int {
	return s.pc
}

// AddPC <=> pc+=n
func (s *LuaState) AddPC(n int) {
	s.pc += n
}

// Fetch current instruction and pc++
func (s *LuaState) Fetch() uint32 {
	code := s.proto.Code[s.pc]
	s.pc++
	return code
}

// GetConst pushes the constant
func (s *LuaState) GetConst(idx int) {
	c := s.proto.Constants[idx]
	s.stack.push(c)
}

// GetRK pushes the constant or the stack value
func (s *LuaState) GetRK(rk int) {
	if rk > 0xFF { // constant
		s.GetConst(rk & 0xFF)
	} else { // register
		s.PushValue(rk + 1) // rk 从 0 开始，而栈索引是从 1 开始
	}
}
