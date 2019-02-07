package codegen

import (
	"vczn/luago/compiler/ast"
	"vczn/luago/compiler/lexer"
	"vczn/luago/vm"
)

type funcInfo struct {
	parent     *funcInfo
	subFuncs   []*funcInfo
	usedRegs   int
	maxRegs    int
	scopeDepth int                  // start from 0
	constants  map[interface{}]int  // constants table
	upvalues   map[string]upvalInfo // upvalue table
	locVars    []*locVarInfo
	locNames   map[string]*locVarInfo
	labels     map[string]labelInfo
	breaks     [][]int // break table, [depth][pc]
	gotos      []*gotoInfo
	insts      []uint32
	lines      []uint32
	firstLine  int
	lastLine   int
	numParams  int
	isVararg   bool
}

type locVarInfo struct {
	prev       *locVarInfo // linked list
	name       string
	scopeDepth int // start from 0
	slot       int // index of binding local variable
	startPC    int
	endPC      int
	captured   bool
}

type labelInfo struct {
	line       int
	pc         int
	scopeDepth int // start from 0
}

type gotoInfo struct {
	jmpPC      int
	scopeDepth int
	label      string
	pending    bool
}

type upvalInfo struct {
	locVarSlot int // 如果捕获的是直接外围函数的局部变量，记录该局部变量占用的寄存器 slot
	upvalIndex int // 如果是已经被直接外围函数捕获的变量，记录该 upvalue 在直接外围函数 upvalue 表中的 index
	index      int // index of capture
}

func newFuncInfo(parent *funcInfo, fd *ast.FuncDefExp) *funcInfo {
	return &funcInfo{
		parent:    parent,
		subFuncs:  []*funcInfo{},
		constants: make(map[interface{}]int),
		upvalues:  make(map[string]upvalInfo),
		locVars:   make([]*locVarInfo, 0, 8),
		locNames:  make(map[string]*locVarInfo),
		labels:    make(map[string]labelInfo),
		breaks:    make([][]int, 1),
		gotos:     nil,
		insts:     make([]uint32, 0, 8),
		lines:     make([]uint32, 0, 8),
		firstLine: fd.FirstLine,
		lastLine:  fd.LastLine,
		isVararg:  fd.IsVararg,
		numParams: len(fd.ParList),
	}
}

func (fi *funcInfo) indexOfConstant(k interface{}) int {
	if idx, found := fi.constants[k]; found {
		return idx
	}
	idx := len(fi.constants)
	fi.constants[k] = idx
	return idx
}

// allocate a register and return the index of the register
func (fi *funcInfo) allocReg() int {
	fi.usedRegs++
	if fi.usedRegs >= 255 {
		panic("function or expression needs too many registers")
	}

	if fi.usedRegs > fi.maxRegs {
		fi.maxRegs = fi.usedRegs
	}

	return fi.usedRegs - 1
}

// free a register
func (fi *funcInfo) freeReg() {
	if fi.usedRegs <= 0 {
		panic("usedRegs <= 0 !")
	}
	fi.usedRegs--
}

// allocate n registers and return the index of the first register
func (fi *funcInfo) allocRegs(n int) int {
	if n <= 0 {
		panic("allocate registers n <= 0!")
	}

	for i := 0; i < n; i++ {
		fi.allocReg()
	}
	return fi.usedRegs - n
}

// free n registers
func (fi *funcInfo) freeRegs(n int) {
	if n < 0 {
		panic("free registers n < 0!")
	}

	for i := 0; i < n; i++ {
		fi.freeReg()
	}
}

func (fi *funcInfo) enterScope(canBreak bool) {
	// looping block can break
	// includes for, repeat and while
	fi.scopeDepth++
	if canBreak {
		fi.breaks = append(fi.breaks, []int{}) // looping block
	} else {
		fi.breaks = append(fi.breaks, nil)
	}
}

