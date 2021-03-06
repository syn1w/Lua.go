package vm

import "luago/api"

// forloop A sBx | for numerical
// for index, limit, step ...
// R(A)+=R(A+2)
// if R(A) <= R(A+1) { pc+=sBx; R(A+3)=R(A) }
// R(A) index, R(A+1) limit, R(A+2) step, R(A+3) i
func luaForLoop(inst Instruction, vm api.ILuaVM) {
	a, sBx := inst.AsBx()
	a++
	// R(A)+=R(A+2)
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arithmetic(api.LuaOpAdd)
	vm.Replace(a)
	// if R(A) <= R(A+1) { pc+=sBx; R(A+3)=R(A) }
	isPositiveStep := vm.ToNumber(a+2) >= 0
	if isPositiveStep && vm.Compare(a, a+1, api.LuaOpLe) ||
		!isPositiveStep && vm.Compare(a+1, a, api.LuaOpLe) {
		vm.AddPC(sBx)
		vm.Copy(a, a+3)
	}
}

// forprep A sBx
// R(A)-=R(A+2); PC+=sBx
func luaForPrep(inst Instruction, vm api.ILuaVM) {
	a, sBx := inst.AsBx()
	a++
	// R(A)-=R(A+2) <=> index -= step
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arithmetic(api.LuaOpSub)
	vm.Replace(a)
	// pc+=sBx
	vm.AddPC(sBx)
}

// tforloop A sBx |
// if R(A+1) != nil then {
//	   R(A)=R(A+1); pc += sBx
// }
// for each, R(A) iterator function, R(A+1) state, R(A+2) enumeration index, R(A+3), ...
// loop variables onwards
// update _var -> nextkey
func luaTForLoop(inst Instruction, vm api.ILuaVM) {
	a, sBx := inst.AsBx()
	a++
	if !vm.IsNil(a + 1) {
		vm.Copy(a+1, a)
		vm.AddPC(sBx)
	}
}
