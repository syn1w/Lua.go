package state

import (
	"fmt"
	"vczn/luago/api"
)

// LuaState impl api.ILuaState
type LuaState struct {
	stack *LuaStack
}

// NewLuaState new a LuaState
func NewLuaState() *LuaState {
	return &LuaState{
		stack: newLuaStack(20), // TODO
	}
}

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

// PushValue pushs the element with idx
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

// ------------------------------------
//    access methods(stack -> go)
// ------------------------------------

// TypeName returns the name of tp
func (s *LuaState) TypeName(t api.LuaType) string {
	switch t {
	case api.LuaTNone:
		return "no value"
	case api.LuaTNil:
		return "nil"
	case api.LuaTBoolean:
		return "boolean"
	case api.LuaTLightUserData:
		return "userdata"
	case api.LuaTNumber:
		return "number"
	case api.LuaTString:
		return "string"
	case api.LuaTTable:
		return "table"
	case api.LuaTFunction:
		return "function"
	case api.LuaTUserData:
		return "userdata"
	case api.LuaTThread:
		return "thread"
	default:
		panic("Type Error")
	}
}

// Type returns the type
func (s *LuaState) Type(idx int) api.LuaType {
	if !s.stack.isValid(idx) {
		return api.LuaTNone
	}

	val := s.stack.get(idx)
	return typeOf(val)
}

// IsNone returns if it is none
func (s *LuaState) IsNone(idx int) bool {
	return s.Type(idx) == api.LuaTNone
}

// IsNil returns if it is nil
func (s *LuaState) IsNil(idx int) bool {
	return s.Type(idx) == api.LuaTNil
}

// IsNoneOrNil returns if it is none nil
func (s *LuaState) IsNoneOrNil(idx int) bool {
	return s.Type(idx) <= api.LuaTNil
}

// IsBoolean returns if it is boolean
func (s *LuaState) IsBoolean(idx int) bool {
	return s.Type(idx) == api.LuaTBoolean
}

// IsString returns if it is string or number
func (s *LuaState) IsString(idx int) bool {
	t := s.Type(idx)
	return t == api.LuaTString || t == api.LuaTNumber // convertable to string
}

// IsNumber returns if it is number or can be converted to number
func (s *LuaState) IsNumber(idx int) bool {
	_, ok := s.ToNumberX(idx)
	return ok
}

// IsInteger returns if it is integer
func (s *LuaState) IsInteger(idx int) bool {
	val := s.stack.get(idx)
	_, ok := val.(int64) // BUG, string value is converted to integer
	return ok
}

// ToBoolean returns boolean value, nil and false are false, other are true
func (s *LuaState) ToBoolean(idx int) bool {
	val := s.stack.get(idx)
	return convertToBoolean(val)
}

// ToNumberX returns number, ok
// TODO
func (s *LuaState) ToNumberX(idx int) (float64, bool) {
	val := s.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
		// TODO
	default:
		return 0, false
	}
}

// ToNumber returns number when successful, returns 0.0 when fail
func (s *LuaState) ToNumber(idx int) float64 {
	n, _ := s.ToNumberX(idx)
	return n
}

// ToIntegerX returns int value, ok
func (s *LuaState) ToIntegerX(idx int) (int64, bool) {
	val := s.stack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

// ToInteger returns int value when successful, return 0 when fails
func (s *LuaState) ToInteger(idx int) int64 {
	i, _ := s.ToIntegerX(idx)
	return i
}

// ToStringX returns string, ok
func (s *LuaState) ToStringX(idx int) (string, bool) {
	val := s.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		str := fmt.Sprintf("%v", x)
		s.stack.set(idx, str) // NOTE: modify stack
		return str, true
	default:
		return "", false
	}
}

// ToString returns string value when successful, return "" when fails
func (s *LuaState) ToString(idx int) string {
	str, _ := s.ToStringX(idx)
	return str
}

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
