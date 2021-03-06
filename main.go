package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"luago/api"
	"luago/binchunk"
	"luago/compiler/lexer"
	"luago/compiler/parser"
	"luago/state"
	"luago/vm"
	"os"
)

// luatype      golangtype
// luaByte      byte
// cint         uint32?
// size_t       uint64
// luaint       int64
// luafloat     float64
// string
// table(list)

// table
// n | ptr

// chunk 内部
// 指令表、常量表、子函数原型等信息都是 list 存储的.

func list(proto *binchunk.ProtoType) {
	printHeader(proto)
	printCode(proto)
	printDetail(proto)
	for _, p := range proto.Protos {
		list(p)
	}
}

func printHeader(proto *binchunk.ProtoType) {
	funcName := "main"
	if proto.LineDefined > 0 {
		funcName = "function"
	}
	varargFlag := ""
	if proto.IsVararg > 0 {
		varargFlag = "+"
	}
	fmt.Printf("\n%s <%s:%d, %d> (%d instructions)\n", funcName,
		proto.Source, proto.LineDefined, proto.LastLineDefined, len(proto.Code))
	fmt.Printf("%d%s params, %d slots, %d upvalues, %d locals, %d constants, %d functions\n",
		proto.NumParams, varargFlag, proto.MaxStackSize, len(proto.Upvalues),
		len(proto.LocVars), len(proto.Constants), len(proto.Protos))
}

func printOperands(i vm.Instruction) {
	switch i.OpMode() {
	case vm.IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a)
		if i.BMode() != vm.OpArgN {
			if b > 0xFF {
				fmt.Printf(" %d", -1-b&0xFF) // constants table index
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.CMode() != vm.OpArgN {
			if c > 0xFF {
				fmt.Printf(" %d", -1-c&0xFF) // constants table index
			} else {
				fmt.Printf(" %d", c)
			}
		}
	case vm.IABx:
		a, bx := i.ABx()
		fmt.Printf("%d", a)
		if i.BMode() == vm.OpArgK { // constants table index
			fmt.Printf(" %d", -1-bx)
		} else if i.BMode() == vm.OpArgU {
			fmt.Printf(" %d", bx)
		}
	case vm.IAsBx:
		a, sBx := i.AsBx()
		fmt.Printf("%d %d", a, sBx)
	case vm.IAx:
		ax := i.Ax()
		fmt.Printf("%d", -1-ax)
	}
}

func printCode(proto *binchunk.ProtoType) {
	for pc, c := range proto.Code {
		line := "-"
		if len(proto.LineInfo) > 0 {
			line = fmt.Sprintf("%d", proto.LineInfo[pc])
		}
		i := vm.Instruction(c)
		fmt.Printf("\t%d\t[%s]\t%s\t", pc+1, line, i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

func constToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}

func upvalueName(proto *binchunk.ProtoType, i int) string {
	if len(proto.UpvalueNames) > 0 {
		return proto.UpvalueNames[i]
	}
	return "-"
}

func printDetail(proto *binchunk.ProtoType) {
	// constants table
	fmt.Printf("constants(%d):\n", len(proto.Constants))
	for i, k := range proto.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constToString(k))
	}

	// locals table
	fmt.Printf("locals(%d):\n", len(proto.LocVars))
	for i, locVar := range proto.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}

	// upvalues table
	fmt.Printf("upvalues(%d):\n", len(proto.Upvalues))
	for i, upvalue := range proto.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalueName(proto, i),
			upvalue.Instack, upvalue.Idx)
	}
}

func print(ls api.ILuaState) int {
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) { // string and number
			fmt.Printf("%s", ls.ToString(i))
		} else {
			fmt.Printf(ls.TypeName(ls.Type(i)))
		}

		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}

func testChunkDump() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage binchunk <outfiles...>")
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	proto := binchunk.Undump(data)
	list(proto)
}

