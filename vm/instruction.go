package vm

// Instruction is Lua instruction code
type Instruction uint32

const (
	maxArgBx  = 1<<18 - 1     // 2^18-1 = 262143
	maxArgSBx = maxArgBx >> 1 // 262143 >> 1 = 131071
)

// Opcode extracts opcode from instruction
func (i Instruction) Opcode() int {
	return int(i & 0x3F)
}

// ABC extracts A, B, C operands from instruction
func (i Instruction) ABC() (a, b, c int) {
	a = int((i >> 6) & 0xFF)
	c = int((i >> 14) & 0x1FF)
	b = int((i >> 23) & 0x1FF)
	return
}

// ABx extracts A, Bx operands from instruction
func (i Instruction) ABx() (a, bx int) {
	//  Bx 0               131071              1<<18-1
	// sBx -131071           0               (1<<18-1)<<1
	//     |-----------------|---------------------|

	a = int((i >> 6) & 0xFF)
	bx = int((i >> 14))
	return
}

// AsBx extracts A, sBx operands from instruction
func (i Instruction) AsBx() (a, sBx int) {
	a, bx := i.ABx()
	sBx = bx - maxArgSBx
	return
}

// Ax extracts Ax operand from instruction
func (i Instruction) Ax() (ax int) {
	ax = int(i >> 6)
	return
}
