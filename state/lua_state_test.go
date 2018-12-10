package state

import (
	"fmt"
	"log"
	"testing"
	"vczn/luago/api"
)

func printStack(ls *LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case api.LuaTBoolean:
			log.Printf("[%t]", ls.ToBoolean(i))
		case api.LuaTNumber:
			log.Printf("[%g]", ls.ToNumber(i))
		case api.LuaTString:
			log.Printf("[%q]", ls.ToString(i))
		default:
			log.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

func TestStateStack(t *testing.T) {
	ls := NewLuaState(20, nil)
	ls.PushBoolean(true)
	printStack(ls)
	ls.PushInteger(10)
	printStack(ls)
	ls.PushNil()
	printStack(ls)
	ls.PushString("hello")
	printStack(ls)
	ls.PushValue(-4)
	printStack(ls)
	ls.Replace(3)
	printStack(ls)
	ls.SetTop(6)
	printStack(ls)
	ls.Remove(-3)
	printStack(ls)
	ls.SetTop(-5)
	printStack(ls)
}

func TestStateOperator(t *testing.T) {
	ls := NewLuaState(20, nil)
	ls.PushInteger(1)
	ls.PushString("2.0")
	ls.PushString("3.0")
	ls.PushNumber(4.0)
	printStack(ls) // [1] ["2.0"] ["3.0"] [4.0]
	ls.Arithmetic(api.LuaOpAdd)
	printStack(ls) // [1] ["2.0"] [7]
	if !ls.IsNumber(-1) {
		t.Error("add error1")
	}
	if ls.ToNumber(-1) != 7.0 {
		t.Error("add error2")
	}
	ls.Arithmetic(api.LuaOpBNot)
	printStack(ls) // [1] ["2.0"] [-8]
	ls.Len(2)
	printStack(ls) // [1] ["2.0"] [-8] [3]
	if !ls.IsNumber(-1) {
		t.Error("len error1")
	}
	if ls.ToNumber(-1) != 3.0 {
		t.Error("len error2")
	}
	ls.Concat(3)
	printStack(ls) // [1] ["2.0-83"]
	if !ls.IsString(-1) {
		t.Error("concat error1")
	}
	if ls.ToString(-1) != "2.0-83" {
		t.Error("concat error2")
	}
	ls.PushBoolean(ls.Compare(1, 2, api.LuaOpEq))
	printStack(ls) // [1] ["2.0-83"] [false]
	if !ls.IsBoolean(-1) {
		t.Error("compare error1")
	}
	if ls.ToBoolean(-1) != false {
		t.Error("compare error2")
	}
}
