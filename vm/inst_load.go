package vm

import "vczn/luago/api"

// loadnil A B | R(A), R(A+1), ..., R(A+B) := nil
func loadNil(i Instruction, vm api.ILuaVM) {
	a, b, _ := i.ABC()
	a++
	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

// loadbool A B C | R(A) := (bool)B; if (C) pc++
func loadBool(i Instruction, vm api.ILuaVM) {
	a, b, c := i.ABC()
	a++
	vm.PushBoolean(b != 0)
	vm.Replace(a)
	if c != 0 {
		vm.AddPC(1)
	}
}

// loadk A, Bx | R(A) := Kst(Bx)
func loadk(i Instruction, vm api.ILuaVM) {
	a, bx := i.ABx()
	a++
	vm.GetConst(bx)
	vm.Replace(a)
}

// loadkx A, R(A) := Kst(extra arg)
func loadkx(i Instruction, vm api.ILuaVM) {
	a, _ := i.ABx()
	a++
	ax := Instruction(vm.Fetch()).Ax()
	vm.GetConst(ax)
	vm.Replace(a)
}
