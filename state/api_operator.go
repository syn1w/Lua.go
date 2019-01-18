package state

import (
	"math"
	"vczn/luago/api"
	"vczn/luago/number"
)

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
	metaMethod  string
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var operators = []operator{
	operator{"__add" /* */, iadd, fadd},
	operator{"__sub" /* */, isub, fsub},
	operator{"__mul" /* */, imul, fmul},
	operator{"__mod" /* */, imod, fmod},
	operator{"__pow" /* */, nil, pow},
	operator{"__div" /* */, nil, div},
	operator{"__idiv" /**/, iidiv, fidiv},
	operator{"__band" /**/, band, nil},
	operator{"__bor" /* */, bor, nil},
	operator{"__nxor" /**/, bxor, nil},
	operator{"__shl" /* */, shl, nil},
	operator{"__shr" /* */, shr, nil},
	operator{"__unm" /* */, iunm, funm},
	operator{"__bnot" /**/, bnot, nil},
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
		return
	}
	mm := oper.metaMethod
	if result, ok := callMetaMethod(a, b, mm, s); ok {
		s.stack.push(result)
		return
	}

	panic("arithmetic error")
}

func callMetaMethod(a, b LuaValue, mmName string, state *LuaState) (LuaValue, bool) {
	var mm LuaValue
	if mm = getMetaField(a, mmName, state); mm == nil {
		if mm = getMetaField(b, mmName, state); mm == nil {
			return nil, false
		}
	}

	state.stack.check(4)
	state.stack.push(mm)
	state.stack.push(a)
	state.stack.push(b)
	state.Call(2, 1)
	return state.stack.pop(), true
}

func getMetaField(val LuaValue, fieldName string, state *LuaState) LuaValue {
	if mt := getMetaTable(val, state); mt != nil {
		return mt.get(fieldName)
	}
	return nil
}

// ------------------------------------
//        compare methods
// ------------------------------------

func (s *LuaState) luaEqual(a, b LuaValue) bool {
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
	case *LuaTable:
		if y, ok := b.(*LuaTable); ok && x != y && s != nil {
			if result, ok := callMetaMethod(x, y, "__eq", s); ok {
				return convertToBoolean(result)
			}
		}
		return a == b
	default:
		return a == b
	}
}

func (s *LuaState) luaLt(a, b LuaValue) bool {
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

	if result, ok := callMetaMethod(a, b, "__lt", s); ok {
		return convertToBoolean(result)
	}
	panic("comparison error")
}

// why not: not(b < a), such as NaN
func (s *LuaState) luaLe(a, b LuaValue) bool {
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

	if result, ok := callMetaMethod(a, b, "__le", s); ok {
		return convertToBoolean(result)
	}

	// __le is not found, call !__lt(b, a)
	if result, ok := callMetaMethod(b, a, "__lt", s); ok {
		return !convertToBoolean(result)
	}

	panic("comparison error")
}

// Compare stack.get(idx1) op stack.get(idx2)
func (s *LuaState) Compare(idx1, idx2 int, op api.CompareOp) bool {
	a := s.stack.get(idx1)
	b := s.stack.get(idx2)
	switch op {
	case api.LuaOpEq:
		return s.luaEqual(a, b)
	case api.LuaOpLt:
		return s.luaLt(a, b)
	case api.LuaOpLe:
		return s.luaLe(a, b)
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
	} else if t, ok := val.(*LuaTable); ok {
		s.stack.push(int64(t.len()))
	} else if result, ok := callMetaMethod(val, val, "__len", s); ok {
		s.stack.push(result)
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

			b := s.stack.pop()
			a := s.stack.pop()
			if result, ok := callMetaMethod(a, b, "__concat", s); ok {
				s.stack.push(result)
				continue
			}

			panic("concatenation error")
		}
	}
	// n == 1 do nothing
}

// RawLen returns the raw length of the value at stack[idx]
func (s *LuaState) RawLen(idx int) uint {
	val := s.stack.get(idx)
	switch x := val.(type) {
	case string:
		return uint(len(x))
	case *LuaTable:
		return uint(x.len())
	default:
		return 0
	}
}

// RawEqual returns two values are primitively equal(that is without calling the `__eq`)
func (s *LuaState) RawEqual(idx1, idx2 int) bool {
	if !s.stack.isValid(idx1) || !s.stack.isValid(idx2) {
		return false
	}

	val1, val2 := s.stack.get(idx1), s.stack.get(idx2)
	return s.luaEqual(val1, val2)
}
