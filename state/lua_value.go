package state

import (
	"fmt"
	"luago/api"
	"luago/number"
)

// LuaValue is the type of Lua value
type LuaValue interface{}

func typeOf(val LuaValue) api.LuaType {
	switch val.(type) {
	case nil:
		return api.LuaTNil
	case bool:
		return api.LuaTBoolean
	case int64:
		return api.LuaTNumber
	case float64:
		return api.LuaTNumber
	case string:
		return api.LuaTString
	case *LuaTable:
		return api.LuaTTable
	case *luaClosure:
		return api.LuaTFunction
	default:
		panic("TODO")
	}
}

// nil and false return false, other return true
func convertToBoolean(val LuaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

func convertToFloat(val LuaValue) (float64, bool) {
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case string:
		return number.ParseFloat(x)
	default:
		return 0.0, false
	}
}

func stringToInteger(s string) (int64, bool) {
	if i, ok := number.ParseInteger(s); ok {
		return i, true
	}

	if f, ok := number.ParseFloat(s); ok {
		return number.FloatToInteger(f)
	}

	return 0, false
}

func convertToInteger(val LuaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return number.FloatToInteger(x)
	case string:
		return stringToInteger(x)
	default:
		return 0, false
	}
}

func setMetaTable(val LuaValue, mt *LuaTable, state *LuaState) {
	if t, ok := val.(*LuaTable); ok { // val is a table
		t.metatable = mt
		return
	}

	// val is not a table, register into registry table
	key := fmt.Sprintf("_MT%d", typeOf(val))
	state.registry.put(key, val)
}

func getMetaTable(val LuaValue, state *LuaState) *LuaTable {
	if t, ok := val.(*LuaTable); ok {
		return t.metatable
	}

	key := fmt.Sprintf("_MT%d", typeOf(val))
	if mt := state.registry.get(key); mt != nil {
		return mt.(*LuaTable)
	}

	return nil
}
