package vm

import "luago/api"

// R(A) := RK(B) op RK(C)
func binaryArith(inst Instruction, vm api.ILuaVM, op api.ArithmeticOp) {
	a, b, c := inst.ABC()
	a++
	vm.GetRK(b)
	vm.GetRK(c)
	vm.Arithmetic(op)
	vm.Replace(a)
}

// R(A) := op R(B)
func unaryArith(inst Instruction, vm api.ILuaVM, op api.ArithmeticOp) {
	a, b, _ := inst.ABC()
	a++
	b++
	vm.PushValue(b)
	vm.Arithmetic(op)
	vm.Replace(a)
}

func add(inst Instruction, vm api.ILuaVM) { // +
	binaryArith(inst, vm, api.LuaOpAdd)
}

func sub(inst Instruction, vm api.ILuaVM) { // -
	binaryArith(inst, vm, api.LuaOpSub)
}

func mul(inst Instruction, vm api.ILuaVM) { // *
	binaryArith(inst, vm, api.LuaOpMul)
}

func mod(inst Instruction, vm api.ILuaVM) { // %
	binaryArith(inst, vm, api.LuaOpMod)
}

func pow(inst Instruction, vm api.ILuaVM) { // ^
	binaryArith(inst, vm, api.LuaOpPow)
}

func div(inst Instruction, vm api.ILuaVM) { // /
	binaryArith(inst, vm, api.LuaOpDiv)
}

func idiv(inst Instruction, vm api.ILuaVM) { // //
	binaryArith(inst, vm, api.LuaOpIDiv)
}

func bAnd(inst Instruction, vm api.ILuaVM) { // &
	binaryArith(inst, vm, api.LuaOpBAnd)
}

func bOr(inst Instruction, vm api.ILuaVM) { // |
	binaryArith(inst, vm, api.LuaOpBOr)
}

func bXor(inst Instruction, vm api.ILuaVM) { // ~
	binaryArith(inst, vm, api.LuaOpBXor)
}

func shl(inst Instruction, vm api.ILuaVM) { // <<
	binaryArith(inst, vm, api.LuaOpBShl)
}

func shr(inst Instruction, vm api.ILuaVM) { // >>
	binaryArith(inst, vm, api.LuaOpBShr)
}

func unm(inst Instruction, vm api.ILuaVM) { // -
	unaryArith(inst, vm, api.LuaOpUnm)
}

func bNot(inst Instruction, vm api.ILuaVM) { // ~
	unaryArith(inst, vm, api.LuaOpBNot)
}

// \#
// R(A) := len(R(B))
func luaLen(inst Instruction, vm api.ILuaVM) {
	a, b, _ := inst.ABC()
	a++
	b++
	vm.Len(b)
	vm.Replace(a)
}

// ..
// R(A) := R(B).. ... ..R(C)
func luaConcat(inst Instruction, vm api.ILuaVM) {
	a, b, c := inst.ABC()
	a++
	b++
	c++
	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}
	vm.Concat(n)
	vm.Replace(a)
}

// if (RK(B) op RK(C) ~= A) then pc++
func luaCompare(inst Instruction, vm api.ILuaVM, op api.CompareOp) {
	a, b, c := inst.ABC()
	vm.GetRK(b)
	vm.GetRK(c)
	if vm.Compare(-2, -1, op) != (a != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}

// eq 0 b c, if b == c then pc++
func luaEq(inst Instruction, vm api.ILuaVM) { // ==
	luaCompare(inst, vm, api.LuaOpEq)
}

func luaLt(inst Instruction, vm api.ILuaVM) { // <
	luaCompare(inst, vm, api.LuaOpLt)
}

func luaLe(inst Instruction, vm api.ILuaVM) { // <=
	luaCompare(inst, vm, api.LuaOpLe)
}

// R(A) := not R(B)
func luaNot(inst Instruction, vm api.ILuaVM) {
	a, b, _ := inst.ABC()
	a++
	b++
	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

// if ((R(B) == C)) R(A) := R(B) else pc++. C 0 and, C 1 or
func luaTestset(inst Instruction, vm api.ILuaVM) {
	a, b, c := inst.ABC()
	a++
	b++
	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}

// if not (R(A) == C) then pc++
func luaTest(inst Instruction, vm api.ILuaVM) {
	a, _, c := inst.ABC()
	a++
	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}
