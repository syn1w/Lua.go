package state

import (
	"fmt"
	"vczn/luago/binchunk"
	"vczn/luago/vm"
)

// Load chunk from binary or text file(compile)
// mode: b(binary), t(text file), bt
// return status code, 0 is ok, 1 is error
func (s *LuaState) Load(chunk []byte, chunkName, mode string) int {
	// TODO: load from text file
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	s.stack.push(c)
	return 0
}

// Call function in stack top
func (s *LuaState) Call(nArgs, nResults int) {
	// push args
	val := s.stack.get(-(nArgs + 1))
	if c, ok := val.(*luaClosure); ok {
		fmt.Printf("call %s(%d, %d)\n", c.proto.Source,
			c.proto.LineDefined, c.proto.LastLineDefined) // debug info
		s.callLuaClosure(nArgs, nResults, c)
	} else {
		panic("no function!")
	}
}

func (s *LuaState) callLuaClosure(nArgs, nResults int, c *luaClosure) {
	// prepare
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams) // fixed parameters
	isVararg := c.proto.IsVararg == 1
	newStack := newLuaStack(nRegs + 20)
	newStack.closure = c

	funcAndArgs := s.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	s.pushLuaStack(newStack)
	s.runLuaClosure()
	s.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		s.stack.check(len(results))
		s.stack.pushN(results, nResults)
	}
}

func (s *LuaState) runLuaClosure() {
	for {
		inst := vm.Instruction(s.Fetch())
		inst.Execute(s)

		if inst.Opcode() == vm.OpRETURN {
			break
		}
	}
}
