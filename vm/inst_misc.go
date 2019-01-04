package vm

import "vczn/luago/api"

// move A B | R(A) := R(B)
func move(i Instruction, vm api.ILuaVM) {
	a, b, _ := i.ABC()
	a++ // registers index begin with 0, stack index begin with 1
	b++
	vm.Copy(b, a)
}

// jmp A, sBx | pc+=sBx; if (A) close all upvalues >= R(A - 1)
func jmp(i Instruction, vm api.ILuaVM) {
	a, sBx := i.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		panic("TODO")
	}
}
