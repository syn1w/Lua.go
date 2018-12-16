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

func (s *LuaState) getTable(t, key LuaValue) api.LuaType {
	if tb, ok := t.(*LuaTable); ok {
		val := tb.get(key)
		s.stack.push(val)
		return typeOf(val)
	}
	panic("not a table")
}

// GetTable pushes the value with key(top) and return type of the value
func (s *LuaState) GetTable(idx int) api.LuaType {
	t := s.stack.get(idx)
	key := s.stack.pop()
	return s.getTable(t, key)
}

// GetField pushes the value with k(param string) and return type of the value
func (s *LuaState) GetField(idx int, k string) api.LuaType {
	t := s.stack.get(idx)
	return s.getTable(t, k)
}

// GetI pushes the value with k(param int64) and return type of the value
func (s *LuaState) GetI(idx int, i int64) api.LuaType {
	t := s.stack.get(idx)
	return s.getTable(t, i)
}

func (s *LuaState) setTable(t, k, v LuaValue) {
	if tb, ok := t.(*LuaTable); ok {
		tb.put(k, v)
		return
	}
	panic("not a table")
}

// SetTable pops the val and pops the key, then puts kv into the table
func (s *LuaState) SetTable(idx int) {
	t := s.stack.get(idx)
	v := s.stack.pop()
	k := s.stack.pop()
	s.setTable(t, k, v)
}

// SetField pops the val, then puts kv into the table
func (s *LuaState) SetField(idx int, k string) {
	t := s.stack.get(idx)
	v := s.stack.pop()
	s.setTable(t, k, v)
}

// SetI pops the val, then puts kv into the table
func (s *LuaState) SetI(idx int, k int64) {
	t := s.stack.get(idx)
	v := s.stack.pop()
	s.setTable(t, k, v)
}
