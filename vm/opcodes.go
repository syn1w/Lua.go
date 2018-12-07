package vm

// reference:
// https://github.com/zxh0/lua.go
// http://luaforge.net/docman/83/98/ANoFrillsIntroToLua51VMInstructions.pdf
// https://the-ravi-programming-language.readthedocs.io/en/latest/lua_bytecode_reference.html#lua-5-3-bytecode-reference
// https://blog.csdn.net/yuanlin2008/column/info/luainternals
// https://cloudwu.github.io/lua53doc/manual.html

// Lua vm instructions
//       31                                   0
//       +-------------------------------------+
// iABC  |   B:9   |   C:9   |  A:8   |opcode:6|
//       +-------------------------------------+
// iABx  |       Bx:18       |  A:8   |opcode:6|
//       +-------------------------------------+
// iAsBx |      sBx:18       |  A:8   |opcode:6|  // signed Bx
//       +-------------------------------------+
// iAx   |          Ax:26             |opcode:6|
//       +-------------------------------------+

// OpMode
const (
	IABC = iota
	IABx
	IAsBx
	IAx
)

// OpName
const (
	OpMOVE     = iota // Copy a value between registers
	OpLOADK           // Load a constant into a register
	OpLOADKX          // Load a constant into a register
	OpLOADBOOL        // Load a boolean into a register
	OpLOADNIL         // Load nil values into a range of registers
	OpGETUPVAL        // Read an upvalue into a register
	OpGETTABUP        // Read a value from table in up-value into a register
	OpGETTABLE        // Read a table element into a register
	OpSETTABUP        // Write a register value into table in up-value
	OpSETUPVAL        // Write a register value into an upvalue
	OpSETTABLE        // Write a register value into a table element
	OpNEWTABLE        // Create a new table
	OpSELF            // Prepare an object method for calling
	OpADD             // Addition operator
	OpSUB             // Subtraction operator
	OpMUL             // Multiplication operator
	OpMOD             // Modulus (remainder) operator
	OpPOW             // Exponentation operator
	OpDIV             // Division operator
	OpIDIV            // Integer division operator
	OpBAND            // Bit-wise AND operator
	OpBOR             // Bit-wise OR operator
	OpBXOR            // Bit-wise Exclusive OR operator
	OpSHL             // Shift bits left
	OpSHR             // Shift bits right
	OpUNM             // Unary minus
	OpBNOT            // Bit-wise NOT operator
	OpNOT             // Logical NOT operator
	OpLEN             // Length operator
	OpCONCAT          // Concatenate a range of registers
	OpJMP             // Unconditional jump
	OpEQ              // Equality test, with conditional jump
	OpLT              // Less than test, with conditional jump
	OpLE              // Less than or equal to test, with conditional jump
	OpTEST            // Boolean test, with conditional jump
	OpTESTSET         // Boolean test, with conditional jump and assignment
	OpCALL            // Call a closure
	OpTAILCALL        // Perform a tail call
	OpRETURN          // Return from function call
	OpFORLOOP         // Iterate a numeric for loop
	OpFORPREP         // Initialization for a numeric for loop
	OpTFORLOOP        // Iterate a generic for loop
	OpTFORCALL        // Initialization for a generic for loop
	OpSETLIST         // Set a range of array elements for a table
	OpCLOSURE         // Create a closure of a function prototype
	OpVARARG          // Assign vararg function arguments to registers
	OpEXTRAARG        // Extra (larger) argument for previous opcode
)

// OpBMode or OPCMode
const (
	OpArgN = iota // argument is not used
	OpArgU        // argument is used
	OpArgR        // argument is a register or a jump offset
	OpArgK        // argument is a constant or register/constant
)

type opcode struct {
	testFlag byte   // operator is test (next instruction must be a jump)
	setAFlag byte   // instruction set regeister
	argBMode byte   // B arg mode
	argCMode byte   // C arg mode
	opMode   byte   // op mode
	name     string // string
}

