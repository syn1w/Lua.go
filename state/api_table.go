package state

import "vczn/luago/api"

// NewTable pushes a empty lua table
func (s *LuaState) NewTable() {
	s.CreateTable(0, 0)
}

// CreateTable pushes a lua table with nArr, nRecord
func (s *LuaState) CreateTable(nArr, nRecord int) {
	t := NewLuaTable(nArr, nRecord)
	s.stack.push(t)
}

func (s *LuaState) getTable(t, key LuaValue, raw bool) api.LuaType {
	if tb, ok := t.(*LuaTable); ok {
		val := tb.get(key)

		if raw || val != nil || !tb.hasMetaField("__index") {
			s.stack.push(val)
			return typeOf(val)
		}
	}

	if !raw {
		if mf := getMetaField(t, "__index", s); mf != nil {
			switch x := mf.(type) {
			case *LuaTable:
				return s.getTable(x, key, false)
			case *luaClosure:
				s.stack.push(mf)
				s.stack.push(t)
				s.stack.push(key)
				s.Call(2, 1)
				v := s.stack.get(-1)
				return typeOf(v)
			}
		}
	}

	panic("index error")
}

// GetTable pushes the value with key(top) and return type of the value
func (s *LuaState) GetTable(idx int) api.LuaType {
	t := s.stack.get(idx)
	key := s.stack.pop()
	return s.getTable(t, key, false)
}

// GetField pushes the value with k(param string) and return type of the value
func (s *LuaState) GetField(idx int, k string) api.LuaType {
	t := s.stack.get(idx)
	return s.getTable(t, k, false)
}

// GetI pushes the value with k(param int64) and return type of the value
func (s *LuaState) GetI(idx int, i int64) api.LuaType {
	t := s.stack.get(idx)
	return s.getTable(t, i, false)
}

func (s *LuaState) setTable(t, key, v LuaValue, raw bool) {
	if tb, ok := t.(*LuaTable); ok {
		if raw || tb.get(key) != nil || !tb.hasMetaField("__newindex") {
			tb.put(key, v)
			return
		}
	}

	if !raw {
		if mf := getMetaField(t, "__newindex", s); mf != nil {
			switch x := mf.(type) {
			case *LuaTable:
				s.setTable(x, key, v, false)
				return
			case *luaClosure:
				s.stack.push(mf)
				s.stack.push(t)
				s.stack.push(key)
				s.stack.push(v)
				s.Call(3, 0)
				return
			}
		}
	}
	panic("index error")
}

// SetTable pops the val and pops the key, then puts kv into the table
func (s *LuaState) SetTable(idx int) {
	t := s.stack.get(idx)
	v := s.stack.pop()
	k := s.stack.pop()
	s.setTable(t, k, v, false)
}

// SetField pops the val, then puts kv into the table
func (s *LuaState) SetField(idx int, k string) {
	t := s.stack.get(idx)
	v := s.stack.pop()
	s.setTable(t, k, v, false)
}

// SetI pops the val, then puts kv into the table
func (s *LuaState) SetI(idx int, k int64) {
	t := s.stack.get(idx)
	v := s.stack.pop()
	s.setTable(t, k, v, false)
}

// GetMetaTable pushes stack[idx].metaTable
func (s *LuaState) GetMetaTable(idx int) bool {
	val := s.stack.get(idx)
	if mt := getMetaTable(val, s); mt != nil {
		s.stack.push(mt)
		return true
	}
	return false
}

// RawGet similar to GetTable, but does a raw access (i.e., without metamethods).
func (s *LuaState) RawGet(idx int) api.LuaType {
	t := s.stack.get(idx)
	k := s.stack.pop()
	return s.getTable(t, k, true)
}

// SetMetaTable sets the stack[idx].mt = stack.pop()
func (s *LuaState) SetMetaTable(idx int) {
	val := s.stack.get(idx)
	mtVal := s.stack.pop()
	if mtVal == nil {
		setMetaTable(val, nil, s)
	} else if mt, ok := mtVal.(*LuaTable); ok {
		setMetaTable(val, mt, s)
	} else {
		panic("table expected")
	}
}

// RawSet similar to SetTable, but does a raw assignment (i.e., without metamethods).
func (s *LuaState) RawSet(idx int) {
	t := s.stack.get(idx)
	v := s.stack.pop()
	key := s.stack.pop()
	s.setTable(t, key, v, true)
}

// RawGetI similar to GetI, but without `__index`
func (s *LuaState) RawGetI(idx int, i int64) api.LuaType {
	t := s.stack.get(idx)
	return s.getTable(t, i, true)
}

// RawSetI similar to SetI, but without `__newindex`
func (s *LuaState) RawSetI(idx int, i int64) {
	t := s.stack.get(idx)
	v := s.stack.pop()
	s.setTable(t, i, v, true)
}
