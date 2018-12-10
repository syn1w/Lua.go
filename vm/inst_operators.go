package vm

import "vczn/luago/api"

// R(A) := RK(B) op RK(C)
func binaryArith(i Instruction, vm api.ILuaVM, op api.ArithmeticOp) {
	a, b, c := i.ABC()
	a++
	vm.GetRK(b)
	vm.GetRK(c)
	vm.Arithmetic(op)
	vm.Replace(a)
}

// R(A) := op R(B)
func unaryArith(i Instruction, vm api.ILuaVM, op api.ArithmeticOp) {
	a, b, _ := i.ABC()
	a++
	b++
	vm.PushValue(b)
	vm.Arithmetic(op)
	vm.Replace(a)
}

func add(i Instruction, vm api.ILuaVM) { // +
	binaryArith(i, vm, api.LuaOpAdd)
}

func sub(i Instruction, vm api.ILuaVM) { // -
	binaryArith(i, vm, api.LuaOpSub)
}

func mul(i Instruction, vm api.ILuaVM) { // *
	binaryArith(i, vm, api.LuaOpMul)
}

func mod(i Instruction, vm api.ILuaVM) { // %
	binaryArith(i, vm, api.LuaOpMod)
}

func pow(i Instruction, vm api.ILuaVM) { // ^
	binaryArith(i, vm, api.LuaOpPow)
}

func div(i Instruction, vm api.ILuaVM) { // /
	binaryArith(i, vm, api.LuaOpDiv)
}

func idiv(i Instruction, vm api.ILuaVM) { // //
	binaryArith(i, vm, api.LuaOpIDiv)
}

func bAnd(i Instruction, vm api.ILuaVM) { // &
	binaryArith(i, vm, api.LuaOpBAnd)
}

func bOr(i Instruction, vm api.ILuaVM) { // |
	binaryArith(i, vm, api.LuaOpBOr)
}

func bXor(i Instruction, vm api.ILuaVM) { // ~
	binaryArith(i, vm, api.LuaOpBXor)
}

func shl(i Instruction, vm api.ILuaVM) { // <<
	binaryArith(i, vm, api.LuaOpBShl)
}

func shr(i Instruction, vm api.ILuaVM) { // >>
	binaryArith(i, vm, api.LuaOpBShr)
}

func unm(i Instruction, vm api.ILuaVM) { // -
	unaryArith(i, vm, api.LuaOpUnm)
}

func bNot(i Instruction, vm api.ILuaVM) { // ~
	unaryArith(i, vm, api.LuaOpBNot)
}

// \#
// R(A) := len(R(B))
func luaLen(i Instruction, vm api.ILuaVM) {
	a, b, _ := i.ABC()
	a++
	b++
	vm.Len(b)
	vm.Replace(a)
}

// ..
// R(A) := R(B).. ... ..R(C)
func luaConcat(i Instruction, vm api.ILuaVM) {
	a, b, c := i.ABC()
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
func luaCompare(i Instruction, vm api.ILuaVM, op api.CompareOp) {
	a, b, c := i.ABC()
	vm.GetRK(b)
	vm.GetRK(c)
	if vm.Compare(-2, -1, op) != (a != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}

func luaEq(i Instruction, vm api.ILuaVM) { // ==
	luaCompare(i, vm, api.LuaOpEq)
}

func luaLt(i Instruction, vm api.ILuaVM) { // <
	luaCompare(i, vm, api.LuaOpLt)
}

func luaLe(i Instruction, vm api.ILuaVM) { // <=
	luaCompare(i, vm, api.LuaOpLe)
}

// R(A) := not R(B)
func luaNot(i Instruction, vm api.ILuaVM) {
	a, b, _ := i.ABC()
	a++
	b++
	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

// if ((R(B) == C)) R(A) := R(B) else pc++. C 0 and, C 1 or
func luaTestset(i Instruction, vm api.ILuaVM) {
	a, b, c := i.ABC()
	a++
	b++
	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}

// if not (R(A) == C) then pc++
func luaTest(i Instruction, vm api.ILuaVM) {
	a, _, c := i.ABC()
	a++
	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}
