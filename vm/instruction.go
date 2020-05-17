package vm

import "luago/api"

// Instruction is Lua instruction code
type Instruction uint32

// bx and sbx max value
const (
	MaxArgBx  = 1<<18 - 1     // 2^18-1 = 262143
	MaxArgSBx = MaxArgBx >> 1 // 262143 >> 1 = 131071
)

// Opcode extracts opcode from instruction
func (inst Instruction) Opcode() int {
	return int(inst & 0x3F)
}

// ABC extracts A, B, C operands from instruction
func (inst Instruction) ABC() (a, b, c int) {
	a = int((inst >> 6) & 0xFF)
	c = int((inst >> 14) & 0x1FF)
	b = int((inst >> 23) & 0x1FF)
	return
}

// ABx extracts A, Bx operands from instruction
func (inst Instruction) ABx() (a, bx int) {
	//  Bx 0               131071              1<<18-1
	// sBx -131071           0               (1<<18-1)>>1
	//     |-----------------|---------------------|

	a = int((inst >> 6) & 0xFF)
	bx = int((inst >> 14))
	return
}

// AsBx extracts A, sBx operands from instruction
func (inst Instruction) AsBx() (a, sBx int) {
	a, bx := inst.ABx()
	sBx = bx - MaxArgSBx
	return
}

// Ax extracts Ax operand from instruction
func (inst Instruction) Ax() (ax int) {
	ax = int(inst >> 6)
	return
}

// Execute an instruction
func (inst Instruction) Execute(vm api.ILuaVM) {
	action := opcodes[inst.Opcode()].action
	if action != nil {
		action(inst, vm)
	} else {
		panic(inst.OpName())
	}
}
