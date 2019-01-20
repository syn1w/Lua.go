package state

import (
	"vczn/luago/api"
)

func (s *LuaState) Error() int {
	err := s.stack.pop()
	panic(err)
}

// PCall catches the error and rethrows the error if there is an error
// pushes the error status code
func (s *LuaState) PCall(nArgs, nResults, msgh int) (status int) {
	caller := s.stack
	status = api.LuaErrRun
	// catch error
	defer func() {
		if err := recover(); err != nil {
			for s.stack != caller {
				s.popLuaStack()
			}
			s.stack.push(err)
		}
	}()

	s.Call(nArgs, nResults)
	status = api.LuaOk
	return
}
