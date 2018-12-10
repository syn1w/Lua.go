package vm

import "vczn/luago/api"

// move A B | R(A) := R(B)
func move(i Instruction, vm api.ILuaVM) {
	a, b, _ := i.ABC()
	a++ // 寄存器索引从 0 开始，栈索引从 1 开始
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
