package state

import "vczn/luago/api"

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