var opcodes = []opcode{
	//     T  A  B       C       mode
	opcode{0, 1, OpArgR, OpArgN, IABC /* */, "MOVE    "}, // A B,   R(A) = R(B)
	opcode{0, 1, OpArgK, OpArgN, IABx /* */, "LOADK   "}, // A Bx,  R(A) = Kst(Bx), K is constant
	opcode{0, 1, OpArgN, OpArgN, IABx /* */, "LOADKX  "}, // A, 	  R(A) := Kst(extra arg)
	opcode{0, 1, OpArgU, OpArgU, IABC /* */, "LOADBOOL"}, // A B C, R(A) := (Bool)B; if (C) pc++
	opcode{0, 1, OpArgU, OpArgN, IABC /* */, "LOADNIL "}, // A B,   R(A), R(A+1), ..., R(A+B) := nil
	opcode{0, 1, OpArgU, OpArgN, IABC /* */, "GETUPVAL"}, // A B,   R(A) := UpValue[B]
	opcode{0, 1, OpArgU, OpArgK, IABC /* */, "GETTABUP"}, // A B C, R(A) := UpValue[B][RK(C)]
	opcode{0, 1, OpArgR, OpArgK, IABC /* */, "GETTABLE"}, // A B C, R(A) := R(B)[RK(C)]
	opcode{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABUP"}, // A B C, UpValue[A][RK(B)] := RK(C)
	opcode{0, 0, OpArgU, OpArgN, IABC /* */, "SETUPVAL"}, // A B,   UpValue[B] := R(A)
	opcode{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABLE"}, // A B C, R(A)[RK(B)] := RK(C)
	opcode{0, 1, OpArgU, OpArgU, IABC /* */, "NEWTABLE"}, // A B C, R(A) := {} (size = B,C)
	opcode{0, 1, OpArgR, OpArgK, IABC /* */, "SELF    "}, // A B C, R(A+1) := R(B); R(A) := R(B)[RK(C)]
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "ADD     "}, // A B C, R(A) := RK(B) + RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SUB     "}, // A B C, R(A) := RK(B) - RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "MUL     "}, // A B C, R(A) := RK(B) * RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "MOD     "}, // A B C, R(A) := RK(B) % RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "POW     "}, // A B C, R(A) := RK(B) ^ RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "DIV     "}, // A B C, R(A) := RK(B) / RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "IDIV    "}, // A B C, R(A) := RK(B) // RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BAND    "}, // A B C, R(A) := RK(B) & RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BOR     "}, // A B C, R(A) := RK(B) | RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "BXOR    "}, // A B C, R(A) := RK(B) ^ RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SHL     "}, // A B C, R(A) := RK(B) << RK(C)
	opcode{0, 1, OpArgK, OpArgK, IABC /* */, "SHR     "}, // A B C, R(A) := RK(B) >> RK(C)
	opcode{0, 1, OpArgR, OpArgN, IABC /* */, "UNM     "}, // A B,   R(A) := -R(B)
	opcode{0, 1, OpArgR, OpArgN, IABC /* */, "BNOT    "}, // A B,   R(A) := ~R(B)
	opcode{0, 1, OpArgR, OpArgN, IABC /* */, "NOT     "}, // A B,   R(A) := !R(B)
	opcode{0, 1, OpArgR, OpArgN, IABC /* */, "LEN     "}, // A B,   R(A) := len(R(B))
	opcode{0, 1, OpArgR, OpArgR, IABC /* */, "CONCAT  "}, // A B C, R(A) := R(B).. ... ..R(C)
	opcode{0, 0, OpArgR, OpArgN, IAsBx /**/, "JMP     "}, // A sBx, pc+=sBx; if (A) close all upvalues >= R(A - 1)
	opcode{1, 0, OpArgK, OpArgK, IABC /* */, "EQ      "}, // A B C, if ((RK(B) == RK(C)) != A) pc++
	opcode{1, 0, OpArgK, OpArgK, IABC /* */, "LT      "}, // A B C, if ((RK(B) <  RK(C)) != A) pc++
	opcode{1, 0, OpArgK, OpArgK, IABC /* */, "LE      "}, // A B C, if ((RK(B) <= RK(C)) != A) pc++
	opcode{1, 0, OpArgN, OpArgU, IABC /* */, "TEST    "}, // A C,   if ((R(A) != C)) pc++
	opcode{1, 1, OpArgR, OpArgU, IABC /* */, "TESTSET "}, // A B C, if ((R(B) == C)) R(A) := R(B) else pc++. C 0 and, C 1 or
	opcode{0, 1, OpArgU, OpArgU, IABC /* */, "CALL    "}, // A B C, R(A), ..., R(A+C-2) := R(A)(R(A+1), ..., R(A+B-1))
	opcode{0, 1, OpArgU, OpArgU, IABC /* */, "TAILCALL"}, // A B C, return R(A)(R(A+1), ..., R(A+B-1))
	opcode{0, 0, OpArgU, OpArgN, IABC /* */, "RETURN  "}, // A B,   return R(A), ..., R(A+B-2)
	opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORLOOP "}, // A sBx, R(A)+=R(A+2); if R(A) <= R(A+1) then { pc+=sBx; R(A+3)=R(A) }, R(A) index, R(A+1) limit, R(A+2) step, R(A+3) i
	opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORPREP "}, // A sBx, R(A)-=R(A+2); PC+=sBx
	opcode{0, 0, OpArgN, OpArgU, IABC /* */, "TFORCALL"}, // A C,   R(A+3), ..., R(A+2+C) := R(A)(R(A+1), R(A+2));
	opcode{0, 1, OpArgR, OpArgN, IAsBx /**/, "TFORLOOP"}, // A sBx, if R(A+1) != nil then { R(A)=R(A+1); pc += sBx }  for each, R(A) iterator function, R(A+1) state, R(A+2) enumeration index, R(A+3), ... loop variables onwards
	opcode{0, 0, OpArgU, OpArgU, IABC /* */, "SETLIST "}, // A B C, R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
	opcode{0, 1, OpArgU, OpArgN, IABx /* */, "CLOSURE "}, // A Bx,  R(A) := closure(KPROTO[Bx])
	opcode{0, 1, OpArgU, OpArgN, IABC /* */, "VARARG  "}, // A B,   R(A), R(A+1), ..., R(A+B-2) = vararg
	opcode{0, 0, OpArgU, OpArgU, IAx /*  */, "EXTRAARG"}, // Ax,    extra (larger) argument for previous opcode
}

// OpName return opname of instruction
func (i Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}

// OpMode return opmode of instruction
func (i Instruction) OpMode() byte {
	return opcodes[i.Opcode()].opMode
}

// BMode return argBMode of instruction
func (i Instruction) BMode() byte {
	return opcodes[i.Opcode()].argBMode
}

// CMode return argCMode of instruction
func (i Instruction) CMode() byte {
	return opcodes[i.Opcode()].argCMode
}
