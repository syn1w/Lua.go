package state

import "luago/api"

// PushGlobalTable pushes the global table into stack
func (s *LuaState) PushGlobalTable() {
	global := s.registry.get(api.LuaRidxGlobals)
	s.stack.push(global)
}

// GetGlobal pushes value and returns typeof(value)
// for name(key) in global table
func (s *LuaState) GetGlobal(name string) api.LuaType {
	t := s.registry.get(api.LuaRidxGlobals)
	return s.getTable(t, name, false)
}

// SetGlobal sets globalTable[name] by stack.top
func (s *LuaState) SetGlobal(name string) {
	t := s.registry.get(api.LuaRidxGlobals)
	val := s.stack.pop()
	s.setTable(t, name, val, false)
}

// Register set globalTable[name] = f
func (s *LuaState) Register(name string, f api.GoFunction) {
	s.PushGoFunction(f)
	s.SetGlobal(name)
}
