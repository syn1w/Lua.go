package state

import (
	"math"
	"vczn/luago/number"
)

// LuaTable { key = value }, key is not nil or NaN
type LuaTable struct {
	arr []LuaValue
	m   map[LuaValue]LuaValue
}

// NewLuaTable new a LuaTable
func NewLuaTable(nArr, nRecord int) *LuaTable {
	t := &LuaTable{}
	if nArr > 0 {
		t.arr = make([]LuaValue, 0, nArr)
	}

	if nRecord > 0 {
		t.m = make(map[LuaValue]LuaValue, nRecord)
	}

	return t
}

func floatToInteger(key LuaValue) LuaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (t *LuaTable) get(key LuaValue) LuaValue {
	key = floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(t.arr)) {
			return t.arr[idx-1]
		}
	}

	return t.m[key]
}

func (t *LuaTable) put(key, val LuaValue) {
	if key == nil {
		panic("table index is nil")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN")
	}

	key = floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(t.arr))
		if idx <= arrLen {
			t.arr[idx-1] = val
			if idx == arrLen && val == nil {
				t.shrinkArray()
			}
			return
		}

		if idx == arrLen+1 {
			delete(t.m, key)
			if val != nil {
				t.arr = append(t.arr, val)
				t.expandArray()
			}
			return
		}

		if val != nil {
			if t.m == nil {
				t.m = make(map[LuaValue]LuaValue)
				t.m[key] = val
			}
		} else {
			delete(t.m, key)
		}
	}
}

func (t *LuaTable) shrinkArray() {
	var i int
	for i = len(t.arr); i > 0; i-- {
		if t.arr[i-1] != nil {
			break
		}
	}
	t.arr = t.arr[0:i]
}

func (t *LuaTable) expandArray() {
	for idx := int64(len(t.arr)) + 1; true; idx++ {
		if val, found := t.m[idx]; found {
			delete(t.m, idx)
			t.arr = append(t.arr, val)
		} else {
			break
		}
	}
}

func (t *LuaTable) len() int {
	return len(t.arr)
}
