package vm

import (
	"vczn/luago/api"
)

// getUpval A B | R(A) := UpValue[B]
// local a, b
// function foo() return a + b end
// getupval 0  0 ; a
// getupval 1  1 ; b
func getUpval(inst Instruction, vm api.ILuaVM) {
	a, b, _ := inst.ABC()
	a++
	b++
	vm.Copy(api.LuaUpvalueIndex(b), a)
}

// setUpval A B | UpValue[B] := R(A)
// local a, b
// function foo() a, b = 1, 2 end
// setupval  1 1     ; b
// setupval  0 0     ; a
func setUpval(inst Instruction, vm api.ILuaVM) {
	a, b, _ := inst.ABC()
	a++
	b++
	vm.Copy(a, api.LuaUpvalueIndex(b))
}

// getTabUp A B C | R(A) := UpValue[B][RK(C)]
// getupvalue + gettable
func getTabUp(inst Instruction, vm api.ILuaVM) {
	a, b, c := inst.ABC()
	a++
	b++

	vm.GetRK(c)
	vm.GetTable(api.LuaUpvalueIndex(b))
	vm.Replace(a)
}

// setTabUp A B C | UpValue[A][RK(B)] := RK(C)
func setTabUp(inst Instruction, vm api.ILuaVM) {
	a, b, c := inst.ABC()
	a++

	vm.GetRK(b)                         // key
	vm.GetRK(c)                         // value
	vm.SetTable(api.LuaUpvalueIndex(a)) // pop value, key and set table
}
