package state

import (
	"vczn/luago/api"
)

//     +-------+
//     |       | 6 invalid  5  <-top
// -1  |   e   | 5          4
// -2  |   d   | 4          3   acceptable
// -3  |   c   | 3          2
// -4  |   b   | 2 valid    1
// -5  |   a   | 1          0
//     +-------+ 0
//              lua         go
// rIdx        absIdx

// LuaStack is the Lua stack struct
type LuaStack struct {
	// virtual stack
	slots []LuaValue
	top   int

	// call info
	closure *luaClosure
	varargs []LuaValue
	pc      int

	// linked list
	prev *LuaStack

	// state
	state *LuaState
}

func newLuaStack(size int, luastate *LuaState) *LuaStack {
	return &LuaStack{
		slots: make([]LuaValue, size),
		top:   0,
		state: luastate,
	}
}

// check if there is enough space to place size elements
func (s *LuaStack) check(size int) {
	free := len(s.slots) - s.top
	if free < size {
		s.slots = append(s.slots, make([]LuaValue, size-free)...)
	}
	// for i := free; i < n; i++ {
	// 	s.slots = append(s.slots, nil)
	// }
}

func (s *LuaStack) empty() bool {
	return s.top == 0
}

func (s *LuaStack) full() bool {
	return s.top == len(s.slots)
}

func (s *LuaStack) push(val LuaValue) {
	if s.top == len(s.slots) {
		panic("stack overflow")
	}
	s.slots[s.top] = val
	s.top++
}

func (s *LuaStack) pushN(vals []LuaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}

	for i := 0; i < n; i++ {
		if i < nVals {
			s.push(vals[i])
		} else {
			s.push(nil)
		}
	}
}

func (s *LuaStack) pop() LuaValue {
	if s.top < 1 {
		panic("stack underflow")
	}
	s.top--
	val := s.slots[s.top]
	s.slots[s.top] = nil
	return val
}

func (s *LuaStack) popN(n int) []LuaValue {
	vals := make([]LuaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = s.pop()
	}

	return vals
}

func (s *LuaStack) absIndex(idx int) int {
	if idx <= api.LuaRegistryIndex { // pseudo index
		return idx
	}

	if idx >= 0 {
		return idx
	}
	return idx + s.top + 1
}

func (s *LuaStack) absIdxIsValid(absIdx int) bool {
	return absIdx > 0 && absIdx <= s.top
}

func (s *LuaStack) isValid(idx int) bool {
	if idx == api.LuaRegistryIndex {
		return true
	}

	absIdx := s.absIndex(idx)
	return s.absIdxIsValid(absIdx)
}

func (s *LuaStack) get(idx int) LuaValue {
	if idx == api.LuaRegistryIndex {
		return s.state.registry
	}

	absIdx := s.absIndex(idx)
	if !s.absIdxIsValid(absIdx) {
		return nil // ?panic
	}
	return s.slots[absIdx-1]
}

func (s *LuaStack) set(idx int, val LuaValue) {
	if idx == api.LuaRegistryIndex {
		s.state.registry = val.(*LuaTable)
		return
	}

	absIdx := s.absIndex(idx)
	if !s.absIdxIsValid(absIdx) {
		panic("invalid index!")
	}
	s.slots[absIdx-1] = val
}

func (s *LuaStack) reverse(from, to int) {
	for from < to {
		s.slots[from], s.slots[to] = s.slots[to], s.slots[from]
		from++
		to--
	}
}
