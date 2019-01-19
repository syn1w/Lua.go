package state

// Next pops a key, and pushes the key-value pair from the table at given index
func (s *LuaState) Next(idx int) bool {
	tbVal := s.stack.get(idx)
	if tb, ok := tbVal.(*LuaTable); ok {
		key := s.stack.pop()
		if nextKey := tb.nextKey(key); nextKey != nil {
			s.stack.push(nextKey)
			s.stack.push(tb.get(nextKey))
			return true
		}
		return false
	}
	panic("api_misc: Next: table expected")
}
