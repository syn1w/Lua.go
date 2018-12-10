package state

import (
	"fmt"
	"vczn/luago/api"
)

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
