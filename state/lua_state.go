package state

import (
	"fmt"
	"math"
	"vczn/luago/api"
	"vczn/luago/number"
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
func (s *LuaState) ToNumberX(idx int) (float64, bool) {
	val := s.stack.get(idx)
	return convertToFloat(val)
}

// ToNumber returns number when successful, returns 0.0 when fail
func (s *LuaState) ToNumber(idx int) float64 {
	n, _ := s.ToNumberX(idx)
	return n
}

// ToIntegerX returns int value, ok
func (s *LuaState) ToIntegerX(idx int) (int64, bool) {
	val := s.stack.get(idx)
	return convertToInteger(val)
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

// ------------------------------------
//        arithmetic methods
// ------------------------------------

var (
	iadd  = func(a, b int64) int64 { return a + b }
	fadd  = func(a, b float64) float64 { return a + b }
	isub  = func(a, b int64) int64 { return a - b }
	fsub  = func(a, b float64) float64 { return a - b }
	imul  = func(a, b int64) int64 { return a * b }
	fmul  = func(a, b float64) float64 { return a * b }
	imod  = number.IMod
	fmod  = number.FMod
	pow   = math.Pow
	div   = func(a, b float64) float64 { return a / b }
	iidiv = number.IFloorDiv
	fidiv = number.FFloorDiv
	band  = func(a, b int64) int64 { return a & b }
	bor   = func(a, b int64) int64 { return a | b }
	bxor  = func(a, b int64) int64 { return a ^ b }
	shl   = number.ShiftLeft
	shr   = number.ShiftRight
	iunm  = func(a, _ int64) int64 { return -a }
	funm  = func(a, _ float64) float64 { return -a }
	bnot  = func(a, _ int64) int64 { return ^a }
)

type operator struct {
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var operators = []operator{
	operator{iadd, fadd},
	operator{isub, fsub},
	operator{imul, fmul},
	operator{imod, fmod},
	operator{nil, pow},
	operator{nil, div},
	operator{iidiv, fidiv},
	operator{band, nil},
	operator{bor, nil},
	operator{bxor, nil},
	operator{shl, nil},
	operator{shr, nil},
	operator{iunm, funm},
	operator{bnot, nil},
}

func luaArith(a, b LuaValue, op operator) LuaValue {
	if op.floatFunc == nil { // bitwise
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else { // arithmetic
		if op.integerFunc != nil { // add, sub, mul, mod, idiv, unm
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}

	return nil
}

// Arithmetic pop one or two operands and push the result
func (s *LuaState) Arithmetic(op api.ArithmeticOp) {
	var a, b LuaValue
	b = s.stack.pop()
	if op != api.LuaOpUnm && op != api.LuaOpBNot {
		a = s.stack.pop()
	} else {
		a = b
	}

	oper := operators[op]
	if result := luaArith(a, b, oper); result != nil {
		s.stack.push(result) // NOTE: modify stack
	} else {
		panic("arithmetic error")
	}
}

// ------------------------------------
//        compare methods
// ------------------------------------

func luaEqual(a, b LuaValue) bool {
	switch x := a.(type) {
	case nil:
		return b == nil // nil == nil, nil ~= other
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x == float64(y)
		case float64:
			return x == y
		default:
			return false
		}
	case string:
		y, ok := b.(string)
		return ok && x == y
	default:
		return a == b
	}
}

func luaLt(a, b LuaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x < float64(y)
		case float64:
			return x < y
		}
	}
	panic("comparison error")
}

// why not not(b < a), such as NaN
func luaLe(a, b LuaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case int64:
			return x <= float64(y)
		case float64:
			return x <= y
		}
	}
	panic("comparison error")
}

// Compare stack.get(idx1) op stack.get(idx2)
func (s *LuaState) Compare(idx1, idx2 int, op api.CompareOp) bool {
	a := s.stack.get(idx1)
	b := s.stack.get(idx2)
	switch op {
	case api.LuaOpEq:
		return luaEqual(a, b)
	case api.LuaOpLt:
		return luaLt(a, b)
	case api.LuaOpLe:
		return luaLe(a, b)
	default:
		panic("invalid compare operator")
	}
}

// ------------------------------------
//             len method
// ------------------------------------

// Len pushes len(stack.get(idx))
// TODO, only string
func (s *LuaState) Len(idx int) {
	val := s.stack.get(idx)
	if str, ok := val.(string); ok {
		s.stack.push(int64(len(str)))
	} else {
		panic("length method error")
	}
}

// ------------------------------------
//           concat method
// ------------------------------------

// Concat pops n elements, concats them, and push the result
func (s *LuaState) Concat(n int) {
	if n == 0 {
		s.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			// TODO: optimization
			if s.IsString(-1) && s.IsString(-2) {
				s2 := s.ToString(-1)
				s1 := s.ToString(-2)
				s.stack.pop()
				s.stack.pop()
				s.stack.push(s1 + s2)
				continue
			}
			panic("concatenation error")
		}
	}
	// n == 1 do nothing
}
