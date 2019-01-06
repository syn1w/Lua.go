package state

// ------------------------------------
//       basic stack manipulation
// ------------------------------------

// GetTop gets stack.top
func (s *LuaState) GetTop() int {
	return s.stack.top
}

// AbsIndex converts idx to absidx
func (s *LuaState) AbsIndex(idx int) int {
	return s.stack.absIndex(idx)
}

// CheckStack avoids stack overflow
func (s *LuaState) CheckStack(n int) bool {
	s.stack.check(n)
	return true // always true
}

// Pop n elements from stack, panic when there are not enough elements
func (s *LuaState) Pop(n int) {
	for i := 0; i < n; i++ {
		s.stack.pop()
	}
}

// Copy <==> stack[toIdx] = stack[fromIdx]
func (s *LuaState) Copy(fromIdx, toIdx int) {
	val := s.stack.get(fromIdx)
	s.stack.set(toIdx, val)
}

// PushValue pushes the element with idx
func (s *LuaState) PushValue(idx int) {
	val := s.stack.get(idx)
	s.stack.push(val)
}

// Replace <==> pop the top element and copy it to idx position
func (s *LuaState) Replace(idx int) {
	val := s.stack.pop()
	s.stack.set(idx, val)
}

// Insert <==> pop the top element and insert it to idx position
func (s *LuaState) Insert(idx int) {
	s.Rotate(idx, 1)
}

// Remove the element in idx position
func (s *LuaState) Remove(idx int) {
	s.Rotate(idx, -1)
	s.Pop(1)
}

// Rotate [idx, top] elements |n| steps(n >= 0 right or up, n < 0 left or down)
func (s *LuaState) Rotate(idx, n int) {
	b := s.stack.absIndex(idx) - 1 // begin
	e := s.stack.top - 1           // end
	var m int                      // middle

	if n >= 0 { // right rotate
		m = e - n
	} else { // left rotate
		m = b - n - 1
	}

	s.stack.reverse(b, m)
	s.stack.reverse(m+1, e)
	s.stack.reverse(b, e)
}

// SetTop set new top, equivalent to push or pop
func (s *LuaState) SetTop(top int) {
	newTop := s.stack.absIndex(top)
	if newTop < 0 {
		panic("stack underflow")
	}
	n := s.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			s.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			s.stack.push(nil)
		}
	}
}