func (fi *funcInfo) exitScope() {
	pendingBreakJmps := fi.breaks[len(fi.breaks)-1]
	fi.breaks = fi.breaks[:len(fi.breaks)-1]

	a := fi.getJmpArgA()
	for _, pc := range pendingBreakJmps {
		sBx := fi.pc() - pc
		inst := (sBx+vm.MaxArgSBx)<<14 | a<<6 | vm.OpJMP
		fi.insts[pc] = uint32(inst)
	}

	fi.fixGotoJmps()

	fi.scopeDepth--
	for _, locVar := range fi.locNames {
		if locVar.scopeDepth > fi.scopeDepth {
			fi.removeLocVar(locVar)
		}
	}
}

func (fi *funcInfo) getJmpArgA() int {
	// TODO
	return 0
}

func (fi *funcInfo) fixGotoJmps() {
	// TODO
}

func (fi *funcInfo) addBreakJmp(pc int) {
	for i := fi.scopeDepth; i >= 0; i-- {
		if fi.breaks[i] != nil { // looping block
			fi.breaks[i] = append(fi.breaks[i], pc)
			return
		}
	}
	panic("<break> at line ? not inside a loop!")
}

// add a local variable and return the index of the variable
func (fi *funcInfo) addLocVar(name string, startPC int) int {
	locVar := &locVarInfo{
		name:       name,
		prev:       fi.locNames[name],
		scopeDepth: fi.scopeDepth,
		slot:       fi.allocReg(),
		startPC:    startPC,
		endPC:      0,
	}

	fi.locVars = append(fi.locVars, locVar)
	fi.locNames[name] = locVar
	return locVar.slot
}

func (fi *funcInfo) removeLocVar(locVar *locVarInfo) {
	fi.freeReg()
	if locVar.prev == nil {
		delete(fi.locNames, locVar.name)
	} else if locVar.prev.scopeDepth == locVar.scopeDepth {
		fi.removeLocVar(locVar.prev)
	} else {
		fi.locNames[locVar.name] = locVar.prev
	}
}

// return the slot of local variable if the local variable bind to a register
// else return -1
func (fi *funcInfo) slotOfLocVar(name string) int {
	if locVar, found := fi.locNames[name]; found {
		return locVar.slot
	}
	return -1
}

func (fi *funcInfo) indexOfUpvalue(name string) int {
	if upval, found := fi.upvalues[name]; found {
		return upval.index
	}

	// try to bind
	if fi.parent != nil {
		if locVar, found := fi.parent.locNames[name]; found {
			idx := len(fi.upvalues)
			fi.upvalues[name] = upvalInfo{locVarSlot: locVar.slot, upvalIndex: -1, index: idx}
			locVar.captured = true
			return idx
		}

		if uvIdx := fi.parent.indexOfUpvalue(name); uvIdx >= 0 {
			idx := len(fi.upvalues)
			fi.upvalues[name] = upvalInfo{locVarSlot: -1, upvalIndex: uvIdx, index: idx}
			return idx
		}
	}

	// binding failed
	return -1
}

func (fi *funcInfo) emitABCInst(line, opcode, a, b, c int) {
	inst := b<<23 | c<<14 | c<<6 | opcode
	fi.insts = append(fi.insts, uint32(inst))
	fi.lines = append(fi.lines, uint32(line))
}

func (fi *funcInfo) emitABxInst(line, opcode, a, bx int) {
	inst := bx<<14 | a<<6 | opcode
	fi.insts = append(fi.insts, uint32(inst))
	fi.lines = append(fi.lines, uint32(line))
}

func (fi *funcInfo) emitAsBxInst(line, opcode, a, sBx int) {
	inst := (sBx+vm.MaxArgSBx)<<14 | a<<6 | opcode
	fi.insts = append(fi.insts, uint32(inst))
	fi.lines = append(fi.lines, uint32(line))
}

func (fi *funcInfo) emitAxInst(line, opcode, ax int) {
	inst := ax<<6 | opcode
	fi.insts = append(fi.insts, uint32(inst))
	fi.lines = append(fi.lines, uint32(line))
}

// r[a] = r[b]
func (fi *funcInfo) emitMove(line, a, b int) {
	fi.emitABCInst(line, vm.OpMOVE, a, b, 0)
}

