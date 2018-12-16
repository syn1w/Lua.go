package vm

import "vczn/luago/api"

// R(A) := {} (size = B, C)
func newTable(i Instruction, vm api.ILuaVM) {
	a, b, c := i.ABC()
	a++
	// b 和 c 最多只有 9 位，最大值是 511，为了防止初始容量不够导致的频繁扩容，使用 Float Point Byte
	// 编码方式，来扩大 b 和 c 表示的值域
	// 编码位 eeeeexxx，当 eeeee == 0 时，那个表示的整数就是 xxx，否则表示 (1xxx) * 2^(eeeee-1)
	vm.CreateTable(Fb2int(b), Fb2int(c))
	vm.Replace(a)
}

// R(A) := R(B)[RK(C)]
func getTable(i Instruction, vm api.ILuaVM) {
	a, b, c := i.ABC()
	a++
	b++
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// R(A)[RK(B)] := RK(C)
func setTable(i Instruction, vm api.ILuaVM) {
	a, b, c := i.ABC()
	a++
	vm.GetRK(b) // key
	vm.GetRK(c) // value
	vm.SetTable(a)
}

const lFieldsPerFlush = 50

// R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
func setList(i Instruction, vm api.ILuaVM) {
	// c 9 bits, 最多直接表示 512 个数组索引，显然不够用，所以分批次来操作，FPF 为 Fields Per Flush
	// 也就是每批次最多处理的元素数量，所以 SETLIST 指令最多可以操作 FPF*512 个元素
	// 如果还有更多元素，使用 EXTRAARG 指令

	a, b, c := i.ABC()
	a++

	if c > 0 {
		c--
	} else {
		c = Instruction(vm.Fetch()).Ax() // ExtraArg
	}

	idx := int64(c * lFieldsPerFlush)
	for i := 1; i <= b; i++ {
		idx++
		vm.PushValue(a + i)
		vm.SetI(a, idx)
	}
}
