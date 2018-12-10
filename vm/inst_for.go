package vm

import "vczn/luago/api"

// for numerical
// for index, limit, step ...
// R(A)+=R(A+2)
// if R(A) <= R(A+1) { pc+=sBx; R(A+3)=R(A) }
// R(A) index, R(A+1) limit, R(A+2) step, R(A+3) i
func luaForLoop(i Instruction, vm api.ILuaVM) {
	a, sBx := i.AsBx()
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

// R(A)-=R(A+2); PC+=sBx
func luaForPrep(i Instruction, vm api.ILuaVM) {
	a, sBx := i.AsBx()
	a++
	// R(A)-=R(A+2) <=> index -= step
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arithmetic(api.LuaOpSub)
	vm.Replace(a)
	// pc+=sBx
	vm.AddPC(sBx)
}