// r[a] = kst[bx]
func (fi *funcInfo) emitLoadk(line, a int, k interface{}) {
	idx := fi.indexOfConstant(k)
	if idx < (1 << 18) {
		fi.emitABxInst(line, vm.OpLOADK, a, idx)
	} else {
		fi.emitABxInst(line, vm.OpLOADKX, a, 0)
		fi.emitAxInst(line, vm.OpEXTRAARG, idx)
	}
}

// r[a] = (bool)b; if (c) pc++
func (fi *funcInfo) emitLoadBool(line, a, b, c int) {
	fi.emitABCInst(line, vm.OpLOADBOOL, a, b, c)
}

// r[a], r[a+1], ..., r[a+b] = nil => b == n-1
func (fi *funcInfo) emitLoadNil(line, a, n int) {
	fi.emitABCInst(line, vm.OpLOADNIL, a, n-1, 0)
}

// r(a) = upvalue[b]
func (fi *funcInfo) emitGetUpval(line, a, b int) {
	fi.emitABCInst(line, vm.OpGETUPVAL, a, b, 0)
}

// upval[b] = r[a]
func (fi *funcInfo) emitSetUpval(line, a, b int) {
	fi.emitABCInst(line, vm.OpSETUPVAL, a, b, 0)
}

// r[a] = upval[b][rk(c)]
func (fi *funcInfo) emitGetTabUp(line, a, b, c int) {
	fi.emitABCInst(line, vm.OpGETTABUP, a, b, c)
}

// upval[a][rk(b)] = rk(c)
func (fi *funcInfo) emitSetTabUp(line, a, b, c int) {
	fi.emitABCInst(line, vm.OpSETTABUP, a, b, c)
}

// r[a] = r[b][rk(c)]
func (fi *funcInfo) emitGetTable(line, a, b, c int) {
	fi.emitABCInst(line, vm.OpGETTABLE, a, b, c)
}

// r[a][rk(b)] = rk(c)
func (fi *funcInfo) emitSetTable(line, a, b, c int) {
	fi.emitABCInst(line, vm.OpSETTABLE, a, b, c)
}

// r(a) = {}, (arr size = b, map size c)
func (fi *funcInfo) emitNewTable(line, a, b, c int) {
	fi.emitABCInst(line, vm.OpNEWTABLE, a, b, c)
}

// r[a+1] := r[b]; r[a] := r[b][rk(c)]
func (fi *funcInfo) emitSelf(line, a, b, c int) {
	fi.emitABCInst(line, vm.OpSELF, a, b, c)
}

var tokenToOpcode = map[int]int{
	lexer.TokenOpAdd:  vm.OpADD,
	lexer.TokenOpSub:  vm.OpSUB,
	lexer.TokenOpMul:  vm.OpMUL,
	lexer.TokenOpMod:  vm.OpMOD,
	lexer.TokenOpPow:  vm.OpPOW,
	lexer.TokenOpDiv:  vm.OpDIV,
	lexer.TokenOpIDiv: vm.OpIDIV,
	lexer.TokenOpBand: vm.OpBAND,
	lexer.TokenOpBor:  vm.OpBOR,
	lexer.TokenOpBxor: vm.OpBXOR,
	lexer.TokenOpShl:  vm.OpSHL,
	lexer.TokenOpShr:  vm.OpSHR,
}

func (fi *funcInfo) emitUnaryOp(line, op, a, b int) {
	switch op {
	case lexer.TokenOpNot:
		fi.emitABCInst(line, vm.OpNOT, a, b, 0)
	case lexer.TokenOpBnot:
		fi.emitABCInst(line, vm.OpBNOT, a, b, 0)
	case lexer.TokenOpLen:
		fi.emitABCInst(line, vm.OpLEN, a, b, 0)
	case lexer.TokenOpUnm:
		fi.emitABCInst(line, vm.OpUNM, a, b, 0)
	}
}

