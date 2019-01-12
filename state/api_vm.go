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
	stack := s.stack
	subProto := stack.closure.proto.Protos[idx]
	closure := newLuaClosure(subProto)
	stack.push(closure)

	for i, upvalInfo := range subProto.Upvalues {
		uvIdx := int(upvalInfo.Idx)

		if upvalInfo.Instack == 1 { // capture local variable in stack of parent function
			if stack.openuvs == nil { // init openUpvalues
				stack.openuvs = map[int]*upvalue{}
			}
			if openuv, found := stack.openuvs[uvIdx]; found { // search upvalue in openUpvalues
				closure.upvals[i] = openuv
			} else { // !found
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				stack.openuvs[uvIdx] = closure.upvals[i]
			}
		} else { //  upvalue has captured by parent
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}
	}
}

// CloseUpvalues closes all upvalues >= a-1
func (s *LuaState) CloseUpvalues(a int) {
	for i, openuv := range s.stack.openuvs {
		if i < a-1 {
			continue
		}

		val := *openuv.val         // copy from register
		openuv.val = &val          // update upvalue
		delete(s.stack.openuvs, i) // close upvalue
	}
}
