package state

// PC returns current pc
func (s *LuaState) PC() int {
	return s.stack.pc
}

// AddPC <=> pc+=n
func (s *LuaState) AddPC(n int) {
	s.stack.pc += n
}

// Fetch current instruction and pc++
func (s *LuaState) Fetch() uint32 {
	code := s.stack.closure.proto.Code[s.stack.pc]
	s.stack.pc++
	return code
}

// GetConst pushes the constant
func (s *LuaState) GetConst(idx int) {
	c := s.stack.closure.proto.Constants[idx]
	s.stack.push(c)
}

// GetRK pushes the constant or the stack value
func (s *LuaState) GetRK(rk int) {
	if rk > 0xFF { // constant
		s.GetConst(rk & 0xFF)
	} else { // register
		s.PushValue(rk + 1) //  rk begin with 0, stack index begin with 1
	}
}

// RegisterCount returns the count the registers
func (s *LuaState) RegisterCount() int {
	return int(s.stack.closure.proto.MaxStackSize)
}

// LoadVararg loads varargs
func (s *LuaState) LoadVararg(n int) {
	if n < 0 {
		n = len(s.stack.varargs)
	}

	s.stack.check(n)
	s.stack.pushN(s.stack.varargs, n)
}

// LoadProto loads function prototype
func (s *LuaState) LoadProto(idx int) {
	proto := s.stack.closure.proto.Protos[idx]
	closure := newLuaClosure(proto)
	s.stack.push(closure)
}