// r(a) = rk(b) op rk(c)
func (fi *funcInfo) emitBinaryOp(line, op, a, b, c int) {
	if opcode, found := tokenToOpcode[op]; found {
		fi.emitABCInst(line, opcode, a, b, c)
		return // arithmetic and bitwise
	}

	// compare operator
	// if ((rk(b) op rk(c)) != a) pc++
	switch op {
	case lexer.TokenOpEq:
		fi.emitABCInst(line, vm.OpEQ, 1, b, c)
	case lexer.TokenOpNe:
		fi.emitABCInst(line, vm.OpEQ, 0, b, c)
	case lexer.TokenOpLt:
		fi.emitABCInst(line, vm.OpLT, 1, b, c)
	case lexer.TokenOpLe:
		fi.emitABCInst(line, vm.OpLE, 1, b, c)
	case lexer.TokenOpGt:
		fi.emitABCInst(line, vm.OpLT, 1, c, b) // c < b <==> b > c
	case lexer.TokenOpGe:
		fi.emitABCInst(line, vm.OpLE, 1, c, b) // c <= b <==> b >= c
	}

	fi.emitJmp(line, 0, 1)         // ----------
	fi.emitLoadBool(line, a, 0, 1) // ------   |
	fi.emitLoadBool(line, a, 1, 0) // <----|----
	// ...                         // <-----
}

// pc+=sBx; if(a) close all upvalues >= r[a-1]
func (fi *funcInfo) emitJmp(line, a, sBx int) int {
	fi.emitAsBxInst(line, vm.OpJMP, a, sBx)
	return len(fi.insts) - 1
}

// if ((r(a) != c)) pc++
func (fi *funcInfo) emitTest(line, a, c int) {
	fi.emitABCInst(line, vm.OpTEST, a, 0, c)
}

// if ((r(b) == c)) r(a) = r(b) else pc++. c == 0 is and, c == 1 is or
func (fi *funcInfo) emitTestSet(line, a, b, c int) {
	fi.emitABCInst(line, vm.OpTESTSET, a, b, c)
}

// Concat

// r[a], ..., r[a+c-2] = r[a](r[a+1], ..., r[a+b-1])
func (fi *funcInfo) emitCall(line, a, nArgs, nRets int) {
	fi.emitABCInst(line, vm.OpCALL, a, nArgs+1, nRets+1)
}

// return r[a](r[a+1], ... ,r[a+b-1])
func (fi *funcInfo) emitTailCall(line, a, nArgs int) {
	fi.emitABCInst(line, vm.OpTAILCALL, a, nArgs+1, 0)
}

// return r(A), r(a+1), ..., r(a+b-2) => b == n+1
func (fi *funcInfo) emitReturn(line, a, n int) {
	fi.emitABCInst(line, vm.OpRETURN, a, n+1, 0)
}

func (fi *funcInfo) emitForPrep(line, a, sBx int) int {
	fi.emitAsBxInst(line, vm.OpFORPREP, a, sBx)
	return len(fi.insts) - 1
}

// r(a)-=r(a+2); pc+=sBx
func (fi *funcInfo) emitForLoop(line, a, sBx int) int {
	fi.emitAsBxInst(line, vm.OpFORLOOP, a, sBx)
	return len(fi.insts) - 1
}

func (fi *funcInfo) emitTForCall(line, a, c int) {
	fi.emitABCInst(line, vm.OpTFORCALL, a, 0, c)
}

func (fi *funcInfo) emitTForLoop(line, a, sBx int) {
	fi.emitAsBxInst(line, vm.OpTFORLOOP, a, sBx)
}

func (fi *funcInfo) emitSetList(line, a, b, c int) {
	fi.emitABCInst(line, vm.OpSETLIST, a, b, c)
}

// r[a] = closure(proto[bx])
func (fi *funcInfo) emitClosure(line, a, bx int) {
	fi.emitABxInst(line, vm.OpCLOSURE, a, bx)
}

// r[a], r[a+1], ..., r[a+b-2] = vararg
// b-1 == n ==> b == n+1
func (fi *funcInfo) emitVararg(line, a, n int) {
	fi.emitABCInst(line, vm.OpVARARG, a, n+1, 0)
}

func (fi *funcInfo) pc() int {
	return len(fi.insts) - 1
}

func (fi *funcInfo) fixSbx(pc, bx int) {
	inst := fi.insts[pc]
	inst = (inst << 18) >> 18                 // clear sBx
	inst = inst | uint32(bx+vm.MaxArgSBx)<<14 // reset bx
	fi.insts[pc] = inst
}
