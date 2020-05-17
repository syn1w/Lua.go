package vm

import (
	"luago/api"
)

// closure, instantiate a subfuntion prototype as a closure
// and put it into R(A)
// A Bx,  R(A) := closure(KPROTO[Bx])
func closure(inst Instruction, vm api.ILuaVM) {
	a, bx := inst.ABx()
	a++
	vm.LoadProto(bx)
	vm.Replace(a)
}

// call
// A B C, R(A), ..., R(A+C-2) := R(A)(R(A+1), ..., R(A+B-1))
func call(inst Instruction, vm api.ILuaVM) {
	// local a, b, c = f(1,2,3,4)
	// CALL 0 5 4

	a, b, c := inst.ABC()
	a++
	nArgs := pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	popResults(a, c, vm)
}

func pushFuncAndArgs(a, b int, vm api.ILuaVM) int {
	if b >= 1 { // func and b - 1 args
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	}
	// b == 0
	// receive all args for another function
	fixStack(a, vm)
	return vm.GetTop() - vm.RegisterCount() - 1
}

func fixStack(a int, vm api.ILuaVM) {
	// second half of the args are already in stack
	// pushes first half of the args
	// and rotate the stack
	dst := int(vm.ToInteger(-1))
	vm.Pop(1)
	vm.CheckStack(dst - a)
	for i := a; i < dst; i++ {
		vm.PushValue(i)
	}

	vm.Rotate(vm.RegisterCount()+1, dst-a)
}

func popResults(a, c int, vm api.ILuaVM) {
	if c == 1 {
		// no result, do nothing
	} else if c > 1 { // c-1 results
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else { // c <= 0
		// leave the return value in stack
		// push destination register
		// to passing intto another function
		// suck as f(1, 2, g())
		// CALL 3 1 0   ; g, c == 0
		// CALL 0 0 1   ; f, b == 0
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}

// A B,   return R(A), ..., R(A+B-2)
func luaReturn(inst Instruction, vm api.ILuaVM) {
	a, b, _ := inst.ABC()
	a++
	if b == 1 {
		// no return values, do nothing
	} else if b > 1 { // b - 1 return values
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else { // b == 0
		fixStack(a, vm)
	}
}

// A B,   R(A), R(A+1), ..., R(A+B-2) = vararg
func vararg(inst Instruction, vm api.ILuaVM) {
	a, b, _ := inst.ABC()
	a++
	if b != 1 {
		vm.LoadVararg(b - 1)
		popResults(a, b, vm)
	}
}

// tailcall A B C | return R(A)(R(A+1), ..., R(A+B-1))
func tailcall(inst Instruction, vm api.ILuaVM) {
	// return f(a, b, c)
	// TAILCALL 0 4 0
	a, b, _ := inst.ABC()
	a++
	c := 0
	nArgs := pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	popResults(a, c, vm)
}

// self A B C | R(A+1) := R(B); R(A) := R(B)[RK(C)]
func luaSelf(inst Instruction, vm api.ILuaVM) {
	// local a, obj; obj:func(a)
	// SELF 2 1 -1 ; func
	a, b, c := inst.ABC()
	a++
	b++
	vm.Copy(b, a+1)
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// tforcall A C | R(A+3), ..., R(A+2+C) := R(A)(R(A+1), R(A+2))
// call next ==> pushes key and value
func luaTForCall(inst Instruction, vm api.ILuaVM) {
	a, _, c := inst.ABC()
	a++
	pushFuncAndArgs(a, 3, vm) // next/inext
	vm.Call(2, c)             // call next/inext
	popResults(a+3, c+1, vm)  // return k, v
}
