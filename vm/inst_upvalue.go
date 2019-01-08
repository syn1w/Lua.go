package vm

import (
	"vczn/luago/api"
)

// gettabup A B C | R(A) := UpValue[B][RK(C)]
func getTableUp(inst Instruction, vm api.ILuaVM) {
	a, _, c := inst.ABC()
	a++

	vm.PushGlobalTable()
	vm.GetRK(c)
	vm.GetTable(-2) // -2 is index of table
	vm.Replace(a)
	vm.Pop(1) // pop _G
}