func testState() {
	ls := state.NewLuaState()
	ls.PushInteger(1)
	ls.PushString("2.0")
	ls.PushString("3.0")
	ls.PushNumber(4.0)
	printStack(ls)

	ls.Arithmetic(api.LuaOpAdd)
	printStack(ls)
	ls.Arithmetic(api.LuaOpBNot)
	printStack(ls)
	ls.Len(2)
	printStack(ls)
	ls.Concat(3)
	printStack(ls)
	ls.PushBoolean(ls.Compare(1, 2, api.LuaOpEq))
	printStack(ls)
}

func printStack(ls *state.LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case api.LuaTBoolean:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case api.LuaTNumber:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case api.LuaTString:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

// func luaMain(proto *binchunk.ProtoType) {
// 	nRegs := int(proto.MaxStackSize)
// 	ls := state.NewLuaState()
// 	ls.SetTop(nRegs)
// 	for {
// 		pc := ls.PC()
// 		inst := vm.Instruction(ls.Fetch())
// 		if inst.Opcode() != vm.OpRETURN {
// 			inst.Execute(ls)
// 			fmt.Printf("[%02d] %s", pc+1, inst.OpName())
// 			printStack(ls)
// 		} else {
// 			break
// 		}
// 	}
// }

// func testVM() {
// 	// if len(os.Args) != 2 {
// 	// 	log.Fatal("Usage main <luac.out> ")
// 	// }
// 	data, err := ioutil.ReadFile("table.out")
// 	if err != nil {
// 		panic(err)
// 	}

// 	proto := binchunk.Undump(data)
// 	luaMain(proto)
// }

func testFunction() {
	data, err := ioutil.ReadFile("err.out")
	if err != nil {
		panic(err)
	}
	ls := state.NewLuaState()
	ls.Register("print", print)
	ls.Register("getmetatable", getMetaTable)
	ls.Register("setmetatable", setMetaTable)
	ls.Register("next", next)
	ls.Register("pairs", pairs)
	ls.Register("ipairs", ipairs)
	ls.Register("error", luaError)
	ls.Register("pcall", pCall)
	ls.Load(data, "table", "b")
	ls.Call(0, 0)
}

// temporary function for testing
func getMetaTable(ls api.ILuaState) int {
	if ls.GetMetaTable(1) {
		ls.PushNil()
	}
	return 1
}

func setMetaTable(ls api.ILuaState) int {
	ls.SetMetaTable(1)
	return 1
}

func next(ls api.ILuaState) int {
	ls.SetTop(2) // 1 table[, 2 key], key can be nil
	if ls.Next(1) {
		return 2 // returns nextKey, nextVal
	}
	ls.PushNil()
	return 1 // returns nil
}

// <==>
// _f, _s, _var = next, t, nil
func pairs(ls api.ILuaState) int {
	ls.PushGoFunction(next) // _f, next
	ls.PushValue(1)         // _s, table
	ls.PushNil()            // _var, key
	return 3
}

func ipairs(ls api.ILuaState) int {
	ls.PushGoFunction(inext)
	ls.PushValue(1)
	ls.PushInteger(0)
	return 3
}

func inext(ls api.ILuaState) int {
	nexti := ls.ToInteger(2) + 1
	ls.PushInteger(nexti)
	if ls.GetI(1, nexti) == api.LuaTNil {
		return 1
	}
	return 2
}

func luaError(ls api.ILuaState) int {
	return ls.Error()
}

func pCall(ls api.ILuaState) int {
	nArgs := ls.GetTop() - 1 // args
	status := ls.PCall(nArgs, -1, 0)
	ls.PushBoolean(status == api.LuaOk)
	ls.Insert(1)
	return ls.GetTop()
}

func testLexer() {
	data := `print("hello world\n", 42, 3.14)`
	lex := lexer.NewLexer(data, "string")
	for {
		line, kind, token := lex.NextToken()
		fmt.Printf("[%2d] [%-10s] %s\n",
			line, lexer.KindToString(kind), token)
		if kind == lexer.TokenEOF {
			break
		}
	}
}

func testParse() {
	data := `print("hello world\n", 42, 3.14)`
	ast := parser.Parse(data, "string")
	b, err := json.Marshal(ast)
	if err != nil {
		panic(err)
	}

	println(string(b))
}

func main() {
	//testChunkDump()
	testState()
	// testVM()
	// testFunction()
	// testParse()
}
