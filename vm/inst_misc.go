package vm

import "luago/api"

// move A B | R(A) := R(B)
func move(inst Instruction, vm api.ILuaVM) {
	a, b, _ := inst.ABC()
	a++ // registers index begin with 0, stack index begin with 1
	b++
	vm.Copy(b, a)
}

// jmp A, sBx | pc+=sBx; if (A) close all upvalues >= R(A - 1)
func jmp(inst Instruction, vm api.ILuaVM) {
	a, sBx := inst.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		vm.CloseUpvalues(a)
	}
}
